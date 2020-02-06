package catalog

import (
	"testing"

	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

const testDir = "../../tests/"

var filenames = []string{"grpc_catalog.sysl"}

func TestCatlog(t *testing.T) {
	for _, filename := range filenames {
		require.NoError(t, loadCatalogFromFile(filename))
	}
}

func loadCatalogFromFile(filename string) error {
	module, err := parse.NewParser().Parse(filename,
		syslutil.NewChrootFs(afero.NewOsFs(), testDir))
	if err != nil {
		return err
	}

	fs := afero.NewMemMapFs()
	log := logrus.New()

	c := Server{
		Host:    ":8080",
		Fs:      fs,
		Log:     log,
		Modules: []*sysl.Module{module},
	}
	return c.RegisterModules()
}
