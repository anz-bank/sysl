package printer_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/alecthomas/assert"

	"github.com/spf13/afero"

	"github.com/anz-bank/sysl/pkg/loader"
	"github.com/anz-bank/sysl/pkg/printer"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/sirupsen/logrus"
)

func TestPrinting(t *testing.T) {
	t.Parallel()

	_, fs := syslutil.WriteToMemOverlayFs("tests")
	log := logrus.Logger{}

	module, _, err := loader.LoadSyslModule("/", "printer.sysl", fs, &log)
	assert.NoError(t, err)

	fileBytes, err := afero.ReadFile(fs, "printer.sysl")
	assert.NoError(t, err)
	var buf bytes.Buffer
	printer.Module(&buf, module)
	if buf.String() != string(fileBytes) {
		fmt.Println(buf.String())
	}
	assert.Equal(t, buf.String(), string(fileBytes))
}

func TestExample(t *testing.T) {
	t.Parallel()

	// Create Sysl sile in Memory file system
	fs := afero.NewMemMapFs()
	f, err := fs.Create("/test.sysl")
	if err != nil {
		panic(err)
	}
	_, err = f.Write([]byte(`
Server[~yay]:
    !type Foo:
        foo <: sequence of string
	Endpoint(req <: Foo):
		return ok <: Foo
`))
	if err != nil {
		panic(err)
	}

	// Load Module
	module, _, err := loader.LoadSyslModule("/", "test.sysl", fs, logrus.New())
	if err != nil {
		panic(err)
	}

	// Make a New printer to os.Stdout (io.Writer) and Module
	printer.Module(os.Stdout, module)
}
