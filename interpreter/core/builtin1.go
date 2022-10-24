package core

import (
	"rxgui/standalone/ctn"
	"rxgui/standalone/qt"
	"rxgui/standalone/util/pseudounion"
)


const NgIndex = 0
const OkIndex = 1

const T_InvalidType = "<InvalidType>"
const T_GenericType = "<GenericType>"
const T_IntoReflectType = "IntoReflectType"
const T_IntoReflectValue = "IntoReflectValue"
const T_DummyReflectType = "DummyReflectType"
const T_Null = "Null"
const T_Error = "Error"
const T_Bool = "Bool"
const T_Int = "Int"
const T_Float = "Float"
const T_Char = "Char"
const T_String = "String"
const T_RegExp = "RegExp"
const T_Bytes = "Bytes"
const T_Asset = "Asset"
const T_List = "List"
const T_Seq = "Seq"
const T_Queue = "Queue"
const T_Heap = "Heap"
const T_Set = "Set"
const T_Map = "Map"
const T_Observable = "$"
const T_Time = "Time"
const T_File = "File"
const T_Lambda = "Lambda"
const T_Pair = "Pair"
const T_Triple = "Triple"
const T_Maybe = "Maybe"
const T_Lens1 = "Lens1"
const T_Lens2 = "Lens2"
const T_Hook = "Hook"

type Lens1 struct {
	Value   Object
	Assign  func(Object) Object
}
type Lens2 struct {
	Value   ctn.Maybe[Object]
	Assign  func(Object) Object
}
type Hook struct {
	Job  Observable
}
func Just(o Object) Object {
	return Obj(Union {
		Index:  OkIndex,
		Object: o,
	})
}
func Nothing() Object {
	return Obj(Union {
		Index:  NgIndex,
		Object: nil,
	})
}
func UnwrapMaybe(o Object) (Object, bool) {
	var u = (*o).(Union)
	if u.Index == OkIndex {
		return u.Object, true
	} else if u.Index == NgIndex {
		return nil, false
	} else {
		panic("something went wrong")
	}
}
func UnwrapLens2(o Object) (Object, bool) {
	var l = FromObject[Lens2](o)
	return l.Value.Value()
}
func MakeHook[T any] (k func()(T,func())) Hook {
	return Hook { Observable(func(pub DataPublisher) {
		var t, c = k()
		var v = ToObject(t)
		if c != nil {
			pub.context.registerCleaner(c)
		}
		pub.observer.value(v)
		pub.observer.complete()
	})}
}
func MakeHookWithEffect[T any] (h RuntimeHandle, k func()(T,Observable,func())) Hook {
	return Hook { ObservableFlattenLast(Observable(func(pub DataPublisher) {
		var t, e, c = k()
		var o = ObservableSyncValue(ToObject(t)).With(e, ErrorLogger(h))
		if c != nil {
			pub.context.registerCleaner(c)
		}
		pub.observer.value(Obj(o))
		pub.observer.complete()
	}))}
}
func MapHook[A any, B any] (h Hook, f func(A)(B)) Hook {
	return Hook { h.Job.Map(func(obj Object) Object {
		var a = FromObject[A](obj)
		var b = f(a)
		return ToObject(b)
	})}
}
func UseHook[A any] (h Hook, f func(A)(Observable)) Hook {
	return Hook { retrieveObject[A](h.Job, nil, func(a A) Observable {
		return f(a)
	})}
}

type Widgets struct {
	pseudounion.Tag
	Widget  Widget
	List    [] Widget
}
func (w Widgets) Value() ([] Widget) {
	switch W := pseudounion.Load(w).(type) {
	case Widget:
		return [] Widget { W }
	case [] Widget:
		return W
	default:
		panic("impossible branch")
	}
}
func (w Widgets) Deref(h RuntimeHandle) ([] qt.Widget) {
	return ctn.MapEach(w.Value(), func(w Widget) qt.Widget {
		return w.Deref(h)
	})
}

type ItemInfo struct {
	Key  string
	Pos  Observable
}
type ItemPos struct {
	Index  int
	Total  int
}
func DistinctUntilItemPosChanged(o Observable) Observable {
	return o.DistinctUntilChanged(func(a Object, b Object) bool {
		var u = FromObject[ItemPos](a)
		var v = FromObject[ItemPos](b)
		return (u == v)
	})
}

type ItemView struct {
	Widgets    Widgets
	Extension  ctn.Maybe[Widget]
}
type ItemViewProvider =
	func(Observable,ItemInfo)(Hook)
//
type ItemEditView struct {
	Widgets    Widgets
	Extension  ctn.Maybe[Widget]
	EditOps    func(key string)(Observable)
}
type ItemEditViewProvider =
	func(Object,ItemInfo)(Hook)
//

type ListEditOperation struct {
	pseudounion.Tag
	Prepend; Append
	Update; Delete
	MoveUp; MoveDown
	MoveTop; MoveBottom
	InsertAbove; InsertBelow
	Reorder
}
type Prepend struct { Value Object }
type Append struct { Value Object }
type Update struct { Key string; Value Object }
type Delete struct { Key ctn.Maybe[string] }
type MoveUp struct { Key ctn.Maybe[string] }
type MoveDown struct { Key ctn.Maybe[string] }
type MoveTop struct { Key ctn.Maybe[string] }
type MoveBottom struct { Key ctn.Maybe[string] }
type InsertAbove struct { PivotKey string; Value Object }
type InsertBelow struct { PivotKey string; Value Object }
type Reorder struct { Reorder func(List)(List) }


