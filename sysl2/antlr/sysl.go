package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"anz-bank/sysl/src/proto"
	"anz-bank/sysl/sysl2/antlr/grammar"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

var (
	root   = flag.String("root", ".", "sysl root directory for input files (default: .)")
	output = flag.String("o", "", "output file name")
)

func init() {
	flag.Parse()
}

// JsonPB ...
func JsonPB(m *sysl.Module, filename string) bool {
	ma := jsonpb.Marshaler{Indent: " ", EmitDefaults: false}
	f, err := os.Create(filename)
	err = ma.Marshal(f, m)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// TextPB ...
func TextPB(m *sysl.Module, filename string) bool {
	if m == nil {
		fmt.Println("input module is nil")
		return false
	}

	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return false
	}
	err = proto.MarshalText(f, m)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// SyslParserErrorListener ...
type SyslParserErrorListener struct {
	*antlr.DefaultErrorListener
	hasErrors bool
}

// SyntaxError ...
func (d *SyslParserErrorListener) SyntaxError(
	recognizer antlr.Recognizer, offendingSymbol interface{},
	line, column int, msg string, e antlr.RecognitionException) {
	d.hasErrors = true
}

// func (d *SyslParserErrorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
// 	// d.hasErrors = false
// }

func getAppName(appname *sysl.AppName) string {
	app_name := appname.Part[0]

	for i := 1; i < len(appname.Part); i++ {
		app_name += " :: " + appname.Part[i]
	}
	return app_name
}

func getApp(appName *sysl.AppName, mod *sysl.Module) *sysl.Application {
	if app, has := mod.Apps[getAppName(appName)]; has {
		return app
	}
	return nil
}

func isAbstract(app *sysl.Application) bool {
	patterns, has := app.Attrs["patterns"]
	if has == false {
		return false
	}
	if x := patterns.GetA(); x != nil {
		for _, y := range x.Elt {
			if y.GetS() == "abstract" {
				return true
			}
		}
	}
	return false
}

func isSameApp(a *sysl.AppName, b *sysl.AppName) bool {
	if len(a.Part) != len(b.Part) {
		return false
	}
	for i := range a.Part {
		if a.Part[i] != b.Part[i] {
			return false
		}
	}
	return true
}

func isSameCall(a *sysl.Call, b *sysl.Call) bool {
	return isSameApp(a.Target, b.Target) && a.Endpoint == b.Endpoint
}

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
			applied = true
			if dst.Attrs == nil {
				dst.Attrs = make(map[string]*sysl.Attribute)
			}
			mergeAttrs(src.Attrs, dst.Attrs)
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
	valid := false
	switch s := dst.GetStmt().(type) {
	case *sysl.Statement_Cond:
		stmts = s.Cond.Stmt
	case *sysl.Statement_Alt:
		for _, c := range s.Alt.Choice {
			for _, ss := range c.Stmt {
				valid = checkCalls(mod, appname, epname, ss)
				if !valid {
					break
				}
			}
		}
		return valid
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
			fmt.Printf("%s::%s calls non-existant App: %s\n", appname, epname, s.Call.Target.Part)
			return false
		}
		_, valid = app.Endpoints[s.Call.Endpoint]
		if !valid {
			fmt.Printf("%s::%s calls non-existant App <- Endpoint (%s <- %s)\n", appname, epname, s.Call.Target.Part, s.Call.Endpoint)
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
		valid = checkCalls(mod, appname, epname, stmt)
		if !valid {
			break
		}
	}
	return valid
}

