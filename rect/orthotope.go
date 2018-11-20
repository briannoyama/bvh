//Copyright 2018 Brian Noyama. Subject to the the Apache License, Version 2.0.
package rect

import (
	"fmt"
	"math"

	disc "github.com/briannoyama/bvh/discreet"
)

const DIMENSIONS int = 3

type Orthotope struct {
	Point [DIMENSIONS]int
	Delta [DIMENSIONS]int
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
func (orth *Orthotope) Intersects(o *Orthotope) int {
	in_t := 0
	out_t := math.MaxInt32
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
			p0_t := ((p0 - orth.Point[index]) << ACCURACY) / orth.Delta[index]
			in_t = disc.Max(in_t, p0_t)

			p1_t := ((p1 - orth.Point[index]) << ACCURACY) / orth.Delta[index]
			out_t = disc.Min(out_t, p1_t)
		}
	}

	if in_t < out_t && in_t >= 0 {
		return in_t
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

func (o *Orthotope) Score() int {
	score := 0
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
