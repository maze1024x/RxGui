package compiler

import (
	"rxgui/util/richtext"
	"rxgui/interpreter/lang/source"
	"rxgui/interpreter/lang/typsys"
	"rxgui/interpreter/lang/textual/ast"
	"rxgui/interpreter/program"
)


type SourceFileGroup struct {
	FileList    [] string
	FileSystem  FileSystem
}
type DebugInfo struct {
	Context     *NsHeaderMap
	Executable  *Executable
}

func Compile(groups ([] SourceFileGroup), meta program.Metadata) (*program.Program, *DebugInfo, richtext.Errors, source.Errors) {
	var ldr_errs richtext.Errors
	var chk_errs source.Errors
	var hdr_list = make([] *Header, 0)
	var impl_list = make([] *Impl, 0)
	for _, group := range groups {
		for _, file := range group.FileList {
			var root, err = Load(file, group.FileSystem)
			if err == nil {
				var hdr, impl = analyze(root, &chk_errs)
				hdr_list = append(hdr_list, hdr)
				impl_list = append(impl_list, impl)
			} else {
				richtext.ErrorsJoin(&ldr_errs, err)
			}
		}
	}
	if ldr_errs != nil {
		return nil, nil, ldr_errs, chk_errs
	}
	var ctx = groupHeaders(hdr_list, &chk_errs)
	var fragments = make([] *Fragment, 0)
	for _, hdr := range hdr_list {
		checkHeader(hdr, ctx, &chk_errs)
	}
	for i := range hdr_list {
		var hdr = hdr_list[i]
		var impl = impl_list[i]
		var fragment = compileImpl(impl, hdr, ctx, &chk_errs)
		fragments = append(fragments, fragment)
	}
	var rtti = ctx.generateTypeInfo()
	var exe = link(fragments)
	var p = &program.Program {
		Metadata:   meta,
		TypeInfo:   rtti,
		Executable: exe,
	}
	var info = &DebugInfo {
		Context:    ctx,
		Executable: exe,
	}
	return p, info, ldr_errs, chk_errs
}

func CompileExpr(node ast.Expr, ns string, expected typsys.Type, bindings ([] *program.Binding), info *DebugInfo) (*program.Expr, *source.Error) {
	var fd = makeFragmentDraft(ns)
	var ec = createExprContext(nil, fd, info.Context)
	for _, b := range bindings {
		ec.addBinding(b)
		ec.useBinding(b.Name)
	}
	var cc = createExprCheckContext(expected, ec)
	var expr, err = checkExpr(node, cc)
	if err != nil { return nil, err }
	if err := ec.unusedBindingError(); err != nil { return nil, err }
	var fragment = fd.content
	linkFragment(info.Executable, fragment)
	return expr, nil
}


