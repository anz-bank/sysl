package syslutil

import (
	"sort"

	sysl "github.com/anz-bank/sysl/pkg/sysl"
)

type StrSet map[string]struct{}

func MakeStrSet(initial ...string) StrSet {
	s := StrSet{}

	for _, v := range initial {
		s.Insert(v)
	}

	return s
}

func MakeStrSetFromAttr(attr string, attrs map[string]*sysl.Attribute) StrSet {
	s := StrSet{}

	if patterns, has := attrs[attr]; has {
		if x := patterns.GetA(); x != nil {
			for _, y := range x.Elt {
				if v, ok := y.Attribute.(*sysl.Attribute_S); ok {
					s.Insert(v.S)
				}
			}
		}
	}

	return s
}

func MakeStrSetFromActionStatement(stmts []*sysl.Statement) StrSet {
	s := StrSet{}
	for _, stmt := range stmts {
		if a, ok := stmt.Stmt.(*sysl.Statement_Action); ok {
			s.Insert(a.Action.Action)
		}
	}

	return s
}

func (s StrSet) Contains(elem string) bool {
	_, ok := s[elem]
	return ok
}

func (s StrSet) Insert(elem string) {
	s[elem] = struct{}{}
}

func (s StrSet) Remove(elem string) {
	delete(s, elem)
}

func (s StrSet) ToSlice() []string {
	o := make([]string, 0, len(s))

	for k := range s {
		o = append(o, k)
	}

	return o
}

func (s StrSet) ToSortedSlice() []string {
	slice := s.ToSlice()
	sort.Strings(slice)

	return slice
}

func (s StrSet) Clone() StrSet {
	out := StrSet{}

	for k := range s {
		out.Insert(k)
	}

	return out
}

func (s StrSet) Union(other StrSet) StrSet {
	out := StrSet{}

	for k := range s {
		out.Insert(k)
	}

	for k := range other {
		out.Insert(k)
	}

	return out
}

func (s StrSet) Intersection(other StrSet) StrSet {
	out := StrSet{}

	for k := range s {
		if other.Contains(k) {
			out.Insert(k)
		}
	}

	return out
}

// Returns the elements that only belong to s. If s is subset of other,
// it would return an empty set. Just be aware that output may be one of
// the inputs and when change the output, the input would also be changed.
func (s StrSet) Difference(other StrSet) StrSet {
	if len(s) == 0 || len(other) == 0 {
		return s
	}
	out := StrSet{}

	for k := range s {
		if !other.Contains(k) {
			out.Insert(k)
		}
	}

	return out
}

// Returns true if child set is a subset of parent set
func (s StrSet) IsSubset(parent StrSet) bool {
	if len(parent) == 0 {
		return len(s) == 0
	}
	if len(parent) < len(s) {
		return false
	}
	for k := range s {
		if !parent.Contains(k) {
			return false
		}
	}
	return true
}
