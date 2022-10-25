package source

import "fmt"


type Span struct {
    Start  int
    End    int
}
func (span Span) Merged(another Span) Span {
    if span == (Span{}) {
        return another
    }
    if another == (Span{}) {
        return span
    }
    var merged = Span { Start: span.Start, End: another.End }
    if !(merged.Start <= merged.End) {
        panic(fmt.Sprintf("invalid span merge: %+v and %+v", span, another))
    }
    return merged
}
func (span Span) Contains(pos int) bool {
    return (span.Start <= pos && pos < span.End)
}
func (span Span) Size() int {
    return (span.End - span.Start)
}


