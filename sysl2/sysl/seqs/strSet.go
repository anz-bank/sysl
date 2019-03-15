package seqs

import (
	"sort"

	"github.com/anz-bank/sysl/src/proto"
)

type (
	nothing struct{}

	strSet struct {
		m map[string]nothing
	}
)

func makeStrSet(initial ...string) *strSet {
	s := &strSet{make(map[string]nothing)}

	for _, v := range initial {
		s.Insert(v)
	}

	return s
}

func makeStrSetFromPatternsAttr(attrs map[string]*sysl.Attribute) *strSet {
	s := &strSet{make(map[string]nothing)}

	if patterns, has := attrs["patterns"]; has {
		if x := patterns.GetA(); x != nil {
			for _, y := range x.Elt {
				s.Insert(y.GetS())
			}
		}
	}

	return s
}

func (s *strSet) Contains(elem string) bool {
	_, ok := s.m[elem]
	return ok
}

func (s *strSet) Insert(elem string) {
	if len(elem) > 0 {
		s.m[elem] = nothing{}
	}
}

func (s *strSet) Remove(elem string) {
	delete(s.m, elem)
}

func (s *strSet) Len() int {
	return len(s.m)
}

func (s *strSet) ToSlice() []string {
	o := make([]string, 0, len(s.m))

	for k := range s.m {
		o = append(o, k)
	}

	sort.Strings(o)

	return o
}

func (s *strSet) Clone() *strSet {
	m := make(map[string]nothing)

	for k := range s.m {
		m[k] = nothing{}
	}

	return &strSet{m}
}

func (s *strSet) Union(other *strSet) *strSet {
	m := make(map[string]nothing)

	for k := range s.m {
		m[k] = nothing{}
	}

	for k := range other.m {
		m[k] = nothing{}
	}

	return &strSet{m}
}

func (s *strSet) Intersection(other *strSet) *strSet {
	m := make(map[string]nothing)

	for k := range s.m {
		if _, ok := other.m[k]; ok {
			m[k] = nothing{}
		}
	}

	return &strSet{m}
}

func (s *strSet) Difference(other *strSet) *strSet {
	m := make(map[string]nothing)

	for k := range s.m {
		if _, ok := other.m[k]; !ok {
			m[k] = nothing{}
		}
	}

	return &strSet{m}
}
