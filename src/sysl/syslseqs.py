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

"""Generate sequence diagrams from sysl modules."""

import collections
import contextlib
import fnmatch
import itertools
import os
import re
import textwrap

from src.sysl import syslalgo
from src.sysl import syslloader
from src.sysl import syslx

from src.util import diagutil
from src.util import writer

import pdb

class _Writer(writer.Writer):
  """Output logic helper."""

  def __init__(self, activations):
    super(_Writer, self).__init__('plantuml')
    self._active = collections.defaultdict(int)
    self._activations = activations

  @contextlib.contextmanager
  def indent(self):
    """Temporarily increase indent level."""
    self._indent += 1
    yield
    self._indent -= 1

  def activate(self, agent):
    """Increase activation level for agent."""
    if self._activations:
      self('activate {}', agent)
    self._active[agent] += 1

  def deactivate(self, agent):
    """Decrease activation level for agent."""
    # TODO: Reinstate this after figuring out why it fires.
    #assert self._active[agent], agent
    if not self._active[agent]:
      return
    self._active[agent] -= 1
    if self._activations:
      self('deactivate {}', agent)

  def deactivate_all(self):
    """Set all activation levels to 0."""
    for (agent, active) in self._active.iteritems():
      for _ in range(active):
        self.deactivate(agent)

  @contextlib.contextmanager
  def activated(self, agent, suppressed=False):
    """Temporarily increase activation level."""
    def deactivate():
      """Callback to decrease level. Allows early deactivation."""
      if active[0]:
        active[0] = False
        self.deactivate(agent)

    if not suppressed:
      self.activate(agent)
    active = [not suppressed]

    yield deactivate

    deactivate()


SequenceDiagParams = collections.namedtuple(
  'SequenceDiagParams', 'endpoints epfmt appfmt activations title blackboxes')


