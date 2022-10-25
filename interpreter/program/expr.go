package program

import (
	"regexp"
	"math/big"
	"rxgui/util/ctn"
	"rxgui/interpreter/core"
	"rxgui/interpreter/lang/source"
	"rxgui/interpreter/lang/typsys"
)


type ExprBasedFunctionValue struct {
	InputArgs  [] *Binding
	InputCtx   [] *Binding
	Output     *Expr
}
func (v *ExprBasedFunctionValue) Call(args ([] core.Object), ctx ([] core.Object), h core.RuntimeHandle) core.Object {
	var eval_ctx = CreateEvalContext(h)
	eval_ctx = eval_ctx.NewCtxBindAll(v.InputArgs, args)
	eval_ctx = eval_ctx.NewCtxBindAll(v.InputCtx, ctx)
	return eval_ctx.Eval(v.Output)
}

type Expr struct {
	Type     typsys.CertainType
	Info     ExprInfo
	Content  ExprContent
}
type ExprInfo struct {
	Location  source.Location
}
func ExprInfoFrom(loc source.Location) ExprInfo {
	return ExprInfo { Location: loc }
}
type ExprContent interface {
	Eval(ctx *EvalContext) core.Object
}

type EvalContext struct {
	bindingValues  ctn.Map[uintptr, core.Object]
	runtimeHandle  core.RuntimeHandle
}
func CreateEvalContext(h core.RuntimeHandle) *EvalContext {
	return &EvalContext {
		bindingValues: makeBindingValuesMap(),
		runtimeHandle: h,
	}
}
func makeBindingValuesMap() ctn.Map[uintptr, core.Object] {
	var cmp = ctn.DefaultCompare[uintptr]
	return ctn.MakeMap[uintptr, core.Object](cmp)
}
func (ctx *EvalContext) Eval(expr *Expr) core.Object {
	return expr.Content.Eval(ctx)
}
func (ctx *EvalContext) EvalAll(exprList ([] *Expr)) ([] core.Object) {
	var objects = make([] core.Object, len(exprList))
	for i := range exprList {
		objects[i] = ctx.Eval(exprList[i])
	}
	return objects
}
func (ctx *EvalContext) EvalAllAsList(exprList ([] *Expr)) core.List {
	var nodes = make([] core.ListNode, len(exprList))
	for i := range exprList {
		nodes[i].Value = ctx.Eval(exprList[i])
	}
	return core.NodesToList(nodes)
}
func (ctx *EvalContext) EvalBinding(binding *Binding) core.Object {
	var ptr = binding.PointerNumber()
	var obj, _ = ctx.bindingValues.Lookup(ptr)
	return obj
}
func (ctx *EvalContext) NewCtx() *EvalContext {
	return &EvalContext {
		bindingValues: makeBindingValuesMap(),
		runtimeHandle: ctx.runtimeHandle,
	}
}
func (ctx *EvalContext) NewCtxMatch(p PatternMatching, obj core.Object) *EvalContext {
	var m = ctx.bindingValues
	for _, item := range p {
		var v core.Object
		if item.Index1 == 0 {
			v = obj
		} else {
			var index = (item.Index1 - 1)
			var r = (*obj).(core.Record)
			v = r.Objects[index]
		}
		var k = item.Binding.PointerNumber()
		m = m.Inserted(k, v)
	}
	return &EvalContext {
		bindingValues: m,
		runtimeHandle: ctx.runtimeHandle,
	}
}
func (ctx *EvalContext) NewCtxUnbind(binding *Binding) *EvalContext {
	var m = ctx.bindingValues
	var ptr = binding.PointerNumber()
	_, m, _ = m.Deleted(ptr)
	return &EvalContext {
		bindingValues: m,
		runtimeHandle: ctx.runtimeHandle,
	}
}
func (ctx *EvalContext) NewCtxBind(binding *Binding, obj core.Object) *EvalContext {
	var m = ctx.bindingValues
	var ptr = binding.PointerNumber()
	m = m.Inserted(ptr, obj)
	return &EvalContext {
		bindingValues: m,
		runtimeHandle: ctx.runtimeHandle,
	}
}
func (ctx *EvalContext) NewCtxBindAll(bindings ([] *Binding), objects ([] core.Object)) *EvalContext {
	var m = ctx.bindingValues
	for i, b := range bindings {
		var ptr = b.PointerNumber()
		var obj = objects[i]
		m = m.Inserted(ptr, obj)
	}
	return &EvalContext {
		bindingValues: m,
		runtimeHandle: ctx.runtimeHandle,
	}
}
func (ctx *EvalContext) NewCtxCapture(bindings ([] *Binding)) *EvalContext {
	var new_ctx = ctx.NewCtx()
	var m = &(new_ctx.bindingValues)
	for _, b := range bindings {
		var ptr = b.PointerNumber()
		var obj, _ = ctx.bindingValues.Lookup(ptr)
		*m = m.Inserted(ptr, obj)
	}
	return new_ctx
}

