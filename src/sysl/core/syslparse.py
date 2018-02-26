# -*- encoding: utf-8 -*-
"""Parse sysl source code."""

# TODO: Allow attr grouping via, e.g.: [~pk]:\n\tfooId...\n\tbarId...

import ast
import contextlib
import getpass
import json
import operator
import os
import sqlite3
import sys
import threading
import time
import urlparse
import re

from sysl.proto import sysl_pb2

from sysl.core import syslx

from sysl.util import simple_parser
from sysl.util import rex
from sysl.util import cache


USER_ENV = 'CONFLUENCE_USER'

REST_METHODS = {key for key in sysl_pb2.Endpoint.RestParams.Method.keys()
                if '_' not in key}

BUILTINS = {
    'true': sysl_pb2.Expr(
        literal=sysl_pb2.Value(b=True),
        type=sysl_pb2.Type(primitive=sysl_pb2.Type.BOOL)),
    'false': sysl_pb2.Expr(
        literal=sysl_pb2.Value(b=False),
        type=sysl_pb2.Type(primitive=sysl_pb2.Type.BOOL)),
    'null': sysl_pb2.Expr(
        literal=sysl_pb2.Value(null=sysl_pb2.Value.Null()),
        type=sysl_pb2.Type(primitive=sysl_pb2.Type.EMPTY)),
}


class SyntaxError(Exception):
    """Exception generated when a syntax error is detected."""

    def __init__(self, parser, fmt, *args, **kwargs):
        source_context = parser.source_context((None, None))
        super(Exception, self).__init__(
            fmt.format(*args, **kwargs), source_context)
        self.source_context = source_context


class _Matcher(object):
    """Match regexes and hold results."""
    groups = property(operator.attrgetter('_groups'))

    def __init__(self, line):
        self._line = line
        self._groups = ()

    def __call__(self, pat):
        match = rex.search(pat, self._line, rex.VERBOSE)
        if match:
            self._groups = tuple(match.groups())
        return match


class _AttrParser(simple_parser.SimpleParser):
    """Parse attribute syntax."""

    def __init__(self, text, attrs, source_context=None):
        super(_AttrParser, self).__init__(text, source_context)
        self._attrs = attrs

    def parse(self):
        """Top-level parse method."""
        if self._attr():
            while self.eat(r'\s*,'):
                self._attr()
            assert not self
        return True

    @staticmethod
    def _populate_attr(attr, value):
        """Populate an attribute based on the type of value."""
        if isinstance(value, basestring):
            attr.s = value
        elif isinstance(value, int):
            attr.i = value
        elif isinstance(value, list):
            for elt in value:
                _AttrParser._populate_attr(attr.a.elt.add(), elt)

    def _attr(self):
        """Parse an attribute."""
        if self.eat(r'\s*([-\.\w]+)\s*=\s*'):
            if not self._value():
                self._syntax_error('Missing attr name', self._text)
            (name, value) = self.pop(2)
            _AttrParser._populate_attr(self._attrs[name], value)
            return True

        # Shorthand: ~foo, ~bar => patterns=["foo", "bar"]
        if self.eat(r'\s*~([^,\s]+)\s*'):
            self._attrs['patterns'].a.elt.add().s = self.pop()
            return True

    def _value(self):
        """Parse a value."""
        if self.eat(r'\s*\[', list):
            arr = self[-1]
            if self._value():
                arr.append(self.pop())
                while self.eat(r'\s*,') and self._value():
                    arr.append(self.pop())
            if not self.eat(r'\s*\]'):
                raise Exception('syntax error' + self._text)
            return True

        def strparse(s):
            return ast.literal_eval('u' + s)
        return (
            self.eat(r'\s*("(?:[^"\\\n]|\\.)*")', strparse) or
            self.eat(r"\s*('(?:[^'\\\n]|\\.)*')", strparse) or
            self.eat(r'\s*([-+]?[1-9]\d*)', int))


class _AttrProcessor(object):
    def __init__(self, attrs=None):
        self.attrs = attrs or []

    def __iadd__(self, attr):
        if attr:
            if isinstance(attr, _AttrProcessor):
                for a in attr.attrs:
                    self.attrs.append(a)
            else:
                self.attrs.append(attr)
        return self

    def __add__(self, attr):
        result = _AttrProcessor(self.attrs[:])
        result += attr
        return result

    def __rshift__(self, field):
        if self.attrs and field is not None:
            _AttrParser(','.join(self.attrs), field)()
        del self.attrs[:]

    def __str__(self):
        return str(self.attrs)


