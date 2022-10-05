package importer

import (
	"github.com/anz-bank/golden-retriever/reader"
	"github.com/anz-bank/sysl/internal/bundles"
	"github.com/sirupsen/logrus"
)

func NewOpenAPIV3Importer(logger *logrus.Logger, reader reader.Reader) Importer {
	return MakeArraiImporterImporter(bundles.OpenAPIImporter, logger, reader)
}
