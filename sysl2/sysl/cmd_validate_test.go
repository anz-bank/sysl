package main

import (
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
	"testing"
)

func TestValidatorDoValidate(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		args     []string
		isErrNil bool
	}{
		"Success": {
			args: []string{
				"sysl2", "validate", "--root-transform", "../tests", "--transform", "transform2.sysl", "--grammar",
				"../tests/grammar.sysl", "--start", "goFile"}, isErrNil: true},
		"Grammar loading fail": {
			args: []string{
				"sysl2", "validate", "--root-transform", "../tests", "--transform", "transform2.sysl", "--grammar",
				"../tests/go.sysl", "--start", "goFile"}, isErrNil: false},
		"Transform loading fail": {
			args: []string{
				"sysl2", "validate", "--root-transform", "../tests", "--transform", "tfm.sysl", "--grammar",
				"../tests/grammar.sysl", "--start", "goFile"}, isErrNil: false},
		"Has validation messages": {
			args: []string{
				"sysl2", "validate", "--root-transform", "../tests", "--transform", "transform1.sysl", "--grammar",
				"../tests/grammar.sysl", "--start", "goFile"}, isErrNil: false},
	}

	for name, tt := range cases {
		args := tt.args
		isErrNil := tt.isErrNil
		t.Run(name, func(t *testing.T) {
			sysl := kingpin.New("sysl", "System Modelling Language Toolkit")
			v := &validateCmd{}
			v.Init(sysl)
			logger, _ := test.NewNullLogger()
			if _, err := sysl.Parse(args[1:]); err != nil {
				assert.FailNow(t, "Failed to parse args")
			}
			err := v.Execute(ExecuteArgs{
				module:     nil,
				modAppName: "",
				fs:         syslutil.NewChrootFs(afero.NewOsFs(), "."),
				logger:     logger,
			})
			if isErrNil {
				assert.Nil(t, err, "Unexpected result")
			} else {
				assert.NotNil(t, err, "Unexpected result")
			}
		})
	}
}

