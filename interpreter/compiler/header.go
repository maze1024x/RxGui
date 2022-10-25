package compiler

import (
    "strings"
    "rxgui/util/ctn"
    "rxgui/interpreter/lang/source"
    "rxgui/interpreter/lang/typsys"
    "rxgui/interpreter/lang/textual/ast"
    "rxgui/interpreter/program"
)


type NsHeaderMap struct { Map ctn.Map[string,*Header] }
func groupHeaders(list ([] *Header), errs *source.Errors) *NsHeaderMap {
    var m0 = ctn.MakeMutMap[string,*([]*Header)](ctn.StringCompare)
    var m1 = ctn.MakeMutMap[string,*Header](ctn.StringCompare)
    for _, hdr := range list {
        var ns = hdr.namespace
        var existing, exists = m0.Lookup(ns)
        if !(exists) {
            var ptr = new([] *Header)
            m0.Insert(ns, ptr)
            existing = ptr
        }
        *existing = append(*existing, hdr)
    }
    m0.ForEach(func(ns string, ptr *([]*Header)) {
        m1.Insert(ns, mergeNsHeaders(ns, *ptr, errs))
    })
    return &NsHeaderMap { Map: m1.Map() }
}
func mergeNsHeaders(ns string, list ([] *Header), errs *source.Errors) *Header {
    var all_als_ns = ctn.MakeMutMap[string,*aliasTarget[string]](ctn.StringCompare)
    var all_als_ref = ctn.MakeMutMap[string,*aliasTarget[source.Ref]](ctn.StringCompare)
    var all_types = ctn.MakeMutMap[string,*typsys.TypeDef](ctn.StringCompare)
    var all_user = ctn.MakeMutMap[string,ctn.Map[string,*funHeader]](ctn.StringCompare)
    var all_gen = ctn.MakeMutMap[genFunKey,*funHeader](genFunKeyCompare)
    for _, hdr := range list {
        if hdr.namespace != ns {
            panic("invalid arguments")
        }
        hdr.nsAliasMap.ForEach(func(name string, item *aliasTarget[string]) {
            var _, duplicate = all_als_ns.Lookup(name)
            if duplicate {
                var loc = item.Location
                source.ErrorsJoin(errs, source.MakeError(loc,
                    E_DuplicateAlias { Name: name }))
            } else {
                all_als_ns.Insert(name, item)
            }
        })
        hdr.refAliasMap.ForEach(func(name string, item *aliasTarget[source.Ref]) {
            var _, duplicate = all_als_ref.Lookup(name)
            if duplicate {
                var loc = item.Location
                source.ErrorsJoin(errs, source.MakeError(loc,
                    E_DuplicateAlias { Name: name }))
            } else {
                all_als_ref.Insert(name, item)
            }
        })
        hdr.typeMap.ForEach(func(name string, item *typsys.TypeDef) {
            var _, duplicate = all_types.Lookup(name)
            if duplicate {
                var loc = item.Info.Location
                source.ErrorsJoin(errs, source.MakeError(loc,
                    E_DuplicateTypeDecl { Name: name }))
            } else {
                all_types.Insert(name, item)
            }
        })
        hdr.userFunMap.ForEach(func(name string, items ctn.Map[string,*funHeader]) {
            var group = (func() ctn.Map[string, *funHeader] {
                var group, exists = all_user.Lookup(name)
                if exists {
                    return group
                } else {
                    return ctn.MakeMap[string,*funHeader](ctn.StringCompare)
                }
            })()
            items.ForEach(func(assoc string, item *funHeader) {
                var duplicate = group.Has(assoc)
                if duplicate {
                    var loc = item.funInfo.Location
                    source.ErrorsJoin(errs, source.MakeError(loc,
                        E_DuplicateFunDecl { Name: name, Assoc: assoc }))
                } else {
                    group = group.Inserted(assoc, item)
                }
            })
            all_user.Insert(name, group)
        })
        hdr.genFunMap.ForEach(func(key genFunKey, item *funHeader) {
            all_gen.Insert(key, item)
        })
    }
    return &Header {
        namespace:   ns,
        nsAliasMap:  all_als_ns.Map(),
        refAliasMap: all_als_ref.Map(),
        typeMap:     all_types.Map(),
        userFunMap:  all_user.Map(),
        genFunMap:   all_gen.Map(),
    }
}
func (ctx *NsHeaderMap) FindType(ref source.Ref) (*typsys.TypeDef, bool) {
    return ctx.lookupType(ref)
}
func (ctx *NsHeaderMap) generateTypeInfo() program.TypeInfo {
    var reg = make(map[source.Ref] *typsys.TypeDef)
    ctx.Map.ForEach(func(ns string, hdr *Header) {
        hdr.typeMap.ForEach(func(name string, def *typsys.TypeDef) {
            var ref = source.MakeRef(ns, name)
            reg[ref] = def
        })
    })
    return program.TypeInfo {
        TypeRegistry: reg,
    }
}
func (ctx* NsHeaderMap) tryRedirectRef(ns string, ref *source.Ref) bool {
    if hdr, ok := ctx.Map.Lookup(ns); ok {
        if target, ok := hdr.nsAliasMap.Lookup(ref.Namespace); ok {
        if target.Valid {
            *ref = source.MakeRef(target.Target, ref.ItemName)
            return true
        }}
        if ref.Namespace == "" {
            if target, ok := hdr.refAliasMap.Lookup(ref.ItemName); ok {
            if target.Valid {
                *ref = target.Target
                return true
            }}
        }
    }
    if ns != "" {
        if ref.Namespace == "" {
            *ref = source.MakeRef(ns, ref.ItemName)
            return true
        }
    }
    return false
}
func (ctx *NsHeaderMap) lookupType(ref source.Ref) (*typsys.TypeDef, bool) {
    if hdr, ok := ctx.Map.Lookup(ref.Namespace); ok {
    if def, ok := hdr.typeMap.Lookup(ref.ItemName); ok {
        return def, true
    }}
    return nil, false
}
func (ctx *NsHeaderMap) lookupMethod(recv_ref source.Ref, name string) (*funHeader, bool) {
    if hdr, ok := ctx.Map.Lookup(recv_ref.Namespace); ok {
    if group, ok := hdr.userFunMap.Lookup(name); ok {
    if f, ok := group.Lookup(recv_ref.ItemName); ok {
        if f.funKind == FK_Method {
            return f, true
        }
    }}}
    return nil, false
}
func (ctx *NsHeaderMap) lookupExactFunction(ref source.Ref, assoc string, kind FunKind) (*funHeader, string, userFunKey, bool) {
    if hdr, ok := ctx.Map.Lookup(ref.Namespace); ok {
    if group, ok := hdr.userFunMap.Lookup(ref.ItemName); ok {
    if f, ok := group.Lookup(assoc); ok {
        if f.funKind == kind {
            var ns = ref.Namespace
            var key = userFunKey {
                name:  ref.ItemName,
                assoc: assoc,
            }
            return f, ns, key, true
        }
    }}}
    return nil, "", userFunKey{}, false
}
func (ctx *NsHeaderMap) lookupFunction(ref source.Ref, r0 ctn.Maybe[source.Ref], kind FunKind) (*funHeader, string, userFunKey, bool) {
    if r0, ok := r0.Value(); ok {
        var ns, item = ref.Namespace, ref.ItemName
        var NS, A = r0.Namespace, r0.ItemName
        if ((NS != "") && (ns == "")) {
            if f, ns, key, ok := ctx.lookupExactFunction(ref, "", kind); ok {
                return f, ns, key, true
            }
            return ctx.lookupExactFunction(source.MakeRef(NS, item), A, kind)
        }
        if (NS == ns) {
            if f, ns, key, ok := ctx.lookupExactFunction(ref, A, kind); ok {
                return f, ns, key, true
            }
            return ctx.lookupExactFunction(ref, "", kind)
        }
    }
    return ctx.lookupExactFunction(ref, "", kind)
}
func (ctx *NsHeaderMap) lookupAssignableKindFunctions(ref source.Ref) ([] func()(*funHeader,string,userFunKey)) {
    var ns = ref.Namespace
    var name = ref.ItemName
    var list = make([] func()(*funHeader,string,userFunKey), 0)
    if hdr, ok := ctx.Map.Lookup(ns); ok {
    if group, ok := hdr.userFunMap.Lookup(name); ok {
        group.ForEach(func(assoc string, f *funHeader) {
            if isAssignableKindFunction(f) {
                list = append(list, func() (*funHeader, string, userFunKey) {
                    var key = userFunKey { name: name, assoc: assoc }
                    return f, ns, key
                })
            }
        })
    }}
    return list
}

