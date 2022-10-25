package syntax

import "strings"


const MAX_NUM_PARTS = 10

type Id int

type Token struct {
    Name     string
    Pattern  Regexp
    Keyword  bool
}

type Rule struct {
    Id         Id
    Nullable   bool
    Branches   [] Branch
    Generated  bool
    GenInfo    RuleGenInfo
}
type Branch struct {
    Parts  [] Part
}
type Part struct {
    Id        Id
    PartType  PartType
    Required  bool
}
type PartType int
const (
    // MatchKeyword stands for matching **parser** keywords.
    // note: two kinds of keywords:
    //   - parser keyword (conditional keyword)
    //   - lexer keyword (strict keyword, has a token definition)
    MatchKeyword PartType = iota
    MatchToken
    Recursive
)
type RuleGenInfo struct {
    Kind  RuleGenKind
    Item  string
    Sep   string
}
type RuleGenKind int
const (
    RuleGenList  RuleGenKind  =  iota
    RuleGenListTail
    RuleGenOptional
)

func GetPartType(name string) PartType {
    var is_keyword = strings.HasPrefix(name, "@") && len(name) > 1
    if is_keyword {
        return MatchKeyword
    } else {
        var t = name[0:1]
        if strings.ToUpper(t) == t {
            // the name starts with capital letter
            return MatchToken
        } else {
            // the name starts with small letter
            return Recursive
        }
    }
}
func EscapePartName(name string) string {
    if strings.HasPrefix(name, "_") && __EscapeMap[name] != "" {
        return __EscapeMap[name]
    } else {
        return name
    }
}

var mapId2Name = make([] string, 0)
func Id2Name(id Id) string {
    return mapId2Name[id]
}
var mapName2Id = make(map[string] Id)
func Name2IdMustExist(str string) Id {
    var id, ok = mapName2Id[str]
    if !(ok) { panic("something went wrong") }
    return id
}
func Name2Id(str string) (Id, bool) {
    var id, ok = mapName2Id[str]
    return id, ok
}
func IsList(name string) bool {
    return strings.ContainsAny(name, "*+")
}
func IsListTail(name string) bool {
    return (IsList(name) && strings.HasSuffix(name, "_tail"))
}
func IsListMustNotEmpty(name string) bool {
    return strings.ContainsAny(name, "+")
}
var mapId2ConditionalKeyword = make(map[Id] ([] rune))
func Id2ConditionalKeyword(id Id) ([] rune) {
    return mapId2ConditionalKeyword[id]
}
var rules = make(map[Id] Rule)
func Id2Rule(id Id) Rule {
    return rules[id]
}
var tokenIndexes = make(map[Id] int)
func Id2Token(id Id) Token {
    return __Tokens[tokenIndexes[id]]
}

func assignId2Name(name string) Id {
    var _, exists = mapName2Id[name]
    if exists {
        panic("conflicting name: " + name)
    }
    var id = Id(len(mapId2Name))
    mapName2Id[name] = id
    mapId2Name = append(mapId2Name, name)
    return id
}
func assignId2Tokens() {
    for i, token := range __Tokens {
        var name = token.Name
        var _, exists = mapName2Id[name]
        if !(exists) {
            var id = assignId2Name(name)
            tokenIndexes[id] = i
        }
    }
}
func assignId2CondKeywords() {
    for _, name := range __ConditionalKeywords {
        var keyword = []rune(strings.TrimLeft(name, "@"))
        if len(keyword) == 0 { panic("empty keyword") }
        var id = assignId2Name(name)
        mapId2ConditionalKeyword[id] = keyword
    }
}
func assignId2Rules() {
    for _, def := range __SyntaxDefinition {
        var t = strings.Split(def, "=")
        var u = strings.Trim(t[0], " ")
        var rule_name = strings.TrimRight(u, "?")
        assignId2Name(rule_name)
    }
}

