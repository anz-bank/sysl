package parse

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/anz-bank/golden-retriever/reader"
	"github.com/anz-bank/golden-retriever/reader/remotefs"
	"github.com/anz-bank/pkg/mod"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"golang.org/x/sync/errgroup"

	"github.com/anz-bank/sysl/pkg/env"
	parser "github.com/anz-bank/sysl/pkg/grammar"
	"github.com/anz-bank/sysl/pkg/importer"
	"github.com/anz-bank/sysl/pkg/msg"
	"github.com/anz-bank/sysl/pkg/pbutil"
	"github.com/anz-bank/sysl/pkg/printer"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
)

type Settings struct {
	MaxImportDepth          int
	OperationSummary        bool
	NoDifferentVersionCheck bool
	NoParsing               bool
}

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
	Settings
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

func parseString(filename string, input antlr.CharStream) (tree parser.ISysl_fileContext, err error) {
	defer func() {
		// recover from panic if one occurred. Set err to nil otherwise.
		if recover() != nil {
			err = syslutil.Exitf(ParseError, fmt.Sprintf("%s has syntax errors\n", filename))
		}
	}()

	errorListener := SyslParserErrorListener{}
	lexer := parser.NewThreadSafeSyslLexer(input)
	defer parser.DeleteLexerState(lexer)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewThreadSafeSyslParser(stream)
	p.GetInterpreter().SetPredictionMode(antlr.PredictionModeSLL)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	p.AddErrorListener(&errorListener)

	p.BuildParseTrees = true
	tree = p.Sysl_file()
	if errorListener.hasErrors {
		return nil, syslutil.Exitf(ParseError, fmt.Sprintf("%s has syntax errors\n", filename))
	}
	return tree, nil
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
	case importer.SyslPB.Name:
		m, err := pbutil.FromPBByteContents(fileName, []byte(file))
		if err != nil {
			return nil, syslutil.Exitf(ParseError, fmt.Sprintf("%s has unknown format: %s", fileName, err))
		}
		var buf bytes.Buffer
		printer.Module(&buf, m)
		output := buf.String()
		return antlr.NewInputStream(output), nil

	case importer.OpenAPI3.Name, importer.OpenAPI2.Name, importer.Protobuf.Name:
		imp, err := importer.Factory(fileName, false, "", []byte(file), logger)
		if err != nil {
			return nil, syslutil.Exitf(ParseError, fmt.Sprintf("%s has unknown format: %s", fileName, err))
		}
		imp, err = imp.Configure(&importer.ImporterArg{AppName: def.appname, PackageName: def.pkg, Imports: ""})
		if err != nil {
			return nil, syslutil.Exitf(ParseError, err.Error())
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
		importer.Protobuf,
	}
	return importer.GuessFileType(fileName, false, file, ParserFormats)
}

func (p *Parser) RestrictToLocalImport() {
	// if root is not defined, only relative imports are allowed
	p.allowAbsoluteImport = false
}

func (p *Parser) Set(settings Settings) {
	p.Settings = settings
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
	r, err := NewReader(fs)
	if err != nil {
		return nil, err
	}
	return p.Parse(filename, r)
}

// ParseFromFsWithVendor parses a sysl definition from an afero filesystem, and vendor remote files in root dir
func (p *Parser) ParseFromFsWithVendor(filename string, fs afero.Fs) (*sysl.Module, error) {
	if p.MaxImportDepth > 0 {
		return nil, fmt.Errorf("can not limit import depth while vendoring")
	}

	r, err := NewReader(fs)
	if err != nil {
		return nil, err
	}

	r.(*remotefs.RemoteFs).Vendor(SyslRootDir(fs))
	return p.Parse(filename, r)
}

type srcInput struct {
	src   importDef
	input string
}

type fileInfo struct {
	imports []importDef
	src     srcInput
}

type retrievedListIndex string

type retrievedList struct {
	l     map[retrievedListIndex]*fileInfo
	mutex sync.Mutex
}

