#!/usr/bin/env python
# -*-: encoding: utf-8 -*-
"""Generate integration views from sysl modules."""

import itertools
import re

from sysl.core import syslloader
from sysl.core import syslx

from sysl.util import diagutil
from sysl.util import writer
from sysl.util.argparse import add_common_diag_options


def yield_call_statements(statements):
    for (_, stmt) in enumerate(statements):
        field = stmt.WhichOneof('stmt')
        if field == 'call':
            yield stmt
        elif field == 'alt':
            for (_, choice) in enumerate(stmt.alt.choice):
                for next_stmt in yield_call_statements(choice.stmt):
                    yield next_stmt
        elif field in {'cond', 'loop', 'loop_n', 'foreach', 'group'}:
            for next_stmt in yield_call_statements(getattr(stmt, field).stmt):
                yield next_stmt


def _generate_view(module, args, integrations, highlights, app, apps, endpt):
    write = writer.Writer('plantuml')

    """Output integration view"""

    restrict_by = endpt.attrs.get(
        'restrict_by').s if endpt.attrs.get('restrict_by') else None

    app_attrs = syslx.View(app.attrs)
    endpt_attrs = syslx.View(endpt.attrs)
    highlight_color = app_attrs['highlight_color'].s
    arrow_color = app_attrs['arrow_color'].s
    indirect_arrow_color = app_attrs['indirect_arrow_color'].s

    diagram_title = ''
    if (app_attrs['title'].s or args.title):
        fmtfn = diagutil.parse_fmt(app_attrs['title'].s or args.title)
        diagram_title = fmtfn(epname=endpt.name, eplongname=endpt.long_name)

    def generate_component_view():

        name_map = {}

        def make_varmgr():
            """Return a variable manager instance."""
            appfmt = syslx.View(module.apps.get(
                args.project).attrs)['appfmt'].s
            appfmt = diagutil.parse_fmt(appfmt)

            def new_var(var, name):
                """Outputs a new definition when VarManager makes a new variable."""

                app = module.apps[name]
                write('[{}] as {}{}',
                      appfmt(appname=name_map.get(name, name),
                             **diagutil.attr_fmt_vars(app.attrs)).replace('\n', r'\n'),
                      var,
                      ' <<highlight>>' if name in highlights else '')

            return diagutil.VarManager(new_var)

        with write.uml():

            if diagram_title:
                write('title ' + diagram_title)

            write('hide stereotype')
            write('scale max 16384 height')

            #write('skinparam nodesep 30')
            #write('skinparam ranksep 120')
            write('skinparam component {{')
            write('  BackgroundColor FloralWhite')
            write('  BorderColor Black')
            write('  ArrowColor Crimson')
            if highlight_color:
                write('  BackgroundColor<<highlight>> ' + highlight_color)
            if arrow_color:
                write('  ArrowColor ' + arrow_color)
            if indirect_arrow_color and indirect_arrow_color != 'none':
                write('  ArrowColor<<indirect>> ' + indirect_arrow_color)
            write('}}')

            var_name = make_varmgr()

            if args.clustered or endpt_attrs['view'].s == 'clustered':
                clusters = diagutil.group_by(
                    apps, key=lambda app: app.partition(' :: ')[0])
                clusters = [
                    (cluster, members)
                    for (cluster, g) in clusters
                    for members in [list(g)]
                    if len(members) > 1]
                name_map = {app: app.partition(' :: ')[-1] or app
                            for (cluster, members) in clusters
                            for app in members}
                for (cluster, cluster_apps) in clusters:
                    write('package "{}" {{', cluster)
                    for app in cluster_apps:
                        var_name(app)
                        # write('  {}', var_name(app))
                    write('}}')

            calls_drawn = set()

            if endpt_attrs['view'].s == 'system':
                for ((app_a, _), (app_b, _), (_, _)) in integrations:
                    direct = {app_a, app_b} & highlights
                    app_a = app_a.partition(' :: ')[0]
                    app_b = app_b.partition(' :: ')[0]
                    if app_a != app_b and (app_a, app_b) not in calls_drawn:
                        if direct or indirect_arrow_color != 'none':
                            write('{} --> {}{}',
                                  var_name(app_a),
                                  var_name(app_b),
                                  '' if direct else ' <<indirect>>')
                            calls_drawn.add((app_a, app_b))
            else:
                for ((app_a, _), (app_b, _)) in integrations:
                    if app_a != app_b and (app_a, app_b) not in calls_drawn:
                        direct = {app_a, app_b} & highlights
                        if direct or indirect_arrow_color != 'none':
                            write('{} --> {}{}',
                                  var_name(app_a),
                                  var_name(app_b),
                                  '' if direct else ' <<indirect>>')
                            calls_drawn.add((app_a, app_b))

                for appname in apps:
                    for mixin in module.apps[appname].mixin2:
                        mixin_name = syslx.fmt_app_name(mixin.name)
                        mixin_app = module.apps[mixin_name]
                        write('{} <|.. {}',
                              var_name(mixin_name),
                              var_name(appname))

    # TODO Some serious refactoring
    def generate_state_view():

        def make_varmgr(istoplevel=False):
            """Return a variable manager instance."""
            appfmt = syslx.View(module.apps.get(
                args.project).attrs)['appfmt'].s
            appfmt = diagutil.parse_fmt(appfmt)

            def new_var(var, name):
                """Outputs a new definition when VarManager makes a new variable."""

                # TODO dodgy, should be using context (look at syslseqs)
                attrs = {}
                state_name = ''

                if istoplevel:
                    template = 'state "{}" as X{}{} {{'
                    attrs = module.apps[name].attrs
                    state_name = name
                else:
                    template = '  state "{}" as {}{}'
                    (app_name, _, ep_name) = name.partition(' : ')
                    attrs = module.apps[app_name].endpoints[ep_name].attrs
                    state_name = ep_name

                write(template,
                      appfmt(appname=state_name,
                             **diagutil.attr_fmt_vars(attrs)).replace('\n', r'\n'),
                      var,
                      ' <<highlight>>' if name in highlights else '')

            return diagutil.VarManager(new_var)

        with write.uml():

            if diagram_title:
                write('title ' + diagram_title)

            write('left to right direction')
            write('scale max 16384 height')
            write('hide empty description')
            write('skinparam state {{')
            write('  BackgroundColor FloralWhite')
            write('  BorderColor Black')
            write('  ArrowColor Crimson')
            if highlight_color:
                write('  BackgroundColor<<highlight>> ' + highlight_color)
            if arrow_color:
                write('  ArrowColor ' + arrow_color)
            if indirect_arrow_color and indirect_arrow_color != 'none':
                write('  ArrowColor<<indirect>> ' + indirect_arrow_color)
                write('  ArrowColor<<internal>> ' + indirect_arrow_color)
            write('}}')

            var_name = make_varmgr()
            tl_var_name = make_varmgr(True)

            clusters = {}

            # group end points and build the declarations
            for (app_a, ep_a), (app_b, ep_b) in integrations:

                if restrict_by and restrict_by not in module.apps[app_a].attrs.keys(
                ) + module.apps[app_b].attrs.keys():
                    continue

                if restrict_by and restrict_by not in module.apps[app_a].endpoints[ep_a].attrs.keys(
                ) + module.apps[app_b].endpoints[ep_b].attrs.keys():
                    continue

                if app_a not in clusters:
                    clusters[app_a] = set()
                if app_b not in clusters:
                    clusters[app_b] = set()

                client_ep = ep_b
                # create clients in the calling app
                clusters[app_a].add(ep_a)
                if app_a != app_b and not module.apps[app_a].endpoints[ep_a].is_pubsub:
                    clusters[app_a].add(client_ep + " client")

                clusters[app_b].add(ep_b)

            for cluster in clusters:
                tl_var_name(cluster)
                for member in clusters[cluster]:
                    var_name(cluster + ' : ' + member)
                write('}}')

            processed = []
            for ((app_a, ep_a), (app_b, ep_b)) in integrations:

                if restrict_by and restrict_by not in module.apps[app_a].attrs.keys(
                ) + module.apps[app_b].attrs.keys():
                    continue

                if restrict_by and restrict_by not in module.apps[app_a].endpoints[ep_a].attrs.keys(
                ) + module.apps[app_b].endpoints[ep_b].attrs.keys():
                    continue

                direct = {app_a, app_b} & highlights

                (match_app, match_ep) = (app_b, ep_b)

                # build the label
                label = ''
                needs_int = app_a != match_app

                pub_sub_src_ptrns = syslx.patterns(
                    module.apps[app_a].endpoints[ep_a].attrs)

                tgt_ptrns = syslx.patterns(module.apps[match_app].endpoints[match_ep].attrs) or \
                    syslx.patterns(module.apps[match_app].attrs)

                for stmt in yield_call_statements(
                        module.apps[app_a].endpoints[ep_a].stmt):

                    # if restrict_by not in stmt.attrs.keys():
                    #  continue

                    app_b_name = ' :: '.join(
                        part for part in stmt.call.target.part)

                    if match_app == app_b_name and match_ep == stmt.call.endpoint:
                        fmt_vars = diagutil.attr_fmt_vars(stmt.attrs)

                        src_ptrns = syslx.patterns(
                            stmt.attrs) or pub_sub_src_ptrns

                        ptrns = u', '.join(sorted(src_ptrns)) + u' → ' + u', '.join(sorted(tgt_ptrns)) \
                            if (bool(src_ptrns) or bool(tgt_ptrns)) else ''

                        label = diagutil.parse_fmt(app.attrs["epfmt"].s)(
                            needs_int=needs_int, patterns=ptrns, **fmt_vars)

                # if not label and middle:
                #  fmt_vars = diagutil.attr_fmt_vars(module.apps[app_b].endpoints[ep_b].attrs)
                #  label = diagutil.parse_fmt(app.attrs["epfmt"].s)(needs_int=needs_int, patterns=u'→' + middle + u'→', **fmt_vars)

                flow = ".".join([app_a, ep_b, app_b, ep_b])
                is_pubsub = module.apps[app_a].endpoints[ep_a].is_pubsub
                ep_b_client = ep_b + " client"
                if app_a != app_b:
                    if is_pubsub:
                        write('{} -{}> {}{}',
                              var_name(app_a + ' : ' + ep_a),
                              '[#blue]',
                              var_name(app_b + ' : ' + ep_b),
                              ' : ' + label if label else '')
                    else:
                        write('{} -{}> {}',
                              var_name(app_a + ' : ' + ep_a),
                              '[#' + indirect_arrow_color +
                              ']-' if indirect_arrow_color else '[#silver]-',
                              var_name(app_a + ' : ' + ep_b_client))
                        if flow not in processed:
                            write('{} -{}> {}{}',
                                  var_name(app_a + ' : ' + ep_b_client),
                                  '[#black]',
                                  var_name(app_b + ' : ' + ep_b),
                                  ' : ' + label if label else '')
                            processed.append(flow)

                else:
                    write('{} -{}> {}{}',
                          var_name(app_a + ' : ' + ep_a),
                          '[#' + indirect_arrow_color +
                          ']-' if indirect_arrow_color else '[#silver]-',
                          var_name(app_b + ' : ' + ep_b),
                          ' : ' + label if label else '')

    if (args.epa) or endpt_attrs['view'].s == 'epa':
        generate_state_view()
    else:
        generate_component_view()

    return str(write)


