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

def export_plantuml(model, with_fields):
  w = writer.Writer('plantuml')

  with w.uml():
    for (tname, ft, t) in wrapped_types(model):
      if t.HasField(u'relation'):
        pkey = primary_key_params(t)
        pkey_fields = {f for (_, f, _) in pkey}
        fkey_fields = {f for (f, _, _) in foreign_keys(t)}

        # Yes, I know; PlantUML isn't Java. But this fits the bill.
        with java.Block(w, 'class {}', tname):
          if with_fields:
            for (fname, f) in t.relation.attr_defs.iteritems():
              if fname in pkey_fields:
                w('{{static}} +{}', java.name(fname))
            for (fname, f) in t.relation.attr_defs.iteritems():
              (_, _, type_) = typeref(f)
              if fname not in pkey_fields:
                if fname in fkey_fields:
                  w('+{}', fname)

    for (tname, t) in sorted(app.types.iteritems()):
      if t.HasField(u'relation'):
        for (fname, _, type_info) in foreign_keys(t):
          w('{} -- {}{}',
            type_info.parent_path,
            '' if {fname} == pkey_fields else '"0..*" ',
            tname)

  diagutil.output_plantuml(
    diagutil.OutputArgs(
      output=root_class + '.png',
      plantuml=None,
      verbose=True,
      expire_cache=False,
      dry_run=False),
    str(w))
