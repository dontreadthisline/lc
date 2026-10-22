package lc

import "strconv"

// 091. 解码方法
// 题面: ../mother-problems/0091-解码方法.md
// 链接: https://leetcode.cn/problems/*/
//12,226
func NumDecodings(s string) int {
	n := len(s)
	toInt := func(s string) int {
		num,err := strconv.ParseInt(s,10,32)
		if err != nil {
			return -1
		}
		return int(num)
	}
	dp := make([]int,n+1)
	for i := 0; i <= n; i++ {
		for j := range 2 {
			if i - j > 0 {
				num := toInt(s[i-j:i])
				if num >= 0 && num <= 26 {
					dp[i] = 0
				}
			}
		}
	}
	return dp[n]
}

