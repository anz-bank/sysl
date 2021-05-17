package relmod

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/arr-ai/arrai/syntax"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/arr-ai/arrai/rel"
	"github.com/stretchr/testify/require"
)

func TestNormalize(t *testing.T) {
	t.Parallel()

	m, err := NormalizeSpec("..", "all.sysl")
	require.NoError(t, err)
	v, err := rel.NewValue(m)
	require.NoError(t, err)
	actual, err := syntax.PrettifyString(v, 0)
	require.NoError(t, err)

	expected, err := ioutil.ReadFile("../out/all_normalize.arrai")
	require.NoError(t, err)
	err = ioutil.WriteFile("../out/tmp_all_normalize.arrai", []byte(actual), os.ModePerm)
	require.NoError(t, err)

	// TODO: Test that the output is correct. For now just test that it didn't error out.
	assert.NotNil(t, actual)
	assert.NotNil(t, expected)
}

func TestNormalize_serializeApp(t *testing.T) {
	t.Parallel()

	appName := []string{"foo"}
	m := &sysl.Module{Apps: map[string]*sysl.Application{"foo": {
		Name: &sysl.AppName{Part: appName},
	}}}

	assertNormalizes(t, m, Schema{App: []App{{AppName: appName}}})
}

func TestNormalize_serializeStmt(t *testing.T) {
	t.Parallel()

	appName := []string{"foo"}
	m := &sysl.Module{Apps: map[string]*sysl.Application{"foo": {
		Name: &sysl.AppName{Part: appName},
		Endpoints: map[string]*sysl.Endpoint{"ep": {
			Name: "ep",
			Stmt: []*sysl.Statement{
				{Stmt: &sysl.Statement_Action{Action: &sysl.Action{Action: "action"}}},
			},
		}},
	}}}

	assertNormalizes(t, m, Schema{
		App:  []App{{AppName: appName}},
		Ep:   []Endpoint{{AppName: appName, EpName: "ep"}},
		Stmt: []Statement{{AppName: appName, EpName: "ep", StmtAction: "action"}},
	})
}

func assertNormalizes(t *testing.T, m *sysl.Module, s Schema) {
	out, err := Normalize(m)
	require.NoError(t, err)
	actual, err := rel.NewValue(out)
	require.NoError(t, err)

	expected, err := rel.NewValue(s)
	require.NoError(t, err)

	rel.AssertEqualValues(t, expected, actual)
}
