package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/andreyvit/diff"
	"github.com/anz-bank/sysl/cmd/syslwbnf/parser"
	"github.com/anz-bank/sysl/pkg/pbutil"

	"github.com/anz-bank/sysl/pkg/parse"

	"github.com/spf13/afero"
)

func main() {
	fs := afero.NewOsFs()
	modAntlr, _, err := parse.LoadAndGetDefaultApp(os.Args[1], fs, parse.NewParser())
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	var buf1 bytes.Buffer
	pbutil.FTextPB(&buf1, modAntlr)
	modWbnf, _, err := parser.LoadAndGetDefaultApp(os.Args[1], fs)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	var buf2 bytes.Buffer
	pbutil.FTextPB(&buf2, modWbnf)

	fmt.Print(diff.LineDiff(buf1.String(), buf2.String()))
}
