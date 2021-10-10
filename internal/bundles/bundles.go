package bundles

import _ "embed"

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

	// SpannerExporter returns the bytes that represents the bundled script for spanner exporter.
	//go:embed assets/spanner_cli.arraiz
	SpannerExporter []byte
)
