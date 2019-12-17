package parse

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/anz-bank/sysl/pkg/msg"
	"github.com/anz-bank/sysl/pkg/pbutil"
	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const mainTestDir = "../../tests/"

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
	if err := proto.UnmarshalText(buf.String(), module); err != nil {
		return nil, errors.Wrapf(err, "Unmarshal proto: %s", filename)
	}
	return module, nil
}

func retainOrRemove(err error, filename string, retainOnError bool) error {
	if err != nil {
		if retainOnError {
			fmt.Printf("%#v retained for checking\n", filename)
		} else {
			os.Remove(filename)
		}
	}
	return err
}

func removeSourceContextFromExpr(expr *sysl.Expr) {
	if expr == nil {
		return
	}
	expr.SourceContext = nil
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
	module, err := NewParser().Parse(filename, syslutil.NewChrootFs(afero.NewOsFs(), root))
	if err != nil {
		return nil, err
	}

	for _, app := range module.Apps {
		for _, t := range app.Views {
			for _, stmt := range t.GetExpr().GetTransform().GetStmt() {
				nullifyType(stmt.GetAssign().GetExpr())
				nullifyType(stmt.GetLet().GetExpr())
			}
		}
	}

	if stripSourceContext {
		// remove stuff that does not match legacy.
		for _, app := range module.Apps {
			app.SourceContext = nil
			for _, ep := range app.Endpoints {
				ep.SourceContext = nil
			}
			for _, t := range app.Types {
				t.SourceContext = nil
			}
			for _, t := range app.Views {
				t.SourceContext = nil
				removeSourceContextFromExpr(t.Expr)
			}
		}
	}

	return module, nil
}

func nullifyType(expr *sysl.Expr) {
	if expr == nil {
		return
	}

	switch t := expr.Expr.(type) {
	case *sysl.Expr_Literal:
		switch t.Literal.Value.(type) {
		case *sysl.Value_Null_, *sysl.Value_B:
			return
		}
		expr.Type = nil
	case *sysl.Expr_List_:
		nullifyType(t.List.Expr[0])
		expr.Type = nil
	case *sysl.Expr_Ifelse:
		nullifyType(t.Ifelse.GetIfTrue())
		nullifyType(t.Ifelse.GetIfFalse())
		expr.Type = nil
	case *sysl.Expr_Binexpr:
		nullifyType(t.Binexpr.GetLhs())
		if t.Binexpr.GetOp().String() != "WHERE" {
			nullifyType(t.Binexpr.GetRhs())
		}
		expr.Type = nil
	case *sysl.Expr_Unexpr:
		nullifyType(expr.GetUnexpr().GetArg())
		expr.Type = nil
	case *sysl.Expr_Transform_:
		for _, stmt := range t.Transform.GetStmt() {
			nullifyType(stmt.GetAssign().GetExpr())
			nullifyType(stmt.GetLet().GetExpr())
		}
	}
}

func parseAndCompare(
	filename, root, golden string,
	goldenProto proto.Message,
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

	if err = pbutil.TextPB(goldenProto, golden, afero.NewMemMapFs()); err != nil {
		return false, err
	}

	pattern := "sysl-test-*.textpb"
	generated, err := ioutil.TempFile("", pattern)
	if err != nil {
		return false, errors.Wrapf(err, "Create tempfile: %s", pattern)
	}
	generatedClosed := false
	defer func() {
		if !generatedClosed {
			generated.Close()
		}
	}()

	if err = retainOrRemove(pbutil.FTextPB(generated, module), generated.Name(), retainOnError); err != nil {
		return false, errors.Wrapf(err, "Generate %#v", generated)
	}
	if err := generated.Close(); err != nil {
		return false, errors.Wrapf(err, "Close %#v", generated)
	}
	generatedClosed = true

	cmd := exec.Command("diff", "-yw", golden, generated.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if retainOnError {
			fmt.Printf("%#v retained for checking\n", generated.Name())
		} else {
			os.Remove(generated.Name())
		}
		return false, errors.Wrapf(err, "diff -yw %#v %#v", golden, generated.Name())
	}

	os.Remove(generated.Name())
	return false, nil
}

func parseAndCompareWithGolden(filename, root string, stripSourceContext bool) (bool, error) {
	goldenFilename := filename
	if !strings.HasSuffix(goldenFilename, syslExt) {
		goldenFilename += syslExt
	}
	goldenFilename += ".golden.textpb"
	golden := path.Join(root, goldenFilename)

	goldenModule, err := readSyslModule(golden)
	if err != nil {
		return false, err
	}
	return parseAndCompare(filename, root, golden, goldenModule, true, stripSourceContext)
}

func testParseAgainstGolden(t *testing.T, filename, root string) {
	equal, err := parseAndCompareWithGolden(filename, root, true)
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
	tmpdir := path.Join(tmproot, dirname)
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

func TestRelational(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/school.sysl", "")
}

func TestImports(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/library.sysl", "")
}

func TestForeignImports(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "tests/foreign_import_swagger.sysl", "")
}

func TestRootArg(t *testing.T) {
	t.Parallel()

	testParseAgainstGolden(t, "school.sysl", "tests")
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

func TestStrings(t *testing.T) {
	t.Parallel()

	testParseAgainstGoldenWithSourceContext(t, "tests/strings_expr.sysl")
}

func TestTypeAlias(t *testing.T) {
	t.Parallel()

	testParseAgainstGoldenWithSourceContext(t, "tests/alias.sysl")
}

func TestUndefinedRootAbsoluteImport(t *testing.T) {
	t.Parallel()

	parser := NewParser()
	parser.RestrictToLocalImport()
	_, err := parser.Parse("absolute_import.sysl", syslutil.NewChrootFs(afero.NewOsFs(), "tests"))
	require.EqualError(t, err, "error importing: importing outside current directory is only allowed when root is defined")
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
