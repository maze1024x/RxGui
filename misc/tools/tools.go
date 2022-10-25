package tools

import (
	"io"
	"os"
	"fmt"
	"os/exec"
	"path/filepath"
	"encoding/json"
	"rxgui/util"
)


func CrashReport(program_path string, html string) {
	var tool_path, err = getToolPath("crash_report")
	if err != nil { logError(err); return }
	var cmd = exec.Command(tool_path, program_path)
	stdin, err := cmd.StdinPipe()
	if err != nil { logError(err); return }
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil { logError(err); return }
	err = cmd.Process.Release()
	if err != nil { logError(err); return }
	_, err = fmt.Fprintln(stdin, html)
	if err != nil { logError(err); return }
	err = stdin.Close()
	if err != nil { logError(err); return }
}

func NaiveDebugger(program_path string) *NaiveDebuggerProcess {
	var tool_path, err = getToolPath("naive_debugger")
	if err != nil { panic(err) }
	var cmd = exec.Command(tool_path, program_path)
	stdin, err := cmd.StdinPipe()
	if err != nil { panic(err) }
	stdout, err := cmd.StdoutPipe()
	if err != nil { panic(err) }
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil { panic(err) }
	err = cmd.Process.Release()
	if err != nil { panic(err) }
	return &NaiveDebuggerProcess { stdin, stdout }
}
type NaiveDebuggerProcess struct {
	stdin   io.Writer
	stdout  io.Reader
}
func (p *NaiveDebuggerProcess) ReceiveCommand(exit chan(struct{})) (string, bool) {
	select {
	case <- exit:
		return "", false
	default:
		var line, _, err = util.WellBehavedFscanln(p.stdout)
		if err == nil {
			return line, true
		} else if err == io.EOF {
			close(exit)
			return "", false
		} else {
			logError(err)
			return "", false
		}
	}
}
type naiveDebuggerMessage struct {
	Category  string   `json:"category"`
	Content   string   `json:"content"`
}
func (p *NaiveDebuggerProcess) SendMessage(category string, content string, exit chan(struct{})) {
	select {
	case <- exit:
		return
	default:
		var msg = naiveDebuggerMessage { category, content }
		var binary, err = json.Marshal(msg)
		if err != nil { logError(err); return }
		_, err = p.stdin.Write(binary)
		if err != nil { logError(err); return }
		_, err = p.stdin.Write([] byte { byte('\n') })
		if err != nil { logError(err); return }
	}
}

func getToolPath(tool_name string) (string, error) {
	var exe_path, err = os.Executable()
	if err != nil { return "", err }
	var exe_dir = filepath.Dir(exe_path)
	var tool_full_name = (tool_name + filepath.Ext(exe_path))
	var tool_path = filepath.Join(exe_dir, tool_full_name)
	return tool_path, nil
}
func logError(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err.Error())
}


