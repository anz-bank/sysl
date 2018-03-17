package main

// Set is Not exported
type Set map[int]struct{}

var empty struct{}

func (s *Set) add(tok int) bool {
    l := len(*s)
    (*s)[tok] = empty
    return l != len(*s)
}

func (s *Set) has(tok int) bool {
    _, has := (*s)[tok]
    return has
}

func (s *Set) union(other *Set) bool {
    l := len(*s)
    for k := range *other {
        (*s)[k] = empty
    }
    return l != len(*s)
}

func (s *Set) clone() *Set {
    ss := make(Set)
    for k := range *s {
        ss[k] = empty
    }
    return &ss
}

func (s *Set) remove(tok int) {
    delete(*s, tok)
}
