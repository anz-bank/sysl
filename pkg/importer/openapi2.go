package importer

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
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
	bs, err := os.ReadFile(path)
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

// Configure allows the imported Sysl application name, package and import directories to be specified.
func (l *OpenAPI2Importer) Configure(arg *ImporterArg) (Importer, error) {
	if arg.AppName == "" {
		return nil, errors.New("application name not provided")
	}
	l.appName = arg.AppName
	l.pkg = arg.PackageName
	return l, nil
}

// convertToOpenAPI3 takes a openapi2 spec and converts it to openapi3
func convertToOpenAPI3(data []byte, uri *url.URL, loader *openapi3.Loader) (*openapi3.T, string, error) {
	var openapiv2 openapi2.T
	jsondata, err := yaml.YAMLToJSON(data)
	if err != nil {
		return nil, "", err
	}
	// kin-openapi has a bug where it tries to unmarshal discriminator as openapi3 spec rather than v2
	// do a find and replace to work around the issue (see https://github.com/getkin/kin-openapi/issues/360)
	reg := regexp.MustCompile(`"discriminator":(\s*)(".*?")`)
	jsondata = reg.ReplaceAll(jsondata, []byte(`"discriminator":$1{"propertyName":$1$2}`))
	err = json.Unmarshal(jsondata, &openapiv2)
	if err != nil {
		return nil, "", err
	}

	if len(openapiv2.Schemes) == 0 {
		openapiv2.Schemes = []string{"https"}
	}
	openapiv3, err := openapi2conv.ToV3WithLoader(&openapiv2, loader, uri)
	if err != nil {
		return nil, "", err
	}

	return openapiv3, openapiv2.BasePath, nil
}
