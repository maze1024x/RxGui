package atom

import (
	"io"
	"fmt"
	"bufio"
	"strings"
	"encoding/json"
	"rxgui/standalone/util"
	"rxgui/lang/textual/ast"
	"rxgui/interpreter"
)

type LangServerContext struct {
	WriteDebugLog   func(info string)
	BuiltinAstData  [] *ast.Root
}
func LangServer(input io.Reader, output io.Writer, debug io.Writer) error {
	input = bufio.NewReader(input)
	var write_line = func(line ([] byte)) error {
		_, err := output.Write(line)
		if err != nil { return err }
		_, err = output.Write(([]byte)("\n"))
		if err != nil { return err }
		return nil
	}
	var builtinAstRootNodes, err = interpreter.ParseBuiltinFiles()
	if err != nil {
		_, _ = fmt.Fprintln(debug, err.Error())
		return err
	}
	var ctx = LangServerContext {
		WriteDebugLog:  func(info string) { _, _ = fmt.Fprintln(debug, info) },
		BuiltinAstData: builtinAstRootNodes,
	}
	for {
		var line_runes, _, err = util.WellBehavedFscanln(input)
		if err != nil { return err }
		var line = string(line_runes)
		var i = strings.Index(line, " ")
		if i == -1 { i = len(line) }
		var cmd = line[:i]
		var arg = strings.TrimPrefix(line[i:], " ")
		switch cmd {
		case "quit":
			return nil
		case "ping":
			err = write_line(([] byte)(arg))
			if err != nil { return err }
		case "lint":
			var raw_req = ([] byte)(arg)
			var req LintRequest
			err := json.Unmarshal(raw_req, &req)
			if err != nil { return err }
			var res = Lint(req, ctx)
			raw_res, err := json.Marshal(&res)
			if err != nil { return err }
			err = write_line(raw_res)
			if err != nil { return err }
		case "autocomplete":
			var raw_req = ([] byte)(arg)
			var req AutoCompleteRequest
			err := json.Unmarshal(raw_req, &req)
			if err != nil { return err }
			var res = AutoComplete(req, ctx)
			raw_res, err := json.Marshal(&res)
			if err != nil { return err }
			err = write_line(raw_res)
			if err != nil { return err }
		}
	}
}


