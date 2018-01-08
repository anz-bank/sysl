#!/usr/local/bin/python
# -*- encoding: utf-8 -*-

import hashlib
import itertools
import logging
import re
import struct
import warnings

from sysl.proto import sysl_pb2

from sysl.core import syslx

from sysl.util import datamodel
from sysl.util import java
from sysl.util import scopes
from sysl.util import writer


def _autoinc_fields(t, module):
    """Discover all autoinc fields, including those transitively so via FKs."""
    result = set()
    for (fname, f) in datamodel.sorted_fields(t):
        jfname = java.name(fname)
        method = java.CamelCase(jfname)
        (java_type, type_info, ftype) = datamodel.typeref(f, module)
        if ftype and 'autoinc' in syslx.patterns(ftype.attrs):
            result.add(fname)
    return result


def export_entity_class(w, tname, t, fk_rsubmap, context):
    '''Export the "Entity" class
    w          : The Writer to use to capture the Java text
    tname      : The "type name" from the Application definition
    t          : The type itself from the Application definition
    fk_rsubmap : Reverse foreign key lookup
    package    : package name
    '''

    # If the Type is a relation or a tuple
    which_type = t.WhichOneof('type')
    if which_type in ('relation', 'tuple'):
        with java.Class(w, tname, context.write_file, package=context.package,
                        static=not context.package,
                        visibility='public' if context.package else 'private'):
            pkey = datamodel.primary_key_params(t, context.module)
            pkey_fields = {f for (_, f, _) in pkey}
            param_defs = [(typ, f) for (typ, _, f) in pkey]
            pk_param = ', '.join(f for (_, _, f) in pkey)

            fkeys = {
                java.name(fname): type_info
                for (fname, _, type_info) in datamodel.foreign_keys(t, context.module)}

            w()
            groups = 0
            for (is_pk, g) in itertools.groupby(
                datamodel.sorted_fields(t),
                    lambda t: t[0] in pkey_fields):
                if is_pk:
                    w('// primary key {{')
                elif groups:
                    w()
                groups += 1
                for (fname, f) in g:
                    jfname = java.name(fname)
                    (ftname, type_info, ftype) = datamodel.typeref(
                        f, context.module)

                    if fname in fkeys:
                        fk_type_info = fkeys[fname]
                        w(u'{} {};  // fk {}.{}',
                          ftname, jfname, fk_type_info.parent_path, fk_type_info.field)
                    else:
                        w(u'{} {};', ftname, jfname)

                if is_pk:
                    w('// }}')

            for (fname, f) in datamodel.sorted_fields(t):
                jfname = java.name(fname)
                method = java.CamelCase(jfname)
                (java_type, type_info, _) = datamodel.typeref(f, context.module)
                java.SeparatorComment(w)

                with java.Method(w, '\npublic', java_type, 'get' + method):
                    w(u'return {};', jfname)

                if fname in pkey_fields:
                    with java.Method(w, '\npublic', tname, '_PRIVATE_set' + method,
                                     [(java_type, jfname)]):
                        w(u'this.{0} = {0};', jfname)
                        attr_def = t.relation.attr_defs[fname]
                        if 'autoinc' in syslx.patterns(attr_def.attrs):
                            with java.If(w, '_model != null && {} != null', jfname):
                                w('_model._PRIVATE_registerId({});', jfname)
                        w('return this;')
                else:
                    if fname in fk_rsubmap:
                        logging.error('%s', (tname, fname, fk_rsubmap[fname]))
                        raise Exception('Error exporting entity class')
                    if jfname in fkeys:
                        fk_type = fkeys[jfname].parent_path
                        fk_field = fkeys[jfname].field
                        if f.type_ref.ref.path[-1:] == [fname]:
                            method_suffix = fk_type
                        else:
                            # TODO: Mention fk field? What if CustCust had no
                            # pk?
                            method_suffix = method + 'From'
                        with java.Method(w, '\npublic', tname, 'set' + method_suffix,
                                         [(fk_type, 'entity')]):
                            w('{} {} = entity == null ? null : entity.{};',
                              java_type, jfname, fk_field)
                            w('Table table = _model.get{}Table();', tname)
                            w()
                            with java.If(w, 'this.{} != null', jfname):
                                w('table.getFk_{0}({0}).remove(this);', jfname)
                            w()
                            w(u'this.{0} = {0};', jfname)
                            w()
                            with java.If(w, 'this.{} != null', jfname):
                                w('table.getFk_{0}({0}).add(this);', jfname)
                            w()
                            w('return this;')

                        with java.Method(w, '\npublic', tname, '_PRIVATE_set' + method,
                                         [(java_type, jfname)]):
                            w(u'this.{0} = {0};', jfname)
                            attr_def = t.relation.attr_defs[fname]
                            if 'autoinc' in syslx.patterns(attr_def.attrs):
                                with java.If(w, '_model != null && {} != null', jfname):
                                    w('_model._PRIVATE_registerId({});', jfname)
                            w('return this;')
                    else:
                        with java.Method(w, '\npublic', tname, 'set' + method,
                                         [(java_type, jfname)]):
                            w(u'this.{0} = {0};', jfname)
                            w('return this;')

                if type_info and type_info.parent_path:
                    fk_type = type_info.parent_path
                    w()
                    if f.type_ref.ref.path[-1:] == [fname]:
                        method_suffix = ''
                    else:
                        method_suffix = 'Via' + method
                    with java.Method(w, 'public', fk_type,
                                     'to' + fk_type + method_suffix):
                        w('return {1} == null ? null : _model.get{0}Table().lookup({1});',
                          fk_type, jfname)

            java.SeparatorComment(w)

            for (fname, f) in datamodel.sorted_fields(t):
                jfname = java.name(fname)
                for (fk_tname, fk_fname) in sorted(fk_rsubmap[fname]):
                    fk_jfname = java.name(fk_fname)
                    fk_method = java.CamelCase(fk_jfname)
                    assert fk_tname in context.app.types, (
                        (fk_tname, context.app.types.keys()))
                    fk_type = context.app.types[fk_tname]
                    assert fk_type.HasField(
                        'relation'), fk_type.WhichOneof('type')
                    fk_pkey = datamodel.primary_key_params(
                        fk_type, context.module)
                    fk_pkey_fields = {f for (_, f, _) in fk_pkey}

                    method_suffix = '' if fk_jfname == jfname else 'Via' + fk_method

                    w()
                    if {fk_fname} == fk_pkey_fields:
                        with java.Method(w, 'public', fk_tname,
                                         'to' + fk_tname + method_suffix):
                            w(('return _model.get{}Table().getBy{}({}).singleOrNull();'),
                              fk_tname, fk_method, jfname)
                    else:
                        with java.Method(w, 'public', fk_tname + '.View',
                                         'to' + fk_tname + 'View' + method_suffix):
                            w('return _model.get{}Table().getBy{}({});',
                              fk_tname, fk_method, jfname)

            with java.Method(w, '\npublic', 'boolean', 'in', [('View', 'view')]):
                w(u'return view.contains(this);')

            if t.HasField('relation'):
                with java.Method(w, '\npublic', 'void', 'delete'):
                    for (fname, f) in datamodel.sorted_fields(t):
                        jfname = java.name(fname)
                        for (fk_tname, fk_fname) in sorted(fk_rsubmap[fname]):
                            fk_jfname = java.name(fk_fname)
                            with java.If(w, '_model.get{}Table().fk_{}.containsKey({})',
                                         fk_tname, fk_jfname, jfname):
                                w(u'throw new {}Exception();', context.model_class)

                    w(u'Table table = _model.get{}Table();', tname)

                    if pkey:
                        w(u'table.items.remove(new Table.Key({}));', pk_param)
                    else:
                        w(u'table.items.remove(this);')

                    for (fname, _, _) in datamodel.foreign_keys(
                            t, context.module):
                        jfname = java.name(fname)
                        with java.If(w, '{0} != null', jfname):
                            w('HashSet<{0}> fk = table.fk_{1}.get({1});',
                              tname, jfname)
                            w(u'fk.remove(this);')
                            with java.If(w, 'fk.size() == 0'):
                                w(u'table.fk_{0}.remove({0});', jfname)

                with java.Method(w, '\npublic', 'View', 'asView'):
                    w('Snapshot result = new Snapshot(_model);')
                    w('result.insert(this);')
                    w('return result;')

            w()
            export_entity_view_class(w, tname, t, fk_rsubmap, context)

            w()
            export_entity_snapshot_class(w, tname, t, context)

            w()
            export_entity_set_class(w, tname, t, context)

            w()
            export_entity_table_class(w, tname, t, fk_rsubmap, context)

            with java.Method(w, '\npublic', 'boolean', 'equals',
                             [('Object', 'object')], override=True):
                with java.If(w, 'this == object'):
                    w('return true;')
                with java.If(w, '!(object instanceof {})', tname):
                    w('return false;')

                w('{0} that = ({0})object;', tname)
                if which_type == 'relation':
                    w('assert _model != null && _model == that._model;')
                if pkey:
                    w('return _key.equals(that._key);', tname)
                else:
                    for (fname, f) in datamodel.sorted_fields(t):
                        jfname = java.name(fname)
                        with java.If(w,
                                     ('this.{0} != that.{0} && '
                                      '(this.{0} == null || !this.{0}.equals(that.{0}))'),
                                     jfname):
                            w('return false;')
                    w('return true;')

            with java.Method(w, '\npublic', 'int', 'hashCode', override=True):
                if pkey:
                    w('return _key == null ? super.hashCode() : _key.hashCode();', tname)
                else:
                    w('int i = 0;')
                    for (fname, f) in datamodel.sorted_fields(t):
                        jfname = java.name(fname)
                        w('i = 3 * i + (this.{0} == null ? 0 : this.{0}.hashCode());',
                          jfname)
                    w('return i;')

            with java.Method(w, '\npublic', 'int', 'canonicalHashCode',
                             [('int[]', 'ids')]):
                w('int i = 0;')
                for (fname, f) in datamodel.sorted_fields(t):
                    jfname = java.name(fname)
                    ref = 'this.' + jfname
                    if 'autoinc' in syslx.patterns(f.attrs):
                        ref = 'ids[{}]'.format(ref)
                        w('i = 3 * i + (this.{} == null ? 0 : {});', jfname, ref)
                    else:
                        w('i = 3 * i + (this.{} == null ? 0 : {}.hashCode());',
                          jfname, ref)
                w('return i;')

            ais = _autoinc_fields(t, context.module)

            with java.Method(w, '\n', 'boolean', 'canonicallyEqual',
                             [(tname, 'a'), ('int[]', 'aIds'),
                              (tname, 'b'), ('int[]', 'bIds')],
                             static=True):
                with java.If(w, 'a == b'):
                    w('return true;')
                if pkey:
                    with java.If(w, 'a._model != null && a._model == b._model'):
                        w('return a._key.equals(b._key);', tname)
                for (fname, f) in datamodel.sorted_fields(t):
                    jfname = java.name(fname)
                    if fname in ais:
                        test = 'aIds[a.{0}] != bIds[b.{0}]'.format(jfname)
                    else:
                        test = '!a.{0}.equals(b.{0})'.format(jfname)
                    with java.If(w,
                                 'a.{0} != b.{0} && (a.{0} == null || {1})',
                                 jfname, test):
                        w('return false;')
                w('return true;')

            if which_type == 'relation':
                with java.Method(w, '\npublic', 'int', 'compareTo', [(tname, 'that')]):
                    with java.If(w, 'this == that'):
                        w('return 0;')
                    if pkey:
                        with java.If(w, '_model != null && _model == that._model'):
                            w('return _key.compareTo(that._key);', tname)
                    w('int c;')
                    for (fname, f) in datamodel.sorted_fields(t):
                        jfname = java.name(fname)
                        with java.If(w,
                                     ('(c = '
                                      '\vthis.{0} == that.{0} ? 0 : '
                                      '\vthis.{0} == null ? -1 : '
                                      '\vthat.{0} == null ? 1 : '
                                      '\vthis.{0}.compareTo(that.{0})) != 0'),
                                     jfname):
                            w('return c;')
                    w('return 0;')

            with java.Method(w, '\npublic', 'String', 'toString', override=True):
                w('StringBuilder sb = new StringBuilder();')
                w('canonicalToString(sb, null);')
                w('return sb.toString();')

            with java.Method(w, '\npublic', 'void', 'canonicalToString',
                             [('StringBuilder', 'sb'), ('int[]', 'ids')]):
                sorted_fields = datamodel.sorted_fields(t)
                if len(sorted_fields) > 1:
                    w('String sep = "";')
                w('sb.append("{}(");', tname)
                for (i, (fname, f)) in enumerate(datamodel.sorted_fields(t)):
                    jfname = java.name(fname)
                    with java.If(w, '{} != null', jfname):
                        if len(sorted_fields) > 1:
                            w('sb.append(sep); sep = ", ";')
                        w('sb.append("{} = ");', jfname)
                        which_ftype = f.WhichOneof('type')
                        if which_ftype == 'primitive' or (
                                which_ftype == 'set' and
                                f.set.WhichOneof('type') == 'primitive'):
                            (_, _, ftype) = datamodel.typeref(f, context.module)
                            if 'autoinc' in syslx.patterns(f.attrs) or (
                                    ftype and 'autoinc' in syslx.patterns(ftype.attrs)):
                                with java.If(w, 'ids == null'):
                                    w('sb.append({});', jfname)
                                    with java.Else(w):
                                        w('sb.append("(");')
                                        w('sb.append(ids[{}]);', jfname)
                                        w('sb.append(")");')
                            else:
                                w('sb.append({});', jfname)

                        elif which_ftype == 'type_ref':
                            (java_type, type_info, fktype) = (
                                datamodel.typeref(f, context.module))
                            if type_info and type_info.parent_path:
                                if fktype and 'autoinc' in syslx.patterns(
                                        fktype.attrs):
                                    with java.If(w, 'ids == null'):
                                        w('sb.append({});', jfname)
                                        with java.Else(w):
                                            w('sb.append("(");')
                                            w('sb.append(ids[{}]);', jfname)
                                            w('sb.append(")");')
                                else:
                                    w('sb.append({});', jfname)
                            else:
                                w('sb.append({}.toString());', jfname)
                        else:
                            w('sb.append({}.toString());', jfname)
                w('sb.append(")");')

            java.SeparatorComment(w)
            w('// non-public')

            w()
            if which_type == 'relation':
                with java.Ctor(w, '', tname, [(context.model_class, 'model')]):
                    w('this._model = model;')

                with java.Ctor(w, 'private', tname):
                    w('_model = null;')

                with java.Method(w, 'public', tname, '_PRIVATE_new', static=True):
                    w('return new {}();', tname)

                with java.Method(w, 'public', tname, '_PRIVATE_new',
                                 [(context.model_class, 'model')], static=True):
                    w('return new {}(model);', tname)

                # The Model class member
                w(u'\nfinal {} _model;', context.model_class)
            else:
                with java.Ctor(w, 'public', tname):
                    pass

                with java.Method(w, '\npublic', tname, '_PRIVATE_new', static=True):
                    w('return new {}();', tname)

            if pkey:
                w(u'Snapshot.Key _key;')
            w()

    elif which_type == 'enum':
        with java.Class(w, tname, context.write_file, type_='enum',
                        package=context.package):
            for (i, (fname, fvalue)) in enumerate(t.enum.items.iteritems()):
                w(u'{}({}){}',
                  fname,
                  fvalue,
                  ',' if i < len(t.enum.items) - 1 else ';')

            w('private int value;')

            with java.Ctor(w, '\nprivate', tname, [('int', 'value')]):
                w('this.value = value;')

            with java.Method(w, '\npublic', 'int', 'getValue'):
                w('return value;')

            with java.Method(w, '\npublic', tname, 'from', [('int', 'i')],
                             static=True):
                with java.Switch(w, 'i'):
                    for (fname, fvalue) in t.enum.items.iteritems():
                        w('case {}: return {};', fvalue, fname)
                w('return null;')

            with java.Method(w, '\npublic', tname, 'from', [('String', 'name')],
                             static=True):
                with java.Switch(w, 'name'):
                    for (fname, fvalue) in t.enum.items.iteritems():
                        w('case "{0}": return {0};', fname)
                w('return null;')

    elif which_type is None:
        warnings.warn('WARNING: empty type: {}'.format(tname))

    else:
        raise RuntimeError('Unexpected type: {}'.format(t.WhichOneof('type')))


