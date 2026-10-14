package lc

/*
给你一个字符串 s，请你将* *s* *分割成一些 子串，使每个子串都是 回文串 。返回 s 所有可能的分割方案。

示例 1：

输入：s = "aab"
输出：[["a","a","b"],["aa","b"]]

示例 2：

输入：s = "a"
输出：[["a"]]

提示：

- 1 <= s.length <= 16

- s 仅由小写英文字母组成
*/

func Partition(s string) [][]string {
	var res [][]string
	var dfs func(int,[]int)
	n := len(s)
	dfs = func(start int,path []int) {
		if start >= n {
		}
		dfs(start+1,append(path,start+1))//切
		dfs(start+1,path)
	}

	dfs(-1,[]int{})
	return res
}

func collect(s string,path []int) []string {
	l := len(path)
	res := make([]string,0,l+1)
	if len(path) == 0 {
		res = append(res,s)
	}
	for i := range l {
		if i == 0 {
			res = append(res,s[:path[0]])
		} else if i == l - 1 {
			res = append(res,s[path[i]:])
		} else  {
			res = append(res,s[i:i+1])
		}
	}
	return res
}
