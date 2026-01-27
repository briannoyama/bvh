package volume

import (
	"fmt"
)

// Orthotope s are N dimensional rectangular polyhedra defined by a point (location) and a delta (width, height, etc)
type Orthotope[K Num] struct {
	P0, P1 [DIM]K
}

// Overlaps returns true if two orthotopes intersect
func (o0 Orthotope[K]) Overlaps(o1 Orthotope[K]) bool {
	intersects := true
	for index, p0 := range o0.P0 {
		p1 := o0.P1[index]
		intersects = intersects && o1.P0[index] <= p1 &&
			p0 <= o1.P1[index]
	}
	return intersects
}

// Contains returns true if all of orth is within the bounds of o0. Ie. the intersection of o0 and 01 is equivalent to o1
func (o0 Orthotope[K]) Contains(o1 Orthotope[K]) bool {
	contains := true
	for index, p0 := range o0.P0 {
		p1 := o0.P1[index]
		contains = contains && o1.P0[index] >= p0 &&
			p1 >= o1.P1[index]
	}
	return contains
}

// TaxiPath returns 0,0,0 if this Orthotope contains the other, else returns the distance to contain it
func (o0 Orthotope[K]) TaxiPath(point [DIM]K) [DIM]K {
	coor := [DIM]K{}
	for index, p0 := range o0.P0 {
		coor[index] = p0 - point[index] // Negative if within
		coor[index] = max(coor[index], 0)

		other := o0.P1[index] - point[index] // Positive if within
		coor[index] += max(other, 0)
	}
	return coor
}

// Intersects return 0 <= t <= 1 for where the orth intersects along the delta, else t = 2 when there's no intersection
func (o0 Orthotope[K]) Intersects(o1 Orthotope[K], d [DIM]K) float32 {
	inT := float32(0)
	outT := float32(1)
	for index, p0 := range o0.P0 {
		p1 := o0.P1[index]

		delta := float32(d[index])
		if delta == 0 {
			if o1.P0[index] > p1 || p0 > o1.P1[index] {
				return 2
			}
		} else {
			p0T := float32(o1.P0[index]-p1) / delta
			p1T := float32(o1.P1[index]-p0) / delta

			if delta < 0 {
				// Swap p0 and p1 for negative directions.
				p0T, p1T = p1T, p0T
			}
			inT = max(inT, p0T)
			outT = min(outT, p1T)
		}
	}

	if inT <= outT {
		return inT
	}
	return 2
}

// MinBounds modifies point and delta such to that the resulting orthotope is the smallest one that can possibly contain
// all others
func (o0 Orthotope[K]) Minbound(o1 Orthotope[K]) Orthotope[K] {
	for index, p0 := range o0.P0 {
		o0.P0[index] = min(p0, o1.P0[index])
		o0.P1[index] = max(o0.P1[index], o1.P1[index])
	}
	return o0
}

// Score adds the lengths of the sides. This is the heuristic used to rebalance collision.BVol objects via swapChecks
func (o0 Orthotope[K]) Score() float32 {
	var score K
	for index, p0 := range o0.P0 {
		score += o0.P1[index] - p0
	}
	return float32(score)
}

// Equals checks if two Orthotopes are equivalent
func (o0 Orthotope[K]) Equals(o1 Orthotope[K]) bool {
	return o0 == o1
}

// String representation of this orth
func (o0 Orthotope[K]) String() string {
	return fmt.Sprintf("P0 %v, P1 %v", o0.P0, o0.P1)
}
