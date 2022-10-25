package interpreter

import (
    "embed"
    "fmt"
    "errors"
    "strconv"
    "rxgui/qt"
    "rxgui/util/richtext"
    "rxgui/interpreter/lang/source"
    "rxgui/interpreter/lang/textual/ast"
    "rxgui/interpreter/builtin"
    "rxgui/interpreter/compiler"
    "rxgui/interpreter/core"
    "rxgui/interpreter/program"
)


func Compile(file string, fs compiler.FileSystem) (*program.Program, *compiler.DebugInfo, error) {
    var p, info, errs1, errs2 = compile(file, fs)
    if errs1 != nil { return nil, nil, errs1 }
    if errs2 != nil { return nil, nil, errs2 }
    return p, info, nil
}
func compile(file string, fs compiler.FileSystem) (*program.Program, *compiler.DebugInfo, richtext.Errors, source.Errors) {
    if file == "" {
        file = dummyProjectFile
        fs = dummyProjectFileSystem
    }
    if fs == nil {
        fs = compiler.RealFileSystem {}
    }
    var q = [] string {
        file,
    }
    var visited = make(map[string] struct{})
    var lists = make([] [] string, 0)
    for len(q) > 0 {
        var current = q[0]
        q = q[1:]
        if _, ok := visited[current]; ok {
            continue
        }
        visited[current] = struct{}{}
        var manifest, err = compiler.ReadManifest(current, fs)
        if err != nil { return nil, nil, richtext.Std2Errors(err), nil }
        lists = append(lists, manifest.ProjectFiles)
        for _, dep := range manifest.DependencyFiles {
            q = append(q, dep)
        }
    }
    var groups = [] compiler.SourceFileGroup {
        builtinSourceFileGroup,
    }
    for _, list := range lists {
        groups = append(groups, compiler.SourceFileGroup {
            FileList:   list,
            FileSystem: fs,
        })
    }
    var meta = program.Metadata {
        ProgramPath: file,
    }
    return compiler.Compile(groups, meta)
}
//go:embed builtin/builtin.km
var builtinRawFileSystem embed.FS
var builtinSourceFileGroup = compiler.SourceFileGroup {
    FileList: [] string {
        "builtin/builtin.km",
    },
    FileSystem: compiler.EmbeddedFileSystem {
        Id: "builtin",
        FS: builtinRawFileSystem,
    },
}
var dummyProjectFile = "dummy.manifest.json"
var dummyProjectFileSystem = compiler.FileSystem(compiler.InlineFileSystem {
    Id:    "dummy",
    Files: map[string] ([] byte) {
        "dummy.km":       ([] byte)("namespace :: entry { $() }"),
        dummyProjectFile: ([] byte)(`{"ProjectFiles":["dummy.km"]}`),
    },
})

func Lint(file string, fs compiler.FileSystem) (richtext.Errors, source.Errors) {
    var _, _, errs1, errs2 = compile(file, fs)
    return errs1, errs2
}
func ParseBuiltinFiles() ([] *ast.Root, richtext.Error) {
    var group = builtinSourceFileGroup
    var result = make([] *ast.Root, len(group.FileList))
    for i, file := range group.FileList {
        var root, err = compiler.Load(file, group.FileSystem)
        if err != nil { return nil, err }
        result[i] = root
    }
    return result, nil
}

func Run(p *program.Program, ns string, args ([] string), d core.Debugger, k func()) error {
    if entry, ok := p.Executable.LookupEntry(ns); ok {
        go (func() {
            run(p, entry, args, d)
            qt.Schedule(k)
        })()
        return nil
    } else {
        return errors.New(fmt.Sprintf(
            `missing entry point for namespace "%s ::"`, ns,
        ))
    }
}
func run(p *program.Program, entry **program.Function, args ([] string), d core.Debugger) {
    var eventloop = core.CreateEventLoop()
    var c = &runtimeContext {
        program:   p,
        arguments: args,
        eventloop: eventloop,
        debugger:  nil,
    }
    var h = &runtimeHandle {
        context:   c,
        frameInfo: core.TopLevelFrameInfo(),
    }
    if d != nil {
        c.debugger = d.CreateInstance(h)
    }
    var ctx = program.CreateEvalContext(h)
    var obj = (program.CallFunction { Callee: entry }).Eval(ctx)
    var o = core.GetObservable(obj)
    var E = make(chan error, 1)
    var T = make(chan bool)
    core.Run(o, eventloop, core.DataSubscriber {
        Error:     E,
        Terminate: T,
    })
    if <- T {
        if c.debugger != nil {
            c.debugger.ExitNotifier().Wait()
        }
    } else {
        var e = <- E
        core.Crash(h, core.ErrorInEntryObservable, e.Error())
    }
}
type runtimeContext struct {
    program    *program.Program
    arguments  [] string
    eventloop  *core.EventLoop
    debugger   core.DebuggerInstance
}
type runtimeHandle struct {
    context    *runtimeContext
    frameInfo  *core.FrameInfo
}
func (h *runtimeHandle) GetFrameInfo() *core.FrameInfo {
    return h.frameInfo
}
func (h *runtimeHandle) WithFrameInfo(info *core.FrameInfo) core.RuntimeHandle {
    return &runtimeHandle {
        context:   h.context,
        frameInfo: info,
    }
}
func (h *runtimeHandle) LibraryNativeFunction(id string) core.NativeFunction {
    var f, exists = builtin.LookupFunction(id)
    if !(exists) {
        var msg = fmt.Sprintf("no such native function: %s", strconv.Quote(id))
        core.Crash(h, core.MissingNativeFunction, msg)
    }
    return f
}
func (h *runtimeHandle) ProgramPath() string {
    return h.context.program.Metadata.ProgramPath
}
func (h *runtimeHandle) ProgramArgs() ([] string) {
    return h.context.arguments
}
func (h *runtimeHandle) SerializationContext() core.SerializationContext {
    return h.context.program.TypeInfo
}
func (h *runtimeHandle) EventLoop() *core.EventLoop {
    return h.context.eventloop
}
func (h *runtimeHandle) Debugger() (core.DebuggerInstance, bool) {
    var d = h.context.debugger
    return d, (d != nil)
}


