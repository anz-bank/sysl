package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"anz-bank/sysl/src/proto"
	"anz-bank/sysl/sysl2/antlr/grammar"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/golang/protobuf/proto"
)

// TextPB ...
func TextPB(m *sysl.Module) {
	fmt.Println(proto.MarshalTextString(m))
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

func postProcess(mod *sysl.Module) {
	for appName, app := range mod.Apps {

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

		for epname, endpoint := range app.Endpoints {

			// add attribtes collected in collector stmts to
			if endpoint.Name == `.. * <- *` {
				for _, collector_stmt := range endpoint.Stmt {
					// var modify_ep *sysl.Endpoint
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
						for call_epname, call_endpoint := range app.Endpoints {
							if call_epname == `.. * <- *` {
								continue
							}

							for _, call_stmt := range call_endpoint.Stmt {
								y := call_stmt.GetCall()
								if y != nil {
									if isSameApp(y.Target, x.Call.Target) && y.Endpoint == x.Call.Endpoint {
										if call_stmt.Attrs == nil {
											call_stmt.Attrs = make(map[string]*sysl.Attribute)
										}
										mergeAttrs(collector_stmt.Attrs, call_stmt.Attrs)
									}
								}
							}
						}
					default:
						panic("unhandled type:")
					}
				}
			}

			if endpoint.Source != nil {
				src_app := getApp(endpoint.Source, mod)
				if src_app != nil {
					eventName := strings.TrimSpace(strings.Split(epname, ">")[1])
					if ep, has := src_app.Endpoints[eventName]; has {
						stmt := &sysl.Statement{
							Stmt: &sysl.Statement_Call{
								Call: &sysl.Call{
									Target:   app.Name,
									Endpoint: epname,
								},
							},
						}
						ep.Stmt = append(ep.Stmt, stmt)
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
		return nil
	}
	root, _ = filepath.Abs(root)

	if !fileExists(filename) {
		if filename[len(filename)-5:] != ".sysl" {
			filename = filename + ".sysl"
		}
		temp := root + "/" + filename

		if !fileExists(temp) {
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
		input := antlr.NewFileStream(filename)
		listener.base = filepath.Dir(filename)
		lexer := parser.NewSyslLexer(input)
		stream := antlr.NewCommonTokenStream(lexer, 0)
		p := parser.NewSyslParser(stream)

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
	// fmt.Printf("Reading file %s\n", os.Args[1])
	mod := Parse(os.Args[1], "")
	TextPB(mod)
}
