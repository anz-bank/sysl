package exporter

import (
	"github.com/anz-bank/sysl/pkg/sysl"
)

// Exporter is an interface implemented by all Sysl exporters.
type Exporter interface {
	// ExportFile reads in a Sysl file from path, converts it to the output format, and writes it to
	// the file system.
	ExportFile(path string) error
	// ExportApp takes a Sysl app, converts it to the output format, and writes it to the file
	// system.
	ExportApp(app *sysl.Application) error
}
