import json
import numbers
import re


JSONNET_KEYWORDS = {
    'else', 'error', 'false', 'for', 'function', 'if', 'import', 'importstr',
    'in', 'local', 'null', 'self', 'super', 'then', 'true',
}
IDENT_RE = re.compile(r'[A-Za-z_]\w*$')


def safe_key(key):
    if IDENT_RE.match(key) and key not in JSONNET_KEYWORDS:
        return key
    return json.dumps(key)


def dumps(obj):
    if isinstance(obj, dict):
        return '{' + ', '.join(
            '{}: {}'.format(safe_key(k), dumps(v))
            for (k, v) in obj.iteritems()) + '}'
    if isinstance(obj, list):
        return '[' + ', '.join(dumps(e) for e in obj) + ']'
    if isinstance(obj, Code):
        return obj.code
    return json.dumps(obj)


class Code(object):
    def __init__(self, code):
        self.code = code


def match(obj, pattern):
    """Match a JSON pattern, binding portions of it to a result object.

    args:
      obj: object to be matched
      pattern: lambda taking binder and returning structure
        structure : pat (>> binding)?

        pat       : <type>                  # instance of type
                  | None                    # equals None
                  | <bool>                  # equals Boolean value
                  | <number>                # equals number
                  | <basestring>            # equals string
                  | [structure]             # list
                  | {key: structure, ...}   # dict
                  | (structure, ...)        # either-or

        key       : <basestring>            # required
                  | <basestring> + "?"      # optional
                  | ()                      # any key

        binding   : <binder>                # binder = value
                  | <binder>.<attr>         # binder.attr = value
                  | <binder>.<attr>.setitem # binder.attr[field] = value
                                            # (only valid for object fields)
                  | binding.as_json         # binding = json(value)
                  | binding.omit_keys(<field>, ...)
                                            # binding = value sans <field>, ...

    Simple example:
    >>> obj = {'x': [2, 3, 4], 'y': 42}
    >>> b = match(obj, lambda b: {'x': [int] >> b.x, 'y': 42 })
    >>> assert b.x == [2, 3, 4]
    """
    path = []
    objstk = []
    patstk = []
    logs = []

    def log(fmt='', *args, **kwargs):
        logs.append(
            (('{}: ' + fmt).format(_pathf(path), *args, **kwargs),
             objstk[-1],
             patstk[-1]))

    def descend(obj, pattern, subpath=None):
        subpath = subpath or []
        path.extend(subpath)
        objstk.append(obj)
        patstk.append(pattern)
        try:
            if isinstance(pattern, _Assign):
                result = descend(obj, pattern.pattern, [pattern])
                if result:
                    pattern(path, result[0])
                return result

            if isinstance(pattern, tuple):
                for pat in pattern:
                    result = descend(obj, pat)
                    if result:
                        return result
                return None

            if isinstance(pattern, type):
                return ([obj] if isinstance(obj, pattern) else
                        log('type(obj) {!r} != expected {!r}', type(obj), type(pattern)))

            if isinstance(pattern, (type(None), bool,
                                    numbers.Number, basestring)):
                return ([obj] if obj == pattern else
                        log('value {!r} != expected {!r}', obj, pattern))

            if isinstance(pattern, dict):
                if not isinstance(obj, dict):
                    return log('{!r} not an object', obj)

                obj_keys = set(obj)
                required = {
                    k for k in pattern
                    if isinstance(k, basestring) and not k.endswith('?')} - {()}
                optional = {
                    k[:-1] for k in pattern
                    if isinstance(k, basestring) and k.endswith('?')}

                if not (required <= obj_keys):
                    return log('missing field(s) {{{}}}', ', '.join(
                        sorted(required)))
                surplus = obj_keys - required - optional
                if surplus:
                    if () not in pattern:
                        return log(
                            'surplus field(s) {{{}}}', ', '.join(sorted(surplus)))
                    else:
                        pat = pattern[()]
                        for k in surplus:
                            result = descend(obj[k], pat, [k + '+'])
                            if not result:
                                return None
                for k in obj_keys & set(required):
                    if not descend(obj[k], pattern[k], [k]):
                        return False
                for k in obj_keys & set(optional):
                    if not descend(obj[k], pattern[k + '?'], [k]):
                        return False
            elif isinstance(pattern, list):
                assert len(pattern) == 1
                if not isinstance(obj, list):
                    return log('{!r} not an array', obj)
                pat = pattern[0]
                for (i, e) in enumerate(obj):
                    if not descend(e, pat, [i]):
                        return False
            return [obj]
        finally:
            objstk.pop()
            patstk.pop()
            del path[len(path) - len(subpath):]

    result = descend(obj, pattern)
    # for (i, (log, _, _)) in enumerate(logs):
    #   import sys; print >>sys.stderr, '    ' if i else '### ', log
    return result


