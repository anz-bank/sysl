package importer

import (
	"bytes"
	"encoding/json"

	"github.com/sirupsen/logrus"

	"github.com/go-openapi/loads"
)

func LoadSwaggerText(args OutputData, text string, logger *logrus.Logger) (out string, err error) {
	doc, err := loads.Analyzed(json.RawMessage(text), "2.0")
	if err != nil {
		logger.Errorf("Failed to load swagger spec: %s\n", err.Error())
		return "", err
	}

	result := &bytes.Buffer{}

	swagger := doc.Spec()
	types := InitTypes(swagger, logger)
	globalParams := buildGlobalParams(swagger.Parameters, types, logger)
	endpoints := InitEndpoints(swagger, types, globalParams, logger)
	info := SyslInfo{
		OutputData:  args,
		Title:       "",
		Description: "",
		OtherFields: []string{},
	}
	if swagger.Info != nil {
		info.Title = swagger.Info.Title
		info.Description = swagger.Info.Description
		values := []string{
			"version", swagger.Info.Version,
			"host", swagger.Host,
			"license", "",
			"termsOfService", swagger.Info.TermsOfService}
		for i := 0; i < len(values); i += 2 {
			key := values[i]
			val := values[i+1]
			if val != "" {
				info.OtherFields = append(info.OtherFields, key, val)
			}
		}
	}

	w := newWriter(result, logger)
	if err := w.Write(info, types, swagger.BasePath, endpoints...); err != nil {
		return "", err
	}

	return result.String(), nil
}
