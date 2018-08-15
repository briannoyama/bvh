package discreet

// Works for 32 bit or smaller ints.
const SHIFT uint = 32

func Min(i int, j int) int {
	return i ^ ((i ^ j) & ((j - i) >> SHIFT))
}

func Max(i int, j int) int {
	return i ^ ((i ^ j) & ((i - j) >> SHIFT))
}
