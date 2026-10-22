package lc

// 1870. 准时到达的列车最小时速
// 题面: ../mother-problems/1870-准时到达的列车最小时速.md
// 链接: https://leetcode.cn/problems/*/

func MinSpeedOnTime(dists []int, hour float64) int {
	l,r,n := 1,10_000_000,len(dists)
	if hour <= float64(n - 1) {
		return -1
	}
	meet := func(speed int) bool {
		h := 0
		total := 0.0
		for i,dist := range dists {
			if i !=  n - 1 {
				h += (dist + speed - 1) / speed
			} else {
				total = float64(h) + float64(dist) / float64(speed)
			}
		}
		return total <= hour
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
