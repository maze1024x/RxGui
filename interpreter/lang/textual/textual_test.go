package textual

import (
	"testing"
	"os"
	"path/filepath"
	"rxgui/lang/textual/syntax"
)


// NOTE: output only (result NOT checked)
func TestOutputSampleParsingResult(t *testing.T) {
	// parseSampleFileDefault(t, "1")
	// parseSampleFileDefault(t, "2")
	// parseSampleFileDefault(t, "3")
	parseSampleFileDefault(t, "4")
	parseSampleFileDefault(t, "5")
}

func parseSampleFile(t *testing.T, name string, root syntax.Id) {
	var exe_path, err = os.Executable()
	if err != nil { panic(err) }
	var project_path = filepath.Dir(filepath.Dir(exe_path))
	var sample_path = filepath.Join(
		project_path,
		"lang", "textual", "test", (name + ".txt"),
	)
	{ var f, err = os.Open(sample_path)
	if err != nil { panic(err) }
	var success = DebugParser(f, sample_path, root)
	if !(success) {
		t.Fail()
	} }
}
func parseSampleFileDefault(t *testing.T, name string) {
	parseSampleFile(t, name, syntax.DefaultRoot())
}


