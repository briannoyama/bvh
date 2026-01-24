// Copyright 2018 Brian Noyama. Subject to the the Apache License, Version 2.0.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"log"
	"math/rand"
	"os"
	"time"

	"github.com/briannoyama/bvh/bvh"
	"github.com/briannoyama/bvh/volume"
)

func main() {
	config := flag.String("config", "testBVH.json",
		"JSON configuration for the test.")
	compare := flag.Bool("compare", false,
		"Compare with Top Down method? Default False.")
	flag.Parse()
	configFile, err := os.Open(*config)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	defer configFile.Close()

	configBytes, _ := io.ReadAll(configFile)

	test := &bvhTest{}
	json.Unmarshal([]byte(configBytes), test)
	if *compare {
		test.comparisonTest()
	} else {
		test.runTest()
	}
}

type bvhTest struct {
	MaxBounds volume.Orthotope[int32]
	MinVol    [volume.DIM]int32
	MaxVol    [volume.DIM]int32
	Additions int
	Removals  int
	Queries   int
	RandSeed  int64
}

func (b *bvhTest) comparisonTest() {
	orths := make([]volume.Orthotope[int32], 0, b.Additions)
	r := rand.New(rand.NewSource(b.RandSeed))
	btree := bvh.NewInt32OrthTree[int](110000)
	for a := 0; a < b.Additions; a += 1 {
		orth := b.makeOrth(r)
		orths = append(orths, orth)

		btree.Add(orth, 0)
		fmt.Printf("%d, %d, %d\n", a, btree.Depth(), int(btree.Score()))
	}
}

func (b *bvhTest) runTest() {
	refs := make([]int, b.Additions)
	orths := make([]volume.Orthotope[int32], 0, b.Additions)
	removed := make(map[int]bool, b.Additions)
	btree := bvh.NewInt32OrthTree[int](110000)
	r := rand.New(rand.NewSource(b.RandSeed))

	if b.Removals > b.Additions {
		fmt.Printf("Incorrect config, removals larger than additions.\n")
		return
	}

	removals := distribute(r, b.Removals, b.Additions)
	queries := distribute(r, b.Queries, b.Additions)
	total := 0
	var addTime, subTime, queTime int64

	for a := 0; a < b.Additions; a += 1 {
		orth := b.makeOrth(r)
		orths = append(orths, orth)

		// Test the addition operation.
		t := time.Now()
		refs[a] = btree.Add(orth, a)
		duration := time.Since(t).Nanoseconds()
		total += 1
		fmt.Printf("add, %d, %d, %d, \"%v\"\n", total, btree.Depth(), duration, orth)
		addTime += duration

		for removal := 0; removal < removals[a]; removal += 1 {
			toRemove := r.Intn(a + 1)
			for ; removed[toRemove] && toRemove <= a; toRemove += 1 {
			}
			if toRemove <= a {
				removed[toRemove] = true

				// Test the removal operation.
				t = time.Now()
				k, _ := btree.Remove(refs[toRemove])
				duration := time.Since(t).Nanoseconds()
				total -= 1
				fmt.Printf("sub, %d, %d, %d, \"%v\"\n", total, btree.Depth(), duration, k)
				subTime += duration
			} else if a+1 < len(removals) {
				removals[a+1] += 1
			}
		}

		count := 0

		for query := 0; query < queries[a]; query += 1 {
			q := b.makeOrth(r)
			count = 0

			// Test the query operation.
			t = time.Now()
			for range btree.Query(q) {
				count += 1
			}
			duration := time.Since(t).Nanoseconds()
			fmt.Printf("que, %d, %d, %d, %d, \"%v\"\n", total, btree.Depth(),
				duration, count, q)
			queTime += duration
		}
	}
	fmt.Printf("score: %d, add: %d, sub: %d, que: %d\n", int(btree.Score()), addTime, subTime, queTime)
	btree.Verify()
}

func distribute(r *rand.Rand, totalEvents int, steps int) []int {
	events := make([]int, steps)
	for e := 0; e < totalEvents; e += 1 {
		events[r.Intn(steps)] += 1
	}

	return events
}

func (b *bvhTest) makeOrth(r *rand.Rand) volume.Orthotope[int32] {
	orth := volume.Orthotope[int32]{}
	for d := 0; d < volume.DIM; d += 1 {
		delta := b.MinVol[d] + r.Int31n(b.MaxVol[d]-b.MinVol[d])

		orth.P0[d] = r.Int31n(b.MaxBounds.P1[d] - delta)
		orth.P1[d] = orth.P0[d] + delta
	}
	return orth
}
