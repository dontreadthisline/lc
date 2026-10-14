package lc_tests

import (
	"math/rand"
	"testing"

	"demo/lc"
)

// ---- 一般有向图 + 简单路径约束: 判定/枚举对比的验证 ----

// refAllPathsSimple: 独立参照 —— 每层复制 visited 的纯函数式(与被测的撤销式刻意不同)
func refAllPathsSimple(adj map[int][]int, start, end int) [][]int {
	var res [][]int
	var walk func(cur int, path []int, visited map[int]bool)
	walk = func(cur int, path []int, visited map[int]bool) {
		if cur == end {
			res = append(res, append([]int(nil), path...))
			return
		}
		for _, next := range adj[cur] {
			if !visited[next] {
				nv := make(map[int]bool, len(visited)+1)
				for k := range visited {
					nv[k] = true
				}
				nv[next] = true
				walk(next, append(path, next), nv)
			}
		}
	}
	if start == end {
		return [][]int{{start}}
	}
	walk(start, []int{start}, map[int]bool{start: true})
	return res
}

// refReachable: BFS 可达性(判定版参照)
func refReachable(adj map[int][]int, start, end int) bool {
	seen := map[int]bool{start: true}
	queue := []int{start}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		if v == end {
			return true
		}
		for _, next := range adj[v] {
			if !seen[next] {
				seen[next] = true
				queue = append(queue, next)
			}
		}
	}
	return false
}

// randomDigraph: 任意方向的随机有向边(可能成环)
func randomDigraph(n int, p float64) map[int][]int {
	adj := map[int][]int{}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j && rand.Float64() < p {
				adj[i] = append(adj[i], j)
			}
		}
	}
	return adj
}

func TestPathSearchSimple(t *testing.T) {
	// 环图: 1→2→3→1 成环, 3→4 出口
	// 判定不死循环; 枚举恰一条简单路径 1→2→3→4
	adj := map[int][]int{
		1: {2},
		2: {3},
		3: {1, 4},
		4: {},
	}
	if !lc.HasPathSimple(adj, 1, 4) {
		t.Errorf("HasPathSimple(1→4) = false, want true")
	}
	got := lc.AllPathsSimple(adj, 1, 4)
	if len(got) != 1 || setKey(got) != setKey([][]int{{1, 2, 3, 4}}) {
		t.Errorf("AllPathsSimple(1→4) = %v, want [[1 2 3 4]]", got)
	}
	if lc.HasPathSimple(adj, 4, 1) {
		t.Errorf("HasPathSimple(4→1) = true, want false(4 无出边)")
	}
	// 共享中转 + 环: 1→{2,3}, 2→4, 3→4, 4→2(回边), 2→5, 4→5
	// 简单路径: 1-2-4-5, 1-2-5, 1-3-4-5, 1-3-4-2-5 (4→2 回边在 2 未用时可用!)
	adj2 := map[int][]int{
		1: {2, 3},
		2: {4, 5},
		3: {4},
		4: {2, 5},
	}
	got2 := lc.AllPathsSimple(adj2, 1, 5)
	if len(got2) != 4 {
		t.Errorf("AllPathsSimple(1→5) = %v, want 4 条(回边 4→2 在 2 未用时可用)", got2)
	}
	// 剪枝版与朴素版一致
	if setKey(lc.AllPathsSimplePruned(adj2, 1, 5)) != setKey(got2) {
		t.Errorf("AllPathsSimplePruned 与朴素版结果不一致")
	}
}

func TestPathSearchSimpleRandom(t *testing.T) {
	for i := 0; i < 1500; i++ {
		n := 2 + rand.Intn(6)
		adj := randomDigraph(n, 0.3)
		start, end := rand.Intn(n), rand.Intn(n)

		want := refAllPathsSimple(adj, start, end)
		var got, gotPruned [][]int
		var has bool
		runSafely(t, func() {
			got = lc.AllPathsSimple(adj, start, end)
			gotPruned = lc.AllPathsSimplePruned(adj, start, end)
			has = lc.HasPathSimple(adj, start, end)
		})
		if setKey(got) != setKey(want) {
			t.Fatalf("AllPathsSimple(%d→%d, %v)\n got  %s\n want %s", start, end, adj, setKey(got), setKey(want))
		}
		if setKey(gotPruned) != setKey(want) {
			t.Fatalf("AllPathsSimplePruned(%d→%d, %v) 与参照不一致", start, end, adj)
		}
		// 判定 ⟺ 枚举非空 (存在路径 ⟺ 存在简单路径)
		if has != refReachable(adj, start, end) {
			t.Fatalf("HasPathSimple 与 BFS 可达性不一致: start=%d end=%d adj=%v", start, end, adj)
		}
	}
}
