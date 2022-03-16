package parse

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/anz-bank/golden-retriever/reader"
	"github.com/anz-bank/golden-retriever/reader/remotefs"
	"github.com/anz-bank/pkg/mod"
	"github.com/anz-bank/sysl/pkg/env"
	parser "github.com/anz-bank/sysl/pkg/grammar"
	"github.com/anz-bank/sysl/pkg/importer"
	"github.com/anz-bank/sysl/pkg/msg"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"
)

// TypeData contains referenced type and actual tuple of referenced type
type TypeData struct {
	RefType *sysl.Type
	Tuple   *sysl.Type
}

// Parser holds assign types and let(variable) definitions in hierarchical order, messages generated while parsing
// and whether a root was found or not
type Parser struct {
	AssignTypes         map[string]TypeData
	LetTypes            map[string]TypeData
	Messages            map[string][]msg.Msg
	allowAbsoluteImport bool
}

// NewParser intializes and returns a new Parser instance
func NewParser() *Parser {
	return &Parser{
		AssignTypes:         map[string]TypeData{},
		LetTypes:            map[string]TypeData{},
		Messages:            map[string][]msg.Msg{},
		allowAbsoluteImport: true,
	}
}

func parseString(filename string, input antlr.CharStream) (parser.ISysl_fileContext, error) {
	errorListener := SyslParserErrorListener{}
	lexer := parser.NewThreadSafeSyslLexer(input)
	defer parser.DeleteLexerState(lexer)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewThreadSafeSyslParser(stream)
	p.GetInterpreter().SetPredictionMode(antlr.PredictionModeSLL)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	p.AddErrorListener(&errorListener)

	p.BuildParseTrees = true
	tree := p.Sysl_file()
	if errorListener.hasErrors {
		return nil, syslutil.Exitf(ParseError, fmt.Sprintf("%s has syntax errors\n", filename))
	}
	return tree, nil
}

// Currently only supports JSON format Proto.
// TODO: Add support for binary and textpb
func importSyslProto(fsinput antlr.CharStream) (*sysl.Module, error) {
	file := fsinput.GetText(0, fsinput.Size())
	syslFile := &sysl.Module{}
	err := protojson.Unmarshal([]byte(file), syslFile)
	if err != nil {
		return nil, err
	}

	return syslFile, nil
}

func importForeign(def importDef, input antlr.CharStream) (antlr.CharStream, error) {
	logger := logrus.StandardLogger()
	fileName, _ := mod.ExtractVersion(def.filename)
	file := input.GetText(0, input.Size())
	fileType, err := detectFileType(fileName, []byte(file))
	if err != nil {
		return nil, err
	}

	switch fileType.Name {
	case importer.SYSL.Name:
		return input, nil
	case importer.OpenAPI3.Name, importer.OpenAPI2.Name:
		imp, err := importer.Factory(fileName, false, "", []byte(file), logger)
		imp.WithAppName(def.appname).WithPackage(def.pkg)
		if err != nil {
			return nil, syslutil.Exitf(ParseError, fmt.Sprintf("%s has unknown format: %s", fileName, err))
		}
		// FIXME: because filepath information is not provided, external references are ignored in OpenAPI3.
		output, err := imp.Load(file)
		if err != nil {
			return nil, syslutil.Exitf(ParseError, fmt.Sprintf("%s has unknown format: %s", fileName, err))
		}
		return antlr.NewInputStream(output), nil
	default:
		return nil, syslutil.Exitf(ParseError, fmt.Sprintf("%s has unknown format", fileName))
	}
}

func detectFileType(fileName string, file []byte) (importer.Format, error) {
	var ParserFormats = []importer.Format{
		importer.OpenAPI3,
		importer.OpenAPI2,
		importer.SYSL,
	}
	return importer.GuessFileType(fileName, false, file, ParserFormats)
}

func (p *Parser) RestrictToLocalImport() {
	// if root is not defined, only relative imports are allowed
	p.allowAbsoluteImport = false
}

