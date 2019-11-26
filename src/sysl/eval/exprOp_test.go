package eval

import (
	"testing"

	"github.com/stretchr/testify/assert"

	sysl "github.com/anz-bank/sysl/src/proto_old"
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

func TestConcat_NilNilPanic(t *testing.T) {
	defer func() {
		r := recover()
		require.NotNil(t, r)
	}()
	_ = concat(nil, nil)
}

func TestConcat_NilListPanic(t *testing.T) {
	defer func() {
		r := recover()
		require.NotNil(t, r)
	}()
	_ = concat(makeNullValueList().GetList(), makeNullValueList().GetList())
}

func TestConcat_EmptyList(t *testing.T) {
	result := concatListList(MakeValueList(), MakeValueList())
	assert.Equal(t, []*sysl.Value{}, result.GetList().Value)
}

func TestConcat_OneItemList(t *testing.T) {
	lhs := MakeValueList(MakeValueBool(true))
	rhs := MakeValueList(MakeValueBool(false))
	result := concatListList(lhs, rhs)
	assert.Equal(t, true, result.GetList().Value[0].GetB())
	assert.Equal(t, false, result.GetList().Value[1].GetB())
}

func TestConcat_OneItemSet(t *testing.T) {
	lhs := MakeValueList(MakeValueBool(true))
	rhs := MakeValueSet()
	rhs.GetSet().Value = append(rhs.GetSet().Value, MakeValueBool(false))

	result := concatListSet(lhs, rhs)
	assert.Equal(t, true, result.GetList().Value[0].GetB())
	assert.Equal(t, false, result.GetList().Value[1].GetB())
}

func TestFlattenSetList_AllNilPanic(t *testing.T) {
	defer func() {
		r := recover()
		require.NotNil(t, r)
	}()
	_ = flattenSetList(newExprEval(nil), Scope{}, nil, "", nil)
}

func TestFlattenSetList_NameExpr(t *testing.T) {
	nameExpr := &sysl.Expr{
		Expr: &sysl.Expr_Name{
			Name: "var",
		},
	}

	s := Scope{}
	s.AddString("var", "name1")

	setOfList := MakeValueSet()
	setOfList.GetSet().Value = append(setOfList.GetSet().Value,
		MakeValueList(MakeValueString("name2")))
	result := flattenSetList(newExprEval(nil), s, setOfList, "var", nameExpr)
	require.Equal(t, "name2", result.GetSet().Value[0].GetS())
}

func TestFlattenListSet_NameExpr(t *testing.T) {
	nameExpr := &sysl.Expr{
		Expr: &sysl.Expr_Name{
			Name: "var",
		},
	}

	s := Scope{}
	s.AddString("var", "name1")

	listOfSet := MakeValueList(MakeValueSet())
	listOfSet.GetList().Value[0].GetSet().Value = append(
		listOfSet.GetList().Value[0].GetSet().Value,
		MakeValueString("name2"))
	result := flattenListSet(newExprEval(nil), s, listOfSet, "var", nameExpr)
	require.Equal(t, "name2", result.GetList().Value[0].GetS())
}
