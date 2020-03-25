package parse

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/pkg/importer"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	parser "github.com/anz-bank/sysl/pkg/grammar"
	"github.com/anz-bank/sysl/pkg/msg"
	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
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
		return nil, Exitf(ParseError, fmt.Sprintf("%s has syntax errors\n", filename))
	}
	return tree, nil
}

func guessMode(filename string) string {
	switch filepath.Ext(filename) {
	case ".sysl":
		return "~sysl"
	case ".yaml", ".json":
		return "~swagger"
	default:
		return ""
	}
}

func importForeign(def importDef, input antlr.CharStream) (antlr.CharStream, error) {
	od := importer.OutputData{
		AppName: def.appname,
		Package: def.pkg,
	}
	var text string
	var err error
	logger := logrus.StandardLogger()

	if def.mode == "" {
		def.mode = guessMode(def.filename)
	}
	od.Mode = def.mode

	switch def.mode {
	case "~sysl":
		return input, nil
	case "~swagger":
		text, err = importer.LoadSwaggerText(od, input.GetText(0, input.Size()), logger)
	default:
		return nil, Exitf(ParseError, fmt.Sprintf("%s has unknown format - (%s)\n", def.filename, def.mode))
	}

	return antlr.NewInputStream(text), err
}

func (p *Parser) RestrictToLocalImport() {
	// if root is not defined, only relative imports are allowed
	p.allowAbsoluteImport = false
}

func (p *Parser) Parse(filename string, fs afero.Fs) (*sysl.Module, error) {
	if !strings.HasSuffix(filename, ".sysl") {
		filename += ".sysl"
	}
	if !fileExists(filename, fs) {
		return nil, Exitf(ImportError, "input file does not exist: %#v", filename)
	}

	imported := map[string]struct{}{}

	var apps = make(map[string]*sysl.Application) // the set of all applications from all files
	var imports []importDef                       // the list of imports yet to be parsed

	source := importDef{
		filename: filename,
	}

	for {

		newImports, err := parseAntlr(source, fs, apps)
		if err != nil {
			return nil, err
		}

		imports = append(imports, newImports...)

		if len(imports) == 0 {
			break
		}

		duplicateImportCheck := func(list []importDef) {
			// assume that a file will have less than 256 imports
			set := make(map[string]byte)
			for _, value := range list {
				if count := set[value.filename]; count == 1 {
					logrus.Warnf("Duplicate import: '%s' in file: '%s'\n", value.filename, source.filename)
				}
				set[value.filename]++
			}
		}
		duplicateImportCheck(newImports)

		imported[source.filename] = struct{}{}

		for len(imports) > 0 {
			source = imports[0]
			imports = imports[1:]
			if _, has := imported[source.filename]; !has {
				if !p.allowAbsoluteImport && strings.HasPrefix(source.filename, "/") {
					return nil, Exitf(2,
						"error importing: importing outside current directory is only allowed when root is defined")
				}
				break
			}
		}

		if _, has := imported[source.filename]; has {
			break
		}
	}

	rootModule := &sysl.Module{Apps: apps}
	p.postProcess(rootModule)
	return rootModule, nil
}

// Parse the given source definition, modifying the app contents in the process.
// Return the list of import definitions declared in the source but not parsed.
func parseAntlr(source importDef, fs afero.Fs, apps map[string]*sysl.Application) ([]importDef, error) {
	filename := source.filename
	logrus.Debugf("Parsing: " + filename)

	fsinput, err := newFSFileStream(filename, fs)
	if err != nil {
		return nil, Exitf(ImportError, fmt.Sprintf("error parsing %#v: %v\n", filename, err))
	}

	listener := NewTreeShapeListener()
	listener.module.Apps = apps
	listener.sc = sourceCtxHelper{source.filename}
	listener.base = filepath.Dir(filename)

	input, err := importForeign(source, fsinput)
	if err != nil {
		return nil, err
	}

	tree, err := parseString(filename, input)
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
			logrus.Warnf("%s::%s calls non-existent App: %s",
				appname, epname, s.Call.Target.Part)
			return false
		}
		_, valid := app.Endpoints[s.Call.Endpoint]
		if !valid {
			logrus.Warnf("%s::%s calls non-existent App <- Endpoint (%s <- %s)",
				appname, epname, s.Call.Target.Part, s.Call.Endpoint)
		}
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
	// logrus.Printf("found anonymous type\n")
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

func (p *Parser) postProcess(mod *sysl.Module) {
	appNames := make([]string, 0, len(mod.Apps))
	for a := range mod.Apps {
		appNames = append(appNames, a)
	}
	sort.Strings(appNames)

	for _, appName := range appNames {
		app := mod.Apps[appName]

		if app.Mixin2 != nil {
			for _, src := range app.Mixin2 {
				srcApp := syslutil.GetApp(src.Name, mod)
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
				attrs = x.Tuple.GetAttrDefs()
			case *sysl.Type_Relation_:
				attrs = x.Relation.GetAttrDefs()
			}
			for fieldname, t := range attrs {
				if x := t.GetTypeRef(); x != nil {
					refApp := app
					refName := x.GetRef().GetPath()[0]
					if refName == "string_8" {
						continue
					}
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
							last := len(x.GetRef().GetPath()) - 1
							field = x.GetRef().GetPath()[last]
							_, has = refApp.Types[field]
						}
						if !has {
							logrus.Warnf("Field %#v (type %#v) refers to Field %#v (type %#v) in app %#v",
								fieldname, typeName, field, refName, appName)
						}
					} else {
						logrus.Warnf("Field %#v (type %#v) refers to type %s in app %#v",
							fieldname, typeName, refName, appName)
					}
				}
			}
		}
		p.inferTypes(mod, appName)
		collectorPubSubCalls(appName, app)
	}
	checkEndpointCalls(mod)
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
	mod, err := p.Parse(model, fs)
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

func NewParser() *Parser {
	return &Parser{
		AssignTypes:         map[string]TypeData{},
		LetTypes:            map[string]TypeData{},
		Messages:            map[string][]msg.Msg{},
		allowAbsoluteImport: true,
	}
}
