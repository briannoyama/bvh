package rect
/*
import (
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
	"testing"
)

func TestBVHString(t *testing.T) {
	tree := &BVOrth{
		orth: &Orthotope{
			point: []int{2, 2},
			delta: []int{8, 8},
		},
		desc: [2]*BVOrth{
			&BVOrth{
				orth: &Orthotope{
					point: []int{2, 2},
					delta: []int{2, 2},
				},
			},
			&BVOrth{
				orth: &Orthotope{
					point: []int{7, 7},
					delta: []int{3, 3},
				},
			},
		},
		depth: 1,
	}
	ansLines := []string{"Point [2 2], Delta [8 8]",
		" Point [2 2], Delta [2 2]",
		" Point [7 7], Delta [3 3]",
		"",
	}

	answer := strings.Join(ansLines, "\n")
	result := tree.String()
	if answer != result {
		t.Errorf("Unexpected String:\n%v\nExpected:\n%v", result, answer)
	}
}

func TestAdd(t *testing.T) {
	leaf := [10]*Orthotope{
		&Orthotope{point: []int{2, 2}, delta: []int{2, 2}},
		&Orthotope{point: []int{7, 7}, delta: []int{3, 3}},
		&Orthotope{point: []int{19, 2}, delta: []int{2, 2}},
		&Orthotope{point: []int{16, 6}, delta: []int{3, 4}},
		&Orthotope{point: []int{10, 11}, delta: []int{2, 2}},
		&Orthotope{point: []int{17, 12}, delta: []int{2, 2}},
		&Orthotope{point: []int{20, 12}, delta: []int{2, 2}},
		&Orthotope{point: []int{4, 16}, delta: []int{6, 6}},
		&Orthotope{point: []int{18, 21}, delta: []int{2, 2}},
		&Orthotope{point: []int{19, 19}, delta: []int{4, 6}},
	}

	tree := &BVOrth{orth: leaf[0]}
	for _, orth := range leaf[1:] {
		tree.Add(orth)
	}

	tree.Group()
	ideal := getIdealTree()
	ideal.Group()
	answer := ideal.String()
	result := tree.String()
	if answer != result {
		t.Errorf("Non-ideal BVH created via add:\n%v\nIdeal:\n%v", result, answer)
	}
}

func TestQuery(t *testing.T) {
}

func TestRemove(t *testing.T) {
}

func getIdealTree() *BVOrth {
	tree := &BVOrth{
		depth: 4,
		orth:  &Orthotope{point: []int{2, 2}, delta: []int{21, 23}},
		desc: [2]*BVOrth{
			&BVOrth{
				depth: 2,
				orth:  &Orthotope{point: []int{2, 2}, delta: []int{10, 20}},
				desc: [2]*BVOrth{
					&BVOrth{
						depth: 1,
						orth:  &Orthotope{point: []int{4, 11}, delta: []int{8, 11}},
						desc: [2]*BVOrth{
							&BVOrth{
								orth: &Orthotope{point: []int{4, 16}, delta: []int{6, 6}},
							},
							&BVOrth{
								orth: &Orthotope{point: []int{10, 11}, delta: []int{2, 2}},
							},
						},
					},
					&BVOrth{
						depth: 1,
						orth:  &Orthotope{point: []int{2, 2}, delta: []int{8, 8}},
						desc: [2]*BVOrth{
							&BVOrth{
								orth: &Orthotope{point: []int{2, 2}, delta: []int{2, 2}},
							},
							&BVOrth{
								orth: &Orthotope{point: []int{7, 7}, delta: []int{3, 3}},
							},
						},
					},
				},
			},
			&BVOrth{
				depth: 3,
				orth:  &Orthotope{point: []int{16, 2}, delta: []int{7, 23}},
				desc: [2]*BVOrth{
					&BVOrth{
						depth: 2,
						orth:  &Orthotope{point: []int{17, 12}, delta: []int{6, 13}},
						desc: [2]*BVOrth{
							&BVOrth{
								depth: 1,
								orth:  &Orthotope{point: []int{17, 12}, delta: []int{5, 2}},
								desc: [2]*BVOrth{
									&BVOrth{
										orth: &Orthotope{point: []int{17, 12}, delta: []int{2, 2}},
									},
									&BVOrth{
										orth: &Orthotope{point: []int{20, 12}, delta: []int{2, 2}},
									},
								},
							},
							&BVOrth{
								depth: 1,
								orth:  &Orthotope{point: []int{18, 19}, delta: []int{5, 6}},
								desc: [2]*BVOrth{
									&BVOrth{
										orth: &Orthotope{point: []int{18, 21}, delta: []int{2, 2}},
									},
									&BVOrth{
										orth: &Orthotope{point: []int{19, 19}, delta: []int{4, 6}},
									},
								},
							},
						},
					},
					&BVOrth{
						depth: 1,
						orth:  &Orthotope{point: []int{16, 2}, delta: []int{5, 8}},
						desc: [2]*BVOrth{
							&BVOrth{
								orth: &Orthotope{point: []int{19, 2}, delta: []int{2, 2}},
							},
							&BVOrth{
								orth: &Orthotope{point: []int{16, 6}, delta: []int{3, 4}},
							},
						},
					},
				},
			},
		},
	}
	return tree
}

func drawBVH(bvOrth *BVOrth, name string) {
	myimage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{25, 25}})
	iter := bvOrth.Iterator()
	for next := iter.Next(); next != nil; next = iter.Next() {
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
*/
