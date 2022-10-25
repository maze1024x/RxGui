package syntax

import (
    "regexp"
    "strings"
)
type Regexp = *regexp.Regexp
func r (pattern string) Regexp { return regexp.MustCompile(`^` + pattern) }


const defaultRoot = "root"
const replRoot = "repl_cmd"
func DefaultRoot() Id { return Name2IdMustExist(defaultRoot) }
func ReplRoot() Id { return Name2IdMustExist(replRoot) }

var __EscapeMap = map[string] string {
    "_bar": "|",
    "_at": "@",
}
var __IgnoreTokens = [...] string {
    "Comment",
    "Blank",
    "LF",
}

const LF = `\n`
const Blanks = ` \t\r　`
const Symbols = `\{\}\[\]\(\)\.,:;@#\|\&\\'"` + "`"
const nameEverywhereDisallow = (Symbols + Blanks + LF)
const nameHeadDisallow = (`0-9` + nameEverywhereDisallow)
const NamePattern = `[^`+nameHeadDisallow+`][^`+nameEverywhereDisallow+`]*`

var __Tokens = [...] Token {
    // doc
    Token { Name: "Doc",      Pattern: r(`///[^`+LF+`]*`) },
    // ignored
    Token { Name: "Comment",  Pattern: r(`//[^`+LF+`]*`) },
    Token { Name: "Blank",    Pattern: r(`[`+Blanks+`]+`) },
    Token { Name: "LF",       Pattern: r(`[`+LF+`]+`) },
    // literals
    Token { Name: "Text",   Pattern: r(`'[^']*'`) },
    Token { Name: "Text",   Pattern: r(`"(\\"|[^"])*"`) },
    Token { Name: "Int",    Pattern: r(`\-?0[xX][0-9A-Fa-f]+`) },
    Token { Name: "Int",    Pattern: r(`\-?0[oO][0-7]+`) },
    Token { Name: "Int",    Pattern: r(`\-?0[bB][01]+`) },
    Token { Name: "Float",  Pattern: r(`\-?\d+(\.\d+)?[Ee][\+\-]?\d+`) },
    Token { Name: "Float",  Pattern: r(`\-?\d+\.\d+`) },
    Token { Name: "Int",    Pattern: r(`\-?\d+`) },
    Token { Name: "Byte",   Pattern: r(`\\x[0-9A-Fa-f][0-9A-Fa-f]`) },
    Token { Name: "Char",   Pattern: r("`.`") },
    Token { Name: "Char",   Pattern: r(`\\u[0-9A-Fa-f]+`) },
    Token { Name: "Char",   Pattern: r(`\\[a-z]`) },
    // symbols
    Token { Name: "(",    Pattern: r(`\(`) },
    Token { Name: ")",    Pattern: r(`\)`) },
    Token { Name: "[",    Pattern: r(`\[`) },
    Token { Name: "]",    Pattern: r(`\]`) },
    Token { Name: "{",    Pattern: r(`\{`) },
    Token { Name: "}",    Pattern: r(`\}`) },
    Token { Name: "...",  Pattern: r(`\.\.\.`) },
    Token { Name: "..",   Pattern: r(`\.\.`) },
    Token { Name: ".",    Pattern: r(`\.`) },
    Token { Name: ",",    Pattern: r(`,`) },
    Token { Name: "::",   Pattern: r(`::`) },
    Token { Name: ":",    Pattern: r(`:`) },
    Token { Name: ";",    Pattern: r(`;`) },
    Token { Name: "@",    Pattern: r(`@`) },
    Token { Name: "#",    Pattern: r(`#`) },
    Token { Name: "|",    Pattern: r(`\|`) },
    Token { Name: "&",    Pattern: r(`\&`) },
    // keywords
    Token { Name: "Const",  Pattern: r(`const`),  Keyword: true },
    Token { Name: "If",     Pattern: r(`if`),     Keyword: true },
    Token { Name: "Else",   Pattern: r(`else`),   Keyword: true },
    Token { Name: "When",   Pattern: r(`when`),   Keyword: true },
    Token { Name: "Each",   Pattern: r(`each`),   Keyword: true },
    Token { Name: "Let",    Pattern: r(`let`),    Keyword: true },
    Token { Name: "New",    Pattern: r(`new`),    Keyword: true },
    Token { Name: "=>",     Pattern: r(`=>`),     Keyword: true },
    Token { Name: "=",      Pattern: r(`=`),      Keyword: true },
    // identifier
    Token { Name: "Name",  Pattern: r(NamePattern) },
}
var nameTokenRegexp = regexp.MustCompile(NamePattern)
func Tokens() ([] Token) { return __Tokens[:] }
func IgnoreTokens() ([] string) { return __IgnoreTokens[:] }
func NameTokenId() Id { return Name2IdMustExist("Name") }
func NameTokenRegexp() Regexp { return nameTokenRegexp }

