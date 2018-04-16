package main

import (
	"fmt"
	"os"
	"path/filepath"

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

func postProcess(mod *sysl.Module) {
	for appName, app := range mod.Apps {

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
					// refApp := x.GetRef().GetAppname()
					// if refApp == nil {
					// }
					refName := x.GetRef().GetPath()[0]
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

						field := x.GetRef().GetPath()[1]
						_, has = ref_attrs[field]

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
	base := filepath.Dir(filename)
	imported := make(map[string]struct{})
	listener := NewTreeShapeListener(base, root)

	for {
		input := antlr.NewFileStream(filename)
		lexer := parser.NewSyslLexer(input)
		stream := antlr.NewCommonTokenStream(lexer, 0)
		p := parser.NewSyslParser(stream)

		p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
		errorListener := SyslParserErrorListener{}
		p.AddErrorListener(&errorListener)
		p.BuildParseTrees = true
		tree := p.Sysl_file()
		if errorListener.hasErrors {
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
