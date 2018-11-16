package rect

import (
	"testing"
)

func TestNext(t *testing.T) {
	bvs := []*BVol{
		&BVol{
			orth: &Orthotope{
				point: [DIMENSIONS]int{2, 2},
				delta: [DIMENSIONS]int{8, 8},
			},
		},
		&BVol{
			orth: &Orthotope{
				point: [DIMENSIONS]int{2, 2},
				delta: [DIMENSIONS]int{2, 2},
			},
		},
		&BVol{
			orth: &Orthotope{
				point: [DIMENSIONS]int{7, 7},
				delta: [DIMENSIONS]int{3, 3},
			},
		},
	}

	bvs[0].vol = [2]*BVol{bvs[1], bvs[2]}
	bvs[0].depth = 1

	iter := bvs[0].Iterator()
	for i := 0; iter.HasNext(); i++ {
		next := iter.Next()

		if bvs[i] != next {
			t.Errorf("Iterator did not return the element %v in order", i)
		}
	}
}

func TestQuery(t *testing.T) {
}

func TestTrace(t *testing.T) {
	//tree := getIdealTree()
}

func TestBVHContains(t *testing.T) {
	tree := getIdealTree()

	var to_check [4]*Orthotope = [4]*Orthotope{
		leaf[2],
		leaf[7],
		&Orthotope{point: [d]int{100, 20}, delta: [d]int{8, 9}},
		&Orthotope{point: [d]int{19, 2}, delta: [d]int{2, 2}}, // Similar to leaf[2]
	}

	contains := [4]bool{true, true, false, false}

	iter := tree.Iterator()
	for index, orth := range to_check {
		if iter.Contains(orth) != contains[index] {
			if contains[index] {
				t.Errorf("Unable to find: %v\n", orth.String())
			} else {
				t.Errorf("Incorrectly found: %v\n", orth.String())
			}
		}
	}
}
