from sysl.proto import sysl_pb2

from sysl.core import syslx

from sysl.util import datamodel
from sysl.util import java
from sysl.util import writer


XML_GEN_MAP = {
    sysl_pb2.Type.ANY: '',
    sysl_pb2.Type.BOOL: '',
    sysl_pb2.Type.INT: '',
    sysl_pb2.Type.FLOAT: '',
    sysl_pb2.Type.DECIMAL: '',
    sysl_pb2.Type.STRING: '',
    sysl_pb2.Type.BYTES: '',
    sysl_pb2.Type.STRING_8: '',
    sysl_pb2.Type.DATE: '',
    sysl_pb2.Type.DATETIME: ', iso8601DateTime',
    sysl_pb2.Type.XML: '',
    sysl_pb2.Type.UUID: '',
}

XML_PARSE_MAP = {
    sysl_pb2.Type.ANY: None,
    sysl_pb2.Type.BOOL: 'BooleanValue',
    sysl_pb2.Type.INT: 'IntValue',
    sysl_pb2.Type.FLOAT: 'DoubleValue',
    sysl_pb2.Type.DECIMAL: 'DecimalValue',
    sysl_pb2.Type.STRING: 'Text',
    sysl_pb2.Type.BYTES: 'BinaryValue',
    sysl_pb2.Type.STRING_8: 'Text',
    sysl_pb2.Type.DATE: ('Text', 'iso8601DateTime'),
    sysl_pb2.Type.DATETIME: ('Text', 'iso8601DateTime'),
    sysl_pb2.Type.XML: 'Text',
    sysl_pb2.Type.UUID: 'Text',
}


