//nolint:lll
package exporter

import (
	"bytes"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/sirupsen/logrus"
)

func TestSpannerExport(t *testing.T) {
	t.Parallel()
	outPath := "test-data/spanner/sample.sql"

	x := MakeTransformExporter(afero.NewOsFs(), &logrus.Logger{}, "", outPath, "spanner")
	b := &bytes.Buffer{}

	syslPath := "exporter/test-data/spanner/sample.sysl"

	fs := syslutil.NewChrootFs(afero.NewOsFs(), "..")
	mod, _, err := parse.LoadAndGetDefaultApp(syslPath, fs, parse.NewParser())
	require.NoError(t, err)

	err = x.ExportToWriter(b, []*sysl.Module{mod}, []string{syslPath})
	require.NoError(t, err)

	expected, err := afero.ReadFile(fs, "exporter/test-data/spanner/sample.sql")
	require.NoError(t, err)

	assert.Equal(t, string(expected), b.String())
}
