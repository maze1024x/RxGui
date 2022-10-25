package program

import (
    "rxgui/util/ctn"
    "rxgui/interpreter/core"
)


type Function struct {
    name   string
    value  core.Function
}
func (f *Function) SetName(name string) {
    f.name = name
}
func (f *Function) SetExprBasedValue(value *ExprBasedFunctionValue) {
    f.value = value
}
func (f *Function) SetFieldValueGetterValueByIndex(i int) {
    f.value = core.FieldValueGetterFunction { Index: i }
}
func (f *Function) SetNativeValueById(id string, const_ bool) {
    if const_ {
        var cache = core.Object(nil)
        var available = false
        f.value = core.NativeFunction(func(args ([] core.Object), ctx ([] core.Object), h core.RuntimeHandle) core.Object {
            if available {
                return cache
            } else {
                if (len(args)|len(ctx) > 0) { panic("something went wrong") }
                cache = h.LibraryNativeFunction(id)(args, ctx, h)
                available = true
                return cache
            }
        })
    } else {
        f.value = core.NativeFunction(func(args ([] core.Object), ctx ([] core.Object), h core.RuntimeHandle) core.Object {
            return h.LibraryNativeFunction(id)(args, ctx, h)
        })
    }
}

type DispatchTable struct {
    name   string
    table  core.DispatchTable
}
func (t *DispatchTable) value() *core.DispatchTable {
    return &(t.table)
}
func (t *DispatchTable) SetName(name string) {
    t.name = name
}
func (t *DispatchTable) SetInterface(id string) {
    t.table.Interface = id
}
func (t *DispatchTable) SetMethods(methods ([] *Function)) {
    t.table.Methods = ctn.MapEach(methods, func(m *Function) *core.Function {
        return &(m.value)
    })
}
func (t *DispatchTable) SetChildren(children ([] *DispatchTable)) {
    t.table.Children = ctn.MapEach(children, func(c *DispatchTable) *core.DispatchTable {
        return &(c.table)
    })
}
func (t *DispatchTable) SetParent(parent *DispatchTable) {
    t.table.Parent = &(parent.table)
}


