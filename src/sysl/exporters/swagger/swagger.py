# -*- encoding: utf-8 -*-
import yaml

from sysl.core import syslx

from sysl.util import datamodel
from sysl.util import java
from sysl.util import rex

from sysl.proto import sysl_pb2


TYPE_MAP = {
    sysl_pb2.Type.ANY: {'type': 'object'},
    sysl_pb2.Type.BOOL: {'type': 'boolean'},
    sysl_pb2.Type.INT: {'type': 'number', 'format': 'integer'},
    sysl_pb2.Type.FLOAT: {'type': 'number', 'format': 'double'},
    sysl_pb2.Type.DECIMAL: {'type': 'number', 'format': 'double'},
    sysl_pb2.Type.STRING: {'type': 'string'},
    sysl_pb2.Type.BYTES: None,
    sysl_pb2.Type.STRING_8: {'type': 'string'},
    sysl_pb2.Type.DATE: {'type': 'string'},
    sysl_pb2.Type.DATETIME: {'type': 'string'},
    sysl_pb2.Type.XML: {'type': 'string'},
}

STATUS_MAP = {
    # 1xx
    100: 'Continue',
    101: 'Switching Protocol',

    # 2xx
    200: 'OK',
    201: 'Created',
    202: 'Accepted',
    203: 'Non-Authoritative Information',
    204: 'No Content',
    205: 'Reset Content',
    206: 'Partial Content',
    207: 'Multi-Status',
    208: 'Already Reported',
    226: 'IM Used',

    # 3xx
    300: 'Multiple Choices',
    301: 'Moved Permanently',
    302: 'Found',
    303: 'See Other',
    304: 'Not Modified',
    305: 'Use Proxy',
    # 306: 'Switch Proxy', No longer used
    307: 'Temporary Redirect',
    308: 'Permanent Redirect',

    # 4xx
    400: 'Bad Request',
    401: 'Unauthorised',
    402: 'Payment Required',
    403: 'Forbidden',
    404: 'Not Found',
    405: 'Method Not Allowed',
    406: 'Not Acceptable',
    407: 'Proxy Authentication Required',
    408: 'Request Timeout',
    409: 'Conflict',
    410: 'Gone',
    411: 'Length Required',
    412: 'Precondition Failed',
    413: 'Payload Too Large',
    414: 'URI Too Long',
    415: 'Unsupported Media Type',
    416: 'Range Not Satisfiable',
    417: 'Expectation Failed',
    418: 'I\'m a teapot',
    421: 'Misdirected Request',
    422: 'Unprocessable Entity',
    423: 'Locked',
    424: 'Failed Dependency',
    426: 'Upgrade Required',
    428: 'Precondition Required',
    429: 'Too Many Requests',
    431: 'Request Header Fields Too Large',
    451: 'Unavailable For Legal Reasons',

    # 5xx
    500: 'Internal Server Error',
    501: 'Not Implemented',
    502: 'Bad Gateway',
    503: 'Service Unavailable',
    504: 'Gateway Timeout',
    505: 'HTTP Version Not Supported',
    506: 'Variant Also Negotiates',
    507: 'Insufficient Storage',
    508: 'Loop Detected',
    510: 'Not Extended',
    511: 'Network Authentication Required',
}


def no_utf8(value):
    if isinstance(value, unicode):
        return value.encode('utf-8')
    if isinstance(value, list):
        return [no_utf8(e) for e in value]
    if isinstance(value, dict):
        return {no_utf8(k): no_utf8(v) for (k, v) in value.iteritems()}
    return value


