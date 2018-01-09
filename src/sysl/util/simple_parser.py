# -*-: encoding: utf-8 -*-
"""Parsing helper class."""

import abc
import operator
import sys

from sysl.proto import sysl_pb2
from sysl.util import rex


class SyntaxError(Exception):
    """Exception generated when a syntax error is detected."""

    def __init__(self, parser, fmt, *args, **kwargs):
        source_context = parser.source_context((None, None))
        super(Exception, self).__init__(
            fmt.format(*args, **kwargs), source_context)
        self.source_context = source_context


class SimpleParser(object):
    """Parsing helper class."""

    __metaclass__ = abc.ABCMeta

    text = property(operator.itemgetter('_text'))
    cur = property(operator.itemgetter('_cur'))
    latest_context = property(operator.itemgetter('_latest_context'))

    def __init__(self, text, source_context=None):
        self._text = text
        self._source_context = source_context

        self._cur = 0
        self._stk = []
        self._regexes = {}
        self._latest_context = source_context

        self.skip_ws = False

    def __repr__(self):
        return u'{}\U0001f595 {}'.format(
            self._text[:self._cur], self._text[self._cur:]).encode('utf-8')

    def __getitem__(self, index):
        return self._stk[index]

    def __delitem__(self, index):
        del self._stk[index]

    def __nonzero__(self):
        return bool(self._stk)

    def __call__(self):
        """Invoke the top-level parse method."""
        try:
            result = self.parse()
            if self._cur < len(self._text.rstrip()):
                raise RuntimeError('input not consumed')
            return result
        except BaseException:
            print >>sys.stderr, "Error @ {!s}🔥 {!s}".format(
                self._text[:self._cur], self._text[self._cur:])
            raise

    def _syntax_error(self, fmt, *args, **kwargs):
        """Handle a syntax error.

        Passes the error to self.report_error, if present, or else raises it.
        """
        err = SyntaxError(self, fmt, *args, **kwargs)
        if hasattr(self, 'report_error'):
            self.report_error(err)
        else:
            raise err

    @abc.abstractmethod
    def parse(self):
        """Top-level parse function. Must be overridden."""
        pass

    def __nonzero__(self):
        return bool(self._text[self._cur:].strip())

    def _regex(self, pat):
        return rex.compile(pat)

    def _eat(self, pat, action=None):
        """Try to eat a pattern, returning Boolean success/fail."""
        m = rex.match(pat, self._text, pos=self._cur)
        if m:
            if action:
                self._stk.append(action(*m.groups()))
            else:
                self._stk.extend(m.groups())
            self._cur += len(m.group(0))
        return bool(m)

    def eat(self, pat, action=None):
        if self.skip_ws:
            pat = r'\s*' + pat
        return self._eat(pat, action)

    def expect(self, pat):
        if not self.eat(pat):
            raise RuntimeError('Syntax error')
        return self

    def push(self, *values):
        """Push terms onto the parse stack."""
        self._stk.extend(values)
        return values

    def pop(self, count=None):
        """Pop a term off the parse stack and return it."""
        if count is None:
            if not self._stk:
                raise Exception('Cannot pop empty parse stack!')
            return self._stk.pop()
        assert count <= len(self._stk)
        result = self._stk[-count:]
        del self[-count:]
        return result

    def top(self):
        if not self._stk:
            raise Exception('No top for empty parse stack!')
        return self._stk[-1]
