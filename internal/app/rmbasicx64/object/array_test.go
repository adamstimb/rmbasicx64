package object

import "testing"

func TestCalculateAddressFromArraySubscripts(t *testing.T) {
	tests := []struct {
		bounds     []int
		subscripts []int
		expected   int
	}{
		{[]int{10}, []int{9}, 10},
		{[]int{10}, []int{0}, 1},
		{[]int{10}, []int{5}, 6},
		{[]int{5, 2}, []int{0, 0}, 0},
		{[]int{5, 2}, []int{4, 1}, 10},
		{[]int{10, 2}, []int{9, 1}, 20},
	}

	for _, tt := range tests {
		actual := calculateAddressFromArraySubscripts(tt.bounds, tt.subscripts)
		if actual != tt.expected {
			t.Errorf("CalculateAddressFromArraySubscripts(%v, %v) expected %v, got %v", tt.bounds, tt.subscripts, tt.expected, actual)
		}
	}
}
