package lc

func longestPalindromicSubStr(s string) string {
	//暴力方法 o(n^3)
	palindromic := func(i, j int) bool {
		for ; i < j; i, j = i+1, j-1 {
			if s[i] != s[j] {
				return false
			}
		}
		return true
	}
	n := len(s)
	max, l, r := 0, 0, 0
	for i := range n {
		for j := i + 1; j < n; j++ {
			if palindromic(i, j) && j-i+1 > max {
				max = j - i + 1
				l, r = i, j
			}
		}
	}
	return s[l : r+1]
}
