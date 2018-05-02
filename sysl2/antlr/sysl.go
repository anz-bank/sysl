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
	// d.hasErrors = true
}

func (d *SyslParserErrorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	// d.hasErrors = false
}

func getAppName(app *sysl.AppName, mod *sysl.Module) *sysl.Application {
	app_name := app.Part[0]

	for i := 1; i < len(app.Part); i++ {
		app_name = " :: " + app.Part[i]
	}
	if app, has := mod.Apps[app_name]; has {
		return app
	}
	return nil
}

func postProcess(mod *sysl.Module) {
	for appName, app := range mod.Apps {

		for epname, endpoint := range app.Endpoints {
			if endpoint.Source != nil {
				src_app := getAppName(endpoint.Source, mod)
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
					// var contextPath string
					var refName string
					// if x.GetContext() != nil && len(x.GetRef().GetPath()) == 1 && len(x.GetContext().Path) > 0 {
					// 	contextPath = strings.Join(x.GetContext().Path, ".")
					// 	refName = contextPath + "."
					// }
					refName = x.GetRef().GetPath()[0]
					if refName == "string_8" {
						continue
					}
					_, has := refApp.Types[refName]
					if has == false {
						fmt.Printf("1:Field %s (type %s) refers to type (%s) in app (%s)\n", fieldname, typeName, refName, appName)
					} else {
						var ref_attrs map[string]*sysl.Type

						switch types.Type.(type) {
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
	}

	postProcess(listener.module)
	return listener.module
}

func main() {
	// fmt.Printf("Reading file %s\n", os.Args[1])
	mod := Parse(os.Args[1], "")
	TextPB(mod)
}
