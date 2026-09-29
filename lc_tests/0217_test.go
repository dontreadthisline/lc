package lc_tests

import (
	"testing"

	"demo/lc"
)

func TestFindMoreThanTwoRepeateNum(t *testing.T) {
	cases := []struct {
		nums []int
		want bool
	}{
		{[]int{1, 2, 3, 1}, true},
		{[]int{1, 2, 3, 4}, false},
		{[]int{1, 1, 1, 3, 3, 4, 3, 2, 4, 2}, true},
		{[]int{}, false},
		{[]int{7}, false},
	}
	for _, c := range cases {
		if got := lc.FindMoreThanTwoRepeateNum(c.nums); got != c.want {
			t.Errorf("Hash(%v) = %v, want %v", c.nums, got, c.want)
		}
		if got := lc.FindMoreThanTwoRepeateNumsSortAlgs(c.nums); got != c.want {
			t.Errorf("Sort(%v) = %v, want %v", c.nums, got, c.want)
		}
	}
}
