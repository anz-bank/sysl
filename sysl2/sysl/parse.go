package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	sysl "github.com/anz-bank/sysl/src/proto"
	parser "github.com/anz-bank/sysl/sysl2/sysl/grammar"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Parse parses a Sysl file under a specified root.
func Parse(filename string, root string) (*sysl.Module, error) {
	if root == "" {
		root = "."
	}
	if !dirExists(root) {
		return nil, exitf(ImportError, "root directory does not exist")
	}
	root, _ = filepath.Abs(root)
	return FSParse(filename, http.Dir(root))
}

// FSParse ...
func FSParse(filename string, fs http.FileSystem) (*sysl.Module, error) {
	if !strings.HasSuffix(filename, ".sysl") {
		filename += ".sysl"
	}
	if !fileExists(filename, fs) {
		return nil, exitf(ImportError, "input file does not exist: %#v", filename)
	}
	imported := map[string]struct{}{}
	listener := NewTreeShapeListener(fs)
	errorListener := SyslParserErrorListener{}

	for {
		logrus.Debugf("Parsing: " + filename)
		input, err := newFSFileStream(filename, fs)
		if err != nil {
			return nil, exitf(ImportError, fmt.Sprintf("error parsing %#v: %v\n", filename, err))
		}
		listener.filename = filename
		listener.base = filepath.Dir(filename)
		lexer := parser.NewSyslLexer(input)
		stream := antlr.NewCommonTokenStream(lexer, 0)
		p := parser.NewSyslParser(stream)
		p.GetInterpreter().SetPredictionMode(antlr.PredictionModeSLL)
		p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
		p.AddErrorListener(&errorListener)

		p.BuildParseTrees = true
		tree := p.Sysl_file()
		if errorListener.hasErrors {
			return nil, exitf(ParseError, fmt.Sprintf("%s has syntax errors\n", filename))
		}

		antlr.ParseTreeWalkerDefault.Walk(listener, tree)
		if len(listener.imports) == 0 {
			break
		}
		imported[filename] = struct{}{}

		for len(listener.imports) > 0 {
			filename = listener.imports[0]
			listener.imports = listener.imports[1:]
			if _, has := imported[filename]; !has {
				break
			}
		}
		if _, has := imported[filename]; has {
			break
		}
	}

	postProcess(listener.module)
	return listener.module, nil
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
		if isSameCall(src.GetCall(), s.Call) {
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
		app := getApp(s.Call.Target, mod)
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

// for nested transform's Type
func inferExprType(mod *sysl.Module,
	appName string,
	expr *sysl.Expr, top bool,
	anonCount int) (*sysl.Type, int) {

	switch t := expr.Expr.(type) {
	case *sysl.Expr_Transform_:
		for _, stmt := range t.Transform.Stmt {
			if stmt.GetLet() != nil {
				_, anonCount = inferExprType(mod, appName, stmt.GetLet().Expr, false, anonCount)
			} else if stmt.GetAssign() != nil {
				_, anonCount = inferExprType(mod, appName, stmt.GetAssign().Expr, false, anonCount)
			}
		}

		if !top && expr.Type == nil {
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
					if aexpr.GetTransform() == nil {
						panic("expression should be of type transform")
					}
					ftype := aexpr.Type
					setof := ftype.GetSet() != nil
					if setof {
						ftype = ftype.GetSet()
					}
					if ftype.GetTypeRef() == nil {
						panic("transform type should be type_ref")
					}
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
		}
	case *sysl.Expr_Relexpr:
		relexpr := t.Relexpr
		if relexpr.Op == sysl.Expr_RelExpr_RANK {
			if !top && expr.Type == nil {
				type1, c := inferExprType(mod, appName, relexpr.Target, true, anonCount)
				anonCount = c
				logrus.Printf(type1.String())
			}
		}
	case *sysl.Expr_Literal:
		expr.Type = valueTypeToSysl(t.Literal)
	case *sysl.Expr_List_:
		expr.Type, _ = inferExprType(mod, appName, t.List.Expr[0], true, anonCount)
	case *sysl.Expr_Ifelse:
		exprTypeIfTrue, _ := inferExprType(mod, appName, t.Ifelse.GetIfTrue(), true, anonCount)
		expr.Type = exprTypeIfTrue
		if t.Ifelse.GetIfFalse() != nil {
			exprTypeIfFalse, _ := inferExprType(mod, appName, t.Ifelse.GetIfFalse(), true, anonCount)
			// TODO if types are not equal, raise an error. Then remove following 3 lines
			if exprTypeIfFalse != nil {
				expr.Type = exprTypeIfFalse
			}
		}
	case *sysl.Expr_Binexpr:
		exprTypeLHS, _ := inferExprType(mod, appName, t.Binexpr.GetLhs(), true, anonCount)
		expr.Type = exprTypeLHS
		exprTypeRHS, _ := inferExprType(mod, appName, t.Binexpr.GetRhs(), true, anonCount)
		// TODO if types are not equal, raise an error. Then remove following 3 lines
		if exprTypeRHS != nil {
			expr.Type = exprTypeRHS
		}

	case *sysl.Expr_Unexpr:
		expr.Type, _ = inferExprType(mod, appName, expr.GetUnexpr().GetArg(), true, anonCount)

	default:
		// TODO Handle expression
		logrus.Debug("[Parse.infer_expr_type] Unhandled type", t)
	}

	return expr.Type, anonCount
}

func inferTypes(mod *sysl.Module, appName string) {
	for _, view := range mod.Apps[appName].Views {
		if HasPattern(view.Attrs, "abstract") {
			continue
		}
		inferExprType(mod, appName, view.Expr, true, 0)
	}
}

func postProcess(mod *sysl.Module) {
	appNames := make([]string, 0, len(mod.Apps))
	for a := range mod.Apps {
		appNames = append(appNames, a)
	}
	sort.Strings(appNames)

	for _, appName := range appNames {
		app := mod.Apps[appName]

		if app.Mixin2 != nil {
			for _, src := range app.Mixin2 {
				srcApp := getApp(src.Name, mod)
				if !HasPattern(srcApp.Attrs, "abstract") {
					logrus.Warnf("mixin App (%s) should be ~abstract", getAppName(src.Name))
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
							k, appName, getAppName(src.Name))
					}
				}
				for k, v := range srcApp.Views {
					if _, has := app.Views[k]; !has {
						app.Views[k] = v
					} else {
						logrus.Warnf("View %s defined in %s and in %s",
							k, appName, getAppName(src.Name))
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
		inferTypes(mod, appName)
		collectorPubSubCalls(appName, app)
	}
	checkEndpointCalls(mod)
}

func valueTypeToSysl(value *sysl.Value) *sysl.Type {
	switch value.Value.(type) {
	case *sysl.Value_B:
		return &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_BOOL,
			},
		}

	case *sysl.Value_I:
		return &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_INT,
			},
		}
	case *sysl.Value_D:
		return &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_FLOAT,
			},
		}
	case *sysl.Value_S:
		return &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_STRING,
			},
		}
	case *sysl.Value_Decimal:
		return &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_DECIMAL,
			},
		}
	case *sysl.Value_Null_:
		return &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_EMPTY,
			},
		}
	default:
		panic(errors.Errorf("valueTypeToSysl: unhandled type: %v", value))

	}
}
