package importer

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_writer_writeUnion(t *testing.T) {
	tests := []struct {
		name     string
		arg      Type
		expected string
	}{
		{"writeUnion",
			&Union{name: "TestUnion", Options: FieldList{{Name: "Apple"}, {Name: "Orange"}}},
			"!union TestUnion:\n    Apple\n    Orange\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writebuf := bytes.NewBuffer([]byte{})
			w := newWriter(writebuf, logrus.New())
			w.writeUnion(tt.arg)                            //nolint:scopelint
			assert.Equal(t, tt.expected, writebuf.String()) //nolint:scopelint
		})
	}
}