type PatternMatching ([] PatternMatchingItem)
type PatternMatchingItem struct {
	Binding  *Binding
	Index1   int  // 0 = whole, 1 = .0
}


type Wrapper struct {
	Inner *Expr
}
func (expr Wrapper) Eval(ctx *EvalContext) core.Object {
	return ctx.Eval(expr.Inner)
}

type CachedExpr struct {
	Expr   *Expr
	Cache  *ExprValueCache
}
type ExprValueCache struct {
	available  bool
	value      core.Object
}
func (expr CachedExpr) Eval(ctx *EvalContext) core.Object {
	if expr.Cache.available {
		return expr.Cache.value
	} else {
		var value = ctx.Eval(expr.Expr)
		expr.Cache.value = value
		expr.Cache.available = true
		return value
	}
}
func Cached(expr *Expr) *Expr {
	return &Expr {
		Type:    expr.Type,
		Info:    expr.Info,
		Content: CachedExpr {
			Expr:  expr,
			Cache: new(ExprValueCache),
		},
	}
}

type FunRef struct {
	Function  **Function
	Context   [] *Expr
	Unpack    bool
}
func (expr FunRef) Eval(ctx *EvalContext) core.Object {
	return core.FunctionToLambdaObject (
		(*expr.Function).value,
		expr.Unpack,
		ctx.EvalAll(expr.Context),
		ctx.runtimeHandle,
	)
}

type LocalRef struct {
	Binding *Binding
}
func (expr LocalRef) Eval(ctx *EvalContext) core.Object {
	return ctx.EvalBinding(expr.Binding)
}

type CallLambda struct {
	Callee    *Expr
	Argument  *Expr
}
func (expr CallLambda) Eval(ctx *EvalContext) core.Object {
	var callee = ctx.Eval(expr.Callee)
	var lambda = (*callee).(core.Lambda)
	var argument = ctx.Eval(expr.Argument)
	return lambda.Call(argument)
}

type CallFunction struct {
	Location   source.Location
	Callee     **Function
	Context    [] *Expr
	Arguments  [] *Expr
}
func (expr CallFunction) Eval(ctx *EvalContext) core.Object {
	var f = *(expr.Callee)
	var f_ctx = ctx.EvalAll(expr.Context)
	var args = ctx.EvalAll(expr.Arguments)
	var h = core.AddFrameInfo(ctx.runtimeHandle, f.name, expr.Location)
	return f.value.Call(args, f_ctx, h)
}

type Interface struct {
	ConcreteValue  *Expr
	DispatchTable  **DispatchTable
}
func (expr Interface) Eval(ctx *EvalContext) core.Object {
	return core.Obj(core.Interface {
		UnderlyingObject: ctx.Eval(expr.ConcreteValue),
		DispatchTable:    (*expr.DispatchTable).value(),
	})
}

