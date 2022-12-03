// Copyright 2018 Brian Noyama. Subject to the the Apache License, Version 2.0.
package rect

import (
	"fmt"
	"math"

	disc "github.com/briannoyama/bvh/discreet"
)

const DIMENSIONS int = 3

type Orthotope struct {
	Point [DIMENSIONS]int32
	Delta [DIMENSIONS]int32
}

var ACCURACY uint = 13

func (o *Orthotope) Overlaps(orth *Orthotope) bool {
	intersects := true
	for index, p0 := range orth.Point {
		p1 := orth.Delta[index] + p0
		intersects = intersects && o.Point[index] <= p1 &&
			p0 <= o.Point[index]+o.Delta[index]
	}
	return intersects
}

func (o *Orthotope) Contains(orth *Orthotope) bool {
	contains := true
	for index, p0 := range o.Point {
		p1 := o.Delta[index] + p0
		contains = contains && orth.Point[index] >= p0 &&
			p1 >= orth.Point[index]+orth.Delta[index]
	}
	return contains
}

/*Let orth represent a direction (a vector where delta defines direction).
 *Return t > 0 for where it intersects, or -1 if it does not intersect.
 */
func (orth *Orthotope) Intersects(o *Orthotope) int32 {
	inT := int32(0)
	outT := int32(math.MaxInt32)
	for index, p0 := range o.Point {
		p1 := o.Delta[index] + p0

		if orth.Delta[index] == 0 {
			if orth.Point[index] < p0 || p1 < orth.Point[index] {
				return -1
			}
		} else {
			if orth.Delta[index] < 0 {
				// Swap p0 and p1 for negative directions.
				p0, p1 = p1, p0
			}
			p0T := ((p0 - orth.Point[index]) << ACCURACY) / orth.Delta[index]
			inT = disc.Max(inT, p0T)

			p1T := ((p1 - orth.Point[index]) << ACCURACY) / orth.Delta[index]
			outT = disc.Min(outT, p1T)
		}
	}

	if inT < outT && inT >= 0 {
		return inT
	}
	return -1
}

func (o *Orthotope) MinBounds(others ...*Orthotope) {
	o.Point = others[0].Point
	o.Delta = others[0].Delta

	for index, p0 := range o.Point {
		p1 := p0 + o.Delta[index]

		for _, other := range others[1:] {
			o.Point[index] = disc.Min(p0, other.Point[index])
			p1 = disc.Max(p1, other.Point[index]+other.Delta[index])
		}
		o.Delta[index] = p1 - o.Point[index]
	}
}

func (o *Orthotope) Volume() int32 {
	v := int32(1)
	for _, d := range o.Delta {
		v *= d
	}
	return v
}

func (o *Orthotope) SurfaceArea() int32 {
	if DIMENSIONS == 1 {
		return 0
	}

	v := o.Volume()
	sa := int32(0)
	for i := 0; i < DIMENSIONS; i++ {
		sa += v / o.Delta[i]
	}
	return 2 * sa
}

func (o *Orthotope) Score() int32 {
	score := int32(0)
	for _, d := range o.Delta {
		score += d
	}
	return score
}

func (o *Orthotope) Equals(other *Orthotope) bool {
	for index, point := range other.Point {
		if o.Point[index] != point {
			return false
		} else if o.Delta[index] != other.Delta[index] {
			return false
		}
	}
	// Return 0 if the orthtopes are equal
	return true
}

// Get a string representation of this orthotope.
func (o *Orthotope) String() string {
	return fmt.Sprintf("Point %v, Delta %v", o.Point, o.Delta)
}
