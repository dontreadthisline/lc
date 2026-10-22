package lc

// 063. 不同路径II
// 题面: ../mother-problems/0063-不同路径II.md
// 链接: https://leetcode.cn/problems/*/

func UniquePathsWithObstacles(grid [][]int) int {
	//这个题目比不同路径1做了如下几点优化
	//1 保持了代码简洁,把dp数组构造收到了for循环
	//2 把边界初始化收纳到了for循环
	//3 逻辑上保持和不同路径1尽量保持一致,只做有效值更新
	m ,n:= len(grid),len(grid[0])
	dp := make([][]int,m)
	for i := range dp {
		dp[i] = make([]int,n)
	}

	for i := range m {
		if grid[i][0] != 0 {
			break
		}
		dp[i][0] = 1
	}

	for i := range n {
		if grid[0][i] != 0 {
			break
		}
		dp[0][i] = 1
	}

	for i := range m {
		for j := range n {
			if i > 0 && j > 0 && grid[i][j] == 0{
				dp[i][j] = dp[i][j-1] + dp[i-1][j]
			}
		}
	}
	return dp[m-1][n-1]
}

