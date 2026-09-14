package algo

import (
	"hash/fnv"
	"sort"
	"strconv"
)

// ConsistentHash 是一致性哈希环。
// 核心思路：每台机器算 replicas 个虚拟节点，虚拟节点的哈希值排进有序环。
// Get(key) 算 key 的哈希，在环上二分找第一个 >= 该哈希的虚拟节点，返回对应机器。
//
// 哈希函数用 fnv-1a + finalizer 搅拌，是确定性的：同一 key 在任意进程算出同一哈希值。
// 这是对分布式场景的关键要求——多个节点必须对 key 的归属达成一致，否则路由会乱。
//
// 纯 fnv64a 不够：它对 "machine-08#0..149" 这类结构化短串雪崩不足，
// 同一台机器的虚拟节点在环上扎堆，导致分布严重不均（某机器拿 21 万 key、另一台 4.6 万）。
// 在 fnv 输出后追加 murmur3 风格的 finalizer（xorshift-multiply）打散位依赖，
// 既保持确定性，又让虚拟节点均匀散开。
type ConsistentHash struct {
	replicas int               // 每台机器的虚拟节点数
	ring     []uint64          // 有序的虚拟节点哈希值，升序
	hashMap  map[uint64]string // 虚拟节点哈希 -> 机器名
}

// NewConsistentHash 创建空环。replicas 是每台机器的虚拟节点数。
func NewConsistentHash(replicas int) *ConsistentHash {
	if replicas <= 0 {
		replicas = 1
	}
	return &ConsistentHash{
		replicas: replicas,
		hashMap:  make(map[uint64]string),
	}
}

// hash 用 fnv-1a 算 64 位哈希，再经 finalizer 搅拌改善雪崩。
// 确定性：不含任何随机源，同 key 同结果，多节点一致。
func (c *ConsistentHash) hash(key string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(key))
	x := h.Sum64()
	// murmur3 风格 finalizer：打散 fnv 对结构化输入的位聚集。
	x ^= x >> 33
	x *= 0xff51afd7ed558ccd
	x ^= x >> 33
	return x
}

// Add 加入若干台机器，为每台生成 replicas 个虚拟节点插入环。
//
// 复杂度：O(n + k log k)，n 是环上已有虚拟节点数，k 是本次新增的虚拟节点数。
// 关键点：进入 Add 时 ring 已经有序，不需要整体重排。把新增的 k 个虚拟节点
// 单独排好（O(k log k)），再和有序的旧 ring 归并（O(n+k)）即可。
// 相比每次整体 sort.Slice（O(n log n)），省掉了对已有序部分的重复排序。
func (c *ConsistentHash) Add(machines ...string) {
	// 收集本次新增的虚拟节点：哈希值 + 它属于哪台机器。
	// 排序后顺序会打乱，所以必须用 map 保留 哈希->机器 的映射，归并时查表写 hashMap。
	type vnode struct {
		h uint64
		m string
	}
	newVnodes := make([]vnode, 0, len(machines)*c.replicas)
	for _, m := range machines {
		for i := 0; i < c.replicas; i++ {
			// 虚拟节点名：机器名 + 分隔符 + 序号
			name := m + "#" + strconv.Itoa(i)
			newVnodes = append(newVnodes, vnode{h: c.hash(name), m: m})
		}
	}
	// 新增段按哈希排序
	sort.Slice(newVnodes, func(i, j int) bool { return newVnodes[i].h < newVnodes[j].h })

	// 旧 ring(有序) 与 newVnodes(有序) 归并成新 ring。
	// 预分配 n+k 容量，避免归并过程中扩容拷贝。
	merged := make([]uint64, 0, len(c.ring)+len(newVnodes))
	i, j := 0, 0
	for i < len(c.ring) && j < len(newVnodes) {
		switch {
		case c.ring[i] < newVnodes[j].h:
			merged = append(merged, c.ring[i])
			i++
		case c.ring[i] > newVnodes[j].h:
			merged = append(merged, newVnodes[j].h)
			c.hashMap[newVnodes[j].h] = newVnodes[j].m
			j++
		default:
			// 哈希碰撞（极罕见）：保留旧节点，丢弃这个新虚拟节点，
			// 避免新机器名错误覆盖 hashMap 里已有的映射。
			i++
			j++
		}
	}
	// 收尾：把剩余段追加进 merged
	merged = append(merged, c.ring[i:]...)
	for ; j < len(newVnodes); j++ {
		merged = append(merged, newVnodes[j].h)
		c.hashMap[newVnodes[j].h] = newVnodes[j].m
	}
	c.ring = merged
}

// Remove 移除机器及其所有虚拟节点。
func (c *ConsistentHash) Remove(machines ...string) {
	// 先标记要删的机器
	drop := make(map[string]bool, len(machines))
	for _, m := range machines {
		drop[m] = true
	}
	// 从 ring 里剔除属于这些机器的虚拟节点
	kept := c.ring[:0]
	for _, h := range c.ring {
		if drop[c.hashMap[h]] {
			delete(c.hashMap, h)
			continue
		}
		kept = append(kept, h)
	}
	c.ring = kept
}

// Get 返回 key 应该落到哪台机器。环为空返回空串。
func (c *ConsistentHash) Get(key string) string {
	if len(c.ring) == 0 {
		return ""
	}
	h := c.hash(key)
	// 二分找第一个 >= h 的虚拟节点；若都小于 h，绕回首部（环的特性）
	idx := sort.Search(len(c.ring), func(i int) bool {
		return c.ring[i] >= h
	})
	if idx == len(c.ring) {
		idx = 0
	}
	return c.hashMap[c.ring[idx]]
}

// Machines 返回当前环上的机器集合（去重），用于外部观察分布。
func (c *ConsistentHash) Machines() []string {
	seen := make(map[string]bool)
	for _, m := range c.hashMap {
		seen[m] = true
	}
	out := make([]string, 0, len(seen))
	for m := range seen {
		out = append(out, m)
	}
	return out
}
