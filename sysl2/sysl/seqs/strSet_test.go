package seqs

import (
	"testing"

	"github.com/anz-bank/sysl/src/proto"
	"github.com/stretchr/testify/assert"
)

func TestMakeStrSet(t *testing.T) {
	a := makeStrSet("a", "b", "c", "e")
	assert.Equal(t, 4, len(a), "Unexpect result")
}

func TestMakeStrSetWithDuplicateInitialValues(t *testing.T) {
	a := makeStrSet("a", "b", "c", "e", "a", "a", "c")
	assert.Equal(t, 4, len(a), "Unexpect result")
}

func TestMakeStrSetWithEmptyStringInitialValues(t *testing.T) {
	a := makeStrSet("a", "b", "c", "e", "a", "a", "c", "", "")
	assert.Equal(t, 5, len(a), "Unexpect result")
}

func TestMakeStrSetWithoutInitialValues(t *testing.T) {
	a := makeStrSet()
	assert.Equal(t, 0, len(a), "Unexpect result")
}

func TestMakeStrSetFromPatternsAttrWithEmptyAttrs(t *testing.T) {
	attrs := map[string]*sysl.Attribute{}

	a := makeStrSetFromPatternsAttr(attrs)
	assert.Equal(t, 0, len(a), "Unexpect result")
}

func TestMakeStrSetFromPatternsAttrWithoutPatternAttr(t *testing.T) {
	attrs := map[string]*sysl.Attribute{
		"test": {Attribute: &sysl.Attribute_S{S: "test"}},
	}

	a := makeStrSetFromPatternsAttr(attrs)
	assert.Equal(t, 0, len(a), "Unexpect result")
}

func TestMakeStrSetFromPatternsAttr(t *testing.T) {
	attrs := map[string]*sysl.Attribute{
		"patterns": {
			Attribute: &sysl.Attribute_A{
				A: &sysl.Attribute_Array{
					Elt: []*sysl.Attribute{
						{Attribute: &sysl.Attribute_S{S: "test"}},
					},
				},
			},
		},
	}

	a := makeStrSetFromPatternsAttr(attrs)
	assert.Equal(t, 1, len(a), "Unexpect result")
}

func TestContains(t *testing.T) {
	a := makeStrSet("a", "b", "c", "e")
	assert.Equal(t, 4, len(a), "Unexpect result")
	assert.True(t, a.Contains("b"), "Unexpect result")
	assert.False(t, a.Contains("d"), "Unexpect result")
}

func TestInsert(t *testing.T) {
	a := makeStrSet("a", "b", "c", "e")
	assert.Equal(t, 4, len(a), "Unexpect result")
	assert.True(t, a.Contains("b"), "Unexpect result")
	assert.False(t, a.Contains("d"), "Unexpect result")

	a.Insert("d")
	assert.Equal(t, 5, len(a), "Unexpect result")
	assert.True(t, a.Contains("b"), "Unexpect result")
	assert.True(t, a.Contains("d"), "Unexpect result")
}

func TestRemove(t *testing.T) {
	a := makeStrSet("a", "b", "c", "e")
	assert.Equal(t, 4, len(a), "Unexpect result")
	assert.True(t, a.Contains("b"), "Unexpect result")
	assert.False(t, a.Contains("d"), "Unexpect result")

	a.Remove("d")
	assert.Equal(t, 4, len(a), "Unexpect result")
	assert.True(t, a.Contains("b"), "Unexpect result")
	assert.False(t, a.Contains("d"), "Unexpect result")

	a.Remove("b")
	assert.Equal(t, 3, len(a), "Unexpect result")
	assert.False(t, a.Contains("b"), "Unexpect result")
	assert.False(t, a.Contains("d"), "Unexpect result")
}

func TestToSlice(t *testing.T) {
	// Given
	a := makeStrSet("c", "b", "a", "e")

	// When
	slice := a.ToSlice()
	sorted := a.ToSortedSlice()

	// Then
	sameValue := func(a, b []string) bool {
		if len(a) != len(b) {
			return false
		}
		set := map[string]struct{}{}
		for _, v := range a {
			set[v] = struct{}{}
		}
		for _, v := range b {
			if _, ok := set[v]; !ok {
				return false
			}
		}
		return true
	}

	assert.True(t, sameValue([]string{"a", "b", "c", "e"}, slice), "Unexpect result")
	assert.Equal(t, []string{"a", "b", "c", "e"}, sorted, "Unexpect result")
}

func TestClone(t *testing.T) {
	a := makeStrSet("c", "b", "a", "e")
	b := a.Clone()
	assert.Equal(t, a, b, "Unexpect result")

	b.Remove("c")
	assert.NotEqual(t, a, b, "Unexpect result")
}

func TestUnion(t *testing.T) {
	a := makeStrSet("c", "b", "a", "e")
	b := makeStrSet("d", "b", "a", "e")

	c := a.Union(b)
	assert.Equal(t, 5, len(c), "Unexpect result")
	assert.Equal(t, []string{"a", "b", "c", "d", "e"}, c.ToSortedSlice(), "Unexpect result")
}

func TestIntersection(t *testing.T) {
	a := makeStrSet("c", "b", "a", "e")
	b := makeStrSet("d", "b", "a", "e")

	c := a.Intersection(b)
	assert.Equal(t, 3, len(c), "Unexpect result")
	assert.Equal(t, []string{"a", "b", "e"}, c.ToSortedSlice(), "Unexpect result")
}

func TestDifference(t *testing.T) {
	a := makeStrSet("c", "b", "a", "e")
	b := makeStrSet("d", "b", "a", "e")

	c := a.Difference(b)
	assert.Equal(t, 1, len(c), "Unexpect result")
	assert.Equal(t, []string{"c"}, c.ToSortedSlice(), "Unexpect result")
}
