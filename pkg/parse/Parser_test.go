package parse

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/anz-bank/golden-retriever/retriever"
	"github.com/pmezard/go-difflib/difflib"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/anz-bank/sysl/pkg/env"
	"github.com/anz-bank/sysl/pkg/msg"
	"github.com/anz-bank/sysl/pkg/pbutil"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

var (
	update = flag.Bool("update", false, "Update golden test files")
)

const (
	randomSha = `1e7c4cecaaa8f76e3c668cebc411f1b03174501f`
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

var mainTestDir = filepath.Join("..", "..", "tests")

func readSyslModule(filename string) (*sysl.Module, error) {
	var buf bytes.Buffer

	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "Open file %#v", filename)
	}
	if _, err := io.Copy(&buf, f); err != nil {
		return nil, err
	}
	if err := f.Close(); err != nil {
		return nil, err
	}

	module := &sysl.Module{}
	if err := prototext.Unmarshal(buf.Bytes(), module); err != nil {
		return nil, errors.Wrapf(err, "Unmarshal proto: %s", filename)
	}
	return module, nil
}

func removeSourceContextFromExpr(expr *sysl.Expr) {
	if expr == nil {
		return
	}
	expr.SourceContext = nil //nolint:staticcheck
	// All the Expr types which contain nested expressions
	switch e := expr.Expr.(type) {
	case *sysl.Expr_GetAttr_:
		removeSourceContextFromExpr(e.GetAttr.Arg)
	case *sysl.Expr_Transform_:
		removeSourceContextFromExpr(e.Transform.Arg)
		for _, stmt := range e.Transform.Stmt {
			switch s := stmt.Stmt.(type) {
			case *sysl.Expr_Transform_Stmt_Assign_:
				removeSourceContextFromExpr(s.Assign.Expr)
			case *sysl.Expr_Transform_Stmt_Let:
				removeSourceContextFromExpr(s.Let.Expr)
			case *sysl.Expr_Transform_Stmt_Inject:
				removeSourceContextFromExpr(s.Inject)
			}
		}
	case *sysl.Expr_Ifelse:
		removeSourceContextFromExpr(e.Ifelse.Cond)
		removeSourceContextFromExpr(e.Ifelse.IfFalse)
		removeSourceContextFromExpr(e.Ifelse.IfTrue)
	case *sysl.Expr_Call_:
		for _, arg := range e.Call.Arg {
			removeSourceContextFromExpr(arg)
		}
	case *sysl.Expr_Unexpr:
		removeSourceContextFromExpr(e.Unexpr.Arg)
	case *sysl.Expr_Binexpr:
		removeSourceContextFromExpr(e.Binexpr.Lhs)
		removeSourceContextFromExpr(e.Binexpr.Rhs)
	case *sysl.Expr_Relexpr:
		for _, arg := range e.Relexpr.Arg {
			removeSourceContextFromExpr(arg)
		}
		removeSourceContextFromExpr(e.Relexpr.Target)
	case *sysl.Expr_Navigate_:
		removeSourceContextFromExpr(e.Navigate.Arg)

	case *sysl.Expr_List_:
		for _, arg := range e.List.Expr {
			removeSourceContextFromExpr(arg)
		}
	case *sysl.Expr_Set:
		for _, arg := range e.Set.Expr {
			removeSourceContextFromExpr(arg)
		}
	case *sysl.Expr_Tuple_:
		for _, arg := range e.Tuple.Attrs {
			removeSourceContextFromExpr(arg)
		}
	}
}

func parseComparable(
	filename, root string,
	stripSourceContext bool,
) (*sysl.Module, error) {
	module, err := NewParser().ParseFromFs(filename, syslutil.NewChrootFs(afero.NewOsFs(), root))
	if err != nil {
		return nil, err
	}

	if stripSourceContext {
		// remove stuff that does not match legacy.
		for _, app := range module.Apps {
			app.SourceContext = nil //nolint:staticcheck
			for _, ep := range app.Endpoints {
				ep.SourceContext = nil //nolint:staticcheck
			}
			for _, t := range app.Types {
				t.SourceContext = nil //nolint:staticcheck
			}
			for _, t := range app.Views {
				t.SourceContext = nil
				removeSourceContextFromExpr(t.Expr)
			}
		}
	}
	return module, nil
}

func parseAndCompareSysl(
	filename1, filename2, root string,
) (bool, error) {
	module1, err := parseComparable(filename1, root, true)
	if err != nil {
		return false, err
	}

	module2, err := parseComparable(filename2, root, true)
	if err != nil {
		return false, err
	}

	if proto.Equal(module1, module2) {
		return true, nil
	}
	return checkSyslEqual(module1, module2), nil
}

func checkSyslEqual(model1 *sysl.Module, model2 *sysl.Module) bool {
	if len(model1.Apps) != len(model2.Apps) {
		return false
	}
	for app1, syslApp1 := range model1.Apps {
		app2, exists := model2.Apps[app1]
		if !exists {
			return false
		}
		for endpoint1, syslEndpoint1 := range syslApp1.Endpoints {
			syslEndpoint2, exists := app2.Endpoints[endpoint1]
			if !exists {
				return false
			}
			return reflect.DeepEqual(syslEndpoint1.Stmt[0].Stmt, syslEndpoint2.Stmt[0].Stmt)
		}
	}
	return false
}

var textProtoExtraSpaceAfterColonRE = regexp.MustCompile(`^([ \t]*: ) `)