def export_entity_view_class(w, tname, t, fk_rsubmap, context):
    if t.WhichOneof('type') in ('relation', 'tuple'):
        with java.Class(w, 'View', context.write_file, static=True, abstract=True,
                        extends='io.sysl.Enumerable<' + tname + '>',
                        implements=['Comparable<View>']):
            pkey = datamodel.primary_key_params(t, context.module)
            pkey_fields = {f for (_, f, _) in pkey}
            param_defs = [(typ, f) for (typ, _, f) in pkey]
            params = ''.join(', ' + f for (_, _, f) in pkey)
            param = ', '.join(f for (_, _, f) in pkey)

            inner_type = (
                'HashMap<Key, {}>' if pkey else 'HashSet<{}>').format(tname)

            if pkey:
                with java.Method(w, '\npublic', tname, 'lookup', param_defs):
                    w(u'return lookup(new Key({}));', param)

                with java.Method(w, '\npublic', 'boolean', 'contains',
                                 [(tname, 'entity')]):
                    w('return lookup(entity._key) != null;')
            else:
                with java.Method(w, '\npublic', 'boolean', 'contains',
                                 [(tname, 'entity')]):
                    with java.For(w, '{} item : this', tname):
                        with java.If(w, 'item.equals(entity)'):
                            w('return true;')
                    w('return false;')

            with java.Method(w, '\npublic', 'boolean', 'contains',
                             [('Object', 'o')]):
                w('return o instanceof {0} && contains(({0})o);', tname)

            with java.Method(w, '\npublic', 'boolean', 'containsAll',
                             [('Collection<? extends {}>'.format(tname), 'c')]):
                with java.For(w, '{} e : c', tname):
                    with java.If(w, '!contains(e)'):
                        w('return false;')
                w('return true;')

            with java.Method(w, '\npublic', 'View', 'first',
                             [('int', 'n'), ('Comparator<' + tname + '>', 'comp')]):
                w('return view(model, '
                  'io.sysl.Enumeration.first(this, n, comp));')

            with java.Method(w, '\npublic', tname, 'singleOrNull',
                             throws=[context.model_class + 'Exception'], override=True):
                with java.Try(w):
                    w('return super.singleOrNull();')
                    with java.Catch(w, 'io.sysl.SyslException', 'ex'):
                        w(u'throw new {}Exception(ex.getMessage());',
                          context.model_class)

            with java.Method(w, '\npublic', tname, 'single',
                             throws=[context.model_class + 'Exception'], override=True):
                with java.Try(w):
                    w('return super.single();')
                    with java.Catch(w, 'io.sysl.SyslException', 'ex'):
                        w(u'throw new {}Exception(ex.getMessage());',
                          context.model_class)

            with java.ViewMethod(w, tname, '\npublic', tname, 'any',
                                 [('final int', 'n')]) as view:
                w('final View input = this;')
                with view() as enumerator_method:
                    with java.Method(w, 'public', 'int', 'size', override=True):
                        w('return input.sizeWithLimit(n);')
                    with enumerator_method():
                        w('return view(model, io.sysl.Enumeration.any(input, n))'
                          '.enumerator();')

            with java.ViewMethod(w, tname, '\npublic', tname, 'where',
                                 [('final io.sysl.Expr<Boolean, ' + tname + '>', 'pred')]
                                 ) as view:
                w('final View input = this;')
                with view() as enumerator_method:
                    if pkey:
                        with java.Method(w, '', tname, 'lookup', [('Key', 'key')]):
                            w(u'{} result = input.lookup(key);', tname)
                            with java.If(w,
                                         'result != null && pred.evaluate(result).booleanValue()'):
                                w(u'return result;')
                            w(u'return null;')
                    else:
                        with java.Method(w, 'public', 'boolean', 'contains',
                                         [(tname, 'entity')], override=True):
                            w('return pred.evaluate(entity).booleanValue() && '
                              '\vView.this.contains(entity);')
                    with enumerator_method():
                        w('return io.sysl.Enumeration.where(input, pred)'
                          '.enumerator();')

            with java.Method(w, '\npublic', 'View', 'and', [('View', 'that')]):
                w(u'return join(that, 2);')

            with java.Method(w, '\npublic', 'View', 'or', [('View', 'that')]):
                w(u'return join(that, 4 | 2 | 1);')

            with java.Method(w, '\npublic', 'View', 'xor', [('View', 'that')]):
                w(u'return join(that, 4 | 1);')

            with java.Method(w, '\npublic', 'View', 'butnot', [('View', 'that')]):
                w(u'return join(that, 4);')

            for (fname, f) in datamodel.sorted_fields(t):
                jfname = java.name(fname)
                method = java.CamelCase(jfname)
                (java_type, type_info, _) = datamodel.typeref(f, context.module)

                # Nested entities or sets thereof
                if (type_info and
                        type_info.type.WhichOneof('type') in ['tuple', 'relation']):
                    if f.WhichOneof('type') == 'set':
                        with java.Method(w, '\npublic', java_type,
                                         'get' + method + 'View'):
                            w(('return {r}.View.view(null, '
                               'io.sysl.Enumeration.flatten(this,\n'
                               '    new io.sysl.Expr'
                               '<io.sysl.Enumerable<{r}>, {t}>() {{\n'
                               '        @Override\n'
                               '        public io.sysl.Enumerable<{r}> '
                               'evaluate({t} t) {{\n'
                               '            return t.get{method}();\n'
                               '        }}\n'
                               '    }}));'),
                              r=re.sub(r'\.View$', '', java_type),
                              t=tname,
                              method=method)
                    else:
                        with java.Method(w, '\npublic', java_type + '.View',
                                         'get' + method + 'View'):
                            w(('return {r}.View.view(null, '
                               'io.sysl.Enumeration.map(this,\n'
                               '    new io.sysl.Expr<{r}, {t}>() {{\n'
                               '        @Override\n'
                               '        public {r} evaluate({t} t) {{\n'
                               '            return t.get{method}();\n'
                               '        }}\n'
                               '    }}));'),
                              r=java_type,
                              t=tname,
                              method=method)

            first_fk = True
            for (fname, f) in datamodel.sorted_fields(t):
                jfname = java.name(fname)
                method = java.CamelCase(jfname)
                (java_type, type_info, _) = datamodel.typeref(f, context.module)

                if type_info and type_info.parent_path:
                    java.SeparatorComment(w)
                    if first_fk:
                        first_fk = False

                    fk_tname = type_info.parent_path
                    if f.type_ref.ref.path[-1:] == [fname]:
                        method_suffix = ''
                    else:
                        method_suffix = 'Via' + method

                    with java.ViewMethod(w, tname, '\n', fk_tname,
                                         'to' + fk_tname + 'ViewInner' + method_suffix) as view:
                        with view() as enumerator_method:
                            with enumerator_method('source') as enumerator:
                                with java.If(w, 'View.this.isEmpty()'):
                                    w(('return io.sysl.Enumeration.<{}>empty()'
                                       '.enumerator();'),
                                      fk_tname)
                                w(u'final {0}.Table index = model.get{0}Table();', fk_tname)
                                with enumerator() as (moveNext, current):
                                    with moveNext():
                                        with java.While(w, 'source.moveNext()'):
                                            with java.If(w, 'source.current().{} != null', jfname):
                                                w('return true;')
                                        w('return false;')
                                    with current():
                                        w('return index.lookup(source.current().{});', jfname)

                    # TODO: Fold the above method into this one.
                    with java.Method(w, '\npublic', fk_tname + '.View',
                                     'to' + fk_tname + 'View' + method_suffix):
                        w(('return {0}.View.view(model, io.sysl.Enumeration.dedup('
                           'to{0}ViewInner{1}()));'),
                          fk_tname, method_suffix)

            for (fname, f) in datamodel.sorted_fields(t):
                jfname = java.name(fname)
                for (fk_tname, fk_fname) in sorted(fk_rsubmap[fname]):
                    fk_jfname = java.name(fk_fname)
                    fk_method = java.CamelCase(fk_jfname)
                    if fk_jfname == jfname:
                        method_suffix = ''
                    else:
                        method_suffix = 'Via' + fk_method

                    with java.Method(w, '\npublic', fk_tname + '.View',
                                     'to' + fk_tname + 'View' + method_suffix):
                        w(('return {0}.View.view(\vmodel, '
                           'io.sysl.Enumeration.flatten(\v'
                           'View.this, new io.sysl.Expr<\v'
                           'io.sysl.Enumerable<{0}>, {1}>() {{'),
                          fk_tname, tname)
                        with w.indent():
                            with java.Method(w, 'public',
                                             'io.sysl.Enumerable<' + fk_tname + '>',
                                             'evaluate',
                                             [(tname, 'entity')]):
                                w('return model.get{0}Table()\v.getBy{1}(\ventity.{2});',
                                  fk_tname, fk_method, jfname)
                        w('}}));')

            with java.Method(w, '\npublic', 'boolean', 'equals',
                             [('Object', 'that')], override=True):
                w('return (that instanceof View && _snapshot().items.equals(\v'
                  '((View)that)._snapshot().items));')

            with java.Method(w, '\npublic', 'int', 'hashCode', override=True):
                w('return _snapshot().items.hashCode();', tname)

            with java.Method(w, '\npublic', 'int', 'compareTo', [('View', 'obj')],
                             override=True):
                w('throw new java.lang.UnsupportedOperationException();')

            with java.Method(w, '\npublic', 'String', 'toString', override=True):
                with java.If(w, 'isEmpty()'):
                    w('return "[]";')
                w('StringBuilder sb = new StringBuilder();')
                w('toString(sb, "");')
                w('return sb.toString();')

            with java.Method(w, '\npublic', 'void', 'toString',
                             [('StringBuilder', 'sb'), ('String', 'indent')]):
                with java.If(w, 'isEmpty()'):
                    w('sb.append("[]");')
                    with java.Else(w):
                        w('sb.append("[\\n");', tname)
                        w('sb.append(indent + "  ");')
                        w('boolean first = true;')
                        w('Iterable<{}> iter;', tname)
                        w('iter = this;')

                        with java.For(w, '{} e : this', tname):
                            with java.If(w, 'first'):
                                w('first = false;')
                                with java.Else(w):
                                    w('sb.append(",\\n  " + indent);')
                            w('sb.append(e.toString());')
                        w('sb.append("\\n");')
                        w('sb.append(indent);')
                        w('sb.append("]");')

            if t.WhichOneof('type') in 'relation':
                with java.Method(w, '\npublic', 'void', 'canonicalToString',
                                 [('StringBuilder', 'sb'), ('String', 'indent'), ('int[]', 'ids')]):
                    with java.If(w, 'isEmpty()'):
                        w('sb.append("[]");')
                        with java.Else(w):
                            w('sb.append("[\\n");', tname)
                            w('sb.append(indent + "  ");')
                            w('boolean first = true;')
                            #w('Iterable<{}> iter;', tname)
                            with java.If(w, 'ids == null'):
                                with java.For(w, '{} e : this', tname):
                                    with java.If(w, 'first'):
                                        w('first = false;')
                                        with java.Else(w):
                                            w('sb.append(",\\n  " + indent);')
                                    w('e.canonicalToString(sb, ids);')

                                with java.Else(w):
                                    #   with w.indent(
                                    #       ('iter = io.sysl.Enumeration.orderBy('
                                    #        'this, new Comparator<{}>() {{'),
                                    #       tname):
                                    #     with java.Method(w, 'public', 'int', 'compare',
                                    #         [(tname, 'a'), (tname, 'b')], override=True):
                                    #       w('return 0;')
                                    #   w('}});')

                                    with java.For(w, '{} e : canonicallySorted(ids)', tname):
                                        with java.If(w, 'first'):
                                            w('first = false;')
                                            with java.Else(w):
                                                w('sb.append(",\\n  " + indent);')
                                        w('e.canonicalToString(sb, ids);')
                            w('sb.append("\\n");')
                            w('sb.append(indent);')
                            w('sb.append("]");')

                with java.Method(w, '\nprivate', 'io.sysl.Enumerable<' + tname + '>',
                                 'canonicallySorted', [('final int[]', 'ids')]):
                    with w.indent(
                        ('io.sysl.Enumerable<{0}> result = io.sysl.Enumeration.orderBy('
                         '\vthis, new Comparator<{0}>() {{'),
                            tname):
                        autoinc = None
                        with java.Method(w, 'public', 'int', 'compare',
                                         [(tname, 'a'), (tname, 'b')], override=True):
                            w('int c;')
                            first_autoinc = True
                            for (fname, f) in datamodel.sorted_fields(t):
                                jfname = java.name(fname)
                                if 'autoinc' in syslx.patterns(f.attrs):
                                    autoinc = fname
                                    continue
                                (_, _, ftype) = datamodel.typeref(
                                    f, context.module)
                                if ftype and 'autoinc' in syslx.patterns(
                                        ftype.attrs):
                                    if first_autoinc:
                                        w('int aId, bId;')
                                        first_autoinc = False
                                    with java.If(w,
                                                 ('(c = '
                                                  '\va.{0} == b.{0} ? 0 : '
                                                  '\va.{0} == null ? -1 : '
                                                  '\vb.{0} == null ? 1 : '
                                                  '\v(aId = ids[a.{0}]) == (bId = ids[b.{0}]) ? 0 : '
                                                  '\vaId < bId ? -1 : 1) != 0'),
                                                 jfname):
                                        w('return c;')
                                else:
                                    with java.If(w,
                                                 ('(c = '
                                                  '\va.{0} == b.{0} ? 0 : '
                                                  '\va.{0} == null ? -1 : '
                                                  '\vb.{0} == null ? 1 : '
                                                  '\va.{0}.compareTo(b.{0})) != 0'),
                                                 jfname):
                                        w('return c;')
                            w('return 0;')
                    w('}});')
                    if autoinc:
                        w('int canonicalId = 0;')
                        with java.For(w, '{} e : result', tname):
                            w('ids[e.{}] = canonicalId++;', autoinc)
                    w('return result;')

                with java.Method(w, '\n', 'boolean', 'canonicallyEqual',
                                 [('View', 'a'), ('int[]', 'aIds'),
                                  ('View', 'b'), ('int[]', 'bIds')],
                                 static=True):
                    with java.If(w, 'a == b'):
                        w('return true;')
                    with java.If(w, 'a == null || b == null'):
                        w('return false;')
                    with java.If(w, 'a.model == b.model'):
                        w('return a.equals(b);')

                    w(('io.sysl.Enumerator<{}> aSorted = '
                       'a.canonicallySorted(aIds).enumerator();'),
                      tname)
                    w(('io.sysl.Enumerator<{}> bSorted = '
                       'b.canonicallySorted(bIds).enumerator();'),
                      tname)
                    with java.While(w, 'aSorted.moveNext() && bSorted.moveNext()'):
                        with java.If(w,
                                     ('!{}.canonicallyEqual('
                                      'aSorted.current(), aIds, bSorted.current(), bIds)'),
                                     tname):
                            w('return false;')
                    w('return !aSorted.moveNext() && !bSorted.moveNext();')

                with java.Method(w, '\npublic', 'int', 'canonicalHashCode',
                                 [('int[]', 'ids')]):
                    w('int h = 0;')
                    with java.For(w, '{} e : canonicallySorted(ids)', tname):
                        w('h = 3 * h + e.canonicalHashCode(ids);')
                    w('return h;')

            with java.Method(w, '\npublic', 'View', 'snapshot'):
                w('return _snapshot();')

            java.SeparatorComment(w)
            w('// non-public')

            with java.Method(w, '\nprivate', 'Snapshot', '_snapshot'):
                w('Snapshot result = new Snapshot(model);')
                with java.For(w, '{} item : this', tname):
                    w('result.insert(item);')
                w('return result;')

            if pkey:
                with java.Method(w, '\n', tname, 'lookup', [('Key', 'key')]):
                    with java.For(w, '{} item : this', tname):
                        with java.If(w, 'item._key.equals(key)'):
                            w('return item;')
                    w('return null;')

            with java.Method(w, '\npublic', 'View', 'view',
                             [(context.model_class, 'model'),
                              ('final io.sysl.Enumerable<' + tname + '>', 'inner')],
                             static=True):
                w('return new View(model) {{')
                with w.indent():
                    with java.EnumeratorMethod(w, 'public', tname):
                        w('return inner.enumerator();')
                w('}};')

            # TODO: implement set operators without pkey
            with java.Method(w, '\nprivate', 'View', 'join',
                             [('View', 'that'), ('int', 'mask')]):
                w('return new Join(this, that, mask);')

            w()
            with java.Class(w, 'Join', context.write_file, visibility='private',
                            static=True, final=True, extends='View'):
                with java.Ctor(w, '', 'Join',
                               [('View', 'a'), ('View', 'b'), ('int', 'mask')]):
                    w('super(a.model);')
                    w('this.a = a;')
                    w('this.b = b;')
                    w('this.mask = mask;')

                if pkey:
                    with java.Method(w, '\n', tname, 'lookup', [('Key', 'key')]):
                        w('{} aLook = a.lookup(key);', tname)
                        w('{} bLook = b.lookup(key);', tname)
                        w('int pattern = (aLook != null ? 7 : 1) & '
                          '(bLook != null ? 7 : 4);')
                        with java.If(w, '(mask & pattern) != 0'):
                            w('return aLook != null ? aLook : bLook;')
                        w('return null;')

                with java.Method(w, 'public',
                                 'io.sysl.Enumerator<' + tname + '>', 'enumerator'):
                    w('return new Enumerator(this);')

                w()
                w('private final View a;')
                w('private final View b;')
                w('private final int mask;')

                w()
                with java.Class(w, 'Enumerator', context.write_file,
                                visibility='private', static=True, final=True,
                                implements=['io.sysl.Enumerator<' + tname + '>']):
                    with java.Ctor(w, '', 'Enumerator', [('Join', 'join')]):
                        w()
                        w('this.join = join;')
                        w('this.lhs = (join.mask & 4) != 0;')
                        w('this.both = (join.mask & 2) != 0;')
                        w('this.rhs = (join.mask & 1) != 0;')
                        w('this.aEnumerator = lhs || both ? join.a.enumerator() : null;')
                        w('this.bEnumerator = both || rhs ? join.b.enumerator() : null;')

                    with java.Method(w, '\npublic', 'boolean', 'moveNext'):
                        with java.If(w, 'aEnumerator != null'):
                            with java.While(w, 'aEnumerator.moveNext()'):
                                with java.If(w,
                                             'join.b.contains(aEnumerator.current()) ? both : lhs'):
                                    w('curr = aEnumerator.current();')
                                    w('return true;')
                            w('aEnumerator = null;')
                            with java.If(w, 'both && !rhs'):
                                w('bEnumerator = null;')
                        with java.If(w, 'bEnumerator != null'):
                            with java.While(w, 'bEnumerator.moveNext()'):
                                with java.If(w, '!join.a.contains(bEnumerator.current())'):
                                    w('curr = bEnumerator.current();')
                                    w('return true;')
                        w('curr = null;')
                        w('return false;')

                    with java.Method(w, '\npublic', tname, 'current'):
                        w('return curr;')

                    w()
                    w('private final Join join;', tname)
                    w('private final boolean lhs;')
                    w('private final boolean both;')
                    w('private final boolean rhs;')
                    w('private io.sysl.Enumerator<{}> aEnumerator;', tname)
                    w('private io.sysl.Enumerator<{}> bEnumerator;', tname)
                    w('private {} curr;', tname)

            if pkey:
                w()
                with java.Class(w, 'Key', context.write_file, static=True, final=True,
                                implements=['Comparable<Key>']):
                    for (typ, _, f) in pkey:
                        w('{} {};', typ, f)

                    with java.Ctor(w, '\n', 'Key', param_defs):
                        for (_, f) in param_defs:
                            with java.If(w, '{} == null', f):
                                w('throw new io.sysl.PrimaryKeyNullException("{}.{}");',
                                  tname, f)
                            w('this.{0} = {0};', f)

                    with java.Method(w, '\npublic', 'boolean', 'equals',
                                     [('Object', 'obj')], override=True):
                        w('Key that = (Key)obj;')
                        w('return ' +
                          u' && \n       '.join(
                              'this.{0}.equals(that.{0})'.format(f)
                              for (typ, _, f) in pkey) +
                          ';')

                    with java.Method(w, '\npublic', 'int', 'hashCode', override=True):
                        for (i, (typ, _, f)) in enumerate(pkey):
                            w('{}{}.hashCode();',
                              ('int h = ' if i == 0 and len(pkey) > 1 else
                               'return ' if i == 0 else
                               'h = 3 * h + ' if i < len(pkey) - 1 else
                               'return 3 * h + '),
                              f)

                    with java.Method(w, '\npublic', 'int', 'compareTo', [('Key', 'that')],
                                     override=True):
                        w('int c;')
                        w('return this == that ? 0 :\n' +
                          u'\n       '.join(
                              '(c = this.{0}.compareTo(that.{0})) != 0 ? c :'.format(
                                  f)
                              for (typ, _, f) in pkey) +
                          '    0;')

            with java.Method(w, '\npublic', context.model_class, 'model'):
                w('return model;')

            w('\nfinal {} model;', context.model_class)

            with java.Ctor(w, '\n', 'View', [(context.model_class, 'model')]):
                w(u'this.model = model;')


