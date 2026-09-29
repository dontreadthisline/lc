package lc_tests

import (
	"demo/lc"
	"reflect"
	"sort"
	"testing"
)

func TestTwoSum(t *testing.T) {
	cases := []struct {
		nums   []int
		target int
		want   []int
	}{
		{[]int{2, 7, 11, 15}, 9, []int{0, 1}},
		{[]int{3, 2, 4}, 6, []int{1, 2}},
		{[]int{3, 3}, 6, []int{0, 1}},
		{[]int{-3, 4, 3, 90}, 0, []int{0, 2}},
	}
	for _, c := range cases {
		want := append([]int(nil), c.want...)
		sort.Ints(want)

		for name, fn := range map[string]func([]int, int) []int{
			"map":  lc.TwoSumMapAlgs,
			"sort": lc.TwoSumSortAlgs,
		} {
			got := fn(c.nums, c.target)
			if len(got) != 2 {
				t.Fatalf("%s(%v, %d) = %v, 应返回两个下标", name, c.nums, c.target, got)
			}
			pair := append([]int(nil), got...)
			sort.Ints(pair)
			if !reflect.DeepEqual(pair, want) {
				t.Errorf("%s(%v, %d) = %v, want %v", name, c.nums, c.target, got, want)
			}
			if c.nums[got[0]]+c.nums[got[1]] != c.target {
				t.Errorf("%s(%v, %d) = %v, 两数之和不为 %d", name, c.nums, c.target, got, c.target)
			}
		}
	}
}