def swagger_file(app, module, root_class, write_file):
    def assign(root, path, value):
        if value is not None:
            for part in path[:-1]:
                root = root.setdefault(part, {})
            root[path[-1]] = value
        return value

    used_definitions = set()
    used_models = set()

    header = {}

    assign(header, ['swagger'], '2.0')
    paths = assign(header, ['paths'], {})
    definitions = assign(header, ['definitions'], {})

    for path in [
        ['version'],
        ['title'],
        ['description'],
        ['contact', 'name'],
            ['contact', 'email']]:
        assign(header, ['info'] + path,
               syslx.View(app.attrs)['.'.join(path)].s)

    assign(header, ['info', 'title'], app.long_name)

    typerefs = set()

    def swagger_type(t):
        if t.WhichOneof('type') == 'set':
            return {'type': 'array', 'items': swagger_type(t.set)}

        # TODO: Probably should use datamodel.typeref
        if t.WhichOneof('type') == 'type_ref':
            # TODO: What is the impact of removing this?
            #assert len(t.type_ref.ref.path) == 1
            model = t.type_ref.ref.path[0]
            if model in module.apps:
                used_models.add(model)
                return {'$ref': '#/definitions/' + model}

        type_ = datamodel.typeref(t, module)[2]
        if type_ is None:
            return None

        which_type = type_.WhichOneof('type')

        if which_type == 'primitive':
            return TYPE_MAP[type_.primitive]
        elif which_type == 'enum':
            return {'type': 'number', 'format': 'integer'}
        elif which_type in ['tuple', 'relation']:
            deftype = '.'.join(t.type_ref.ref.path)
            used_definitions.add(deftype)
            return {'$ref': '#/definitions/' + deftype}
        else:
            # return {'type': 'sysl:' + type_.WhichOneof('type')}
            raise RuntimeError(
                'Unexpected field type for Swagger '
                'export: ' + which_type)

    for (epname, endpt) in (
            sorted(app.endpoints.iteritems(), key=lambda t: t[0].lower())):
        if not endpt.HasField('rest_params'):
            continue

        rp = endpt.rest_params
        path = paths.setdefault(rp.path, {})
        method = path.setdefault(
            sysl_pb2.Endpoint.RestParams.Method.Name(rp.method).lower(), {})
        params = method['parameters'] = []
        for qp in rp.query_param:
            assert qp.loc, 'missing QueryParam.loc'
            params.append({
                'name': qp.name,
                'in': sysl_pb2.Endpoint.RestParams.QueryParam.Loc.Name(qp.loc).lower(),
                'required': not qp.type.opt,
            })
            params[-1].update(TYPE_MAP[qp.type.primitive])

        assign(method, ['description'], endpt.docstring)

        for stmt in endpt.stmt:
            if stmt.WhichOneof('stmt') == 'ret':
                m = rex.match(ur'''
          (?:
            (\d+)
            (?:·\((.*?)\))?
            (?:·:)?
            (?:·<:(·set·of)?·([\w\.]+))?
            (?:·or·{((?:\d+,)*\d+)})?
          ) | (?:
            one•of·{((?:\d+·,·)*\d+·)}
          )
          $
          ''', stmt.ret.payload)
                if m:
                    [status, descr, setof, type_, statuses, statuses2] = m.groups()
                    resp = assign(method, ['responses', status], {})
                    assign(resp, ['description'],
                           descr or STATUS_MAP.get(int(status or '0'), '(no description)'))
                    if type_:
                        schema = assign(resp, ['schema'], {})
                        if setof:
                            schema['type'] = 'array'
                            items = schema['items'] = {}
                        else:
                            items = schema
                        items['$ref'] = '#/definitions/' + type_
                        used_definitions.add(type_)
                    if statuses or statuses2:
                        for s in (statuses or statuses2).split(','):
                            assign(method, ['responses', s, 'description'],
                                   STATUS_MAP[int(s)])
                else:
                    print stmt.ret.payload
                    raise Exception('Bad return statement')

        body = syslx.View(endpt.attrs)['body'].s
        if body:
            params.append({
                'name': 'requestBody',
                'in': 'body',
                'required': True,
                'schema': {'$ref': '#/definitions/' + body}
            })
            if body == 'PoOrder':
                used_models.add(body)
            else:
                used_definitions.add(body)

    used_defs2 = used_definitions.copy()
    used_definitions.clear()

    # Figure out which types are used
    while True:
        for (tname, t) in sorted(app.types.iteritems()):
            if tname in used_defs2:
                if t.WhichOneof('type') in ['relation', 'tuple']:
                    entity = getattr(t, t.WhichOneof('type'))
                    for (fname, f) in entity.attr_defs.iteritems():
                        swagger_type(f)

        used_definitions -= used_defs2
        if not used_definitions:
            break

        used_defs2 |= used_definitions

    for (tname, t) in sorted(app.types.iteritems()):
        if tname not in used_defs2:
            continue

        if t.WhichOneof('type') in ['relation', 'tuple']:
            properties = {}
            definition = definitions[tname] = {
                'properties': properties
            }

            entity = getattr(t, t.WhichOneof('type'))

            for (fname, f) in entity.attr_defs.iteritems():
                jfname = java.name(fname)
                properties[jfname] = swagger_type(f)

    for modelname in used_models:
        defn = definitions[modelname] = {}
        props = defn['properties'] = {}
        for (tname, t) in module.apps[modelname].types.iteritems():
            props[tname] = {
                'type': 'array',
                'items': {'$ref': '#/definitions/' + tname}
            }

        # TODO: dedup with above logic
        for (tname, t) in module.apps[modelname].types.iteritems():
            if t.WhichOneof('type') in ['relation', 'tuple']:
                properties = {}
                definition = definitions[tname] = {
                    'properties': properties
                }

                entity = getattr(t, t.WhichOneof('type'))

                for (fname, f) in entity.attr_defs.iteritems():
                    jfname = java.name(fname)
                    (a, type_info, type_) = datamodel.typeref(f, module)
                    if type_info is None:
                        properties[jfname] = swagger_type(f)
                    else:
                        #assert type_.WhichOneof('type') == 'primitive'
                        properties[jfname] = swagger_type(type_)

    write_file(yaml.dump(no_utf8(header)), root_class + '.swagger.yaml')
