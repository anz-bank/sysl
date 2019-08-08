package main

import (
	"flag"
	"regexp"
	"strings"
	"testing"

	"github.com/anz-bank/sysl/sysl2/sysl/parse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateSequenceDiag(t *testing.T) {
	m, _ := parse.Parse("demo/simple/sysl-sd.sysl", "../../")
	l := &labeler{}
	p := &sequenceDiagParam{}
	p.endpoints = []string{"WebFrontend <- RequestProfile"}
	p.AppLabeler = l
	p.EndpointLabeler = l
	p.title = "Profile"
	r, err := generateSequenceDiag(m, p)

	expected := `''''''''''''''''''''''''''''''''''''''''''
''                                      ''
''  AUTOGENERATED CODE -- DO NOT EDIT!  ''
''                                      ''
''''''''''''''''''''''''''''''''''''''''''

@startuml
control "WebFrontend" as _0
control "Api" as _1
database "Database" as _2
skinparam maxMessageSize 250
title Profile
== WebFrontend <- RequestProfile ==
[->_0 : RequestProfile
activate _0
 _0->_1 : GET /users/{user_id}/profile
 activate _1
  _1->_2 : QueryUser
  activate _2
  _1<--_2 : User
  deactivate _2
 _0<--_1 : UserProfile
 deactivate _1
[<--_0 : Profile Page
deactivate _0
@enduml
`

	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.Nil(t, err)
	assert.Equal(t, expected, r)
}

func TestGenerateSequenceDiagramsToFormatNameAttributes(t *testing.T) {
	m, _ := parse.Parse("sequence_diagram_name_format.sysl", "./tests/")
	al := MakeFormatParser(`%(@status?<color red>%(appname)</color>|%(appname))`)
	el := MakeFormatParser(`%(@status? <color green>%(epname)</color>|%(epname))`)
	p := &sequenceDiagParam{}
	p.endpoints = []string{"User <- Check Balance"}
	p.AppLabeler = al
	p.EndpointLabeler = el
	r, err := generateSequenceDiag(m, p)
	p.title = "Diagram"
	expected := `''''''''''''''''''''''''''''''''''''''''''
''                                      ''
''  AUTOGENERATED CODE -- DO NOT EDIT!  ''
''                                      ''
''''''''''''''''''''''''''''''''''''''''''

@startuml
actor "User" as _0
boundary "MobileApp" as _1
control "<color red>Server</color>" as _2
database "DB" as _3
skinparam maxMessageSize 250
== User <- Check Balance ==
 _0->_1 : Login
 activate _1
  _1->_2 : Login
  activate _2
  _2 -> _2 : do input validation
   _2->_3 :  <color green>Save</color>
  _1<--_2 : success or failure
  deactivate _2
 deactivate _1
 _0->_1 : Check Balance
 activate _1
  _1->_2 : Read User Balance
  activate _2
   _2->_3 :  <color green>Load</color>
  _1<--_2 : balance
  deactivate _2
 deactivate _1
@enduml
`

	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.Nil(t, err)
	assert.Equal(t, expected, r)
}

func TestGenerateSequenceDiagramsToFormatComplexAttributes(t *testing.T) {
	m, _ := parse.Parse("sequence_diagram_name_format.sysl", "./tests/")
	al := MakeFormatParser(`%(@status?<color red>%(appname)</color>|%(appname))`)
	el := MakeFormatParser(`%(@status? <color green>%(epname)</color>|%(epname))`)
	p := &sequenceDiagParam{}
	p.endpoints = []string{"User <- Check Balance"}
	p.AppLabeler = al
	p.EndpointLabeler = el
	r, err := generateSequenceDiag(m, p)
	p.title = "Diagram"
	expected := `''''''''''''''''''''''''''''''''''''''''''
''                                      ''
''  AUTOGENERATED CODE -- DO NOT EDIT!  ''
''                                      ''
''''''''''''''''''''''''''''''''''''''''''

@startuml
actor "User" as _0
boundary "MobileApp" as _1
control "<color red>Server</color>" as _2
database "DB" as _3
skinparam maxMessageSize 250
== User <- Check Balance ==
 _0->_1 : Login
 activate _1
  _1->_2 : Login
  activate _2
  _2 -> _2 : do input validation
   _2->_3 :  <color green>Save</color>
  _1<--_2 : success or failure
  deactivate _2
 deactivate _1
 _0->_1 : Check Balance
 activate _1
  _1->_2 : Read User Balance
  activate _2
   _2->_3 :  <color green>Load</color>
  _1<--_2 : balance
  deactivate _2
 deactivate _1
@enduml
`

	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.Nil(t, err)
	assert.Equal(t, expected, r)
}

