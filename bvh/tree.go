package bvh

import (
	"fmt"
	"iter"
	"math"
	"strings"

	"github.com/briannoyama/go-fast/v2/fast"
)

type VolDepth[K any] struct {
	depth int
	vol   K
}

type Volume[K any, D any] interface {
	Intersects(k1 K, d D) float32
	Minbound(k1 K) K
	Overlaps(k1 K) bool
	Score() float32
}

type Tree[K Volume[K, D], V, D any] struct {
	fast.FTreeMap[VolDepth[K], V]
}

func (b *Tree[K, V, D]) Add(k K, v V) int {
	var ref int
	// Only happens when tree is empty or has 1 element
	if b.Len() < 2 {
		// Second element added
		if b.Len() == 0 {
			ref = b.Add0(VolDepth[K]{vol: k}, v)
		} else {
			ref = b.Add1(VolDepth[K]{vol: k}, v)
			b.minBoundParent(0)
		}
		return ref
	}

	// Find best place to add key
	parent := 0
	next := 0
	current := b.Root()
	for current >= 0 {
		rel := b.Rel(current)
		k0 := b.Key(rel[0]).vol
		k1 := b.Key(rel[1]).vol
		s0 := k0.Minbound(k).Score() - k0.Score()
		s1 := k1.Minbound(k).Score() - k1.Score()
		next = int(math.Float32bits(s1-s0) >> 31)
		parent = current
		current = rel[next]
	}

	// We've reached a leaf node
	ref = b.AddAdj(parent, next, VolDepth[K]{vol: k}, v)

	// Adding inserts a new child below the parent.
	parent = b.Rel(parent)[next]
	b.minBoundParent(parent)

	// Travel up the tree to rebalance
	gparent := b.Rel(parent)[2]
	for parent != b.Root() {
		rel := b.Rel(gparent)
		aunt := rel[0] ^ parent ^ rel[1]
		if b.Key(current).depth > b.Key(aunt).depth {
			b.Swap(current, aunt)
			b.minBoundParent(parent)
		}

		// Bump to the parent node
		current = parent
		parent = gparent
		gparent = rel[2]
		// Try to swap children to optimize tree
		b.redistribute(parent)
		b.minBoundParent(parent)
	}
	return ref
}

func (b *Tree[K, V, D]) Depth() int {
	return b.Key(b.Root()).depth
}

func (b *Tree[K, V, D]) Intersects(k K, d D, t *float32) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		b.VisitAll(func(v *VolDepth[K]) bool {
			*t = k.Intersects(v.vol, d)
			return 0 <= *t && *t <= 1
		}, func(k *VolDepth[K], v *V) bool {
			return yield(k.vol, *v)
		})
	}
}

func (b *Tree[K, V, D]) Remove(ref int) (K, V) {
	if b.Len() == 1 {
		k, v := b.FTreeMap.Remove0(ref)
		return k.vol, v
	}
	k, v, parent := b.FTreeMap.Remove(ref)
	gparent := b.Parent(parent)
	for parent != b.Root() {
		aunts := b.Rel(gparent)
		diff := b.Key(aunts[0]).depth - b.Key(aunts[1]).depth
		// If there's a difference greater than 1
		if max(diff, -diff) > 1 {
			bigger := (diff >> 63) & 1
			cousins := b.Rel(aunts[bigger])
			c := []*VolDepth[K]{b.Key(cousins[0]), b.Key(cousins[1])}
			cbigger := ((c[0].depth - c[1].depth) >> 63) & 1
			b.Swap(cousins[cbigger], aunts[bigger^1])
			b.minBoundParent(aunts[bigger])
		}

		b.minBoundParent(gparent)
		if b.Key(gparent).depth > 1 {
			b.redistribute(gparent)
		}
		parent = gparent
		gparent = b.Parent(parent)
	}
	return k.vol, v
}

