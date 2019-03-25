package seqs

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	// given
	w := MakeSequenceDiagramWriter(true)

	// when
	n, err := w.Write([]byte("test\n"))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, 0, w.ind)
	assert.True(t, w.atBeginOfLine)
}

func TestWriteWithoutln(t *testing.T) {
	// given
	w := MakeSequenceDiagramWriter(true)

	// when
	n, err := w.Write([]byte("test"))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 4, n)
	assert.Equal(t, 0, w.ind)
	assert.False(t, w.atBeginOfLine)
}

func TestWriteWithoutlnInNewLine(t *testing.T) {
	// given
	w := MakeSequenceDiagramWriter(true)
	w.atBeginOfLine = true
	w.ind = 1

	// when
	n, err := w.Write([]byte("test"))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 4, n)
	assert.Equal(t, 5, w.body.Len())
	assert.Equal(t, 1, w.ind)
	assert.False(t, w.atBeginOfLine)
}

func TestWriteMultiLines(t *testing.T) {
	// given
	w := MakeSequenceDiagramWriter(true)
	w.atBeginOfLine = true
	w.ind = 1

	// when
	n, err := w.Write([]byte("line1\nline2"))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 11, n)
	assert.Equal(t, 13, w.body.Len())
	assert.Equal(t, 1, w.ind)
	assert.False(t, w.atBeginOfLine)
}

func TestWriteString(t *testing.T) {
	// given
	w := MakeSequenceDiagramWriter(true)

	// when
	n, err := w.WriteString("test\n")

	// then
	assert.Nil(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, 0, w.ind)
	assert.True(t, w.atBeginOfLine)
}

func TestWriteByte(t *testing.T) {
	// given
	w := MakeSequenceDiagramWriter(true)

	// when
	err := w.WriteByte('a')

	// then
	assert.Nil(t, err)
	assert.Equal(t, 1, w.body.Len())
	assert.False(t, w.atBeginOfLine)
}

func TestWriteByteln(t *testing.T) {
	// given
	w := MakeSequenceDiagramWriter(true)

	// when
	err := w.WriteByte('\n')

	// then
	assert.Nil(t, err)
	assert.Equal(t, 1, w.body.Len())
	assert.True(t, w.atBeginOfLine)
}

func TestWriteHead(t *testing.T) {
	// given
	w := MakeSequenceDiagramWriter(true)

	// when
	w.WriteHead("head")

	// then
	assert.Equal(t, 5, w.head.Len())
}

func TestIndent(t *testing.T) {
	// given
	w := MakeSequenceDiagramWriter(true)

	// when
	w.Indent()

	// then
	assert.Equal(t, 1, w.ind)
}

func TestUnindent(t *testing.T) {
	w := MakeSequenceDiagramWriter(true)

	assert.Panics(t, func() {
		w.Unindent()
	})

	w.Indent()
	assert.Equal(t, 1, w.ind)

	w.Unindent()
	assert.Equal(t, 0, w.ind)

	assert.Panics(t, func() {
		w.Unindent()
	})
}

func TestActivate(t *testing.T) {
	w := MakeSequenceDiagramWriter(true)

	w.Activate("a")
	assert.Equal(t, 11, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{"a": 1}, w.active)

	w.Activate("a")
	assert.Equal(t, 22, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{"a": 2}, w.active)
}

func TestActivated(t *testing.T) {
	w := MakeSequenceDiagramWriter(true)

	d := w.Activated("a", false)
	assert.Equal(t, 11, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{"a": 1}, w.active)

	d()
	assert.Equal(t, 24, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{}, w.active)

	d()
	assert.Equal(t, 24, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{}, w.active)
}

func TestActivatedWithSuppressed(t *testing.T) {
	w := MakeSequenceDiagramWriter(true)

	d := w.Activated("a", true)
	assert.Equal(t, 0, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{}, w.active)

	d()
	assert.Equal(t, 0, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{}, w.active)

	d()
	assert.Equal(t, 0, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{}, w.active)
}

func TestDeactivate(t *testing.T) {
	w := MakeSequenceDiagramWriter(true)

	w.Activate("a")
	assert.Equal(t, 11, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{"a": 1}, w.active)

	w.Deactivate("a")
	assert.Equal(t, 24, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{}, w.active)

	w.Deactivate("a")
	assert.Equal(t, 24, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{}, w.active)

	w.Deactivate("b")
	assert.Equal(t, 24, w.body.Len())
	assert.True(t, w.atBeginOfLine)
	assert.Equal(t, map[string]int{}, w.active)
}

func TestWriteIndent(t *testing.T) {
	w := MakeSequenceDiagramWriter(true)
	w.atBeginOfLine = true
	w.ind = 1

	w.writeIndent()
	assert.Equal(t, 1, w.body.Len())
	assert.False(t, w.atBeginOfLine)

	w.writeIndent()
	assert.Equal(t, 1, w.body.Len())
	assert.False(t, w.atBeginOfLine)
}

func TestStringer(t *testing.T) {
	// Given
	w := MakeSequenceDiagramWriter(false)
	w.WriteHead("head")
	_, err := w.WriteString("body\n")

	// When
	s := w.String()

	// Then
	assert.Nil(t, err)
	assert.Equal(t, "@startuml\nhead\nbody\n@enduml\n", s)
}

func TestStringerEmpty(t *testing.T) {
	w := MakeSequenceDiagramWriter(false)
	_, err := w.WriteString("body\n")

	s := w.String()
	expected := ""

	assert.Nil(t, err)
	assert.Equal(t, expected, s)
}

func TestStringerWithAutogen(t *testing.T) {
	// Given
	w := MakeSequenceDiagramWriter(true)
	w.WriteHead("head")
	_, err := w.WriteString("body\n")

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
	assert.Nil(t, err)
	assert.Equal(t, expected, s)
}

func TestWriteTo(t *testing.T) {
	// Given
	w := MakeSequenceDiagramWriter(false, "properties 1")
	w.WriteHead("head")
	w.WriteString("body\n")

	// When
	var b bytes.Buffer
	n, err := w.WriteTo(&b)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, int64(41), n)
	assert.Equal(t, "@startuml\nhead\nproperties 1\nbody\n@enduml\n", b.String())
}
