package seqs

import (
	"testing"

	"github.com/anz-bank/sysl/src/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockVisitor struct {
	mock.Mock
}

func (m *MockVisitor) Visit(e Element) {
	m.Called(e)
}

type elementSuite struct {
	suite.Suite
	e Element
	v *MockVisitor
}

func (suite *elementSuite) SetupTest() {
	suite.v = new(MockVisitor)
	suite.v.On("Visit", mock.Anything)
}

func (suite *elementSuite) TestElementAccept() {
	t := suite.T()
	suite.e.Accept(suite.v)

	suite.v.AssertCalled(t, "Visit", suite.e)
	suite.v.AssertNumberOfCalls(t, "Visit", 1)
}

func TestElementSuite(t *testing.T) {
	suite.Run(t, &elementSuite{e: &EndpointCollectionElement{}})
	suite.Run(t, &elementSuite{e: &EndpointElement{}})
	suite.Run(t, &elementSuite{e: &StatementElement{}})
}

type MockVarManager struct {
	mock.Mock
}

func (m *MockVarManager) UniqueVarForAppName(appName string) string {
	args := m.Called(appName)

	return args.String(0)
}

func TestMakeEntry(t *testing.T) {
	entry := makeEntry("a <- b [upto b <- c]")

	assert.NotNil(t, entry)
	assert.Equal(t, "a", entry.appName)
	assert.Equal(t, "b", entry.endpointName)
	assert.Equal(t, "b <- c", entry.upto)
}

func TestMakeEndpointCollectionElement(t *testing.T) {
	e := MakeEndpointCollectionElement("title",
		[]string{"a <- b [upto b <- c]"},
		[][]string{
			{},
			{"b <- c"},
			{"c <- d", "test"},
		})

	assert.NotNil(t, e)
	assert.Equal(t, "title", e.title)
}

func TestEndpointElementSender(t *testing.T) {
	m := new(MockVarManager)
	e := &EndpointElement{}
	m.On("UniqueVarForAppName", mock.Anything).Return("a")

	s := e.sender(m)

	m.AssertNumberOfCalls(t, "UniqueVarForAppName", 0)
	assert.Equal(t, "[", s)
}

func TestEndpointElementSenderWith(t *testing.T) {
	m := new(MockVarManager)
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
	m := new(MockVarManager)
	e := &EndpointElement{
		appName: "test",
	}
	m.On("UniqueVarForAppName", mock.Anything).Return("a")

	s := e.agent(m)

	m.AssertCalled(t, "UniqueVarForAppName", "test")
	assert.Equal(t, "a", s)
}

func TestEndpointElementApplication(t *testing.T) {
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
	m := &sysl.Module{}
	e := &EndpointElement{
		appName: "test",
	}

	assert.Panics(t, func() { e.application(m) })
}

func TestEndpointElementEndpoint(t *testing.T) {
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
	m := &sysl.Application{}
	e := &EndpointElement{
		endpointName: "test",
	}

	assert.Panics(t, func() { e.endpoint(m) })
}
