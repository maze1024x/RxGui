package compiler

import (
    "rxgui/interpreter/lang/source"
    "rxgui/interpreter/lang/typsys"
    "rxgui/interpreter/lang/textual/ast"
    "rxgui/interpreter/program"
)


func checkCast(cast ast.Cast, in ast.Expr, cc *exprCheckContext) (*program.Expr, *source.Error) {
    var loc = cast.Location
    var target, err = cc.context.makeType(cast.Target)
    if err != nil { return nil, err }
    { var expr, err = cc.checkChildExpr(target.Type, in)
    if err != nil { return nil, err }
    return cc.assign(target, loc,
        program.Wrapper { Inner: expr })
    }
}

func checkPipeGet(get ast.PipeGet, in ast.Expr, cc *exprCheckContext) (*program.Expr, *source.Error) {
    var key = ast.Id2String(get.Key)
    var loc = get.Location
    var ctx = cc.context
    var in_expr, err = cc.checkChildExpr(nil, in)
    if err != nil { return nil, err }
    var in_t = in_expr.Type
    if record, def, args, ok := getRecord(in_t.Type, ctx); ok {
    if field_index, ok := record.FieldIndexMap[key]; ok {
        var field = record.FieldList[field_index]
        var field_t_ = inflateFieldType(field, def, args)
        var field_t = typsys.CertainType { Type: field_t_ }
        return cc.assign(field_t, loc,
            program.FieldValue {
                Record: in_expr,
                Index:  field_index,
            })
    }}
    if inner_t_, ok := program.T_Observable_(in_t.Type); ok {
    if record, def, args, ok := getRecord(inner_t_, ctx); ok {
    if field_index, ok := record.FieldIndexMap[key]; ok {
        var field = record.FieldList[field_index]
        var field_t_ = inflateFieldType(field, def, args)
        var field_ob_t_ = program.T_Observable(field_t_)
        var field_ob_t = typsys.CertainType { Type: field_ob_t_ }
        return cc.assign(field_ob_t, loc,
            program.ObservableFieldProjection {
                Base:  in_expr,
                Index: field_index,
            })
    }}}
    if method_t, f, index, path, ok := ctx.resolveMethod(in_t, key, nil); ok {
        if f != nil {
            return cc.assign(method_t, loc,
                program.ConcreteMethodValue {
                    Location: loc,
                    This:     in_expr,
                    Path:     path,
                    Method:   f,
                })
        } else {
            return cc.assign(method_t, loc,
                program.AbstractMethodValue {
                    Location:  loc,
                    Interface: in_expr,
                    Path:      path,
                    Index:     index,
                })
        }
    }
    return cc.error(loc,
        E_NoSuchFieldOrMethod {
            FieldName: key,
            TypeDesc:  typsys.DescribeCertain(in_t),
        })
}

func checkPipeInterior(interior ast.PipeInterior, in ast.Expr, cc *exprCheckContext) (*program.Expr, *source.Error) {
    var ref = getRef(interior.RefBase)
    var loc = interior.Location
    var ctx = cc.context
    var in_expr, err = cc.checkChildExpr(nil, in)
    if err != nil { return nil, err }
    var in_t = in_expr.Type
    if ref.Namespace == "" {
        var key = ref.ItemName
        if r, d, a, o, b, ok := getInteriorReferableRecord(in_t, ctx); ok {
            var index, exists = r.FieldIndexMap[key]
            if !(exists) { goto NG }
            var field = r.FieldList[index]
            var field_t_ = inflateFieldType(field, d, a)
            var lens1_t_ = program.T_Lens1(b.Type, field_t_)
            var lens1_t = typsys.CertainType { Type: lens1_t_}
            return cc.assign(lens1_t, loc,
                program.InteriorRef {
                    Base:    in_expr,
                    Index:   index,
                    Kind:    program.RK_RecordField,
                    Operand: o,
                })
        }
        if e, o, b, ok := getInteriorReferableEnum(in_t, ctx); ok {
            var index, exists = e.FieldIndexMap[key]
            if !(exists) { goto NG }
            var null_t = program.T_Null()
            var lens2_t_ = program.T_Lens2(b.Type, null_t)
            var lens2_t = typsys.CertainType { Type: lens2_t_ }
            return cc.assign(lens2_t, loc,
                program.InteriorRef {
                    Base:    in_expr,
                    Index:   index,
                    Kind:    program.RK_EnumItem,
                    Operand: o,
                })
        }
        if u, d, a, o, b, ok := getInteriorReferableUnion(in_t, ctx); ok {
            var index, exists = u.FieldIndexMap[key]
            if !(exists) { goto NG }
            var field = u.FieldList[index]
            var field_t_ = inflateFieldType(field, d, a)
            var lens2_t_ = program.T_Lens2(b.Type, field_t_)
            var lens2_t = typsys.CertainType { Type: lens2_t_ }
            return cc.assign(lens2_t, loc,
                program.InteriorRef {
                    Base:    in_expr,
                    Index:   index,
                    Kind:    program.RK_UnionItem,
                    Operand: o,
                })
        }
    }
    if _, I, a, o, b, ok := getInteriorReferableInterface(in_t, ctx); ok {
        var C, exists = ctx.resolveType(ref)
        if !(exists) { goto NG }
        var _, is_interface = C.Content.(typsys.Interface)
        if is_interface { goto NG }
        var table, ok = ctx.resolveDispatchTable(C, I)
        if !(ok) { goto NG }
        var base_type = b.Type
        var concrete_type = typsys.RefType { Def: C.Ref, Args: a }
        var lens2_t_ = program.T_Lens2(base_type, concrete_type)
        var lens2_t = typsys.CertainType { Type: lens2_t_ }
        return cc.assign(lens2_t, loc,
            program.InteriorRef {
                Base:    in_expr,
                Table:   table,
                Kind:    program.RK_DynamicCast,
                Operand: o,
            })
    }
    NG:
    return cc.error(loc,
        E_InteriorRefUnavailable {
            InteriorRef: ref.String(),
            TypeDesc:    typsys.DescribeCertain(in_t),
        })
}

func checkPipeInfix(infix ast.PipeInfix, in ast.Expr, cc *exprCheckContext) (*program.Expr, *source.Error) {
    var left, err = cc.checkChildExpr(nil, in)
    if err != nil { return nil, err }
    var right = infix.PipeCall
    var op = infix.Callee
    var loc = infix.Location
    return checkCallFunRef(op, left, right, true, loc, cc)
}
func checkInfixTerm(I ast.InfixTerm, cc *exprCheckContext) (*program.Expr, *source.Error) {
    var left, err = cc.checkChildExpr(nil, I.Left)
    if err != nil { return nil, err }
    var right = craftInfixRight(I.Right, I.Node)
    var op = I.Operator
    var loc = I.Location
    return checkCallFunRef(op, left, right, true, loc, cc)
}

func checkPipeCall(args ast.VariousPipeCall, callee ast.Expr, cc *exprCheckContext) (*program.Expr, *source.Error) {
    var loc = callee.Location
    if ref_node, tag, is_new, is_ref := ast.GetStandaloneRef(callee); is_ref {
        if is_new {
            return checkNewRecord(ref_node, args, tag, loc, cc)
        } else {
            var ctx = cc.context
            var ref = getRef(ref_node.Base)
            if ((ref.Namespace == "") && ctx.hasBinding(ref.ItemName)) {
                return checkCallExpr(callee, args, cc)
            } else {
                return checkCallFunRef(ref_node, nil, args, false, loc, cc)
            }
        }
    } else {
        return checkCallExpr(callee, args, cc)
    }
}