func parseAndCompare(
	filename, root, golden string,
	goldenProto protoreflect.ProtoMessage,
	retainOnError bool,
	stripSourceContext bool,
) (bool, error) {
	module, err := parseComparable(filename, root, stripSourceContext)
	if err != nil {
		return false, err
	}

	if proto.Equal(goldenProto, module) {
		return true, nil
	}

	expected := bytes.Buffer{}
	actual := bytes.Buffer{}
	if err = pbutil.FTextPB(&expected, goldenProto); err != nil {
		return false, err
	}
	if err = pbutil.FTextPB(&actual, module); err != nil {
		return false, err
	}
	expectedBytes := textProtoExtraSpaceAfterColonRE.ReplaceAll(expected.Bytes(), []byte(`$1`))
	actualBytes := textProtoExtraSpaceAfterColonRE.ReplaceAll(actual.Bytes(), []byte(`$1`))
	diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(expectedBytes)),
		B:        difflib.SplitLines(string(actualBytes)),
		FromFile: "Expected: " + golden,
		FromDate: "",
		ToFile:   "Actual",
		ToDate:   "",
		Context:  1,
	})
	if err != nil {
		return false, err
	}
	if env.SYSL_DEV_UPDATE_GOLDEN_TESTS.On() {
		if err := ioutil.WriteFile(golden, actualBytes, 0600); err != nil {
			logrus.Errorf("Error updating golden file %q: %v", golden, err)
		} else {
			logrus.Errorf("Updated golden file %q", golden)
		}
	}
	return false, errors.New(diff)
}

func parseAndCompareWithGolden(filename, root string, stripSourceContext bool) (bool, error) {
	goldenFilename := filename
	if !strings.HasSuffix(goldenFilename, syslExt) {
		goldenFilename += syslExt
	}
	goldenFilename += ".golden.textpb"
	golden := filepath.Join(root, goldenFilename)

	if *update {
		updated := bytes.Buffer{}
		if !strings.HasSuffix(filename, syslExt) {
			filename += syslExt
		}
		// update test files
		mod, err := NewParser().ParseFromFs(filename, syslutil.NewChrootFs(afero.NewOsFs(), "."))
		if err != nil {
			return false, err
		}
		if err = pbutil.FTextPB(&updated, mod); err != nil {
			return false, err
		}
		err = ioutil.WriteFile(goldenFilename, updated.Bytes(), 0600)
		if err != nil {
			return false, err
		}
	}

	goldenModule, err := readSyslModule(golden)
	if err != nil {
		return false, err
	}
	return parseAndCompare(filename, root, golden, goldenModule, true, stripSourceContext)
}

func testParseAgainstGolden(t *testing.T, filename, root string) {
	equal, err := parseAndCompareWithGolden(filename, root, false)
	if assert.NoError(t, err) {
		assert.True(t, equal, "%#v %#v", root, filename)
	}
}

func testParseAgainstGoldenWithSourceContext(t *testing.T, filename string) {
	equal, err := parseAndCompareWithGolden(filename, "", false)
	if assert.NoError(t, err) {
		assert.True(t, equal, "%#v", filename)
	}
}

func TestParseBadRoot(t *testing.T) {
	t.Parallel()

	_, err := parseComparable("dontcare.sysl", "NON-EXISTENT-ROOT", false)
	assert.Error(t, err)
}

func TestParseMissingFile(t *testing.T) {
	t.Parallel()

	_, err := parseComparable("doesn't.exist.sysl", "tests", false)
	assert.Error(t, err)
}

func TestParseDirectoryAsFile(t *testing.T) {
	t.Parallel()

	dirname := "not-a-file.sysl"
	tmproot := os.TempDir()
	tmpdir := filepath.Join(tmproot, dirname)
	require.NoError(t, os.Mkdir(tmpdir, 0755))
	defer os.Remove(tmpdir)
	_, err := parseComparable(dirname, tmproot, false)
	assert.Error(t, err)
}

func TestParseBadFile(t *testing.T) {
	t.Parallel()

	_, err := parseAndCompareWithGolden("sysl.go", "", false)
	assert.Error(t, err)
}

func TestSimpleEP(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/test1.sysl", "")
}

func TestSimpleEPNoSuffix(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/test1", "")
}

func TestAttribs(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/attribs.sysl", "")
}

func TestIfElse(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/if_else.sysl", "")
}

func TestArgs(t *testing.T) {
	t.Parallel()

	testParseAgainstGoldenWithSourceContext(t, "tests/args.sysl")
}

func TestTypeAsFieldName(t *testing.T) {
	t.Parallel()

	testParseAgainstGoldenWithSourceContext(t, "tests/type_as_field.sysl")
}

func TestSimpleEPWithSpaces(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/with_spaces.sysl", "")
}

func TestSimpleEP2(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/test4.sysl", "")
}

func TestUnion(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/union.sysl", "")
}

func TestSimpleEndpointParams(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/ep_params.sysl", "")
}

func TestOneOfStatements(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/oneof.sysl", "")
}

func TestDuplicateEndpoints(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/duplicate.sysl", "")
}

func TestEventing(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/eventing.sysl", "")
}

func TestCollector(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/collector.sysl", "")
}

func TestPubSubCollector(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/pubsub_collector.sysl", "")
}

func TestDocstrings(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/docstrings.sysl", "")
}

func TestMixins(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/mixin.sysl", "")
}

func TestNamespaceTypes(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/namespace_types.sysl", "")
}

func TestForLoops(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/for_loop.sysl", "")
}

func TestGroupStmt(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/group_stmt.sysl", "")
}

