package eval

import (
	"testing"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/stretchr/testify/require"
)

func makeNullValue() *sysl.Value {
	nullVal := sysl.Value{
		Value: &sysl.Value_Null_{
			Null: &sysl.Value_Null{},
		},
	}
	return &nullVal
}

func makeNullValueList() *sysl.Value {
	nullValueList := sysl.Value{
		Value: &sysl.Value_List_{
			List: nil,
		},
	}
	return &nullValueList
}

func makeNonNullValueList() *sysl.Value {
	nonNullValueList := sysl.Value{
		Value: &sysl.Value_List_{
			List: &sysl.Value_List{
				Value: []*sysl.Value{{
					Value: &sysl.Value_S{
						S: "",
					}},
				},
			},
		},
	}
	return &nonNullValueList
}

func TestCmpListNull_BothNull(t *testing.T) {
	lhs := makeNullValue()
	rhs := makeNullValue()

	require.Equal(t, true, cmpListNull(lhs, rhs).GetB())
}

func TestCmpListNull_RhsNull(t *testing.T) {
	lhs := makeNullValueList()
	rhs := makeNullValue()

	require.Equal(t, true, cmpListNull(lhs, rhs).GetB())
}

func TestCmpListNull_LhsNull(t *testing.T) {
	lhs := makeNullValue()
	rhs := makeNullValueList()

	require.Equal(t, true, cmpListNull(lhs, rhs).GetB())
}

func TestCmpListNull_LhsValid(t *testing.T) {
	lhs := makeNonNullValueList()
	rhs := makeNullValue()

	require.Equal(t, false, cmpListNull(lhs, rhs).GetB())
}

func TestCmpListNull_RhsValid(t *testing.T) {
	lhs := makeNullValue()
	rhs := makeNonNullValueList()

	require.Equal(t, false, cmpListNull(lhs, rhs).GetB())
}
