package exporter

import (
	"io/ioutil"
	"testing"

	"github.com/anz-bank/sysl/sysl2/sysl/parse"
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

// nolint:gochecknoglobals
var fileList = []string{
	"exporter/test-data/EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH",
	"exporter/test-data/EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_200_RESPONSE_DESCRIPTION_ONLY",
	"exporter/test-data/EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_201_LOCATION_HEADER_RESPONSE",
	"exporter/test-data/EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_201_RESPONSE_DESCRIPTION_ONLY",
	"exporter/test-data/EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_BODY_PARAMETER",
	"exporter/test-data/EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_DEFAULT_RESPONSE",
	"exporter/test-data/EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_ERROR_RESPONSE",
	"exporter/test-data/EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_PATH_WITH_X_DASH_WHATEVER_RESPONSE",
	"exporter/test-data/EXAMPLE_SWAGGER_SPEC_WITH_ENDPOINT_RETURNING_ARRAY_OF_DEFINED_OBJECT_TYPE",
	"exporter/test-data/SIMPLE_SWAGGER_EXAMPLE",
	"exporter/test-data/SWAGGER_ENSURE_HEADER_FIELDS_IN_INSERTION_ORDER",
	"exporter/test-data/SWAGGER_HEADER_AND_BODY_PARAM_EXAMPLE",
	"exporter/test-data/SWAGGER_QUERY_PARAM_EXAMPLE",
	"exporter/test-data/SWAGGER_REQUIRED_AND_OPTIONAL_FIELDS_EXAMPLE",
	"exporter/test-data/SWAGGER_TOP_LEVEL_ARRAY_EXAMPLE",
	"exporter/test-data/SWAGGER_WITH_ENUMS",
	"exporter/test-data/SWAGGER_WITH_HEADER_VAR_OVERRIDDEN_IN_METHOD",
	"exporter/test-data/SWAGGER_WITH_PATHS_VAR_REFERRING_GLOBAL_PARAMS_OBJECT",
	"exporter/test-data/SWAGGER_WITH_PATH_VAR_TYPE_IN_API",
	"exporter/test-data/SWAGGER_WITH_PATH_VAR_TYPE_IN_GLOBAL_PARAMETERS",
	"exporter/test-data/SWAGGER_WITH_PATH_VAR_TYPE_OVERRIDDEN_IN_SECOND_METHOD",
	"exporter/test-data/SWAGGER_WITH_PATHS_VAR_REFERRING_GLOBAL_PARAMS_OBJECT",
}

func TestExportAll(t *testing.T) {
	t.Parallel()
	modelParser := parse.NewParser()
	for _, filename := range fileList {
		mod, _, err1 := parse.LoadAndGetDefaultApp(filename+`.sysl`, syslutil.NewChrootFs(afero.NewOsFs(), ".."), modelParser)
		require.NoError(t, err1)
		if err1 != nil {
			t.Errorf("Error reading sysl %s", fileList[0])
		}
		swaggerExporter := MakeSwaggerExporter(mod.GetApps()["testapp"], logrus.StandardLogger())
		err2 := swaggerExporter.GenerateSwagger()
		require.NoError(t, err2)
		out, err := swaggerExporter.SerializeToYaml()
		require.NoError(t, err)

		yamlFileBytes, err := ioutil.ReadFile("../" + filename + `.yaml`)
		require.NoError(t, err)
		if string(yamlFileBytes) != string(out) {
			t.Errorf("Content mismatched\n%s\n*******\n%s", string(yamlFileBytes), string(out))
		}
	}
}
