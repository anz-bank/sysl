package seqs

import (
	"github.com/anz-bank/sysl/src/proto"
)

type Visitor interface {
	Visit(Element)
}

type Element interface {
	Accept(Visitor)
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

type Labeler interface {
	AppLabeler
	EndpointLabeler
}

type VarManager interface {
	UniqueVarForAppName(appName string) string
}
