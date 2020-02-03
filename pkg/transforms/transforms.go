package transforms

import (
	"fmt"

	"github.com/anz-bank/sysl/pkg/eval"
	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
)

type Worker interface {
	Apply(mod *sysl.Module, appNames ...string) map[string]*sysl.Value
}

func NewWorker(transformMod *sysl.Module, appName, viewName string) (Worker, error) {
	app, has := transformMod.Apps[appName]
	if !has {
		return nil, fmt.Errorf("app '%s' not found in transform module", appName)
	}
	view, has := app.Views[viewName]
	if !has {
		return nil, fmt.Errorf("view '%s' not found in transform app", viewName)
	}
	b := base{
		mod:  transformMod,
		app:  app,
		view: view,
	}

	if len(view.Param) == 1 {
		_, detail := syslutil.GetTypeDetail(view.Param[0].Type)
		if detail == templateInputType {
			return &templated{base: b}, nil
		}
	}
	filenames, has := app.Views["filename"]
	if !has {
		return nil, fmt.Errorf("view '%s' not found in transform app", "filename")
	}
	return &semantic{base: b, filenames: filenames}, nil
}

type base struct {
	mod  *sysl.Module
	app  *sysl.Application
	view *sysl.View
}

func (b *base) eval(view *sysl.View, scope eval.Scope) *sysl.Value {
	if view.Expr.Type == nil {
		view.Expr.Type = view.RetType
	}
	return eval.EvaluateApp(b.app, view, scope)
}
