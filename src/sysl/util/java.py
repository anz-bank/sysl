# -*- encoding: utf-8 -*-

import contextlib
import decimal
import itertools
import json
import re
import struct
import sys
import uuid

from sysl.proto import sysl_pb2

from sysl.core import syslx

import datamodel
import scopes
import writer


GET_DELIM = '<<' + uuid.uuid4().hex + '>>'
CODE_MARKER = '<<' + uuid.uuid4().hex + '>>'
CODE_MARKER2 = '<<' + uuid.uuid4().hex + '>>'

CAMEL_CASE_RE = re.compile('_+([A-Za-z0-9])')

BOOL_TYPE = sysl_pb2.Type(primitive=sysl_pb2.Type.BOOL)

TYPE_MAP = {
    sysl_pb2.Type.ANY: u'Object',
    sysl_pb2.Type.BOOL: u'Boolean',
    sysl_pb2.Type.INT: u'Integer',
    sysl_pb2.Type.FLOAT: u'Double',
    sysl_pb2.Type.DECIMAL: u'BigDecimal',
    sysl_pb2.Type.STRING: u'String',
    sysl_pb2.Type.BYTES: u'byte[]',
    sysl_pb2.Type.STRING_8: u'String',
    sysl_pb2.Type.DATE: u'LocalDate',
    sysl_pb2.Type.DATETIME: u'DateTime',
    sysl_pb2.Type.XML: u'String',
    sysl_pb2.Type.UUID: u'UUID',
}

JAVA_KEYWORDS = frozenset(

    'abstract assert boolean break byte case catch char class const continue '
    'default do double else extends false final finally float for goto if '
    'implements import instanceof int interface long native new null package '
    'private protected public return short static super switch synchronized '
    'this throw throws transient true try void volatile while'

    .split()
)

# TODO: replace (2, '...') keys with just '...'. Ditto (1, '...').
PRECEDENCE = {key: i
              for (i, (c, ops)) in enumerate([
                  (3, '?:'),
                  (2, 'COALESCE'),
                  (2, 'BUTNOT'),
                  (2, 'OR'),
                  (2, 'AND'),
                  (2, 'BITOR'),
                  (2, 'BITXOR'),
                  (2, 'BITAND'),
                  (2, 'EQ NE'),
                  (2, 'LT LE GT GE IN CONTAINS NOT_IN NOT_CONTAINS'),
                  (2, 'ADD SUB'),
                  (2, 'MUL DIV MOD'),
                  (1, 'POS NEG NOT INV'),
                  (2, 'POW'),
                  (0, ['.',
                       (1, sysl_pb2.Expr.UnExpr.SINGLE),
                       (1, sysl_pb2.Expr.UnExpr.SINGLE_OR_NULL),
                       (2, sysl_pb2.Expr.BinExpr.WHERE),
                       (2, sysl_pb2.Expr.BinExpr.FLATTEN),
                       (2, sysl_pb2.Expr.BinExpr.TO_MATCHING),
                       (2, sysl_pb2.Expr.BinExpr.TO_NOT_MATCHING),
                       ]),
                  #(9, 'SUM MIN MAX AVERAGE'),
                  (0, ''),
              ])
              for pc in [{
                  0: lambda op: [op],
                  1: lambda op: [(1, op), (1, sysl_pb2.Expr.UnExpr.Op.Value(op))],
                  2: lambda op: [(2, op), (2, sysl_pb2.Expr.BinExpr.Op.Value(op))],
                  3: lambda op: [op],
              }[c]]
              for op in (ops.split() if isinstance(ops, basestring) else ops)

              # + [i] maps precedence to itself
              for key in pc(op) + [i]
              }
PRECEDENCE[-1] = -1

CLASS_STACK = []


def CamelCase(name):
    return name[:1].upper() + name[1:]


def mixedCase(name):
    return name[:1].lower() + name[1:]


def underscore_to_capscase(name):
    if not name:
        return name
    parts = CAMEL_CASE_RE.split(name)
    return u''.join(p.upper() if i % 2 else p for (i, p) in enumerate(parts))


def safe(name):
    if name in JAVA_KEYWORDS:
        return name + '_'
    return name


def name(fieldname):
    return safe(underscore_to_capscase(fieldname))


@contextlib.contextmanager
def Block(w, _control=None, *args, **kwargs):
    if _control is None:
        assert not kwargs
        w('{{')
    else:
        w(_control + u' {{', *args, **kwargs)
    with w.indent():
        yield
    w(u'}}')


def If(w, _control, *args, **kwargs):
    return Block(w, 'if (' + _control + ')', *args, **kwargs)


@contextlib.contextmanager
def ElseIf(w, _control, *args, **kwargs):
    with w.indent(-w.increment):
        w('}} else if (' + _control + ') {{', *args, **kwargs)
    yield


@contextlib.contextmanager
def Else(w):
    with w.indent(-w.increment):
        w('}} else {{')
    yield


def For(w, _control, *args, **kwargs):
    return Block(w, 'for (' + _control + ')', *args, **kwargs)


def While(w, _control, *args, **kwargs):
    return Block(w, 'while (' + _control + ')', *args, **kwargs)


def Switch(w, _control, *args, **kwargs):
    return Block(w, 'switch (' + _control + ')', *args, **kwargs)


@contextlib.contextmanager
def Case(w, fmt, *args, **kwargs):
    w(u'case ' + fmt + ':', *args, **kwargs)
    with w.indent():
        yield


@contextlib.contextmanager
def Default(w):
    w(u'default:')
    with w.indent():
        yield


def Try(w):
    return Block(w, 'try')


@contextlib.contextmanager
def Catch(w, etype, evar):
    with w.indent(-w.increment):
        w('}} catch ({} {}) {{', etype, evar)
    yield


def Import(w, cls):
    w.head(u'import {};', cls)


def StandardImports(w):
    Import(w, 'java.lang.Iterable')
    Import(w, 'java.lang.Boolean')
    Import(w, 'java.lang.StringBuffer')
    Import(w, 'java.lang.reflect.InvocationTargetException')
    w.head()
    Import(w, 'java.math.BigDecimal')
    w.head()
    Import(w, 'java.text.SimpleDateFormat')
    w.head()
    Import(w, 'java.util.ArrayList')
    Import(w, 'java.util.Arrays')
    Import(w, 'java.util.Comparator')
    Import(w, 'java.util.Collection')
    Import(w, 'java.util.Collections')
    Import(w, 'java.util.HashMap')
    Import(w, 'java.util.HashSet')
    Import(w, 'java.util.Iterator')
    Import(w, 'java.util.NoSuchElementException')
    Import(w, 'java.util.UUID')
    w.head()
    Import(w, 'org.apache.commons.lang3.StringUtils')
    w.head()
    Import(w, 'org.joda.time.DateTime')
    Import(w, 'org.joda.time.LocalDate')
    w.head()


