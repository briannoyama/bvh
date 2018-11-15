package rect

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func drawBVH(BVol *BVol, name string) {
	myimage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{25, 25}})
	iter := BVol.Iterator()
	for iter.HasNext() {
		next := iter.Next()

		c := color.RGBA{uint8(255 / (next.depth + 1)), uint8(255 / (2*next.depth + 1)),
			uint8(255), 255}
		for y := next.orth.point[1]; y < next.orth.point[1]+next.orth.delta[1]; y += 1 {
			myimage.Set(next.orth.point[0], y, c)
			myimage.Set(next.orth.point[0]+next.orth.delta[0]-1, y, c)
		}
		for x := next.orth.point[0]; x < next.orth.point[0]+next.orth.delta[0]; x += 1 {
			myimage.Set(x, next.orth.point[1], c)
			myimage.Set(x, next.orth.point[1]+next.orth.delta[1]-1, c)
		}
	}
	myfile, _ := os.Create(name)
	png.Encode(myfile, myimage)
}

var leaf [10]*Orthotope = [10]*Orthotope{
	&Orthotope{point: [d]int{2, 2}, delta: [d]int{2, 2}},
	&Orthotope{point: [d]int{7, 7}, delta: [d]int{3, 3}},
	&Orthotope{point: [d]int{19, 2}, delta: [d]int{2, 2}},
	&Orthotope{point: [d]int{16, 6}, delta: [d]int{3, 4}},
	&Orthotope{point: [d]int{10, 11}, delta: [d]int{2, 2}},
	&Orthotope{point: [d]int{17, 12}, delta: [d]int{2, 2}},
	&Orthotope{point: [d]int{20, 12}, delta: [d]int{2, 2}},
	&Orthotope{point: [d]int{4, 16}, delta: [d]int{6, 6}},
	&Orthotope{point: [d]int{18, 21}, delta: [d]int{2, 2}},
	&Orthotope{point: [d]int{19, 19}, delta: [d]int{4, 6}},
}

func TestAdd(t *testing.T) {
	scores := [9]int{26, 57, 77, 100, 120, 135, 188, 218, 247}

	tree := &BVol{orth: leaf[0]}
	for index, orth := range leaf[1:] {
		if !tree.Add(orth) {
			t.Errorf("Unable to add: %v\n", orth.String())
		}
		if scores[index] != tree.Score() {
			drawBVH(tree, "error_tree.png")
			t.Errorf("Unexpected score: %d\nExpected: %d\nTree:\n%v", tree.Score(),
				scores[index], tree.String())
		}
	}

	if tree.Add(leaf[0]) {
		t.Errorf("Incorrectly added existing volume: %v\n", leaf[0].String())
	}

	/*
		for index, orth := range leaf[:9] {
			fmt.Printf("%d-r\n", index)
			tree.Remove(orth)
			drawBVH(tree, fmt.Sprintf("test%d-r.png", index))
			fmt.Printf("%v\n", tree.String())
		}
	*/

	ideal := getIdealTree()
	if !ideal.Equals(tree) {
		t.Errorf("Non-ideal BVH created via add:\n%v\nIdeal:\n%v", tree.String(),
			ideal.String())
	}
}

func TestRemove(t *testing.T) {
	tree := getIdealTree()
	for index, orth := range leaf[:9] {
		fmt.Printf("%d-r, %v\n", index, orth)
		tree.Remove(orth)
		drawBVH(tree, fmt.Sprintf("test%d-r.png", index))
	}
}

func TestQuery(t *testing.T) {
}

func getIdealTree() *BVol {
	tree := &BVol{depth: 4,
		orth: &Orthotope{point: [d]int{2, 2}, delta: [d]int{21, 23}},
		vol: [2]*BVol{
			&BVol{depth: 3,
				orth: &Orthotope{point: [d]int{16, 2}, delta: [d]int{7, 23}},
				vol: [2]*BVol{
					&BVol{depth: 1,
						orth: &Orthotope{point: [d]int{18, 19}, delta: [d]int{5, 6}},
						vol: [2]*BVol{
							&BVol{orth: leaf[8]},
							&BVol{orth: leaf[9]},
						},
					},
					&BVol{depth: 2,
						orth: &Orthotope{point: [d]int{16, 2}, delta: [d]int{6, 12}},
						vol: [2]*BVol{
							&BVol{depth: 1,
								orth: &Orthotope{point: [d]int{16, 2}, delta: [d]int{5, 8}},
								vol: [2]*BVol{
									&BVol{orth: leaf[2]},
									&BVol{orth: leaf[3]},
								},
							},
							&BVol{depth: 1,
								orth: &Orthotope{point: [d]int{17, 12}, delta: [d]int{5, 2}},
								vol: [2]*BVol{
									&BVol{orth: leaf[6]},
									&BVol{orth: leaf[5]},
								},
							},
						},
					},
				},
			},
			&BVol{depth: 2,
				orth: &Orthotope{point: [d]int{2, 2}, delta: [d]int{10, 20}},
				vol: [2]*BVol{
					&BVol{depth: 1,
						orth: &Orthotope{point: [d]int{4, 11}, delta: [d]int{8, 11}},
						vol: [2]*BVol{
							&BVol{orth: leaf[4]},
							&BVol{orth: leaf[7]},
						},
					},
					&BVol{depth: 1,
						orth: &Orthotope{point: [d]int{2, 2}, delta: [d]int{8, 8}},
						vol: [2]*BVol{
							&BVol{orth: leaf[1]},
							&BVol{orth: leaf[0]},
						},
					},
				},
			},
		},
	}
	return tree
}
