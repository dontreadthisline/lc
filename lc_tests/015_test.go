package lc_tests

import (
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	"demo/lc"
)

// 归一化：三元组内排序 + 组间排序 + 去重，让任意顺序的产出可比
func normTriplets(v [][]int) [][]int {
	out := make([][]int, 0, len(v))
	seen := map[string]bool{}
	for _, t := range v {
		c := append([]int(nil), t...)
		sort.Ints(c)
		k := key(c)
		if !seen[k] {
			seen[k] = true
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		for k := 0; k < 3; k++ {
			if out[i][k] != out[j][k] {
				return out[i][k] < out[j][k]
			}
		}
		return false
	})
	return out
}

func key(t []int) string {
	parts := make([]string, len(t))
	for i, v := range t {
		parts[i] = strconv.Itoa(v)
	}
	return strings.Join(parts, ",")
}

func TestThreeSum(t *testing.T) {
	cases := []struct {
		nums []int
		want [][]int
	}{
		{[]int{-1, 0, 1, 2, -1, -4}, [][]int{{-1, -1, 2}, {-1, 0, 1}}},
		{[]int{0, 1, 1}, [][]int{}},
		{[]int{0, 0, 0}, [][]int{{0, 0, 0}}},
	}
	for _, c := range cases {
		want := normTriplets(c.want)
		for name, fn := range map[string]func([]int) [][]int{
			"BruteForce": lc.ThreeSumBruteForce,
			"SortAlgs":   lc.ThreeSumSortAlgs,
		} {
			got := normTriplets(fn(c.nums))
			if !reflect.DeepEqual(got, want) {
				t.Errorf("%s(%v) = %v, want %v", name, c.nums, got, want)
			}
		}
	}
}