func TestUntilLoop(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/until_loop.sysl", "")
}

func TestTuple(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/test2.sysl", "")
}

func TestInplaceTuple(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/inplace_tuple.sysl", "")
}

func TestNumbersFormat(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/numbers_format.sysl", "")
}

func TestImports(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/library.sysl", "")
}

func TestForeignImports(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/foreign_import_swagger.sysl", "")
}

func TestForeignImportsWithNamespace(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/foreign_import_swagger_namespace.sysl", "")
}

func TestRootArgAndRelational(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/school.sysl", "")
}

func TestSequenceType(t *testing.T) {
	t.Parallel()

	testParseAgainstGoldenWithSourceContext(t, "tests/sequence_type.sysl")
}

func TestRestApi(t *testing.T) {
	t.Parallel()

	testParseAgainstGoldenWithSourceContext(t, "tests/test_rest_api.sysl")
}

func TestRestApiQueryParams(t *testing.T) {
	t.Parallel()

	testParseAgainstGoldenWithSourceContext(t, "tests/rest_api_query_params.sysl")
}

func TestSimpleProject(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/project.sysl", "")
}

func TestUrlParamOrder(t *testing.T) {
	t.Parallel()

	testParseAgainstGoldenWithSourceContext(t, "tests/rest_url_params.sysl")
}

func TestRestApi_WrongOrder(t *testing.T) {
	t.Parallel()

	testParseAgainstGoldenWithSourceContext(t, "tests/bad_order.sysl")
}

func TestTransform(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/transform.sysl", "")
}

func TestTypeRefs(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/type_refs.sysl", "")
}

func TestImpliedDot(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/implied.sysl", "")
}

func TestStmts(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/stmts.sysl", "")
}

func TestMath(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/math.sysl", "")
}

func TestTableof(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/tableof.sysl", "")
}

func TestRank(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/rank.sysl", "")
}

func TestMatching(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/matching.sysl", "")
}

func TestNavigate(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/navigate.sysl", "")
}

func TestFuncs(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/funcs.sysl", "")
}

func TestPetshop(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/petshop.sysl", "")
}

func TestCrash(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/crash.sysl", "")
}

func TestAlias(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/alias_inline.sysl", "")
}

func TestEscapedEndpoints(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/endpoints.sysl", "")
}

func TestStrings(t *testing.T) {
	t.Parallel()

	testParseAgainstGoldenWithSourceContext(t, "tests/strings_expr.sysl")
}

func TestTypeAlias(t *testing.T) {
	t.Parallel()

	testParseAgainstGoldenWithSourceContext(t, "tests/alias.sysl")
}

func TestEnum(t *testing.T) {
	t.Parallel()

	testParseAgainstGoldenWithSourceContext(t, "tests/enum.sysl")
}

func TestMergeAttrs(t *testing.T) {
	t.Parallel()

	testParseAgainstGoldenWithSourceContext(t, "tests/merge_attrs.sysl")
}

func TestMergeAttrsTable(t *testing.T) {
	t.Parallel()

	testParseAgainstGoldenWithSourceContext(t, "tests/merge_attrs_table.sysl")
}

func TestOpenAPI3(t *testing.T) {
	t.Parallel()

	testParseAgainstGoldenWithSourceContext(t, "tests/openapi3.sysl")
}

func TestImportProtoJSON(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/import_proto_JSON.sysl", "")
}

func TestImportProtoJSONConflict(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/import_proto_JSON_conflict.sysl", "")
}

func TestImportProtoJSONMerge(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/import_proto_JSON_merge.sysl", "")
}

func TestImportProtoJSONEmpty(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/import_proto_JSON_empty.sysl", "")
}

func TestImportProtoJSONDep(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/import_proto_JSON_dep.sysl", "")
}

func TestParseOptionalViewParams(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/optional_params_view.sysl", "")
}

func TestImportProtoJSONNonExistent(t *testing.T) {
	t.Parallel()
	_, err := parseComparable("tests/import_proto_JSON_nonexistent.sysl", "tests", false)
	assert.Error(t, err)
}

func TestImportProtoJSONInvalid(t *testing.T) {
	t.Parallel()
	_, err := parseComparable("tests/import_proto_JSON_invalid.sysl", "tests", false)
	assert.Error(t, err)
}

func TestAttrScope(t *testing.T) {
	t.Parallel()

	testParseAgainstGoldenWithSourceContext(t, "tests/attr_scope.sysl")
}

func TestAnnotationMerge(t *testing.T) {
	t.Parallel()

	testParseAgainstGoldenWithSourceContext(t, "tests/anno_merge.sysl")
}

func TestEmptyTable(t *testing.T) {
	t.Parallel()

	testParseAgainstGoldenWithSourceContext(t, "tests/empty_table.sysl")
}

func TestViewAttr(t *testing.T) {
	t.Parallel()

	testParseAgainstGoldenWithSourceContext(t, "tests/view_attr.sysl")
}

func TestUndefinedRootAbsoluteImport(t *testing.T) {
	t.Parallel()

	parser := NewParser()
	parser.RestrictToLocalImport()
	_, err := parser.ParseFromFs("absolute_import.sysl", syslutil.NewChrootFs(afero.NewOsFs(), "tests"))
	assert.NoError(t, err)
}

func TestDefinedRootAbsoluteImport(t *testing.T) {
	t.Parallel()

	parser := NewParser()
	parser.RestrictToLocalImport()

	_, err := parser.ParseFromFs("subfolder/absolute_import.sysl", syslutil.NewChrootFs(afero.NewOsFs(), "tests"))
	require.NoError(t, err)
}

