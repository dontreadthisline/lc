package lc

import (
	"slices"
	"sort"
)

/*
给你一个整数数组 nums ，判断是否存在三元组 [nums[i], nums[j], nums[k]] 满足 i != j、i != k 且 j != k ，同时还满足 nums[i] + nums[j] + nums[k] == 0 。请你返回所有和为 0 且不重复的三元组。

注意：答案中不可以包含重复的三元组。

示例 1：

输入：nums = [-1,0,1,2,-1,-4]
输出：[[-1,-1,2],[-1,0,1]]
解释：
nums[0] + nums[1] + nums[2] = (-1) + 0 + 1 = 0 。
nums[1] + nums[2] + nums[4] = 0 + 1 + (-1) = 0 。
nums[0] + nums[3] + nums[4] = (-1) + 2 + (-1) = 0 。
不同的三元组是 [-1,0,1] 和 [-1,-1,2] 。
注意，输出的顺序和三元组的顺序并不重要。

示例 2：

输入：nums = [0,1,1]
输出：[]
解释：唯一可能的三元组和不为 0 。

示例 3：

输入：nums = [0,0,0]
输出：[[0,0,0]]
解释：唯一可能的三元组和为 0 。

提示：

- 3 <= nums.length <= 3000

- -10^5 <= nums[i] <= 10^5
*/

func ThreeSum(nums []int) [][]int {
	return ThreeSumSortAlgs(nums)
}

func ThreeSumBruteForce(nums []int) [][]int {
	res := make([][]int, 0)
	n := len(nums)
	for i := range n {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				if nums[i]+nums[j]+nums[k] == 0 {
					t := []int{nums[i], nums[j], nums[k]}
					slices.Sort(t)
					res = append(res, t)
				}
			}
		}
	}
	// 去重（暴力枚举的值级重复）
	dedup := res[:0]
	for idx, t := range res {
		if idx == 0 || t[0] != res[idx-1][0] || t[1] != res[idx-1][1] || t[2] != res[idx-1][2] {
			_ = idx
			dedup = append(dedup, t)
		}
	}
	return dedup
}

func ThreeSumSortAlgs(nums []int) [][]int {
	res, n := make([][]int, 0), len(nums)
	sort.Ints(nums)
	for i := range n {
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		j, k := i+1, n-1
		for j < k {
			s := nums[i] + nums[j] + nums[k]
			if s == 0 {
				res = append(res, []int{nums[i], nums[j], nums[k]})
				j, k = j+1, k-1
				for ; j < k && nums[j] == nums[j-1]; j++ {
				}
				for ; j < k && nums[k] == nums[k+1]; k-- {
				}
			} else if s < 0 {
				j++
			} else {
				k--
			}
		}
	}
	return res
}
