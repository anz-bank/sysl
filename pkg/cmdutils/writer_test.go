package cmdutils

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrite(t *testing.T) {
	t.Parallel()

	// given
	w := MakeSequenceDiagramWriter(true)

	// when
	n, err := w.Write([]byte("test\n"))

	// then
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Zero(t, w.Ind)
	assert.True(t, w.AtBeginOfLine)
}

func TestWriteWithoutln(t *testing.T) {
	t.Parallel()

	// given
	w := MakeSequenceDiagramWriter(true)

	// when
	n, err := w.Write([]byte("test"))

	// then
	assert.NoError(t, err)
	assert.Equal(t, 4, n)
	assert.Zero(t, w.Ind)
	assert.False(t, w.AtBeginOfLine)
}

func TestWriteWithoutlnInNewLine(t *testing.T) {
	t.Parallel()

	// given
	w := MakeSequenceDiagramWriter(true)
	w.AtBeginOfLine = true
	w.Ind = 1

	// when
	n, err := w.Write([]byte("test"))

	// then
	assert.NoError(t, err)
	assert.Equal(t, 4, n)
	assert.Equal(t, 5, w.Body.Len())
	assert.Equal(t, 1, w.Ind)
	assert.False(t, w.AtBeginOfLine)
}

func TestWriteMultiLines(t *testing.T) {
	t.Parallel()

	// given
	w := MakeSequenceDiagramWriter(true)
	w.AtBeginOfLine = true
	w.Ind = 1

	// when
	n, err := w.Write([]byte("line1\nline2"))

	// then
	assert.NoError(t, err)
	assert.Equal(t, 11, n)
	assert.Equal(t, 13, w.Body.Len())
	assert.Equal(t, 1, w.Ind)
	assert.False(t, w.AtBeginOfLine)
}

func TestWriteString(t *testing.T) {
	t.Parallel()

	// given
	w := MakeSequenceDiagramWriter(true)

	// when
	n, err := w.WriteString("test\n")

	// then
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Zero(t, w.Ind)
	assert.True(t, w.AtBeginOfLine)
}

func TestWriteByte(t *testing.T) {
	t.Parallel()

	// given
	w := MakeSequenceDiagramWriter(true)

	// when
	err := w.WriteByte('a')

	// then
	assert.NoError(t, err)
	assert.Equal(t, 1, w.Body.Len())
	assert.False(t, w.AtBeginOfLine)
}

func TestWriteByteln(t *testing.T) {
	t.Parallel()

	// given
	w := MakeSequenceDiagramWriter(true)

	// when
	err := w.WriteByte('\n')

	// then
	assert.NoError(t, err)
	assert.Equal(t, 1, w.Body.Len())
	assert.True(t, w.AtBeginOfLine)
}

func TestWriteHead(t *testing.T) {
	t.Parallel()

	// given
	w := MakeSequenceDiagramWriter(true)

	// when
	_, err := w.WriteHead("head")
	require.NoError(t, err)

	// then
	assert.Equal(t, 5, w.Head.Len())
}

func TestIndent(t *testing.T) {
	t.Parallel()

	// given
	w := MakeSequenceDiagramWriter(true)

	// when
	w.Indent()

	// then
	assert.Equal(t, 1, w.Ind)
}

func TestUnindent(t *testing.T) {
	t.Parallel()

	w := MakeSequenceDiagramWriter(true)

	assert.Panics(t, func() {
		w.Unindent()
	})

	w.Indent()
	assert.Equal(t, 1, w.Ind)

	w.Unindent()
	assert.Zero(t, w.Ind)

	assert.Panics(t, func() {
		w.Unindent()
	})
}

func TestActivate(t *testing.T) {
	t.Parallel()

	w := MakeSequenceDiagramWriter(true)

	w.Activate("a")
	assert.Equal(t, 11, w.Body.Len())
	assert.True(t, w.AtBeginOfLine)
	assert.Equal(t, map[string]int{"a": 1}, w.Active)

	w.Activate("a")
	assert.Equal(t, 22, w.Body.Len())
	assert.True(t, w.AtBeginOfLine)
	assert.Equal(t, map[string]int{"a": 2}, w.Active)
}

