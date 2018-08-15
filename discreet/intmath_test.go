package discreet

import (
	"testing"
)

func TestMin(t *testing.T) {
  expected := -5
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
  expected := -5
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
