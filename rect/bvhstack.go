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
	bvol, index := s.peek()
	bvolPrev := bvol

	// Trace up tree
	for bvol.desc == nil || bvol.desc.len <= index {
		s.pop()

		// The end of the stack.
		if !s.HasNext() {
			return bvolPrev
		}

		// Check out next child.
		s.intStack[len(s.intStack)-1]++
		bvol, index = s.peek()
	}

	bvol = bvol.desc.vol[index]
	s.append(bvol, 0)
	return bvolPrev
}

/*
 * Trace performs ray tracing on the BVH returning an orthotope.
 * Its depth in the tree, and the distance from the beginning of the vector, o,
 * passed in.
 */
func (s *orthStack) Trace(o *Orthotope) (*Orthotope, int, int) {
	bvol, distance := s.pop()
	var depth int

	if bvol.desc != nil {
		// Get the depth to return.
		depth = bvol.desc.depth

		// Find the distances for each child, if there's a collision.
		distances := []int{}
		volumes := []*BVol{}
		for i := 0; i < bvol.desc.len; i++ {
			distance := o.Intersects(bvol.desc.vol[i].orth)
			if distance >= 0 {
				distances = append(distances, distance)
				volumes = append(volumes, bvol.desc.vol[i])

				// Insertion sort, to return values that are closest first.
				for j := i; j >= 1; j-- {
					if distances[j] > distances[j-1] {
						distances[j], distances[j-1] = distances[j-1], distances[j]
						volumes[j], volumes[j-1] = volumes[j-1], volumes[j]
					}
				}
			}
		}

		for index, vol := range volumes {
			s.append(vol, distances[index])
		}
	}

	return bvol.orth, depth, distance
}

/*
 * Query looks for intersections between the orthotope, o, and the BVH
 * returning one intersection at a time.
 */
func (s *orthStack) Query(o *Orthotope) *Orthotope {
	bvol, index := s.peek()

	for bvol.desc != nil {
		bvol, index = bvol.desc.vol[index], 0
		s.append(bvol, index)
	}

	bvolPrev := bvol

	// Trace up tree
	for bvol.desc == nil || bvol.desc.len <= index {
		s.pop()

		// The end of the stack.
		if !s.HasNext() {
			return bvolPrev.orth
		}

		// Check out next child.
		s.intStack[len(s.intStack)-1]++
		bvol, index = s.peek()
	}

	// Trace down the tree
	for bvol.desc != nil {
		bvol, index = bvol.desc.vol[index], 0
		s.append(bvol, index)
	}

	return bvolPrev.orth
}

// Attempt rebalancing when the depth of the tree has potentially changed.
// Hypothetically, have one rebalance method to rule them all.
func (s *orthStack) rebalanceAdd(o *Orthotope) {
	gParent, gIndex := s.pop()
	last_depth := 0
	for s.HasNext() {
		parent, pIndex := gParent, gIndex
		gParent, gIndex = s.pop()
		for index, aunt := range gParent.desc.Slice() {
			auntDepth := 0
			if aunt.desc != nil {
				auntDepth = aunt.desc.depth + 1
			}
			if auntDepth < last_depth {
				// parent depth is 2 above, so swap.
				parent.desc.vol[pIndex], gParent.desc.vol[index] =
					gParent.desc.vol[index], parent.desc.vol[pIndex]

				gParent.redistribute(gIndex, index)
				last_depth--
				break
			}
		}
		parent.minBounds()
		parent.desc.depth = last_depth
		last_depth++
	}
	gParent.minBounds()
}
