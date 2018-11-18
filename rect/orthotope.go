package rect

import (
	"fmt"
	"math"

	disc "github.com/briannoyama/bvh/discreet"
)

const DIMENSIONS int = 2

type Orthotope struct {
	point [DIMENSIONS]int
	delta [DIMENSIONS]int
}

var ACCURACY uint = 13

func (o *Orthotope) Overlaps(orth *Orthotope) bool {
	intersects := true
	for index, p0 := range orth.point {
		p1 := orth.delta[index] + p0
		intersects = intersects && o.point[index] < p1 &&
			p0 < o.point[index]+o.delta[index]
	}
	return intersects
}

func (o *Orthotope) Contains(orth *Orthotope) bool {
	contains := true
	for index, p0 := range o.point {
		p1 := o.delta[index] + p0
		contains = contains && orth.point[index] >= p0 &&
			p1 >= orth.point[index]+orth.delta[index]
	}
	return contains
}

/*Let orth represent a direction (a vector where delta defines direction).
 *Return t > 0 for where it intersects, or -1 if it does not intersect.
 */
func (orth *Orthotope) Intersects(o *Orthotope) int {
	in_t := 0
	out_t := math.MaxInt32
	for index, p0 := range o.point {
		p1 := o.delta[index] + p0

		if orth.delta[index] == 0 {
			if orth.point[index] < p0 || p1 < orth.point[index] {
				return -1
			}
		} else {
			if orth.delta[index] < 0 {
				// Swap p0 and p1 for negative directions.
				p0, p1 = p1, p0
			}
			p0_t := ((p0 - orth.point[index]) << ACCURACY) / orth.delta[index]
			in_t = disc.Max(in_t, p0_t)

			p1_t := ((p1 - orth.point[index]) << ACCURACY) / orth.delta[index]
			out_t = disc.Min(out_t, p1_t)
		}
	}

	if in_t < out_t && in_t >= 0 {
		return in_t
	}
	return -1
}

func (o *Orthotope) MinBounds(others ...*Orthotope) {
	o.point = others[0].point
	o.delta = others[0].delta

	for index, p0 := range o.point {
		p1 := p0 + o.delta[index]

		for _, other := range others[1:] {
			o.point[index] = disc.Min(p0, other.point[index])
			p1 = disc.Max(p1, other.point[index]+other.delta[index])
		}
		o.delta[index] = p1 - o.point[index]
	}
}

func (o *Orthotope) Score() int {
	score := 0
	for _, d := range o.delta {
		score += d
	}
	return score
}

func (o *Orthotope) Equals(other *Orthotope) bool {
	for index, point := range other.point {
		if o.point[index] != point {
			return false
		} else if o.delta[index] != other.delta[index] {
			return false
		}
	}
	// Return 0 if the orthtopes are equal
	return true
}

// Get a string representation of this orthotope.
func (o *Orthotope) String() string {
	return fmt.Sprintf("Point %v, Delta %v", o.point, o.delta)
}
