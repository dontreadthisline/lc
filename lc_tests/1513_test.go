package lc_tests

import (
	"math/rand"
	"strings"
	"testing"

	"demo/lc"
)

// ---- 1513. 仅含 1 的子串数 ----

// refNumOfOnesSubStr: 遍历每段连续 1, 长 L 的段贡献 L*(L+1)/2 个子串, 取模
func refNumOfOnesSubStr(s string) int {
	const mod = 1_000_000_007
	total, cur := 0, 0
	for _, c := range s {
		if c == '1' {
			cur++
			total = (total + cur) % mod // 以当前位结尾的 1 子串数 = cur, 累加即段内组合数
		} else {
			cur = 0
		}
	}
	return total
}

func TestNumOfOnesSubStr(t *testing.T) {
	cases := []struct {
		s    string
		want int
	}{
		{"0110111", 9},
		{"101", 2},
		{"111111", 21},
		{"000", 0},
		{"", 0},
		{"1", 1},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.NumOfOnesSubStr(c.s) })
		if got != c.want {
			t.Errorf("NumOfOnesSubStr(%q) = %d, want %d", c.s, got, c.want)
		}
	}
}

func TestNumOfOnesSubStrRandom(t *testing.T) {
	for i := 0; i < 5000; i++ {
		n := rand.Intn(40)
		var sb strings.Builder
		for j := 0; j < n; j++ {
			if rand.Intn(2) == 0 {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		}
		s := sb.String()
		want := refNumOfOnesSubStr(s)
		var got int
		runSafely(t, func() { got = lc.NumOfOnesSubStr(s) })
		if got != want {
			t.Fatalf("NumOfOnesSubStr(%q) = %d, want %d", s, got, want)
		}
	}
}
