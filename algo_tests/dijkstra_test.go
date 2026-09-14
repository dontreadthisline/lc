package algo_tests

import "demo/algo"

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// 经典用例，手算期望值：
//
//	0 →1 (10)   0 →2 (3)   2 →1 (4)   1 →3 (1)
//	2 →3 (8)    2 →4 (2)   3 →4 (1)
//
// dist[1]=7  走 0→2→1 (3+4)，首条 0→1 (10) 成为陈旧条目
// dist[3]=8  走 0→2→1→3 (7+1)，首条 0→2→3 (11) 成为陈旧条目
// dist[4]=5  走 0→2→4 (3+2)
//
// 本用例恰好覆盖两次"陈旧条目跳过"路径。
func classicGraph() *algo.Graph {
	g := algo.NewGraph(5)
	g.AddEdge(0, 1, 10)
	g.AddEdge(0, 2, 3)
	g.AddEdge(2, 1, 4)
	g.AddEdge(1, 3, 1)
	g.AddEdge(2, 3, 8)
	g.AddEdge(2, 4, 2)
	g.AddEdge(3, 4, 1)
	return g
}

func TestDijkstraClassic(t *testing.T) {
	got := algo.Dijkstra(classicGraph(), 0)
	if want := []int{0, 7, 3, 8, 5}; !slices.Equal(got, want) {
		t.Fatalf("dist = %v, 期望 %v", got, want)
	}
}

func TestDijkstraNaiveClassic(t *testing.T) {
	got := algo.DijkstraNaive(classicGraph(), 0)
	if want := []int{0, 7, 3, 8, 5}; !slices.Equal(got, want) {
		t.Fatalf("dist = %v, 期望 %v", got, want)
	}
}

func TestDijkstraUnreachable(t *testing.T) {
	g := algo.NewGraph(4)
	g.AddEdge(0, 1, 5)
	// 顶点 2、3 无入边；同时覆盖朴素版的"剩余全不可达提前收工"分支

	for _, dij := range []func(*algo.Graph, int) []int{algo.Dijkstra, algo.DijkstraNaive} {
		got := dij(g, 0)
		if want := []int{0, 5, -1, -1}; !slices.Equal(got, want) {
			t.Fatalf("dist = %v, 期望 %v", got, want)
		}
	}
}

// 零权边合法（非负即可），最短路可以由零权边组成。
func TestDijkstraZeroWeight(t *testing.T) {
	g := algo.NewGraph(3)
	g.AddEdge(0, 1, 0)
	g.AddEdge(1, 2, 0)

	got := algo.Dijkstra(g, 0)
	if want := []int{0, 0, 0}; !slices.Equal(got, want) {
		t.Fatalf("dist = %v, 期望 %v", got, want)
	}
}

// 无向图 = 每条边加两次；换源点后距离数组应镜像对称。
func TestDijkstraUndirected(t *testing.T) {
	g := algo.NewGraph(4)
	for _, e := range [][3]int{{0, 1, 1}, {1, 2, 2}, {2, 3, 1}, {0, 3, 5}} {
		g.AddEdge(e[0], e[1], e[2])
		g.AddEdge(e[1], e[0], e[2])
	}

	if got := algo.Dijkstra(g, 0); !slices.Equal(got, []int{0, 1, 3, 4}) {
		t.Fatalf("从 0 出发 dist = %v", got) // 0→3 直达 5，绕行 1+2+1=4 更短
	}
	if got := algo.Dijkstra(g, 3); !slices.Equal(got, []int{4, 3, 1, 0}) {
		t.Fatalf("从 3 出发 dist = %v", got)
	}
}

// bellmanFord 是对照实现：O(V·E) 朴素松弛，完全不碰堆，
// 与被测代码无共享逻辑，交叉验证才有意义。
func bellmanFord(g *algo.Graph, src int) []int {
	n := g.N()
	dist := make([]int, n)
	for i := range dist {
		dist[i] = math.MaxInt
	}
	dist[src] = 0

	type wedge struct{ from, to, w int }
	var es []wedge
	for from := 0; from < g.N(); from++ {
		adj := g.Edges(from)
		for _, e := range adj {
			es = append(es, wedge{from, e.To, e.Weight})
		}
	}

	for range n - 1 {
		changed := false
		for _, e := range es {
			if dist[e.from] != math.MaxInt && dist[e.from]+e.w < dist[e.to] {
				dist[e.to] = dist[e.from] + e.w
				changed = true
			}
		}
		if !changed {
			break
		}
	}

	for i, d := range dist {
		if d == math.MaxInt {
			dist[i] = -1
		}
	}
	return dist
}

// 随机对照：随机图上 algo.Dijkstra 的结果必须与 Bellman-Ford 一致。
func TestDijkstraRandomCrossCheck(t *testing.T) {
	rnd := rand.New(rand.NewPCG(7, 42))
	for round := range 200 {
		n := 1 + rnd.IntN(12)
		g := algo.NewGraph(n)
		for from := range n {
			for to := range n {
				// 允许平行边、允许零权，只排除负权和自环
				if from != to && rnd.IntN(3) == 0 {
					g.AddEdge(from, to, rnd.IntN(11))
				}
			}
		}
		src := rnd.IntN(n)
		want := bellmanFord(g, src)

		// 两个实现跑同一批随机图，都必须与对照一致
		for _, dij := range []func(*algo.Graph, int) []int{algo.Dijkstra, algo.DijkstraNaive} {
			if got := dij(g, src); !slices.Equal(got, want) {
				t.Fatalf("round %d (n=%d, src=%d):\n got  %v\n want %v", round, n, src, got, want)
			}
		}
	}
}
