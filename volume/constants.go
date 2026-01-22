package volume

type Num interface {
	~int | ~int32 | ~int16 | ~float64 | ~float32
}

const DIM = 3