def sequence_diag(module, params, log_integration=None):
  """Generate a sequence diagram.

  Params:
    module: sysl_pb2.Module
    params: SequenceDiagParams
    log_call: invoked for each visited call.
      def(app=sysl_pb2.Application(),
        stmt=sysl_pb2.Stmt(),
        patterns={str})
  """
  blackboxes = params.blackboxes or []
  already_visited = collections.defaultdict(int)

  write = _Writer(params.activations)
  var_names = []

  def new_var(var, appname):
    """Outputs a new definition when VarManager makes a new variable."""
    app = module.apps[appname]
    has_category = syslx.patterns(app.attrs) & {'human', 'cron', 'db', 'external', 'ui'}
    assert len(has_category) <= 1
    (order, agent) = {
      'human': (0, 'actor'),
      'ui' : (1, 'boundary'),
      'cron': (2, 'control'),
      'db': (4, 'database'),
      'external': (5, 'control'),
    }.get(
      ''.join(has_category),
      (3, 'control'))
    label = params.appfmt(
      appname=appname,
      **diagutil.attr_fmt_vars(app.attrs)
      ).replace(u'\n', ur'\n')
    var_names.append(((order, int(var[1:])), u'{} "{}" as {}'.format(agent, label, var)))

  var_name = diagutil.VarManager(new_var)

  def visit_endpoint(
      from_app,
      appname,
      epname,
      uptos,
      sender_patterns,
      stmt=None,
      deactivate=None):
    """Recursively visit an endpoint."""
    if from_app:
      sender = var_name(syslx.fmt_app_name(from_app.name))
    else:
      sender = '['
    agent = var_name(appname)
    app = module.apps.get(appname)
    endpt = app.endpoints.get(epname)
    assert endpt

    def visit_stmts(stmts, deactivate, last_parent_stmt):
      """Recursively visit a stmt list."""
      def block(last_stmt, block_stmts, fmt, *args):
        """Output a compound block."""
        write(fmt, *args)
        with write.indent():
          return visit_stmts(block_stmts, deactivate, last_stmt)

      def block_with_end(last_stmt, block_stmts, fmt, *args):
        """Output a compound block, including the 'end' clause."""
        payload = block(last_stmt, block_stmts, fmt, *args)
        write('end')
        return payload

      payload = None

      for (i, stmt) in enumerate(stmts):
        last_stmt = last_parent_stmt and i == len(stmts) - 1
        if stmt.HasField('call'):
          with write.indent():
            payload = visit_endpoint(
              app,
              syslx.fmt_app_name(stmt.call.target),
              stmt.call.endpoint,
              uptos,
              app_patterns,
              stmt,
              last_stmt and deactivate)
        elif stmt.HasField('action'):
          write('{0} -> {0} : {1}', agent, r'\n'.join(textwrap.wrap(stmt.action.action, 40)))
        elif stmt.HasField('cond'):
          payload = block_with_end(last_stmt, stmt.cond.stmt,
                       'opt {}',
                       stmt.cond.test)
        elif stmt.HasField('loop'):
          payload = block_with_end(last_stmt, stmt.loop.stmt,
                       'loop {} {}',
                       stmt.loop.Mode.Name(stmt.loop.mode),
                       stmt.loop.criterion)
        elif stmt.HasField('loop_n'):
          payload = block_with_end(last_stmt, stmt.loop_n.stmt,
                       'loop {} times',
                       stmt.loop_n.count)
        elif stmt.HasField('foreach'):
          payload = block_with_end(last_stmt, stmt.foreach.stmt,
                       'loop for each {}',
                       stmt.foreach.collection)
        elif stmt.HasField('group'):
          payload = block_with_end(last_stmt, stmt.group.stmt,
                       'group {}',
                       stmt.group.title)
        elif stmt.HasField('alt'):
          prefix = 'alt'
          for (j, choice) in enumerate(stmt.alt.choice):
            last_alt_stmt = last_stmt and j == len(stmt.alt.choice) - 1
            payload = block(last_alt_stmt, choice.stmt,
                    '{} {}', prefix, choice.cond)
            prefix = 'else'
          write('end')
        elif stmt.HasField('ret'):
          write('{}<--{} : {}', sender, agent, stmt.ret.payload)
        else:
          raise Exception('No statement!')

      return payload

    app_patterns = syslx.patterns(app.attrs)
    target_patterns = syslx.patterns(endpt.attrs)

    patterns = target_patterns

    human = 'human' in app_patterns
    human_sender = 'human' in sender_patterns
    cron = 'cron' in sender_patterns
    needs_int = not (human or human_sender or cron) and sender != agent
    # pdb.set_trace()
    label = re.sub(ur'^.*? -> ', u' ⬄ ', r'\n'.join(textwrap.wrap(unicode(epname), 40)))

    cron = 'cron' in app_patterns

    if stmt:
      assert stmt.HasField('call')

      # pdb.set_trace()
      patterns = syslx.patterns(stmt.attrs) or patterns


      label = params.epfmt(
        epname=label,
        human='human' if human else '',
        human_sender='human sender' if human_sender else '',
        needs_int='needs_int' if needs_int else '',
        args=', '.join(p.name for p in stmt.call.arg),
        patterns=', '.join(sorted(patterns)),
        **diagutil.attr_fmt_vars(stmt.attrs)
        ).replace('\n', r'\n')

    if not ((human and sender == '[') or cron):
      ep_patterns = syslx.patterns(endpt.attrs)

      icon = '<&timer> ' if 'cron' in ep_patterns else ''
      # pdb.set_trace()
      write('{}->{} : {}{}', sender, agent, icon, label)
      if log_integration and stmt:
        log_integration(app=from_app, stmt=stmt, patterns=patterns)

    payload = syslalgo.return_payload(endpt.stmt)
    calling_self = from_app and syslx.fmt_app_name(from_app.name) == appname
    if not calling_self and not payload and deactivate:
      deactivate()

    if len(endpt.stmt):
      hit_blackbox = False
      for (upto, comment) in itertools.chain(uptos.iteritems(), already_visited.keys()):
        # Compare the common prefix of the current endpt and upto.
        upto_parts = upto.split(' <- ')
        if [appname, epname][:len(upto_parts)] == upto_parts:
          hit_blackbox = True
          if payload:
            write.activate(agent)
            if comment is not None:
              write('note over {}: {}', agent, comment or 'see below')
          else:
            if comment is not None:
              write('note {}: {}',
                  'left' if sender > agent else 'right',
                  comment or 'see below')
          if payload:
            write('{}<--{} : {}', sender, agent, payload)
            write.deactivate(agent)
          break

      if not hit_blackbox:
        with write.activated(agent, human or cron) as deactivate:
          visiting = (appname + ' <- ' + epname, None) #'see above')
          already_visited[visiting] += 1
          try:
            return visit_stmts(endpt.stmt, deactivate, True)
          finally:
            already_visited[visiting] -= 1
            if not already_visited[visiting]:
              del already_visited[visiting]

  with write.uml():
    #write('scale max 8192 height')
    if params.title:
      write('title {}', params.title)

    app_eps = [re.match(r'(.*?)\s*<-\s*(.*?)(?:\s*\[upto\s+(.*)\])*$', endpt).groups()
           for endpt in params.endpoints]

    # Treat each endpoint as a blackbox for the other endpoints.
    uptos = {appname + ' <- ' + epname for (appname, epname, _) in app_eps}

    # Global blackboxes
    blackboxes = {app: comment
            for (app, comment) in blackboxes}

    for (appname, epname, upto) in app_eps:
      write('== {} <- {} ==', appname, epname)
      bbs = blackboxes.copy()
      for bbox in ({upto} if upto else set()) | uptos - {appname + ' <- ' + epname}:
        bbs[bbox] = 'see below'
      already_visited.clear()
      visit_endpoint(None, appname, epname, bbs, [])
      write.deactivate_all()

    for (_, var) in sorted(var_names):
      write.head('{}', var)

  return str(write)