type InterfaceTransformUpward struct {
	Arg   *Expr
	Path  [] int
}
func (expr InterfaceTransformUpward) Eval(ctx *EvalContext) core.Object {
	var arg = ctx.Eval(expr.Arg)
	var I = (*arg).(core.Interface)
	var table = I.DispatchTable
	for _, index := range expr.Path {
		table = table.Children[index]
	}
	return core.Obj(core.Interface {
		UnderlyingObject: I.UnderlyingObject,
		DispatchTable:    table,
	})
}

type InterfaceTransformDownward struct {
	Arg     *Expr
	Depth   int
	Target  string
}
func (expr InterfaceTransformDownward) Eval(ctx *EvalContext) core.Object {
	var arg = ctx.Eval(expr.Arg)
	var I = (*arg).(core.Interface)
	var table = I.DispatchTable
	for i := 0; i < expr.Depth; i += 1 {
		table = table.Parent
		if table == nil {
			break
		}
	}
	if ((table != nil) && (table.Interface == expr.Target)) {
		return core.Just(core.Obj(core.Interface {
			UnderlyingObject: I.UnderlyingObject,
			DispatchTable:    table,
		}))
	} else {
		return core.Nothing()
	}
}

type InterfaceFromSamValue struct {
	Value  *Expr
}
func (expr InterfaceFromSamValue) Eval(ctx *EvalContext) core.Object {
	var value = ctx.Eval(expr.Value)
	return core.Obj(core.CraftSamInterface(value))
}

type Null struct {}
func (expr Null) Eval(_ *EvalContext) core.Object {
	return nil
}

type Enum int
func (expr Enum) Eval(_ *EvalContext) core.Object {
	var index = int(expr)
	return core.Obj(core.Enum(index))
}

type EnumToInt struct {
	EnumValue  *Expr
}
func (expr EnumToInt) Eval(ctx *EvalContext) core.Object {
	var enum_obj = ctx.Eval(expr.EnumValue)
	var n = big.NewInt(int64(int((*enum_obj).(core.Enum))))
	return core.Obj(core.Int { Value: n })
}

type Union struct {
	Index  int
	Value  *Expr
}
func (expr Union) Eval(ctx *EvalContext) core.Object {
	return core.Obj(core.Union {
		Index:  expr.Index,
		Object: ctx.Eval(expr.Value),
	})
}

type Record struct {
	Values  [] *Expr
}
func (expr Record) Eval(ctx *EvalContext) core.Object {
	return core.Obj(core.Record {
		Objects: ctx.EvalAll(expr.Values),
	})
}

type FieldValue struct {
	Record  *Expr
	Index   int
}
func (expr FieldValue) Eval(ctx *EvalContext) core.Object {
	var obj = ctx.Eval(expr.Record)
	var record = (*obj).(core.Record)
	return record.Objects[expr.Index]
}

type ObservableFieldProjection struct {
	Base   *Expr
	Index  int
}
func (expr ObservableFieldProjection) Eval(ctx *EvalContext) core.Object {
	var index = expr.Index
	var base = ctx.Eval(expr.Base)
	var raw = core.GetObservable(base).Map(func(obj core.Object) core.Object {
		var record = (*obj).(core.Record)
		return record.Objects[index]
	})
	return core.Obj(raw.DistinctUntilObjectChanged())
}

type ConcreteMethodValue struct {
	Location  source.Location
	This      *Expr
	Path      [] int
	Method    **Function
}
func (expr ConcreteMethodValue) Eval(ctx *EvalContext) core.Object {
	var this = ctx.Eval(expr.This)
	if len(expr.Path) > 0 {
		var I = (*this).(core.Interface)
		var obj = I.UnderlyingObject
		var table = I.DispatchTable
		for _, index := range expr.Path {
			table = table.Children[index]
		}
		this = core.Obj(core.Interface {
			UnderlyingObject: obj,
			DispatchTable:    table,
		})
	}
	var args = [] core.Object { this }
	var f = (*expr.Method)
	var h = core.AddFrameInfo(ctx.runtimeHandle, f.name, expr.Location)
	return f.value.Call(args, nil, h)
}

