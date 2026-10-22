package lc

import (
	"slices"
)

// 1482. 制作m束花所需的最少天数
// 题面: ../mother-problems/1482-制作m束花所需的最少天数.md
// 链接: https://leetcode.cn/problems/*/

func MinDays(days []int, m int, k int) int {
	l,r ,n := 1,slices.Max(days),len(days)
	if m * k > n {
		return -1
	}
	meet := func(d int) bool {
		var total ,streak int
		for _,day := range days {
			if d >= day {
				streak++
				if streak == k {
					total++
					streak = 0
				}
			} else {
				streak = 0
			}
		}
		return total >= m
	}
	for l < r {
		mid := l + (r - l) / 2
		if meet(mid) {
			r = mid
		} else {
			l = mid + 1
		}
	}
	return l
}
