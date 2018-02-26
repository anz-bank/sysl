"""sysl module loader"""

import codecs
import collections
import os
import re
import sys

from sysl.proto import sysl_pb2

from sysl.core import syslalgo
from sysl.core import syslparse
from sysl.core import syslx


def fmt_app_name(appname):
    """Format a sysl_pb2.AppName as a syntactically valid string."""
    return ' :: '.join(appname.part)


def fmt_call(call):
    """Format a sysl_pb2.Call as a syntactically valid string."""
    return fmt_app_name(call.target) + ' <- ' + call.endpoint


def _resolve_mixins(module):
    """Resolve mixin references.

    Copies endpoints from mixed-in apps.
    """
    # Detect cycles
    edges = {
        (appname, syslx.fmt_app_name(mixin.name))
        for (appname, app) in module.apps.iteritems()
        for mixin in app.mixin2}
    while True:
        more_edges = {
            (a, c)
            for (a, b1) in edges
            for (b2, c) in edges
            if b1 == b2
        } - edges
        if not more_edges:
            break
        edges |= more_edges
    self_edges = {(a, b) for (a, b) in edges if a == b}
    if self_edges:
        raise RuntimeError(
            "mixin cycle(s) detected involving: {}".format(
                ', '.join(a for (a, _) in self_edges)))

    # recursively inject mixins, avoiding double-injection
    injected = set()

    def inject(appname):
        app = module.apps[appname]
        if appname not in injected:
            for mixin in app.mixin2:
                mixin_app = inject(syslx.fmt_app_name(mixin.name))

                # Check for ~abstract
                assert 'abstract' in syslx.patterns(mixin_app.attrs), (
                    "mixin {} must be ~abstract".format(syslx.fmt_app_name(mixin.name)))

                for (epname, endpt) in mixin_app.endpoints.iteritems():
                    app.endpoints[epname].CopyFrom(endpt)

                for (tname, t) in mixin_app.types.iteritems():
                    app.types[tname].CopyFrom(t)

                for (vname, v) in mixin_app.views.iteritems():
                    app.views[vname].CopyFrom(v)

            injected.add(appname)

        return app

    for appname in module.apps:
        inject(appname)


def _check_deps(module, validate):
    """Check app:endpoint dependencies."""
    deps = set()
    errors = []

    for (appname, app) in module.apps.iteritems():
        for epname in app.endpoints:
            endpt = app.endpoints[epname]

            for (_, call) in syslalgo.enumerate_calls(endpt.stmt):
                targetname = syslx.fmt_app_name(call.target)
                if targetname not in module.apps:
                    errors.append('{} <- {}: calls non-existent app {}'.format(
                        appname, epname, targetname))
                else:
                    target = module.apps[targetname]
                    assert 'abstract' not in syslx.patterns(target.attrs), (
                        "call target '{}' must not be ~abstract".format(targetname))
                    if call.endpoint not in target.endpoints:
                        errors.append(
                            '{} <- {}: calls non-existent endpoint {} -> {}'.format(
                                appname, epname, targetname, call.endpoint))
                    else:
                        deps.add(
                            ((appname, epname), (targetname, call.endpoint)))

    if errors and validate:
        raise Exception('broken deps:\n  ' + '\n  '.join(errors))

    return deps


def _map_subscriptions(module):
    """Map pubsub subscriptions into direct calls."""
    for appname in module.apps:
        app = module.apps[appname]

        if 'abstract' in syslx.patterns(app.attrs):
            continue

        for epname in app.endpoints:
            endpt = app.endpoints[epname]

            if endpt.HasField('source'):
                src_app = module.apps[syslx.fmt_app_name(endpt.source)]
                src_ep_name = endpt.name.split(' -> ')[1]
                assert src_ep_name in src_app.endpoints, (
                    appname, epname, src_ep_name, str(src_app))
                src_ep = src_app.endpoints[src_ep_name]

                # Add call to pubsub endpoint.
                stmt = src_ep.stmt.add()
                call = stmt.call
                call.target.CopyFrom(app.name)
                call.endpoint = endpt.name

                # Maybe add ret.
                ret_payload = syslalgo.return_payload(endpt.stmt)
                if ret_payload:
                    stmt = src_ep.stmt.add()
                    stmt.ret.payload = ret_payload


