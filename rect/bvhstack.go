// Copyright 2018 Brian Noyama. Subject to the the Apache License, Version 2.0.
package rect

import (
	"math"
)


// OrthStack gives methods for working with Orthotope BVol.
type OrthStack[T Num] interface {
	Reset()
	HasNext() bool
	Next() *BVol[T]
	Trace(o *Orthotope[T]) (*Orthotope[T], T)
	Query(o *Orthotope[T]) *Orthotope[T]
	Add(orth *Orthotope[T]) bool
	Contains(orth *Orthotope[T]) bool
	Remove(o *Orthotope[T]) bool
}

type orthStack[T Num] struct {
	bvh      *BVol[T]
	bvStack  []*BVol[T]
	intStack []int
}

// Resets the stack.
func (s *orthStack[T]) Reset() {
	s.intStack = s.intStack[:0]
	s.bvStack = s.bvStack[:0]
	s.bvStack = append(s.bvStack, s.bvh)
	s.intStack = append(s.intStack, 0)
}

func (s *orthStack[T]) HasNext() bool {
	return len(s.bvStack) > 0
}

func (s *orthStack[T]) append(bvol *BVol[T], index int) {
	s.bvStack = append(s.bvStack, bvol)
	s.intStack = append(s.intStack, index)
}

func (s *orthStack[T]) peek() (*BVol[T], T) {
	return s.bvStack[len(s.bvStack)-1], T(s.intStack[len(s.intStack)-1])
}

func (s *orthStack[T]) pop() (*BVol[T], T) {
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
func (s *orthStack[T]) Next() *BVol[T] {
	bvolPrev, _ := s.peek()

	if s.traceUp() {
		bvol, index := s.peek()
		bvol = bvol.desc[int(index)]
		s.append(bvol, 0)
	}

	return bvolPrev
}

/*
 * Trace performs ray tracing on the BVH returning an orthotope.
 * Its depth in the tree, and the distance from the beginning of the vector, o,
 * passed in.
 */
func (s *orthStack[T]) Trace(o *Orthotope[T]) (*Orthotope[T], T) {
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
					s.append(bvol.desc[0], int(distance0))
					bvol, distance = bvol.desc[1], distance1
				} else {
					s.append(bvol.desc[1], int(distance1))
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
func (s *orthStack[T]) traceUp() bool {
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

func (s *orthStack[T]) queryNext(o *Orthotope[T]) *BVol[T] {
	bvol, index := s.peek()
	for bvol.depth > 0 {
		if index >= 2 {
			if !s.traceUp() {
				break
			}
		} else {
			if bvol.desc[int(index)].vol.Overlaps(o) {
				s.append(bvol.desc[int(index)], 0)
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
func (s *orthStack[T]) Query(o *Orthotope[T]) *Orthotope[T] {
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

func (s *orthStack[T]) path(o *Orthotope[T]) *BVol[T] {
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
				if bvol.desc[int(index)].vol.Contains(o) {
					s.append(bvol.desc[int(index)], 0)
				} else {
					s.intStack[len(s.intStack)-1]++
				}
			}
			bvol, index = s.peek()
		}
	}
	return bvol
}

func (s *orthStack[T]) Contains(o *Orthotope[T]) bool {
	s.Reset()
	bvol := s.path(o)

	// Check that the orthotope is the last thing from the path.
	return o == bvol.vol
}

// Add an orthotope to a Bounding Volume Hierarchy. Only add to root volume.
func (s *orthStack[T]) Add(orth *Orthotope[T]) bool {
	s.Reset()
	bvol := s.bvh
	if bvol.vol == nil {
		// Add by setting the vol when there is no volumes.
		bvol.vol = orth
	}
	comp := Orthotope[T]{}
	lowIndex := -1

	for next := bvol; next.vol != orth; next = next.desc[lowIndex] {
		if next.depth == 0 {
			// We've reached a leaf node, and we need to insert a parent node.
			next.desc[0] = &BVol[T]{vol: orth}
			next.desc[1] = &BVol[T]{vol: next.vol}
			next.depth = 1
			comp = *next.vol
			next.vol = &comp
			lowIndex = 0
		} else {
			// We cannot add the orthotope here. Descend.
			smallestScore := T(math.MaxInt32)

			for index, vol := range next.desc {
				comp.MinBounds(orth, vol.vol)

				if vol.vol == orth {
					// The volume has already been added.
					return false
				}

				score := comp.Score() - vol.vol.Score()
				if score < smallestScore {
					lowIndex = index
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
func (s *orthStack[T]) Remove(o *Orthotope[T]) bool {
	s.Reset()
	bvol := s.path(o)
	if o == bvol.vol {
		s.pop()
		if s.HasNext() {
			parent, pIndex := s.pop()
			if s.HasNext() {
				gParent, gIndex := s.peek()
				// Delete the node by replacing the parent.
				gParent.desc[int(gIndex)] = parent.desc[int(pIndex)^1]
				s.rebalanceRemove()
			} else {
				// Delete the node by replacing the volume and children with cousin.
				cousin := parent.desc[int(pIndex)^1]
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
func (s *orthStack[T]) Score() T {
	s.Reset()
	score := T(0)

	for s.HasNext() {
		score += s.Next().vol.Score()
	}
	return score
}

func (s *orthStack[T]) SAH(cInternal, cLeaves, cOverlap float64) float64 {
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
func (s *orthStack[T]) rebalanceAdd() {
	gParent, gIndex := s.pop()
	for s.HasNext() {
		parent, pIndex := gParent, gIndex
		gParent, gIndex = s.pop()

		aIndex := int(gIndex) ^ 1

		if gParent.desc[aIndex].depth < parent.desc[int(pIndex)].depth {
			// Swap to fix balance.
			parent.desc[int(pIndex)], gParent.desc[aIndex] =
				gParent.desc[aIndex], parent.desc[int(pIndex)]
			parent.redepth()
		}
		gParent.redistribute()
		// Found that gParent was not consistently getting minBound after redistribute.
		gParent.minBound()
	}
	gParent.minBound()
}

// Attempt rebalancing when the depth of the tree has potentially decreased.
func (s *orthStack[T]) rebalanceRemove() {
	for s.HasNext() {
		parent, pIndex := s.pop()

		cIndex := int(pIndex) ^ 1
		cousin := parent.desc[cIndex]
		depth := parent.desc[int(pIndex)].depth

		if cousin.depth > depth+1 {
			swap := 0
			// Swap to fix balance. Try to minimize hierarchy with swap.
			if cousin.desc[1].depth == depth+1 {
				if cousin.desc[0].depth == depth+1 {
					cousin.vol.MinBounds(cousin.desc[1].vol, parent.desc[int(pIndex)].vol)
					score := cousin.vol.Score() - cousin.desc[1].vol.Score()
					cousin.vol.MinBounds(cousin.desc[0].vol, parent.desc[int(pIndex)].vol)
					if score < cousin.vol.Score()-cousin.desc[0].vol.Score() {
						swap = 1
					}
				} else {
					swap = 1
				}
			}
			parent.desc[int(pIndex)], cousin.desc[swap] =
				cousin.desc[swap], parent.desc[int(pIndex)]
			cousin.redepth()
			cousin.minBound()
		}
		parent.minBound()
		parent.redistribute()
	}
}
