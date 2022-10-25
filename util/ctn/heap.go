package ctn

type Heap[T any] struct {
    LTT   *LTT[T]
    LtOp  Less[T]
}
func MakeHeap[T any] (lt Less[T]) Heap[T] {
    return Heap[T] { LTT: nil, LtOp: lt }
}
func (h Heap[T]) from(t *LTT[T]) Heap[T] {
    return Heap[T] { LTT: t, LtOp: h.LtOp }
}
func (h Heap[T]) ForEach(f func(v T)) {
    var current = h
    for { if v, next, ok := current.Shifted(); ok {
        f(v)
        current = next
    } else {
        break
    }}
}
func (h Heap[T]) Inserted(v T) Heap[T] {
    return h.from(h.LTT.Pushed(v, h.LtOp))
}
func (h Heap[T]) Shifted() (T, Heap[T], bool) {
    var popped, rest, exists = h.LTT.Popped(h.LtOp)
    return popped, h.from(rest), exists
}
func (h Heap[T]) First() (T, bool) {
    return h.LTT.Top()
}
func (h Heap[T]) IsEmpty() bool {
    return (h.LTT == nil)
}
func (h Heap[T]) Size() int {
    return int(h.LTT.GetSize())
}
func (h Heap[T]) MutHeap() MutHeap[T] {
    var ptr = new(Heap[T])
    *ptr = h
    return MutHeap[T] { ptr }
}

type MutHeap[T any] struct { ptr *Heap[T] }
func MakeMutHeap[T any] (lt Less[T]) MutHeap[T] {
    var h = MakeHeap(lt)
    return MutHeap[T] { &h }
}
func (mh MutHeap[T]) Heap() Heap[T] {
    return *(mh.ptr)
}
func (mh MutHeap[T]) ForEach(f func(v T)) {
    mh.ptr.ForEach(f)
}
func (mh MutHeap[T]) Insert(v T) {
    var pushed = mh.ptr.Inserted(v)
    *(mh.ptr) = pushed
}
func (mh MutHeap[T]) Shift() (T, bool) {
    var v, popped, ok = mh.ptr.Shifted()
    *(mh.ptr) = popped
    return v, ok
}
func (mh MutHeap[T]) First() (T, bool) {
    return mh.ptr.First()
}
func (mh MutHeap[T]) IsEmpty() bool {
    return mh.ptr.IsEmpty()
}
func (mh MutHeap[T]) Size() int {
    return mh.ptr.Size()
}