func TestOutsideOfRootImport(t *testing.T) {
	t.Parallel()

	parser := NewParser()
	parser.RestrictToLocalImport()
	_, err := parser.ParseFromFs("outsideroot_import.sysl", syslutil.NewChrootFs(afero.NewOsFs(), "tests"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "error reading \"../alias.sysl\": \npermission denied, file outside root\n")
}

func TestDuplicateImport(t *testing.T) {
	t.Parallel()

	file1 := "tests/duplicate_import.sysl"
	file2 := "tests/duplicate_import_single.sysl"
	res, err := parseAndCompareSysl(file1, file2, "")
	if assert.NoError(t, err) {
		assert.True(t, res, "%s and %s proto output diff", file1, file2)
	}
}

func TestCircularImport(t *testing.T) {
	t.Parallel()

	circularImportA := `
import circular_import_b
One:
	...
`
	circularImportB := `
import circular_import_c
Two:
	...
`
	circularImportC := `
import circular_import_a
Three:
	...
`
	r := mockReader{contents: map[string]mockContent{
		"circular_import_a.sysl": {circularImportA, retriever.ZeroHash, ""},
		`circular_import_b.sysl`: {circularImportB, retriever.ZeroHash, ""},
		`circular_import_c.sysl`: {circularImportC, retriever.ZeroHash, ""},
	}}

	c, err := NewParser().Parse("circular_import_a", r)

	require.NoError(t, err)
	require.Equal(t, 3, len(c.Apps))
	require.NotNil(t, c.Apps["One"])
	require.NotNil(t, c.Apps["Two"])
	require.NotNil(t, c.Apps["Three"])
}

func TestImportWithMaxDepth(t *testing.T) {
	t.Parallel()

	one := `
import two
One:
	...
`
	two := `
import three
Two:
	...
`
	three := `
Three:
	...
`
	r := mockReader{contents: map[string]mockContent{
		"one.sysl":   {one, retriever.ZeroHash, ""},
		`two.sysl`:   {two, retriever.ZeroHash, ""},
		`three.sysl`: {three, retriever.ZeroHash, ""},
	}}

	p := NewParser()

	// Depth not specified
	c, err := p.Parse("one", r)
	require.NoError(t, err)
	require.Equal(t, 3, len(c.Apps))
	require.NotNil(t, c.Apps["One"])
	require.NotNil(t, c.Apps["Two"])
	require.NotNil(t, c.Apps["Three"])

	// Depth of 1
	p.SetMaxImportDepth(1)
	c, err = p.Parse("one", r)
	require.NoError(t, err)
	require.Equal(t, 1, len(c.Apps))
	require.NotNil(t, c.Apps["One"])

	// Depth of 2
	p.SetMaxImportDepth(2)
	c, err = p.Parse("one", r)
	require.NoError(t, err)
	require.Equal(t, 2, len(c.Apps))
	require.NotNil(t, c.Apps["One"])
	require.NotNil(t, c.Apps["Two"])
}

func TestLintValid(t *testing.T) {
	assertLintLogs(t,
		"tests/lint_valid.sysl", "")
}

func TestCaseSensitiveRedefinition(t *testing.T) {
	assertLintLogs(t,
		"tests/case_sensitive_redefinition.sysl",
		`lint: case-sensitive redefinitions detected:\n`+
			`ApP:tests/case_sensitive_redefinition.sysl:5:4\n`+
			`App:tests/case_sensitive_redefinition.sysl:2:4\n`+
			`aPP:tests/case_sensitive_redefinition.sysl:8:4`,
	)
}

func TestCaseSensitiveRedefinitionImport(t *testing.T) {
	assertLintLogs(t,
		"tests/case_sensitive_redef_import1.sysl",
		`lint: case-sensitive redefinitions detected:\n`+
			`Sensitive:tests/case_sensitive_redef_import1.sysl:3:4\n`+
			`sEnsitive:tests/case_sensitive_redef_import2.sysl:2:4`,
	)
}

func TestReturnLint(t *testing.T) {
	assertLintLogs(t,
		"tests/lint_return.sysl",
		`lint tests/lint_return.sysl:5:12: 'return some_type' not supported, use 'return ok <: some_type' instead`,
	)
}

func TestEndpointDoesNotExistLint(t *testing.T) {
	assertLintLogs(t,
		"tests/invalid_call_endpoint.sysl",
		`lint tests/invalid_call_endpoint.sysl:4:12: Endpoint '/hello' does not exist for call 'Call <- GET /hello'`)
	assertLintLogs(t,
		"tests/invalid_call_endpoint_import.sysl",
		`lint tests/invalid_call_endpoint_import.sysl:6:12: Endpoint '/hello' does not exist for call 'Call <- GET /hello'`)
	assertLintLogs(t,
		"tests/invalid_call_self_endpoint.sysl",
		`lint tests/invalid_call_self_endpoint.sysl:4:12: Endpoint '/hello' does not exist for call 'App <- POST /hello'`,
	)
	assertLintLogs(t,
		"tests/invalid_call_simple_endpoint.sysl",
		`lint tests/invalid_call_simple_endpoint.sysl:3:8: Endpoint 'Endpoint' does not exist for call 'Call <- Endpoint'`,
	)
}

func TestMethodDoesNotExistLint(t *testing.T) {
	assertLintLogs(t,
		"tests/invalid_call_method.sysl",
		`lint tests/invalid_call_method.sysl:4:12: Method 'POST' does not exist for call 'Call <- POST /hello'`)
	assertLintLogs(t,
		"tests/invalid_call_method_import.sysl",
		`lint tests/invalid_call_method_import.sysl:6:12: Method 'POST' does not exist for call 'Call <- POST /hello'`,
	)
	assertLintLogs(t,
		"tests/invalid_call_nested.sysl",
		`lint tests/invalid_call_nested.sysl:4:12: Method 'POST' does not exist for call 'Call <- POST /hi/hello/yo'`,
	)
	assertLintLogs(t,
		"tests/invalid_call_self_method.sysl",
		`lint tests/invalid_call_self_method.sysl:4:12: Method 'POST' does not exist for call 'App <- POST /hi'`,
	)
}

func TestAppDoesNotExistLint(t *testing.T) {
	assertLintLogs(t,
		"tests/invalid_call_app.sysl",
		`lint tests/invalid_call_app.sysl:4:12: Application 'Call' does not exist for call 'Call <- GET /hello'`,
	)
	assertLintLogs(t,
		"tests/invalid_call_simple_app.sysl",
		`lint tests/invalid_call_simple_app.sysl:3:8: Application 'Call' does not exist for call 'Call <- End'`,
	)
	assertLintLogs(t, "tests/app_names_with_spaces.sysl", ``)
}

func assertLintLogs(t *testing.T, file, logMsg string) {
	var buf bytes.Buffer
	// FIXME: using logrus global logger makes it impossible to parallelize log tests
	logrus.SetOutput(&buf)
	_, err := NewParser().ParseFromFs(file, syslutil.NewChrootFs(afero.NewOsFs(), ""))
	require.NoError(t, err)
	logrus.SetOutput(os.Stderr)
	if logMsg == "" {
		assert.Equal(t, "", buf.String())
	} else {
		assert.Contains(t, buf.String(), fmt.Sprintf("level=warning msg=\"%s\"", logMsg))
	}
}

func TestInferExprTypeNonTransform(t *testing.T) {
	t.Parallel()

	type expectedData struct {
		exprType     *sysl.Type
		inferredType *sysl.Type
		messages     map[string][]msg.Msg
	}

	memFs, fs := syslutil.WriteToMemOverlayFs(mainTestDir)
	parser := NewParser()
	expressions := map[string]*sysl.Expr{}
	transform, appName, err := LoadAndGetDefaultApp("transform1.sysl", fs, parser)
	require.NoError(t, err)
	syslutil.AssertFsHasExactly(t, memFs)
	viewName := "inferExprTypeNonTransform"
	viewRetType := transform.GetApps()[appName].Views[viewName].GetRetType()

	for _, stmt := range transform.GetApps()[appName].Views[viewName].GetExpr().GetTransform().GetStmt() {
		expressions[stmt.GetAssign().GetName()] = stmt.GetAssign().GetExpr()
	}

	cases := map[string]struct {
		input    *sysl.Expr
		expected expectedData
	}{
		"String": {
			input: expressions["stringType"],
			expected: expectedData{
				exprType: syslutil.TypeString(), inferredType: syslutil.TypeString(), messages: map[string][]msg.Msg{}}},
		"Int": {
			input: expressions["intType"],
			expected: expectedData{
				exprType: syslutil.TypeInt(), inferredType: syslutil.TypeInt(), messages: map[string][]msg.Msg{}}},
		"Bool": {
			input: expressions["boolType"],
			expected: expectedData{
				exprType: syslutil.TypeBool(), inferredType: syslutil.TypeBool(), messages: map[string][]msg.Msg{}}},
		"Valid bool unary result": {
			input: expressions["unaryResultValidBool"],
			expected: expectedData{
				exprType: syslutil.TypeBool(), inferredType: syslutil.TypeBool(), messages: map[string][]msg.Msg{}}},
		"Valid int unary result": {
			input: expressions["unaryResultValidInt"],
			expected: expectedData{
				exprType: syslutil.TypeInt(), inferredType: syslutil.TypeInt(), messages: map[string][]msg.Msg{}}},
		"Invalid unary result bool": {
			input: expressions["unaryResultInvalidBool"],
			expected: expectedData{
				exprType: syslutil.TypeBool(), inferredType: syslutil.TypeBool(),
				messages: map[string][]msg.Msg{viewName: {
					{MessageID: msg.ErrInvalidUnary, MessageData: []string{viewName, "STRING"}}}}}},
		"Invalid unary result int": {
			input: expressions["unaryResultInvalidInt"],
			expected: expectedData{
				exprType: syslutil.TypeInt(), inferredType: syslutil.TypeInt(),
				messages: map[string][]msg.Msg{viewName: {
					{MessageID: msg.ErrInvalidUnary, MessageData: []string{viewName, "STRING"}}}}}},
	}

	for name, test := range cases {
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			newParser := NewParser()
			exprType, _, inferredType := newParser.inferExprType(nil, "", input, false, 0,
				viewName, viewName, viewRetType)
			assert.Equal(t, expected.exprType, exprType)
			assert.Equal(t, expected.inferredType, inferredType)
			assert.Equal(t, expected.messages, newParser.GetMessages())
		})
	}
}

