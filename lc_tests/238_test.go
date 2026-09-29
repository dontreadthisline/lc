package lc_tests

import (
	"demo/lc"
	"math/rand"
	"slices"
	"testing"
)

// 238. 除自身以外数组的乘积
func TestProductExceptSelf(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want []int
	}{
		{"官方用例1", []int{1, 2, 3, 4}, []int{24, 12, 8, 6}},
		{"官方用例2_含零", []int{-1, 1, 0, -3, 3}, []int{0, 0, 9, 0, 0}},
		{"两个元素", []int{3, 2}, []int{2, 3}},
		{"零在前", []int{0, 4}, []int{4, 0}},
		{"零在后", []int{1, 0}, []int{0, 1}},
		{"两个零", []int{0, 0}, []int{0, 0}},
		{"负数", []int{2, -3}, []int{-3, 2}},
		{"五元素", []int{1, 2, 3, 4, 5}, []int{120, 60, 40, 30, 24}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := lc.ProductExceptSelf(c.nums); !slices.Equal(got, c.want) {
				t.Errorf("ProductExceptSelf(%v) = %v, want %v", c.nums, got, c.want)
			}
		})
	}
}

// 暴力对拍：O(n^2) 枚举 vs 前后缀积实现
func TestProductExceptSelfBruteCrossCheck(t *testing.T) {
	rng := rand.New(rand.NewSource(23))
	for iter := 0; iter < 5000; iter++ {
		n := 2 + rng.Intn(11)
		nums := make([]int, n)
		for i := range nums {
			nums[i] = rng.Intn(7) - 3
		}
		want := make([]int, n)
		for i := range want {
			p := 1
			for j, x := range nums {
				if j != i {
					p *= x
				}
			}
			want[i] = p
		}
		if got := lc.ProductExceptSelf(nums); !slices.Equal(got, want) {
			t.Fatalf("iter=%d nums=%v: got %v want %v", iter, nums, got, want)
		}
	}
}
