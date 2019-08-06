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
    "date": "date",
    "datetime": "dateTime",
    "string[]": "string[]",
    "date-only": "date-only",
    "array": "array",
    "number[]": "number[]"
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
    p = argparse.ArgumentParser(description='Converts RAML+JSON Schema documents to a Sysl spec')
    p.add_argument('raml_path', help='path of input raml document')
    p.add_argument('appname', help='appname')
    p.add_argument('package', help='package')
    p.add_argument('outfile', help='path of output file')
    return p.parse_args(args=argv[1:])


def make_default_logger():
    logger = logging.getLogger('import_raml')
    logger.setLevel(logging.WARN)
    return logger


def load_vocabulary(words_fn='/usr/share/dict/words'):
    if not os.path.exists(words_fn):
        return []
    else:
        return (w.strip() for w in open(words_fn))


class RamlTranslator:
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

    def translate_path_template_params(self, path, params):
        if params:
            parts = re.split(r'({[^/]*?})', path)
            if len(parts) > 1:
                pathParams = {}
                for param in params:
                    if '$ref' in param:
                        param = raml['parameters'][param['$ref'].rpartition('/')[2]]

                    if param['in'] == 'path':
                        pathParams['{' + param['name'] + '}'] = param
                if len(pathParams) > 0:
                    for (i, p) in enumerate(parts):
                        if parts[i] in pathParams:
                            param = pathParams[parts[i]]
                            parts[i] = "{" + param["name"] + "<:" + TYPE_MAP[param['type']] + "}"
                    return "".join(parts)
        return re.sub(r'({[^/]*?})', self.javaParam, path)

    def translate(self, raml, appname, package, raml_path, w):

        #TODO(kirkpatg) deal with headers
        #TODO(kirkpatg) deal with traits (unhappy paths)
        #TODO(kirkpatg) fix repated inclusion of "!alias EXTERNAL_ErrorItem_context_obj:"
        hasInfo = False
        title = raml.pop('title', '')

        w(u'{}{} [package={}]:',
            appname, title and ' ' + json.dumps(title), json.dumps(package))

        if 'baseUri' in raml:
            baseUri = raml['baseUri']
            if 'version' in raml:
                version = raml['version']
                if r'{version}' in baseUri:
                    baseUri = re.sub(r'\{version\}', version, baseUri)
            w()
            with w.indent():
                w('/{}:'.format(baseUri))
                with w.indent():
                    self.writeEndpoints(raml, raml_path, appname, w)
        else:
            self.writeEndpoints(raml, raml_path, appname, w)

        self.writeDefs(raml, raml_path, w)

    def _buildEnumConstraints(self, constraints):
        constraintsValue = []
        for param in constraints:
            if constraints[param].get('enum', None):
                constraintsValue.append([param])
                constraintsValue[-1].extend(constraints[param]['enum'])
        return constraintsValue

    def writeEndpoints(self, raml, raml_path, appname, w):

        def _addConstraints(ramlSnippet):
            if ramlSnippet.get('queryParameters', None):
                constraintsValue = self._buildEnumConstraints(ramlSnippet['queryParameters'])
                if len(constraintsValue) > 0:
                    w('@paramEnumConstraints = {}', json.dumps(constraintsValue))

        def _nextMethod(raml):
            for key in raml:
                if key in {'get', 'put', 'post', 'patch', 'delete'}:
                    yield key, raml[key]

        def _buildQueryParams(methodParamsRaml):
            paramLine = ' ?'
            for param in methodParamsRaml:
                if paramLine != ' ?':
                    paramLine += '&'
                if methodParamsRaml[param].get('type', None):
                    paramLine += param + '=' + TYPE_MAP[methodParamsRaml[param]['type']]
                else:
                    paramLine += param + '=string'
                if not methodParamsRaml[param].get('required', None):
                    paramLine += '?'
            return paramLine

        def _addAttributeIfTagExists(tag, ramlSnippet):
            if ramlSnippet.get(tag, None):
                w(u'@{} = "{}"', tag, ramlSnippet[tag])

        def _addDisplayName(ramlSnippet):
            _addAttributeIfTagExists('displayName', ramlSnippet)

        def _addDescription(ramlSnippet):
            _addAttributeIfTagExists('description', ramlSnippet)

        def _addVerbAndParams(method, methodRaml):
            if methodRaml.get('queryParameters', None):
                queryParams = _buildQueryParams(methodRaml['queryParameters'])
                w(u'{}{}:', method.upper(), queryParams)
            else:
                w(u'{}:', method.upper())

        def _addReturn(methodRaml):
            def _buildReturn(responsesRaml):
                for rcName, rcSpec in methodRaml['responses'].iteritems():
                        if rcSpec:
                            if rcSpec.get('body', None) and \
                            rcSpec.get('body').get('application/json', None) and \
                            rcSpec.get('body').get('application/json').get('type', None):
                                dataType = rcSpec['body']['application/json']['type']
                                if isinstance(dataType, IncludeTag):
                                    dataType = self.schemaFileNameToType(str(dataType))

                                w(u'return {} <: {}.{}', rcName, appname, dataType)
                            else:
                                w(u'return {}', rcName)

            if methodRaml.get('responses', None ):
                responsesRaml = methodRaml['responses']
                if len(responsesRaml) > 1:
                    w(u'one of:')
                    with w.indent():
                        _buildReturn(responsesRaml)
                else:
                    _buildReturn(responsesRaml)
                    
            else:
                w(u'return')

        def _addRequestHeaders(ramlSnippet):
            requestHeaders = ramlSnippet.get('headers', None)
            if requestHeaders is not None:
                headerArray = _getHeaders(requestHeaders)
                if len(headerArray) > 0:
                    w('@requestHeaders = {}', json.dumps(headerArray))

        def _addResponseHeaders(ramlSnippet):
            if ramlSnippet.get('responses', None ):
                for rcName, rcSpec in ramlSnippet['responses'].iteritems():
                    if rcSpec:
                        responseHeaders = rcSpec.get('headers', None)
                        if responseHeaders is not None:
                            headerArray = _getHeaders(responseHeaders, rcName)
                            if len(headerArray) > 0:
                                w('@responseHeaders = {}', json.dumps(headerArray))

        def _getHeaders(ramlSnippet, responseCode = None):
            headerArray = []
            for hName, hSpec in ramlSnippet.iteritems():
                headerArray.append([['name', hName]])
                if responseCode:
                    headerArray[-1].append(['responseCode', responseCode])
                for fName, fValue in hSpec.iteritems():
                    headerArray[-1].append([fName, str(fValue)])
            return headerArray

        def _decodeEndpoint(ep, epRaml):
            w(u'{}:', re.sub(r'\{(.*?)\}', r'{\1<:string}', ep))
            with w.indent():
                _addDisplayName(epRaml)
                for (method, methodRaml) in _nextMethod(epRaml):
                    _addVerbAndParams(method, methodRaml)

                    with w.indent():
                        _addDisplayName(methodRaml)
                        _addDescription(methodRaml)
                        _addConstraints(methodRaml)
                        _addRequestHeaders(methodRaml)
                        _addResponseHeaders(methodRaml)
                        _addReturn(methodRaml)

        def _isEndPoint(name):
            return name.startswith('/')

        def _visitAndDecodeAllEndpoints(ramlSnippet):
            for ep in ramlSnippet:
                if _isEndPoint(ep):
                    _decodeEndpoint(ep, ramlSnippet[ep])
                    with w.indent():
                        _visitAndDecodeAllEndpoints(ramlSnippet[ep])

        def _gatherInlineTraits(ramlSnippet, traitSource, traits, traitPrefix=''):
            if ramlSnippet.get('traits', None) is not None:
                for k, v in ramlSnippet['traits'].iteritems():
                    keyName = k
                    if traitPrefix:
                        keyName = traitPrefix+'.'+k
                    traits[keyName] = v
                    traits[keyName]['traitSource'] = traitSource

        def _gatherIncludedTraits(ramlSnippet, raml_path, traits):
            if ramlSnippet.get('uses', None) is not None:
                for traitName, traitSpec in ramlSnippet['uses'].iteritems():
                    fileName = os.path.join(os.path.dirname(raml_path), ramlSnippet['uses'][traitName])
                    with open(fileName, 'r') as f:
                        traitsRaml = yaml.safe_load(f)
                        _gatherInlineTraits(traitsRaml, os.path.dirname(fileName), traits, traitName)

        def _addTraitsToRaml(ramlSnippet, basePath, traits, inheritedTraits):
            for epName, epValue in ramlSnippet.iteritems():
                if _isEndPoint(epName):
                    epTraits = epValue.get('is', [])
                    # get end points
                    for (method, methodRaml) in _nextMethod(epValue):
                        methodTraits = methodRaml.get('is', [])
                        for trait in epTraits + inheritedTraits + methodTraits:
                            # some traits are dictionaries
                            traitName = ''
                            if isinstance(trait, dict):
                                assert len(trait) ==1
                                traitName = trait.keys()[0]
                            else:
                                traitName = trait
                            for typeName, typeValue in traits[traitName].iteritems():
                                if typeName in ['description', 'traitSource']:
                                    continue
                                for valueName, valueValue in typeValue.iteritems():
                                    if valueValue.get('body', None):
                                        if valueValue['body'].get('application/json', None):
                                            if valueValue['body']['application/json'].get('type', None):
                                                if isinstance(valueValue['body']['application/json']['type'], IncludeTag):
                                                    valueValue['body']['application/json']['type'] = \
                                                        traits[trait]['traitSource'].replace(os.path.dirname(basePath) + '/', '') + \
                                                            '/' + str(valueValue['body']['application/json']['type'])
                                    if methodRaml.get(typeName, None):
                                        methodRaml[typeName][valueName] = valueValue
                                    else:
                                        methodRaml[typeName] = {valueName : valueValue}
                    # add traits
                    with w.indent():
                        _addTraitsToRaml(epValue, basePath, traits, inheritedTraits)

        def _expandTraits(ramlSnippet, raml_path):
            traits = dict()
            _gatherInlineTraits(ramlSnippet, raml_path, traits)
            _gatherIncludedTraits(ramlSnippet, raml_path, traits)
            _addTraitsToRaml(ramlSnippet, raml_path, traits, [])

        _expandTraits(raml, raml_path)
        _visitAndDecodeAllEndpoints(raml)

    def schemaFileNameToType(self, fname):
        if fname[0].isupper():
            return fname
        fname = os.path.basename(fname)
        fname = fname.replace(".json", "")
        fname = fname.replace("-", " ")
        fname = fname.title()
        fname = fname.replace(" ", "")
        return fname

    def writeDefs(self, raml, raml_path, w):
        #TODO(kirkpatg) budgeting-presentation: ErrorItem was not dragged in
        w()
        w('#' + '-' * 75)
        w('# definitions')

        def _schemaFileNameToType(fname):
            fname = os.path.basename(fname)
            fname = fname.replace(".json", "")
            fname = fname.replace("-", " ")
            fname = fname.title()
            fname = fname.replace(" ", "")
            return fname

        def _yieldBangIncludedSchemas(raml):
            for k, v in raml.iteritems():
                if k == 'type' and isinstance(v, IncludeTag):
                    yield str(v)
                if isinstance(v, dict):
                    for includedSchema in _yieldBangIncludedSchemas(v):
                        yield includedSchema

        def _includeAllSchemas(typeName, pathToSchema, includedSchemas, referringPath):
            pathToSchema = os.path.abspath(pathToSchema)
            # work out what sort of path it is
            anchor = []
            fileName = pathToSchema
            if '#/' in pathToSchema:
                splitPath = pathToSchema.split('#/')
                anchor = splitPath[-1].split('/')
                fileName = splitPath[0]
            
            # deal with circular includes
            if typeName in includedSchemas:
                if includedSchemas[typeName] != pathToSchema:
                    self.warn('%s type has multiple definitions:\nold\t%s\nnew\t%s' % (typeName, includedSchemas[typeName], pathToSchema))
                return

            includedSchemas[typeName] = pathToSchema

            def _callInclude(fspec):
                # An internal reference
                if fspec.startswith('#'):
                    _includeAllSchemas(fspec[fspec.rfind('/')+1:], pathToSchema + fspec, \
                        includedSchemas, pathToSchema)
                # An externally referenced name
                elif '#' in fspec:
                    _includeAllSchemas(fspec[fspec.rfind('/')+1:], 
                        os.path.join(os.path.dirname(pathToSchema), fspec), 
                        includedSchemas, pathToSchema)
                # An externally referenced file
                else:
                    _includeAllSchemas(self.schemaFileNameToType(os.path.basename(fspec)), \
                        os.path.join(os.path.dirname(pathToSchema), fspec), 
                        includedSchemas, pathToSchema)

            def _includeArraySchema(tspec):
                assert 'items' in tspec, "This is not an array: {}".format(tspec)

                if '$ref' in tspec['items']:
                    _callInclude(tspec['items']['$ref'])
                    
            with open(fileName, 'r') as f:
                tspec = json.load(f)
                
                # Schemas embedded in properties
                if tspec.get('type', None) == 'object':
                    properties = dict()
                    if len(anchor) == 0:
                        properties = tspec.get('properties', None)
                    elif len(anchor) == 1:
                        properties = tspec[anchor[0]].get('properties', None)
                    elif len(anchor) == 2:
                        properties = tspec[anchor[0]][anchor[1]].get('properties', None)
                    else:
                        assert False
                    if properties:
                        for (fname, fspec) in properties.iteritems():
                            if '$ref' in fspec:
                                _callInclude(fspec['$ref'])
                            
                            # Schemas embedded in arrays in properties
                            if fspec.get('type', None) == 'array':
                                _includeArraySchema(fspec)
                            
                # Schemas embedded in top-level arrays
                if tspec.get('type', None) == 'array':
                    _includeArraySchema(tspec)
                
        def _processSchema(pathToSchema, tname, w):

            def getfields(fname, ftype, required, isArray=False):
                fieldTemplate = '{} <: {}:'
                ftypeSyntax = ftype

                if isArray:
                    fieldTemplate = '{}{}'
                elif not required:
                    ftypeSyntax = ftype + '?'

                return fieldTemplate.format(
                    fname if fname not in SYSL_TYPES else fname + '_',
                    ftypeSyntax)

            def _extractArrayItems(tname, arrayName, arraySchema, requiredFields):
                assert arraySchema['type'] == 'array'

                # Items are defined by reference
                if '$ref' in arraySchema['items']:
                    (ftype, fdescr) = self.parse_typespec(arraySchema, arrayName, tname)
                    w(getfields(arrayName, ftype, (arrayName in requiredFields)))
                elif 'type' in arraySchema['items']:
                    # Items are defined by an object with properties
                    if arraySchema['items']['type'] == 'object':
                        _extractProperties(arraySchema['items'])
                    # Items are defined by a single type
                    else:
                        _extractSingleProperty(arraySchema['items'], 'items', tname, ['items'])
                else:
                    assert False

            def _extractSingleProperty(fspec, fname, tname, requiredFields):
                niceNameFspec = fspec
                if '$ref' in fspec:
                    niceNameFspec = dict(fspec)
                    niceNameFspec['$ref'] = self.schemaFileNameToType(fspec['$ref'])
                
                (ftype, fdescr) = self.parse_typespec(niceNameFspec, fname, tname)
                
                w(getfields(fname, ftype, (fname in requiredFields)))
                with w.indent():
                    w(self.getTag(fname, 'json_tag'))
                    if fdescr is not None and len(fdescr) > 0:
                        w(self.getTag(fdescr, 'description'))

                    # Add everything else
                    for k, v in fspec.iteritems():
                        if k in {'$ref', 'description', 'type'}:
                            continue
                        else:
                            w('@{} = "{}"', k, v)
                            
            def _extractProperties(propSchema):
                # A list of properties/fields that are required to be present in this schema
                requiredFields = []
                if 'required' in propSchema:
                    if isinstance(propSchema['required'], list):
                        requiredFields = propSchema['required']
                w('!type {}:', tname)

                with w.indent():
                    if 'extends' in propSchema:
                        w(getfields('extends', \
                            self.schemaFileNameToType(propSchema['extends']['$ref']), True))
                        with w.indent():
                            if 'description' in propSchema:
                                w('@description = "{}"', propSchema['description'])

                    if 'properties' not in propSchema:
                        return

                    for (fname, fspec) in sorted(propSchema['properties'].iteritems()):
                        _extractSingleProperty(fspec, fname, tname, requiredFields)

            splitPath = pathToSchema.split('#/')
            with open(splitPath[0], 'r') as f:
                loadedSchema = json.load(f)
                
                # There are four ways a schema can be represented
                # 1. An object with properties (name comes from filename)
                # 2. One or more named types (name come from type name)
                # 3. A named collection of named types (name comes from type name)
                # 4. An array (name comes from 'sequence of' filename)
                # TODO(gbk) handle enums properly

                def _loadTypes(loadedSchema):
                    if loadedSchema['type'] == 'object':
                        _extractProperties(loadedSchema)
                    elif loadedSchema['type'] == 'array':
                        w('!type {}:', tname)
                        with w.indent():
                            _extractArrayItems(tname, 'items', loadedSchema, [])
                    elif loadedSchema['type'] == 'string' and 'enum' in loadedSchema:
                        w('!alias {}:', tname)
                        with w.indent():
                            _extractSingleProperty(loadedSchema, 'value', tname, ['value'])
                    else:
                        _extractSingleProperty(loadedSchema, 'value', tname, [])

                # Types 1 or 4:
                if 'type' in loadedSchema:
                    _loadTypes(loadedSchema)
                elif 'properties' in loadedSchema:
                    self.warn("Properties element without accompanying type: %r" % loadedSchema)
                    _extractProperties(loadedSchema)
                else:
                    # Type 2
                    if tname in loadedSchema:
                        if 'type' in loadedSchema[tname]:
                            _loadTypes(loadedSchema[tname])
                        elif 'properties' in loadedSchema:
                            self.warn("Properties element without accompanying type: %r" % loadedSchema[tname])
                            _extractProperties(loadedSchema[tname])
                        else:
                            _extractSingleProperty(loadedSchema[tname], 'value', tname, ['value'])
                    # Type 3
                    else:
                        collection = splitPath[-1].split('/')
                        assert tname == collection[1]
                        if 'type' in loadedSchema[collection[0]][tname]:
                            _loadTypes(loadedSchema[collection[0]][tname])
                        elif 'properties' in loadedSchema:
                            self.warn("Properties element without accompanying type: %r" % loadedSchema[collection[0]][tname])
                            _extractProperties(loadedSchema[collection[0]][tname])
                        else:
                            _extractSingleProperty(loadedSchema[collection[0]][tname], 'value', tname, ['value'])
                        

            self.writeAlias(externAlias, w)

        def _processTypeSpecList():
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


        includedSchemas = dict()
        # types section
        for (tname, includedSchema) in sorted(raml.get('types', {}).iteritems()):
            # sometimes bang-includes are schemas (e.g. with a type and example)
            schemaName = ''
            if isinstance(includedSchema, dict):
                schemaName = includedSchema['type']
            else:
                schemaName = includedSchema

            pathToSchema = os.path.join(os.path.dirname(raml_path), str(schemaName))
            _includeAllSchemas(tname, pathToSchema, includedSchemas, raml_path)

        # when there is no 'types' collection
        for schema in _yieldBangIncludedSchemas(raml):
            pathToSchema = os.path.join(os.path.dirname(raml_path), str(schema))
            tname = self.schemaFileNameToType(os.path.basename(schema))
            _includeAllSchemas(tname, pathToSchema, includedSchemas, raml_path)

        # build all the types
        for tname, schema in sorted(includedSchemas.iteritems()):
            w()
            _processSchema(schema, tname, w)

        _processTypeSpecList()

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
            assert False
            self.warn('Ignoring unexpected "items". Schema has "items" but did not have defined "type". Note: %r' % (tspec, ))
            del tspec['items']

        descr = tspec.pop('description', None)
        if typ == 'array':
            # assert not (set(tspec.keys()) - {'$schema', 'type', 'items', 'example', 'uniqueItems', 'minItems', 'maxItems', 'required'}), tspec

            (itype, idescr) = self.parse_typespec(tspec['items'], parentRef, typeRef)

            if '$ref' in tspec['items']:
                itype = self.schemaFileNameToType(itype)
            return (sysl_array_type_of(itype), descr)

        def r(t):
            return (t, descr)

        fmt = tspec.get('format')
        ref = tspec.get('$ref')

        if fmt is not None and fmt not in SWAGGER_FORMATS:
            self.error('Invalid format: %s at %s -> %s' % (fmt, typeRef, parentRef))

        if ref:
            assert not set(tspec.keys()) - {'$ref', 'type'}, tspec
            return r(ref)
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



class IncludeTag(yaml.YAMLObject):
        yaml_tag = u'!include'
        def __init__(self, included_path):
            self.included_path = included_path
        def __repr__(self):
            return self.included_path

        @classmethod
        def from_yaml(cls, loader, node):
            return IncludeTag(node.value)


def main():
    args = parse_args(sys.argv)

    yaml.SafeLoader.add_constructor('!include', IncludeTag.from_yaml)

    with open(args.raml_path, 'r') as f:
        raml = yaml.safe_load(f)

    w = writer.Writer('sysl')

    translator = RamlTranslator(logger=make_default_logger())
    translator.translate(raml, args.appname, args.package, args.raml_path, w=w)

    with open(args.outfile, 'w') as f_out:
        f_out.write(str(w))


if __name__ == '__main__':
    debug.init()
    main()