class _ExprParser(simple_parser.SimpleParser):
    def __init__(self, text, children, source_context=None):
        super(_ExprParser, self).__init__(text, source_context)
        self.children = children
        self.skip_ws = True

    def parse(self):
        if not self._expr():
            self.__init__(self._text, self._source_context)
            self._expr()
            raise RuntimeError('syntax error')
        return self.pop()

    def _expr(self):
        return self._ifelse()

    def _unop(self, regex, opmap, descend):
        if not self.eat(regex):
            return descend()
        op = opmap[self.pop()]
        if not descend():
            raise RuntimeError('syntax error')
        arg = self.pop()
        [expr] = self.push(sysl_pb2.Expr())
        expr.unexpr.op = op
        expr.unexpr.arg.CopyFrom(arg)
        return True

    def _binop(self, regex, opmap, descend_lhs, descend_rhs=None):
        if descend_rhs is None:
            descend_rhs = descend_lhs

        if not descend_lhs():
            return False

        while self.eat(regex):
            op = opmap[self.pop()]
            lhs = self.pop()
            if not descend_rhs():
                raise RuntimeError('syntax error')
            rhs = self.pop()
            self.push(sysl_pb2.Expr(
                binexpr=sysl_pb2.Expr.BinExpr(
                    op=op, lhs=lhs, rhs=rhs)))

        return True

    def _ifelse(self):
        if self.eat(r'if\b'):
            if not self._expr():
                raise RuntimeError('syntax error')
            nullsafe = self.eat('\?')
            self.expect(r'then\b')
            if not self._expr():
                raise RuntimeError('syntax error')
            self.expect(r'else\b')
            if not self._expr():
                raise RuntimeError('syntax error')

            (a, b, c) = self.pop(3)
            self.push(sysl_pb2.Expr(
                ifelse=sysl_pb2.Expr.IfElse(
                    cond=a, if_true=b, if_false=c, nullsafe=nullsafe)))

            return True

        return self._coalesce()

    def _coalesce(self):
        return self._binop(
            r'(\?\?)', {'??': sysl_pb2.Expr.BinExpr.COALESCE}, self._butnot)

    def _butnot(self):
        return self._binop(
            r'(but\ not\b)', {
                'but not': sysl_pb2.Expr.BinExpr.BUTNOT,
            }, self._or)

    def _or(self):
        return self._binop(
            r'(\|\|)', {'||': sysl_pb2.Expr.BinExpr.OR}, self._and)

    def _and(self):
        return self._binop(
            r'(&&)', {'&&': sysl_pb2.Expr.BinExpr.AND}, self._bitor)

    def _bitor(self):
        return self._binop(
            r'(\|(?!\|)|or\b)', {
                '|': sysl_pb2.Expr.BinExpr.BITOR,
                'or': sysl_pb2.Expr.BinExpr.BITOR,
            }, self._bitxor)

    def _bitxor(self):
        return self._binop(
            r'(\^|xor)', {'^': sysl_pb2.Expr.BinExpr.BITXOR}, self._bitand)

    def _bitand(self):
        return self._binop(
            r'(&(?!&)|and)', {'&': sysl_pb2.Expr.BinExpr.BITAND}, self._rel)

    def _rel(self):
        return self._binop(
            r'(!?in|!?contains\b|!=|==|<=?|>=?)',
            {'==': sysl_pb2.Expr.BinExpr.EQ,
             '!=': sysl_pb2.Expr.BinExpr.NE,
             '<': sysl_pb2.Expr.BinExpr.LT,
             '<=': sysl_pb2.Expr.BinExpr.LE,
             '>': sysl_pb2.Expr.BinExpr.GT,
             '>=': sysl_pb2.Expr.BinExpr.GE,
             'in': sysl_pb2.Expr.BinExpr.IN,
             'contains': sysl_pb2.Expr.BinExpr.CONTAINS,
             '!in': sysl_pb2.Expr.BinExpr.NOT_IN,
             '!contains': sysl_pb2.Expr.BinExpr.NOT_CONTAINS,
             },
            self._shift)

    def _shift(self):
        return self._addsub()

    def _addsub(self):
        return self._binop(
            r'(\+|-(?!>))',
            {'-': sysl_pb2.Expr.BinExpr.SUB,
             '+': sysl_pb2.Expr.BinExpr.ADD},
            self._muldiv)

    def _muldiv(self):
        return self._binop(
            r'([/%]|\*(?!\*))',
            {'*': sysl_pb2.Expr.BinExpr.MUL,
             '/': sysl_pb2.Expr.BinExpr.DIV,
             '%': sysl_pb2.Expr.BinExpr.MOD},
            self._unadd)

    def _unadd(self):
        return self._unop(
            r'([+~]|-(?!>)|!(?!=))',
            {'-': sysl_pb2.Expr.UnExpr.NEG,
             '+': sysl_pb2.Expr.UnExpr.POS,
             '~': sysl_pb2.Expr.UnExpr.INV,
             '!': sysl_pb2.Expr.UnExpr.NOT},
            self._pow)

    def _pow(self):
        return self._binop(
            r'(\*\*)',
            {'**': sysl_pb2.Expr.BinExpr.POW},
            self._call,
            self._unadd)

    RELOPS = (
        r'where|singleOrNull|single|count|any|flatten|'
        r'sum|min|max|average|rank|first|snapshot'
    )

    def _call(self):
        # foo.y
        # foo?.y
        # foo -> y
        # foo where(true)
        if self._atom():
            return self._call_tail()

        # ...
        if self.eat(r'\.\.\.'):
            self.push(sysl_pb2.Expr())
            return True

        # (.) = "implied dot"
        # (.)?.y
        # (.) -> y
        # (.) ~> expr
        # (.) !~> expr
        if self.eat(ur'(?=\?\.|->|!?~(?:\[(?:(?:\w+,)*\w+)\])?>)'):
            self.push(sysl_pb2.Expr(name='.'))
            return self._call_tail()

        # . ->/where()/...
        if self.eat(r'\.(?=\s*(?:\??->|{})\b)'.format(self.RELOPS)):
            self.push(sysl_pb2.Expr(name='.'))
            return self._call_tail()

        # .
        # (.).y
        cur = self._cur
        if self.eat(r'\.'):
            self.push(sysl_pb2.Expr(name='.'))
            if self._ident():
                self.pop(2)
                self._cur = cur
            return self._call_tail()

        return False

    TAIL_PATTERN = r'''
        (
          (?:(\??)(?:\[|\(|->|!?~(?:\[((?:\w+,)*\w+)\])?>|\.))|
          (?:{})\b
        )
        '''.format(RELOPS)

    def _call_tail(self, single=False):
        while self.eat(self.TAIL_PATTERN):
            (op, nullsafe, squiggly_args) = self.pop(3)
            if nullsafe:
                op = op[1:]

            if op == '[':
                raise RuntimeError('[...] obsolete; use "~>" instead')

            elif op == '(':
                arg = self.pop()
                [lhs] = self.push(sysl_pb2.Expr(
                    call=sysl_pb2.Expr.Call(func=arg.name)))
                self._list(lhs.call.arg, r'\)')

            # -> <ident>(...
            # -> (...
            elif op == '->' and self.eat(r'(?:[<(])'):
                if self._type_spec():
                    (out_type, setof) = self.pop(2)
                else:
                    (out_type, setof) = (None, False)
                    self.expect(r'\(')

                lhs = self.pop()
                [expr] = self.push(sysl_pb2.Expr())
                transform = expr.transform
                transform.arg.CopyFrom(lhs)
                if self_.ident():
                    (ident, scopevar_setof) = self.pop(2)
                    assert not scopevar_setof
                    transform.scopevar = ident
                else:
                    transform.scopevar = '.'
                transform.nullsafe = bool(nullsafe)
                if self.children:
                    if self.eat(r'$'):
                        raise RuntimeError('syntax error')
                    break
                else:
                    self._transform_list(transform)
                    self.expect('\)')

                if out_type:
                    target_type = sysl_pb2.Type()
                    target_path = out_type.split('.')
                    target_type.type_ref.ref.appname.part.extend(
                        target_path[:1])
                    target_type.type_ref.ref.path.extend(target_path[1:])
                    if setof:
                        expr.type.set.CopyFrom(target_type)
                    else:
                        expr.type.CopyFrom(target_type)
                elif setof:
                    expr.type.set.CopyFrom(sysl_pb2.Type())

            elif op == '->':
                lhs = self.pop()
                if not self.eat(r'(?:\b(set)\s+of\s+)?(\.?\w+)'):
                    raise RuntimeError('syntax error')
                (setof, ident) = self.pop(2)
                via = self.pop() if self.eat(r'via\s+(\w+)\b') else ''
                self.push(sysl_pb2.Expr(
                    navigate=sysl_pb2.Expr.Navigate(
                        arg=lhs,
                        attr=ident,
                        nullsafe=bool(nullsafe),
                        setof=bool(setof),
                        via=via)))

            elif op[:2] in ['~>', '~[', '!~'] and op[-1:] == '>':
                lhs = self.pop()

                # Treat implied dot as an rhs expression

                # (.)?.y        # (.) = "implied dot"
                # (.) -> y
                if self.eat(r'(?=\?\.|->)'):
                    self.push(sysl_pb2.Expr(name='.'))
                    if not self._call_tail(single=True):
                        raise RuntimeError('syntax error')

                # NOT . ->/where()/...
                elif self.eat(r'\.(?!\s*(?:\??->|{})\b)'.format(self.RELOPS)):
                    # .
                    # (.).y
                    cur = self._cur - 1
                    self.push(sysl_pb2.Expr(name='.'))
                    if self._ident():
                        self.pop(2)
                        self._cur = cur
                        if not self._call_tail(single=True):
                            raise RuntimeError('syntax error')

                elif not self._atom():
                    raise RuntimeError('missing atom after "~>"')

                rhs = self.pop()
                self.push(sysl_pb2.Expr(
                    binexpr=sysl_pb2.Expr.BinExpr(
                        op=(
                          sysl_pb2.Expr.BinExpr.TO_NOT_MATCHING if op[:1] == '!' else
                          sysl_pb2.Expr.BinExpr.TO_MATCHING),
                        lhs=lhs,
                        rhs=rhs,
                        attr_name=rex.split(ur'·,·', squiggly_args) if squiggly_args else ['*'])))

            elif op in ['where', 'flatten']:
                self.expect(r'\(')
                scopevar = self.pop() if self.eat(r'(\w*):') else '.'

                if not self._expr():
                    raise RuntimeError('missing expr')
                pred = self.pop()

                self.expect(r'\)')

                self.push(sysl_pb2.Expr(
                    binexpr=sysl_pb2.Expr.BinExpr(
                        op=sysl_pb2.Expr.BinExpr.Op.Value(op.upper()),
                        lhs=self.pop(),
                        rhs=pred,
                        scopevar=scopevar)))

            elif op in ['sum', 'min', 'max', 'average']:
                self.expect(r'\(')
                scopevar = self.pop() if self.eat(r'(\w*):') else '.'
                if not self._expr():
                    raise RuntimeError('missing expr')

                args = [self.pop()]
                while self.eat(',') and self._expr():
                    args.append(self.pop())

                self.expect(r'\)')

                [expr] = self.push(sysl_pb2.Expr(
                    relexpr=sysl_pb2.Expr.RelExpr(
                        op=sysl_pb2.Expr.RelExpr.Op.Value(op.upper()),
                        target=self.pop(),
                        scopevar=scopevar)))

                for arg in args:
                    expr.relexpr.arg.add().CopyFrom(arg)

            elif op in ['single', 'singleOrNull']:
                ornull = op == 'singleOrNull'
                lhs = self.pop()
                [expr] = self.push(sysl_pb2.Expr())
                expr.unexpr.op = (
                    sysl_pb2.Expr.UnExpr.SINGLE_OR_NULL if ornull else
                    sysl_pb2.Expr.UnExpr.SINGLE)
                expr.unexpr.arg.CopyFrom(lhs)

            elif op == 'count':
                lhs = self.pop()
                [call] = self.push(sysl_pb2.Expr())
                call.call.func = '.count'
                call.call.arg.add().CopyFrom(lhs)

            elif op == 'any':
                self.expect(r'\(')
                arg = self.pop()
                limit = self._expr() and self.pop()
                self.expect(r'\)')
                [lhs] = self.push(sysl_pb2.Expr(
                    call=sysl_pb2.Expr.Call(
                        func='.any')))
                lhs.call.arg.add().CopyFrom(arg)
                if limit:
                    lhs.call.arg.add().CopyFrom(limit)

            elif op == 'rank':
                target = self.pop()

                if self._type_spec():
                    (out_type, setof) = self.pop(2)
                else:
                    (out_type, setof) = (None, False)

                assert not setof, '"set of" not required (or allowed) in rank<...>()'

                self.expect(r'\(')
                scopevar = self.pop() if self.eat(r'(\w*):') else '.'
                if not self._expr():
                    raise RuntimeError('missing expr')
                args = [self.pop()]
                desc = [self.eat(r'(asc|desc)\b') and self.pop() == 'desc']
                while self.eat(r','):
                    if not self._expr():
                        raise RuntimeError('missing expr')
                    args.append(self.pop())
                    desc.append(self.eat(r'(asc|desc)\b') and
                                self.pop() == 'desc')
                self.expect(ur'as\s+(\w+)·\)')
                rank_attr = self.pop()
                [expr] = self.push(sysl_pb2.Expr(
                    relexpr=sysl_pb2.Expr.RelExpr(
                        op=sysl_pb2.Expr.RelExpr.RANK,
                        target=target,
                        arg=args,
                        scopevar=scopevar,
                        descending=desc,
                        attr_name=[rank_attr])))

                for arg in args:
                    expr.relexpr.arg.add().CopyFrom(arg)

                if out_type:
                    target_type = sysl_pb2.Type()
                    target_path = out_type
                    target_type.type_ref.ref.appname.part.extend(
                        target_path[:1])
                    target_type.type_ref.ref.path.extend(target_path[1:])
                    expr.type.set.CopyFrom(target_type)
                else:
                    expr.type.set.CopyFrom(sysl_pb2.Type())

            elif op == 'first':
                target = self.pop()

                args = (
                    [self.pop()] if self._expr() else
                    [sysl_pb2.Expr(
                        literal=sysl_pb2.Value(
                            null=sysl_pb2.Value.Null()))])

                self.expect(ur'by·\(')

                scopevar = self.pop() if self.eat(r'(\w*):') else '.'
                if not self._expr():
                    raise RuntimeError('missing expr')
                args.append(self.pop())
                desc = [self.eat(r'(asc|desc)\b') and self.pop() == 'desc']
                while self.eat(r','):
                    if not self._expr():
                        raise RuntimeError('missing expr')
                    args.append(self.pop())
                    desc.append(self.eat(r'(asc|desc)\b') and
                                self.pop() == 'desc')
                self.expect(ur'\)')
                [expr] = self.push(sysl_pb2.Expr(
                    relexpr=sysl_pb2.Expr.RelExpr(
                        op=sysl_pb2.Expr.RelExpr.FIRST_BY,
                        target=target,
                        arg=args,
                        scopevar=scopevar,
                        descending=desc)))

                for arg in args:
                    expr.relexpr.arg.add().CopyFrom(arg)

            elif op == 'snapshot':
                self.push(sysl_pb2.Expr(
                    relexpr=sysl_pb2.Expr.RelExpr(
                        op=sysl_pb2.Expr.RelExpr.SNAPSHOT,
                        target=self.pop())))

            elif op == '.':
                if not self._ident():
                    raise RuntimeError('syntax error')
                (ident, setof) = self.pop(2)
                if setof == 'set':
                    raise RuntimeError('syntax error')
                expr = sysl_pb2.Expr.GetAttr(
                    arg=self.pop(), attr=ident, nullsafe=bool(nullsafe),
                    setof=bool(setof))
                self.push(sysl_pb2.Expr(get_attr=expr))

            else:
                assert not "Oops!", op

            dot_prefix = False

            if single:
                break

        return True

    def _transform_list(self, transform):
        if not self._transform_item(transform):
            return False

        while self.eat(r',') and self._transform_item(transform):
            pass
        return True

    def _transform_item(self, transform):
        if self.eat(r'\*(?:\s*-\s*(?:(\w+)|{((?:\w+\s*,\s*)*\w+)}))?'):
            [exclude1, exclude2] = self.pop(2)
            transform.all_attrs = True
            transform.except_attrs.extend(
                (([exclude1] if exclude1 else []) +
                 (rex.split(r'\s*,\s*', exclude2) if exclude2 else [])))
            return True
        return False

    def _atom(self):
        if self.eat(r'\('):
            if not self._expr():
                raise RuntimeError('syntax error')
            self.expect(r'\)')
            return True

        if self.eat(r'\['):
            [expr] = self.push(sysl_pb2.Expr())
            return self._list(expr.list.expr, r'\]')

        if self.eat(r'{'):
            [expr] = self.push(sysl_pb2.Expr())
            if self.eat(ur':·}'):
                expr.tuple.CopyFrom(sysl_pb2.Expr.Tuple())
                expr.type.tuple.CopyFrom(sysl_pb2.Type.Tuple())
                return True
            if self._list(expr.set.expr, r'}'):
                exprs = expr.set.expr
                if exprs and all(e.type == exprs[0].type for e in exprs):
                    expr.type.set.CopyFrom(exprs[0].type)
                else:
                    expr.type.set.primitive = expr.type.ANY
                    if not expr.set.expr:
                        expr.set.CopyFrom(expr.List())
                return True
            return False

        if self.eat(r'(\d+(?:\.\d+)?(?:[Ee][-+]?\d+)?f?)'):
            value = self.pop()
            [expr] = self.push(sysl_pb2.Expr())
            if value.endswith('f'):
                expr.literal.d = float(value[:-1])
            elif rex.match(r'\d+$', value):
                expr.literal.i = int(value)
            else:
                expr.literal.decimal = value
            return True

        if self.eat(r'("(?:[^"]|\\.)*"|\'(?:[^\']|\\.)*\')'):
            value = ast.literal_eval(self.pop())
            [expr] = self.push(sysl_pb2.Expr())
            expr.literal.s = value
            return True

        if self._ident():
            (ident, setof) = self.pop(2)
            if setof:
                raise RuntimeError('syntax error')
            builtin = BUILTINS.get(ident)
            if builtin:
                self.push(builtin)
            else:
                self.push(sysl_pb2.Expr(name=ident))
            return True
        return False

    def _ident(self, action=None):
        if self.eat(r'(?:\b(set|table)\s+of\s+)?(\w+)'):
            [setof, ident] = self.pop(2)
            if action:
                self.push(action(ident, setof))
            else:
                self.push(ident, setof)
            return True
        return False

    def _list(self, target, delim):
        while self._expr():
            target.add().CopyFrom(self.pop())
            if not self.eat(r','):
                break
        self.expect(delim)
        return True

    def _type_spec(self):
        if self.eat(r'<'):
            if self._ident():
                (out_type, setof) = self.pop(2)
                out_type = [out_type]
                setof = bool(setof)
                while self.eat(r'\.'):
                    if not self.eat(r'(\w+)'):
                        raise RuntimeError('syntax error')
                    out_type.append(self.pop())
                self.push(out_type, setof)
            elif self.eat(r'set\s+of\b'):
                setof = True
                self.push(None, True)
            else:
                raise RuntimeError('syntax error')
            self.expect(r'>')
            return True
        return False


