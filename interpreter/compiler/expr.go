package compiler

import (
	"rxgui/util/ctn"
	"rxgui/interpreter/lang/source"
	"rxgui/interpreter/lang/typsys"
	"rxgui/interpreter/lang/textual/ast"
	"rxgui/interpreter/program"
)


type exprContext struct {
	scope     ctn.MutMap[string,*program.Binding]
	local     ctn.MutSet[*program.Binding]
	used      ctn.MutSet[*program.Binding]
	unused    ctn.MutSet[*program.Binding]
	params    [] string
	fragment  fragmentDraft
	context   *NsHeaderMap
}
func createExprContext(params ([] string), fd fragmentDraft, ctx *NsHeaderMap) *exprContext {
	return &exprContext {
		scope:    ctn.MakeMutMap[string,*program.Binding](ctn.StringCompare),
		local:    ctn.MakeMutSet[*program.Binding](program.BindingCompare),
		used:     ctn.MakeMutSet[*program.Binding](program.BindingCompare),
		unused:   ctn.MakeMutSet[*program.Binding](program.BindingCompare),
		params:   params,
		fragment: fd,
		context:  ctx,
	}
}
func (ctx *exprContext) unusedBindingError() *source.Error {
	if ctx.unused.Size() == 0 {
		return nil
	}
	var bindings = make([] *program.Binding, 0)
	ctx.unused.ForEach(func(item *program.Binding) {
		bindings = append(bindings, item)
	})
	bindings, _ = ctn.StableSorted(bindings, func(a *program.Binding, b *program.Binding) bool {
		return a.Location.Pos.Span.Start < b.Location.Pos.Span.Start
	})
	var first = bindings[0]
	return source.MakeError(first.Location,
		E_UnusedBinding {
			BindingName: first.Name,
		})
}
func (ctx *exprContext) fragmentNamespace() string {
	return ctx.fragment.namespace()
}
func (ctx *exprContext) tryRedirectRef(ref *source.Ref) bool {
	return ctx.context.tryRedirectRef(ctx.fragmentNamespace(), ref)
}
func (ctx *exprContext) createBinding(name string, t typsys.CertainType, loc source.Location) *program.Binding {
	return ctx.createBinding4(name, t, loc, false)
}
func (ctx *exprContext) createBinding4(name string, t typsys.CertainType, loc source.Location, const_ bool) *program.Binding {
	var binding = &program.Binding {
		Name:     name,
		Type:     t,
		Location: loc,
		Constant: const_,
	}
	ctx.addBinding(binding)
	return binding
}
func (ctx *exprContext) addBinding(binding *program.Binding) {
	ctx.scope.Insert(binding.Name, binding)
	ctx.local.Insert(binding)
	ctx.unused.Insert(binding)
}
func (ctx *exprContext) hasBinding(name string) bool {
	var _, exists = ctx.scope.Lookup(name)
	return exists
}
func (ctx *exprContext) useBinding(name string) (*program.Binding, bool) {
	var binding, exists = ctx.scope.Lookup(name)
	if exists {
		ctx.used.Insert(binding)
		ctx.unused.Delete(binding)
	}
	return binding, exists
}
func (ctx *exprContext) copyWithNewConstScope() *exprContext {
	var child = ctx.scope.FilterClone(func(_ string, binding *program.Binding) bool {
		return binding.Constant
	})
	var new_ctx exprContext
	new_ctx = *ctx
	new_ctx.scope = child
	return &new_ctx
}
func (ctx *exprContext) copyWithNewBlockScope() *exprContext {
	var child = ctx.scope.Clone()
	var new_ctx exprContext
	new_ctx = *ctx
	new_ctx.scope = child
	return &new_ctx
}
func (ctx *exprContext) copyWithNewCaptureScope() (*exprContext, func()([]*program.Binding)) {
	var child = ctx.scope.Clone()
	var local = ctn.MakeMutSet[*program.Binding](program.BindingCompare)
	var used = ctn.MakeMutSet[*program.Binding](program.BindingCompare)
	var new_ctx exprContext
	new_ctx = *ctx
	new_ctx.scope = child
	new_ctx.local = local
	new_ctx.used = used
	return &new_ctx, func() ([] *program.Binding) {
		var capture = make([] *program.Binding, 0)
		used.ForEach(func(binding *program.Binding) {
			if !(local.Has(binding)) {
				capture = append(capture, binding)
				ctx.used.Insert(binding)
			}
		})
		return capture
	}
}
func (ctx *exprContext) makeType(node ast.Type) (typsys.CertainType, *source.Error) {
	var cons_ctx = createTypeConsContext(ctx.params)
	var t = makeType(node, cons_ctx)
	var loc = node.Location
	var fragment_ns = ctx.fragmentNamespace()
	var ns_header_map = ctx.context
	var err = validateAlterType(&t, loc, fragment_ns, ns_header_map)
	if err != nil { return typsys.CertainType {}, err }
	return typsys.CertainType { Type: t }, nil
}
func (ctx *exprContext) resolveType(ref source.Ref) (*typsys.TypeDef, bool) {
	var def, exists = ctx.context.lookupType(ref)
	if !(exists) {
		if ctx.tryRedirectRef(&ref) {
			def, exists = ctx.context.lookupType(ref)
		}
	}
	return def, exists
}
func (ctx *exprContext) resolveDispatchTable(d *typsys.TypeDef, Id *typsys.TypeDef) (**program.DispatchTable, bool) {
	if _, is_interface := d.Content.(typsys.Interface); is_interface {
		panic("invalid argument")
	}
	if _, ok := getInterfacePath(d, Id, ctx, nil); ok {
		var ns = d.Ref.Namespace
		var key = dispatchKey {
			ConcreteType:  d.Ref.ItemName,
			InterfaceType: Id.Ref,
		}
		var table = ctx.fragment.referDispatchTable(ns, key)
		return table, true
	} else {
		return nil, false
	}
}
func (ctx *exprContext) resolveMethod(t typsys.CertainType, name string, base ([] int)) (typsys.CertainType, **program.Function, int, ([] int), bool) {
	var def, args, ok = getTypeDefArgs(t.Type, ctx)
	if !(ok) { goto NG }
	if program.T_InvalidType_(t.Type) { goto NG }
	if hdr, ok := ctx.context.lookupMethod(def.Ref, name); ok {
		var out_t_ = typsys.Inflate(hdr.output.type_, def.Parameters, args)
		var out_t = typsys.CertainType { Type: out_t_ }
		var ns = def.Ref.Namespace
		var key = userFunKey {
			name:  name,
			assoc: def.Ref.ItemName,
		}
		var f = ctx.fragment.referUserFun(ns, key)
		return out_t, f, -1, nil, true
	}
	if I, ok := def.Content.(typsys.Interface); ok {
		var index, exists = I.FieldIndexMap[name]
		if exists {
			var field = I.FieldList[index]
			var field_t_ = inflateFieldType(field, def, args)
			var field_t = typsys.CertainType { Type: field_t_ }
			return field_t, nil, index, nil, true
		} else {
			for i, child_ref := range def.Interfaces {
				var child_t = typsys.CertainType {
					Type: typsys.RefType {
						Def:  child_ref,
						Args: args,
					},
				}
				var method_t, f, index, path, ok =
					ctx.resolveMethod(child_t, name, append(base, i))
				if ok {
					return method_t, f, index, path, true
				}
			}
		}
	}
	NG:
	return typsys.CertainType {}, nil, -1, nil, false
}

