package textual

import (
	"io"
	"fmt"
	"reflect"
	"rxgui/lang/source"
	"rxgui/lang/textual/scanner"
	"rxgui/lang/textual/syntax"
	"rxgui/lang/textual/parser"
	"rxgui/lang/textual/transformer"
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


