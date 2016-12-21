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

# TODO: Add no_except param to sysl_model.

VALID_SERIALIZERS = set(['json_in', 'json_out', 'xml_in', 'xml_out'])


def _inouts(prefix, serializers):
  if serializers == None:
    serializers = ["*_*"]

  inouts = []
  for s in serializers or []:
    (fmt, dirn) = s.split('_')
    fmts = ['json', 'xml'] if fmt == '*' else [fmt]
    dirns = ['in', 'out'] if dirn == '*' else [dirn]
    for fmt in fmts:
      for dirn in dirns:
        inouts.append(fmt + "_" + dirn)
  bogus_serializers = [
    s for s in inouts
    if s not in VALID_SERIALIZERS]
  if bogus_serializers:
    fail("invalid serializers: " + bogus_serializers, attr="serializer")

  files = (
    ([prefix + "JsonDeserializer.java"] if 'json_in' in inouts else []) +
    ([prefix + "JsonSerializer.java"] if 'json_out' in inouts else []) +
    ([prefix + "XmlDeserializer.java"] if 'xml_in' in inouts else []) +
    ([prefix + "XmlSerializer.java"] if 'xml_out' in inouts else [])
  )
  return (inouts, files)


def sysl_model(
    name, srcs, root, module, app, entities, package,
    deps=None,
    serializers=None,  # <xml|json|*>_<in|out|*>
    visibility=None):
  outpath = package.replace(".", "/") + "/"
  prefix = outpath + app

  (inouts, serializer_files) = _inouts(prefix, serializers)

  _sysl_tool(
    name = name,
    tool = "model",
    srcs = srcs,
    root = root,
    module = module,
    app = app,
    package = package,
    entities = entities,
    jar = True,
    outs = (
      [
        prefix + ".java",
        prefix + "Exception.java",
      ] +
      serializer_files +
      [outpath + e + ".java" for e in entities]
    ),
    deps = [
      "//java/io/sysl",
      "@org_apache_commons_commons_lang3//jar",
      "@com_fasterxml_jackson_core_jackson_core//jar",
      "@com_fasterxml_jackson_core_jackson_databind//jar",
      "@joda_time_joda_time//jar",
    ] + (deps or []),
    flags = [
      '--serializers=' + ','.join(inouts),
    ],
    visibility = visibility,
  )


def sysl_facade(name, srcs, root, module, app, model, package,
    serializers=None,  # <xml|json|*>_<in|out|*>
    visibility=None):
  outpath = package.replace(".", "/") + "/"
  prefix = outpath + app

  (inouts, serializer_files) = _inouts(prefix, serializers)

  _sysl_tool(
    name = name,
    tool = "facade",
    srcs = srcs,
    root = root,
    module = module,
    app = app,
    package = package,
    outs = [prefix + ".java"] + serializer_files,
    deps = [
      "//java/io/sysl",
      "@com_fasterxml_jackson_core_jackson_core//jar",
      "@com_fasterxml_jackson_core_jackson_databind//jar",
      "@joda_time_joda_time//jar",
      "@org_apache_commons_commons_lang3//jar",
      model
    ],
    flags = [
      '--serializers=' + ','.join(inouts),
    ],
    jar = True,
    visibility = visibility,
  )


def sysl_xsd(name, srcs, root, module, app, visibility=None):
  _sysl_tool(
    name = name,
    tool = "xsd",
    srcs = srcs,
    root = root,
    module = module,
    app = app,
    package = '',
    outs = [app + ".xsd"],
    visibility = visibility,
  )


def sysl_swagger(name, srcs, root, module, app, package,
    outbase=None,
    visibility=None):
  _sysl_tool(
    name = name,
    tool = "swagger",
    srcs = srcs,
    root = root,
    module = module,
    app = app,
    package = package,
    outs = [app.replace(' :: ', '_').replace(' ', '') + ".swagger.yaml"],
    visibility = visibility,
  )


def sysl_spring_rest_service(name, srcs, root, module, app, package,
    entities=None,
    outbase=None,
    visibility=None):
  entities = entities or []
  outpath = package.replace(".", "/") + "/"

  _sysl_tool(
    name = name,
    tool = "spring-rest-service",
    srcs = srcs,
    root = root,
    module = module,
    app = app,
    package = package,
    entities = entities,
    jar = True,
    outs = [
      outpath + app + "Controller.java",
    ] + [outpath + e + ".java" for e in entities],
    deps = [
      # TODO: external deps: lombok, spring {web,beans}, swagger annotations
    ],
    visibility = visibility,
  )


# private

def _sysl_tool(
    name, tool, srcs, root, module, app, package,
    entities=None,
    jar=False,
    outs=None,
    deps=None,
    flags=None,
    visibility=None):

  if entities:
    native.filegroup(
      name = name + "_entities",
      srcs = outs,
      visibility = visibility,
    )
  else:
    entities = []

  native.genrule(
    name = name + "_dummy",
    outs = ["." + name + "_dummy"],
    cmd = "touch '$@'",
  )

  native.genrule(
    name = name + "_java" if jar else name,
    srcs = srcs + [":" + name + "_dummy"],
    outs = outs,
    cmd = " ".join([
        "$(location //src/exporters:reljam)",
        "--root '%s'" % root,
        "--out '$(location :%s_dummy)/..'" % name,
        "--package '%s'" % package,
      ] + (["--entities '%s'" % ','.join(entities)] if entities else []) +
      (flags or []) + [
        tool,
        "'%s'" % module,
        "'%s'" % app,
      ]),
    tools = ["//src/exporters:reljam"],
    visibility = visibility,
  )

  if jar:
    native.java_library(
      name = name,
      srcs = outs + ([name + "_entities"] if entities else []),
      deps = deps or [],
      visibility = visibility,
    )
