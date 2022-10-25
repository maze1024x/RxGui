package core

import (
	"strings"
	"rxgui/util/richtext"
	"rxgui/interpreter/lang/source"
)


type RuntimeHandle interface {
	GetFrameInfo() *FrameInfo
	WithFrameInfo(info *FrameInfo) RuntimeHandle
	LibraryNativeFunction(id string) NativeFunction
	ProgramPath() string
	ProgramArgs() ([] string)
	SerializationContext() SerializationContext
	EventLoop() *EventLoop
	Debugger() (DebuggerInstance, bool)
}

type FrameInfo struct {
	Base      *FrameInfo
	Callee    string
	CallSite  source.Location
}
var topLevelFrameInfo = &FrameInfo {
	Callee: "[top-level]",
}
func TopLevelFrameInfo() *FrameInfo {
	return topLevelFrameInfo
}
func AddFrameInfo(h RuntimeHandle, callee string, loc source.Location) RuntimeHandle {
	var base = h.GetFrameInfo()
	return h.WithFrameInfo(&FrameInfo {
		Base:     base,
		Callee:   callee,
		CallSite: loc,
	})
}
func (info *FrameInfo) DescribeCaller() string {
	if info.Base != nil {
		return describeFunctionName(info.Base.Callee)
	} else {
		return "[nil]"
	}
}
func (info *FrameInfo) DescribeStackTrace() richtext.Block {
	var b richtext.Block
	var loc = source.Location {}
	for i := info; i != topLevelFrameInfo; i = i.Base {
		var callee_desc = describeFunctionName(i.Callee)
		b.WriteSpan(callee_desc, richtext.TAG_B)
		if loc.File != nil {
			var file_desc = loc.FileDesc()
			var pos_desc = loc.PosDesc()
			b.WriteSpan(file_desc)
			b.WriteSpan(pos_desc)
		}
		b.WriteLineFeed()
		loc = i.CallSite
	}
	return b
}
func describeFunctionName(name string) string {
	if name == "" {
		return "[entry]"
	} else if strings.HasSuffix(name, "::") {
		var ns = strings.TrimSuffix(name, "::")
		return (ns + "::[entry]")
	} else {
		return name
	}
}


