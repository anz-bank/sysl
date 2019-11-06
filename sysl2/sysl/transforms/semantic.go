package transforms

import (
	"sort"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/anz-bank/sysl/sysl2/sysl/eval"
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
)

type semantic struct {
	base
	filenames *sysl.View
}

// assume args are
//  app <: sysl.App and
//  type <: sysl.Type
//  typeName <: string
//  module <: sysl.Module

func (t *semantic) Apply(mod *sysl.Module, appNames ...string) map[string]*sysl.Value {
	output := map[string]*sysl.Value{}
	for _, name := range appNames {
		filenamesData := createScope(t.filenames, mod, mod.Apps[name])
		outputData := createScope(t.view, mod, mod.Apps[name])

		var result, filenames *sysl.Value
		if tp := outputData.tp; tp != nil {
			types := mod.Apps[name].Types
			filenames = eval.MakeValueList()
			result = eval.MakeValueList()
			var tNames []string
			for tName := range types {
				tNames = append(tNames, tName)
			}
			sort.Strings(tNames)
			for _, tName := range tNames {
				outputData.scope[tp.name] = eval.MakeValueString(tName)
				outputData.scope[tp.t] = eval.TypeToValue(types[tName])
				eval.AppendItemToValueList(filenames.GetList(), t.eval(t.filenames, outputData.scope))
				eval.AppendItemToValueList(result.GetList(), t.eval(t.view, outputData.scope))
			}
		} else {
			result = t.eval(t.view, outputData.scope)
			filenames = t.eval(t.filenames, filenamesData.scope)
		}
		// filename should either be a map with a single key:"filename", or a list of the same length as the result value
		switch {
		case filenames.GetMap() != nil:
			filename := filenames.GetMap().Items["filename"].GetS()
			output[filename] = result
		case filenames.GetList() != nil && result.GetList() != nil:
			fileValues := filenames.GetList().Value
			for i, v := range result.GetList().Value {
				filename := fileValues[i].GetMap().Items["filename"].GetS()
				output[filename] = v
			}
		}
	}
	return output
}

type typeParams struct {
	t    string
	name string
}
type transformData struct {
	scope eval.Scope
	tp    *typeParams
}

func createScope(view *sysl.View, mod *sysl.Module, app *sysl.Application) transformData {
	res := transformData{
		scope: eval.Scope{},
	}
	tp := typeParams{}
	for _, param := range view.Param {
		_, detail := syslutil.GetTypeDetail(param.Type)
		switch detail {
		case "App", "Application":
			res.scope.AddApp(param.Name, app)
		case "Module":
			res.scope.AddModule(param.Name, mod)
		case "Type":
			tp.t = param.Name
		case "STRING": // Assuming it must be the type name
			tp.name = param.Name
		}
	}
	if tp.name != "" && tp.t != "" {
		res.tp = &tp
	}
	return res
}
