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

"""Install exception hook for easier debugging."""

import logging
import platform
import sys


HAS_ANSI = platform.system() != 'Windows'


def ansi(code):
  """Conditionally return ANSI-escaped codes."""
  return '\033[' + code if HAS_ANSI else code


def color256(r, g, b):
  assert {r, g, b} <= set(range(6))

  # https://en.wikipedia.org/wiki/ANSI_escape_code#Colors
  code = 16 + 36 * r + 6 * g + b
  return '\033[38;5;{}m'.format(code)


def _hook(type_, value, tback):
  """Exception hook callback."""
  if hasattr(sys, 'ps1') or not sys.stderr.isatty():
  # we are in interactive mode or we don't have a tty-like
  # device, so we call the default hook
    sys.__excepthook__(type_, value, tback)
  else:
    import traceback, pdb
    # we are NOT in interactive mode, print the exception...
    traceback.print_exception(type_, value, tback)

    # Dirty hack because Py27 doesn't chain exceptions
    if value.args:
      tb2 = value.args[-1]
      if isinstance(tb2, type(tback)):
        ex = value.args[-2]
        print >>sys.stderr, '{}Caused by{} '.format(ansi('1;35m'), ansi('0m')),
        traceback.print_exception(type_(ex), ex, tb2)

      print
    # ...then start the debugger in post-mortem mode.
    # pdb.pm() # deprecated
    pdb.post_mortem(tback) # more "modern"


sys.excepthook = _hook


if not HAS_ANSI:
  # Try to fake ANSI escapes.
  try:
    # pylint: disable=import-error
    import colorama
    colorama.init()
    HAS_ANSI = True
  except:
    try:
      # pylint: disable=import-error, unused-import
      import tendo.ansiterm
    except:
      pass


def init():
  """Just gives main modules an excuse to reference debug."""
  logging.basicConfig(
    format="%(levelname)s:%(name)s %(pathname)s:%(lineno)d %(message)s")
