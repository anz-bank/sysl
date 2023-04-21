package bundles

import "embed"

// TODO: refactor these into transform plugins and move them into exporters or importers.
var (
	// OpenAPIImporter returns the bytes that represents the bundled script for openapi importer.
	//go:embed assets/import_openapi_cli.arraiz
	OpenAPIImporter []byte

	// SQLImporter returns the bytes that represents the bundled script for SQL importer.
	//go:embed assets/import_sql_cli.arraiz
	SQLImporter []byte

	// ProtoImporter returns the bytes that represents the bundled script for SQL importer.
	//go:embed assets/import_proto_cli.arraiz
	ProtoImporter []byte

	// Transformer returns the bytes that represents the bundled script for avro importer.
	//go:embed assets/transformer_cli.arraiz
	Transformer []byte
)

// BundlesFs is a filesystem that can be used to dynamically load the available bundles.
//
//go:embed *
var BundlesFs embed.FS

func MustRead(path string) []byte {
	b, err := BundlesFs.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return b
}
