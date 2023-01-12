package core

import (
    "fmt"
    "strconv"
    "rxgui/util/ctn"
    "rxgui/interpreter/lang/typsys"
)


type Logger struct {
    FrameInfo  *FrameInfo
    Debugger   ctn.Maybe[DebuggerInstance]
}
func MakeLogger(h RuntimeHandle) Logger {
    return Logger {
        FrameInfo: h.GetFrameInfo(),
        Debugger:  ctn.MakeMaybe(h.Debugger()),
    }
}
func ErrorLogger(h RuntimeHandle) func(error) {
    return MakeLogger(h).LogError
}

func (lg Logger) LogError(err error) {
    print("[Error] ", err.Error(), "\n")
    if d, ok := lg.Debugger.Value(); ok {
        d.InspectError(err, lg.FrameInfo)
    }
}
func (lg Logger) Inspect(v Object, t typsys.CertainType, hint string) {
    if d, ok := lg.Debugger.Value(); ok {
        d.InspectValue(v, t, lg.FrameInfo, hint)
    }
}
func (lg Logger) Expose(name string, v Object, t typsys.CertainType) Observable {
    if d, ok := lg.Debugger.Value(); ok {
        return Observable(func(pub DataPublisher) {
            var ctx, ob = pub.useInheritedContext()
            var unset = d.ReplExposeValue(name, v, t, lg.FrameInfo)
            ctx.registerCleaner(unset)
            ob.value(nil)
            ob.complete()
        })
    } else {
        return doSync(func() {})
    }
}
var logProxyNextId = uint64(1000)
func (lg Logger) Trace(o Observable, inner_t typsys.CertainType, hint string) Observable {
    if d, ok := lg.Debugger.Value(); ok {
        return Observable(func(pub DataPublisher) {
            var ctx, ob = pub.useInheritedContext()
            var id = logProxyNextId
            logProxyNextId++
            var sub = fmt.Sprintf("%d (%p)", id, ctx)
            var info = lg.FrameInfo
            d.InspectObSub(sub, info, hint)
            var proxy = &observer {
                value: func(v Object) {
                    d.InspectObValue(v, inner_t, sub, info, hint)
                    ob.value(v)
                },
                error: func(err error) {
                    d.InspectObError(err, sub, info, hint)
                    ob.error(err)
                },
                complete: func() {
                    d.InspectObCompletion(sub, info, hint)
                    ob.complete()
                },
            }
            ctx.registerCleaner(func() {
                d.InspectObContextDisposal(sub, info, hint)
            })
            pub.run(o, ctx, proxy)
        })
    } else {
        return o
    }
}
func (lg Logger) Watch(o Observable, inner_t typsys.CertainType, hint string) Observable {
    return lg.Trace(o, inner_t, hint).MapTo(nil).TakeLast()
}

func (lg Logger) LogRequest(req *Request) {
    if d, ok := lg.Debugger.Value(); ok {
        var method = string(req.Method)
        var endpoint = req.Endpoint.String()
        var token = strconv.Quote(req.AuthToken)
        var body = strconv.Quote((string)(req.BodyContent))
        d.LogRequest(method, endpoint, token, body, lg.FrameInfo)
    }
}
func (lg Logger) LogFileRead(f File) {
    if d, ok := lg.Debugger.Value(); ok {
        var file = f.Path
        d.LogFileRead(file, lg.FrameInfo)
    }
}
func (lg Logger) LogFileWrite(f File, content_ ([] byte)) {
    if d, ok := lg.Debugger.Value(); ok {
        var file = f.Path
        var content = strconv.Quote((string)(content_))
        d.LogFileWrite(file, content, lg.FrameInfo)
    }
}


