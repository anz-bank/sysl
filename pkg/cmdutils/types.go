package cmdutils

import (
	"io"

	"github.com/alecthomas/kingpin/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/sysl"
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
	DBType    string
}

type ExecuteArgs struct {
	Command        string
	Modules        []*sysl.Module
	Filesystem     afero.Fs
	Logger         *logrus.Logger
	DefaultAppName string
	ModulePaths    []string
	Root           string
	Stdin          io.Reader
	Stdout         io.Writer
}

type DiagramCmd struct {
	AppName            string
	Endpoint           string
	Output             string
	IntegrationDiagram bool
	SequenceDiagram    bool
	EndpointAnalysis   bool
	DataDiagram        bool
}

type Command interface {
	Configure(*kingpin.Application) *kingpin.CmdClause
	Execute(ExecuteArgs) error
	Name() string
	MaxSyslModule() int
}

type PreExecuteCommand interface {
	PreExecute(*parse.Settings) error
}
