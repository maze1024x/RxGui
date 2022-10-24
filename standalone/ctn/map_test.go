package ctn

import (
	"testing"
	"math/rand"
)

func TestMapBasic(t *testing.T) {
	var m = MakeMutMap[string,int](StringCompare,
		MakePair("a",1),
		MakePair("b",2),
	)
	m.Insert("c", 3)
	m.Insert("c", 4)
	m.Delete("b")
	{ var v, found = m.Lookup("a")
	if !(found && v == 1) {
		t.Fatalf("wrong behavior for a")
	} }
	if m.Has("b") {
		t.Fatalf("wrong behavior for b")
	}
	{ var v, found = m.Lookup("c")
	if !(found && v == 4) {
		t.Fatalf("wrong behavior for c")
	} }
	m.ForEach(func(key string, value int) {
		t.Log(key, value)
	})
}

// TODO: test for MergedWith, ...

func TestToOrderedMap(t *testing.T) {
	var ordered_keys = [] string { "1", "2", "3", "4", "5", "6", "7", "8", "9" }
	for trial := 0; trial < 100; trial += 1 {
		var h = make(map[string] int)
		var keys = make([] string, len(ordered_keys))
		for i, key := range ordered_keys {
			keys[i] = key
		}
		rand.Shuffle(len(keys), func(i, j int) {
			var tmp = keys[i]
			keys[i] = keys[j]
			keys[j] = tmp
		})
		for i, key := range keys {
			h[key] = (i+1)
		}
		var i = 0
		ToOrderedMap(h).ForEach(func(key string, value int) {
			if key != ordered_keys[i] {
				t.Fatalf("keys are not ordered (trial %d)", trial)
			}
			i += 1
		})
	}
}