type Header struct {
    namespace    string
    nsAliasMap   ctn.Map[string, *aliasTarget[string]]
    refAliasMap  ctn.Map[string, *aliasTarget[source.Ref]]
    typeMap      ctn.Map[string, *typsys.TypeDef]
    userFunMap   ctn.Map[string, ctn.Map[string, *funHeader]]
    genFunMap    ctn.Map[genFunKey, *funHeader]
}
type aliasTarget[T any] struct {
    Valid     bool
    Target    T
    Location  source.Location
}
type funHeader struct {
    funKind         FunKind
    funInfo         FunInfo
    variadic        bool
    typeParams      [] string
    output          funOutput
    inputsExplicit  *typsys.Fields
    inputsImplicit  *typsys.Fields
}
type funOutput struct {
    type_  typsys.Type
    loc    source.Location
}
type FunKind int
const (
    FK_Generated FunKind = iota
    FK_Ordinary
    FK_Operator
    FK_Method
    FK_Const
    FK_Entry
)
type FunInfo struct {
    Location  source.Location
    Document  string
}
func isAssignableKindFunction(f *funHeader) bool {
    switch f.funKind {
    case FK_Ordinary, FK_Operator, FK_Const:
        return true
    }
    return false
}
func someFunctionKindAssignable(group ctn.Map[string,*funHeader]) bool {
    var result = false
    group.ForEach(func(_ string, f *funHeader) {
        result = (result || isAssignableKindFunction(f))
    })
    return result
}
type genFunKey struct {
    typeName   string
    funName    string
    fieldName  string
}
func (key genFunKey) Describe(ns string) string {
    var elements = [] string { ns, key.typeName, key.funName, key.fieldName }
    var buf strings.Builder
    buf.WriteRune('[')
    buf.WriteString(strings.Join(elements, ","))
    buf.WriteRune(']')
    return buf.String()
}
func genFunKeyCompare(a genFunKey, b genFunKey) ctn.Ordering {
    var o = ctn.StringCompare(a.typeName, b.typeName)
    if o != ctn.Equal {
        return o
    } else {
        var o = ctn.StringCompare(a.funName, b.funName)
        if o != ctn.Equal {
            return o
        } else {
            return ctn.StringCompare(a.fieldName, b.fieldName)
        }
    }
}

