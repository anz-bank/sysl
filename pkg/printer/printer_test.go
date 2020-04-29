package printer

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/alecthomas/assert"

	"github.com/spf13/afero"

	"github.com/anz-bank/sysl/pkg/loader"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/sirupsen/logrus"
)

func TestPrinting(t *testing.T) {
	_, fs := syslutil.WriteToMemOverlayFs("tests")
	log := logrus.Logger{}

	module, _, err := loader.LoadSyslModule("/", "printer.sysl", fs, &log)
	assert.NoError(t, err)

	fileBytes, err := afero.ReadFile(fs, "printer.sysl")
	assert.NoError(t, err)
	var buf bytes.Buffer
	pr := NewPrinter(&buf)
	//pr := NewPrinter(os.Stdout)
	pr.PrintModule(module)
	if buf.String() != string(fileBytes) {
		fmt.Println(buf.String())
	}
	assert.Equal(t, buf.String(), string(fileBytes))
}


func TestExample(t *testing.T) {

	// Create Sysl sile in Memory file system
	fs := afero.NewMemMapFs()
	f, _ := fs.Create("/test.sysl")
	f.Write([]byte(`
Server[~yay]:
    !type Foo:
        foo <: sequence of string
	Endpoint(req <: Foo):
		return ok <: Foo
`))

	// Load Module
	module, _, _ := loader.LoadSyslModule("/", "test.sysl", fs,logrus.New())

	// Make a New printer to os.Stdout (io.Writer) and PrintModule
	NewPrinter(os.Stdout).PrintModule(module)
}