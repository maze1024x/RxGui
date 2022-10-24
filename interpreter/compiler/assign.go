package compiler

import (
	"rxgui/lang/source"
	"rxgui/lang/typsys"
	"rxgui/interpreter/program"
)


func assign(expected typsys.Type, expr *program.Expr, ctx *exprContext, s0 *typsys.InferringState) (*program.Expr, *typsys.InferringState, *source.Error) {
	var loc = expr.Info.Location
	var t = expr.Type.Type
	if certain, ok := typsys.GetCertainOrInferred(expected, s0); ok {
		expected = certain.Type
	}
	// direct
	if ok, s1 := typsys.Match(expected, t, s0); ok {
		return expr, s1, nil
	}
	// convert to interface
	if I, Id, Ia, ok := getInterface(expected, ctx); ok {
		var d, a, ok = getTypeDefArgs(t, ctx)
		if !(ok) { goto NG }
		if _, is_interface := d.Content.(typsys.Interface); is_interface {
			// interface -> interface (upward)
			var ok, s1 = typsys.MatchAll(Ia, a, s0)
			if !(ok) { goto NG }
			{ var It, ok = typsys.GetCertainOrInferred(expected, s1)
			if !(ok) { goto NG }
			{ var up, ok = getInterfacePath(d, Id, ctx, nil)
			if !(ok) { goto NG }
			var converted = &program.Expr {
				Type:    It,
				Info:    expr.Info,
				Content: program.InterfaceTransformUpward {
					Arg:  expr,
					Path: up,
				},
			}
			return converted, s1, nil }}
		} else {
			// concrete -> interface
			var table, ok = ctx.resolveDispatchTable(d, Id)
			if !(ok) {
				if isSamInterface(I, Id) {
					var field = I.FieldList[0]
					var field_t = inflateFieldType(field, Id, Ia)
					var ok, s1 = typsys.Match(field_t, t, s0)
					if !(ok) { goto NG }
					{ var It, ok = typsys.GetCertainOrInferred(expected, s1)
					if !(ok) { goto NG }
					var converted = &program.Expr {
						Type: It,
						Info: expr.Info,
						Content: program.InterfaceFromSamValue {
							Value: expr,
						},
					}
					return converted, s1, nil }
				}
				goto NG
			}
			{ var ok, s1 = typsys.MatchAll(Ia, a, s0)
			if !(ok) { goto NG }
			{ var It, ok = typsys.GetCertainOrInferred(expected, s1)
			if !(ok) { goto NG }
			var converted = &program.Expr {
				Type:    It,
				Info:    expr.Info,
				Content: program.Interface {
					ConcreteValue: expr,
					DispatchTable: table,
				},
			}
			return converted, s1, nil }}
		}
	}
	if e, ok := program.T_Maybe_(expected); ok {
	if _, Id, Ia, ok := getInterface(e, ctx); ok {
	if d, a, ok := getTypeDefArgs(t, ctx); ok {
	if ((d != Id) && !(program.T_Null_(t))) {
		if _, is_interface := d.Content.(typsys.Interface); is_interface {
			// interface -> interface (downward)
			var ok, s1 = typsys.MatchAll(Ia, a, s0)
			if !(ok) { goto NG }
			{ var maybe_It, ok = typsys.GetCertainOrInferred(expected, s1)
			if !(ok) { goto NG }
			{ var down, ok = getInterfacePath(Id, d, ctx, nil)
			if !(ok) { goto NG }
			var converted = &program.Expr {
				Type:    maybe_It,
				Info:    expr.Info,
				Content: program.InterfaceTransformDownward {
					Arg:    expr,
					Depth:  len(down),
					Target: Id.Ref.String(),
				},
			}
			return converted, s1, nil }}
		} else {
			// interface -> concrete
			// *** implemented as lens2 interior reference ***
			goto NG
		}
	}}}}
	// convert to union
	if U, Ud, Ua, ok := getUnion(expected, ctx); ok {
		var converted *program.Expr
		var s2 *typsys.InferringState
		var key string
		var found = false
		for i, field := range U.FieldList {
			var index = i
			var field_t = inflateFieldType(field, Ud, Ua)
			if ok, s1 := typsys.Match(field_t, t, s0); ok {
			if certain, ok := typsys.GetCertainOrInferred(expected, s1); ok {
				if found {
					var this_key = field.Name
					return nil, nil, source.MakeError(loc,
						E_AmbiguousAssignmentToUnion {
							Union: Ud.Ref.String(),
							Key1:  key,
							Key2:  this_key,
						})
				}
				found = true
				key = field.Name
				s2 = s1
				converted = &program.Expr {
					Type:    certain,
					Info:    expr.Info,
					Content: program.Union {
						Index: index,
						Value: expr,
					},
				}
			}}
		}
		if found {
			return converted, s2, nil
		} else {
			goto NG
		}
	}
	// convert to Int
	if program.T_Int_(expected) {
		var _, _, ok = getEnum(t, ctx)
		if !(ok) { goto NG }
		var converted = &program.Expr {
			Type:    typsys.CertainType { Type: program.T_Int() },
			Info:    expr.Info,
			Content: program.EnumToInt { EnumValue: expr },
		}
		return converted, s0, nil
	}
	// convert to T_IntoReflectType(e)
	if e, ok := program.T_IntoReflectType_(expected); ok {
		{ var ok = program.T_DummyReflectType_(expr.Type.Type)
		if !(ok) { goto NG }
		{ var certain, ok = typsys.GetCertainOrInferred(e, s0)
		if !(ok) { goto NG }
		var rt = program.MakeReflectType_(certain)
		var converted_t = program.T_IntoReflectType(certain.Type)
		var converted = &program.Expr {
			Type:    typsys.CertainType { Type: converted_t },
			Info:    expr.Info,
			Content: program.ReflectType { Type: rt },
		}
		return converted, s0, nil }}
	}
	// convert to T_IntoReflectValue(e)
	if e, ok := program.T_IntoReflectValue_(expected); ok {
		var expr, s1, err = assign(e, expr, ctx, s0)
		if err != nil { return nil, nil, err }
		var rt = program.MakeReflectType_(expr.Type)
		var converted_t = program.T_IntoReflectValue(expr.Type.Type)
		var converted = &program.Expr {
			Type:    typsys.CertainType { Type: converted_t },
			Info:    expr.Info,
			Content: program.ReflectValue {
				Type:  rt,
				Value: expr,
			},
		}
		return converted, s1, nil
	}
	NG:
	return nil, nil, source.MakeError(loc,
		E_NotAssignable {
			From: typsys.Describe(t),
			To:   typsys.DescribeWithInferringState(expected, s0),
		})
}


