package main

import (
	"os"
	"testing"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	"github.com/golang/protobuf/jsonpb"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type mockVisitor struct {
	mock.Mock
}

func (m *mockVisitor) Visit(e Element) error {
	m.Called(e)
	return nil
}

type elementSuite struct {
	suite.Suite
	e Element
	v *mockVisitor
}

func (suite *elementSuite) SetupTest() {
	suite.v = new(mockVisitor)
	suite.v.On("Visit", mock.Anything)
}

func (suite *elementSuite) TestElementAccept() {
	t := suite.T()
	require.NoError(t, suite.e.Accept(suite.v))

	suite.v.AssertCalled(t, "Visit", suite.e)
	suite.v.AssertNumberOfCalls(t, "Visit", 1)
}

func TestElementSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, &elementSuite{e: &EndpointCollectionElement{}})
	suite.Run(t, &elementSuite{e: &EndpointElement{}})
	suite.Run(t, &elementSuite{e: &StatementElement{}})
}

type mockVarManager struct {
	mock.Mock
}

func (m *mockVarManager) UniqueVarForAppName(appName string) string {
	args := m.Called(appName)

	return args.String(0)
}

func TestMakeEntry(t *testing.T) {
	t.Parallel()

	entry := makeEntry("a <- b [upto b <- c]")

	assert.NotNil(t, entry)
	assert.Equal(t, "a", entry.appName)
	assert.Equal(t, "b", entry.endpointName)
	assert.Equal(t, "b <- c", entry.upto)
}

func TestMakeEndpointCollectionElement(t *testing.T) {
	t.Parallel()

	e := MakeEndpointCollectionElement("title",
		[]string{"a <- b [upto b <- c]"},
		map[string]*Upto{"": {Comment: ""},
			"b <- c": {Comment: ""},
			"c <- d": {Comment: "test"}})
	assert.NotNil(t, e)
	assert.Equal(t, "title", e.title)
}

func TestEndpointElementSender(t *testing.T) {
	t.Parallel()

	m := new(mockVarManager)
	e := &EndpointElement{}
	m.On("UniqueVarForAppName", mock.Anything).Return("a")

	s := e.sender(m)

	m.AssertNumberOfCalls(t, "UniqueVarForAppName", 0)
	assert.Equal(t, "[", s)
}

func TestEndpointElementSenderWith(t *testing.T) {
	t.Parallel()

	m := new(mockVarManager)
	e := &EndpointElement{
		fromApp: &sysl.AppName{
			Part: []string{"test"},
		},
	}
	m.On("UniqueVarForAppName", mock.Anything).Return("a")

	s := e.sender(m)

	m.AssertCalled(t, "UniqueVarForAppName", "test")
	assert.Equal(t, "a", s)
}

func TestEndpointElementAgent(t *testing.T) {
	t.Parallel()

	m := new(mockVarManager)
	e := &EndpointElement{
		appName: "test",
	}
	m.On("UniqueVarForAppName", mock.Anything).Return("a")

	s := e.agent(m)

	m.AssertCalled(t, "UniqueVarForAppName", "test")
	assert.Equal(t, "a", s)
}

func TestEndpointElementApplication(t *testing.T) {
	t.Parallel()

	m := &sysl.Module{
		Apps: map[string]*sysl.Application{
			"test": {},
		},
	}
	e := &EndpointElement{
		appName: "test",
	}

	s := e.application(m)

	assert.NotNil(t, s)
}

func TestEndpointElementApplicationPanic(t *testing.T) {
	t.Parallel()

	m := &sysl.Module{}
	e := &EndpointElement{
		appName: "test",
	}

	assert.Panics(t, func() { e.application(m) })
}

func TestEndpointElementEndpoint(t *testing.T) {
	t.Parallel()

	m := &sysl.Application{
		Endpoints: map[string]*sysl.Endpoint{
			"test": {},
		},
	}
	e := &EndpointElement{
		endpointName: "test",
	}

	s := e.endpoint(m)

	assert.NotNil(t, s)
}

func TestEndpointElementEndpointPanic(t *testing.T) {
	t.Parallel()

	m := &sysl.Application{}
	e := &EndpointElement{
		endpointName: "test",
	}

	assert.Panics(t, func() { e.endpoint(m) })
}

type mockEndpointLabeler struct {
	mock.Mock
}

func (m *mockEndpointLabeler) LabelEndpoint(p *EndpointLabelerParam) string {
	args := m.Called(p)

	return args.String(0)
}

func TestEndpointElementEndpointLabel(t *testing.T) {
	t.Parallel()

	// Given
	l := new(mockEndpointLabeler)
	m := &sysl.Module{}
	e := &EndpointElement{
		endpointName: "a -> b",
	}

	// When
	actual := e.label(l, m, &sysl.Endpoint{}, syslutil.MakeStrSet(), false, false, false)

	// Then
	assert.Equal(t, " ⬄ b", actual)
}

