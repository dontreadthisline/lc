package lc

import "slices"

/*
给定一个不含重复数字的数组 nums ，返回其 *所有可能的全排列* 。你可以 按任意顺序 返回答案。

示例 1：

输入：nums = [1,2,3]
输出：[[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]

示例 2：

输入：nums = [0,1]
输出：[[0,1],[1,0]]

示例 3：

输入：nums = [1]
输出：[[1]]

提示：

- 1 <= nums.length <= 6

- -10 <= nums[i] <= 10

- nums 中的所有整数 互不相同
*/

func Permute(nums []int) [][]int {
	var res [][]int
	n := len(nums)
	used := make([]bool, n) // 排列必须 used: 候选全域开放, 同一值可能被再次选中
	var backtrack func([]int)
	backtrack = func(path []int) {
		if len(path) == n {
			res = append(res, slices.Clone(path)) // 收获必须拷贝: path 底层数组共享
			return                                // 叶子处 used 全占用, 不 return 也空转; 留着语义更清晰
		}
		for i := range n {
			if !used[i] {
				used[i] = true
				backtrack(append(path, nums[i])) // append 传参: 本帧 path 未被赋值, 无需回退
				used[i] = false                  // 撤销: used 是共享数组, 必须复原(全函数唯一的可变共享状态)
			}
		}
	}
	backtrack(nil)
	return res
}

func Factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * Factorial(n-1)
}
