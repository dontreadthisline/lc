package lc

import "slices"

// DAG 路径搜索对比: 判定版 vs 枚举版
// 对照讨论: 127 单词接龙(无向有环词典图) 与 77/78(递增 DAG) 的理论
// 约定: adj 为有向邻接表, 结构上保证无环(DAG), 否则两个函数都需要补 visited

// HasPathDAG 判定版: 存在一条 start→end 即返回。
//
// 整个函数没有 visited —— DAG 无环, 永远走不回自己, 防环账本天然免了。
// 朴素版的代价: 同一节点可被不同前缀重复探索(菱形结构 1→2→4 / 1→3→4 里的 4 被搜两次),
// 最坏代价与「路径条数」同阶。
/*
./path_search.go
*/

func HasPathDAG(adj map[int][]int, start, end int) bool {
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
		for _, next := range adj[v] {
			dfs(next)
		}
	}
	dfs(start)
	return found
}

// HasPathDAGMemo 判定版 + 记忆化: 「从 v 到不了 end」是 v 的固有属性(与来路无关)。
//
// dead 是永久缓存, 不撤销 —— 与 LadderLength 判定版同款语义(失败缓存)。
// 复杂度从「路径条数级」降到 O(V+E)。

func HasPathDAGMemo(adj map[int][]int, start, end int) bool {
	dead := map[int]bool{} // 失败缓存: dead[v]=true 表示从 v 出发到不了 end
	var found bool
	var dfs func(int)
	dfs = func(v int) {
		if found {
			return
		}
		if v == end {
			found = true
			return
		}
		if dead[v] {
			return // 这个节点没前途, 之前探明过, 直接跳过
		}
		for _, next := range adj[v] {
			dfs(next)
		}
		if !found {
			dead[v] = true // 只有确认此路不通才缓存; found 时 v 明明有活路, 不能标死
		}
	}
	dfs(start)
	return found
}

// AllPathsDAG 枚举版: 所有 start→end 路径。
//
// 同样没有 visited、没有撤销 —— DAG 的方向蕴含了「路径天然简单」,
// 与 77/78 组合/子集的递增 DAG 完全同构: 路的修法本身就是剪枝。
// append 传参: 本帧 path 从未被赋值, 无需回退。
func AllPathsDAG(adj map[int][]int, start, end int) [][]int {
	var res [][]int
	var dfs func(cur int, path []int)
	dfs = func(cur int, path []int) {
		if cur == end {
			res = append(res, slices.Clone(path)) // 收获必须拷贝
			return
		}
		for _, next := range adj[cur] {
			dfs(next, append(path, next)) // 零共享可变状态, 零撤销
		}
	}
	dfs(start, []int{start})
	return res
}
