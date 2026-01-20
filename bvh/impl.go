package bvh

import (
	"github.com/briannoyama/bvh/volume"
	"github.com/briannoyama/go-fast/fast"
)

type Int32Node[V any] = Tree[volume.Orthotope[int32], V, [volume.DIM]int32]

func NewInt32Node[V any]() Int32Node[V] {
	return Int32Node[V]{
		FTreeMap:   fast.NewFTreeMap[VolDepth[volume.Orthotope[int32]], V](),
		contains:   volume.Contains[int32],
		equals:     volume.Equals[int32],
		intersects: volume.Intersects[int32],
		minbound:   volume.MinBounds[int32],
		overlaps:   volume.Overlaps[int32],
		score:      volume.Score[int32],
	}
}