func checkExpr(expr ast.Expr, cc *exprCheckContext) (*program.Expr, *source.Error) {
	var pipeline = make([] ast.VariousPipe, 0)
	pipeline = append(pipeline, expr.Pipeline...)
	for _, cast := range expr.Casts {
		pipeline = append(pipeline, ast.VariousPipe {
			Node: cast.Node,
			Pipe: ast.PipeCast {
				Node: cast.Node,
				Cast: cast,
			},
		})
	}
	var n = len(pipeline)
	if n == 0 {
		return checkTerm(expr.Term, cc)
	} else {
		var last_index = (n - 1)
		var before_last = pipeline[:last_index]
		var last = pipeline[last_index]
		var rest = ast.Expr {
			Node:     expr.Node, // use a location that spans the whole expr
			Term:     expr.Term,
			Pipeline: before_last,
		}
		return checkPipe(rest, last, cc)
	}
}
func checkTerm(term ast.VariousTerm, cc *exprCheckContext) (*program.Expr, *source.Error) {
	switch T := term.Term.(type) {
	case ast.Lambda:
		return checkLambda(T, cc)
	case ast.Block:
		return checkBlock(T, cc)
	case ast.InfixTerm:
		return checkInfixTerm(T, cc)
	case ast.RefTerm:
		return checkRefTerm(T, cc)
	case ast.ImplicitRefTerm:
		return checkImplicitRefTerm(T, cc)
	case ast.If:
		return checkIf(T, cc)
	case ast.When:
		return checkWhen(T, cc)
	case ast.Each:
		return checkEach(T, cc)
	case ast.Int:
		return checkInt(T, cc)
	case ast.Float:
		return checkFloat(T, cc)
	case ast.Char:
		return checkChar(T, cc)
	case ast.Bytes:
		return checkBytes(T, cc)
	case ast.String:
		return checkString(T, cc)
	default:
		panic("impossible branch")
	}
}
func checkPipe(in ast.Expr, pipe ast.VariousPipe, cc *exprCheckContext) (*program.Expr, *source.Error) {
	switch P := pipe.Pipe.(type) {
	case ast.PipeCast:
		return checkCast(P.Cast, in, cc)
	case ast.PipeGet:
		return checkPipeGet(P, in, cc)
	case ast.PipeInterior:
		return checkPipeInterior(P, in, cc)
	case ast.PipeInfix:
		return checkPipeInfix(P, in, cc)
	case ast.VariousPipeCall:
		return checkPipeCall(P, in, cc)
	default:
		panic("impossible branch")
	}
}

