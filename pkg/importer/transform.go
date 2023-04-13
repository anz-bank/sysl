package importer

import (
	"os"
	"path"
	"strings"

	"github.com/anz-bank/sysl/internal/bundles"
	"github.com/anz-bank/sysl/pkg/arrai/transform"
	"github.com/arr-ai/arrai/rel"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func MakeTransformImporter(logger *logrus.Logger, transformName string) *TransformImporter {
	return &TransformImporter{
		transformName: transformName,
		logger:        logger,
	}
}

// TransformImporter enables importing from various formats by running embedded arr.ai transform scripts that convert
// data into Sysl.
type TransformImporter struct {
	appName       string
	logger        *logrus.Logger
	transformName string
}

// Configure allows the imported Sysl application name, package and import directories to be specified.
func (i *TransformImporter) Configure(arg *ImporterArg) (Importer, error) {
	if arg.AppName == "" {
		return nil, errors.New("application name not provided")
	}
	i.appName = arg.AppName
	return i, nil
}

// LoadFile generates a Sysl spec by invoking the arr.ai transform.
func (i *TransformImporter) LoadFile(path string) (string, error) {
	bs, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return i.Load(string(bs))
}

// Load generates a Sysl spec given the content of an input file.
func (i *TransformImporter) Load(content string) (string, error) {
	input, err := rel.NewValue(map[string]interface{}{
		"appName": i.appName,
		"content": content,
	})
	if err != nil {
		return "", err
	}

	scriptPath := path.Join("importers", strings.ToLower(i.transformName), "transform.arraiz")
	result, err := transform.EvalWithParam(bundles.MustRead(scriptPath), scriptPath, input)
	if result != nil {
		return result.String(), err
	}
	return "", err
}
