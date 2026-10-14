package lc

import "math"

/*
给你一个整数数组 nums ，请你找出一个具有最大和的连续子数组（子数组最少包含一个元素），返回其最大和。

子数组 是数组中的一个连续部分。

示例 1：

输入：nums = [-2,1,-3,4,-1,2,1,-5,4]
输出：6
解释：连续子数组 [4,-1,2,1] 的和最大，为 6 。

示例 2：

输入：nums = [1]
输出：1

示例 3：

输入：nums = [5,4,-1,7,8]
输出：23

提示：

- 1 <= nums.length <= 10^5

- -10^4 <= nums[i] <= 10^4

进阶：如果你已经实现复杂度为 O(n) 的解法，尝试使用更为精妙的 分治法 求解。
*/

func MaxSubArray(nums []int) int {
	panic("TODO")
}

func NaviWayMaxSubArray(nums []int) int {
	res := math.MinInt
	n := len(nums)
	var sum int
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			sum = 0
			for k := i; k <= j; k++ {
				sum += nums[k]
			}
			res = max(res,sum)
		}
	}
	return res
}

func PrefixSumMaxSubArray(nums []int) int {
	res := math.MinInt
	n := len(nums)
	prefix := make([]int,n)
	prefix[0] = nums[0]
	for i := 1; i < n; i++ {
		prefix[i] = prefix[i-1] + nums[i]
	}
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			res = max(res,prefix[j]- prefix[i]+nums[i])
		}
	}
	return res
}

func GreedDeltaSumMaxSubArray(nums []int) int {
	res := math.MinInt
	
	return res
}

func DivideConquerMaxSubArray(nums []int) int {
	panic("TODO") // 半成品临时补齐以便编译, 待实现
}
