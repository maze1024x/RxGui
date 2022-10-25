package source

import (
	"fmt"
	"rxgui/util/ctn"
	"rxgui/util/richtext"
)


type Location struct {
	File  File
	Pos   Position
}
// note: File should be comparable
type File interface {
	GetKey() FileKey
	DescribePosition(Position) string
	FormatMessage(Position, richtext.Block) richtext.Block
}
type FileKey struct {
	Context  string
	Path     string
}
type Position struct {
	Span
}
func (key FileKey) String() string {
	return fmt.Sprintf("%s[%s]", key.Context, key.Path)
}
func (l Location) FileDesc() string {
	if l.File != nil {
		return l.File.GetKey().String()
	} else {
		return "[]"
	}
}
func (l Location) PosDesc() string {
	if l.File != nil {
		return l.File.DescribePosition(l.Pos)
	} else {
		return "()"
	}
}
func (l Location) FormatMessage(b richtext.Block) richtext.Block {
	if l.File != nil {
		return l.File.FormatMessage(l.Pos, b)
	} else {
		return b
	}
}
func (l Location) FormatMessageLite(b richtext.Block) richtext.Block {
	if l.File != nil {
		var desc = fmt.Sprintf("%s %s", l.FileDesc(), l.PosDesc())
		return b.WithLeadingLine(desc, richtext.TAG_B)
	} else {
		return b
	}
}
func CompareFileKey(a FileKey, b FileKey) ctn.Ordering {
	var ord = ctn.StringCompare(a.Context, b.Context)
	if ord == ctn.Equal {
		return ctn.StringCompare(a.Path, b.Path)
	} else {
		return ord
	}
}


