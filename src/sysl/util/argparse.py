def add_modules_option(argp):
    argp.add_argument('modules', nargs='+', help='\n'.join([
        'input files without .sysl extension and with leading /',
        'eg: /project_dir/my_models',
        'combine with --root if needed']))


def add_output_option(argp):
    argp.add_argument('--output', '-o', help='output file', required=True)


def add_common_diag_options(argp):
    """Add common diagramming options to a subcommand parser."""
    argp.add_argument(
        '--title', '-t', type=lambda s: unicode(s, 'utf8'),
        help='diagram title')
    argp.add_argument(
        '--plantuml', '-p',
        help=('\n'.join(['base url of plantuml server ',
                         '(default: $SYSL_PLANTUML or http://localhost:8080/plantuml ',
                         'see http://plantuml.com/server.html#install for more info)'])))
    argp.add_argument(
        '--verbose', '-v', action='store_true',
        help='Report each output.')
    argp.add_argument(
        '--expire-cache', action='store_true',
        help='Expire cache entries to force checking against real destination')
    argp.add_argument(
        '--dry-run', action='store_true',
        help="Don't perform confluence uploads, but show what would have happened")
    argp.add_argument(
        '--filter',
        help="Only generate diagrams whose output paths match a pattern")

    add_modules_option(argp)
    add_output_option(argp)
