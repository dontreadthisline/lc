package lc_tests

import (
	"math/rand"
	"slices"
	"testing"

	"demo/lc"
)

func bruteContainsDuplicate(nums []int) bool {
	for i := range nums {
		for j := i + 1; j < len(nums); j++ {
			if nums[i] == nums[j] {
				return true
			}
		}
	}
	return false
}

func TestContainsDuplicate(t *testing.T) {
	cases := []struct {
		nums []int
		want bool
	}{
		{[]int{}, false},
		{[]int{1}, false},
		{[]int{1, 2, 3, 1}, true},
		{[]int{1, 2, 3, 4}, false},
		{[]int{1, 1, 1, 3, 3, 4, 3, 2, 4, 2}, true},
		{[]int{-5, -5}, true},
	}
	for _, c := range cases {
		if got := lc.ContainsDuplicate(c.nums); got != c.want {
			t.Errorf("ContainsDuplicate(%v) = %v, want %v", c.nums, got, c.want)
		}
	}
}

func TestContainsDuplicateRandom(t *testing.T) {
	for i := 0; i < 20000; i++ {
		n := rand.Intn(31)
		nums := make([]int, n)
		for j := range nums {
			nums[j] = rand.Intn(31) - 15
		}
		want := bruteContainsDuplicate(nums)
		if got := lc.ContainsDuplicate(nums); got != want {
			t.Fatalf("ContainsDuplicate(%v) = %v, want %v", nums, got, want)
		}
		// Sort 版会原地排序, 必须传副本
		cp := slices.Clone(nums)
		if got := lc.ContainsDuplicateSortAlgs(cp); got != want {
			t.Fatalf("ContainsDuplicateSortAlgs(%v) = %v, want %v", nums, got, want)
		}
	}
}
