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

# -*- mode: utf-8 -*-

import collections

def topo_sort(G):
  """Topologically sort a graph.

  params:
    G: A dict of nodes to sets of target nodes.
  yields:
    (depth, node) pairs, in topological order
  """

  S = {(0, a) for (a, bb) in G.iteritems()}

  # Count incoming edges per node and exclude target nodes from S.
  R = collections.defaultdict(int)
  for bb in G.itervalues():
    for b in bb:
      R[b] += 1
      S.discard((0, b))

  while S:
    (n, a) = S.pop()
    yield (n, a)
    for b in G.pop(a, ()):
      R[b] -= 1
      if not R[b]:
        S.add((n + 1, b))

  if G:
    raise RuntimeError('cycles detected')
