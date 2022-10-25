package ctn

import "reflect"


type Pair[A any, B any] func()(A,B)
func MakePair[A any, B any] (a A, b B) Pair[A,B] {
    return func() (A,B) {
        return a, b
    }
}
func ReflectMakePair(a reflect.Value, b reflect.Value) reflect.Value {
    return reflect.MakeFunc(ReflectTypePair(a.Type(), b.Type()), func([] reflect.Value) ([] reflect.Value) {
        return [] reflect.Value { a, b }
    })
}

func (p Pair[A,B]) First() A {
    var a, _ = p()
    return a
}
func (p Pair[A,B]) Second() B {
    var _, b = p()
    return b
}
func (p Pair[K,V]) Key() K {
    return p.First()
}
func (p Pair[K,V]) Value() V {
    return p.Second()
}
func ReflectPairUnpack(v reflect.Value) (reflect.Value, reflect.Value) {
    var ret = v.Call(nil)
    return ret[0], ret[1]
}

func ReflectTypePair(a reflect.Type, b reflect.Type) reflect.Type {
    return reflect.FuncOf (
        [] reflect.Type {},
        [] reflect.Type { a, b },
        false,
    )
}
func ReflectTypeMatchPair(t reflect.Type) (reflect.Type, reflect.Type, bool) {
    if !(t.Kind() == reflect.Func) { return nil, nil, false }
    if !(!(t.IsVariadic())) { return nil, nil, false }
    if !(t.NumIn() == 0) { return nil, nil, false }
    if !(t.NumOut() == 2) { return nil, nil, false }
    return t.Out(0), t.Out(1), true
}


