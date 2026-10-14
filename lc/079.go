package lc

/*
给定一个 m x n 二维字符网格 board 和一个字符串单词 word 。如果 word 存在于网格中，返回 true ；否则，返回 false 。

单词必须按照字母顺序，通过相邻的单元格内的字母构成，其中“相邻”单元格是那些水平相邻或垂直相邻的单元格。同一个单元格内的字母不允许被重复使用。

示例 1：

输入：board = [['A','B','C','E'],['S','F','C','S'],['A','D','E','E']], word = "ABCCED"
输出：true

示例 2：

输入：board = [['A','B','C','E'],['S','F','C','S'],['A','D','E','E']], word = "SEE"
输出：true

示例 3：

输入：board = [['A','B','C','E'],['S','F','C','S'],['A','D','E','E']], word = "ABCB"
输出：false

提示：

- m == board.length

- n = board[i].length

- 1 <= m, n <= 6

- 1 <= word.length <= 15

- board 和 word 仅由大小写英文字母组成

进阶：你可以使用搜索剪枝的技术来优化解决方案，使其在 board 更大的情况下可以更快解决问题？
../lc/200.go
图类型: 普通图
*/

func Exist(board [][]byte, word string) bool {
	m, n := len(board), len(board[0])
	dirs := [][]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
	visited := make([][]bool, m)
	for i := range m {
		visited[i] = make([]bool, n)
	}
	var found bool
	var backtrack func(int, int, int)
	backtrack = func(r, c int, start int) {
		if board[r][c] != word[start] { //没有路径
			return
		}
		if found || start >= len(word) - 1 {
			found = true
			return
		}
		visited[r][c] = true
		for _, dir := range dirs {
			r1, c1 := r+dir[0], c+dir[1]
			if r1 < 0 || c1 < 0 || r1 >= m || c1 >= n || visited[r1][c1] {
				continue
			}
			backtrack(r1, c1, start+1)
		}
		visited[r][c] = false //在这里撤销回溯
	}

	for i, row := range board {
		for j := range row {
			backtrack(i, j, 0)
			if found {
				return found
			}
		}
	}
	return false
}