func analyze(root *ast.Root, errs *source.Errors) (*Header,*Impl) {
    var ns = ast.MaybeId2String(root.Namespace)
    var als_ns = ctn.MakeMutMap[string,*aliasTarget[string]](ctn.StringCompare)
    var als_ref = ctn.MakeMutMap[string,*aliasTarget[source.Ref]](ctn.StringCompare)
    var types = ctn.MakeMutMap[string,*typsys.TypeDef](ctn.StringCompare)
    var types_gen = ctn.MakeMutMap[string,generatedConstList](ctn.StringCompare)
    var user_hdr = ctn.MakeMutMap[string,ctn.Map[string,*funHeader]](ctn.StringCompare)
    var user_impl = ctn.MakeMutMap[userFunKey,funImpl](userFunKeyCompare)
    var user_gen = ctn.MakeMutMap[userFunKey,generatedConstList](userFunKeyCompare)
    var gen_hdr = ctn.MakeMutMap[genFunKey,*funHeader](genFunKeyCompare)
    var gen_impl = ctn.MakeMutMap[genFunKey,funImpl](genFunKeyCompare)
    for _, alias := range root.Aliases {
        if alias.Off {
            continue
        }
        var loc = alias.Location
        switch T := alias.Target.AliasTarget.(type) {
        case ast.AliasToNamespace:
            var target = ast.Id2String(T.Namespace)
            var name = ast.MaybeId2String(alias.Name)
            var duplicate = als_ns.Has(name)
            if duplicate {
                source.ErrorsJoin(errs, source.MakeError(loc,
                    E_DuplicateAlias {
                        Name: name,
                    }))
                continue
            }
            als_ns.Insert(name, &aliasTarget[string] {
                Valid:    true,
                Target:   target,
                Location: loc,
            })
        case ast.AliasToRefBase:
            var target = getRef(T.RefBase)
            var name = (func() string {
                if name := ast.MaybeId2String(alias.Name); name != "" {
                    return name
                } else {
                    return target.ItemName
                }
            })()
            var duplicate = als_ref.Has(name)
            if duplicate {
                source.ErrorsJoin(errs, source.MakeError(loc,
                    E_DuplicateAlias {
                        Name: name,
                    }))
                continue
            }
            als_ref.Insert(name, &aliasTarget[source.Ref] {
                Valid:    true,
                Target:   target,
                Location: loc,
            })
        default:
            panic("impossible branch")
        }
    }
    var pm = makeParamMapping(root.Statements)
    for _, stmt := range root.Statements {
        switch S := stmt.Statement.(type) {
        case ast.DeclType:
            if S.Off {
                continue
            }
            var name = ast.Id2String(S.Name)
            if types.Has(name) {
                source.ErrorsJoin(errs, source.MakeError(S.Location,
                    E_DuplicateTypeDecl {
                        Name: name,
                    }))
                continue
            }
            var loc = S.Location
            var doc = ast.GetDocContent(S.Docs)
            var info = typsys.Info { Location: loc, Document: doc }
            var ref = source.MakeRef(ns, name)
            var ifs = ctn.MapEach(S.Implements, getRef)
            var params = ctn.MapEach(S.TypeParams, ast.Id2String)
            var ctx = createTypeConsContext(params)
            var gen = make(generatedConstList, 0)
            var content = makeTypeDefContent(S.TypeDef, ctx, &gen, errs)
            types_gen.Insert(name, gen)
            types.Insert(name, &typsys.TypeDef {
                Info:       info,
                Ref:        ref,
                Interfaces: ifs,
                Parameters: params,
                Content:    content,
            })
        default:
            var u = unifyFunDecl(S, pm, ns)
            if u.off {
                continue
            }
            var name = ast.Id2String(u.name)
            var group = (func() ctn.Map[string, *funHeader] {
                var group, exists = user_hdr.Lookup(name)
                if exists {
                    return group
                } else {
                    return ctn.MakeMap[string,*funHeader](ctn.StringCompare)
                }
            })()
            var params = ctn.MapEach(u.sig.TypeParams, ast.Id2String)
            var ctx = createTypeConsContext(params)
            var kind = u.kind
            var assoc = getAssocTypeName(u.sig.Inputs, kind, ns, ctx)
            var duplicate = group.Has(assoc)
            if duplicate {
                source.ErrorsJoin(errs, source.MakeError(u.loc,
                    E_DuplicateFunDecl {
                        Name:  name,
                        Assoc: assoc,
                    }))
                continue
            }
            var key = userFunKey { name: name, assoc: assoc }
            var impl = u.impl
            var loc = u.loc
            var doc = ast.GetDocContent(u.doc)
            var info = FunInfo { Location: loc, Document: doc }
            var va = u.va
            var gen = make(generatedConstList, 0)
            var in_exp_ = u.sig.Inputs
            var in_imp_ = u.sig.Implicit
            var in_exp = makeFunInputFields(in_exp_, nil, ctx, &gen, errs)
            var in_imp = makeFunInputFields(in_imp_, in_exp, ctx, &gen, errs)
            var out_t = makeType(u.sig.Output, ctx)
            var out = funOutput { type_: out_t, loc: u.sig.Output.Location }
            user_gen.Insert(key, gen)
            user_hdr.Insert(name, group.Inserted(assoc, &funHeader {
                funKind:        kind,
                funInfo:        info,
                variadic:       va,
                typeParams:     params,
                output:         out,
                inputsExplicit: in_exp,
                inputsImplicit: in_imp,
            }))
            user_impl.Insert(key, impl)
        }
    }
    types_gen.ForEach(func(src_type_name string, list generatedConstList) {
        var params = (func() ([] string) {
            var src_type_def, ok = types.Lookup(src_type_name)
            if !(ok) { panic("something went wrong") }
            return src_type_def.Parameters
        })()
        list.consume(params, func(item_name string, item_hdr *funHeader, item_impl funImpl) {
            var key = genFunKey {
                typeName:  src_type_name,
                fieldName: item_name,
            }
            if gen_hdr.Has(key) { panic("something went wrong") }
            gen_hdr.Insert(key, item_hdr)
            gen_impl.Insert(key, item_impl)
        })
    })
    user_gen.ForEach(func(fun_key userFunKey, list generatedConstList) {
        var params = (func() ([] string) {
            var group = (func() ctn.Map[string, *funHeader] {
                var group, ok = user_hdr.Lookup(fun_key.name)
                if !(ok) { panic("something went wrong") }
                return group
            })()
            var f_hdr, ok = group.Lookup(fun_key.assoc)
            if !(ok) { panic("something went wrong") }
            return f_hdr.typeParams
        })()
        list.consume(params, func(item_name string, item_hdr *funHeader, item_impl funImpl) {
            var key = genFunKey {
                funName:   fun_key.name,
                typeName:  fun_key.assoc,
                fieldName: item_name,
            }
            if gen_hdr.Has(key) { panic("something went wrong") }
            gen_hdr.Insert(key, item_hdr)
            gen_impl.Insert(key, item_impl)
        })
    })
    var hdr = &Header {
        namespace:   ns,
        nsAliasMap:  als_ns.Map(),
        refAliasMap: als_ref.Map(),
        typeMap:     types.Map(),
        userFunMap:  user_hdr.Map(),
        genFunMap:   gen_hdr.Map(),
    }
    var impl = &Impl {
        namespace:  ns,
        userFunMap: user_impl.Map(),
        genFunMap:  gen_impl.Map(),
    }
    return hdr, impl
}
type generatedConst struct {
    name  string
    val   funImpl
    typ   typsys.Type
    loc   source.Location
    doc   string
}
type generatedConstList ([] generatedConst)
func (list generatedConstList) consume(params ([] string), k func(string,*funHeader,funImpl)) {
    for _, item := range list {
        var info = FunInfo {
            Location: item.loc,
            Document: item.doc,
        }
        var out = funOutput {
            type_: item.typ,
            loc:   item.loc,
        }
        var item_hdr = &funHeader {
            funKind:    FK_Generated,
            funInfo:    info,
            typeParams: params,
            output:     out,
            inputsExplicit: &typsys.Fields {},
            inputsImplicit: &typsys.Fields {},
        }
        var item_impl = item.val
        k(item.name, item_hdr, item_impl)
    }
}

