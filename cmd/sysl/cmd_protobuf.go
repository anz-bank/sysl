package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/pbutil"
	"github.com/anz-bank/sysl/pkg/sysl"

	"gopkg.in/alecthomas/kingpin.v2"
)

type protobufCmd struct {
	output        string
	mode          string
	compact       bool
	filter        string
	splitappspath string
}

func (p *protobufCmd) Name() string       { return "protobuf" }
func (p *protobufCmd) MaxSyslModule() int { return 1 }

func (p *protobufCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate textpb/json/binary").Alias("pb")
	cmd.Flag("output", "output file name").Short('o').Default("-").StringVar(&p.output)
	cmd.Flag(
		"split-apps",
		`Splits Applications into their own directory structure under the given base path.
		Ignored if output is set. This BREAKS referencing.`).
		Short('s').
		Default("-").
		StringVar(&p.splitappspath)
	opts := []string{"textpb", "json", "pb"}
	cmd.Flag("mode", fmt.Sprintf("output mode: [%s]", strings.Join(opts, ","))).
		Default(opts[0]).
		EnumVar(&p.mode, opts...)
	cmd.Flag("compact", "Output without newlines and indentations").
		Default("false").
		BoolVar(&p.compact)
	cmd.Flag("filter", "filter applications").StringVar(&p.filter)
	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *protobufCmd) Execute(args cmdutils.ExecuteArgs) error {
	args.Logger.Debugf("Protobuf: %+v", *p)

	if p.splitappspath != "-" {
		args.Logger.Warn("Using --split-apps to split input into multiple files will BREAK references")
	}

	p.output = strings.TrimSpace(p.output)
	p.mode = strings.TrimSpace(p.mode)

	toJSON := p.mode == "json" || p.mode == "" && strings.HasSuffix(p.output, ".json")

	opt := pbutil.OutputOptions{Compact: p.compact}

	m := args.Modules[0]

	if p.filter != "" {
		apps := make(map[string]*sysl.Application)

		for app, v := range m.Apps {
			if strings.HasPrefix(app, p.filter) {
				apps[app] = v
			}
		}

		m.Apps = apps
		m.Imports = nil
	}

	if toJSON {
		if p.splitappspath != "-" {
			return pbutil.OutputSplitApplications(m, "json", opt, p.splitappspath, "data.json", args.Filesystem)
		} else if p.output == "-" {
			return pbutil.FJSONPBWithOpt(os.Stdout, m, opt)
		}
		return pbutil.JSONPBWithOpt(m, p.output, args.Filesystem, opt)
	}

	if p.mode == "" || p.mode == "textpb" {
		if p.output == "-" && p.splitappspath == "-" {
			return pbutil.FTextPBWithOpt(os.Stdout, m, opt)
		}
		if p.splitappspath != "-" {
			return pbutil.OutputSplitApplications(m, p.mode, opt, p.splitappspath, "data.textpb", args.Filesystem)
		} else {
			return pbutil.TextPBWithOpt(m, p.output, args.Filesystem, opt)
		}
	}

	// output format is binary
	if p.output == "-" && p.splitappspath == "-" {
		return pbutil.GeneratePBBinaryMessage(os.Stdout, m)
	}

	if p.splitappspath != "-" {
		return pbutil.OutputSplitApplications(m, p.mode, opt, p.splitappspath, "data.pb", args.Filesystem)
	}
	return pbutil.GeneratePBBinaryMessageFile(m, p.output, args.Filesystem)
}