// ParseString parses a sysl definition in string form.
func (p *Parser) ParseString(content string) (*sysl.Module, error) {
	fs := afero.NewMemMapFs()
	file, err := fs.Create("temp.sysl")
	if err != nil {
		return nil, err
	}
	if _, err := file.Write([]byte(content)); err != nil {
		return nil, err
	}
	return p.ParseFromFs("temp.sysl", fs)
}

// ParseFromFs parses a sysl definition from an afero filesystem
func (p *Parser) ParseFromFs(filename string, fs afero.Fs) (*sysl.Module, error) {
	reader, err := NewReader(fs)
	if err != nil {
		return nil, err
	}
	return p.Parse(filename, reader)
}

// ParseFromFsWithVendor parses a sysl definition from an afero filesystem, and vendor remote files in root dir
func (p *Parser) ParseFromFsWithVendor(filename string, fs afero.Fs) (*sysl.Module, error) {
	reader, err := NewReader(fs)
	if err != nil {
		return nil, err
	}

	reader.(*remotefs.RemoteFs).Vendor(SyslRootDir(fs))
	return p.Parse(filename, reader)
}

type srcInput struct {
	src   importDef
	input string
}

type srcInputs struct {
	inputs []srcInput
	ch     chan struct{}
}

type retrievedList struct {
	l     map[string]*srcInputs
	mutex sync.Mutex
}

// Parse parses a sysl definition from an retriever interface
func (p *Parser) Parse(resource string, reader reader.Reader) (*sysl.Module, error) {
	listener := NewTreeShapeListener()
	listener.lint()

	if filepath.Ext(resource) == "" {
		resource += syslExt
	}

	retrieved := retrievedList{make(map[string]*srcInputs), sync.Mutex{}}

	if err := collectSpecs(
		context.Background(),
		importDef{filename: resource},
		reader,
		&retrieved,
	); err != nil {
		return nil, err
	}

	return p.parseSpecs(retrieved.l[resource], listener)
}

func (p *Parser) parseSpecs(specs *srcInputs, listener *TreeShapeListener) (*sysl.Module, error) { //nolint:funlen
	for _, v := range specs.inputs {
		src := v.src
		logrus.Debug("Parsing: ", src.filename)

		fsinput := &fsFileStream{antlr.NewInputStream(v.input), src.filename}

		version := ""
		if syslutil.IsRemoteImport(src.filename) {
			v := strings.Split(src.filename, "@")
			version = v[len(v)-1]
		}
		srcCtxFile := cleanImportFilename(src.filename)
		listener.base = importDir(src.filename)
		// FIXME: listener.sc.version and listener.version maybe duplicated
		listener.sc = sourceCtxHelper{srcCtxFile, version}
		listener.version = version

		// Import Sysl Proto
		if strings.HasSuffix(src.filename, ".sysl.pb.json") {
			syslProtoImport, err := importSyslProto(fsinput)
			if err != nil {
				return nil, fmt.Errorf("error parsing %s: %w", src.filename, err)
			}
			if syslProtoImport != nil {
				// Merge structs recursively
				//nolint:govet
				if err := mergo.Merge(listener.module, *syslProtoImport); err != nil {
					return nil, err
				}
			}
			break
		}

		str, err := importForeign(src, fsinput)
		if err != nil {
			return nil, err
		}

		tree, err := parseString(src.filename, str)
		if err != nil {
			return nil, err
		}

		walker := antlr.NewParseTreeWalker()
		walker.Walk(listener, tree)
	}

	listener.lintAppDefs()
	listener.lintEndpoint()
	p.postProcess(listener.module)
	return listener.module, nil
}

