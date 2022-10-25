package compiler

import (
	"regexp"
	"strings"
	"strconv"
	"unicode/utf8"
	"rxgui/util"
	"rxgui/interpreter/lang/source"
	"rxgui/interpreter/lang/typsys"
	"rxgui/interpreter/lang/textual/ast"
	"rxgui/interpreter/program"
)


func checkInt(I ast.Int, cc *exprCheckContext) (*program.Expr, *source.Error) {
	if typsys.Equal(cc.expected, program.T_Float()) {
		return cc.forwardTo(ast.WrapTermAsExpr(ast.VariousTerm {
			Node: I.Node,
			Term: ast.Float {
				Node:  I.Node,
				Value: I.Value,
			},
		}))
	}
	var loc = I.Location
	var value, ok = util.WellBehavedParseInteger(source.CodeToChars(I.Value))
	if !(ok) {
		panic("something went wrong")
	}
	var int_t = typsys.CertainType { Type: program.T_Int() }
	return cc.assign(int_t, loc,
		program.IntLiteral {
			Value: value,
		})
}

func checkFloat(F ast.Float, cc *exprCheckContext) (*program.Expr, *source.Error) {
	var loc = F.Location
	var value, ok = util.ParseDouble(source.CodeToChars(F.Value))
	if !(ok) {
		return cc.error(loc, E_InvalidFloat {})
	}
	if !(util.IsNormalFloat(value)) {
		panic("something went wrong")
	}
	var float_t = typsys.CertainType { Type: program.T_Float() }
	return cc.assign(float_t, loc,
		program.FloatLiteral {
			Value: value,
		})
}

func checkChar(C ast.Char, cc *exprCheckContext) (*program.Expr, *source.Error) {
	var loc = C.Location
	var value, err = parseCharLiteral(C.Value)
	if err != nil { return cc.error(loc, err) }
	var char_t = typsys.CertainType { Type: program.T_Char() }
	return cc.assign(char_t, loc,
		program.CharLiteral {
			Value: value,
		})
}

func checkBytes(B ast.Bytes, cc *exprCheckContext) (*program.Expr, *source.Error) {
	var loc = B.Location
	var value = make([] byte, len(B.Bytes))
	for i, b := range B.Bytes {
		var x = hexDigitCharToByte(b.Value[2])
		var y = hexDigitCharToByte(b.Value[3])
		value[i] = ((x << 4) | y)
	}
	var bytes_t = typsys.CertainType { Type: program.T_Bytes() }
	return cc.assign(bytes_t, loc,
		program.BytesLiteral {
			Value: value,
		})
}
func hexDigitCharToByte(d uint16) byte {
	if d >= 'a' {
		return byte(10 + (d - 'a'))
	} else if d >= 'A' {
		return byte(10 + (d - 'A'))
	} else {
		return byte(d - '0')
	}
}

func checkString(S ast.String, cc *exprCheckContext) (*program.Expr, *source.Error) {
	var buf strings.Builder
	var first = S.First
	var s, err = normalizeText(first.Value)
	if err != nil { return cc.error(first.Location, err) }
	buf.WriteString(s)
	for _, part := range S.Parts {
		switch P := part.Content.StringPartContent.(type) {
		case ast.Text:
			var s, err = normalizeText(P.Value)
			if err != nil { return cc.error(P.Location, err) }
			buf.WriteString(s)
		case ast.Char:
			var r, err = parseCharLiteral(P.Value)
			if err != nil { return cc.error(P.Location, err) }
			buf.WriteRune(r)
		default:
			panic("impossible branch")
		}
	}
	var value = buf.String()
	var loc = S.Location
	if typsys.Equal(cc.expected, program.T_RegExp()) {
		var value, err = compileRegexpLiteral(value)
		if err != nil { return cc.error(loc, err) }
		var regexp_t = typsys.CertainType { Type: program.T_RegExp() }
		return cc.assign(regexp_t, loc,
			program.RegexpLiteral {
				Value: value,
			})
	}
	var string_t = typsys.CertainType { Type: program.T_String() }
	return cc.assign(string_t, loc,
		program.StringLiteral {
			Value: value,
		})
}

func normalizeText(code source.Code) (string, source.ErrorContent) {
	var L = len(code)
	if code[0] == '\'' && code[L-1] == '\'' {
		return source.CodeToString(code[1:(L-1)]), nil
	} else if code[0] == '"' && code[L-1] == '"' {
		var u, err = strconv.Unquote(source.CodeToString(code))
		if err != nil {
			return "", E_InvalidText {}
		}
		return u, nil
	} else {
		panic("impossible branch")
	}
}
func parseCharLiteral(code source.Code) (rune, source.ErrorContent) {
	var value, ok = util.ParseRune(source.CodeToChars(code))
	if ok {
		return value, nil
	} else {
		return utf8.RuneError, E_InvalidChar {}
	}
}
func compileRegexpLiteral(str string) (*regexp.Regexp, source.ErrorContent) {
	var value, err = regexp.Compile(str)
	if err == nil {
		return value, nil
	} else {
		return nil, E_InvalidRegexp {
			Detail: err.Error(),
		}
	}
}


