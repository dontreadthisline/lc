package lc_tests

import (
	"math/rand"
	"testing"

	"demo/lc"
)

// ---- DAG 路径搜索对比的验证 ----

// refAllPathsDAG: 独立参照 —— 递归拼接式(与被测的 path 传参式刻意不同)
// allPaths(v) = { [v] + p : p ∈ allPaths(next) for next ∈ adj[v] }, allPaths(end) = { [end] }
func refAllPathsDAG(adj map[int][]int, start, end int) [][]int {
	var build func(v int) [][]int
	build = func(v int) [][]int {
		if v == end {
			return [][]int{{v}}
		}
		var out [][]int
		for _, next := range adj[v] {
			for _, p := range build(next) {
				out = append(out, append([]int{v}, p...))
			}
		}
		return out
	}
	return build(start)
}

// randomDAG: 节点 0..n-1, 边 i→j 仅当 i<j(天然无环), 概率 p
func randomDAG(n int, p float64) map[int][]int {
	adj := map[int][]int{}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rand.Float64() < p {
				adj[i] = append(adj[i], j)
			}
		}
	}
	return adj
}

func TestDAGPathSearch(t *testing.T) {
	// 菱形 + 汇聚: 4 是共享中转, 2→3 使路径可以 1→2→3→4→5
	adj := map[int][]int{
		1: {2, 3},
		2: {3, 4},
		3: {4},
		4: {5},
	}
	if !lc.HasPathDAG(adj, 1, 5) {
		t.Errorf("HasPathDAG(1→5) = false, want true")
	}
	if !lc.HasPathDAGMemo(adj, 1, 5) {
		t.Errorf("HasPathDAGMemo(1→5) = false, want true")
	}
	if lc.HasPathDAG(adj, 5, 1) {
		t.Errorf("HasPathDAG(5→1) = true, want false(有向图, 逆方向无路)")
	}
	if got := lc.AllPathsDAG(adj, 1, 5); len(got) != 3 {
		t.Errorf("AllPathsDAG(1→5) = %v, want 3 条", got)
	}
	// 自环点: start == end
	if !lc.HasPathDAG(adj, 3, 3) {
		t.Errorf("HasPathDAG(3→3) = false, want true")
	}
	if got := lc.AllPathsDAG(adj, 3, 3); len(got) != 1 || got[0][0] != 3 {
		t.Errorf("AllPathsDAG(3→3) = %v, want [[3]]", got)
	}
}

func TestDAGPathSearchRandom(t *testing.T) {
	for i := 0; i < 1500; i++ {
		n := 2 + rand.Intn(9)
		adj := randomDAG(n, 0.35)
		start, end := rand.Intn(n), rand.Intn(n)

		want := refAllPathsDAG(adj, start, end)
		var gotPaths [][]int
		runSafely(t, func() { gotPaths = lc.AllPathsDAG(adj, start, end) })
		if setKey(gotPaths) != setKey(want) {
			t.Fatalf("AllPathsDAG(%d→%d, %v)\n got  %s\n want %s", start, end, adj, setKey(gotPaths), setKey(want))
		}
		// 判定版与枚举版一致性: 有路径 ⟺ 枚举非空
		var has, hasMemo bool
		runSafely(t, func() {
			has = lc.HasPathDAG(adj, start, end)
			hasMemo = lc.HasPathDAGMemo(adj, start, end)
		})
		if has != (len(want) > 0) || hasMemo != (len(want) > 0) {
			t.Fatalf("判定版不一致: naive=%v memo=%v, 枚举=%d 条 (start=%d end=%d adj=%v)",
				has, hasMemo, len(want), start, end, adj)
		}
	}
}
