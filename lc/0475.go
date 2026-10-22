package lc

import "slices"

// 0475. 供暖器
// 题面: ../mother-problems/0475-供暖器.md
// 链接: https://leetcode.cn/problems/*/
func abs(a int) int {
	if a >= 0 {
		return a
	}
	return -a
}
func FindRadius(houses []int, heaters []int) int {
	slices.Sort(heaters) //排序,便于二分
	n,l,r := len(heaters),0,len(houses)
	meet := func(radius int) bool {
		for _,house := range houses {
			//查询一个当前离自己最近的加热器 右端
			//距离大于radius 直接return false
			lo := LowerBound(heaters,house)
			if lo >= n {
				lo = n - 1
			}
			var t1,t2 bool
			t1 = abs(house - heaters[lo]) <= radius //右边界
			if lo > 0 {
				t2 = abs(house - heaters[lo-1]) <= radius //左边界
			}
			if !t1 && !t2 {
				return false
			}
		}
		return true
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