func TestInferExprTypeTransform(t *testing.T) {
	t.Parallel()

	type expectedData struct {
		exprType     *sysl.Type
		inferredType *sysl.Type

		letTypeRef   *sysl.Type
		letTypeTuple *sysl.Type
		letTypeScope string
		messages     map[string][]msg.Msg
	}

	memFs, fs := syslutil.WriteToMemOverlayFs(mainTestDir)
	parser := NewParser()
	transform, appName, err := LoadAndGetDefaultApp("transform1.sysl", fs, parser)
	require.NoError(t, err)
	syslutil.AssertFsHasExactly(t, memFs)
	views := transform.GetApps()[appName].Views

	cases := map[string]struct {
		viewName string
		expected expectedData
	}{
		"Transform type assign": {
			viewName: "inferExprTypeTransformAssign",
			expected: expectedData{
				exprType: syslutil.TypeString(),
				inferredType: &sysl.Type{
					Type: &sysl.Type_Tuple_{
						Tuple: &sysl.Type_Tuple{AttrDefs: map[string]*sysl.Type{"transformTypeAssign": {
							Type: &sysl.Type_Tuple_{
								Tuple: &sysl.Type_Tuple{AttrDefs: map[string]*sysl.Type{"bar": syslutil.TypeString()}}}}}}}},
				letTypeScope: "inferExprTypeTransformAssign:transformTypeAssign:foo",
				letTypeRef:   syslutil.TypeString(),
				letTypeTuple: syslutil.TypeString(),
				messages:     map[string][]msg.Msg{}}},
		"Nested transform type assign": {
			viewName: "inferExprTypeTransformNestedAssign",
			expected: expectedData{
				exprType: syslutil.TypeString(),
				inferredType: &sysl.Type{
					Type: &sysl.Type_Tuple_{
						Tuple: &sysl.Type_Tuple{AttrDefs: map[string]*sysl.Type{"nestedTransformTypeAssignTfm": {
							Type: &sysl.Type_Tuple_{
								Tuple: &sysl.Type_Tuple{AttrDefs: map[string]*sysl.Type{"bar": {
									Type: &sysl.Type_Tuple_{
										Tuple: &sysl.Type_Tuple{AttrDefs: map[string]*sysl.Type{"assign": syslutil.TypeString()}},
									}}}}}}}}}},
				letTypeScope: "inferExprTypeTransformNestedAssign:nestedTransformTypeAssignTfm:bar:variable",
				letTypeRef:   syslutil.TypeInt(),
				letTypeTuple: syslutil.TypeInt(),
				messages:     map[string][]msg.Msg{}}},
		"Nested transform type let": {
			viewName: "inferExprTypeTransformNestedLet",
			expected: expectedData{
				exprType: syslutil.TypeString(),
				inferredType: &sysl.Type{
					Type: &sysl.Type_Tuple_{
						Tuple: &sysl.Type_Tuple{AttrDefs: map[string]*sysl.Type{"nestedTransformTypeLetTfm": {
							Type: &sysl.Type_Tuple_{
								Tuple: &sysl.Type_Tuple{AttrDefs: map[string]*sysl.Type{"foo": syslutil.TypeNone()}}}}}}}},
				letTypeScope: "inferExprTypeTransformNestedLet:nestedTransformTypeLetTfm:bar:variable",
				letTypeRef:   syslutil.TypeInt(),
				letTypeTuple: syslutil.TypeInt(),
				messages:     map[string][]msg.Msg{}}},
	}

	for name, test := range cases {
		viewName := test.viewName
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			newParser := NewParser()
			_, _, inferredType := newParser.inferExprType(nil, "", views[viewName].GetExpr(), true, 0,
				viewName, viewName, views[viewName].GetRetType())

			assert.Equal(t, expected.inferredType, inferredType)
			assert.Equal(t, newParser.GetAssigns()[viewName].Tuple, inferredType)
			assert.Equal(t, newParser.GetLets()[expected.letTypeScope].Tuple, expected.letTypeTuple)
			assert.Equal(t, newParser.GetLets()[expected.letTypeScope].RefType, expected.letTypeRef)
			assert.Equal(t, expected.messages, newParser.GetMessages())
		})
	}
}

