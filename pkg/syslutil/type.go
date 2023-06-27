package syslutil

import (
	"sort"

	sysl "github.com/anz-bank/sysl/pkg/sysl"
)

type NamedType struct {
	Name string
	Type *sysl.Type
}

type NamedTypes []NamedType

// Len is the number of elements in the collection.
func (t NamedTypes) Len() int {
	return len(t)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (t NamedTypes) Less(i, j int) bool {
	a := t[i].Type.SourceContexts[0].Start.Line
	b := t[j].Type.SourceContexts[0].Start.Line
	return a < b
}

// Swap swaps the elements with indexes i and j.
func (t NamedTypes) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t NamedTypes) Where(p NamedTypePredicate) NamedTypes {
	result := NamedTypes{}
	for _, namedType := range t {
		if p(namedType) {
			result = append(result, namedType)
		}
	}
	return result
}

// NamedTypesInSourceOrder a type map's names in the order the types appear in
// the source file.
func NamedTypesInSourceOrder(types map[string]*sysl.Type) NamedTypes {
	namedTypes := make(NamedTypes, 0, len(types))
	for name, t := range types {
		namedTypes = append(namedTypes, NamedType{Name: name, Type: t})
	}
	sort.Sort(namedTypes)
	return namedTypes
}

type NamedTypePredicate func(nt NamedType) bool

func NamedTypeAll(NamedType) bool {
	return true
}

func NamedTypeNot(p NamedTypePredicate) NamedTypePredicate {
	return func(nt NamedType) bool {
		return !p(nt)
	}
}
