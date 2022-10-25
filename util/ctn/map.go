package ctn

import "rxgui/standalone/util/backports/constraints"


type Map[K any, V any] struct {
	AVL  *AVL[Pair[K,V]]
	Cmp  Compare[Pair[K,V]]
}
func MakeMap[K any, V any] (cmp Compare[K], entries ...Pair[K,V]) Map[K,V] {
	var m = Map[K,V] {
		AVL: nil,
		Cmp: func(left Pair[K,V], right Pair[K,V]) Ordering {
			return cmp(left.Key(), right.Key())
		},
	}
	for _, entry := range entries {
		m = m.Inserted(entry.Key(), entry.Value())
	}
	return m
}
func ToOrderedMap[K constraints.Ordered, V any] (h (map[K] V)) Map[K,V] {
	var m = MakeMutMap[K,V](func(left K, right K) Ordering {
		if left < right {
			return Smaller
		} else if left == right {
			return Equal
		} else {
			return Bigger
		}
	})
	for k, v := range h {
		m.Insert(k, v)
	}
	return m.Map()
}

func (m Map[K,V]) IsEmpty() bool {
	return m.AVL == nil
}
func (m Map[K,V]) Size() int {
	return int(m.AVL.GetSize())
}
func (m Map[K,V]) from(a *AVL[Pair[K,V]]) Map[K,V] {
	return Map[K,V] { AVL: a, Cmp: m.Cmp }
}

func (m Map[K,V]) ForEach(f func(K,V)) {
	m.AVL.Walk(func(entry Pair[K,V]) {
		f(entry.Key(), entry.Value())
	})
}
func (m Map[K,V]) Has(k K) bool {
	var _, found = m.Lookup(k)
	return found
}
func (m Map[K,V]) Lookup(k K) (V, bool) {
	var entry, exists = m.AVL.Lookup(MakePair(k, zero[V]()), m.Cmp)
	if exists {
		return entry.Value(), true
	} else {
		return zero[V](), false
	}
}
func (m Map[K,V]) Inserted(k K, v V) Map[K,V] {
	var arg = MakePair(k, v)
	var inserted, _ = m.AVL.Inserted(arg, m.Cmp)
	return m.from(inserted)
}
func (m Map[K,V]) Deleted(k K) (V, Map[K,V], bool) {
	var arg = MakePair(k, zero[V]())
	var entry, deleted, exists = m.AVL.Deleted(arg, m.Cmp)
	if exists {
		return entry.Value(), m.from(deleted), true
	} else {
		return zero[V](), m, false
	}
}
func (m Map[K,V]) MergedWith(another Map[K,V]) Map[K,V] {
	var draft = m
	another.ForEach(func(k K, v V) {
		draft = draft.Inserted(k, v)
	})
	return draft
}
func (m Map[K,V]) MutMap() MutMap[K,V] {
	var ptr = new(Map[K,V])
	*ptr = m
	return MutMap[K,V] { ptr }
}

type MutMap[K any, V any] struct { ptr *Map[K,V] }
func MakeMutMap[K any, V any] (cmp Compare[K], entries ...Pair[K,V]) MutMap[K,V] {
	var m = MakeMap[K,V](cmp, entries...)
	return MutMap[K,V] { &m }
}
func (mm MutMap[K,V]) Map() Map[K,V] {
	return *(mm.ptr)
}
func (mm MutMap[K,V]) ForEach(f func(K,V)) {
	mm.ptr.ForEach(f)
}
func (mm MutMap[K,V]) Has(k K) bool {
	return mm.ptr.Has(k)
}
func (mm MutMap[K,V]) Lookup(k K) (V, bool) {
	return mm.ptr.Lookup(k)
}
func (mm MutMap[K,V]) Insert(k K, v V) {
	var inserted = mm.ptr.Inserted(k, v)
	*(mm.ptr) = inserted
}
func (mm MutMap[K,V]) Delete(k K) (V, bool) {
	var value, deleted, ok = mm.ptr.Deleted(k)
	*(mm.ptr) = deleted
	return value, ok
}
func (mm MutMap[K,V]) Size() int {
	return mm.ptr.Size()
}
func (mm MutMap[K,V]) Clone() MutMap[K,V] {
	return (*(mm.ptr)).MutMap()
}
func (mm MutMap[K,V]) FilterClone(f func(K,V)(bool)) MutMap[K,V] {
	var m = *(mm.ptr)
	mm.ptr.ForEach(func(k K, v V) {
		if !(f(k,v)) {
			_, m, _ = m.Deleted(k)
		}
	})
	return m.MutMap()
}

func Keys[K any, V any] (f func(func(K,V))) ([] K) {
	var keys = make([] K, 0)
	f(func(k K, _ V) {
		keys = append(keys, k)
	})
	return keys
}
func Values[K any, V any] (f func(func(K,V))) ([] V) {
	var values = make([] V, 0)
	f(func(_ K, v V) {
		values = append(values, v)
	})
	return values
}
func Entries[K any, V any] (f func(func(K,V))) ([] Pair[K,V]) {
	var entries = make([] Pair[K,V], 0)
	f(func(k K, v V) {
		entries = append(entries, MakePair(k, v))
	})
	return entries
}


