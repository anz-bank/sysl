package printer

import (
	"bytes"
	"testing"

	"github.com/alecthomas/assert"

	"github.com/spf13/afero"

	"github.com/anz-bank/sysl/pkg/loader"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/sirupsen/logrus"
)

func TestPrinting(t *testing.T) {
	_, fs := syslutil.WriteToMemOverlayFs("../../tests")
	log := logrus.Logger{}

	module, _, _ := loader.LoadSyslModule("/", "printer.sysl", fs, &log)

	fileBytes, _ := afero.ReadFile(fs, "printer.sysl")

	var buf bytes.Buffer
	pr := NewPrinter(&buf)
	pr.PrintModule(module)
	assert.Equal(t, buf.String(), string(fileBytes))
}
