package algo_tests

import "demo/algo"

import (
	"reflect"
	"testing"
)

// 图结构：
//
//	0 → 1 → 3 ← 2    5 → 6      （5,6 与 0..4 不连通）
//	│↖───┘  │
//	└→ 2    └→ 4
var reachCases = []struct {
	start int
	want  map[int]bool
	dist  map[int]int // BFS 最短距离（DFS 不检查这列）
}{
	{0, map[int]bool{0: true, 1: true, 2: true, 3: true, 4: true},
		map[int]int{0: 0, 1: 1, 2: 1, 3: 2, 4: 2}},
	{5, map[int]bool{5: true, 6: true}, map[int]int{5: 0, 6: 1}},
	{6, map[int]bool{6: true}, map[int]int{6: 0}},
	{4, map[int]bool{4: true}, map[int]int{4: 0}}, // 有向图：4 没有出边
}

func buildGraph() *algo.Graph {
	g := algo.NewGraph(7)
	g.AddEdge(0, 1, 1)
	g.AddEdge(0, 2, 1)
	g.AddEdge(1, 3, 1)
	g.AddEdge(2, 3, 1)
	g.AddEdge(3, 0, 1) // 有环
	g.AddEdge(1, 4, 1)
	g.AddEdge(5, 6, 1) // 独立分量
	return g
}

// 初始化 visited 和 paths（-1 = 无前驱哨兵，不能用零值：0 是合法父顶点）
func initTools(n int) (map[int]bool, []int) {
	visited := map[int]bool{}
	paths := make([]int, n)
	for i := range paths {
		paths[i] = -1
	}
	return visited, paths
}

// 沿 paths 回溯得到路径长度
func pathLen(paths []int, v int) int {
	l := 0
	for paths[v] != -1 {
		v = paths[v]
		l++
	}
	return l
}

func TestDfs(t *testing.T) {
	g := buildGraph()
	for _, c := range reachCases {
		visited, _ := initTools(g.N())
		algo.Dfs(g, c.start, visited, make([]int, g.N()))
		if !reflect.DeepEqual(visitedKeys(visited), c.want) && !equalSets(visited, c.want) {
			t.Errorf("Dfs Reach(%d) = %v, want %v", c.start, visited, c.want)
		}
	}
}

func TestDfsIterWay(t *testing.T) {
	g := buildGraph()
	for _, c := range reachCases {
		visited, _ := initTools(g.N())
		algo.Dfs(g, c.start, visited, make([]int, g.N()))
		if !equalSets(visited, c.want) {
			t.Errorf("DfsIterWay Reach(%d) = %v, want %v", c.start, visited, c.want)
		}
	}
}

// BFS 的契约：可达 + paths 是最短路径树（每列都要验）
func TestBfs(t *testing.T) {
	g := buildGraph()
	for _, c := range reachCases {
		visited, paths := initTools(g.N())
		algo.Bfs(g, c.start, visited, paths)
		if !equalSets(visited, c.want) {
			t.Errorf("Bfs Reach(%d) = %v, want %v", c.start, visited, c.want)
		}
		for v, d := range c.dist {
			if got := pathLen(paths, v); got != d {
				t.Errorf("Bfs start=%d: 路径长(%d) = %d, want %d（paths=%v）",
					c.start, v, got, d, paths)
			}
		}
	}
}

// 专治“忘记预标记 start”：0→1, 1→0 的环回指根
func TestNoPreMarkDemo(t *testing.T) {
	g := algo.NewGraph(2)
	g.AddEdge(0, 1, 1)
	g.AddEdge(1, 0, 1)

	// 正确版
	visited, paths := initTools(2)
	algo.Bfs(g, 0, visited, paths)
	t.Logf("正确版:   paths=%v（paths[0]=-1，根无父亲）", paths)

	// 去掉预标记的版本
	badVisited, badPaths := initTools(2)
	badPaths[0] = -1
	expansions := map[int]int{}
	queue := []int{0} // ← 唯一区别：没标就入队
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		expansions[node]++
		for _, adj := range g.Edges(node) {
			if adj.To >= 0 && adj.To < g.N() && !badVisited[adj.To] {
				badVisited[adj.To] = true
				badPaths[adj.To] = node
				queue = append(queue, adj.To)
			}
		}
	}
	t.Logf("无预标记: paths=%v（paths[0]=%d，根被写入！）, 0 被展开 %d 次",
		badPaths, badPaths[0], expansions[0])
}

