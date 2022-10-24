package core

import (
	"fmt"
	"bytes"
	"errors"
	"strings"
	"strconv"
	"reflect"
	"math/big"
	"encoding/json"
	"rxgui/lang/source"
	"rxgui/lang/typsys"
)


type ReflectType struct { type_ typsys.CertainType }
func MakeReflectType(t typsys.CertainType) ReflectType {
	var u_ = convertToReflectiveType(t.Type)
	var u = typsys.CertainType { Type: u_ }
	return ReflectType { u }
}
func (rt ReflectType) CertainType() typsys.CertainType {
	return rt.type_
}
func (rt ReflectType) Type() typsys.Type {
	return rt.CertainType().Type
}
var genericType = typsys.Type(typsys.RefType {
	Def:  source.MakeRef("", T_GenericType),
	Args: [] typsys.Type {},
})
func convertToReflectiveType(t typsys.Type) typsys.Type {
	return typsys.Transform(t, func(t typsys.Type) (typsys.Type, bool) {
		switch t.(type) {
		case typsys.InferringType:
			return genericType, true
		case typsys.ParameterType:
			return genericType, true
		case typsys.RefType:
			return nil, false
		default:
			panic("impossible branch")
		}
	})
}

type ReflectValue struct {
	type_  ReflectType
	value  Object
}
func AssumeValidReflectValue(t ReflectType, v Object) ReflectValue {
	return ReflectValue {
		type_: t,
		value: v,
	}
}
func (rv ReflectValue) Type() ReflectType {
	return rv.type_
}
func (rv ReflectValue) Value() Object {
	return rv.value
}

type ReflectObservable struct {
	innerType        ReflectType
	observableValue  Observable
}
func (rv ReflectValue) CastToReflectObservable() ReflectObservable {
	if T, ok := rv.Type().Type().(typsys.RefType); ok {
	if T.Def.Namespace == "" && T.Def.ItemName == T_Observable {
	if len(T.Args) == 1 {
		var t_ = T.Args[0]
		var t = typsys.CertainType { Type: t_ }
		return ReflectObservable {
			innerType:       MakeReflectType(t),
			observableValue: GetObservable(rv.Value()),
		}
	}}}
	panic("cannot cast non-observable object to ReflectObservable")
}
func (ro ReflectObservable) InnerType() ReflectType {
	return ro.innerType
}
func (ro ReflectObservable) ObservableValue() Observable {
	return ro.observableValue
}

