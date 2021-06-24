package importer

import (
	"sort"
	"strings"

	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/sirupsen/logrus"
)

type Func func(args OutputData, text string, logger *logrus.Logger) (out string, err error)

type Type interface {
	Name() string
}

type StandardType struct {
	name       string
	Properties FieldList
	Attrs      []string
}

func (s *StandardType) Name() string { return s.name }

// Union represents a union type https://sysl.io/docs/lang-spec#union
type Union struct {
	name    string
	Options FieldList
}

func (u *Union) Name() string { return u.name }

type SyslBuiltIn struct {
	name syslutil.BuiltInType
}

func (s *SyslBuiltIn) Name() string { return s.name }

var StringAlias = &SyslBuiltIn{name: syslutil.Type_STRING}

// !alias type without the EXTERNAL_ prefix
type Alias struct {
	name   string
	Target Type
	Attrs  []string
}

func (s *Alias) Name() string { return s.name }

type ExternalAlias struct {
	name   string
	Target Type
	Attrs  []string
}

func NewStringAlias(name string, attrs ...string) Type {
	return &ExternalAlias{
		name:   name,
		Target: StringAlias,
		Attrs:  attrs,
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
	Attrs []string
}

func (s *Array) Name() string { return s.name }

type Enum struct {
	name  string
	Attrs []string
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
	Name     string
	Type     Type
	Optional bool
	Attrs    []string
	SizeSpec *sizeSpec
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

func (props FieldList) Sort() {
	sort.SliceStable(props, func(i, j int) bool {
		return strings.Compare(props[i].Name, props[j].Name) < 0
	})
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
	for _, i := range item {
		if i.Name() != "" {
			t.types = append(t.types, i)
		}
	}
}

func (t *TypeList) AddAndRet(item Type) Type {
	if item.Name() != "" {
		t.types = append(t.types, item)
	}
	return item
}

func checkBuiltInTypes(name string) (Type, bool) {
	if syslutil.IsBuiltIn(name) {
		return &SyslBuiltIn{name: name}, true
	}
	return &StandardType{}, false
}