type AbstractMethodValue struct {
	Location   source.Location
	Interface  *Expr
	Path       [] int
	Index      int
}
func (expr AbstractMethodValue) Eval(ctx *EvalContext) core.Object {
	var interface_object = ctx.Eval(expr.Interface)
	var I = (*interface_object).(core.Interface)
	var this = I.UnderlyingObject
	var table = I.DispatchTable
	for _, index := range expr.Path {
		table = table.Children[index]
	}
	var f = *(table.Methods[expr.Index])
	var args = [] core.Object { this }
	var h = core.AddFrameInfo(ctx.runtimeHandle, "(dynamic)", expr.Location)
	return f.Call(args, nil, h)
}

type List struct {
	Items  [] *Expr
}
func (expr List) Eval(ctx *EvalContext) core.Object {
	return core.Obj(ctx.EvalAllAsList(expr.Items))
}

type InteriorRef struct {
	Base     *Expr
	Index    int
	Table    **DispatchTable  // optional, used in dynamic cast
	Kind     InteriorRefKind
	Operand  InteriorRefOperand
}
type InteriorRefKind int
const (
	RK_RecordField InteriorRefKind = iota
	RK_EnumItem
	RK_UnionItem
	RK_DynamicCast
)
type InteriorRefOperand int
const (
	RO_Direct InteriorRefOperand = iota
	RO_Lens1
	RO_Lens2
)
func (expr InteriorRef) Eval(ctx *EvalContext) core.Object {
	var base = expr.Base.Content.Eval(ctx)
	switch expr.Kind {
	case RK_RecordField:
		i := expr.Index
		switch expr.Operand {
		case RO_Direct: return core.Lens1FromRecord(base, i)
		case RO_Lens1:  return core.Lens1FromRecordLens1(base, i)
		}
	case RK_EnumItem:
		i := expr.Index
		switch expr.Operand {
		case RO_Direct: return core.Lens2FromEnum(base, i)
		case RO_Lens1:  return core.Lens2FromEnumLens1(base, i)
		case RO_Lens2:  return core.Lens2FromEnumLens2(base, i)
		}
	case RK_UnionItem:
		i := expr.Index
		switch expr.Operand {
		case RO_Direct: return core.Lens2FromUnion(base, i)
		case RO_Lens1:  return core.Lens2FromUnionLens1(base, i)
		case RO_Lens2:  return core.Lens2FromUnionLens2(base, i)
		}
	case RK_DynamicCast:
		t := (*expr.Table).value()
		switch expr.Operand {
		case RO_Direct: return core.Lens2FromInterface(base, t)
		case RO_Lens1:  return core.Lens2FromInterfaceLens1(base, t)
		case RO_Lens2:  return core.Lens2FromInterfaceLens2(base, t)
		}
	}
	panic("something went wrong")
}

type IntLiteral struct {
	Value *big.Int
}
func (expr IntLiteral) Eval(_ *EvalContext) core.Object {
	return core.Obj(core.Int { Value: expr.Value })
}

type CharLiteral struct {
	Value rune
}
func (expr CharLiteral) Eval(_ *EvalContext) core.Object {
	return core.Obj(core.Char(expr.Value))
}

type FloatLiteral struct {
	Value float64
}
func (expr FloatLiteral) Eval(_ *EvalContext) core.Object {
	return core.Obj(core.Float(expr.Value))
}

type BytesLiteral struct {
	Value  [] byte
}
func (expr BytesLiteral) Eval(_ *EvalContext) core.Object {
	return core.Obj(core.Bytes(expr.Value))
}

type StringLiteral struct {
	Value  string
}
func (expr StringLiteral) Eval(_ *EvalContext) core.Object {
	return core.Obj(core.String(expr.Value))
}

type RegexpLiteral struct {
	Value  *regexp.Regexp
}
func (expr RegexpLiteral) Eval(_ *EvalContext) core.Object {
	return core.Obj(core.RegExp { Value: expr.Value })
}

