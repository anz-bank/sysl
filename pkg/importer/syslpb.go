package importer

import (
	"bytes"

	"github.com/anz-bank/sysl/pkg/pbutil"
	"github.com/anz-bank/sysl/pkg/printer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type SyslPBImporter struct {
	logger *logrus.Logger
}

// MakeSyslPBImporter is a factory method for creating new SyslPB importer.
func MakeSyslPBImporter(logger *logrus.Logger) *SyslPBImporter {
	return &SyslPBImporter{logger: logger}
}

func (i *SyslPBImporter) LoadFile(path string) (string, error) {
	m, err := pbutil.FromPB(path, afero.NewOsFs())
	if err != nil {
		return "", errors.Wrap(err, "Failed to load Sysl PB from "+path)
	}
	var buf bytes.Buffer
	printer.Module(&buf, m)
	return buf.String(), nil
}

// Load returns a Sysl spec equivalent to protoSpec.
func (i *SyslPBImporter) Load(content string) (string, error) {
	m, err := pbutil.FromPBStringContents("import.pb", content)
	if err != nil {
		return "", errors.Wrap(err, "Failed to load Sysl PB from bytes")
	}
	var buf bytes.Buffer
	printer.Module(&buf, m)
	return buf.String(), nil
}

// Configure for Sysl PB files is a no-op; they are simply imported as-is.
func (i *SyslPBImporter) Configure(arg *ImporterArg) (Importer, error) {
	return i, nil
}
