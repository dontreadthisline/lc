package lc

import "slices"

// 217. 存在重复元素
// 题面: ../mother-problems/0217-存在重复元素.md
func containsDuplicate(nums []int) bool {
	m := make(map[int]int, len(nums))
	for _, num := range nums {
		m[num] += 1
		if m[num] >= 2 {
			return true
		}
	}
	return false
}

func containsDuplicateSortAlgs(nums []int) bool {
	n := len(nums)
	slices.Sort(nums)
	for i := 1; i < n; i++ {
		if nums[i] == nums[i-1] {
			return true
		}
	}
	return false
}
