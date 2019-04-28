package seqs

import (
	"reflect"
	"testing"

	"github.com/anz-bank/sysl/src/proto"
	"github.com/stretchr/testify/assert"
)

func TestTransformBlackBoxes(t *testing.T) {
	type args struct {
		blackboxes []*sysl.Attribute
	}

	eltFirst := []*sysl.Attribute{
		{
			Attribute: &sysl.Attribute_S{
				S: "Value A",
			},
		},
		{
			Attribute: &sysl.Attribute_S{
				S: "Value B",
			},
		},
	}
	attrFirst := &sysl.Attribute{
		Attribute: &sysl.Attribute_A{
			A: &sysl.Attribute_Array{
				Elt: eltFirst,
			},
		},
	}
	eltSecond := []*sysl.Attribute{
		{
			Attribute: &sysl.Attribute_S{
				S: "Value C",
			},
		},
		{
			Attribute: &sysl.Attribute_S{
				S: "Value D",
			},
		},
	}
	attrSecond := &sysl.Attribute{
		Attribute: &sysl.Attribute_A{
			A: &sysl.Attribute_Array{
				Elt: eltSecond,
			},
		},
	}

	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "Case-Null",
			args: args{blackboxes: []*sysl.Attribute{}},
			want: [][]string{},
		},
		{
			name: "Case-ConvertSuccess",
			args: args{blackboxes: []*sysl.Attribute{attrFirst, attrSecond}},
			want: [][]string{{"Value A", "Value B"}, {"Value C", "Value D"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TransformBlackBoxes(tt.args.blackboxes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransformBlackBoxes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseBlackBoxesFromArgument(t *testing.T) {
	type args struct {
		blackboxFlags []string
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "Case-Null",
			args: args{[]string{}},
			want: [][]string{},
		},
		{
			name: "Case-ConvertSuccess",
			args: args{[]string{"Value A,Value B", "Value C,Value D"}},
			want: [][]string{{"Value A", "Value B"}, {"Value C", "Value D"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseBlackBoxesFromArgument(tt.args.blackboxFlags); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseBlackBoxesFromArgument() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeAttributes(t *testing.T) {
	type args struct {
		app   map[string]*sysl.Attribute
		edpnt map[string]*sysl.Attribute
	}

	appAttr := &sysl.Attribute{
		Attribute: &sysl.Attribute_S{
			S: "Value A",
		},
	}
	appMap := map[string]*sysl.Attribute{
		"app": appAttr,
	}
	epAttr := &sysl.Attribute{
		Attribute: &sysl.Attribute_S{
			S: "Value B",
		},
	}
	epMap := map[string]*sysl.Attribute{
		"ep": epAttr,
	}
	tests := []struct {
		name string
		args args
		want map[string]*sysl.Attribute
	}{
		{
			"Case-Null",
			args{},
			map[string]*sysl.Attribute{},
		},
		{
			"Case-Merge app",
			args{appMap, map[string]*sysl.Attribute{}},
			map[string]*sysl.Attribute{
				"app": appAttr,
			},
		},
		{
			"Case-Merge ep",
			args{map[string]*sysl.Attribute{}, epMap},
			map[string]*sysl.Attribute{
				"ep": epAttr,
			},
		},
		{
			"Case-Merge app and ep",
			args{appMap, epMap},
			map[string]*sysl.Attribute{
				"app": appAttr,
				"ep":  epAttr,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeAttributes(tt.args.app, tt.args.edpnt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeAttributes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopyBlackboxes(t *testing.T) {
	// given
	m := map[string]string{
		"keyA": "value A",
	}

	// when
	r := copyBlackboxes(m)

	// then
	assert.Equal(t, m, r, "unexpected map")
}

func TestCopyBlackboxesByNil(t *testing.T) {
	// given
	var m map[string]string

	// when
	r := copyBlackboxes(m)

	// then
	assert.NotNil(t, r)
}

func TestGetAppName(t *testing.T) {
	// given
	a := &sysl.AppName{
		Part: []string{"test", "name"},
	}

	// when
	actual := getAppName(a)

	// then
	assert.Equal(t, "test :: name", actual, "unexpected result")
}

func TestGetAppAttr(t *testing.T) {
	// given
	attr := map[string]*sysl.Attribute{
		"attr1": {},
	}
	m := &sysl.Module{
		Apps: map[string]*sysl.Application{
			"test": {Attrs: attr},
		},
	}

	// when
	actual := getApplicationAttrs(m, "test")

	// then
	assert.Equal(t, attr, actual)
}

func TestGetAppAttrWhenAppNotExist(t *testing.T) {
	// given
	m := &sysl.Module{
		Apps: make(map[string]*sysl.Application),
	}

	// when
	actual := getApplicationAttrs(m, "test")

	// then
	assert.Nil(t, actual)
}

func TestSortedISOCtrlSlice(t *testing.T) {
	// given
	attrs := map[string]*sysl.Attribute{
		"iso_ctrl_11_txt": {},
		"iso_ctrl_12_txt": {},
		"iso_ctrl_5_txt":  {},
	}

	// when
	actual := getSortedISOCtrlSlice(attrs)

	// then
	assert.Equal(t, []string{"11", "12", "5"}, actual)
}

func TestSortedISOCtrlSliceEmpty(t *testing.T) {
	// given
	attrs := make(map[string]*sysl.Attribute)

	// when
	actual := getSortedISOCtrlSlice(attrs)

	// then
	assert.Equal(t, []string{}, actual)
}

func TestSortedISOCtrlStr(t *testing.T) {
	// given
	attrs := map[string]*sysl.Attribute{
		"iso_ctrl_11_txt": {},
		"iso_ctrl_12_txt": {},
		"iso_ctrl_5_txt":  {},
	}

	// when
	actual := getSortedISOCtrlStr(attrs)

	// then
	assert.Equal(t, "11, 12, 5", actual)
}

func TestSortedISOCtrlStrEmpty(t *testing.T) {
	// given
	attrs := make(map[string]*sysl.Attribute)

	// when
	actual := getSortedISOCtrlStr(attrs)

	// then
	assert.Equal(t, "", actual)
}

func TestFormatArgs(t *testing.T) {
	// given
	m := &sysl.Module{
		Apps: map[string]*sysl.Application{
			"test": {
				Types: map[string]*sysl.Type{
					"User": {
						Attrs: map[string]*sysl.Attribute{
							"iso_conf": {
								Attribute: &sysl.Attribute_S{
									S: "Red",
								},
							},
							"iso_integ": {
								Attribute: &sysl.Attribute_S{
									S: "I",
								},
							},
						},
					},
				},
			},
		},
	}

	// when
	actual := formatArgs(m, "test", "User")

	assert.Equal(t, "<color blue>test.User</color> <<color red>R, I</color>>", actual)
}

func TestFormatArgsWithoutIsoInteg(t *testing.T) {
	// given
	m := &sysl.Module{
		Apps: map[string]*sysl.Application{
			"test": {
				Types: map[string]*sysl.Type{
					"User": {
						Attrs: map[string]*sysl.Attribute{
							"iso_conf": {
								Attribute: &sysl.Attribute_S{
									S: "Red",
								},
							},
						},
					},
				},
			},
		},
	}

	// when
	actual := formatArgs(m, "test", "User")

	assert.Equal(t, "<color blue>test.User</color> <<color red>R, ?</color>>", actual)
}

func TestFormatArgsWithoutIsoConf(t *testing.T) {
	// given
	m := &sysl.Module{
		Apps: map[string]*sysl.Application{
			"test": {
				Types: map[string]*sysl.Type{
					"User": {
						Attrs: map[string]*sysl.Attribute{
							"iso_integ": {
								Attribute: &sysl.Attribute_S{
									S: "I",
								},
							},
						},
					},
				},
			},
		},
	}

	// when
	actual := formatArgs(m, "test", "User")

	assert.Equal(t, "<color blue>test.User</color> <<color green>?, I</color>>", actual)
}

func TestFormatArgsWithoutAttrs(t *testing.T) {
	// given
	m := &sysl.Module{
		Apps: map[string]*sysl.Application{
			"test": {
				Types: map[string]*sysl.Type{
					"User": {
						Attrs: make(map[string]*sysl.Attribute),
					},
				},
			},
		},
	}

	// when
	actual := formatArgs(m, "test", "User")

	assert.Equal(t, "<color blue>test.User</color> <<color green>?, ?</color>>", actual)
}

func TestFormatArgsWithoutParameterTypeName(t *testing.T) {
	// given
	m := &sysl.Module{
		Apps: map[string]*sysl.Application{
			"test": {
				Types: map[string]*sysl.Type{
					"User": {
						Attrs: make(map[string]*sysl.Attribute),
					},
				},
			},
		},
	}

	// when
	actual := formatArgs(m, "test", "")

	assert.Equal(t, "<color blue>test.</color> <<color green>?, ?</color>>", actual)
}

func TestFormatArgsWithoutAppName(t *testing.T) {
	// given
	m := &sysl.Module{
		Apps: map[string]*sysl.Application{
			"test": {
				Types: map[string]*sysl.Type{
					"User": {
						Attrs: make(map[string]*sysl.Attribute),
					},
				},
			},
		},
	}

	// when
	actual := formatArgs(m, "", "User")

	assert.Equal(t, "<color blue>.User</color> <<color green>?, ?</color>>", actual)
}

func TestFormatArgsWithoutAppNameAndParameterTypeName(t *testing.T) {
	// given
	m := &sysl.Module{
		Apps: map[string]*sysl.Application{
			"test": {
				Types: map[string]*sysl.Type{
					"User": {
						Attrs: make(map[string]*sysl.Attribute),
					},
				},
			},
		},
	}

	// when
	actual := formatArgs(m, "", "")

	assert.Equal(t, "<color blue>.</color> <<color green>?, ?</color>>", actual)
}

func TestFormatReturnParam(t *testing.T) {
	// given
	m := &sysl.Module{
		Apps: map[string]*sysl.Application{
			"test": {
				Types: map[string]*sysl.Type{
					"User": {
						Attrs: make(map[string]*sysl.Attribute),
					},
				},
			},
		},
	}

	// when
	actual := formatReturnParam(m, "test.User")

	assert.Equal(t, []string{"<color blue>test.User</color> <<color green>?, ?</color>>"}, actual)
}

func TestFormatReturnParamSplit(t *testing.T) {
	// given
	m := &sysl.Module{
		Apps: map[string]*sysl.Application{
			"test": {
				Types: map[string]*sysl.Type{
					"User": {
						Attrs: make(map[string]*sysl.Attribute),
					},
				},
			},
		},
	}

	// when
	actual := formatReturnParam(m, "test.User,profile<:test.User,Bool,set of test.User, one of {test.User,ab}")

	expected := []string{
		"<color blue>test.User</color> <<color green>?, ?</color>>",
		"<color blue>test.User</color> <<color green>?, ?</color>>",
		"<color blue>test.User</color> <<color green>?, ?</color>>",
		"<color blue>test.User</color> <<color green>?, ?</color>>",
		"ab",
	}

	assert.Equal(t, expected, actual)
}

func TestGetReturnPayload(t *testing.T) {
	stmts := []*sysl.Statement{
		{
			Stmt: &sysl.Statement_Call{},
		},
		{
			Stmt: &sysl.Statement_Action{},
		},
		{
			Stmt: &sysl.Statement_Ret{
				Ret: &sysl.Return{
					Payload: "test",
				},
			},
		},
	}

	actual := getReturnPayload(stmts)

	assert.Equal(t, "test", actual)
}

func TestGetReturnPayloadWithAlt(t *testing.T) {
	stmts := []*sysl.Statement{
		{
			Stmt: &sysl.Statement_Alt{
				Alt: &sysl.Alt{
					Choice: []*sysl.Alt_Choice{
						{
							Cond: "cond 1",
							Stmt: []*sysl.Statement{},
						},
						{
							Cond: "cond 2",
							Stmt: []*sysl.Statement{
								{
									Stmt: &sysl.Statement_Ret{
										Ret: &sysl.Return{
											Payload: "test",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	actual := getReturnPayload(stmts)

	assert.Equal(t, "test", actual)
}

func TestGetReturnPayloadWithCond(t *testing.T) {
	stmts := []*sysl.Statement{
		{
			Stmt: &sysl.Statement_Cond{
				Cond: &sysl.Cond{
					Test: "cond 1",
					Stmt: []*sysl.Statement{
						{
							Stmt: &sysl.Statement_Ret{
								Ret: &sysl.Return{
									Payload: "test",
								},
							},
						},
					},
				},
			},
		},
	}

	actual := getReturnPayload(stmts)

	assert.Equal(t, "test", actual)
}

func TestGetReturnPayloadWithLoop(t *testing.T) {
	stmts := []*sysl.Statement{
		{
			Stmt: &sysl.Statement_Loop{
				Loop: &sysl.Loop{
					Mode:      sysl.Loop_WHILE,
					Criterion: "criterion",
					Stmt: []*sysl.Statement{
						{
							Stmt: &sysl.Statement_Ret{
								Ret: &sysl.Return{
									Payload: "test",
								},
							},
						},
					},
				},
			},
		},
	}

	actual := getReturnPayload(stmts)

	assert.Equal(t, "test", actual)
}

func TestGetReturnPayloadWithLoopN(t *testing.T) {
	stmts := []*sysl.Statement{
		{
			Stmt: &sysl.Statement_LoopN{
				LoopN: &sysl.LoopN{
					Count: 10,
					Stmt: []*sysl.Statement{
						{
							Stmt: &sysl.Statement_Ret{
								Ret: &sysl.Return{
									Payload: "test",
								},
							},
						},
					},
				},
			},
		},
	}

	actual := getReturnPayload(stmts)

	assert.Equal(t, "test", actual)
}

func TestGetReturnPayloadWithForeach(t *testing.T) {
	stmts := []*sysl.Statement{
		{
			Stmt: &sysl.Statement_Foreach{
				Foreach: &sysl.Foreach{
					Collection: "collection 1",
					Stmt: []*sysl.Statement{
						{
							Stmt: &sysl.Statement_Ret{
								Ret: &sysl.Return{
									Payload: "test",
								},
							},
						},
					},
				},
			},
		},
	}

	actual := getReturnPayload(stmts)

	assert.Equal(t, "test", actual)
}

func TestGetReturnPayloadWithGroup(t *testing.T) {
	stmts := []*sysl.Statement{
		{
			Stmt: &sysl.Statement_Group{
				Group: &sysl.Group{
					Title: "group 1",
					Stmt: []*sysl.Statement{
						{
							Stmt: &sysl.Statement_Ret{
								Ret: &sysl.Return{
									Payload: "test",
								},
							},
						},
					},
				},
			},
		},
	}

	actual := getReturnPayload(stmts)

	assert.Equal(t, "test", actual)
}

func TestGetAndFmtParam(t *testing.T) {
	m := &sysl.Module{
		Apps: map[string]*sysl.Application{
			"test": {
				Types: map[string]*sysl.Type{
					"User": {
						Attrs: map[string]*sysl.Attribute{
							"iso_conf": {
								Attribute: &sysl.Attribute_S{
									S: "Red",
								},
							},
							"iso_integ": {
								Attribute: &sysl.Attribute_S{
									S: "I",
								},
							},
						},
					},
				},
			},
		},
	}

	p := []*sysl.Param{
		{
			Name: "profile",
			Type: &sysl.Type{
				Type: &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{
						Ref: &sysl.Scope{
							Appname: &sysl.AppName{
								Part: []string{"test"},
							},
							Path: []string{"User"},
						},
					},
				},
			},
		},
	}

	actual := getAndFmtParam(m, p)

	assert.Equal(t, []string{"<color blue>test.User</color> <<color red>R, I</color>>"}, actual)
}

func TestNormalizeEndpointName(t *testing.T) {
	actual := normalizeEndpointName("a -> b")

	assert.Equal(t, " ⬄ b", actual)
}
