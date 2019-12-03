package importer

import (
	"io"
	"strings"
)

type IndentWriter struct {
	current int
	text    string

	io.Writer
}

func (i *IndentWriter) Push() {
	i.current++
}

func (i *IndentWriter) Pop() {
	if i.current > 0 {
		i.current--
	}
}

func (i *IndentWriter) Write() error {
	_, err := i.Writer.Write([]byte(strings.Repeat(i.text, i.current)))
	return err
}

func (i *IndentWriter) CurrentIndentLen() int {
	return len(strings.Repeat(i.text, i.current))
}

func NewIndentWriter(text string, out io.Writer) *IndentWriter {
	return &IndentWriter{
		current: 0,
		text:    text,
		Writer:  out,
	}
}
