package exporter

import (
	"bytes"
	"io"
	"os"
	"path"
	"strings"

	"github.com/anz-bank/sysl/internal/bundles"
	"github.com/anz-bank/sysl/pkg/arrai/transform"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func MakeTransformExporter(
	fs afero.Fs,
	logger *logrus.Logger,
	rootPath, outPath,
	transformName string,
) *TransformExporter {
	return &TransformExporter{
		fs:            fs,
		log:           logger,
		rootPath:      rootPath,
		outPath:       outPath,
		transformName: transformName,
	}
}

// TransformExporter enables exporting into various formats by running embedded arr.ai transform scripts that convert
// Sysl.
type TransformExporter struct {
	fs            afero.Fs
	log           *logrus.Logger
	rootPath      string
	outPath       string
	transformName string
}

// ExportFile reads in a Sysl file from path, converts it to the output format, and writes it to
// the file system.
func (e *TransformExporter) ExportFile(modules []*sysl.Module, modulePaths []string) error {
	b := &bytes.Buffer{}
	if err := e.ExportToWriter(b, modules, modulePaths); err != nil {
		return err
	}
	return afero.WriteFile(e.fs, e.outPath, b.Bytes(), os.ModePerm)
}

func (e *TransformExporter) ExportToWriter(w io.Writer, modules []*sysl.Module, modulePaths []string) error {
	input, err := transform.BuildTransformInput(modules, modulePaths)
	if err != nil {
		return err
	}

	scriptPath := path.Join("exporters", strings.ToLower(e.transformName), "transform.arraiz")
	result, err := transform.EvalWithParam(bundles.MustRead(scriptPath), scriptPath, input)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(result.String()))
	return err
}

func (e *TransformExporter) ExportApp(*sysl.Application) error {
	return errors.Errorf("ExportApp not implemented for transform exporter")
}
