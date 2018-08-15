package rect

type orthStack struct {
	bvStack  []*BVol
	intStack []int
}

func (s *orthStack) hasNext() bool {
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
 * organized such that the returned value from peek reflects the last value
 * returned by Next.
 */
func (s *orthStack) Next() *BVol {
	bvol, index := s.peek()

	if index == WIDTH {
		s.pop()
		s.intStack[len(s.intStack)-1] += 1
		bvol, index = s.peek()
	}

	if index != WIDTH {
		child := bvol.desc.vol[index]
		if child == nil {
			s.intStack[len(s.intStack)-1] = WIDTH
		} else {
			for child.desc != nil {
				s.append(child, 0)
				bvol = child
				child = bvol.desc.vol[index]
			}
			s.intStack[len(s.intStack)-1] = WIDTH
		}
	}
	return bvol
}

/*
// Todo, make the trace method for raycasting.
func (s *orthStack) Trace(o *Orthotope) *Orthotope {
	bvOrth, index := s.peek()
	orth := bvOrth.orth

	if bvOrth.desc[index] == nil {
		bvOrth, index = s.pop()
		for index == 1 {
			if !s.hasNext() {
				return nil
			}

			bvOrth, index = s.pop()
		}

		if index == 0 {
			// Change direction of parent.
			s.append(bvOrth, 1)
			s.append(bvOrth.desc[1], 0)
		}
	} else {
		s.append(bvOrth, 0)
	}
	return orth
}
*/
