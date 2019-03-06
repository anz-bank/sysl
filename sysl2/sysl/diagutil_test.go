package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

var testPlantumlInput = `
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

func TestDeflateAndEncode(t *testing.T) {
	//Given
	const expected = "UDfSaKbhmp0GXU_pAnwvYqY6NaniKkXoAgGRFUGW9l4qY7gh99SkzByN9GvnUfBGzmrwZw5bYEpZqDIqxThekngp5zdS-AwDqbOpS83L9tRPkyEReOeZRpW8PbVZxK0o2c-kxTbpWuO_xoG4ticZ-nPa5vgYYxLWvRjNLmiL1IRVOQ7m8E-3X3WAA0fQgz9gvFy8yJQw3uwIyi5gLLg37BVNJvWFGNoO_wJ3kkftteyZECqO0gnHfSsGutuG__KSn1CcIhPN5ahjdH5NSYPOdRWP-J7QcMLedPpKu5XgnJkXgQDfAMsLjl07003__m400F__Rhg13W00"

	//When
	actual := DeflateAndEncode([]byte(testPlantumlInput))

	//Then
	require.Equal(t, expected, actual, "Unexpected output")
}

func TestOutputPlantumlWithPng(t *testing.T) {
	//Given
	output := "test.png"
	plantuml := "http://www.plantuml.com/plantuml"
	umlInput := testPlantumlInput

	//When
	OutputPlantuml(output, plantuml, umlInput)

	//Then
	_, err := os.Stat(output)
	assert.False(t, os.IsNotExist(err))
}

func TestOutputPlantumlWithSvg(t *testing.T) {
	//Given
	output := "test.svg"
	plantuml := "http://www.plantuml.com/plantuml"
	umlInput := testPlantumlInput

	//When
	OutputPlantuml(output, plantuml, umlInput)

	//Then
	_, err := os.Stat(output)
	assert.False(t, os.IsNotExist(err))
}

func TestOutputPlantumlWithUml(t *testing.T) {
	//Given
	output := "test.uml"
	plantuml := "http://www.plantuml.com/plantuml"
	umlInput := testPlantumlInput

	//When
	OutputPlantuml(output, plantuml, umlInput)

	//Then
	_, err := os.Stat(output)
	assert.False(t, os.IsNotExist(err))
}

func TestOutputPlantumlWithoutUrl(t *testing.T) {
	//Given
	output := "test.uml"
	plantuml := ""
	umlInput := testPlantumlInput

	//When
	OutputPlantuml(output, plantuml, umlInput)

	//Then
	_, err := os.Stat(output)
	assert.False(t, os.IsNotExist(err))
}
