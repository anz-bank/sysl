from sysl.proto import sysl_pb2

from sysl.core import syslx

from sysl.util import datamodel
from sysl.util import java
from sysl.util import writer


JSON_GEN_MAP = {
    sysl_pb2.Type.ANY: ('Object', '{}'),
    sysl_pb2.Type.BOOL: ('Boolean', '{}'),
    sysl_pb2.Type.INT: ('Number', '{}'),
    sysl_pb2.Type.FLOAT: ('Number', '{}'),
    sysl_pb2.Type.DECIMAL: ('Number', '{}'),
    sysl_pb2.Type.STRING: ('String', '{}'),
    sysl_pb2.Type.BYTES: ('Binary', '{}'),
    sysl_pb2.Type.STRING_8: ('String', '{}'),
    sysl_pb2.Type.DATE: ('LocalDate', '{}'),
    sysl_pb2.Type.DATETIME: ('DateTime', '{}'),
    sysl_pb2.Type.XML: ('String', '{}'),
    sysl_pb2.Type.UUID: ('String', '{}'),
}

JSON_PARSE_MAP = {
    sysl_pb2.Type.ANY: None,
    sysl_pb2.Type.BOOL: 'BooleanValue',
    sysl_pb2.Type.INT: 'IntValue',
    sysl_pb2.Type.FLOAT: 'DoubleValue',
    sysl_pb2.Type.DECIMAL: 'DecimalValue',
    sysl_pb2.Type.STRING: 'Text',
    sysl_pb2.Type.BYTES: 'BinaryValue',
    sysl_pb2.Type.STRING_8: 'Text',
    sysl_pb2.Type.DATE: ('Text', 'iso8601DateTime.parseLocalDate({})'),
    sysl_pb2.Type.DATETIME: ('Text', 'iso8601DateTime.parseDateTime({})'),
    sysl_pb2.Type.XML: 'Text',
    sysl_pb2.Type.UUID: ('Text', 'UUID.fromString({})'),
}


