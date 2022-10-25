package main

/*
#include <string.h>
typedef const char* string_t;
*/
import "C"
import (
	"os"
	"bytes"
	"unsafe"
	"reflect"
	"runtime"
	"os/exec"
	"archive/zip"
	"path/filepath"
	"rxgui/qt"
	"rxgui/util/fatal"
	"rxgui/interpreter"
	"rxgui/interpreter/compiler"
)


var coproc *os.Process

//export StartCoproc
func StartCoproc(interpreter_cmd_ C.string_t, file_rel_path C.string_t, argc C.int, argv *C.string_t) {
	var interpreter_cmd = getString(interpreter_cmd_)
	var file = getAbsPath(getString(file_rel_path))
	var args = getStrings(argc, argv)
	if interpreter_cmd != "" {
		args = append([] string { file }, args...)
		file = interpreter_cmd
	}
	var a_out, a_in, err = os.Pipe()
	if err != nil { fatal.ThrowError(err) }
	{ var b_out, b_in, err = os.Pipe()
	if err != nil { fatal.ThrowError(err) }
	os.Stdin = a_out
	os.Stdout = b_in
	var cmd = exec.Command(file, args...)
	cmd.Stdin = b_out
	cmd.Stdout = a_in
	cmd.Stderr = os.Stderr
	{ var err = cmd.Start()
	if err != nil { fatal.ThrowError(err) }
	coproc = cmd.Process }}
}

//export CompileAndRun
func CompileAndRun(archive_addr *C.char, archive_size C.size_t, file_name C.string_t, entry_ns C.string_t) {
	runtime.LockOSThread()
	var archive = getBytes(archive_addr, archive_size)
	var r, err = zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
	if err != nil { fatal.ThrowError(err) }
	var fs = compiler.ZipFilesystem { Id: "default", Reader: r }
	var file = getString(file_name)
	{ var p, _, err = interpreter.Compile(file, fs)
	if err != nil { fatal.ThrowError(err) }
	var args = make([] string, len(os.Args))
	copy(args, os.Args)
	if len(args) > 0 {
		args = args[1:]
	}
	var p_ns = getString(entry_ns)
	var p_args = args
	qt.Init()
	{ var err = interpreter.Run(p, p_ns, p_args, nil, func() { qt.Exit(0) })
	if err != nil { fatal.ThrowError(err) }
	qt.Main() }}
}

func getBytes(ptr *C.char, size C.size_t) ([] byte) {
	var h = &reflect.SliceHeader {
		Data: uintptr(unsafe.Pointer(ptr)),
		Len:  int(size),
		Cap:  int(size),
	}
	return *(*([] byte))(unsafe.Pointer(h))
}
func getString(raw C.string_t) string {
	var h = &reflect.StringHeader {
		Data: uintptr(unsafe.Pointer(raw)),
		Len:  int(C.strlen(raw)),
	}
	return *(*string)(unsafe.Pointer(h))
}
func getStrings(n C.int, ptr *C.string_t) ([] string) {
	var h = &reflect.SliceHeader {
		Data: uintptr(unsafe.Pointer(ptr)),
		Len:  int(n),
		Cap:  int(n),
	}
	var t = *(*([] C.string_t))(unsafe.Pointer(h))
	var all = make([] string, len(t))
	for i := range t {
		all[i] = getString(t[i])
	}
	return all
}
func getAbsPath(rel_path string) string {
	var exe_file, err = os.Executable()
	if err != nil { fatal.ThrowError(err) }
	var exe_dir = filepath.Dir(exe_file)
	return filepath.Join(exe_dir, rel_path)
}

// dummy, needed by Go compiler
func main() {}


