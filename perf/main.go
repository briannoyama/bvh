//Copyright 2018 Brian Noyama. Subject to the the Apache License, Version 2.0.
package perf

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

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
	randSeed  int
}

func (*bvhTest) runTest() {

}