def serializer(context):
    (app, module, package, model_class, write_file, _, _) = context

    facade = bool(context.wrapped_model)

    w = writer.Writer('java')

    java.Package(w, package)

    java.StandardImports(w)

    java.Import(w, 'java.io.IOException')
    w.head()
    java.Import(w, 'java.text.ParseException')
    w.head()
    java.Import(w, 'com.fasterxml.jackson.core.JsonGenerator')
    java.Import(w, 'com.fasterxml.jackson.databind.JsonSerializer')
    java.Import(w, 'com.fasterxml.jackson.databind.SerializerProvider')
    w.head()
    java.Import(w, 'org.joda.time.format.DateTimeFormatter')
    java.Import(w, 'org.joda.time.format.ISODateTimeFormat')

    if facade:
        model_name = syslx.fmt_app_name(context.wrapped_model.name)
        modelpkg = syslx.View(context.wrapped_model.attrs)['package'].s
        w.head()
        java.Import(w, modelpkg + '.*')

    w()
    with java.Class(w, model_class + 'JsonSerializer', write_file,
                    package=package, extends='JsonSerializer<' + model_class + '>'):
        w()
        with java.Method(w, 'public', 'void', 'serialize',
                         [(model_class, 'facade' if facade else 'model'),
                          ('JsonGenerator', 'g'),
                          ('SerializerProvider', 'provider')],
                         throws=['IOException'],
                         override=True):
            if facade:
                w('{} model = facade.getModel();', model_name)

            w(u'g.writeStartObject();')
            for (tname, t) in sorted(app.types.iteritems()):
                if t.WhichOneof('type') == 'relation':
                    w(u'serialize{0}View(g, model.get{0}Table());', tname)
            w(u'g.writeEndObject();')

        for (tname, t) in sorted(app.types.iteritems()):
            if not t.WhichOneof('type') in ['relation', 'tuple']:
                continue
            java.SeparatorComment(w)
            if t.WhichOneof('type') == 'relation':
                w()
                with java.Method(w, 'public', 'void', 'serialize' + tname + 'View',
                                 [('JsonGenerator', 'g'),
                                  (tname + '.View', 'view')],
                                 throws=['IOException']):
                    with java.If(w, 'view.isEmpty()'):
                        w('return;')
                    w('g.writeArrayFieldStart("{}");', tname)
                    with java.For(w, '{} item : view', tname):
                        w('g.writeStartObject();')
                        for (fname, f) in datamodel.sorted_fields(t):
                            jfname = java.name(fname)
                            method = java.CamelCase(jfname)
                            type_ = datamodel.typeref(f, module)[2]
                            extra = '{}'
                            which_type = type_.WhichOneof('type')
                            if which_type == 'primitive':
                                (jsontype,
                                 extra) = JSON_GEN_MAP[type_.primitive]
                                if type_.primitive == type_.DECIMAL:
                                    for c in type_.constraint:
                                        if c.scale:
                                            access = (
                                                '{0} == null ? null : item.{0}.setScale({1}, '
                                                'java.math.RoundingMode.HALF_UP)').format(
                                                jfname, c.scale)
                                            break
                                else:
                                    access = jfname
                            elif which_type == 'enum':
                                jsontype = 'Number'
                                access = jfname + '.getValue()'
                            elif which_type == 'tuple':
                                access = jfname
                            else:
                                raise RuntimeError(
                                    'Unexpected field type for JSON export: ' + which_type)
                            w(u'writeField(g, "{}", {});',
                              jfname, extra.format('item.' + access))
                        w(u'g.writeEndObject();')
                    w(u'g.writeEndArray();')
            else:
                # fieldname, entity
                with java.Method(w, 'public', 'void', 'serialize',
                                 [('JsonGenerator', 'g'),
                                  ('String', 'fieldname'),
                                  (tname, 'entity')],
                                 throws=['IOException']):
                    with java.If(w, 'entity == null'):
                        w(u'return;')
                    w(u'g.writeFieldName(fieldname);')
                    w(u'serialize(g, entity);')

                # entity
                with java.Method(w, 'public', 'void', 'serialize',
                                 [('JsonGenerator', 'g'),
                                  (tname, 'entity')],
                                 throws=['IOException']):
                    w(u'g.writeStartObject();')
                    for (fname, f) in datamodel.sorted_fields(t):
                        jfname = java.name(fname)
                        method = java.CamelCase(jfname)
                        type_ = datamodel.typeref(f, module)[2]
                        which_type = ''
                        if type_ is None:
                            if f.WhichOneof('type') == 'set' and f.set.HasField('primitive'):
                                which_type = 'tuple'
                        else:
                            which_type = type_.WhichOneof('type')
                        extra = '{}'
                        if which_type == 'primitive':
                            access = 'entity.get{}()'.format(method)
                            (jsontype, extra) = JSON_GEN_MAP[type_.primitive]
                            if type_.primitive == type_.DECIMAL:
                                for c in type_.constraint:
                                    if c.scale:
                                        access = (
                                            'entity.{} == null ? null : entity.get{}().setScale({}, '
                                            'java.math.RoundingMode.HALF_UP)').format(
                                            jfname, method, c.scale)
                                        break
                            w(u'writeField(g, "{}", {});', jfname, access)
                        elif which_type == 'enum':
                            access = 'entity.{} == null ? null : entity.get{}().getValue()'.format(jfname, method)
                            w(u'writeField(g, "{}", {});', jfname, access)
                        elif which_type == 'tuple':
                            w(u'serialize(g, "{}", entity.get{}());'.format(
                                fname, method))
                        else:
                            raise RuntimeError(
                                'Unexpected field type for JSON export: ' + which_type)
                    # end for
                    w(u'g.writeEndObject();')
                w()

                # view, fieldname
                with java.Method(w, 'public', 'void', 'serialize',
                                 [('JsonGenerator', 'g'), ('String', 'fieldname'),
                                  (tname + '.View', 'view')],
                                 throws=['IOException']):
                    with java.If(w, 'view == null || view.isEmpty()'):
                        w(u'return;')
                    w(u'g.writeArrayFieldStart(fieldname);')
                    with java.For(w, '{} item : view', tname):
                        w(u'serialize(g, item);')
                    w(u'g.writeEndArray();')

        for t in ['Boolean', 'String']:
            w()
            with java.Method(w, 'private', 'void', 'writeField',
                             [('JsonGenerator', 'g'),
                              ('String', 'fieldname'), (t, 'value')],
                             throws=['IOException']):
                with java.If(w, 'value != null'):
                    w('g.write{}Field(fieldname, value);', t)

        # TODO(sahejsingh): implement for ['Boolean', 'Integer', 'Double', 'BigDecimal', 'DateTime', 'LocalDate', 'UUID']
        w()
        with java.Method(w, 'private', 'void', 'serialize',
                         [('JsonGenerator', 'g'),
                          ('String', 'fieldname'), ('HashSet<String>', 'value')],
                         throws=['IOException']):
            with java.If(w, 'value == null || value.isEmpty()'):
                w(u'return;')
            w(u'g.writeArrayFieldStart(fieldname);')
            with java.For(w, 'String item : value', t):
                w(u'g.writeString(item);')
            w(u'g.writeEndArray();')

        w()
        with java.Method(w, 'private', 'void', 'writeField',
                         [('JsonGenerator', 'g'), ('String',
                                                   'fieldname'), ('Integer', 'value')],
                         throws=['IOException']):
            with java.If(w, 'value != null'):
                w('g.writeNumberField(fieldname, value.intValue());')

        w()
        with java.Method(w, 'private', 'void', 'writeField',
                         [('JsonGenerator', 'g'), ('String',
                                                   'fieldname'), ('Double', 'value')],
                         throws=['IOException']):
            with java.If(w, 'value != null'):
                w('g.writeNumberField(fieldname, value.doubleValue());')

        w()
        with java.Method(w, 'private', 'void', 'writeField',
                         [('JsonGenerator', 'g'),
                          ('String', 'fieldname'),
                          ('BigDecimal', 'value')],
                         throws=['IOException']):
            with java.If(w, 'value != null'):
                w('g.writeNumberField(fieldname, value);')

        w()
        with java.Method(w, 'private', 'void', 'writeField',
                         [('JsonGenerator', 'g'),
                          ('String', 'fieldname'),
                          ('DateTime', 'value')],
                         throws=['IOException']):
            with java.If(w, 'value != null'):
                w('g.writeStringField(fieldname, iso8601DateTime.print(value));')

        w()
        with java.Method(w, 'private', 'void', 'writeField',
                         [('JsonGenerator', 'g'),
                          ('String', 'fieldname'),
                          ('LocalDate', 'value')],
                         throws=['IOException']):
            with java.If(w, 'value != null'):
                w('g.writeStringField(fieldname, iso8601Date.print(value));')

        w()
        with java.Method(w, 'private', 'void', 'writeField',
                         [('JsonGenerator', 'g'), ('String',
                                                   'fieldname'), ('UUID', 'value')],
                         throws=['IOException']):
            with java.If(w, 'value != null'):
                w('g.writeStringField(fieldname, value.toString());')

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

    java.Import(w, 'java.io.IOException')
    w.head()
    java.Import(w, 'java.text.ParseException')
    w.head()
    java.Import(w, 'com.fasterxml.jackson.core.JsonParseException')
    java.Import(w, 'com.fasterxml.jackson.core.JsonParser')
    java.Import(w, 'com.fasterxml.jackson.core.JsonToken')
    java.Import(w, 'com.fasterxml.jackson.databind.JsonDeserializer')
    java.Import(w, 'com.fasterxml.jackson.databind.DeserializationContext')
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

    has_tables = any(
        t.HasField('relation')
        for (tname, t) in sorted(app.types.iteritems()))

    w()
    with java.Class(w, model_class + 'JsonDeserializer', write_file,
                    package=package,
                    extends='JsonDeserializer<' + model_class + '>'):
        w()
        with java.Method(w, 'public', model_class, 'deserialize',
                            [('JsonParser', 'p'),
                             ('DeserializationContext', 'provider')],
                            throws=['IOException', 'JsonParseException'], override=True):
            w('{0} model = new {0}();', model_name)
            if facade:
                w('{0} facade = new {0}(model);', model_class)
            w()
            w('eatToken(p, JsonToken.START_OBJECT);')
            with java.While(w, 'p.getCurrentToken() != JsonToken.END_OBJECT'):
                with java.Switch(w, 'eatName(p)'):
                    for (tname, t) in sorted(app.types.iteritems()):
                        if t.HasField('relation'):
                            w(('case "{0}": deserialize{0}Table(p, '
                                'model.get{0}Table()); break;'),
                                tname)
                    w('default: skipJson(p);')
            w('expectToken(p, JsonToken.END_OBJECT);')
            w()
            if facade:
                w('return facade;')
            else:
                w('return model;')

        for (tname, t) in sorted(app.types.iteritems()):
            if t.HasField('tuple'):
                # HashSet<User defined type>
                with java.Method(w, 'private', tname + '.View', 'deserialize' + tname + 'View',
                                    [('JsonParser', 'p')],
                                    throws=['IOException', 'JsonParseException']):
                    w()
                    w('{0}.Set view = new {0}.Set();', tname)
                    w('eatToken(p, JsonToken.START_ARRAY);')
                    with java.While(w, 'p.getCurrentToken() != JsonToken.END_ARRAY'):
                        w('{0} entity = {0}._PRIVATE_new();', tname)
                        w('deserialize(p, entity);')
                        w('view.add(entity);')
                    w('p.nextToken();')
                    w('return view;')

                with java.Method(w, 'public', tname, 'deserialize',
                                    [('JsonParser', 'p'), (tname, 'entity')],
                                    throws=['IOException', 'JsonParseException'], override=False):
                    w()
                    w('eatToken(p, JsonToken.START_OBJECT);')
                    with java.If(w, 'entity == null'):
                        w('entity = {0}._PRIVATE_new();', tname)
                    with java.While(w, 'p.getCurrentToken() != JsonToken.END_OBJECT'):
                        with java.Switch(w, 'eatName(p)'):
                            for (fname, f) in datamodel.sorted_fields(t):
                                jfname = java.name(fname)
                                (typename, _, type_) = datamodel.typeref(
                                    f, module)
                                extra = '{}'
                                set_with_primitive = False

                                if type_ is None:
                                    if f.WhichOneof('type') == 'set' and f.set.HasField('primitive'):
                                        which_type = 'tuple'
                                        set_with_primitive = True
                                        type_ = f.set
                                else:
                                    which_type = type_.WhichOneof('type')

                                if which_type == 'primitive':
                                    jsontype = JSON_PARSE_MAP[type_.primitive]
                                    if isinstance(jsontype, tuple):
                                        (jsontype, extra) = jsontype
                                elif which_type == 'enum':
                                    jsontype = 'IntValue'
                                    extra = typename + '.from({})'
                                elif which_type == 'tuple':
                                    if set_with_primitive:
                                        extra = 'deserializeArray(p)'
                                    elif f.WhichOneof('type') == 'set':
                                        extra = 'deserialize{}View(p)'.format(
                                            f.set.type_ref.ref.path[-1])
                                    elif f.WhichOneof('type') == 'type_ref':
                                        extra = 'deserialize(p, entity.get{}())'.format(
                                            java.CamelCase(jfname))
                                    jsontype = ''
                                else:
                                    raise RuntimeError(
                                        'Unexpected field type for JSON export: ' + which_type)
                                private = ''
                                if type_.primitive in [
                                        sysl_pb2.Type.DATE, sysl_pb2.Type.DATETIME]:
                                    with java.Case(w, '"{}"', jfname):
                                        w(('entity.{}set{}('
                                            'p.getCurrentToken() == JsonToken.VALUE_NULL'
                                            ' ? null : {}); p.nextToken(); break;'),
                                            private,
                                            java.CamelCase(jfname),
                                            extra.format('p.get{}()'.format(jsontype)))
                                else:
                                    w(('case "{}": entity.{}set{}(p.getCurrentToken() == '
                                        'JsonToken.VALUE_NULL ? null : {}); {} break;'),
                                        jfname,
                                        private,
                                        java.CamelCase(jfname),
                                        extra.format(
                                            'p.get{}()'.format(jsontype)),
                                        '' if which_type == 'tuple' else 'p.nextToken();')

                            w('default: skipJson(p);')
                    w('p.nextToken();')
                    w('return entity;')

            if not t.HasField('relation'):
                continue

            pkey = datamodel.primary_key_params(t, context.module)
            pkey_fields = {f for (_, f, _) in pkey}
            fkeys = {
                java.name(fname): type_info
                for (fname, _, type_info) in datamodel.foreign_keys(t, context.module)}

            private_setters = pkey_fields | set(fkeys.iterkeys())

            w()
            with java.Method(w, 'private', 'void', 'deserialize' + tname + 'Table',
                             [('JsonParser', 'p'), (tname + '.Table', 'table')],
                             throws=['IOException', 'JsonParseException']):
                w('eatToken(p, JsonToken.START_ARRAY);')
                with java.While(w, 'p.getCurrentToken() != JsonToken.END_ARRAY'):
                    w('eatToken(p, JsonToken.START_OBJECT);')
                    w('{0} entity = {0}._PRIVATE_new(table.model());', tname)
                    with java.While(w, 'p.getCurrentToken() != JsonToken.END_OBJECT'):
                        with java.Switch(w, u'eatName(p)'):
                            for (fname, f) in datamodel.sorted_fields(t):
                                jfname = java.name(fname)
                                (typename, _, type_) = datamodel.typeref(f, module)
                                extra = '{}'
                                which_type = type_.WhichOneof('type')
                                if which_type == 'primitive':
                                    jsontype = JSON_PARSE_MAP[type_.primitive]
                                    if isinstance(jsontype, tuple):
                                        (jsontype, extra) = jsontype
                                elif which_type == 'enum':
                                    jsontype = 'IntValue'
                                    extra = typename + '.from({})'
                                else:
                                    raise RuntimeError(
                                        'Unexpected field type for JSON export: ' + which_type)
                                private = '_PRIVATE_' if jfname in private_setters else ''
                                if type_.primitive in [
                                        sysl_pb2.Type.DATE, sysl_pb2.Type.DATETIME]:
                                    with java.Case(w, '"{}"', jfname):
                                        w(('entity.{}set{}('
                                           'p.getCurrentToken() == JsonToken.VALUE_NULL'
                                           ' ? null : {}); break;'),
                                          private,
                                          java.CamelCase(jfname),
                                          extra.format('p.get{}()'.format(jsontype)))
                                else:
                                    w(('case "{}": entity.{}set{}(p.getCurrentToken() == '
                                       'JsonToken.VALUE_NULL ? null : {}); break;'),
                                      jfname,
                                      private,
                                      java.CamelCase(jfname),
                                      extra.format('p.get{}()'.format(jsontype)))
                            with java.Default(w):
                                w('skipJson(p);')
                                w('continue;')
                        w('p.nextToken();')
                    w('p.nextToken();')
                    w()
                    w('table.insert(entity);')
                w('p.nextToken();')

        # HashSet<Primitive Type>
        with java.Method(w, 'private', 'HashSet<String>', 'deserializeArray',
                            [('JsonParser', 'p')],
                            throws=['IOException', 'JsonParseException']):
            w()
            w('HashSet<String> view = new HashSet<String>();')
            w('eatToken(p, JsonToken.START_ARRAY);')
            with java.While(w, 'p.getCurrentToken() != JsonToken.END_ARRAY'):
                w('expectToken(p, JsonToken.VALUE_STRING);')
                w('view.add(p.getText());')
                w('p.nextToken();')
            w('p.nextToken();')
            w('return view;')

        with java.Method(w, '\nprivate', 'void', 'eatToken',
                         [('JsonParser', 'p'), ('JsonToken', 'token')],
                         throws=['IOException']):
            w(u'expectToken(p, token);')
            w(u'p.nextToken();')

        with java.Method(w, '\nprivate', 'void', 'expectToken',
                         [('JsonParser', 'p'), ('JsonToken', 'token')]):
            with java.If(w, 'p.getCurrentToken() != token'):
                w(('System.err.printf("<<Unexpected token: %s (expecting %s)>>\\n", '
                   'tokenName(p.getCurrentToken()), tokenName(token));'))
                w('throw new {}Exception();', model_name)

        with java.Method(w, '\nprivate', 'String', 'eatName', [('JsonParser', 'p')],
                         throws=['IOException']):
            w('expectToken(p, JsonToken.FIELD_NAME);')
            w('String name = p.getCurrentName();')
            w('p.nextToken();')
            w('return name;')

        with java.Method(w, '\nprivate', 'String', 'tokenName',
                         [('JsonToken', 'token')]):
            with java.If(w, 'token == null'):
                w('return "null";')
            with java.Switch(w, 'token'):
                for tok in (
                    'END_ARRAY END_OBJECT FIELD_NAME NOT_AVAILABLE START_ARRAY '
                    'START_OBJECT VALUE_EMBEDDED_OBJECT VALUE_FALSE VALUE_NULL '
                    'VALUE_NUMBER_FLOAT VALUE_NUMBER_INT VALUE_STRING VALUE_TRUE'
                ).split():
                    w('case {0}: return "{0}";', tok)
            w('return "";')

        # TODO: refactor into base class
        # TODO: replace recursion with depth counter
        with java.Method(w, '\nprivate', 'void', 'skipJson', [('JsonParser', 'p')],
                         throws=['IOException']):
            w('JsonToken tok = p.getCurrentToken();')
            w('p.nextToken();')
            with java.Switch(w, 'tok'):
                for tok in (
                    'END_ARRAY END_OBJECT FIELD_NAME NOT_AVAILABLE '
                    'VALUE_EMBEDDED_OBJECT VALUE_FALSE VALUE_NULL '
                        'VALUE_NUMBER_FLOAT VALUE_NUMBER_INT VALUE_STRING VALUE_TRUE').split():
                    w('case {}: break;', tok)
                with java.Case(w, 'START_ARRAY'):
                    with java.While(w, 'p.getCurrentToken() != JsonToken.END_ARRAY'):
                        w('skipJson(p);')
                    w('p.nextToken();')
                    w('break;')
                with java.Case(w, 'START_OBJECT'):
                    with java.While(w, 'p.getCurrentToken() != JsonToken.END_OBJECT'):
                        w('p.nextToken();')
                        w('skipJson(p);')
                    w('p.nextToken();')
                    w('break;')

        # TODO: Is permissive dateTimeParser OK for date types?
        w('\nprivate final DateTimeFormatter iso8601DateTime = '
          'ISODateTimeFormat.dateTimeParser();')

        w()
