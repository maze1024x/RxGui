package cst

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
	"rxgui/standalone/util/richtext"
	"rxgui/lang/source"
	"rxgui/lang/textual/scanner"
)


func FormatError (
	file   string,
	code   source.Code,
	info   scanner.RowColInfo,
	spans  scanner.RowSpanMap,
	point  scanner.Point,
	spot   scanner.Span,
	fov    uint,
	desc   richtext.Block,
) richtext.Block {
	var nh_rows = make([]scanner.Span, 0)
	var i = point.Row
	var j = point.Row
	for i > 1 && uint(point.Row - i) < (fov/2) {
		i -= 1
	}
	var last_row int
	if len(code) > 0 {
		last_row = info[len(code)-1].Row
	} else {
		last_row = 1
	}
	for j < last_row && uint(j - point.Row) < (fov/2) {
		j += 1
	}
	var start_row = i
	var end_row = j
	for r := start_row; r <= end_row; r += 1 {
		if r >= len(spans) { break }
		nh_rows = append(nh_rows, spans[r])
	}
	var expected_width = len(strconv.Itoa(end_row))
	var align = func(num int) string {
		var num_str = strconv.Itoa(num)
		var num_width = len(num_str)
		var buf strings.Builder
		buf.WriteString(num_str)
		for i := num_width; i < expected_width; i += 1 {
			buf.WriteRune(' ')
		}
		return buf.String()
	}
	var src richtext.Block
	var src_title = fmt.Sprintf (
		"----- (row %d, column %d) %s",
		point.Row, point.Col, file,
	)
	src.WriteLine(src_title, richtext.TAG_B)
	src.WriteLineFeed()
	var style = richtext.TAG_NORMAL
	for i, row := range nh_rows {
		var current_row = (start_row + i)
		src.WriteSpan(fmt.Sprintf("  %s | ", align(current_row)), style)
		var buf strings.Builder
		var commit = func() {
			var content = buf.String()
			buf.Reset()
			if len(content) > 0 {
				src.WriteRawSpan(content, style)
			}
		}
		for j, char := range code[row.Start: row.End] {
			var pos = (row.Start + j)
			if pos == spot.Start {
				commit()
				style = richtext.TAG_HIGHLIGHT
			}
			if pos == spot.End {
				commit()
				style = richtext.TAG_NORMAL
			}
			if utf16.IsSurrogate(rune(char)) {
				buf.WriteRune(unicode.ReplacementChar)
			} else {
				buf.WriteRune(rune(char))
			}
		}
		if row.End == spot.Start {
			commit()
			style = richtext.TAG_HIGHLIGHT
		}
		if row.End == spot.End {
			commit()
			style = richtext.TAG_NORMAL
		}
		commit()
		src.WriteLineFeed()
	}
	var msg richtext.Block
	msg.Append(src)
	msg.Append(desc)
	return msg
}