// collectSpecs retrieves the contents for a sourceFile. It then parses the source to find all imports and recursively
// (and in parallel) retrieves those as well. A slice of all the files and their contents (in the order they
// appear in the source file) are placed into the retrievedList.
func collectSpecs(ctx context.Context, source importDef, reader reader.Reader, retrieved *retrievedList) error {
	retrieved.mutex.Lock()
	if inputs, has := retrieved.l[source.filename]; has {
		logrus.Warnf("Duplicate import: '%s'\n", cleanImportFilename(source.filename))
		retrieved.mutex.Unlock()
		// Wait for result to be ready before returning
		<-inputs.ch

		return nil
	}
	inputs := srcInputs{nil, make(chan struct{})}
	defer close(inputs.ch)
	retrieved.l[source.filename] = &inputs
	retrieved.mutex.Unlock()

	content, hash, err := reader.ReadHash(ctx, source.filename)
	if err != nil {
		return syslutil.Exitf(ImportError, fmt.Sprintf(
			"error reading %#v: \n%v\n", cleanImportFilename(source.filename), err,
		))
	}

	version := hash.String()

	importsInput := extractImports(source.filename, content)

	sources := []srcInput{{source, string(content)}}

	if importsInput.Len() == 0 {
		retrieved.mutex.Lock()
		defer retrieved.mutex.Unlock()
		inputs.inputs = sources

		return nil
	}

	children, err := parseImports(source, sourceCtxHelper{source.filename, version}, importsInput.String())
	if err != nil {
		return err
	}

	g := new(errgroup.Group)
	for _, c := range children {
		c := c
		g.Go(func() error {
			return collectSpecs(ctx, c, reader, retrieved)
		})
	}

	err = g.Wait()
	if err != nil {
		return syslutil.Exitf(ImportError, fmt.Sprintf(
			"error reading %#v: \n%v", cleanImportFilename(source.filename), err,
		))
	}

	retrieved.mutex.Lock()
	defer retrieved.mutex.Unlock()
	alreadyAdded := map[string]struct{}{
		source.filename: {},
	}
	// Go through and add all the sources from each of the children while checking that each source is only added once
	for _, c := range children {
		for _, src := range retrieved.l[c.filename].inputs {
			if _, has := alreadyAdded[src.src.filename]; !has {
				sources = append(sources, src)
				alreadyAdded[src.src.filename] = struct{}{}
			}
		}
	}
	inputs.inputs = sources

	return nil
}

var importStmtPrefix = []byte("import ")

func extractImports(filename string, content []byte) (importsInput bytes.Buffer) {
	// non-sysl specs remote reference file fetching is not yet supported.
	if !strings.Contains(filename, syslExt) {
		return
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if bytes.HasPrefix(scanner.Bytes(), importStmtPrefix) {
			importsInput.Write(scanner.Bytes())
			importsInput.WriteByte('\n')
		}
	}

	return
}

// parseImports parses string with only import statements.
func parseImports(parent importDef, src sourceCtxHelper, input string) ([]importDef, error) {
	listener := NewTreeShapeListener()
	listener.lint()

	fsinput := &fsFileStream{antlr.NewInputStream(input), parent.filename}

	listener.sc = src
	listener.base = importDir(parent.filename)
	listener.version = src.version

	tree, err := parseString(parent.filename, fsinput)
	if err != nil {
		return nil, err
	}

	walker := antlr.NewParseTreeWalker()
	walker.Walk(listener, tree)

	return listener.imports, nil
}

// apply attributes from src to dst statement and all of its
// child statements as well (e.g. For / Loop statements).
func applyAttributes(src *sysl.Statement, dst *sysl.Statement) bool {
	var stmts []*sysl.Statement
	applied := false
	switch s := dst.GetStmt().(type) {
	case *sysl.Statement_Cond:
		stmts = s.Cond.Stmt
	case *sysl.Statement_Alt:
		for _, c := range s.Alt.Choice {
			for _, ss := range c.Stmt {
				applied = applyAttributes(src, ss) || applied
			}
		}
		return applied
	case *sysl.Statement_Group:
		stmts = s.Group.Stmt
	case *sysl.Statement_Loop:
		stmts = s.Loop.Stmt
	case *sysl.Statement_LoopN:
		stmts = s.LoopN.Stmt
	case *sysl.Statement_Foreach:
		stmts = s.Foreach.Stmt
	case *sysl.Statement_Call:
		if syslutil.IsSameCall(src.GetCall(), s.Call) {
			if dst.Attrs == nil {
				dst.Attrs = map[string]*sysl.Attribute{}
			}
			mergeAttrs(src.Attrs, dst.Attrs)
			applied = true
		}
		return applied
	case *sysl.Statement_Action:
		return applied
	case *sysl.Statement_Ret:
		return applied
	default:
		panic("collector: unhandled type")
	}

	for _, stmt := range stmts {
		applied = applyAttributes(src, stmt) || applied
	}
	return applied
}

