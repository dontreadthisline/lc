package algo

// MultiSourceBfs 多源 BFS：k 个源点同时入队（dist=0），一次遍历得到每个
// 顶点到"最近源点"的跳数——等价于从虚拟超级源点跑单源 BFS（虚拟点向所有
// 源点各连一条边），层序单调对所有源统一成立，"入队即敲定"原样继承。
//
// dist 兼任 visited：-1 表示未发现，首次赋值即终值（同 dist 语义的哨兵复用，
// 省掉单独的 map）。paths 为最近源点方向的前驱树：从 v 回溯到根即它最近的
// 源点，pathTo 可直接复用。
//
// 契约：dist、paths 长度均为 len(g.adj) 且预填 -1（paths 可为 nil 省略）。
// 非法源点跳过；重复源点天然去重（dist 已非 -1）。
//
// 复杂度 O(V+E)，与源点个数无关——逐源跑 k 遍是 O(k(V+E)) 还要事后取 min。
func MultiSourceBfs(g *Graph, sources []int, dist []int, paths []int) {
	queue := NewQueue[int]()
	for _, s := range sources {
		if s < 0 || s >= len(g.adj) || dist[s] != -1 {
			continue
		}
		dist[s] = 0
		queue.Push(s)
	}
	for queue.Len() > 0 {
		u := queue.Pop()
		for _, e := range g.adj[u] {
			if e.To >= 0 && e.To < len(g.adj) && dist[e.To] == -1 {
				dist[e.To] = dist[u] + 1
				if paths != nil {
					paths[e.To] = u
				}
				queue.Push(e.To)
			}
		}
	}
}
