package seqs

import (
	"sort"

	"github.com/anz-bank/sysl/src/proto"
)

type strSet map[string]struct{}

func makeStrSet(initial ...string) strSet {
	s := strSet{}

	for _, v := range initial {
		s.Insert(v)
	}

	return s
}

func makeStrSetFromPatternsAttr(attrs map[string]*sysl.Attribute) strSet {
	s := strSet{}

	if patterns, has := attrs["patterns"]; has {
		if x := patterns.GetA(); x != nil {
			for _, y := range x.Elt {
				if v := y.GetS(); len(v) > 0 {
					s.Insert(y.GetS())
				}
			}
		}
	}

	return s
}

func (s strSet) Contains(elem string) bool {
	_, ok := s[elem]
	return ok
}

func (s strSet) Insert(elem string) {
	s[elem] = struct{}{}
}

func (s strSet) Remove(elem string) {
	delete(s, elem)
}

func (s strSet) ToSlice() []string {
	o := make([]string, 0, len(s))

	for k := range s {
		o = append(o, k)
	}

	return o
}

func (s strSet) ToSortedSlice() []string {
	slice := s.ToSlice()
	sorted := make([]string, len(slice))
	copy(sorted, slice)

	sort.Strings(sorted)

	return sorted
}

func (s strSet) Clone() strSet {
	out := strSet{}

	for k := range s {
		out.Insert(k)
	}

	return out
}

func (s strSet) Union(other strSet) strSet {
	out := strSet{}

	for k := range s {
		out.Insert(k)
	}

	for k := range other {
		out.Insert(k)
	}

	return out
}

func (s strSet) Intersection(other strSet) strSet {
	out := strSet{}

	for k := range s {
		if other.Contains(k) {
			out.Insert(k)
		}
	}

	return out
}

func (s strSet) Difference(other strSet) strSet {
	out := strSet{}

	for k := range s {
		if !other.Contains(k) {
			out.Insert(k)
		}
	}

	return out
}