func checkCalls(mod *sysl.Module, appname string, epname string, dst *sysl.Statement) bool {
	var stmts []*sysl.Statement
	switch s := dst.GetStmt().(type) {
	case *sysl.Statement_Cond:
		stmts = s.Cond.Stmt
	case *sysl.Statement_Alt:
		for _, c := range s.Alt.Choice {
			for _, ss := range c.Stmt {
				if !checkCalls(mod, appname, epname, ss) {
					return false
				}
			}
		}
		return true
	case *sysl.Statement_Group:
		stmts = s.Group.Stmt
	case *sysl.Statement_Loop:
		stmts = s.Loop.Stmt
	case *sysl.Statement_LoopN:
		stmts = s.LoopN.Stmt
	case *sysl.Statement_Foreach:
		stmts = s.Foreach.Stmt
	case *sysl.Statement_Call:
		app := syslutil.GetApp(s.Call.Target, mod)
		if app == nil {
			return false
		}
		_, valid := app.Endpoints[s.Call.Endpoint]
		return valid
	case *sysl.Statement_Action:
		return true
	case *sysl.Statement_Ret:
		return true
	default:
		panic("collector: unhandled type")
	}

	for _, stmt := range stmts {
		if !checkCalls(mod, appname, epname, stmt) {
			return false
		}
	}
	return true
}

func collectorPubSubCalls(appName string, app *sysl.Application) {
	endpoint := app.Endpoints[`.. * <- *`]
	if endpoint == nil {
		return
	}

	for _, collectorStmt := range endpoint.Stmt {
		switch x := collectorStmt.Stmt.(type) {
		case *sysl.Statement_Action:
			modifyEP := app.Endpoints[x.Action.Action]
			if modifyEP == nil {
				logrus.Errorf("App (%s) calls non-existent endpoint (%s)\n",
					appName, x.Action.Action)
				continue
			}
			if modifyEP.Attrs == nil {
				modifyEP.Attrs = map[string]*sysl.Attribute{}
			}
			mergeAttrs(collectorStmt.Attrs, modifyEP.Attrs)
		case *sysl.Statement_Call:
			applied := false

			for callEPName, callEndpoint := range app.Endpoints {
				if callEPName == `.. * <- *` {
					continue
				}
				for _, callStmt := range callEndpoint.Stmt {
					applied = applyAttributes(collectorStmt, callStmt) || applied
				}
			}
			if !applied {
				logrus.Errorf("Unused template (%s <- %s) in app %s\n",
					x.Call.Target.Part, x.Call.Endpoint, appName)
			}
		default:
			panic("unhandled type:")
		}
	}
}

func checkEndpointCalls(mod *sysl.Module) bool {
	valid := false
	for appName, app := range mod.Apps {
		for epname, ep := range app.Endpoints {
			for _, stmt := range ep.Stmt {
				valid = checkCalls(mod, appName, epname, stmt)
				if !valid {
					return valid
				}
			}
		}
	}
	return valid
}

