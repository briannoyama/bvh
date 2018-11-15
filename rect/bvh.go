package rect

import (
	"strings"

	disc "github.com/briannoyama/bvh/discreet"
)

//"math"
//"strings"

//A Bounding Volume for orthotopes. Wraps the orthotope and .
type BVol struct {
	orth  *Orthotope
	vol   [2]*BVol
	depth int
}

func (bvol *BVol) minBound() {
	if bvol.depth > 0 {
		bvol.orth.MinBounds(bvol.vol[0].orth, bvol.vol[1].orth)
	}
}

func (bvol *BVol) redepth() {
	bvol.depth = disc.Max(bvol.vol[0].depth, bvol.vol[1].depth) + 1
}

// Get an iterator for each volume in a Bounding Volume Hierarhcy.
func (bvol *BVol) Iterator() *orthStack {
	stack := &orthStack{bvh: bvol, bvStack: []*BVol{bvol}, intStack: []int{0}}
	return stack
}

// Add an orthotope to a Bounding Volume Hierarchy. Only add to root volume.
func (bvol *BVol) Add(orth *Orthotope) bool {
	s := bvol.Iterator()
	return s.Add(orth)
}

func (bvol *BVol) Remove(orth *Orthotope) bool {
	s := bvol.Iterator()
	return s.Remove(orth)
}

func (bvol *BVol) Score() int {
	s := bvol.Iterator()
	return s.Score()
}

// Rebalances the children of a given volume.
func (bvol *BVol) redistribute() {
	if bvol.vol[1].depth > bvol.vol[0].depth {
		swapCheck(bvol.vol[1], bvol, 0)
	} else if bvol.vol[1].depth < bvol.vol[0].depth {
		swapCheck(bvol.vol[0], bvol, 1)
	} else if bvol.vol[1].depth > 0 {
		swapCheck(bvol.vol[0], bvol.vol[1], 1)
	}
	bvol.redepth()
}

func swapCheck(first *BVol, second *BVol, secIndex int) {
	first.minBound()
	second.minBound()
	minScore := first.orth.Score() + second.orth.Score()
	minIndex := -1

	for index := 0; index < 2; index++ {
		first.vol[index], second.vol[secIndex] =
			second.vol[secIndex], first.vol[index]

			// Ensure that swap did not unbalance second.
		if disc.Abs(second.vol[0].depth-second.vol[1].depth) < 2 {
			// Score first then second, since first may be a child of second.
			first.minBound()
			second.minBound()
			score := first.orth.Score() + second.orth.Score()
			if score < minScore {
				// Update the children with the best split
				minScore = score
				minIndex = index
			}
		}
	}

	if minIndex < 1 {
		first.vol[minIndex+1], second.vol[secIndex] =
			second.vol[secIndex], first.vol[minIndex+1]

		// Recalculate bounding volume
		first.minBound()
		second.minBound()
	}

	// Recalculate depth
	first.redepth()
	second.redepth()
}

//Recursive algorithm for comparing BVHs
func (bvh *BVol) Equals(other *BVol) bool {
	return (bvh.depth == 0 && other.depth == 0 && bvh.orth == other.orth) ||
		(bvh.depth > 0 && other.depth > 0 && bvh.orth.Equals(other.orth) &&
			((bvh.vol[0].Equals(other.vol[0]) && bvh.vol[1].Equals(other.vol[1])) ||
				(bvh.vol[1].Equals(other.vol[0]) && bvh.vol[0].Equals(other.vol[1]))))

}

// An indented string representation of the BVH (helps for debugging)
func (bvh *BVol) String() string {
	iter := bvh.Iterator()
	maxDepth := bvh.depth
	toPrint := []string{}

	for iter.HasNext() {
		next := iter.Next()
		toPrint = append(toPrint, strings.Repeat(" ", int(maxDepth-next.depth)))
		toPrint = append(toPrint, next.orth.String()+"\n")
	}

	return strings.Join(toPrint, "")
}
