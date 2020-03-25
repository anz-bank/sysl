package cmdutils

import (
	"github.com/anz-bank/sysl/pkg/parse"
	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Visitor interface {
	Visit(Element) error
}

type Element interface {
	Accept(Visitor) error
}

type EndpointLabelerParam struct {
	EndpointName string
	Human        string
	HumanSender  string
	NeedsInt     string
	Args         string
	Patterns     string
	Controls     string
	Attrs        map[string]*sysl.Attribute
}

type EndpointLabeler interface {
	LabelEndpoint(*EndpointLabelerParam) string
}

type AppLabeler interface {
	LabelApp(appName, controls string, attrs map[string]*sysl.Attribute) string
}

type ClassLabeler interface {
	LabelClass(className string) string
}

type VarManager interface {
	UniqueVarForAppName(appName string) string
}

type CmdContextParamCodegen struct {
	RootTransform    string
	Transform        string
	Grammar          string
	Start            string
	DepPath          string
	BasePath         string
	DisableValidator bool
}

type CmdContextParamSeqgen struct {
	EndpointFormat string
	AppFormat      string
	Title          string
	Output         string
	EndpointsFlag  []string
	AppsFlag       []string
	BlackboxesFlag map[string]string
	Blackboxes     [][]string
	Group          string
}

type CmdContextParamIntgen struct {
	Title     string
	Output    string
	Project   string
	Filter    string
	Exclude   []string
	Clustered bool
	EPA       bool
}

type CmdContextParamDatagen struct {
	Title       string
	Output      string
	Project     string
	Direct      bool
	Filter      string
	ClassFormat string
}

type CmdDatabaseScriptParams struct {
	Title     string
	OutputDir string
	AppNames  string
	DbType    string
}

type ExecuteArgs struct {
	Command        string
	Modules        []*sysl.Module
	Filesystem     afero.Fs
	Logger         *logrus.Logger
	DefaultAppName string
	ParserType     parse.ParserType
}

type Command interface {
	Configure(*kingpin.Application) *kingpin.CmdClause
	Execute(ExecuteArgs) error
	Name() string
	MaxSyslModule() int
}
