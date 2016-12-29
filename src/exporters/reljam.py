#!/usr/local/bin/python
# -*- encoding: utf-8 -*-

# Copyright 2016 The Sysl Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License."""Super smart code writer."""

import argparse
import collections
import os
import re

from src.sysl import syslloader
from src.sysl import syslx

from src.util import datamodel
from src.util import debug
from src.util import java
from src.util import writer
import src.util.file

import src.exporters.api.spring_rest
import src.exporters.java.model as java_model
import src.exporters.java.facade as java_facade
import src.exporters.json_out.serializer as json_export
import src.exporters.swagger.swagger as swagger_export
import src.exporters.xml.serializer as xml_export
import src.exporters.xml.xsd as xsd_export


Context = collections.namedtuple(
  'Context', 'app module package model_class write_file appname wrapped_model')


def export(
    mode, module, appname, outdir, expected_package, entities, serializers):
  app = module.apps.get(appname)
  assert app, appname

  package = syslx.View(app.attrs)['package'].s
  if mode != 'xsd':
    assert package == expected_package, (package, expected_package)

  model_class = '_'.join(app.name.part).replace(' ', '')

  write_file = src.util.file.FileWriter(outdir, package, entities)

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
  bogus_serializers = serializers - {'json_in', 'json_out', 'xml_in', 'xml_out'}
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
    entities |= { appname, appname + 'Exception'} | serializer_entities()

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
        java_model.export_entity_class(w, tname, t, fk_rmap[tname], context)

    java_model.export_model_class(fk_rmap, context)
    java_model.export_exception_class(context)
    export_serializers()

  elif mode == 'facade':
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
      for endpt in app.endpoints.itervalues()}
    assert None not in interfaces, '\n' + '\n'.join(sorted([
      endpt.name
      for endpt in app.endpoints.itervalues()
      for i in [endpt.attrs['interface'].s]
      if not i],
      key=lambda name: reversed(name.split())))

    entities |= {model_class + 'Controller'} | set(interfaces)
    src.exporters.api.spring_rest.service(interfaces, context)

  elif mode == 'view':
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


def main():
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
  args = argp.parse_args()

  out = os.path.normpath(args.out)

  (module, _, _) = syslloader.load(args.module, True, args.root)

  entities = set(args.entities.split(',')) if args.entities else set()

  export(args.mode, module, args.app, out, args.package, entities,
    args.serializers.split(',') if args.serializers else [])

  if entities:
    raise RuntimeError('Some entities not output as expected: ' +
      ', '.join(sorted(entities)))


if __name__ == '__main__':
  debug.init()
  main()