func fileNameToIndex(filename string) retrievedListIndex {
	ret := cleanImportFilename(filename)
	// Only use a single version of each file, trim off any version for indexing
	i := strings.Index(ret, "@")
	if i > -1 {
		ret = ret[:i]
	}

	return retrievedListIndex(ret)
}

// Parse parses a sysl definition from an retriever interface
func (p *Parser) Parse(resource string, reader reader.Reader) (*sysl.Module, error) {
	listener := NewTreeShapeListener()
	listener.lint()

	if filepath.Ext(resource) == "" {
		resource += syslExt
	}

	retrieved := retrievedList{make(map[retrievedListIndex]*fileInfo), sync.Mutex{}}

	if err := p.collectSpecs(
		context.Background(),
		newImportDef(resource),
		reader,
		&retrieved,
		p.MaxImportDepth,
		0,
	); err != nil {
		return nil, err
	}

	specs := []srcInput{}
	flattenSpecs(&specs, resource, &retrieved)

	if p.OperationSummary {
		out := struct {
			FilesProcessed []string `json:"filesProcessed"`
		}{
			FilesProcessed: make([]string, len(specs)),
		}
		for i := range specs {
			out.FilesProcessed[i] = specs[i].src.filename
		}
		b, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			return nil, err
		}
		_, _ = os.Stdout.Write(b)
		_, _ = os.Stdout.Write([]byte("\n"))
	}

	if p.NoParsing {
		return listener.module, nil
	}

	return p.parseSpecs(specs, listener)
}

