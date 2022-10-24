package util

import (
	"os"
	"path/filepath"
)

var exeDir = (func() string {
	var exe_path, err = os.Executable()
	if err != nil { panic(err) }
	return filepath.Dir(exe_path)
})()
func ExeDir() string { return exeDir }