type loadAppArgs struct {
	root   string
	models string
}

func TestLoadAppReturnError(t *testing.T) {
	test := loadAppArgs{
		"../../demo/simple/", "",
	}
	_, err := loadApp(test.root, test.models)
	assert.Error(t, err)
}

func TestLoadApp(t *testing.T) {
	test := loadAppArgs{
		"./tests/", "sequence_diagram_test.sysl",
	}
	mod, err := loadApp(test.root, test.models)
	require.NoError(t, err)
	assert.NotNil(t, mod)
	apps := mod.GetApps()
	app := apps["Database"]

	assert.Equal(t, []string{"Database"}, app.GetName().GetPart())

	appPatternsAttr := app.GetAttrs()["patterns"].GetA().GetElt()
	patterns := make([]string, 0, len(appPatternsAttr))
	for _, val := range appPatternsAttr {
		patterns = append(patterns, val.GetS())
	}
	assert.Equal(t, []string{"db"}, patterns)

	queryUserParams := app.GetEndpoints()["QueryUser"].GetParam()
	params := make([]string, 0, len(queryUserParams))
	for _, val := range queryUserParams {
		params = append(params, val.GetName())
	}
	assert.Equal(t, []string{"user_id"}, params)
}

type sdArgs struct {
	rootModel      string
	endpointFormat string
	appFormat      string
	title          string
	output         string
	endpoints      []string
	apps           []string
	modules        string
	blackboxes     [][]string
	groupbox       string
}

func TestDoConstructSequenceDiagramsNoSyslSdFiltersWithoutEndpoints(t *testing.T) {
	// Given
	args := &sdArgs{
		rootModel: "./tests/",
		modules:   "sequence_diagram_test.sysl",
	}
	expected := make(map[string]string)

	// When
	result, err := DoConstructSequenceDiagrams(args.rootModel, args.endpointFormat, args.appFormat,
		args.title, args.output, args.modules, args.endpoints, args.apps, args.blackboxes, args.groupbox)
	require.NoError(t, err)

	// Then
	assert.Equal(t, expected, result)
}

func TestDoConstructSequenceDiagramsMissingFile(t *testing.T) {
	// Given
	args := &sdArgs{
		rootModel: "./tests/",
		modules:   "MISSING_FILE.sysl",
	}

	// When
	_, err := DoConstructSequenceDiagrams(args.rootModel, args.endpointFormat, args.appFormat,
		args.title, args.output, args.modules, args.endpoints, args.apps, args.blackboxes, args.groupbox)
	assert.Error(t, err)
}

func TestDoConstructSequenceDiagramsNoSyslSdFilters(t *testing.T) {
	// Given
	args := &sdArgs{
		rootModel: "./tests/",
		modules:   "sequence_diagram_test.sysl",
		endpoints: []string{"QueryUser"},
		output:    "_.png",
	}

	// When
	assert.Panics(t, func() {
		_, err := DoConstructSequenceDiagrams(args.rootModel, args.endpointFormat, args.appFormat,
			args.title, args.output, args.modules, args.endpoints, args.apps, args.blackboxes, args.groupbox)
		assert.NoError(t, err)
	})
}

