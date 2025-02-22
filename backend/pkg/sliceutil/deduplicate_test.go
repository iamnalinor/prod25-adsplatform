package sliceutil

import (
	"reflect"
	"testing"
)

func TestDeduplicateLast(t *testing.T) {
	type testCase struct {
		arg  []int
		want []int
	}
	tests := []testCase{
		{[]int{}, []int{}},
		{[]int{1, 2, 3}, []int{1, 2, 3}},
		{[]int{1, 2, 4, 2, 3}, []int{1, 4, 2, 3}},
		{[]int{5, 5, 5, 5, 42, 42, 42}, []int{5, 42}},
	}
	for _, tt := range tests {
		if got := DeduplicateLast(tt.arg); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("DeduplicateLast() = %v, want %v", got, tt.want)
		}
	}
}
