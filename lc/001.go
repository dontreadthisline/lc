package lc

import (
	"sort"
)

func TwoSum(nums []int, target int) []int {
	return TwoSumSortAlgs(nums, target)
}

func TwoSumMapAlgs(nums []int, target int) []int {
	m := make(map[int]int, len(nums))
	for i, num := range nums {
		if j, ok := m[target-num]; ok {
			return []int{i, j} //先查后插,能够保证target-num查到的一定不是自己
		}
		m[num] = i
	}
	return []int{-1, -1}
}

func TwoSumSortAlgs(nums []int, target int) []int {
	res := []int{-1, -1}
	n := len(nums)
	idxs := make([]int, n)
	for i := range n {
		idxs[i] = i
	}
	sort.Slice(idxs, func(i int, j int) bool {
		return nums[idxs[i]] < nums[idxs[j]]
	})
	i, j := 0, n-1
	for i < j {
		s := nums[idxs[i]] + nums[idxs[j]]
		if s == target {
			res = []int{idxs[i], idxs[j]}
			break
		} else if s > target {
			j -= 1
		} else {
			i += 1
		}
	}
	return res
}
