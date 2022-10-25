package core

import "rxgui/util/ctn"


func Lens1FromRecord(object Object, index int) Object {
	return ToObject(lens1fromRecord(object, index))
}
func Lens2FromEnum(object Object, index int) Object {
	return ToObject(lens2fromEnum(ctn.Just(object), index))
}
func Lens2FromUnion(object Object, index int) Object {
	return ToObject(lens2fromUnion(ctn.Just(object), index))
}
func Lens2FromInterface(object Object, t *DispatchTable) Object {
	return ToObject(lens2fromInterface(ctn.Just(object), t))
}
func Lens1FromRecordLens1(object Object, index int) Object {
	var ab = FromObject[Lens1](object)
	return ToObject(lens1compose1(ab, func(b Object) Lens1 {
		return lens1fromRecord(b, index)
	}))
}
func Lens2FromEnumLens1(object Object, index int) Object {
	var ab = FromObject[Lens1](object)
	return ToObject(lens1compose2(ab, func(b Object) Lens2 {
		return lens2fromEnum(ctn.Just(b), index)
	}))
}
func Lens2FromUnionLens1(object Object, index int) Object {
	var ab = FromObject[Lens1](object)
	return ToObject(lens1compose2(ab, func(b Object) Lens2 {
		return lens2fromUnion(ctn.Just(b), index)
	}))
}
func Lens2FromEnumLens2(object Object, index int) Object {
	var ab = FromObject[Lens2](object)
	return ToObject(lens2compose(ab, func(opt_b ctn.Maybe[Object]) Lens2 {
		return lens2fromEnum(opt_b, index)
	}))
}
func Lens2FromUnionLens2(object Object, index int) Object {
	var ab = FromObject[Lens2](object)
	return ToObject(lens2compose(ab, func(opt_b ctn.Maybe[Object]) Lens2 {
		return lens2fromUnion(opt_b, index)
	}))
}
func Lens2FromInterfaceLens1(object Object, t *DispatchTable) Object {
	var ab = FromObject[Lens1](object)
	return ToObject(lens1compose2(ab, func(b Object) Lens2 {
		return lens2fromInterface(ctn.Just(b), t)
	}))
}
func Lens2FromInterfaceLens2(object Object, t *DispatchTable) Object {
	var ab = FromObject[Lens2](object)
	return ToObject(lens2compose(ab, func(opt_b ctn.Maybe[Object]) Lens2 {
		return lens2fromInterface(opt_b, t)
	}))
}

func lens1fromRecord(object Object, index int) Lens1 {
	var r = (*object).(Record)
	var objects = r.Objects
	if !(index < len(objects)) {
		panic("invalid argument")
	}
	return Lens1 {
		Value:  objects[index],
		Assign: func(new_object Object) Object {
			var new_objects = make([] Object, len(objects))
			copy(new_objects, objects)
			new_objects[index] = new_object
			var o = ObjectImpl(Record { new_objects })
			return &o
		},
	}
}
func lens2fromEnum(opt_object ctn.Maybe[Object], index int) Lens2 {
	var value = ctn.Nothing[Object]()
	if object, ok := opt_object.Value(); ok {
		var u = (*object).(Enum)
		if int(u) == index {
			value = ctn.Just[Object](nil)
		}
	}
	return Lens2 {
		Value:  value,
		Assign: func(_ Object) Object {
			var o = ObjectImpl(Enum(index))
			return &o
		},
	}
}
func lens2fromUnion(opt_object ctn.Maybe[Object], index int) Lens2 {
	var value = ctn.Nothing[Object]()
	if object, ok := opt_object.Value(); ok {
		var u = (*object).(Union)
		if u.Index == index {
			value = ctn.Just(u.Object)
		}
	}
	return Lens2 {
		Value:  value,
		Assign: func(new_object Object) Object {
			var o = ObjectImpl(Union {
				Index:  index,
				Object: new_object,
			})
			return &o
		},
	}
}
func lens2fromInterface(opt_object ctn.Maybe[Object], t *DispatchTable) Lens2 {
	var value = ctn.Nothing[Object]()
	if object, ok := opt_object.Value(); ok {
		var u = (*object).(Interface)
		if u.DispatchTable == t {
			value = ctn.Just(u.UnderlyingObject)
		}
	}
	return Lens2 {
		Value:  value,
		Assign: func(new_object Object) Object {
			var o = ObjectImpl(Interface {
				UnderlyingObject: new_object,
				DispatchTable:    t,
			})
			return &o
		},
	}
}
func lens1compose1(ab Lens1, f func(Object)(Lens1)) Lens1 {
	return Lens1 {
		Value:  f(ab.Value).Value,
		Assign: func(c Object) Object {
			return ab.Assign(f(ab.Value).Assign(c))
		},
	}
}
func lens1compose2(ab Lens1, f func(Object)(Lens2)) Lens2 {
	return Lens2 {
		Value:  f(ab.Value).Value,
		Assign: func(c Object) Object {
			return ab.Assign(f(ab.Value).Assign(c))
		},
	}
}
func lens2compose(ab Lens2, f func(ctn.Maybe[Object])(Lens2)) Lens2 {
	return Lens2 {
		Value:  f(ab.Value).Value,
		Assign: func(c Object) Object {
			return ab.Assign(f(ab.Value).Assign(c))
		},
	}
}


