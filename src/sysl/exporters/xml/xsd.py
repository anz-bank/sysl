import contextlib

from sysl.proto import sysl_pb2

from sysl.core import syslx

from sysl.util import datamodel
from sysl.util import java
from sysl.util import writer


XSD_TYPE_MAP = {
    sysl_pb2.Type.ANY: None,
    sysl_pb2.Type.BOOL: u'xs:boolean',
    sysl_pb2.Type.INT: u'xs:int',
    sysl_pb2.Type.FLOAT: None,
    sysl_pb2.Type.DECIMAL: u'xs:decimal',
    sysl_pb2.Type.STRING: u'xs:string',
    sysl_pb2.Type.BYTES: None,
    sysl_pb2.Type.STRING_8: u'xs:string',
    sysl_pb2.Type.DATE: u'xs:date',
    sysl_pb2.Type.DATETIME: u'xs:dateTime',
    sysl_pb2.Type.XML: u'xs:string',
    sysl_pb2.Type.UUID: u'xs:string',
}


def xsd(context):

    w = writer.Writer('xml')
    w.increment = 2
    xsd_separator = '<!-- ' + '=' * 55 + ' -->'

    def e(_name, **attrs):
        if _name.endswith('/'):
            w(u'<{}{}/>',
              _name[:-1],
              u''.join(u' {}="{}"'.format(k.replace('_', ':'), v)
                       for (k, v) in attrs.iteritems()))
        else:
            @contextlib.contextmanager
            def f():
                w(u'<{}{}>',
                  _name,
                  u''.join(u' {}="{}"'.format(k.replace('_', ':'), v)
                           for (k, v) in attrs.iteritems()))
                with w.indent():
                    yield
                w(u'</{}>', _name)
            return f()

    def xs(_name, **attrs):
        return e('xs:' + _name, **attrs)

    def build_element(attr_fname, attr_f, is_set, is_attr):
        jfname = java.name(attr_fname)
        method = java.CamelCase(jfname)
        type_ = datamodel.typeref(attr_f, context.module)[2]
        which_type = type_.WhichOneof('type')
        if which_type == 'primitive':
            xsdtype = XSD_TYPE_MAP[type_.primitive]
        elif which_type == 'enum':
            xsdtype = 'xs:int'
        elif which_type == 'tuple':
            offset = -1
            if is_set:
                offset = -2
            xsdtype = datamodel.typeref(attr_f, context.module)[
                0].split('.')[offset]
        else:
            raise RuntimeError(
                'Unexpected field type for XSD '
                'export: ' + which_type)
        if is_attr:
            xs('attribute/', name=jfname, type=xsdtype, use='optional')
        else:
            xs('element/', name=jfname, type=xsdtype, minOccurs=0)

    def build_relational_xsd():
        # top level element contains all entities for
        # relational schemas but will only contain the
        # root entity for hierarchical
        with xs('element', name=context.model_class):
            with xs('complexType'):
                with xs('sequence', minOccurs=1, maxOccurs=1):
                    # each "relation" is a list of things
                    for (tname, ft, t) in syslx.sorted_types(context):
                        if t.HasField('relation'):
                            xs('element/', name=tname + 'List', type=tname + 'List',
                               minOccurs=0)

            # build keys and key refs
            for (tname, ft, t) in syslx.sorted_types(context):
                if t.HasField('relation'):
                    pkey = datamodel.primary_key_params(t, context.module)
                    pkey_fields = {f for (_, f, _) in pkey}
                    has_content = False

                    def xsd_key_header(msg):
                        w('{}', '<!-- ' + msg.center(55) + ' -->')

                    if pkey:
                        if not has_content:
                            has_content = True
                            w()
                            w(xsd_separator)
                            xsd_key_header(tname + ' keys')

                        with xs('key', name='key_' + tname):
                            xs('selector/',
                               xpath='./{0}List/{0}'.format(tname))
                            for f in sorted(pkey_fields):
                                xs('field/', xpath=f)

                    for (i, (fname, _, type_info)) in enumerate(sorted(
                            datamodel.foreign_keys(t, context.module))):
                        if not has_content:
                            has_content = True
                            w()
                            w(xsd_separator)
                            xsd_key_header(tname + ' keyrefs')
                        elif pkey and i == 0:
                            xsd_key_header('keyrefs')
                        fk_type = type_info.parent_path
                        fk_field = type_info.field
                        with xs('keyref',
                                name='keyref_{}_{}'.format(tname, fname),
                                refer='key_' + fk_type):
                            xs('selector/',
                               xpath='./{0}List/{0}'.format(tname))
                            xs('field/', xpath=fname)
                    if has_content:
                        w(xsd_separator)
            w()

        w()
        # construct the entities
        for (tname, ft, t) in syslx.sorted_types(context):
            if t.HasField('relation'):
                with xs('complexType', name=tname + 'List'):
                    with xs('sequence', maxOccurs='unbounded'):
                        with xs('element', name=tname):
                            with xs('complexType'):
                                with xs('all'):
                                    for (fname, f) in sorted(
                                            t.relation.attr_defs.iteritems()):
                                        if 'xml_attribute' not in syslx.patterns(
                                                f.attrs):
                                            build_element(
                                                fname, f, False, False)

                                # attributes second
                                for (fname, f) in sorted(
                                        t.relation.attr_defs.iteritems()):
                                    if 'xml_attribute' in syslx.patterns(
                                            f.attrs):
                                        build_element(fname, f, False, True)

    def build_hierarchical_xsd():

        # Top level element
        with xs('element', name=context.model_class):
            with xs('complexType'):
                with xs('sequence', minOccurs=1, maxOccurs=1):
                    for (tname, ft, t) in syslx.sorted_types(context):
                        if 'xml_root' in syslx.patterns(t.attrs):
                            xs('element/', name=tname, type=tname, minOccurs=0)
                            break

        w(xsd_separator)

        for (tname, ft, t) in syslx.sorted_types(context):
            with xs('complexType', name=tname):

                with xs('all'):
                    for (fname, f) in sorted(t.tuple.attr_defs.iteritems()):
                        if 'xml_attribute' not in syslx.patterns(f.attrs):
                            if f.HasField('set'):
                                with xs('element', name=fname + 'List'):
                                    with xs('complexType'):
                                        with xs('sequence', maxOccurs='unbounded'):
                                            build_element(
                                                fname, f, True, False)
                            else:
                                build_element(fname, f, False, False)

                # attributes second
                for (fname, f) in sorted(t.tuple.attr_defs.iteritems()):
                    if 'xml_attribute' in syslx.patterns(f.attrs):
                        build_element(fname, f, False, True)

    def build_xsd():
        for (_, _, t) in syslx.sorted_types(context):
            if t.HasField('relation'):
                build_relational_xsd()
                return
        build_hierarchical_xsd()

    with e('xs:schema',
           xmlns_xs='http://www.w3.org/2001/XMLSchema',
           attributeFormDefault='unqualified',
           elementFormDefault='qualified',
           version='1.0'):
        build_xsd()

    context.write_file(w, context.model_class + '.xsd')
