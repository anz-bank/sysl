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


SWAGGER_FORMATS = {'int32', 'int64', 'float', 'double', 'date', 'date-time'}

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


def is_sysl_array_type(ftype):
    return ftype.startswith('sequence of ')


def type_as_key(swagt):
    if isinstance(swagt, dict):
        return frozenset(sorted(swagt.iteritems()))
    assert isinstance(swagt, basestring)
    return swagt


def extract_properties(tspec):
    # Some schemas for object types might lack a properties attribute.
    # As per Open API Spec 2.0, partially defined via JSON Schema draft 4,
    # if properties is missing we assume it is the empty object.
    # https://tools.ietf.org/html/draft-fge-json-schema-validation-00#section-5.4
    return tspec.get('properties', {})


METHOD_ORDER = {
    m: i for (i, m) in enumerate('get put post delete patch parameters'.split())}


def parse_args(argv):
    p = argparse.ArgumentParser(description='Converts Swagger (aka Open API Specification) documents to a Sysl spec')
    p.add_argument('swagger_path', help='path of input swagger document')
    p.add_argument('appname', help='appname')
    p.add_argument('package', help='package')
    p.add_argument('outfile', help='path of output file')
    return p.parse_args(args=argv[1:])


def make_default_logger():
    logger = logging.getLogger('import_swagger')
    logger.setLevel(logging.WARN)
    return logger


def load_vocabulary(words_fn='/usr/share/dict/words'):
    if not os.path.exists(words_fn):
        return []
    else:
        return (w.strip() for w in open(words_fn))


