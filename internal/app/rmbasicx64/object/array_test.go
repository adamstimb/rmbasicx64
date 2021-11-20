package object

import "testing"

func TestCalculateAddressFromArraySubscripts(t *testing.T) {
	tests := []struct {
		subscripts []int
		bounds     []int
		expected   int
	}{
		{[]int{10}, []int{10}, 9},
		{[]int{10}, []int{1}, 0},
		{[]int{10}, []int{5}, 4},
		{[]int{5, 2}, []int{4, 1}, 9},
		{[]int{5, 2}, []int{4, 0}, 8},
		{[]int{5, 2, 2}, []int{4, 1, 1}, 19},
	}

	for _, tt := range tests {
		actual := calculateAddressFromArraySubscripts(tt.subscripts, tt.bounds)
		if actual != tt.expected {
			t.Errorf("CalculateAddressFromArraySubscripts(%v, %v) expected %v, got %v", tt.subscripts, tt.bounds, tt.expected, actual)
		}
	}
}
