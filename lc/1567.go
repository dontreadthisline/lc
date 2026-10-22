package lc

// 1567. 乘积为正数的最长子数组长度
// 题面: ../mother-problems/1567-乘积为正数的最长子数组长度.md
// 链接: https://leetcode.cn/problems/*/

func GetMaxLen(nums []int) int {
	res,n := 0,len(nums)
	dp := make([][2]int,n+1)
	for i := 1; i <= n; i++ {
		if num := nums[i-1]; num > 0 {

		} else if num < 0 {

		}
	}
	_ = dp // 半成品临时占位, 待实现
	return res
}
