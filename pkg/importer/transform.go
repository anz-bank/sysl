package importer

import (
	"io/ioutil"
	"path"

	"github.com/anz-bank/sysl/pkg/arrai/transform"
	"github.com/anz-bank/sysl/transforms"
	"github.com/arr-ai/arrai/rel"
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

// WithAppName allows the imported Sysl application name to be specified.
func (i *TransformImporter) WithAppName(appName string) Importer {
	i.appName = appName
	return i
}

// WithPackage allows the imported Sysl package attribute to be specified.
func (i *TransformImporter) WithPackage(_ string) Importer {
	return i
}

// WithImports allows the imported Sysl import paths attribute to be specified
func (i *TransformImporter) WithImports(_ string) Importer {
	return i
}

// LoadFile generates a Sysl spec by invoking the arr.ai transform.
func (i *TransformImporter) LoadFile(path string) (string, error) {
	bs, err := ioutil.ReadFile(path)
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

	transformBytes, err := transforms.EmbedFs.ReadFile(path.Join("importers", i.transformName, "/transform.arrai"))
	if err != nil {
		return "", err
	}

	result, err := transform.EvalWithParam(string(transformBytes), input)

	if result != nil {
		return result.String(), err
	}
	return "", err
}
