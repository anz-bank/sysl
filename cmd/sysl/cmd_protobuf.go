package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/anz-bank/sysl/pkg/cmdutils"

	"github.com/anz-bank/sysl/pkg/pbutil"
	"gopkg.in/alecthomas/kingpin.v2"
)

type protobuf struct {
	output string
	mode   string
}

func (p *protobuf) Name() string       { return "protobuf" }
func (p *protobuf) MaxSyslModule() int { return 1 }

func (p *protobuf) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate textpb/json/binary").Alias("pb")
	cmd.Flag("output", "output file name").Short('o').Default("-").StringVar(&p.output)
	opts := []string{"textpb", "json", "pb"}
	cmd.Flag("mode", fmt.Sprintf("output mode: [%s]", strings.Join(opts, ","))).
		Default(opts[0]).
		EnumVar(&p.mode, opts...)
	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *protobuf) Execute(args cmdutils.ExecuteArgs) error {
	args.Logger.Debugf("Protobuf: %+v", *p)

	p.output = strings.TrimSpace(p.output)
	p.mode = strings.TrimSpace(p.mode)

	toJSON := p.mode == "json" || p.mode == "" && strings.HasSuffix(p.output, ".json")

	if toJSON {
		if p.output == "-" {
			return pbutil.FJSONPB(os.Stdout, args.Modules[0])
		}
		return pbutil.JSONPB(args.Modules[0], p.output, args.Filesystem)
	}

	if p.mode == "" || p.mode == "textpb" {
		if p.output == "-" {
			return pbutil.FTextPB(os.Stdout, args.Modules[0])
		}
		return pbutil.TextPB(args.Modules[0], p.output, args.Filesystem)
	}

	// output format is binary
	if p.output == "-" {
		return pbutil.GeneratePBBinaryMessage(os.Stdout, args.Modules[0])
	}
	return pbutil.GeneratePBBinaryMessageFile(args.Modules[0], p.output, args.Filesystem)
}
