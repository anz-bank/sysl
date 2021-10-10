package importer

import (
	"github.com/anz-bank/sysl/internal/bundles"
	"github.com/sirupsen/logrus"
)

func NewOpenAPIV3Importer(logger *logrus.Logger) Importer {
	return MakeArraiImporterImporter(bundles.OpenAPIImporter, logger)
}