func TestEndpointElementEndpointLabelWithValidStmt(t *testing.T) {
	t.Parallel()

	// Given
	l := new(mockEndpointLabeler)
	m := &sysl.Module{}
	e := &EndpointElement{
		endpointName: "a -> b",
		stmt: &sysl.Statement{
			Stmt: &sysl.Statement_Call{
				Call: &sysl.Call{},
			},
		},
		senderEndpointPatterns: syslutil.MakeStrSet(),
	}
	l.On("LabelEndpoint", mock.Anything).Return("test")

	// When
	actual := e.label(l, m, &sysl.Endpoint{}, syslutil.MakeStrSet("a"), false, true, false)

	// Then
	l.AssertNumberOfCalls(t, "LabelEndpoint", 1)
	assert.Equal(t, "test", actual)
}

func TestEndpointElementEndpointLabelWithValidStmtAndEmptyPatterns(t *testing.T) {
	t.Parallel()

	// Given
	l := new(mockEndpointLabeler)
	m := &sysl.Module{}
	e := &EndpointElement{
		endpointName: "a -> b",
		stmt: &sysl.Statement{
			Stmt: &sysl.Statement_Call{
				Call: &sysl.Call{},
			},
		},
		senderEndpointPatterns: syslutil.MakeStrSet(),
	}
	l.On("LabelEndpoint", mock.Anything).Return("test")

	// When
	actual := e.label(l, m, &sysl.Endpoint{}, syslutil.MakeStrSet(), false, true, false)

	// Then
	l.AssertNumberOfCalls(t, "LabelEndpoint", 1)
	assert.Equal(t, "test", actual)
}

func TestStatementElementIsLastStmt(t *testing.T) {
	t.Parallel()

	e := &StatementElement{
		isLastParentStmt: true,
		stmts:            []*sysl.Statement{{}},
	}
	assert.True(t, e.isLastStmt(0))
}

func readModule(p string) (*sysl.Module, error) {
	m := &sysl.Module{}
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	if err := jsonpb.Unmarshal(f, m); err != nil {
		return nil, err
	}
	return m, nil
}

type labeler struct{}

func (l *labeler) LabelApp(appName, controls string, attrs map[string]*sysl.Attribute) string {
	return appName
}

func (l *labeler) LabelEndpoint(p *EndpointLabelerParam) string {
	return p.EndpointName
}

func TestSequenceDiagramVisitorVisit(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	l := &labeler{}
	w := MakeSequenceDiagramWriter(true, "skinparam maxMessageSize 250")
	m, err := readModule("./tests/sequence_diagram_project.golden.json")
	require.NoError(t, err)
	v := MakeSequenceDiagramVisitor(l, l, w, m, "appname", "", logger)
	e := MakeEndpointCollectionElement("Profile", []string{"WebFrontend <- RequestProfile"}, map[string]*Upto{
		"Frontend <- Profile": {
			ValueType:  BBEndpointCollection,
			Comment:    "see below",
			VisitCount: 0,
		},
		"ApplicationFrontend <- AppProfile": {
			ValueType:  BBApplication,
			Comment:    "see below",
			VisitCount: 0,
		},
		"Commandline <- CommandlineApp": {
			ValueType:  BBCommandLine,
			Comment:    "see below",
			VisitCount: 0,
		},
	})

	require.NoError(t, e.Accept(v))

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

	assert.Equal(t, expected, w.String())
}

func TestSequenceDiagramToFormatNameAttributesVisitorVisit(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	al := MakeFormatParser(`%(@status?<color red>%(appname)</color>|%(appname))`)
	el := MakeFormatParser(`%(@status? <color green>%(epname)</color>|%(epname))`)
	w := MakeSequenceDiagramWriter(true, "skinparam maxMessageSize 250")
	m, err := readModule("./tests/sequence_diagram_name_format.golden.json")
	require.NoError(t, err)
	v := MakeSequenceDiagramVisitor(al, el, w, m, "appname", "", logger)
	e := MakeEndpointCollectionElement("Diagram", []string{"User <- Check Balance"}, map[string]*Upto{})

	require.NoError(t, e.Accept(v))

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
title Diagram
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

	assert.Equal(t, expected, w.String())
}

func TestVisitStatement(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	l := &labeler{}
	w := MakeSequenceDiagramWriter(true, "skinparam maxMessageSize 250")
	m, err := readModule("./tests/sequence_diagram_project.golden.json")
	require.NoError(t, err)
	v := MakeSequenceDiagramVisitor(l, l, w, m, "appname", "", logger)
	stmt := []*sysl.Statement{
		{Stmt: &sysl.Statement_Loop{Loop: &sysl.Loop{Stmt: []*sysl.Statement{}}}},
		{Stmt: &sysl.Statement_LoopN{LoopN: &sysl.LoopN{Stmt: []*sysl.Statement{}}}},
		{Stmt: &sysl.Statement_Foreach{Foreach: &sysl.Foreach{Stmt: []*sysl.Statement{}}}},
		{Stmt: &sysl.Statement_Alt{Alt: &sysl.Alt{Choice: []*sysl.Alt_Choice{{Stmt: []*sysl.Statement{}}}}}},
		{Stmt: &sysl.Statement_Cond{Cond: &sysl.Cond{Stmt: []*sysl.Statement{}}}},
		{Stmt: &sysl.Statement_Group{Group: &sysl.Group{Stmt: []*sysl.Statement{}}}},
	}
	e := &StatementElement{
		stmts: stmt,
	}
	error := v.visitStatment(e)
	assert.Nil(t, error)
}
