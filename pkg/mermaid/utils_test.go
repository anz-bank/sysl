package mermaid

import "testing"

func TestCleanString(t *testing.T) {
	type args struct {
		temp string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"RemoveSpaces", args{"a b"}, "ab"},
		{"CurlyBraces", args{"{b}"}, "_b_"},
		{"SquareBraces", args{"[b]"}, "_b_"},
		{"DoubleQuote", args{`"b"`}, "b"},
		{"Tilda", args{"~"}, ""},
		{"Colon", args{"a:b"}, "a_b"},
		{"LessThan", args{"a<b"}, "ab"},
		{"URLUnencode:", args{"%2E"}, "."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//nolint:scopelint
			if got := CleanString(tt.args.temp); got != tt.want {
				t.Errorf("CleanString() = %v, want %v", got, tt.want)
			}
		})
	}
}
