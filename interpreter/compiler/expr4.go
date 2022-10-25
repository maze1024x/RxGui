package compiler

import (
	"rxgui/util/ctn"
	"rxgui/interpreter/lang/source"
	"rxgui/interpreter/lang/typsys"
	"rxgui/interpreter/lang/textual/ast"
	"rxgui/interpreter/program"
)


func checkIf(I ast.If, cc *exprCheckContext) (*program.Expr, *source.Error) {
	var return_type = cc.expected
	var branches = make([] program.IfBranch, 0)
	var add = func(c ([] ast.Cond), b ast.Block) *source.Error {
		{ var cc = cc.withBlockScope()
		var conds = make([] program.Cond, len(c))
		for i, c := range c {
			var value, err = cc.checkChildExpr(nil, c.Expr)
			if err != nil { return err }
			var inner_t, kind, ok = getCondInnerType(value.Type)
			if !(ok) {
				return source.MakeError(c.Expr.Location,
					E_InvalidCondType {
						TypeDesc: typsys.DescribeCertain(value.Type),
					})
			}
			if kind == program.CK_Bool {
			if pattern, ok := c.Pattern.(ast.VariousPattern); ok {
				return source.MakeError(pattern.Location,
					E_InvalidCondPattern {})
			}}
			{ var match, err = cc.match(c.Pattern, inner_t)
			if err != nil { return err }
			conds[i] = program.Cond {
				Kind:  kind,
				Match: match,
				Value: value,
			}}
		}
		var b = ast.WrapTermAsExpr(ast.VariousTerm { Node: b.Node, Term: b })
		var value, err = cc.checkChildExpr(return_type, b)
		if err != nil { return err }
		if return_type == nil {
			return_type = value.Type.Type
		}
		branches = append(branches, program.IfBranch {
			Conds: conds,
			Value: value,
		})
		return nil }
	}
	var err = add(I.Conds, I.Yes)
	if err != nil { return nil, err }
	for _, elif := range I.ElIfs {
		var err = add(elif.Conds, elif.Yes)
		if err != nil { return nil, err }
	}
	{ var err = add(nil, I.No)
	if err != nil { return nil, err } }
	var t, ok = cc.getCertainOrInferred(return_type)
	if !(ok) { panic("something went wrong") }
	var loc = I.Location
	return cc.assign(t, loc,
		program.If {
			Branches: branches,
		})
}

func checkWhen(W ast.When, cc *exprCheckContext) (*program.Expr, *source.Error) {
	var ctx = cc.context
	var operand, err = cc.checkChildExpr(nil, W.Operand)
	if err != nil { return nil, err }
	var operand_t = operand.Type
	var return_type = cc.expected
	if f, d, a, k, ok := getUnionOrEnumFields(operand_t.Type, ctx); ok {
		var branches = make([] *program.WhenBranch, len(f.FieldList))
		for _, c := range W.Cases {
		for _, n := range c.Names {
			if c.Off {
				break
			}
			var in = c.InputPattern
			var out = c.OutputExpr
			var key = ast.Id2String(n)
			var not_default = (key != Underscore)
			var cc, match, index, err = (func() (*exprCheckContext, program.PatternMatching, int, *source.Error) {
				if not_default {
					var index, exists = f.FieldIndexMap[key]
					if !(exists) {
						var loc = n.Location
						return nil, nil, -1, source.MakeError(loc,
							E_NoSuchCase {
								CaseName: key,
							})
					}
					if branches[index] != nil {
						var loc = n.Location
						return nil, nil, -1, source.MakeError(loc,
							E_DuplicateCase {
								CaseName: key,
							})
					}
					if k == program.UE_Union {
						var field = f.FieldList[index]
						var in_t_ = inflateFieldType(field, d, a)
						var in_t = typsys.CertainType { Type: in_t_ }
						var cc = cc.withBlockScope()
						var match, err = cc.match(in, in_t)
						if err != nil { return nil, nil, -1, err }
						return cc, match, index, nil
					} else {
						if pattern, ok := in.(ast.VariousPattern); ok {
							var loc = pattern.Location
							return nil, nil, -1, source.MakeError(loc,
								E_InvalidCasePattern {})
						}
						return cc, nil, index, nil
					}
				} else {
					if pattern, ok := in.(ast.VariousPattern); ok {
						var loc = pattern.Location
						return nil, nil, -1, source.MakeError(loc,
							E_InvalidCasePattern {})
					}
					return cc, nil, -1, nil
				}
			})()
			if err != nil { return nil, err }
			{ var value, err = cc.checkChildExpr(return_type, out)
			if err != nil { return nil, err }
			if return_type == nil {
				return_type = value.Type.Type
			}
			var branch = &program.WhenBranch {
				Match: match,
				Value: value,
			}
			if not_default {
				branches[index] = branch
			} else {
				var ok = false
				for i := range branches {
					if branches[i] == nil {
						branches[i] = branch
						ok = true
					}
				}
				if !(ok) {
					var loc = n.Location
					return cc.error(loc, E_SuperfluousDefaultCase {})
				}
			}}
		}}
		for i := range branches {
			if branches[i] == nil {
				var missing = f.FieldList[i].Name
				var loc = W.Location
				return cc.error(loc,
					E_MissingCase {
						CaseName: missing,
					})
			}
		}
		var t, ok = cc.getCertainOrInferred(return_type)
		if !(ok) { panic("something went wrong") }
		var operand = program.WhenOperand {
			Kind:  k,
			Value: operand,
		}
		var loc = W.Location
		return cc.assign(t, loc,
			program.When {
				Operand:  operand,
				Branches: branches,
			})
	}
	var loc = W.Operand.Location
	return cc.error(loc,
		E_InvalidWhenOperand {
			TypeDesc: typsys.DescribeCertain(operand_t),
		})
}

