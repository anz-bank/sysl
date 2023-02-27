package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/anz-bank/sysl/pkg/cmdutils"

	"github.com/anz-bank/sysl/pkg/pbutil"
	"gopkg.in/alecthomas/kingpin.v2"
)

type protobufCmd struct {
	output  string
	mode    string
	compact bool
}

func (p *protobufCmd) Name() string       { return "protobuf" }
func (p *protobufCmd) MaxSyslModule() int { return 1 }

func (p *protobufCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate textpb/json/binary").Alias("pb")
	cmd.Flag("output", "output file name").Short('o').Default("-").StringVar(&p.output)
	opts := []string{"textpb", "json", "pb"}
	cmd.Flag("mode", fmt.Sprintf("output mode: [%s]", strings.Join(opts, ","))).
		Default(opts[0]).
		EnumVar(&p.mode, opts...)
	cmd.Flag("compact", "Output without newlines and indentations").
		Default("false").
		BoolVar(&p.compact)
	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *protobufCmd) Execute(args cmdutils.ExecuteArgs) error {
	args.Logger.Debugf("Protobuf: %+v", *p)

	p.output = strings.TrimSpace(p.output)
	p.mode = strings.TrimSpace(p.mode)

	toJSON := p.mode == "json" || p.mode == "" && strings.HasSuffix(p.output, ".json")

	opt := pbutil.OutputOptions{Compact: p.compact}

	if toJSON {
		if p.output == "-" {
			return pbutil.FJSONPBWithOpt(os.Stdout, args.Modules[0], opt)
		}
		return pbutil.JSONPBWithOpt(args.Modules[0], p.output, args.Filesystem, opt)
	}

	if p.mode == "" || p.mode == "textpb" {
		if p.output == "-" {
			return pbutil.FTextPBWithOpt(os.Stdout, args.Modules[0], opt)
		}
		return pbutil.TextPBWithOpt(args.Modules[0], p.output, args.Filesystem, opt)
	}

	// output format is binary
	if p.output == "-" {
		return pbutil.GeneratePBBinaryMessage(os.Stdout, args.Modules[0])
	}
	return pbutil.GeneratePBBinaryMessageFile(args.Modules[0], p.output, args.Filesystem)
}
