package bundles

import (
	"embed"
	"sync"
)

var (
	OpenAPIImporter = &BundledFile{path: "assets/import_openapi_cli.arraiz"}
	SQLImporter     = &BundledFile{path: "assets/import_sql_cli.arraiz"}
	ProtoImporter   = &BundledFile{path: "assets/import_proto_cli.arraiz"}
	Transformer     = &BundledFile{path: "assets/transformer_cli.arraiz"}
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

type BundledFile struct {
	path  string
	once  sync.Once
	bytes []byte
}

func (p *BundledFile) Bytes() []byte {
	p.once.Do(func() {
		p.bytes = MustRead(p.path)
	})
	return p.bytes
}
