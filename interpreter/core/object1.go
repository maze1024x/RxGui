package core

import (
	"time"
	"regexp"
	"math/big"
	"rxgui/util/ctn"
)


type Object =
	*ObjectImpl
//
type ObjectImpl interface {
	impl(*ObjectImpl)
}
func Obj(o ObjectImpl) Object {
	return &o
}


// Reflect
func (ReflectType) impl(Object) {}
func (ReflectValue) impl(Object) {}


// Primitive

func (Bool) impl(Object) {}
func (Int) impl(Object) {}
func (Float) impl(Object) {}
func (Bytes) impl(Object) {}
func (String) impl(Object) {}
func (Char) impl(Object) {}
func (RegExp) impl(Object) {}
func (Enum) impl(Object) {}
func (Time) impl(Object) {}
func (File) impl(Object) {}
func (Error) impl(Object) {}

type Bool bool
type Int struct { Value *big.Int }
type Float float64
type Bytes ([] byte)
type String string
type Char rune
type RegExp struct { Value *regexp.Regexp }
type Enum int
type Time time.Time
type File struct { Path string }
type Error struct { Value error }


// Rx

func (Observable) impl(Object) {}
func (Subject) impl(Object) {}


// Interface and Lambda

func (Interface) impl(Object) {}
func (Lambda) impl(Object) {}

type Interface struct {
	UnderlyingObject  Object
	DispatchTable     *DispatchTable
}
type Lambda struct { Call func(Object)(Object) }


// Container

func (List) impl(Object) {}
func (Seq) impl(Object) {}

func (Queue) impl(Object) {}
func (Heap) impl(Object) {}
func (Set) impl(Object) {}
func (Map) impl(Object) {}

type Queue ctn.Queue[Object]
type Heap ctn.Heap[Object]
type Set ctn.Set[Object]
type Map ctn.Map[Object,Object]


// Algebraic

func (Record) impl(Object) {}
func (Union) impl(Object) {}

type Record struct {
	Objects  [] Object
}
type Union struct {
	Index   int
	Object  Object
}


// GUI

func (Action) impl(Object) {}
func (Widget) impl(Object) {}
func (Signal) impl(Object) {}
func (Events) impl(Object) {}
func (Prop) impl(Object) {}


