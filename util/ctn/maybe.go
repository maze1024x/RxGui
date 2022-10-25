package ctn

import "reflect"


type Maybe[T any] func()(T,bool)
func MakeMaybe[T any] (t T, ok bool) Maybe[T] {
	if ok {
		return Just(t)
	} else {
		return nil
	}
}
func Just[T any] (t T) Maybe[T] {
	return func() (T,bool) {
		return t, true
	}
}
func Nothing[T any] () Maybe[T] {
	return nil
}
func ReflectJust(v reflect.Value) reflect.Value {
	return reflect.MakeFunc(ReflectTypeMaybe(v.Type()), func([] reflect.Value) ([] reflect.Value) {
		return [] reflect.Value { v, reflect.ValueOf(true) }
	})
}
func ReflectNothing(t reflect.Type) reflect.Value {
	return reflect.Zero(ReflectTypeMaybe(t))
}

func (opt Maybe[T]) Value() (T,bool) {
	if opt != nil {
		return opt()
	} else {
		return zero[T](), false
	}
}
func ReflectMaybeValue(v reflect.Value) (reflect.Value, bool) {
	if v.IsNil() {
		return reflect.Value {}, false
	} else {
		return v.Call(nil)[0], true
	}
}

func ReflectTypeMaybe(t reflect.Type) reflect.Type {
	return reflect.FuncOf (
		[] reflect.Type {},
		[] reflect.Type { t, reflect.TypeOf(false) },
		false,
	)
}
func ReflectTypeMatchMaybe(t reflect.Type) (reflect.Type, bool) {
	if !(t.Kind() == reflect.Func) { return nil, false }
	if !(!(t.IsVariadic())) { return nil, false }
	if !(t.NumIn() == 0) { return nil, false }
	if !(t.NumOut() == 2) { return nil, false }
	if !(t.Out(1) == reflect.TypeOf(false)) { return nil, false }
	return t.Out(0), true
}


