package syslutil

import (
	"testing"

	sysl "github.com/anz-bank/sysl/src/proto_old"
	"github.com/stretchr/testify/assert"
)

func TestMakeStrSet(t *testing.T) {
	t.Parallel()

	a := MakeStrSet("a", "b", "c", "e")
	assert.Equal(t, 4, len(a))
}

func TestMakeStrSetWithDuplicateInitialValues(t *testing.T) {
	t.Parallel()

	a := MakeStrSet("a", "b", "c", "e", "a", "a", "c")
	assert.Equal(t, 4, len(a))
}

func TestMakeStrSetWithEmptyStringInitialValues(t *testing.T) {
	t.Parallel()

	a := MakeStrSet("a", "b", "c", "e", "a", "a", "c", "", "")
	assert.Equal(t, 5, len(a))
}

func TestMakeStrSetWithoutInitialValues(t *testing.T) {
	t.Parallel()

	a := MakeStrSet()
	assert.Equal(t, 0, len(a))
}

func TestMakeStrSetFromSpecificAttrWithEmptyAttrs(t *testing.T) {
	t.Parallel()

	attrs := map[string]*sysl.Attribute{}

	a := MakeStrSetFromAttr("patterns", attrs)
	assert.Equal(t, 0, len(a))
}

func TestMakeStrSetFromSpecificAttrWithoutPatternAttr(t *testing.T) {
	t.Parallel()

	attrs := map[string]*sysl.Attribute{
		"test": {Attribute: &sysl.Attribute_S{S: "test"}},
	}

	a := MakeStrSetFromAttr("patterns", attrs)
	assert.Equal(t, 0, len(a))
}

func TestMakeStrSetFromPatternsAttr(t *testing.T) {
	t.Parallel()

	attrs := map[string]*sysl.Attribute{
		"patterns": {Attribute: &sysl.Attribute_A{A: &sysl.Attribute_Array{
			Elt: []*sysl.Attribute{{Attribute: &sysl.Attribute_S{S: "test"}}},
		}}},
	}

	a := MakeStrSetFromAttr("patterns", attrs)
	assert.Equal(t, 1, len(a))
}

func TestMakeStrSetFromActionStatement(t *testing.T) {
	t.Parallel()

	stmts := []*sysl.Statement{
		{Stmt: &sysl.Statement_Action{Action: &sysl.Action{Action: "AppA"}}},
	}

	a := MakeStrSetFromActionStatement(stmts)
	assert.Equal(t, 1, len(a))
}

func TestContains(t *testing.T) {
	t.Parallel()

	a := MakeStrSet("a", "b", "c", "e")
	assert.Equal(t, 4, len(a))
	assert.True(t, a.Contains("b"))
	assert.False(t, a.Contains("d"))
}

func TestInsert(t *testing.T) {
	t.Parallel()

	a := MakeStrSet("a", "b", "c", "e")
	assert.Equal(t, 4, len(a))
	assert.True(t, a.Contains("b"))
	assert.False(t, a.Contains("d"))

	a.Insert("d")
	assert.Equal(t, 5, len(a))
	assert.True(t, a.Contains("b"))
	assert.True(t, a.Contains("d"))
}

func TestRemove(t *testing.T) {
	t.Parallel()

	a := MakeStrSet("a", "b", "c", "e")
	assert.Equal(t, 4, len(a))
	assert.True(t, a.Contains("b"))
	assert.False(t, a.Contains("d"))

	a.Remove("d")
	assert.Equal(t, 4, len(a))
	assert.True(t, a.Contains("b"))
	assert.False(t, a.Contains("d"))

	a.Remove("b")
	assert.Equal(t, 3, len(a))
	assert.False(t, a.Contains("b"))
	assert.False(t, a.Contains("d"))
}

func TestToSlice(t *testing.T) {
	t.Parallel()

	// Given
	a := MakeStrSet("c", "b", "a", "e")

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

	assert.True(t, sameValue([]string{"a", "b", "c", "e"}, slice))
	assert.Equal(t, []string{"a", "b", "c", "e"}, sorted)
}

func TestClone(t *testing.T) {
	t.Parallel()

	a := MakeStrSet("c", "b", "a", "e")
	b := a.Clone()
	assert.Equal(t, a, b)

	b.Remove("c")
	assert.NotEqual(t, a, b)
}

func TestStrSetUnion(t *testing.T) {
	t.Parallel()

	a := MakeStrSet("d", "b", "e", "f")
	b := a.Union(MakeStrSet("d", "f", "g", "z"))
	c := MakeStrSet("d", "b", "e", "f", "g", "z")

	assert.Equal(t, b, c)
}

func TestIntersection(t *testing.T) {
	t.Parallel()

	a := MakeStrSet("c", "b", "a", "e")
	b := MakeStrSet("d", "b", "a", "e")

	c := a.Intersection(b)
	assert.Equal(t, 3, len(c))
	assert.Equal(t, []string{"a", "b", "e"}, c.ToSortedSlice())
}

func TestDifference(t *testing.T) {
	t.Parallel()

	a := MakeStrSet("c", "b", "a", "e")
	b := MakeStrSet("d", "b", "a", "e")

	c := a.Difference(b)
	assert.Equal(t, 1, len(c))
	assert.Equal(t, []string{"c"}, c.ToSortedSlice())
}

func TestDifferenceWhenEmptySet(t *testing.T) {
	t.Parallel()

	a := MakeStrSet()
	b := MakeStrSet("a", "z", "y")
	c := a.Difference(b)

	assert.Equal(t, c, a)
}

func TestSubWhenParentAndChildEmpty(t *testing.T) {
	t.Parallel()

	// Given
	c := MakeStrSet()
	p := MakeStrSet()
	expected := true

	// When
	actual := c.IsSubset(p)

	// Then
	assert.Equal(t, expected, actual)
}

func TestSubWhenParentEmpty(t *testing.T) {
	t.Parallel()

	// Given
	c := MakeStrSet("A")
	p := MakeStrSet()
	expected := false

	// When
	actual := c.IsSubset(p)

	// Then
	assert.Equal(t, expected, actual)
}

func TestSubWhenChildEmpty(t *testing.T) {
	t.Parallel()

	// Given
	c := MakeStrSet()
	p := MakeStrSet("A")
	expected := true

	// When
	actual := c.IsSubset(p)

	// Then
	assert.Equal(t, expected, actual)
}

func TestSubSetWhenParentLessThanChild(t *testing.T) {
	t.Parallel()

	// Given
	c := MakeStrSet("a", "z", "y")
	p := MakeStrSet("a", "y")
	expected := false

	// When
	actual := c.IsSubset(p)

	// Then
	assert.Equal(t, expected, actual)
}

func TestSubSetWhenDifferent(t *testing.T) {
	t.Parallel()

	assert.False(t, MakeStrSet("a", "z", "f").IsSubset(MakeStrSet("a", "z", "y", "d")))
}
