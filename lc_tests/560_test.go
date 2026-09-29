package lc_tests

import (
	"demo/lc"
	"math/rand"
	"testing"
)

// 560. 和为 K 的子数组（含负数，滑窗无资格的主场）
func TestSubarraySum(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		{"官方用例1", []int{1, 1, 1}, 2, 2},
		{"官方用例2", []int{1, 2, 3}, 3, 2},
		{"单元素等于k", []int{2}, 2, 1},
		{"单元素不等于k", []int{2}, 3, 0},
		{"含负数", []int{2, -2, 2}, 2, 3},
		{"k为零_正负抵消", []int{1, -1, 1, -1}, 0, 4},
		{"k为零_空段配对", []int{0, 0}, 0, 3},
		{"负k", []int{-1, -1, 1}, -2, 1},
		{"无解", []int{1, 2, 3}, 100, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := lc.SubarraySumHashway(c.nums, c.k); got != c.want {
				t.Errorf("SubarraySumHashway(%v, %d) = %d, want %d", c.nums, c.k, got, c.want)
			}
			if got := lc.SubarraySum(c.nums, c.k); got != c.want {
				t.Errorf("SubarraySum(暴力, %v, %d) = %d, want %d", c.nums, c.k, got, c.want)
			}
		})
	}
}

// 暴力对拍：O(n^2) 枚举 vs 前缀+哈希，含负数
func TestSubarraySumBruteCrossCheck(t *testing.T) {
	rng := rand.New(rand.NewSource(5))
	for iter := 0; iter < 20000; iter++ {
		n := rng.Intn(13)
		nums := make([]int, n)
		for i := range nums {
			nums[i] = rng.Intn(11) - 5
		}
		k := rng.Intn(13) - 6
		want := lc.SubarraySum(nums, k)
		if got := lc.SubarraySumHashway(nums, k); got != want {
			t.Fatalf("iter=%d nums=%v k=%d: 哈希=%d 暴力=%d", iter, nums, k, got, want)
		}
	}
}
