//Copyright 2018 Brian Noyama. Subject to the the Apache License, Version 2.0.
package rect

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
	"testing"
)

func TestTopDownBVH(t *testing.T) {
	orths := make([]*Orthotope, len(leaf))
	copy(orths, leaf[:])
	tree := TopDownBVH(orths)
	if tree.Score() > 262 {
		t.Errorf("Ineeficient BVH created via TopDown:\n%v", tree.String())
		drawBVH(tree, "error_ideal_tree.png")
	}
}

func TestAdd(t *testing.T) {
	scores := [10]int{4, 26, 57, 77, 100, 120, 135, 188, 218, 247}

	tree := &BVol{}
	for index, orth := range leaf {
		if !tree.Add(orth) {
			t.Errorf("Unable to add: %v\n", orth.String())
		}
		if scores[index] != tree.Score() {
			drawBVH(tree, "error_add_tree.png")
			t.Errorf("Unexpected score: %d\nExpected: %d\nTree:\n%v", tree.Score(),
				scores[index], tree.String())
		}
	}

	if tree.Add(leaf[0]) {
		t.Errorf("Incorrectly added existing volume: %v\n", leaf[0].String())
	}

	ideal := getIdealTree()
	if !ideal.Equals(tree) {
		t.Errorf("Non-ideal BVH created via add:\n%v\nIdeal:\n%v", tree.String(),
			ideal.String())
	}
}

func TestRemove(t *testing.T) {
	tree := getIdealTree()

	// Reordering leaves to remove to test edge cases.
	var to_remove [9]*Orthotope = [9]*Orthotope{
		leaf[8],
		leaf[0],
		leaf[2],
		leaf[1],
		leaf[3],
		leaf[4],
		leaf[6],
		leaf[5],
		leaf[7],
	}

	scores := [9]int{233, 196, 173, 152, 112, 97, 77, 50, 10}

	for index, orth := range to_remove {
		if !tree.Remove(orth) {
			t.Errorf("Unable to remove: %v\n", orth.String())
		}
		if scores[index] != tree.Score() {
			drawBVH(tree, "error_remove_tree.png")
			t.Errorf("Unexpected score: %d\nExpected: %d\nTree:\n%v", tree.Score(),
				scores[index], tree.String())
		}
	}

	if !tree.Remove(leaf[9]) {
		t.Errorf("Unable to remove: %v\n", leaf[9].String())
	}

	if tree.Remove(leaf[0]) {
		t.Errorf("Incorrectly removing non-existing volume: %v\n", leaf[0].String())
	}

}

func TestString(t *testing.T) {
	tree := getIdealTree()
	expectedString :=
		"Point [2 2], Delta [21 23]\n" +
			" Point [16 2], Delta [7 23]\n" +
			"   Point [18 19], Delta [5 6]\n" +
			"    Point [18 21], Delta [2 2]\n" +
			"    Point [19 19], Delta [4 6]\n" +
			"  Point [16 2], Delta [6 12]\n" +
			"   Point [16 2], Delta [5 8]\n" +
			"    Point [19 2], Delta [2 2]\n" +
			"    Point [16 6], Delta [3 4]\n" +
			"   Point [17 12], Delta [5 2]\n" +
			"    Point [20 12], Delta [2 2]\n" +
			"    Point [17 12], Delta [2 2]\n" +
			"  Point [2 2], Delta [10 20]\n" +
			"   Point [4 11], Delta [8 11]\n" +
			"    Point [10 11], Delta [2 2]\n" +
			"    Point [4 16], Delta [6 6]\n" +
			"   Point [2 2], Delta [8 8]\n" +
			"    Point [7 7], Delta [3 3]\n" +
			"    Point [2 2], Delta [2 2]\n"
	actual := tree.String()
	if strings.Replace(actual, " 0", "", -1) != expectedString {
		t.Errorf("Actual string:\n%v\n...doesn't match expected:\n%v\n",
			actual, expectedString)
	}
}