def export_entity_snapshot_class(w, tname, t, context):
    if t.WhichOneof('type') in ('relation', 'tuple'):
        with java.Class(w, 'Snapshot', context.write_file, static=True,
                        extends='View'):
            pkey = datamodel.primary_key_params(t, context.module)
            pkey_fields = {f for (_, f, _) in pkey}
            param_defs = [(typ, f) for (typ, _, f) in pkey]
            params = ''.join(', ' + f for (_, _, f) in pkey)
            param = ', '.join(f for (_, _, f) in pkey)

            inner_type = (
                'HashMap<Key, {}>' if pkey else 'HashSet<{}>').format(tname)

            if pkey:
                with java.Method(w, '\npublic', tname, 'lookup', param_defs,
                                 override=True):
                    w(u'return lookup(new Key({}));', param)
            # else:
            #   with java.Method(w, '\npublic', 'boolean', 'contains',
            #       [(tname, 'entity')], override=True):
            #     w('return items.contains(entity);')

            with java.Method(w, '\npublic',
                             'io.sysl.Enumerator<' + tname + '>', 'enumerator',
                             override=True):
                w(('return io.sysl.Enumeration.enumerator(\v'
                   'items{}.iterator());'),
                  '.values()' if pkey else '')

            with java.Method(w, '\npublic', 'boolean', 'isEmpty', override=True):
                w(u'return items.isEmpty();')

            with java.Method(w, '\npublic', 'int', 'size', override=True):
                w('return items.size();')

            with java.Method(w, '\npublic', tname, 'singleOrNull', override=True):
                with java.Switch(w, u'size()'):
                    w('case 0: return null;')
                    w('case 1: return iterator().next();')
                    w('default: throw new {}Exception("size() == " + size() + " > 1");',
                      context.model_class)

            with java.Method(w, '\npublic', tname, 'single', override=True):
                with java.If(w, 'size() != 1'):
                    w(u'throw new {}Exception("size() == " + size() + " ≠ 1");',
                      context.model_class)
                w(u'return iterator().next();')

            with java.Method(w, '\npublic', 'boolean', 'equals',
                             [('Object', 'that')], override=True):
                w(u'return that instanceof Snapshot &&')
                w(u'       items.equals(((Snapshot)that).items);')

            with java.Method(w, '\npublic', 'int', 'hashCode', override=True):
                w(u'return items.hashCode();', tname)

            java.SeparatorComment(w)
            w('// non-public')

            if pkey:
                with java.Method(w, '\n', tname, 'lookup', [('Key', 'key')],
                                 override=True):
                    w(u'return items.get(key);', param)

            w()
            w(u'final {} items;', inner_type)

            with java.Ctor(w, '\n', 'Snapshot', [(context.model_class, 'model')]):
                w(u'super(model);')
                w(u'this.items = new {}();', inner_type)

            with java.Method(w, '\n', 'boolean', 'insert', [(tname, 'entity')]):
                if pkey:
                    with java.If(w, 'entity._key == null'):
                        w('entity._key = new Key({});',
                          ', '.join('entity.' + f for (_, _, f) in pkey))
                    with java.If(w, 'items.containsKey(entity._key)'):
                        w('return false;')
                    w('items.put(entity._key, entity);')
                    w('return true;')
                else:
                    w('return items.add(entity);')

            w()


