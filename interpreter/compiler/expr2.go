package compiler

import (
	"strings"
	"rxgui/standalone/ctn"
	"rxgui/lang/source"
	"rxgui/lang/typsys"
	"rxgui/lang/textual/ast"
	"rxgui/interpreter/program"
)


func checkCallFunRef(callee ast.Ref, arg0 *program.Expr, args ast.VariousPipeCall, infix bool, loc source.Location, cc *exprCheckContext) (*program.Expr, *source.Error) {
	var ctx = cc.context
	var ref = getRef(callee.Base)
	var t0 = getOptionalExprType(arg0)
	var ref_loc = callee.Location
	var hdr, f, ns, key, err = determineCallee(ref, t0, infix, ref_loc, ctx)
	if err != nil { return nil, err }
	var Tp = hdr.typeParams
	{ var Ta, err = getTypeArgs(callee.TypeArgs, ctx)
	if err != nil { return nil, err }
	var O = hdr.output.type_
	var E = hdr.inputsExplicit
	var I = hdr.inputsImplicit
	var dvg = makeCallFunRefDefaultValueGetter(ns, key, ctx)
	var va = hdr.variadic
	return cc.infer(Tp, Ta, O, loc, func(cc *exprCheckContext) (program.ExprContent, *source.Error) {
		var Ea_, Ia_ = splitRawArgs(args, I)
		var Ea = adaptRawArgs(arg0, Ea_)
		var Ia = adaptRawArgs(nil, Ia_)
		var f_args, err = checkArguments(cc, Ea, E, Tp, dvg, va, loc)
		if err != nil { return nil, err }
		{ var f_ctx, err = checkArguments(cc, Ia, I, Tp, dvg, false, loc)
		if err != nil { return nil, err }
		return program.CallFunction {
			Location:  loc,
			Callee:    f,
			Context:   f_ctx,
			Arguments: f_args,
		}, nil }
	})}
}

func checkNewRecord(record ast.Ref, args ast.VariousPipeCall, tag string, loc source.Location, cc *exprCheckContext) (*program.Expr, *source.Error) {
	var ctx = cc.context
	var ref = getRef(record.Base)
	var def, exists = ctx.resolveType(ref)
	if !(exists) { return cc.error(loc, E_NoSuchType { ref.String() }) }
	var R_, ok = def.Content.(typsys.Record)
	if !(ok) { return cc.error(loc, E_NotRecord { ref.String() }) }
	var rt, tag_valid = getRecordTransform(tag)
	if !(tag_valid) { return cc.error(loc, E_InvalidRecordTag { tag }) }
	var R = rt.fields.apply(R_.Fields)
	var Ra = adaptRawArgs(nil, args)
	var Tp = def.Parameters
	var Ta, err = getTypeArgs(record.TypeArgs, ctx)
	if err != nil { return nil, err }
	var output = rt.recordType.apply(typsys.DefType(def))
	var dvt = makeDefaultValueTransform(rt)
	var dvg = makeNewRecordDefaultValueGetter(ref, dvt, ctx)
	return cc.infer(Tp, Ta, output, loc, func(cc *exprCheckContext) (program.ExprContent, *source.Error) {
		var values, err = checkArguments(cc, Ra, R, Tp, dvg, false, loc)
		if err != nil { return nil, err }
		return rt.recordValue.apply(program.Record { Values: values }), nil
	})
}

