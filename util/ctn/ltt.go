package ctn


type LTT[T any] struct {
    Value  T
    Left   *LTT[T]
    Right  *LTT[T]
    Dist   uint64
    Size   uint64
}
func LttNode[T any] (v T, left *LTT[T], right *LTT[T]) *LTT[T] {
    var ld = left.GetDist()
    var rd = right.GetDist()
    if !(ld >= rd) { panic("violation of leftist property") }
    return &LTT[T] {
        Value: v,
        Left:  left,
        Right: right,
        Dist:  (1 + rd),
        Size:  (1 + left.GetSize() + right.GetSize()),
    }
}
func LttLeaf[T any] (v T) *LTT[T] {
    return &LTT[T] {
        Value: v,
        Left:  nil,
        Right: nil,
        Dist:  1,
        Size:  1,
    }
}

func (node *LTT[T]) GetDist() uint64 {
    if node != nil {
        return node.Dist
    } else {
        return 0
    }
}
func (node *LTT[T]) GetSize() uint64 {
    if node != nil {
        return node.Size
    } else {
        return 0
    }
}
func (node *LTT[T]) Merge(another *LTT[T], lt Less[T]) *LTT[T] {
    if node == nil { return another }
    if another == nil { return node }
    var smaller *LTT[T]
    var bigger *LTT[T]
    if !(lt(another.Value, node.Value)) {
        smaller = node
        bigger = another
    } else {
        bigger = node
        smaller = another
    }
    var v = smaller.Value
    var a = smaller.Left
    var b = smaller.Right.Merge(bigger, lt)
    if a.GetDist() >= b.GetDist() {
        return LttNode(v, a, b)
    } else {
        return LttNode(v, b, a)
    }
}

func (node *LTT[T]) Top() (T, bool) {
    if node != nil {
        return node.Value, true
    } else {
        return zero[T](), false
    }
}

func (node *LTT[T]) Popped(lt Less[T]) (T, *LTT[T], bool) {
    if node != nil {
        var rest = node.Left.Merge(node.Right, lt)
        return node.Value, rest, true
    } else {
        return zero[T](), nil, false
    }
}

func (node *LTT[T]) Pushed(v T, lt Less[T]) *LTT[T] {
    return node.Merge(LttLeaf(v), lt)
}


