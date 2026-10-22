package lc

import (
	"math"
)

// 2594. 修车的最少时间
// 题面: ../mother-problems/2594-修车的最少时间.md
// 链接: https://leetcode.cn/problems/*/

func RepairCars(ranks []int, cars int) int64 {
	var l,r int64
	l,r = 1,math.MaxInt64
	meet := func(minute float64) bool {
		sum := 0
		for _,rank:= range ranks {
			sum += int(math.Sqrt(minute/float64(rank)))
		}
		return sum >= cars
	}
	for l < r {
		mid := l + (r - l)  / 2
		if meet(float64(mid)) {
			r = mid
		} else {
			l = mid + 1
		}
	}
	return l
}