func TestParseSysl(t *testing.T) {
	t.Parallel()

	content := `
App:
	Endpoint:
		...
`
	_, err := NewParser().ParseString(content)
	assert.Nil(t, err)
}

type mockReader struct {
	// This is added so that it implements afero.Fs which golden retriever needs.
	// Currently unused
	afero.Fs
	contents map[string]mockContent
}

type mockContent struct {
	content string
	hash    retriever.Hash
	branch  string
}

func (r mockReader) Read(_ context.Context, resource string) ([]byte, error) {
	if c, ok := r.contents[resource]; ok {
		return []byte(c.content), nil
	}
	return nil, fmt.Errorf("file not found")
}

func (r mockReader) ReadHash(_ context.Context, resource string) ([]byte, retriever.Hash, error) {
	if c, ok := r.contents[resource]; ok {
		return []byte(c.content), r.contents[resource].hash, nil
	}
	return nil, retriever.Hash{}, fmt.Errorf("file not found")
}

func (r mockReader) ReadHashBranch(_ context.Context, resource string) ([]byte, retriever.Hash, string, error) {
	if c, ok := r.contents[resource]; ok {
		return []byte(c.content), r.contents[resource].hash, r.contents[resource].branch, nil
	}
	return nil, retriever.ZeroHash, "", fmt.Errorf("file not found")
}

