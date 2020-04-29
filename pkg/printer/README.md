# printer
prints out sysl proto data models back to source code; a "reverse parser"

Usage:
```go
package main
import (
"os"

"github.com/spf13/afero"
"github.com/anz-bank/sysl/pkg/printer"
"github.com/anz-bank/sysl/pkg/loader"
"github.com/sirupsen/logrus"
)
func main(){
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
	printer.NewPrinter(os.Stdout).PrintModule(module)
}

```
And on stdout should print:
```
Server[~yay]:
    !type Foo:
        foo <: sequence of string
    Endpoint(req <: Foo):
        return ok <: Foo
```
