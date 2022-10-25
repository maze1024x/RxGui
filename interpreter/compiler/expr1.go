package compiler

import (
    "rxgui/interpreter/lang/textual/ast"
    "rxgui/interpreter/lang/source"
    "rxgui/interpreter/lang/typsys"
    "rxgui/interpreter/program"
)


func checkLambda(L ast.Lambda, cc *exprCheckContext) (*program.Expr, *source.Error) {
    { var cc, capture = cc.withCaptureScope()
    var loc = L.Location
    var in_t, out_t_, err = cc.getExpectedLambda(loc)
    if err != nil { return nil, err }
    var self_binding *program.Binding
    if self_name, ok := L.SelfRefName.(ast.Identifier); ok {
        var self_pattern = craftPatternSingle(self_name)
        var lambda_t_ = program.T_Lambda(in_t.Type, out_t_)
        var self_t, ok = cc.getCertainOrInferred(lambda_t_)
        if !(ok) { return cc.error(loc, E_ExpectExplicitTypeCast {}) }
        var self_pm, err = cc.match(self_pattern, self_t)
        if err != nil { panic("something went wrong") }
        if len(self_pm) != 1 { panic("something went wrong") }
        self_binding = self_pm[0].Binding
    }
    { var in_pm, err = cc.match(L.InputPattern, in_t)
    if err != nil { return nil, err }
    { var out_expr, err = cc.checkChildExpr(out_t_, L.OutputExpr)
    if err != nil { return nil, err }
    var out_t = out_expr.Type
    var lambda_t_ = program.T_Lambda(in_t.Type, out_t.Type)
    var lambda_t = typsys.CertainType { Type: lambda_t_ }
    var ctx = capture()
    return cc.assign(lambda_t, loc,
        program.Lambda {
            Ctx:  ctx,
            In:   in_pm,
            Out:  out_expr,
            Self: self_binding,
        })
    }}}
}
func (cc *exprCheckContext) getExpectedLambda(loc source.Location) (typsys.CertainType, typsys.Type, *source.Error) {
    var c_nil typsys.CertainType
    if cc.expected == nil {
        return c_nil, nil, source.MakeError(loc,
            E_ExpectExplicitTypeCast {})
    }
    var t = (func() typsys.Type {
        if certain, ok := cc.getExpectedCertainOrInferred(); ok {
            return certain.Type
        } else {
            return cc.expected
        }
    })()
    var raw_input, output, ok = program.T_Lambda_(t)
    if !(ok) {
        var ctx = cc.context
        if t, sam_ok := getSamInterfaceMethodType(t, ctx); sam_ok {
            raw_input, output, ok = program.T_Lambda_(t)
        }
    }
    if !(ok) {
        return c_nil, nil, source.MakeError(loc,
            E_LambdaAssignedToIncompatibleType {
                TypeDesc: cc.describeExpected(),
            })
    }
    { var input, ok = cc.getCertainOrInferred(raw_input)
    if !(ok) {
        return c_nil, nil, source.MakeError(loc,
            E_ExpectExplicitTypeCast {})
    }
    return input, output, nil }
}
func craftPatternSingle(name ast.Identifier) ast.VariousPattern {
    return ast.VariousPattern {
        Node:    name.Node,
        Pattern: ast.PatternSingle {
            Node: name.Node,
            Name: name,
        },
    }
}

func checkBlock(B ast.Block, cc *exprCheckContext) (*program.Expr, *source.Error) {
    if len(B.Bindings) == 0 {
        return cc.forwardTo(B.Return)
    } else {
        var first, rest = cutAstBlock(B)
        switch B := first.Binding.(type) {
        case ast.BindingPlain:
            if B.Off {
                return cc.forwardTo(rest)
            }
            var val_expr, err = cc.checkChildExpr3(nil, B.Value, B.Const)
            if err != nil { return nil, err }
            var val_t = val_expr.Type
            { var cc = cc.withBlockScope()
            { var pm, err = cc.match3(B.Pattern, val_t, B.Const)
            if err != nil { return nil, err }
            { var rest_expr, err = cc.forwardTo(rest)
            if err != nil { return nil, err }
            var block_t = rest_expr.Type
            return cc.assign(block_t, B.Location,
                program.Let {
                    In:  pm,
                    Arg: val_expr,
                    Out: rest_expr,
                })
            }}}
        case ast.BindingCps:
            if B.Off {
                return cc.forwardTo(rest)
            }
            var k = ast.WrapTermAsExpr(ast.VariousTerm {
                Node: B.Node,
                Term: ast.Lambda {
                    Node:         B.Node,
                    InputPattern: B.Pattern,
                    OutputExpr:   rest,
                },
            })
            var left, err = cc.checkChildExpr(nil, B.Value)
            if err != nil { return nil, err }
            var right = craftInfixRight(k, B.Node)
            var op = B.Callee
            var loc = B.Location
            return checkCallFunRef(op, left, right, true, loc, cc)
        default:
            panic("impossible branch")
        }
    }
}
func cutAstBlock(B ast.Block) (ast.VariousBinding, ast.Expr) {
    var first = B.Bindings[0]
    var rest = B.Bindings[1:]
    var rest_node = (func() ast.Node {
        if len(rest) > 0 {
            return rest[0].Node
        } else {
            return B.Return.Node
        }
    })()
    var rest_block = ast.WrapBlockAsExpr(ast.Block {
        Node:     rest_node,
        Bindings: rest,
        Return:   B.Return,
    })
    return first, rest_block
}