func (p *Parser) parseSpecs(specs []srcInput, listener *TreeShapeListener) (*sysl.Module, error) { //nolint:funlen
	// Import all foreign types in parallel
	type syslInput struct {
		src             importDef
		syslProtoImport *sysl.Module
		err             error
		str             antlr.CharStream
	}
	syslInputs := make([]syslInput, len(specs))

	g := new(errgroup.Group)
	for i := range specs {
		v := &specs[i]
		out := &syslInputs[i]
		g.Go(func() error {
			out.src = v.src

			// Import Sysl Proto
			out.syslProtoImport, out.err = pbutil.FromPBStringContents(v.src.filename, v.input)
			if !errors.Is(pbutil.ErrUnknownExtension, out.err) {
				return nil
			}

			fsinput := &fsFileStream{antlr.NewInputStream(v.input), v.src.filename}
			var err error
			out.str, err = importForeign(v.src, fsinput)
			if err != nil {
				return err
			}

			return nil
		})
	}
	gerr := g.Wait()
	if gerr != nil {
		return nil, gerr
	}

	for _, v := range syslInputs {
		src := v.src
		logrus.Debug("Parsing: ", src.filename)

		version := ""
		if syslutil.IsRemoteImport(src.filename) {
			v := strings.Split(src.filename, "@")
			version = v[len(v)-1]
		}
		srcCtxFile := cleanImportFilename(src.filename)
		listener.base = importDir(src.filename)
		listener.sc = sourceCtxHelper{srcCtxFile, version}

		// Import Sysl Proto
		if !errors.Is(pbutil.ErrUnknownExtension, v.err) {
			if v.err != nil {
				return nil, fmt.Errorf("error parsing %s: %w", src.filename, v.err)
			}
			if v.syslProtoImport != nil {
				// Merge structs recursively
				if err := mergo.Merge(listener.module, v.syslProtoImport); err != nil {
					return nil, err
				}
			}

			continue
		}

		tree, err := parseString(src.filename, v.str)
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

// Takes a starting file and flattens all the imports that were already retrieved into an ordered list (recursively)
func flattenSpecs(specs *[]srcInput, filename string, retrieved *retrievedList) {
	// Only add each file once
	filenameIndex := fileNameToIndex(filename)
	for _, si := range *specs {
		if fileNameToIndex(si.src.filename) == filenameIndex {
			return
		}
	}

	fi, found := retrieved.l[filenameIndex]
	// If it wasn't found it must have been ignored due to depth
	if found {
		*specs = append(*specs, fi.src)
		for _, v := range fi.imports {
			flattenSpecs(specs, v.filename, retrieved)
		}
	}
}

// collectSpecs retrieves the contents for a sourceFile. It then parses the source to find all imports and recursively
// (and in parallel) retrieves those as well. All the results are placed into the retrievedList.
func (p *Parser) collectSpecs(
	ctx context.Context,
	source importDef,
	reader reader.Reader,
	retrieved *retrievedList,
	maxImportDepth, currentImportDepth int,
) error {
	if maxImportDepth > 0 && currentImportDepth >= maxImportDepth {
		return nil
	}

	filenameIndex := fileNameToIndex(source.filename)
	retrieved.mutex.Lock()
	if fi, has := retrieved.l[filenameIndex]; has {
		retrieved.mutex.Unlock()

		if !p.NoDifferentVersionCheck {
			appname1 := strings.ReplaceAll(fi.src.src.appname, " :: ", "::")
			appname2 := strings.ReplaceAll(source.appname, " :: ", "::")
			if appname1 != appname2 {
				return syslutil.Exitf(ImportError, fmt.Sprintf(
					"%#v imported as different appnames: '%v' and '%v'", filenameIndex, appname1, appname2,
				))
			}

			ver1 := ""
			i := strings.Index(fi.src.src.filename, "@")
			if i > -1 {
				ver1 = fi.src.src.filename[i+1:]
			}

			ver2 := ""
			i = strings.Index(source.filename, "@")
			if i > -1 {
				ver2 = source.filename[i+1:]
			}

			// treat master/main/develop as default
			switch ver1 {
			case "master", "main", "develop":
				ver1 = ""
			}
			switch ver2 {
			case "master", "main", "develop":
				ver2 = ""
			}

			if ver1 != ver2 {
				return syslutil.Exitf(ImportError, fmt.Sprintf(
					"%#v imported as different versions: '%v' and '%v'", filenameIndex, ver1, ver2,
				))
			}
		}

		return nil
	}
	fi := &fileInfo{}
	fi.src.src = source
	retrieved.l[filenameIndex] = fi
	retrieved.mutex.Unlock()

	content, hash, branch, err := reader.ReadHashBranch(ctx, source.filename)
	if err != nil {
		return syslutil.Exitf(ImportError, fmt.Sprintf(
			"error reading %#v: \n%v\n", source.filename, err,
		))
	}
	fi.src.input = string(content)

	importsInput := extractImports(source.filename, content)

	if importsInput.Len() == 0 {
		return nil
	}

	version := branch
	if version == "" {
		version = hash.String()
	}

	children, err := parseImports(source, sourceCtxHelper{source.filename, version}, importsInput.String())
	if err != nil {
		return err
	}

	fi.imports = children

	g := new(errgroup.Group)
	for _, c := range children {
		c := c
		g.Go(func() error {
			return p.collectSpecs(ctx, c, reader, retrieved, maxImportDepth, currentImportDepth+1)
		})
	}

	err = g.Wait()
	if err != nil {
		return syslutil.Exitf(ImportError, fmt.Sprintf(
			"error reading %#v: \n%v", source.filename, err,
		))
	}

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
	if len(appPath) == 0 && len(typePath) > 1 {
		// local field ref
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
					continue
				}
				attrs = x.Tuple.GetAttrDefs()
			case *sysl.Type_Relation_:
				if len(x.Relation.AttrDefs) == 0 {
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
		splitPath := strings.Split(typeName, ".")
		if len(splitPath) > 1 && injectType(app.Types[typeName], app.Types, splitPath) {
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
	m, err := p.ParseFromFs(model, fs)
	if err != nil {
		return nil, "", err
	}
	modelAppName := getDefaultAppName(m)
	return m, modelAppName, nil
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
