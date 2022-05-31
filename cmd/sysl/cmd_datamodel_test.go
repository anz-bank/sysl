package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/anz-bank/sysl/pkg/diagrams"

	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/anz-bank/sysl/pkg/datamodeldiagram"
)

func TestDataModel(t *testing.T) {
	t.Parallel()
	checkPlantUML(t)

	p := &datamodelCmd{}
	p.Output = "whatever.svg"
	p.Project = "Project"          //nolint:goconst
	p.ClassFormat = "%(classname)" //nolint:goconst
	p.Direct = false
	fs := afero.NewOsFs()
	filename := filepath.Join(testDir, "multiple-app-datamodel.sysl")
	m, err := parse.NewParser().ParseFromFs(filename, fs)
	if err != nil {
		panic(err)
	}
	outmap, err := datamodeldiagram.GenerateDataModels(&p.CmdContextParamDatagen, m, logrus.New())
	assert.Nil(t, err, "Generating the data diagrams failed")
	err = p.GenerateFromMap(outmap, afero.NewMemMapFs())
	assert.Nil(t, err, "Generating the data diagrams failed")

	expected := map[string]string{"whatever.svg": `@startuml
''''''''''''''''''''''''''''''''''''''''''
''                                      ''
''  AUTOGENERATED CODE -- DO NOT EDIT!  ''
''                                      ''
''''''''''''''''''''''''''''''''''''''''''

class "AnotherApp.AnotherType" as _0 << (D,orchid) >> {
}
class "App.User" as _1 << (D,orchid) >> {
+ beep : **Server.ftgyhb**
+ weuiyfgwihe : **AnotherApp.AnotherType**
+ whatever : int
}
class "Server.ftgyhb" as _2 << (D,orchid) >> {
+ blah : int
}
_1 *-- "1..1 " _0
_1 *-- "1..1 " _2
@enduml
`}
	assert.Equal(t, outmap, expected)
}

func TestHTML(t *testing.T) {
	t.Parallel()
	checkPlantUML(t)

	p := &datamodelCmd{}
	p.Output = "whatever.html"
	p.Project = "Project"          //nolint:goconst
	p.ClassFormat = "%(classname)" //nolint:goconst
	p.Direct = false
	fs := afero.NewOsFs()
	plantuml := os.Getenv("SYSL_PLANTUML")
	filename := filepath.Join(testDir, "multiple-app-datamodel.sysl")
	m, err := parse.NewParser().ParseFromFs(filename, fs)
	if err != nil {
		panic(err)
	}
	fs = afero.NewMemMapFs()
	outmap, err := datamodeldiagram.GenerateDataModels(&p.CmdContextParamDatagen, m, logrus.New())
	assert.Nil(t, err, "Generating the data diagrams failed")
	err = p.GenerateFromMap(outmap, fs)
	assert.NoError(t, err)
	err = diagrams.OutputPlantuml(p.Output, plantuml, outmap["whatever.html"], fs)
	assert.NoError(t, err)
	file, err := fs.Open(p.Output)
	assert.NoError(t, err)
	html, err := ioutil.ReadAll(file)
	assert.NoError(t, err)
	expected := fmt.Sprintf(
		"<img src=\"%s/svg/UDgCaK5hmZ0SnU_v56-"+
			"zTAjhi1v5nBA4iOk5BPvBrB-cqDMGHAMCVVUn"+
			"LJh65FfE8VpUumV_XG_QXUDxpUB1ON6CGRcW-"+
			"KeLpt8fNtCb1PuA8P6c40MMXO8KB-gkHmUl3d"+
			"PbcrfxZoXl3i6GowtbbwTgBKNG7kKOindknUF"+
			"1RKorVS1yZW_ssJUjvIjFhcEpQ-m8QoABAPBa"+
			"ZTo97D-5VMlMIS96EDEnQdVxSsNeXxXkqg561"+
			"pgHmnHL4tuL_ens7fCR7hKsVRlCaAGfeepp31"+
			"7AyR-V2LjGi_q-_rS0003__ykTVD00\" alt=\"plantuml\">"+
			"\n", plantuml)
	assert.Equal(t, expected, string(html))
}

