package ctn

import (
	"testing"
	"sort"
)


func TestMapFilterReduce(t *testing.T) {
	s := [] int { 1, 2, 3, 4, 5, 6, 7, 8, 9 }
	s = MapEach(s, func(i int) int {
		return i*10
	})
	s = Filter(s, func(i int) bool {
		return (i <= 40)
	})
	{ sum := Reduce(s, int(1), func(sum int, i int) int {
		return (sum + i)
	})
	if sum != 101 {
		t.Fatalf("wrong result: %d", sum)
	} }
}

func TestMapEachDeflate(t *testing.T) {
	s := [] int { 1, 5, 8, 3, 6, 4, 9, 2 }
	expected := [] int { -1, -3, -4, -2 }
	s = MapEachDeflate(s, func(a int) (int, bool) {
		if (a <= 4) {
			return -a, true
		}  else {
			return 0, false
		}
	})
	if len(s) != len(expected) {
		t.Fatalf("wrong result: %+v", s)
	}
	for i := range s {
		if s[i] != expected[i] {
			t.Fatalf("wrong result: %+v", s)
		}
	}
}

func TestRemoveFrom(t *testing.T) {
	s := [] int { 1, 1, 2, 3, 5, 8, 1, 1, 2, 3, 5, 8 }
	expected := [] int { 2, 3, 5, 8, 2, 3, 5, 8 }
	s = RemoveFrom(s, 1)
	if len(s) != len(expected) {
		t.Fatalf("wrong result: %+v", s)
	}
	for i := range s {
		if s[i] != expected[i] {
			t.Fatalf("wrong result: %+v", s)
		}
	}
}

func TestStableSort(t *testing.T) {
	u := [] int { 50, 55, 60, 65, 20, 30, 35, 80, 85, 90, 10, 12, 18, 70, 40 }
	v, _ := (StableSorted(u, func(a int, b int) bool {
		return ((a / 10) < (b / 10))
	}))
	sort.Slice(u, func(i, j int) bool {
		return (u[i] < u[j])
	})
	for i := range u {
		if u[i] != v[i] {
			t.Fatalf("wrong result: %+v", v)
		}
	}
}


