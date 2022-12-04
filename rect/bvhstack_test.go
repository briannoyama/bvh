// Copyright 2018 Brian Noyama. Subject to the the Apache License, Version 2.0.
package rect

import (
	"testing"
)

func TestNext(t *testing.T) {
	bvs := []*BVol{
		{
			vol: &Orthotope{
				Point: [DIMENSIONS]int32{2, 2},
				Delta: [DIMENSIONS]int32{8, 8},
			},
		},
		{
			vol: &Orthotope{
				Point: [DIMENSIONS]int32{2, 2},
				Delta: [DIMENSIONS]int32{2, 2},
			},
		},
		{
			vol: &Orthotope{
				Point: [DIMENSIONS]int32{7, 7},
				Delta: [DIMENSIONS]int32{3, 3},
			},
		},
	}

	bvs[0].desc = [2]*BVol{bvs[1], bvs[2]}
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
	tree := getIdealTree()
	query := [5]*Orthotope{
		{Point: [d]int32{11, 12}, Delta: [d]int32{0, 0}},
		{Point: [d]int32{14, 15}, Delta: [d]int32{0, 0}},
		{Point: [d]int32{-2, -2}, Delta: [d]int32{30, 30}},
		{Point: [d]int32{30, 30}, Delta: [d]int32{30, 30}},
		{Point: [d]int32{17, 9}, Delta: [d]int32{5, 5}},
	}
	results := [5][]*Orthotope{
		{leaf[4]},
		{},
		make([]*Orthotope, len(leaf)),
		{},
		{leaf[3], leaf[5], leaf[6]},
	}
	// For the second test, copy contents of leaf.
	copy(results[2], leaf[:])

	for in, q := range query {
		iter := tree.Iterator()
		iter.Reset()
		for r := iter.Query(q); r != nil; r = iter.Query(q) {
			found := 0
			for _, orth := range results[in] {
				if r == orth {
					break
				}
				found++
			}
			if found < len(results[in]) {
				results[in] = append(results[in][:found], results[in][found+1:]...)
			} else {
				t.Errorf("Querying %v returned unexpected value: %v\n",
					q.String(), r.String())
			}
		}
		for _, orth := range results[in] {
			t.Errorf("Querying %v did not return %v\n", q.String(), orth.String())
		}
	}
	iter := (&BVol{}).Iterator()
	if iter.Query(leaf[0]) != nil {
		t.Errorf("Querying an empty hierarchy returned non nil value!\n")
	}
}

func TestTrace(t *testing.T) {
	tree := getIdealTree()
	query := [5]*Orthotope{
		{Point: [d]int32{-2, 0}, Delta: [d]int32{4, 2}},
		{Point: [d]int32{14, 11}, Delta: [d]int32{-1, 0}},
		{Point: [d]int32{7, 20}, Delta: [d]int32{4, -5}},
		{Point: [d]int32{30, 30}, Delta: [d]int32{-1, -1}},
		{Point: [d]int32{0, 40}, Delta: [d]int32{5, -1}},
	}
	results := [5][]*Orthotope{
		{leaf[0], leaf[3]},
		{leaf[4]},
		{leaf[7], leaf[3], leaf[2]},
		{leaf[9], leaf[4], leaf[1], leaf[0]},
		{},
	}

	for in, q := range query {
		iter := tree.Iterator()
		iter.Reset()
		prevDist := int32(0)
		for r, dist := iter.Trace(q); r != nil; r, dist = iter.Trace(q) {
			if dist < prevDist {
				t.Errorf("Tracing %v returned a distance out of order: %d < %d\n",
					q.String(), dist, prevDist)
			}
			prevDist = dist
			if len(results[in]) > 0 && results[in][0] == r {
				results[in] = results[in][1:]
			} else {
				t.Errorf("Tracing %v returned unexpected or out of order value: %v\n",
					q.String(), r.String())
			}
		}
		for _, orth := range results[in] {
			t.Errorf("Tracing %v did not return %v\n", q.String(), orth.String())
		}
	}
	iter := (&BVol{}).Iterator()
	if r, _ := iter.Trace(leaf[0]); r != nil {
		t.Errorf("Tracing an empty hierarchy returned non nil value!\n")
	}
}

func TestBVHContains(t *testing.T) {
	tree := getIdealTree()

	toCheck := [4]*Orthotope{
		leaf[2],
		leaf[7],
		{Point: [d]int32{100, 20}, Delta: [d]int32{8, 9}},
		{Point: [d]int32{19, 2}, Delta: [d]int32{2, 2}}, // Similar to leaf[2]
	}

	contains := [4]bool{true, true, false, false}

	iter := tree.Iterator()
	for index, orth := range toCheck {
		if iter.Contains(orth) != contains[index] {
			if contains[index] {
				t.Errorf("Unable to find: %v\n", orth.String())
			} else {
				t.Errorf("Incorrectly found: %v\n", orth.String())
			}
		}
	}
}
