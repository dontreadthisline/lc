package algo_tests

import "demo/algo"

import (
	"math"
	"math/rand"
	"testing"
)

// initDist 预填 -1 的 dist/paths（多源 BFS 的哨兵约定）
func initDist(n int) (dist []int, paths []int) {
	dist = make([]int, n)
	paths = make([]int, n)
	for i := range dist {
		dist[i] = -1
		paths[i] = -1
	}
	return
}

// 手算用例：
//
//	0 → 1 → 2 ↘
//	              3    （3 从两侧可达：0 侧 3 跳，5 侧 2 跳 → min = 2）
//	5 → 4 ↗
//
// 源点 {0, 5}：dist = [0 1 2 2 1 0]
func TestMultiSourceBfsFixed(t *testing.T) {
	g := algo.NewGraph(6)
	g.AddEdge(0, 1, 1)
	g.AddEdge(1, 2, 1)
	g.AddEdge(2, 3, 1)
	g.AddEdge(5, 4, 1)
	g.AddEdge(4, 3, 1)

	dist, paths := initDist(g.N())
	algo.MultiSourceBfs(g, []int{0, 5}, dist, paths)

	want := []int{0, 1, 2, 2, 1, 0}
	for v, w := range want {
		if dist[v] != w {
			t.Errorf("dist[%d] = %d, want %d", v, dist[v], w)
		}
	}
	// 顶点 3 的最近源应是 5（2 跳），不是 0（3 跳）：回溯验证
	p := algo.PathTo(5, 3, paths)
	if len(p) != 3 || p[0] != 5 || p[2] != 3 {
		t.Errorf("3 的最近源路径 = %v, want [5 4 3]", p)
	}
}

// 交叉验证：多源一遍 vs 单源逐跑取 min（oracle 与被测代码零共享逻辑）
func TestMultiSourceBfsCrossCheck(t *testing.T) {
	r := rand.New(rand.NewSource(2024))
	for round := 0; round < 200; round++ {
		n := 2 + r.Intn(20)
		g := algo.NewGraph(n)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if r.Intn(4) == 0 {
					g.AddEdge(i, j, 1)
				}
			}
		}
		// 随机源集，故意混入重复和非法值
		k := 1 + r.Intn(n)
		sources := make([]int, k)
		for i := range sources {
			switch r.Intn(6) {
			case 0:
				sources[i] = -1
			case 1:
				sources[i] = n + r.Intn(3)
			default:
				sources[i] = r.Intn(n)
			}
		}

		// oracle：逐源单源 BFS，dist 取 pathLen，整体取 min
		want := make([]int, n)
		for i := range want {
			want[i] = math.MaxInt
		}
		for _, s := range sources {
			visited, paths := initTools(n)
			algo.Bfs(g, s, visited, paths)
			for v := 0; v < n; v++ {
				if visited[v] {
					if d := pathLen(paths, v); d < want[v] {
						want[v] = d
					}
				}
			}
		}

		dist, paths := initDist(n)
		algo.MultiSourceBfs(g, sources, dist, paths)

		for v := 0; v < n; v++ {
			switch {
			case want[v] == math.MaxInt && dist[v] != -1:
				t.Fatalf("round %d: %d 不可达但 dist=%d", round, v, dist[v])
			case want[v] != math.MaxInt && dist[v] != want[v]:
				t.Fatalf("round %d: dist[%d] = %d, want %d", round, v, dist[v], want[v])
			case dist[v] >= 0 && pathLen(paths, v) != dist[v]:
				t.Fatalf("round %d: paths 回溯长 %d ≠ dist[%d] = %d", round, pathLen(paths, v), v, dist[v])
			}
		}
	}
}

// 边界：空源集、全非法源、重复源
func TestMultiSourceBfsEdge(t *testing.T) {
	g := algo.NewGraph(3)
	g.AddEdge(0, 1, 1)

	dist, _ := initDist(3)
	algo.MultiSourceBfs(g, nil, dist, nil)
	for v, d := range dist {
		if d != -1 {
			t.Errorf("空源集 dist[%d] = %d, want -1", v, d)
		}
	}

	dist, _ = initDist(3)
	algo.MultiSourceBfs(g, []int{-1, 99}, dist, nil)
	for v, d := range dist {
		if d != -1 {
			t.Errorf("全非法源 dist[%d] = %d, want -1", v, d)
		}
	}

	dist, _ = initDist(3)
	algo.MultiSourceBfs(g, []int{0, 0, 0}, dist, nil)
	if dist[0] != 0 || dist[1] != 1 || dist[2] != -1 {
		t.Errorf("重复源去重错误: %v", dist)
	}
}
