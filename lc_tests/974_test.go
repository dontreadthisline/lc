package lc_tests

import (
	"demo/lc"
	"math/rand"
	"testing"
)

// 974. 和可被 K 整除的子数组
func TestSubarraysDivByK(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		{"官方用例1", []int{4, 5, 0, -2, -3, 1}, 5, 7},
		{"官方用例2_单元素不整除", []int{5}, 9, 0},
		{"单元素整除", []int{5}, 5, 1},
		{"负数同余", []int{-1, 2, 9}, 2, 2},
		{"全负数", []int{-6, -3}, 3, 3},
		{"k为1_全部子数组", []int{1, 2}, 1, 3},
		{"零元素", []int{0, 0}, 4, 3},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := lc.SubarraysDivByK(c.nums, c.k); got != c.want {
				t.Errorf("SubarraysDivByK(%v, %d) = %d, want %d", c.nums, c.k, got, c.want)
			}
		})
	}
}

// 暴力对拍：O(n^2) 枚举 vs 前缀取模配对，含负数
func bruteDivByK(nums []int, k int) int {
	cnt := 0
	for i := 0; i < len(nums); i++ {
		s := 0
		for j := i; j < len(nums); j++ {
			s += nums[j]
			if s%k == 0 {
				cnt++
			}
		}
	}
	return cnt
}

func TestSubarraysDivByKBruteCrossCheck(t *testing.T) {
	rng := rand.New(rand.NewSource(9))
	for iter := 0; iter < 20000; iter++ {
		n := rng.Intn(13)
		nums := make([]int, n)
		for i := range nums {
			nums[i] = rng.Intn(11) - 5
		}
		k := 1 + rng.Intn(7)
		if want := bruteDivByK(nums, k); lc.SubarraysDivByK(nums, k) != want {
			t.Fatalf("iter=%d nums=%v k=%d: 实现=%d 暴力=%d", iter, nums, k, lc.SubarraysDivByK(nums, k), want)
		}
	}
}
