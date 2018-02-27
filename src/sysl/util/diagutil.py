"""Diagramming utilities"""

import collections
import cStringIO
import itertools
import os
import re
import sys

import plantuml
import requests

import cache
import confluence
import simple_parser


def group_by(src, key):
    """Apply sorting and grouping in a single operation."""
    return itertools.groupby(sorted(src, key=key), key=key)


def fmt_app_name(appname):
    """Format an app name as a single string."""
    return ' :: '.join(appname.part)


OutputArgs = collections.namedtuple('OutputArgs',
                                    'output plantuml verbose expire_cache dry_run')


def output_plantuml(args, puml_input):
    """Output a PlantUML diagram."""
    ext = os.path.splitext(args.output or '')[-1][1:]
    SUPPORTED_MODES = {'png': 'img', 'svg': 'svg', 'uml': None}
    if ext not in SUPPORTED_MODES:
        raise Exception('Extension "{}" not supported. Valid extensions: {}.'.format(
            ext, ', '.join(SUPPORTED_MODES)))
    mode = {'png': 'img', 'svg': 'svg', 'uml': None, '': None}[ext]
    server = (args.plantuml or
              os.getenv('SYSL_PLANTUML', 'http://localhost:8080/plantuml'))
    if mode:
        def calc():
            data = plantuml.deflate_and_encode(puml_input)
            response = requests.get('{}/{}/{}'.format(server, mode, data))
            response.raise_for_status()
            return response.content
        out = cache.get(mode + ':' + puml_input, calc)

    useConfluence = args.output.startswith('confluence://')

    if args.verbose:
        print args.output + '...' * useConfluence,
        sys.stdout.flush()

    if useConfluence:
        if confluence.upload_attachment(
            args.output, cStringIO.StringIO(
                out), args.expire_cache, args.dry_run
        ) is None:
            if args.verbose:
                print '\033[1;30m(no change)\033[0m',
        else:
            if args.verbose:
                print '\033[1;32muploaded\033[0m',
                if args.dry_run:
                    print '... not really (dry-run)',
    else:
        (open(args.output, 'wb') if args.output else sys.stdout).write(out)
        # Uncomment this to print out Plant UML
        #(open(args.output + '.puml', 'w') if args.output else sys.stdout).write(puml_input)

    if args.verbose:
        print


class VarManager(object):
    """Synthesise a mapping from names to variables.

    This class is used to map arbitrary names, which may not be valid in some
    syntactic contexts, to more uniform names.
    """

    def __init__(self, newvar):
        self._symbols = {}
        self._newvar = newvar

    def __call__(self, name):
        """Return a variable name for a given name.

        Make sure the same name always maps to the same variable."""
        if name in self._symbols:
            return self._symbols[name]

        var = '_{}'.format(len(self._symbols))
        self._newvar(var, name)
        self._symbols[name] = var
        return var


class _FmtParser(simple_parser.SimpleParser):
    """Parse format strings used in project .sysl files."""
    # TODO: Document the format string sublanguage.

    def parse(self):
        """Top-level parse function."""
        if self.expansions():
            code = 'lambda **vars: ' + self.pop()
            return eval(code)  # pylint: disable=eval-used

    def expansions(self, term=u'$'):
        """Parse expansions."""

        result = [repr(u'')]
        while self.eat(ur'((?:[^%]|%[^(\n]|\n)*?)(?=' + term + ur'|%\()'):
            prefix = self.pop()
            prefix = re.sub(
                u'%(.)', ur'\1', prefix.replace(u'%%', u'\1')
            ).replace(u'\1', u'%')
            if prefix:
                result.append(repr(prefix))

            if self.eat(ur'%\('):
                if not self.eat(ur'(@?\w+)'):
                    raise Exception('missing variable reference')
                var = cond = u"vars.get({!r}, '')".format(self.pop())

                # conditionals!
                if self.eat(ur'([!=]=)'):
                    cond = var + " {} ".format(self.pop())
                    if not self.eat(ur'\'([\w ]+)\''):
                        raise Exception('missing conditional value')
                    cond = cond + u"{!r}".format(self.pop())

                if self.eat(ur'~/([^/]+)/'):
                    cond = u're.search({!r}, {})'.format(
                        self.pop().replace('\b', r'\b'), var)

                have = self.eat(ur'[=?]')
                if have:
                    if not self.expansions(ur'$|[|)]'):
                        raise Exception('wat?')
                    have = self.pop()

                have_not = self.eat(ur'\|')
                if have_not:
                    if not self.expansions(ur'$|\)'):
                        raise Exception('wat?')
                    have_not = self.pop()

                if not self.eat(ur'\)'):
                    raise Exception('unclosed expansion')

                result.append(u"({} if {} else {})".format(
                    have or var, cond, have_not or repr('')))
            else:
                self.push(u'(' + u' + '.join(result) + u')')
                return True


def parse_fmt(text):
    """Parse a format string."""
    return _FmtParser(text)()


def attr_fmt_vars(*attrses, **kwargs):
    """Return a dict based attrs that is suitable for use in parse_fmt()."""
    fmt_vars = {}

    for attrs in attrses:
        if type(attrs).__name__ in ['MessageMap', 'MessageMapContainer']:
            for (name, attr) in attrs.iteritems():
                if attr.WhichOneof('attribute'):
                    fmt_vars['@' +
                             name] = getattr(attr, attr.WhichOneof('attribute'))
                else:
                    fmt_vars['@' + name] = ''
        else:
            fmt_vars.update(attrs)

    fmt_vars.update(kwargs)

    return fmt_vars
