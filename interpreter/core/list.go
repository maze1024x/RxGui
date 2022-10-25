package core

import "rxgui/util/ctn"


type List struct { Head *ListNode; Size uint }

type ListNode struct {
	Value  Object
	Next   *ListNode
}
func NodesToList(nodes ([] ListNode)) List {
	var n = len(nodes)
	if n > 0 {
		for i := 0; i < (n - 1); i += 1 {
			nodes[i].Next = &(nodes[i+1])
		}
		return List { &(nodes[0]), uint(n) }
	} else {
		return List {}
	}
}
func ToObjectList[T any] (slice ([] T)) List {
	var nodes = make([] ListNode, len(slice))
	for i := range slice {
		nodes[i].Value = ToObject[T](slice[i])
	}
	return NodesToList(nodes)
}
func EmptyList() List {
	return (List {})
}
func Cons(value Object, l List) List {
	return List {
		Head: &ListNode {
			Value: value,
			Next:  l.Head,
		},
		Size: (l.Size + 1),
	}
}
func Count(n int) List {
	if (n < 0) { n = 0 }
	var nodes = make([] ListNode, n)
	for i := 0; i < n; i += 1 {
		nodes[i].Value = ToObject(i)
	}
	return NodesToList(nodes)
}

type ListBuilder struct {
	nodes ([] ListNode)
}
func (buf *ListBuilder) Append(item Object) {
	buf.nodes = append(buf.nodes, ListNode { Value: item })
}
func (buf *ListBuilder) Collect() List {
	return NodesToList(buf.nodes)
}

type Seq struct {
	list  List
}
func (s Seq) Empty() bool {
	return s.list.Empty()
}
func (s Seq) Length() int {
	return s.list.Length()
}
func (s Seq) ToList() List {
	return s.list.Reversed()
}
func (s Seq) Last() (Object, bool) {
	return s.list.First()
}
func (s Seq) Appended(item Object) Seq {
	return Seq { Cons(item, s.list) }
}
func (s Seq) Sorted(lt ctn.Less[Object]) Seq {
	return Seq { s.list.Sorted(lt).Reversed() }
}
func (s Seq) Filter(f func(Object)(bool)) Seq {
	return Seq { s.list.Filter(f) }
}

func (l List) ToSeq() Seq {
	return Seq { l.Reversed() }
}
func (l List) Length() int {
	return int(l.Size)
}
func (l List) ForEach(f func(Object)) {
	var node = l.Head
	for node != nil {
		f(node.Value)
		node = node.Next
	}
}
func (l List) ForEachWithIndex(f func(int,Object)) int {
	var node = l.Head
	var i = 0
	for node != nil {
		f(i, node.Value)
		i += 1
		node = node.Next
	}
	return i
}
func (l List) Reversed() List {
	var items = make([] Object, 0)
	l.ForEach(func(item Object) {
		items = append(items, item)
	})
	var L = len(items)
	var buf ListBuilder
	for i := (L-1); i >= 0; i -= 1 {
		buf.Append(items[i])
	}
	return buf.Collect()
}
//go:noinline
func (l List) Empty() bool {
	return (l.Head == nil)
}
func (l List) First() (Object, bool) {
	if l.Head != nil {
		return l.Head.Value, true
	} else {
		return nil, false
	}
}
func (l List) Shifted() (Object, List, bool) {
	if l.Head != nil {
		var value = l.Head.Value
		var new_head = l.Head.Next
		var new_size = (l.Size - 1)
		return value, List{new_head,new_size}, true
	} else {
		return nil, List{}, false
	}
}
func (l List) Sorted(lt ctn.Less[Object]) List {
	var nodes = make([] ListNode, 0)
	l.ForEach(func(item Object) {
		nodes = append(nodes, ListNode { Value: item })
	})
	nodes, _ = ctn.StableSorted(nodes, func(a ListNode, b ListNode) bool {
		return lt(a.Value, b.Value)
	})
	return NodesToList(nodes)
}
func (l List) Take(limit int) List {
	if limit <= 0 {
		return List {}
	}
	var buf ListBuilder
	var count = 0
	var node = l.Head
	for node != nil {
		buf.Append(node.Value)
		count++
		node = node.Next
		if count == limit {
			break
		}
	}
	return buf.Collect()
}
func (l List) WithIndex() List {
	var nodes = make([] ListNode, 0)
	l.ForEachWithIndex(func(i int, a Object) {
		var pair = ToObject(ctn.MakePair(a, i))
		nodes = append(nodes, ListNode { Value: pair })
	})
	return NodesToList(nodes)
}
func (l List) Map(f func(Object)(Object)) List {
	var nodes = make([] ListNode, 0)
	l.ForEach(func(a Object) {
		var b = f(a)
		nodes = append(nodes, ListNode { Value: b })
	})
	return NodesToList(nodes)
}
func (l List) DeflateMap(f func(Object)(ctn.Maybe[Object])) List {
	var nodes = make([] ListNode, 0)
	l.ForEach(func(a Object) {
		if b, ok := f(a).Value(); ok {
			nodes = append(nodes, ListNode { Value: b })
		}
	})
	return NodesToList(nodes)
}
func (l List) FlatMap(f func(Object)(List)) List {
	var nodes = make([] ListNode, 0)
	l.ForEach(func(a Object) {
		var b = f(a)
		b.ForEach(func(b Object) {
			nodes = append(nodes, ListNode { Value: b })
		})
	})
	return NodesToList(nodes)
}
func (l List) ZipMap(m List, f func(Object,Object)(Object)) List {
	var u = l.Head
	var v = m.Head
	var nodes = make([] ListNode, 0)
	for {
		if u == nil || v == nil {
			return NodesToList(nodes)
		} else {
			var a, b = u.Value, v.Value
			var c = f(a, b)
			nodes = append(nodes, ListNode { Value: c })
			u = u.Next
			v = v.Next
		}
	}
}
func (l List) Filter(f func(Object)(bool)) List {
	var nodes = make([] ListNode, 0)
	l.ForEach(func(a Object) {
		if f(a) {
			nodes = append(nodes, ListNode { Value: a })
		}
	})
	return NodesToList(nodes)
}
func (l List) Scan(b Object, f func(Object,Object)(Object)) List {
	var nodes = make([] ListNode, 0)
	l.ForEach(func(a Object) {
		b = f(b, a)
		nodes = append(nodes, ListNode { Value: b })
	})
	return NodesToList(nodes)
}
func (l List) Fold(b Object, f func(Object,Object)(Object)) Object {
	l.ForEach(func(a Object) {
		b = f(b, a)
	})
	return b
}


