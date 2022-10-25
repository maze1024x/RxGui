package core

import (
	"os"
	"fmt"
	"rxgui/util/richtext"
	"rxgui/misc/tools"
)


type CrashKind string
const (
	ProgrammedCrash CrashKind = "crash"
	ValueUndefined CrashKind = "undefined"
	AssertionFailed CrashKind = "assertion failed"
	InvariantViolation CrashKind = "invariant violation"
	InvalidArgument CrashKind = "invalid argument"
	MissingNativeFunction CrashKind = "missing native function"
	ErrorInEntryObservable CrashKind = "error occurred in entry observable"
	UseAfterFree CrashKind = "use after free of external resources"
	FailedToGenerateRandomNumber CrashKind = "failed to generate random number"
)

func Crash(h RuntimeHandle, kind CrashKind, msg string) {
	print("[Crash] ", kind, ": ", msg, "\n")
	var info = h.GetFrameInfo()
	var f = info.DescribeCaller()
	var loc = info.CallSite
	var b richtext.Block
	b.WriteSpan("Crashed in ", richtext.TAG_ERR)
	b.WriteSpan(f, richtext.TAG_ERR_INLINE)
	b.WriteLineFeed()
	b.Append(loc.FormatMessage(describeCrash(kind, msg)))
	b.Append(info.DescribeStackTrace().WithLeadingLine("***"))
	if d, ok := h.Debugger(); ok && d.NotifyProgramCrash(b) {
		// do nothing
	} else {
		var program_path = h.ProgramPath()
		var html = b.RenderHtml()
		tools.CrashReport(program_path, html)
	}
	os.Exit(2)
	// noinspection ALL
	panic("process should have exited")
}
func describeCrash(kind CrashKind, msg string) richtext.Block {
	var b richtext.Block
	b.WriteSpan(fmt.Sprintf("%s: ", kind), richtext.TAG_ERR)
	b.WriteSpan(msg, richtext.TAG_ERR_NOTE)
	return b
}

func Crash1[A any] (h RuntimeHandle, kind CrashKind, msg string) A {
	var a A
	Crash(h, kind, msg)
	return a
}
func Crash2[A any, B any] (h RuntimeHandle, kind CrashKind, msg string) (A,B) {
	var a A
	var b B
	Crash(h, kind, msg)
	return a, b
}


