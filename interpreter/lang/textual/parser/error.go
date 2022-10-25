package parser

import (
    "rxgui/standalone/util/richtext"
    "rxgui/lang/textual/cst"
    "rxgui/lang/textual/syntax"
)

type Error struct {
    IsScannerError   bool
    ScannerError     error
    Tree             *cst.Tree
    HasExpectedPart  bool
    ExpectedPart     syntax.Id
    NodeIndex        int  // may be bigger than the index of last token
}
func (err *Error) IsEmptyTree() bool {
    return err.NodeIndex < 0 || len(err.Tree.Tokens) == 0
}
func (err *Error) Desc() richtext.Block {
    if err.IsScannerError {
        var desc richtext.Block
        desc.WriteLine(err.ScannerError.Error(), richtext.TAG_ERR)
        return desc
    }
    if err.IsEmptyTree() {
        var desc richtext.Block
        desc.WriteLine("empty input", richtext.TAG_ERR)
        return desc
    }
    var node = err.Tree.Nodes[err.NodeIndex]
    var got string
    if node.Pos >= len(err.Tree.Tokens) {
        got = "EOF"
    } else {
        got = syntax.Id2Name(err.Tree.Tokens[node.Pos].Id)
    }
    var desc richtext.Block
    if err.HasExpectedPart {
        desc.WriteSpan("Syntax unit", richtext.TAG_ERR)
        desc.WriteSpan(syntax.Id2Name(err.ExpectedPart), richtext.TAG_ERR_INLINE)
        desc.WriteSpan("expected", richtext.TAG_ERR)
    } else {
        desc.WriteSpan("Parser stuck", richtext.TAG_ERR)
    }
    desc.WriteSpan("(", richtext.TAG_ERR)
    desc.WriteSpan("got", richtext.TAG_ERR)
    desc.WriteSpan(got, richtext.TAG_ERR_INLINE)
    desc.WriteSpan(")", richtext.TAG_ERR)
    return desc
}
func (err *Error) Message() richtext.Block {
    if err.IsScannerError {
        return err.Desc()
    }
    if err.IsEmptyTree() {
        return err.Desc()
    }
    var tree = err.Tree
    var token = cst.GetNodeFirstToken(tree, err.NodeIndex)
    var point = tree.Info[token.Span.Start]
    var desc = err.Desc()
    const FOV = 5
    return cst.FormatError (
        tree.Name, tree.Code, tree.Info, tree.SpanMap,
        point, token.Span, FOV, desc,
    )
}
func (err *Error) Error() string {
    return err.Message().RenderConsole()
}


