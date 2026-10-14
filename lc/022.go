package lc

import (
	"slices"
)

/*
数字 n 代表生成括号的对数，请你设计一个函数，用于能够生成所有可能的并且 有效的 括号组合。

示例 1：

输入：n = 3
输出：["((()))","(()())","(())()","()(())","()()()"]

示例 2：

输入：n = 1
输出：["()"]

提示：

- 1 <= n <= 8
*/

func GenerateParenthesis(n int) []string {
	var res []string
	var backtrack func([]byte)
	//这样其实和dfs无差异,并且无剪枝,复杂度就是完全的2^n
	chars := []byte{'(',')'}
	backtrack = func(path []byte) {
		if len(path) >= 2*n {
			if valid(path) {
				res = append(res,string(slices.Clone(path)))
			}
			return
		}
		//这里可以提前剪枝
		for _,c := range chars {
			if check(path,c) {
				backtrack(append(path,c))
			}
		}
	}
	backtrack([]byte{})
	return res
}

func check(path []byte, c byte) bool {
	return true
}

func valid(path []byte) bool {
	stack := make([]byte,0,len(path)*2)
	for _,c := range path {
		if c == '(' {
			stack = append(stack,c)
		} else {
			if len(stack) == 0 || stack[len(stack)-1] != '(' {
				return false
			} else {
				stack = stack[0:len(stack)-1]
			}
		}
	}
	return len(stack) == 0
}