func checkHeader(hdr *Header, ctx *NsHeaderMap, errs *source.Errors) {
    validateAlterHeaderAliases(hdr, ctx, errs)
    validateAlterHeaderTypes(hdr, ctx, errs)
    validateHeaderFunctions(hdr, ctx, errs)
}
func validateAlterHeaderAliases(hdr *Header, ctx *NsHeaderMap, errs *source.Errors) {
    var hasRelevantRef = func(ref source.Ref) bool {
        if ref_hdr, ok := ctx.Map.Lookup(ref.Namespace); ok {
            if ref_hdr.typeMap.Has(ref.ItemName) {
                return true
            }
            if group, ok := ref_hdr.userFunMap.Lookup(ref.ItemName); ok {
            if someFunctionKindAssignable(group) {
                return true
            }}
        }
        return false
    }
    var setInvalidAndThrow = func(valid *bool, err *source.Error) {
        *valid = false
        source.ErrorsJoin(errs, err)
    }
    hdr.nsAliasMap.ForEach(func(name string, target *aliasTarget[string]) {
        var v, loc = &(target.Valid), target.Location
        if ctx.Map.Has(name) {
            setInvalidAndThrow(v, source.MakeError(loc,
                E_InvalidAlias {
                    Name: name,
                }))
        }
        if !(ctx.Map.Has(target.Target)) {
            setInvalidAndThrow(v, source.MakeError(loc,
                E_AliasTargetNotFound {
                    Target: ("namespace " + target.Target),
                }))
        }
    })
    hdr.refAliasMap.ForEach(func(name string, target *aliasTarget[source.Ref]) {
        var v, loc = &(target.Valid), target.Location
        if hasRelevantRef(source.MakeRef("", name)) {
            setInvalidAndThrow(v, source.MakeError(loc,
                E_InvalidAlias {
                    Name: name,
                }))
        }
        if !(hasRelevantRef(target.Target)) {
            setInvalidAndThrow(v, source.MakeError(loc,
                E_AliasTargetNotFound {
                    Target: target.Target.String(),
                }))
        }
    })
}
func validateAlterHeaderTypes(hdr *Header, ctx *NsHeaderMap, errs *source.Errors) {
    var ns = hdr.namespace
    var validateAlterSingleType = func(t *typsys.Type, loc source.Location) {
        var err = validateAlterType(t, loc, ns, ctx)
        source.ErrorsJoin(errs, err)
    }
    var validateAlterFieldTypes = func(fields *typsys.Fields) {
        for i := range fields.FieldList {
            var field = &(fields.FieldList[i])
            if field.Type != nil {
                var loc = field.Info.Location
                var field_t = &(field.Type)
                validateAlterSingleType(field_t, loc)
            }
        }
    }
    hdr.typeMap.ForEach(func(_ string, def *typsys.TypeDef) {
        for i := range def.Interfaces {
            var ref = &(def.Interfaces[i])
            var loc = def.Info.Location
            var err = validateAlterInterfaceTypeRef(ref, loc, ns, ctx)
            source.ErrorsJoin(errs, err)
        }
        var fields, ok = (func() (*typsys.Fields, bool) {
            switch content := def.Content.(type) {
                case typsys.Record:    return content.Fields, true
                case typsys.Union:     return content.Fields, true
                case typsys.Interface: return content.Fields, true
                default:               return nil, false
            }
        })()
        if ok {
            validateAlterFieldTypes(fields)
        }
    })
    hdr.userFunMap.ForEach(func(_ string, group ctn.Map[string,*funHeader]) {
    group.ForEach(func(_ string, f *funHeader) {
        validateAlterFieldTypes(f.inputsExplicit)
        validateAlterFieldTypes(f.inputsImplicit)
        validateAlterSingleType(&(f.output.type_), f.output.loc)
    })})
    hdr.genFunMap.ForEach(func(_ genFunKey, f *funHeader) {
        validateAlterFieldTypes(f.inputsExplicit)
        validateAlterFieldTypes(f.inputsImplicit)
        validateAlterSingleType(&(f.output.type_), f.output.loc)
    })
}
func validateAlterType(t *typsys.Type, loc source.Location, ns string, ctx *NsHeaderMap) *source.Error {
    var first_error *source.Error = nil
    var throw = func(err *source.Error) {
        if first_error == nil {
            first_error = err
        }
    }
    *t = typsys.Transform(*t, func(t typsys.Type) (typsys.Type, bool) {
        if ref_type, ok := t.(typsys.RefType); ok {
            var ref = ref_type.Def
            var def, exists = ctx.lookupType(ref)
            var ref_altered = false
            if !(exists) {
                if ctx.tryRedirectRef(ns, &ref) {
                    ref_altered = true
                    def, exists = ctx.lookupType(ref)
                }
            }
            if !(exists) {
                throw(source.MakeError(loc,
                    E_NoSuchType {
                        ref.String(),
                    }))
                return program.T_InvalidType(), true
            }
            var required = len(def.Parameters)
            var given = len(ref_type.Args)
            if given != required {
                throw(source.MakeError(loc,
                    E_TypeArgsWrongQuantity {
                        Type:     ref.String(),
                        Given:    given,
                        Required: required,
                    }))
                return program.T_InvalidType(), true
            }
            if ref_altered {
                return typsys.RefType {
                    Def:  ref,
                    Args: ref_type.Args,
                }, true
            }
        }
        return nil, false
    })
    return first_error
}
func validateAlterInterfaceTypeRef(ref *source.Ref, loc source.Location, ns string, ctx *NsHeaderMap) *source.Error {
    // NOTE: this function requires T_InvalidType to be an interface type
    var def, exists = ctx.lookupType(*ref)
    if !(exists) {
        if ctx.tryRedirectRef(ns, ref) {
            def, exists = ctx.lookupType(*ref)
        }
    }
    if !(exists) {
        var err = source.MakeError(loc, E_NoSuchType { (*ref).String() })
        *ref = program.T_InvalidType().(typsys.RefType).Def
        return err
    }
    var _, ok = def.Content.(typsys.Interface)
    if !(ok) {
        var err = source.MakeError(loc, E_NotInterface { (*ref).String() })
        *ref = program.T_InvalidType().(typsys.RefType).Def
        return err
    }
    return nil
}
func validateHeaderFunctions(hdr *Header, ctx *NsHeaderMap, errs *source.Errors) {
    hdr.userFunMap.ForEach(func(name string, group ctn.Map[string,*funHeader]) {
    group.ForEach(func(assoc string, f *funHeader) {
        if f.funKind == FK_Method {
            var ns = hdr.namespace
            var type_ref = source.MakeRef(ns, assoc)
            if def, ok := ctx.lookupType(type_ref); ok {
                var conflict = (func() bool {
                    if R, ok := def.Content.(typsys.Record); ok {
                        var _, conflict = R.FieldIndexMap[name]
                        return conflict
                    }
                    if I, ok := def.Content.(typsys.Interface); ok {
                        var _, conflict = I.FieldIndexMap[name]
                        return conflict
                    }
                    return false
                })()
                if conflict {
                    var err = source.MakeError(f.funInfo.Location,
                        E_MethodNameUnavailable {
                            Name: name,
                        })
                    source.ErrorsJoin(errs, err)
                }
            }
        }
        if f.funKind == FK_Operator {
            var in_exp = f.inputsExplicit.FieldList
            if len(in_exp) == 0 {
                var loc = f.funInfo.Location
                source.ErrorsJoin(errs, source.MakeError(loc,
                    E_MissingOperatorParameter {}))
            } else if in_exp[0].Info.HasDefaultValue {
                var loc = f.funInfo.Location
                source.ErrorsJoin(errs, source.MakeError(loc,
                    E_OperatorFirstParameterHasDefaultValue {}))
            }
        }
        if f.variadic {
            var in_exp = f.inputsExplicit.FieldList
            var arity = len(in_exp)
            if arity == 0 {
                var loc = f.funInfo.Location
                source.ErrorsJoin(errs, source.MakeError(loc,
                    E_MissingVariadicParameter {}))
            } else {
                var last_index = (arity - 1)
                var last = in_exp[last_index]
                if _, ok := program.T_List_(last.Type); ok {
                    // OK
                } else {
                    var loc = f.funInfo.Location
                    source.ErrorsJoin(errs, source.MakeError(loc,
                        E_InvalidVariadicParameter {}))
                }
            }
        }
        for _, imp_field := range f.inputsImplicit.FieldList {
            if P, _, ok := ast.SplitImplicitRef(imp_field.Name); ok {
                var ok = false
                for _, p := range f.typeParams {
                    if p == P {
                        ok = true
                        break
                    }
                }
                if !(ok) {
                    var loc = imp_field.Info.Location
                    source.ErrorsJoin(errs, source.MakeError(loc,
                        E_NoSuchTypeParameter { P }))
                }
            }
        }
    })})
}


