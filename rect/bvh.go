package rect

import (
	"math"

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
	bvol.orth.MinBounds(bvol.vol[0].orth, bvol.vol[1].orth)
}

func (bvol *BVol) redepth() {
	bvol.depth = disc.Max(bvol.vol[0].depth, bvol.vol[1].depth) + 1
}

// Get an iterator for each volume in a Bounding Volume Hierarhcy.
func (bvol *BVol) Iterator() *orthStack {
	stack := &orthStack{bvStack: []*BVol{bvol}, intStack: []int{0}}
	return stack
}

// Add an orthotope to a Bounding Volume Hierarchy. Only add to root volume.
func (bvol *BVol) Add(orth *Orthotope) {
	comp := Orthotope{}
	lowIndex := -1
	s := &orthStack{bvStack: []*BVol{}, intStack: []int{}}

	for next := bvol; next.orth != orth; next = next.vol[lowIndex] {
		if next.depth == 0 {
			// We've reached a leaf node, and we need to insert a parent node.
			next.vol[0] = &BVol{orth: orth}
			next.vol[1] = &BVol{orth: next.orth}
			comp = *next.orth
			next.orth = &comp
			lowIndex = 0
		} else {
			// We cannot add the orthotope here. Descend.
			smallestScore := math.MaxInt32

			for index, vol := range next.vol {
				comp.MinBounds(orth, vol.orth)
				score := comp.Score()
				if score < smallestScore {
					lowIndex = index
					smallestScore = score
				}
			}
		}
		s.append(next, lowIndex)
	}
	// Orthotope has been added, but tree needs to be rebalanced.
	s.rebalanceAdd(orth)
}

func (bvol *BVol) redistribute() {
	// If the depth is 0 then we need to do a tri swap.
	if bvol.vol[0].depth == 0 {
		swapCheck(bvol.vol[1], bvol, 0)
	} else if bvol.vol[1].depth == 0 {
		swapCheck(bvol.vol[0], bvol, 1)
	} else {
		swapCheck(bvol.vol[0], bvol.vol[1], 1)
	}
}

func swapCheck(first *BVol, second *BVol, secIndex int) {
	minScore := first.orth.Score() + second.orth.Score()
	minIndex := -1

	for index := 0; index < 2; index++ {
		first.vol[index].orth, second.vol[secIndex].orth =
			second.vol[secIndex].orth, first.vol[index].orth
		//Score first then second, since first may be a child of second.
		first.minBound()
		second.minBound()
		score := first.orth.Score() + second.orth.Score()
		if score < minScore {
			// Update the children with the best split
			minScore = score
			minIndex = index
		}
	}

	if minIndex < 1 {
		first.vol[minIndex+1].orth, second.vol[secIndex].orth =
			second.vol[secIndex].orth, first.vol[minIndex+1].orth
		first.minBound()
		second.minBound()
		first.depth = disc.Max(first.vol[0].depth, first.vol[1].depth) + 1
		first.depth = disc.Max(first.vol[0].depth, first.vol[1].depth) + 1
	}

}

