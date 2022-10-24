package core

import (
	"rxgui/standalone/util"
	"rxgui/standalone/util/richtext"
	"rxgui/lang/typsys"
)

type Debugger interface {
	CreateInstance(h RuntimeHandle) DebuggerInstance
}
type DebuggerInstance interface {
	ExitNotifier() util.ExitNotifier
	NotifyProgramCrash(msg richtext.Block) bool
	ReplInterface
	Inspector
	InputOutputActionMonitor
}
type ReplInterface interface {
	ReplEcho(cmd string, id int)
	ReplRespondExprValue(v Object, t typsys.CertainType, id int)
	ReplRespondExprError(msg richtext.Block, id int)
	ReplNotifyObValue(v Object, t typsys.CertainType, id int)
	ReplNotifyObError(err error, id int)
	ReplNotifyObCompletion(id int)
	ReplExposeValue(name string, v Object, t typsys.CertainType, info *FrameInfo) func()
}
type Inspector interface {
	InspectValue(v Object, t typsys.CertainType, info *FrameInfo, hint string)
	InspectError(err error, info *FrameInfo)
	InspectObSub(sub string, info *FrameInfo, hint string)
	InspectObValue(v Object, t typsys.CertainType, sub string, info *FrameInfo, hint string)
	InspectObError(err error, sub string, info *FrameInfo, hint string)
	InspectObCompletion(sub string, info *FrameInfo, hint string)
	InspectObContextDisposal(sub string, info *FrameInfo, hint string)
}
type InputOutputActionMonitor interface {
	LogRequest(method string, endpoint string, token string, body string, info *FrameInfo)
	LogFileRead(file string, info *FrameInfo)
	LogFileWrite(file string, content string, info *FrameInfo)
}


