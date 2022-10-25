package compiler

import (
	"strings"
	"rxgui/util/ctn"
	"rxgui/interpreter/lang/source"
	"rxgui/interpreter/lang/typsys"
	"rxgui/interpreter/lang/textual/ast"
	"rxgui/interpreter/program"
)


type Impl struct {
	namespace   string
	userFunMap  ctn.Map[userFunKey, funImpl]
	genFunMap   ctn.Map[genFunKey, funImpl]
}
type funImpl interface { impl(funImpl) }
type userFunKey struct {
	name   string
	assoc  string
}
func (key userFunKey) Describe(ns string) string {
	var buf strings.Builder
	if ns != "" {
		buf.WriteString(ns)
		buf.WriteString("::")
	}
	buf.WriteString(key.name)
	if key.assoc != "" {
		buf.WriteRune('(')
		buf.WriteString(key.assoc)
		buf.WriteRune(')')
	}
	return buf.String()
}
func userFunKeyCompare(a userFunKey, b userFunKey) ctn.Ordering {
	var o = ctn.StringCompare(a.name, b.name)
	if o != ctn.Equal {
		return o
	} else {
		return ctn.StringCompare(a.assoc, b.assoc)
	}
}
type external[Key any] struct {
	namespace    string
	internalKey  Key
}
func externalCompare[Key any] (cmp ctn.Compare[Key]) func(external[Key],external[Key])(ctn.Ordering) {
	return func(a external[Key], b external[Key]) ctn.Ordering {
		var o = ctn.StringCompare(a.namespace, b.namespace)
		if o != ctn.Equal {
			return o
		} else {
			return cmp(a.internalKey, b.internalKey)
		}
	}
}

func (funImplAstExpr) impl(funImpl) {}
type funImplAstExpr struct {
	Expr  ast.Expr
}
func (funImplLibraryNative) impl(funImpl) {}
type funImplLibraryNative struct {
	Id  string
}
// func (funImplLoadedAsset) impl(funImpl) {}
// type funImplLoadedAsset struct {
//     Path  string
//     Data  [] byte
// }

func funImplFromAstBody(body ast.VariousBody) funImpl {
	switch B := body.Body.(type) {
	case ast.Block:
		return funImplAstExpr {
			Expr: ast.WrapBlockAsExpr(B),
		}
	case ast.NativeBody:
		var id, err = normalizeText(B.Id.Value)
		if err != nil {
			id = "<InvalidNativeFuncId>"
		}
		return funImplLibraryNative {
			Id: id,
		}
	// case ast.LoadedAsset:
	//     return funImplLoadedAsset {
	//         Path: B.Path,
	//         Data: B.Data,
	//     }
	default:
		panic("impossible branch")
	}
}

