package lc

/*
在给定的 m x n 网格 grid 中，每个单元格可以有以下三个值之一：

- 值 0 代表空单元格；

- 值 1 代表新鲜橘子；

- 值 2 代表腐烂的橘子。

每分钟，腐烂的橘子 周围 4 个方向上相邻 的新鲜橘子都会腐烂。

返回 *直到单元格中没有新鲜橘子为止所必须经过的最小分钟数。如果不可能，返回 -1* 。

示例 1：

输入：grid = [[2,1,1],[1,1,0],[0,1,1]]
输出：4

示例 2：

输入：grid = [[2,1,1],[0,1,1],[1,0,1]]
输出：-1
解释：左下角的橘子（第 2 行， 第 0 列）永远不会腐烂，因为腐烂只会发生在 4 个方向上。

示例 3：

输入：grid = [[0,2]]
输出：0
解释：因为 0 分钟时已经没有新鲜橘子了，所以答案就是 0 。

提示：

- m == grid.length

- n == grid[i].length

- 1 <= m, n <= 10

- grid[i][j] 仅为 0、1 或 2
*/

//OrangesRotting 多源BFS
func OrangesRotting(grid [][]int) int {
	m,n := len(grid),len(grid[0])
	dirs := [][]int{
		{-1,0},{0,1},{1,0},{0,-1},
	}
	queue := make([][]int,0)
	for i := range m {
		for j := range n {
			if grid[i][j] == 2 { //已腐烂,多源bfs
				queue = append(queue,[]int{i,j})
			}
		}
	}
	var cnt int
	for len(queue) > 0 {
		r,c := queue[0][0],queue[0][1]
		for _,dir := range dirs {
			r1,c1 := r + dir[0],c + dir[1]
			if r1 < 0 || c1 < 0 || r1 >= m || c1 >= n {
				continue
			}
			if grid[r1][c1] != 1{ //已经腐烂或者本身就是空位
				continue
			}
			grid[r1][c1] = 2
			cnt += 1
		}
		queue = queue[1:len(queue)]
	}
	return cnt
}
