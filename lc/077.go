package lc

import "slices"

/*
给定两个整数 n 和 k，返回范围 [1, n] 中所有可能的 k 个数的组合。

你可以按 任何顺序 返回答案。

示例 1：

输入：n = 4, k = 2
输出：
[
  [2,4],
  [3,4],
  [2,3],
  [1,2],
  [1,3],
  [1,4],
]

示例 2：

输入：n = 1, k = 1
输出：[[1]]

提示：

- 1 <= n <= 20

- 1 <= k <= n
*/

func Combine(n int, k int) [][]int {
	var res [][]int
	var backtrack func(int,[]int)
	backtrack = func(start int,path []int) {
		if len(path) == k {
			res = append(res,slices.Clone(path))
			return// 这里的return与否其实都不影响
		}
		for i:= start+1; i <= n; i++ {
			backtrack(i,append(path, i))
		}
	}
	backtrack(0,[]int{})
	return res
}

func Permunate(n int ) [][]int {
	var res [][]int
	var backtrack func([]int)
	used := make([]bool,n+1)
	backtrack = func(path []int) {
		if len(path) == n {
			item := make([]int,n)
			copy(item,path)
			res = append(res,item)
			return//这里的return与否都不影响
		}
		for i:= 1; i <= n; i++ {
			if !used[i] {
				used[i] = true
				backtrack(append(path, i))
				used[i] = false
			}
		}
	}
	backtrack([]int{})
	return res
}