type SerializationContext interface {
	LookupType(rt ReflectType) (*typsys.TypeDef, source.Ref, ([] typsys.Type), bool)
}
func Marshal(rv ReflectValue, ctx SerializationContext) ([] byte, error) {
	var m = jsonMarshaler { rv, ctx }
	var b, err = m.MarshalJSON()
	if err != nil {
		if cause := errors.Unwrap(err); cause != nil {
			return nil, cause
		}
		return nil, err
	}
	return b, nil
}
func Unmarshal(b ([] byte), rt ReflectType, ctx SerializationContext) (Object, error) {
	var m = jsonUnmarshaler { rt, nil, ctx }
	var err = json.Unmarshal(b, &m)
	if err != nil { return nil, err }
	return m.value, nil
}
type jsonMarshaler struct {
	value    ReflectValue
	context  SerializationContext
}
func (m jsonMarshaler) MarshalJSON() ([] byte, error) {
	var rt = m.value.Type()
	var obj = m.value.Value()
	var ctx = m.context
	var def, ref, args, exists = ctx.LookupType(rt)
	if !(exists) { panic("something went wrong") }
	switch content := def.Content.(type) {
	case typsys.NativeContent:
		if ref.Namespace == "" {
			switch ref.ItemName {
			case T_Null:
				return ([] byte)("null"), nil
			case T_Bool:
				return json.Marshal(GetBool(obj))
			case T_Int:
				return json.Marshal(GetIntAsRawBigInt(obj))
			case T_Float:
				// NOTE: error on NaN, +Inf, -Inf
				return json.Marshal(GetFloat(obj))
			case T_String:
				return json.Marshal(GetString(obj))
			case T_Bytes:
				return json.Marshal(GetBytes(obj))
			case T_List:
				var l = GetList(obj)
				var slice = make([] jsonMarshaler, l.Length())
				var item_rt = firstArgType(args)
				l.ForEachWithIndex(func(i int, item_v Object) {
					var item_rv = AssumeValidReflectValue(item_rt, item_v)
					slice[i] = jsonMarshaler { item_rv, ctx }
				})
				return json.Marshal(slice)
			}
		}
	case typsys.Record:
		var r = (*obj).(Record)
		var buf jsonMarshalerRecordStructBuilder
		for i, field := range content.FieldList {
			var field_name = field.Name
			var field_rt = inflateFieldType(field, def, args)
			var field_v = r.Objects[i]
			var field_rv = AssumeValidReflectValue(field_rt, field_v)
			buf.append(field_name, field_rv, ctx)
		}
		var struct_ = buf.collect()
		return json.Marshal(struct_)
	case typsys.Union:
		if ref == source.MakeRef("", T_Maybe) {
			var inner_rt = firstArgType(args)
			var inner_v, ok = UnwrapMaybe(obj)
			if ok {
				var inner_rv = AssumeValidReflectValue(inner_rt, inner_v)
				var inner = jsonMarshaler { inner_rv, ctx }
				return json.Marshal(inner)
			} else {
				return ([] byte)("null"), nil
			}
		} else {
			var u = (*obj).(Union)
			var index = u.Index
			var inner_v = u.Object
			var field = content.FieldList[index]
			var inner_rt = inflateFieldType(field, def, args)
			var inner_rv = AssumeValidReflectValue(inner_rt, inner_v)
			var struct_ = jsonMarshalerUnionStruct {
				Index:  u.Index,
				Object: jsonMarshaler { inner_rv, ctx },
			}
			return json.Marshal(struct_)
		}
	case typsys.Enum:
		var e = (*obj).(Enum)
		var index = int(e)
		var field = content.FieldList[index]
		var field_name = field.Name
		return json.Marshal(field_name)
	}
	return nil, errors.New(fmt.Sprintf(
		"serialization is unavailable for type %s",
		typsys.Describe(rt.Type()),
	))
}
type jsonUnmarshaler struct {
	type_    ReflectType
	value    Object
	context  SerializationContext
}
func (m *jsonUnmarshaler) UnmarshalJSON(b ([] byte)) error {
	if len(b) == 0 {
		return errors.New("empty content")
	}
	var rt = m.type_
	var v = &(m.value)
	var ctx = m.context
	var def, ref, args, exists = ctx.LookupType(rt)
	if !(exists) { panic("something went wrong") }
	switch content := def.Content.(type) {
	case typsys.NativeContent:
		if ref.Namespace == "" {
			switch ref.ItemName {
			case T_Null:
				return nil
			case T_Bool:
				return jsonUnmarshalPrimitive[bool] (b,v)
			case T_Int:
				return jsonUnmarshalPrimitive[*big.Int] (b,v)
			case T_Float:
				return jsonUnmarshalPrimitive[float64] (b,v)
			case T_String:
				return jsonUnmarshalPrimitive[string] (b,v)
			case T_Bytes:
				return jsonUnmarshalPrimitive[[]byte] (b,v)
			case T_List:
				var segments ([] json.RawMessage)
				var err = json.Unmarshal(b, &segments)
				if err != nil { return err }
				var nodes = make([] ListNode, len(segments))
				var item_rt = firstArgType(args)
				for i, s := range segments {
					var item = &(nodes[i].Value)
					var err = jsonUnmarshalObject(s, item, item_rt, ctx)
					if err != nil { return err }
				}
				var l = NodesToList(nodes)
				*v = Obj(l)
				return nil
			}
		}
	case typsys.Record:
		var buf jsonUnmarshalerRecordStructBuilder
		for _, field := range content.FieldList {
			var field_name = field.Name
			var field_rt = inflateFieldType(field, def, args)
			buf.append(field_name, field_rt, ctx)
		}
		var objects, err = buf.collect(b)
		if err != nil { return err }
		*v = Obj(Record { objects })
		return nil
	case typsys.Union:
		if ref == source.MakeRef("", T_Maybe) {
			if bytes.Equal(b, ([] byte)("null")) {
				*v = Nothing()
				return nil
			} else {
				var inner_rt = firstArgType(args)
				var inner = new(Object)
				var err = jsonUnmarshalObject(b, inner, inner_rt, ctx)
				if err != nil { return err }
				*v = Just(*inner)
				return nil
			}
		} else {
			var struct_ = jsonUnmarshalerUnionStruct { Index: -1 }
			var err1 = json.Unmarshal(b, &struct_)
			if err1 != nil { return err1 }
			var index = struct_.Index
			var inner_raw = struct_.Object
			var L = len(content.FieldList)
			if !(0 <= index && index < L) { return errors.New("invalid union") }
			var field = content.FieldList[index]
			var inner_rt = inflateFieldType(field, def, args)
			var inner = new(Object)
			var err2 = jsonUnmarshalObject(inner_raw, inner, inner_rt, ctx)
			if err2 != nil { return err2 }
			*v = Obj(Union { Index: index, Object: *inner })
			return nil
		}
	case typsys.Enum:
		var field_name string
		var err = json.Unmarshal(b, &field_name)
		if err != nil { return err }
		var index = content.FieldIndexMap[field_name]
		var e = Enum(index)
		*v = Obj(e)
		return nil
	}
	return errors.New(fmt.Sprintf(
		"serialization is unavailable for type %s",
		typsys.Describe(rt.Type()),
	))
}
func firstArgType(args ([] typsys.Type)) ReflectType {
	var t_ = args[0]
	var t = typsys.CertainType { Type: t_ }
	var rt = MakeReflectType(t)
	return rt
}
func inflateFieldType(field typsys.Field, def *typsys.TypeDef, args ([] typsys.Type)) ReflectType {
	var t_ = typsys.Inflate(field.Type, def.Parameters, args)
	var t = typsys.CertainType { Type: t_ }
	var rt = MakeReflectType(t)
	return rt
}
func jsonUnmarshalPrimitive[T any] (b ([] byte), v *Object) error {
	var temp T
	var err = json.Unmarshal(b, &temp)
	if err != nil { return err }
	*v = ToObject(temp)
	return nil
}
func jsonUnmarshalObject(b ([] byte), v *Object, rt ReflectType, ctx SerializationContext) error {
	var m = jsonUnmarshaler { rt, nil, ctx }
	var err = json.Unmarshal(b, &m)
	if err != nil { return err }
	*v = m.value
	return nil
}
type jsonMarshalerUnionStruct struct {
	Index   int
	Object  jsonMarshaler
}
type jsonMarshalerRecordStructBuilder struct {
	fields  [] reflect.StructField
	values  [] jsonMarshaler
}
func (buf *jsonMarshalerRecordStructBuilder) append(name string, rv ReflectValue, ctx SerializationContext) {
	var field = jsonStructField (
		name, reflect.TypeOf(jsonMarshaler {}), len(buf.fields),
	)
	var value = jsonMarshaler {
		rv, ctx,
	}
	buf.fields = append(buf.fields, field)
	buf.values = append(buf.values, value)
}
func (buf *jsonMarshalerRecordStructBuilder) collect() interface{} {
	var ptr_rv = reflect.New(reflect.StructOf(buf.fields))
	var rv = ptr_rv.Elem()
	for i := 0; i < len(buf.fields); i += 1 {
		var m = buf.values[i]
		rv.Field(i).Set(reflect.ValueOf(m))
	}
	var v = rv.Interface()
	return v
}
type jsonUnmarshalerUnionStruct struct {
	Index   int
	Object  json.RawMessage
}
type jsonUnmarshalerRecordStructBuilder struct {
	fields  [] reflect.StructField
	values  [] jsonUnmarshaler
}
func (buf *jsonUnmarshalerRecordStructBuilder) append(name string, type_ ReflectType, ctx SerializationContext) {
	var field = jsonStructField (
		name, reflect.TypeOf(jsonUnmarshaler {}), len(buf.fields),
	)
	var value = jsonUnmarshaler {
		type_, nil, ctx,
	}
	buf.fields = append(buf.fields, field)
	buf.values = append(buf.values, value)
}
func (buf *jsonUnmarshalerRecordStructBuilder) collect(b ([] byte)) ([] Object, error) {
	var L = len(buf.fields)
	var ptr_rv = reflect.New(reflect.StructOf(buf.fields))
	var rv = ptr_rv.Elem()
	for i := 0; i < L; i += 1 {
		var m = buf.values[i]
		rv.Field(i).Set(reflect.ValueOf(m))
	}
	var ptr = ptr_rv.Interface()
	var err = json.Unmarshal(b, ptr)
	if err != nil { return nil, err }
	var objects = make([] Object, L)
	for i := 0; i < L; i += 1 {
		var field_value = rv.Field(i).Interface().(jsonUnmarshaler).value
		if field_value == nil {
			var m = buf.values[i]
			var field_type = m.type_
			var _, ref, _, exists = m.context.LookupType(field_type)
			if !(exists) { panic("something went wrong") }
			if ref.Namespace == "" {
				switch ref.ItemName {
				case T_Null:
					goto OK
				case T_Maybe:
					field_value = Nothing()
					goto OK
				}
			}
			var tag = rv.Type().Field(i).Tag
			var field_name = strings.TrimPrefix(string(tag), "json:")
			return nil, errors.New("missing field: " + field_name)
		}
		OK:
		objects[i] = field_value
	}
	return objects, nil
}
func jsonStructField(name string, type_ reflect.Type, n int) reflect.StructField {
	return reflect.StructField {
		Name: fmt.Sprintf("Export_%d", n),
		Tag:  reflect.StructTag(fmt.Sprintf("json:%s", strconv.Quote(name))),
		Type: type_,
	}
}


