package lc_tests

import (
	"sort"
	"testing"

	"demo/lc"
)

func anagramKey(s string) string {
	b := []byte(s)
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	return string(b)
}

func TestGroupAnagrams(t *testing.T) {
	cases := [][]string{
		{"eat", "tea", "tan", "ate", "nat", "bat"},
		{""},
		{"a"},
		{"ab", "ba", "ac"}, // 比较器传递性违例的对手盘
	}
	for _, strs := range cases {
		groups := lc.GroupAnagrams(strs)
		seen := 0
		for _, g := range groups {
			if len(g) == 0 {
				t.Errorf("出现空组: %v", groups)
				continue
			}
			key := anagramKey(g[0])
			for _, s := range g {
				if anagramKey(s) != key {
					t.Errorf("组内混入非异位词: %v", g)
				}
				seen++
			}
		}
		if seen != len(strs) {
			t.Errorf("分组覆盖不全: %d/%d", seen, len(strs))
		}
	}
}
