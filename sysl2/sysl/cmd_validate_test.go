package main

import (
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/alecthomas/kingpin.v2"
)

func TestValidatorDoValidate(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		args     []string
		isErrNil bool
	}{
		"Success": {
			args: []string{
				"sysl2", "validate", "--root-transform", "./tests", "--transform", "transform2.sysl", "--grammar",
				"./tests/grammar.sysl", "--start", "goFile"}, isErrNil: true},
		"Grammar loading fail": {
			args: []string{
				"sysl2", "validate", "--root-transform", "./tests", "--transform", "transform2.sysl", "--grammar",
				"./tests/go.sysl", "--start", "goFile"}, isErrNil: false},
		"Transform loading fail": {
			args: []string{
				"sysl2", "validate", "--root-transform", "./tests", "--transform", "tfm.sysl", "--grammar",
				"./tests/grammar.sysl", "--start", "goFile"}, isErrNil: false},
		"Has validation messages": {
			args: []string{
				"sysl2", "validate", "--root-transform", "./tests", "--transform", "transform1.sysl", "--grammar",
				"./tests/grammar.sysl", "--start", "goFile"}, isErrNil: false},
	}

	for name, tt := range cases {
		args := tt.args
		isErrNil := tt.isErrNil
		t.Run(name, func(t *testing.T) {
			sysl := kingpin.New("sysl", "System Modelling Language Toolkit")

			cmd := &validateCmd{}
			require.NotNil(t, cmd.Configure(sysl))

			var selectedCmd string
			var err error
			if selectedCmd, err = sysl.Parse(args[1:]); err != nil {
				assert.FailNow(t, "Failed to parse args")
			}
			require.Equal(t, cmd.Name(), selectedCmd)
			l, _ := test.NewNullLogger()
			execArgs := ExecuteArgs{
				Module:        nil,
				ModuleAppName: "",
				Filesystem:    afero.NewOsFs(),
				Logger:        l,
			}
			err = cmd.Execute(execArgs)
			if isErrNil {
				assert.Nil(t, err, "Unexpected result")
			} else {
				assert.NotNil(t, err, "Unexpected result")
			}
		})
	}
}