type exprCheckContext struct {
	expected   typsys.Type
	inferring  **typsys.InferringState
	context    *exprContext
}
func createExprCheckContext(expected typsys.Type, context *exprContext) *exprCheckContext {
	return &exprCheckContext {
		expected:  expected,
		inferring: nil,
		context:   context,
	}
}
func (cc *exprCheckContext) describeExpected() string {
	if cc.expected == nil {
		panic("invalid argument")
	} else {
		var s = cc.getInferringState()
		return typsys.DescribeWithInferringState(cc.expected, s)
	}
}
func (cc *exprCheckContext) getCertainOrInferred(t typsys.Type) (typsys.CertainType, bool) {
	if t == nil {
		panic("invalid argument")
	} else {
		var s = cc.getInferringState()
		return typsys.GetCertainOrInferred(t, s)
	}
}
func (cc *exprCheckContext) getExpectedCertainOrInferred() (typsys.CertainType, bool) {
	var nil_t typsys.CertainType
	if cc.expected == nil {
		return nil_t, false
	} else {
		return cc.getCertainOrInferred(cc.expected)
	}
}
func (cc *exprCheckContext) error(loc source.Location, e source.ErrorContent) (*program.Expr, *source.Error) {
	return nil, source.MakeError(loc, e)
}
func (cc *exprCheckContext) expect(expected typsys.Type) *exprCheckContext {
	return &exprCheckContext {
		expected:  expected,
		inferring: cc.inferring,
		context:   cc.context,
	}
}
func (cc *exprCheckContext) use(context *exprContext) *exprCheckContext {
	return &exprCheckContext {
		expected:  cc.expected,
		inferring: cc.inferring,
		context:   context,
	}
}
func (cc *exprCheckContext) try() (*exprCheckContext, **typsys.InferringState) {
	var S **typsys.InferringState
	if cc.inferring != nil {
		S = new(*typsys.InferringState)
		*(S) = *(cc.inferring)
	}
	return &exprCheckContext {
		expected:  cc.expected,
		inferring: S,
		context:   cc.context,
	}, S
}
func (cc *exprCheckContext) withConstScope() *exprCheckContext {
	return cc.use(cc.context.copyWithNewConstScope())
}
func (cc *exprCheckContext) withBlockScope() *exprCheckContext {
	return cc.use(cc.context.copyWithNewBlockScope())
}
func (cc *exprCheckContext) withCaptureScope() (*exprCheckContext, func()([]*program.Binding)) {
	var new_ctx, capture = cc.context.copyWithNewCaptureScope()
	return cc.use(new_ctx), capture
}
func (cc *exprCheckContext) forwardTo(expr ast.Expr) (*program.Expr, *source.Error) {
	return checkExpr(expr, cc)
}
func (cc *exprCheckContext) getInferringState() *typsys.InferringState {
	if cc.inferring != nil {
		return *(cc.inferring)
	} else {
		return nil
	}
}
func (cc *exprCheckContext) updateInferringState(s *typsys.InferringState) {
	if cc.inferring != nil {
		*(cc.inferring) = s
	}
}
func (cc *exprCheckContext) loadTrialInferringState(S **typsys.InferringState) {
	if cc.inferring != nil {
	if S != nil {
		*(cc.inferring) = *S
	}}
}
func (cc *exprCheckContext) checkChildExpr (
	expected  typsys.Type,
	expr_     ast.Expr,
) (*program.Expr, *source.Error) {
	return cc.checkChildExpr3(expected, expr_, false)
}
func (cc *exprCheckContext) checkChildExpr3 (
	expected  typsys.Type,
	expr_     ast.Expr,
	const_    bool,
) (*program.Expr, *source.Error) {
	if const_ {
		// noinspection ALL
		cc = cc.withConstScope()
	}
	var expr, err = checkExpr(expr_, cc.expect(expected))
	if err != nil { return nil, err }
	if const_ {
		expr = program.Cached(expr)
	}
	return expr, nil
}
func (cc *exprCheckContext) assignChildExpr (
	expected  typsys.Type,
	expr      *program.Expr,
) (*program.Expr, *source.Error) {
	if expected == nil {
		return expr, nil
	} else {
		var ctx = cc.context
		var s0 = cc.getInferringState()
		var expr, s1, err = assign(expected, expr, ctx, s0)
		if err != nil { return nil, err }
		cc.updateInferringState(s1)
		return expr, nil
	}
}
func (cc *exprCheckContext) assign (
	t        typsys.CertainType,
	loc      source.Location,
	content  program.ExprContent,
) (*program.Expr, *source.Error) {
	var info = program.ExprInfoFrom(loc)
	var expr0 = &program.Expr {
		Type:    t,
		Info:    info,
		Content: content,
	}
	return cc.assignChildExpr(cc.expected, expr0)
}
func (cc *exprCheckContext) tryAssign (
	t        typsys.CertainType,
	loc      source.Location,
	content  program.ExprContent,
) (*program.Expr, **typsys.InferringState, *source.Error) {
	var trial, S = cc.try()
	var expr, err = trial.assign(t, loc, content)
	return expr, S, err
}
func (cc *exprCheckContext) infer (
	params  [] string,
	args    [] typsys.CertainType,
	output  typsys.Type,
	loc     source.Location,
	check   func(cc *exprCheckContext) (program.ExprContent, *source.Error),
) (*program.Expr, *source.Error) {
	var s = typsys.Infer(params)
	var inner_cc = &exprCheckContext {
		inferring: &s,
		context:   cc.context,
	}
	output = typsys.ToInferring(output, params)
	if len(args) > len(params) {
		return cc.error(loc, E_TooManyTypeArgs {})
	}
	for i := range args {
		var A = args[i].Type
		var P = typsys.Type(typsys.ParameterType { Name: params[i] })
		P = typsys.ToInferring(P, params)
		var ok, s_ = typsys.Match(P, A, s)
		if !(ok) { panic("something went wrong") }
		s = s_
	}
	if expected_t, ok := cc.getExpectedCertainOrInferred(); ok {
		if ok, s_ := typsys.Match(expected_t.Type, output, s); ok {
			s = s_
			var content, err = check(inner_cc)
			if err != nil { return nil, err }
			return &program.Expr {
				Type:    expected_t,
				Info:    program.ExprInfoFrom(loc),
				Content: content,
			}, nil
		} else {
			var content, err = check(inner_cc)
			if err != nil { return nil, err }
			var output_t, ok = typsys.GetCertainOrInferred(output, s)
			if !(ok) {
				var ctx = cc.context
				var m, sam_ok = getSamInterfaceMethodType(expected_t.Type, ctx)
				if sam_ok {
				if sam_ok, s_ := typsys.Match(m, output, s); sam_ok {
					s = s_
					output_t, ok = typsys.CertainType { Type: m }, true
				}}
			}
			if ok {
				return cc.assign(output_t, loc, content)
			} else {
				return cc.error(loc, E_ExpectSufficientTypeArguments {})
			}
		}
	} else {
		var content, err = check(inner_cc)
		if err != nil { return nil, err }
		if output_t, ok := typsys.GetCertainOrInferred(output, s); ok {
			return cc.assign(output_t, loc, content)
		} else {
			return cc.error(loc, E_ExpectExplicitTypeCast {})
		}
	}
}
func (cc *exprCheckContext) tryInfer (
	params  [] string,
	args    [] typsys.CertainType,
	output  typsys.Type,
	loc     source.Location,
	check   func(cc *exprCheckContext) (program.ExprContent, *source.Error),
) (*program.Expr, **typsys.InferringState, *source.Error) {
	var trial, S = cc.try()
	var expr, err = trial.infer(params, args, output, loc, check)
	return expr, S, err
}


