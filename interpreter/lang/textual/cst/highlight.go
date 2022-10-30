package cst

import (
    "rxgui/interpreter/lang/textual/scanner"
    "rxgui/interpreter/lang/textual/syntax"
)


type Color struct { R,G,B int }
type ColorKey int
type ColorScheme (map[ColorKey] Color)
const (
    CK_NULL ColorKey = iota
    CK_Keyword
    CK_Type
    CK_Const
    CK_Callable
    CK_Text
    CK_Number
)
var defaultColorScheme = ColorScheme {
    CK_Keyword: Color { 0, 51, 179 },
    CK_Type: Color { 51, 110, 204 },
    CK_Const: Color { 135, 16, 148 },
    CK_Callable: Color { 40, 109, 115 },
    CK_Text: Color { 6, 125, 23 },
    CK_Number: Color { 23, 80, 235 },
}
func DefaultColorScheme() ColorScheme { return defaultColorScheme }

type HighlightRange struct {
    Span   scanner.Span
    Color  ColorKey
}
type HighlightRanges ([] HighlightRange)

func Highlight(tree *Tree) HighlightRanges {
    var out = make(HighlightRanges, 0)
    for ptr := 0; ptr < len(tree.Nodes); ptr += 1 {
        var key = highlight(ptr, tree.Nodes)
        if key != CK_NULL {
            out = append(out, HighlightRange {
                Span:  tree.Nodes[ptr].Span,
                Color: key,
            })
        }
    }
    return out
}
func highlight(ptr int, nodes ([] TreeNode)) ColorKey {
    var node = &nodes[ptr]
    switch node.Part.PartType {
    case syntax.MatchKeyword:
        return CK_Keyword
    case syntax.MatchToken:
        var id = node.Part.Id
        switch id {
        case __at:
            return CK_Keyword
        case __Text:
            return CK_Text
        case __Int, __Float, __Byte, __Char:
            return CK_Number
        default:
            var token = syntax.Id2Token(id)
            if (token.Keyword && (id != __equal) && (id != __derive)) {
                return CK_Keyword
            }
        }
    case syntax.Recursive:
        var id = node.Part.Id
        switch id {
        case __ref_base:
            if matchAncestors(&ptr, nodes, __ref) {
                if matchAncestors(&ptr, nodes, __type) {
                    return CK_Type
                }
                if matchAncestors(&ptr, nodes, __ref_term, __term) {
                    var ref_term_ptr = nodes[ptr].Children[0]
                    var new_ptr = nodes[ref_term_ptr].Children[0]
                    if (nodes[new_ptr].Length > 0) {
                        return CK_Type
                    }
                    var pipes = nodes[ptr+1]
                    if pipes.Length > 0 {
                        var pipe = nodes[pipes.Children[0]]
                        var pipe_content = nodes[pipe.Children[0]]
                        if pipe_content.Part.Id == __pipe_call {
                            return CK_Callable
                        }
                    }
                }
                if matchAncestors(&ptr, nodes, __pipe_infix) {
                    return CK_Callable
                }
                if matchAncestors(&ptr, nodes, __infix_term) {
                    return CK_Callable
                }
                if matchAncestors(&ptr, nodes, __binding_cps) {
                    return CK_Keyword
                }
            }
        case __name:
            if matchAncestors(&ptr, nodes, __ref_base) {
                break
            }
            skipListIntermediateAncestors(&ptr, nodes)
            if matchAncestors(&ptr, nodes, __case) {
                return CK_Type
            }
            { ptr = nodes[ptr].Parent
            if matchAncestors(&ptr, nodes, __stmt) {
                var stmt_content = nodes[node.Parent]
                switch stmt_content.Part.Id {
                case __decl_func:
                    return CK_Callable
                case __decl_const:
                    return CK_Const
                }
            }
            if matchAncestors(&ptr, nodes, __pattern) {
                var pattern_index = nodes[ptr].Index
                if matchAncestors(&ptr, nodes, __binding_plain) {
                    var prev = (pattern_index - 1)
                    var let_ptr = nodes[ptr].Children[prev]
                    var kw_ptr = nodes[let_ptr].Children[0]
                    if nodes[kw_ptr].Part.Id == __Const {
                        return CK_Const
                    }
                }
            }}
        }
    }
    return CK_NULL
}
var __at = syntax.Name2IdMustExist("@")
var __Text = syntax.Name2IdMustExist("Text")
var __Int = syntax.Name2IdMustExist("Int")
var __Float = syntax.Name2IdMustExist("Float")
var __Byte = syntax.Name2IdMustExist("Byte")
var __Char = syntax.Name2IdMustExist("Char")
var __equal = syntax.Name2IdMustExist("=")
var __derive = syntax.Name2IdMustExist("=>")
var __ref_base = syntax.Name2IdMustExist("ref_base")
var __ref = syntax.Name2IdMustExist("ref")
var __type = syntax.Name2IdMustExist("type")
var __ref_term = syntax.Name2IdMustExist("ref_term")
var __term = syntax.Name2IdMustExist("term")
var __pipe_call = syntax.Name2IdMustExist("pipe_call")
var __pipe_infix = syntax.Name2IdMustExist("pipe_infix")
var __infix_term = syntax.Name2IdMustExist("infix_term")
var __binding_cps = syntax.Name2IdMustExist("binding_cps")
var __name = syntax.Name2IdMustExist("name")
var __stmt = syntax.Name2IdMustExist("stmt")
var __pattern = syntax.Name2IdMustExist("pattern")
var __binding_plain = syntax.Name2IdMustExist("binding_plain")
var __Const = syntax.Name2IdMustExist("Const")
var __case = syntax.Name2IdMustExist("case")
var __decl_func = syntax.Name2IdMustExist("decl_func")
var __decl_const = syntax.Name2IdMustExist("decl_const")
func matchAncestors(ptr *int, nodes ([] TreeNode), ancestors ...syntax.Id) bool {
    for {
        if len(ancestors) > 0 {
            var id = ancestors[0]
            ancestors = ancestors[1:]
            var parent = nodes[*ptr].Parent
            if parent >= 0 {
                if nodes[parent].Part.Id == id {
                    *ptr = parent
                    continue
                } else {
                    return false
                }
            } else {
                return false
            }
        } else {
            return true
        }
    }
}
func skipListIntermediateAncestors(ptr *int, nodes ([] TreeNode)) {
    for {
        var parent = nodes[*ptr].Parent
        if parent >= 0 {
            var name = syntax.Id2Name(nodes[parent].Part.Id)
            if syntax.IsList(name) {
                *ptr = parent
                continue
            }
        }
        break
    }
}