type Fragment struct {
	namespace  string
	internal   internalMaps
	external   externalMaps
	// TODO: usage map (usageKey -> source.Location)
}
type internalMaps struct {
	dspMap      ctn.Map[dispatchKey, **program.DispatchTable]
	userFunMap  ctn.Map[userFunKey, **program.Function]
	genFunMap   ctn.Map[genFunKey, **program.Function]
}
type externalMaps struct {
	dspMap      ctn.Map[external[dispatchKey], **program.DispatchTable]
	userFunMap  ctn.Map[external[userFunKey], **program.Function]
	genFunMap   ctn.Map[external[genFunKey], **program.Function]
}
type Executable struct {
	content  ctn.Map[string,*internalMaps]
}
func (exe *Executable) LookupEntry(ns string) (**program.Function, bool) {
	if group, ok := exe.content.Lookup(ns); ok {
		var entry_key = userFunKey { name: "" }
		if item, ok := group.userFunMap.Lookup(entry_key); ok {
			return item, true
		}
	}
	return nil, false
}
func linkFragment(exe *Executable, fragment *Fragment) {
	var m = exe.content
	fragment.external.dspMap.ForEach(func(k external[dispatchKey], v **program.DispatchTable) {
		var group, exists = m.Lookup(k.namespace)
		if !(exists) { panic("something went wrong") }
		var item, ok = group.dspMap.Lookup(k.internalKey)
		if !(ok) { panic("something went wrong") }
		*v = *item
	})
	fragment.external.userFunMap.ForEach(func(k external[userFunKey], v **program.Function) {
		var group, exists = m.Lookup(k.namespace)
		if !(exists) { panic("something went wrong") }
		var item, ok = group.userFunMap.Lookup(k.internalKey)
		if !(ok) { panic("something went wrong") }
		*v = *item
	})
	fragment.external.genFunMap.ForEach(func(k external[genFunKey], v **program.Function) {
		var group, exists = m.Lookup(k.namespace)
		if !(exists) { panic("something went wrong") }
		var item, ok = group.genFunMap.Lookup(k.internalKey)
		if !(ok) { panic("something went wrong") }
		*v = *item
	})
}
func link(all ([] *Fragment)) *Executable {
	var m = ctn.MakeMutMap[string,*internalMaps](ctn.StringCompare)
	for _, fragment := range all {
		var ns = fragment.namespace
		var existing, exists = m.Lookup(ns)
		if !(exists) {
			var ptr = &internalMaps {
				dspMap:     ctn.MakeMap[dispatchKey,**program.DispatchTable](dispatchKeyCompare),
				userFunMap: ctn.MakeMap[userFunKey,**program.Function](userFunKeyCompare),
				genFunMap:  ctn.MakeMap[genFunKey,**program.Function](genFunKeyCompare),
			}
			m.Insert(ns, ptr)
			existing = ptr
		}
		fragment.internal.dspMap.ForEach(func(k dispatchKey, v **program.DispatchTable) {
			existing.dspMap = existing.dspMap.Inserted(k, v)
		})
		fragment.internal.userFunMap.ForEach(func(k userFunKey, v **program.Function) {
			existing.userFunMap = existing.userFunMap.Inserted(k, v)
		})
		fragment.internal.genFunMap.ForEach(func(k genFunKey, v **program.Function) {
			existing.genFunMap = existing.genFunMap.Inserted(k, v)
		})
	}
	var exe = &Executable { m.Map() }
	for _, fragment := range all {
		linkFragment(exe, fragment)
	}
	return exe
}
type fragmentDraft struct { content *Fragment }
func makeFragmentDraft(ns string) fragmentDraft {
	var internal_maps = internalMaps {
		dspMap:     ctn.MakeMap[dispatchKey,**program.DispatchTable](dispatchKeyCompare),
		userFunMap: ctn.MakeMap[userFunKey,**program.Function](userFunKeyCompare),
		genFunMap:  ctn.MakeMap[genFunKey,**program.Function](genFunKeyCompare),
	}
	var external_maps = externalMaps {
		dspMap:     ctn.MakeMap[external[dispatchKey],**program.DispatchTable](externalCompare(dispatchKeyCompare)),
		userFunMap: ctn.MakeMap[external[userFunKey],**program.Function](externalCompare(userFunKeyCompare)),
		genFunMap:  ctn.MakeMap[external[genFunKey],**program.Function](externalCompare(genFunKeyCompare)),
	}
	var fragment = &Fragment {
		namespace: ns,
		internal:  internal_maps,
		external:  external_maps,
	}
	return fragmentDraft { content: fragment }
}
func (draft fragmentDraft) namespace() string {
	return draft.content.namespace
}
func (draft fragmentDraft) createDispatchTable(pair dispatchKey) (*program.DispatchTable, bool) {
	var m = &(draft.content.internal.dspMap)
	var _, exists = m.Lookup(pair)
	if exists {
		return nil, false
	} else {
		var t = new(program.DispatchTable)
		var ptr = new(*program.DispatchTable)
		*ptr = t
		*m = m.Inserted(pair, ptr)
		return t, true
	}
}
func (draft fragmentDraft) createOrGetUserFun(key userFunKey) *program.Function {
	var m = &(draft.content.internal.userFunMap)
	var existing, exists = m.Lookup(key)
	if exists {
		return *existing
	} else {
		var f = new(program.Function)
		var ptr = new(*program.Function)
		*ptr = f
		*m = m.Inserted(key, ptr)
		return f
	}
}
func (draft fragmentDraft) createOrGetGenFun(key genFunKey) *program.Function {
	var m = &(draft.content.internal.genFunMap)
	var existing, exists = m.Lookup(key)
	if exists {
		return *existing
	} else {
		var f = new(program.Function)
		var ptr = new(*program.Function)
		*ptr = f
		*m = m.Inserted(key, ptr)
		return f
	}
}
func (draft fragmentDraft) referDispatchTable(ns string, key dispatchKey) **program.DispatchTable {
	var fragment = draft.content
	var ext = external[dispatchKey] {
		namespace:   ns,
		internalKey: key,
	}
	if table, ok := fragment.external.dspMap.Lookup(ext); ok {
		return table
	} else {
		var table = new(*program.DispatchTable)
		fragment.external.dspMap =
			fragment.external.dspMap.Inserted(ext, table)
		return table
	}
}
func (draft fragmentDraft) referUserFun(ns string, key userFunKey) **program.Function {
	var fragment = draft.content
	var ext = external[userFunKey] {
		namespace:   ns,
		internalKey: key,
	}
	if f, ok := fragment.external.userFunMap.Lookup(ext); ok {
		return f
	} else {
		var f = new(*program.Function)
		fragment.external.userFunMap =
			fragment.external.userFunMap.Inserted(ext, f)
		return f
	}
}
func (draft fragmentDraft) referGenFun(ns string, key genFunKey) **program.Function {
	var fragment = draft.content
	var ext = external[genFunKey] {
		namespace:   ns,
		internalKey: key,
	}
	if f, ok := fragment.external.genFunMap.Lookup(ext); ok {
		return f
	} else {
		var f = new(*program.Function)
		fragment.external.genFunMap =
			fragment.external.genFunMap.Inserted(ext, f)
		return f
	}
}
func compileImpl(impl *Impl, hdr *Header, ctx *NsHeaderMap, errs *source.Errors) *Fragment {
	if impl.namespace != hdr.namespace { panic("invalid arguments") }
	var ns = impl.namespace
	var fd = makeFragmentDraft(ns)
	fillDispatchInfo(fd, hdr, ctx, errs)
	fillFunctions(fd, impl, hdr, ctx, errs)
	var fragment = fd.content
	return fragment
}
func fillFunctions(fd fragmentDraft, impl *Impl, hdr *Header, ctx *NsHeaderMap, errs *source.Errors) {
	var ns = fd.namespace()
	hdr.userFunMap.ForEach(func(name string, group ctn.Map[string,*funHeader]) {
	group.ForEach(func(assoc string, f_hdr *funHeader) {
		var key = userFunKey { name: name, assoc: assoc }
		var f_impl, exists = impl.userFunMap.Lookup(key)
		if !(exists) { panic("missing function implementation") }
		var f = fd.createOrGetUserFun(key)
		f.SetName(key.Describe(ns))
		compileFunction(f, f_hdr, f_impl, fd, ctx, errs)
	})})
	hdr.genFunMap.ForEach(func(key genFunKey, f_hdr *funHeader) {
		var f_impl, exists = impl.genFunMap.Lookup(key)
		if !(exists) { panic("missing function implementation") }
		var f = fd.createOrGetGenFun(key)
		f.SetName(key.Describe(ns))
		compileFunction(f, f_hdr, f_impl, fd, ctx, errs)
	})
}

