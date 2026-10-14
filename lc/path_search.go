package lc

import "slices"

// 一般有向图(可能有环) + 简单路径约束(不允许重复节点)的判定/枚举对比。
// 对照: lc/dag_path.go(DAG 版, 零 visited) —— 去掉无环保证后, 环把两件事逼出来了:
//   判定版需要 visited 防无限递归; 枚举版的简单性从「图结构免费」变成「解的定义」。

// HasPathSimple 判定版: 存在 start→end 路径即返回。
//
// visited 永久, 不撤销。免撤销的严格理由是单调性:
// visited 只增不减 → v 首次被探索时世界最宽松(封锁最少), 首败即终局 ——
// 之后任何时刻重试 v, 世界只会更严, 重试不可能比首试更好, 缓存正确。
// (通俗版: 从 v 到不了 end 是 v 的固有属性, 与来路无关)
// 另外两重职责: 防环(不标记会无限递归) + 失败缓存。
//
// 为什么判定版可以无视「简单路径」约束? 因为 存在路径 ⟺ 存在简单路径:
// 任何带环路径裁掉环段后仍是合法路径(裁环桥)。约束在存在性问题里是免费的。
// 整体 O(V+E), 与 DAG 记忆化版同阶。
/*
*
./dag_path.go
*/
func HasPathSimple(adj map[int][]int, start, end int) bool {
	visited := map[int]bool{}
	var found bool
	var dfs func(int)
	dfs = func(v int) {
		if found {
			return // 已找到, 剪掉剩余探索(短路)
		}
		if v == end {
			found = true
			return
		}
		if visited[v] {
			return // 防环 + 失败缓存, 一行两职
		}
		visited[v] = true
		for _, next := range adj[v] {
			dfs(next)
		}
	}
	dfs(start)
	return found
}

// AllPathsSimple 枚举版: 所有 start→end 简单路径。
//
// visited 随栈借还 —— 与判定版方向相反, 原因在时间结构:
// 判定版的世界单调收紧(v 首败即终局); 枚举版的世界随栈呼吸 ——
// 路径 A 在 v 失败后退栈归还, 路径 B 到达 v 时封锁已解除, 面对的世界可能
// 比 A 首试时更宽松, 首败可以翻案 —— 所以失败不能缓存, v 必须归还。
// 此时 visited 的身份是「当前路径的占用表」(路径属性), 不再是「全局死亡名单」(固有属性)。
// 不撤销则: 第一条 DFS 足迹封锁词典, 后续路径饿死(漏解)。
func AllPathsSimple(adj map[int][]int, start, end int) [][]int {
	visited := map[int]bool{}
	var res [][]int
	var dfs func(cur int, path []int)
	dfs = func(cur int, path []int) {
		if cur == end {
			res = append(res, slices.Clone(path)) // 收获必须拷贝
			return
		}
		visited[cur] = true // 占用: 路径属性, 不是节点固有属性
		for _, next := range adj[cur] {
			if !visited[next] {
				dfs(next, append(path, next)) // append 传参: path 无需回退
			}
		}
		visited[cur] = false // 撤销: 归还给下一条候选路径
	}
	dfs(start, []int{start})
	return res
}

// AllPathsSimplePruned 枚举版 + 预处理剪枝。
//
// 简单性虽是路径属性(不能随栈外缓存), 但它的弱化版是固有属性:
// 「从 v 完全不可达 end」(不带简单性约束的反向可达) 与来路无关, 可全局预处理。
// 反向 BFS 一次, 不可达 end 的节点整棵子树剪掉 —— 对应 127 讨论中
// 「马尔可夫可达性缓存仍可保留」的具体化: 只缓存绝对死, 不误伤换路可活的。
func AllPathsSimplePruned(adj map[int][]int, start, end int) [][]int {
	// 反向可达: 从 end 沿反向边能到的所有节点
	radj := map[int][]int{}
	for v, nexts := range adj {
		for _, to := range nexts {
			radj[to] = append(radj[to], v)
		}
	}
	canReachEnd := map[int]bool{end: true}
	queue := []int{end}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, prev := range radj[v] {
			if !canReachEnd[prev] {
				canReachEnd[prev] = true
				queue = append(queue, prev)
			}
		}
	}
	if !canReachEnd[start] {
		return nil
	}
	visited := map[int]bool{}
	var res [][]int
	var dfs func(cur int, path []int)
	dfs = func(cur int, path []int) {
		if cur == end {
			res = append(res, slices.Clone(path))
			return
		}
		visited[cur] = true
		for _, next := range adj[cur] {
			// 随栈占用(路径属性) + 永久可达(固有属性): 两种约束各司其职
			if !visited[next] && canReachEnd[next] {
				dfs(next, append(path, next))
			}
		}
		visited[cur] = false
	}
	dfs(start, []int{start})
	return res
}