func checkCallExpr(callee ast.Expr, args ast.VariousPipeCall, cc *exprCheckContext) (*program.Expr, *source.Error) {
	var loc = args.Location
	var callee_expr, err = cc.checkChildExpr(nil, callee)
	if err != nil { return nil, err }
	if in_t, out_t, sam, ok := cc.getCallable(callee_expr.Type); ok {
		if sam {
			var lambda_t_ = program.T_Lambda(in_t.Type, out_t.Type)
			var lambda_t = typsys.CertainType { Type: lambda_t_ }
			callee_expr = &program.Expr {
				Type:    lambda_t,
				Info:    callee_expr.Info,
				Content: program.AbstractMethodValue {
					Interface: callee_expr,
					Index:     0,
				},
			}
		}
		var arg_t = in_t
		var arg_spec = getLambdaParameters(arg_t)
		var arg_content, err = checkLambdaArguments(arg_spec, args, cc)
		if err != nil { return nil, err }
		var arg_expr = &program.Expr {
			Type:    arg_t,
			Info:    program.ExprInfoFrom(loc),
			Content: arg_content,
		}
		return cc.assign(out_t, loc,
			program.CallLambda {
				Callee:   callee_expr,
				Argument: arg_expr,
			})
	} else {
		return cc.error(loc,
			E_NotCallable {
				TypeDesc: typsys.DescribeCertain(callee_expr.Type),
			})
	}
}
func (cc *exprCheckContext) getCallable(t typsys.CertainType) (typsys.CertainType, typsys.CertainType, bool, bool) {
	var sam = false
	var in, out, ok = program.T_Lambda_(t.Type)
	if !(ok) {
		if t, sam_ok := getSamInterfaceMethodType(t.Type, cc.context); sam_ok {
			sam = true
			in, out, ok = program.T_Lambda_(t)
		}
	}
	if ok {
		var in_t = typsys.CertainType { Type: in }
		var out_t = typsys.CertainType { Type: out }
		return in_t, out_t, sam, true
	} else {
		var nil_t typsys.CertainType
		return nil_t, nil_t, false, false
	}
}

func checkRefTerm(ref_term ast.RefTerm, cc *exprCheckContext) (*program.Expr, *source.Error) {
	var ctx = cc.context
	var loc = ref_term.Location
	var ref_node = ref_term.Ref
	var _, is_new = ref_term.New.(ast.New)
	if is_new {
		return cc.error(loc,
			E_InvalidConstructorUsage {})
	}
	var ref = getRef(ref_node.Base)
	if ref.Namespace == "" {
		if binding, ok := ctx.useBinding(ref.ItemName); ok {
			if len(ref_node.TypeArgs) > 0 {
				return cc.error(loc,
					E_SuperfluousTypeArgs {})
			}
			return cc.assign(binding.Type, loc,
				program.LocalRef {
					Binding: binding,
				})
		}
		if t, ok := cc.getExpectedCertainOrInferred(); ok {
		if e, _, ok := getEnum(t.Type, ctx); ok {
		if index, ok := e.FieldIndexMap[ref.ItemName]; ok {
			return cc.assign(t, loc,
				program.Enum (
					index,
				))
		}}}
	}
	var Ta, err = getTypeArgs(ref_term.Ref.TypeArgs, ctx)
	if err != nil { return nil, err }
	var options = findFunRefOptions(ref, ctx)
	if len(options) == 0 {
		return cc.error(loc, E_NoSuchThing { ref.String() })
	}
	var names = ctn.MapEach(options, func(option func()(*funHeader,**program.Function,string)) string {
		var _, _, name = option()
		return name
	})
	var trials = ctn.MapEach(options, func(option func()(*funHeader,**program.Function,string)) func()(*program.Expr,**typsys.InferringState,*source.Error) {
		var hdr, f, _ = option()
		return func() (*program.Expr, **typsys.InferringState, *source.Error) {
			var Tp = hdr.typeParams
			var O = hdr.output.type_
			var E = hdr.inputsExplicit
			var I = hdr.inputsImplicit
			if hdr.funKind == FK_Const {
				if len(Tp) > 0 {
					panic("something went wrong")
				}
				if len(Ta) > 0 {
					return nil, nil, source.MakeError(loc,
						E_SuperfluousTypeArgs {})
				}
				if len(E.FieldList) > 0 || len(I.FieldList) > 0 {
					panic("something went wrong")
				}
				var t = typsys.CertainType { Type: O }
				return cc.tryAssign(t, loc,
					program.CallFunction {
						Location: loc,
						Callee:   f,
					})
			} else {
				if len(I.FieldList) > 0 {
					return nil, nil, source.MakeError(loc,
						E_CannotAssignFunctionRequiringImplicitInput {})
				}
				var L, unpack, ok = makeLambdaType(E, O)
				if !(ok) {
					return nil, nil, source.MakeError(loc,
						E_UnableToUseAsLambda { describeFunInOut(hdr) })
				}
				return cc.tryInfer(Tp, Ta, L, loc, func(cc *exprCheckContext) (program.ExprContent, *source.Error) {
					return program.FunRef {
						Function: f,
						Unpack:   unpack,
					}, nil
				})
			}
		}
	})
	var ref_expr *program.Expr
	var ref_S **typsys.InferringState
	var found = false
	var details = make([] NoneAssignableErrorDetail, 0)
	for i, trial := range trials {
		var expr, S, err = trial()
		if err == nil {
			if found {
				return cc.error(loc,
					E_MultipleAssignable { ref.String() })
			} else {
				ref_expr = expr
				ref_S = S
				found = true
			}
		} else {
			details = append(details, NoneAssignableErrorDetail {
				ItemName:     names[i],
				ErrorContent: err.Content,
			})
		}
	}
	if found {
		cc.loadTrialInferringState(ref_S)
		return ref_expr, nil
	} else {
		return cc.error(loc,
			E_NoneAssignable { ref.String(), details })
	}
}

