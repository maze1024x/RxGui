package debugging

import (
	"fmt"
	"sync"
	"rxgui/standalone/util"
	"rxgui/standalone/util/richtext"
	"rxgui/lang/source"
	"rxgui/lang/typsys"
	"rxgui/lang/textual/ast"
	"rxgui/lang/textual/syntax"
	"rxgui/lang/textual/parser"
	"rxgui/lang/textual/transformer"
	"rxgui/interpreter/core"
	"rxgui/interpreter/program"
	"rxgui/interpreter/compiler"
)


type Repl struct {
	mutex       sync.Mutex
	debugInfo   *compiler.DebugInfo
	evaluator   *program.EvalContext
	bindings    [] *program.Binding
	eventloop   *core.EventLoop
}
func CreateRepl(info *compiler.DebugInfo, h core.RuntimeHandle) *Repl {
	var evaluator = program.CreateEvalContext(h)
	var bindings = make([] *program.Binding, 0)
	var eventloop = h.EventLoop()
	return &Repl {
		debugInfo:  info,
		evaluator:  evaluator,
		bindings:   bindings,
		eventloop:  eventloop,
	}
}
const ReplHelp = "<pre><b>*** REPL ***</b>\n" +
	"<b>Usage:</b> EXPR | NAME = EXPR | :run EXPR</pre>"

func (ctx *Repl) HandleCommand(cmd string, id int, d core.ReplInterface, exit util.ExitNotifier) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	d.ReplEcho(cmd, id)
	var code = source.StringToCode(cmd)
	var path = fmt.Sprintf("%d", id)
	var root = syntax.ReplRoot()
	var cmd_cst, err = parser.Parse(code, root, path)
	if err != nil {
		d.ReplRespondExprError(err.Message(), id)
		return
	}
	var key = source.FileKey {
		Context: "repl",
		Path:    path,
	}
	var name = fmt.Sprintf("temp%d", id)
	var cmd_ast = transformer.Transform(cmd_cst, key).(ast.VariousReplCmd)
	var expr_node = ast.ReplCmdGetExpr(cmd_ast.ReplCmd)
	{ var expr, err = compiler.CompileExpr(expr_node, "", nil, ctx.bindings, ctx.debugInfo)
	if err != nil {
		d.ReplRespondExprError(err.Message(), id)
		return
	}
	var value = ctx.evaluator.Eval(expr)
	ctx.bindExprValue(name, expr, value)
	d.ReplRespondExprValue(value, expr.Type, id)
	switch C := cmd_ast.ReplCmd.(type) {
	case ast.ReplAssign:
		var specified_name = ast.Id2String(C.Name)
		ctx.bindExprValue(specified_name, expr, value)
	case ast.ReplRun:
		if o, ok := (*value).(core.Observable); ok {
			var inner_t_, _ = program.T_Observable_(expr.Type.Type)
			var inner_t = typsys.CertainType { Type: inner_t_}
			go ctx.subscribe(o, inner_t, id, d, exit)
		} else {
			var msg richtext.Block
			msg.WriteLine("run: expect Observable", richtext.TAG_ERR)
			d.ReplRespondExprError(msg, id)
		}
	default:
		// do nothing
	}}
}
func (ctx *Repl) Expose(name string, v core.Object, t typsys.CertainType, loc source.Location) func() {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	b := ctx.bind(name, t, loc, v)
	return func() {
		ctx.unbind(b)
	}
}

func (ctx *Repl) bindExprValue(name string, expr *program.Expr, value core.Object) {
	ctx.bind(name, expr.Type, expr.Info.Location, value)
}
func (ctx *Repl) bind(name string, t typsys.CertainType, loc source.Location, value core.Object) *program.Binding {
	var binding = &program.Binding {
		Name:     name,
		Type:     t,
		Location: loc,
	}
	ctx.bindings = append(ctx.bindings, binding)
	ctx.evaluator = ctx.evaluator.NewCtxBind(binding, value)
	return binding
}
func (ctx *Repl) unbind(target *program.Binding) bool {
	var L = len(ctx.bindings)
	var index, found = -1, false
	for i := (L-1); i >= 0; i -= 1 {
		if ctx.bindings[i] == target {
			index, found = i, true
			ctx.evaluator = ctx.evaluator.NewCtxUnbind(ctx.bindings[i])
		}
	}
	if found {
		var new_bindings = make([] *program.Binding, 0, (L-1))
		for i, b := range ctx.bindings {
			if i != index {
				new_bindings = append(new_bindings, b)
			}
		}
		ctx.bindings = new_bindings
		return true
	} else {
		return false
	}
}
func (ctx *Repl) subscribe(o core.Observable, inner_t typsys.CertainType, id int, d core.ReplInterface, exit util.ExitNotifier) {
	var ch_values = make(chan core.Object, 256)
	var ch_error = make(chan error, 1)
	core.Run(o, ctx.eventloop, core.DataSubscriber {
		Values: ch_values,
		Error:  ch_error,
	})
	var ch_exit = exit.Signal()
	loop: for {
		select {
		case v, not_close := <- ch_values:
			if not_close {
				d.ReplNotifyObValue(v, inner_t, id)
			} else {
				d.ReplNotifyObCompletion(id)
				break loop
			}
		case err, not_close := <- ch_error:
			if not_close {
				d.ReplNotifyObError(err, id)
			} else {
				break loop
			}
		case <- ch_exit:
			break loop
		}
	}
}


