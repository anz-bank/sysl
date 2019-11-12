package eval

import (
	"testing"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/stretchr/testify/require"
)

func TestUnarySingle_NilPanic(t *testing.T) {
	require.Panics(t, func() { _ = unarySingle(nil) })
}

func TestUnarySingle_NotCollectionPanic(t *testing.T) {
	require.Panics(t, func() { _ = unarySingle(MakeValueBool(true)) })
}

func TestUnarySingle_NilListValuePanic(t *testing.T) {
	myVal := &sysl.Value{
		Value: &sysl.Value_List_{
			List: &sysl.Value_List{
				Value: nil,
			},
		},
	}
	require.Panics(t, func() { _ = unarySingle(myVal) })
}
func TestUnarySingle_NilListPanic(t *testing.T) {
	myVal := &sysl.Value{
		Value: &sysl.Value_List_{
			List: nil,
		},
	}
	require.Panics(t, func() { _ = unarySingle(myVal) })
}
func TestUnarySingle_OneValueOK(t *testing.T) {
	myList := MakeValueList(MakeValueBool(true))
	result := unarySingle(myList)
	require.Equal(t, result.GetB(), true)
}

func TestUnarySingle_TwoValuePanic(t *testing.T) {
	myList := MakeValueList(MakeValueBool(true), MakeValueBool(false))
	require.Panics(t, func() { _ = unarySingle(myList) })
}