def _apply_call_templates(app):
    """Apply call templates found in '.. * <- *' | '*' pseudo-endpoints.

    Project-specific metadata may be applied as follows:

      MyApp:
        .. * <- *:
          Foo <- bar [myproj='XYZ-007']

    In the above example, whenever MyApp calls Foo <- bar,
    _apply_call_templates will attach the attribute myproj='XYZ-007' to the
    call.

    It will also validate that all templates are applied at least once.
    """

    # Look for the pseudo endpoint.
    pseudos = {name for name in app.endpoints
               if re.match(r'(\.\.\s*\*\s*<-\s*\*|\*)', name)}
    if not pseudos:
        return
    if len(pseudos) > 1:
        raise Exception('Too many call templates: {}'.format(
            ', '.join(repr(p) for p in pseudos)))

    pseudo = app.endpoints[pseudos.pop()]

    templates = {}

    def call_templates():
        # Discover templates.
        for stmt in pseudo.stmt:
            if stmt.HasField('call'):
                templates[fmt_call(stmt.call)] = [stmt.attrs, 0]

        # Apply templates.
        endpoints = [endpt for (_, endpt) in app.endpoints.iteritems(
        ) if not re.match(r'(\.\.\s*\*\s*<-\s*\*|\*)', endpt.name)]

        for endpt in endpoints:
            for (stmt, call) in syslalgo.enumerate_calls(endpt.stmt):
                fmtcall = fmt_call(call)
                template = templates.get(fmtcall)
                if template is not None:
                    for (name, attr) in template[0].iteritems():
                        stmt.attrs[name].CopyFrom(attr)
                    template[1] += 1

    def ep_templates():
        # Discover templates.
        for stmt in pseudo.stmt:
            if stmt.HasField('action'):
                templates[stmt.action.action] = [stmt.attrs, 0]

        # Apply templates.
        endpoints = [endpt for (_, endpt) in app.endpoints.iteritems(
        ) if not re.match(r'(\.\.\s*\*\s*<-\s*\*|\*)', endpt.name)]

        for endpt in endpoints:
            template = templates.get(endpt.name)
            if template is not None:
                for (name, attr) in template[0].iteritems():
                    endpt.attrs[name].CopyFrom(attr)
                template[1] += 1

    call_templates()
    ep_templates()

    # Error on unused templates, in case of typos.
    call = None  # In case of empty loop
    unused = {
        call
        for (call, n) in templates.iteritems()
        if n[1] == 0}

    # TODO: add better message
    # App is referring to an unused app-endpoint
    if unused:
        raise RuntimeError('Unused templates in {}: {}', fmt_app_name(
            app.name), ', '.join(repr(c) for c in unused))


def _infer_types(app):
    """Infer types of views and expressions from their bodies.

    Synthesize types for anonymous transforms.
    """

    for (vname, v) in app.views.iteritems():
        assert (
            (v.expr.WhichOneof('expr') == 'transform') ^
            ('abstract' in syslx.patterns(v.attrs))
        ), '{}: {}'.format(vname, v.expr)

        if v.ret_type.WhichOneof('type') is None:
            assert v.expr.type.WhichOneof('type')
            v.ret_type.CopyFrom(v.expr.type)

        nAnons = [0]

        def infer_expr_type(expr, top=True):
            which_expr = expr.WhichOneof('expr')

            if which_expr == 'transform':
                transform = expr.transform

                # Must recurse first
                for stmt in transform.stmt:
                    which_stmt = stmt.WhichOneof('stmt')
                    if which_stmt in ['assign', 'let']:
                        infer_expr_type(getattr(stmt, which_stmt).expr, False)

                if not top and not expr.type.WhichOneof('type'):
                    tname = 'AnonType_{}__'.format(nAnons[0])
                    nAnons[0] += 1

                    newt = app.types[tname].tuple

                    for stmt in transform.stmt:
                        which_stmt = stmt.WhichOneof('stmt')
                        if which_stmt == 'assign':
                            assign = stmt.assign
                            aexpr = assign.expr
                            assert aexpr.WhichOneof(
                                'expr') == 'transform', aexpr
                            ftype = aexpr.type
                            setof = ftype.WhichOneof('type') == 'set'
                            if setof:
                                ftype = ftype.set
                            assert ftype.WhichOneof('type') == 'type_ref'

                            t = sysl_pb2.Type()
                            tr = t.type_ref
                            tr.context.appname.CopyFrom(app.name)
                            tr.context.path.append(tname)
                            tr.ref.CopyFrom(ftype.type_ref.ref)

                            if setof:
                                t = sysl_pb2.Type(set=t)

                            newt.attr_defs[assign.name].CopyFrom(t)

                    tr = expr.type.set.type_ref
                    tr.context.appname.CopyFrom(app.name)
                    tr.ref.appname.CopyFrom(app.name)
                    tr.ref.path.append(tname)

            elif which_expr == 'relexpr':
                relexpr = expr.relexpr

                if relexpr.op == relexpr.RANK:
                    if not top and not expr.type.WhichOneof('type'):
                        raise RuntimeError(
                            "rank() type inference not implemented")
                        expr.type.CopyFrom(infer_expr_type(relexpr.target))
                        rank = expr.type.add()
                        rank.primitive = rank.INT

            return expr.type

        infer_expr_type(v.expr)


def postprocess(module):
    _resolve_mixins(module)
    _map_subscriptions(module)
    for (appname, app) in module.apps.iteritems():
        _apply_call_templates(app)
        _infer_types(app)


def load(names, validate, root):
    """Load a sysl module."""
    if isinstance(names, basestring):
        names = [names]

    module = sysl_pb2.Module()
    imports = set()

    def do_import(name, indent="-"):
        """Import a sysl module and its dependencies."""
        imports.add(name)
        # print indent, name
        (basedir, _) = os.path.split(name)
        new_imports = {
            root + i if i[:1] == '/' else os.path.join(basedir, i)
            for i in syslparse.Parser().parse(
                codecs.open(name + '.sysl', 'r', 'utf8'), name + '.sysl', module)
        } - imports
        #print >>sys.stderr, '+++++++++++', new_imports
        while new_imports:
            do_import(new_imports.pop(), indent + "-")
            new_imports -= imports

    for name in names:
        if name not in imports:
            if name[:1] != '/':
                raise RuntimeError('module ref must start with "/"')
            do_import(root + name)

    try:
        postprocess(module)
        deps = _check_deps(module, validate)
    except RuntimeError as ex:
        raise Exception('load({!r})'.format(names), ex, sys.exc_info()[2])

    return (module, deps, imports)