func inferAnonymousType(mod *sysl.Module, appName string, anonCount int, t *sysl.Expr_Transform_, expr *sysl.Expr) int {
	newType := &sysl.Type{
		Type: &sysl.Type_Tuple_{
			Tuple: &sysl.Type_Tuple{
				AttrDefs: map[string]*sysl.Type{},
			},
		},
	}
	typeName := fmt.Sprintf("AnonType_%d__", anonCount)
	anonCount++
	if mod.Apps[appName].Types == nil {
		mod.Apps[appName].Types = map[string]*sysl.Type{}
	}
	mod.Apps[appName].Types[typeName] = newType
	attr := newType.GetTuple().AttrDefs
	for _, stmt := range t.Transform.Stmt {
		if stmt.GetAssign() != nil {
			assign := stmt.GetAssign()
			aexpr := assign.Expr
			syslutil.Assert(aexpr.GetTransform() != nil, "expression should be of type transform")
			ftype := aexpr.Type
			setof := ftype.GetSet() != nil
			if setof {
				ftype = ftype.GetSet()
			}
			syslutil.Assert(ftype.GetTypeRef() != nil, "transform type should be type_ref")
			t1 := &sysl.Type{
				Type: &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{
						Context: &sysl.Scope{
							Appname: mod.Apps[appName].Name,
							Path:    []string{typeName},
						},
						Ref: ftype.GetTypeRef().Ref,
					},
				},
			}
			if setof {
				t1 = &sysl.Type{
					Type: &sysl.Type_Set{
						Set: t1,
					},
				}
			}
			attr[assign.Name] = t1
		}
	}
	expr.Type = &sysl.Type{
		Type: &sysl.Type_Set{
			Set: &sysl.Type{
				Type: &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{
						Context: &sysl.Scope{
							Appname: mod.Apps[appName].Name,
						},
						Ref: &sysl.Scope{
							Appname: mod.Apps[appName].Name,
							Path:    []string{typeName},
						},
					},
				},
			},
		},
	}

	return anonCount
}

