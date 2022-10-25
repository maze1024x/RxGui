package syntax

import (
    "testing"
    "fmt"
    "sort"
    "strings"
    "unicode"
)


func TestOutputIntellijBnfGrammar(_ *testing.T) {
    outputGeneratedIntellijBnfGrammar()
}
func TestOutputTreeSitterGrammar(_ *testing.T) {
    outputGeneratedTreeSitterGrammar()
}

func outputGeneratedIntellijBnfGrammar() {
    var lex_mapping = make(map[string] [] string)
    var tokens_mapping = make(map[string] string)
    for _, token := range __Tokens {
        var name = token.Name
        if ((name == "Blank") || (name == "LF")) {
            name = "Blank"
        }
        var output_name string
        if name == "Blank" {
            output_name = "WHITE_SPACE"
        } else if !(unicode.IsLetter(rune(name[0]))) {
            var buf strings.Builder
            for _, ch := range name {
                buf.WriteString(fmt.Sprintf("%d00", ch))
            }
            var suffix = buf.String()
            output_name = ("SYM" + suffix)
        } else {
            if name == "Text" {
                // workaround: getText() conflicts with inherited method
                output_name = "TOKEN_TEXT"
            } else if name == "Name" {
                // workaround: getName() conflicts with inherited method
                output_name = "TOKEN_NAME"
            } else {
                output_name = strings.ToUpper(name)
            }
        }
        tokens_mapping[name] = output_name
        var content = strings.TrimPrefix(token.Pattern.String(), "^")
        lex_mapping[name] = append(lex_mapping[name], content)
    }
    println("tokens=[")
    var lex_visited = make(map[string] bool)
    for _, token := range __Tokens {
        var name = token.Name
        if lex_visited[name] { continue }; lex_visited[name] = true
        var contents, exists = lex_mapping[name]
        if !(exists) { continue }
        var output_name = tokens_mapping[name]
        var tmp = strings.Join(contents, "|")
        tmp = strings.ReplaceAll(tmp, "\"", "\\\"")
        tmp = strings.ReplaceAll(tmp, "'", "\\'")
        var wrapped = fmt.Sprintf("\"regexp:%s\"", tmp)
        println(fmt.Sprintf("    %s=%s", output_name, wrapped))
    }
    println("]")
    const replPartsPrefix = "repl_"
    var id_list = make([] int, 0)
    for id := range rules {
        id_list = append(id_list, int(id))
    }
    sort.Ints(id_list)
    for _, id := range id_list {
        var id = Id(id)
        var rule = rules[id]
        var rule_name = mapId2Name[id]
        if rule_name == "" { panic("something went wrong") }
        if rule.Generated {
            continue
        }
        if strings.HasPrefix(rule_name, replPartsPrefix) {
            continue
        }
        var conflict = func(rule_name string) bool {
            // conflict with keyword part name (leading @ stripped)
            var _, a = mapName2Id["@" + rule_name]
            // conflict with token part name (in case-insensitive context)
            var _, b = mapName2Id[strings.ToUpper(rule_name[:1]) + rule_name[1:]]
            return (a || b || (rule_name == "string"))
        }
        const conflict_prefix = "node_"
        if conflict(rule_name) {
            rule_name = (conflict_prefix + rule_name)
        }
        var refer_part = func(part_name string) string {
            switch GetPartType(part_name) {
            case MatchKeyword:
                return strings.TrimPrefix(part_name, "@")
            case MatchToken:
                return tokens_mapping[part_name]
            case Recursive:
                if conflict(part_name) {
                    return (conflict_prefix + part_name)
                } else {
                    return part_name
                }
            default:
                panic("impossible branch")
            }
        }
        var branch_expr = make([] string, len(rule.Branches))
        for i, b := range rule.Branches {
            var part_expr = make([] string, len(b.Parts))
            for j, p := range b.Parts {
                var part_name = mapId2Name[p.Id]
                var part_rule = rules[p.Id]
                switch p.PartType {
                case MatchKeyword, MatchToken:
                    part_expr[j] = refer_part(part_name)
                case Recursive:
                    if part_rule.Generated {
                        var info = part_rule.GenInfo
                        switch info.Kind {
                        case RuleGenList:
                            var sep = info.Sep
                            var nullable = part_rule.Nullable
                            var ref = refer_part(info.Item)
                            if sep == "" {
                                if nullable { ref += "*" } else { ref += "+" }
                            } else {
                                // assume sep is a keyword for now
                                var sep_output_name = tokens_mapping[sep]
                                var tail = (sep_output_name + " " + ref)
                                ref += (" { " + tail + " }*")
                                if nullable { ref = ("[ " + ref + " ]") }
                            }
                            part_expr[j] = ref
                        case RuleGenListTail:
                            panic("something went wrong")
                        case RuleGenOptional:
                            var ref = refer_part(info.Item)
                            ref += "?"
                            part_expr[j] = ref
                        default:
                            panic("impossible branch")
                        }
                    } else {
                        var ref = refer_part(part_name)
                        if part_rule.Nullable {
                            ref += "?"
                        }
                        part_expr[j] = ref
                    }
                default:
                    panic("impossible branch")
                }
            }
            branch_expr[i] = strings.Join(part_expr, " ")
        }
        var expr = strings.Join(branch_expr, " | ")
        println(fmt.Sprintf("%s ::= %s", rule_name, expr))
    }
}
func outputGeneratedTreeSitterGrammar() {
    println("// TOKENS")
    var lex_mapping = make(map[string] [] string)
    var lex_kw_mapping = make(map[string] string)
    var is_all_symbols = func(str string) bool {
        for _, char := range str {
            if !(unicode.IsSymbol(char) || unicode.IsPunct(char)) {
                return false
            }
        }
        return true
    }
    for _, token := range __Tokens {
        var name = token.Name
        var content = strings.TrimPrefix(token.Pattern.String(), "^")
        if token.Keyword {
            lex_kw_mapping[name] = content
        } else if is_all_symbols(name) {
            lex_kw_mapping[name] = name
        } else {
            lex_mapping[name] = append(lex_mapping[name], content)
        }
    }
    var lex_visited = make(map[string] bool)
    for _, token := range __Tokens {
        var name = token.Name
        var values, in_lex_mapping = lex_mapping[name]
        if !(in_lex_mapping) { continue }
        if lex_visited[name] { continue }; lex_visited[name] = true
        var buf strings.Builder
        buf.WriteString("const ")
        buf.WriteString(name)
        buf.WriteString(" = ")
        if token.Keyword {
            var str = fmt.Sprintf("'%s'", values[0])
            buf.WriteString(str)
        } else {
            buf.WriteString("/")
            if len(values) == 0 {
                panic("something went wrong")
            } else if len(values) == 1 {
                var raw = values[0]
                var escaped = strings.ReplaceAll(raw, "/", "\\/")
                buf.WriteString(escaped)
            } else {
                var raw = strings.Join(values, "|")
                var escaped = strings.ReplaceAll(raw, "/", "\\/")
                buf.WriteString(escaped)
            }
            buf.WriteString("/")
        }
        println(buf.String())
    }
    println("// RULES")
    const replPartsPrefix = "repl_"
    var keyword = func(content string) string {
        return fmt.Sprintf("'%s'", strings.ReplaceAll(content, "\\", "\\\\"))
    }
    var ref = func(name string) string {
        return fmt.Sprintf("$.%s", name)
    }
    var choice = func(list ...string) string {
        return fmt.Sprintf("choice(%s)", strings.Join(list, ", "))
    }
    var seq = func(list ...string) string {
        return fmt.Sprintf("seq(%s)", strings.Join(list, ", "))
    }
    var optional = func(content string, is_optional bool) string {
        if is_optional {
            return fmt.Sprintf("optional(%s)", content)
        } else {
            return content
        }
    }
    var repeat = func(content string, nullable bool) string {
        if nullable {
            return fmt.Sprintf("repeat(%s)", content)
        } else {
            return fmt.Sprintf("repeat1(%s)", content)
        }
    }
    var match_kw_opt_rule = func(rule Rule) (string, bool) {
        var num_branches = len(rule.Branches)
        if num_branches == 0 { panic("something went wrong") }
        var options = make([] string, num_branches)
        for i, b := range rule.Branches {
            var num_parts = len(b.Parts)
            if num_parts != 1 {
                return "", false
            }
            var part = b.Parts[0]
            var name = mapId2Name[part.Id]
            if name == "" { panic("something went wrong") }
            var token_kw_content, is_lex_kw = lex_kw_mapping[name]
            var is_cond_kw = (part.PartType == MatchKeyword)
            if is_lex_kw {
                options[i] = keyword(token_kw_content)
            } else if is_cond_kw {
                var cond_kw_content = strings.TrimPrefix(name, "@")
                options[i] = keyword(cond_kw_content)
            } else {
                return "", false
            }
        }
        if len(options) == 1 {
            return optional(options[0], rule.Nullable), true
        } else {
            return optional(choice(options...), rule.Nullable), true
        }
    }
    var id_list = make([] int, 0)
    for id := range rules {
        id_list = append(id_list, int(id))
    }
    sort.Ints(id_list)
    for _, id := range id_list {
        var id = Id(id)
        var rule = rules[id]
        var rule_name = mapId2Name[id]
        if rule_name == "" { panic("something went wrong") }
        if rule.Generated {
            continue
        }
        if _, is_kw_opt_rule := match_kw_opt_rule(rule); is_kw_opt_rule {
            continue
        }
        if strings.HasPrefix(rule_name, replPartsPrefix) {
            continue
        }
        var buf strings.Builder
        if id == DefaultRoot() {
            buf.WriteString("source_file")
        } else {
            buf.WriteString(rule_name)
        }
        buf.WriteString(": $ => ")
        var branch_contents = make([] string, len(rule.Branches))
        for i, b := range rule.Branches {
            var part_contents = make([] string, len(b.Parts))
            for j, part := range b.Parts {
                var part_name = mapId2Name[part.Id]
                if part_name == "" { panic("something went wrong") }
                var part_rule = rules[part.Id]
                switch part.PartType {
                case MatchKeyword:
                    // cond kw
                    var content = strings.TrimPrefix(part_name, "@")
                    part_contents[j] = keyword(content)
                case MatchToken:
                    var kw_value, as_keyword = lex_kw_mapping[part_name]
                    if as_keyword {
                        // lex kw
                        part_contents[j] = keyword(kw_value)
                    } else {
                        part_contents[j] = part_name
                    }
                case Recursive:
                    if part_rule.Generated {
                        var info = part_rule.GenInfo
                        switch info.Kind {
                        case RuleGenList:
                            var item = info.Item
                            var sep = info.Sep
                            var nullable = part_rule.Nullable
                            if sep == "" {
                                part_contents[j] = repeat(ref(item), nullable)
                            } else {
                                // assume sep is a keyword for now
                                var tail = seq(keyword(sep), ref(item))
                                var list = seq(ref(item), repeat(tail, true))
                                part_contents[j] = optional(list, nullable)
                            }
                        case RuleGenListTail:
                            panic("something went wrong")
                        case RuleGenOptional:
                            var item = info.Item
                            part_contents[j] = optional(ref(item), true)
                        default:
                            panic("impossible branch")
                        }
                    } else if content, ok := match_kw_opt_rule(part_rule); ok {
                        part_contents[j] = content
                    } else {
                        var nullable = part_rule.Nullable
                        part_contents[j] = optional(ref(part_name), nullable)
                    }
                default:
                    panic("impossible branch")
                }
            }
            if len(part_contents) == 0 {
                panic("something went wrong")
            } else if len(part_contents) == 1 {
                branch_contents[i] = part_contents[0]
            } else {
                branch_contents[i] = seq(part_contents...)
            }
        }
        if len(branch_contents) == 0 {
            panic("something went wrong")
        } else if len(branch_contents) == 1 {
            buf.WriteString(branch_contents[0])
        } else {
            buf.WriteString(choice(branch_contents...))
        }
        buf.WriteString(",")
        println(buf.String())
    }
}


