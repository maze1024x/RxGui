package ast

import (
    "strings"
    "rxgui/interpreter/lang/source"
)


func GetDocContent(raw ([] Doc)) string {
    var buf strings.Builder
    for _, line := range raw {
        var t = source.CodeToString(line.RawContent)
        t = strings.TrimPrefix(t, "///")
        t = strings.TrimPrefix(t, " ")
        t = strings.TrimRight(t, " \r")
        buf.WriteString(t)
        buf.WriteRune('\n')
    }
    return buf.String()
}

func MaybeId2String(maybe_id MaybeIdentifier) string {
    if id, ok := maybe_id.(Identifier); ok {
        return Id2String(id)
    } else {
        return ""
    }
}
func Id2String(id Identifier) string {
    return source.CodeToString(id.Name)
}
func String2Id(s string, node Node) Identifier {
    return Identifier {
        Node: node,
        Name: source.StringToCode(s),
    }
}
func String2Ref(s string, node Node) Ref {
    return Ref {
        Node: node,
        Base: RefBase {
            Node: node,
            Item: String2Id(s, node),
        },
    }
}
func Strings2Ref(ns string, item string, node Node) Ref {
    return Ref {
        Node: node,
        Base: RefBase {
            Node: node,
            NS:   String2Id(ns, node),
            Item: String2Id(item, node),
        },
    }
}

func WrapTermAsExpr(term VariousTerm) Expr {
    return Expr {
        Node:     term.Node,
        Term:     term,
        Pipeline: nil,
    }
}
func WrapBlockAsExpr(block Block) Expr {
    return WrapTermAsExpr(VariousTerm {
        Node: block.Node,
        Term: block,
    })
}
func WrapExprIntoTerm(expr Expr) VariousTerm {
    return VariousTerm {
        Node: expr.Node,
        Term: Block {
            Node:   expr.Node,
            Return: expr,
        },
    }
}
func WrapExprIntoBlock(expr Expr) Block {
    return Block {
        Node:     expr.Node,
        Bindings: [] VariousBinding {},
        Return:   expr,
    }
}

func GetStandaloneTerm(expr Expr) (VariousTerm, bool) {
    if (len(expr.Casts) == 0) && (len(expr.Pipeline) == 0) {
        return expr.Term, true
    } else {
        return VariousTerm {}, false
    }
}
func GetStandaloneRef(expr Expr) (Ref, string, bool, bool) {
    if term, ok := GetStandaloneTerm(expr); ok {
        if ref_term, ok := term.Term.(RefTerm); ok {
            if new_, ok := ref_term.New.(New); ok {
                var tag = MaybeId2String(new_.Tag)
                return ref_term.Ref, tag, true, true
            } else {
                return ref_term.Ref, "", false, true
            }
        }
    }
    return Ref{}, "", false, false
}

func SplitImplicitRef(name string) (string, string, bool) {
    const sep = "/"
    return strings.Cut(name, sep)
}

func ReplCmdGetExpr(cmd ReplCmd) Expr {
    switch cmd := cmd.(type) {
    case ReplAssign:
        return cmd.Expr
    case ReplRun:
        return cmd.Expr
    case ReplEval:
        return cmd.Expr
    default:
        panic("impossible branch")
    }
}


