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

func TestTrace(t *testing.T) {
	//tree := getIdealTree()
}
