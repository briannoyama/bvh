//Copyright 2018 Brian Noyama. Subject to the the Apache License, Version 2.0.
package perf

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
	configFile, err := os.Open(*config)
	if err != nil {
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
	maxBounds *rect.Orthotope
	minVol    *[rect.DIMENSIONS]int
	maxVol    *[rect.DIMENSIONS]int
	additions int
	removals  int
	queries   int
	randSeed  int64
}

func (b *bvhTest) runTest() {
	orths := make([]*rect.Orthotope, 0, b.additions)
	removals := make(map[int]bool)
	bvol := &rect.BVol{}
	iter := bvol.Iterator()
	r := rand.New(rand.NewSource(b.randSeed))
	for a := 0; a < b.additions; a += 1 {
		orth := b.makeOrth(r)
		orths = append(orths, orth)

		// Test the addition operation.
		t := time.Now()
		bvol.Add(orth)
		duration := t.Sub(time.Now())

		fmt.Printf("add, %d, %d, %d\n", a-b.removals, bvol.GetDepth(), duration)
		if b.removals > 0 && b.removals > r.Intn(b.additions) {
			toRemove := r.Intn(a)
			_, exists := removals[toRemove]
			for attempts := 100; exists && attempts > 0; attempts -= 1 {
				toRemove = r.Intn(a)
				_, exists = removals[toRemove]
			}
			removals[toRemove] = true

			// Test the removal operation.
			t = time.Now()
			bvol.Remove(orths[toRemove])
			duration = t.Sub(time.Now())

			b.removals -= 1
			fmt.Printf("sub, %d, %d, %d\n", a-b.removals, bvol.GetDepth(), duration)
		}
		if b.queries > 0 && b.queries > r.Intn(b.additions) {
			q := b.makeOrth(r)
			iter.Reset()
			count := 0

			// Test the query operation.
			t = time.Now()
			for r := iter.Query(q); r != nil; r = iter.Query(q) {
				count += 1
			}
			duration = t.Sub(time.Now())

			fmt.Printf("que, %d, %d, %d, %d\n", a-b.removals, bvol.GetDepth(),
				duration, count)
		}
	}
}

func (b *bvhTest) makeOrth(r *rand.Rand) *rect.Orthotope {
	orth := &rect.Orthotope{}
	for d := 0; d < rect.DIMENSIONS; d += 1 {
		orth.Delta[d] = b.minVol[d] + r.Intn(b.maxVol[d]-b.minVol[d])
		orth.Point[d] = b.maxBounds.Point[d] + r.Intn(b.maxBounds.Delta[d]-
			orth.Delta[d])
	}
	return orth
}
