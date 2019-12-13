package importer

import (
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

type Func func(args OutputData, text string, logger *logrus.Logger) (out string, err error)

type Type interface {
	Name() string
}

type StandardType struct {
	name       string
	Properties FieldList
	Attributes []string
}

func (s *StandardType) Name() string { return s.name }

type Union struct {
	name    string
	Options FieldList
}

func (u *Union) Name() string { return u.name }

type SyslBuiltIn struct {
	name string
}

func (s *SyslBuiltIn) Name() string { return s.name }

var StringAlias = &SyslBuiltIn{name: StringTypeName}

// !alias type without the EXTERNAL_ prefix
type Alias struct {
	name   string
	Target Type
}

func (s *Alias) Name() string { return s.name }

type ExternalAlias struct {
	name   string
	Target Type
}

const (
	StringTypeName = "string"
	ObjectTypeName = "object"
	ArrayTypeName  = "array"
)

func NewStringAlias(name string) Type {
	return &ExternalAlias{
		name:   name,
		Target: StringAlias,
	}
}

func (s *ExternalAlias) Name() string { return s.name }

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

type maxType int

const (
	MinOnly maxType = iota
	MaxSpecified
	OpenEnded
)

type sizeSpec struct {
	Min     int
	Max     int
	MaxType maxType
}
type Field struct {
	Name       string
	Type       Type
	Optional   bool
	Attributes []string
	SizeSpec   *sizeSpec
}

type TypeList struct {
	types []Type
}

func (t TypeList) Items() []Type {
	return t.types
}

func (t TypeList) Sort() {
	sort.SliceStable(t.types, func(i, j int) bool {
		a := t.types[i].Name()
		b := t.types[j].Name()
		return strings.Compare(a, b) < 0
	})
}

type FieldList []Field

// nolint:gochecknoglobals
var builtIns = []string{"int32", "int64", "int", "float", "string", "date", "bool", "decimal", "datetime", "xml"}

func IsBuiltIn(name string) bool {
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

func (t *TypeList) AddAndRet(item Type) Type {
	t.types = append(t.types, item)
	return item
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
