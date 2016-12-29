# -*- encoding: utf-8 -*-

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

import re


I           = re.I
IGNORECASE  = re.IGNORECASE

L           = re.L
LOCALE      = re.LOCALE

M           = re.M
MULTILINE   = re.MULTILINE

S           = re.S
DOTALL      = re.DOTALL

U           = re.U
UNICODE     = re.UNICODE

X           = re.X
VERBOSE     = re.VERBOSE


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