/* TestParseSyslRetriever tests that a file can be imported */
func TestParseSyslRetriever(t *testing.T) {
	t.Parallel()

	one := `
import two.sysl
One:
	_:
		...
`
	two := `
Two:
	...
`
	nottwo := `
NotTwo:
	...
`

	r := mockReader{contents: map[string]mockContent{
		"./one.sysl": {one, retriever.ZeroHash, ""},
		"two.sysl":   {two, retriever.ZeroHash, ""},
		"./two.sysl": {nottwo, retriever.ZeroHash, ""},
	}}

	p := NewParser()

	c, err := p.Parse("./one", r)
	require.NoError(t, err)
	require.Equal(t, 2, len(c.Apps))
	require.NotNil(t, c.Apps["One"])
	require.NotNil(t, c.Apps["Two"])
}

const three = `
Three:
	_:
		...
`

/* TestParseSyslRetrieverRemote tests a remote import */
func TestParseSyslRetrieverRemote(t *testing.T) {
	t.Parallel()

	one := `
import //github.com/org/repo/two.sysl@master
import 3.sysl
One:
	_:
		...
`
	two := `
Two:
	...
`

	sha := randomSha
	h, err := retriever.NewHash(sha)
	require.NoError(t, err)
	r := mockReader{contents: map[string]mockContent{
		"./one.sysl":                            {one, retriever.ZeroHash, ""},
		"//github.com/org/repo/two.sysl@master": {two, h, "master"},
		`3.sysl`:                                {three, retriever.ZeroHash, ""},
	}}

	p := NewParser()

	c, err := p.Parse("./one", r)
	require.NoError(t, err)
	require.Equal(t, 3, len(c.Apps))
	require.NotNil(t, c.Apps["One"])
	require.NotNil(t, c.Apps["Two"])
	require.NotNil(t, c.Apps["Three"])
}

/* TestParseSyslRetrieverRemoteImport tests a remote file that imports from its own repository */
func TestParseSyslRetrieverRemoteImport(t *testing.T) {
	t.Parallel()

	one := `
import //github.com/org/repo/two.sysl@master
One:
	_:
		...
`
	two := `
import three.sysl
Two:
	...
`

	sha := randomSha
	h, err := retriever.NewHash(sha)
	require.NoError(t, err)
	r := mockReader{contents: map[string]mockContent{
		"./one.sysl":                              {one, retriever.ZeroHash, ""},
		"//github.com/org/repo/two.sysl@master":   {two, h, "master"},
		"//github.com/org/repo/three.sysl@master": {three, h, "master"},
	}}

	p := NewParser()

	c, err := p.Parse("./one", r)
	require.NoError(t, err)
	require.Equal(t, 3, len(c.Apps))
	require.NotNil(t, c.Apps["One"])
	require.NotNil(t, c.Apps["Two"])
	require.NotNil(t, c.Apps["Three"])
}

/* TestParseSyslRetrieverRemoteImportFromRoot tests a remote file that imports from the root of its own repository */
func TestParseSyslRetrieverRemoteImportFromRoot(t *testing.T) {
	t.Parallel()

	one := `
import //github.com/org/repo/subdir/two.sysl@master
One:
	_:
		...
`
	two := `
import /three.sysl
Two:
	...
`

	sha := randomSha
	h, err := retriever.NewHash(sha)
	require.NoError(t, err)
	r := mockReader{contents: map[string]mockContent{
		"./one.sysl": {one, retriever.ZeroHash, ""},
		"//github.com/org/repo/subdir/two.sysl@master": {two, h, "master"},
		"//github.com/org/repo/three.sysl@master":      {three, h, "master"},
	}}

	p := NewParser()

	c, err := p.Parse("./one", r)
	require.NoError(t, err)
	require.Equal(t, 3, len(c.Apps))
	require.NotNil(t, c.Apps["One"])
	require.NotNil(t, c.Apps["Two"])
	require.NotNil(t, c.Apps["Three"])
}

/*
TestParseSyslRetrieverRemoteImportNotMaster tests that a remote file imported not from master branch will import

	a local file from the same branch
*/
func TestParseSyslRetrieverRemoteImportNotMaster(t *testing.T) {
	t.Parallel()

	one := `
import //github.com/org/repo/two.sysl@branch
One:
	_:
		...
`
	two := `
import three.sysl
Two:
	...
`
	three := `
Three:
	...
`

	sha := randomSha
	h, err := retriever.NewHash(sha)
	require.NoError(t, err)
	r := mockReader{contents: map[string]mockContent{
		"./one.sysl":                              {one, retriever.ZeroHash, ""},
		"//github.com/org/repo/two.sysl@branch":   {two, h, "branch"},
		"//github.com/org/repo/three.sysl@branch": {three, h, "branch"},
	}}

	p := NewParser()

	c, err := p.Parse("./one", r)
	require.NoError(t, err)
	require.Equal(t, 3, len(c.Apps))
	require.NotNil(t, c.Apps["One"])
	require.NotNil(t, c.Apps["Two"])
	require.NotNil(t, c.Apps["Three"])
}

/* TestParseSyslRetrieverRemoteFail tests an invalid remote import (error includes file that tried to import it)*/
func TestParseSyslRetrieverRemoteFail(t *testing.T) {
	t.Parallel()

	one := `
import //github.com/org/repo/two.sysl@master
`
	r := mockReader{contents: map[string]mockContent{
		"./one.sysl": {one, retriever.ZeroHash, ""},
	}}

	p := NewParser()

	_, err := p.Parse("./one", r)
	require.Error(t, err)
	require.Contains(t, err.Error(), "one.sysl")
}

