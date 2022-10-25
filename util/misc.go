package util

import (
	"strings"
	"unicode/utf8"
	"encoding/json"
)

func UnmarshalJsonAllowEmpty(content string, value_ptr interface{}) error {
	if content != "" {
		return json.Unmarshal(([] byte)(content), value_ptr)
	} else {
		return nil
	}
}

// ensures invalid bytes are decoded into utf8.RuneError
func WellBehavedDecodeUtf8(b ([] byte)) string {
	var buf strings.Builder
	for len(b) > 0 {
		var char, size = utf8.DecodeRune(b)
		buf.WriteRune(char)
		b = b[size:]
	}
	return buf.String()
}
// can fail
func WellBehavedTryDecodeUtf8(b ([] byte)) (string,bool) {
	var buf strings.Builder
	for len(b) > 0 {
		var char, size = utf8.DecodeRune(b)
		if ((char == utf8.RuneError) && (size == 1)) {
			return "", false
		}
		buf.WriteRune(char)
		b = b[size:]
	}
	return buf.String(), true
}


