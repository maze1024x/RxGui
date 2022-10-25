package pseudounion

import "reflect"


type Tag int
var tagType = reflect.TypeOf(Tag(0))

func Load(u interface{}) interface{} {
    var ru = reflect.ValueOf(u)
    if ru.Kind() != reflect.Struct {
        checkKind(ru, reflect.Pointer)
        ru = ru.Elem()
        checkKind(ru, reflect.Struct)
    }
    var index = int(assertTagType(ru.Field(0)).Int())
    if index <= 0 {
        panic("invalid pseudo-union index: zero or negative")
    }
    if index > ru.NumField() {
        panic("invalid pseudo-union index: too big")
    }
    return ru.Field(index).Interface()
}

func Store[T any] (v interface{}) T {
    var u T
    var ru = reflect.ValueOf(&u).Elem()
    if ru.Kind() != reflect.Struct {
        checkKind(ru, reflect.Pointer)
        ru.Set(reflect.New(ru.Type().Elem()))
        ru = ru.Elem()
        checkKind(ru, reflect.Struct)
    }
    var rv = reflect.ValueOf(v)
    for i := 1; i < ru.NumField(); i += 1 {
        if rv.Type().AssignableTo(ru.Type().Field(i).Type) {
            ru.Field(i).Set(rv)
            assertTagType(ru.Field(0)).SetInt(int64(i))
            return u
        }
    }
    panic("invalid pseudo-union value: none of fields assignable")
}

func checkKind(ru reflect.Value, expected reflect.Kind) {
    if ru.Kind() != expected {
        panic("invalid pseudo-union type: " + ru.Type().String())
    }
}
func assertTagType(rv reflect.Value) reflect.Value {
    if rv.Type() != tagType {
        panic("invalid pseudo-union")
    }
    return rv
}


