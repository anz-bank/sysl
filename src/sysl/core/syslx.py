import collections

from sysl.proto import sysl_pb2


def fmt_app_name(appname):
    """Format a sysl_pb2.AppName as a syntactically valid string."""
    return ' :: '.join(appname.part)


def patterns(attrs):
    return frozenset(attr.s for attr in View(attrs)['patterns'].a.elt)


def wrapped_facade_types(context):
    for (tname, ft) in sorted(context.app.wrapped.types.iteritems()):
        assert tname in context.wrapped_model.types, tname
        t = context.wrapped_model.types[tname]
        yield (tname, ft, t)


def sorted_types(context):
    if context.app.HasField('wrapped'):
        return wrapped_facade_types(context)
    return (
        (tname, t, t)
        for (tname, t) in sorted(context.app.types.iteritems()))


class TypeInfoByRef(object):
    """Convenience class to fetch types by sysl_pb2.Type.TypeRef."""

    TYPE_INFO = collections.namedtuple('_TypeInfoByRef__TYPE_INFO',
                                       'path type parent_path parent_type field app')

    def __init__(self, module):
        self._module = module

    def __getitem__(self, typeref):
        for i in range(len(typeref.context.appname.part), -1, -1):
            ref = typeref.ref
            appname = sysl_pb2.AppName()
            appname.part.extend(typeref.context.appname.part[:i])
            appname.part.extend(ref.appname.part)
            app = AppByName(self._module).get(appname)
            if app is None:
                continue

            if not ref.path:
                return self.TYPE_INFO(None, None, None, None, None, app)

            path = u'.'.join(ref.path)
            type_ = app.types.get(path)

            if type_ is not None:
                return self.TYPE_INFO(path, type_, None, None, None, app)

            # typerefs can be type.type.type or type.type.field. The latter
            # implies a foreign key. If this is the case, return the field's
            # type.
            parent_path = u'.'.join(ref.path[:-1])
            parent_type = (
                app.types.get(parent_path) if parent_path else
                app.types.get(ref.path[-1]))
            leaf = ref.path[-1]
            if parent_type.WhichOneof('type') not in ('relation', 'tuple'):
                raise RuntimeError(u'{} . {} not a table or entity'.format(
                    fmt_app_name(ref.appname), parent_path))
            entity = getattr(parent_type, parent_type.WhichOneof('type'))
            type_ = entity.attr_defs[leaf]
            return self.TYPE_INFO(None, type_, parent_path,
                                  parent_type, leaf, app)

        raise RuntimeError(u'typeref not found: ' + str(typeref))


class AppByName(object):
    """Convenience class to fetch apps by sysl_pb2.AppName."""

    def __init__(self, module):
        self._module = module

    def __getitem__(self, appname):
        return self._module.apps[fmt_app_name(appname)]

    def get(self, appname, default=None):
        return self._module.apps.get(fmt_app_name(appname), default)


class View(object):
    """Convenience class to fetch missing attributes."""

    def __init__(self, attrs):
        self._attrs = attrs

    def __getitem__(self, attrname):
        attr = self._attrs.get(attrname)
        return _DefaultAttr(attr)


# TODO: Recursive defaulting.
class _DefaultAttr(object):
    """Placeholder for missing attribute.

    Could have just used an empty sysl_pb2.Attribute, but that would have
    permitted accidentally setting attributes. This class provides read-only
    fields to avoid this.
    """

    def __init__(self, attr):
        self._attr = attr

    # pylint: disable=invalid-name,missing-docstring
    @property
    def s(self):
        return self._attr.s if self._attr else None

    @property
    def i(self):
        return self._attr.i if self._attr else None

    @property
    def n(self):
        return self._attr.n if self._attr else None

    @property
    def a(self):
        return self._attr.a if self._attr else sysl_pb2.Attribute().a
