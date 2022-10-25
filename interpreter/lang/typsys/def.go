package typsys

import "rxgui/interpreter/lang/source"


type TypeDef struct {
	Info        Info
	Ref         source.Ref
	Interfaces  [] source.Ref
	Parameters  [] string
	Content     TypeDefContent
}
type Fields struct {
	FieldIndexMap  map[string] int
	FieldList      [] Field
}
type Field struct {
	Info  FieldInfo
	Name  string
	Type  Type  // nil in enum
}
type Info struct {
	Location  source.Location
	Document  string
}
type FieldInfo struct {
	Info
	HasDefaultValue  bool
}
type TypeDefContent interface { impl(TypeDefContent) }

func (NativeContent) impl(TypeDefContent) {}
type NativeContent struct {}

func (Record) impl(TypeDefContent) {}
type Record struct {
	*Fields
}

func (Interface) impl(TypeDefContent) {}
type Interface struct {
	*Fields
}

func (Union) impl(TypeDefContent) {}
type Union struct {
	*Fields
}

func (Enum) impl(TypeDefContent) {}
type Enum struct {
	*Fields
}


