package lc_tests

import (
	"testing"

	"demo/lc"
)

func TestLongestConsecutive(t *testing.T) {
	cases := []struct {
		nums []int
		want int
	}{
		{[]int{100, 4, 200, 1, 3, 2}, 4},
		{[]int{0, 3, 7, 2, 5, 8, 4, 6, 0, 1}, 9},
		{[]int{1, 2, 2, 3}, 3}, // 重复元素
		{[]int{5}, 1},
		{[]int{}, 0}, // 空数组
	}
	for _, c := range cases {
		if got := lc.LongestConsecutive(c.nums); got != c.want {
			t.Errorf("Hash(%v) = %d, want %d", c.nums, got, c.want)
		}
		if got := lc.LongestConsecutiveSortAlgs(c.nums); got != c.want {
			t.Errorf("Sort(%v) = %d, want %d", c.nums, got, c.want)
		}
	}
}