var __ConditionalKeywords = [...] string {
    "@namespace", "@using",
    "@run",
    "@entry", "@type", "@function", "@operator", "@method",
    "@native", "@record", "@interface", "@union", "@enum",
    "@variadic",
}
func GetKeywordList() ([] string) {
    var list = make([] string, 0)
    for _, v := range __ConditionalKeywords {
        var kw = strings.TrimPrefix(v, "@")
        list = append(list, kw)
    }
    for _, t := range __Tokens {
        if t.Keyword {
            var kw = strings.TrimPrefix(t.Pattern.String(), "^")
            list = append(list, kw)
        }
    }
    return list
}

var __SyntaxDefinition = [...] string {
    "root = ns! alias* stmt*",
      "ns = @namespace name? ::! ",
        "name = Name",
      "alias = off @using alias_name alias_target",
        "off? = #",
        "alias_name? = name = ",
        "alias_target = alias_to_ns | alias_to_ref_base",
          "alias_to_ns = @namespace name!",
          "alias_to_ref_base = ref_base!",
      "stmt = decl_entry | decl_type | decl_func | decl_const | decl_method",
    "decl_entry = docs off @entry block!",
      "docs? = doc+",
        "doc = Doc",
    "decl_type = docs off @type name! type_params impl type_def!",
      "type_params? = [ name*, ]!",
      "impl? = ( ref_base*, )!",
        "ref_base = ns_prefix name",
          "ns_prefix? = name :: ",
      "type_def = native_type_def | record | interface | union | enum",
        "native_type_def = @native",
        "record = @record record_def!",
          "record_def = { field*, }!",
            "field = docs name type! field_default",
            "field_default? = ( expr! )!",
        "interface = @interface {! method*, }!",
          "method = docs name type!",
        "union = @union {! type+,! }!",
        "enum = @enum {! enum_item+,! }!",
          "enum_item = docs name",
    "decl_func = docs off function variadic name! sig! body!",
      "function = @function | @operator",
      "variadic? = @variadic",
      "sig = type_params inputs! implicit output!",
        "inputs = record_def",
        "implicit? = inputs",
        "output = type",
      "body = native_body | block!",
        "native_body = @native (! text! )!",
    "decl_const = docs off Const name! type! body!",
    "decl_method = docs off @method receiver! .! name! type! body!",
      "receiver = name",
    "type = ref",
      "ref = ref_base type_args",
        "type_args? = [ type*, ]!",
    "expr = cast* term pipe*",
      "cast = ( [ type! ]! )!",
      "pipe = pipe_call | pipe_infix | pipe_cast | pipe_get | pipe_interior",
        "pipe_call = call_ordered | call_unordered",
          "call_ordered = ( expr*, )!",
          "call_unordered = { arg_mapping*, }!",
            "arg_mapping = name arg_mapping_to",
              "arg_mapping_to? = : expr!",
        "pipe_infix = off _bar ref! pipe_call!",
        "pipe_cast = . cast",
        "pipe_get = . name",
        "pipe_interior = . ( ref_base! )!",
    "term = infix_term | lambda | if | when | each | block | ref_term " +
        "| int | float | char | bytes | string ",
      "infix_term = ( infix_left operator infix_right! )!",
        "infix_left = expr",
        "operator = ref",
        "infix_right = expr",
      "lambda = { pattern? lambda_self => expr! }!",
        "lambda_self? = & name!",
        "pattern = pattern_single | pattern_multiple",
          "pattern_single = name",
          "pattern_multiple = ( name+, )",
      "if = If (! cond+,! )! if_yes elif* Else! if_no",
        "cond = cond_pattern expr!",
          "cond_pattern? = Let pattern! =! ",
        "if_yes = block!",
        "if_no = block!",
        "elif = If (! cond+,! )! block!",
      "when = When (! expr! )! {! case+,! }!",
        "case = off name+_bar pattern? =>! expr!",
      "each = Each (! type! )! {! case+,! }!",
      "block = { binding* expr! }!",
        "binding = binding_plain | binding_cps",
          "binding_plain = off let pattern! =! expr! ,!",
            "let = Let | Const",
          "binding_cps = off _at ref! cps_pattern expr! ,!",
            "cps_pattern? = pattern = ",
      "ref_term = new ref",
        "new? = New new_tag",
          "new_tag? = : name!",
      "int = Int",
      "float = Float",
      "char = Char",
      "bytes = byte+",
        "byte = Byte",
      "string = text string_part* ",
        "text = Text",
        "string_part = .. string_part_content!",
          "string_part_content = text | char",
    "repl_cmd = repl_assign | repl_run | repl_eval",
      "repl_assign = name = expr!",
      "repl_run = : @run expr!",
      "repl_eval = expr!",
}


