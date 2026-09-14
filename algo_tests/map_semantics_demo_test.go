package algo_tests

import "demo/algo"

import "testing"

// 现象2的道具：函数内部试图 make 补救
func rescueFix(visited map[int]bool) {
	visited = make(map[int]bool) // 只改了本地那份"指针拷贝"的指向
	visited[42] = true
}

func TestMapSemantics(t *testing.T) {
	// ── 现象1：nil map，读不炸、写炸 ──
	var nilMap map[int]bool
	t.Log("读 nil map 返回零值:", nilMap[1]) // false，安全
	func() {
		defer func() { t.Log("写 nil map 的下场:", recover()) }()
		nilMap[1] = true
	}()

	// ── 现象2：callee 里 make 补救，caller 拿不到 ──
	var v map[int]bool // 零值 = nil
	rescueFix(v)
	t.Log("调完 rescueFix 后 caller 的 v 仍是 nil:", v == nil)

	// ── 现象3：同一个 map 用两次，串味 ──
	g := algo.NewGraph(7)
	g.AddEdge(0, 1, 1)
	g.AddEdge(0, 2, 1)
	g.AddEdge(5, 6, 1)
	visited := map[int]bool{}
	algo.Dfs(g, 0, visited, make([]int, 7))
	t.Log("第一次 Reach(0):", len(visited), "个顶点") // 3
	algo.Dfs(g, 5, visited, make([]int, 7))
	t.Log("第二次 Reach(5):", len(visited), "个顶点 ← 期望 2，实际 5") // 残留！
}
