//Copyright 2018 Brian Noyama. Subject to the the Apache License, Version 2.0.
package rect

import (
	"testing"
)

func TestNext(t *testing.T) {
	bvs := []*BVol{
		&BVol{
			vol: &Orthotope{
				Point: [DIMENSIONS]int{2, 2},
				Delta: [DIMENSIONS]int{8, 8},
			},
		},
		&BVol{
			vol: &Orthotope{
				Point: [DIMENSIONS]int{2, 2},
				Delta: [DIMENSIONS]int{2, 2},
			},
		},
		&BVol{
			vol: &Orthotope{
				Point: [DIMENSIONS]int{7, 7},
				Delta: [DIMENSIONS]int{3, 3},
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
		&Orthotope{Point: [d]int{11, 12}, Delta: [d]int{0, 0}},
		&Orthotope{Point: [d]int{14, 15}, Delta: [d]int{0, 0}},
		&Orthotope{Point: [d]int{-2, -2}, Delta: [d]int{30, 30}},
		&Orthotope{Point: [d]int{30, 30}, Delta: [d]int{30, 30}},
		&Orthotope{Point: [d]int{17, 9}, Delta: [d]int{5, 5}},
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
}

func TestTrace(t *testing.T) {
	tree := getIdealTree()
	query := [5]*Orthotope{
		&Orthotope{Point: [d]int{-2, 0}, Delta: [d]int{4, 2}},
		&Orthotope{Point: [d]int{14, 11}, Delta: [d]int{-1, 0}},
		&Orthotope{Point: [d]int{7, 20}, Delta: [d]int{4, -5}},
		&Orthotope{Point: [d]int{30, 30}, Delta: [d]int{-1, -1}},
		&Orthotope{Point: [d]int{0, 40}, Delta: [d]int{5, -1}},
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
		prev_dist := 0
		for r, dist := iter.Trace(q); r != nil; r, dist = iter.Trace(q) {
			if dist < prev_dist {
				t.Errorf("Tracing %v returned a distance out of order: %d < %d\n",
					q.String(), dist, prev_dist)
			}
			prev_dist = dist
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
}

func TestBVHContains(t *testing.T) {
	tree := getIdealTree()

	var to_check [4]*Orthotope = [4]*Orthotope{
		leaf[2],
		leaf[7],
		&Orthotope{Point: [d]int{100, 20}, Delta: [d]int{8, 9}},
		&Orthotope{Point: [d]int{19, 2}, Delta: [d]int{2, 2}}, // Similar to leaf[2]
	}

	contains := [4]bool{true, true, false, false}

	iter := tree.Iterator()
	for index, orth := range to_check {
		if iter.Contains(orth) != contains[index] {
			if contains[index] {
				t.Errorf("Unable to find: %v\n", orth.String())
			} else {
				t.Errorf("Incorrectly found: %v\n", orth.String())
			}
		}
	}
}
