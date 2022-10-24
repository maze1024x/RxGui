package ctn

import (
	"testing"
	"sort"
)

func TestHeapBasic(t *testing.T) {
	var numbers = [] float64 { 7, 2, 3, 5, 4, 6, 1 }
	var h = MakeMutHeap[float64](func(x float64, y float64) bool { return x < y })
	for _, x := range numbers {
		h.Insert(x)
	}
	sort.Float64s(numbers)
	for _, x := range numbers[:(len(numbers) - 1)] {
		if value, ok := h.Shift(); ok {
			if value != x {
				t.Fatalf("wrong Shift result")
			}
		} else {
			t.Fatalf("wrong Shift/Insert behavior")
		}
	}
	if value, ok := h.First(); ok {
		if value != numbers[len(numbers) - 1] {
			t.Fatalf("wrong First result")
		}
	} else {
		t.Fatalf("wrong Shift/Insert behavior")
	}
	h.Shift()
	if _, ok := h.First(); ok {
		t.Fatalf("wrong Shift/Insert behavior")
	}
}


