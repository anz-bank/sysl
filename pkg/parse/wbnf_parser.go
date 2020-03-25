package parse

import (
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/spf13/afero"
)

// Parse the given source definition, modifying the app contents in the process.
// Return the list of import definitions declared in the source but not parsed.
func parseWbnf(source importDef, fs afero.Fs, apps map[string]*sysl.Application) ([]importDef, error) {
	return nil, Exitf(ParseError, "wbnf parser unimplemented")
}
