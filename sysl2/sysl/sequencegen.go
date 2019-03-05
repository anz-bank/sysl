package main

import (
	"flag"
	"io"
	"strings"

	"github.com/sirupsen/logrus"
)

// DoGenerateSequenceDiagrams generate sequence diagrams for the given model
func DoGenerateSequenceDiagrams(stdout, stderr io.Writer, flags *flag.FlagSet, args []string) error {
	endpoint := flags.String("endpoint", "", "Include endpoint in sequence diagram")
	app := flags.String("app", "", "Include all endpoints for app in sequence diagram (currently " +
		"only works with templated --output). Use SYSL_SD_FILTERS env (a " +
		"comma-list of shell globs) to limit the diagrams generated")
	no_activations := flags.Bool("no-activations", true, "Suppress sequence diagram activation bars(default: true)")
	endpoint_format := flags.String("endpoint_format", "%(epname)", "Specify the format string for sequence diagram endpoints. " +
		"May include %%(epname), %%(eplongname) and %%(@foo) for attribute foo(default: %(epname))")
	app_format := flags.String("app_format", "%(appname)", "Specify the format string for sequence diagram participants. " +
		"May include %%(appname) and %%(@foo) for attribute foo(default: %(appname))")
	blackbox := flags.String("blackbox", "", "Apps to be treated as black boxes")
	title := flags.String("title", "", "diagram title")
	plantuml := flags.String("plantuml", "", strings.Join([]string{"base url of plantuml server",
		"(default: $SYSL_PLANTUML or http://localhost:8080/plantuml",
		"see http://plantuml.com/server.html#install for more info)"}, "\n"))
	verbose := flags.Bool("verbose", false, "Report each output(default: false)")
	expire_cache := flags.Bool("expire-cache", false, "Expire cache entries to force checking against real destination(default: false)")
	dry_run := flags.Bool("dry-run", false, "Don't perform confluence uploads, but show what would have happened(default: false)")
	filter := flags.String("filter", "", "Only generate diagrams whose output paths match a pattern")
	modules := flags.String("modules", "", strings.Join([]string{"input files without .sysl extension and with leading /",
		"eg: /project_dir/my_models",
		"combine with --root if needed"}, "\n"))
	output := flags.String("output", "%(epname).png", "output file(default: %(epname).png)")

	err := flags.Parse(args[1:])
	if err != nil {
		return err
	}
	logrus.Warnf("endpoint: %s\n", *endpoint)
	logrus.Warnf("app: %s\n", *app)
	logrus.Warnf("no_activations: %t\n", *no_activations)
	logrus.Warnf("endpoint_format: %s\n", *endpoint_format)
	logrus.Warnf("app_format: %s\n", *app_format)
	logrus.Warnf("blackbox: %s\n", *blackbox)
	logrus.Warnf("title: %s\n", *title)
	logrus.Warnf("app: %s\n", *app)
	logrus.Warnf("plantuml: %s\n", *plantuml)
	logrus.Warnf("verbose: %t\n", *verbose)
	logrus.Warnf("expire_cache: %t\n", *expire_cache)
	logrus.Warnf("dry_run: %t\n", *dry_run)
	logrus.Warnf("filter: %s\n", *filter)
	logrus.Warnf("modules: %s\n", *modules)
	logrus.Warnf("output: %s\n", *output)

	return nil
}
