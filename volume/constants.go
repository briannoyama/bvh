package volume

type Num interface {
	~int | ~int32 | ~int16 | ~float64 | ~float32
}

const DIM = 3

type Volume[K any, D any] interface {
	Contains(k1 K) bool
	Intersects(k1 K, d D) float32
	Minbound(k1 K) K
	Overlaps(k1 K) bool
	Score() float32
}
