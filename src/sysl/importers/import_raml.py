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
import filecmp

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
    def __init__(self, element, parentRef, typeRef, path):
        self.element   = element
        self.parentRef = parentRef
        self.typeRef   = typeRef
        self.path      = path



typeSpecList = []
externAlias  = {}


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
    def __init__(self, logger, raml_path, vocabulary_factory=None):
        if vocabulary_factory is None:
            vocabulary_factory = load_vocabulary
        self._logger = logger
        self._param_cache = {}
        self._words = set()
        self._vocabulary_factory = vocabulary_factory
        self._raml_path = raml_path
        self._base_dir = os.path.abspath(os.path.dirname(raml_path))

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

    def translate(self, raml, appname, package, w):

        # Inject all the traits before we do anything
        self.incorporateTraits(raml)

        # Build up all the types so we can use them throughout the processing
        typeLookup, includedSchemas = self.processSchemas(raml)
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
                    self.writeEndpoints(raml, appname, typeLookup, w)
        else:
            self.writeEndpoints(raml, appname, typeLookup, w)

        self.writeDefs(raml, includedSchemas, w)

    def _buildEnumConstraints(self, constraints):
        constraintsValue = []
        for param in constraints:
            if constraints[param].get('enum', None):
                constraintsValue.append([param])
                constraintsValue[-1].extend(constraints[param]['enum'])
        return constraintsValue

    def incorporateTraits(self, raml):
        def _gatherInlineTraits(ramlSnippet, traitSource, traits, traitPrefix=''):
            if ramlSnippet.get('traits', None) is not None:
                for k, v in ramlSnippet['traits'].iteritems():
                    keyName = k
                    if traitPrefix:
                        keyName = traitPrefix+'.'+k
                    traits[keyName] = v
                    traits[keyName]['traitSource'] = traitSource

        def _gatherIncludedTraits(ramlSnippet, traits):
            if ramlSnippet.get('uses', None) is not None:
                for traitName, traitSpec in ramlSnippet['uses'].iteritems():
                    fileName = os.path.join(self._base_dir, ramlSnippet['uses'][traitName])
                    with open(fileName, 'r') as f:
                        traitsRaml = yaml.safe_load(f)
                        _gatherInlineTraits(traitsRaml, fileName, traits, traitName)

        def _addTraitsToRaml(ramlSnippet, traits, inheritedTraits):
            for epName, epValue in ramlSnippet.iteritems():
                if self._isEndPoint(epName):
                    epTraits = epValue.get('is', [])
                    # get end points
                    for (method, methodRaml) in self.nextMethod(epValue):
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
                                                    valueValue['body']['application/json']['schema'] = \
                                                        os.path.abspath(os.path.join(os.path.dirname(traits[trait]['traitSource']), \
                                                            str(valueValue['body']['application/json']['type']))).replace(self._base_dir + '/', '')
                                                    valueValue['body']['application/json']['type'] = \
                                                        self.schemaFileNameToType(str(valueValue['body']['application/json']['type']), traits[trait]['traitSource'])
                                    if methodRaml.get(typeName, None):
                                        methodRaml[typeName][valueName] = valueValue
                                    else:
                                        methodRaml[typeName] = {valueName : valueValue}

                    _addTraitsToRaml(epValue, traits, inheritedTraits)

        def _expandTraits(ramlSnippet):
            traits = dict()
            _gatherInlineTraits(ramlSnippet, self._raml_path, traits)
            _gatherIncludedTraits(ramlSnippet, traits)
            _addTraitsToRaml(ramlSnippet, traits, [])

        _expandTraits(raml)

    def _isEndPoint(self, name):
        return name.startswith('/')

    def nextMethod(self, ramlSnippet):
        for key in ramlSnippet:
            if key in {'get', 'put', 'post', 'patch', 'delete'}:
                yield key, ramlSnippet[key]

    def writeEndpoints(self, raml, appname, typeLookup, w):

        def _addConstraints(ramlSnippet):
            if ramlSnippet.get('queryParameters', None):
                constraintsValue = self._buildEnumConstraints(ramlSnippet['queryParameters'])
                if len(constraintsValue) > 0:
                    w('@paramEnumConstraints = {}', json.dumps(constraintsValue))

        

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
                body = _getRequestBody(methodRaml)
                if body:
                    w(u'{} {}:', method.upper(), body)
                else:
                    w(u'{}:', method.upper())

        def _addReturn(methodRaml):
            def _buildReturn(responsesRaml):
                for rcName, rcSpec in sorted(methodRaml['responses'].iteritems()):
                    if rcSpec:
                        if rcSpec.get('body', None) and \
                            rcSpec.get('body').get('application/json', None):
                            typeTag = ''
                            if rcSpec.get('body').get('application/json').get('type', None):
                                typeTag = 'type'
                            elif rcSpec.get('body').get('application/json').get('schema', None):
                                typeTag = 'schema'
                            if typeTag:
                                dataType = rcSpec['body']['application/json'][typeTag]
                                
                                if isinstance(dataType, IncludeTag):
                                    dataType = self.schemaFileNameToType(str(dataType), self._raml_path)
                                
                                if dataType not in ['string']:
                                    dataType = typeLookup[dataType]
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

        def _getRequestBody(methodRaml):
            if methodRaml.get('body', None):
                if not methodRaml['body'].get('application/json', None):
                    self.error("Expecting application/json but got: {}".format(methodRaml['body'].keys()))
                    return None
                if not methodRaml['body']['application/json'].get('type', None):
                    self.error("Expecting type but got: {}".format(methodRaml['body']['application/json'].keys()))
                    return None
                typeName = self.schemaFileNameToType(str(methodRaml['body']['application/json']['type']), self._raml_path)
                typeName = typeLookup[typeName]
                return '(request <: {}.{})'.format(appname, typeName)
            else:
                return None
        
        def _decodeEndpoint(ep, epRaml):
            w(u'{}:', re.sub(r'\{(.*?)\}', r'{\1<:string}', ep))
            with w.indent():
                _addDisplayName(epRaml)
                for (method, methodRaml) in self.nextMethod(epRaml):
                    _addVerbAndParams(method, methodRaml)

                    with w.indent():
                        _addDisplayName(methodRaml)
                        _addDescription(methodRaml)
                        _addConstraints(methodRaml)
                        _addRequestHeaders(methodRaml)
                        _addResponseHeaders(methodRaml)
                        _addReturn(methodRaml)

        def _visitAndDecodeAllEndpoints(ramlSnippet):
            for ep in ramlSnippet:
                if self._isEndPoint(ep):
                    _decodeEndpoint(ep, ramlSnippet[ep])
                    with w.indent():
                        _visitAndDecodeAllEndpoints(ramlSnippet[ep])

        _visitAndDecodeAllEndpoints(raml)

    def schemaFileNameToType(self, fname, path):
        if not fname:
            return None
        
        fullPath = os.path.abspath(
                        os.path.join(os.path.dirname(path.split('#')[0]), \
                            os.path.dirname(fname.split('#')[0])))
        prefix = fullPath.replace(self._base_dir,'').replace('/','.')
        if len(prefix) > 1:
            if prefix[0] == '.':
                prefix = prefix[1:]

        anchor = fname.split('#')
        if len(anchor) == 2 and anchor[0] and '.json' not in anchor[0]:
            anchorBase = anchor[1].split('/')
            if len(anchorBase) > 1:
                prefix += '.' + '.'.join(anchorBase[1:-1]) + '.'
            else:
                if prefix:
                    prefix += '.'
        else:
            if prefix:
                prefix += '.'

        if fname.endswith('.json'):
            fname = os.path.basename(fname)
            fname = fname.replace(".json", "")
            fname = fname.replace("-", " ")
            fname = fname.title()
            fname = fname.replace(" ", "")
            return prefix + fname

        if '/' in fname:
            if '#' in fname:
                return prefix + fname[(fname.rfind('/')+1):]
            else:
                assert False
        else:
            return fname

    def processSchemas(self, raml):

        includedSchemas = dict()
        typeLookup = dict()

        def _includeAllSchemas(typeName, pathToSchema):
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
                assert includedSchemas[typeName] == pathToSchema
                return

            assert pathToSchema.count('#') <= 1, "Malformed schema path: {}".format(pathToSchema)
            
            includedSchemas[typeName] = pathToSchema

            def _nextInclude(fspec):
                tn = self.schemaFileNameToType(fspec, pathToSchema.split('#')[0])
                pts = os.path.join(os.path.dirname(pathToSchema.split('#')[0]), fspec)
                
                if fspec.startswith('#'):
                    pts = pathToSchema.split('#')[0] + fspec

                _includeAllSchemas(tn, pts)

            with open(fileName, 'r') as f:
                tspec = json.load(f)
                
                baseSpec = tspec
                if len(anchor) == 1:
                    baseSpec = tspec[anchor[0]]
                elif len(anchor) == 2:
                    baseSpec = tspec[anchor[0]][anchor[1]]
                elif len(anchor) != 0:
                    assert False

                def _descendTheSpec(spec):
                    for k, v in spec.iteritems():
                        if k == '$ref':
                            _nextInclude(v)
                        elif isinstance(v, dict):
                            _descendTheSpec(v)
                _descendTheSpec(baseSpec)
        
        # recurse through included raml files and add types to the includedSchemas dict
        def _getIncludesFromRamlFiles(ramlSnippet, path, pk):
            
            def _allRaml(rs, pk):
                for (k,v) in rs.iteritems():
                    if isinstance(v, dict):
                        pk.append(k)
                        _allRaml(v, pk)
                        pk.pop()
                    elif isinstance(v, IncludeTag) or isinstance(v, basestring):
                        if k not in ['example','description'] and '.json' in str(v):
                            pathToSchema = os.path.join(os.path.dirname(path), str(v))
                            tname = self.schemaFileNameToType(str(v), path)
                            if 'types' in pk or 'schemas' in pk:
                                if pk[-1] in ['types','schemas']:
                                    prefix = '.'.join(pk[:-1])
                                    if prefix:
                                        typeLookup[prefix + '.' + k] = tname
                                    typeLookup[k] = tname
                                elif pk[-1] != 'examples':
                                    typeLookup[pk[-1]] = tname
                                        
                                    #typeLookup[tname[tname.rfind('.')+1:]] = tname
                            
                            typeLookup[tname] = tname
                            _includeAllSchemas(tname, pathToSchema)
                
            _allRaml(ramlSnippet, pk)
            
            if ramlSnippet.get('uses', None):
                for k, v in ramlSnippet['uses'].iteritems():
                    nextPath = os.path.join(os.path.dirname(path), v)
                    with open(nextPath, 'r') as f:
                        nextRaml = yaml.safe_load(f)
                        pk.append(k)
                        _getIncludesFromRamlFiles(nextRaml, os.path.abspath(nextPath), pk)
                        pk.pop()

        _getIncludesFromRamlFiles(raml, self._raml_path, [])

        return typeLookup, includedSchemas

    def writeDefs(self, raml, includedSchemas, w):
        #TODO(kirkpatg) budgeting-presentation: ErrorItem was not dragged in
        w()
        w('#' + '-' * 75)
        w('# definitions')
                
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
                    (ftype, fdescr) = self.parse_typespec(arraySchema, pathToSchema, arrayName, tname)
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
                    niceNameFspec['$ref'] = self.schemaFileNameToType(fspec['$ref'], pathToSchema.split('#')[0])
                
                (ftype, fdescr) = self.parse_typespec(niceNameFspec, pathToSchema, fname, tname)
                w(getfields(fname, ftype, (fname in requiredFields)))
                with w.indent():
                    w(self.getTag(fname, 'json_tag'))
                    if fdescr is not None and len(fdescr) > 0:
                        w(self.getTag(fdescr, 'description'))

                    # Add everything else
                    for k, v in fspec.iteritems():
                        if k in {'$ref', 'description', 'type', 'items', 'properties'}:
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
                        w('@extends = "{}"', self.schemaFileNameToType(propSchema['extends']['$ref'], pathToSchema.split('#')[0]))
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
                        w('!alias {}:', tname)
                        with w.indent():
                            _extractSingleProperty(loadedSchema, 'value', tname, ['value'])

                # Types 1 or 4:
                if 'type' in loadedSchema:
                    _loadTypes(loadedSchema)
                elif 'properties' in loadedSchema:
                    self.warn("Properties element without accompanying type: %r" % loadedSchema)
                    _extractProperties(loadedSchema)
                else:
                    stName = tname.split('.')[-1]
                    # Type 2
                    if stName in loadedSchema:
                        if 'type' in loadedSchema[stName]:
                            _loadTypes(loadedSchema[stName])
                        elif 'properties' in loadedSchema:
                            self.warn("Properties element without accompanying type: %r" % loadedSchema[stName])
                            _extractProperties(loadedSchema[stName])
                        else:
                            _extractSingleProperty(loadedSchema[stName], 'value', tname, ['value'])
                    # Type 3
                    else:
                        collection = splitPath[-1].split('/')
                        assert stName == collection[1]
                        if 'type' in loadedSchema[collection[0]][stName]:
                            _loadTypes(loadedSchema[collection[0]][stName])
                        elif 'properties' in loadedSchema:
                            self.warn("Properties element without accompanying type: %r" % loadedSchema[collection[0]][stName])
                            _extractProperties(loadedSchema[collection[0]][stName])
                        else:
                            w('!alias {}:', tname)
                            with w.indent():
                                _extractSingleProperty(loadedSchema[collection[0]][stName], 'value', tname, ['value'])

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
                            ftypeSyntax = self.parse_typespec(v, typeSpecList[index].path, k, typeRef)[0]
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

        # build all the types
        for tname, schema in sorted(includedSchemas.iteritems()):
            if 'example' not in tname:
                w()
                _processSchema(schema, tname, w)

        _processTypeSpecList()

    def writeAlias(self, alias, w):
        for key in alias:
            w()
            w('!alias {}:', key)
            with w.indent():
                w(alias[key])

    def parse_typespec(self, tspec, path='', parentRef='', typeRef=''):

        typ = tspec.get('type')

        # skip invalid arrays
        if not typ and 'items' in tspec:
            assert False
            self.warn('Ignoring unexpected "items". Schema has "items" but did not have defined "type". Note: %r' % (tspec, ))
            del tspec['items']

        descr = tspec.pop('description', None)
        if typ == 'array':
            # assert not (set(tspec.keys()) - {'$schema', 'type', 'items', 'example', 'uniqueItems', 'minItems', 'maxItems', 'required'}), tspec

            (itype, idescr) = self.parse_typespec(tspec['items'], path, parentRef, typeRef)

            if '$ref' in tspec['items']:
                itype = self.schemaFileNameToType(itype, path)
            return (sysl_array_type_of(itype), descr)

        def r(t):
            return (t, descr)

        fmt = tspec.get('format')
        ref = self.schemaFileNameToType(tspec.get('$ref'), path)

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
            typeSpecList.append(TypeSpec(tspec, parentRef, typeRef, path))
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

    translator = RamlTranslator(logger=make_default_logger(), raml_path=args.raml_path)
    translator.translate(raml, args.appname, args.package, w=w)

    with open(args.outfile, 'w') as f_out:
        f_out.write(str(w))


if __name__ == '__main__':
    debug.init()
    main()
