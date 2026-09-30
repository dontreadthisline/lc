package lc_tests

import (
	"demo/lc"
	"testing"
)

// 005. 最长回文子串（多个合法答案时按长度+回文性验收）
func TestLongestPalindromicSubStr(t *testing.T) {
	isPal := func(s string) bool {
		for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
			if s[i] != s[j] {
				return false
			}
		}
		return true
	}
	cases := []struct {
		name      string
		s         string
		wantLen   int
		mustEqual string // 非空时要求精确匹配
	}{
		{"官方用例1", "babad", 3, ""},
		{"官方用例2", "cbbd", 2, "bb"},
		{"单字符", "a", 1, "a"},
		{"无重复", "ac", 1, ""},
		{"全同", "aaaa", 4, "aaaa"},
		{"尾部回文", "abcddcba", 8, "abcddcba"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got string
			panicked := safeCall(func() { got = lc.LongestPalindromicSubStr(c.s) })
			if panicked {
				t.Errorf("LongestPalindromicSubStr(%q) panic", c.s)
				return
			}
			if !isPal(got) {
				t.Errorf("结果 %q 不是回文", got)
				return
			}
			if len(got) != c.wantLen {
				t.Errorf("LongestPalindromicSubStr(%q) = %q (len %d), want len %d", c.s, got, len(got), c.wantLen)
				return
			}
			if c.mustEqual != "" && got != c.mustEqual {
				t.Errorf("唯一解场景 = %q, want %q", got, c.mustEqual)
			}
		})
	}
}