class SwaggerTranslator:
    def __init__(self, logger, vocabulary_factory=None):
        if vocabulary_factory is None:
            vocabulary_factory = load_vocabulary
        self._logger = logger
        self._param_cache = {}
        self._words = set()
        self._vocabulary_factory = vocabulary_factory

    def warn(self, msg):
        self._logger.warn(msg)

    def error(self, msg):
        self._logger.error(msg)

    def words(self):
        """Lazy-load WORDS."""
        if not self._words:
            self._words.update(self._vocabulary_factory())
            if not self._words:
                self.warn("could not load any vocabulary, janky environment-specific heuristics for renaming path template names may fail")
        return self._words

    def javaParam(self, match):
        # TODO(anz-rfc) this seems janky and fragile.
        param = match.group(1)
        ident = self._param_cache.get(param)

        if ident is None:
            # foo-bar to fooBar
            ident = re.sub(r'-(\w?)', lambda m: m.group(1).upper(), param)

            # {fooid} -> fooId (only if foo is in WORDS but fooid isn't)
            m = re.match(r'{([a-z]+)id}$', ident)
            if m:
                word = m.group(1)
                if word in self.words() and word + 'id' not in self.words():
                    ident = '{' + word + 'Id}'

            ident = ident.replace('}', '<:string}')

            self._param_cache[param] = ident

        return ident

    def translate_path_template_params(self, path):
        # OAS 2.0 supports path templating.
        # ref: https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#pathTemplating
        # For some reason, we rename these params via javaParam
        return re.sub(r'({[^/]*?})', self.javaParam, path)

    def translate(self, swag, appname, package, w):
        hasInfo = False

        hasInfo = 'info' in swag
        title = swag['info'].pop('title', '') if hasInfo else ''

        w(u'{}{} [package={}]:',
            appname, title and ' ' + json.dumps(title), json.dumps(package))

        with w.indent():
            if hasInfo:
                def info_attrs(info, prefix=''):
                    for (name, value) in sorted(info.iteritems()):
                        if isinstance(value, dict):
                            info_attrs(value, prefix + name + '.')
                        elif name != 'description':
                            w('@{}{} = {}', prefix, name, json.dumps(value))

                info_attrs(swag['info'])

            if 'host' in swag:
                w('@host = {}', json.dumps(swag['host']))

            w(u'@description =:')
            with w.indent():
                w(u'| {}', swag['info'].get('description', 'No description.').replace("\n", "\n|"))

            if 'basePath' in swag:
                w()
                w('{}:'.format(swag['basePath']))
                with w.indent():
                    self.writeEndpoints(swag, w)
            else:
                self.writeEndpoints(swag, w)

            self.writeDefs(swag, w)

    def writeEndpoints(self, swag, w):
        for (path, api) in sorted(swag['paths'].iteritems()):
            w(u'\n{}:', self.translate_path_template_params(path))
            with w.indent():
                if 'parameters' in api:
                    del api['parameters']
                for (i, (method, body)) in enumerate(sorted(api.iteritems(),
                                                            key=lambda t: METHOD_ORDER[t[0]])):
                    params = {where: [] for where in ['query', 'body', 'header']}

                    if 'parameters' in body:
                        for param in body['parameters']:
                            paramIn = param.get('in')
                            if paramIn and paramIn != 'path':
                                params[paramIn].append(param)
                            if '$ref' in param:
                                params['header'].append(swag['parameters'][param['$ref'].rpartition('/')[2]])

                    header = self.getHeaders(params['header'])
                    methodBody = self.getBody(params['body'])

                    if len(header) > 0 and len(methodBody) > 0:
                        paramStr = ' ({})'.format(methodBody + ', ' + header)
                    elif len(header) == 0 and len(methodBody) == 0:
                        paramStr = ''
                    else:
                        paramStr = ' ({})'.format(methodBody + header)

                    w(u'{}{}{}{}:',
                        method.upper(),
                        ' ?' if params['query'] else '',
                        '&'.join(
                            '{}={}{}'.format(p['name'], TYPE_MAP[p['type']], '' if p['required'] else '?')
                            for p in params['query']
                        ),
                        paramStr)
                    with w.indent():
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

                        returnValues = OrderedDict()
                        for key in sorted(responses):
                            schema = responses.get(key, {}).get('schema')
                            if schema is not None:
                                if schema.get('type') == 'array':
                                    returnValues['sequence of ' + schema['items']['$ref'].rpartition('/')[2]] = True
                                else:
                                    returnValues[schema['$ref'].rpartition('/')[2]] = True
                            else:
                                returnValues[', '.join(responses.keys())] = True
                            if key == 'default' or key.startswith('x-'):
                                self.warn('default responses and x-* responses are not implemented')

                        w(u'return {}'.format(', '.join(returnValues)))

                    if i < len(api) - 1:
                        w()

    def writeDefs(self, swag, w):
        w()
        w('#' + '-' * 75)
        w('# definitions')

        for (tname, tspec) in sorted(swag.get('definitions', {}).iteritems()):
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
                        w(getfields(fname, ftype))
                        with w.indent():
                            w(self.getTag(fname, 'json_tag'))
                            if fdescr is not None and len(fdescr) > 0:
                                w(self.getTag(fdescr, 'description'))
                # handle top-level arrays
                elif tspec.get('type') == 'array':
                    (ftype, fdescr) = self.parse_typespec(tspec, '', tname)
                    w(getfields('', ftype, True))
                    with w.indent():
                        if fdescr is not None and len(fdescr) > 0:
                            w(self.getTag(fdescr, 'description'))
                elif 'enum' in tspec:
                    w(tspec.get('type'))
                else:
                    externAlias['EXTERNAL_' + tname] = 'string'

        index = 0
        while len(typeSpecList) > index:
            w()
            typeName = '{}_obj'.format(
                typeSpecList[index].typeRef + ('_' + typeSpecList[index].parentRef if len(
                    typeSpecList[index].parentRef) > 0 else '')
            )

            if 'properties' in typeSpecList[index].element:
                w('!type {}:', typeName)
                with w.indent():
                    for (k, v) in extract_properties(typeSpecList[index].element).iteritems():
                        typeRef = typeSpecList[index].typeRef + '_' + typeSpecList[index].parentRef
                        ftypeSyntax = self.parse_typespec(v, k, typeRef)[0]
                        if 'required' not in typeSpecList[index].element or k not in typeSpecList[index].element['required']:
                            ftypeSyntax = ftypeSyntax + '?'
                        fields = '{} <: {}:'.format(
                            k if k not in SYSL_TYPES else k + '_',
                            ftypeSyntax)
                        w(fields)
                        with w.indent():
                            w(self.getTag(k, 'json_tag'))
            else:
                w('!alias EXTERNAL_{}:', typeName)
                with w.indent():
                    w('string')

            index += 1

        self.writeAlias(externAlias, w)

    def writeAlias(self, alias, w):
        for key in alias:
            w()
            w('!alias {}:', key)
            with w.indent():
                w(alias[key])

    def parse_typespec(self, tspec, parentRef='', typeRef=''):

        typ = tspec.get('type')

        # skip invalid arrays
        if not typ and 'items' in tspec:
            self.warn('Ignoring unexpected "items". Schema has "items" but did not have defined "type". Note: %r' % (tspec, ))
            del tspec['items']

        descr = tspec.pop('description', None)

        if typ == 'array':
            assert not (set(tspec.keys()) - {'type', 'items', 'example'}), tspec

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

        if fmt is not None and fmt not in SWAGGER_FORMATS:
            self.error('Invalid format: %s at %s -> %s' % (fmt, typeRef, parentRef))

        if ref:
            assert not set(tspec.keys()) - {'$ref'}, tspec
            return r(ref.lstrip('#/definitions/'))
        if (typ, fmt) == ('string', None):
            return r('string')
        if (typ, fmt) == ('boolean', None):
            return r('bool')
        elif (typ, fmt) == ('string', 'date'):
            return r('date')
        elif (typ, fmt) == ('string', 'date-time'):
            return r('datetime')
        elif (typ, fmt) == ('integer', None):
            return r('int')
        elif (typ, fmt) == ('integer', 'int32'):
            return r('int32')
        elif (typ, fmt) == ('integer', 'int64'):
            return r('int64')
        elif (typ, fmt) == ('number', 'double'):
            return r('float')
        elif (typ, fmt) == ('number', 'float'):
            return r('float')
        elif (typ, fmt) == ('number', None):
            return r('float')
        elif (typ, fmt) == ('object', None):
            typeSpecList.append(TypeSpec(tspec, parentRef, typeRef))
            retVal = typeRef + ('_' + parentRef if len(parentRef) > 0 else '') + '_obj'
            if 'properties' not in tspec:
                return r('EXTERNAL_' + retVal)
            return r(retVal)

        else:
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

            if 'type' in headerParam:
                typeName = TYPE_MAP[headerParam['type']]

            headerList.append('{} <: {} [~header, {}, name="{}"]'.format(paramName.replace('-', '_').lower(), typeName, necessity, paramName))

        return ', '.join(headerList)

    def getBody(self, bodyParams):
        if len(bodyParams) == 0:
            return ''

        paramList = []
        for param in bodyParams:
            paramList.append('{} <: {} [~body]'.format(param['name'], param['schema']['$ref'].rpartition('/')[2]))

        return ', '.join(paramList)


def main():
    args = parse_args(sys.argv)

    with open(args.swagger_path, 'r') as f:
        swag = yaml.load(f)

    w = writer.Writer('sysl')

    translator = SwaggerTranslator(logger=make_default_logger())
    translator.translate(swag, args.appname, args.package, w=w)

    with open(args.outfile, 'w') as f_out:
        f_out.write(str(w))


if __name__ == '__main__':
    debug.init()
    main()
