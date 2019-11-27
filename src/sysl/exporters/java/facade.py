from sysl.core import syslx

from sysl.util import datamodel
from sysl.util import java
from sysl.util import writer


# Facade class does reference the Model class via the getModel public method.
def export_facade_class(context):
    model_name = syslx.fmt_app_name(context.wrapped_model.name)
    modelpkg = syslx.View(context.wrapped_model.attrs)['package'].s

    w = writer.Writer('java')

    java.Package(w, context.package)
    java.StandardImports(w)
    java.Import(w, modelpkg + '.' + model_name)

    w()
    with java.Class(w, context.appname, context.write_file,
                    package=context.package):
        with java.Ctor(w, '\npublic', context.appname, [(model_name, 'model')]):
            w('this.model = model;')
            for (tname, _, _) in syslx.wrapped_facade_types(context):
                w('this.{}Facade = new {}Facade();',
                  java.safe(tname[:1].lower() + tname[1:]),
                  tname)

        with java.Method(w, '\npublic', model_name, 'getModel'):
            w('return model;')

        java.SeparatorComment(w)

        for (tname, ft, t) in syslx.wrapped_facade_types(context):
            if t.HasField('relation'):
                pkey = datamodel.primary_key_params(t, context.module)
                pkey_fields = {f for (_, f, _) in pkey}
                param_defs = [
                    (typ, jfname) for (typ, fname, jfname) in pkey
                    if 'autoinc' not in syslx.patterns(t.relation.attr_defs[fname].attrs)]
                params = ''.join(', ' + f for (_, f) in param_defs)
                param = ', '.join(f for (_, f) in param_defs)

                fkeys = {
                    java.name(fname): type_info
                    for (fname, _, type_info) in datamodel.foreign_keys(t, context.module)}

                inner_type = ('HashMap<Key, {}>' if pkey else 'HashSet<{}>'
                              ).format(tname)

                required = []
                for fname in sorted(ft.relation.attr_defs.keys()):
                    f = t.relation.attr_defs.get(fname)
                    if ('required' in syslx.patterns(f.attrs) or
                        'required' in syslx.patterns(
                            ft.relation.attr_defs.get(fname).attrs)):
                        jfname = java.name(fname)
                        method = java.CamelCase(jfname)
                        (java_type, type_info, _) = datamodel.typeref(
                            f, context.module)
                        if java_type == 'Object':
                            datamodel.typeref(f, context.module)
                        required.append((fname, java_type))

                with java.Method(w, '\npublic', tname + 'Facade', 'get' + tname):
                    w('return {}Facade;', java.safe(
                        tname[:1].lower() + tname[1:]))

                w()
                with java.Class(w, tname + 'Facade', context.write_file):
                    with java.Method(w, 'public', '{}.{}.Table'.format(modelpkg, tname),
                                     'table'):
                        w('return model.get{}Table();', tname, param)

                    w()
                    if param_defs or required:
                        with java.Method(w, 'public', 'Builder0', 'build'):
                            w('return new Builder0();')

                        keytypes = {f: kt for (kt, f) in param_defs}
                        keys = sorted(keytypes)
                        keyindices = {k: i for (i, k) in enumerate(keys)}

                        if len(keys) > 3:
                            # 4 perms yields 16 builders with 32 setters.
                            logging.error('OUCH! Too many primary key fields')
                            raise Exception('Too many primary key fields')
                        for perm in range(2**len(keys)):
                            bname = 'Builder' + str(perm)
                            w()
                            with java.Class(w, bname, context.write_file):
                                done = [k for (i, k) in enumerate(
                                    keys) if 2**i & perm]
                                remaining = [k for k in keys if k not in done]

                                with java.Ctor(w, '', bname,
                                               [(keytypes[k], k) for k in done]):
                                    for k in done:
                                        w('this.{0} = {0};', k)
                                    if required and not remaining:
                                        w('this._pending = {};',
                                          hex(2**len(required) - 1).rstrip('L'))

                                for fname in remaining:
                                    f = t.relation.attr_defs[fname]
                                    jfname = java.name(fname)
                                    method = java.CamelCase(jfname)
                                    (java_type, type_info, _) = datamodel.typeref(
                                        f, context.module)
                                    next_bname = 'Builder{}'.format(
                                        perm | 2**keyindices[fname])
                                    w()

                                    if jfname in fkeys:
                                        fk_type = fkeys[jfname].parent_path
                                        fk_field = fkeys[jfname].field
                                        if f.type_ref.ref.path[-1:] == [fname]:
                                            method_suffix = fk_type
                                        else:
                                            method_suffix = method + 'From'
                                        with java.Method(w, 'public', next_bname,
                                                         'with' + method_suffix,
                                                         [(modelpkg + '.' + fk_type, 'entity')]):
                                            w('{} {} = entity == null ? null : entity.get{}();',
                                              java_type, jfname, java.CamelCase(fk_field))
                                            w('return new {}({});',
                                              next_bname,
                                              ', '.join(k for k in keys if k == fname or k in done))
                                    else:
                                        with java.Method(w, 'public', next_bname,
                                                         'with' +
                                                         java.CamelCase(fname),
                                                         [(keytypes[fname], fname)]):
                                            w('return new {}({});',
                                              next_bname,
                                              ', '.join(k for k in keys if k == fname or k in done))

                                if not remaining:
                                    for (i, (r, rtype)) in enumerate(required):
                                        method = java.CamelCase(r)
                                        w()

                                        # TODO: jfname set in a previous loop??
                                        if jfname in fkeys:
                                            fk_type = fkeys[jfname].parent_path
                                            fk_field = fkeys[jfname].field
                                            if f.type_ref.ref.path[-1:] == [fname]:
                                                method = fk_type
                                            else:
                                                method += 'From'
                                            with java.Method(w, 'public', bname, 'with' + method,
                                                             [(modelpkg + '.' + fk_type,
                                                               java.mixedCase(fk_type))]):
                                                with java.If(w, '(_pending & {}) == 0', hex(2**i)):
                                                    # TODO: More specific
                                                    # exception
                                                    w('throw new RuntimeException();')
                                                w('this.{0} = {0};',
                                                  java.mixedCase(fk_type))
                                                w('_pending &= ~{};', hex(2**i))
                                                w('return this;')
                                        else:
                                            with java.Method(w, 'public', bname, 'with' + method,
                                                             [(rtype, r)]):
                                                with java.If(w, '(_pending & {}) == 0', hex(2**i)):
                                                    # TODO: More specific
                                                    # exception
                                                    w('throw new RuntimeException();')
                                                w('this.{0} = {0};', r)
                                                w('_pending &= ~{};', hex(2**i))
                                                w('return this;')

                                    with java.Method(w, '\npublic', modelpkg + '.' + tname,
                                                     'insert'):
                                        if required:
                                            with java.If(w, '_pending != 0'):
                                                # TODO: More specific exception
                                                w('throw new RuntimeException();')
                                            w(u'{} result = table()._PRIVATE_insert({});',
                                              modelpkg + '.' + tname, param)
                                            for (r, rtype) in required:
                                                if jfname in fkeys:
                                                    fk_field = fkeys[jfname].field
                                                    w('result.set{}({});',
                                                      fk_type, java.mixedCase(fk_type))
                                                else:
                                                    w('result.set{}({});',
                                                      java.CamelCase(r), r)
                                            w('return result;')
                                        else:
                                            w(u'return table()._PRIVATE_insert({});', param)

                                if done:
                                    w()
                                    w('// primary key')
                                    for d in done:
                                        w('private final {} {};',
                                          keytypes[d], d)

                                if required and not remaining:
                                    w()
                                    w('// required fields')
                                    for (r, rtype) in required:
                                        if jfname in fkeys:
                                            fk_type = fkeys[jfname].parent_path
                                            fk_field = fkeys[jfname].field
                                            w('private {} {};',
                                              modelpkg + '.' + fk_type, java.mixedCase(fk_type))
                                        else:
                                            w('private {} {};', rtype, r)
                                    w()
                                    w('private int _pending;')
                    else:
                        with java.Method(w, 'public', modelpkg + '.' + tname, 'insert'):
                            w('{} result = table()._PRIVATE_insert();',
                              modelpkg + '.' + tname)
                            w('return result;')

                java.SeparatorComment(w)

        for (tname, _, t) in syslx.wrapped_facade_types(context):
            if t.HasField('relation'):
                w('private final {}Facade {}Facade;',
                  tname,
                  java.safe(tname[:1].lower() + tname[1:]))

        w()
        w('private final {}.{} model;', modelpkg, model_name)
        w()
