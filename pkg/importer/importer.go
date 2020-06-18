package importer

// Importer is an interface implemented by all sysl importers
type Importer interface {
	// Load takes in a string in a format supported by an the importer
	// It returns the converted Sysl as a string
	Load(file string) (string, error)
	// WithAppName allows the exported Sysl application name to be specified
	WithAppName(appName string) Importer
	// WithPackage allows the exported Sysl package attribute to be specified
	WithPackage(packageName string) Importer
}