def export_entity_set_class(w, tname, t, context):
    if t.WhichOneof('type') in ('relation', 'tuple'):
        with java.Class(w, 'Set', context.write_file, static=True,
                        extends='Snapshot'):
            pkey = datamodel.primary_key_params(t, context.module)
            pkey_fields = {f for (_, f, _) in pkey}
            param_defs = [(typ, f) for (typ, _, f) in pkey]
            params = ''.join(', ' + f for (_, _, f) in pkey)
            param = ', '.join(f for (_, _, f) in pkey)

            inner_type = (
                'HashMap<Key, {}>' if pkey else 'HashSet<{}>').format(tname)

            with java.Ctor(w, 'public', 'Set'):
                w('super(null);')

            with java.Method(w, '\npublic', 'boolean', 'add', [(tname, 'e')]):
                w('return insert(e);')

            with java.Method(w, '\npublic', 'boolean', 'addAll',
                             [('Collection<{}>'.format(tname), 'c')]):
                w('boolean changed = false;')
                with java.For(w, '{} e : c', tname):
                    w('insert(e);')
                    w('changed = true;')
                w('return changed;')

            with java.Method(w, '\npublic', 'void', 'clear'):
                w('items.clear();')

            with java.Method(w, '\npublic', 'boolean', 'remove', [('Object', 'o')]):
                if pkey:
                    w(('return o instanceof {0} && \vitems.remove((({0})o)._key)'
                       ' != null;'),
                      tname)
                else:
                    w('return items.remove(o);')

            with java.Method(w, '\npublic', 'boolean', 'removeAll',
                             [('Collection<{}>'.format(tname), 'c')]):
                if pkey:
                    w('boolean changed = false;')
                    with java.For(w, '{} e : this', tname):
                        with java.If(w, 'contains(e)'):
                            w('remove(e);')
                    w('return true;')
                else:
                    w('return items.removeAll(c);')

            with java.Method(w, '\npublic', 'boolean', 'retainAll',
                             [('Collection<{}>'.format(tname), 'c')]):
                if pkey:
                    w('int old_size = items.size();')
                    w('ArrayList<{0}> buf = new ArrayList<{0}>();', tname)
                    w('buf.ensureCapacity(items.size());')
                    with java.For(w, '{} e : c', tname):
                        with java.If(w, 'contains(e)'):
                            w('buf.add(e);')
                    w('items.clear();')
                    with java.For(w, '{} e : buf', tname):
                        w('items.put(e._key, e);')
                    w()
                    w('return items.size() < old_size;')
                else:
                    w('return items.retainAll(c);')

            with java.Method(w, '\npublic', 'Object', 'toArray'):
                if pkey:
                    w('{0}[] result = new {0}[size()];', tname)
                    w('int i = 0;')
                    with java.For(w, '{} e : this', tname):
                        w('result[i] = e;')
                        w('i++;')
                    w('return result;')
                else:
                    w('return items.toArray();')


