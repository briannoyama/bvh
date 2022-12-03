// Copyright 2018 Brian Noyama. Subject to the the Apache License, Version 2.0.
package rect

import (
	"math"
)

// OrthStack gives methods for working with Orthotope BVol.
type OrthStack interface {
	Reset()
	HasNext() bool
	Next() *BVol
	Trace(o *Orthotope) (*Orthotope, int32)
	Query(o *Orthotope) *Orthotope
	Add(orth *Orthotope) bool
	Contains(orth *Orthotope) bool
	Remove(o *Orthotope) bool
}

type orthStack struct {
	bvh      *BVol
	bvStack  []*BVol
	intStack []int32
}

// Resets the stack.
func (s *orthStack) Reset() {
	s.intStack = s.intStack[:0]
	s.bvStack = s.bvStack[:0]
	s.bvStack = append(s.bvStack, s.bvh)
	s.intStack = append(s.intStack, 0)
}

func (s *orthStack) HasNext() bool {
	return len(s.bvStack) > 0
}

func (s *orthStack) append(bvol *BVol, index int32) {
	s.bvStack = append(s.bvStack, bvol)
	s.intStack = append(s.intStack, index)
}

func (s *orthStack) peek() (*BVol, int32) {
	return s.bvStack[len(s.bvStack)-1], s.intStack[len(s.intStack)-1]
}

func (s *orthStack) pop() (*BVol, int32) {
	bvol, index := s.peek()
	s.bvStack = s.bvStack[:len(s.bvStack)-1]
	s.intStack = s.intStack[:len(s.intStack)-1]
	return bvol, index
}

/*
 * Iterates through the tree by modifying the stack in place. The stack will be
 * organized such that peek reflects the next value that will be returned.
 * In this way, next pops off an element while traversing the tree in pre-order.
 */
func (s *orthStack) Next() *BVol {
	bvolPrev, _ := s.peek()

	if s.traceUp() {
		bvol, index := s.peek()
		bvol = bvol.desc[index]
		s.append(bvol, 0)
	}

	return bvolPrev
}

/*
 * Trace performs ray tracing on the BVH returning an orthotope.
 * Its depth in the tree, and the distance from the beginning of the vector, o,
 * passed in.
 */
func (s *orthStack) Trace(o *Orthotope) (*Orthotope, int32) {
	if !s.HasNext() {
		return nil, -1
	}
	bvol, distance := s.pop()

	for bvol.depth > 0 {

		// Find the distances for each child, if there's a collision.
		distance0 := o.Intersects(bvol.desc[0].vol)
		distance1 := o.Intersects(bvol.desc[1].vol)

		if distance0 >= 0 {
			if distance1 >= 0 {
				if distance1 < distance0 {
					s.append(bvol.desc[0], distance0)
					bvol, distance = bvol.desc[1], distance1
				} else {
					s.append(bvol.desc[1], distance1)
					bvol, distance = bvol.desc[0], distance0
				}
			} else {
				bvol, distance = bvol.desc[0], distance0
			}
		} else if distance1 >= 0 {
			bvol, distance = bvol.desc[1], distance1
		} else if s.HasNext() {
			bvol, distance = s.pop()
		} else {
			return nil, -1
		}
	}

	return bvol.vol, distance
}

/* Goes up the tree until it finds the next unvisited child index, after
 * looking at parents.
 */
func (s *orthStack) traceUp() bool {
	bvol, index := s.peek()
	for bvol.depth == 0 || index >= 2 {
		s.pop()

		// The end of the stack.
		if !s.HasNext() {
			return false
		}

		// Check out next child.
		s.intStack[len(s.intStack)-1]++
		bvol, index = s.peek()
	}
	return true
}

func (s *orthStack) queryNext(o *Orthotope) *BVol {
	bvol, index := s.peek()
	for bvol.depth > 0 {
		if index >= 2 {
			if !s.traceUp() {
				break
			}
		} else {
			if bvol.desc[index].vol.Overlaps(o) {
				s.append(bvol.desc[index], 0)
			} else {
				s.intStack[len(s.intStack)-1]++
			}
		}
		bvol, index = s.peek()
	}
	return bvol
}

/*
 * Query looks for intersections between the orthotope, o, and the BVH
 * returning one intersection at a time.
 */
func (s *orthStack) Query(o *Orthotope) *Orthotope {
	// When the stack is empty, there are no more volumes to return.
	if !s.HasNext() {
		return nil
	}
	bvol := s.queryNext(o)
	if !s.HasNext() {
		return nil
	}

	// Use trace up to get the next possible branch.
	if s.traceUp() {
		//s.queryNext(o)
	}
	return bvol.vol
}

func (s *orthStack) path(o *Orthotope) *BVol {
	bvol, index := s.peek()
	for bvol.vol != o && s.HasNext() {
		if bvol.depth == 0 {
			if !s.traceUp() {
				break
			}
			bvol, index = s.peek()
		}
		for bvol.depth > 0 {
			if index >= 2 {
				if !s.traceUp() {
					break
				}
			} else {
				if bvol.desc[index].vol.Contains(o) {
					s.append(bvol.desc[index], 0)
				} else {
					s.intStack[len(s.intStack)-1]++
				}
			}
			bvol, index = s.peek()
		}
	}
	return bvol
}