func checkEach(E ast.Each, cc *exprCheckContext) (*program.Expr, *source.Error) {
	var ctx = cc.context
	var operand_t, err = ctx.makeType(E.Operand)
	if err != nil { return nil, err }
	var item_type = (func() typsys.Type {
		if cc.expected == nil {
			return nil
		}
		if item_type, ok := program.T_List_(cc.expected); ok {
			return item_type
		}
		return nil
	})()
	if f, d, a, k, ok := getUnionOrEnumFields(operand_t.Type, ctx); ok {
		var values = make([] program.EachValue, 0)
		var occurred = make(map[int] struct{})
		for _, c := range E.Cases {
		for _, n := range c.Names {
			if c.Off {
				break
			}
			var in = c.InputPattern
			var out = c.OutputExpr
			var key = ast.Id2String(n)
			var is_default = (key == Underscore)
			if is_default {
				var loc = n.Location
				return cc.error(loc, E_SuperfluousDefaultCase {})
			}
			var cc, match, index, err = (func() (*exprCheckContext, program.PatternMatching, int, *source.Error) {
				var index, exists = f.FieldIndexMap[key]
				if !(exists) {
					var loc = n.Location
					return nil, nil, -1, source.MakeError(loc,
						E_NoSuchCase {
							CaseName: key,
						})
				}
				var _, duplicate = occurred[index]
				occurred[index] = struct{}{}
				if duplicate {
					var loc = n.Location
					return nil, nil, -1, source.MakeError(loc,
						E_DuplicateCase {
							CaseName: key,
						})
				}
				if k == program.UE_Union {
					var field = f.FieldList[index]
					var field_t_ = inflateFieldType(field, d, a)
					var union_t_ = operand_t.Type
					var in_t_ = program.T_Lambda(field_t_, union_t_)
					var in_t = typsys.CertainType { Type: in_t_ }
					var cc = cc.withBlockScope()
					var match, err = cc.match(in, in_t)
					if err != nil { return nil, nil, -1, err }
					return cc, match, index, nil
				} else {
					var enum_t = operand_t
					var in_t = enum_t
					var cc = cc.withBlockScope()
					var match, err = cc.match(in, in_t)
					if err != nil { return nil, nil, -1, err }
					return cc, match, index, nil
				}
			})()
			if err != nil { return nil, err }
			{ var value, err = cc.checkChildExpr(item_type, out)
			if err != nil { return nil, err }
			if item_type == nil {
				item_type = value.Type.Type
			}
			{ var value = program.EachValue {
				Kind:  k,
				Index: index,
				Match: match,
				Value: value,
			}
			values = append(values, value) }}
		}}
		for i, field := range f.FieldList {
			var index = i
			var _, ok = occurred[index]
			if !(ok) {
				var missing = field.Name
				var loc = E.Location
				return cc.error(loc,
					E_MissingCase {
						CaseName: missing,
					})
			}
		}
		var item_t, ok = cc.getCertainOrInferred(item_type)
		if !(ok) { panic("something went wrong") }
		var t_ = program.T_List(item_t.Type)
		var t = typsys.CertainType { Type: t_ }
		var items = ctn.MapEach(values, func(value program.EachValue) *program.Expr {
			return &program.Expr {
				Type:    item_t,
				Info:    value.Value.Info,
				Content: value,
			}
		})
		var loc = E.Location
		return cc.assign(t, loc,
			program.List {
				Items: items,
			})
	}
	var loc = E.Operand.Location
	return cc.error(loc,
		E_InvalidEachOperand {
			TypeDesc: typsys.DescribeCertain(operand_t),
		})
}


func getCondInnerType(t typsys.CertainType) (typsys.CertainType, program.CondKind, bool) {
	var nil_t typsys.CertainType
	if program.T_Bool_(t.Type) {
		return nil_t, program.CK_Bool, true
	}
	if inner, ok := program.T_Maybe_(t.Type); ok {
		var inner_t = typsys.CertainType { Type: inner }
		return inner_t, program.CK_Maybe, true
	}
	if _, inner, ok := program.T_Lens2_(t.Type); ok {
		var inner_t = typsys.CertainType { Type: inner }
		return inner_t, program.CK_Lens2, true
	}
	return nil_t, program.CondKind(-1), false
}
func getUnionOrEnumFields(t typsys.Type, ctx *exprContext) (*typsys.Fields, *typsys.TypeDef, ([] typsys.Type), program.UnionOrEnum, bool) {
	if u, d, a, ok := getUnion(t, ctx); ok {
		return u.Fields, d, a, program.UE_Union, true
	}
	if e, d, ok := getEnum(t, ctx); ok {
		return e.Fields, d, nil, program.UE_Enum, true
	}
	return nil, nil, nil, program.UnionOrEnum(-1), false
}


