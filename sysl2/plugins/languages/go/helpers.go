package helpers

import (
	g "github.com/anz-bank/sysl/sysl2/codegen/golang"
	"github.com/anz-bank/sysl/sysl2/plugins"
	"github.com/sirupsen/logrus"
)

// SingleGoFileGenerateCodeResponse creates a GenerateCodeResponse for a single
// Go file represented as an AST.
func SingleGoFileGenerateCodeResponse(goFile g.File) (*plugins.GenerateCodeResponse, error) {
	src, err := g.Format("", &goFile)
	if err != nil {
		logrus.Errorf("------------\n%s\n--------------", goFile)
		return nil, err
	}

	return plugins.SingleFileGenerateCodeResponse("", string(src)), nil
}
