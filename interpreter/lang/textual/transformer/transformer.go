package transformer

import "fmt"
import "reflect"
import "strings"
import "rxgui/interpreter/lang/source"
import "rxgui/interpreter/lang/textual/cst"
import "rxgui/interpreter/lang/textual/ast"
import "rxgui/interpreter/lang/textual/syntax"
import "rxgui/interpreter/lang/textual/parser"


/**
 *  Syntax Tree Transformer
 *
 *  This package is responsible for transforming the CSTs output
 *    by the `parser` package into typed and well-structured ASTs
 *    according to the declarative configurations defined in
 *    the `ast` package.
 */

type Tree = *cst.Tree
type Pointer = int
type Context = map[string] interface{}
type Transformer = func(Tree, Pointer) reflect.Value

func Transform(tree Tree, key source.FileKey) interface{} {
    var dive func(Tree, Pointer, ([] syntax.Id)) (Pointer, bool)
    dive = func (tree Tree, ptr Pointer, path ([] syntax.Id)) (Pointer, bool) {
        if len(path) == 0 {
            return ptr, true
        }
        var parser_node = &tree.Nodes[ptr]
        var L = parser_node.Length
        for i := 0; i < L; i += 1 {
            var child_ptr = parser_node.Children[i]
            var child = &tree.Nodes[child_ptr]
            if child.Part.Id == path[0] {
                if child.Part.PartType == syntax.Recursive && child.Length == 0 {
                    return -1, false
                } else {
                    return dive(tree, child_ptr, path[1:])
                }
            }
        }
        return -1, false
    }
    var transform Transformer
    transform = func (tree Tree, ptr Pointer) reflect.Value {
        var parser_node = &tree.Nodes[ptr]
        var info = getNodeInfoById(parser_node.Part.Id)
        var node = reflect.New(info.Type)
        node.Elem().FieldByName("Node").Set(reflect.ValueOf(ast.Node {
            Location: source.Location {
                File: makeFile(tree, key),
                Pos:  source.Position { Span: parser_node.Span },
            },
        }))
        var transform_dived = func (child_info *NodeChildInfo, f Transformer) {
            var field_index = child_info.FieldIndex
            var field = info.Type.Field(field_index)
            var field_value = node.Elem().Field(field_index)
            var path = child_info.DivePath
            var dived_ptr, exists = dive(tree, ptr, path)
            if exists {
                field_value.Set(f(tree, dived_ptr))
            } else {
                if child_info.Fallback != -1 {
                    var fallback_id = child_info.Fallback
                    var L = parser_node.Length
                    for i := 0; i < L; i += 1 {
                        var child_ptr = parser_node.Children[i]
                        var child_parser_node = &tree.Nodes[child_ptr]
                        if child_parser_node.Part.Id == fallback_id {
                            var value = transform(tree, child_ptr)
                            if field.Type.Kind() == reflect.Slice {
                                var t = reflect.MakeSlice(field.Type, 1, 1)
                                t.Index(0).Set(value)
                                field_value.Set(t)
                            } else {
                                field_value.Set(value)
                            }
                            break
                        }
                    }
                } else if child_info.Optional {
                    if field.Type.Kind() == reflect.Slice {
                        var empty_slice = reflect.MakeSlice(field.Type, 0, 0)
                        field_value.Set(empty_slice)
                    } else if field.Type.Kind() == reflect.Bool {
                        field_value.Set(reflect.ValueOf(false))
                    }
                } else {
                    var path_strlist = make([] string, len(path))
                    for i, segment := range path {
                        path_strlist[i] = syntax.Id2Name(segment)
                    }
                    var path_str = strings.Join(path_strlist, ".")
                    panic(fmt.Sprintf (
                        "transform(): path `%v` cannot be found in `%v`",
                        path_str, syntax.Id2Name(parser_node.Part.Id),
                    ))
                }
            }
        }
        if info.First != -1 {
            if parser_node.Length > 0 {
                var field_value = node.Elem().Field(info.First)
                field_value.Set(transform(tree, parser_node.Children[0]))
            } else {
                panic(fmt.Sprintf (
                    "transform(): cannot get first child of empty node `%v`",
                    syntax.Id2Name(parser_node.Part.Id),
                ))
            }
        }
        if info.Last != -1 {
            var L = parser_node.Length
            if L > 0 {
                var field_value = node.Elem().Field(info.Last)
                field_value.Set(transform(tree, parser_node.Children[L-1]))
            } else {
                panic(fmt.Sprintf (
                    "transform(): cannot get last child of empty node `%v`",
                    syntax.Id2Name(parser_node.Part.Id),
                ))
            }
        }
        for _, child_info_group := range info.Children {
        for _, child_info := range child_info_group {
            transform_dived(&child_info, transform)
        }}
        for _, child_info_group := range info.Strings {
        for _, child_info := range child_info_group {
            transform_dived (
                &child_info,
                func (tree Tree, dived_ptr Pointer) reflect.Value {
                    var dived_node = &tree.Nodes[dived_ptr]
                    if dived_node.Part.PartType == syntax.MatchToken {
                        var content = getTokenContent(tree, dived_ptr)
                        return reflect.ValueOf(content)
                    } else {
                        panic(fmt.Sprintf (
                            "cannot get token content of non-token part %v",
                            syntax.Id2Name(dived_node.Part.Id),
                        ))
                    }
                },
            )
        }}
        for _, child_info_group := range info.Options {
        for _, child_info := range child_info_group {
            transform_dived (
                &child_info,
                func (tree Tree, _ Pointer) reflect.Value {
                    return reflect.ValueOf(true)
                },
            )
        }}
        for _, list_info_group := range info.Lists {
        for _, list_info := range list_info_group {
            transform_dived (
                &list_info.NodeChildInfo,
                func (tree Tree, dived_ptr Pointer) reflect.Value {
                    var item_id = list_info.ItemId
                    var tail_id = list_info.TailId
                    var item_ptrs = flatSubTree(tree, dived_ptr, item_id, tail_id)
                    var field_index = list_info.FieldIndex
                    var field = info.Type.Field(field_index)
                    if field.Type.Kind() != reflect.Slice {
                        panic("cannot transform list to non-slice field")
                    }
                    var N = len(item_ptrs)
                    var slice = reflect.MakeSlice(field.Type, N, N)
                    for i, item_ptr := range item_ptrs {
                        slice.Index(i).Set(transform(tree, item_ptr))
                    }
                    return slice
                },
            )
        }}
        return node.Elem()
    }
    return transform(tree, 0).Interface()
}

