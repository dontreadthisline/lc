package lc_tests

import (
	"demo/lc"
	"testing"
)

// 278. 第一个错误的版本（mock：isBadVersion 恒返回 true，所以答案永远是 1）
func TestFirstBadVersion(t *testing.T) {
	for _, n := range []int{1, 2, 5, 10, 100} {
		var got int
		panicked := safeCall(func() { got = lc.FirstBadVersion(n) })
		if panicked {
			t.Errorf("FirstBadVersion(%d) panic", n)
			continue
		}
		// 当前 mock 是 isBadVersion(n) = true，所有版本都是坏的，第一个坏版本 = 1
		if got != 1 {
			t.Errorf("FirstBadVersion(%d) = %d, want 1（mock 全坏）", n, got)
		}
	}
}
