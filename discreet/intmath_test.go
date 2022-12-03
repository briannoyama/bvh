// Copyright 2018 Brian Noyama. Subject to the the Apache License, Version 2.0.
package discreet

import (
	"testing"
)

func TestMin(t *testing.T) {
	expected := int32(-5)
	actual := Min(expected, 10)
	if actual != expected {
		t.Errorf("Expected %d, got %d.", expected, actual)
	}

	expected = 20
	actual = Min(24, expected)
	if actual != expected {
		t.Errorf("Expected %d, got %d.", expected, actual)
	}
}

func TestMax(t *testing.T) {
	expected := int32(-5)
	actual := Max(expected, -7)
	if actual != expected {
		t.Errorf("Expected %d, got %d.", expected, actual)
	}

	expected = 20
	actual = Max(4, expected)
	if actual != expected {
		t.Errorf("Expected %d, got %d.", expected, actual)
	}
}

func TestAbs(t *testing.T) {
	expected := int32(20)
	actual := Abs(-20)
	if actual != expected {
		t.Errorf("Expected %d, got %d.", expected, actual)
	}
}

func TestPow(t *testing.T) {
	expected := int32(177147)
	actual := Pow(3, 11)
	if actual != expected {
		t.Errorf("Expected %d, got %d.", expected, actual)
	}
}
