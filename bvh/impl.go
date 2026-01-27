package bvh

import (
	"github.com/briannoyama/bvh/volume"
	"github.com/briannoyama/go-fast/v2/fast"
)

func NewTree[K volume.Volume[K, D], V, D any](len int) Tree[K, V, D] {
	return Tree[K, V, D]{
		FTreeMap: fast.NewPreAllocFTreeMap[VolDepth[K], V](len),
	}
}

// Int32OrthTree: A BVH tree mapping int32 Orthotope to V
type Int32OrthTree[V any] = Tree[volume.Orthotope[int32], V, [volume.DIM]int32]

// NewInt32OrthTree: Create a tree of preallocated len
func NewInt32OrthTree[V any](len int) Int32OrthTree[V] {
	return NewTree[volume.Orthotope[int32], V, [volume.DIM]int32](len)
}

// IntOrthTree: A BVH tree mapping int Orthotope to V
type IntOrthTree[V any] = Tree[volume.Orthotope[int], V, [volume.DIM]int]

// NewIntOrthTree: Create a tree of preallocated len
func NewIntOrthTree[V any](len int) IntOrthTree[V] {
	return NewTree[volume.Orthotope[int], V, [volume.DIM]int](len)
}

// Float32OrthTree: A BVH tree mapping float32 Orthotope to V
type Float32OrthTree[V any] = Tree[volume.Orthotope[float32], V, [volume.DIM]float32]

// NewFloat32OrthTree: Create a tree of preallocated len
func NewFloat32OrthTree[V any](len int) Float32OrthTree[V] {
	return NewTree[volume.Orthotope[float32], V, [volume.DIM]float32](len)
}

// Int32SphereTree: A BVH tree mapping int32 Sphere to V
type Int32SphereTree[V any] = Tree[volume.Sphere[int32], V, [volume.DIM]int32]

// NewInt32SphereTree: Create a tree of preallocated len
func NewInt32SphereTree[V any](len int) Int32SphereTree[V] {
	return NewTree[volume.Sphere[int32], V, [volume.DIM]int32](len)
}

// IntSphereTree: A BVH tree mapping int Sphere to V
type IntSphereTree[V any] = Tree[volume.Sphere[int], V, [volume.DIM]int]

// NewIntSphereTree: Create a tree of preallocated len
func NewIntSphereTree[V any](len int) IntSphereTree[V] {
	return NewTree[volume.Sphere[int], V, [volume.DIM]int](len)
}

// Float32SphereTree: A BVH tree mapping float32 Sphere to V
type Float32SphereTree[V any] = Tree[volume.Sphere[float32], V, [volume.DIM]float32]

// NewFloat32SphereTree: Create a tree of preallocated len
func NewFloat32SphereTree[V any](len int) Float32SphereTree[V] {
	return NewTree[volume.Sphere[float32], V, [volume.DIM]float32](len)
}
