package lc

/*
给你一个由 '1'（陆地）和 '0'（水）组成的的二维网格，请你计算网格中岛屿的数量。

岛屿总是被水包围，并且每座岛屿只能由水平方向和/或竖直方向上相邻的陆地连接形成。

此外，你可以假设该网格的四条边均被水包围。

示例 1：

输入：grid = [
  ['1','1','1','1','0'],
  ['1','1','0','1','0'],
  ['1','1','0','0','0'],
  ['0','0','0','0','0']
]
输出：1

示例 2：

输入：grid = [
  ['1','1','0','0','0'],
  ['1','1','0','0','0'],
  ['0','0','1','0','0'],
  ['0','0','0','1','1']
]
输出：3

提示：

- m == grid.length

- n == grid[i].length

- 1 <= m, n <= 300

- grid[i][j] 的值为 '0' 或 '1'
*/

func NumIslands(grid [][]byte) int {
	m,n := len(grid),len(grid[0])
	dirs := [][]int{
		{-1,0},{0,1},{1,0},{0,-1},
	}

	visited := make([][]bool,m)
	for i := range m {
		visited[i] = make([]bool,n)
	}

	var dfs func(r int,c int)
	dfs = func(r int,c int) {

		visited[r][c] = true
		for _,dir := range dirs {
			r1,c1 := dir[0] + r,dir[1] + c
			//边界
			if r1 < 0 || c1 < 0 || r1 >= m || c1 >= n {
				continue
			}
			//已访问或者本身不可达
			if visited[r1][c1] || grid[r1][c1] == '0' {
				continue
			}
			dfs(r1,c1)
		}
	}

	cnt := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == '1' && !visited[i][j] {
				dfs(i,j)
				cnt++
			}
		}
	}

	return cnt
}

