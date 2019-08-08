package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testPlantumlInput = `
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

const plantumlDotCom = "http://www.plantuml.com/plantuml"

func TestDeflateAndEncode(t *testing.T) {
	//Given
	const expected = "UDfSaKbhmp0GXU_pAnwvYqY6NaniKkXoAgGRFUGW9l4qY7gh99SkzByN9GvnUfBGzmrwZw5bYE" +
		"pZqDIqxThekngp5zdS-AwDqbOpS83L9tRPkyEReOeZRpW8PbVZxK0o2c-kxTbpWuO_xoG4ticZ-nPa5vgYYxLWv" +
		"RjNLmiL1IRVOQ7m8E-3X3WAA0fQgz9gvFy8yJQw3uwIyi5gLLg37BVNJvWFGNoO_wJ3kkftteyZECqO0gnHfSsG" +
		"utuG__KSn1CcIhPN5ahjdH5NSYPOdRWP-J7QcMLedPpKu5XgnJkXgQDfAMsLjl0N003__swwWGu0"

	//When
	actual, err := DeflateAndEncode([]byte(testPlantumlInput))
	require.NoError(t, err)

	//Then
	assert.Equal(t, expected, actual, "Unexpected output")
}

func testOutputPlantuml(t *testing.T, output, output2 string) {
	require.NoError(t, OutputPlantuml(output, plantumlDotCom, testPlantumlInput))
	_, err := os.Stat(output2)
	require.NoError(t, err)
	os.Remove(output2)
}

func TestOutputPlantumlWithPng(t *testing.T) {
	testOutputPlantuml(t, "test.png", "test.png")
}

func TestOutputPlantumlWithSvg(t *testing.T) {
	testOutputPlantuml(t, "test.svg", "test.svg")
}

func TestOutputPlantumlWithUml(t *testing.T) {
	testOutputPlantuml(t, "test.uml", "test.puml")
}

func TestEncode6bit(t *testing.T) {
	data := []struct {
		input    byte
		expected byte
	}{
		{0, 48},  // 0
		{63, 95}, // _
		{24, 79}, // O
	}

	for _, v := range data {
		v := v
		t.Run(string(int(v.input)), func(tt *testing.T) {
			actual := encode6bit(v.input)
			assert.Equal(tt, v.expected, actual)
		})
	}
}

func TestEncode6bitPanic(t *testing.T) {
	// Given
	b := byte(255)

	// Then
	assert.Panics(t, func() {
		encode6bit(b)
	}, "unexpected character!")
}

func TestOutPutWithWrongFormat(t *testing.T) {
	//Given
	output := "test.wrong"
	umlInput := testPlantumlInput

	//When
	require.Error(t, OutputPlantuml(output, plantumlDotCom, umlInput))

	//Then
	_, err := os.Stat(output)
	if !assert.True(t, os.IsNotExist(err)) {
		assert.NoError(t, os.Remove(output))
	}
}

func TestWrongHttpRequest(t *testing.T) {
	//Given
	url := "ww.plantuml.co"

	//When
	out, err := sendHTTPRequest(url)

	//Then
	assert.Nil(t, out)
	assert.NotNil(t, err)
}

func TestEncode(t *testing.T) {
	//Given
	data := []byte{'a'}

	//When
	r := encode(data)

	//Then
	assert.NotNil(t, r)
}
