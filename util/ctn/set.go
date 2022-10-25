package ctn


type Set[T any] struct {
    AVL  *AVL[T]
    Cmp  Compare[T]
}
func MakeSet[T any] (cmp Compare[T], items ...T) Set[T] {
    var s = Set[T] { AVL: nil, Cmp: cmp }
    for _, item := range items {
        s = s.Inserted(item)
    }
    return s
}

func (s Set[T]) IsEmpty() bool {
    return s.AVL == nil
}
func (s Set[T]) Size() int {
    return int(s.AVL.GetSize())
}
func (s Set[T]) from(a *AVL[T]) Set[T] {
    return Set[T] { AVL: a, Cmp: s.Cmp }
}

func (s Set[T]) ForEach(f func(T)) {
    s.AVL.Walk(f)
}
func (s Set[T]) Has(v T) bool {
    var _, found = s.AVL.Lookup(v, s.Cmp)
    return found
}
func (s Set[T]) Inserted(v T) Set[T] {
    var inserted, _ = s.AVL.Inserted(v, s.Cmp)
    return s.from(inserted)
}
func (s Set[T]) Deleted(v T) (Set[T], bool) {
    var _, deleted, exists = s.AVL.Deleted(v, s.Cmp)
    return s.from(deleted), exists
}
func (s Set[T]) MutSet() MutSet[T] {
    var ptr = new(Set[T])
    *ptr = s
    return MutSet[T] { ptr }
}

type MutSet[T any] struct { ptr *Set[T] }
func MakeMutSet[T any] (cmp Compare[T], items ...T) MutSet[T] {
    var s = MakeSet[T](cmp, items...)
    return MutSet[T] { &s }
}
func (ms MutSet[T]) Set() Set[T] {
    return *(ms.ptr)
}
func (ms MutSet[T]) ForEach(f func(T)) {
    ms.ptr.ForEach(f)
}
func (ms MutSet[T]) Has(v T) bool {
    return ms.ptr.Has(v)
}
func (ms MutSet[T]) Insert(v T) {
    var inserted = ms.ptr.Inserted(v)
    *(ms.ptr) = inserted
}
func (ms MutSet[T]) Delete(v T) bool {
    var deleted, ok = ms.ptr.Deleted(v)
    *(ms.ptr) = deleted
    return ok
}
func (ms MutSet[T]) Size() int {
    return ms.ptr.Size()
}
func (ms MutSet[T]) Clone() MutSet[T] {
    return (*(ms.ptr)).MutSet()
}

func Items[T any] (f func(func(T))) ([] T) {
    var items = make([] T, 0)
    f(func(t T) {
        items = append(items, t)
    })
    return items
}