func getIdealTree() *BVol {
	tree := &BVol{depth: 4,
		vol: &Orthotope{Point: [d]int{2, 2}, Delta: [d]int{21, 23}},
		desc: [2]*BVol{
			&BVol{depth: 3,
				vol: &Orthotope{Point: [d]int{16, 2}, Delta: [d]int{7, 23}},
				desc: [2]*BVol{
					&BVol{depth: 1,
						vol: &Orthotope{Point: [d]int{18, 19}, Delta: [d]int{5, 6}},
						desc: [2]*BVol{
							&BVol{vol: leaf[8]},
							&BVol{vol: leaf[9]},
						},
					},
					&BVol{depth: 2,
						vol: &Orthotope{Point: [d]int{16, 2}, Delta: [d]int{6, 12}},
						desc: [2]*BVol{
							&BVol{depth: 1,
								vol: &Orthotope{Point: [d]int{16, 2}, Delta: [d]int{5, 8}},
								desc: [2]*BVol{
									&BVol{vol: leaf[2]},
									&BVol{vol: leaf[3]},
								},
							},
							&BVol{depth: 1,
								vol: &Orthotope{Point: [d]int{17, 12}, Delta: [d]int{5, 2}},
								desc: [2]*BVol{
									&BVol{vol: leaf[6]},
									&BVol{vol: leaf[5]},
								},
							},
						},
					},
				},
			},
			&BVol{depth: 2,
				vol: &Orthotope{Point: [d]int{2, 2}, Delta: [d]int{10, 20}},
				desc: [2]*BVol{
					&BVol{depth: 1,
						vol: &Orthotope{Point: [d]int{4, 11}, Delta: [d]int{8, 11}},
						desc: [2]*BVol{
							&BVol{vol: leaf[4]},
							&BVol{vol: leaf[7]},
						},
					},
					&BVol{depth: 1,
						vol: &Orthotope{Point: [d]int{2, 2}, Delta: [d]int{8, 8}},
						desc: [2]*BVol{
							&BVol{vol: leaf[1]},
							&BVol{vol: leaf[0]},
						},
					},
				},
			},
		},
	}
	return tree
}

func drawBVH(BVol *BVol, name string) {
	myimage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{25, 25}})
	iter := BVol.Iterator()
	for iter.HasNext() {
		next := iter.Next()

		c := color.RGBA{uint8(255 / (next.depth + 1)), uint8(255 / (2*next.depth + 1)),
			uint8(255), 255}
		for y := next.vol.Point[1]; y < next.vol.Point[1]+next.vol.Delta[1]; y += 1 {
			myimage.Set(next.vol.Point[0], y, c)
			myimage.Set(next.vol.Point[0]+next.vol.Delta[0]-1, y, c)
		}
		for x := next.vol.Point[0]; x < next.vol.Point[0]+next.vol.Delta[0]; x += 1 {
			myimage.Set(x, next.vol.Point[1], c)
			myimage.Set(x, next.vol.Point[1]+next.vol.Delta[1]-1, c)
		}
	}
	myfile, _ := os.Create(name)
	png.Encode(myfile, myimage)
}

var leaf [10]*Orthotope = [10]*Orthotope{
	&Orthotope{Point: [d]int{2, 2}, Delta: [d]int{2, 2}},
	&Orthotope{Point: [d]int{7, 7}, Delta: [d]int{3, 3}},
	&Orthotope{Point: [d]int{19, 2}, Delta: [d]int{2, 2}},
	&Orthotope{Point: [d]int{16, 6}, Delta: [d]int{3, 4}},
	&Orthotope{Point: [d]int{10, 11}, Delta: [d]int{2, 2}},
	&Orthotope{Point: [d]int{17, 12}, Delta: [d]int{2, 2}},
	&Orthotope{Point: [d]int{20, 12}, Delta: [d]int{2, 2}},
	&Orthotope{Point: [d]int{4, 16}, Delta: [d]int{6, 6}},
	&Orthotope{Point: [d]int{18, 21}, Delta: [d]int{2, 2}},
	&Orthotope{Point: [d]int{19, 19}, Delta: [d]int{4, 6}},
}