func (s *orthStack) Contains(o *Orthotope) bool {
	s.Reset()
	bvol := s.path(o)

	// Check that the orthotope is the last thing from the path.
	return o == bvol.vol
}

// Add an orthotope to a Bounding Volume Hierarchy. Only add to root volume.
func (s *orthStack) Add(orth *Orthotope) bool {
	s.Reset()
	bvol := s.bvh
	if bvol.vol == nil {
		// Add by setting the vol when there is no volumes.
		bvol.vol = orth
	}
	comp := Orthotope{}
	lowIndex := int32(-1)

	for next := bvol; next.vol != orth; next = next.desc[lowIndex] {
		if next.depth == 0 {
			// We've reached a leaf node, and we need to insert a parent node.
			next.desc[0] = &BVol{vol: orth}
			next.desc[1] = &BVol{vol: next.vol}
			next.depth = 1
			comp = *next.vol
			next.vol = &comp
			lowIndex = int32(0)
		} else {
			// We cannot add the orthotope here. Descend.
			smallestScore := int32(math.MaxInt32)

			for index, vol := range next.desc {
				comp.MinBounds(orth, vol.vol)

				if vol.vol == orth {
					// The volume has already been added.
					return false
				}

				score := comp.Score()
				if score < smallestScore {
					lowIndex = int32(index)
					smallestScore = score
				}
			}
		}
		s.append(next, lowIndex)
	}
	// Orthotope has been added, but tree needs to be rebalanced.

	s.rebalanceAdd()
	return true
}

// Remove an orthotope from the BVH associated with this stack.
func (s *orthStack) Remove(o *Orthotope) bool {
	s.Reset()
	bvol := s.path(o)
	if o == bvol.vol {
		s.pop()
		if s.HasNext() {
			parent, pIndex := s.pop()
			if s.HasNext() {
				gParent, gIndex := s.peek()
				// Delete the node by replacing the parent.
				gParent.desc[gIndex] = parent.desc[pIndex^1]
				s.rebalanceRemove()
			} else {
				// Delete the node by replacing the volume and children with cousin.
				cousin := parent.desc[pIndex^1]
				parent.vol = cousin.vol
				parent.desc = cousin.desc
				parent.depth = cousin.depth
			}
		} else {
			// For depths of 0, delete by removing the volume.
			bvol.vol = nil
		}
		return true
	}
	return false
}

// Returns the total score by using the volumes Score method for each volume.
func (s *orthStack) Score() int32 {
	s.Reset()
	score := int32(0)

	for s.HasNext() {
		score += s.Next().vol.Score()
	}
	return score
}

func (s *orthStack) SAH(cInternal, cLeaves, cOverlap float64) float64 {
	s.Reset()

	var ci, cl, co float64

	for s.HasNext() {
		n := s.Next()
		if n.depth == 0 {
			cl += float64(n.vol.SurfaceArea())
			// This BVH implementatino only has a single AABB in its
			// leaf node.
			co += float64(n.vol.SurfaceArea())
		} else {
			ci += float64(n.vol.SurfaceArea())
		}
	}
	return (cInternal*ci + cLeaves*cl + cOverlap*co) / float64(s.bvh.vol.SurfaceArea())
}

// Attempt rebalancing when the depth of the tree has potentially increased.
func (s *orthStack) rebalanceAdd() {
	gParent, gIndex := s.pop()
	for s.HasNext() {
		parent, pIndex := gParent, gIndex
		gParent, gIndex = s.pop()

		aIndex := gIndex ^ 1

		if gParent.desc[aIndex].depth < parent.desc[pIndex].depth {
			// Swap to fix balance.
			parent.desc[pIndex], gParent.desc[aIndex] =
				gParent.desc[aIndex], parent.desc[pIndex]
			parent.redepth()
		}
		gParent.redistribute()
		// Found that gParent was not consistently getting minBound after redistribute.
		gParent.minBound()
	}
	gParent.minBound()
}

// Attempt rebalancing when the depth of the tree has potentially decreased.
func (s *orthStack) rebalanceRemove() {
	for s.HasNext() {
		parent, pIndex := s.pop()

		cIndex := pIndex ^ 1
		cousin := parent.desc[cIndex]
		depth := parent.desc[pIndex].depth

		if cousin.depth > depth+1 {
			swap := 0
			// Swap to fix balance. Try to minimize hierarchy with swap.
			if cousin.desc[1].depth == depth+1 {
				if cousin.desc[0].depth == depth+1 {
					cousin.vol.MinBounds(cousin.desc[1].vol, parent.desc[pIndex].vol)
					score := cousin.vol.Score()
					cousin.vol.MinBounds(cousin.desc[0].vol, parent.desc[pIndex].vol)
					if score < cousin.vol.Score() {
						swap = 1
					}
				} else {
					swap = 1
				}
			}
			parent.desc[pIndex], cousin.desc[swap] =
				cousin.desc[swap], parent.desc[pIndex]
			cousin.redepth()
			cousin.minBound()
		}
		parent.minBound()
		parent.redistribute()
	}
}
