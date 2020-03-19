package integrationdiagram

import (
	"regexp"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/sirupsen/logrus"
)

func GenerateIntegrations(intgenParams *cmdutils.CmdContextParamIntgen,
	model *sysl.Module,
	logger *logrus.Logger) (map[string]string, error) {
	r := make(map[string]string)

	logger.Debugf("project: %v\n", intgenParams.Project)
	logger.Debugf("clustered: %t\n", intgenParams.Clustered)
	logger.Debugf("exclude: %s\n", intgenParams.Exclude)
	logger.Debugf("epa: %t\n", intgenParams.EPA)
	logger.Debugf("title: %s\n", intgenParams.Title)
	logger.Debugf("filter: %s\n", intgenParams.Filter)
	logger.Debugf("output: %s\n", intgenParams.Output)

	if len(intgenParams.Exclude) == 0 && intgenParams.Project != "" {
		intgenParams.Exclude = []string{intgenParams.Project}
	}
	excludeStrSet := syslutil.MakeStrSet(intgenParams.Exclude...)

	// The "project" app that specifies the required view of the integration
	app := model.GetApps()[intgenParams.Project]
	of := cmdutils.MakeFormatParser(intgenParams.Output)
	// Iterate over each endpoint within the selected project
	for epname, endpt := range app.GetEndpoints() {
		outputDir := of.FmtOutput(intgenParams.Project, epname, endpt.GetLongName(), endpt.GetAttrs())
		if intgenParams.Filter != "" {
			re := regexp.MustCompile(intgenParams.Filter)
			if !re.MatchString(outputDir) {
				continue
			}
		}
		excludes := syslutil.MakeStrSetFromAttr("exclude", endpt.GetAttrs())
		passthroughs := syslutil.MakeStrSetFromAttr("passthrough", endpt.GetAttrs())
		b := MakeBuilderfromStmt(model, endpt.GetStmt(), excludeStrSet.Union(excludes), passthroughs)
		intsParam := &IntsParam{b.FinalApps, b.SeedAppsMap, b.DepsOut, app, endpt}
		args := &Args{intgenParams.Title, intgenParams.Project, intgenParams.Clustered, intgenParams.EPA}
		r[outputDir] = GenerateView(args, intsParam, model)
	}

	return r, nil
}
