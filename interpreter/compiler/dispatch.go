package compiler

import (
	"fmt"
	"rxgui/standalone/ctn"
	"rxgui/lang/source"
	"rxgui/lang/typsys"
	"rxgui/interpreter/program"
)


func fillDispatchInfo(fd fragmentDraft, hdr *Header, ctx *NsHeaderMap, errs *source.Errors) {
	hdr.typeMap.ForEach(func(_ string, C *typsys.TypeDef) {
		var _, C_is_interface = C.Content.(typsys.Interface)
		for _, I_ := range C.Interfaces {
			var I, exists = ctx.lookupType(I_)
			if !(exists) { panic("something went wrong") }
			if C_is_interface {
				// do nothing here (case handled elsewhere)
			} else {
				var _, err = createDispatchTable(C, I, fd, ctx)
				source.ErrorsJoin(errs, err)
			}
		}
	})
}
func createDispatchTable (
	C    *typsys.TypeDef,
	I    *typsys.TypeDef,
	fd   fragmentDraft,
	ctx  *NsHeaderMap,
) (*program.DispatchTable, *source.Error) {
	var loc = C.Info.Location
	var spec, ok = I.Content.(typsys.Interface)
	if !(ok) { panic("something went wrong") }
	var pair = dispatchKey {
		ConcreteType:  C.Ref.ItemName,
		InterfaceType: I.Ref,
	}
	{ var t, ok = fd.createDispatchTable(pair)
	if !(ok) {
		return nil, source.MakeError(loc,
			E_DuplicateInterface { I.Ref.String() })
	}
	if !(haveIdenticalTypeParams(C, I)) {
		return nil, source.MakeError(loc,
			E_TypeParamsNotIdentical {
				Concrete:  C.Ref.String(),
				Interface: I.Ref.String(),
			})
	}
	var methods = make([] *program.Function, len(spec.FieldList))
	for i, field := range spec.FieldList {
		var method_name = field.Name
		var expected_type = field.Type
		var f, ok = ctx.lookupMethod(C.Ref, method_name)
		if !(ok) {
			if record, is_record := C.Content.(typsys.Record); is_record {
			if index, found := record.FieldIndexMap[method_name]; found {
			if typsys.Equal(record.FieldList[index].Type, expected_type) {
				var getter = new(program.Function)
				getter.SetName(fmt.Sprintf("[field_%d_%p]", index, getter))
				getter.SetFieldValueGetterValueByIndex(index)
				methods[i] = getter
				continue
			}}}
			return nil, source.MakeError(loc,
				E_MissingMethod {
					Concrete:  C.Ref.String(),
					Interface: I.Ref.String(),
					Method:    method_name,
				})
		}
		var actual_type = f.output.type_
		if !(typsys.Equal(actual_type, expected_type)) {
			return nil, source.MakeError(loc,
				E_WrongMethodType {
					Concrete:  C.Ref.String(),
					Interface: I.Ref.String(),
					Method:    method_name,
					Expected:  typsys.Describe(expected_type),
					Actual:    typsys.Describe(actual_type),
				})
		}
		var key = userFunKey {
			name:  method_name,
			assoc: C.Ref.ItemName,
		}
		methods[i] = fd.createOrGetUserFun(key)
	}
	var children = make([] *program.DispatchTable, len(I.Interfaces))
	for i, II_ := range I.Interfaces {
		var II, exists = ctx.lookupType(II_)
		if !(exists) { panic("something went wrong") }
		var child, err = createDispatchTable(C, II, fd, ctx)
		if err != nil { return nil, err }
		child.SetParent(t)
		children[i] = child
	}
	t.SetName(pair.Describe(fd.namespace()))
	t.SetInterface(I.Ref.String())
	t.SetMethods(methods)
	t.SetChildren(children)
	return t, nil }
}

type dispatchKey struct {
	ConcreteType   string
	InterfaceType  source.Ref
}
func (pair dispatchKey) Describe(ns string) string {
	var C = source.MakeRef(ns, pair.ConcreteType)
	var I = pair.InterfaceType
	return fmt.Sprintf("(%s,%s)", C, I)
}
func dispatchKeyCompare(a dispatchKey, b dispatchKey) ctn.Ordering {
	var o = ctn.StringCompare(a.ConcreteType, b.ConcreteType)
	if o != ctn.Equal {
		return o
	} else {
		return source.RefCompare(a.InterfaceType, b.InterfaceType)
	}
}


