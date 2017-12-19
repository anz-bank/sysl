# -*- encoding: utf-8 -*-
import re


I = re.I
IGNORECASE = re.IGNORECASE

L = re.L
LOCALE = re.LOCALE

M = re.M
MULTILINE = re.MULTILINE

S = re.S
DOTALL = re.DOTALL

U = re.U
UNICODE = re.UNICODE

X = re.X
VERBOSE = re.VERBOSE


CACHE = {}


def cache(pattern, flags):
    flags |= VERBOSE
    key = (pattern, flags)
    c = CACHE.get(key)
    if c is None:
        c = CACHE[key] = (
            re.compile(pattern.replace(u'·', ur'\s*').replace(u'•', ur'\s+'), flags))
    return c


def cached(f):
    def g(pattern, string, flags=0, *args, **kwargs):
        return f(cache(pattern, flags), string, *args, **kwargs)
    return g


_cre = type(re.compile(''))

search = cached(_cre.search)
match = cached(_cre.match)
split = cached(_cre.split)
findall = cached(_cre.findall)
finditer = cached(_cre.finditer)


def sub(pattern, repl, string, count=0, flags=0):
    return cache(pattern, flags).sub(repl, string, count=count)


def subn(pattern, repl, string, count=0, flags=0):
    return cache(pattern, flags).subn(repl, string, count=count)


def escape(string):
    return re.escape(string)


def purge():
    re.purge()


error = re.error