// for nested transform's Type
func (p *Parser) inferExprType(mod *sysl.Module,
	appName string,
	expr *sysl.Expr, top bool,
	anonCount int,
	viewName string,
	scope string,
	refType *sysl.Type) (exprType *sysl.Type, anonymousCount int, inferredType *sysl.Type) {
	switch t := expr.Expr.(type) {
	case *sysl.Expr_Transform_:
		newTuple := &sysl.Type{
			Type: &sysl.Type_Tuple_{
				Tuple: &sysl.Type_Tuple{
					AttrDefs: map[string]*sysl.Type{},
				},
			},
		}

		attrDefs := newTuple.GetTuple().AttrDefs

		for _, stmt := range t.Transform.Stmt {
			switch s := stmt.Stmt.(type) {
			case *sysl.Expr_Transform_Stmt_Assign_:
				varName := s.Assign.GetName()

				if _, exists := attrDefs[varName]; exists {
					p.Messages[viewName] = append(p.Messages[viewName], *msg.NewMsg(msg.ErrRedefined, []string{scope,
						varName}))
				} else {
					var inferredType *sysl.Type
					_, anonCount, inferredType = p.inferExprType(mod, appName, s.Assign.Expr, false, anonCount,
						viewName, scope+":"+varName, refType)
					attrDefs[varName] = inferredType
					p.AssignTypes[viewName] = TypeData{RefType: refType, Tuple: newTuple}
				}
			case *sysl.Expr_Transform_Stmt_Let:
				varName := s.Let.GetName()
				if _, exists := p.LetTypes[scope+":"+varName]; exists {
					p.Messages[viewName] = append(p.Messages[viewName], *msg.NewMsg(msg.ErrRedefined, []string{scope, varName}))
				} else {
					var inferredType *sysl.Type
					_, anonCount, inferredType = p.inferExprType(mod, appName, s.Let.Expr, false, anonCount,
						viewName, scope+":"+varName, s.Let.GetExpr().GetType())
					p.LetTypes[scope+":"+varName] = TypeData{RefType: s.Let.GetExpr().GetType(), Tuple: inferredType}
				}
			}
		}

		if !top && expr.Type == nil {
			anonCount = inferAnonymousType(mod, appName, anonCount, t, expr)
		}

		return expr.Type, anonCount, newTuple
	case *sysl.Expr_Relexpr:
		if t.Relexpr.Op == sysl.Expr_RelExpr_RANK {
			if !top && expr.Type == nil {
				type1, c, _ := p.inferExprType(mod, appName, t.Relexpr.Target, true, anonCount, viewName, scope, refType)
				anonCount = c
				logrus.Printf(type1.String())
			}
		}
	case *sysl.Expr_Literal:
		expr.Type = valueTypeToSysl(t.Literal)
	case *sysl.Expr_List_:
		expr.Type, _, _ = p.inferExprType(mod, appName, t.List.Expr[0], true, anonCount, viewName, scope, refType)
	case *sysl.Expr_Ifelse:
		exprTypeIfTrue, _, _ := p.inferExprType(mod, appName, t.Ifelse.GetIfTrue(), true, anonCount,
			viewName, scope, refType)
		expr.Type = exprTypeIfTrue
		if t.Ifelse.GetIfFalse() != nil {
			exprTypeIfFalse, _, _ := p.inferExprType(mod, appName, t.Ifelse.GetIfFalse(), true, anonCount,
				viewName, scope, refType)
			// TODO if exprTypeIfTrue != exprTypeIfFalse, raise an error. Then remove following 3 lines
			if exprTypeIfFalse != nil {
				expr.Type = exprTypeIfFalse
			}
		}
	case *sysl.Expr_Binexpr:
		exprTypeLHS, _, _ := p.inferExprType(mod, appName, t.Binexpr.GetLhs(), true, anonCount, viewName, scope, refType)
		expr.Type = exprTypeLHS
		exprTypeRHS, _, _ := p.inferExprType(mod, appName, t.Binexpr.GetRhs(), true, anonCount, viewName, scope, refType)
		// TODO if exprTypeRHS != exprTypeLHS, raise an error. Then remove following 3 lines
		if exprTypeRHS != nil {
			expr.Type = exprTypeRHS
		}

	case *sysl.Expr_Unexpr:
		varType, _, _ := p.inferExprType(mod, appName, expr.GetUnexpr().GetArg(), true, anonCount, viewName, scope, refType)
		switch t.Unexpr.GetOp() {
		case sysl.Expr_UnExpr_NOT, sysl.Expr_UnExpr_INV:
			if !syslutil.HasSameType(varType, syslutil.TypeBool()) {
				_, typeDetail := syslutil.GetTypeDetail(varType)
				p.Messages[viewName] = append(p.Messages[viewName],
					*msg.NewMsg(msg.ErrInvalidUnary, []string{viewName, typeDetail}))
			}
			expr.Type = syslutil.TypeBool()
		case sysl.Expr_UnExpr_NEG, sysl.Expr_UnExpr_POS:
			if !syslutil.HasSameType(varType, syslutil.TypeInt()) {
				_, typeDetail := syslutil.GetTypeDetail(varType)
				p.Messages[viewName] = append(p.Messages[viewName],
					*msg.NewMsg(msg.ErrInvalidUnary, []string{viewName, typeDetail}))
			}
			expr.Type = syslutil.TypeInt()
		default:
			expr.Type = syslutil.TypeNone()
		}

	default:
		// TODO Handle expression
		logrus.Debug("[Parse.infer_expr_type] Unhandled type", t)
		return syslutil.TypeNone(), anonCount, syslutil.TypeNone()
	}

	return expr.Type, anonCount, expr.Type
}

func (p *Parser) inferTypes(mod *sysl.Module, appName string) {
	for viewName, view := range mod.Apps[appName].Views {
		if syslutil.HasPattern(view.Attrs, "abstract") {
			continue
		}
		if view.Expr.GetTransform() == nil {
			logrus.Warnf("view %s expression should be of type transform", viewName)
			continue
		}
		p.inferExprType(mod, appName, view.Expr, true, 0, viewName, viewName, view.GetRetType())
	}
}

