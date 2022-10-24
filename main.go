package main

import (
    "fmt"
    "os"
    "runtime"
    "rxgui/standalone/qt"
    "rxgui/standalone/util"
    "rxgui/standalone/util/argv"
    "rxgui/standalone/util/fatal"
    "rxgui/lang/textual"
    "rxgui/lang/textual/syntax"
    "rxgui/interpreter"
    "rxgui/interpreter/core"
    "rxgui/interpreter/debugging"
    "rxgui/support/atom"
)


const Version = "0.0.0 experimental"
const SourceFilePathPrompt = "Input a source file path or press Enter to start REPL:"

type Args struct {
    Positional  [] string  `arg:"positional" hint:"[PATH [ARGUMENT]...]"`
    Command     string     `arg:"command" key:"help-0; version-0; atom-0; parse; run" default:"run" desc:"show this help; show version; start plugin backend service for the Atom Editor; parse the file at PATH; run the file or directory at PATH"`
    EntryNS     string     `arg:"value-string" key:"entry" hint:"NS" desc:"namespace of entry point"`
    Debug       bool       `arg:"flag-enable" key:"debug" desc:"enable debugger"`
}

var Commands = map[string] func(*Args) {
    "version": func(_ *Args) {
        fmt.Println(Version)
    },
    "atom": func(_ *Args) {
        var err = atom.LangServer(os.Stdin, os.Stdout, os.Stderr)
        if err != nil { fatal.ThrowError(err) }
    },
    "parse": func(args *Args) {
        var L = len(args.Positional)
        if L == 0 {
            textual.DebugParser(os.Stdin, "(stdin)", syntax.ReplRoot())
        } else if L >= 1 {
            for _, file := range args.Positional {
                f, err := os.Open(file)
                if err != nil { fatal.ThrowError(err) }
                textual.DebugParser(f, f.Name(), syntax.DefaultRoot())
                err = f.Close()
                if err != nil { fatal.ThrowError(err) }
            }
        }
    },
    "run": func(args *Args) {
        var file, file_not_specified, p_args = (func() (string, bool, ([] string)) {
            var L = len(args.Positional)
            if L == 0 {
                var file = ""
                return file, true, nil
            } else {
                var file = args.Positional[0]
                return file, false, args.Positional[1:]
            }
        })()
        if file_not_specified {
            fmt.Fprintf(os.Stderr, "%s\n", SourceFilePathPrompt)
            var line, _, err = util.WellBehavedFscanln(os.Stdin)
            if err != nil { fatal.ThrowError(err) }
            if line != "" {
                file = line
                file_not_specified = false
            }
        }
        var p, info, err = interpreter.Compile(file, nil)
        if err != nil { fatal.ThrowError(err) }
        var p_ns = args.EntryNS
        var d = (func() core.Debugger {
            var debugger_enabled = (args.Debug || file_not_specified)
            if debugger_enabled {
                return debugging.CreateNaiveDebugger(info)
            } else {
                return nil
            }
        })()
        qt.Init()
        { var err = interpreter.Run(p, p_ns, p_args, d, func() { qt.Exit(0) })
        if err != nil { fatal.ThrowError(err) }
        qt.Main() }
    },
}

func main() {
    runtime.LockOSThread()
    var args, help, err = argv.ParseArgs[Args](os.Args)
    if err != nil { fatal.ThrowBadArgsError(err, help) }
    if args.Command == "help" {
        fmt.Println(help)
    } else {
        Commands[args.Command](&args)
    }
}


