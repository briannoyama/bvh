package bvh

import (
	"github.com/briannoyama/bvh/volume"
	"github.com/briannoyama/go-fast/v2/fast"
)

type Int32Node[V any] = Tree[volume.Orthotope[int32], V, [volume.DIM]int32]

func NewInt32Node[V any](len int) Int32Node[V] {
	return Int32Node[V]{
		FTreeMap: fast.NewPreAllocFTreeMap[VolDepth[volume.Orthotope[int32]], V](len),
	}
}