/*
//Remove an orthtope from the BVH. Only remove from root volume.
func (bvh *BVOrth) Remove(orth *Orthotope) bool {
	low_index := -1
	var stack *orthStack

	for next := bvh; next.orth != orth; next = next.desc[low_index] {
		if next.orth.Contains(orth) && next.desc[0] != nil {
			low_index = 0
		} else {
			next, low_index = stack.pop()
			for low_index == 1 {
				if !stack.hasNext() {
					return false
				}

				next, low_index = stack.pop()
			}
			low_index = 1
		}
		stack.append(next, low_index)
	}

	stack.pop()
	parent, index := stack.peek()
	parent.orth = parent.desc[index^1].orth
	parent.desc[0], parent.desc[1] = nil, nil
	rebalance1(stack)
	return true
}

//Rebalance the tree after a delete using the path to the deleted volume.
func rebalance1(stack *orthStack) {
	for stack.hasNext() {
		parent, _ := stack.pop()
		gParent, gIndex := stack.peek()
		gParent.rebound()
		if parent.depth == gParent.depth-3 {
			//One node guaranteed can be swapped
			gParent.sahReduce1(gIndex)
		} else if gParent.desc[gIndex^1].depth == gParent.depth-2 {
			gParent.depth -= 1
		} else {
			gParent.sahReduce2(gIndex)
		}
	}
}

//Swap the optimal node for rebalancing a tree after a delete.
//Returns true if gParent's depth is changed.
func (bvh *BVOrth) sahReduce1(parent int) {
	sibling := bvh.desc[parent^1]
	comp := bvh.orth.DimensionCopy()
	minSA := math.MinInt32
	minIndex := -1
	for index, niece := range sibling.desc {
		//If the niece has a difference of 2 it cannot be swapped.
		if niece.depth == bvh.desc[parent].depth-2 {
			bvh.desc[parent].depth -= 1
			bvh.depth -= 1
			minIndex = index ^ 1
			break
		} else {
			comp.MinBounds(bvh.desc[parent].orth, niece.orth)
			sa := comp.Score()
			if sa < minSA {
				minSA = sa
				minIndex = index
			}
		}
	}
	bvh.desc[parent], sibling.desc[minIndex] = sibling.desc[minIndex],
		bvh.desc[parent]
}

//Swap a node without rebalancing after a delete.
func (bvh *BVOrth) sahReduce2(parent int) {
	sibling := bvh.desc[parent^1]
	comp := bvh.orth.DimensionCopy()
	minSA := sibling.orth.Score()
	minIndex := -1
	for index, niece := range sibling.desc {
		comp.MinBounds(bvh.desc[parent].orth, niece.orth)
		sa := comp.Score()
		if sa < minSA {
			minSA = sa
			minIndex = index
		}
	}
	if minIndex != -1 {
		bvh.desc[parent], sibling.desc[minIndex] = sibling.desc[minIndex],
			bvh.desc[parent]
	}
}

func (bvh *BVOrth) rebound() {
	bvh.orth.MinBounds(bvh.desc[0].orth, bvh.desc[1].orth)
}

//Rebalance the tree after an add by following the path to the added volume.
func rebalance0(stack *orthStack, orth *Orthotope) {
	parent, index := stack.pop()
	for stack.hasNext() {
		gParent, gIndex := stack.peek()
		//newDepth := parent.desc[index].depth
		gParent.orth.MinBounds(gParent.orth, orth)
		if gParent.desc[gIndex^1].depth+1 == gParent.depth &&
			gParent.depth != parent.depth+1 {
			//No Sa reduction needed here. Add ensures that volume is smallest.
			gParent.depth = parent.depth + 1
		} else if gParent.depth != parent.depth+1 {
			//Swap to balance the tree
			gParent.desc[gIndex^1], parent.desc[index] = parent.desc[index],
				gParent.desc[gIndex^1]
			parent.depth = gParent.depth - 1
			//Compare grand-children to make child volumes smaller.
			gParent.sahReduce0()
		} else {
			//Compare grand-children to make child volumes smaller.
			gParent.sahReduce0()
		}
		parent, index = stack.pop()
	}
}

//Swap grandchild nodes for when a volume is added.
func (bvh *BVOrth) sahReduce0() {
	sibling := bvh.desc
	comp := bvh.orth.DimensionCopy()

	sibling[0].rebound()
	sibling[1].rebound()
	minSA := sibling[0].orth.Score() + sibling[1].orth.Score()

	minIndex := -1
	for index, niece := range sibling[1].desc {
		comp.MinBounds(sibling[0].desc[0].orth, niece.orth)
		sa := comp.Score()
		comp.MinBounds(sibling[0].desc[1].orth, sibling[1].desc[index^1].orth)
		sa += comp.Score()
		if sa < minSA {
			minSA = sa
			minIndex = index
		}
	}
	if minIndex != -1 {
		sibling[0].desc[minIndex].orth, sibling[1].desc[minIndex^1].orth =
			sibling[1].desc[minIndex^1].orth, sibling[0].desc[minIndex].orth
	}
	//Ensure siblings get volume updated.
	sibling[0].rebound()
	sibling[1].rebound()
}

//Return the group representation for equivalent BVHs
func (bvh *BVOrth) Group() {
	stack := []*BVOrth{bvh}

	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if cur.desc[0] != nil {
			if cur.desc[0].orth.Compare(cur.desc[1].orth) > 0 {
				// Swap if out of order.
				cur.desc[0], cur.desc[1] = cur.desc[1], cur.desc[0]
			}
			// Add both BVHs to a stack
			stack = append(stack, cur.desc[0])
			stack = append(stack, cur.desc[1])
		}
	}
}

//Print and indented string representation of the BVH
func (bvh *BVOrth) String() string {
	iter := bvh.Iterator()
	maxDepth := bvh.depth
	toPrint := []string{}

	for next := iter.Next(); next != nil; next = iter.Next() {
		toPrint = append(toPrint, strings.Repeat(" ", int(maxDepth-next.depth)))
		toPrint = append(toPrint, next.orth.String()+"\n")
	}

	return strings.Join(toPrint, "")
}
*/