def export_entity_table_class(w, tname, t, fk_rsubmap, context):
    if t.WhichOneof('type') == 'relation':
        with java.Class(w, 'Table', context.write_file, static=True,
                        extends='Snapshot'):
            pkey = datamodel.primary_key_params(t, context.module)
            pkey_fields = {f for (_, f, _) in pkey}
            param_defs = [(typ, jfname)
                          for (typ, fname, jfname) in pkey
                          if ('autoinc' not in
                              syslx.patterns(t.relation.attr_defs[fname].attrs))]

            fkeys = {
                java.name(fname): type_info
                for (fname, _, type_info) in datamodel.foreign_keys(t, context.module)}

            inner_type = (
                'HashMap<Key, {}>' if pkey else 'HashSet<{}>').format(tname)

            with java.Method(w, '\npublic', 'void', 'clear'):
                for (fname, f) in datamodel.sorted_fields(t):
                    jfname = java.name(fname)
                    for (fk_tname, fk_fname) in sorted(fk_rsubmap[fname]):
                        fk_jfname = java.name(fk_fname)
                        with java.If(w, '!model.get{}Table().fk_{}.isEmpty()',
                                     fk_tname, fk_jfname):
                            w(u'throw new {}Exception();', context.model_class)
                w(u'items.clear();')

            java.SeparatorComment(w)
            w(u'// non-public')

            for _ in datamodel.foreign_keys(t, context.module):
                w()
                break
            for (fname, java_type, _) in datamodel.foreign_keys(
                    t, context.module):
                jfname = java.name(fname)
                w(u'final HashMap<{}, HashSet<{}>> fk_{};', java_type, tname, jfname)

            with java.Ctor(w, '\n', 'Table', [(context.model_class, 'model')]):
                w(u'super(model);')

                for (fname, java_type, _) in datamodel.foreign_keys(
                        t, context.module):
                    jfname = java.name(fname)
                    w(u'this.fk_{} = new HashMap<{}, HashSet<{}>>();',
                      jfname, java_type, tname)

            # TODO: make internal
            with java.Method(w, '\npublic', 'boolean', 'insert', [(tname, 'entity')]):
                fks = list(datamodel.foreign_keys(t, context.module))
                if pkey:
                    for (_, fname, jfname) in pkey:
                        ftype = t.relation.attr_defs[fname]
                        if 'autoinc' in syslx.patterns(ftype.attrs):
                            with java.If(w, 'entity.{} == null', jfname):
                                w('entity.{} = model.nextId;', jfname)
                                w('model.nextId++;')
                if fks:
                    with java.If(w, 'super.insert(entity)'):
                        for (fname, java_type, type_info) in fks:
                            jfname = java.name(fname)
                            method = java.CamelCase(jfname)
                            type_ = type_info.parent_path
                            with java.If(w, 'entity.{} != null', jfname):
                                w(u'getFk_{0}(entity.{0}).add(entity);', jfname)
                        w('return true;')
                    w('return false;')
                else:
                    w('return super.insert(entity);')

            for (fname, java_type, type_info) in (
                    datamodel.foreign_keys(t, context.module)):
                jfname = java.name(fname)
                method = java.CamelCase(jfname)
                type_ = type_info.parent_path

                with java.Method(w, '\n', 'HashSet<' + tname + '>', 'getFk_' + jfname,
                                 [(java_type, jfname)]):
                    w('HashSet<{0}> set = fk_{1}.get({1});', tname, jfname)
                    with java.If(w, 'set == null'):
                        w('set = new HashSet<{}>();', tname)
                        w('fk_{0}.put({0}, set);', jfname)
                    w('return set;')

                with java.ViewMethod(w, tname, '\n', tname, 'getBy' + method,
                                     [('final ' + java_type, jfname)]) as view:
                    with view() as enumerator_method:
                        with enumerator_method():
                            w('HashSet<{0}> set = fk_{1}.get({1});',
                              tname, jfname)
                            with java.If(w, 'set != null'):
                                w('return io.sysl.Enumeration.enumerator(set.iterator());')
                            w('return io.sysl.Enumeration.<{}>empty().enumerator();', tname)

            with java.Method(w, '\npublic', 'void', '_PRIVATE_insert',
                             [(tname, 'entity')]):
                w('insert(entity);')

            with java.Method(w, '\npublic', tname, '_PRIVATE_insert',
                             param_defs):
                w('{0} entity = new {0}(model);', tname)
                for (_, fname, jfname) in pkey:
                    ftype = t.relation.attr_defs[fname]
                    if 'autoinc' in syslx.patterns(ftype.attrs):
                        w('entity.{0} = model.nextId;', jfname)
                    else:
                        w('entity.{0} = {0};', jfname)
                w('insert(entity);')
                for (_, fname, jfname) in pkey:
                    ftype = t.relation.attr_defs[fname]
                    if 'autoinc' in syslx.patterns(ftype.attrs):
                        w('model.nextId++;', jfname, tname)
                w('return entity;')

            required_fields = [
                (fname, f)
                for (fname, f) in datamodel.sorted_fields(t)
                if 'required' in syslx.patterns(f.attrs)]

            constrained_fields = any(
                (f.primitive == f.STRING and
                    any(c.precision for c in f.constraint))
                for (_, f) in datamodel.sorted_fields(t)
            )

            if required_fields or fkeys or constrained_fields:
                with java.Method(w, '\npublic', 'void', 'validate'):
                    with java.For(w, '{} e : this', tname):
                        for (fname, f) in required_fields:
                            w('assert e.{} != null;', java.name(fname))

                        for (fname, f) in datamodel.sorted_fields(t):
                            jfname = java.name(fname)
                            if jfname in fkeys:
                                fk_type = fkeys[jfname].parent_path
                                with java.If(w, 'e.{} != null', jfname):
                                    w('HashSet<{0}> fk = getFk_{1}(e.{1});',
                                      tname, jfname)
                                    w('assert fk != null;')
                                    if pkey:
                                        w('assert fk.contains(e);')
                                    w('assert model.get{}Table().lookup(e.{}) != null;',
                                      fk_type, jfname)

                            if f.primitive == f.STRING:
                                for c in f.constraint:
                                    # TODO: Use length, not precision.
                                    if c.precision:
                                        # TODO: Special-case for required
                                        # fields.
                                        w('assert e.{0} == null || e.{0}.length() <= {1};',
                                          jfname, c.precision)

            with java.Method(w, '\npublic', 'boolean', 'equals', [('Object', 'obj')]):
                with java.If(w, 'this == obj'):
                    w('return true;')
                with java.If(w, '!(obj instanceof {})', tname):
                    w('return false;')
                w('Table that = (Table)obj;')
                w(('assert model == that.model : '
                   '"inter-model {}.Table.equals() undefined";'),
                  tname)

                if pkey:
                    w('return this.items.equals(that.items);')
                else:
                    # HashSet usage is borked because entities are inserted before
                    # their fields are populated. We skirt around this here by building
                    # temporary HashSets for comparison purposes.
                    # TODO: Revisit after introducing transactional updates.
                    for v in ['a', 'b']:
                        w('HashSet<{0}> {1} = new HashSet<{0}>();', tname, v)
                        # Dunno if clone() unborks, so just elt-wise copy.
                        with java.For(w, '{} e : this', tname):
                            w('{}.add(e);', v)
                    w('return a.equals(b);')

            with java.Method(w, '\npublic', 'int', 'hashCode'):
                w('int h = 0;')
                with java.For(w, '{} e : this', tname):
                    w('h ^= e.hashCode();')
                w('return h;')

            w()


