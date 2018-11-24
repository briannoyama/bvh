//Copyright 2018 Brian Noyama. Subject to the the Apache License, Version 2.0.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/briannoyama/bvh/rect"
)

func main() {
	config := flag.String("config", "test.json", "JSON configuration for the test.")
	flag.Parse()
	configFile, err := os.Open(*config)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	defer configFile.Close()

	configBytes, _ := ioutil.ReadAll(configFile)

	test := &bvhTest{}
	json.Unmarshal([]byte(configBytes), test)
	test.runTest()
}

type operation struct {
	orth   *rect.Orthotope
	opcode int
}

type bvhTest struct {
	MaxBounds *rect.Orthotope
	MinVol    *[rect.DIMENSIONS]int
	MaxVol    *[rect.DIMENSIONS]int
	Additions int
	Removals  int
	Queries   int
	RandSeed  int64
}

func (b *bvhTest) runTest() {
	orths := make([]*rect.Orthotope, 0, b.Additions)
	removed := make(map[int]bool, b.Additions)
	bvol := &rect.BVol{}
	iter := bvol.Iterator()
	r := rand.New(rand.NewSource(b.RandSeed))

	if b.Removals > b.Additions {
		fmt.Printf("Incorrect config, removals larger than additions.\n")
		return
	}

	removals := *distribute(r, b.Removals, b.Additions)
	queries := *distribute(r, b.Queries, b.Additions)
	total := 0

	for a := 0; a < b.Additions; a += 1 {
		orth := b.makeOrth(r)
		orths = append(orths, orth)

		// Test the addition operation.
		t := time.Now()
		iter.Add(orth)
		duration := time.Now().Sub(t).Nanoseconds()
		total += 1
		fmt.Printf("add, %d, %d, %d\n", total, bvol.GetDepth(), duration)

		for removal := 0; removal < removals[a]; removal += 1 {
			toRemove := r.Intn(a + 1)
			for ; removed[toRemove] && toRemove <= a; toRemove += 1 {
			}
			if toRemove <= a {
				removed[toRemove] = true

				// Test the removal operation.
				t = time.Now()
				iter.Remove(orths[toRemove])
				duration := time.Now().Sub(t).Nanoseconds()
				total -= 1
				fmt.Printf("sub, %d, %d, %d\n", total, bvol.GetDepth(), duration)
			} else if a+1 < len(removals) {
				removals[a+1] += 1
			}
		}
		for query := 0; query < queries[a]; query += 1 {
			q := b.makeOrth(r)
			iter.Reset()
			count := 0

			// Test the query operation.
			t = time.Now()
			for r := iter.Query(q); r != nil; r = iter.Query(q) {
				count += 1
			}
			duration := time.Now().Sub(t).Nanoseconds()
			fmt.Printf("que, %d, %d, %d, %d\n", total, bvol.GetDepth(),
				duration, count)
		}
	}
}

func distribute(r *rand.Rand, totalEvents int, steps int) *[]int {
	events := make([]int, steps)
	for e := 0; e < totalEvents; e += 1 {
		events[r.Intn(steps)] += 1
	}

	return &events
}

func (b *bvhTest) makeOrth(r *rand.Rand) *rect.Orthotope {
	orth := &rect.Orthotope{}
	for d := 0; d < rect.DIMENSIONS; d += 1 {
		orth.Delta[d] = b.MinVol[d] + r.Intn(b.MaxVol[d]-b.MinVol[d])
		orth.Point[d] = b.MaxBounds.Point[d] + r.Intn(b.MaxBounds.Delta[d]-
			orth.Delta[d])
	}
	return orth
}