@contextlib.contextmanager
def Class(w, name,
          write_file=None,
          final=False,
          implements=None,
          type_='class',
          visibility='public',
          extends=None,
          static=False,
          abstract=False,
          package=None):
    if name.startswith('\n'):
        w()
        name = name[1:]

    CLASS_STACK.append(name)

    w(u'{}{}{}{}{} {}{}{} {{',
      visibility + ' ' if visibility else '',
      'static ' if static else '',
      'abstract ' if abstract else '',
      'final ' if final else '',
      type_,
      name,
      ' extends ' + extends if extends else '',
      ' implements ' + ', '.join(implements) if implements else '')

    with w.indent():
        yield

    w(u'}}')
    CLASS_STACK.pop()
    assert bool(CLASS_STACK) ^ bool(package), (CLASS_STACK, package)
    if not CLASS_STACK:
        assert write_file
        write_file.java(w, name, package)


@contextlib.contextmanager
def Method(w, visibility, return_type, name, params=(), throws=(),
           override=False, static=False):
    if visibility.startswith('\n'):
        w()
        visibility = visibility[1:]

    if override:
        w('@Override')

    w(u'{}{}{}{}(\v{}){} {{',
      visibility and visibility + u' ',
      u'static ' if static else '',
      return_type and return_type + u' ',
      name,
      u', \v'.join(u'{} {}'.format(t, p) for (t, p) in params),
      (u'\n        throws ' + ', '.join(ex for ex in throws)
       if throws else ''))

    with w.indent():
        yield

    w(u'}}')


def Ctor(w, visibility, name, params=(), throws=()):
    assert name == CLASS_STACK[-1]
    return Method(w, visibility, '', name, params, throws)


def SeparatorComment(w):
    w(u'\n// ' + '-' * (68 - w.indent_level))


def Package(w, name):
    w.head(u'package {};\n', name)


@contextlib.contextmanager
def EnumeratorMethod(w, visibility, out_tname, name='enumerator',
                     params=(), throws=(), override=True, source=None):
    if source:
        (src_var, src_expr, src_tname) = source
    else:
        (src_var, src_expr, src_tname) = (None, None, None)

    with Method(w, visibility, 'io.sysl.Enumerator<' + out_tname + '>',
                name, params, throws, override):
        if src_var:
            w('final io.sysl.Enumerator<{}> {} = \v{}.enumerator();',
              src_tname, src_var, src_expr)

        @contextlib.contextmanager
        def enumerator():
            w('return new io.sysl.Enumerator<{}>() {{', out_tname)

            called = {'moveNext': False, 'current': False}

            with w.indent():
                @contextlib.contextmanager
                def moveNext():
                    with Method(w, 'public', 'boolean', 'moveNext',
                                override=True):
                        yield
                    called['moveNext'] = True

                @contextlib.contextmanager
                def current():
                    with Method(w, '\npublic', out_tname, 'current',
                                override=True):
                        yield
                    called['current'] = True

                yield (moveNext, current)

                if not called['moveNext']:
                    assert src_var
                    with moveNext():
                        w('return {}.moveNext();', src_var)

                if not called['current']:
                    assert src_var
                    with current():
                        w('return {}.current();', src_var)

            w('}};')

        yield enumerator


@contextlib.contextmanager
def ViewMethod(w, tname, visibility, out_tname, name, params=(), throws=(),
               override=False):
    view_tname = 'View' if tname == out_tname else out_tname + '.View'

    with Method(w, visibility, view_tname, name, params, throws,
                override):
        @contextlib.contextmanager
        def view():
            w('return new {}(model) {{', view_tname)
            with w.indent():
                @contextlib.contextmanager
                def enumerator_method(src_var=None):
                    with EnumeratorMethod(w, 'public', out_tname,
                                          source=src_var and (src_var, 'View.this', tname)) as enumerator:
                        yield enumerator

                yield enumerator_method

        yield view

        w('}};')


def codeForType(t, scope):
    if isinstance(t, sysl_pb2.Application):
        assert len(t.name.part) == 1
        return t.name.part[0]

    which_type = t.WhichOneof('type')

    if which_type == 'primitive':
        return TYPE_MAP[t.primitive]

    elif which_type == 'type_ref':
        ref = t.type_ref.ref

        if len(ref.path) > 1:
            raise NotImplementedError('Nested types')

        if ref.HasField('appname'):
            app = syslx.AppByName(scope.module).get(ref.appname)
            assert app, 'Appname {} not found'.format(
                '.'.join(ref.appname.part))
            package = syslx.View(app.attrs)['package'].s
            assert package, ref
            is_anon = ref.path and re.match(
                r'AnonType_\d+__$', '.'.join(ref.path))
            return '{}{}{}'.format(
                package,
                ('' if ref.path and not is_anon else
                 ''.join('.' + p for p in ref.appname.part)),
                ''.join('.' + p for p in ref.path))

        info = syslx.TypeInfoByRef(scope.module)[t.type_ref]
        return info.app.attrs['package'].s + '.' + ref.path[0]

    if which_type == 'set':
        if t.set.WhichOneof('type') == 'tuple':
            return 'io.sysl.Enumerable<' + codeForType(t.set, scope) + '>'
        return codeForType(t.set, scope) + '.View'

    if which_type == 'tuple':
        assert not t.tuple.attr_defs, 'Only empty tuple types supported for now'
        return 'EmptyTuple'

    elif which_type == 'no_type':
        return 'void'

    elif which_type is None:
        raise RuntimeError("codeForType(): Type not set")

    raise NotImplementedError("Unable to generate code for " + which_type)


def literalToJava(literal):
    def f(code, primitive):
        return (code, sysl_pb2.Type(primitive=primitive))

    which_lit = literal.WhichOneof('value')
    if which_lit == 'b':
        return f('true' if literal.b else 'false', sysl_pb2.Type.BOOL)
    if which_lit == 'i':
        return f(str(literal.i), sysl_pb2.Type.INT)
    if which_lit == 'd':
        return f(str(literal.d), sysl_pb2.Type.FLOAT)
    if which_lit == 's':
        # TODO: Ensure correct Java syntax
        return f(json.dumps(literal.s), sysl_pb2.Type.STRING)
    if which_lit == 'decimal':
        return f('new BigDecimal("{}")'.format(literal.decimal),
                 sysl_pb2.Type.DECIMAL)
    if which_lit == 'uuid':
        return f(
            'new UUID({}L, {}L)'.format(*struct.unpack('>qq', literal.uuid)),
            sysl_pb2.Type.UUID)
    if which_lit == 'null':
        return f('null', sysl_pb2.Type.EMPTY)
    raise RuntimeError(
        'Cannot convert literal.{} to Java'.format(which_lit))


lastNewvar = [0]


