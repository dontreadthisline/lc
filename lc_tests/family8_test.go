package lc_tests

import (
	"demo/lc"
	"testing"
)

// 994. 腐烂的橘子（多源 BFS）
func TestOrangesRotting(t *testing.T) {
	cases := []struct {
		name string
		grid [][]int
		want int
	}{
		{"官方1_两分钟", [][]int{{2, 1, 1}, {1, 1, 0}, {0, 1, 1}}, 4},
		{"官方2_不可达", [][]int{{2, 1, 1}, {0, 1, 1}, {1, 0, 1}}, -1},
		{"官方3_全空", [][]int{{0, 2}}, 0},
		{"无新鲜", [][]int{{2, 2}, {2, 2}}, 0},
		{"单格新鲜", [][]int{{1}}, -1},
		{"单格腐烂", [][]int{{2}}, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got int
			panicked := safeCall(func() { got = lc.OrangesRotting(c.grid) })
			if panicked {
				t.Errorf("OrangesRotting panic")
				return
			}
			if got != c.want {
				t.Errorf("OrangesRotting(%v) = %d, want %d", c.grid, got, c.want)
			}
		})
	}
}

// 127. 单词接龙
func TestLadderLength(t *testing.T) {
	cases := []struct {
		name     string
		begin    string
		end      string
		wordList []string
		want     int
	}{
		{"官方1", "hit", "cog", []string{"hot", "dot", "dog", "lot", "log", "cog"}, 5},
		{"官方2_end不在列表", "hit", "cog", []string{"hot", "dot", "dog", "lot", "log"}, 0},
		{"begin等于end", "a", "a", []string{"a", "b", "c"}, 1},
		{"一步到位", "a", "b", []string{"a", "b"}, 2},
		{"begin不在列表但有路径", "hit", "hot", []string{"hot", "dot"}, 2},
		{"无路径", "aaa", "zzz", []string{"aaa", "bbb", "ccc", "zzz"}, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got int
			panicked := safeCall(func() { got = lc.LadderLength(c.begin, c.end, c.wordList) })
			if panicked {
				t.Errorf("LadderLength panic")
				return
			}
			if got != c.want {
				t.Errorf("LadderLength(%q, %q, %v) = %d, want %d", c.begin, c.end, c.wordList, got, c.want)
			}
		})
	}
}

// 695. 岛屿的最大面积（已修过，复验）
func TestMaxAreaOfIslandV2(t *testing.T) {
	if got := lc.MaxAreaOfIsland([][]int{{1}}); got != 1 {
		t.Errorf("单格岛 = %d, want 1", got)
	}
	if got := lc.MaxAreaOfIsland([][]int{{1, 1}}); got != 2 {
		t.Errorf("双格岛 = %d, want 2", got)
	}
	if got := lc.MaxAreaOfIsland([][]int{{0, 0}, {0, 0}}); got != 0 {
		t.Errorf("全零 = %d, want 0", got)
	}
}