func parseRules() {
    for _, def := range __SyntaxDefinition {
        var p = strings.Index(def, "=")
        if (p == -1) { panic(def + ": invalid rule: missing =") }
        // name = ...
        var name_def = strings.Trim(def[:p], " ")
        var name = strings.TrimSuffix(name_def, "?")
        var nullable = strings.HasSuffix(name_def, "?")
        var id, exists = mapName2Id[name]
        if (!exists) { panic("undefined rule name: " + name) }
        // ... = branches
        var branches_def = strings.Trim(def[p+1:], " ")
        if (branches_def == "") { panic(name + ": missing rule definition") }
        var branch_defs = strings.Split(branches_def, " | ")
        var n_branches = len(branch_defs)
        var branches_part_defs = make([][] string, n_branches)
        for i, branch_def := range branch_defs {
            branches_part_defs[i] = strings.Split(branch_def, " ")
        }
        var branches = make([] Branch, n_branches)
        for i, branch_part_defs := range branches_part_defs {
            var num_parts = len(branch_part_defs)
            branches[i].Parts = make([] Part, num_parts)
            if num_parts == 0 {
                panic(name + ": zero parts")
            }
            if num_parts > MAX_NUM_PARTS {
                panic(name + ": too many parts")
            }
            for j, part_def := range branch_part_defs {
                // check if valid
                if part_def == "" {
                    panic("redundant blank in definition of " + name_def)
                }
                // extract part name
                var required = strings.HasSuffix(part_def, "!")
                var part_name = strings.TrimRight(part_def, "!")
                part_name = EscapePartName(part_name)
                // add to list if it is a keyword
                var part_type = GetPartType(part_name)
                var id, exists = mapName2Id[part_name]
                if (!exists) {
                    if u := strings.Index(part_name, "*"); u != -1 {
                        var item = EscapePartName(part_name[:u])
                        var sep = EscapePartName(part_name[(u+1):])
                        id = generateListRule(part_name, item, sep, true)
                    } else if v := strings.Index(part_name, "+"); v != -1 {
                        var item = EscapePartName(part_name[:v])
                        var sep = EscapePartName(part_name[(v + 1):])
                        id = generateListRule(part_name, item, sep, false)
                    } else if strings.HasSuffix(part_name, "?") {
                        var item = EscapePartName(
                            strings.TrimSuffix(part_name, "?"),
                        )
                        id = generateOptionalRule(part_name, item)
                    } else {
                        panic("undefined part: " + part_name)
                    }
                }
                branches[i].Parts[j] = Part {
                    Id: id, Required: required, PartType: part_type,
                }
            }
        }
        rules[id] = Rule {
            Id: id,
            Nullable: nullable,
            Branches: branches,
        }
    }
}
func generateListRule(list string, item string, sep string, nullable bool) Id {
    var gen_part = func(name string, required bool) Part {
        var id, exists = mapName2Id[name]
        if !(exists) {
            panic("undefined part: " + name)
        }
        return Part {
            Id:       id,
            PartType: GetPartType(name),
            Required: required,
        }
    }
    var list_id = assignId2Name(list)
    var tail = (list + "_tail")
    var tail_id = assignId2Name(tail)
    rules[list_id] = Rule {
        Id: list_id,
        Nullable: nullable,
        Branches: [] Branch { { [] Part {
            gen_part(item, false),
            gen_part(tail, false),
        } } },
        Generated: true,
        GenInfo: RuleGenInfo { RuleGenList, item, sep },
    }
    rules[tail_id] = Rule {
        Id: tail_id,
        Nullable: true,
        Branches: [] Branch { { (func() ([] Part) {
            if sep != "" { return [] Part {
                gen_part(sep, false),
                gen_part(item, true),
                gen_part(tail, false),
            } } else { return [] Part {
                gen_part(item, false),
                gen_part(tail, false),
            } }
        })() } },
        Generated: true,
        GenInfo: RuleGenInfo { RuleGenListTail, item, sep },
    }
    return list_id
}
func generateOptionalRule(opt string, item string) Id {
    var item_id, exists = mapName2Id[item]
    if !(exists) {
        panic("undefined part: " + item)
    }
    var item_type = GetPartType(item)
    var opt_id = assignId2Name(opt)
    rules[opt_id] = Rule {
        Id: opt_id,
        Nullable: true,
        Branches: [] Branch { { [] Part {
            Part {
                Id:       item_id,
                PartType: item_type,
                Required: false,
            },
        } } },
        Generated: true,
        GenInfo: RuleGenInfo { RuleGenOptional, item, "" },
    }
    return opt_id
}

func init() {
    // we assume id starts from tokens
    assignId2Tokens()
    assignId2CondKeywords()
    assignId2Rules()
    parseRules()
}


