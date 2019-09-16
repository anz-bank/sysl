package commands

import (
	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

type ExecuteArgs struct {
	Module        *sysl.Module
	ModuleAppName string
	Filesystem    afero.Fs
	Logger        *logrus.Logger
}

type Command interface {
	Init(*kingpin.Application) *kingpin.CmdClause
	Execute(ExecuteArgs) error
	Name() string
	RequireSyslModule() bool
}
