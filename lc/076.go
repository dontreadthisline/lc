package lc

// 76. 最小覆盖子串
// 题面: ../mother-problems/0076-最小覆盖子串.md
// 滑窗（流派 B：求最短合法窗口）：合法 → 记录 + 左缩找更短
// 判定：窗口内每个字符 c 的出现次数 ≥ t 中 c 的出现次数

func MinWindow(s string, t string) string {
	need := make(map[byte]int, len(t))
	for i := 0; i < len(t); i++ {
		need[t[i]]++
	}

	cnt := make(map[byte]int, len(s))
	covered := func() bool {
		for c, n := range need {
			if cnt[c] < n {
				return false
			}
		}
		return true
	}

	start, res := -1, len(s)+1
	i := 0
	for j := 0; j < len(s); j++ {
		cnt[s[j]]++
		for covered() && i <= j {
			if j-i+1 < res {
				res = j - i + 1
				start = i
			}
			cnt[s[i]]--
			i++
		}
	}
	if start == -1 {
		return ""
	}
	return s[start : start+res]
}
