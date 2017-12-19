#!/usr/local/bin/python
# -*- encoding: utf-8 -*-

'''Generate JavaScript encoding of model.'''

import argparse
import cStringIO
import json
import numbers
import os
import re

from sysl.proto import sysl_pb2

from sysl.core import syslloader
from sysl.core import syslx

from sysl.util import datamodel
from sysl.util import debug
from sysl.util import java


def js_dump(appname, obj, indent, outfile):
    def write_obj(obj):
        if isinstance(obj, (bool, numbers.Number, type(None), basestring)):
            json.dump(obj, outfile)
        elif isinstance(obj, dict):
            outfile.write('{')
            for (i, (key, value)) in enumerate(sorted(obj.iteritems())):
                if i:
                    outfile.write(',')
                outfile.write(key if re.match(r'\w+$', key)
                              else json.dumps(key))
                outfile.write(':')
                write_obj(value)
            outfile.write('}')
        else:
            outfile.write('[')
            for (i, value) in enumerate(obj):
                if i:
                    outfile.write(',')
                write_obj(value)
            outfile.write(']')

    # Rationale: http://goo.gl/7Bej8a
    outfile.write('''\
((root, factory) => {{
  if (typeof exports === 'object') {{
    module.exports = root.{0} = factory(require('app/sysl'));  // CommonJS
  }} else if (typeof define === 'function' && define.amd) {{
    define(['sysl'], sysl => root.{0} = factory(sysl));  // AMD
  }} else {{
    root.{0} = factory(root.sysl);  // Browser
  }}
}})(this, sysl => {{'use strict';return sysl.createModelClass('''.format(
        appname))

    if indent not in (None, 'compact'):
        guts = json.dumps(obj, indent=indent, separators=',:', sort_keys=True)
        guts = re.sub(r'(?m)^ {1,4}', r'', guts)
        guts = re.sub(r'(?m)^(\S.*{\n)(?=\S)', r'\1\n', guts)
        guts = re.sub(r'(?<=}\n})\n(?=})', r'', guts)
        outfile.write(guts)
    else:
        write_obj(obj)

    outfile.write(');});')


def export_model_js(module, appname, outpath, indent):
    '''The Model as JSON'''
    app = module.apps.get(appname)
    assert app, appname

    model = {}
    model['_version'] = '0.1'
    model['model'] = '_'.join(app.name.part)
    tables = model['types'] = {}

    for (tname, t) in app.types.iteritems():
        titem = tables[tname] = {}
        tmeta = titem['_'] = {}
        if t.HasField('relation'):
            tmeta['rel'] = True

        pkey = datamodel.primary_key_params(t, module)
        pkey_fields = {f for (_, f, _) in pkey}
        param_defs = [(typ, f) for (typ, _, f) in pkey]
        pk_param = ', '.join(f for (_, _, f) in pkey)

        fkeys = {java.name(fname): type_info
                 for (fname, _, type_info) in datamodel.foreign_keys(t, module)}

        for (fname, f) in datamodel.sorted_fields(t):
            jfname = java.name(fname)
            (java_type, type_info, ftype) = datamodel.typeref(f, module)

            which_f = f.WhichOneof('type')
            if which_f == 'primitive':
                ftname = ftype.primitive
            elif fname in fkeys:
                fk_type_info = fkeys[fname]
                ftname = (
                    fk_type_info.parent_path +
                    ('.' + fk_type_info.field if fk_type_info.field != jfname else '='))

            if fname in pkey_fields:
                jfname += '*'
            elif 'required' in syslx.patterns(f.attrs):
                jfname += '!'

            fitem = titem[jfname] = [ftname]

            if len(fitem) == 1:
                fitem = titem[jfname] = fitem[0]

    out = cStringIO.StringIO()
    js_dump(appname, model, indent, out)
    open(outpath, 'w').write(out.getvalue())


def wrapped_types(model, context):
    for (tname, ft) in sorted(app.wrapped.types.iteritems()):
        assert tname in model.types, tname
        t = model.types[tname]
        yield (tname, ft, t)


# Facade class does reference the Model class via the getModel public method.
def export_facade_class(model, appname, context):
    raise NotImplementedError()


def main():
    argp = argparse.ArgumentParser(
        description='sysl relational JavaScript Model exporter')

    argp.add_argument('--root', '-r', default='.',
                      help='sysl system root directory')
    argp.add_argument('--out', '-o', required=True, help='Output file')
    argp.add_argument('--indent', type=int, help='Indentation level')
    argp.add_argument('module', help='Module to load')
    argp.add_argument('app', help='Application to export')
    args = argp.parse_args()

    out = os.path.normpath(args.out)

    (module, _, _) = syslloader.load(args.module, True, args.root)

    export_model_js(module, args.app, out, args.indent)


if __name__ == '__main__':
    debug.init()
    main()