func checkImplicitRefTerm(I ast.ImplicitRefTerm, cc *exprCheckContext) (*program.Expr, *source.Error) {
	var name = ast.Id2String(I.Name)
	var ref_node = (func() ast.Ref {
		if P, item, ok := ast.SplitImplicitRef(name); ok {
		if s := cc.getInferringState(); s != nil {
		if ns, ok := s.GetInferredParameterNamespace(P); ok {
			return ast.Strings2Ref(ns, item, I.Node)
		}}}
		return ast.String2Ref(name, I.Node)
	})()
	var ref_term = ast.RefTerm {
		Node: ref_node.Node,
		Ref:  ref_node,
	}
	return checkRefTerm(ref_term, cc)
}


type arguments ([] argument)
type argument struct {
	name  string
	expr  *program.Expr
	node  *ast.Expr
}
func (arg argument) getName() (string, bool) {
	return arg.name, (arg.name != "")
}
func (arg argument) getExpr() (*program.Expr, bool) {
	return arg.expr, (arg.expr != nil)
}
func (arg argument) getNode() (*ast.Expr, bool) {
	return arg.node, (arg.node != nil)
}
func (arg argument) getLocation() source.Location {
	if expr, ok := arg.getExpr(); ok {
		return expr.Info.Location
	} else if node, ok := arg.getNode(); ok {
		return node.Location
	} else {
		panic("impossible branch")
	}
}
func (arg argument) check(cc *exprCheckContext, expected typsys.Type) (*program.Expr, *source.Error) {
	if expr, ok := arg.getExpr(); ok {
		return cc.assignChildExpr(expected, expr)
	} else if node, ok := arg.getNode(); ok {
		return cc.checkChildExpr(expected, *node)
	} else {
		panic("impossible branch")
	}
}
func adaptRawArgs(raw0 *program.Expr, raw ast.VariousPipeCall) arguments {
	var result = make(arguments, 0)
	if raw0 != nil {
		result = append(result, argument {
			expr: raw0,
		})
	}
	switch R := raw.PipeCall.(type) {
	case ast.CallOrdered:
		for i := range R.Arguments {
			result = append(result, argument {
				node: &(R.Arguments[i]),
			})
		}
	case ast.CallUnordered:
		for _, mapping := range R.Mappings {
			var name = ast.Id2String(mapping.Name)
			var node = getArgMappingExprNode(mapping.Name, mapping.Value)
			result = append(result, argument {
				name: name,
				node: node,
			})
		}
	}
	return result
}
func getArgMappingExprNode(k ast.Identifier, v ast.MaybeExpr) *ast.Expr {
	if node, ok := v.(ast.Expr); ok {
		return &node
	} else {
		var node = ast.WrapTermAsExpr(ast.VariousTerm {
			Node: k.Node,
			Term: ast.RefTerm {
				Node: k.Node,
				Ref:  ast.Ref {
					Node: k.Node,
					Base: ast.RefBase {
						Node: k.Node,
						Item: k,
					},
				},
			},
		})
		return &node
	}
}
func splitRawArgs(raw ast.VariousPipeCall, imp_spec *typsys.Fields) (ast.VariousPipeCall, ast.VariousPipeCall) {
	var imp_mappings = make([] ast.ArgumentMapping, len(imp_spec.FieldList))
	var imp_is_set = make([] bool, len(imp_spec.FieldList))
	var explicit = (func() ast.VariousPipeCall {
		switch R := raw.PipeCall.(type) {
		case ast.CallOrdered:
			return raw
		case ast.CallUnordered:
			var exp_mappings = make([] ast.ArgumentMapping, 0)
			for _, mapping := range R.Mappings {
				var name = ast.Id2String(mapping.Name)
				var index, is_imp = imp_spec.FieldIndexMap[name]
				if is_imp {
					imp_mappings[index] = mapping
					imp_is_set[index] = true
				} else {
					exp_mappings = append(exp_mappings, mapping)
				}
			}
			return ast.VariousPipeCall {
				Node:     raw.Node,
				PipeCall: ast.CallUnordered {
					Node:     raw.Node,
					Mappings: exp_mappings,
				},
			}
		default:
			panic("impossible branch")
		}
	})()
	var implicit = (func() ast.VariousPipeCall {
		for i := range imp_mappings {
			if !(imp_is_set[i]) {
				var name = imp_spec.FieldList[i].Name
				imp_mappings[i] = ast.ArgumentMapping {
					Node:  raw.Node,
					Name:  ast.String2Id(name, raw.Node),
					Value: ast.WrapTermAsExpr(ast.VariousTerm {
						Node: raw.Node,
						Term: ast.ImplicitRefTerm {
							Node: raw.Node,
							Name: ast.String2Id(name, raw.Node),
						},
					}),
				}
			}
		}
		return ast.VariousPipeCall {
			Node:     raw.Node,
			PipeCall: ast.CallUnordered {
				Node:     raw.Node,
				Mappings: imp_mappings,
			},
		}
	})()
	return explicit, implicit
}
type defaultValueGetter struct {
	fragment   fragmentDraft
	namespace  string
	funName    string
	typeName   string
	transform  defaultValueTransform
}
func makeCallFunRefDefaultValueGetter(f_ns string, f_key userFunKey, ctx *exprContext) defaultValueGetter {
	return defaultValueGetter {
		fragment:  ctx.fragment,
		namespace: f_ns,
		funName:   f_key.name,
		typeName:  f_key.assoc,
	}
}
func makeNewRecordDefaultValueGetter(ref source.Ref, dvt defaultValueTransform, ctx *exprContext) defaultValueGetter {
	return defaultValueGetter {
		fragment:  ctx.fragment,
		namespace: ref.Namespace,
		typeName:  ref.ItemName,
		transform: dvt,
	}
}
func (dvg defaultValueGetter) get(field_name string) **program.Function {
	var ns = dvg.namespace
	var key = genFunKey {
		typeName:  dvg.typeName,
		funName:   dvg.funName,
		fieldName: field_name,
	}
	return dvg.fragment.referGenFun(ns, key)
}
func checkArguments (
	cc      *exprCheckContext,
	args    arguments,
	fields  *typsys.Fields,
	inf     ([] string),
	dvg     defaultValueGetter,
	va      bool,
	loc     source.Location,
) ([] *program.Expr, *source.Error) {
	var arity = len(fields.FieldList)
	var last = (arity - 1)
	var result = make([] *program.Expr, arity)
	var ordered_count = 0
	var unordered = make(map[string] int)
	var used = make([] bool, len(args))
	for i, arg := range args {
		if arg.name == "" {
			if len(unordered) > 0 { panic("something went wrong") }
			ordered_count += 1
		} else {
			var name = arg.name
			if _, duplicate := unordered[name]; duplicate {
				return nil, source.MakeError(arg.getLocation(),
					E_DuplicateArgument {
						Name: name,
					})
			}
			unordered[name] = i
		}
	}
	for i, field := range fields.FieldList {
		var field_name = field.Name
		var field_type = typsys.ToInferring(field.Type, inf)
		var field_has_default = field.Info.HasDefaultValue
		if va && (len(unordered) == 0) && (i == last) {
			var pos0 = i
			var va_type, ok = program.T_List_(field_type)
			if !(ok) {
				// error should have been reported on function signature
				break
			}
			var va_list = (func() arguments {
				if pos0 < len(args) {
					return args[pos0:]
				} else {
					return nil
				}
			})()
			var va_expr_list = make([] *program.Expr, len(va_list))
			for j, va_item := range va_list {
				var expr, err = va_item.check(cc, va_type)
				if err != nil { return nil, err }
				va_expr_list[j] = expr
			}
			for pos := pos0; pos < len(args); pos += 1 {
				used[pos] = true
			}
			if va_certain, ok := cc.getCertainOrInferred(va_type); ok {
				result[i] = &program.Expr {
					Type:    va_certain,
					Info:    program.ExprInfoFrom(loc),
					Content: program.List { Items: va_expr_list },
				}
			} else {
				return nil, source.MakeError(loc,
					E_UnableToInferVaType {})
			}
		} else {
			if pos, ok := (func() (int, bool) {
				if i < ordered_count {
					return i, true
				} else if pos, ok := unordered[field_name]; ok {
					return pos, true
				} else {
					return -1, false
				}
			})(); ok {
				var arg = args[pos]
				if used[pos] {
					return nil, source.MakeError(arg.getLocation(),
						E_DuplicateArgument {
							Name: field_name,
						})
				}
				used[pos] = true
				var expr, err = arg.check(cc, field_type)
				if err != nil { return nil, err }
				result[i] = expr
			} else {
				if field_has_default {
					if field_t, ok := cc.getCertainOrInferred(field_type); ok {
						var f = dvg.get(field_name)
						var expr = &program.Expr {
							Type:    field_t,
							Info:    program.ExprInfoFrom(loc),
							Content: dvg.transform.apply(program.CallFunction {
								Location: loc,
								Callee:   f,
							}),
						}
						result[i] = expr
					} else {
						return nil, source.MakeError(loc,
							E_UnableToInferDefaultValueType {
								ArgName: field_name,
							})
					}
				} else {
					return nil, source.MakeError(loc,
						E_MissingArgument {
							Name: field_name,
						})
				}
			}
		}
	}
	for pos, used := range used {
		if !(used) {
			var arg = args[pos]
			return nil, source.MakeError(arg.getLocation(),
				E_SuperfluousArgument {})
		}
	}
	return result, nil
}
func makeLambdaType(inputs *typsys.Fields, output typsys.Type) (typsys.Type, bool, bool) {
	const keep = false
	const unpack = true
	var args = inputs.FieldList
	var ret = output
	var arity = len(args)
	if arity == 1 {
		var arg = args[0].Type
		return program.T_Lambda(arg, ret), keep, true
	} else if arity == 2 {
		var arg = program.T_Pair(args[0].Type, args[1].Type)
		return program.T_Lambda(arg, ret), unpack, true
	} else if arity == 3 {
		var arg = program.T_Triple(args[0].Type, args[1].Type, args[2].Type)
		return program.T_Lambda(arg, ret), unpack, true
	} else {
		return nil, false, false
	}
}
func getLambdaParameters(in_t typsys.CertainType) ([] typsys.CertainType) {
	if program.T_Null_(in_t.Type) {
		return [] typsys.CertainType {}
	}
	if a_, b_, ok := program.T_Pair_(in_t.Type); ok {
		var a = typsys.CertainType { Type: a_ }
		var b = typsys.CertainType { Type: b_ }
		return [] typsys.CertainType { a, b }
	}
	if a_, b_, c_, ok := program.T_Triple_(in_t.Type); ok {
		var a = typsys.CertainType { Type: a_ }
		var b = typsys.CertainType { Type: b_ }
		var c = typsys.CertainType { Type: c_ }
		return [] typsys.CertainType { a, b, c }
	}
	return [] typsys.CertainType { in_t }
}
func checkLambdaArguments(params ([] typsys.CertainType), args ast.VariousPipeCall, cc *exprCheckContext) (program.ExprContent, *source.Error) {
	var loc = args.Location
	switch args := args.PipeCall.(type) {
	case ast.CallOrdered:
		var given = len(args.Arguments)
		var required = len(params)
		if given != required {
			return nil, source.MakeError(loc,
				E_LambdaCallWrongArgsQuantity {
					Given:    given,
					Required: required,
				})
		}
		var expr_list = make([] *program.Expr, len(params))
		for i, p := range params {
			var expr, err = cc.checkChildExpr(p.Type, args.Arguments[i])
			if err != nil { return nil, err }
			expr_list[i] = expr
		}
		var n = len(expr_list)
		if n == 0 {
			return program.Null {}, nil
		} else if n == 1 {
			return expr_list[0].Content, nil
		} else {
			return program.Record { Values: expr_list }, nil
		}
	case ast.CallUnordered:
		return nil, source.MakeError(loc,
			E_LambdaCallUnorderedArgs {})
	default:
		panic("impossible branch")
	}
}