def serializer(context):
    (app, module, package, model_class, write_file, _, _) = context

    facade = bool(context.wrapped_model)

    w = writer.Writer('java')

    java.Package(w, package)

    java.StandardImports(w)

    java.Import(w, 'java.text.ParseException')
    w.head()
    java.Import(w, 'javax.xml.stream.XMLStreamException')
    java.Import(w, 'javax.xml.stream.XMLStreamWriter')
    w.head()
    java.Import(w, 'org.joda.time.format.DateTimeFormatter')
    java.Import(w, 'org.joda.time.format.ISODateTimeFormat')

    if facade:
        model_name = syslx.fmt_app_name(context.wrapped_model.name)
        modelpkg = syslx.View(context.wrapped_model.attrs)['package'].s
        w.head()
        java.Import(w, modelpkg + '.*')
    else:
        modelpkg = package

    with java.Class(w, '\n{}XmlSerializer'.format(model_class), write_file,
                    package=package):
        with java.Method(w, '\npublic', 'void', 'serialize',
                         [(model_class, 'facade' if facade else 'model'),
                          ('XMLStreamWriter', 'xsw')],
                         throws=['XMLStreamException']):
            if facade:
                w('{} model = facade.getModel();', model_name)

            w('xsw.writeStartElement("{}");', model_class)
            for (tname, ft, t) in syslx.sorted_types(context):
                if t.WhichOneof('type') == 'relation':
                    w('serialize(model.get{0}Table(), xsw, "{0}");', tname)
            w('xsw.writeEndElement();')

        for (tname, ft, t) in syslx.sorted_types(context):
            java.SeparatorComment(w)

            if t.WhichOneof('type') in ['relation', 'tuple']:
                with java.Method(w, '\npublic', 'void', 'serialize',
                                 [(modelpkg + '.' + tname + '.View', 'view'),
                                  ('XMLStreamWriter', 'xsw'),
                                  ('String', 'tag')],
                                 throws=['XMLStreamException']):
                    with java.If(w, 'view == null || view.isEmpty()'):
                        w('return;')
                    if t.WhichOneof('type') == 'relation':
                        w('xsw.writeStartElement("{}List");', tname)
                    order = [o.s for o in syslx.View(
                        t.attrs)['xml_order'].a.elt]
                    if order:
                        order_logic = []
                        for o in order:
                            order_logic.append(
                                '            i = io.sysl.ExprHelper.doCompare'
                                '(a.get{0}(), b.get{0}());\n'
                                '            if (i != 0) return i;\n'.format(java.CamelCase(o)))
                        order_by = (
                            '.orderBy(\n'
                            '    new java.util.Comparator<{0}> () {{\n'
                            '        @Override\n'
                            '        public int compare({0} a, {0} b) {{\n'
                            '            int i;\n'
                            '{1}'
                            '            return 0;\n'
                            '        }}\n'
                            '    }})\n').format(tname, ''.join(order_logic))
                    else:
                        order_by = ''

                    with java.For(w, '{} item : view{}',
                                  modelpkg + '.' + tname, order_by):
                        w('serialize(item, xsw, tag);', tname)
                    if t.WhichOneof('type') == 'relation':
                        w('xsw.writeEndElement();')

                with java.Method(w, '\npublic', 'void', 'serialize',
                                 [(modelpkg + '.' + tname, 'entity'),
                                  ('XMLStreamWriter', 'xsw')],
                                 throws=['XMLStreamException']):
                    w('serialize(entity, xsw, "{}");', tname)

                with java.Method(w, '\npublic', 'void', 'serialize',
                                 [(modelpkg + '.' + tname, 'entity'),
                                  ('XMLStreamWriter', 'xsw'),
                                  ('String', 'tag')],
                                 throws=['XMLStreamException']):
                    with java.If(w, 'entity == null'):
                        w('return;')
                    w('xsw.writeStartElement(tag);', tname)
                    for wantAttrs in [True, False]:
                        for (fname, f) in datamodel.sorted_fields(t):
                            jfname = java.name(fname)
                            method = java.CamelCase(jfname)
                            tref = datamodel.typeref(f, module)
                            type_ = tref[2]
                            if not type_:
                                continue
                            extra = ''
                            which_type = type_.WhichOneof('type')
                            if which_type == 'primitive':
                                extra = XML_GEN_MAP[type_.primitive]
                                if type_.primitive == type_.DECIMAL:
                                    for c in type_.constraint:
                                        if c.scale:
                                            access = 'round(entity.get{}(), {})'.format(
                                                method, c.scale)
                                            break
                                    else:
                                        access = 'round(entity.get{}(), 0)'.format(
                                            method)
                                else:
                                    access = 'entity.get{}()'.format(method)
                            elif which_type == 'enum':
                                access = 'entity.get{}().getValue()'.format(method)
                            elif which_type == 'tuple':
                                access = 'entity.get{}()'.format(method)
                            else:
                                raise RuntimeError(
                                    'Unexpected field type for XML export: ' + which_type)
                            if wantAttrs == (
                                    'xml_attribute' in syslx.patterns(f.attrs)):
                                if wantAttrs:
                                    w('serializeAttr("{}", {}{}, xsw);',
                                      jfname, access, extra)
                                else:
                                    # TODO: Why is 'index' bad? (NB: This was a debug hook, but the author forgot why.)
                                    assert jfname != 'index'
                                    w('serializeField("{}", {}{}, xsw);',
                                      jfname, access, extra)
                    w('xsw.writeEndElement();')

                with java.Method(w, '\npublic', 'void', 'serializeField',
                                 [('String', 'fieldname'),
                                  (modelpkg + '.' + tname + '.View', 'view'),
                                  ('XMLStreamWriter', 'xsw')],
                                 throws=['XMLStreamException']):
                    with java.If(w, 'view != null && !view.isEmpty()'):
                        w('serialize(view, xsw, fieldname);')

                with java.Method(w, '\npublic', 'void', 'serializeField',
                                 [('String', 'fieldname'), (modelpkg + '.' + tname, 'entity'),
                                  ('XMLStreamWriter', 'xsw')],
                                 throws=['XMLStreamException']):
                    with java.If(w, 'entity != null'):
                        w('serialize(entity, xsw, fieldname);')

        def serialize(t, access, extra_args=None):
            with java.Method(w, '\nprivate', 'void', 'serialize',
                             ([(t, 'value')] + (extra_args or []) +
                              [('XMLStreamWriter', 'xsw')]),
                             throws=['XMLStreamException']):
                w('xsw.writeCharacters({});', access)

        serialize('String', 'value')
        serialize('Boolean', 'value.toString()')
        serialize('Integer', 'value.toString()')
        serialize('Double', 'value.toString()')
        serialize('BigDecimal', 'value.toString()')
        serialize('UUID', 'value.toString()')
        serialize('DateTime', 'fmt.print(value)',
                  [('DateTimeFormatter', 'fmt')])
        serialize('LocalDate', 'iso8601Date.print(value)')

        def serializeField(t, access, extra_args=None):
            with java.Method(w, '\nprivate', 'void', 'serializeField',
                             ([('String', 'fieldname'), (t, 'value')] + (extra_args or []) +
                              [('XMLStreamWriter', 'xsw')]),
                             throws=['XMLStreamException']):
                with java.If(w, 'value != null'):
                    w('xsw.writeStartElement(fieldname);')
                    w('serialize(value{}, xsw);',
                      ''.join(', ' + arg for (_, arg) in extra_args or []))
                    w('xsw.writeEndElement();')

            with java.Method(w, '\nprivate', 'void', 'serializeAttr',
                             ([('String', 'fieldname'), (t, 'value')] + (extra_args or []) +
                              [('XMLStreamWriter', 'xsw')]),
                             throws=['XMLStreamException']):
                with java.If(w, 'value != null'):
                    w('xsw.writeAttribute(fieldname, {});', access)

        serializeField('String', 'value')
        serializeField('Boolean', 'value.toString()')
        serializeField('Integer', 'value.toString()')
        serializeField('Double', 'value.toString()')
        serializeField('BigDecimal', 'value.toString()')
        serializeField('UUID', 'value.toString()')
        serializeField('DateTime', 'fmt.print(value)',
                       [('DateTimeFormatter', 'fmt')])
        serializeField('LocalDate', 'iso8601Date.print(value)')

        with java.Method(w, '\nprivate final', 'BigDecimal', 'round',
                         [('BigDecimal', 'd'), ('int', 'scale')]):
            w('return d == null ? d : '
              'd.setScale(scale, java.math.RoundingMode.HALF_UP);')

        w('\nprivate final DateTimeFormatter iso8601Date = '
          'ISODateTimeFormat.date();')
        w('private final DateTimeFormatter iso8601DateTime = '
          'ISODateTimeFormat.dateTime();')


