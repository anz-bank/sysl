package exporter

import (
	"fmt"
	"os"

	"github.com/anz-bank/sysl/internal/bundles"
	"github.com/anz-bank/sysl/pkg/arrai"
	"github.com/anz-bank/sysl/pkg/loader"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type SpannerExporter struct {
	fs       afero.Fs
	log      *logrus.Logger
	rootPath string
	outPath  string
}

// MakeSpannerExporter returns a new SpannerExporter.
func MakeSpannerExporter(fs afero.Fs, logger *logrus.Logger, rootPath, outPath string) *SpannerExporter {
	return &SpannerExporter{fs: fs, log: logger, rootPath: rootPath, outPath: outPath}
}

// ExportFile reads in a Sysl file from path, converts it to the output format, and writes it to
// the file system.
func (s *SpannerExporter) ExportFile(inPath string) error {
	err := loader.EnsureSyslPb(s.fs, s.rootPath, inPath)
	if err != nil {
		return err
	}

	out, err := s.transform(inPath)
	if err != nil {
		return err
	}
	err = afero.WriteFile(s.fs, s.outPath, []byte(out), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// transform converts the Sysl at path to an SQL output string.
func (s *SpannerExporter) transform(path string) (string, error) {
	b := bundles.SpannerExporter
	out, err := arrai.EvaluateBundle(b, path)
	if err != nil {
		return "", errors.Wrap(arrai.ExecutionError{
			Context:  fmt.Sprintf("export(`%s`)", path),
			Err:      err,
			ShortMsg: "Error executing SQL exporter",
		}, "Executing arrai transform failed")
	}
	return out.String(), nil
}

func (s *SpannerExporter) ExportApp(*sysl.Application) error {
	return errors.Errorf("ExportApp not implemented for Spanner exporter")
}