type transform[T any] struct {
	f  func(T) T
}
func makeTransform[T any] (f func(T)(T)) transform[T] {
	return transform[T] { f }
}
func (t transform[T]) apply(v T) T {
	if t.f != nil {
		return t.f(v)
	} else {
		return v
	}
}
type defaultValueTransform struct {
	transform[program.ExprContent]
}
func makeDefaultValueTransform(rt recordTransform) defaultValueTransform {
	return defaultValueTransform { rt.fieldDefault }
}
type recordTransform struct {
	fields        transform[*typsys.Fields]
	recordType    transform[typsys.Type]
	fieldDefault  transform[program.ExprContent]
	recordValue   transform[program.ExprContent]
}
func getRecordTransform(tag string) (recordTransform, bool) {
	switch tag {
	case "":
		return recordTransform {}, true
	case program.T_Observable_Tag:
		return recordTransform {
			fields: makeTransform(func(fields *typsys.Fields) *typsys.Fields {
				return mapFieldsType(fields, program.T_Observable)
			}),
			recordType: makeTransform(func(t typsys.Type) typsys.Type {
				return program.T_Observable(t)
			}),
			fieldDefault: makeTransform(func(e program.ExprContent) program.ExprContent {
				return program.ObservableRecordDefault {
					Content: e,
				}
			}),
			recordValue: makeTransform(func(e program.ExprContent) program.ExprContent {
				return program.ObservableRecord {
					Record: e.(program.Record),
				}
			}),
		}, true
	case program.T_Observable_Maybe_Tag:
		return recordTransform {
			fields: makeTransform(func(fields *typsys.Fields) *typsys.Fields {
				return mapFieldsType(mapFieldsType(fields,
					program.T_Maybe), program.T_Observable)
			}),
			recordType: makeTransform(func(t typsys.Type) typsys.Type {
				return program.T_Observable(program.T_Maybe(t))
			}),
			fieldDefault: makeTransform(func(e program.ExprContent) program.ExprContent {
				return program.ObservableMaybeRecordDefault {
					Content: e,
				}
			}),
			recordValue: makeTransform(func(e program.ExprContent) program.ExprContent {
				return program.ObservableMaybeRecord {
					Record: e.(program.Record),
				}
			}),
		}, true
	case program.T_Hook_Tag:
		return recordTransform {
			fields: makeTransform(func(fields *typsys.Fields) *typsys.Fields {
				return mapFieldsType(fields, program.T_Hook)
			}),
			recordType: makeTransform(func(t typsys.Type) typsys.Type {
				return program.T_Hook(t)
			}),
			fieldDefault: makeTransform(func(e program.ExprContent) program.ExprContent {
				return program.HookRecordDefault {
					Content: e,
				}
			}),
			recordValue: makeTransform(func(e program.ExprContent) program.ExprContent {
				return program.HookRecord {
					Record: e.(program.Record),
				}
			}),
		}, true
	default:
		return recordTransform {}, false
	}
}
func mapFieldsType(fields *typsys.Fields, f func(typsys.Type)(typsys.Type)) *typsys.Fields {
	var m = make(map[string] int)
	for k, v := range fields.FieldIndexMap {
		m[k] = v
	}
	var l = ctn.MapEach(fields.FieldList, func(field typsys.Field) typsys.Field {
		return typsys.Field {
			Info: field.Info,
			Name: field.Name,
			Type: f(field.Type),
		}
	})
	return &typsys.Fields {
		FieldIndexMap: m,
		FieldList:     l,
	}
}

