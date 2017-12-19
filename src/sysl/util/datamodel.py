import collections

from sysl.proto import sysl_pb2

from sysl.core import syslx

from sysl.util import algo
from sysl.util import java


def typeref(type_, module):
    which_type = type_.WhichOneof('type')
    if which_type == 'primitive':
        return (java.TYPE_MAP[type_.primitive], None, type_)

    if which_type == 'type_ref':
        type_ref = type_.type_ref
        # TODO handle packages and subtypes
        type_info = syslx.TypeInfoByRef(module)[type_ref]
        which_type = type_info.type.WhichOneof('type')
        if which_type == 'enum':
            return (type_info.path, type_info, type_info.type)
        if which_type == 'primitive':
            return (
                java.TYPE_MAP[type_info.type.primitive],
                type_info,
                type_info.type)
        if which_type in ['tuple', 'relation']:
            if (type_ref.ref.appname and
                    type_ref.context.appname != type_ref.ref.appname):
                prefix = type_info.app.attrs['package'].s + '.'
            else:
                prefix = ''
            return (prefix + type_info.path, type_info, type_info.type)
        if which_type == 'type_ref':
            (tr_type, _, tr_info_type) = typeref(type_info.type, module)
            return (tr_type, type_info, tr_info_type)
        raise RuntimeError(
            'type_ref must refer to primitive, enum, tuple, or relation, not ' +
            str(which_type))

    if which_type == 'set':
        s = type_.set
        (set_type, set_info, set_info_type) = typeref(s, module)
        if set_info:
            return (set_type + '.View', set_info, set_info_type)
        elif s.HasField('primitive'):
            return ('HashSet<' + java.TYPE_MAP[s.primitive] + '>', None, None)
        assert False

    return ('Object', None, None)


def sorted_by_source_line(items):
    return sorted(items.iteritems(),
                  key=lambda t: t[1].source_context.start.line)


def sorted_fields(t):
    if isinstance(t, sysl_pb2.Type.Relation):
        return sorted_fields_relation(t)
    if isinstance(t, sysl_pb2.Type.Tuple):
        return sorted_fields_tuple(t)

    which_t = t.WhichOneof('type')

    if which_t == 'relation':
        return sorted_fields_relation(t.relation)
    elif which_t == 'tuple':
        return sorted_fields_tuple(t.tuple)
    else:
        return []


def sorted_fields_relation(r):
    pkey_fields = (r.primary_key.attr_name
                   if isinstance(r, sysl_pb2.Type.Relation)
                   else [])
    return sorted(r.attr_defs.iteritems(),
                  key=lambda e: (
        e[0] not in pkey_fields,
        (e[0] not in pkey_fields) ^ e[1].HasField('type_ref'),
        e[0].lower()))


def sorted_fields_tuple(t):
    return sorted(t.attr_defs.iteritems(),
                  key=lambda e: e[0].lower())


def primary_key_params(t, module):
    if t.HasField('tuple'):
        return []
    if t.HasField('relation'):
        pkey = t.relation.primary_key.attr_name
        return [
            (typeref(f, module)[0], fname, java.name(fname))
            for (fname, f) in sorted_fields(t)
            if fname in pkey
        ]


def foreign_keys(t, module):
    for (fname, f) in sorted_fields(t):
        (java_type, type_info, _) = typeref(f, module)
        if type_info and type_info.parent_path:
            yield (fname, java_type, type_info)


def fk_topo_sort(types, module):
    G = {}
    for (tname, t) in types.iteritems():
        G[tname] = {ti.parent_path for (_, _, ti) in foreign_keys(t, module)}
    return ((p, types[p])
            for (_, p) in sorted(algo.topo_sort(G), reverse=True))


def build_fk_reverse_map(app, module):
    fk_rmap = collections.defaultdict(lambda: collections.defaultdict(set))

    for (tname, t) in app.types.iteritems():
        if t.HasField('relation'):
            for (fname, java_type, type_info) in foreign_keys(t, module):
                fk_type = type_info.parent_path
                fk_field = type_info.field
                fk_rmap[fk_type][fk_field].add((tname, fname))

    return fk_rmap


def all_fields():
    dictdictset = (
        lambda: collections.defaultdict(
            lambda: collections.defaultdict(set)))
    result = dictdictset()
    ambiguity_check = dictdictset()

    for (tname, t) in app.types.iteritems():
        if t.WhichOneof('type') in ['relation', 'tuple']:
            for (fname, f) in t.relation.attr_defs.iteritems():
                (java_type, _, _) = typeref(f, module)
                result[fname][java_type].add(tname)
                ambiguity_check[fname][re.sub(
                    '\.View$', '', java_type)].add(tname)

    ambiguous = {fname: ftypes
                 for (fname, ftypes) in ambiguity_check.iteritems()
                 if len(ftypes) > 1}
    if ambiguous:
        err('\033[1;33mRELJAM: WARNING\033[0m: Ambiguously typed fields:')
        for (fname, ftypes) in sorted(ambiguous.iteritems()):
            err('  {}:', fname)
            for (ftype, tnames) in sorted(ftypes.iteritems()):
                err('    \033[1;36m{}\033[0m: {}', ftype, ', '.join(tnames))

    return {fname: next(iter(ftypes))
            for (fname, ftypes) in result.iteritems()}
