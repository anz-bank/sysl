package parse

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type fsFileStream struct {
	*antlr.InputStream
	filename string
}