func flatSubTree (
    tree     Tree,       ptr   Pointer,
    extract  syntax.Id,  next  syntax.Id,
) ([] Pointer) {
    var sequence = make([] int, 0)
    for notEmpty(tree, ptr) {
        var extract_ptr = -1
        var next_ptr = -1
        var node = &tree.Nodes[ptr]
        var L = node.Length
        for i := 0; i < L; i += 1 {
            var child_ptr = node.Children[i]
            var child = &tree.Nodes[child_ptr]
            if child.Part.Id == extract {
                extract_ptr = child_ptr
            }
            if child.Part.Id == next {
                next_ptr = child_ptr
            }
        }
        if extract_ptr == -1 {
            panic("cannot extract part " + syntax.Id2Name(extract))
        }
        sequence = append(sequence, extract_ptr)
        if next_ptr == -1 {
            panic("next part " + syntax.Id2Name(next) + " not found")
        }
        ptr = next_ptr
    }
    return sequence
}
func getTokenContent(tree Tree, ptr int) source.Code {
    var node = &tree.Nodes[ptr]
    return tree.Tokens[node.Pos + node.Amount - 1].Content
}
func notEmpty(tree Tree, ptr Pointer) bool {
    return tree.Nodes[ptr].Length > 0
}

func printNodeRecursively (
    buf *strings.Builder,
    node reflect.Value, name string, depth int, is_last ([] bool),
) {
    var code_zero source.Code
    var code_rt = reflect.TypeOf(code_zero)
    var T = node.Type()
    if T.Kind() == reflect.Interface {
        printNodeRecursively(buf, node.Elem(), name, depth, is_last)
        return
    }
    const INC = 2
    const SPACE = " "
    parser.Repeat(depth+1, func (i int) {
        if depth > 0 && i < depth {
            if is_last[i] {
                parser.Fill(buf, INC, "", SPACE)
            } else {
                parser.Fill(buf, INC, "│", SPACE)
            }
        } else {
            if is_last[depth] {
                parser.Fill(buf, INC, "└", "─")
            } else {
                parser.Fill(buf, INC, "├", "─")
            }
        }
    })
    if T.Kind() == reflect.Struct && T.NumField() > 0 {
        buf.WriteString("┬─")
    } else if T.Kind() == reflect.Slice && !(T.AssignableTo(code_rt)) {
        if node.Len() > 0 {
            buf.WriteString("┬─")
        } else {
            buf.WriteString("──")
        }
    } else {
        buf.WriteString("──")
    }
    fmt.Fprintf(buf, "\033[1m\033[%vm", parser.GetANSIColor(depth))
    fmt.Fprintf(buf, "[%v] %v", name, T.String())
    fmt.Fprintf(buf, "\033[0m")
    fmt.Fprintf(buf, "\033[%vm", parser.GetANSIColor(depth))
    buf.WriteRune(' ')
    var newline = func () {
        fmt.Fprintf(buf, "\033[0m\n")
    }
    switch T.Kind() {
    case reflect.Bool:
        fmt.Fprintf(buf, "<%v>", node.Interface().(bool))
        newline()
    case reflect.Slice:
        if T.AssignableTo(code_rt) {
            fmt.Fprintf(buf, "'%v'",
                source.CodeToString(node.Interface().(source.Code)))
            newline()
        } else {
            var L = node.Len()
            if L == 0 {
                buf.WriteString("(empty)")
            }
            fmt.Fprintf(buf, "\033[0m")
            buf.WriteRune('\n')
            for i := 0; i < L; i += 1 {
                var child = node.Index(i)
                is_last = append(is_last, i == L-1)
                printNodeRecursively(
                    buf, child, fmt.Sprintf("%d", i),
                    (depth + 1), is_last,
                )
                is_last = is_last[:len(is_last)-1]
            }
        }
    case reflect.Struct:
        var L = node.NumField()
        newline()
        for i := 0; i < L; i += 1 {
            var field = node.Field(i)
            var field_info = T.Field(i)
            if field_info.Name == "Node" {
                continue
            }
            if field_info.Type.Kind() == reflect.Interface && field.IsNil() {
                continue
            }
            is_last = append(is_last, i == L-1)
            printNodeRecursively(
                buf, field, field_info.Name,
                (depth + 1), is_last,
            )
            is_last = is_last[:len(is_last)-1]
        }
    default:
        newline()
    }
}
func PrintNode(node reflect.Value) {
    var buf strings.Builder
    var is_last = make([] bool, 0, 1000)
    is_last = append(is_last, true)
    printNodeRecursively(&buf, node, "Root", 0, is_last)
    fmt.Println(buf.String())
}


