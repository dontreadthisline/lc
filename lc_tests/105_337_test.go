package lc_tests

import (
	"math/rand"
	"testing"

	"demo/lc"
)

// ---- 树序列化辅助: 前序+中序 两元组唯一确定一棵树 ----

func treePreOrder(root *lc.TreeNode, out *[]int) {
	if root == nil {
		return
	}
	*out = append(*out, root.Val)
	treePreOrder(root.Left, out)
	treePreOrder(root.Right, out)
}

func treeInOrder(root *lc.TreeNode, out *[]int) {
	if root == nil {
		return
	}
	treeInOrder(root.Left, out)
	*out = append(*out, root.Val)
	treeInOrder(root.Right, out)
}

// ---- 105. 从前序与中序遍历序列构造二叉树 ----

func TestBuildTree(t *testing.T) {
	cases := []struct {
		pre, in []int
	}{
		{[]int{3, 9, 20, 15, 7}, []int{9, 3, 15, 20, 7}},
		{[]int{-1}, []int{-1}},
		{[]int{1, 2}, []int{2, 1}},
	}
	for _, c := range cases {
		var pre2, in2 []int
		runSafely(t, func() {
			root := lc.BuildTree(c.pre, c.in)
			treePreOrder(root, &pre2)
			treeInOrder(root, &in2)
		})
		if !equalInts(pre2, c.pre) || !equalInts(in2, c.in) {
			t.Errorf("BuildTree(%v, %v) 重建后 pre=%v in=%v", c.pre, c.in, pre2, in2)
		}
	}
}

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestBuildTreeRandom(t *testing.T) {
	for i := 0; i < 500; i++ {
		n := 1 + rand.Intn(40)
		// 互异值随机树 -> 序列化 -> 重建 -> 再序列化比对
		tree := genRandomTreeWithDistinct(n)
		var pre, in []int
		treePreOrder(tree, &pre)
		treeInOrder(tree, &in)
		var pre2, in2 []int
		runSafely(t, func() {
			root := lc.BuildTree(pre, in)
			treePreOrder(root, &pre2)
			treeInOrder(root, &in2)
		})
		if !equalInts(pre2, pre) || !equalInts(in2, in) {
			t.Fatalf("BuildTree 随机树重建失败 pre=%v in=%v -> pre2=%v in2=%v", pre, in, pre2, in2)
		}
	}
}

// genRandomTreeWithDistinct: 随机插入互异值构造随机形态二叉树(非 BST, 结构随机)
func genRandomTreeWithDistinct(n int) *lc.TreeNode {
	vals := randomShuffled(n)
	var insert func(node *lc.TreeNode, v int, goLeft bool) *lc.TreeNode
	insert = func(node *lc.TreeNode, v int, goLeft bool) *lc.TreeNode {
		if node == nil {
			return &lc.TreeNode{Val: v}
		}
		if goLeft {
			node.Left = insert(node.Left, v, rand.Intn(2) == 0)
		} else {
			node.Right = insert(node.Right, v, rand.Intn(2) == 0)
		}
		return node
	}
	var root *lc.TreeNode
	for _, v := range vals {
		if root == nil {
			root = &lc.TreeNode{Val: v}
		} else {
			insert(root, v, rand.Intn(2) == 0)
		}
	}
	return root
}

// ---- 337. 打家劫舍 III ----

// refRobIII: 暴力枚举 —— 树形独立集: 每个节点选/不选, 父子不能同选, 最大化权和
func refRobIII(root *lc.TreeNode) int {
	// 展开成数组 + 父指针
	type info struct {
		val    int
		parent int
	}
	nodes := []info{}
	var walk func(node *lc.TreeNode, parent int) int
	walk = func(node *lc.TreeNode, parent int) int {
		if node == nil {
			return -1
		}
		id := len(nodes)
		nodes = append(nodes, info{node.Val, parent})
		walk(node.Left, id)
		walk(node.Right, id)
		return id
	}
	walk(root, -1)
	n := len(nodes)
	best := 0
	for mask := 0; mask < 1<<n; mask++ {
		ok, sum := true, 0
		for i := 0; i < n && ok; i++ {
			if mask&(1<<i) != 0 {
				sum += nodes[i].val
				if nodes[i].parent >= 0 && mask&(1<<nodes[i].parent) != 0 {
					ok = false
				}
			}
		}
		if ok && sum > best {
			best = sum
		}
	}
	return best
}

func TestRobIII(t *testing.T) {
	cases := []struct {
		vals []*int
		want int
	}{
		{[]*int{ip(3), ip(2), ip(3), nil, ip(3), nil, ip(1)}, 7},
		{[]*int{ip(3), ip(4), ip(5), ip(1), ip(3), nil, ip(1)}, 9},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.RobIII(buildTree(c.vals)) })
		if got != c.want {
			t.Errorf("RobIII(%v) = %d, want %d", c.vals, got, c.want)
		}
	}
}

func TestRobIIIRandom(t *testing.T) {
	for i := 0; i < 300; i++ {
		n := 1 + rand.Intn(12)
		// 复用随机树, 覆盖值
		tree := genRandomTreeWithDistinct(n)
		var rewrite func(*lc.TreeNode)
		rewrite = func(node *lc.TreeNode) {
			if node == nil {
				return
			}
			node.Val = 1 + rand.Intn(20)
			rewrite(node.Left)
			rewrite(node.Right)
		}
		rewrite(tree)
		want := refRobIII(tree)
		var got int
		runSafely(t, func() { got = lc.RobIII(tree) })
		if got != want {
			t.Fatalf("RobIII 随机树 = %d, want %d(枚举)", got, want)
		}
	}
}
