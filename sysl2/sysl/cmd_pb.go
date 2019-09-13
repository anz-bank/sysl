package main

import (
	"fmt"
	"github.com/anz-bank/sysl/sysl2/sysl/pbutil"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
)

type pbCommand struct {
	output string
	mode string
}

func (p* pbCommand) Name() string { return "pb"}
func (p *pbCommand) RequireSyslModule() bool { return true }

func (p* pbCommand) Init(app *kingpin.Application)  *kingpin.CmdClause{

	cmd := app.Command(p.Name(), "Generate textpb/json")
	cmd.Flag("output", "output file name").Short('o').Default("-").StringVar(&p.output)
	opts := []string{"textpb", "json"}
	cmd.Flag("mode", fmt.Sprintf("output mode: [%s]", strings.Join(opts, ","))).
		Default(opts[0]).
		EnumVar(&p.mode, opts...)
	return cmd
}

func (p* pbCommand) Execute(args ExecuteArgs) error {

	args.logger.Infof("Protobuf: %+v", *p)

	toJSON := p.mode == "json" || p.mode == "" && strings.HasSuffix(p.output, ".json")

	if toJSON {
		if p.output == "-" {
			return pbutil.FJSONPB(args.logger.Out, args.module)
		}
		return pbutil.JSONPB(args.module, p.output, args.fs)
	}
	if p.output == "-" {
		return pbutil.FTextPB(logrus.StandardLogger().Out, args.module)
	}
	return pbutil.TextPB(args.module, p.output, args.fs)
}
