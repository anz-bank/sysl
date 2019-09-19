package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gopkg.in/alecthomas/kingpin.v2"
)

func TestEnsureFlagsNonEmpty_AllowsExcludes(t *testing.T) {

	sysl := kingpin.New("sysl", "System Modelling Language Toolkit")
	cmd := sysl.Command("foo", "")
	_ = cmd.Flag("bar", "").Default("foo").String()
	_ = cmd.Flag("other", "").Default("foo").String()

	EnsureFlagsNonEmpty(cmd, "bar")

	args := []string{"foo", "--bar", ""}
	selected, err := sysl.Parse(args)
	assert.Equal(t, "foo", selected)
	assert.NoError(t, err)
}

func TestEnsureFlagsNonEmpty(t *testing.T) {

	sysl := kingpin.New("sysl", "System Modelling Language Toolkit")
	cmd := sysl.Command("foo", "")
	cmd.Flag("bar", "").Default("foo")

	EnsureFlagsNonEmpty(cmd)

	args := []string{"foo", "--bar", ""}
	_, err := sysl.ParseContext(args)
	assert.Error(t, err)
}