class _Result(object):
    def __repr__(self):
        return '<>'

    def __getattr__(self, _):
        return None


class _JsonGetter(object):
    def __init__(self, parent):
        self._parent = parent

    def __getattr__(self, name):
        return json.dumps(getattr(self._parent, name))


class Matcher(_Result):
    def __init__(self):
        self._result = None

    def __call__(self, obj, pattern):
        self._result = _Result()
        binder = _Binder(self._result)
        return bool(match(obj, pattern(binder)))

    def __getattr__(self, name):
        return getattr(self._result, name)


class _Assign(object):
    def __init__(self, binding, pattern):
        self._binding = binding
        self.pattern = pattern

    def __call__(self, path, value):
        return self._binding(path, value)

    def __repr__(self):
        return '{} >> {}'.format(self.pattern, self._binding)


class _Binding(object):
    def __rrshift__(self, pattern):
        return _Assign(self, pattern)

    @property
    def as_json(self):
        return _AsJson(self)

    def __sub__(self, keys):
        return _OmitKeys(self, keys)


class _BindAttr(_Binding):
    def __init__(self, binder, name):
        self._binder = binder
        self._name = name

    @property
    def setitem(self):
        return _SetItem(self)

    def __call__(self, path, value):
        setattr(self._binder._result, self._name, value)
        return value

    def __repr__(self):
        return '{}.{}'.format(self._binder, self._name)


class _SetItem(_Binding):
    def __init__(self, bindattr):
        self._bindattr = bindattr
        result = self._bindattr._binder._result
        name = self._bindattr._name
        if not result.name:
            self._obj = {}
            setattr(result, name, self._obj)
        else:
            self._obj = getattr(result, name)

    def __call__(self, path, value):
        key = path[-1]
        if key.endswith('+'):
            key = key[:-1]
        self._obj[key] = value
        return value

    def __repr__(self):
        return '{}-{{{}}}'.format(self._binding, ', '.join(sorted(self._keys)))


class _AsJson(_Binding):
    def __init__(self, binding):
        self._binding = binding

    def __call__(self, path, value):
        return self._binding(path, dumps(value))

    def __repr__(self):
        return '{}.as_json'.format(self._binding)


class _OmitKeys(_Binding):
    def __init__(self, binding, keys):
        self._binding = binding
        self._keys = keys

    def __call__(self, path, value):
        result = {}
        for attr in value:
            if not attr.startswith('_') and attr not in self._keys:
                result[attr] = value[attr]
        return self._binding(path, result)

    def __repr__(self):
        return '{}-{{{}}}'.format(self._binding, ', '.join(sorted(self._keys)))


class _Binder(object):
    def __init__(self, result):
        self._result = result

    def __getattr__(self, name):
        return _BindAttr(self, name)

    def __repr__(self):
        return repr(self._result)


def _pathf(path):
    """Format a JSON access path."""
    return 'obj' + ''.join(
        '[{}]'.format(p) if isinstance(p, int) else
        '{{{}}}'.format(p._binding) if isinstance(p, _Assign) else
        '.{}'.format(p)
        for p in path)
