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
from argparse import RawTextHelpFormatter

# sysl.util.debug as debug
from sysl.core import sysldata
from sysl.core import syslints
from sysl.core import syslloader
from sysl.core import syslseqs
from sysl.util.argparse import add_modules_option, add_output_option
from sysl.__version__ import __version__


def _pb_sub_parser(subparser):
    """Setup proto subcommand."""
    argp = subparser.add_parser('pb',
                                description='Create binary protobuf output',
                                formatter_class=RawTextHelpFormatter)

    def cmd(args):
        """Handle subcommand."""
        (module, _, _) = syslloader.load(
            args.modules, args.validations, args.root)
        out = module.SerializeToString()

        (open(args.output, 'wb') if args.output else sys.stdout).write(out)

    argp.set_defaults(func=cmd)

    add_modules_option(argp)
    add_output_option(argp)


def _textpb_sub_parser(subparser):
    """Setup proto subcommand."""
    argp = subparser.add_parser('textpb',
                                description='Create text protobuf output',
                                formatter_class=RawTextHelpFormatter)

    def cmd(args):
        """Handle subcommand."""
        (module, _, _) = syslloader.load(
            args.modules, args.validations, args.root)
        out = str(module)

        (open(args.output, 'w') if args.output else sys.stdout).write(out)

    argp.set_defaults(func=cmd)

    add_modules_option(argp)
    add_output_option(argp)


def _deps_sub_parser(subparser):
    """Setup deps subcommand."""
    argp = subparser.add_parser('deps',
                                description='Create module dependency output',
                                formatter_class=RawTextHelpFormatter)

    def cmd(args):
        """Handle subcommand."""
        out = cStringIO.StringIO()
        fmt = args.target + ' : {}\n'

        for module in args.modules:
            (_, _, modules) = syslloader.load(
                [module], args.validations, args.root)

            print >>out, fmt.format(
                module, ' '.join(m + '.sysl' for m in modules))

        (open(args.output, 'w') if args.output else sys.stdout).write(out.getvalue())

    argp.set_defaults(func=cmd)

    argp.add_argument('--target', '-t', default='{}',
                      help='format string for target spec')
    add_output_option(argp)
    add_modules_option(argp)


def main(input_args=sys.argv[1:]):
    """Main function."""
    argp = argparse.ArgumentParser(
        description='System Modelling Language Toolkit',
        formatter_class=RawTextHelpFormatter)

    argp.add_argument('--no-validations', '--nv', dest='validations',
                      action='store_false', default=True,
                      help='suppress validations')
    argp.add_argument('--root', '-r', default='.',
                      help='sysl root directory for input files (default: .)')
    argp.add_argument('--version', '-v',
                      help='show version number (semver.org standard)',
                      action='version', version='%(prog)s ' + __version__)
    argp.add_argument('--trace', '-t',
                      action='store_true', default=False)

    subparser = argp.add_subparsers(help='\n'.join(['sub-commands',
                                                    'more help with: sysl <sub-command> --help', 'eg: sysl pb --help']))
    _pb_sub_parser(subparser)
    _textpb_sub_parser(subparser)
    _deps_sub_parser(subparser)
    sysldata.add_subparser(subparser)
    syslints.add_subparser(subparser)
    syslseqs.add_subparser(subparser)

    args = argp.parse_args(input_args)

    # Ensure the output directory exists.
    if 'output' in args and os.path.dirname(args.output):
        try:
            os.makedirs(os.path.dirname(args.output))
        except OSError as exc:
            if exc.errno != errno.EEXIST or os.path.isdir(args.output):
                raise
    try:
        args.func(args)
    except Exception as e:
        if args.trace:
            raise
        else:
            print e
            sys.exit(1)


if __name__ == '__main__':
    main()
    sys.exit(0)
