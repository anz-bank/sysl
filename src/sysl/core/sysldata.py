# -*-: encoding: utf-8 -*-
"""Generate data views from sysl modules."""

import collections
import re
from argparse import RawTextHelpFormatter

from sysl.util import diagutil
from sysl.util import writer
from sysl.util.argparse import add_common_diag_options
from sysl.core import syslloader


def _attr_sort_key(field):
    return field[1].source_context.start.line


def _make_varmgr(module, appname, write):
    """Return a variable manager instance."""
    def new_var(var, name):
        """Outputs a new definition when VarManager makes a new variable."""
        app = module.apps[name]
        write('class "{}" as {} << (D,orchid) >> {{', name, var)
        typespec = module.apps.get(appname).types.get(name)
        assert typespec.WhichOneof('type') == 'tuple'
        fields = sorted(typespec.tuple.attr_defs.iteritems(),
                        key=_attr_sort_key)
        for (fieldname, fieldtype) in fields:
            which = fieldtype.WhichOneof('type')
            suffix = ''
            prefix = ''
            while which == 'list':
                fieldtype = fieldtype.list.type
                which = fieldtype.WhichOneof('type')
                suffix = '[]' + suffix

            if which == 'set':
                fieldtype = fieldtype.set
                which = fieldtype.WhichOneof('type')
                prefix = 'set <'
                suffix = '>' + suffix

            bold = False
            if which == 'primitive':
                typestr = fieldtype.Primitive.Name(fieldtype.primitive).lower()
            elif which == 'type_ref':
                typestr = '.'.join(fieldtype.type_ref.ref.path)
                bold = True
            else:
                typestr = '<color red>**{}**</color>'.format(which)
            typestr = prefix + typestr + suffix
            if bold:
                typestr = '**{}**'.format(typestr)
            write('+ {} : {}', fieldname, typestr)
        write('}}')

    return diagutil.VarManager(new_var)


def _generate_view(module, appname, types):
    """Output integration view"""
    write = writer.Writer('plantuml')

    var_name = _make_varmgr(module, appname, write)

    with write.uml():
        for (appname, name, typespec) in types:
            var_name(name)

            link_sets = collections.defaultdict(
                lambda: collections.defaultdict(list))

            fields = sorted(
                typespec.tuple.attr_defs.iteritems(), key=_attr_sort_key)
            for (fieldname, fieldtype) in fields:
                cardinality = u' '
                while fieldtype.WhichOneof('type') == 'list':
                    fieldtype = fieldtype.list.type
                    cardinality = u'0..*'

                if fieldtype.WhichOneof('type') == 'set':
                    fieldtype = fieldtype.set
                    cardinality = u'0..*'

                if fieldtype.WhichOneof('type') == 'type_ref':
                    ref = u'.'.join(fieldtype.type_ref.ref.path)
                    # Hacky!
                    if ref.startswith(u'Common Data.'):
                        continue

                    refs = [n for (_, n, _) in types if n.endswith(ref)]

                    line_template = u'{} {{}} *-- "{}" {}'.format(
                        var_name(name), cardinality, var_name(refs[0]) if refs else ref)
                    link_sets[ref][line_template].append(fieldname)

            for (_, line_templates) in link_sets.iteritems():
                if len(line_templates) > 1:
                    for (line_template, fieldnames) in line_templates.iteritems():
                        for fieldname in fieldnames:
                            write(line_template, '"' + fieldname + '"')
                else:
                    for (line_template, fieldnames) in line_templates.iteritems():
                        for _ in fieldnames:
                            write(line_template, '')

    return str(write)


def dataviews(module, args):
    """Generate a set of data views."""
    out = []

    parts = re.match(ur'^(.*?)(?:\s*<-\s*(.*?))?$', args.project)
    (appname, epname) = parts.groups()

    app = module.apps.get(appname)
    if not app:
        raise Exception('Invalid project "{}"'.format(args.project))
    if epname is not None:
        endpts = [app.endpoints.get(epname)]
    else:
        endpts = app.endpoints.itervalues()

    out_fmt = diagutil.parse_fmt(args.output)

    for endpt in endpts:
        types = []

        for stmt in endpt.stmt:
            appname = stmt.action.action
            if not module.apps.get(appname):
                continue
            for (name, typespec) in module.apps.get(appname).types.iteritems():
                if typespec.WhichOneof('type') == 'tuple':
                    types.append((appname, name, typespec))

        args.output = out_fmt(
            appname=appname,
            epname=endpt.name,
            eplongname=endpt.long_name,
            **diagutil.attr_fmt_vars(app.attrs, endpt.attrs))

        if args.filter and not re.match(args.filter, args.output):
            continue

        out.append(_generate_view(module, appname, types))

        diagutil.output_plantuml(args, out[-1])

    return out


def add_subparser(subp):
    """Setup data subcommand."""
    argp = subp.add_parser('data', formatter_class=RawTextHelpFormatter)

    def cmd(args):
        """Handle subcommand."""
        (module, _, _) = syslloader.load(
            args.modules, args.validations, args.root)

        out = dataviews(module, args)

    argp.set_defaults(func=cmd)

    argp.add_argument('--project', '-j',
                      help='project pseudo-app to render \n' +
                      'eg: PetShopModel, "ATM <- GetBalance"',
                      required=True)

    add_common_diag_options(argp)