// d(4)=2（0→2→4），lazy 版会把 paths[4] 覆盖成 3（路径长 3）
func TestBfsPathsOverwrite(t *testing.T) {
	g := algo.NewGraph(5)
	g.AddEdge(0, 1, 1)
	g.AddEdge(0, 2, 1)
	g.AddEdge(1, 3, 1)
	g.AddEdge(2, 4, 1)
	g.AddEdge(3, 4, 1)

	visited, paths := initTools(5)
	algo.Bfs(g, 0, visited, paths)
	if got := pathLen(paths, 4); got != 2 {
		t.Errorf("paths[4] 重建路径长 = %d, want 2（paths=%v）", got, paths)
	}
}

// 非法输入：不 panic、不污染
func TestInvalidStart(t *testing.T) {
	g := buildGraph()
	for _, s := range []int{-1, 7} {
		visited, paths := initTools(g.N())
		algo.Dfs(g, s, visited, paths)
		algo.Dfs(g, s, visited, paths)
		algo.Bfs(g, s, visited, paths)
		if len(visited) != 0 {
			t.Errorf("start=%d 污染了 visited: %v", s, visited)
		}
	}
}

func equalSets(got, want map[int]bool) bool {
	if len(got) != len(want) {
		return false
	}
	for k := range want {
		if !got[k] {
			return false
		}
	}
	return true
}

func hasEdge(g *algo.Graph, from, to int) bool {
	for _, e := range g.Edges(from) {
		if e.To == to {
			return true
		}
	}
	return false
}

func TestPathTo(t *testing.T) {
	g := buildGraph()

	// BFS：0→4 的最短路 [0 1 4]
	_, paths := initTools(g.N())
	visited, _ := initTools(g.N())
	algo.Bfs(g, 0, visited, paths)
	if got := algo.PathTo(0, 4, paths); !reflect.DeepEqual(got, []int{0, 1, 4}) {
		t.Errorf("Bfs algo.PathTo(0,4) = %v, want [0 1 4]", got)
	}
	// 不可达 → nil（链断在 -1，没碰到 start）
	if got := algo.PathTo(0, 6, paths); got != nil {
		t.Errorf("algo.PathTo(0,6) = %v, want nil", got)
	}
	// start == end → [start]
	if got := algo.PathTo(0, 0, paths); !reflect.DeepEqual(got, []int{0}) {
		t.Errorf("algo.PathTo(0,0) = %v, want [0]", got)
	}
	// 起点不匹配：这份 paths 的根是 0，不是 5
	if got := algo.PathTo(5, 4, paths); got != nil {
		t.Errorf("algo.PathTo(5,4) = %v, want nil（paths 的根不是 5）", got)
	}
	// 非法 end
	if got := algo.PathTo(0, -1, paths); got != nil {
		t.Errorf("algo.PathTo(0,-1) = %v, want nil", got)
	}

	// DFS：路径合法但未必最短——首尾正确、相邻顶点之间都有边即可
	visited2, paths2 := initTools(g.N())
	algo.Dfs(g, 0, visited2, paths2)
	p := algo.PathTo(0, 4, paths2)
	if len(p) < 2 || p[0] != 0 || p[len(p)-1] != 4 {
		t.Fatalf("Dfs algo.PathTo(0,4) = %v, want 至少 [0 ... 4]", p)
	}
	for i := 0; i+1 < len(p); i++ {
		if !hasEdge(g, p[i], p[i+1]) {
			t.Errorf("边 %d→%d 不存在, path=%v", p[i], p[i+1], p)
		}
	}
}

func visitedKeys(m map[int]bool) map[int]bool { return m }

//2,1,1
//0,1,1
//0,1,2
