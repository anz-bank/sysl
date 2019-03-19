package seqs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/anz-bank/sysl/src/proto"
)

func TestMakeFormatParser(t *testing.T) {
	fp := MakeFormatParser("%(appname)")

	assert.NotNil(t, fp)
	assert.Equal(t, "%(appname)", fp.self)
}

func TestLabelEndpoint(t *testing.T) {
	// Given
	p := &EndpointLabelerParam{
		EndpointName: "Login",
	}
	fp := MakeFormatParser("%(epname)")

	// When
	formatStr := fp.LabelEndpoint(p)

	// Then
	assert.Equal(t, p.EndpointName, formatStr)
}

func TestLabelApp(t *testing.T) {
	// Given
	appName := "Project"
	attrs := map[string]*sysl.Attribute{}
	fp := MakeFormatParser("%(appname)")

	// When
	formatStr := fp.LabelApp(appName, "", attrs)

	// Then
	assert.Equal(t, appName, formatStr)
}

func TestFmtSeq(t *testing.T) {
	// Given
	fp := MakeFormatParser("%(seqtitle)")
	seqtitleAttr := &sysl.Attribute{
		Attribute: &sysl.Attribute_S{
			S: "Diagram",
		},
	}
	seqtitleMap := map[string]*sysl.Attribute{
		"seqtitle": seqtitleAttr,
	}

	// When
	formatStr := fp.FmtSeq("", "", seqtitleMap)

	// Then
	assert.Equal(t, "Diagram", formatStr)
}

func TestFmtOutput(t *testing.T) {
	// Given
	fp := MakeFormatParser("%(epname).png")
	endpointName := "Login"

	// When
	formatStr := fp.FmtOutput("", endpointName, "", map[string]*sysl.Attribute{})

	// Then
	assert.Equal(t, "Login.png", formatStr)
}

func TestParse(t *testing.T) {
	// Given
	attrs := map[string]string{}
	fp := MakeFormatParser("1ba%%%%(DT?%(@c2?//%(@c4?--%(cc?dd|edd)--|bc)//\n|cc)|bb)**%(appname)**")

	// When
	formatStr := fp.Parse(attrs)

	// Then
	assert.Equal(t, "1ba%%(DT?cc|bb)****", formatStr)
}

func TestParseUnclosedExpansion(t *testing.T) {
	// Given
	attrs := map[string]string{}
	fp := MakeFormatParser("1ba%%%%(DT?%(@c2?//%(@c4?--%(cc?dd|edd)--|bc)//\n|cc)|bb)**%(appname**")

	// Then
	assert.Panics(t, func() {
		fp.Parse(attrs)
	}, "unclosed expansion")
}

func TestParseMissingVariable(t *testing.T) {
	// Given
	attrs := map[string]string{}
	fp := MakeFormatParser("1ba%%%%(DT?%(@c2?//%(@c4?--%(cc?dd|edd)--|bc)//\n|cc)|bb)**%()**")

	// Then
	assert.Panics(t, func() {
		fp.Parse(attrs)
	}, "missing variable reference")
}

func TestParseMissingConditionValue(t *testing.T) {
	// Given
	attrs := map[string]string{}
	fp := MakeFormatParser("1ba%%%%(DT?%(@c2==?//%(@c4?--%(cc?dd|edd)--|bc)//\n|cc)|bb)**%(appname)**")

	// Then
	assert.Panics(t, func() {
		fp.Parse(attrs)
	}, "missing conditional value")
}

func TestParseWithEqualConditionValue(t *testing.T) {
	// Given
	attrs := map[string]string{
		"c2": "aa",
	}
	fp := MakeFormatParser("1ba%%%%(DT?%(@c2=='aa'?//%(@c4?--%(cc?dd|edd)--|bc)//\n|cc)|bb)**%(appname)**")

	// When
	formatStr := fp.Parse(attrs)

	// Then
	assert.Equal(t, `1ba%%(DT?cc|bb)****`, formatStr)
}

func TestParseWithNotEqualConditionValue(t *testing.T) {
	// Given
	attrs := map[string]string{
		"c2": "ab",
	}
	fp := MakeFormatParser("1ba%%%%(DT?%(@c2!='aa'?//%(@c4?--%(cc?dd|edd)--|bc)//\n|cc)|bb)**%(appname)**")

	// When
	formatStr := fp.Parse(attrs)

	// Then
	assert.Equal(t, `1ba%%(DT?//bc//\n|bb)****`, formatStr)
}

func TestParseSearchValue(t *testing.T) {
	// Given
	attrs := map[string]string{
		"c2": "ab",
	}
	fp := MakeFormatParser("1ba%%%%(DT?%(@c2~/ab/?//%(@c4?--%(cc?dd|edd)--|bc)//\n|cc)|bb)**%(appname)**")

	// When
	formatStr := fp.Parse(attrs)

	// Then
	assert.Equal(t, `1ba%%(DT?cc|bb)****`, formatStr)
}

func TestMergeAttributesMap(t *testing.T) {
	// Given
	seqtitleAttr := &sysl.Attribute{
		Attribute: &sysl.Attribute_S{
			S: "Diagram",
		},
	}
	seqtitleMap := map[string]*sysl.Attribute{
		"seqtitle": seqtitleAttr,
	}
	valMap := map[string]string{
		"appname": "Project",
	}

	// When
	mergeAttributesMap(valMap, seqtitleMap)

	// Then
	assert.Equal(t, map[string]string{"appname": "Project", "@seqtitle": "Diagram"}, valMap)
}
