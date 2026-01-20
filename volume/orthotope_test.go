// Copyright 2018 Brian Noyama. Subject to the the Apache License, Version 2.0.
package volume

import (
	"testing"
)

func TestOverlaps(t *testing.T) {
	o1 := Orthotope[int32]{P0: [DIM]int32{10, -20}, P1: [DIM]int32{40, 10}}
	o2 := Orthotope[int32]{P0: [DIM]int32{-10, 5}, P1: [DIM]int32{20, 35}}
	o3 := Orthotope[int32]{P0: [DIM]int32{-10, 25}, P1: [DIM]int32{20, 55}}

	overlaps := Overlaps(o1, o2)
	if !overlaps {
		t.Errorf("Expected orthtopes to overlap. Got %v.", overlaps)
	}

	overlaps = Overlaps(o1, o3)
	if overlaps {
		t.Errorf("Expected orthtopes to not overlap. Got %v.", overlaps)
	}
}

func TestContains(t *testing.T) {
	o1 := Orthotope[int32]{P0: [DIM]int32{10, -20}, P1: [DIM]int32{40, 10}}
	o2 := Orthotope[int32]{P0: [DIM]int32{15, -20}, P1: [DIM]int32{35, 0}}
	o3 := Orthotope[int32]{P0: [DIM]int32{-10, 5}, P1: [DIM]int32{20, 35}}

	contains := Contains(o1, o2)
	if !contains {
		t.Errorf("Expected orthtope to contain other. Got %v.", contains)
	}

	contains = Contains(o2, o1)
	if contains {
		t.Errorf("Expected orthtope to not contain other. Got %v.", contains)
	}

	contains = Contains(o1, o3)
	if contains {
		t.Errorf("Expected orthtope to not contain other. Got %v.", contains)
	}
}

func TestScore(t *testing.T) {
	o := Orthotope[int32]{P0: [DIM]int32{10, -20}, P1: [DIM]int32{40, -5}}

	score := Score(o)
	expected := float32(45)
	if score != expected {
		t.Errorf("Expected %v, got %v.", expected, score)
	}
}

func TestIntersects(t *testing.T) {
	o1 := Orthotope[int32]{P0: [DIM]int32{10, 15}, P1: [DIM]int32{30, 25}}
	o2 := Orthotope[int32]{P0: [DIM]int32{55, 65}, P1: [DIM]int32{75, 85}}
	o3 := Orthotope[int32]{P0: [DIM]int32{-20, 25}, P1: [DIM]int32{10, 45}}

	t1 := Intersects(o1, o2, [DIM]int32{25, 30}, [DIM]int32{-10, -30})
	t2 := Intersects(o1, o3, [DIM]int32{-5, 20}, [DIM]int32{0, 0})
	t3 := Intersects(o2, o3, [DIM]int32{0, -20}, [DIM]int32{0, 20})
	expected := float32(2)
	if t3 != expected {
		t.Errorf("Expected %v, got %v.", expected, t3)
	}
	if t1 == expected {
		t.Errorf("Expected something greater than %v, got %v.", expected, t1)
	}
	if t1 <= t2 {
		t.Errorf("Expected distance to be greater than %v, got %v.", t1, t2)
	}
}

func TestMinBounds(t *testing.T) {
	o1 := Orthotope[int32]{P0: [DIM]int32{10, -20}, P1: [DIM]int32{40, 10}}
	o2Orig := Orthotope[int32]{P0: [DIM]int32{15, -20}, P1: [DIM]int32{35, 0}}
	o2 := Orthotope[int32]{P0: [DIM]int32{15, -20}, P1: [DIM]int32{35, 0}}
	o3 := Orthotope[int32]{P0: [DIM]int32{-10, 5}, P1: [DIM]int32{20, 35}}

	o1 = MinBounds(o1, o2, o3)
	expected := Orthotope[int32]{P0: [DIM]int32{-10, -20}, P1: [DIM]int32{40, 35}}

	if o1 != expected {
		t.Errorf("Expected %v and %v doesn't match.", o1,
			expected)
	}
	if o2 != o2Orig {
		t.Errorf("Orthotope %v unintenitionally modified to %v.", o2Orig, o2)
	}
}

func TestOrthString(t *testing.T) {
	o1 := Orthotope[int32]{P0: [DIM]int32{10, -20}, P1: [DIM]int32{30, 30}}

	if String(o1) != "P0 [10 -20 0], P1 [30 30 0]" {
		t.Errorf("String method not working: %v", String(o1))
	}
}

func TestOrthEquals(t *testing.T) {
	o1 := Orthotope[int32]{P0: [DIM]int32{10, -20}, P1: [DIM]int32{40, 10}}
	o2 := Orthotope[int32]{P0: [DIM]int32{10, -20}, P1: [DIM]int32{40, 10}}
	o3 := Orthotope[int32]{P0: [DIM]int32{10, -5}, P1: [DIM]int32{40, 15}}
	o4 := Orthotope[int32]{P0: [DIM]int32{10, -5}, P1: [DIM]int32{40, 20}}

	if !Equals(o1, o2) {
		t.Errorf("%v should equal %v", o1, o2)
	}

	if Equals(o1, o3) {
		t.Errorf("%v should not equal %v", o1, o2)
	}

	if Equals(o4, o3) {
		t.Errorf("%v should not equal %v", o1, o2)
	}
}
