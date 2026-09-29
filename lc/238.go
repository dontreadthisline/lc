package lc

/*
给你一个整数数组 nums，返回 数组 answer ，其中 answer[i] 等于 nums 中除了 nums[i] 之外其余各元素的乘积 。

题目数据 保证 数组 nums之中任意元素的全部前缀元素和后缀的乘积都在  32 位 整数范围内。

请 不要使用除法，且在 O(n) 时间复杂度内完成此题。

示例 1:

输入: nums = [1,2,3,4]
输出: [24,12,8,6]

示例 2:

输入: nums = [-1,1,0,-3,3]
输出: [0,0,9,0,0]

提示：

- 2 <= nums.length <= 10^5

- -30 <= nums[i] <= 30

- 输入 保证 数组 answer[i] 在  32 位 整数范围内

进阶：你可以在 O(1) 的额外空间复杂度内完成这个题目吗？（ 出于对空间复杂度分析的目的，输出数组 不被视为 额外空间。）
*/
//left[i] = nums[0] * ... * nums[i]

func ProductExceptSelf(nums []int) []int {
	n := len(nums)
	left, right := make([]int, n), make([]int, n)
	left[0], right[n-1] = nums[0], nums[n-1]
	for i, j := 1, n-2; i <= n && j >= 0; i, j = i+1, j-1 {
		left[i] = left[i-1] * nums[i]
		right[j] = right[j+1] * nums[j]
	}

	res := make([]int, n)
	for i := range n {
		if i == 0 {
			res[i] = right[1]
		} else if i == n-1 {
			res[i] = left[n-2]
		} else {
			res[i] = left[i-1] * right[i+1]
		}
	}
	return res
}
