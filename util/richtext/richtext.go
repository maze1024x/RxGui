package richtext


type Block struct {
    ClassList  [] string
    Lines      [] Line
    Children   [] Block
}
func (b *Block) AddClass(classes ...string) {
    b.ClassList = append(b.ClassList, classes...)
}
func (b *Block) WriteLine(content string, tags ...string) {
    b.Lines = append(b.Lines, Line { [] Span { {
        Content: content,
        Tags:    tags,
    } } })
}
func (b *Block) WriteSpan(content string, tags ...string) {
    b.writeSpan(content, false, tags)
}
func (b *Block) WriteRawSpan(content string, tags ...string) {
    b.writeSpan(content, true, tags)
}
func (b *Block) writeSpan(content string, raw bool, tags ([] string)) {
    var span = Span {
        Content: content,
        Tags:    tags,
        Raw:     raw,
    }
    if len(b.Lines) > 0 {
        var last = &(b.Lines[len(b.Lines) - 1])
        last.Spans = append(last.Spans, span)
    } else {
        b.Lines = append(b.Lines, Line { Spans: [] Span { span } })
    }
}
func (b *Block) WriteLineFeed() {
    b.Lines = append(b.Lines, Line { Spans: [] Span {} })
}
func (b *Block) Append(another Block) {
    b.Children = append(b.Children, another)
}
func (b Block) WithoutLeadingSpan() Block {
    if len(b.Lines) == 0 {
        return b
    } else {
        if len(b.Lines[0].Spans) == 0 {
            return b
        } else {
            var lines = make([] Line, len(b.Lines))
            copy(lines, b.Lines)
            lines[0].Spans = lines[0].Spans[1:]
            return Block {
                ClassList: b.ClassList,
                Lines:     lines,
                Children:  b.Children,
            }
        }
    }
}
func (b Block) WithLeadingSpan(content string, tags ...string) Block {
    if len(b.Lines) == 0 {
        return b.WithLeadingLine(content, tags...)
    } else {
        var s = Span {
            Content: content,
            Tags:    tags,
        }
        var lines = make([] Line, len(b.Lines))
        copy(lines, b.Lines)
        lines[0].Spans = append([] Span { s }, lines[0].Spans...)
        return Block {
            ClassList: b.ClassList,
            Lines:     lines,
            Children:  b.Children,
        }
    }
}
func (b Block) WithLeadingLine(content string, tags ...string) Block {
    var l = Line { [] Span { {
        Content: content,
        Tags:    tags,
    } } }
    return Block {
        ClassList: b.ClassList,
        Lines:     append([] Line { l }, b.Lines...),
        Children:  b.Children,
    }
}

type Line struct {
    Spans  [] Span
}

type Span struct {
    Content  string
    Tags     [] string
    Raw      bool
    Link     MaybeLink
}
const (
    TAG_NORMAL       = "null"
    TAG_B            = "b"
    TAG_I            = "i"
    TAG_DEL          = "del"
    TAG_DBG_TYPE     = "dbg-type"
    TAG_DBG_FIELD    = "dbg-field"
    TAG_DBG_STRING   = "dbg-string"
    TAG_DBG_NUMBER   = "dbg-number"
    TAG_DBG_CONSTANT = "dbg-constant"
    TAG_HIGHLIGHT    = "highlight"
    TAG_ERR          = "error"
    TAG_ERR_NOTE     = "note"
    TAG_ERR_INLINE   = "inline"
    TAG_INPUT        = "input"
    TAG_OUTPUT       = "output"
    TAG_SUCCESS      = "success"
    TAG_FAILURE      = "failure"
)

type MaybeLink interface { Maybe(Link, MaybeLink) }
func (Link) Maybe(Link, MaybeLink) {}
type Link struct {
    Page    string
    Anchor  string
}


