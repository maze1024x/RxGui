package typsys

import "rxgui/lang/source"


type Type interface { impl(Type) }
type CertainType struct { Type Type }

func (InferringType) impl(Type) {}
type InferringType struct {
	Id  string
}

func (ParameterType) impl(Type) {}
type ParameterType struct {
	Name  string
}

func (RefType) impl(Type) {}
type RefType struct {
	Def   source.Ref
	Args  [] Type
}

func DefType(def *TypeDef) Type {
	return DefType2(def.Ref, def.Parameters)
}
func DefType2(ref source.Ref, params ([] string)) Type {
	return RefType {
		Def:  ref,
		Args: (func() ([] Type) {
			var args = make([] Type, len(params))
			for i, p := range params {
				args[i] = ParameterType { p }
			}
			return args
		})(),
	}
}


