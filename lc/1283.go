package lc

import "slices"

// 1283. 使结果不超过阈值的最小除数
// 题面: ../mother-problems/1283-使结果不超过阈值的最小除数.md
// 链接: https://leetcode.cn/problems/*/
// 这个题目的根本难点在于向上/向下取整,使得必须每次找到divisor必须分开算
// 否则的话,如果是浮点数,就是一道简单的数学题,直接套公式算就行了 S / x <= threshold x >= S  / threshold
// 所以x的最小值就是 S /threshold

func SmallestDivisor(nums []int, threshold int) int {
	l,r := 1,slices.Max(nums)
	meet := func(divisor int) bool {
		s := 0
		for _,num := range nums {
			s += (num + divisor - 1) / divisor
		}
		return s <= threshold
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