const Underscore = "_"

func (cc *exprCheckContext) match (
	pattern  ast.MaybePattern,
	value_t  typsys.CertainType,
) (program.PatternMatching, *source.Error) {
	return cc.match3(pattern, value_t, false)
}

func (cc *exprCheckContext) match3 (
	pattern  ast.MaybePattern,
	value_t  typsys.CertainType,
	const_   bool,
) (program.PatternMatching, *source.Error) {
	var result = make(program.PatternMatching, 0)
	if pattern, ok := pattern.(ast.VariousPattern); ok {
		switch P := pattern.Pattern.(type) {
		case ast.PatternSingle:
			var name = ast.Id2String(P.Name)
			var loc = P.Location
			if name == Underscore {
				break
			}
			var b = cc.context.createBinding4(name, value_t, loc, const_)
			result = append(result, program.PatternMatchingItem {
				Binding: b,
				Index1:  0,
			})
		case ast.PatternMultiple:
			var loc = P.Location
			var record, def, args, ok = getRecord(value_t.Type, cc.context)
			if !(ok) {
				return nil, source.MakeError(loc,
					E_CannotMatchRecord {
						TypeDesc: typsys.DescribeCertain(value_t),
					})
			}
			var record_size = len(record.FieldList)
			var pattern_arity = len(P.Names)
			if pattern_arity != record_size {
				return nil, source.MakeError(loc,
					E_RecordSizeNotMatching {
						PatternArity: pattern_arity,
						RecordSize:   record_size,
						RecordDesc:   typsys.DescribeCertain(value_t),
					})
			}
			var occurred = make(map[string] struct{})
			for i, binding_node := range P.Names {
				var field = record.FieldList[i]
				var field_t_ = inflateFieldType(field, def, args)
				var field_t = typsys.CertainType { Type: field_t_ }
				var name = ast.Id2String(binding_node)
				var loc = binding_node.Location
				if name == Underscore {
					continue
				}
				var _, duplicate = occurred[name]
				occurred[name] = struct{}{}
				if duplicate {
					return nil, source.MakeError(loc,
						E_DuplicateBinding { name })
				}
				var b = cc.context.createBinding4(name, field_t, loc, const_)
				result = append(result, program.PatternMatchingItem {
					Binding: b,
					Index1:  (1 + i),
				})
			}
		default:
			panic("impossible branch")
		}
	}
	return result, nil
}


