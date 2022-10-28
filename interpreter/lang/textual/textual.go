package textual

import (
    "io"
    "fmt"
    "strconv"
    "reflect"
    "rxgui/interpreter/lang/source"
    "rxgui/interpreter/lang/textual/ast"
    "rxgui/interpreter/lang/textual/scanner"
    "rxgui/interpreter/lang/textual/syntax"
    "rxgui/interpreter/lang/textual/parser"
    "rxgui/interpreter/lang/textual/transformer"
)

func DebugParser(file io.Reader, name string, root syntax.Id) bool {
    fmt.Println("------------------------------------------------------")
    fmt.Printf("\033[1m%s\033[0m\n", name)
    var code_bytes, e = io.ReadAll(file)
    if e != nil { panic(e) }
    var code = source.DecodeUtf8ToCode(code_bytes)
    var tokens, info, _, s_err = scanner.Scan(code)
    if s_err != nil {
        fmt.Println(s_err.Error())
        return false
    }
    fmt.Println("------------------------------------------------------")
    fmt.Println("Tokens:")
    for i, token := range tokens {
        fmt.Printf(
            "(%v) at [%v, %v] (%v, %v) %v: %v\n",
            i,
            token.Span.Start,
            token.Span.End,
            info[token.Span.Start].Row,
            info[token.Span.Start].Col,
            syntax.Id2Name(token.Id),
            source.CodeToString(token.Content),
        )
    }
    var tree, err = parser.Parse(code, root, name)
    fmt.Println("------------------------------------------------------")
    fmt.Println("CST Nodes:")
    parser.PrintFlatTree(tree.Nodes)
    fmt.Println("------------------------------------------------------")
    fmt.Println("CST:")
    parser.PrintTree(tree)
    if err != nil {
        var msg = err.Message()
        fmt.Println(msg.RenderConsole())
        return false
    } else {
        fmt.Println("------------------------------------------------------")
        fmt.Println("AST:")
        var key = source.FileKey {
            Context: "debug",
            Path:    name,
        }
        var transformed = transformer.Transform(tree, key)
        transformer.PrintNode(reflect.ValueOf(transformed))
        return true
    }
}

func OutputAdHocIndex(code_bytes ([] byte), name string) {
    var root = syntax.DefaultRoot()
    var code = source.DecodeUtf8ToCode(code_bytes)
    var tree, err = parser.Parse(code, root, name)
    if err != nil {
        var msg = err.Message()
        fmt.Println(msg.RenderConsole())
    }
    var key = source.FileKey {
        Context: "indexing",
        Path:    name,
    }
    var transformed = transformer.Transform(tree, key)
    var root_node = transformed.(ast.Root)
    var node_str = func(node ast.Node) string {
        var span = node.Location.Pos.Span
        return source.CodeToString(tree.Code[span.Start: span.End])
    }
    var ns = ast.MaybeId2String(root_node.Namespace)
    var stmts = root_node.Statements
    fmt.Println("{")
    fmt.Printf(`"ns": %s`, strconv.Quote(ns))
    fmt.Println(",")
    fmt.Println(`"items": [`)
    var first = true
    var add = func(fields ...string) {
        if first {
            first = false
        } else {
            fmt.Println(",")
        }
        fmt.Print("[")
        for i, field := range fields {
            fmt.Print(strconv.Quote(field))
            if i != len(fields)-1 {
                fmt.Print(", ")
            }
        }
        fmt.Print("]")
    }
    for _, stmt := range stmts {
        switch S := stmt.Statement.(type) {
        case ast.DeclType:
            var type_name = node_str(S.Name.Node)
            add("type", type_name, node_str(S.TypeDef.Node))
            switch D := S.TypeDef.TypeDef.(type) {
            case ast.Enum:
                for _, item := range D.Items {
                    add("enum-value", node_str(item.Node), type_name)
                }
            case ast.Record:
                for _, field := range D.Def.Fields {
                    add("field", node_str(field.Name.Node), node_str(field.Type.Node))
                }
            case ast.Interface:
                for _, method := range D.Methods {
                    add("abstract-method", node_str(method.Name.Node), node_str(method.Type.Node))
                }
            }
        case ast.DeclConst:
            add("const", node_str(S.Name.Node), node_str(S.Type.Node))
        case ast.DeclFunction:
            var kind = (func() string {
                if S.Operator { return "operator" } else { return "function" }
            })()
            add(kind, node_str(S.Name.Node), node_str(S.Signature.Output.Node), node_str(S.Signature.Node))
        case ast.DeclMethod:
            add("method", node_str(S.Name.Node), node_str(S.Type.Node))
        }
    }
    fmt.Println("]")
    fmt.Println("}")
}