func collectorPubSubCalls(mod *sysl.Module) {
	for appName, app := range mod.Apps {
		// add attribtes collected in collector stmts to
		if endpoint, has := app.Endpoints[`.. * <- *`]; has {
			for _, collector_stmt := range endpoint.Stmt {
				switch x := collector_stmt.Stmt.(type) {
				case *sysl.Statement_Action:
					modify_ep := app.Endpoints[x.Action.Action]
					if modify_ep == nil {
						fmt.Printf("App (%s) calls non-existant endpoint (%s)\n", appName, x.Action.Action)
						continue
					}
					if modify_ep.Attrs == nil {
						modify_ep.Attrs = make(map[string]*sysl.Attribute)
					}
					mergeAttrs(collector_stmt.Attrs, modify_ep.Attrs)
				case *sysl.Statement_Call:
					applied := false

					for call_epname, call_endpoint := range app.Endpoints {
						if call_epname == `.. * <- *` {
							continue
						}
						for _, call_stmt := range call_endpoint.Stmt {
							applied = applyAttributes(collector_stmt, call_stmt) || applied
						}
					}
					if !applied {
						fmt.Printf("Unused template (%s <- %s) in app %s\n", x.Call.Target.Part, x.Call.Endpoint, appName)
					}
				default:
					panic("unhandled type:")
				}
			}
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

func postProcess(mod *sysl.Module) {
	var appNames []string
	for a := range mod.Apps {
		appNames = append(appNames, a)
	}
	sort.Strings(appNames)

	for _, appName := range appNames {
		app := mod.Apps[appName]

		if app.Mixin2 != nil {
			for _, src := range app.Mixin2 {
				src_app := getApp(src.Name, mod)
				if isAbstract(src_app) == false {
					fmt.Printf("mixin App (%s) should be ~abstract\n", getAppName(src.Name))
					continue
				}
				if app.Types == nil {
					app.Types = make(map[string]*sysl.Type)
				}
				for k, v := range src_app.Types {
					if _, has := app.Types[k]; !has {
						app.Types[k] = v
					} else {
						fmt.Printf("Type %s defined in %s and in %s\n", k, appName, getAppName(src.Name))
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
					var refName string
					refName = x.GetRef().GetPath()[0]
					if refName == "string_8" {
						continue
					}
					refType, has := refApp.Types[refName]
					if has == false {
						fmt.Printf("1:Field %s (type %s) refers to type (%s) in app (%s)\n", fieldname, typeName, refName, appName)
					} else {
						var ref_attrs map[string]*sysl.Type

						switch refType.Type.(type) {
						case *sysl.Type_Tuple_:
							refType, _ := refApp.Types[refName].Type.(*sysl.Type_Tuple_)
							ref_attrs = refType.Tuple.GetAttrDefs()
						case *sysl.Type_Relation_:
							refType, _ := refApp.Types[refName].Type.(*sysl.Type_Relation_)
							ref_attrs = refType.Relation.GetAttrDefs()
						}
						var field string
						var has bool
						if len(x.GetRef().GetPath()) > 1 {
							last := len(x.GetRef().GetPath()) - 1
							field = x.GetRef().GetPath()[last]
							_, has = ref_attrs[field]
						} else if len(x.GetRef().GetPath()) == 1 {
							last := len(x.GetRef().GetPath()) - 1
							field = x.GetRef().GetPath()[last]
							_, has = refApp.Types[field]
						}
						if has == false {
							fmt.Printf("2:Field %s (type %s) refers to Field (%s) in app (%s)/type (%s)\n", fieldname, typeName, field, appName, refName)
						}
					}
				}
			}
		}
	}
	checkEndpointCalls(mod)
	collectorPubSubCalls(mod)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
func dirExists(dirName string) bool {
	info, err := os.Stat(dirName)
	return err == nil && info.IsDir()

}

// Parse ...
func Parse(filename string, root string) *sysl.Module {
	if root == "" {
		root = "."
	}
	if !dirExists(root) {
		fmt.Println("root directory does not exist")
		return nil
	}
	root, _ = filepath.Abs(root)

	if !fileExists(filename) {
		if filename[len(filename)-5:] != ".sysl" {
			filename = filename + ".sysl"
		}
		temp := root + "/" + filename

		if !fileExists(temp) {
			fmt.Printf("input file does not exist\nFilename: %s\n", temp)
			return nil
		}
		filename = temp
	}
	var empty struct{}
	filename, _ = filepath.Abs(filename)
	imported := make(map[string]struct{})
	listener := NewTreeShapeListener(root)
	errorListener := SyslParserErrorListener{}

	for {
		fmt.Println(filename)
		input, err := antlr.NewFileStream(filename)
		if err != nil {
			fmt.Printf("%v,\n%s has errors\n", err, filename)
			break
		}
		listener.base = filepath.Dir(filename)
		lexer := parser.NewSyslLexer(input)
		stream := antlr.NewCommonTokenStream(lexer, 0)
		p := parser.NewSyslParser(stream)
		// p.GetInterpreter().SetPredictionMode(antlr.PredictionModeLLExactAmbigDetection)
		p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
		p.AddErrorListener(&errorListener)

		p.BuildParseTrees = true
		tree := p.Sysl_file()
		if errorListener.hasErrors {
			fmt.Printf("%s has errors\n", filename)
			break
		}

		antlr.ParseTreeWalkerDefault.Walk(listener, tree)
		if len(listener.imports) == 0 {
			break
		}
		imported[filename] = empty

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
	return listener.module
}

func main() {
	fmt.Printf("Args: %v\n", flag.Args())
	fmt.Printf("Root: %s\n", *root)
	fmt.Printf("Module: %s\n", flag.Arg(0))
	mod := Parse(flag.Arg(0), *root)
	if mod != nil {
		TextPB(mod, *output)
	} else {
		os.Exit(1)
	}
	// JsonPB(mod, *output)
}