func TestDoConstructSequenceDiagrams(t *testing.T) {
	// Given
	args := &sdArgs{
		rootModel: "./tests/",
		modules:   "sequence_diagram_project.sysl",
		output:    "%(epname).png",
		apps:      []string{"Project"},
	}
	expectContent := `''''''''''''''''''''''''''''''''''''''''''
''                                      ''
''  AUTOGENERATED CODE -- DO NOT EDIT!  ''
''                                      ''
''''''''''''''''''''''''''''''''''''''''''

@startuml
control "" as _0
control "" as _1
database "" as _2
skinparam maxMessageSize 250
title Profile
== WebFrontend <- RequestProfile ==
[->_0 : RequestProfile
activate _0
 _0->_1 :` + " " + `
 activate _1
  _1->_2 :` + " " + `
  activate _2
  _1<--_2 : User
  deactivate _2
 _0<--_1 : UserProfile
 deactivate _1
[<--_0 : Profile Page
deactivate _0
@enduml
`
	expected := map[string]string{
		"_.png": expectContent,
	}

	// When
	result, err := DoConstructSequenceDiagrams(args.rootModel, args.endpointFormat, args.appFormat,
		args.title, args.output, args.modules, args.endpoints, args.apps, args.blackboxes, args.groupbox)
	require.NoError(t, err)

	// Then
	assert.Equal(t, expected, result)
}

func TestDoConstructSequenceDiagramWithBlackbox(t *testing.T) {
	// Given
	args := &sdArgs{
		rootModel:  "./tests/",
		modules:    "call.sysl",
		output:     "tests/call.png",
		endpoints:  []string{"MobileApp <- Login"},
		blackboxes: [][]string{{"Server <- Login", "call to database"}},
	}

	// When
	expectContent := `''''''''''''''''''''''''''''''''''''''''''
''                                      ''
''  AUTOGENERATED CODE -- DO NOT EDIT!  ''
''                                      ''
''''''''''''''''''''''''''''''''''''''''''

@startuml
control "" as _0
control "" as _1
control "" as _2
skinparam maxMessageSize 250
== MobileApp <- Login ==
[->_0 : Login
activate _0
 _0->_1 :` + " " + `
 activate _1
  _1->_2 :` + " " + `
  activate _2
  note over _2: call to database
  _1<--_2 : <color blue>Server.LoginResponse</color> <<color green>?, ?</color>>
  deactivate _2
 _0<--_1 : <color blue>APIGateway.LoginResponse</color> <<color green>?, ?</color>>
 deactivate _1
deactivate _0
@enduml
`
	result, err := DoConstructSequenceDiagrams(args.rootModel, args.endpointFormat, args.appFormat,
		args.title, args.output, args.modules, args.endpoints, args.apps, args.blackboxes, args.groupbox)
	require.NoError(t, err)
	expected := map[string]string{"tests/call.png": expectContent}
	// Then
	assert.Equal(t, expected, result)
}

func TestDoConstructSequenceDiagramsToFormatComplexName(t *testing.T) {
	// Given
	args := &sdArgs{
		rootModel: "./tests/",
		modules:   "sequence_diagram_complex_format.sysl",
		output:    "%(epname).png",
		apps:      []string{"Project"},
	}
	expectContent := `''''''''''''''''''''''''''''''''''''''''''
''                                      ''
''  AUTOGENERATED CODE -- DO NOT EDIT!  ''
''                                      ''
''''''''''''''''''''''''''''''''''''''''''

@startuml
control "//te//\n<color grey>Ex e</color>\n**User**" as _0
control "**MobileApp**" as _1
skinparam maxMessageSize 250
title Diagram
== User <- Check Balance ==
[->_0 : Check Balance
activate _0
 _0->_1 : //«hello»//** <color red>pat?</color>**aa Login
 deactivate _0
@enduml
`
	expected := map[string]string{
		"Seq.png": expectContent,
	}

	// When
	result, err := DoConstructSequenceDiagrams(args.rootModel, args.endpointFormat, args.appFormat,
		args.title, args.output, args.modules, args.endpoints, args.apps, args.blackboxes, args.groupbox)
	require.NoError(t, err)

	// Then
	assert.Equal(t, expected, result)
}

