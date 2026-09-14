package algo

import "math"

// Edge 邻接表里的一条有向边。
type Edge struct {
	To     int
	Weight int
}

// Graph 是有向非负权图，邻接表存储。
type Graph struct {
	adj [][]Edge
}

// NewGraph 创建 n 个顶点（编号 0..n-1）的空图。
func NewGraph(n int) *Graph {
	return &Graph{adj: make([][]Edge, n)}
}

// AddEdge 加一条有向边 from→to；无向图正反各调一次。
// 权重必须非负：Dijkstra 的正确性依赖"弹出距离单调不减"，
// 负权边会破坏这一性质（有负权要用 Bellman-Ford）。
func (g *Graph) AddEdge(from, to, weight int) {
	g.adj[from] = append(g.adj[from], Edge{To: to, Weight: weight})
}

// visit 是堆里的条目：一条"到 node 的候选距离"。
type visit struct {
	node int
	dist int
}

// Dijkstra 单源最短路，返回 src 到各顶点的最短距离，不可达为 -1。
// 复杂度 O(E log E)，即 O(E log V)。
//
// 实现走"惰性删除"（lazy deletion）路线：
// 松弛成功时不修改堆中旧条目，而是压入更小的新条目；
// 旧条目留在堆里，弹出时因 dist 落后于 dist[] 被跳过。
// 代价是堆内最多同时存 O(E) 条目，换来完全不需要维护
// "元素 -> 堆内下标"的映射（eager decrease-key 需要它，
// 但 Heap 没有 swap 钩子也做不了，见 genericheap.go 的讨论）。
//
// 正确性要点：非负权保证堆顶弹出的距离单调不减，
// 因此每个节点第一次被弹出时，其距离已是最终值。
func Dijkstra(g *Graph, src int) []int {
	n := len(g.adj)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = math.MaxInt
	}
	dist[src] = 0

	h := NewHeap(func(a, b visit) bool { return a.dist < b.dist })
	h.Push(visit{node: src, dist: 0})

	for h.Len() > 0 {
		u := h.Pop()
		// 惰性删除核心：弹出的不是该节点最新的距离，说明是陈旧条目。
		if u.dist > dist[u.node] {
			continue
		}
		for _, e := range g.adj[u.node] {
			if nd := u.dist + e.Weight; nd < dist[e.To] {
				dist[e.To] = nd
				h.Push(visit{node: e.To, dist: nd}) // 旧条目不删，留给弹出时跳过
			}
		}
	}

	for i, d := range dist {
		if d == math.MaxInt {
			dist[i] = -1
		}
	}
	return dist
}

func DijkstraNaive(g *Graph, src int) []int {
	n := len(g.adj)
	dist := make([]int, n)
	done := make([]bool, n)
	for i := range dist {
		dist[i] = math.MaxInt
	}
	dist[src] = 0

	for range n {
		u := argmin(done, dist) //从src到u的未计算的最短路径,-1不存在
		if u == -1 || dist[u] == math.MaxInt {
			break
		}
		done[u] = true
		for _, e := range g.adj[u] {
			if nd := dist[u] + e.Weight; nd < dist[e.To] {
				dist[e.To] = nd
			}
		}
	}

	// 不可达归一化为 -1，与堆版约定一致
	for i, d := range dist {
		if d == math.MaxInt {
			dist[i] = -1
		}
	}
	return dist
}

func argmin(done []bool, dist []int) int {
	u := -1
	for v := range len(dist) {
		if !done[v] && (u == -1 || dist[v] < dist[u]) {
			u = v
		}
	}
	return u
}