type typeConsContext struct {
    parameters  map[string] struct{}
}
func createTypeConsContext(params ([] string)) *typeConsContext {
    return &typeConsContext {
        parameters: (func() (map[string] struct{}) {
            var set = make(map[string] struct{})
            for _, p := range params {
                set[p] = struct{}{}
            }
            return set
        })(),
    }
}
func (ctx *typeConsContext) hasTypeParameter() bool {
    return (len(ctx.parameters) > 0)
}
func (ctx *typeConsContext) isTypeParameter(name string) bool {
    var _, is_parameter = ctx.parameters[name]
    return is_parameter
}
func makeType(t ast.Type, ctx *typeConsContext) typsys.Type {
    var param_name, is_param = getTypeParamName(t.Ref.Base, ctx)
    if is_param {
        return typsys.ParameterType {
            Name: param_name,
        }
    } else {
        var def_ref = getRef(t.Ref.Base)
        var args = make([] typsys.Type, len(t.Ref.TypeArgs))
        for i := range t.Ref.TypeArgs {
            args[i] = makeType(t.Ref.TypeArgs[i], ctx)
        }
        return typsys.RefType {
            Def:  def_ref,
            Args: args,
        }
    }
}
func makeTypeDefContent (
    content  ast.VariousTypeDef,
    ctx      *typeConsContext,
    gen      *generatedConstList,
    errs     *source.Errors,
) typsys.TypeDefContent {
    type fieldNode struct {
        loc    source.Location
        doc    [] ast.Doc
        name   ast.Identifier
        type_  ctn.Maybe[ast.Type]
        dv     ast.MaybeExpr
    }
    var collect = func(kind string, nodes ([] fieldNode)) typsys.Fields {
        var fields = make([] typsys.Field, len(nodes))
        var index = make(map[string] int)
        for i, node := range nodes {
            var name = ast.Id2String(node.name)
            if !(isValidFieldName(name)) {
                source.ErrorsJoin(errs, source.MakeError(
                    node.loc,
                    E_InvalidFieldName {
                        FieldKind: kind,
                        FieldName: name,
                    },
                ))
                continue
            }
            if _, duplicate := index[name]; duplicate {
                source.ErrorsJoin(errs, source.MakeError(
                    node.loc,
                    E_DuplicateField {
                        FieldKind: kind,
                        FieldName: name,
                    },
                ))
                continue
            }
            var loc = node.loc
            var doc = ast.GetDocContent(node.doc)
            var type_ typsys.Type
            if type_node, ok := node.type_.Value(); ok {
                type_ = makeType(type_node, ctx)
            }
            var default_, has_default = node.dv.(ast.Expr)
            if has_default {
                *gen = append(*gen, generatedConst {
                    name: name,
                    val:  funImplAstExpr { default_ },
                    typ:  type_,
                    loc:  default_.Location,
                    doc:  "",
                })
            }
            var info = typsys.FieldInfo {
                Info: typsys.Info { Location: loc, Document: doc },
                HasDefaultValue: has_default,
            }
            fields[i] = typsys.Field {
                Info: info,
                Name: name,
                Type: type_,
            }
            index[name] = i
        }
        return typsys.Fields {
            FieldIndexMap: index,
            FieldList:     fields,
        }
    }
    switch C := content.TypeDef.(type) {
    case ast.Record:
        var nodes = ctn.MapEach(C.Def.Fields, func(raw ast.Field) fieldNode {
            return fieldNode {
                loc:   raw.Location,
                doc:   raw.Docs,
                name:  raw.Name,
                type_: ctn.Just(raw.Type),
                dv:    raw.Default,
            }
        })
        var fields = collect("field", nodes)
        return typsys.Record {
            Fields: &fields,
        }
    case ast.Interface:
        var nodes = ctn.MapEach(C.Methods, func(raw ast.Method) fieldNode {
            return fieldNode {
                loc:   raw.Location,
                doc:   raw.Docs,
                name:  raw.Name,
                type_: ctn.Just(raw.Type),
                dv:    nil,
            }
        })
        var methods = collect("method", nodes)
        return typsys.Interface {
            Fields: &methods,
        }
    case ast.Union:
        var nodes = make([] fieldNode, 0)
        for _, item := range C.Items {
            nodes = append(nodes, fieldNode {
                loc:   item.Location,
                doc:   nil,
                name:  item.Ref.Base.Item,
                type_: ctn.Just(item),
                dv:    nil,
            })
        }
        var items = collect("item", nodes)
        return typsys.Union {
            Fields: &items,
        }
    case ast.Enum:
        if ctx.hasTypeParameter() {
            source.ErrorsJoin(errs, source.MakeError(
                content.Location,
                E_GenericEnum {},
            ))
        }
        var nodes = make([] fieldNode, 0)
        for _, item := range C.Items {
            nodes = append(nodes, fieldNode {
                loc:   item.Location,
                doc:   item.Docs,
                name:  item.Name,
                type_: nil,
                dv:    nil,
            })
        }
        var items = collect("item", nodes)
        return typsys.Enum {
            Fields: &items,
        }
    case ast.NativeTypeDef:
        return typsys.NativeContent {}
    default:
        panic("impossible branch")
    }
}
func makeFunInputFields (
    inputs  ast.MaybeInputs,
    saved   *typsys.Fields,
    ctx     *typeConsContext,
    gen     *generatedConstList,
    errs    *source.Errors,
) *typsys.Fields {
    if inputs, ok := inputs.(ast.Inputs); ok {
        var nodes = inputs.Content.Fields
        var fields = make([] typsys.Field, len(nodes))
        var index = make(map[string] int)
        for i, node := range nodes {
            var name = ast.Id2String(node.Name)
            if !(isValidFieldName(name)) {
                source.ErrorsJoin(errs, source.MakeError(
                    node.Location,
                    E_InvalidFieldName {
                        FieldKind: "input",
                        FieldName: name,
                    },
                ))
                continue
            }
            var _, duplicate = index[name]
            if !(duplicate) {
                if saved != nil {
                    _, duplicate = saved.FieldIndexMap[name]
                }
            }
            if duplicate {
                source.ErrorsJoin(errs, source.MakeError(
                    node.Location,
                    E_DuplicateField {
                        FieldKind: "input",
                        FieldName: name,
                    },
                ))
                continue
            }
            var loc = node.Location
            var doc = ast.GetDocContent(node.Docs)
            var type_ = makeType(node.Type, ctx)
            var default_, has_default = node.Default.(ast.Expr)
            if has_default {
                *gen = append(*gen, generatedConst {
                    name: name,
                    val:  funImplAstExpr { default_ },
                    typ:  type_,
                    loc:  default_.Location,
                    doc:  "",
                })
            }
            var info = typsys.FieldInfo {
                Info: typsys.Info { Location: loc, Document: doc },
                HasDefaultValue: has_default,
            }
            fields[i] = typsys.Field {
                Info: info,
                Name: name,
                Type: type_,
            }
            index[name] = i
        }
        return &typsys.Fields {
            FieldIndexMap: index,
            FieldList:     fields,
        }
    } else {
        return &typsys.Fields {}
    }
}
func getTypeParamName(rb ast.RefBase, ctx *typeConsContext) (string, bool) {
    var _, ns_specified = rb.NS.(ast.Identifier)
    var ns_not_specified = !(ns_specified)
    var item_name = ast.Id2String(rb.Item)
    if ns_not_specified && ctx.isTypeParameter(item_name) {
        return item_name, true
    } else {
        return "", false
    }
}
func getRef(rb ast.RefBase) source.Ref {
    var ns = ast.MaybeId2String(rb.NS)
    var item = ast.Id2String(rb.Item)
    return source.MakeRef(ns, item)
}
func isValidFieldName(name string) bool {
    return (name != Underscore)
}

