package lc_tests

import (
	"demo/lc"
	"testing"
)

func TestLongtestSubStr(t *testing.T) {
	strs := []string{"abcabcbb", "bbbbbb", "", "pwwkew"}
	wants := []int{3, 1, 0, 3}
	for i, str := range strs {
		got := lc.NoRepeateLongestSubStr(str)
		want := wants[i]
		if got != want {
			t.Fatal(got, want)
		}
	}
}