def integration_views(module, deps, args):
    """Generate an integration view."""

    def find_matching_apps(integrations):
        """Yield all apps that match the integrations."""
        exclude = set(args.exclude)
        app_re = re.compile(
            r'^(?:{})(?: *::|$)'.format('|'.join(integrations)))
        for ((app1, _), (app2, _)) in deps:
            if not ({app1, app2} & exclude):
                for app in [app1, app2]:
                    if app_re.match(app) and 'human' not in syslx.patterns(
                            module.apps[app].attrs):
                        yield app

    def find_apps(matching_apps):
        """Yield all apps that are relevant to this view."""
        exclude = set(args.exclude)
        for ((app1, _), (app2, _)) in deps:
            if not ({app1, app2} & exclude) and ({app1, app2} & matching_apps):
                for app in [app1, app2]:
                    if ('human' not in syslx.patterns(module.apps[app].attrs)):
                        yield app

    out_fmt = diagutil.parse_fmt(args.output)

    # The "project" app that specifies the required view of the integrations
    app = module.apps[args.project]

    out = []

    # Interate over each endpoint within the selected project
    for endpt in app.endpoints.itervalues():

        # build the set of excluded items
        exclude_attr = endpt.attrs.get('exclude')
        exclude = set(
            e.s for e in exclude_attr.a.elt) if exclude_attr else set()

        passthrough_apps_attr = endpt.attrs.get('passthrough')
        passthrough_apps = set(
            e.s for e in passthrough_apps_attr.a.elt) if passthrough_apps_attr else set()

        integrations = []

        # endpt.stmt's "action" will conatain the "apps" whose integration is to be drawn
        # each one of these will be placed into the "integrations" list
        for s in endpt.stmt:
            assert s.WhichOneof('stmt') == 'action', str(s)
            integrations.append(s.action.action)

        # include the requested "app" and all the apps upon which the requested "app"
        # depends in the app set
        matching_apps = set(find_matching_apps(integrations))
        apps = set(find_apps(matching_apps)) - exclude - passthrough_apps

        def find_integrations():
            """Return all integrations between relevant apps."""

            integration_set = {((app1, ep1), (app2, ep2))
                               for ((app1, ep1), (app2, ep2)) in deps
                               if ({app1, app2} <= apps and not ({app1, app2} & exclude)) and not {ep1, ep2} & {'.. * <- *', '*'}}

            if passthrough_apps:

                def passthrough_walk(dep, integration_set):
                    """ visit passthrough + 1: dep tuple, integration_set tuples """

                    # add to integration set
                    if (not ({dep[1][0]} & exclude) and not (
                            {dep[1][1]} & {'.. * <- *', '*'})):
                        integration_set.add(dep)

                    # find the next outbond dep
                    if dep[1][0] in passthrough_apps:
                        # call visit_endpoints for all calls
                        for stmt in yield_call_statements(
                                module.apps[dep[1][0]].endpoints[dep[1][1]].stmt):
                            app_2_name = ' :: '.join(
                                part for part in stmt.call.target.part)
                            ep_2_name = stmt.call.endpoint
                            passthrough_walk(
                                ((dep[1]), (app_2_name, ep_2_name)), integration_set)

                    return

                # collect outbound dependencies
                outbound_deps = {((app1, ep1), (app2, ep2))
                                 for ((app1, ep1), (app2, ep2)) in deps
                                 if (({app1, app2} <= apps) or ({app1} <= apps and {app2} <= passthrough_apps)) and not
                                    ({app1, app2} & exclude) and not ({ep1, ep2} & {'.. * <- *', '*'})}

                # collect outbound pubsub dependencies
                # inbound_deps

                # add passthroughs
                for dep in outbound_deps:
                    passthrough_walk(dep, integration_set)

            return integration_set

        args.output = out_fmt(
            appname=args.project,
            epname=endpt.name,
            eplongname=endpt.long_name,
            **diagutil.attr_fmt_vars(app.attrs, endpt.attrs))

        if args.filter and not re.match(args.filter, args.output):
            continue

        # print args.project, matching_apps
        out.append(_generate_view(module, args, find_integrations(),
                                  matching_apps, app, apps, endpt))

        diagutil.output_plantuml(args, out[-1])

    return out


def add_subparser(subp):
    """Setup ints subcommand."""
    argp = subp.add_parser('ints')

    def cmd(args):
        """Handle subcommand."""
        (module, deps, _) = syslloader.load(
            args.modules, args.validations, args.root)

        if not args.exclude and args.project:
            args.exclude = {args.project}

        integration_views(module, deps, args)

    argp.set_defaults(func=cmd)

    argp.add_argument('--clustered', '-c', action='store_true',
                      help='group integration components into clusters')
    argp.add_argument('--exclude', '-e', action='append',
                      help='apps to exclude')
    argp.add_argument('--project', '-j',
                      help='project pseudo-app to render')

    argp.add_argument('--epa', '-epa', action='store_true',
                      help='produce and EPA integration view')

    add_common_diag_options(argp)
