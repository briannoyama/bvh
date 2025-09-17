// Copyright 2018 Brian Noyama. Subject to the the Apache License, Version 2.0.
package rect

import (
	"fmt"
	"math"

)

const DIMENSIONS int = 3

type Num interface {
	int32 | int64 | float32 | float64
}
type Orthotope[T Num] struct {
	Point [DIMENSIONS]T
	Delta [DIMENSIONS]T
}

var ACCURACY uint = 13

func (o *Orthotope[T]) Overlaps(orth *Orthotope[T]) bool {
	intersects := true
	for index, p0 := range orth.Point {
		p1 := orth.Delta[index] + p0
		intersects = intersects && o.Point[index] <= p1 &&
			p0 <= o.Point[index]+o.Delta[index]
	}
	return intersects
}

func (o *Orthotope[T]) Contains(orth *Orthotope[T]) bool {
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
func (orth *Orthotope[T]) Intersects(o *Orthotope[T]) T {
	inT := T(0)
	outT := T(math.MaxInt32)
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
			p0subtract := p0 - orth.Point[index]
			p1subtract := p1 - orth.Point[index]
			pow2 := math.Pow(2, float64(ACCURACY))

			p0T := (p0subtract * T(pow2)) / orth.Delta[index]
			inT = max(inT, p0T)

			p1T := (p1subtract * T(pow2)) / orth.Delta[index]
			outT = min(outT, p1T)
		}
	}

	if inT < outT && inT >= 0 {
		return inT
	}
	return -1
}

func (o *Orthotope[T]) MinBounds(others ...*Orthotope[T]) {
	o.Point = others[0].Point
	o.Delta = others[0].Delta

	for index, p0 := range o.Point {
		p1 := p0 + o.Delta[index]

		for _, other := range others[1:] {
			o.Point[index] = min(p0, other.Point[index])
			p1 = max(p1, other.Point[index]+other.Delta[index])
		}
		o.Delta[index] = p1 - o.Point[index]
	}
}

func (o *Orthotope[T]) Volume() T {
	v := T(1)
	for _, d := range o.Delta {
		v *= d
	}
	return v
}

func (o *Orthotope[T]) SurfaceArea() T {
	if DIMENSIONS == 1 {
		return 0
	}

	v := o.Volume()
	sa := T(0)
	for i := 0; i < DIMENSIONS; i++ {
		sa += v / o.Delta[i]
	}
	return 2 * sa
}

func (o *Orthotope[T]) Score() T {
	score := T(0)
	for _, d := range o.Delta {
		score += d
	}
	return score
}

func (o *Orthotope[T]) Equals(other *Orthotope[T]) bool {
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
func (o *Orthotope[T]) String() string {
	return fmt.Sprintf("Point %v, Delta %v", o.Point, o.Delta)
}
