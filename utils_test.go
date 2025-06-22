package main

import (
	"reflect"
	"testing"
)

func TestDivisionBoundaries(t *testing.T) {
	tests := []struct {
		length   int
		parts    int
		expected [][]int
	}{
		{
			length: 10, parts: 2,
			expected: [][]int{{0, 4}, {5, 9}},
		},
		{
			length: 9, parts: 3,
			expected: [][]int{{0, 2}, {3, 5}, {6, 8}},
		},
		{
			length: 5, parts: 5,
			expected: [][]int{{0, 0}, {1, 1}, {2, 2}, {3, 3}, {4, 4}},
		},
		{
			length: 8, parts: 4,
			expected: [][]int{{0, 1}, {2, 3}, {4, 5}, {6, 7}},
		},
		{
			length: 1, parts: 1,
			expected: [][]int{{0, 0}},
		},
	}

	for _, tt := range tests {
		result := divisionBoundaries(tt.length, tt.parts)

		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf(
				"divisionBoundaries(%d, %d) = %v; want %v",
				tt.length,
				tt.parts,
				result,
				tt.expected,
			)
		}
	}
}

func TestDivisionBoundaries_InvalidParts(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("divisionBoundaries did not panic with parts=0")
		}
	}()

	divisionBoundaries(10, 0)
}
