package richtext

import (
    "fmt"
    "html"
    "strings"
)


type renderer interface {
    WrapSpanContent(content string, tags ([] string)) string
}
type plainTextRenderer struct {}
type htmlRenderer struct {
    CssPalette  CssPalette
}
type consoleRenderer struct {
    AnsiPalette  AnsiPalette
}
func (opts plainTextRenderer) WrapSpanContent(content string, _ ([] string)) string {
    return content
}
func (opts htmlRenderer) WrapSpanContent(content string, tags ([] string)) string {
    var palette = opts.CssPalette
    var style = htmlStyleAttr(palette(tags))
    var element = htmlElement("span", style, content)
    return element
}
func (opts consoleRenderer) WrapSpanContent(content string, tags ([] string)) string {
    var palette = opts.AnsiPalette
    var buf strings.Builder
    for _, tag := range tags {
        buf.WriteString(palette(tag))
    }
    buf.WriteString(content)
    buf.WriteString(reset)
    return buf.String()
}

func (b Block) RenderPlainText() string {
    return render(b, plainTextRenderer {})
}
func (b Block) RenderHtml() string {
    return fmt.Sprintf("<pre>%s</pre>", render(b, htmlRenderer {
        CssPalette: defaultCssPalette,
    }))
}
func (b Block) RenderConsole() string {
    return render(b, consoleRenderer {
        AnsiPalette: defaultAnsiPalette,
    })
}

func render(b Block, r renderer) string {
    if len(b.Lines) == 0 && len(b.Children) == 1 {
        return render(b.Children[0], r)
    }
    return getBlockText(b, false, r)
}
func getBlockText(b Block, indent bool, r renderer) string {
    return withIndention(indent, func() string {
        var buf strings.Builder
        for _, l := range b.Lines {
            for i, span := range l.Spans {
                buf.WriteString(r.WrapSpanContent(span.Content, span.Tags))
                if spanNeedTrailingSpace(i, l.Spans) {
                    buf.WriteRune(' ')
                }
            }
            if len(l.Spans) > 0 {
                buf.WriteRune('\n')
            }
        }
        var indent_child = (len(b.Lines) > 0)
        for _, child := range b.Children {
            buf.WriteString(getBlockText(child, indent_child, r))
        }
        return buf.String()
    })
}
func spanNeedTrailingSpace(i int, spans ([] Span)) bool {
    if (i + 1) < len(spans) {
        var current = spans[i]
        var next = spans[i + 1]
        if !(current.Raw) &&
        !(strings.HasSuffix(current.Content, " ")) &&
        !(strings.HasPrefix(next.Content, " ")) &&
        next.Content != "," &&
        next.Content != ":" {
            return true
        }
    }
    return false
}
func withIndention(indent bool, k func()(string)) string {
    if !(indent) {
        return k()
    }
    const indention = "    "
    var buf strings.Builder
    var lines = strings.Split(strings.TrimSuffix(k(), "\n"), "\n")
    for i, line := range lines {
        if i > 0 {
            buf.WriteRune('\n')
        }
        buf.WriteString(indention)
        buf.WriteString(line)
    }
    buf.WriteRune('\n')
    return buf.String()
}

func htmlStyleAttr(style ([] func()(string,string))) string {
    if len(style) == 0 {
        return ""
    }
    var buf strings.Builder
    var occurred = make(map[string] struct{})
    buf.WriteString("style=\"")
    for _, pair := range style {
        var key, value = pair()
        if key != "" {
            var _, duplicate = occurred[key]
            if !(duplicate) {
                occurred[key] = struct{}{}
                buf.WriteString(fmt.Sprintf("%s:%s;", key, value))
            }
        }
    }
    buf.WriteString("\"")
    return buf.String()
}
func htmlElement(element string, attrs string, content string) string {
    var buf strings.Builder
    buf.WriteString("<")
    buf.WriteString(element)
    if attrs != "" {
        buf.WriteString(" ")
        buf.WriteString(attrs)
    }
    buf.WriteString(">")
    buf.WriteString(html.EscapeString(content))
    buf.WriteString("</")
    buf.WriteString(element)
    buf.WriteString(">")
    return buf.String()
}