func (b *Tree[K, V, D]) Query(k K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		b.VisitAll(func(v *VolDepth[K]) bool {
			return k.Overlaps(v.vol)
		}, func(k *VolDepth[K], v *V) bool {
			return yield(k.vol, *v)
		})
	}
}

func (b *Tree[K, V, D]) Score() float32 {
	score := float32(0)
	b.VisitAll(func(v *VolDepth[K]) bool {
		score += v.vol.Score()
		return true
	}, func(*VolDepth[K], *V) bool { return true })
	return score
}

func (b *Tree[K, V, D]) String() string {
	maxDepth := b.Key(b.Root()).depth
	lines := []string{}
	b.VisitAll(func(v *VolDepth[K]) bool {
		lines = append(lines, fmt.Sprintf("%s%v", strings.Repeat(" ", maxDepth-v.depth), v.vol))
		return true
	}, func(_ *VolDepth[K], v *V) bool {
		lines[len(lines)-1] += fmt.Sprintf(": %v", v)
		return true
	})
	return strings.Join(lines, "\n")
}

func (b *Tree[K, V, D]) minBoundParent(ref int) {
	rel := b.Rel(ref)
	k0, k1 := b.Key(rel[0]), b.Key(rel[1])
	*b.Key(ref) = VolDepth[K]{
		vol:   k0.vol.Minbound(k1.vol),
		depth: max(k0.depth+1, k1.depth+1),
	}
}

func (b *Tree[K, V, D]) redistribute(ref int) {
	children := b.Rel(ref)
	d0, d1 := b.Key(children[0]).depth, b.Key(children[1]).depth
	if d0 == d1 && d0 > 0 {
		b.swapcheck0(children[0], children[1])
	} else {
		// Need & 1 since shifting negative int will result in -1
		left := int((d0-d1)>>63) & 1
		b.swapcheck1(children[left], children[1^left])
	}
}

func (b *Tree[K, V, D]) swapcheck0(first, second int) {
	r0, r1 := b.Rel(first), b.Rel(second)
	k00, k01, k10, k11 := b.Key(r0[0]), b.Key(r0[1]), b.Key(r1[0]), b.Key(r1[1])

	// Base case minIndex = -1
	minScore := k00.vol.Minbound(k01.vol).Score() + k10.vol.Minbound(k11.vol).Score()

	// Swap minIndex = 0
	score := k10.vol.Minbound(k01.vol).Score() + k00.vol.Minbound(k11.vol).Score()
	minIndex := 0 - int32(math.Float32bits(minScore-score)>>31)
	minScore = min(score, minScore)

	// Swap minIndex = 1
	score = k10.vol.Minbound(k00.vol).Score() + k01.vol.Minbound(k11.vol).Score()
	flag := int32(math.Float32bits(minScore-score) >> 31)
	minIndex = (flag * minIndex) | (flag ^ 1)

	// Only swap for cases where swap has a smaller score.
	if minIndex >= 0 {
		b.Swap(r0[minIndex], r1[0])
		b.minBoundParent(first)
		b.minBoundParent(second)
	}
}

func (b *Tree[K, V, D]) swapcheck1(first, second int) {
	r0 := b.Rel(first)
	k00, k01 := b.Key(r0[0]), b.Key(r0[1])
	k1 := b.Key(second)

	// Base case minIndex = -1
	minScore := k00.vol.Minbound(k01.vol).Score()

	// Swap minIndex = 0
	score := k1.vol.Minbound(k01.vol).Score()
	minIndex := 0 - int32(math.Float32bits(minScore-score)>>31)
	minScore = min(score, minScore)

	// Swap minIndex = 1
	score = k1.vol.Minbound(k00.vol).Score()
	flag := int32(math.Float32bits(minScore-score) >> 31)
	minIndex = (flag * minIndex) | (flag ^ 1)

	// Only swap for cases where swap has a smaller score.
	if minIndex >= 0 {
		b.Swap(r0[minIndex], second)
		b.minBoundParent(first)
	}
}