func TestActivated(t *testing.T) {
	t.Parallel()

	w := MakeSequenceDiagramWriter(true)

	d := w.Activated("a", false)
	assert.Equal(t, 11, w.Body.Len())
	assert.True(t, w.AtBeginOfLine)
	assert.Equal(t, map[string]int{"a": 1}, w.Active)

	d()
	assert.Equal(t, 24, w.Body.Len())
	assert.True(t, w.AtBeginOfLine)
	assert.Equal(t, map[string]int{}, w.Active)

	d()
	assert.Equal(t, 24, w.Body.Len())
	assert.True(t, w.AtBeginOfLine)
	assert.Equal(t, map[string]int{}, w.Active)
}

func TestActivatedWithSuppressed(t *testing.T) {
	t.Parallel()

	w := MakeSequenceDiagramWriter(true)

	d := w.Activated("a", true)
	assert.Zero(t, w.Body.Len())
	assert.True(t, w.AtBeginOfLine)
	assert.Equal(t, map[string]int{}, w.Active)

	d()
	assert.Zero(t, w.Body.Len())
	assert.True(t, w.AtBeginOfLine)
	assert.Equal(t, map[string]int{}, w.Active)

	d()
	assert.Zero(t, w.Body.Len())
	assert.True(t, w.AtBeginOfLine)
	assert.Equal(t, map[string]int{}, w.Active)
}

func TestDeactivate(t *testing.T) {
	t.Parallel()

	w := MakeSequenceDiagramWriter(true)

	w.Activate("a")
	assert.Equal(t, 11, w.Body.Len())
	assert.True(t, w.AtBeginOfLine)
	assert.Equal(t, map[string]int{"a": 1}, w.Active)

	w.Deactivate("a")
	assert.Equal(t, 24, w.Body.Len())
	assert.True(t, w.AtBeginOfLine)
	assert.Equal(t, map[string]int{}, w.Active)

	w.Deactivate("a")
	assert.Equal(t, 24, w.Body.Len())
	assert.True(t, w.AtBeginOfLine)
	assert.Equal(t, map[string]int{}, w.Active)

	w.Deactivate("b")
	assert.Equal(t, 24, w.Body.Len())
	assert.True(t, w.AtBeginOfLine)
	assert.Equal(t, map[string]int{}, w.Active)
}

func TestWriteIndent(t *testing.T) {
	t.Parallel()

	w := MakeSequenceDiagramWriter(true)
	w.AtBeginOfLine = true
	w.Ind = 1

	w.WriteIndent()
	assert.Equal(t, 1, w.Body.Len())
	assert.False(t, w.AtBeginOfLine)

	w.WriteIndent()
	assert.Equal(t, 1, w.Body.Len())
	assert.False(t, w.AtBeginOfLine)
}

func TestStringer(t *testing.T) {
	t.Parallel()

	// Given
	w := MakeSequenceDiagramWriter(false)
	_, err := w.WriteHead("head")
	require.NoError(t, err)
	_, err = w.WriteString("body\n")
	require.NoError(t, err)

	// When
	s := w.String()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "@startuml\nhead\nbody\n@enduml\n", s)
}

func TestStringerEmpty(t *testing.T) {
	t.Parallel()

	w := MakeSequenceDiagramWriter(false)
	_, err := w.WriteString("body\n")

	s := w.String()
	expected := ""

	assert.NoError(t, err)
	assert.Equal(t, expected, s)
}

func TestStringerWithAutogen(t *testing.T) {
	t.Parallel()

	// Given
	w := MakeSequenceDiagramWriter(true)
	_, err := w.WriteHead("head")
	require.NoError(t, err)
	_, err = w.WriteString("body\n")
	require.NoError(t, err)

	// When
	s := w.String()

	// Then
	expected := `''''''''''''''''''''''''''''''''''''''''''
''                                      ''
''  AUTOGENERATED CODE -- DO NOT EDIT!  ''
''                                      ''
''''''''''''''''''''''''''''''''''''''''''

@startuml
head
body
@enduml
`
	assert.NoError(t, err)
	assert.Equal(t, expected, s)
}

func TestWriteTo(t *testing.T) {
	t.Parallel()

	// Given
	w := MakeSequenceDiagramWriter(false, "properties 1")
	_, err := w.WriteHead("head")
	require.NoError(t, err)
	_, err = w.WriteString("body\n")
	require.NoError(t, err)

	// When
	var b bytes.Buffer
	n, err := w.WriteTo(&b)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, int64(41), n)
	assert.Equal(t, "@startuml\nhead\nproperties 1\nbody\n@enduml\n", b.String())
}
