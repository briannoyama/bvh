// Copyright 2018 Brian Noyama. Subject to the the Apache License, Version 2.0.
package volume

import (
	"testing"
)

func TestOverlaps(t *testing.T) {
	o1 := Orthotope[int32]{P0: [DIM]int32{10, -20}, P1: [DIM]int32{40, 10}}
	o2 := Orthotope[int32]{P0: [DIM]int32{-10, 5}, P1: [DIM]int32{20, 35}}
	o3 := Orthotope[int32]{P0: [DIM]int32{-10, 25}, P1: [DIM]int32{20, 55}}

	overlaps := o1.Overlaps(o2)
	if !overlaps {
		t.Errorf("Expected orthtopes to overlap. Got %v.", overlaps)
	}

	overlaps = o1.Overlaps(o3)
	if overlaps {
		t.Errorf("Expected orthtopes to not overlap. Got %v.", overlaps)
	}
}

func TestContains(t *testing.T) {
	o1 := Orthotope[int32]{P0: [DIM]int32{10, -20}, P1: [DIM]int32{40, 10}}
	o2 := Orthotope[int32]{P0: [DIM]int32{15, -20}, P1: [DIM]int32{35, 0}}
	o3 := Orthotope[int32]{P0: [DIM]int32{-10, 5}, P1: [DIM]int32{20, 35}}

	contains := o1.Contains(o2)
	if !contains {
		t.Errorf("Expected orthtope to contain other. Got %v.", contains)
	}

	contains = o2.Contains(o1)
	if contains {
		t.Errorf("Expected orthtope to not contain other. Got %v.", contains)
	}

	contains = o1.Contains(o3)
	if contains {
		t.Errorf("Expected orthtope to not contain other. Got %v.", contains)
	}
}

func TestScore(t *testing.T) {
	o := Orthotope[int32]{P0: [DIM]int32{10, -20}, P1: [DIM]int32{40, -5}}

	score := o.Score()
	expected := float32(45)
	if score != expected {
		t.Errorf("Expected %v, got %v.", expected, score)
	}
}

func TestIntersects(t *testing.T) {
	o1 := Orthotope[int32]{P0: [DIM]int32{10, 15}, P1: [DIM]int32{30, 25}}
	o2 := Orthotope[int32]{P0: [DIM]int32{55, 65}, P1: [DIM]int32{75, 85}}
	o3 := Orthotope[int32]{P0: [DIM]int32{-20, 25}, P1: [DIM]int32{10, 45}}

	t1 := o1.Intersects(o2, [DIM]int32{35, 60})
	t2 := o1.Intersects(o3, [DIM]int32{-5, 20})
	t3 := o2.Intersects(o3, [DIM]int32{0, -40})
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
	o2 := Orthotope[int32]{P0: [DIM]int32{15, -20}, P1: [DIM]int32{35, 0}}
	o3 := Orthotope[int32]{P0: [DIM]int32{-10, 5}, P1: [DIM]int32{20, 35}}

	o2 = o2.Minbound(o3)
	expected := Orthotope[int32]{P0: [DIM]int32{-10, -20}, P1: [DIM]int32{35, 35}}

	if o2 != expected {
		t.Errorf("Expected %v and %v doesn't match.", o2,
			expected)
	}
}

func TestOrthString(t *testing.T) {
	o1 := Orthotope[int32]{P0: [DIM]int32{10, -20}, P1: [DIM]int32{30, 30}}

	if o1.String() != "P0 [10 -20 0], P1 [30 30 0]" {
		t.Errorf("String method not working: %v", o1.String())
	}
}

func TestOrthEquals(t *testing.T) {
	o1 := Orthotope[int32]{P0: [DIM]int32{10, -20}, P1: [DIM]int32{40, 10}}
	o2 := Orthotope[int32]{P0: [DIM]int32{10, -20}, P1: [DIM]int32{40, 10}}
	o3 := Orthotope[int32]{P0: [DIM]int32{10, -5}, P1: [DIM]int32{40, 15}}
	o4 := Orthotope[int32]{P0: [DIM]int32{10, -5}, P1: [DIM]int32{40, 20}}

	if !o1.Equals(o2) {
		t.Errorf("%v should equal %v", o1, o2)
	}

	if o1.Equals(o3) {
		t.Errorf("%v should not equal %v", o1, o2)
	}

	if o4.Equals(o3) {
		t.Errorf("%v should not equal %v", o1, o2)
	}
}
