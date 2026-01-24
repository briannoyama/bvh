package bvh

import (
	"github.com/briannoyama/bvh/volume"
	"github.com/briannoyama/go-fast/v2/fast"
)

func NewTree[K Volume[K, D], V, D any](len int) Tree[K, V, D] {
	return Tree[K, V, D]{
		FTreeMap: fast.NewPreAllocFTreeMap[VolDepth[K], V](len),
	}
}

type Int32OrthTree[V any] = Tree[volume.Orthotope[int32], V, [volume.DIM]int32]

func NewInt32OrthTree[V any](len int) Int32OrthTree[V] {
	return NewTree[volume.Orthotope[int32], V, [volume.DIM]int32](len)
}

type IntOrthTree[V any] = Tree[volume.Orthotope[int], V, [volume.DIM]int]

func NewIntOrthTree[V any](len int) IntOrthTree[V] {
	return NewTree[volume.Orthotope[int], V, [volume.DIM]int](len)
}

type Float32OrthTree[V any] = Tree[volume.Orthotope[float32], V, [volume.DIM]float32]

func NewFloat32OrthTree[V any](len int) Float32OrthTree[V] {
	return NewTree[volume.Orthotope[float32], V, [volume.DIM]float32](len)
}
