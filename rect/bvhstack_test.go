package rect
/*
import (
	"testing"
)

func TestNext(t *testing.T) {
	bvs := []*BVOrth{
		&BVOrth{
			orth: &Orthotope{
				point: []int{2, 2},
				delta: []int{8, 8},
			},
			desc:  [2]*BVOrth{},
			depth: 1,
		},
		&BVOrth{
			orth: &Orthotope{
				point: []int{2, 2},
				delta: []int{2, 2},
			},
		},
		&BVOrth{
			orth: &Orthotope{
				point: []int{7, 7},
				delta: []int{3, 3},
			},
		},
	}

	bvs[0].desc[0] = bvs[1]
	bvs[0].desc[1] = bvs[2]

	iter := bvs[0].Iterator()
	for i, next := 0, iter.Next(); next != nil; i, next = i+1, iter.Next() {
		if bvs[i] != next {
			t.Errorf("Iterator did not return the %vth element in order", i)
		}
	}
}

func TestTrace(t *testing.T) {
	//tree := getIdealTree()
}
*/