type Lambda struct {
	Ctx   [] *Binding
	In    PatternMatching
	Out   *Expr
	Self  *Binding
}
func (expr Lambda) Eval(ctx *EvalContext) core.Object {
	var capture = ctx.NewCtxCapture(expr.Ctx)
	var input_pattern = expr.In
	var output_expr = expr.Out
	var lambda = core.Obj(core.Lambda { Call: func(arg core.Object) core.Object {
		var inner = capture.NewCtxMatch(input_pattern, arg)
		return inner.Eval(output_expr)
	}})
	if expr.Self != nil {
		capture = capture.NewCtxBind(expr.Self, lambda)
	}
	return lambda
}

type Let struct {
	In   PatternMatching
	Arg  *Expr
	Out  *Expr
}
func (expr Let) Eval(ctx *EvalContext) core.Object {
	var arg = ctx.Eval(expr.Arg)
	var inner = ctx.NewCtxMatch(expr.In, arg)
	return inner.Eval(expr.Out)
}

type If struct {
	Branches  [] IfBranch
}
type IfBranch struct {
	Conds  [] Cond
	Value  *Expr
}
type Cond struct {
	Kind   CondKind
	Match  PatternMatching
	Value  *Expr
}
type CondKind int
const (
	CK_Bool CondKind = iota
	CK_Maybe
	CK_Lens2
)
func (expr If) Eval(ctx *EvalContext) core.Object {
	for _, branch := range expr.Branches {
		var branch_ctx = ctx
		var branch_ok = true
		for _, cond := range branch.Conds {
			var v = branch_ctx.Eval(cond.Value)
			switch cond.Kind {
			case CK_Bool:
				if core.GetBool(v) {
					continue
				}
			case CK_Maybe:
				if v, ok := core.UnwrapMaybe(v); ok {
					branch_ctx = branch_ctx.NewCtxMatch(cond.Match, v)
					continue
				}
			case CK_Lens2:
				if v, ok := core.UnwrapLens2(v); ok {
					branch_ctx = branch_ctx.NewCtxMatch(cond.Match, v)
					continue
				}
			default:
				panic("impossible branch")
			}
			branch_ok = false
			break
		}
		if branch_ok {
			return branch_ctx.Eval(branch.Value)
		}
	}
	panic("bad if expression")
}

type When struct {
	Operand   WhenOperand
	Branches  [] *WhenBranch
}
type WhenBranch struct {
	Match  PatternMatching
	Value  *Expr
}
type WhenOperand struct {
	Kind   UnionOrEnum
	Value  *Expr
}
type UnionOrEnum int
const (
	UE_Union UnionOrEnum = iota
	UE_Enum
)
func (expr When) Eval(ctx *EvalContext) core.Object {
	switch expr.Operand.Kind {
	case UE_Union:
		var operand = ctx.Eval(expr.Operand.Value)
		var u = (*operand).(core.Union)
		var branch = expr.Branches[u.Index]
		var inner = ctx.NewCtxMatch(branch.Match, u.Object)
		return inner.Eval(branch.Value)
	case UE_Enum:
		var operand = ctx.Eval(expr.Operand.Value)
		var index = int((*operand).(core.Enum))
		var branch = expr.Branches[index]
		return ctx.Eval(branch.Value)
	default:
		panic("impossible branch")
	}
}

type EachValue struct {
	Kind   UnionOrEnum
	Index  int
	Match  PatternMatching
	Value  *Expr
}
func (expr EachValue) Eval(ctx *EvalContext) core.Object {
	switch expr.Kind {
	case UE_Union:
		var index = expr.Index
		var conv = core.Obj(core.Lambda { Call: func(arg core.Object) core.Object {
			return core.Obj(core.Union {
				Index:  index,
				Object: arg,
			})
		}})
		var inner = ctx.NewCtxMatch(expr.Match, conv)
		return inner.Eval(expr.Value)
	case UE_Enum:
		var index = expr.Index
		var enum_obj = core.Obj(core.Enum(index))
		var inner = ctx.NewCtxMatch(expr.Match, enum_obj)
		return inner.Eval(expr.Value)
	default:
		panic("impossible branch")
	}
}

