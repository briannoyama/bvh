package main

import (
	"github.com/briannoyama/bvh/rect"
	"testing"
)

func TestExample(t *testing.T) {

	// Change the DIMENSIONS constant in orthotope.go for your use case.
	orth := &rect.Orthotope{Point: [3]int32{10, -20, 10}, Delta: [3]int32{30, 30, 30}}
	bvol := &rect.BVol{}

	// Convenience method for adding/removing orthotopes.
	bvol.Add(orth)
	bvol.Remove(orth)

	// Use an iterator to reduce the amount of Garbage Collection
	iter := bvol.Iterator()

	// Iterators can all Add/Remove Orthotopes.
	iter.Add(orth)

	// You can add identical Orthotopes. BVol differentiates identical orthotope objects by their address
	orth2 := &rect.Orthotope{Point: [3]int32{10, -20, 10}, Delta: [3]int32{30, 30, 30}}
	iter.Add(orth2)

	// Use iterators to query for overlapping Orthotopes
	q := &rect.Orthotope{Point: [3]int32{0, -10, 10}, Delta: [3]int32{20, 20, 20}}
	iter.Reset()
	for r := iter.Query(q); r != nil; r = iter.Query(q) {
		t.Logf("Orthtope: %d @%p", r, r)
	}

	iter.Remove(orth2)
	orth2.Point[0] = 15
	iter.Add(orth2)

	// Use iterators to Trace rays through Orthotopes
	iter.Reset()
	for r, d := iter.Trace(q); r != nil; r, d = iter.Trace(q) {
		// Distances are arbitrary/relative. Objects farther away will have higher distances.
		t.Logf("Orthtope: %d @%p w/ Distance: %d", r, r, d)
	}

}
