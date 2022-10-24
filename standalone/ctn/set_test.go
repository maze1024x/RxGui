package ctn

import "testing"

func TestSetBasic(t *testing.T) {
	var s = MakeMutSet[string](StringCompare,
		"foo",
		"bar",
	)
	s.Insert("baz")
	s.Delete("foo")
	if s.Has("foo") {
		t.Fatalf("wrong behavior for foo")
	}
	if !(s.Has("bar")) {
		t.Fatalf("wrong behavior for bar")
	}
	if !(s.Has("baz")) {
		t.Fatalf("wrong behavior for baz")
	}
	s.ForEach(func(item string) {
		t.Log(item)
	})
}