// fixTypeRefScope fixes the local and full type references so that it points to the correct type.
func fixTypeRefScope(mod *sysl.Module, currApp string, ref *sysl.Scope) {
	if ref == nil {
		return
	}

	appPath := ref.GetAppname().GetPart()
	if len(appPath) > 1 {
		// multiple app parts means full type specification
		return
	}

	// type ref must have at least app and type
	typePath := ref.GetPath()
	if len(typePath) == 0 {
		// impossible
		return
	}

	if len(appPath) == 0 && len(typePath) == 1 {
		// local type ref
		return
	}

	appName, typeName := appPath[0], typePath[0]
	if currApp == appName {
		// same app
		return
	}

	if app, exists := mod.Apps[appName]; exists {
		if _, exists := app.Types[typeName]; exists {
			// full type ref
			return
		}
	}

	if app, exists := mod.Apps[currApp]; exists {
		if _, exists := app.Types[appName]; exists {
			// local type ref using deep ref e.g A.B.C.D
			ref.Appname = nil
			ref.Path = append([]string{appName}, typePath...)
			return
		}
	}
	// TODO: type not found, return an error?
}

func fixParamTypeRef(mod *sysl.Module, app *sysl.Application, appName string) {
	for _, ep := range app.GetEndpoints() {
		for _, param := range ep.GetParam() {
			t := param.GetType()
			if ref := t.GetTypeRef(); ref != nil {
				fixTypeRefScope(mod, appName, ref.GetRef())
			}
		}
	}
}

func (p *Parser) postProcess(mod *sysl.Module) { // nolint:funlen
	appNames := make([]string, 0, len(mod.Apps))
	for a := range mod.Apps {
		appNames = append(appNames, a)
	}
	sort.Strings(appNames)

	for _, appName := range appNames {
		app := mod.Apps[appName]
		fixParamTypeRef(mod, app, appName)

		if app.Mixin2 != nil {
			for _, src := range app.Mixin2 {
				srcApp := syslutil.GetApp(src.Name, mod)
				if srcApp == nil {
					logrus.Warnf("mixin App (%s) not found", syslutil.GetAppName(src.Name))
					continue
				}
				if !syslutil.HasPattern(srcApp.Attrs, "abstract") {
					logrus.Warnf("mixin App (%s) should be ~abstract", syslutil.GetAppName(src.Name))
					continue
				}
				if srcApp.Types != nil && app.Types == nil {
					app.Types = map[string]*sysl.Type{}
				}
				if srcApp.Views != nil && app.Views == nil {
					app.Views = map[string]*sysl.View{}
				}
				for k, v := range srcApp.Types {
					if _, has := app.Types[k]; !has {
						app.Types[k] = v
					} else {
						logrus.Warnf("Type %s defined in %s and in %s",
							k, appName, syslutil.GetAppName(src.Name))
					}
				}
				for k, v := range srcApp.Views {
					if _, has := app.Views[k]; !has {
						app.Views[k] = v
					} else {
						logrus.Warnf("View %s defined in %s and in %s",
							k, appName, syslutil.GetAppName(src.Name))
					}
				}
			}
		}

		for typeName, types := range app.Types {
			var attrs map[string]*sysl.Type

			switch x := types.Type.(type) {
			case *sysl.Type_Tuple_:
				if len(x.Tuple.AttrDefs) == 0 {
					types.Type = nil
					continue
				}
				attrs = x.Tuple.GetAttrDefs()
			case *sysl.Type_Relation_:
				if len(x.Relation.AttrDefs) == 0 {
					types.Type = nil
					continue
				}
				attrs = x.Relation.GetAttrDefs()
			}

			for fieldname, t := range attrs {
				if x := t.GetTypeRef(); x != nil {
					refApp := app
					refName := x.GetRef().GetPath()[0]
					if refName == "string_8" {
						continue
					}
					fixTypeRefScope(mod, appName, x.GetRef())
					if refType, has := refApp.Types[refName]; has {
						var refAttrs map[string]*sysl.Type

						switch refType.Type.(type) {
						case *sysl.Type_Tuple_:
							refType, _ := refApp.Types[refName].Type.(*sysl.Type_Tuple_)
							refAttrs = refType.Tuple.GetAttrDefs()
						case *sysl.Type_Relation_:
							refType, _ := refApp.Types[refName].Type.(*sysl.Type_Relation_)
							refAttrs = refType.Relation.GetAttrDefs()
						}
						var field string
						var has bool
						if len(x.GetRef().GetPath()) > 1 {
							last := len(x.GetRef().GetPath()) - 1
							field = x.GetRef().GetPath()[last]
							_, has = refAttrs[field]
						} else if len(x.GetRef().GetPath()) == 1 {
							field = x.GetRef().GetPath()[0]
							_, has = refApp.Types[field]
						}
						if !has {
							logrus.Infof("Field %#v (type %#v) refers to Field %#v (type %#v) in app %#v",
								fieldname, typeName, field, refName, appName)
						}
					} else {
						logrus.Infof("Field %#v (type %#v) refers to type %s in app %#v",
							fieldname, typeName, refName, appName)
					}
				}
			}
		}
		p.inferTypes(mod, appName)
		collectorPubSubCalls(appName, app)
		renestTypes(app)
	}
	checkEndpointCalls(mod)
}

