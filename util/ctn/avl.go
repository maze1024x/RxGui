package ctn


type AVL[T any] struct {
    Value   T
    Left    *AVL[T]
    Right   *AVL[T]
    Size    uint64
    Height  uint64
}
func AvlNode[T any] (v T, left *AVL[T], right *AVL[T]) *AVL[T] {
    var node = &AVL[T] {
        Value:  v,
        Left:   left,
        Right:  right,
        Size:   1 + left.GetSize() + right.GetSize(),
        Height: 1 + max(left.GetHeight(), right.GetHeight()),
    }
    return node.balanced()
}
func AvlLeaf[T any] (v T) *AVL[T] {
    return &AVL[T] {
        Value:  v,
        Left:   nil,
        Right:  nil,
        Size:   1,
        Height: 1,
    }
}

func (node *AVL[T]) IsLeaf() bool {
    return (node.Left == nil && node.Right == nil)
}
func (node *AVL[T]) GetSize() uint64 {
    if node == nil {
        return 0
    } else {
        return node.Size
    }
}
func (node *AVL[T]) GetHeight() uint64 {
    if node == nil {
        return 0
    } else {
        return node.Height
    }
}

func (node *AVL[T]) Walk(f func(T)) {
    if node == nil { return }
    node.Left.Walk(f)
    f(node.Value)
    node.Right.Walk(f)
}
func (node *AVL[T]) Lookup(target T, cmp Compare[T]) (T, bool) {
    if node == nil {
        return zero[T](), false
    } else {
        switch cmp(target, node.Value) {
        case Smaller:
            return node.Left.Lookup(target, cmp)
        case Bigger:
            return node.Right.Lookup(target, cmp)
        case Equal:
            return node.Value, true
        default:
            panic("impossible branch")
        }
    }
}
func (node *AVL[T]) Inserted(inserted T, cmp Compare[T]) (*AVL[T], bool) {
    if node == nil {
        return AvlLeaf(inserted), false
    } else {
        var value = node.Value
        var left = node.Left
        var right = node.Right
        switch cmp(inserted, value) {
        case Smaller:
            var left_inserted, override = left.Inserted(inserted, cmp)
            return AvlNode(value, left_inserted, right), override
        case Bigger:
            var right_inserted, override = right.Inserted(inserted, cmp)
            return AvlNode(value, left, right_inserted), override
        case Equal:
            return AvlNode(inserted, left, right), true
        default:
            panic("impossible branch")
        }
    }
}
func (node *AVL[T]) Deleted(target T, cmp Compare[T]) (T, *AVL[T], bool) {
    if node == nil {
        return zero[T](), nil, false
    } else {
        var value = node.Value
        var left = node.Left
        var right = node.Right
        switch cmp(target, value) {
        case Smaller:
            var deleted, rest, found = left.Deleted(target, cmp)
            if found {
                return deleted, AvlNode(value, rest, right), true
            } else {
                return zero[T](), node, false
            }
        case Bigger:
            var deleted, rest, found = right.Deleted(target, cmp)
            if found {
                return deleted, AvlNode(value, left, rest), true
            } else {
                return zero[T](), node, false
            }
        case Equal:
            if left == nil {
                return value, right, true
            } else if right == nil {
                return value, left, true
            } else {
                var node_state, _ = node.GetBalanceState()
                if node_state == RightTaller {
                    var successor, rest_right, found = right.DeleteMin()
                    assert(found, "right subtree should not be empty")
                    return value, AvlNode(successor, left, rest_right), true
                } else {
                    var prior, rest_left, found = left.DeletedMax()
                    assert(found, "left subtree should not be empty")
                    return value, AvlNode(prior, rest_left, right), true
                }
            }
        default:
            panic("impossible branch")
        }
    }
}

func (node *AVL[T]) DeleteMin() (T, *AVL[T], bool) {
    if node == nil {
        return zero[T](), nil, false
    } else {
        var value = node.Value
        var left = node.Left
        var right = node.Right
        var deleted, rest, found = left.DeleteMin()
        if found {
            return deleted, AvlNode(value, rest, right), true
        } else {
            return value, right, true
        }
    }
}
func (node *AVL[T]) DeletedMax() (T, *AVL[T], bool) {
    if node == nil {
        return zero[T](), nil, false
    } else {
        var value = node.Value
        var left = node.Left
        var right = node.Right
        var deleted, rest, found = right.DeletedMax()
        if found {
            return deleted, AvlNode(value, left, rest), true
        } else {
            return value, left, true
        }
    }
}

type BalanceState int
const (
    LeftTaller BalanceState =  iota
    RightTaller
    NeitherTaller
)
func (node *AVL[T]) GetBalanceState() (BalanceState, uint) {
    var L = node.Left.GetHeight()
    var R = node.Right.GetHeight()
    if L > R {
        return LeftTaller, uint(L - R)
    } else if L < R {
        return RightTaller, uint(R - L)
    } else {
        return NeitherTaller, 0
    }
}
func (node *AVL[T]) balanced() *AVL[T] {
    var current = node
    var current_state, diff = current.GetBalanceState()
    if current_state == NeitherTaller || diff == 1 {
        return current
    } else {
        assert(diff == 2, "invalid usage of balanced()")
        switch current_state {
        case LeftTaller:
            var left = current.Left
            var left_state, _ = left.GetBalanceState()
            switch left_state {
            case LeftTaller, NeitherTaller:
                var new_right = AvlNode(current.Value, left.Right, current.Right)
                var new_current = AvlNode(left.Value, left.Left, new_right)
                return new_current
            case RightTaller:
                var middle = left.Right
                var new_left = AvlNode(left.Value, left.Left, middle.Left)
                var new_right = AvlNode(current.Value, middle.Right, current.Right)
                var new_current = AvlNode(middle.Value, new_left, new_right)
                return new_current
            default:
                panic("impossible branch")
            }
        case RightTaller:
            var right = current.Right
            var right_state, _ = right.GetBalanceState()
            switch right_state {
            case LeftTaller:
                var middle = right.Left
                var new_left = AvlNode(current.Value, current.Left, middle.Left)
                var new_right = AvlNode(right.Value, middle.Right, right.Right)
                var new_current = AvlNode(middle.Value, new_left, new_right)
                return new_current
            case RightTaller, NeitherTaller:
                var new_left = AvlNode(current.Value, current.Left, right.Left)
                var new_current = AvlNode(right.Value, new_left, right.Right)
                return new_current
            default:
                panic("impossible branch")
            }
        default:
            panic("impossible branch")
        }
    }
}


