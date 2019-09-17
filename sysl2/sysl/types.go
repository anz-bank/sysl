package main

import (
	sysl "github.com/anz-bank/sysl/src/proto"
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
	rootModel     *string
	model         *string
	rootTransform *string
	transform     *string
	grammar       *string
	start         *string
	outDir        *string
}

type CmdContextParamSeqgen struct {
	endpointFormat string
	appFormat      string
	title          string
	output         string
	endpointsFlag  []string
	appsFlag       []string
	blackboxesFlag map[string]string
	blackboxes     [][]string
	group          string
}

type CmdContextParamIntgen struct {
	root      *string
	title     *string
	output    *string
	project   *string
	filter    *string
	modules   *string
	exclude   *[]string
	clustered *bool
	epa       *bool
	plantuml  *string
}

type CmdContextParamDatagen struct {
	root        *string
	title       *string
	output      *string
	project     *string
	filter      *string
	modules     *string
	plantuml    *string
	classFormat *string
}

type ExecuteArgs struct {
	Module        *sysl.Module
	ModuleAppName string
	Filesystem    afero.Fs
	Logger        *logrus.Logger
}

type Command interface {
	Configure(*kingpin.Application) *kingpin.CmdClause
	Execute(ExecuteArgs) error
	Name() string
	RequireSyslModule() bool
}