func TestLinkOutput(t *testing.T) {
	t.Parallel()
	checkPlantUML(t)

	p := &datamodelCmd{}
	p.Output = "whatever.html"
	p.Project = "Project"          //nolint:goconst
	p.ClassFormat = "%(classname)" //nolint:goconst
	p.Direct = false
	plantuml := os.Getenv("SYSL_PLANTUML")
	fs := afero.NewOsFs()
	filename := filepath.Join(testDir, "multiple-app-datamodel.sysl")
	m, err := parse.NewParser().ParseFromFs(filename, fs)
	if err != nil {
		panic(err)
	}
	fs = afero.NewMemMapFs()
	outmap, err := datamodeldiagram.GenerateDataModels(&p.CmdContextParamDatagen, m, logrus.New())
	assert.Nil(t, err, "Generating the data diagrams failed")
	err = p.GenerateFromMap(outmap, fs)
	assert.NoError(t, err)
	err = diagrams.OutputPlantuml(p.Output, plantuml, outmap["whatever.html"], fs)
	assert.NoError(t, err)
	file, err := fs.Open(p.Output)
	assert.NoError(t, err)
	link, err := ioutil.ReadAll(file)
	assert.NoError(t, err)
	expected := fmt.Sprintf(
		"<img src=\"%s/svg/UDgCaK5hmZ0SnU_v56-"+
			"zTAjhi1v5nBA4iOk5BPvBrB-cqDMGHAMCVVUn"+
			"LJh65FfE8VpUumV_XG_QXUDxpUB1ON6CGRcW-"+
			"KeLpt8fNtCb1PuA8P6c40MMXO8KB-gkHmUl3d"+
			"PbcrfxZoXl3i6GowtbbwTgBKNG7kKOindknUF"+
			"1RKorVS1yZW_ssJUjvIjFhcEpQ-m8QoABAPBa"+
			"ZTo97D-5VMlMIS96EDEnQdVxSsNeXxXkqg561"+
			"pgHmnHL4tuL_ens7fCR7hKsVRlCaAGfeepp31"+
			"7AyR-V2LjGi_q-_rS0003__ykTVD00\" alt=\"plantuml\">"+
			"\n", plantuml)
	assert.Equal(t, expected, string(link))
}

func TestSequence(t *testing.T) {
	t.Parallel()
	checkPlantUML(t)

	p := &datamodelCmd{}
	p.Output = "whatever.puml"
	p.Project = "sequence"
	p.ClassFormat = "%(classname)" //nolint:goconst
	p.Direct = false
	plantuml := os.Getenv("SYSL_PLANTUML")
	fs := afero.NewOsFs()
	filename := filepath.Join(testDir, "sequence.sysl")
	m, err := parse.NewParser().ParseFromFs(filename, fs)
	if err != nil {
		panic(err)
	}
	fs = afero.NewMemMapFs()
	outmap, err := datamodeldiagram.GenerateDataModels(&p.CmdContextParamDatagen, m, logrus.New())
	assert.Nil(t, err, "Generating the data diagrams failed")
	err = p.GenerateFromMap(outmap, fs)
	assert.NoError(t, err)
	err = diagrams.OutputPlantuml(p.Output, plantuml, outmap["whatever.puml"], fs)
	assert.NoError(t, err)
	file, err := fs.Open(p.Output)
	assert.NoError(t, err)
	link, err := ioutil.ReadAll(file)
	assert.NoError(t, err)
	expected := `@startuml
''''''''''''''''''''''''''''''''''''''''''
''                                      ''
''  AUTOGENERATED CODE -- DO NOT EDIT!  ''
''                                      ''
''''''''''''''''''''''''''''''''''''''''''

class "App.foo" as _0 << (D,orchid) >> {
+ content : **App2.bar**
}
class "App2.bar" as _1 << (D,orchid) >> {
+ content : **Sequence <bar2>**
}
class "App2.bar2" as _2 << (D,orchid) >> {
+ content : **Sequence <App3.ifhu>**
}
class "App3.ifhu" as _3 << (D,orchid) >> {
+ content : string
}
_0 *-- "1..1 " _1
_1 *-- "0..*" _2
_2 *-- "0..*" _3
@enduml

`
	assert.Equal(t, expected, string(link), string(link))
}
