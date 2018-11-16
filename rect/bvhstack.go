package rect

import "math"

type orthStack struct {
	bvh      *BVol
	bvStack  []*BVol
	intStack []int
}

func (s *orthStack) Reset() {
	s.intStack = s.intStack[:0]
	s.bvStack = s.bvStack[:0]
	s.bvStack = append(s.bvStack, s.bvh)
	s.intStack = append(s.intStack, 0)
}

func (s *orthStack) HasNext() bool {
	return len(s.bvStack) > 0
}

func (s *orthStack) append(bvol *BVol, index int) {
	s.bvStack = append(s.bvStack, bvol)
	s.intStack = append(s.intStack, index)
}

func (s *orthStack) peek() (*BVol, int) {
	return s.bvStack[len(s.bvStack)-1], s.intStack[len(s.intStack)-1]
}

func (s *orthStack) pop() (*BVol, int) {
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
		bvol = bvol.vol[index]
		s.append(bvol, 0)
	}

	return bvolPrev
}

/*
 * Trace performs ray tracing on the BVH returning an orthotope.
 * Its depth in the tree, and the distance from the beginning of the vector, o,
 * passed in.
 */
func (s *orthStack) Trace(o *Orthotope) (*Orthotope, int, int) {
	bvol, distance := s.pop()

	if bvol.depth > 0 {

		// Find the distances for each child, if there's a collision.
		distance0 := o.Intersects(bvol.vol[0].orth)
		distance1 := o.Intersects(bvol.vol[1].orth)

		if distance0 >= 0 {
			if distance1 >= 0 {
				if distance1 < distance0 {
					s.append(bvol.vol[0], distance0)
					s.append(bvol.vol[1], distance1)
				} else {
					s.append(bvol.vol[1], distance1)
					s.append(bvol.vol[0], distance0)
				}
			} else {
				s.append(bvol.vol[0], distance0)
			}
		} else if distance1 >= 0 {
			s.append(bvol.vol[1], distance1)
		}
	}

	return bvol.orth, bvol.depth, distance
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
			if bvol.vol[index].orth.Overlaps(o) {
				s.append(bvol.vol[index], 0)
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
	bvol := s.queryNext(o)

	// Use trace up to get the next possible branch.
	if s.traceUp() {
		s.queryNext(o)
	}
	return bvol.orth
}

func (s *orthStack) path(o *Orthotope) *BVol {
	bvol, index := s.peek()
	for bvol.depth > 0 {
		if index >= 2 {
			if !s.traceUp() {
				break
			}
		} else {
			if bvol.vol[index].orth.Contains(o) {
				s.append(bvol.vol[index], 0)
			} else {
				s.intStack[len(s.intStack)-1]++
			}
		}
		bvol, index = s.peek()
	}
	return bvol
}

func (s *orthStack) Contains(o *Orthotope) bool {
	s.Reset()
	bvol := s.path(o)

	// Check that the orthotope is the last thing from the path.
	return o == bvol.orth
}

// Add an orthotope to a Bounding Volume Hierarchy. Only add to root volume.
func (s *orthStack) Add(orth *Orthotope) bool {
	s.Reset()
	bvol := s.bvh
	comp := Orthotope{}
	lowIndex := -1

	for next := bvol; next.orth != orth; next = next.vol[lowIndex] {
		if next.depth == 0 {
			// We've reached a leaf node, and we need to insert a parent node.
			next.vol[0] = &BVol{orth: orth}
			next.vol[1] = &BVol{orth: next.orth}
			next.depth = 1
			comp = *next.orth
			next.orth = &comp
			lowIndex = 0
		} else {
			// We cannot add the orthotope here. Descend.
			smallestScore := math.MaxInt32

			for index, vol := range next.vol {
				comp.MinBounds(orth, vol.orth)

				if vol.orth == orth {
					// The volume has already been added.
					return false
				}

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

	s.rebalanceAdd()
	return true
}

func (s *orthStack) Remove(o *Orthotope) bool {
	s.Reset()
	bvol := s.path(o)
	if o == bvol.orth {
		s.pop()
		if s.HasNext() {
			parent, pIndex := s.pop()
			if s.HasNext() {
				gParent, gIndex := s.peek()
				// Delete the node by replacing the parent.
				gParent.vol[gIndex] = parent.vol[pIndex^1]
				s.rebalanceRemove()
			} else {
				// Delete the node by replacing the volume and children with cousin.
				cousin := parent.vol[pIndex^1]
				parent.orth = cousin.orth
				parent.vol = cousin.vol
				parent.depth = cousin.depth
			}
			return true
		}
	}
	return false
}

func (s *orthStack) Score() int {
	s.Reset()
	score := 0

	for s.HasNext() {
		score += s.Next().orth.Score()
	}
	return score
}

/* Attempt rebalancing when the depth of the tree has changed.
 * Hypothetically, have one rebalance method to rule them all.
 */
func (s *orthStack) rebalanceAdd() {
	gParent, gIndex := s.pop()
	for s.HasNext() {
		parent, pIndex := gParent, gIndex
		gParent, gIndex = s.pop()

		aIndex := gIndex ^ 1

		if gParent.vol[aIndex].depth < parent.vol[pIndex].depth {
			// Swap to fix balance.
			parent.vol[pIndex], gParent.vol[aIndex] =
				gParent.vol[aIndex], parent.vol[pIndex]
			parent.redepth()
		}
		gParent.redistribute()
	}
	gParent.minBound()
}

func (s *orthStack) rebalanceRemove() {
	for s.HasNext() {
		parent, pIndex := s.pop()

		cIndex := pIndex ^ 1
		cousin := parent.vol[cIndex]
		depth := parent.vol[pIndex].depth

		if cousin.depth > depth+1 {
			swap := 0
			// Swap to fix balance. Try to minimize hierarchy with swap.
			if cousin.vol[1].depth == depth+1 {
				if cousin.vol[0].depth == depth+1 {
					cousin.orth.MinBounds(cousin.vol[1].orth, parent.vol[pIndex].orth)
					score := cousin.orth.Score()
					cousin.orth.MinBounds(cousin.vol[0].orth, parent.vol[pIndex].orth)
					if score < cousin.orth.Score() {
						swap = 1
					}
				} else {
					swap = 1
				}
			}
			parent.vol[pIndex], cousin.vol[swap] =
				cousin.vol[swap], parent.vol[pIndex]
			cousin.redepth()
			cousin.minBound()
		}
		//parent.redepth()
		parent.minBound()
		parent.redistribute()
	}
}