def deserializer(context):
    (app, module, package, model_class, write_file, _, _) = context

    facade = bool(context.wrapped_model)

    w = writer.Writer('java')

    java.Package(w, package)

    java.StandardImports(w)

    java.Import(w, 'javax.xml.stream.XMLStreamConstants')
    java.Import(w, 'javax.xml.stream.XMLStreamException')
    java.Import(w, 'javax.xml.stream.XMLStreamReader')
    w.head()
    java.Import(w, 'java.text.ParseException')
    w.head()
    java.Import(w, 'org.joda.time.format.DateTimeFormatter')
    java.Import(w, 'org.joda.time.format.ISODateTimeFormat')

    if facade:
        model_name = syslx.fmt_app_name(context.wrapped_model.name)
        modelpkg = syslx.View(context.wrapped_model.attrs)['package'].s
        w.head()
        java.Import(w, modelpkg + '.*')
    else:
        model_name = model_class
        modelpkg = package

    with java.Class(w, '\n{}XmlDeserializer'.format(model_class), write_file,
                    package=package):
        with java.Method(w, 'public', 'void', 'deserialize',
                         [(model_class, 'facade' if facade else 'model'),
                          ('XMLStreamReader', 'xsr')],
                         throws=[model_name + 'Exception', 'XMLStreamException']):
            if facade:
                w('{} model = facade.getModel();', model_name)
            w('expect(XMLStreamConstants.START_ELEMENT, xsr.next());')
            with java.While(w, 'xsr.next() == XMLStreamConstants.START_ELEMENT'):
                with java.Switch(w, 'xsr.getLocalName()'):
                    for (tname, ft, t) in syslx.sorted_types(context):
                        if t.HasField('relation'):
                            w('case "{0}List": deserialize(model.get{0}Table(), xsr); break;',
                              tname)
                            w('case "{0}": deserializeOne(model.get{0}Table(), xsr); break;',
                              tname)
                    w('default: skipElement(xsr);')
            w('expect(XMLStreamConstants.END_ELEMENT, xsr.getEventType());')

        for (tname, ft, t) in syslx.sorted_types(context):
            if not t.HasField('relation'):
                continue
            pkey = datamodel.primary_key_params(t, context.module)
            pkey_fields = {f for (_, f, _) in pkey}
            fkeys = {
                java.name(fname): type_info
                for (fname, _, type_info) in datamodel.foreign_keys(t, context.module)}

            private_setters = pkey_fields | set(fkeys.iterkeys())

            with java.Method(w, '\nprivate', 'void', 'deserialize',
                             [(modelpkg + '.' + tname + '.Table', 'table'),
                              ('XMLStreamReader', 'xsr')],
                             throws=['XMLStreamException']):
                with java.While(w, 'xsr.next() == XMLStreamConstants.START_ELEMENT'):
                    w('deserializeOne(table, xsr);')
                w('expect(XMLStreamConstants.END_ELEMENT, xsr.getEventType());')

            with java.Method(w, '\nprivate', 'void', 'deserializeOne',
                             [(modelpkg + '.' + tname + '.Table', 'table'),
                              ('XMLStreamReader', 'xsr')],
                             throws=['XMLStreamException']):
                w('{0} entity = {0}._PRIVATE_new(table.model());',
                  modelpkg + '.' + tname)
                with java.While(w, 'xsr.next() == XMLStreamConstants.START_ELEMENT'):
                    with java.Switch(w, 'xsr.getLocalName()'):
                        for (fname, f) in datamodel.sorted_fields(t):
                            jfname = java.name(fname)
                            (typename, _, type_) = datamodel.typeref(f, module)
                            extra = '{}'
                            which_type = type_.WhichOneof('type')
                            if which_type == 'primitive':
                                xmltype = XML_PARSE_MAP[type_.primitive]
                                if isinstance(xmltype, tuple):
                                    (xmltype, extra) = xmltype
                            elif which_type == 'enum':
                                xmltype = 'IntValue'
                                extra = typename + '.from({})'
                            else:
                                raise RuntimeError(
                                    'Unexpected field type for XML export: ' + which_type)

                            private = '_PRIVATE_' if jfname in private_setters else ''

                            if type_.primitive in [
                                    sysl_pb2.Type.DATE, sysl_pb2.Type.DATETIME]:
                                w('case "{}": entity.{}set{}(read{}(xsr, {})); break;',
                                  jfname, private, java.CamelCase(jfname), typename, extra)
                            else:
                                w('case "{}": entity.{}set{}(read{}(xsr)); break;',
                                  jfname, private, java.CamelCase(jfname), typename)
                        w('default: skipField(xsr);')
                w('table.insert(entity);')
                w('expect(XMLStreamConstants.END_ELEMENT, xsr.getEventType());')

        with java.Method(w, '\nprivate', 'void', 'expect',
                         [('int', 'expected'), ('int', 'got')]):
            with java.If(w, 'got != expected'):
                w('System.err.printf(\v'
                  '"<<Unexpected token: %s (expecting %s)>>\\n", \v'
                  'tokenName(got), tokenName(expected));')
                w('throw new {}Exception();', model_name)

        with java.Method(w, '\nprivate', 'String', 'tokenName', [('int', 'token')]):
            with java.Switch(w, 'token'):
                for tok in ('ATTRIBUTE CDATA CHARACTERS COMMENT DTD END_DOCUMENT '
                            'END_ELEMENT ENTITY_DECLARATION ENTITY_REFERENCE NAMESPACE '
                            'NOTATION_DECLARATION PROCESSING_INSTRUCTION SPACE START_DOCUMENT '
                            'START_ELEMENT'.split()):
                    w('case XMLStreamConstants.{0}: return "{0}";', tok)
            w('return new Integer(token).toString();')

        def read(t, access, extra_args=None):
            with java.Method(w, '\nprivate', t, 'read' + t,
                             [('XMLStreamReader', 'xsr')],
                             throws=['XMLStreamException']):
                with java.If(w, '!getCharacters(xsr)'):
                    w('return null;')
                w('{} result = {};', t, access)
                w('expect(XMLStreamConstants.END_ELEMENT, xsr.next());')
                w('return result;')

        read('String', 'xsr.getText()')
        read('Boolean', 'Boolean.parseBoolean(xsr.getText())')
        read('Integer', 'Integer.parseInt(xsr.getText())')
        read('Double', 'Double.parseDouble(xsr.getText())')
        read('BigDecimal', 'new BigDecimal(xsr.getText())')
        read('UUID', 'UUID.fromString(xsr.getText())')

        with java.Method(w, '\nprivate', 'DateTime', 'readDateTime',
                         [('XMLStreamReader', 'xsr'),
                          ('DateTimeFormatter', 'fmt')],
                         throws=['XMLStreamException']):
            with java.If(w, '!getCharacters(xsr)'):
                w('return null;')
            w('DateTime result = fmt.parseDateTime(xsr.getText());')
            w('expect(XMLStreamConstants.END_ELEMENT, xsr.next());')
            w('return result;')

        with java.Method(w, '\nprivate', 'LocalDate', 'readLocalDate',
                         [('XMLStreamReader', 'xsr'),
                          ('DateTimeFormatter', 'fmt')],
                         throws=['XMLStreamException']):
            with java.If(w, '!getCharacters(xsr)'):
                w('return null;')
            w('LocalDate result = fmt.parseLocalDate(xsr.getText());')
            w('expect(XMLStreamConstants.END_ELEMENT, xsr.next());')
            w('return result;')

        with java.Method(w, '\nprivate', 'void', 'skipField',
                         [('XMLStreamReader', 'xsr')], throws=['XMLStreamException']):
            with java.If(w, 'getCharacters(xsr)'):
                w('expect(XMLStreamConstants.END_ELEMENT, xsr.next());')

        with java.Method(w, '\nprivate', 'boolean', 'getCharacters',
                         [('XMLStreamReader', 'xsr')], throws=['XMLStreamException']):
            w('int tok = xsr.next();')
            with java.If(w, 'tok == XMLStreamConstants.END_ELEMENT'):
                w('return false;')
            w('expect(XMLStreamConstants.CHARACTERS, tok);')
            w('return true;')

        with java.Method(w, '\nprivate', 'void', 'skipElement',
                         [('XMLStreamReader', 'xsr')], throws=['XMLStreamException']):
            with java.While(w, 'xsr.next() != XMLStreamConstants.END_ELEMENT'):
                with java.If(w,
                             'xsr.getEventType() == XMLStreamConstants.START_ELEMENT'):
                    w('skipElement(xsr);')

        # TODO: Is permissive dateTimeParser OK for date types?
        w('\nprivate final DateTimeFormatter iso8601DateTime = '
          'ISODateTimeFormat.dateTimeParser();')

        w()
