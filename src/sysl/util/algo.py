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
