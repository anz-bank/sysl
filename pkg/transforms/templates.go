package transforms

import (
	"fmt"
	"log"
	"sort"

	"github.com/anz-bank/sysl/pkg/eval"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
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
	logrus.Tracef("Apply evalRes: %s", evalRes.String())

	fn := func(v *sysl.Value) (filename string, data string, err error) {
		valMap := v.GetMap().Items["apps"].GetSet().Value[0].GetMap().Items["app"].GetMap()
		if valMap == nil {
			return "", "", fmt.Errorf("incorrect return type")
		}
		val, ok := valMap.Items["Filename"]
		if !ok {
			return "", "", fmt.Errorf("'Filename' not set")
		}
		filename = val.GetS()
		val, ok = valMap.Items["Data"]
		if !ok {
			return "", "", fmt.Errorf("'Data' not set")
		}

		data = eval.UnaryString(val).GetS()
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
	logrus.Tracef("Apply result: %+v", result)
	return result
}

const templateInputType = "sysl.TemplateInput"