type ObservableRecord struct {
	Record  Record
}
type ObservableMaybeRecord struct {
	Record  Record
}
func (expr ObservableRecord) Eval(ctx *EvalContext) core.Object {
	var values = ctx.EvalAll(expr.Record.Values)
	var observables = ctn.MapEach(values, core.GetObservable)
	return core.Obj(core.CombineLatest(observables...).Map(func(combined core.Object) core.Object {
		var objects = make([] core.Object, len(values))
		core.GetList(combined).ForEachWithIndex(func(i int, obj core.Object) {
			objects[i] = obj
		})
		return core.Obj(core.Record { Objects: objects })
	}))
}
func (expr ObservableMaybeRecord) Eval(ctx *EvalContext) core.Object {
	var combined = (ObservableRecord { expr.Record }).Eval(ctx)
	return core.Obj(core.GetObservable(combined).Map(func(obj core.Object) core.Object {
		var r = (*obj).(core.Record)
		var unwrapped = make([] core.Object, len(r.Objects))
		for i := range r.Objects {
			if inner, ok := core.UnwrapMaybe(r.Objects[i]); ok {
				unwrapped[i] = inner
			} else {
				return core.Nothing()
			}
		}
		return core.Just(core.Obj(core.Record { Objects: unwrapped }))
	}))
}
type ObservableRecordDefault struct {
	Content  ExprContent
}
type ObservableMaybeRecordDefault struct {
	Content  ExprContent
}
func (expr ObservableRecordDefault) Eval(ctx *EvalContext) core.Object {
	var v = expr.Content.Eval(ctx)
	return core.Obj(core.ObservableSyncValue(v))
}
func (expr ObservableMaybeRecordDefault) Eval(ctx *EvalContext) core.Object {
	var v = expr.Content.Eval(ctx)
	return core.Obj(core.ObservableSyncValue(core.Just(v)))
}

type HookRecord struct {
	Record  Record
}
func (expr HookRecord) Eval(ctx *EvalContext) core.Object {
	var values = ctx.EvalAll(expr.Record.Values)
	var jobs = ctn.MapEach(values, func(v core.Object) core.Observable {
		return core.FromObject[core.Hook](v).Job
	})
	var result = core.ObservableSyncValue(core.Obj(core.EmptyList()))
	var L = len(jobs)
	for i := (L-1); i >= 0; i -= 1 {
		result = (func(current core.Observable, next core.Observable) core.Observable {
			return current.AwaitNoexcept(ctx.runtimeHandle, func(head core.Object) core.Observable {
			return next.AwaitNoexcept(ctx.runtimeHandle, func(tail_obj core.Object) core.Observable {
				var tail = core.GetList(tail_obj)
				var list = core.Cons(head, tail)
				return core.ObservableSyncValue(core.Obj(list))
			})})
		})(jobs[i], result)
	}
	{ var result = result.Map(func(obj core.Object) core.Object {
		var list = core.GetList(obj)
		var objects = make([] core.Object, list.Length())
		list.ForEachWithIndex(func(i int, object core.Object) {
			objects[i] = object
		})
		return core.Obj(core.Record { Objects: objects })
	})
	return core.ToObject(core.Hook { Job: result }) }
}
type HookRecordDefault struct {
	Content  ExprContent
}
func (expr HookRecordDefault) Eval(ctx *EvalContext) core.Object {
	var v = expr.Content.Eval(ctx)
	var o = core.ObservableSyncValue(v)
	return core.ToObject(core.Hook { Job: o })
}

type ReflectType struct {
	Type  ReflectType_
}
func (expr ReflectType) Eval(_ *EvalContext) core.Object {
	var t = core.ReflectType(expr.Type)
	return core.Obj(t)
}

type ReflectValue struct {
	Type   ReflectType_
	Value  *Expr
}
func (expr ReflectValue) Eval(ctx *EvalContext) core.Object {
	var t = core.ReflectType(expr.Type)
	var v = ctx.Eval(expr.Value)
	return core.Obj(core.AssumeValidReflectValue(t, v))
}


