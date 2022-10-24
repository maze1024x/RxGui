package debugging

import (
	"fmt"
	"time"
	"strconv"
	"rxgui/standalone/util"
	"rxgui/standalone/util/richtext"
	"rxgui/lang/typsys"
	"rxgui/interpreter/core"
	"rxgui/interpreter/compiler"
	"rxgui/support/tools"
)


func CreateNaiveDebugger(info *compiler.DebugInfo) core.Debugger {
	return naiveDebugger { info }
}

type naiveDebugger struct { info *compiler.DebugInfo }
func (d naiveDebugger) CreateInstance(h core.RuntimeHandle) core.DebuggerInstance {
	var ictx = MakeInspectContext(d.info.Context)
	var repl = CreateRepl(d.info, h)
	var instance = &naiveDebuggerInstance {
		proc: tools.NaiveDebugger(h.ProgramPath()),
		exit: make(chan struct{}),
		next: 1,
		ictx: ictx,
		repl: repl,
	}
	instance.proc.SendMessage("repl", ReplHelp, instance.exit)
	go instance.receiveCommands()
	return instance
}
type naiveDebuggerInstance struct {
	proc  *tools.NaiveDebuggerProcess
	exit  chan struct{}
	next  int
	ictx  InspectContext
	repl  *Repl
}
func (d *naiveDebuggerInstance) getCmdId() int {
	var id = d.next
	d.next += 1
	return id
}
func (d *naiveDebuggerInstance) receiveCommands() {
	loop: for {
		if line, ok := d.proc.ReceiveCommand(d.exit); ok {
			if line != "" {
				var id = d.getCmdId()
				d.repl.HandleCommand(line, id, d, d.ExitNotifier())
			}
		} else {
			break loop
		}
	}
}
func (d *naiveDebuggerInstance) sendEcho(cmd string, id int) {
	var b richtext.Block
	b.WriteSpan(fmt.Sprintf("[%d] ", id), richtext.TAG_INPUT)
	b.WriteSpan(cmd)
	d.proc.SendMessage("repl", b.RenderHtml(), d.exit)
}
func (d *naiveDebuggerInstance) sendReply(msg richtext.Block, id int, header string, tags ...string) {
	var b richtext.Block
	if header != "" { header = (" " + header) }
	b.WriteLine(fmt.Sprintf("(%d)%s", id, header), tags...)
	b.Append(msg)
	d.proc.SendMessage("repl", b.RenderHtml(), d.exit)
}
func (d *naiveDebuggerInstance) sendNotice(title string, msg richtext.Block) {
	var b richtext.Block
	b.WriteLine(fmt.Sprintf("[notice] %s", title), richtext.TAG_B)
	b.Append(msg)
	d.proc.SendMessage("repl", b.RenderHtml(), d.exit)
}
func (d *naiveDebuggerInstance) ExitNotifier() util.ExitNotifier {
	return util.MakeExitNotifier(d.exit)
}
func (d *naiveDebuggerInstance) NotifyProgramCrash(msg richtext.Block) bool {
	d.proc.SendMessage("*", msg.RenderHtml(), d.exit)
	d.proc.SendMessage("<control>", "PROGRAM_CRASH", d.exit)
	return true
}
func (d *naiveDebuggerInstance) ReplEcho(cmd string, id int) {
	d.sendEcho(cmd, id)
}
func (d *naiveDebuggerInstance) ReplRespondExprValue(v core.Object, t typsys.CertainType, id int) {
	d.sendReply(Inspect(v, t.Type, d.ictx), id, "", richtext.TAG_OUTPUT)
}
func (d *naiveDebuggerInstance) ReplRespondExprError(msg richtext.Block, id int) {
	d.sendReply(msg, id, "", richtext.TAG_FAILURE)
}
func (d *naiveDebuggerInstance) ReplNotifyObValue(v core.Object, t typsys.CertainType, id int) {
	d.sendReply(Inspect(v, t.Type, d.ictx), id, "<value>", richtext.TAG_OUTPUT)
}
func (d *naiveDebuggerInstance) ReplNotifyObError(err error, id int) {
	var msg richtext.Block
	msg.WriteLine(err.Error(), richtext.TAG_ERR)
	d.sendReply(msg, id, "<error>", richtext.TAG_FAILURE)
}
func (d *naiveDebuggerInstance) ReplNotifyObCompletion(id int) {
	var msg richtext.Block
	msg.WriteLine("|")
	d.sendReply(msg, id, "<complete>", richtext.TAG_SUCCESS)
}
func (d *naiveDebuggerInstance) ReplExposeValue(name string, v core.Object, t typsys.CertainType, info *core.FrameInfo) func() {
	var loc = info.CallSite
	var t_desc = typsys.Describe(t.Type)
	unset := d.repl.Expose(name, v, t, loc)
	var tip richtext.Block
	tip.WriteSpan(name, richtext.TAG_DBG_FIELD)
	tip.WriteSpan(t_desc, richtext.TAG_DBG_TYPE)
	var content = loc.FormatMessageLite(tip)
	d.sendNotice("exposed", content)
	return func() {
		unset()
		var tip richtext.Block
		tip.WriteSpan((name + " " + t_desc), richtext.TAG_DEL)
		var content = loc.FormatMessageLite(tip)
		d.sendNotice("unset (exposed value)", content)
	}
}
func (d *naiveDebuggerInstance) InspectValue(v core.Object, t typsys.CertainType, info *core.FrameInfo, hint string) {
	var f = info.DescribeCaller()
	var loc = info.CallSite
	var header = fmt.Sprintf("--- %s in %s", strconv.Quote(hint), f)
	var content = Inspect(v, t.Type, d.ictx)
	var msg richtext.Block
	msg.WriteLine(header, richtext.TAG_B)
	msg.Append(loc.FormatMessageLite(content))
	d.proc.SendMessage("inspection", msg.RenderHtml(), d.exit)
}
func (d *naiveDebuggerInstance) InspectError(err error, info *core.FrameInfo) {
	var f = info.DescribeCaller()
	var loc = info.CallSite
	var header = fmt.Sprintf("*** Error in %s", f)
	var content richtext.Block
	content.WriteLine(err.Error(), richtext.TAG_ERR)
	var msg richtext.Block
	msg.WriteLine(header, richtext.TAG_B)
	msg.Append(loc.FormatMessage(content))
	d.proc.SendMessage("inspection", msg.RenderHtml(), d.exit)
}
func (d *naiveDebuggerInstance) InspectObSub(sub string, info *core.FrameInfo, hint string) {
	var f = info.DescribeCaller()
	var loc = info.CallSite
	var header1 = fmt.Sprintf("--- %s subscribed in %s", strconv.Quote(hint), f)
	var header2 = fmt.Sprintf("    %s", sub)
	var content richtext.Block
	var msg richtext.Block
	msg.WriteLine(header1, richtext.TAG_B)
	msg.WriteLine(header2)
	msg.Append(loc.FormatMessageLite(content))
	d.proc.SendMessage("inspection", msg.RenderHtml(), d.exit)
}
func (d *naiveDebuggerInstance) InspectObValue(v core.Object, t typsys.CertainType, sub string, _ *core.FrameInfo, hint string) {
	var header1 = fmt.Sprintf(">>> %s <value>", strconv.Quote(hint))
	var header2 = fmt.Sprintf("    %s", sub)
	var content = Inspect(v, t.Type, d.ictx)
	var msg richtext.Block
	msg.WriteLine(header1, richtext.TAG_B, richtext.TAG_OUTPUT)
	msg.WriteLine(header2, richtext.TAG_OUTPUT)
	msg.Append(content)
	d.proc.SendMessage("inspection", msg.RenderHtml(), d.exit)
}
func (d *naiveDebuggerInstance) InspectObError(err error, sub string, _ *core.FrameInfo, hint string) {
	var header1 = fmt.Sprintf(">>> %s <error>", strconv.Quote(hint))
	var header2 = fmt.Sprintf("    %s", sub)
	var content richtext.Block
	content.WriteLine(err.Error(), richtext.TAG_ERR)
	var msg richtext.Block
	msg.WriteLine(header1, richtext.TAG_B, richtext.TAG_FAILURE)
	msg.WriteLine(header2, richtext.TAG_FAILURE)
	msg.Append(content)
	d.proc.SendMessage("inspection", msg.RenderHtml(), d.exit)
}
func (d *naiveDebuggerInstance) InspectObCompletion(sub string, _ *core.FrameInfo, hint string) {
	var header1 = fmt.Sprintf(">>> %s <complete>", strconv.Quote(hint))
	var header2 = fmt.Sprintf("    %s", sub)
	var content richtext.Block
	content.WriteLine("|")
	var msg richtext.Block
	msg.WriteLine(header1, richtext.TAG_B, richtext.TAG_SUCCESS)
	msg.WriteLine(header2, richtext.TAG_SUCCESS)
	msg.Append(content)
	d.proc.SendMessage("inspection", msg.RenderHtml(), d.exit)
}
func (d *naiveDebuggerInstance) InspectObContextDisposal(sub string, _ *core.FrameInfo, hint string) {
	var header1 = fmt.Sprintf(">>> %s <context-disposal>", strconv.Quote(hint))
	var header2 = fmt.Sprintf("    %s", sub)
	var content richtext.Block
	var msg richtext.Block
	msg.WriteLine(header1, richtext.TAG_B, richtext.TAG_I)
	msg.WriteLine(header2, richtext.TAG_I)
	msg.Append(content)
	d.proc.SendMessage("inspection", msg.RenderHtml(), d.exit)
}
func (d *naiveDebuggerInstance) LogRequest(method string, endpoint string, _ string, _ string, _ *core.FrameInfo) {
	var t = time.Now().Format("15:04:05")
	var msg richtext.Block
	msg.WriteLine(fmt.Sprintf("[%s] %s %s", t, method, endpoint))
	d.proc.SendMessage("io", msg.RenderHtml(), d.exit)
}
func (d *naiveDebuggerInstance) LogFileRead(file string, _ *core.FrameInfo) {
	var t = time.Now().Format("15:04:05")
	var msg richtext.Block
	msg.WriteLine(fmt.Sprintf("[%s] READ %s", t, file))
	d.proc.SendMessage("io", msg.RenderHtml(), d.exit)
}
func (d *naiveDebuggerInstance) LogFileWrite(file string, _ string, _ *core.FrameInfo) {
	var t = time.Now().Format("15:04:05")
	var msg richtext.Block
	msg.WriteLine(fmt.Sprintf("[%s] WRITE %s", t, file))
	d.proc.SendMessage("io", msg.RenderHtml(), d.exit)
}


