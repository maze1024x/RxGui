package compiler

import (
    "rxgui/interpreter/lang/typsys"
    "rxgui/interpreter/program"
)


func getTypeDefArgs(t typsys.Type, ctx *exprContext) (*typsys.TypeDef, ([] typsys.Type), bool) {
    if ref_type, ok := t.(typsys.RefType); ok {
        var def, exists = ctx.context.lookupType(ref_type.Def)
        if !(exists) { panic("something went wrong") }
        var args = ref_type.Args
        return def, args, true
    } else {
        return nil, nil, false
    }
}

func inflateFieldType(field typsys.Field, def *typsys.TypeDef, args ([] typsys.Type)) typsys.Type {
    var raw_type = field.Type
    var params = def.Parameters
    var inflated_type = typsys.Inflate(raw_type, params, args)
    return inflated_type
}
func haveIdenticalTypeParams(a *typsys.TypeDef, b *typsys.TypeDef) bool {
    if len(a.Parameters) != len(b.Parameters) {
        return false
    }
    var n = len(a.Parameters)
    for i := 0; i < n; i += 1 {
        if a.Parameters[i] != b.Parameters[i] {
            return false
        }
    }
    return true
}

func getRecord(t typsys.Type, ctx *exprContext) (typsys.Record, *typsys.TypeDef, ([] typsys.Type), bool) {
    { var def, args, ok = getTypeDefArgs(t, ctx)
    if !(ok) { goto NG }
    { var record, ok = def.Content.(typsys.Record)
    if !(ok) { goto NG }
    return record, def, args, true } }
    NG:
    return typsys.Record {}, nil, nil, false
}
func getUnion(t typsys.Type, ctx *exprContext) (typsys.Union, *typsys.TypeDef, ([] typsys.Type), bool) {
    var def, args, ok = getTypeDefArgs(t, ctx)
    if !(ok) { goto NG }
    { var union, ok = def.Content.(typsys.Union)
    if !(ok) { goto NG }
    return union, def, args, true }
    NG:
    return typsys.Union {}, nil, nil, false
}
func getEnum(t typsys.Type, ctx *exprContext) (typsys.Enum, *typsys.TypeDef, bool) {
    var def, _, ok = getTypeDefArgs(t, ctx)
    if !(ok) { goto NG }
    { var enum, ok = def.Content.(typsys.Enum)
    if !(ok) { goto NG }
    return enum, def, true }
    NG:
    return typsys.Enum {}, nil, false
}
func getInterface(t typsys.Type, ctx *exprContext) (typsys.Interface, *typsys.TypeDef, ([] typsys.Type), bool) {
    var def, args, ok = getTypeDefArgs(t, ctx)
    if !(ok) { goto NG }
    { var interface_, ok = def.Content.(typsys.Interface)
    if !(ok) { goto NG }
    return interface_, def, args, true }
    NG:
    return typsys.Interface {}, nil, nil, false
}
func isSamInterface(I typsys.Interface, Id *typsys.TypeDef) bool {
    return (len(Id.Interfaces) == 0) && (len(I.FieldList) == 1)
}
func getSamInterfaceMethodType(t typsys.Type, ctx *exprContext) (typsys.Type, bool) {
    // note: behavior of this function should be in sync with assign()
    var I, Id, Ia, is_interface = getInterface(t, ctx)
    if is_interface {
        if isSamInterface(I, Id) {
            var field = I.FieldList[0]
            var field_type = inflateFieldType(field, Id, Ia)
            return field_type, true
        }
    }
    return nil, false
}
func getInterfacePath(root *typsys.TypeDef, descendant *typsys.TypeDef, ctx *exprContext, base ([] int)) ([] int, bool) {
    for i, child_ref := range root.Interfaces {
        var child, exists = ctx.context.lookupType(child_ref)
        if !(exists) { panic("something went wrong") }
        var current = append(base, i)
        if child == descendant {
            return current, true
        } else {
            var path, ok = getInterfacePath(child, descendant, ctx, current)
            if ok {
                return path, true
            }
        }
    }
    return nil, false
}