def newvar():
    lastNewvar[0] += 1
    return 'tmp_' + str(lastNewvar[0]) + '__'


class NoExpr(object):
    def __str__(self):
        raise NotImplementedError("Cannot stringify non-expr")


def types_match(a, b):
    which_a = a.WhichOneof('type')
    which_b = b.WhichOneof('type')
    return (
        getattr(a, which_a) == getattr(b, which_b) if which_a and which_b else
        which_a == which_b)


def oneline(msg):
    return '\033[1;30m{{\033[0m{}\033[1;30m}}\033[0m'.format(
        re.sub(r'\s+', ' ', str(msg).strip()))


def coerce_to_bool(expr):
    if expr.WhichOneof('expr') == 'call' and expr.call.func == 'bool':
        return expr
    return sysl_pb2.Expr(call=sysl_pb2.Expr.Call(func='bool', arg=[expr]))


def codeForExpr(w_, expr, scope, module, let=None):
    W = [w_]
    scope_stk = [scope]
    let_stk = []

    def cur_scope():
        return scope_stk[-1]

    @contextlib.contextmanager
    def new_scope(dot_=None, **kwargs):
        inner_scope = scopes.Scope(cur_scope(), dot_, **kwargs)
        scope_stk.append(inner_scope)
        yield inner_scope
        scope_stk.pop()

    @contextlib.contextmanager
    def let_scope(let):
        if let:
            let_stk.append(let)
        yield
        if let:
            let_stk.pop()

    @contextlib.contextmanager
    def inner_java_scope():
        outer_w = W[0]
        W[0] = writer.Writer()
        yield W[0]
        W[0] = outer_w

    def tmp(var_type):
        var = let_stk[-1] if let_stk else newvar()
        W[0]('{} {};', var_type, var)
        return var

    def tmpForExpr(expr):
        (expr_code, expr_type) = E(-1, expr)
        assert expr_type, (expr, expr_code)

        if expr.WhichOneof('expr') == 'name':
            return (expr_code, expr_code, expr_type)

        expr_var = tmp(codeForType(expr_type, cur_scope()))
        return ('({} = {})'.format(expr_var, expr_code), expr_var, expr_type)

    def tmpForType(expr_type):
        var = tmp(codeForType(expr_type, cur_scope()))
        return var

    def E(parent_pc, expr, let=None):
        parent_pc = PRECEDENCE[parent_pc]

        def protect(pc, expr):
            pc = PRECEDENCE[pc]
            return '(' + expr + ')' if pc < parent_pc else expr

        which_expr = expr.WhichOneof('expr')

        if which_expr == 'name':
            if expr.name:
                return (
                    ((cur_scope().get('__dot__') or 'item') if expr.name == '.' else
                     safe(expr.name)),
                    cur_scope()[expr.name])
            raise RuntimeError('Missing expr.name')

        elif which_expr == 'literal':
            return literalToJava(expr.literal)

        # elif which_expr == 'transform':

        elif which_expr == 'ifelse':
            return ifelseExprToJava(protect, expr.ifelse)

        elif which_expr == 'unexpr':
            unexpr = expr.unexpr
            pc = PRECEDENCE[(1, unexpr.op)]
            op = UNOP_JAVA[unexpr.op]
            is_not = unexpr.op == unexpr.NOT
            if isinstance(op, basestring):
                (arg, arg_type) = E(
                    pc, coerce_to_bool(unexpr.arg) if is_not else unexpr.arg)
                return (
                    protect(pc, re.sub(r'^(!!)+', '', '{}{}'.format(op, arg))),
                    arg_type)
            return op(protect, unexpr.op, unexpr.arg)

        elif which_expr == 'tuple':
            tup = expr.tuple
            assert not tup.attrs, 'Only empty tuple supported for now'
            return ('EmptyTuple.theOne', expr.type)

        elif which_expr == 'binexpr':
            binexpr = expr.binexpr
            pc = PRECEDENCE[(2, binexpr.op)]
            op = BINOP_JAVA[binexpr.op]
            if op is None:
                raise RuntimeError(
                    '{} operator to Java not supported'.format(
                        sysl_pb2.Expr.BinExpr.Op.Name(binexpr.op)))

            if isinstance(op, str):
                if binexpr.scopevar:
                    with new_scope(binexpr.scopevar):
                        (lhs, lhs_type) = E(pc, binexpr.lhs)
                        (rhs, rhs_type) = E(pc, binexpr.rhs)
                else:
                    (lhs, lhs_type) = E(pc, binexpr.lhs)
                    (rhs, rhs_type) = E(pc, binexpr.rhs)
                return (protect(pc, '{} {} {}'.format(lhs, op, rhs)), lhs_type)

            return op(protect, binexpr)

        elif which_expr == 'relexpr':
            relexpr = expr.relexpr
            pc = PRECEDENCE[(2, relexpr.op)]
            op = RELOP_JAVA[relexpr.op]
            if op is None:
                raise RuntimeError(
                    '{} operator to Java not supported'.format(
                        sysl_pb2.Expr.BinExpr.Op.Name(relexpr.op)))

            return op(protect, expr)

        elif which_expr in ['get_attr', 'navigate']:
            is_get_attr = which_expr == 'get_attr'
            is_navigate = which_expr == 'navigate'

            ga = getattr(expr, which_expr)

            if_false_code = 'null'

            if ga.nullsafe:
                (arg_init_code, arg_code, arg_type_) = tmpForExpr(ga.arg)
                pc = '?:'
                code = '{} == null ? {} : {}.{}{}()'.format(
                    arg_init_code, CODE_MARKER2, arg_code, GET_DELIM, CODE_MARKER)
            else:
                pc = '.'
                (arg, arg_type_) = E(pc, ga.arg)
                code = '{}.{}{}()'.format(arg, GET_DELIM, CODE_MARKER)

            (app, arg_type) = cur_scope().resolve(arg_type_)

            suffix = ''

            if isinstance(arg_type, sysl_pb2.Application):
                assert is_get_attr
                arg_setof = False
                if not ga.setof:
                    raise RuntimeError('<{0}>.{1} must be <{0}>.table of {1}'.format(
                        diagutil.fmt_app_name(arg_type.name), ga.attr))
                attr_type = sysl_pb2.Type()
                ref = attr_type.type_ref.ref
                ref.appname.CopyFrom(arg_type.name)
                ref.path.append(ga.attr)
                expr_type = sysl_pb2.Type(set=attr_type)
                code = code.replace(GET_DELIM, 'get')
                suffix = 'Table'
            elif isinstance(arg_type, sysl_pb2.Type):
                which_type = arg_type.WhichOneof('type')
                arg_setof = which_type == 'set'
                if arg_setof:
                    arg_type = arg_type.set
                    which_type = arg_type.WhichOneof('type')

                if is_get_attr:
                    if which_type == 'relation':
                        entity_type = arg_type.relation
                    elif which_type == 'tuple':
                        entity_type = arg_type.tuple
                    else:
                        raise RuntimeError('syntax error', which_type)
                    expr_type = entity_type.attr_defs.get(ga.attr)
                    assert expr_type, (
                        '\033[1;31mField \033[37m{}.{}\033[31m not found\033[0m: {}'.format(
                            ga.arg.name, ga.attr, ga))
                    if expr_type.WhichOneof('type') == 'type_ref':
                        (_, expr_type, _, _, _, _) = (
                            syslx.TypeInfoByRef(scope.module)[expr_type.type_ref])
                    code = code.replace(GET_DELIM, 'get')
                elif ga.attr.startswith('.'):
                    expr_type = sysl_pb2.Type()
                    ref = expr_type.type_ref.ref
                    ref.appname.CopyFrom(app.name)
                    ref.path.append(ga.attr)
                    if ga.setof:
                        expr_type = sysl_pb2.Type(set=expr_type)
                        suffix = 'View'
                    code = code.replace(GET_DELIM, 'get')
                else:
                    expr_type = sysl_pb2.Type()
                    ref = expr_type.type_ref.ref
                    ref.appname.CopyFrom(app.name)
                    ref.path.append(ga.attr)
                    if ga.setof:
                        if_false_code = 'new {}.Set()'.format(
                            codeForType(expr_type, cur_scope()))
                        expr_type = sysl_pb2.Type(set=expr_type)
                        suffix = 'View'
                    code = code.replace(GET_DELIM, 'to')
            else:
                raise RuntimeError('syntax error')

            via = 'Via' + CamelCase(ga.via) if is_navigate and ga.via else ''

            return (
                protect(pc,
                        code.replace(
                            CODE_MARKER,
                            CamelCase(ga.attr.lstrip('.')) + suffix + via
                        ).replace(
                            CODE_MARKER2,
                            if_false_code
                        )),
                expr_type)

        elif which_expr == 'call':
            call = expr.call
            pc = PRECEDENCE['.']
            func = call.func
            args_and_types = [E(pc, arg) for arg in call.arg]
            args = [at[0] for at in args_and_types]
            if func == 'bool':
                [(arg, arg_type)] = args_and_types
                if arg_type.primitive == arg_type.BOOL:
                    return (protect(pc, arg), arg_type)
                if arg_type.WhichOneof('type') == 'set':
                    return (protect((1, 'NOT'), '!' +
                                    arg + '.isEmpty()'), BOOL_TYPE)
            elif func == 'str':
                assert len(args) == 1
                return (
                    protect(pc, 'String.valueOf(' + args[0] + ')'),
                    sysl_pb2.Type(primitive=sysl_pb2.Type.STRING))
            elif func == 'int':
                assert len(args) in [1, 2]
                if len(args) == 2:
                    return (protect(pc, 'io.sysl.ExprHelper.toInteger({}, {})'
                                    .format(*args)),
                            sysl_pb2.Type(primitive=sysl_pb2.Type.INT))
                return (protect(pc, 'Integer.parseInt({})'.format(args[0])),
                        sysl_pb2.Type(primitive=sysl_pb2.Type.INT))
            elif func == 'formatDate':
                assert len(args) == 2
                return (
                    protect(
                        pc, 'DateTimeFormat.forPattern({1}).print(new DateTime({0}))'.format(*args)),
                    sysl_pb2.Type(primitive=sysl_pb2.Type.STRING))

            fmt = '{1}{0}({2})' if call.func.startswith('.') else '{0}({1}{3})'
            type_ = {
                '.any': args_and_types[0][1] if args_and_types else None,
                '.count': sysl_pb2.Type(primitive=sysl_pb2.Type.INT),
                'autoinc': sysl_pb2.Type(primitive=sysl_pb2.Type.INT),
                'concat': sysl_pb2.Type(primitive=sysl_pb2.Type.STRING),
                'lstrip': sysl_pb2.Type(primitive=sysl_pb2.Type.STRING),
                'log': sysl_pb2.Type(primitive=sysl_pb2.Type.BOOL),
                'now': sysl_pb2.Type(primitive=sysl_pb2.Type.DATETIME),
                'regsub': sysl_pb2.Type(primitive=sysl_pb2.Type.STRING),
                'rstrip': sysl_pb2.Type(primitive=sysl_pb2.Type.STRING),
                'strip': sysl_pb2.Type(primitive=sysl_pb2.Type.STRING),
                'substr': sysl_pb2.Type(primitive=sysl_pb2.Type.STRING),
                'today': sysl_pb2.Type(primitive=sysl_pb2.Type.DATE),
                'to_date': sysl_pb2.Type(primitive=sysl_pb2.Type.DATE),
            }.get(call.func)

            if type_ is None:
                type_ = scope[call.func]

            return (
                protect(pc,
                        fmt.format(
                            call.func,
                            args[0] if args else '',
                            ', '.join(args[1:]),
                            ''.join(', ' + a for a in args[1:]))),
                type_)

        elif which_expr in ['list', 'set']:
            exprs = getattr(expr, which_expr).expr
            (elems, types) = zip(*[E(-1, e)
                                   for e in exprs]) if exprs else [(), ()]
            type_codes = {codeForType(t, cur_scope()) for t in types}
            [type_code] = type_codes if len(type_codes) == 1 else ['Object']
            (pc, fmt) = (
                (-1, 'new {}[]{{{}}}') if which_expr != 'set' else
                ('.', 'new io.sysl.EnumerableSet(Arrays.asList({1}))'))
            expr_type = sysl_pb2.Type()
            if len(type_codes) == 1:
                expr_type.set.CopyFrom(next(iter(types)))
            else:
                expr_type.set.primitive = sysl_pb2.Type.ANY
            return (protect(pc, fmt.format(
                type_code, ', \v'.join(elems))), expr_type)

        elif which_expr == 'transform':
            return transformToJava(protect, expr)

        elif which_expr is None:
            return (NoExpr(), None)

        else:
            raise RuntimeError(
                'expr.{} to Java not supported'.format(which_expr))

    def ifelseExprToJava(protect, ifelse):
        def commonLhsEqLiteral(ifelse):
            """Return True if ifelse chain has same expr == same literal-type"""
            cond = ifelse.cond
            if not ifelse.nullsafe and cond.WhichOneof('expr') == 'binexpr':
                binexpr = cond.binexpr
                if (binexpr.op == sysl_pb2.Expr.BinExpr.EQ and
                    binexpr.rhs.WhichOneof('expr') == 'literal' and
                        binexpr.rhs.literal.WhichOneof('value') == 's'):
                    lel = (
                        binexpr.lhs,
                        binexpr.rhs.literal.WhichOneof('value'))
                    if (ifelse.if_false.WhichOneof('expr') != 'ifelse' or
                            lel == commonLhsEqLiteral(ifelse.if_false.ifelse)):
                        return lel
            return (None, None)

        (cond, cond_type) = E('?:', coerce_to_bool(ifelse.cond))
        (if_true, if_true_type) = E('?:', ifelse.if_true)

        (clhs, _) = commonLhsEqLiteral(ifelse)
        if clhs:
            var = tmpForType(if_true_type)

            (control_init_code, control_code, _) = tmpForExpr(clhs)

            cases = []
            while True:
                rhs = E(-1, ifelse.cond.binexpr.rhs)[0]
                if_true = E(-1, ifelse.if_true)[0]
                cases.append((if_true, rhs))
                if ifelse.if_false.WhichOneof('expr') != 'ifelse':
                    break
                ifelse = ifelse.if_false.ifelse
            default = E(-1, ifelse.if_false)[0]

            with Switch(W[0],
                        '{} == null ? "47989B2E-7895-4063-B5F3-AC2076A43745" : {}',
                        control_init_code, control_code):
                gg = [
                    (rhs, list(g))
                    for (rhs, g) in itertools.groupby(sorted(cases), lambda t: t[0])]
                for (i, (rhs, g)) in enumerate(gg):
                    for (_, cond) in g:
                        W[0]('case {}:', cond)
                    with W[0].indent():
                        if i < len(gg) - 1 or rhs != default:
                            W[0]('{} = \v{};', var, rhs)
                            W[0]('break;')
                with Default(W[0]):
                    W[0]('{} = \v{};', var, default)

            return (var, if_true_type)

        (if_false, if_false_type) = E('?:', ifelse.if_false)

        # TODO: edge-cases other than null?
        if if_true_type.primitive == if_true_type.EMPTY:
            result_type = if_false_type
        else:
            result_type = if_true_type

        if ifelse.nullsafe:
            (cond_init_code, cond_code, _) = tmpForExpr(
                coerce_to_bool(ifelse.cond))
            return (
                protect('?:', '{} == null ? null : \v{} ? \v{} : \v{}'.format(
                    cond_init_code, cond_code, if_true, if_false)),
                result_type)

        return (
            protect('?:', '{} \v? {} \v: \v{}'.format(
                cond, if_true, if_false)),
            result_type)

    def whereToJava(protect, binexpr):
        pc = PRECEDENCE['.']
        (arg, arg_type) = E('.', binexpr.lhs)
        assert arg_type.WhichOneof(
            'type') == 'set', arg_type.WhichOneof('type')
        arg_name = 'dot_' if binexpr.scopevar == '.' else binexpr.scopevar
        with new_scope(**{arg_name: arg_type.set}):
            with inner_java_scope() as inner_w:
                (pred, pred_type) = E(pc, coerce_to_bool(binexpr.rhs))
                which_pred_type = pred_type.WhichOneof('type')
                is_primitive = which_pred_type == 'primitive'
