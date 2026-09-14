package algo

func DfsRec(g *Graph, start int, visited map[int]bool, paths []int) {
	if start < 0 || start >= len(g.adj) {
		return
	}
	visited[start] = true
	for _, adj := range g.adj[start] {
		if !visited[adj.To] {
			paths[adj.To] = start
			DfsRec(g, adj.To, visited, paths)
		}
	}
}

func Dfs(g *Graph, start int, visited map[int]bool, paths []int) {
	if start < 0 || start >= len(g.adj) {
		return
	}
	visited[start] = true
	stack := []int{start}
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, adj := range g.adj[node] {
			if adj.To >= 0 && adj.To < len(g.adj) && !visited[adj.To] {
				visited[adj.To] = true
				paths[adj.To] = node
				stack = append(stack, adj.To)
			}
		}
	}
}

func Bfs(g *Graph, start int, visited map[int]bool, paths []int) {
	if start < 0 || start >= len(g.adj) {
		return
	}
	visited[start] = true
	queue := []int{start}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		for _, adj := range g.adj[node] {
			if adj.To >= 0 && adj.To < len(g.adj) && !visited[adj.To] {
				visited[adj.To] = true
				paths[adj.To] = node
				queue = append(queue, adj.To)
			}
		}
	}
}

func PathTo(start, end int, paths []int) []int {
	if end < 0 || end >= len(paths) {
		return nil
	}
	way := make([]int, 0, len(paths))
	for v := end; v != -1; v = paths[v] {
		way = append(way, v)
		if v == start {
			for i, j := 0, len(way)-1; i < j; i, j = i+1, j-1 {
				way[i], way[j] = way[j], way[i]
			}
			return way
		}
		if len(way) > len(paths) {
			return nil
		}
	}
	return nil
}

// N 返回顶点数（0..n-1）。
func (g *Graph) N() int { return len(g.adj) }

// Edges 返回 from 的全部出边（只读视图，调用方不得修改）。
func (g *Graph) Edges(from int) []Edge { return g.adj[from] }
