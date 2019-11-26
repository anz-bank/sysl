package main

import (
	"fmt"
	"strings"

	"github.com/anz-bank/sysl/src/sysl/pbutil"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

type protobuf struct {
	output string
	mode   string
}

func (p *protobuf) Name() string            { return "protobuf" }
func (p *protobuf) RequireSyslModule() bool { return true }

func (p *protobuf) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate textpb/json").Alias("pb")
	cmd.Flag("output", "output file name").Short('o').Default("-").StringVar(&p.output)
	opts := []string{"textpb", "json"}
	cmd.Flag("mode", fmt.Sprintf("output mode: [%s]", strings.Join(opts, ","))).
		Default(opts[0]).
		EnumVar(&p.mode, opts...)
	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *protobuf) Execute(args ExecuteArgs) error {
	args.Logger.Debugf("Protobuf: %+v", *p)

	p.output = strings.TrimSpace(p.output)
	p.mode = strings.TrimSpace(p.mode)

	toJSON := p.mode == "json" || p.mode == "" && strings.HasSuffix(p.output, ".json")

	if toJSON {
		if p.output == "-" {
			return pbutil.FJSONPB(args.Logger.Out, args.Module)
		}
		return pbutil.JSONPB(args.Module, p.output, args.Filesystem)
	}
	if p.output == "-" {
		return pbutil.FTextPB(logrus.StandardLogger().Out, args.Module)
	}
	return pbutil.TextPB(args.Module, p.output, args.Filesystem)
}
