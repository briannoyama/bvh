package rect

type orthStack struct {
	bvStack  []*BVol
	intStack []int
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

func (s *orthStack) popIfExists() (*BVol, int) {
	if len(s.bvStack) == 0 {
		return nil, -1
	}
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

func (s *orthStack) upNext() bool {
	s.pop()

	// The end of the stack.
	if s.HasNext() {
		s.intStack[len(s.intStack)-1]++
		return true
	}
	return false
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

/* Attempt rebalancing when the depth of the tree has changed.
 * Hypothetically, have one rebalance method to rule them all.
 */
func (s *orthStack) rebalanceAdd(o *Orthotope) {
	gParent, gIndex := s.pop()
	lastDepth := 1
	for s.HasNext() {
		parent, pIndex := gParent, gIndex
		gParent, gIndex = s.pop()

		aIndex := gIndex ^ 1
		aunt := gParent.vol[aIndex]
		if aunt.depth+1 < lastDepth {
			// parent depth is 2 above, so swap.
			parent.vol[pIndex], gParent.vol[aIndex] =
				gParent.vol[aIndex], parent.vol[pIndex]

			lastDepth--
		}

		parent.depth = lastDepth
		lastDepth++
		gParent.redistribute()
		// Investigate another swap
	}
	gParent.depth = lastDepth
	gParent.minBound()
}