func getInteriorReferableRecord(base_t typsys.CertainType, ctx *exprContext) (typsys.Record, *typsys.TypeDef, ([] typsys.Type), program.InteriorRefOperand, typsys.CertainType, bool) {
    if record, def, args, ok := getRecord(base_t.Type, ctx); ok {
        return record, def, args, program.RO_Direct, base_t, true
    }
    if base_t_, field_t, ok := program.T_Lens1_(base_t.Type); ok {
        var base_t = typsys.CertainType { Type: base_t_ }
        if record, def, args, ok := getRecord(field_t, ctx); ok {
            return record, def, args, program.RO_Lens1, base_t, true
        }
    }
    return typsys.Record{}, nil, nil, -1, typsys.CertainType{}, false
}
func getInteriorReferableEnum(base_t typsys.CertainType, ctx *exprContext) (typsys.Enum, program.InteriorRefOperand, typsys.CertainType, bool) {
    if enum, _, ok := getEnum(base_t.Type, ctx); ok {
        return enum, program.RO_Direct, base_t, true
    }
    if base_t_, field_t, ok := program.T_Lens1_(base_t.Type); ok {
        var base_t = typsys.CertainType { Type: base_t_ }
        if enum, _, ok := getEnum(field_t, ctx); ok {
            return enum, program.RO_Lens1, base_t, true
        }
    }
    if base_t_, field_t, ok := program.T_Lens2_(base_t.Type); ok {
        var base_t = typsys.CertainType { Type: base_t_ }
        if enum, _, ok := getEnum(field_t, ctx); ok {
            return enum, program.RO_Lens2, base_t, true
        }
    }
    return typsys.Enum{}, -1, typsys.CertainType{}, false
}
func getInteriorReferableUnion(base_t typsys.CertainType, ctx *exprContext) (typsys.Union, *typsys.TypeDef, ([] typsys.Type), program.InteriorRefOperand, typsys.CertainType, bool) {
    if union, def, args, ok := getUnion(base_t.Type, ctx); ok {
        return union, def, args, program.RO_Direct, base_t, true
    }
    if base_t_, field_t, ok := program.T_Lens1_(base_t.Type); ok {
        var base_t = typsys.CertainType { Type: base_t_ }
        if union, def, args, ok := getUnion(field_t, ctx); ok {
            return union, def, args, program.RO_Lens1, base_t, true
        }
    }
    if base_t_, field_t, ok := program.T_Lens2_(base_t.Type); ok {
        var base_t = typsys.CertainType { Type: base_t_ }
        if union, def, args, ok := getUnion(field_t, ctx); ok {
            return union, def, args, program.RO_Lens2, base_t, true
        }
    }
    return typsys.Union{}, nil, nil, -1, typsys.CertainType{}, false
}
func getInteriorReferableInterface(base_t typsys.CertainType, ctx *exprContext) (typsys.Interface, *typsys.TypeDef, ([] typsys.Type), program.InteriorRefOperand, typsys.CertainType, bool) {
    if interface_, def, args, ok := getInterface(base_t.Type, ctx); ok {
        return interface_, def, args, program.RO_Direct, base_t, true
    }
    if base_t_, field_t, ok := program.T_Lens1_(base_t.Type); ok {
        var base_t = typsys.CertainType { Type: base_t_ }
        if interface_, def, args, ok := getInterface(field_t, ctx); ok {
            return interface_, def, args, program.RO_Lens1, base_t, true
        }
    }
    if base_t_, field_t, ok := program.T_Lens2_(base_t.Type); ok {
        var base_t = typsys.CertainType { Type: base_t_ }
        if interface_, def, args, ok := getInterface(field_t, ctx); ok {
            return interface_, def, args, program.RO_Lens2, base_t, true
        }
    }
    return typsys.Interface{}, nil, nil, -1, typsys.CertainType{}, false
}


