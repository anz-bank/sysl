#!/usr/bin/env python
# -*- encoding: utf-8 -*-

import argparse
import json
import logging
import re
import sys
import textwrap
import os.path
import yaml

from collections import OrderedDict
from sysl.util import debug
from sysl.util import writer

from sysl.proto import sysl_pb2


OPENAPI_FORMATS = {'int32', 'int64', 'float', 'double', 'date', 'date-time', 'byte', 'binary'}

SYSL_TYPES = {
    'int32', 'int64', 'int', 'float', 'string',
    'date', 'bool', 'decimal', 'datetime', 'xml'
}

TYPE_MAP = {
    "string": "string",
    "boolean": "bool",
    "number": "float",
    "integer": "int",
    "date": "date"
}


class TypeSpec:
    def __init__(self, element, parentRef, typeRef):
        self.element = element
        self.parentRef = parentRef
        self.typeRef = typeRef


typeSpecList = []
externAlias = {}


def sysl_array_type_of(itemtype):
    return 'sequence of ' + itemtype


def extract_properties(tspec):
    # Some schemas for object types might lack a properties attribute.
    # As per Open API Spec 2.0, partially defined via JSON Schema draft 4,
    # if properties is missing we assume it is the empty object.
    # https://tools.ietf.org/html/draft-fge-json-schema-validation-00#section-5.4
    return tspec.get('properties', {})


METHOD_ORDER = {
    m: i for (i, m) in enumerate('get put post delete patch parameters'.split())}


def parse_args(argv):
    p = argparse.ArgumentParser(description='Converts Open API Specification documents to a Sysl spec')
    p.add_argument('openapi_path', help='path of input openapi document')
    p.add_argument('appname', help='appname')
    p.add_argument('package', help='package')
    p.add_argument('outfile', help='path of output file')
    return p.parse_args(args=argv[1:])


def make_default_logger():
    logger = logging.getLogger('import_openapi')
    logger.setLevel(logging.WARN)
    return logger


