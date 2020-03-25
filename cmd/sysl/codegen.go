package main

import (
	"io"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/mod"

	"github.com/anz-bank/sysl/pkg/ebnfparser"

	"github.com/anz-bank/sysl/pkg/eval"
	"github.com/anz-bank/sysl/pkg/msg"
	"github.com/anz-bank/sysl/pkg/parse"
	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/anz-bank/sysl/pkg/validate"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// Node can be string or node
type Node []interface{}

type CodeGenOutput struct {
	filename string
	output   Node
}

/*
params[0] - modelAppName
params[1] - transformAppName
params[2] - depPath
params[3] - viewName
params[4] - basePath
*/
// applyTranformToModel loads applies the transform to input model
func applyTranformToModel(model, transform *sysl.Module, params ...string) (*sysl.Value, error) {
	modelApp, has := model.Apps[params[0]]
	if !has {
		var apps []string
		for k := range model.Apps {
			apps = append(apps, k)
		}
		sort.Strings(apps)
		return nil, errors.Errorf("app %s does not exist in model, available Apps: [%s]", params[0], strings.Join(apps, ", "))
	}
	view := transform.Apps[params[1]].Views[params[3]]
	if view == nil {
		return nil, errors.Errorf("Cannot execute missing view: %s, in app %s", params[3], params[1])
	}
	s := eval.Scope{}
	s.AddApp("app", modelApp)
	s.AddModule("module", model)
	s.AddString("depPath", params[2])
	s["basePath"] = eval.MakeValueString(params[4])
	var result *sysl.Value

	if perTypeTransform(view.Param) {
		result = eval.MakeValueList()
		var tNames []string
		for tName := range modelApp.Types {
			tNames = append(tNames, tName)
		}
		sort.Strings(tNames)
		for _, tName := range tNames {
			t := modelApp.Types[tName]
			s["typeName"] = eval.MakeValueString(tName)
			s["type"] = eval.TypeToValue(t)
			eval.AppendItemToValueList(result.GetList(), eval.EvaluateView(transform, params[1], params[3], s))
		}
	} else {
		result = eval.EvaluateView(transform, params[1], params[3], s)
	}

	return result, nil
}

func perTypeTransform(params []*sysl.Param) bool {
	paramMap := make(map[string]struct{})

	for _, p := range params {
		paramMap[p.Name] = struct{}{}
	}

	if _, has := paramMap["app"]; has {
		if _, has := paramMap["type"]; has {
			return true
		}
	} else {
		panic("Expecting at least an app <: sysl.App")
	}
	return false
}

// Serialize serializes node to string
func Serialize(w io.Writer, delim string, node Node) error {
	for _, n := range node {
		switch x := n.(type) {
		case string:
			if _, err := io.WriteString(w, x+delim); err != nil {
				return err
			}
		case Node:
			if err := Serialize(w, delim, x); err != nil {
				return err
			}
		}
	}
	return nil
}

// GenerateCode transform input sysl model to code in the target language described by
// grammar and a sysl transform
func GenerateCode(
	codegenParams *cmdutils.CmdContextParamCodegen,
	model *sysl.Module, modelAppName string,
	fs afero.Fs, logger *logrus.Logger, parserType parse.ParserType) ([]*CodeGenOutput, error) {
	var codeOutput []*CodeGenOutput
	depPath := codegenParams.DepPath
	basePath := codegenParams.BasePath

	logger.Debugf("root-transform: %s\n", codegenParams.RootTransform)
	logger.Debugf("transform: %s\n", codegenParams.Transform)
	logger.Debugf("dep-path: %s\n", codegenParams.DepPath)
	logger.Debugf("grammar: %s\n", codegenParams.Grammar)
	logger.Debugf("start: %s\n", codegenParams.Start)
	logger.Debugf("basePath: %s\n", codegenParams.BasePath)

	var transformFs afero.Fs
	transformFs = syslutil.NewChrootFs(fs, codegenParams.RootTransform)
	if mod.SyslModules {
		transformFs = mod.NewFs(transformFs)
	}
	tfmParser := parse.NewParserWithParserType(parserType)
	tx, transformAppName, err := parse.LoadAndGetDefaultApp(codegenParams.Transform, transformFs, tfmParser)
	if err != nil {
		return nil, err
	}

	if mod.SyslModules {
		fs = mod.NewFs(fs)
	}
	g, err := ebnfparser.ReadGrammar(fs, codegenParams.Grammar, codegenParams.Start)
	if err != nil {
		return nil, err
	}

	if !codegenParams.DisableValidator {
		grammarSysl, err := validate.LoadGrammarWithParserType(codegenParams.Grammar, fs, parserType)
		if err != nil {
			msg.NewMsg(msg.WarnValidationSkipped, []string{err.Error()}).LogMsg()
		} else {
			validator := validate.NewValidator(grammarSysl, tx.GetApps()[transformAppName], tfmParser)
			validator.Validate(codegenParams.Start, codegenParams.DepPath, codegenParams.BasePath)
			validator.LogMessages()
		}
	}

	fileNames, err := applyTranformToModel(model, tx, modelAppName, transformAppName,
		depPath, "filename", basePath)
	if err != nil {
		return nil, err
	}
	result, err := applyTranformToModel(model, tx, modelAppName, transformAppName,
		depPath, g.Start, basePath)
	if err != nil {
		return nil, err
	}
	switch {
	case fileNames.GetMap() != nil:
		filename := fileNames.GetMap().Items["filename"].GetS()
		logger.Println(filename)

		if result.GetMap() != nil {
			codeOutput = appendCodeOutput(g, result, logger, codeOutput, filename)
		} else if result.GetList() != nil {
			for _, v := range result.GetList().Value {
				codeOutput = appendCodeOutput(g, v, logger, codeOutput, filename)
			}
		}
	case fileNames.GetList() != nil && result.GetList() != nil:
		fileValues := fileNames.GetList().Value
		for i, v := range result.GetList().Value {
			filename := fileValues[i].GetMap().Items["filename"].GetS()
			codeOutput = appendCodeOutput(g, v, logger, codeOutput, filename)
		}
	default:
		panic("Unexpected combination for filenames and transformation results")
	}

	return codeOutput, nil
}

func appendCodeOutput(g *ebnfparser.EbnfGrammar, v *sysl.Value,
	logger *logrus.Logger, codeOutput []*CodeGenOutput, filename string) []*CodeGenOutput {
	r, err := ebnfparser.GenerateOutput(g, v, logger)
	if err != nil {
		return nil
	}
	codeOutput = append(codeOutput, &CodeGenOutput{filename, Node{r}})
	return codeOutput
}

func outputToFiles(output []*CodeGenOutput, fs afero.Fs) error {
	for _, o := range output {
		f, err := fs.Create(o.filename)
		if err != nil {
			return errors.Wrapf(err, "unable to create %q", o.filename)
		}
		logrus.Infoln("Writing file: " + f.Name())
		if err := Serialize(f, " ", o.output); err != nil {
			return errors.Wrapf(err, "error writing to %q", o.filename)
		}
		if err := f.Close(); err != nil {
			return errors.Wrapf(err, "error closing %q", o.filename)
		}
	}
	return nil
}
