package seqs

import (
	"testing"

	"github.com/anz-bank/sysl/src/proto"
	"github.com/stretchr/testify/assert"
)

func TestMakeStrSet(t *testing.T) {
	a := makeStrSet("a", "b", "c", "e")
	assert.Equal(t, 4, a.Len(), "Unexpect result")
}

func TestMakeStrSetWithDuplicateInitialValues(t *testing.T) {
	a := makeStrSet("a", "b", "c", "e", "a", "a", "c")
	assert.Equal(t, 4, a.Len(), "Unexpect result")
}

func TestMakeStrSetWithEmptyStringInitialValues(t *testing.T) {
	a := makeStrSet("a", "b", "c", "e", "a", "a", "c", "", "")
	assert.Equal(t, 4, a.Len(), "Unexpect result")
}

func TestMakeStrSetWithoutInitialValues(t *testing.T) {
	a := makeStrSet()
	assert.Equal(t, 0, a.Len(), "Unexpect result")
}

func TestMakeStrSetFromPatternsAttrWithEmptyAttrs(t *testing.T) {
	attrs := make(map[string]*sysl.Attribute)

	a := makeStrSetFromPatternsAttr(attrs)
	assert.Equal(t, 0, a.Len(), "Unexpect result")
}

func TestMakeStrSetFromPatternsAttrWithoutPatternAttr(t *testing.T) {
	attrs := make(map[string]*sysl.Attribute)
	attrs["test"] = &sysl.Attribute{Attribute: &sysl.Attribute_S{S: "test"}}

	a := makeStrSetFromPatternsAttr(attrs)
	assert.Equal(t, 0, a.Len(), "Unexpect result")
}

func TestMakeStrSetFromPatternsAttr(t *testing.T) {
	attrs := make(map[string]*sysl.Attribute)
	attrs["patterns"] = &sysl.Attribute{
		Attribute: &sysl.Attribute_A{
			A: &sysl.Attribute_Array{
				Elt: []*sysl.Attribute{
					{Attribute: &sysl.Attribute_S{S: "test"}},
				},
			},
		},
	}

	a := makeStrSetFromPatternsAttr(attrs)
	assert.Equal(t, 1, a.Len(), "Unexpect result")
}

func TestContains(t *testing.T) {
	a := makeStrSet("a", "b", "c", "e")
	assert.Equal(t, 4, a.Len(), "Unexpect result")
	assert.True(t, a.Contains("b"), "Unexpect result")
	assert.False(t, a.Contains("d"), "Unexpect result")
}

func TestInsert(t *testing.T) {
	a := makeStrSet("a", "b", "c", "e")
	assert.Equal(t, 4, a.Len(), "Unexpect result")
	assert.True(t, a.Contains("b"), "Unexpect result")
	assert.False(t, a.Contains("d"), "Unexpect result")

	a.Insert("d")
	assert.Equal(t, 5, a.Len(), "Unexpect result")
	assert.True(t, a.Contains("b"), "Unexpect result")
	assert.True(t, a.Contains("d"), "Unexpect result")
}

func TestRemove(t *testing.T) {
	a := makeStrSet("a", "b", "c", "e")
	assert.Equal(t, 4, a.Len(), "Unexpect result")
	assert.True(t, a.Contains("b"), "Unexpect result")
	assert.False(t, a.Contains("d"), "Unexpect result")

	a.Remove("d")
	assert.Equal(t, 4, a.Len(), "Unexpect result")
	assert.True(t, a.Contains("b"), "Unexpect result")
	assert.False(t, a.Contains("d"), "Unexpect result")

	a.Remove("b")
	assert.Equal(t, 3, a.Len(), "Unexpect result")
	assert.False(t, a.Contains("b"), "Unexpect result")
	assert.False(t, a.Contains("d"), "Unexpect result")
}

func TestLen(t *testing.T) {
	a := makeStrSet("a", "b", "c", "e")
	assert.Equal(t, 4, a.Len(), "Unexpect result")
}

func TestToSlice(t *testing.T) {
	a := makeStrSet("c", "b", "a", "e")
	assert.Equal(t, []string{"a", "b", "c", "e"}, a.ToSlice(), "Unexpect result")
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
	assert.Equal(t, 5, c.Len(), "Unexpect result")
	assert.Equal(t, []string{"a", "b", "c", "d", "e"}, c.ToSlice(), "Unexpect result")
}

func TestIntersection(t *testing.T) {
	a := makeStrSet("c", "b", "a", "e")
	b := makeStrSet("d", "b", "a", "e")

	c := a.Intersection(b)
	assert.Equal(t, 3, c.Len(), "Unexpect result")
	assert.Equal(t, []string{"a", "b", "e"}, c.ToSlice(), "Unexpect result")
}

func TestDifference(t *testing.T) {
	a := makeStrSet("c", "b", "a", "e")
	b := makeStrSet("d", "b", "a", "e")

	c := a.Difference(b)
	assert.Equal(t, 1, c.Len(), "Unexpect result")
	assert.Equal(t, []string{"c"}, c.ToSlice(), "Unexpect result")
}
