package integrationdiagram

import (
	"fmt"

	"github.com/anz-bank/sysl/pkg/sysl"
)

type AppElement struct {
	Name     string
	Endpoint string
}

type AppDependency struct {
	Self, Target AppElement
	Statement    *sysl.Statement
}

func (dep *AppDependency) String() string {
	return fmt.Sprintf("%s:%s:%s:%s", dep.Self.Name, dep.Self.Endpoint, dep.Target.Name, dep.Target.Endpoint)
}

type IntsArg struct {
	RootModel string
	Title     string
	Output    string
	Project   string
	Filter    string
	Modules   string
	Exclude   []string
	Clustered bool
	Epa       bool
}