func renestTypes(app *sysl.Application) {
	mode := env.SYSL_DEV_RENEST_FLATTENED_TYPES.Value()
	if mode == "off" {
		return
	}

	typeNames := make([]string, 0, len(app.Types))
	for typeName := range app.Types {
		typeNames = append(typeNames, typeName)
	}
	sort.Strings(typeNames)

	for _, typeName := range typeNames {
		path := strings.Split(typeName, ".")
		if len(path) > 1 && injectType(app.Types[typeName], app.Types, path) {
			if mode == "move" {
				delete(app.Types, typeName)
			}
		}
	}
}

func injectType(leaf *sysl.Type, attrs map[string]*sysl.Type, path []string) bool {
	if attrs == nil {
		return false
	}
	if len(path) > 1 {
		return injectType(leaf, subTuple(attrs, path[0]), path[1:])
	}
	attrs[path[0]] = leaf
	return true
}

func subTuple(attrs map[string]*sysl.Type, elem string) map[string]*sysl.Type {
	if attr, has := attrs[elem]; has {
		if tuple, is := attr.Type.(*sysl.Type_Tuple_); is {
			return tuple.Tuple.AttrDefs
		}
	}
	return nil
}

func valueTypeToSysl(value *sysl.Value) *sysl.Type {
	switch value.Value.(type) {
	case *sysl.Value_B:
		return syslutil.TypeBool()
	case *sysl.Value_I:
		return syslutil.TypeInt()
	case *sysl.Value_D:
		return syslutil.TypeFloat()
	case *sysl.Value_S:
		return syslutil.TypeString()
	case *sysl.Value_Decimal:
		return syslutil.TypeDecimal()
	case *sysl.Value_Null_:
		return syslutil.TypeEmpty()
	default:
		panic(errors.Errorf("valueTypeToSysl: unhandled type: %v", value))
	}
}

// return the one and only app defined in the module
func getDefaultAppName(mod *sysl.Module) string {
	for app := range mod.Apps {
		return app
	}
	return ""
}

func LoadAndGetDefaultApp(model string, fs afero.Fs, p *Parser) (*sysl.Module, string, error) {
	// Model we want to generate code for
	mod, err := p.ParseFromFs(model, fs)
	if err != nil {
		return nil, "", err
	}
	modelAppName := getDefaultAppName(mod)
	return mod, modelAppName, nil
}

func (p *Parser) GetAssigns() map[string]TypeData {
	return p.AssignTypes
}

func (p *Parser) GetLets() map[string]TypeData {
	return p.LetTypes
}

func (p *Parser) GetMessages() map[string][]msg.Msg {
	return p.Messages
}
