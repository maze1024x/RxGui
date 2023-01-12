package compiler

import (
    "os"
    "io"
    "fmt"
    "embed"
    "io/fs"
    "errors"
    "archive/zip"
    "path/filepath"
)


type FileSystem interface {
    Open(path string) (File, error)
    Key() string
}
type File interface {
    Close() error
    ReadContent() ([] byte, error)
}
func ReadFile(path string, fs FileSystem) ([] byte, error) {
    var f, err = fs.Open(path)
    if err != nil { return nil, err }
    { var content, err = f.ReadContent()
    if err != nil { return nil, err }
    { var err = f.Close()
    if err != nil { return nil, err }
    return content, nil } }
}

type RealFileSystem struct {}
type RealFile struct {
    Fd  *os.File
}
func (_ RealFileSystem) Open(path string) (File, error) {
    var fd, err = os.Open(path)
    if err != nil { return nil, err }
    return RealFile{fd}, nil
}
func (_ RealFileSystem) Key() string {
    return ""
}
func (f RealFile) Close() error {
    return f.Fd.Close()
}
func (f RealFile) ReadContent() ([] byte, error) {
    return io.ReadAll(f.Fd)
}

type InlineFileSystem struct {
    Id     string
    Files  map[string] ([] byte)
}
type InlineFile struct {
    Content  [] byte
}
func (fs InlineFileSystem) Open(path string) (File, error) {
    var content, exists = fs.Files[path]
    if !(exists) {
        return nil, errors.New("no such file: " + path)
    }
    return InlineFile{content}, nil
}
func (fs InlineFileSystem) Key() string {
    return fmt.Sprintf("(inline.%s)", fs.Id)
}
func (f InlineFile) Close() error {
    return nil
}
func (f InlineFile) ReadContent() ([] byte, error) {
    return f.Content, nil
}

type EmbeddedFileSystem struct {
    Id  string
    FS  embed.FS
}
type StdFile struct {
    File  fs.File
}
func (fs EmbeddedFileSystem) Open(path string) (File, error) {
    var f, err = fs.FS.Open(filepath.ToSlash(path))
    if err != nil { return nil, err }
    return StdFile{f}, nil
}
func (fs EmbeddedFileSystem) Key() string {
    return fmt.Sprintf("(embedded.%s)", fs.Id)
}
func (f StdFile) Close() error {
    return f.File.Close()
}
func (f StdFile) ReadContent() ([]byte, error) {
    return io.ReadAll(f.File)
}

type ZipFilesystem struct {
    Id      string
    Reader  *zip.Reader
}
func (fs ZipFilesystem) Open(path string) (File, error) {
    var f, err = fs.Reader.Open(filepath.ToSlash(path))
    if err != nil { return nil, err }
    return StdFile{f}, nil
}
func (fs ZipFilesystem) Key() string {
    return fmt.Sprintf("(zip.%s)", fs.Id)
}


