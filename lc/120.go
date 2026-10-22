package lc

import "math"

// 120. 三角形最小路径和
// 题面: ../mother-problems/0120-三角形最小路径和.md
// 链接: https://leetcode.cn/problems/*/

func MinimumTotal(triangle [][]int) int {
	n := len(triangle)
	dp := make([][]int,n)
	for i := range n {
		dp[i] = make([]int,n)
	}
	res := math.MaxInt
	for i := range n {
		if i == 0 {
			dp[i] = triangle[i]
			continue
		}
		for j := range i {
			if j == 0 {
				dp[i][j] = dp[i-1][j] + triangle[i][j]
			} else {
				dp[i][j] = min(dp[i-1][j],dp[i-1][j-1]) + triangle[i][j]
			}
			if i == n -1 {
				res = min(res,dp[i][j])
			}
		}
	}
	return res
}