class OpenApiTranslator:
    def __init__(self, logger):
        self._logger = logger
        self._param_cache = {}

    def warn(self, msg):
        self._logger.warn(msg)

    def error(self, msg):
        self._logger.error(msg)

    def javaParam(self, param):
        # TODO(anz-rfc) this seems janky and fragile.
        param = param[1:-1]
        ident = self._param_cache.get(param)

        if ident is None:
            # foo-bar to fooBar
            ident = re.sub(r'-(\w?)', lambda m: m.group(1).upper(), param)
            self._param_cache[param] = ident

        return ident

    def translate_path_template_params(self, path, params):
        parts = re.split(r'({[^/]*?})', path)
        if len(parts) > 1 and not params:
            self.warn("not enough path params path: %s" % (path))
        pathParams = {}
        for param in params:
            assert param['in'] == 'path'
            pathParams['{' + param['name'] + '}'] = param
        for (i, p) in enumerate(parts):
            if parts[i] in pathParams:
                param = pathParams[parts[i]]
                parts[i] = "{" + self.javaParam(parts[i]) + "<:" + TYPE_MAP[param['schema']['type']] + "}"
            elif parts[i] != "" and parts[i][0] == "{":
                self.warn("could not find type for path param: %s in params%s" % (parts[i], params))
                parts[i] = "{" + self.javaParam(parts[i]) + "<:string}"
        return "".join(parts)

    def translate(self, oaSpec, appname, package, w):
        hasInfo = 'info' in oaSpec
        title = oaSpec['info'].pop('title', '') if hasInfo else ''

        with w.indent(u'{}{} [package={}]:', appname, title and ' ' + json.dumps(title), json.dumps(package)):
            if hasInfo:
                def info_attrs(info, prefix=''):
                    for (name, value) in sorted(info.iteritems()):
                        if isinstance(value, dict):
                            info_attrs(value, prefix + name + '.')
                        elif name != 'description':
                            w('@{}{} = {}', prefix, name, json.dumps(value))

                info_attrs(oaSpec['info'])

            if 'host' in oaSpec:
                w('@host = {}', json.dumps(oaSpec['host']))

            with w.indent(u'@description =:'):
                w(u'| {}', oaSpec['info'].get('description', 'No description.').replace("\n", "\n|"))

            if 'basePath' in oaSpec:
                w()
                with w.indent('{}:', oaSpec['basePath']):
                    self.writeEndpoints(oaSpec, w)
            else:
                self.writeEndpoints(oaSpec, w)
            self.writeDefs(oaSpec, w)

    def mergeParams(self, oaSpec, srcParams, overrideParams):
        params = {
            'path': OrderedDict(),
            'header': OrderedDict(),
            'query': OrderedDict(),
        }

        for param in srcParams + overrideParams:
            if '$ref' in param:
                param = oaSpec['parameters'][param['$ref'].rpartition('/')[2]]
            params[param['in']][param["name"]] = param

        return [params[i].values() for i in ('path', 'header', 'query')]

    def writeEndpoints(self, oaSpec, w):

        def includeResponses(responses):
            returnValues = OrderedDict()

            for rc, rspec in sorted(responses.iteritems()):

                if rc == 'default' or rc.startswith('x-'):
                    self.warn('default and x-* responses are not implemented')
                    continue

                schema = None
                if rspec.get('content', None):
                    schema = rspec['content']
                    if schema.get('application/json', None):
                        schema = schema['application/json']
                        if schema.get('schema', None):
                            schema = schema['schema']

                if schema is not None:
                    varName = ""
                    if rc.startswith('2'):
                        varName = "ok <: "
                    elif rc.startswith('4') or rc.startswith('5'):
                        varName = "error <: "

                    if varName != "":
                        if schema.get('type') == 'array':
                            returnValues[varName + 'sequence of ' + schema['items']['$ref'].rpartition('/')[2]] = True
                        else:
                            returnValues[varName + schema['$ref'].rpartition('/')[2]] = True

            if len(returnValues) > 2:
                self.error('invalid return value set:' + json.dumps(returnValues))

            if len(returnValues) == 0:
                w(u'return')
            else:
                for rv in returnValues:
                    w(u'return {}', rv)

        for (path, api) in sorted(oaSpec['paths'].iteritems()):
            apiParams = api.pop('parameters') if 'parameters' in api else []

            for (i, (method, body)) in enumerate(sorted(api.iteritems(), key=lambda t: METHOD_ORDER[t[0]])):
                bodyContent = OrderedDict()
                bodyParams = body['parameters'] if 'parameters' in body else []
                pathParams, headerParams, queryParams = self.mergeParams(oaSpec, apiParams, bodyParams)

                with w.indent(u'\n{}:', self.translate_path_template_params(path, pathParams)):
                    header = self.getHeaders(headerParams)
                    appJson = body.get('requestBody', {}).get('content', {}).get('application/json', None)
                    jsonBody = [appJson] if appJson else []
                    methodBody = self.getBody(jsonBody)
                    params = [p for p in [methodBody, header] if p]
                    paramStr = ' ({})'.format(', '.join(params)) if params else ''

                    with w.indent(u'{}{}{}{}:',
                                  method.upper(), paramStr, ' ?' if queryParams else '',
                                  '&'.join('{}={}{}'.format(p['name'], TYPE_MAP[p['schema']['type']], '' if p.get('required', None) else '?')
                                           for p in queryParams)):
                        for line in textwrap.wrap(
                                body.get('description', 'No description.').strip(), 64):
                            w(u'| {}', line)

                        # Backwards compat: support integer response keys.
                        responses = {str(k): v for (k, v) in body['responses'].iteritems()}

                        # Valid keys of responses can be:
                        # - http status codes (I believe they must be string values, NOT integers, but spec is unclear).
                        # - "default"
                        # - any string matching the pattern "^x-(.*)$""
                        # ref: https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#responses-object

                        includeResponses(responses)

    def writeDefs(self, oaSpec, w):
        w()
        w('#' + '-' * 75)
        w('# definitions')

        schemas = oaSpec.get('components', {}).get('schemas', None)
        if not schemas:
            return

        def _collectExternalsSoReferencesAreCorrect():
            for (tname, tspec) in sorted(schemas.iteritems()):
                properties = extract_properties(tspec)

                if properties or tspec.get('type') == 'array' or 'enum' in tspec:
                    continue
                externAlias['EXTERNAL_' + tname] = 'string'

        _collectExternalsSoReferencesAreCorrect()

        for (tname, tspec) in sorted(schemas.iteritems()):
            properties = extract_properties(tspec)
            w()

            if properties:
                w('!type {}:', tname)
            elif tspec.get('type') == 'array' or 'enum' in tspec:
                w('!alias {}:', tname)

            requiredFields = []
            if 'required' in tspec:
                requiredFields = tspec['required']

            with w.indent():
                def getfields(fname, ftype, isArray=False):
                    fieldTemplate = '{} <: {}:'
                    ftypeSyntax = ftype

                    if isArray:
                        fieldTemplate = '{}{}'
                    elif fname not in requiredFields:
                        ftypeSyntax = ftype + '?'

                    return fieldTemplate.format(
                        fname if fname not in SYSL_TYPES else fname + '_',
                        ftypeSyntax)

                if properties:
                    for (fname, fspec) in sorted(properties.iteritems()):
                        (ftype, fdescr) = self.parse_typespec(fspec, fname, tname)
                        with w.indent(getfields(fname, ftype)):
                            w(self.getTag(fname, 'json_tag'))
                            if fdescr and len(fdescr) > 0:
                                w(self.getTag(fdescr, 'description'))
                # handle top-level arrays
                elif tspec.get('type') == 'array':
                    (ftype, fdescr) = self.parse_typespec(tspec, '', tname)
                    with w.indent(getfields('', ftype, True)):
                        if fdescr and len(fdescr) > 0:
                            w(self.getTag(fdescr, 'description'))
                elif 'enum' in tspec:
                    w(TYPE_MAP[tspec.get('type')])

        for typeSpec in typeSpecList:
            w()
            typeName = '{}_obj'.format(
                typeSpec.typeRef + ('_' + typeSpec.parentRef if len(
                    typeSpec.parentRef) > 0 else '')
            )

            if 'properties' in typeSpec.element:
                with w.indent('!type {}:', typeName):
                    for (k, v) in extract_properties(typeSpec.element).iteritems():
                        typeRef = typeSpec.typeRef + '_' + typeSpec.parentRef
                        ftypeSyntax = self.parse_typespec(v, k, typeRef)[0]
                        if 'required' not in typeSpec.element or k not in typeSpec.element['required']:
                            ftypeSyntax = ftypeSyntax + '?'
                        fields = '{} <: {}:'.format(k if k not in SYSL_TYPES else k + '_', ftypeSyntax)
                        with w.indent(fields):
                            w(self.getTag(k, 'json_tag'))
            else:
                with w.indent('!alias EXTERNAL_{}:', typeName):
                    w('string')

        self.writeAlias(externAlias, w)

    def writeAlias(self, alias, w):
        for key in alias:
            w()
            with w.indent('!alias {}:', key):
                w(alias[key])

    def parse_typespec(self, tspec, parentRef='', typeRef=''):

        typ = tspec.get('type')

        # skip invalid arrays
        if not typ and 'items' in tspec:
            self.warn('Ignoring unexpected "items". Schema has "items" but did not have defined "type". Note: %r' % (tspec, ))
            del tspec['items']

        descr = tspec.pop('description', None)

        if typ == 'array':
            assert not set(tspec.keys()) - {'type', 'items', 'example'}, tspec

            # skip invalid type
            if '$ref' in tspec['items'] and 'type' in tspec['items']:
                self.warn('Ignoring unexpected "type". Schema has "$ref" but also has unexpected "type". Note: %r' % (tspec, ))
                del tspec['items']['type']

            (itype, idescr) = self.parse_typespec(tspec['items'], parentRef, typeRef)
            return (sysl_array_type_of(itype), descr)

        def r(t):
            return (t, descr)

        fmt = tspec.get('format')

        ref = tspec.get('$ref')

        if fmt is not None and fmt not in OPENAPI_FORMATS:
            self.error('Invalid format: %s at %s -> %s' % (fmt, typeRef, parentRef))

        typMap = {
            ('string', None): r('string'),
            ('boolean', None): r('bool'),
            ('string', 'date'): r('date'),
            ('string', 'date-time'): r('datetime'),
            ('string', 'byte'): r('string'),
            ('integer', None): r('int'),
            ('integer', 'int32'): r('int32'),
            ('integer', 'int64'): r('int64'),
            ('number', 'double'): r('float'),
            ('number', 'float'): r('float'),
            ('number', None): r('float'),
        }

        if ref:
            assert not set(tspec.keys()) - {'$ref'}, tspec
            t = ref.replace('#/components/schemas/', '')
            if externAlias.get('EXTERNAL_' + t, None):
                return r('EXTERNAL_' + t)
            if externAlias.get('EXTERNAL_' + t + '_obj', None):
                return r('EXTERNAL_' + t + '_obj')
            return r(t)

        if (typ, fmt) in typMap:
            return typMap[(typ, fmt)]

        if (typ, fmt) == ('object', None):
            typeSpecList.append(TypeSpec(tspec, parentRef, typeRef))
            retVal = typeRef + ('_' + parentRef if len(parentRef) > 0 else '') + '_obj'
            if 'properties' not in tspec:
                return r('EXTERNAL_' + retVal)
            return r(retVal)

        aliasName = 'EXTERNAL_' + typeRef + ('_' + parentRef if len(parentRef) > 0 else '') + '_obj'
        externAlias[aliasName] = 'string'
        return r(aliasName)

    def getTag(self, tagName, tagType):
        if tagType is None or tagName is None:
            return ''
        return '@{} = "{}"'.format(tagType, tagName)

    def getHeaders(self, headerParams):
        headerList = []

        for headerParam in headerParams:
            paramName = headerParam['name']
            typeName = 'string'
            necessity = '~optional'

            if headerParam.get('required'):
                necessity = '~required'

            if 'schema' in headerParam:
                typeName = TYPE_MAP[headerParam['schema']['type']]

            headerList.append('{} <: {} [~header, {}, name="{}"]'.format(paramName.replace('-', '_').lower(), typeName, necessity, paramName))

        return ', '.join(headerList)

    def getBody(self, bodyParams):
        if len(bodyParams) == 0:
            return ''

        paramList = []
        for param in bodyParams:
            name = param.get('name', None)
            if not name:
                name = param['schema']['$ref'].rpartition('/')[2] + 'Request'
            paramList.append('{} <: {} [~body]'.format(name, param['schema']['$ref'].rpartition('/')[2]))

        return ', '.join(paramList)


def main():
    args = parse_args(sys.argv)

    with open(args.openapi_path, 'r') as f:
        oaSpec = yaml.load(f)

    w = writer.Writer('sysl')

    translator = OpenApiTranslator(logger=make_default_logger())
    translator.translate(oaSpec, args.appname, args.package, w)

    with open(args.outfile, 'w') as f_out:
        f_out.write(str(w))


if __name__ == '__main__':
    debug.init()
    main()
