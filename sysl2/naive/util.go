package parser

type intSet map[int]struct{}

func (s *intSet) add(tok int) bool {
	l := len(*s)
	(*s)[tok] = struct{}{}
	return l != len(*s)
}

func (s *intSet) has(tok int) bool {
	_, has := (*s)[tok]
	return has
}

func (s *intSet) union(other *intSet) bool {
	l := len(*s)
	for k := range *other {
		(*s)[k] = struct{}{}
	}
	return l != len(*s)
}

func (s *intSet) clone() *intSet {
	ss := make(intSet)
	for k := range *s {
		ss[k] = struct{}{}
	}
	return &ss
}

func (s *intSet) remove(tok int) {
	delete(*s, tok)
}
