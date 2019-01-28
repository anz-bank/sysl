package sysl

import "sort"

type typesSourceOrder []*Type

// Len is the number of elements in the collection.
func (t typesSourceOrder) Len() int {
	return len(t)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (t typesSourceOrder) Less(i, j int) bool {
	a := t[i].SourceContext
	b := t[j].SourceContext
	if a == nil {
		return b != nil
	}
	if b == nil {
		return false
	}
	return a.Start.Line < b.Start.Line
}

// Swap swaps the elements with indexes i and j.
func (t typesSourceOrder) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// TypesInSourceOrder a type map's names in the order the types appear in the
// source file.
func TypesInSourceOrder(types []*Type) []*Type {
	tmp := make([]*Type, len(types))
	copy(tmp, types)
	sort.Sort(typesSourceOrder(tmp))
	return tmp
}

type typeNamesSourceOrder struct {
	types map[string]*Type
	names []string
}

// Len is the number of elements in the collection.
func (t *typeNamesSourceOrder) Len() int {
	return len(t.names)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (t *typeNamesSourceOrder) Less(i, j int) bool {
	a := t.types[t.names[i]].SourceContext.Start.Line
	b := t.types[t.names[j]].SourceContext.Start.Line
	return a < b
}

// Swap swaps the elements with indexes i and j.
func (t *typeNamesSourceOrder) Swap(i, j int) {
	t.names[i], t.names[j] = t.names[j], t.names[i]
}

// TypeNamesInSourceOrder a type map's names in the order the types appear in
// the source file.
func TypeNamesInSourceOrder(types map[string]*Type) []string {
	names := make([]string, 0, len(types))
	for name := range types {
		names = append(names, name)
	}
	sort.Sort(&typeNamesSourceOrder{types, names})
	return names
}

// TypeNamesInAlphabeticalOrder returns a type map's names in alphabetical
// order.
func TypeNamesInAlphabeticalOrder(types map[string]*Type) []string {
	names := make([]string, 0, len(types))
	for name := range types {
		names = append(names, name)
	}
	sort.Sort(&typeNamesSourceOrder{types, names})
	return names
}
