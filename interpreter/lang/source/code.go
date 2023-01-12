package source

import (
    "io"
    "strings"
    "unicode"
    "unicode/utf8"
    "unicode/utf16"
)


type Code ([] uint16)

func CodeReader(code Code, pos int) io.RuneReader {
    return &codeReader { src: code, pos: pos }
}
type codeReader struct {
    src Code
    pos int
}
func (r *codeReader) ReadRune() (rune, int, error) {
    if r.pos >= len(r.src) {
        return -1, 0, io.EOF
    }
    var a = r.src[r.pos]
    if utf16.IsSurrogate(rune(a)) {
        if (r.pos + 1) < len(r.src) {
            var b = r.src[r.pos + 1]
            var next = utf16.DecodeRune(rune(a), rune(b))
            r.pos += 2
            return next, 2, nil
        } else {
            var next = unicode.ReplacementChar
            r.pos += 1
            return next, 1, nil
        }
    } else {
        var next = rune(a)
        r.pos += 1
        return next, 1, nil
    }
}

type codeBuilder struct {
    current Code
}
func (buf *codeBuilder) Append(char rune) {
    var code = &(buf.current)
    if 0 <= char && char < unicode.MaxRune {
        if char < 0x10000 {
            *code = append(*code, uint16(char))
        } else {
            var a, b = utf16.EncodeRune(char)
            *code = append(*code, uint16(a), uint16(b))
        }
    } else {
        *code = append(*code, unicode.ReplacementChar)
    }
}
func (buf *codeBuilder) Collect() Code {
    return buf.current
}

func CodeToChars(code Code) ([] rune) {
    return utf16.Decode(code)
}
func CodeToString(code Code) string {
    var buf strings.Builder
    var reader = codeReader { src: code, pos: 0 }
    for { if char, _, err := reader.ReadRune(); err != io.EOF {
        buf.WriteRune(char)
    } else {
        break
    }}
    return buf.String()
}
func StringToCode(s string) Code {
    var buf codeBuilder
    for _, char := range s {
        buf.Append(char)
    }
    return buf.Collect()
}
func DecodeUtf8ToCode(b ([] byte)) Code {
    var buf codeBuilder
    for len(b) > 0 {
        var char, size = utf8.DecodeRune(b)
        buf.Append(char)
        b = b[size:]
    }
    return buf.Collect()
}

const placeholder = "<?>"
func Placeholder() Code {
    var snippet = make(Code, len(placeholder))
    for i := range placeholder {
        snippet[i] = uint16(rune(placeholder[i]))
    }
    return snippet
}
func IsPlaceholder(snippet Code) bool {
    if len(snippet) == len(placeholder) {
        for i := range placeholder {
            if rune(snippet[i]) != rune(placeholder[i]) {
                return false
            }
        }
        return true
    } else {
        return false
    }
}