def export_model_class(fk_rmap, context):
    '''The Model class'''

    # Create a new "Writer" that will be passed around (especially to the "java"
    # methods) and be updated with the Java cass text
    w = writer.Writer('java')

    java.Package(w, context.package)

    java.StandardImports(w)
    java.Import(w, 'java.lang.reflect.InvocationTargetException')
    w()
    java.Import(w, 'org.joda.time.DateTime')
    java.Import(w, 'org.joda.time.LocalDate')

    is_abstract = any(
        'abstract' in syslx.patterns(v.attrs)
        for (_, v) in context.app.views.iteritems())  # TODO: what is this structure?

    has_tables = any(
        t.WhichOneof('type') == 'relation'
        for (_, t) in context.app.types.iteritems())

    w()
    model_version = syslx.View(context.app.attrs)['version'].s
    if model_version:
        w('@io.sysl.Version("{}")', model_version)
    with java.Class(w, context.model_class, context.write_file,
                    package=context.package, extends='io.sysl.ModelView',
                    abstract=is_abstract):
        w()
        for (tname, t) in sorted(context.app.types.iteritems()):
            if t.WhichOneof('type') == 'relation':
                w('{0}.Table tbl_{0};', tname)

        for (tname, t) in sorted(context.app.types.iteritems()):
            if t.HasField('relation'):
                with java.Method(w, '\npublic', tname + '.Table',
                                 'get' + tname + 'Table'):
                    with java.If(w, 'tbl_{} == null', tname):
                        w('tbl_{0} = new {0}.Table(this);', tname)
                    w('return tbl_{};', tname)

        with java.Ctor(w, '\npublic', context.model_class):
            pass

        if has_tables:
            with java.Method(w, '\npublic', 'String', 'toString'):
                w('StringBuilder sb = new StringBuilder();')
                w('canonicalToString(sb, "", null);')
                w('return sb.toString();')

            with java.Method(w, '\npublic', 'String', 'canonicalToString'):
                w('StringBuilder sb = new StringBuilder();')
                w('int[] ids = new int[nextId];')
                w('canonicalToString(sb, "", ids);')
                w('return sb.toString();')

            with java.Method(w, '\nprivate', 'void', 'canonicalToString',
                             [('StringBuilder', 'sb'), ('String', 'indent'), ('int[]', 'ids')]):
                w('sb.append(indent);', tname)
                w('sb.append("{} {{\\n");', context.model_class)
                for (tname, t) in sorted(context.app.types.iteritems()):
                    if t.WhichOneof('type') == 'relation':
                        with java.If(w, 'tbl_{0} != null && !tbl_{0}.isEmpty()', tname):
                            w('sb.append(indent);')
                            w('sb.append("  table of {} ");', tname)
                            w('tbl_{}.canonicalToString(sb, indent + "  ", ids);', tname)
                            w('sb.append(",\\n");')
                w('sb.append(indent);', tname)
                w('sb.append("}}");')

        export_view_class_body(w, context)

        for (tname, t) in sorted(context.app.types.iteritems()):
            if re.match(r'AnonType_\d+__$', tname):
                export_entity_class(w, tname, t, fk_rmap[tname],
                                    context._replace(package=None))

        with w.transaction() as rollback:
            empty = True
            with java.Method(w, '\npublic', 'void', 'validate'):
                for (tname, t) in sorted(context.app.types.iteritems()):
                    if t.WhichOneof('type') == 'relation':
                        pkey = datamodel.primary_key_params(t, context.module)
                        pkey_fields = {f for (_, f, _) in pkey}
                        param_defs = [(typ, jfname)
                                      for (typ, fname, jfname) in pkey
                                      if ('autoinc' not in
                                          syslx.patterns(t.relation.attr_defs[fname].attrs))]

                        required_fields = [
                            (fname, f)
                            for (fname, f) in datamodel.sorted_fields(t)
                            if 'required' in syslx.patterns(f.attrs)]
                        fkeys = {
                            java.name(fname): type_info
                            for (fname, _, type_info) in (
                                datamodel.foreign_keys(t, context.module))}
                        if required_fields or fkeys:
                            with java.If(w, 'tbl_{} != null', tname):
                                w('tbl_{}.validate();', tname)
                            empty = False
            if empty:
                rollback()

        with java.Method(w, '\npublic', 'boolean', 'equals',
                         [('Object', 'obj')], override=True):
            with java.If(w, 'this == obj'):
                w('return true;')
            with java.If(w, '!(obj instanceof {})', context.model_class):
                w('return false;')
            w('{0} that = ({0})obj;', context.model_class)
            w('int[] thisIds = new int[nextId];')
            w('int[] thatIds = new int[that.nextId];')
            w('return equals(this, thisIds, that, thatIds);')

        with java.Method(w, '\nprivate', 'boolean', 'equals',
                         [(context.model_class, 'a'), ('int[]', 'aIds'),
                          (context.model_class, 'b'), ('int[]', 'bIds')],
                         static=True):
            with java.If(w, 'a == b'):
                w('return true;')

            for (tname, t) in datamodel.fk_topo_sort(
                    context.app.types, context.module):
                if t.WhichOneof('type') == 'relation':
                    with java.If(w, 'a.tbl_{0} == null', tname):
                        with java.If(w, 'b.tbl_{0} != null', tname):
                            w('return false;')
                        with java.ElseIf(w,
                                         '!{0}.Table.canonicallyEqual(a.tbl_{0}, aIds, b.tbl_{0}, bIds)',
                                         tname):
                            w('return false;')
                w()
            w('return true;')

        with java.Method(w, '\npublic', 'int', 'hashCode'):
            # Lets go for a large odd pseudo-random multiplier seeded by the
            # type.
            h = hashlib.sha256(
                '\0'.join([context.package, context.model_class]))
            m = struct.unpack('>i', h.digest()[-4:])[0] & 0x7fffffff | 1
            w('final int m = {};', m)
            w('int h = 0;', m)

            w('int[] ids = new int[nextId];')
            for (tname, t) in datamodel.fk_topo_sort(
                    context.app.types, context.module):
                if t.WhichOneof('type') == 'relation':
                    w(('h = m * h + (tbl_{0} == null ? 0 :'
                       ' tbl_{0}.canonicalHashCode(ids));'),
                      tname)

            w('return h;')

        with java.Method(w, '\npublic', 'void', '_PRIVATE_registerId',
                         [('int', 'id')]):
            with java.If(w, 'nextId <= id'):
                w('nextId = id + 1;')

        with java.Method(w, '\npublic', 'int', '_PRIVATE_getNextId'):
            w('return nextId;')

        w()
        w('int nextId = 1;')


