#!/usr/bin/env python
# -*-: encoding: utf-8 -*-
"""sysl.py

Sysl compiler and toolkit.
"""

import argparse
import cStringIO
import errno
import os
import sys

import src.util.debug as debug

from sysl.sysl import sysldata
from sysl.sysl import syslints
from sysl.sysl import syslloader
from sysl.sysl import syslseqs


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

  # Ensure the output directory exists.
  if 'output' in args:
    try:
      os.makedirs(os.path.dirname(args.output))
    except OSError as exc:
      if exc.errno != errno.EEXIST or os.path.isdir(args.output):
        raise

  args.func(args)


if __name__ == '__main__':
  debug.init()
  main()
