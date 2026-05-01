package bvh

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"testing"

	"github.com/briannoyama/bvh/volume"
)

type o32 = volume.Orthotope[int32]

var keys = [10]o32{
	{P0: [volume.DIM]int32{2, 2}, P1: [volume.DIM]int32{4, 4}},
	{P0: [volume.DIM]int32{7, 7}, P1: [volume.DIM]int32{10, 10}},
	{P0: [volume.DIM]int32{19, 2}, P1: [volume.DIM]int32{21, 4}},
	{P0: [volume.DIM]int32{16, 6}, P1: [volume.DIM]int32{19, 10}},
	{P0: [volume.DIM]int32{10, 11}, P1: [volume.DIM]int32{12, 13}},
	{P0: [volume.DIM]int32{17, 12}, P1: [volume.DIM]int32{19, 14}},
	{P0: [volume.DIM]int32{20, 12}, P1: [volume.DIM]int32{22, 14}},
	{P0: [volume.DIM]int32{4, 16}, P1: [volume.DIM]int32{10, 22}},
	{P0: [volume.DIM]int32{18, 21}, P1: [volume.DIM]int32{20, 23}},
	{P0: [volume.DIM]int32{19, 19}, P1: [volume.DIM]int32{23, 25}},
}

func TestAll(t *testing.T) {
	// Add
	scores := [10]float32{4, 26, 57, 77, 100, 120, 135, 188, 218, 247}
	refs := [10]int{}
	tree := NewInt32OrthTree[string](10)
	for i, v := range keys {
		refs[i] = tree.Add(v, v.String())
		if scores[i] != tree.Score() {
			t.Errorf("Unoptimal score (add). Expected %f, got %f\n", scores[i], tree.Score())
		}
	}

	for i, r := range refs {
		if keys[i] != tree.Key(r) {
			t.Errorf("Key does not match volume. Expected %v, got %v\n", keys[i], tree.Key(r))
		}
		k, v := tree.KeyVal(r)
		if k.String() != v {
			t.Errorf("Key does not map to expected value. Expected %v, got %s\n", k, v)
		}
	}

	// Query
	query(t, tree, o32{P0: [volume.DIM]int32{11, 12}, P1: [volume.DIM]int32{11, 12}}, keys[4:5])
	query(t, tree, o32{P0: [volume.DIM]int32{14, 15}, P1: [volume.DIM]int32{14, 15}}, keys[0:0])
	query(t, tree, o32{P0: [volume.DIM]int32{-2, -2}, P1: [volume.DIM]int32{28, 28}}, keys[:])
	query(t, tree, o32{P0: [volume.DIM]int32{30, 30}, P1: [volume.DIM]int32{60, 60}}, keys[0:0])
	query(t, tree, o32{P0: [volume.DIM]int32{17, 9}, P1: [volume.DIM]int32{22, 14}},
		[]o32{keys[3], keys[5], keys[6]})

	// Intersects
	intersect(t, tree, o32{P0: [volume.DIM]int32{-2, 0}, P1: [volume.DIM]int32{2, 2}},
		[volume.DIM]int32{14, 4}, []float32{0, 1}, []o32{keys[0], keys[3]})
	intersect(t, tree, o32{P0: [volume.DIM]int32{14, 11}, P1: [volume.DIM]int32{15, 12}},
		[volume.DIM]int32{-4, 0}, []float32{0.5}, keys[4:5])
	intersect(t, tree, o32{P0: [volume.DIM]int32{7, 20}, P1: [volume.DIM]int32{9, 22}},
		[volume.DIM]int32{20, -25}, []float32{0, 0.4, 0.64, 0.4}, []o32{keys[7], keys[3], keys[2], keys[5]})
	intersect(t, tree, o32{P0: [volume.DIM]int32{30, 30}, P1: [volume.DIM]int32{31, 31}},
		[volume.DIM]int32{-40, -40}, []float32{0.175, 0.45, 0.5, 0.65, 0.25}, []o32{keys[9], keys[4], keys[1], keys[0], keys[8]})
	intersect(t, tree, o32{P0: [volume.DIM]int32{0, 40}, P1: [volume.DIM]int32{1, 41}},
		[volume.DIM]int32{50, -10}, []float32{}, keys[0:0])

	// Remove
	scores = [10]float32{232, 195, 173, 152, 112, 97, 77, 50, 10, 0}
	toRemove := []int{8, 0, 2, 1, 3, 4, 6, 5, 7, 9}
	for i, k := range toRemove {
		tree.Remove(refs[k])
		if scores[i] != tree.Score() {
			t.Errorf("Unoptimal score (remove). Expected %f, got %f\n", scores[i], tree.Score())
		}
	}
}

func intersect(t *testing.T, tree Int32OrthTree[string], q o32, d [volume.DIM]int32, rd []float32, rO32 []o32) {
	results := map[o32]string{}
	for i, r := range rO32 {
		results[r] = fmt.Sprintf("%s : %f", r.String(), rd[i])
	}
	td := float32(0)
	for k, v := range tree.Intersects(q, d, &td) {
		checkAgainst(t, k, fmt.Sprintf("%s : %f", v, td), results)
	}
	checkAllValues(t, results)
}

func query(t *testing.T, tree Int32OrthTree[string], q o32, rO32 []o32) {
	t.Helper()
	results := map[o32]string{}
	for _, r := range rO32 {
		results[r] = r.String()
	}
	for k, v := range tree.Query(q) {
		checkAgainst(t, k, v, results)
	}
	checkAllValues(t, results)
}

func checkAgainst(t *testing.T, k o32, v string, results map[o32]string) {
	t.Helper()
	resultV, exists := results[k]
	if !exists {
		t.Errorf("Key %v not expected", k)
	} else {
		delete(results, k)
		if resultV != v {
			t.Errorf("Value \"%s\" not expected. Expected: \"%s\"", v, resultV)
		}
	}
}

func checkAllValues(t *testing.T, results map[o32]string) {
	t.Helper()
	if len(results) > 0 {
		t.Errorf("Values missing:\n  %s", strings.Join(slices.Collect(maps.Values(results)), "\n  "))
	}
}