func compileFunction (
	f     *program.Function,
	hdr   *funHeader,
	impl  funImpl,
	fd    fragmentDraft,
	ctx   *NsHeaderMap,
	errs  *source.Errors,
) {
	switch I := impl.(type) {
	case funImplAstExpr:
		var v = compileFunImplAstExpr(I, hdr, fd, ctx, errs)
		f.SetExprBasedValue(v)
	case funImplLibraryNative:
		f.SetNativeValueById(I.Id, (hdr.funKind == FK_Const))
	// case funImplLoadedAsset:
	//     panic("not implemented")
	default:
		panic("impossible branch")
	}
}
func compileFunImplAstExpr (
	impl  funImplAstExpr,
	hdr   *funHeader,
	fd    fragmentDraft,
	ctx   *NsHeaderMap,
	errs  *source.Errors,
) *program.ExprBasedFunctionValue {
	var in_exp = hdr.inputsExplicit
	var in_imp = hdr.inputsImplicit
	var out = hdr.output
	var params = hdr.typeParams
	var ec = createExprContext(params, fd, ctx)
	var in = [...] [] *program.Binding { nil, nil }
	for i, fields := range ([...] *typsys.Fields { in_exp, in_imp }) {
		in[i] = make([] *program.Binding, len(fields.FieldList))
		for j, field := range fields.FieldList {
			var name = field.Name
			var t = typsys.CertainType { Type: field.Type }
			var loc = field.Info.Location
			in[i][j] = ec.createBinding(name, t, loc)
		}
	}
	var in_args, in_ctx = in[0], in[1]
	var expr = impl.Expr
	var expected = out.type_
	var cc = createExprCheckContext(expected, ec)
	var output, err = checkExpr(expr, cc)
	if err != nil {
		source.ErrorsJoin(errs, err)
		return nil
	}
	if err := ec.unusedBindingError(); err != nil {
		source.ErrorsJoin(errs, err)
		return nil
	}
	if hdr.funKind == FK_Const {
		output = program.Cached(output)
	}
	return &program.ExprBasedFunctionValue {
		InputArgs: in_args,
		InputCtx:  in_ctx,
		Output:    output,
	}
}


