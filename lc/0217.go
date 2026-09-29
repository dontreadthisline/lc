package lc

/*
给你一个整数数组 nums 。如果任一值在数组中出现 至少两次 ，返回 true ；如果数组中每个元素互不相同，返回 false 。

示例 1：

输入：nums = [1,2,3,1]

输出：true

解释：

元素 1 在下标 0 和 3 出现。

示例 2：

输入：nums = [1,2,3,4]

输出：false

解释：

所有元素都不同。

示例 3：

输入：nums = [1,1,1,3,3,4,3,2,4,2]

输出：true

提示：

- 1 <= nums.length <= 10^5

- -10^9 <= nums[i] <= 10^9
*/

import "slices"

func FindMoreThanTwoRepeateNum(nums []int) bool {
	m := make(map[int]int, len(nums))
	for _, num := range nums {
		m[num] += 1
		if m[num] >= 2 {
			return true
		}
	}
	return false
}

func FindMoreThanTwoRepeateNumsSortAlgs(nums []int) bool {
	slices.Sort(nums)
	for i, num := range nums {
		//热点路径,让判断有价值一些
		if i >= 1 {
			if num == nums[i-1] {
				return true
			}
		}
	}
	return false
}
