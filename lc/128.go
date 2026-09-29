package lc

import "slices"

/*
给定一个未排序的整数数组 nums ，找出数字连续的最长序列（不要求序列元素在原数组中连续）的长度。

请你设计并实现时间复杂度为 O(n)* *的算法解决此问题。

示例 1：

输入：nums = [100,4,200,1,3,2]
输出：4
解释：最长数字连续序列是 [1, 2, 3, 4]。它的长度为 4。

示例 2：

输入：nums = [0,3,7,2,5,8,4,6,0,1]
输出：9

示例 3：

输入：nums = [1,0,1,2]
输出：3

提示：

- 0 <= nums.length <= 10^5

- -10^9 <= nums[i] <= 10^9
*/

//去重后排序,并查集,hashset

func LongestConsecutiveSortAlgs(nums []int) int {
	slices.Sort(nums)
	best, cur := 0, 0
	for i, v := range nums {
		if i > 0 && v == nums[i-1] {
			continue // 重复：原地踏步
		}
		if i > 0 && v == nums[i-1]+1 {
			cur++ // 连续：链延长
		} else {
			cur = 1 // 断链：新起点
		}
		best = max(best, cur)
	}
	return best
}

func LongestConsecutive(nums []int) int {
	return longestConsecutiveHashsetAlgs(nums)
}

func longestConsecutiveHashsetAlgs(nums []int) int {
	n, res := len(nums), 0
	m := make(map[int]int, n)
	for _, num := range nums {
		m[num] = 1
	}
	//1,2,4
	for num := range m {
		if m[num-1] <= 0 {
			j := 0
			for ; m[num+j] > 0; j++ {
			}
			res = max(res, j)
		}
	}
	return res
}
