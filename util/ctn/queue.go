package ctn

type Queue[T any] struct {
    Heap  Heap[Pair[uint64,T]]
    Next  uint64
}
func MakeQueue[T any] () Queue[T] {
    return Queue[T] {
        Heap: MakeHeap(func(p1 Pair[uint64,T], p2 Pair[uint64,T]) bool {
            return (p1.Key() < p2.Key())
        }),
        Next: 0,
    }
}
func (q Queue[T]) ForEach(f func(v T)) {
    q.Heap.ForEach(func(p Pair[uint64,T]) {
        f(p.Value())
    })
}
func (q Queue[T]) Appended(v T) Queue[T] {
    var n = q.Next
    return Queue[T] {
        Heap: q.Heap.Inserted(MakePair(n, v)),
        Next: (n + 1),
    }
}
func (q Queue[T]) Shifted() (T, Queue[T], bool) {
    var p, rest, ok = q.Heap.Shifted()
    if ok {
        return p.Value(), Queue[T] {
            Heap: rest,
            Next: q.Next,
        }, true
    } else {
        return zero[T](), q, false
    }
}
func (q Queue[T]) First() (T, bool) {
    var p, ok = q.Heap.First()
    if ok {
        return p.Value(), true
    } else {
        return zero[T](), false
    }
}
func (q Queue[T]) IsEmpty() bool {
    return q.Heap.IsEmpty()
}
func (q Queue[T]) Size() int {
    return q.Heap.Size()
}

type MutQueue[T any] struct { ptr *Queue[T] }
func MakeMutQueue[T any] () MutQueue[T] {
    var q = MakeQueue[T]()
    return MutQueue[T] { &q }
}
func (mq MutQueue[T]) Queue() Queue[T] {
    return *(mq.ptr)
}
func (mq MutQueue[T]) ForEach(f func(v T)) {
    mq.ptr.ForEach(f)
}
func (mq MutQueue[T]) Append(v T) {
    var appended = mq.ptr.Appended(v)
    *(mq.ptr) = appended
}
func (mq MutQueue[T]) Shift() (T, bool) {
    var v, shifted, ok = mq.ptr.Shifted()
    *(mq.ptr) = shifted
    return v, ok
}
func (mq MutQueue[T]) First() (T, bool) {
    return mq.ptr.First()
}
func (mq MutQueue[T]) IsEmpty() bool {
    return mq.ptr.IsEmpty()
}
func (mq MutQueue[T]) Size() int {
    return mq.ptr.Size()
}


