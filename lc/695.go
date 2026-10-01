package lc

/*
给你一个大小为 m x n 的二进制矩阵 grid 。

岛屿 是由一些相邻的 1 (代表土地) 构成的组合，这里的「相邻」要求两个 1 必须在 水平或者竖直的四个方向上 相邻。你可以假设 grid 的四个边缘都被 0（代表水）包围着。

岛屿的面积是岛上值为 1 的单元格的数目。

计算并返回 grid 中最大的岛屿面积。如果没有岛屿，则返回面积为 0 。

示例 1：

输入：grid = [[0,0,1,0,0,0,0,1,0,0,0,0,0],[0,0,0,0,0,0,0,1,1,1,0,0,0],[0,1,1,0,1,0,0,0,0,0,0,0,0],[0,1,0,0,1,1,0,0,1,0,1,0,0],[0,1,0,0,1,1,0,0,1,1,1,0,0],[0,0,0,0,0,0,0,0,0,0,1,0,0],[0,0,0,0,0,0,0,1,1,1,0,0,0],[0,0,0,0,0,0,0,1,1,0,0,0,0]]
输出：6
解释：答案不应该是 11 ，因为岛屿只能包含水平或垂直这四个方向上的 1 。

示例 2：

输入：grid = [[0,0,0,0,0,0,0,0]]
输出：0

提示：

- m == grid.length

- n == grid[i].length

- 1 <= m, n <= 50

- grid[i][j] 为 0 或 1
*/

func MaxAreaOfIsland(grid [][]int) int {
	m,n := len(grid),len(grid[0])
	dirs := [][]int{
		{-1,0},{0,1},{1,0},{0,-1},
	}

	visited := make([][]bool,m)
	for i := range m {
		visited[i] = make([]bool,n)
	}

	var dfs func(r,c int,cnt *int)
	dfs = func(r,c int,cnt *int) {
		visited[r][c] = true
		*cnt = *cnt + 1
		for _,dir := range dirs {
			r1,c1 := dir[0] + r,dir[1] + c
			if r1 < 0 || c1 < 0 || r1 >= m || c1 >= n {
				continue
			}
			if visited[r1][c1] || grid[r1][c1] == 0 {
				continue
			}
			dfs(r1,c1,cnt)
		}
	}

	res := 0
	for i := range m {
		for j := range n {
			if grid[i][j] == 1 && !visited[i][j] {
				cnt := 0
				dfs(i,j,&cnt)
				res = max(res,cnt)
			}
		}
	}
	return res
}
