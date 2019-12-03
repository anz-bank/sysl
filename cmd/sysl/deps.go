package main

import (
	"fmt"

	sysl "github.com/anz-bank/sysl/pkg/proto_old"
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