const This = "this"
type unifiedFunDecl struct {
    kind  FunKind
    va    bool
    name  ast.Identifier
    sig   ast.FunctionSignature
    loc   source.Location
    doc   [] ast.Doc
    off   bool
    impl  funImpl
}
type paramMapping (map[string] ([] ast.Identifier))
func makeParamMapping(stmts ([] ast.VariousStatement)) paramMapping {
    var pm = make(paramMapping)
    for _, stmt := range stmts {
        switch S := stmt.Statement.(type) {
        case ast.DeclType:
            var name = ast.Id2String(S.Name)
            var _, exists = pm[name]
            if !(exists) {
                pm[name] = S.TypeParams
            }
        }
    }
    return pm
}
func unifyFunDecl(stmt ast.Statement, pm paramMapping, ns string) *unifiedFunDecl {
    switch S := stmt.(type) {
    case ast.DeclFunction:
        var kind = (func() FunKind {
            if S.Operator {
                return FK_Operator
            } else {
                return FK_Ordinary
            }
        })()
        return &unifiedFunDecl {
            kind: kind,
            va:   S.Variadic,
            name: S.Name,
            sig:  S.Signature,
            loc:  S.Location,
            doc:  S.Docs,
            off:  S.Off,
            impl: funImplFromAstBody(S.Body),
        }
    case ast.DeclConst:
        var sig = ast.FunctionSignature {
            Node:   S.Type.Node,
            Output: S.Type,
        }
        return &unifiedFunDecl {
            kind: FK_Const,
            name: S.Name,
            sig:  sig,
            loc:  S.Location,
            doc:  S.Docs,
            off:  S.Off,
            impl: funImplFromAstBody(S.Body),
        }
    case ast.DeclMethod:
        var recv = S.Receiver
        var recv_type, recv_params = getReceiverTypeAndParams(recv, pm, ns)
        var recv_inputs = getReceiverInputs(recv_type)
        var sig = ast.FunctionSignature {
            Node:       S.Node,
            TypeParams: recv_params,
            Inputs:     recv_inputs,
            Output:     S.Type,
        }
        return &unifiedFunDecl {
            kind: FK_Method,
            name: S.Name,
            sig:  sig,
            loc:  S.Location,
            doc:  S.Docs,
            off:  S.Off,
            impl: funImplFromAstBody(S.Body),
        }
    case ast.DeclEntry:
        var empty = ast.String2Id("", S.Node)
        var expr = ast.WrapBlockAsExpr(S.Content)
        var ob_null_t = program.AT_Observable(program.AT_Null())(S.Node)
        var sig = ast.FunctionSignature {
            Node:   S.Node,
            Output: ob_null_t,
        }
        return &unifiedFunDecl {
            kind: FK_Entry,
            name: empty,
            sig:  sig,
            loc:  S.Location,
            doc:  S.Docs,
            off:  S.Off,
            impl: funImplAstExpr { expr },
        }
    default:
        panic("impossible branch")
    }
}
func getReceiverTypeAndParams(recv ast.Identifier, pm paramMapping, ns string) (ast.Type, ([] ast.Identifier)) {
    var recv_name = ast.Id2String(recv)
    var recv_params = pm[recv_name]
    var recv_type = ast.Type {
        Node: recv.Node,
        Ref:  ast.Ref {
            Node: recv.Node,
            Base: ast.RefBase {
                Node: recv.Node,
                NS:   ast.String2Id(ns, recv.Node),
                Item: recv,
            },
            TypeArgs: ctn.MapEach(recv_params, func(param ast.Identifier) ast.Type {
                return ast.Type {
                    Node: param.Node,
                    Ref:  ast.Ref {
                        Node: param.Node,
                        Base: ast.RefBase {
                            Node: param.Node,
                            Item: param,
                        },
                    },
                }
            }),
        },
    }
    return recv_type, recv_params
}
func getReceiverInputs(recv_type ast.Type) ast.Inputs {
    var content = ast.RecordDef {
        Node:   recv_type.Node,
        Fields: [] ast.Field { {
            Node: recv_type.Node,
            Name: ast.String2Id(This, recv_type.Node),
            Type: recv_type,
        } },
    }
    return ast.Inputs {
        Node:    recv_type.Node,
        Content: content,
    }
}
func getAssocTypeName(inputs ast.Inputs, kind FunKind, ns string, ctx *typeConsContext) string {
    switch kind {
    case FK_Operator, FK_Method:
        var fields = inputs.Content.Fields
        if len(fields) > 0 {
            var first = fields[0]
            var ref = getRef(first.Type.Ref.Base)
            if ref.Namespace == ns {
            if !((ref.Namespace == "") && ctx.isTypeParameter(ref.ItemName)) {
                return ref.ItemName
            }}
        }
    }
    return ""
}


