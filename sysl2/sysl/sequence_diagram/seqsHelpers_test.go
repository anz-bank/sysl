package sd

import (
	"reflect"
	"testing"

	"github.com/anz-bank/sysl/src/proto"
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

func TestFindMatchItems(t *testing.T) {
	type args struct {
		origin string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Case-Null",
			args: args{""},
			want: nil,
		},
		{
			name: "Case-Convert Success",
			args: args{"%(appname)"},
			want: []string{"%(appname)"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindMatchItems(tt.args.origin); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindMatchItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveWrapper(t *testing.T) {
	type args struct {
		origin string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Case-Null",
			args: args{""},
			want: "",
		},
		{
			name: "Case-Convert Success",
			args: args{"%(appname)"},
			want: "appname",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveWrapper(tt.args.origin); got != tt.want {
				t.Errorf("RemoveWrapper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemovePercentSymbol(t *testing.T) {
	type args struct {
		origin string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Case-Null",
			args{""},
			"",
		},
		{
			"Case-Remove Percent",
			args{"%VariableA, %VariableB"},
			"VariableA, VariableB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemovePercentSymbol(tt.args.origin); got != tt.want {
				t.Errorf("RemovePercentSymbol() = %v, want %v", got, tt.want)
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
				"ep": epAttr,
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
