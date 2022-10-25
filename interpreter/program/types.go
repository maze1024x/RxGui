package program

import (
    "rxgui/util/ctn"
    "rxgui/interpreter/core"
    "rxgui/interpreter/lang/source"
    "rxgui/interpreter/lang/typsys"
    "rxgui/interpreter/lang/textual/ast"
)


const T_Observable_Tag = core.T_Observable
const T_Observable_Maybe_Tag = (core.T_Observable + "?")
const T_Hook_Tag = core.T_Hook

type AstType func(ast.Node)(ast.Type)
func AT_Null() AstType {
    return makeAstType(core.T_Null)
}
func AT_Observable(item_t AstType) AstType {
    return makeAstType(core.T_Observable, item_t)
}
func makeAstType(item string, args ...AstType) AstType {
    return func(node ast.Node) ast.Type {
        return ast.Type {
            Node: node,
            Ref:  ast.Ref {
                Node:     node,
                Base:     ast.RefBase {
                    Node: node,
                    Item: ast.String2Id(item, node),
                },
                TypeArgs: ctn.MapEach(args, func(arg AstType) ast.Type {
                    return arg(node)
                }),
            },
        }
    }
}

func T_InvalidType() typsys.Type {
    return makeType(core.T_InvalidType)
}
func T_InvalidType_(t typsys.Type) bool {
    return matchType0(core.T_InvalidType, t)
}
func T_IntoReflectType(t typsys.Type) typsys.Type {
    return makeType(core.T_IntoReflectType, t)
}
func T_IntoReflectType_(t typsys.Type) (typsys.Type, bool) {
    return matchType1(core.T_IntoReflectType, t)
}
func T_IntoReflectValue(t typsys.Type) typsys.Type {
    return makeType(core.T_IntoReflectValue, t)
}
func T_IntoReflectValue_(t typsys.Type) (typsys.Type, bool) {
    return matchType1(core.T_IntoReflectValue, t)
}
func T_DummyReflectType_(t typsys.Type) bool {
    return matchType0(core.T_DummyReflectType, t)
}
func T_Null() typsys.Type {
    return makeType(core.T_Null)
}
func T_Null_(t typsys.Type) bool {
    return matchType0(core.T_Null, t)
}
func T_Bool_(t typsys.Type) bool {
    return matchType0(core.T_Bool, t)
}
func T_Maybe(t typsys.Type) typsys.Type {
    return makeType(core.T_Maybe, t)
}
func T_Maybe_(t typsys.Type) (typsys.Type, bool) {
    return matchType1(core.T_Maybe, t)
}
func T_Int() typsys.Type {
    return makeType(core.T_Int)
}
func T_Int_(t typsys.Type) bool {
    return matchType0(core.T_Int, t)
}
func T_Float() typsys.Type {
    return makeType(core.T_Float)
}
func T_Char() typsys.Type {
    return makeType(core.T_Char)
}
func T_Bytes() typsys.Type {
    return makeType(core.T_Bytes)
}
func T_String() typsys.Type {
    return makeType(core.T_String)
}
func T_RegExp() typsys.Type {
    return makeType(core.T_RegExp)
}
func T_Lambda(in typsys.Type, out typsys.Type) typsys.Type {
    return makeType(core.T_Lambda, in, out)
}
func T_Lambda_(t typsys.Type) (typsys.Type, typsys.Type, bool) {
    return matchType2(core.T_Lambda, t)
}
func T_List(t typsys.Type) typsys.Type {
    return makeType(core.T_List, t)
}
func T_List_(t typsys.Type) (typsys.Type, bool) {
    return matchType1(core.T_List, t)
}
func T_Pair(a typsys.Type, b typsys.Type) typsys.Type {
    return makeType(core.T_Pair, a, b)
}
func T_Pair_(t typsys.Type) (typsys.Type, typsys.Type, bool) {
    return matchType2(core.T_Pair, t)
}
func T_Triple(a typsys.Type, b typsys.Type, c typsys.Type) typsys.Type {
    return makeType(core.T_Triple, a, b, c)
}
func T_Triple_(t typsys.Type) (typsys.Type, typsys.Type, typsys.Type, bool) {
    return matchType3(core.T_Triple, t)
}
func T_Lens1(a typsys.Type, b typsys.Type) typsys.Type {
    return makeType(core.T_Lens1, a, b)
}
func T_Lens1_(t typsys.Type) (typsys.Type, typsys.Type, bool) {
    return matchType2(core.T_Lens1, t)
}
func T_Lens2(a typsys.Type, b typsys.Type) typsys.Type {
    return makeType(core.T_Lens2, a, b)
}
func T_Lens2_(t typsys.Type) (typsys.Type, typsys.Type, bool) {
    return matchType2(core.T_Lens2, t)
}
func T_Observable(t typsys.Type) typsys.Type {
    return makeType(core.T_Observable, t)
}
func T_Observable_(t typsys.Type) (typsys.Type, bool) {
    return matchType1(core.T_Observable, t)
}
func T_Hook(t typsys.Type) typsys.Type {
    return makeType(core.T_Hook, t)
}
func makeType(item string, args ...typsys.Type) typsys.Type {
    return typsys.RefType {
        Def:  source.MakeRef("", item),
        Args: args,
    }
}
// TODO: makeType(variadic->slice), makeType0, makeType1, makeType2, ...
func matchType(item string, t typsys.Type) ([] typsys.Type, bool) {
    if ref_type, ok := t.(typsys.RefType); ok {
        var ref = ref_type.Def
        var args = ref_type.Args
        if (ref.Namespace == "") && (ref.ItemName == item) {
            return args, true
        }
    }
    return nil, false
}
func matchType0(item string, t typsys.Type) bool {
    var args, ok = matchType(item, t)
    if ok {
        if len(args) != 0 { panic("something went wrong") }
        return true
    } else {
        return false
    }
}
func matchType1(item string, t typsys.Type) (typsys.Type, bool) {
    var args, ok = matchType(item, t)
    if ok {
        if len(args) != 1 { panic("something went wrong") }
        return args[0], true
    } else {
        return nil, false
    }
}
func matchType2(item string, t typsys.Type) (typsys.Type, typsys.Type, bool) {
    var args, ok = matchType(item, t)
    if ok {
        if len(args) != 2 { panic("something went wrong") }
        return args[0], args[1], true
    } else {
        return nil, nil, false
    }
}
func matchType3(item string, t typsys.Type) (typsys.Type, typsys.Type, typsys.Type, bool) {
    var args, ok = matchType(item, t)
    if ok {
        if len(args) != 3 { panic("something went wrong") }
        return args[0], args[1], args[2], true
    } else {
        return nil, nil, nil, false
    }
}

type ReflectType_ core.ReflectType
func MakeReflectType_(t typsys.CertainType) ReflectType_ {
    return ReflectType_(core.MakeReflectType(t))
}


