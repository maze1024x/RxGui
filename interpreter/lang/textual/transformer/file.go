package transformer

import (
    "fmt"
    "rxgui/util/richtext"
    "rxgui/interpreter/lang/source"
    "rxgui/interpreter/lang/textual/cst"
    "rxgui/interpreter/lang/textual/scanner"
)

type File struct {
    Key      source.FileKey
    Code     source.Code
    Info     scanner.RowColInfo
    SpanMap  scanner.RowSpanMap
}
func makeFile(tree *cst.Tree, key source.FileKey) *File {
    return &File {
        Key:     key,
        Code:    tree.Code,
        Info:    tree.Info,
        SpanMap: tree.SpanMap,
    }
}
func (f *File) GetKey() source.FileKey {
    return f.Key
}
func (f *File) DescribePosition(pos source.Position) string {
    var point = f.Info[pos.Start]
    return fmt.Sprintf("(row %d, column %d)", point.Row, point.Col)
}
func (f *File) FormatMessage(pos source.Position, desc richtext.Block) richtext.Block {
    var span = pos.Span
    var point scanner.Point
    if span.Start < len(f.Info) {
        point = f.Info[span.Start]
    }
    var name = f.Key.String()
    const FOV = 5
    return cst.FormatError (
        name, f.Code, f.Info, f.SpanMap,
        point, span, FOV, desc,
    )
}


