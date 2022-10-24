package ctn

import "rxgui/standalone/util/backports/constraints"


type Compare[T any] (
	func(left T, right T)(Ordering))

type Less[T any] (
	func(left T, right T)(bool))

func StringCompare(a string, b string) Ordering {
	if a == b {
		return Equal
	} else if a < b {
		return Smaller
	} else {
		return Bigger
	}
}
func DefaultCompare[T constraints.Ordered] (left T, right T) Ordering {
	if left == right {
		return Equal
	} else if left < right {
		return Smaller
	} else {
		return Bigger
	}
}

type Ordering int
const (
	Smaller Ordering = iota
	Equal
	Bigger
)
func (o Ordering) String() string {
	switch o {
	case Smaller: return "L<R"
	case Bigger:  return "L>R"
	case Equal:   return "L=R"
	default:
		panic("impossible branch")
	}
}
func (o Ordering) Reversed() Ordering {
	switch o {
	case Smaller: return Bigger
	case Bigger:  return Smaller
	case Equal:   return Equal
	default:
		panic("impossible branch")
	}
}
func (cmp Compare[T]) OrderingReversed() Compare[T] {
	return func(a T, b T) Ordering {
		return cmp(a, b).Reversed()
	}
}


