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

import os


class FileWriter(object):
  def __init__(self, out_dir, package, entities):
    self.out_dir = out_dir
    self.package = package
    self.entities = entities

  def __call__(self, w, out_path):
    out_path = os.path.join(self.out_dir, out_path)
    try:
      os.makedirs(os.path.dirname(out_path))
    except:
      pass
    open(out_path, 'w').write(str(w))

  def java(self, w, name, package):
    assert name
    assert name in self.entities, ('Unexpected entity generated: ' + name +
      ' (check BUILD has sysl_model(..., entities=[..., "' + name + ', ...], ...))')
    self.entities.remove(name)
    self(w, os.path.join(self.package.replace('.', '/'), name + u'.java'))
