package lc

// 746. 使用最小花费爬楼梯
// 题面: ../mother-problems/0746-使用最小花费爬楼梯.md
// 链接: https://leetcode.cn/problems/*/

func MinCostClimbingStairs(cost []int) int {
	n := len(cost)
	dp := make([]int,n+1)
	//踩在第i台阶上所需要的最小花费,最后的答案一定是dp[n]
	//因为无论最后怎么到的终点（一步或者是两步) 最后的状态一定是踩在第n个台阶上(第n个台阶是我们加上的虚拟台阶)
	for i := 0; i <= n; i++ {
		if i == 0 {
			dp[i] = cost[0]
		} else if i == 1 {
			dp[i] = min(cost[0],cost[1])
		} else {
			dp[i] = min(dp[i-1],dp[i-2]+cost[i-1])
		}
	}
	return dp[n]
}
