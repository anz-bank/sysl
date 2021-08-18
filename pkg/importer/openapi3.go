package importer

import (
	"github.com/sirupsen/logrus"
)

func NewOpenAPIV3Importer(logger *logrus.Logger) Importer {
	return MakeArraiImporterImporter(ArraiImporterDir+"/openapi/import_cli.arraiz", logger)
}
