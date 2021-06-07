package relmod

import (
	"context"
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

	ctx := context.Background()

	m, err := NormalizeSpec(ctx, "..", "all.sysl")
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

func TestNormalize_serializeStmtAction(t *testing.T) {
	t.Parallel()

	appName := []string{"foo"}
	action := sysl.Statement_Action{Action: &sysl.Action{Action: "action"}}
	m := buildStmt(appName, &sysl.Statement{Stmt: &action})

	assertNormalizes(t, m, Schema{
		App:  []App{{AppName: appName}},
		Ep:   []Endpoint{{AppName: appName, EpName: "ep"}},
		Stmt: []Statement{{AppName: appName, EpName: "ep", StmtAction: "action", StmtIndex: []int{0}}},
	})
}

func TestNormalize_serializeStmtCall(t *testing.T) {
	t.Parallel()

	appName := []string{"foo"}
	targetApp := []string{"bar"}
	targetEp := "targetEp"
	m := buildStmt(appName, &sysl.Statement{Stmt: &sysl.Statement_Call{Call: &sysl.Call{
		Target:   &sysl.AppName{Part: targetApp},
		Endpoint: targetEp,
	}}})

	assertNormalizes(t, m, Schema{
		App: []App{{AppName: appName}},
		Ep:  []Endpoint{{AppName: appName, EpName: "ep"}},
		Stmt: []Statement{{AppName: appName, EpName: "ep", StmtIndex: []int{0}, StmtCall: tuple{
			"appName": targetApp,
			"epName":  targetEp,
		}}},
	})
}

func TestNormalize_serializeStmtAlt(t *testing.T) {
	t.Parallel()

	appName := []string{"foo"}
	m := &sysl.Module{Apps: map[string]*sysl.Application{"foo": {
		Name: &sysl.AppName{Part: appName},
		Endpoints: map[string]*sysl.Endpoint{"ep": {
			Name: "ep",
			Stmt: []*sysl.Statement{
				{Stmt: &sysl.Statement_Alt{Alt: &sysl.Alt{Choice: []*sysl.Alt_Choice{
					{
						Cond: "if",
						Stmt: []*sysl.Statement{
							{Stmt: &sysl.Statement_Action{Action: &sysl.Action{Action: "0.0"}}},
							{Stmt: &sysl.Statement_Action{Action: &sysl.Action{Action: "0.1"}}},
						},
					},
					{
						Cond: "else",
						Stmt: []*sysl.Statement{
							{Stmt: &sysl.Statement_Action{Action: &sysl.Action{Action: "1.0"}}},
							{Stmt: &sysl.Statement_Action{Action: &sysl.Action{Action: "1.1"}}},
						},
					},
				}}}},
			},
		}},
	}}}

	assertNormalizes(t, m, Schema{
		App: []App{{AppName: appName}},
		Ep:  []Endpoint{{AppName: appName, EpName: "ep"}},
		Stmt: []Statement{
			{AppName: appName, EpName: "ep", StmtAlt: tuple{"choice": "if"}, StmtIndex: []int{0, 0}},
			{AppName: appName, EpName: "ep", StmtAction: "0.0", StmtIndex: []int{0, 0, 0}},
			{AppName: appName, EpName: "ep", StmtAction: "0.1", StmtIndex: []int{0, 0, 1}},
			{AppName: appName, EpName: "ep", StmtAlt: tuple{"choice": "else"}, StmtIndex: []int{0, 1}},
			{AppName: appName, EpName: "ep", StmtAction: "1.0", StmtIndex: []int{0, 1, 0}},
			{AppName: appName, EpName: "ep", StmtAction: "1.1", StmtIndex: []int{0, 1, 1}},
		},
	})
}

func buildStmt(appName []string, stmt *sysl.Statement) *sysl.Module {
	return &sysl.Module{Apps: map[string]*sysl.Application{"foo": {
		Name: &sysl.AppName{Part: appName},
		Endpoints: map[string]*sysl.Endpoint{"ep": {
			Name: "ep",
			Stmt: []*sysl.Statement{stmt},
		}},
	}}}
}

func assertNormalizes(t *testing.T, m *sysl.Module, s Schema) {
	ctx := context.Background()

	out, err := Normalize(ctx, m)
	require.NoError(t, err)
	actual, err := rel.NewValue(out)
	require.NoError(t, err)

	expected, err := rel.NewValue(s)
	require.NoError(t, err)

	rel.AssertEqualValues(t, expected, actual)
}
