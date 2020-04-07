package mermaidsequencediagram

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/spf13/afero"
)

const previousapp = "..."

func TestBadInputsToGenerateMermaidSequenceDiagram(t *testing.T) {
	t.Parallel()
	appname := "wrongname"
	epname := "wrongep"
	m, err := parse.NewParser().Parse("demo/simple/sysl-sd.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateMermaidSequenceDiagram(m, appname, epname, previousapp, 1, true)
	assert.NotNil(t, m)
	assert.Empty(t, r)
	assert.Error(t, err)
}

func TestGenerateMermaidSequenceDiagram(t *testing.T) {
	t.Parallel()
	appname := "WebFrontend"
	epname := "RequestProfile"
	m, err := parse.NewParser().Parse("demo/simple/sysl-sd.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateMermaidSequenceDiagram(m, appname, epname, previousapp, 1, true)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.Nil(t, err)
	expected := `%% AUTOGENERATED CODE -- DO NOT EDIT!

sequenceDiagram
 ...->>WebFrontend: RequestProfile
 WebFrontend->>+Api: GET /users/{user_id}/profile
 Api->>+Database: QueryUser
 Database-->>-Api: User
 Api-->>-WebFrontend: UserProfile
 WebFrontend-->>...: Profile Page
`
	assert.Equal(t, expected, r)
}

func TestGenerateMermaidSequenceDiagram2(t *testing.T) {
	t.Parallel()
	appname := "WebFrontend"
	epname := "RequestProfile"
	m, err := parse.NewParser().Parse("demo/simple/sysl-sd2.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateMermaidSequenceDiagram(m, appname, epname, previousapp, 1, true)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.Nil(t, err)
	expected := `%% AUTOGENERATED CODE -- DO NOT EDIT!

sequenceDiagram
 ...->>WebFrontend: RequestProfile
 WebFrontend->>+Api: GET /users/{user_id}/profile
 Api->>+Database: QueryUser
 Database-->>-Api: User [~y, x="1"]
 Api-->>-WebFrontend: UserProfile
 WebFrontend->>+WebFrontend: FooBar
 WebFrontend-->>...: Profile Page
`
	assert.Equal(t, expected, r)
}

func TestGenerateMermaidSequenceDiagram3(t *testing.T) {
	t.Parallel()
	appname := "MobileApp"
	epname := "Login"
	m, err := parse.NewParser().Parse("demo/simple/sysl-app-hyperlink.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateMermaidSequenceDiagram(m, appname, epname, previousapp, 1, true)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.Nil(t, err)
	expected := `%% AUTOGENERATED CODE -- DO NOT EDIT!

sequenceDiagram
 ...->>MobileApp: Login
 MobileApp->>+Server: LoginRequest
 Server-->>-MobileApp: MobileApp.LoginResponse
`
	assert.Equal(t, expected, r)
}

func TestGenerateMermaidSequenceDiagramWithIfElseLoopActionAndGroupStatements(t *testing.T) {
	t.Parallel()
	appname := "BatEater"
	epname := "EatBat"
	m, err := parse.NewParser().Parse("demo/simple/best-ever-sysl-example.sysl",
		syslutil.NewChrootFs(afero.NewOsFs(), projDir))
	if err != nil {
		t.Error(err)
	}
	r, err := GenerateMermaidSequenceDiagram(m, appname, epname, previousapp, 1, true)
	assert.NotNil(t, m)
	assert.NotNil(t, r)
	assert.Nil(t, err)
	expected := `%% AUTOGENERATED CODE -- DO NOT EDIT!

sequenceDiagram
 ...->>BatEater: EatBat
 BatEater->>+Actions: EatBat
 Actions-->>-BatEater: Coronavirus
 alt BatEater got Coronavirus
  BatEater->>+Actions: SpreadCoronavirus
  Actions->>+TheWorld: WTFBro
  TheWorld->>TheWorld: DoNothing
  Actions->>+ChineseGovernment: FirstResponse
  ChineseGovernment->>ChineseGovernment: DenyEverything
  loop until futhur notice
   ChineseGovernment->>+Actions: SilenceJournalists
   Actions->>Actions: ...
   ChineseGovernment->>+Actions: BuildHospitalsVeryQuickly
   Actions->>Actions: HowdYouDoThat
   Actions-->>-ChineseGovernment: Hospital
  end
  ChineseGovernment-->>-Actions: Nothing
  Actions->>+ChineseGovernment: SecondResponse
  alt Coronavirus still exists
   ChineseGovernment->>ChineseGovernment: Lockdown
  end
  ChineseGovernment-->>-Actions: Insult
  Actions->>+TheWorld: WereFucked
  TheWorld->>+Actions: Giveup
  Actions->>Actions: PlayVideoGames
  Actions->>Actions: EatJunkFood
  Actions->>Actions: GameOver
  Actions-->>-BatEater: Coronavirus
 else if sdsd
  BatEater->>BatEater: norhinf
 else
  BatEater->>BatEater: go back to normal
 end
 BatEater-->>...: Coronavirus
`
	assert.Equal(t, expected, r)
}