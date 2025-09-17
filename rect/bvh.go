// Copyright 2018 Brian Noyama. Subject to the the Apache License, Version 2.0.
package rect

import (
	"math"
	"sort"
	"strings"

	"golang.org/x/exp/constraints"
)

func GenericAbs[T constraints.Signed | constraints.Float](val T) T {
	if val < 0 {
		return -val
	}
	return val
}

// A Bounding Volume for orthotopes. Wraps the orthotope and contains descendents.
type BVol[T Num] struct {
	vol   *Orthotope[T]
	desc  [2]*BVol[T]
	depth T
}

func (bvol *BVol[T]) minBound() {
	if bvol.depth > 0 {
		bvol.vol.MinBounds(bvol.desc[0].vol, bvol.desc[1].vol)
	}
}

func (bvol *BVol[T]) redepth() {
	bvol.depth = max(bvol.desc[0].depth, bvol.desc[1].depth) + 1
}

type byDimension[T Num] struct {
	orths     []*Orthotope[T]
	dimension int
}

func (d byDimension[T]) Len() int {
	return len(d.orths)
}

func (d byDimension[T]) Swap(i, j int) {
	d.orths[i], d.orths[j] = d.orths[j], d.orths[i]
}

// Compare the midpoints along a dimension.
func (d byDimension[T]) Less(i, j int) bool {
	return (d.orths[i].Point[d.dimension] +
		d.orths[i].Delta[d.dimension]) <
		(d.orths[j].Point[d.dimension] +
			d.orths[j].Delta[d.dimension])
}

// Creates a balanced BVH by recursively halving, sorting and comparing vols.
func TopDownBVH[T Num](orths []*Orthotope[T]) *BVol[T] {
	if len(orths) == 1 {
		return &BVol[T]{vol: orths[0]}
	}
	comp1 := &Orthotope[T]{}
	comp2 := &Orthotope[T]{}
	mid := len(orths) / 2

	lowDim := 0
	lowScore := T(math.MaxInt32)
	for d := 0; d < DIMENSIONS; d++ {
		sort.Sort(byDimension[T]{orths: orths, dimension: d})
		comp1.MinBounds(orths[:mid]...)
		comp2.MinBounds(orths[mid:]...)
		score := comp1.Score() + comp2.Score()
		if score < lowScore {
			lowScore = score
			lowDim = d
		}
	}
	if lowDim < DIMENSIONS-1 {
		sort.Sort(byDimension[T]{orths: orths, dimension: lowDim})
	}
	bvol := &BVol[T]{vol: comp1,
		desc: [2]*BVol[T]{TopDownBVH(orths[:mid]), TopDownBVH(orths[mid:])}}
	bvol.redepth()
	bvol.minBound()
	return bvol
}

func (bvol *BVol[T]) GetDepth() T {
	return bvol.depth
}

// Get an iterator for each volume in a Bounding Volume Hierarhcy.
func (bvol *BVol[T]) Iterator() *orthStack[T] {
	stack := &orthStack[T]{bvh: bvol, bvStack: []*BVol[T]{bvol}, intStack: []int{0}}
	return stack
}

// Add an orthotope to a Bounding Volume Hierarchy. Only add to root volume.
func (bvol *BVol[T]) Add(orth *Orthotope[T]) bool {
	s := bvol.Iterator()
	return s.Add(orth)
}

func (bvol *BVol[T]) Remove(orth *Orthotope[T]) bool {
	s := bvol.Iterator()
	return s.Remove(orth)
}

func (bvol *BVol[T]) Score() T {
	s := bvol.Iterator()
	return s.Score()
}

// SAH is a surface area heuristic as defined by MacDonald and Booth, 1990
// (https://doi.org/10.1007/BF01911006). This is an estimate of the overall tree
// quality.
func (bvol *BVol[T]) SAH() float64 {
	return bvol.Iterator().SAH(1.0, 1.2, 0)
}

// Rebalances the children of a given volume.
func (bvol *BVol[T]) redistribute() {
	if bvol.desc[1].depth > bvol.desc[0].depth {
		swapCheck(bvol.desc[1], bvol, 0)
	} else if bvol.desc[1].depth < bvol.desc[0].depth {
		swapCheck(bvol.desc[0], bvol, 1)
	} else if bvol.desc[1].depth > 0 {
		swapCheck(bvol.desc[0], bvol.desc[1], 1)
	}
	bvol.redepth()
}

func swapCheck[T Num](first *BVol[T], second *BVol[T], secIndex int) {
	first.minBound()
	second.minBound()
	minScore := first.vol.Score() + second.vol.Score()
	minIndex := -1

	for index := 0; index < 2; index++ {
		first.desc[index], second.desc[secIndex] =
			second.desc[secIndex], first.desc[index]

		// Ensure that swap did not unbalance second.
		if GenericAbs(second.desc[0].depth-second.desc[1].depth) < 2 {
			// Score first then second, since first may be a child of second.
			first.minBound()
			second.minBound()
			score := first.vol.Score() + second.vol.Score()
			if score < minScore {
				// Update the children with the best split
				minScore = score
				minIndex = index
			}
		}
	}

	if minIndex < 1 {
		first.desc[minIndex+1], second.desc[secIndex] =
			second.desc[secIndex], first.desc[minIndex+1]

		// Recalculate bounding volume
		first.minBound()
		second.minBound()
	}

	// Recalculate depth
	first.redepth()
	second.redepth()
}

// Recursive algorithm for comparing BVHs
func (bvh *BVol[T]) Equals(other *BVol[T]) bool {
	return (bvh.depth == 0 && other.depth == 0 && bvh.vol == other.vol) ||
		(bvh.depth > 0 && other.depth > 0 && bvh.vol.Equals(other.vol) &&
			((bvh.desc[0].Equals(other.desc[0]) && bvh.desc[1].Equals(other.desc[1])) ||
				(bvh.desc[1].Equals(other.desc[0]) && bvh.desc[0].Equals(other.desc[1]))))

}

// An indented string representation of the BVH (helps for debugging)
func (bvh *BVol[T]) String() string {
	iter := bvh.Iterator()
	maxDepth := bvh.depth
	toPrint := []string{}

	for iter.HasNext() {
		next := iter.Next()
		toPrint = append(toPrint, strings.Repeat(" ", int(maxDepth-next.depth)))
		toPrint = append(toPrint, next.vol.String()+"\n")
	}

	return strings.Join(toPrint, "")
}
