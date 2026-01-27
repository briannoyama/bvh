package volume

import (
	"fmt"
	"math"
)

type Sphere[K Num] struct {
	P [DIM]K
	R K
}

// Overlaps returns true if two spheres intersect
func (s0 Sphere[K]) Overlaps(s1 Sphere[K]) bool {
	_, distance := s0.diff(s1)
	bound := s0.R + s1.R
	return bound*bound > distance
}

// Contains returns true if all of sphere is within the bounds of o0. Ie. the intersection of o0 and 01 is equivalent to o1
func (s0 Sphere[K]) Contains(s1 Sphere[K]) bool {
	_, distance := s0.diff(s1)
	bound := s0.R - s1.R
	return bound >= 0 && bound*bound >= distance
}

// Intersects return 0 <= t <= 1 for where the sphere intersects along the delta, else t = 2 when there's no intersection
func (s0 Sphere[K]) Intersects(s1 Sphere[K], d [DIM]K) float32 {
	if s0.Overlaps(s1) {
		return 0
	}
	var a, b, c float64
	for index, p := range s0.P {
		pd := s1.P[index] - p
		a += float64(d[index] * d[index])
		b += float64(d[index] * pd)
		c += float64(pd * pd)
	}
	rd := float64(s0.R + s1.R)
	b *= -2
	c -= rd * rd

	if b*b > 4*a*c {
		sqrt := math.Sqrt(b*b - 4*a*c)
		t0 := (-b - sqrt) / (2 * a)
		t1 := (-b + sqrt) / (2 * a)
		t0 = min(t0, t1)
		if t0 >= 0 && t0 <= 1 {
			return float32(t0)
		}
	}

	return 2
}

// MinBounds modifies point and delta such to that the resulting Sphere is the smallest one that can possibly contain
// all others
func (s0 Sphere[K]) Minbound(s1 Sphere[K]) Sphere[K] {
	diff, distance := s0.diff(s1)
	distance = K(math.Sqrt(float64(distance)))
	midOffset := s1.R - s0.R
	s := Sphere[K]{}
	for index, p := range s0.P {
		diff[index] += diff[index] * midOffset / distance
		s.P[index] = p + diff[index]/2
	}
	// Add 1 to radius to ensure volume is bounded (prevents rouding error)
	s.R = (distance+s1.R+s0.R)/2 + 1
	return s
}

func (s0 Sphere[K]) diff(s1 Sphere[K]) ([DIM]K, K) {
	diff := [DIM]K{}
	var distance K
	for index, p := range s0.P {
		diff[index] = s1.P[index] - p
		distance += diff[index] * diff[index]
	}
	return diff, distance
}

// Score adds the lengths of the sides. This is the heuristic used to rebalance collision.BVol objects via swapChecks
func (s0 Sphere[K]) Score() float32 {
	return float32(s0.R)
}

// Equals checks if two Spheres are equivalent
func (s0 Sphere[K]) Equals(s1 Sphere[K]) bool {
	return s0 == s1
}

// String representation of this sphere
func (s0 Sphere[K]) String() string {
	return fmt.Sprintf("P %v, R %v", s0.P, s0.R)
}
