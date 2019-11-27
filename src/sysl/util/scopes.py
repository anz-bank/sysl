from sysl.proto import sysl_pb2

from sysl.core import syslx


class Scope(object):
    def __init__(self, _parent, _dot=None, **kwargs):
        if isinstance(_parent, Scope):
            self.parent = _parent
            self.module = _parent.module
        else:
            self.parent = None
            self.module = _parent

        self.types = {}

        for (name, type_) in kwargs.iteritems():
            self[name] = type_

        if _dot is not None:
            self['.'] = _dot

    def __getitem__(self, name):
        out = self._get(name)
        if not out:
            raise RuntimeError('Symbol not found in any scope: ' + name)
        return out[0]

    def __setitem__(self, name, type_):
        assert name not in self.types, (name, type_, self.types.keys())
        if name != '__dot__':
            assert isinstance(
                type_, (sysl_pb2.Type, sysl_pb2.Application, type(None)))
        self.types[name] = type_

    def _get(self, name):
        if name in self.types:
            return [self.types[name]]
        if self.parent:
            return self.parent._get(name)
        return None

    def get(self, name):
        out = self._get(name)
        return out[0] if out else None

    def resolve(self, t):
        if isinstance(t, str):
            return self.resolve(self[t])

        which_type = t.WhichOneof('type')

        if which_type == 'set':
            (app, t) = self.resolve(t.set)
            return (app, sysl_pb2.Type(set=t))

        if which_type == 'type_ref':
            ref = t.type_ref.ref

            type_info = syslx.TypeInfoByRef(self.module)[t.type_ref]
            app = type_info.app
            assert app, str(t)

            if not ref.path:
                return (app, app)

            if len(ref.path) > 1:
                raise NotImplementedError('Nested types')

            return (app, type_info.type)

        return (None, t)
