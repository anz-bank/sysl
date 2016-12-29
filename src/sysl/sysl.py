#!/usr/bin/env python
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

"""sysl.py

Sysl compiler and toolkit.
"""

import argparse
import cStringIO
import sys

import src.util.debug as debug

from src.sysl import sysldata
from src.sysl import syslints
from src.sysl import syslloader
from src.sysl import syslseqs


def _pb_sub_parser(subparser):
  """Setup proto subcommand."""
  argp = subparser.add_parser('pb')

  def cmd(args):
    """Handle subcommand."""
    (module, _, _) = syslloader.load(args.modules, args.validations, args.root)
    out = module.SerializeToString()

    (open(args.output, 'wb') if args.output else sys.stdout).write(out)

  argp.set_defaults(func=cmd)

  argp.add_argument('--output', '-o',
            help='output file')
  argp.add_argument('modules', nargs='+',
            help='modules')


def _textpb_sub_parser(subparser):
  """Setup proto subcommand."""
  argp = subparser.add_parser('textpb')

  def cmd(args):
    """Handle subcommand."""
    (module, _, _) = syslloader.load(args.modules, args.validations, args.root)
    out = str(module)

    (open(args.output, 'w') if args.output else sys.stdout).write(out)

  argp.set_defaults(func=cmd)

  argp.add_argument('--output', '-o',
            help='output file')
  argp.add_argument('modules', nargs='+',
            help='modules')


def _deps_sub_parser(subparser):
  """Setup deps subcommand."""
  argp = subparser.add_parser('deps')

  def cmd(args):
    """Handle subcommand."""
    out = cStringIO.StringIO()
    fmt = args.target + ' : {}\n'

    for module in args.modules:
      (_, _, modules) = syslloader.load([module], args.validations, args.root)

      print >>out, fmt.format(module, ' '.join(m + '.sysl' for m in modules))

    (open(args.output, 'w') if args.output else sys.stdout).write(out.getvalue())


  argp.set_defaults(func=cmd)

  argp.add_argument('--target', '-t', default='{}',
            help='format string for target spec')
  argp.add_argument('--output', '-o',
            help='output file')
  argp.add_argument('modules', nargs='+',
            help='modules')


def main():
  """Main function."""
  argp = argparse.ArgumentParser(
    description='System Modelling Language Toolkit')

  argp.add_argument('--no-validations', '--nv', dest='validations',
            action='store_false', default=True,
            help='suppress validations')
  argp.add_argument('--root', '-r', default='.',
            help='sysl system root directory')

  subparser = argp.add_subparsers(help='sub-commands')
  _pb_sub_parser(subparser)
  _textpb_sub_parser(subparser)
  _deps_sub_parser(subparser)
  sysldata.add_subparser(subparser)
  syslints.add_subparser(subparser)
  syslseqs.add_subparser(subparser)

  args = argp.parse_args()
  args.func(args)


if __name__ == '__main__':
  debug.init()
  main()