#        if which_pred_type == 'set':

                if (not is_primitive or pred_type.primitive != sysl_pb2.Type.BOOL):
                    pass
                    # print >>sys.stderr, (
                    #   'WARNING: "where" predicate must be bool, not {}: {}'.format(
                    #     (sysl_pb2.Type.Primitive.Name(pred_type.primitive).lower()
                    #      if is_primitive else which_pred_type),
                    #     oneline(pred_type)))
                if codeForType(arg_type, cur_scope()).endswith('.View'):
                    fmt = (
                        '{arg}.where(\n'
                        '    new io.sysl.Expr<Boolean, {type}>() {{\n'
                        '        @Override\n'
                        '        public Boolean evaluate(final {type} {scopevar}) {{\n'
                        '            {extra_code}return {pred};\n'
                        '        }}\n'
                        '    }})')
                else:
                    fmt = (
                        'io.sysl.Enumeration.where(\v{arg},\n'
                        '    new io.sysl.Expr<Boolean, {type}>() {{\n'
                        '        @Override\n'
                        '        public Boolean evaluate(final {type} {scopevar}) {{\n'
                        '            {extra_code}return {pred};\n'
                        '        }}\n'
                        '    }})')

                code = fmt.format(
                    arg=arg,
                    type=codeForType(arg_type.set, cur_scope()),
                    scopevar='item' if binexpr.scopevar == '.' else binexpr.scopevar,
                    pred=pred,
                    extra_code=CODE_MARKER)
                code = code.replace(
                    CODE_MARKER, str(inner_w).replace('\n', '\n            '))
            return (protect(pc, code), arg_type)

    def flattenToJava(protect, binexpr):
        pc = PRECEDENCE['.']
        (target, target_type) = E(-1, binexpr.lhs)
        arg_name = 'dot_' if binexpr.scopevar == '.' else binexpr.scopevar
        with new_scope(**{arg_name: target_type}):
            (item, item_type) = E(pc, binexpr.rhs)
            with inner_java_scope() as inner_w:
                code = (
                    '{item_type}.View.view(\vnull, io.sysl.Enumeration.flatten(\n'
                    '    {target},\n'
                    '    new io.sysl.Expr<\n'
                    '            io.sysl.Enumerable<{item_type}>,\n'
                    '            {type}>() {{\n'
                    '        @Override\n'
                    '        public {item_type}.View evaluate(final {type} {scopevar}) {{\n'
                    '            {extra_code}return {item};\n'
                    '        }}\n'
                    '    }}))'.format(
                        target=target,
                        type=codeForType(target_type.set, cur_scope()),
                        item_type=codeForType(item_type.set, cur_scope()),
                        scopevar='item' if binexpr.scopevar == '.' else binexpr.scopevar,
                        item=item,
                        extra_code=CODE_MARKER))
                code = code.replace(
                    CODE_MARKER, str(inner_w).replace('\n', '\n            '))
            return (protect(pc, code), target_type)

    def toMatchingToJava(positive):
        def toJava(protect, binexpr):
            # Rewrite "a ~> b" as "b where(b: a where(a: a.fld1 == b.fld1 ...))", where
            # fld1... are fields in common between a and b.
            rewrite = sysl_pb2.Expr()

            where = rewrite.binexpr
            where.op = where.WHERE
            where.lhs.CopyFrom(binexpr.rhs)
            where.scopevar = 'b'

            if positive:
                pred = where.rhs.binexpr
            else:
                where.rhs.unexpr.op = where.rhs.unexpr.NOT
                pred = where.rhs.unexpr.arg.binexpr
            pred.op = pred.WHERE
            pred.lhs.CopyFrom(binexpr.lhs)
            pred.scopevar = 'a'

            test = None

            def attr_names(arg):
                # We just want type info, so we'll use inner scope to trap any
                # codegen.
                with inner_java_scope():
                    (_, arg_type) = E('.', arg)
                    (_, resolved_arg_type) = cur_scope().resolve(arg_type)
                    if resolved_arg_type.WhichOneof('type') == 'set':
                        resolved_arg_type = resolved_arg_type.set
                    which_type = resolved_arg_type.WhichOneof('type')
                    assert which_type in ['relation', 'tuple'], which_type
                    attr_defs = getattr(resolved_arg_type,
                                        which_type).attr_defs
                    return set(attr_defs.keys())

            lhs_attr_names = attr_names(binexpr.lhs)
            rhs_attr_names = attr_names(binexpr.rhs)
            available_attr_names = lhs_attr_names & rhs_attr_names
            attr_names = set(binexpr.attr_name)
            if attr_names == {'*'}:
                attr_names = available_attr_names
            else:
                assert attr_names <= available_attr_names, u'{{{}}} ⊈ {{{}}}'.format(
                    ','.join(attr_names - available_attr_names),
                    ','.join(available_attr_names)).encode('utf-8')

            for fname in attr_names:
                if test:
                    test = sysl_pb2.Expr(
                        binexpr=sysl_pb2.Expr.BinExpr(
                            op=sysl_pb2.Expr.BinExpr.AND,
                            rhs=test))
                    comp = test.binexpr.lhs.binexpr
                else:
                    test = sysl_pb2.Expr()
                    comp = test.binexpr

                comp.op = comp.EQ
                lhs = comp.lhs.get_attr
                lhs.arg.name = 'a'
                lhs.attr = fname
                rhs = comp.rhs.get_attr
                rhs.arg.name = 'b'
                rhs.attr = fname

            pred.rhs.CopyFrom(test)

            return E('.', rewrite)

        return toJava

    SINGLE_OPS = [
        sysl_pb2.Expr.UnExpr.SINGLE,
        sysl_pb2.Expr.UnExpr.SINGLE_OR_NULL,
    ]

    def unexprToJava(prec, fmt):
        def b2j(protect, op, a):
            (arg, arg_type) = E(PRECEDENCE['.'], a)
            return (
                protect(prec or PRECEDENCE[op], fmt.format(arg)),
                arg_type.set if op in SINGLE_OPS else arg_type,
            )
        return b2j

    def binexprToJava(prec, fmt, out_type=None):
        prec = prec and PRECEDENCE[prec]

        def b2j(protect, binexpr):
            (lhs, lhs_type) = E(PRECEDENCE['.'], binexpr.lhs)
            (rhs, rhs_type) = E(PRECEDENCE['.'], binexpr.rhs)
            return (
                protect(
                    prec or PRECEDENCE[(2, binexpr.op)], fmt.format(lhs, rhs)),
                out_type or lhs_type)
        return b2j

    def aggexprToJava(protect, expr):
        relexpr = expr.relexpr

        pc = PRECEDENCE['.']
        (target, target_type) = E('.', relexpr.target)
        arg_name = 'dot_' if relexpr.scopevar == '.' else relexpr.scopevar
        with new_scope(**{arg_name: target_type}):
            assert len(relexpr.arg) == 1, 'Unsupported relexpr arg list'
            with inner_java_scope() as inner_w:
                (item, item_type) = E(-1, relexpr.arg[0])
            which_item = item_type.WhichOneof('type')
            is_primitive = which_item == 'primitive'
            agg = sysl_pb2.Expr.RelExpr.Op.Name(relexpr.op).lower()
            item_type_code = codeForType(item_type, cur_scope())
            item_type_suffix = '' if agg in ['min', 'max'] else item_type_code
            code = (
                '{target}.{agg}{item_type_suffix}(new io.sysl.Expr<{item_type_code}, {type}>() {{\n'
                '        @Override\n'
                '        public {item_type_code} evaluate(final {type} {scopevar}) {{\n'
                '            {extra_code}return {item};\n'
                '        }}\n'
                '    }})'.format(
                    target=target,
                    agg=agg,
                    type=codeForType(target_type.set, cur_scope()),
                    item_type_code=item_type_code,
                    item_type_suffix=item_type_suffix,
                    scopevar='item' if relexpr.scopevar == '.' else relexpr.scopevar,
                    item=item,
                    extra_code=str(inner_w).replace('\n', '\n            ')))
            return (protect(pc, code), item_type)

    def rankToJava(protect, expr):
        relexpr = expr.relexpr

        (in_code, in_type) = E(-1, relexpr.target)
        scope_name = 'dot_' if relexpr.scopevar == '.' else relexpr.scopevar
        scopevar = 'item' if relexpr.scopevar == '.' else relexpr.scopevar
        with new_scope(**{scope_name: in_type}):
            compare_code = ''
            for (arg, desc) in zip(relexpr.arg, relexpr.descending):
                with inner_java_scope() as inner_w:
                    with new_scope(__dot__='b' if desc else 'a'):
                        (a_item, _) = E('.', arg)
                    with new_scope(__dot__='a' if desc else 'b'):
                        (b_item, _) = E('.', arg)
                compare_code += str(inner_w)
                compare_code += 'if ((i = {}.compareTo({})) != 0) return i;\n'.format(
                    a_item, b_item)

            out_type = expr.type
            out_attrs = {
                attr_name
                for (attr_name, _) in datamodel.sorted_fields(
                    scope.resolve(out_type.set)[1])} - {relexpr.attr_name[0]}
            assign_code = ''
            (tuple_app, tuple_type) = scope.resolve(in_type.set)
            for (attr_name, attr) in datamodel.sorted_fields(tuple_type):
                if attr_name in out_attrs:
                    assign_code += 'out.set{0}({1}.get{0}());\n'.format(
                        CamelCase(name(attr_name)), scopevar)
            assign_code += 'out.set{}(r);\n'.format(
                CamelCase(name(relexpr.attr_name[0])))

            code = (
                '{out_type}.View.view(null, \vio.sysl.Enumeration.rank(\v{in_code}, \vnew io.sysl.Enumeration.Ranker<\v{type}, {out_type}>() {{\n'
                '        @Override\n'
                '        public int compare(final {type} a, final {type} b) {{\n'
                '            int i;\n'
                '            {compare_code}return 0;\n'
                '        }}\n'
                '        @Override\n'
                '        public {out_type} ranked(final {type} {scopevar}, final int r) {{\n'
                '            {out_type} out =\n'
                '                new {out_type}();\n'
                '            {assign_code}return out;\n'
                '        }}\n'
                '    }}))'.format(
                    in_code=in_code,
                    agg=sysl_pb2.Expr.RelExpr.Op.Name(relexpr.op).lower(),
                    type=codeForType(in_type.set, cur_scope()),
                    out_type=codeForType(out_type.set, cur_scope()),
                    scopevar=scopevar,
                    compare_code=compare_code.replace('\n', '\n' + '    ' * 3),
                    assign_code=assign_code.replace('\n', '\n' + '    ' * 3),
                ))
            return (protect('.', code), expr.type)

    def firstByToJava(protect, expr):
        relexpr = expr.relexpr

        pc = PRECEDENCE['.']
        (in_code, in_type) = E(-1, relexpr.target)
        scope_name = 'dot_' if relexpr.scopevar == '.' else relexpr.scopevar
        scopevar = 'item' if relexpr.scopevar == '.' else relexpr.scopevar
        (n_code, _) = E(pc, relexpr.arg[0])
        with new_scope(**{scope_name: in_type}):
            compare_code = ''
            for (arg, desc) in zip(relexpr.arg[1:], relexpr.descending):
                with inner_java_scope() as inner_w:
                    with new_scope(__dot__='b' if desc else 'a'):
                        (a_item, _) = E(pc, arg)
                    with new_scope(__dot__='a' if desc else 'b'):
                        (b_item, _) = E(pc, arg)
                compare_code += str(inner_w)
                compare_code += 'if ((i = {}.compareTo({})) != 0) return i;\n'.format(
                    a_item, b_item)

            code = (
                '{type}.View.view(null, \vio.sysl.Enumeration.first(\v'
                '{in_code}, {n_code}, \vnew java.util.Comparator<'
                '\v{type}>() {{\n'
                '        @Override\n'
                '        public int compare(final {type} a, final {type} b) {{\n'
                '            int i;\n'
                '            {compare_code}return 0;\n'
                '        }}\n'
                '    }}))'.format(
                    in_code=in_code,
                    n_code=n_code,
                    agg=sysl_pb2.Expr.RelExpr.Op.Name(relexpr.op).lower(),
                    type=codeForType(in_type.set, cur_scope()),
                    scopevar=scopevar,
                    compare_code=compare_code.replace('\n', '\n' + '    ' * 3),
                ))
            return (protect(pc, code), in_type)

    def snapshotToJava(protect, expr):
        relexpr = expr.relexpr

        pc = PRECEDENCE['.']
        (target, target_type) = E(pc, relexpr.target)
        return (protect(pc, '{}.snapshot()'.format(target)), target_type)

    def transformToJava(protect, expr):
        transform = expr.transform
        (target, target_type) = E(-1, transform.arg)

        (ret_app, ret_type) = cur_scope().resolve(expr.type)
        ret_type_code = codeForType(expr.type, cur_scope())

        is_model = isinstance(ret_type, sysl_pb2.Application)
        is_relation = (
            not is_model and
            (ret_type.WhichOneof('type') == 'relation' or
             (ret_type.WhichOneof('type') == 'set' and
                ret_type.set.WhichOneof('type') == 'relation')))

        is_model = isinstance(ret_type, sysl_pb2.Application)
        is_view = not is_model and ret_type.WhichOneof('type') == 'set'
        out_type = expr.type.set if is_view else expr.type
        out_type_code = codeForType(out_type, scope)
        with new_scope(out=out_type):
            (_, resolved_out_type) = scope.resolve(out_type)
            pkey = ([] if is_model else
                    datamodel.primary_key_params(resolved_out_type, module))
            pkey_fields = {jfname for (_, _, jfname) in pkey}

            fkeys = (set() if is_model else {
                name(fname): type_info
                for (fname, _, type_info) in datamodel.foreign_keys(
                    resolved_out_type, module)
            })
            fkey_fields = set(fkeys)

            special_fields = pkey_fields | fkey_fields

            scopevar = 'dot_' if transform.scopevar == '.' else transform.scopevar
            with new_scope(**{scopevar: target_type}) as transform_scope:
                with inner_java_scope() as w:
                    for stmt in transform.stmt:
                        which_stmt = stmt.WhichOneof('stmt')

                        if which_stmt == 'assign':
                            assign = stmt.assign
                            if assign.name in special_fields:
                                set_prefix = '_PRIVATE_'
                            else:
                                set_prefix = ''

                            if assign.expr.WhichOneof('expr'):
                                (code, _) = codeForExpr(
                                    w, assign.expr, transform_scope, module)

                                w('out.{}set{}\037 ({});',
                                  set_prefix, CamelCase(assign.name), code)
                            else:
                                w('// out.{}set{}(...);',
                                  set_prefix, CamelCase(assign.name))

                        elif which_stmt == 'let':
                            let = stmt.let
                            if let.expr.WhichOneof('expr'):
                                (code, type_) = codeForExpr(
                                    w, let.expr, transform_scope, module, let=let.name)
                                transform_scope[let.name] = type_
                                if code:
                                    w('final {} {} = \v{};',
                                      codeForType(type_, transform_scope),
                                      let.name,
                                      code)
                            else:
                                w('// <type> {} = ...;', code)

                        elif which_stmt == 'inject':
                            w('{};', codeForExpr(w, stmt.inject,
                                                 transform_scope, module)[0])

                        else:
                            raise RuntimeError(
                                'Unexpected stmt type: ' + which_stmt)

                if is_view:
                    is_table = scope.resolve(
                        out_type)[1].WhichOneof('type') == 'relation'
                    code = (
                        '{out_type}.View.view(null, \n'
                        '    {target}.map(new io.sysl.Expr<{out_type}, {type}>() {{\n'
                        '        @Override\n'
                        '        public {out_type} evaluate(final {type} {scopevar}) {{\n'
                        '            {out_type} out = {ctor};\n'
                        '            {extra_code}return out;\n'
                        '        }}\n'
                        '    }}))'.format(
                            target=target,
                            type=codeForType(target_type.set, cur_scope()),
                            out_type=out_type_code,
                            scopevar=(
                                'item' if transform.scopevar == '.' else transform.scopevar),
                            ctor=(
                                ('{}._PRIVATE_new()' if is_table else 'new {}()').format(
                                    out_type_code)),
                            extra_code=str(w).replace('\n', '\n            ')))
                else:
                    # TODO: Figure out why target_type is a set.
                    if target_type.WhichOneof('type') == 'set':
                        target_type = target_type.set

                    code = (
                        'new io.sysl.Expr<{out_type}, {type}>() {{\n'
                        '        @Override\n'
                        '        public {out_type} evaluate(final {type} {scopevar}) {{\n'
                        '            {out_type} out = {out_type}._PRIVATE_new();\n'
                        '            {extra_code}return out;\n'
                        '        }}\n'
                        '    }}.evaluate({target})'.format(
                            target=target,
                            type=codeForType(target_type, cur_scope()),
                            out_type=out_type_code,
                            scopevar=(
                                'item' if transform.scopevar == '.' else transform.scopevar),
                            extra_code=str(w).replace('\n', '\n            ')))

                return (protect('.', code), expr.type)

    def unopJava():
        # TODO: decimal
        return {
            sysl_pb2.Expr.UnExpr.NEG:
            unexprToJava('.', 'io.sysl.ExprHelper.minus({})'),
            sysl_pb2.Expr.UnExpr.POS:
            unexprToJava('.', 'io.sysl.ExprHelper.plus({})'),
            sysl_pb2.Expr.UnExpr.NOT: '!',
            sysl_pb2.Expr.UnExpr.INV: '~',
            sysl_pb2.Expr.UnExpr.SINGLE: unexprToJava('.', '{}.single()'),
            sysl_pb2.Expr.UnExpr.SINGLE_OR_NULL:
            unexprToJava('.', '{}.singleOrNull()'),
        }
    UNOP_JAVA = unopJava()

    def setop(protect, binexpr):
        prec = PRECEDENCE[(2, binexpr.op)]
        (lhs, lhs_type) = E(prec, binexpr.lhs)
        (rhs, rhs_type) = E(prec, binexpr.rhs)
        if lhs_type.WhichOneof('type') == 'set':
            method = {
                sysl_pb2.Expr.BinExpr.BUTNOT: 'butnot',
                sysl_pb2.Expr.BinExpr.BITAND: 'and',
                sysl_pb2.Expr.BinExpr.BITOR: 'or',
                sysl_pb2.Expr.BinExpr.BITXOR: 'xor',
            }[binexpr.op]
            return (
                protect(PRECEDENCE['.'], '{}.{}({})'.format(lhs, method, rhs)),
                lhs_type)
        elif lhs_type.WhichOneof('type') == 'primitive':
            op = ({
                sysl_pb2.Expr.BinExpr.ADD: '{}.add({})',
                sysl_pb2.Expr.BinExpr.SUB: '{}.subtract({})',
                sysl_pb2.Expr.BinExpr.MUL: '{}.multiply({})',
                sysl_pb2.Expr.BinExpr.DIV: (
                  '{}.divide({}, java.math.RoundingMode.HALF_UP)'),
                sysl_pb2.Expr.BinExpr.MOD: '{}.remainder({})',
                sysl_pb2.Expr.BinExpr.POW: '{}.pow({})',
            } if lhs_type.primitive == lhs_type.DECIMAL else {
                sysl_pb2.Expr.BinExpr.ADD: '{} + {}',
                sysl_pb2.Expr.BinExpr.SUB: '{} - {}',
                sysl_pb2.Expr.BinExpr.MUL: '{} * {}',
                sysl_pb2.Expr.BinExpr.DIV: '{} / {}',
                sysl_pb2.Expr.BinExpr.MOD: '{} % {}',
            })[binexpr.op]
            return (protect(prec, op.format(lhs, rhs)), lhs_type)
        else:
            raise RuntimeError('Unsupported type')

    def coalesce(protect, binexpr):
        prec = PRECEDENCE[(2, binexpr.op)]
        (lhs_init_code, lhs_code, lhs_type) = tmpForExpr(binexpr.lhs)
        (rhs, rhs_type) = E(prec, binexpr.rhs)
        if not types_match(lhs_type, rhs_type):
            print >>sys.stderr, 'WARNING: Mismatched binop types: {} ≠ {}'.format(
                oneline(lhs_type), oneline(rhs_type))
        return (
            protect(prec, '{} == null \v? {} \v: {}').format(
                lhs_init_code, rhs, lhs_code),
            lhs_type)

    def binopJava():
        def comp(prec, expr, out_type=None, prefix=''):
            return binexprToJava(
                prec, prefix + 'io.sysl.ExprHelper.' + expr, out_type)

        return {
            sysl_pb2.Expr.BinExpr.EQ: comp('.', 'areEqual({}, {})', BOOL_TYPE),
            sysl_pb2.Expr.BinExpr.NE: comp(
                (1, sysl_pb2.Expr.UnExpr.NOT),
                'areEqual({}, {})',
                BOOL_TYPE,
                prefix='!'),
            sysl_pb2.Expr.BinExpr.LT: comp(None, 'doCompare({}, {}) < 0', BOOL_TYPE),
            sysl_pb2.Expr.BinExpr.LE: comp(None, 'doCompare({}, {}) <= 0', BOOL_TYPE),
            sysl_pb2.Expr.BinExpr.GT: comp(None, 'doCompare({}, {}) > 0', BOOL_TYPE),
            sysl_pb2.Expr.BinExpr.GE: comp(None, 'doCompare({}, {}) >= 0', BOOL_TYPE),
            sysl_pb2.Expr.BinExpr.IN: (
                binexprToJava(None, '{1}.contains({0})', BOOL_TYPE)),
            sysl_pb2.Expr.BinExpr.CONTAINS: (
                binexprToJava(None, '{}.contains({})', BOOL_TYPE)),
            sysl_pb2.Expr.BinExpr.NOT_IN: (
                binexprToJava(None, '!(({1}).contains({0}))', BOOL_TYPE)),
            sysl_pb2.Expr.BinExpr.NOT_CONTAINS: (
                binexprToJava(None, '!(({}).contains({}))', BOOL_TYPE)),

            sysl_pb2.Expr.BinExpr.ADD: setop,
            sysl_pb2.Expr.BinExpr.SUB: setop,
            sysl_pb2.Expr.BinExpr.MUL: setop,
            sysl_pb2.Expr.BinExpr.DIV: setop,
            sysl_pb2.Expr.BinExpr.MOD: setop,
            sysl_pb2.Expr.BinExpr.POW: setop,

            sysl_pb2.Expr.BinExpr.AND: '&&',
            sysl_pb2.Expr.BinExpr.OR: '||',

            sysl_pb2.Expr.BinExpr.BUTNOT: setop,
            sysl_pb2.Expr.BinExpr.BITAND: setop,
            sysl_pb2.Expr.BinExpr.BITOR: setop,
            sysl_pb2.Expr.BinExpr.BITXOR: setop,

            sysl_pb2.Expr.BinExpr.COALESCE: coalesce,

            sysl_pb2.Expr.BinExpr.WHERE: whereToJava,
            sysl_pb2.Expr.BinExpr.FLATTEN: flattenToJava,
            sysl_pb2.Expr.BinExpr.TO_MATCHING: toMatchingToJava(True),
            sysl_pb2.Expr.BinExpr.TO_NOT_MATCHING: toMatchingToJava(False),
        }
    BINOP_JAVA = binopJava()

    def relopJava():
        return {
            sysl_pb2.Expr.RelExpr.MIN: aggexprToJava,
            sysl_pb2.Expr.RelExpr.MAX: aggexprToJava,
            sysl_pb2.Expr.RelExpr.SUM: aggexprToJava,
            sysl_pb2.Expr.RelExpr.AVERAGE: aggexprToJava,
            sysl_pb2.Expr.RelExpr.RANK: rankToJava,
            sysl_pb2.Expr.RelExpr.SNAPSHOT: snapshotToJava,
            sysl_pb2.Expr.RelExpr.FIRST_BY: firstByToJava,
        }
    RELOP_JAVA = relopJava()

    return E(-1, expr, let)
