// Package discreet contains convenience operations for working with integers.
// Copyright 2018 Brian Noyama. Subject to the the Apache License, Version 2.0.
package discreet

// SHIFT sets the number of bits for which Min and Max work.
const SHIFT uint = 32

// Min of i and j
func Min(i int, j int) int {
	return i ^ ((i ^ j) & ((j - i) >> SHIFT))
}

// Max of i and j
func Max(i int, j int) int {
	return i ^ ((i ^ j) & ((i - j) >> SHIFT))
}

// Abs absolute value of i
func Abs(i int) int {
	mask := (i >> SHIFT)
	return mask ^ (mask + i)
}

// Pow efficiently computes n^p
func Pow(n, p int) int {
	result := 1
	for 0 != p {
		if 0 != (p & 1) {
			result *= n
		}
		p >>= 1
		n *= n
	}
	return result
}
