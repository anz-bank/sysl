#!/usr/local/bin/python
# -*- encoding: utf-8 -*-

import argparse
import collections
import os
import re
import sys

from sysl.core import syslloader
from sysl.core import syslx

from sysl.util import datamodel
from sysl.util import java
from sysl.util import writer
import sysl.util.file

import sysl.exporters.api.spring_rest
import sysl.exporters.java.model as java_model
import sysl.exporters.java.facade as java_facade
import sysl.exporters.json_out.serializer as json_export
import sysl.exporters.swagger.swagger as swagger_export
import sysl.exporters.xml.serializer as xml_export
import sysl.exporters.xml.xsd as xsd_export
from sysl.__version__ import __version__


Context = collections.namedtuple(
    'Context', 'app module package model_class write_file appname wrapped_model')


def export(
        mode, module, appname, outdir, expected_package, entities, serializers):
    app = module.apps.get(appname)
    assert app, appname

    package = syslx.View(app.attrs)['package'].s
    if mode != 'xsd' and expected_package is not None:
        assert package == expected_package, (package, expected_package)

    model_class = '_'.join(app.name.part).replace(' ', '')

    write_file = sysl.util.file.FileWriter(outdir, package, entities)

    inouts = []
    for s in serializers:
        assert '_' in s, s
        (fmt, dirn) = s.split('_')
        fmts = ['json', 'xml'] if fmt == '*' else [fmt]
        dirns = ['in', 'out'] if dirn == '*' else [dirn]
        for fmt in fmts:
            for dirn in dirns:
                inouts.append(fmt + "_" + dirn)
    serializers = set(inouts)
    bogus_serializers = serializers - \
        {'json_in', 'json_out', 'xml_in', 'xml_out'}
    assert not bogus_serializers, bogus_serializers

    if app.HasField('wrapped'):
        model_name = syslx.fmt_app_name(app.wrapped.name)
        assert model_name in module.apps, (
            'missing app: ' + repr(model_name),
            module.apps.keys(),
            app.wrapped.endpoints.keys())
        wrapped_model = module.apps[model_name]
    else:
        wrapped_model = None

    context = Context(
        app, module, package, model_class, write_file, appname, wrapped_model)

    def serializer_entities():
        return (
            ({appname + 'JsonDeserializer'} if 'json_in' in serializers else set()) |
            ({appname + 'JsonSerializer'} if 'json_out' in serializers else set()) |
            ({appname + 'XmlDeserializer'} if 'xml_in' in serializers else set()) |
            ({appname + 'XmlSerializer'} if 'xml_out' in serializers else set())
        )

    def export_serializers():
        if 'json_in' in serializers:
            json_export.deserializer(context)
        if 'json_out' in serializers:
            json_export.serializer(context)
        if 'xml_in' in serializers:
            xml_export.deserializer(context)
        if 'xml_out' in serializers:
            xml_export.serializer(context)

    if mode == 'model':
        if entities:
            entities |= {appname, appname +
                         'Exception'} | serializer_entities()

        # Build a foreign key reverse map to enable efficeint navigation
        # in the generated classes
        fk_rmap = datamodel.build_fk_reverse_map(app, module)

        # For each of the "types" in the Application message from the
        # protocol buffer represenation of the sysl generate an
        # "Entity" class
        for (tname, t) in sorted(app.types.iteritems()):
            if not re.match(r'AnonType_\d+__$', tname):
                w = writer.Writer('java')
                java.Package(w, package)
                java.StandardImports(w)
                java_model.export_entity_class(
                    w, tname, t, fk_rmap[tname], context)

        java_model.export_model_class(fk_rmap, context)
        java_model.export_exception_class(context)
        export_serializers()

    elif mode == 'facade':
        if entities:
            entities |= {appname} | serializer_entities()
        java_facade.export_facade_class(context)
        export_serializers()

    elif mode == 'xsd':
        xsd_export.xsd(context)

    elif mode == 'swagger':
        swagger_export.swagger_file(app, module, model_class, write_file)

    elif mode == 'spring-rest-service':
        interfaces = {
            endpt.attrs['interface'].s
            for endpt in app.endpoints.itervalues()
            if endpt.attrs['interface'].s != ''}
        assert None not in interfaces, '\n' + '\n'.join(sorted([
            endpt.name
            for endpt in app.endpoints.itervalues()
            for i in [endpt.attrs['interface'].s]
            if not i],
            key=lambda name: reversed(name.split())))

        if entities:
            entities |= {model_class + 'Controller'} | set(interfaces)
        sysl.exporters.api.spring_rest.service(interfaces, context)

    elif mode == 'view':
        if entities:
            entities |= {appname}

        w = writer.Writer('java')
        java.Package(w, package)
        java.StandardImports(w)
        w()
        java.Import(w, 'org.joda.time.DateTime')
        java.Import(w, 'org.joda.time.format.DateTimeFormat')
        java.Import(w, 'org.joda.time.format.DateTimeFormatter')
        w()
        java_model.export_view_class(w, context)


def main(input_args=sys.argv[1:]):
    argp = argparse.ArgumentParser(
        description='sysl relational Java Model exporter')

    argp.add_argument(
        '--root', '-r', default='.',
        help='sysl system root directory')
    argp.add_argument(
        '--out', '-o', default='.',
        help='Output root directory')
    argp.add_argument(
        '--entities',
        help=(
            'Commalist of entities that are expected to have corresponding '
            'output files generated.  This is for verification only.  It doesn’t '
            'determine which files are output.'))
    argp.add_argument(
        '--package',
        help=(
            'Package expected to be used for generated classes. This is for '
            'verification only.  It doesn’t determine the package used.'))
    argp.add_argument(
        '--serializers', default='*_*',
        help='Control output of XML and JSON serialization code.')
    argp.add_argument(
        'mode',
        choices=[
            'model',
            'facade',
            'view',
            'xsd',
            'swagger',
            'spring-rest-service',
        ],
        help='Code generation mode')
    argp.add_argument(
        'module',
        help='Module to load')
    argp.add_argument(
        'app',
        help='Application to export')
    argp.add_argument('--version', '-v',
                      help='show version number (semver.org standard)',
                      action='version', version='%(prog)s ' + __version__)

    args = argp.parse_args(input_args)

    out = os.path.normpath(args.out)

    (module, _, _) = syslloader.load(args.module, True, args.root)

    entities = set(args.entities.split(',')) if args.entities else None

    export(args.mode, module, args.app, out, args.package, entities,
           args.serializers.split(',') if args.serializers else [])

    if entities:
        raise RuntimeError('Some entities not output as expected: ' +
                           ', '.join(sorted(entities)))
