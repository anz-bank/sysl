#!/usr/bin/env python
# -*- encoding: utf-8 -*-

import argparse
import json
import logging
import re
import sys
import textwrap

import yaml

from sysl.util import debug
from sysl.util import writer

from sysl.proto import sysl_pb2


# TODO: dedup with //src/exporters/swagger/swagger.py
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

WORDS = set()

PARAM_CACHE = {}


def words():
    """Lazy-load WORDS."""
    if not WORDS:
        WORDS.update(w.strip() for w in open('/usr/share/dict/words'))
    return WORDS


def type_as_key(swagt):
    if isinstance(swagt, dict):
        return frozenset(sorted(swagt.iteritems()))
    assert isinstance(swagt, basestring)
    return swagt


def javaParam(match):
    param = match.group(1)
    ident = PARAM_CACHE.get(param)

    if ident is None:
        # foo-bar to fooBar
        ident = re.sub(r'-(\w?)', lambda m: m.group(1).upper(), param)

        # {fooid} -> fooId (only if foo is in WORDS but fooid isn't)
        m = re.match(r'{([a-z]+)id}$', ident)
        if m:
            word = m.group(1)
            if word in words() and word + 'id' not in words():
                ident = '{' + word + 'Id}'

        ident = ident.replace('}', '<:string}')

        PARAM_CACHE[param] = ident

    return ident


def extract_properties(tspec):
    # Some schemas for object types might lack a properties attribute.
    # As per Open API Spec 2.0, partially defined via JSON Schema draft 4,
    # if properties is missing we assume it is the empty object.
    # https://tools.ietf.org/html/draft-fge-json-schema-validation-00#section-5.4
    return tspec.get('properties', {})


# Swagger offers multiple ways to describe the same type, hence this mess.
SWAGGER_TYPE_MAP = dict({
    (type_as_key(swagt), sysl_pb2.Type.Primitive.Name(syslt).lower())
    for (syslt, swagt) in TYPE_MAP.iteritems()
    if swagt
} | {
    (swagt['format'], sysl_pb2.Type.Primitive.Name(syslt).lower())
    for (syslt, swagt) in TYPE_MAP.iteritems()
    if swagt and swagt.get('type') == 'number'
} | {
    ('int32', sysl_pb2.Type.INT),
    ('uint32', sysl_pb2.Type.INT),
    ('int64', sysl_pb2.Type.INT),
    ('uint64', sysl_pb2.Type.INT),
})


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


class SwaggerTranslator:
    def __init__(self, logger):
        self._logger = logger

    def warn(self, msg):
        self._logger.warn(msg)

    def translate(self, swag, appname, package, w):

        if 'info' in swag:
            def info_attrs(info, prefix=''):
                for (name, value) in sorted(info.iteritems()):
                    if isinstance(value, dict):
                        info_attrs(value, prefix + name + '.')
                    else:
                        w('@{}{} = {}', prefix, name, json.dumps(value))

            title = swag['info'].pop('title', '')
            info_attrs(swag['info'])
        else:
            title = ''

        if 'host' in swag:
            w('@host = {}', json.dumps(swag['host']))

        w(u'{}{} [package={}]:',
            appname, title and ' ' + json.dumps(title), json.dumps(package))

        with w.indent():
            w(u'| {}', swag['info'].get('description', 'No description.'))

            for (path, api) in sorted(swag['paths'].iteritems()):
                # {foo-bar} to {fooBar}
                w(u'\n{}:', re.sub(r'({[^/]*?})', javaParam, path))
                with w.indent():
                    if 'parameters' in api:
                        del api['parameters']
                    for (i, (method, body)) in enumerate(sorted(api.iteritems(),
                                                                key=lambda t: METHOD_ORDER[t[0]])):
                        qparams = dict()

                        if 'parameters' in body and 'in' in body['parameters']:
                            qparams = [p for p in body['parameters']
                                       if p['in'] == 'query']
                        w(u'{}{}{}:',
                          method.upper(),
                          ' ?' if qparams else '',
                          '&'.join(('{}={}{}'.format(p['name'],
                                                     SWAGGER_TYPE_MAP[p['type']],
                                                     '' if p['required'] else '?')
                                    if p['type'] != 'string' else '{name}=string'.format(**p))
                                   for p in qparams))
                        with w.indent():
                            for line in textwrap.wrap(
                                    body.get('description', 'No description.').strip(), 64):
                                w(u'| {}', line)

                            responses = body['responses']
                            errors = ','.join(
                                sorted(str(e) for e in responses if e >= 400))

                            if 200 in responses:
                                r200 = responses[200]
                                if 'schema' in r200:
                                    ok = r200['schema']
                                    if ok.get('type') == 'array':
                                        items = ok['items']
                                        if '$ref' in items:
                                            itemtype = items['$ref'][
                                                len('#/definitions/'):]
                                            ret = ': <: set of ' + itemtype
                                        else:
                                            ret = ': <: ...'
                                    elif '$ref' in ok:
                                        ret = ': <: ' + ok['$ref'][
                                            len('#/definitions/'):]
                                    else:
                                        ret = ' (' + r200.get('description') + ')'
                                else:
                                    ret = ' (' + r200.get('description') + ')'
                                w(u'return 200{} or {{{}}}', ret, errors)
                            elif 201 in responses:
                                r201 = responses[201]
                                if 'headers' in r201:
                                    ok = r201['headers']
                                    w(u'return 201 ({}) or {{{}}}',
                                      ok['Location']['description'],
                                      errors)
                                else:
                                    w(u'return 201 ({})',
                                      r201['description'])

                        if i < len(api) - 1:
                            w()

            w()
            w('#' + '-' * 75)
            w('# definitions')

            for (tname, tspec) in sorted(swag['definitions'].iteritems()):
                w()
                w('!type {}:', tname)

                with w.indent():
                    properties = extract_properties(tspec)
                    if properties:
                        for (fname, fspec) in sorted(properties.iteritems()):
                            (ftype, fdescr) = self.parse_typespec(fspec)
                            w('{} <: {}{}',
                              fname,
                              ftype if ftype.startswith('set of ') or ftype.endswith('*') else ftype + '?',
                              ' "' + fdescr + '"' if fdescr else '')
                    # handle top-level arrays
                    elif tspec.get('type') == 'array':

                        (ftype, fdescr) = self.parse_typespec(tspec)
                        w('{} <: {}{}',
                          fname,
                          ftype if ftype.startswith('set of ') or ftype.endswith('*') else ftype + '?',
                          ' "' + fdescr + '"' if fdescr else '')
                    else:
                        assert True, tspec

    def parse_typespec(self, tspec):

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

            (itype, idescr) = self.parse_typespec(tspec['items'])
            assert idescr is None
            return ('set of ' + itype, descr)

        def r(t):
            return (t, descr)

        fmt = tspec.get('format')
        ref = tspec.get('$ref')

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
            descrs = {
                k: descr
                for (k, v) in extract_properties(tspec).iteritems()
                for descr in [self.parse_typespec(v)[1]]
                if descr
            }
            assert not descrs, descrs
            fields = ('{} <: {}'.format(k, self.parse_typespec(v)[0])
                      for (k, v) in sorted(extract_properties(tspec).iteritems()))
            return r('{' + ', '.join(fields) + '}')
        else:
            return r(str(tspec))


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
