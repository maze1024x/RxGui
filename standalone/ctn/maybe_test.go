package ctn

import (
	"testing"
	"reflect"
)

func TestMaybe(t *testing.T) {
	var a_opt = Just[int](1)
	var b_opt = Nothing[string]()
	if a, some := a_opt.Value(); some {
		if a != 1 {
			t.Fatalf("wrong value")
		}
	} else {
		t.Fatalf("wrong behavior")
	}
	if _, some := b_opt.Value(); some {
		t.Fatalf("wrong behavior")
	}
	if ar, some := ReflectMaybeValue(reflect.ValueOf(a_opt)); some {
		var a = ar.Interface().(int)
		if a != 1 {
			t.Fatalf("wrong reflect value")
		}
	} else {
		t.Fatalf("wrong reflect behavior")
	}
	if _, some := ReflectMaybeValue(reflect.ValueOf(b_opt)); some {
		t.Fatalf("wrong reflect behavior")
	}
	{ var u = reflect.TypeOf(struct { a int } { a: 1 })
	var v, ok = ReflectTypeMatchMaybe(ReflectTypeMaybe(u))
	if !(ok) || u != v {
		t.Fatalf("reflect type operations inconsistent")
	} }
}


