package volume

import (
	"testing"
)

func TestSphereOverlaps(t *testing.T) {
	s1 := Sphere[int32]{P: [DIM]int32{3, 3}, R: 5}
	s2 := Sphere[int32]{P: [DIM]int32{2, -1}, R: 3}
	s3 := Sphere[int32]{P: [DIM]int32{8, 8}, R: 2}

	overlaps := s1.Overlaps(s2)
	if !overlaps {
		t.Errorf("Expected orthtopes to overlap. Got %v.", overlaps)
	}

	overlaps = s1.Overlaps(s3)
	if overlaps {
		t.Errorf("Expected orthtopes to not overlap. Got %v.", overlaps)
	}
}

func TestSphereContains(t *testing.T) {
	s1 := Sphere[int32]{P: [DIM]int32{3, 3}, R: 5}
	s2 := Sphere[int32]{P: [DIM]int32{1, 1}, R: 2}
	s3 := Sphere[int32]{P: [DIM]int32{8, 8}, R: 2}

	contains := s1.Contains(s2)
	if !contains {
		t.Errorf("Expected orthtope to contain other. Got %v.", contains)
	}

	contains = s2.Contains(s1)
	if contains {
		t.Errorf("Expected orthtope to not contain other. Got %v.", contains)
	}

	contains = s1.Contains(s3)
	if contains {
		t.Errorf("Expected orthtope to not contain other. Got %v.", contains)
	}
}

func TestSphereIntersects(t *testing.T) {
	s1 := Sphere[int32]{P: [DIM]int32{3, 3}, R: 5}
	s2 := Sphere[int32]{P: [DIM]int32{2, -1}, R: 3}
	s3 := Sphere[int32]{P: [DIM]int32{8, 8}, R: 2}

	t1 := s1.Intersects(s2, [DIM]int32{35, 60})
	t2 := s2.Intersects(s3, [DIM]int32{12, 18})
	t3 := s2.Intersects(s3, [DIM]int32{12, 1})

	if t1 != 0 {
		t.Errorf("Expected 0, got %v.", t1)
	}
	if t2 >= 1 {
		t.Errorf("Expected distance to be less than 1 got %v.", t2)
	}
	if t3 != 2 {
		t.Errorf("Expected 2, got %v.", t3)
	}
}

func TestSphereMinbound(t *testing.T) {
	s2 := Sphere[int32]{P: [DIM]int32{2, -1}, R: 6}
	s3 := Sphere[int32]{P: [DIM]int32{8, 8}, R: 2}
	s1 := s2.Minbound(s3)
	expected := Sphere[int32]{P: [DIM]int32{4, 2}, R: 10}

	if s1 != expected {
		t.Errorf("Expected %v and %v doesn't match.", s1, expected)
	}

	if !s1.Contains(s2) || !s1.Contains(s3) {
		t.Errorf("Volume %v does not contain %v and %v", s1, s2, s3)
	}
}

func TestSphereScore(t *testing.T) {
	s := Sphere[int32]{P: [DIM]int32{3, 3}, R: 5}
	score := s.Score()
	expected := float32(5)
	if score != expected {
		t.Errorf("Expected %v, got %v.", expected, score)
	}
}

func TestSphereEquals(t *testing.T) {
	s1 := Sphere[int32]{P: [DIM]int32{3, 3}, R: 5}
	s2 := Sphere[int32]{P: [DIM]int32{2, -1}, R: 3}
	s3 := Sphere[int32]{P: [DIM]int32{3, 3}, R: 5}

	if !s1.Equals(s3) {
		t.Errorf("%v should equal %v", s1, s2)
	}

	if s1.Equals(s2) {
		t.Errorf("%v should not equal %v", s1, s2)
	}
}

func TestSphereString(t *testing.T) {
	s1 := Sphere[int32]{P: [DIM]int32{3, 3}, R: 5}

	if s1.String() != "P [3 3 0], R 5" {
		t.Errorf("String method not working: %v", s1.String())
	}
}
