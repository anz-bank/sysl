package main

import (
	"encoding/json"
	"sort"

	"github.com/anz-bank/sysl/src/proto"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func addInt64(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueI64(lhs.GetI() + rhs.GetI())
}

func gtInt64(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueBool(lhs.GetI() > rhs.GetI())
}

func ltInt64(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueBool(lhs.GetI() < rhs.GetI())
}

func geInt64(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueBool(lhs.GetI() >= rhs.GetI())
}

func leInt64(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueBool(lhs.GetI() <= rhs.GetI())
}

func not(f func(lhs, rhs *sysl.Value) *sysl.Value) func(lhs, rhs *sysl.Value) *sysl.Value {
	return func(lhs, rhs *sysl.Value) *sysl.Value {
		return MakeValueBool(!f(lhs, rhs).GetB())
	}
}

func subInt64(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueI64(lhs.GetI() - rhs.GetI())
}

func mulInt64(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueI64(lhs.GetI() * rhs.GetI())
}

func divInt64(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueI64(lhs.GetI() / rhs.GetI())
}

func modInt64(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueI64(lhs.GetI() % rhs.GetI())
}

func addString(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueString(lhs.GetS() + rhs.GetS())
}

func cmpString(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueBool(lhs.GetS() == rhs.GetS())
}

func cmpInt(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueBool(lhs.GetI() == rhs.GetI())
}

func cmpBool(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueBool(lhs.GetB() == rhs.GetB())
}

func andBool(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueBool(lhs.GetB() && rhs.GetB())
}

func cmpNullTrue(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueBool(true)
}

func cmpNullFalse(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueBool(false)
}

func flattenListMap(txApp *sysl.Application, assign Scope, list *sysl.Value, scopeVar string, rhs *sysl.Expr) *sysl.Value {
	listResult := MakeValueList()
	for _, l := range list.GetList().Value {
		assign[scopeVar] = l
		appendItemToValueList(listResult.GetList(), Eval(txApp, assign, rhs))
	}
	return listResult
}

func flattenListList(txApp *sysl.Application, assign Scope, list *sysl.Value, scopeVar string, rhs *sysl.Expr) *sysl.Value {
	listResult := MakeValueList()
	for _, l := range list.GetList().Value {
		for _, ll := range l.GetList().Value {
			assign[scopeVar] = ll
			appendItemToValueList(listResult.GetList(), Eval(txApp, assign, rhs))
		}
	}
	return listResult
}

func flattenListSet(txApp *sysl.Application, assign Scope, list *sysl.Value, scopeVar string, rhs *sysl.Expr) *sysl.Value {
	listResult := MakeValueList()
	for _, l := range list.GetList().Value {
		for _, ll := range l.GetSet().Value {
			assign[scopeVar] = ll
			appendItemToValueList(listResult.GetList(), Eval(txApp, assign, rhs))
		}
	}
	return listResult
}

func flattenSetMap(txApp *sysl.Application, assign Scope, list *sysl.Value, scopeVar string, rhs *sysl.Expr) *sysl.Value {
	setResult := MakeValueSet()
	for _, l := range list.GetSet().Value {
		if l.GetMap() == nil {
			panic(errors.Errorf("flattenSetMap: expecting map instead of %v ", l))
		}
		assign[scopeVar] = l
		res := Eval(txApp, assign, rhs)
		switch x := res.Value.(type) {
		case *sysl.Value_Set:
			for _, ll := range x.Set.Value {
				appendItemToValueList(setResult.GetSet(), ll)
			}
		default:
			appendItemToValueList(setResult.GetSet(), res)
		}
	}
	return setResult
}

func flattenSetSet(txApp *sysl.Application, assign Scope, list *sysl.Value, scopeVar string, rhs *sysl.Expr) *sysl.Value {
	setResult := MakeValueSet()
	for _, l := range list.GetSet().Value {
		for _, ll := range l.GetSet().Value {
			assign[scopeVar] = ll
			appendItemToValueList(setResult.GetSet(), Eval(txApp, assign, rhs))
		}
	}
	return setResult
}

func concatList(lhs, rhs *sysl.Value) *sysl.Value {
	list := MakeValueList()

	list.GetList().Value = append(lhs.GetList().Value, rhs.GetList().Value...)
	logrus.Printf("concatList: lhs %d, rhs %d res %d\n", len(lhs.GetList().Value), len(rhs.GetList().Value), len(list.GetList().Value))
	return list
}

func setUnion(lhs, rhs *sysl.Value) *sysl.Value {
	itemType := getContainedType(lhs)
	if itemType == VALUE_NO_ARG {
		itemType = getContainedType(rhs)
	}

	if itemType == VALUE_NO_ARG {
		return MakeValueSet()
	}

	switch itemType {
	case VALUE_INT:
		unionSet := unionIntSets(intSet(lhs.GetSet().Value), intSet(rhs.GetSet().Value))
		logrus.Printf("Union set: lhs %d, rhs %d res %d\n", len(lhs.GetSet().Value), len(rhs.GetSet().Value), len(unionSet))
		return intSetToValueSet(unionSet)
	case VALUE_STRING:
		unionSet := unionStringSets(stringSet(lhs.GetSet().Value), stringSet(rhs.GetSet().Value))
		logrus.Printf("Union set: lhs %d, rhs %d res %d\n", len(lhs.GetSet().Value), len(rhs.GetSet().Value), len(unionSet))
		return stringSetToValueSet(unionSet)
	case VALUE_MAP:
		unionSet := unionMapSets(mapSet(lhs.GetSet().Value), mapSet(rhs.GetSet().Value))
		logrus.Printf("Union set: lhs %d, rhs %d res %d\n", len(lhs.GetSet().Value), len(rhs.GetSet().Value), len(unionSet))
		return mapSetToValueSet(unionSet)
	}
	panic(errors.Errorf("union of itemType: %s not supported", itemType.String()))
}

func stringInSet(lhs, rhs *sysl.Value) *sysl.Value {
	str := lhs.GetS()
	for _, v := range rhs.GetSet().Value {
		if str == v.GetS() {
			return MakeValueBool(true)
		}
	}
	return MakeValueBool(false)
}

func intSet(list []*sysl.Value) map[int64]struct{} {
	m := map[int64]struct{}{}
	for _, item := range list {
		m[item.GetI()] = struct{}{}
	}
	return m
}

func unionIntSets(lhs, rhs map[int64]struct{}) map[int64]struct{} {
	for key := range rhs {
		lhs[key] = struct{}{}
	}
	return lhs
}

func intSetToValueSet(lhs map[int64]struct{}) *sysl.Value {
	m := MakeValueSet()
	var keys []int

	for key := range lhs {
		keys = append(keys, int(key))
	}
	sort.Ints(keys)

	for _, key := range keys {
		appendItemToValueList(m.GetSet(), MakeValueI64(int64(key)))
	}
	return m
}

func stringSet(list []*sysl.Value) map[string]struct{} {
	m := map[string]struct{}{}

	for _, item := range list {
		m[item.GetS()] = struct{}{}
	}
	return m
}

func unionStringSets(lhs, rhs map[string]struct{}) map[string]struct{} {
	for key := range rhs {
		lhs[key] = struct{}{}
	}
	return lhs
}

func stringSetToValueSet(lhs map[string]struct{}) *sysl.Value {
	m := MakeValueSet()

	// for stable output
	var keys []string
	for key := range lhs {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		appendItemToValueList(m.GetSet(), MakeValueString(key))
	}
	return m
}

func mapSet(m []*sysl.Value) map[string]*sysl.Value {
	resultMap := map[string]*sysl.Value{}
	for _, item := range m {
		// Marshal() sorts the keys, so should get stable output.
		bytes, err := json.Marshal(item.GetMap().Items)
		if err == nil {
			resultMap[string(bytes)] = item
		}
	}
	return resultMap
}

func unionMapSets(lhs, rhs map[string]*sysl.Value) map[string]*sysl.Value {
	for key, val := range rhs {
		lhs[key] = val
	}
	return lhs
}

func mapSetToValueSet(lhs map[string]*sysl.Value) *sysl.Value {
	m := MakeValueSet()
	var keys []string
	for key := range lhs {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		appendItemToValueList(m.GetSet(), lhs[key])
	}
	return m
}

func stringInNull(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueBool(false)
}

func stringInList(lhs, rhs *sysl.Value) *sysl.Value {
	str := lhs.GetS()
	for _, v := range rhs.GetList().Value {
		if str == v.GetS() {
			return MakeValueBool(true)
		}
	}
	return MakeValueBool(false)
}

func whereSet(txApp *sysl.Application, assign Scope, list *sysl.Value, scopeVar string, rhs *sysl.Expr) *sysl.Value {
	setResult := MakeValueSet()
	for _, l := range list.GetSet().Value {
		assign[scopeVar] = l
		predicate := Eval(txApp, assign, rhs)
		if predicate.GetB() {
			appendItemToValueList(setResult.GetSet(), l)
		}
	}
	return setResult
}

func whereList(txApp *sysl.Application, assign Scope, list *sysl.Value, scopeVar string, rhs *sysl.Expr) *sysl.Value {
	listResult := MakeValueList()
	logrus.Printf("scope: %s, list len: %d", scopeVar, len(list.GetList().Value))
	for _, l := range list.GetList().Value {
		assign[scopeVar] = l
		predicate := Eval(txApp, assign, rhs)
		if predicate.GetB() {
			appendItemToValueList(listResult.GetList(), l)
		}
	}
	return listResult
}

func whereMap(txApp *sysl.Application, assign Scope, map_ *sysl.Value, scopeVar string, rhs *sysl.Expr) *sysl.Value {
	mapResult := MakeValueMap()
	logrus.Printf("scope: %s, list len: %d", scopeVar, len(map_.GetMap().Items))
	for key, val := range map_.GetMap().Items {
		m := MakeValueMap()
		addItemToValueMap(m, "key", MakeValueString(key))
		addItemToValueMap(m, "value", val)
		assign[scopeVar] = m
		predicate := Eval(txApp, assign, rhs)
		if predicate.GetB() {
			addItemToValueMap(mapResult, key, val)
		}
	}
	return mapResult
}
