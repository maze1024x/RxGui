package atom

import (
	"sort"
	"strings"
	"rxgui/lang/textual/ast"
	"rxgui/lang/textual/syntax"
	"rxgui/interpreter/compiler"
)


var IdentifierRegexp = syntax.NameTokenRegexp()
var Keywords = syntax.GetKeywordList()

type AutoCompleteRequest struct {
	PrecedingText  string     `json:"precedingText"`
	LocalBindings  [] string  `json:"localBindings"`
	CurrentPath    string     `json:"currentPath"`
}

type AutoCompleteResponse struct {
	Suggestions  [] AutoCompleteSuggestion   `json:"suggestions"`
}

type AutoCompleteSuggestion struct {
	Text     string   `json:"text"`
	Replace  string   `json:"replacementPrefix"`
	Type     string   `json:"type"`
	Display  string   `json:"displayText,omitempty"`
}

func AutoComplete(req AutoCompleteRequest, ctx LangServerContext) AutoCompleteResponse {
	const double_colon = "::"
	var get_search_text = func() (string, string) {
		var raw_text = req.PrecedingText
		var ranges = IdentifierRegexp.FindAllStringIndex(raw_text, -1)
		if len(ranges) == 0 {
			return "", ""
		}
		var last = ranges[len(ranges)-1]
		var lo = last[0]
		var hi = last[1]
		if hi != len(raw_text) {
			if strings.HasSuffix(raw_text, double_colon) &&
				hi == (len(raw_text) - len(double_colon)) {
				return "", raw_text[lo:hi]
			} else {
				return "", ""
			}
		}
		var text = raw_text[lo:hi]
		var ldc = len(double_colon)
		if len(ranges) >= ldc {
			var preceding = ranges[len(ranges)-ldc]
			var p_lo = preceding[0]
			var p_hi = preceding[1]
			if (p_hi + ldc) == lo && raw_text[p_hi:lo] == double_colon {
				var text_mod = raw_text[p_lo:p_hi]
				return text, text_mod
			} else {
				return text, ""
			}
		} else {
			return text, ""
		}
	}
	var input, input_ns = get_search_text()
	if input == "" && input_ns == "" {
		return AutoCompleteResponse {}
	}
	var quick_check = func(id ast.Identifier) bool {
		if !(len(id.Name) > 0) { panic("something went wrong") }
		if len(input) > 0 {
			var first_char = rune(id.Name[0])
			if first_char < 128 &&
				first_char != rune(input[0]) &&
				(first_char + ('a' - 'A')) != rune(input[0]) {
				return false
			}
		}
		return true
	}
	var suggestions = make([] AutoCompleteSuggestion, 0)
	if len(input) > 0 && input_ns == "" {
		for _, binding := range req.LocalBindings {
			if strings.HasPrefix(binding, input) {
				suggestions = append(suggestions, AutoCompleteSuggestion {
					Text:    binding,
					Replace: input,
					Type:    "variable",
				})
			}
		}
	}
	var suggested_function_names = make(map[string] bool)
	var process_statement func(ast.Statement, string)
	process_statement = func(stmt ast.Statement, ns string) {
		var name, type_, ok = (func() (string, string, bool) {
			switch s := stmt.(type) {
			case ast.DeclType:
				if !(quick_check(s.Name)) { return "", "", false }
				return ast.Id2String(s.Name), "type", true
			case ast.DeclConst:
				if !(quick_check(s.Name)) { return "", "", false }
				return ast.Id2String(s.Name), "constant", true
			case ast.DeclFunction:
				if !(quick_check(s.Name)) { return "", "", false }
				return ast.Id2String(s.Name), "function", true
			default:
				return "", "", false
			}
		})()
		if ok {
			var name_lower = strings.ToLower(name)
			if strings.HasPrefix(name, input) || strings.HasPrefix(name_lower, input) {
				if type_ == "function" {
					if suggested_function_names[name] {
						goto skip
					} else {
						suggested_function_names[name] = true
					}
				}
				suggestions = append(suggestions, AutoCompleteSuggestion {
					Text:    name,
					Replace: input,
					Type:    type_,
				})
			}
			skip:
		}
	}
	if (input_ns == "" && len(input) >= 2) || input_ns != "" {
		var fs = compiler.RealFileSystem {}
		var script_ast, err = compiler.Load(req.CurrentPath, fs)
		if err != nil { goto keywords }
		var script_ns = ast.MaybeId2String(script_ast.Namespace)
		for _, item := range script_ast.Statements {
			process_statement(item.Statement, script_ns)
		}
		for _, builtin_ast := range ctx.BuiltinAstData {
			var ns = ast.MaybeId2String(builtin_ast.Namespace)
			for _, item := range builtin_ast.Statements {
				process_statement(item.Statement, ns)
			}
		}
		if input_ns == "" {
			if strings.HasPrefix(script_ns, input) {
				suggestions = append(suggestions, AutoCompleteSuggestion {
					Text:    script_ns,
					Replace: input,
					Type:    "import",
				})
			}
		}
	}
	keywords:
	if len(input) > 0 && input_ns == "" {
		for _, kw := range Keywords {
			if len(kw) > 1 && ('a' <= kw[0] && kw[0] <= 'z') &&
				strings.HasPrefix(kw, input) {
				suggestions = append(suggestions, AutoCompleteSuggestion {
					Text:    kw,
					Replace: input,
					Type:    "keyword",
				})
			}
		}
	}
	sort.SliceStable(suggestions, func(i, j int) bool {
		var a = suggestions[i]
		var b = suggestions[j]
		if a.Type == b.Type {
			return a.Text < b.Text
		} else if a.Type == "keyword" {
			return false
		} else if b.Type == "keyword" {
			return true
		} else if a.Type == "function" {
			return false
		} else if b.Type == "function" {
			return true
		} else {
			return a.Text < b.Text
		}
	})
	return AutoCompleteResponse { Suggestions: suggestions }
}

