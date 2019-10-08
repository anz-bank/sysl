package importer

import (
	"sort"
	"strings"
)

type Type interface {
	Name() string
}

type StandardType struct {
	name       string
	Properties []Field
	Attributes []string
}

func (s *StandardType) Name() string { return s.name }

type SyslBuiltIn struct {
	name string
}

func (s *SyslBuiltIn) Name() string { return s.name }

type Alias struct {
	name   string
	Target Type
}

func NewStringAlias(name string) Type {
	return &Alias{
		name:   name,
		Target: &SyslBuiltIn{name: "string"},
	}
}

func (s *Alias) Name() string { return s.name }

type ImportedBuiltInAlias struct {
	name   string // input language type name
	Target Type
}

func (s *ImportedBuiltInAlias) Name() string { return s.name }

type Array struct {
	name  string
	Items Type
}

func (s *Array) Name() string { return s.name }

type Enum struct {
	name string
}

func (s *Enum) Name() string { return s.name }

type Field struct {
	Name       string
	Type       Type
	Optional   bool
	Attributes []string
}

type TypeList struct {
	types []Type
}

func (t TypeList) Items() []Type {
	return t.types
}

func (t TypeList) Sort() {
	sort.SliceStable(t.types, func(i, j int) bool {
		return strings.Compare(t.types[i].Name(), t.types[j].Name()) < 0
	})
}

type FieldList []Field

// nolint:gochecknoglobals
var builtIns = []string{"int32", "int64", "int", "float", "string", "date", "bool", "decimal", "datetime", "xml"}

func IsKeyword(name string) bool {
	for _, kw := range builtIns {
		if name == kw {
			return true
		}
	}
	return false
}

func (t TypeList) Find(name string) (Type, bool) {
	if builtin, ok := checkBuiltInTypes(name); ok {
		return builtin, ok
	}

	for _, n := range t.Items() {
		if n.Name() == name {
			if importAlias, ok := n.(*ImportedBuiltInAlias); ok {
				return importAlias.Target, true
			}
			return n, true
		}
	}
	return &StandardType{}, false
}

func (t *TypeList) Add(item ...Type) {
	t.types = append(t.types, item...)
}

func checkBuiltInTypes(name string) (Type, bool) {
	if contains(name, builtIns) {
		return &SyslBuiltIn{name: name}, true
	}
	return &StandardType{}, false
}

func contains(needle string, haystack []string) bool {
	for _, x := range haystack {
		if x == needle {
			return true
		}
	}
	return false
}