def export_view_class(w, context):
    w()

    is_abstract = any(
        'abstract' in syslx.patterns(v.attrs)
        for (_, v) in context.app.views.iteritems())

    with java.Class(w, context.model_class, context.write_file,
                    package=context.package, extends='io.sysl.ModelView',
                    abstract=is_abstract):
        export_view_class_body(w)


def export_view_class_body(w, context):
    global_scope = scopes.Scope(context.module)
    for (vname, v) in sorted(context.app.views.iteritems()):
        global_scope[vname] = v.ret_type

    for (vname, v) in sorted(context.app.views.iteritems()):
        w()
        scope = scopes.Scope(
            global_scope,
            _app=context.app,
            **{param.name: param.type for param in v.param})

        (ret_app, ret_type) = global_scope.resolve(v.ret_type)
        ret_type_code = java.codeForType(v.ret_type, global_scope)

        params = [
            ('final ' + java.codeForType(param.type, global_scope),
             java.safe(param.name))
            for param in v.param]

        use_out_param = 'partial' in syslx.patterns(v.attrs)
        is_model = isinstance(ret_type, sysl_pb2.Application)
        is_relation = (
            not is_model and
            (ret_type.WhichOneof('type') == 'relation' or
             (ret_type.WhichOneof('type') == 'set' and
                ret_type.set.WhichOneof('type') == 'relation')))

        if use_out_param:
            params.append(('final ' + ret_type_code, 'out'))
            ret_type_code = 'void'

        is_abstract = 'abstract' in syslx.patterns(v.attrs)
        assert (v.expr.WhichOneof('expr') == 'transform') ^ is_abstract, (
            'View must be abstract or root expr must be transform')

        transform = v.expr.transform

        is_view = not is_model and ret_type.WhichOneof('type') == 'set'
        out_type = v.ret_type.set if is_view else v.ret_type
        out_type_code = java.codeForType(out_type, scope)
        out_package = ''.join(out_type_code.rpartition('.')[:2])
        scope['out'] = out_type

        (_, resolved_out_type) = scope.resolve(out_type)
        is_primitive = out_type.WhichOneof('type') == 'primitive'
        pkey = ([] if is_model or is_primitive else
                datamodel.primary_key_params(resolved_out_type, context.module))
        pkey_fields = {jfname for (_, _, jfname) in pkey}

        fkeys = (set() if is_model or is_primitive else {
            java.name(fname): type_info
            for (fname, _, type_info) in datamodel.foreign_keys(
                resolved_out_type, context.module)
        })
        fkey_fields = set(fkeys)

        special_fields = pkey_fields | fkey_fields

        if is_abstract:
            w('public abstract {} {}({});', ret_type_code, vname,
              ', '.join('{} {}'.format(t, p) for (t, p) in params))
            continue

        with java.Method(w, 'public', ret_type_code, vname, params):
            if is_view:
                w('{0}.Set result = new {0}.Set();', out_type_code)

            (item_access, item_type) = java.codeForExpr(
                w, transform.arg, scope, context.module)
            if not is_view and item_access == param.name:
                item_type = param.type

            if transform.scopevar and transform.scopevar != '.':
                transform_scope = scopes.Scope(scope, __dot__=transform.scopevar,
                                               **{transform.scopevar: item_type})
            else:
                transform_scope = scopes.Scope(
                    scope, item_type, __dot__='item')

            def transform_item():
                w()
                with w.table():
                    for stmt in transform.stmt:
                        which_stmt = stmt.WhichOneof('stmt')

                        if which_stmt == 'assign':
                            assign = stmt.assign
                            if assign.name in special_fields:
                                set_prefix = '_PRIVATE_'
                            else:
                                set_prefix = ''

                            if assign.expr.WhichOneof('expr'):
                                (code, _) = java.codeForExpr(
                                    w, assign.expr, transform_scope, context.module)

                                if is_model:
                                    with java.Block(w):
                                        w('{0}{1}.Table table = \vout.get{1}Table();',
                                          out_package, assign.name)
                                        with java.For(w, '{}{} e : \v{}',
                                                      out_package, assign.name, code):
                                            w('table._PRIVATE_insert(e);')
                                else:
                                    w('out.{}set{}\037 ({});',
                                      set_prefix, java.CamelCase(assign.name), code)
                            else:
                                w('// out.{}set{}(...);',
                                  set_prefix, java.CamelCase(assign.name))
                        elif which_stmt == 'let':
                            let = stmt.let
                            if let.expr.WhichOneof('expr'):
                                (code, type_) = java.codeForExpr(
                                    w, let.expr, transform_scope, context.module, let=let.name)
                                transform_scope[let.name] = type_
                                if code:
                                    w('final {} {} = \v{};',
                                      java.codeForType(type_, transform_scope), let.name, code)
                            else:
                                w('// <type> {} = ...;', code)
                        elif which_stmt == 'inject':
                            w('{};', java.codeForExpr(
                              w, stmt.inject, transform_scope, context.module)[0])
                if is_view:
                    w('\nresult.add(out);')

            if is_view:
                with java.For(w, 'final {} {} : {}',
                              java.codeForType(item_type.set, transform_scope),
                              transform_scope['__dot__'],
                              item_access):
                    if not use_out_param:
                        if is_relation:
                            w('{0} out = \v{0}._PRIVATE_new();', out_type_code)
                        else:
                            w('{0} out = \vnew {0}();', out_type_code)
                    transform_item()
            else:
                if re.search(r'\bpackage\b', item_access):
                    ta = transform.arg
                    java.codeForExpr(w, ta, transform_scope, context.module)
                w('final {} item = {};', java.codeForType(item_type, transform_scope),
                  item_access)
                if not use_out_param:
                    if out_type_code.endswith('.View') or not is_relation:
                        w('{0} out = \vnew {0}();',
                          re.sub(r'\.View$', '.Set', out_type_code))
                    else:
                        w('{0} out = \v{0}._PRIVATE_new();', out_type_code)
                transform_item()
            if not use_out_param:
                if is_view:
                    w('\nreturn result;')
                else:
                    w('\nreturn out;')

    w()


def export_exception_class(context):
    w = writer.Writer('java')

    java.Package(w, context.package)

    java.Import(w, 'java.lang.RuntimeException')

    with java.Class(w, context.model_class + 'Exception', context.write_file,
                    extends='RuntimeException',
                    package=context.package):
        with java.Ctor(w, '\npublic', context.model_class + 'Exception'):
            pass
        with java.Ctor(w, '\npublic', context.model_class + 'Exception',
                       [('String', 'message')]):
            w('super(message);')
        w()
