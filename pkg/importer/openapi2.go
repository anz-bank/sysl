package importer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/spf13/afero"

	"github.com/ghodss/yaml"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/sirupsen/logrus"
)

// https://swagger.io/blog/api-strategy/difference-between-swagger-and-openapi/
//   OpenAPI = Specification
//   Swagger = Tools for implementing the specification

type OpenAPI2Importer struct {
	openAPI2Spec string
	filepath     string
	*OpenAPI3Importer
}

func MakeOpenAPI2Importer(logger *logrus.Logger, basePath string, filePath string) *OpenAPI2Importer {
	return &OpenAPI2Importer{OpenAPI3Importer: &OpenAPI3Importer{
		appName:  "",
		pkg:      "",
		basePath: basePath,
		logger:   logger,
		fs:       afero.NewOsFs(),
		types:    TypeList{},
	},
		filepath: filePath}
}

func (l *OpenAPI2Importer) LoadFile(path string) (string, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return l.Load(string(bs))
}

func (l *OpenAPI2Importer) Load(oas2spec string) (string, error) {
	u, err := pathToURL(l.filepath)
	if err != nil {
		return "", err
	}
	oas3spec, basePath, err := convertToOpenAPI3([]byte(oas2spec), u, NewOpenAPI3Loader(l.logger, l.fs))
	if err != nil {
		return "", fmt.Errorf("error converting openapi 2:%w", err)
	}
	l.openAPI2Spec = oas2spec
	l.basePath = basePath
	return l.convertSpec(oas3spec)
}

// Set the AppName of the imported app
func (l *OpenAPI2Importer) WithAppName(appName string) Importer {
	l.appName = appName
	return l
}

// Set the package attribute of the imported app
func (l *OpenAPI2Importer) WithPackage(pkg string) Importer {
	l.pkg = pkg
	return l
}

// Set the importPaths attribute of the imported app
func (l *OpenAPI2Importer) WithImports(_ string) Importer {
	return l
}

// convertToOpenAPI3 takes a openapi2 spec and converts it to openapi3
func convertToOpenAPI3(data []byte, uri *url.URL, loader *openapi3.SwaggerLoader) (*openapi3.Swagger, string, error) {
	var openapiv2 openapi2.Swagger
	jsondata, err := yaml.YAMLToJSON(data)
	if err != nil {
		return nil, "", err
	}
	err = json.Unmarshal(jsondata, &openapiv2)
	if err != nil {
		return nil, "", err
	}

	if len(openapiv2.Schemes) == 0 {
		openapiv2.Schemes = []string{"https"}
	}
	openapiv3, err := openapi2conv.ToV3Swagger(&openapiv2)
	if err != nil {
		return nil, "", err
	}
	// openapi2conv doesnt handle external references correctly so to avoid that we need to serialise the converted
	// v3 doc back to text and manually replace #/components/* references (easier to do it via text instead of walking
	// the whole object tree again
	if oa3text, err := json.Marshal(openapiv3); err == nil {
		replacer := strings.NewReplacer(
			"#/definitions/", "#/components/schemas/",
			"#/responses/", "#/components/responses/",
			"#/parameters/", "#/components/parameters/",
		)
		oa3replaced := replacer.Replace(string(oa3text))
		_ = json.Unmarshal([]byte(oa3replaced), &openapiv3) //nolint: errcheck // we dont care about the error here
	}

	err = loader.ResolveRefsIn(openapiv3, uri)
	if err != nil {
		return nil, "", err
	}

	return openapiv3, openapiv2.BasePath, nil
}
