package parse

import "testing"

func TestMustUnescape(t *testing.T) {
	type args struct {
		endpoint string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{".", args{"%2E"}, "."},
		{":", args{"%3A"}, ":"},
		{"+", args{"%2B"}, "+"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//nolint:scopelint
			if got := MustUnescape(tt.args.endpoint); got != tt.want {
				t.Errorf("MustUnescape() = %v, want %v", got, tt.want)
			}
		})
	}
}