class Parser(object):
    def source_context(self, start=None, end=None):
        if start is None:
            (start, end) = self.context

        sctx = sysl_pb2.SourceContext()

        def set_pos(pos, field):
            if pos:
                if pos[0] is not None:
                    field.line = pos[0]
                if pos[1] is not None:
                    field.col = pos[1]

        set_pos(start, sctx.start)
        set_pos(end, sctx.end)

        return sctx

    def _syntax_error(self, fmt, *args, **kwargs):
        """Handle a syntax error.

        Passes the error to self.report_error, if present, or else raises it.
        """
        err = SyntaxError(self, fmt, *args, **kwargs)
        if hasattr(self, 'report_error'):
            self.report_error(err)
        else:
            raise err

    def _parse_indent(self, src):
        """Parse text indentation to simplify semantic parsing."""
        try:
            # Ugly hack to combine multiline strings
            lines = []
            accum = None
            accum_indent = None
            line_no = -1
            for (line_no, line) in enumerate(src):
                if rex.match(ur'·\#', line):
                    continue
                (indent, pipe, text) = rex.match(
                    ur'^(·)(\|?)(.*?)·$', line).groups()
                if pipe:
                    assert text[:1] in (' ', ''), (
                        'lines starting with "|" must also have a '
                        'space immediately after or be empty.')
                    text = text[1:].replace('\\n', '\n').replace('\\\\', '\\')
                    if accum and indent != accum_indent:
                        lines.append(accum)
                        accum = None
                    if not accum:
                        accum = [line_no, indent, '']
                        accum_indent = indent

                    # Empty lines (after the "|") are paragraph markers.
                    if text:
                        if accum[2][-1:] not in ('\n', ''):
                            accum[2] += ' '
                        accum[2] += text
                    else:
                        accum[2] += '\n\n'
                else:
                    if accum:
                        lines.append(accum[:2] + ['| ' + accum[2]])
                        accum = None
                    lines.append((line_no, indent, text))
            if accum:
                lines.append(accum)
                accum = None

            end_line_no = line_no

            def postprocess():
                """Detect, validate and output indent information."""
                indent_stk = ['']
                prev_indent = ''

                for (line_no, indent, line) in lines:
                    if line:
                        if indent != prev_indent:
                            if indent.startswith(prev_indent):
                                indent_stk.append(indent)
                            elif prev_indent.startswith(indent):
                                while indent_stk[-1] > indent:
                                    indent_stk.pop()
                                if indent_stk[-1] != indent:
                                    raise Exception(
                                        'outdentation mismatch', line_no + 1, line)
                            else:
                                raise Exception(
                                    'indentation mismatch', line_no + 1, line)
                        prev_indent = indent
                    yield (line_no + 1, len(indent_stk) if line else None, line)

                # Ensure an empty trailing line to help self._group_indent().
                yield (end_line_no + 1, 0, '')

            return postprocess()
        except BaseException:
            print >>sys.stderr, 'ERROR @ line ', line_no + 1, line
            raise

    def _group_indent(self, src, path):
        """Group indented lines to simplify semantic parsing."""
        src = list(self._parse_indent(src))

        def inner(i=0, outer=0):
            """Internal function to perform recursion."""
            result = []
            last_nonempty_line_no = None
            while i < len(src):
                (line_no, indent, line) = src[i]
                if indent is None:
                    i += 1
                    continue
                if not line.startswith('#'):
                    last_nonempty_line_no = line_no
                if indent <= outer:
                    return (i, result, last_nonempty_line_no)
                (new_i, children, end_line_no) = inner(i + 1, indent)
                if new_i is None:
                    return (i, result, last_nonempty_line_no)
                i = new_i
                if not line.startswith('#'):
                    result.append((line_no, end_line_no, line, children))

            if last_nonempty_line_no is None:
                return (None, None, None)
            return (i, result, last_nonempty_line_no)

        def coalesce(lines):
            return ' '.join(
                line + ' ' + coalesce(children) for (_, _, line, children) in lines)

        def check(lines):
            result = []
            i = 0
            while i < len(lines):
                (line_no, end_line_no, line, children) = lines[i]

                if line.endswith(':') and not line.startswith('|'):
                    if not children:
                        raise RuntimeError(
                            "{}: lines ending in ':' must have children: {}"
                            .format(line_no, line))
                    line = line[:-1]
                    children = check(children)
                    if line.startswith('@'):
                        if not line.endswith('='):
                            raise RuntimeError("@attr: must be @attr =:")
                        if not all(c.startswith('|')
                                   for (_, _, c, _) in children):
                            raise RuntimeError(
                                "@attr =: children must be multiline string")
                        no_space = [
                            c for (_, _, c, _) in children
                            if not c.startswith('| ')]
                        if no_space:
                            raise RuntimeError(
                                "@attr: multiline row(s) missing space after '|': " + no_space)
                        cc = json.dumps(' '.join(c[2:]
                                                 for (_, _, c, _) in children))
                        line += cc
                        children = []

                elif line.startswith('@'):
                    # Coalesce @-lines with their descendants.
                    cc = coalesce(children)
                    line += cc and (' ' + cc)
                    children = []
                    if i + 1 < len(lines):
                        (line_no2, end_line_no2, line2,
                         children2) = lines[i + 1]
                        # Also coalesce closing brackets following @-line.
                        if rex.match('^[\]\}\)]+$', line2):
                            end_line_no = end_line_no2
                            line += ' ' + line2
                            i += 1

                elif children:
                    raise RuntimeError(
                        "{}: lines with children must match /^@|:$/"
                        .format(line_no))

                result.append((line_no, end_line_no, line, children))

                i += 1

            return result

        (_, result, end_line_no) = inner()
        return check(result)

    def _parse_app_name(self, field, name, app):
        """Parse an app name."""
        if name == '.':
            field.CopyFrom(app.name)
        else:
            del field.part[:]
            field.part.extend(part.strip() for part in name.split('::'))

    def _parse_attrs(self, attrs, field):
        """Parse an attribute spec."""
        if attrs:
            _AttrParser(attrs, field)()

    def _parse_rest_api(
            self, app, oldpath, newpath, extra_attrs, path_params, children):
        """Parse ReST path hierarchy, triggered by '/...' endpoint name."""

        path_params = path_params.copy()

        def prepare_endpoint(method, params, grandchildren):
            """Prepare ReST endpoint."""
            epname = method + ' ' + path
            endpt = app.endpoints[epname]
            endpt.name = epname

            if params is not None:
                for param_pair in rex.split(r'\s*,\s*', params):
                    param = endpt.param.add()
                    if param_pair.count('<:') == 1:
                        (pname, ptname) = rex.split(r'\s*<:\s*', param_pair)
                        param.name = pname
                        self._parseGlobalRef(ptname, param.type, app)

            endpt.attrs['patterns'].a.elt.add().s = 'rest'
            (extra_attrs + attrs) >> endpt.attrs
            if grandchildren[0][2].startswith('| '):
                endpt.docstring = grandchildren[0][2][2:]
                del grandchildren[0]
            self._parse_stmts(endpt, grandchildren, app)
            return endpt

        def prepare_rest_params(method, endpt, line_no):
            """Prepare endpoint rest_params."""
            rest = endpt.rest_params
            rest.method = getattr(sysl_pb2.Endpoint.RestParams, method)
            rest.path = path

            qparams = []
            if query is not None:
                q = query.lstrip('?')

                for (name, values) in (
                        urlparse.parse_qs(q, strict_parsing=True).iteritems()):
                    if name in path_params:
                        raise RuntimeError('qparam {} duplicates path param: {}'.format(
                            name, path.sub('{' + name + '}', r'\033[1;33m\1\033[0m')))
                    [type_] = values
                    qparams.append(
                        (name, (type_, sysl_pb2.Endpoint.RestParams.QueryParam.QUERY)))

            for (name, (type_, loc)) in path_params.items() + qparams:
                qp = rest.query_param.add()
                qp.name = name
                qp.loc = loc
                self._parse_type(app, line_no, '', type_, qp.type)

        attrp = _AttrProcessor()

        newpath2 = ''
        parts = newpath.lstrip('/').split('/')
        for (i, part) in enumerate(parts):
            if '{' in part:
                m = rex.match(ur'{([-\w]+)·<:·(.*?)}$', part)
                if not m:
                    raise RuntimeError('rogue "{" found in rest path: ' +
                                       '/'.join(
                                           parts[:i] +
                                           ['\033[1;31m' + parts[i] + '\033[0m'] +
                                           parts[i + 1:]))
                [name, type_] = m.groups()
                if name in path_params:
                    raise RuntimeError('duplicate path param: {}'.format(
                        rex.sub('({' + name + '.*?})', r'\033[1;33m\1\033[0m', oldpath)))
                path_params[name] = (
                    type_, sysl_pb2.Endpoint.RestParams.QueryParam.PATH)
                newpath2 += '/{' + name + '}'
            else:
                newpath2 += '/' + part

        path = oldpath + newpath2

        for (line_no, _, line, grandchildren) in children:
            match = _Matcher(line)

            if match(r'^@(.*)'):
                attrp += match.groups[0]
            elif match(ur'^(/\S*)(?:·\[(.*)\])?$'):
                (subpath, extra_extra_attrs) = match.groups
                self._parse_rest_api(
                    app, path, subpath, attrp + extra_attrs + extra_extra_attrs,
                    path_params, grandchildren)
            elif match(ur'^(\w+)·(\([^)]*\))?(?:\s+(\?\S*))?(?:·\[(.*)\])?·$'):
                (method, params, query, attrs) = match.groups
                (attrp + attrs) >> None
                if method not in REST_METHODS:
                    raise RuntimeError(
                        'Disallowed ReST method (expecting {}): {}'.format(
                            '|'.join(REST_METHODS), method))

                if params is not None:
                    params = params[1:-1]
                endpt = prepare_endpoint(method, params, grandchildren)
                prepare_rest_params(method, endpt, line_no)
            else:
                raise RuntimeError('Invalid ReST syntax: ' + line)

    def _parse_type(self, app, line_no, name, typespec, type_):
        """Parse type syntax."""
        match = _Matcher(typespec)

        if match(ur'''
        ^
        (?:(set|list)\s+of\s)?
        ·([^(?]*?[^\s?])
        (?:
          ·\(·
            (?:
              (?:(\d*)\.\.)?
              (\d*)
              (?:[,\.](\d+))?
            )
          ·\)
        )?
        ·(\??)·
        $
        '''):
            (collection, ident, lo, hi, scale, optional) = match.groups
            assert not (collection and optional), (collection, optional, type_)

            if optional:
                type_.opt = True
            elif collection:
                type_ = getattr(type_, collection)

            primitive = ('_' not in ident) and getattr(
                type_, ident.upper(), None)

            if primitive:
                type_.primitive = primitive
                precision = hi
                if lo or hi or precision or scale:
                    c = type_.constraint.add()
                    if lo:
                        c.length.min = int(lo)
                    if hi:
                        c.length.max = int(hi)
                    if scale:
                        c.precision = int(precision)
                        c.scale = int(scale)
                    elif primitive in [type_.STRING, type_.STRING_8, type_.BYTES]:
                        c.length.min = 0
                        if hi:
                            c.length.max = int(hi)
            else:
                m = rex.match(ur'^(?i)(n)?(var)?char$', ident)
                if m:
                    (n, var) = m.groups()
                    type_.primitive = (
                        sysl_pb2.Type.STRING if n else
                        sysl_pb2.Type.STRING_8)
                    if not hi:
                        raise RuntimeError(ident + ' must have lenspec')
                    hi = int(hi)
                    lo = int(lo) if lo else 0 if var else hi
                    if not var and hi != lo:
                        raise RuntimeError('{} must have fixed lenspect, not {}..{}'
                                           .format(ident, lo, hi))
                    c = type_.constraint.add()
                    c.length.min = lo
                    c.length.max = hi
                else:
                    m = rex.match(r'^int(32|64)$', ident)
                    if m:
                        type_.primitive = sysl_pb2.Type.INT
                        c = type_.constraint.add()

                        size = int(m.group(1))
                        hi = 2**(size - 1) - 1
                        lo = ~hi
                        c.range.min.i = lo
                        c.range.max.i = hi
                    else:
                        tr = type_.type_ref
                        tr.context.appname.CopyFrom(app.name)
                        tr.context.path.extend(
                            s.strip()
                            for s in (name and name.split('.')))
                        (appname, _, path) = ident.rpartition('::')
                        if appname:
                            tr.ref.appname.part.extend(appname.split('::'))
                        tr.ref.path.extend(s.strip() for s in path.split('.'))

            type_.source_context.start.line = line_no
        else:
            raise RuntimeError('Unrecognised typespec: ' + typespec)

    def _parse_enum_defn(self, app, name, children):
        """Parse enum definitions."""
        typedecl = app.types[name]
        enum = typedecl.enum

        for (line_no, end_line_no, line, grandchildren) in children:
            match = _Matcher(line)

            if match(ur'(\w+)\s*=\s*(-?\d+)$'):
                [field, value] = match.groups
                enum.items[field] = int(value)

    def _parse_type_decl(self, app, keyword, name, children):
        """Parse type declarations."""
        typedecl = app.types[name]
        if keyword == 'type':
            type_ = typedecl.tuple
        elif keyword == 'table':
            type_ = typedecl.relation
        else:
            raise RuntimeError('keyword must be type or table, not ' + keyword)

        attrp = _AttrProcessor()

        for (line_no, end_line_no, line, grandchildren) in children:
            match = _Matcher(line)

            if match(r'^@(.*)'):
                attr = (
                    match.groups[0] +
                    (json.dumps(' '.join(x[2:] for (_, _, x, _) in grandchildren))
                     if grandchildren else ''))
                (_AttrProcessor() + attr) >> typedecl.attrs
                attrp += match.groups[0]
            elif match(r'^!type\s+(\w+)(?:\s+\[(.*)\])?$'):
                [subname, attrs] = match.groups
                subdecl = self._parse_type_decl(
                    app, 'type', name + '.' + subname, grandchildren)
                self._parse_attrs(attrs, subdecl.attrs)
            elif match(ur'!enum\s+(.*)$'):
                [subname] = match.groups
                self._parse_enum_defn(app, name + '.' + subname, grandchildren)
            elif match(ur'''
          ^
          (\w*?)  # attrname
          (?:·=·(\d+))?  # fieldid
          (?:·\s*\((?:(\d*)\.\.(\d*))\))?  # repeat
          (?:·<(?::·([^\["]*?))?)?  # typespec
          (?:\s+"([^"]*)")?  # docstring
          (?:\s+\[(.*)\]·)?  # attrs
          $
          '''):
                (attrname, fieldid, min_repeat, max_repeat, typespec, docstring, attrs
                 ) = match.groups
                if min_repeat is None:
                    min_repeat = max_repeat = 1
                elif max_repeat == '':
                    max_repeat = float('inf')

                attr_attrs = sysl_pb2.Application().attrs
                (attrp + attrs) >> attr_attrs

                if 'pk' in syslx.patterns(attr_attrs):
                    assert keyword == 'table'
                    type_.primary_key.attr_name.append(attrname)

                if attrname:
                    attrtype = type_.attr_defs[attrname]
                    if grandchildren and not typespec:
                        self._parse_type_decl(
                            app, 'type', name + '.' + attrname, grandchildren)
                        if max_repeat > 1:
                            attrtype = attrtype.list.type
                        attrtype.type_ref.ref.path.append(attrname)
                    else:
                        if typespec:
                            if max_repeat > 1:
                                attrtype = attrtype.list.type
                            self._parse_type(
                                app, line_no, name, typespec, attrtype)
                            for (line_no, end_line_no, line,
                                 ggrandchildren) in grandchildren:
                                attr = (
                                    line[1:] +
                                    (json.dumps(
                                        ' '.join(x[2:] for (_, _, x, _) in ggrandchildren))
                                        if ggrandchildren else ''))
                                (_AttrProcessor() + attr) >> attrtype.attrs
                        for (aname, a) in attr_attrs.iteritems():
                            attrtype.attrs[aname].CopyFrom(a)

                    attrtype.source_context.start.line = line_no
                    if docstring:
                        attrtype.docstring = docstring

                else:
                    if grandchildren:
                        self._parse_type_decl(
                            app, keyword, attrname, grandchildren)
                    else:
                        raise NotImplemented(
                            'name missing (and no field injection yet)')

            elif match(ur'\.\.\.$'):
                pass
            else:
                self._syntax_error('invalid type decl ' + line)

        return typedecl

    def _parse_stmts(self, field, stmts, app):  # pylint: disable=too-many-branches
        """Parse a list of statements."""
        if stmts == ['...']:
            return

        attrp = _AttrProcessor()

        for (_, _, stmt_line, children) in stmts:
            stmt = None
            match = _Matcher(stmt_line)

            # First, extract attributes from eol.
            if not stmt_line.startswith('@'):
                if match(r'^(.*?)(?:\s*\[(.*)\])?\s*$'):
                    stmt = field.stmt.add()
                    (stmt_line, attrs) = match.groups
                    (attrp + attrs) >> stmt.attrs
                    match = _Matcher(stmt_line)

            if match(r'^@(.*)'):
                attrp += match.groups[0]
                attrp >> field.attrs
            elif match(r'^(?:(\w+)\s*=\s*)?(.*)\s+<-\s+(.*?)\s*(?:\(([^)]*)\))?$'):
                (_, target, stmt.call.endpoint, args) = match.groups
                self._parse_app_name(stmt.call.target, target, app)
                if args:
                    for arg in rex.split(ur'\s*,\s*', args):
                        stmt.call.arg.add().name = arg
            elif match(r'(?i)^If\s+(.*)$'):
                (stmt.cond.test,) = match.groups
                self._parse_stmts(stmt.cond, children, app)
            elif match(r'(?i)^Return\s+(.*)$'):
                (stmt.ret.payload,) = match.groups
                assert not children
            elif match(r'(?i)^(?:Loop\s+)?(while|until)\s+(.*)$'):
                (mode, stmt.loop.criterion) = match.groups
                stmt.loop.mode = stmt.loop.Mode.Value(mode.upper())
                self._parse_stmts(stmt.loop, children, app)
            elif match(r'(?i)^Loop\s+(\d+)\s+times$'):
                (stmt.loop_n.count,) = match.groups
                self._parse_stmts(stmt.loop_n, children, app)
            elif match(r'(?i)^For\s*each\s+(.*?)\s*$'):
                (stmt.foreach.collection,) = match.groups
                self._parse_stmts(stmt.foreach, children, app)
            elif match(r'(?i)^One\s*of$'):
                for (_, _, cond, children) in children:
                    choice = stmt.alt.choice.add()
                    choice.cond = cond
                    self._parse_stmts(choice, children, app)
            elif children:
                stmt.group.title = stmt_line
                self._parse_stmts(stmt.group, children, app)
            else:
                stmt.action.action = stmt_line

    def _parse_expr(self, expr, children, field):
        field.CopyFrom(_ExprParser(expr, children)())

    def _parse_expr_block(self, app, expr_field, children):
        for (line_no, _, line, grandchildren) in children:

            # TODO: Check closing braces, don't just ignore them.
            line = line.lstrip(')')
            if not line:
                assert not grandchildren, (line, grandchildren)
                continue
            match = _Matcher(line)

            if match(
                    ur'^(.*?)·(->(?:·<·(set\s+of\b·)?([\w\.]*)·>)?)·\((\w*)$'):
                [expr, transform, setof, out_type, scopevar] = match.groups
                if not expr:
                    expr = '.'
                self._parse_transform(
                    app, expr_field, expr, bool(
                        setof), out_type, scopevar or '.',
                    grandchildren)
            elif match(r'^if\b(?:\s*(.+)\s*==\s*)?$'):
                [ifexpr] = match.groups
                if ifexpr is None:
                    ifexpr = 'true'
                self._parse_if(app, expr_field, ifexpr, grandchildren)
            else:
                assert not grandchildren, line
                self._parse_expr(line, [], expr_field)

    def _parse_transform(
            self, app, expr_field, expr, setof, out_type, scopevar, children):
        transform = expr_field.transform
        if scopevar:
            transform.scopevar = scopevar or '.'
        expr_field2 = transform.arg

        if out_type:
            target_type = sysl_pb2.Type()
            target_path = out_type.split('.')
            tr = target_type.type_ref
            tr.context.appname.CopyFrom(app.name)
            tr.ref.appname.part.extend(target_path[:1])
            tr.ref.path.extend(target_path[1:])
            if setof:
                expr_field.type.set.CopyFrom(target_type)
            else:
                expr_field.type.CopyFrom(target_type)

        self._parse_expr(expr, children, expr_field2)

        for (_, _, line, grandchildren) in children:
            match = _Matcher(line)

            if match(r'^(?:(let|table\s+of)\s+)?(\w+)\s*=\s*(.*?)$'):
                [let, name, expr] = match.groups
                stmt = transform.stmt.add()
                assign = stmt.let if let == 'let' else stmt.assign
                assign.name = name
                if rex.match(r'table\s+of$', let or ''):
                    assign.table = True
                self._parse_expr_block(
                    app, assign.expr, [(None, None, expr, grandchildren)])
            elif match(r'^\*(?:\s*-\s*(?:(\w+)|{((?:\w+\s*,\s*)*\w+)}))$'):
                [exclude1, exclude2] = match.groups
                transform.all_attrs = True
                transform.except_attrs.extend(
                    (([exclude1] if exclude1 else []) +
                     (rex.split(r'\s*,\s*', exclude2) if exclude2 else [])))
            elif match(r'^(.+)\s*\.\s*\*$'):
                [expr] = match.groups
                assert not grandchildren
                inject = transform.stmt.add().inject
                self._parse_expr(expr, [], inject)
                assert inject.HasField('call')
                tailarg = inject.call.arg.add()
                tailarg.name = 'out'
            elif match(r'^(\w*)\.(\w+)$'):
                [name, attr] = match.groups
                assign = transform.stmt.add().assign
                assign.name = attr
                get_attr = assign.expr.get_attr
                get_attr.arg.name = name or '.'
                get_attr.attr = attr
            elif match(r'^\)$'):
                pass
            else:
                raise RuntimeError('Unrecognised transform entry: ' + line)

    def _parse_if(self, app, expr_field, ifexpr, children):
        ifelse = expr_field.ifelse
        iftail = expr_field

        ifvar = sysl_pb2.Expr()
        self._parse_expr(ifexpr, children, ifvar)

        for (_, _, line, grandchildren) in children:
            match = _Matcher(line)

            # TODO: Don't match, parse.
            if match(r'^((?:[^=]|=(?!>))+)\s*=>\s*(.*?)$'):
                assert iftail
                [controls, expr] = match.groups
                for control in rex.split(r'\s*,\s*', controls):
                    if ifexpr == 'true':
                        callexpr = iftail.ifelse.cond.call
                        callexpr.func = 'bool'
                        self._parse_expr(
                            control, grandchildren, callexpr.arg.add())
                    else:
                        binexpr = iftail.ifelse.cond.binexpr
                        binexpr.op = sysl_pb2.Expr.BinExpr.EQ
                        binexpr.lhs.CopyFrom(ifvar)
                        self._parse_expr(control, grandchildren, binexpr.rhs)
                    self._parse_expr_block(
                        app, iftail.ifelse.if_true,
                        [(None, None, expr, grandchildren)])
                    iftail = iftail.ifelse.if_false
            elif match(r'^else\s*(.*?)$'):
                [expr] = match.groups
                self._parse_expr_block(
                    app, iftail, [(None, None, expr, grandchildren)])
                iftail = None
            elif match(r'^\)$'):
                pass
            else:
                raise RuntimeError('Unrecognised if-else entry: ' + line)

    PRIMITIVES = {
        p.lower()
        for p in set(sysl_pb2.Type.Primitive.keys()) - {'NO_Primitive'}
    }

    def _parseGlobalRef(self, ref, typ, app):

        if not ref:
            return

        if ref in self.PRIMITIVES:
            typ.primitive = sysl_pb2.Type.Primitive.Value(ref.upper())
            return

        m = rex.match(r'set\s+of\s+([\w\.]+)$', ref)
        if m:
            ref = m.group(1)
            typ = typ.set

        path = ref.split('.')

        # set up the type ref
        sref = typ.type_ref
        # TODO: Fill in scopedRef.context
        sref.ref.appname.part.extend(rex.split(r'\s*::\s*', path[0]))
        sref.ref.path.extend(path[1:])

    def _parse_view(self, parent, line, children):
        """Parse a view specification."""

        m = rex.match(r'''(?x)
      ^!view\s+
      ([\w${}]+)
      (?:\s*\(\s*(.*?)\s*\))?
      (?:\s*->\s*((?:set\s+of\s+)?[\w\.]+)(\??))?
      (?:\s+\[(.*?)\])?  # attrs
      $
      ''',
                      line)
        if not m:
            print '>>>', line
            raise RuntimeError('Unmatched view syntax')
        (name, params, ret_type, optional, attrs) = m.groups()
        view = parent.views[name]

        if params:
            for param in rex.split(r'\s*,\s*', params):
                if param.count('<:') != 1:
                    raise RuntimeError('Parameter syntax error')
                (pname, ptname) = rex.split(r'\s*<:\s*', param)
                param = view.param.add()
                param.name = pname
                self._parseGlobalRef(ptname, param.type, parent)

        self._parse_attrs(attrs, view.attrs)

        if ret_type:
            self._parseGlobalRef(ret_type, view.ret_type, parent)

        if 'abstract' not in syslx.patterns(view.attrs):
            self._parse_expr_block(parent, view.expr, children)

    def _parse_app(self, module, app, name, long_name, attrs, endpoints):
        """Parse an app."""
        if app is None:
            app = module.apps[name]

        self._parse_app_name(app.name, name, app)

        if long_name is not None:
            app.long_name = long_name

        if attrs:
            attrs >> app.attrs

        for (_, _, endpoint_line, stmts) in endpoints:
            if endpoint_line.endswith(': ...'):
                endpoint_line = endpoint_line[:-5]

            match = _Matcher(endpoint_line)

            attrp = _AttrProcessor()

            if match(r'^@(.*)'):
                (_AttrProcessor() + match.groups[0]) >> app.attrs
                attrp += match.groups[0]
                stmts = []
            elif match(r'^(/\S*)(?:\s*\[(.*)\])?$'):
                (path, extra_attrs) = match.groups
                self._parse_rest_api(
                    app, '', path, attrp + extra_attrs, {}, stmts)
                stmts = []
            elif match(r'^!(type|table)\s+(\w+)(?:\s+\[(.*)\])?$'):
                [keyword, name, attrs] = match.groups
                subdecl = self._parse_type_decl(app, keyword, name, stmts)
                self._parse_attrs(attrs, subdecl.attrs)
                stmts = []
            elif match(r'^!enum\s+(.*)$'):
                [name] = match.groups
                self._parse_enum_defn(app, name, stmts)
                stmts = []
            elif match(r'^!view\b'):
                self._parse_view(app, endpoint_line, stmts)
                stmts = []
            elif match(r'^(?:-\|>)\s+(.*)$'):
                [name] = match.groups
                self._parse_app(module, app.mixin2.add(),
                                name, None, None, stmts)
                stmts = []
            elif match(r'^!wrap\s+(.*)$'):
                [name] = match.groups
                assert not app.HasField('wrapped'), "can only wrap one app"
                self._parse_app(module, app.wrapped, name, None, None, stmts)
                stmts = []
            elif match(r'^(.*?)(\s+->\s+)(.*?)(?:\s+\[(.*)\])?$'):
                [source, op, name, child_attrs] = match.groups
                endpt = app.endpoints[source + op + name]
                endpt.name = source + op + name
                self._parse_app_name(endpt.source, source, app)
                (attrp + child_attrs) >> endpt.attrs
            elif match(r'''
          (<->\ )?
          (.*?)
          (?:\s*\(([^)]*)\))?
          (?:\s+"(.*)")?
          (?:\s+\[(.*)\])?
          $
          '''):
                (pubsub, name, param_string, long_name, child_attrs) = match.groups

                # TODO: grok params
                endpt = app.endpoints[name]
                endpt.name = name
                if param_string:
                    for param_pair in rex.split(r'\s*,\s*', param_string):
                        param = endpt.param.add()
                        if param_pair.count('<:') == 1:
                            (pname, ptname) = rex.split(
                                r'\s*<:\s*', param_pair)
                            param.name = pname
                            self._parseGlobalRef(ptname, param.type, app)
                        else:
                            param.name = param_pair
                            param.type.no_type.SetInParent()

                if long_name is not None:
                    endpt.long_name = long_name

                (attrp + child_attrs) >> endpt.attrs

                endpt.is_pubsub = bool(pubsub)
            else:
                assert False, "Should be unreachable! " + endpoint_line

            if stmts:
                self._parse_stmts(endpt, stmts, app)

    def parse(self, src, path, module):
        """Parse source into a sysl_pb2.Module."""
        imports = set()
        attrp = _AttrProcessor()

        # with self.dynamic_scope(context=source_context((0, None))
        for (line_no, end_line_no, app_line, endpoints) in (
                self._group_indent(src, path)):
            try:
                if endpoints:
                    match = _Matcher(app_line)

                    if match(
                            r'^(.*?)\s*(\([^)]*\))?(?:\s+"(.*)")?(?:\s*\[(.*)\])?$'):
                        (name, _, long_name, attrs) = match.groups
                        # print name, long_name
                        self._parse_app(module, None, name, long_name,
                                        attrp + attrs, endpoints)

                elif app_line.startswith('@'):
                    raise RuntimeError("@... not allowed at top level")

                else:
                    assert app_line.startswith('import '), app_line
                    imports.add(app_line[7:])
            except BaseException:
                print >>sys.stderr, "Error @ {}".format(line_no)
                raise

        return imports
