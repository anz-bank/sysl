package transforms

import (
	"fmt"
	"log"
	"sort"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/anz-bank/sysl/sysl2/sysl/eval"
)

type templated struct {
	base
}

func (t *templated) Apply(mod *sysl.Module, appNames ...string) map[string]*sysl.Value {
	input := eval.Scope{}
	input.AddModule("Module", mod)

	apps := eval.MakeValueList()
	var aNames []string
	if len(appNames) == 0 {
		for app := range mod.Apps {
			aNames = append(aNames, app)
		}
	} else {
		aNames = append(aNames, appNames...)
	}
	sort.Strings(aNames)
	for _, app := range aNames {
		s := eval.Scope{}
		s.AddApp(app, mod.Apps[app])
		eval.AppendItemToValueList(apps.GetList(), s[app])
	}
	input["Apps"] = apps

	s := eval.Scope{}
	s[t.view.Param[0].Name] = input.ToValue()

	evalRes := t.eval(t.view, s)

	fn := func(v *sysl.Value) (filename string, data string, err error) {
		if m := v.GetMap(); m == nil {
			return "", "", fmt.Errorf("incorrect return type")
		}
		val, ok := v.GetMap().Items["Filename"]
		if !ok {
			return "", "", fmt.Errorf("'Filename' not set")
		}
		filename = val.GetS()
		val, ok = v.GetMap().Items["Data"]
		if !ok {
			return "", "", fmt.Errorf("'Data' not set")
		}
		data = val.GetS()
		return
	}

	result := map[string]*sysl.Value{}
	if evalRes.GetList() != nil {
		for _, val := range evalRes.GetList().Value {
			if fname, data, err := fn(val); err == nil {
				result[fname] = eval.MakeValueString(data)
			}
		}
	} else if fname, data, err := fn(evalRes); err == nil {
		result[fname] = eval.MakeValueString(data)
	}
	log.Printf("%+v", result)
	return result
}

const templateInputType = "sysl.TemplateInput"