func TestDoGenerateSequenceDiagrams(t *testing.T) {
	type args struct {
		flags *flag.FlagSet
		args  []string
	}
	argsData := []string{"sd"}
	tests := []struct {
		name string
		args args
	}{
		{
			"Case-Do generate sequence diagrams",
			args{
				flag.NewFlagSet(argsData[0], flag.PanicOnError),
				argsData,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Figure out why this is broken.
			// require.NoError(t, DoGenerateSequenceDiagrams(tt.args.args))
		})
	}
}

func TestDoConstructSequenceDiagramWithGroupingCommandline(t *testing.T) {
	// Given
	args := &sdArgs{
		rootModel: "./tests/",
		modules:   "call.sysl",
		output:    "tests/call.png",
		endpoints: []string{"MobileApp <- Login"},
		groupbox:  "owner",
	}
	var boxPresent bool
	var err error

	// When
	boxServer := `box "server" #LightBlue
	participant _\d
	participant _\d
end box`
	boxClient := `box "client" #LightBlue
	participant _\d
	participant _\d
end box`
	result, err := DoConstructSequenceDiagrams(args.rootModel, args.endpointFormat, args.appFormat,
		args.title, args.output, args.modules, args.endpoints, args.apps, args.blackboxes, args.groupbox)
	require.NoError(t, err)

	// Then
	boxPresent, err = regexp.MatchString(boxServer, result["tests/call.png"])
	assert.Nil(t, err, "Error compiling regular expression")
	assert.True(t, boxPresent)
	boxPresent, err = regexp.MatchString(boxClient, result["tests/call.png"])
	assert.Nil(t, err, "Error compiling regular expression")
	assert.True(t, boxPresent)
	assert.Equal(t, 4, strings.Count(result["tests/call.png"], "participant"))
}

func TestDoConstructSequenceDiagramWithGroupingSysl(t *testing.T) {
	// Given
	args := &sdArgs{
		rootModel: "./tests/",
		modules:   "groupby.sysl",
		output:    "%(epname).png",
		endpoints: []string{"SEQ-One"},
		apps:      []string{"Project :: Sequences"},
	}
	var boxPresent bool
	var err error

	// When
	boxOnpremise := `box "onpremise" #LightBlue
	participant _\d
	participant _\d
end box`
	boxCloud := `box "cloud" #LightBlue
	participant _\d
	participant _\d
end box`
	result, err := DoConstructSequenceDiagrams(args.rootModel, args.endpointFormat, args.appFormat,
		args.title, args.output, args.modules, args.endpoints, args.apps, args.blackboxes, args.groupbox)
	require.NoError(t, err)

	// Then
	boxPresent, err = regexp.MatchString(boxOnpremise, result["SEQ-One.png"])
	assert.Nil(t, err, "Error compiling regular expression")
	assert.True(t, boxPresent)
	boxPresent, err = regexp.MatchString(boxCloud, result["SEQ-One.png"])
	assert.Nil(t, err, "Error compiling regular expression")
	assert.True(t, boxPresent)
}

func TestDoConstructSequenceDiagramWithOneEntityBox(t *testing.T) {
	// Given
	args := &sdArgs{
		rootModel: "./tests/",
		modules:   "groupby.sysl",
		output:    "%(epname).png",
		endpoints: []string{"SEQ-Two"},
		apps:      []string{"Project :: Sequences"},
		groupbox:  "location",
	}

	var boxPresent bool
	var err error

	// When
	boxCloud := `box "cloud" #LightBlue
	participant _\d
end box`
	result, err := DoConstructSequenceDiagrams(args.rootModel, args.endpointFormat, args.appFormat,
		args.title, args.output, args.modules, args.endpoints, args.apps, args.blackboxes, args.groupbox)
	require.NoError(t, err)

	// Then
	boxPresent, err = regexp.MatchString(boxCloud, result["SEQ-Two.png"])
	assert.Nil(t, err, "Error compiling regular expression")
	assert.True(t, boxPresent)
}