func TestFixTypeRefScope(t *testing.T) {
	t.Parallel()
	module := getTestModule()

	fields := module.Apps["App"].Types["Type"].Type.(*sysl.Type_Tuple_).Tuple.AttrDefs
	// full type ref
	ref := fields["a"].GetTypeRef().Ref
	fixTypeRefScope(module, "App", ref)
	assert.Equal(t, []string{"Types"}, ref.GetAppname().Part)
	assert.Equal(t, []string{"Type2"}, ref.GetPath())

	// local type ref to a field
	ref = fields["b"].GetTypeRef().Ref
	fixTypeRefScope(module, "App", ref)
	assert.True(t, ref.GetAppname() == nil)
	assert.Equal(t, []string{"Type", "a"}, ref.GetPath())

	// local type ref
	ref = fields["c"].GetTypeRef().Ref
	fixTypeRefScope(module, "App", ref)
	assert.True(t, ref.GetAppname() == nil)
	assert.Equal(t, []string{"Type"}, ref.GetPath())

	// full type ref with namespace
	ref = fields["d"].GetTypeRef().Ref
	fixTypeRefScope(module, "App", ref)
	assert.Equal(t, []string{"A", "B", "C"}, ref.GetAppname().Part)
	assert.Equal(t, []string{"Type"}, ref.GetPath())

	// full type ref to a local type
	ref = fields["e"].GetTypeRef().Ref
	fixTypeRefScope(module, "App", ref)
	assert.Equal(t, []string{"App"}, ref.GetAppname().Part)
	assert.Equal(t, []string{"Type"}, ref.GetPath())

	// nil ref
	ref = fields["f"].GetTypeRef().Ref
	fixTypeRefScope(module, "App", ref)
	assert.True(t, ref == nil)

	// empty type path
	ref = fields["g"].GetTypeRef().Ref
	fixTypeRefScope(module, "App", ref)
	assert.True(t, len(ref.GetPath()) == 0)

	// params
	fixParamTypeRef(module, module.Apps["App"], "App")
	params := module.Apps["App"].Endpoints["Endpoint"].Param

	ref = params[0].Type.GetTypeRef().Ref
	assert.True(t, ref.GetAppname() == nil)
	assert.Equal(t, []string{"Type", "a"}, ref.GetPath())

	ref = params[1].Type.GetTypeRef().Ref
	assert.Equal(t, []string{"Types"}, ref.GetAppname().Part)
	assert.Equal(t, []string{"Type2"}, ref.GetPath())
}

func getTestModule() *sysl.Module {
	return &sysl.Module{
		Apps: map[string]*sysl.Application{
			"App": {
				Name: &sysl.AppName{Part: []string{"App"}},
				Endpoints: map[string]*sysl.Endpoint{
					"Endpoint": {
						Name: "Endpoint",
						Param: []*sysl.Param{
							{
								Name: "a",
								Type: &sysl.Type{
									Type: &sysl.Type_TypeRef{
										TypeRef: &sysl.ScopedRef{
											Ref: &sysl.Scope{
												// refer to the field Type.x
												Appname: &sysl.AppName{Part: []string{"Type"}},
												Path:    []string{"a"},
											},
										},
									},
								},
							},
							{
								Name: "b",
								Type: &sysl.Type{
									Type: &sysl.Type_TypeRef{
										TypeRef: &sysl.ScopedRef{
											Ref: &sysl.Scope{
												Appname: &sysl.AppName{Part: []string{"Types"}},
												Path:    []string{"Type2"},
											},
										},
									},
								},
							},
						},
					},
				},
				Types: map[string]*sysl.Type{
					"Type": {
						Type: &sysl.Type_Tuple_{
							Tuple: &sysl.Type_Tuple{
								AttrDefs: map[string]*sysl.Type{
									"a": {
										Type: &sysl.Type_TypeRef{
											TypeRef: &sysl.ScopedRef{
												Ref: &sysl.Scope{
													Appname: &sysl.AppName{Part: []string{"Types"}},
													Path:    []string{"Type2"},
												},
											},
										},
									},
									"b": {
										Type: &sysl.Type_TypeRef{
											TypeRef: &sysl.ScopedRef{
												Ref: &sysl.Scope{
													// refer to the field Type.x
													Appname: &sysl.AppName{Part: []string{"Type"}},
													Path:    []string{"a"},
												},
											},
										},
									},
									"c": {
										Type: &sysl.Type_TypeRef{
											TypeRef: &sysl.ScopedRef{
												Ref: &sysl.Scope{
													Path: []string{"Type"},
												},
											},
										},
									},
									"d": {
										Type: &sysl.Type_TypeRef{
											TypeRef: &sysl.ScopedRef{
												Ref: &sysl.Scope{
													Appname: &sysl.AppName{Part: []string{"A", "B", "C"}},
													Path:    []string{"Type"},
												},
											},
										},
									},
									"e": {
										Type: &sysl.Type_TypeRef{
											TypeRef: &sysl.ScopedRef{
												Ref: &sysl.Scope{
													Appname: &sysl.AppName{Part: []string{"App"}},
													Path:    []string{"Type"},
												},
											},
										},
									},
									"f": {
										Type: &sysl.Type_TypeRef{
											TypeRef: &sysl.ScopedRef{
												Ref: nil,
											},
										},
									},
									"g": {
										Type: &sysl.Type_TypeRef{
											TypeRef: &sysl.ScopedRef{
												Ref: &sysl.Scope{Path: []string{}},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"Types": {
				Name: &sysl.AppName{Part: []string{"Types"}},
				Types: map[string]*sysl.Type{
					"Type2": {
						Type: &sysl.Type_Tuple_{},
					},
				},
			},
		},
	}
}
