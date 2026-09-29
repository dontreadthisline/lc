package lc

// 滑动窗口万能模板：两个流派，骨架一样，方向相反
//
// 流派 A：求最长合法窗口（3、424、904）
//   不合法 → 左缩（把窗口缩回合法范围）
//   合法   → 记录答案
//
// 流派 B：求最短合法窗口（209、76）
//   合法   → 记录答案
//   不合法 → 右扩（加元素凑够）
//
// 区别只在收缩条件是"不合法触发缩"还是"合法触发缩"。

// ---- 流派 A：求最长合法窗口 ----
// 扩 → 不合法时缩 → 合法时记
func MaxValidWindow(s string, isInvalid func(lo int, hi int) bool) {
	lo := 0
	for hi := 0; hi < len(s); hi++ {
		// 扩：右指针前进一格
		_ = s[hi]

		// 缩：不合法时左缩
		for isInvalid(lo, hi) {
			lo++
		}

		// 合法：更新答案（此处窗口 [lo, hi] 满足条件）
		_ = hi - lo + 1
	}
}

// ---- 流派 B：求最短合法窗口 ----
// 扩 → 合法时缩（找更短）→ 不合法时停
func MinValidWindow(s string, isValid func(lo int, hi int) bool) int {
	lo := 0
	res := len(s) + 1
	for hi := 0; hi < len(s); hi++ {
		// 扩：右指针前进一格
		_ = s[hi]

		// 缩：合法时左缩（找更短的合法窗口）
		for isValid(lo, hi) {
			if hi-lo+1 < res {
				res = hi - lo + 1
			}
			lo++
		}
	}
	if res > len(s) {
		return 0
	}
	return res
}
