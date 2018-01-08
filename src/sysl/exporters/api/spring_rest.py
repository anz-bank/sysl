# -*- encoding: utf-8 -*-

import collections
import json

from sysl.util import datamodel
from sysl.util import java
from sysl.util import rex
from sysl.util import scopes
from sysl.util import writer

from sysl.proto import sysl_pb2

from sysl.core import syslx


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

# TODO: merge with //src/exporters/swagger:swagger.STATUS_MAP
STATUS_MAP = {
    400: 'Bad Request',
    401: 'Unauthorised',
    403: 'Forbidden',
    404: 'Not Found',
    412: '???',
    500: 'Internal Server Error',
}


def codeForParams(query_params, scope):
    return [
        (java.codeForType(qp.type, scope), qp.name)
        for qp in query_params]


def controller(interfaces, context):
    (app, module, package, model_class, write_file, _, _) = context

    w = writer.Writer('java')

    java.Package(w, package)

    java.Import(w, 'io.swagger.annotations.Api')
    java.Import(w, 'io.swagger.annotations.ApiOperation')
    java.Import(w, 'io.swagger.annotations.ApiParam')
    java.Import(w, 'io.swagger.annotations.ApiResponse')
    java.Import(w, 'io.swagger.annotations.ApiResponses')
    java.Import(w, 'lombok.extern.slf4j.Slf4j')
    java.Import(w, 'org.springframework.beans.factory.annotation.Autowired')
    java.Import(w, 'org.springframework.web.bind.annotation.PathVariable')
    java.Import(w, 'org.springframework.web.bind.annotation.RequestHeader')
    java.Import(w, 'org.springframework.web.bind.annotation.RequestMapping')
    java.Import(w, 'org.springframework.web.bind.annotation.RequestMethod')
    java.Import(w, 'org.springframework.web.bind.annotation.RequestParam')
    java.Import(w, 'org.springframework.web.bind.annotation.ResponseBody')
    java.Import(w, 'org.springframework.web.bind.annotation.RestController')
    w()

    w('@Api(value = {}, description = {}, position = {})',
      json.dumps(app.long_name),
      json.dumps(syslx.View(app.attrs)['description'].s),
      1)
    w('@Slf4j')
    w('@RestController')
    w('@ResponseBody')
    w('@RequestMapping(value = {}, produces = {})',
      json.dumps('/'),
      json.dumps('application/json;version=1.0;charset=UTF-8;'))

    with java.Class(w, model_class + 'Controller', write_file,
                    visibility='public', package=package):
        endpts = app.endpoints.itervalues()
        interfaces = {endpt.attrs['interface'].s for endpt in endpts}
        for (i, interface) in enumerate(sorted(interfaces)):
            if not interface:
                endpts = app.endpoints.itervalues()
                endpts_no_interface = [endpt.name for endpt in endpts if not
                                       endpt.attrs['interface'].s]
                print 'No interfaces for\n' + ('\n').join(endpts_no_interface)
                continue
            w('\n@Autowired'[not i:])
            w('private {} {};', interface, java.mixedCase(interface))

        scope = scopes.Scope(module)

        for (epname, endpt) in (
            sorted(
                app.endpoints.iteritems(),
                key=lambda t: (t[1].rest_params.path, t[0]))):
            if not endpt.HasField('rest_params'):
                continue

            rp = endpt.rest_params
            rest_method = rp.Method.Name(rp.method)
            method_name = (
                rest_method +
                rex.sub(
                    r'{(\w+)}',
                    lambda m: m.group(1).upper(),
                    rex.sub(r'[-/]', '_', rp.path)))

            def responses(stmts, result=None, cond=''):
                if result is None:
                    result = collections.defaultdict(list)

                for stmt in stmts:
                    which_stmt = stmt.WhichOneof('stmt')
                    if which_stmt == 'cond':
                        responses(
                            stmt.cond.stmt, result, (cond and cond + ' & ') + stmt.cond.test)
                    elif which_stmt == 'ret':
                        m = rex.match(ur'''
              (?:
                (?:
                  (\d+)·
                  (\([^\)]+\))?·
                  (\w+)?·
                  (?:
                    <:·
                    (empty\s+)?
                    (set\s+of\s+)?
                    ([\w.]+|\.\.\.)
                  )?
                )
                |
                (?:
                  one\s+of·{((?:\d+·,·)*\d+·)}
                )
              )?
              $
              ''', stmt.ret.payload)
                        if m:
                            [status, descr, expr, empty, setof,
                                type_, statuses] = m.groups()
                            for status in rex.split(
                                    ur'·,·', status or statuses):
                                status = int(status)
                                result[int(status)].append(
                                    cond or descr or STATUS_MAP.get(int(status)) or '???')
                        else:
                            print repr(stmt.ret.payload)
                            raise Exception('Bad return statement')
                return result

            w()
            w('@RequestMapping(method = RequestMethod.{}, \vpath = {})',
              rest_method, json.dumps(rp.path))
            w('@ApiOperation(value = {})', json.dumps(endpt.docstring))
            w('@ApiResponses({{')
            with w.indent():
                for (status, conds) in sorted(
                        responses(endpt.stmt).iteritems()):
                    w('@ApiResponse(code = {}, message =', status)
                    with w.indent():
                        for (i, cond) in enumerate(conds):
                            w('"<p style=\\"white-space:nowrap\\">{}</p>"{}',
                              cond, ' +' if i < len(conds) - 1 else '')
                    w('),')
            w('}})')

            params = codeForParams(rp.query_param, scope)
            with java.Method(w, 'public', 'Object', method_name, params):
                w('return {}.{}({});',
                  java.mixedCase(endpt.attrs['interface'].s),
                  method_name,
                  ', '.join('\v' + p for (_, p) in params))


def interface(interfaces, context):
    (app, module, package, model_class, write_file, _, _) = context

    for interface in interfaces:
        w = writer.Writer('java')

        java.Package(w, package)

        with java.Class(w, interface, write_file, type_='interface',
                        visibility='public', package=package):
            scope = scopes.Scope(module)

            for (epname, endpt) in (
                sorted(
                    app.endpoints.iteritems(),
                    key=lambda t: (t[1].rest_params.path, t[0]))):
                if endpt.attrs['interface'].s != interface:
                    continue
                if not endpt.HasField('rest_params'):
                    continue

                rp = endpt.rest_params
                rest_method = rp.Method.Name(rp.method)
                method_name = (
                    rest_method +
                    rex.sub(
                        r'{(\w+)}',
                        lambda m: m.group(1).upper(),
                        rex.sub(r'[-/]', '_', rp.path)))

                w('public Object {}({});',
                  method_name,
                  ', '.join('\v{} {}'.format(t, p)
                            for (t, p) in codeForParams(rp.query_param, scope)))


def service(interfaces, context):
    interface(interfaces, context)
    controller(interfaces, context)