func findFunRefOptions(ref source.Ref, ctx *exprContext) ([] func()(*funHeader,**program.Function,string)) {
	var all = ctx.context.lookupAssignableKindFunctions(ref)
	if len(all) == 0 {
		if ctx.tryRedirectRef(&ref) {
			all = ctx.context.lookupAssignableKindFunctions(ref)
		}
	}
	return ctn.MapEach(all, func(item func()(*funHeader,string,userFunKey)) func()(*funHeader,**program.Function,string) {
		var hdr, ns, key = item()
		var f = ctx.fragment.referUserFun(ns, key)
		var name = key.Describe(ns)
		var option = func() (*funHeader, **program.Function, string) {
			return hdr, f, name
		}
		return option
	})
}
func determineCallee(ref source.Ref, t0 ctn.Maybe[typsys.CertainType], infix bool, loc source.Location, ctx *exprContext) (*funHeader, **program.Function, string, userFunKey, *source.Error) {
	var r0 = getTypeAssocRef(t0)
	var kind = (func() FunKind {
		if infix {
			return FK_Operator
		} else {
			return FK_Ordinary
		}
	})()
	var hdr, ns, key, ok = ctx.context.lookupFunction(ref, r0, kind)
	if !(ok) {
		if ctx.tryRedirectRef(&ref) {
			hdr, ns, key, ok = ctx.context.lookupFunction(ref, r0, kind)
		}
	}
	if ok {
		return hdr, ctx.fragment.referUserFun(ns, key), ns, key, nil
	} else {
		var kind_desc = (func() string {
			if kind == FK_Operator {
				return "operator"
			} else {
				return "function"
			}
		})()
		var buf strings.Builder
		buf.WriteString(ref.String())
		if r0, ok := r0.Value(); ok {
		if r0.Namespace == ref.Namespace {
			buf.WriteRune('(')
			buf.WriteString(r0.ItemName)
			buf.WriteRune(')')
		}}
		var name_desc = buf.String()
		return nil, nil, "", userFunKey{}, source.MakeError(loc,
			E_NoSuchFunction {
				FunKindDesc: kind_desc,
				FunNameDesc: name_desc,
			})
	}
}
func getOptionalExprType(expr *program.Expr) ctn.Maybe[typsys.CertainType] {
	if expr != nil {
		return ctn.Just(expr.Type)
	} else {
		return nil
	}
}
func getTypeAssocRef(t ctn.Maybe[typsys.CertainType]) ctn.Maybe[source.Ref] {
	if t, ok := t.Value(); ok {
		if ref_type, ok := t.Type.(typsys.RefType); ok {
			return ctn.Just(ref_type.Def)
		}
	}
	return nil
}
func getTypeArgs(nodes ([] ast.Type), ctx *exprContext) ([] typsys.CertainType, *source.Error) {
	var args = make([] typsys.CertainType, len(nodes))
	for i, node := range nodes {
		var t, err = ctx.makeType(node)
		if err != nil { return nil, err }
		args[i] = t
	}
	return args, nil
}
func craftInfixRight(r ast.Expr, node ast.Node) ast.VariousPipeCall {
	return ast.VariousPipeCall {
		Node:     node,
		PipeCall: ast.CallOrdered {
			Node:      node,
			Arguments: [] ast.Expr { r },
		},
	}
}
func describeFunInOut(hdr *funHeader) string {
	var inputs = hdr.inputsExplicit.FieldList
	var output = hdr.output
	var buf strings.Builder
	buf.WriteRune('{')
	if len(inputs) > 0 {
		buf.WriteRune(' ')
		for i, item := range inputs {
			buf.WriteString(item.Name)
			buf.WriteRune(' ')
			buf.WriteString(typsys.Describe(item.Type))
			if i != (len(inputs) - 1) {
				buf.WriteRune(',')
				buf.WriteRune(' ')
			}
		}
		buf.WriteRune(' ')
	}
	buf.WriteRune('}')
	buf.WriteRune(' ')
	buf.WriteString(typsys.Describe(output.type_))
	return buf.String()
}


