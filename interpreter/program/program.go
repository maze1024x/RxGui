package program

import (
	"rxgui/lang/source"
	"rxgui/lang/typsys"
	"rxgui/interpreter/core"
)


type Program struct {
	Metadata    Metadata
	TypeInfo    TypeInfo
	Executable  Executable
}
type Metadata struct {
	ProgramPath  string
}
type Executable interface {
	LookupEntry(ns string) (**Function, bool)
}
type TypeInfo struct {
	TypeRegistry  map[source.Ref] *typsys.TypeDef
}
func (rtti TypeInfo) LookupType(rt core.ReflectType) (*typsys.TypeDef, source.Ref, ([] typsys.Type), bool) {
	var t = rt.Type()
	if T, ok := t.(typsys.RefType); ok {
	if def, ok := rtti.TypeRegistry[T.Def]; ok {
		return def, T.Def, T.Args, true
	}}
	return nil, source.Ref{}, nil, false
}