def add_subparser(subp):
  """Setup seqdiags subcommand."""
  argp = subp.add_parser('sd')

  def cmd(args):
    """Handle subcommand."""
    (module, _, _) = syslloader.load(args.modules, args.validations, args.root)

    def output_sd(args, params):
      """Generate and output a sequence diagram."""
      # pdb.set_trace()
      out = sequence_diag(module, params)

      diagutil.output_plantuml(args, out)

    epfilters = os.getenv('SYSL_SD_FILTERS', '*').split(',')

    # TODO: Find a cleaner way to trigger multi-output.
    if '%(epname)' in args.output:
      out_fmt = diagutil.parse_fmt(args.output)
      for appname in args.app:
        app = module.apps[appname]

        bbs = [[e.s for e in bbox.a.elt]
             for bbox in syslx.View(app.attrs)['blackboxes'].a.elt]

        seqtitle = diagutil.parse_fmt(syslx.View(app.attrs)['seqtitle'].s or args.seqtitle)
        epfmt = diagutil.parse_fmt(
          syslx.View(app.attrs)['epfmt'].s or args.endpoint_format)
        appfmt = diagutil.parse_fmt(
          syslx.View(app.attrs)['appfmt'].s or args.app_format)


        for (name, endpt) in sorted(app.endpoints.iteritems(), key=lambda kv: kv[1].name):
          if not any(fnmatch.fnmatch(name, filt) for filt in epfilters):
            continue

          attrs = {u'@' + name: value.s
               for (name, value) in endpt.attrs.iteritems()}
          args.output = out_fmt(
            appname=appname,
            epname=name,
            eplongname=endpt.long_name,
            **attrs)

          bbs2 = [[e.s for e in bbox.a.elt]
              for bbox in syslx.View(endpt.attrs)['blackboxes'].a.elt]

          varrefs = diagutil.attr_fmt_vars(
            app.attrs,
            endpt.attrs,
            epname=endpt.name,
            eplongname=endpt.long_name)

          out = sequence_diag(module, SequenceDiagParams(
            endpoints=[' :: '.join(s.call.target.part) + ' <- ' + s.call.endpoint
                   for s in endpt.stmt
                   if s.WhichOneof('stmt') == 'call'],
            epfmt=epfmt,
            appfmt=appfmt,
            activations=args.activations,
            title=seqtitle(**varrefs).replace('\n', r'\n'),
            blackboxes=bbs + bbs2))
          diagutil.output_plantuml(args, out)

    else:
      out = sequence_diag(module, SequenceDiagParams(
        endpoints=[i[0] for i in args.endpoint],  # -s builds list of lists (idkw).
        epfmt=diagutil.parse_fmt(args.endpoint_format),
        appfmt=diagutil.parse_fmt(args.app_format),
        activations=args.activations,
        title=args.title,
        blackboxes=args.blackboxes))
      diagutil.output_plantuml(args, out)

  argp.set_defaults(func=cmd)

  # Sequence diagrams
  argp.add_argument(
    '--endpoint', '-s', action='append',
    help='Include endpoint in sequence diagram.')
  argp.add_argument(
    '--app', '-a', action='append',
    help=('Include all endpoints for app in sequence diagram (currently '
        'only works with templated --output). Use SYSL_SD_FILTERS env (a '
        'comma-list of shell globs) to limit the diagrams generated'))
  argp.add_argument(
    '--no-activations', '--na', dest='activations', action='store_false', default=True,
    help='Suppress sequence diagram activation bars.')
  argp.add_argument(
    '--endpoint_format', '--epfmt', default='%(epname)',
    help=('Specify the format string for sequence diagram endpoints. '
        'May include %%(epname), %%(eplongname) and %%(@foo) for attribute foo.'))
  argp.add_argument(
    '--app_format', '--appfmt', default='%(appname)',
    help=('Specify the format string for sequence diagram participants. '
        'May include %%(appname) and %%(@foo) for attribute foo.'))
  argp.add_argument(
    '--blackbox', '--bb', action='append',
    help='Apps to be treated as black boxes.')

  diagutil.add_common_diag_options(argp)
