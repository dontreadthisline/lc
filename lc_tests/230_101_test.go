package lc_tests

import (
	"math/rand"
	"slices"
	"testing"

	"demo/lc"
)

// ---- 测试辅助(与 family7_tree_test.go 的 buildTree/ip 不冲突) ----

// bstFromSorted: 有序数组递归建平衡 BST
func bstFromSorted(arr []int) *lc.TreeNode {
	if len(arr) == 0 {
		return nil
	}
	mid := len(arr) / 2
	return &lc.TreeNode{Val: arr[mid], Left: bstFromSorted(arr[:mid]), Right: bstFromSorted(arr[mid+1:])}
}

// bstNaiveInsert: 测试内自带的朴素插入, 用于构造随机形态 BST(不依赖被测的 InsertNode)
func bstNaiveInsert(root *lc.TreeNode, key int) *lc.TreeNode {
	if root == nil {
		return &lc.TreeNode{Val: key}
	}
	if key < root.Val {
		root.Left = bstNaiveInsert(root.Left, key)
	} else {
		root.Right = bstNaiveInsert(root.Right, key)
	}
	return root
}

func inorderVals(root *lc.TreeNode) []int {
	var res []int
	var walk func(*lc.TreeNode)
	walk = func(n *lc.TreeNode) {
		if n == nil {
			return
		}
		walk(n.Left)
		res = append(res, n.Val)
		walk(n.Right)
	}
	walk(root)
	return res
}

// randomShuffled: 返回 [1..n] 的随机排列
func randomShuffled(n int) []int {
	arr := make([]int, n)
	for i := range arr {
		arr[i] = i + 1
	}
	rand.Shuffle(n, func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	return arr
}

// isSymRef: 独立参照实现(递归镜像比较)
func isSymRef(p, q *lc.TreeNode) bool {
	if p == nil || q == nil {
		return p == q
	}
	return p.Val == q.Val && isSymRef(p.Left, q.Right) && isSymRef(p.Right, q.Left)
}

// genRandomTree: 随机形态随机值的树(几乎必不对称, 但以参照实现对拍为准)
func genRandomTree(depth int) *lc.TreeNode {
	if depth == 0 || rand.Intn(4) == 0 {
		return nil
	}
	return &lc.TreeNode{
		Val:   1 + rand.Intn(1000),
		Left:  genRandomTree(depth - 1),
		Right: genRandomTree(depth - 1),
	}
}

// genMirrorTree: 构造一棵必对称的树(随机左子树 + 其镜像)
func genMirrorTree(depth int) *lc.TreeNode {
	var mirror func(*lc.TreeNode) *lc.TreeNode
	mirror = func(n *lc.TreeNode) *lc.TreeNode {
		if n == nil {
			return nil
		}
		return &lc.TreeNode{Val: n.Val, Left: mirror(n.Right), Right: mirror(n.Left)}
	}
	root := 1 + rand.Intn(1000)
	left := genRandomTree(depth)
	return &lc.TreeNode{Val: root, Left: left, Right: mirror(left)}
}

// ---- 230. 二叉搜索树中第K小的元素 ----

func TestKthSmallest(t *testing.T) {
	// 题面示例: root = [3,1,4,null,2], k = 1 -> 1
	root1 := buildTree([]*int{ip(3), ip(1), ip(4), nil, ip(2)})
	if got := lc.KthSmallest(root1, 1); got != 1 {
		t.Errorf("KthSmallest([3,1,4,null,2], 1) = %d, want 1", got)
	}
	// 题面示例: root = [5,3,6,2,4,null,null,1], k = 3 -> 3
	root2 := buildTree([]*int{ip(5), ip(3), ip(6), ip(2), ip(4), nil, nil, ip(1)})
	if got := lc.KthSmallest(root2, 3); got != 3 {
		t.Errorf("KthSmallest([5,3,6,2,4,null,null,1], 3) = %d, want 3", got)
	}
}

func TestKthSmallestRandom(t *testing.T) {
	for i := 0; i < 300; i++ {
		n := 1 + rand.Intn(80)
		vals := randomShuffled(n)
		var root *lc.TreeNode
		if i%2 == 0 {
			sorted := slices.Clone(vals)
			slices.Sort(sorted)
			root = bstFromSorted(sorted) // 平衡树
		} else {
			for _, v := range vals {
				root = bstNaiveInsert(root, v) // 随机形态退化树
			}
		}
		in := inorderVals(root)
		if !slices.IsSorted(in) {
			t.Fatalf("构造的树不是 BST: %v", in)
		}
		for k := 1; k <= n; k++ {
			if got := lc.KthSmallest(root, k); got != in[k-1] {
				t.Fatalf("KthSmallest k=%d = %d, want %d (inorder=%v)", k, got, in[k-1], in)
			}
		}
	}
}

// ---- 101. 对称二叉树 ----

func TestIsSymmetric(t *testing.T) {
	cases := []struct {
		vals []*int
		want bool
	}{
		{[]*int{ip(1), ip(2), ip(2), ip(3), ip(4), ip(4), ip(3)}, true},
		{[]*int{ip(1), ip(2), ip(2), nil, ip(3), nil, ip(3)}, false},
		{[]*int{ip(1)}, true},
		{[]*int{ip(1), ip(2), ip(2), ip(3), nil, nil, ip(3)}, true},
		{[]*int{ip(1), ip(2), ip(3)}, false},
	}
	for _, c := range cases {
		root := buildTree(c.vals)
		if got := lc.IsSymmetric(root); got != c.want {
			t.Errorf("IsSymmetric(%v) = %v, want %v", c.vals, got, c.want)
		}
	}
}

func TestIsSymmetricRandom(t *testing.T) {
	for i := 0; i < 5000; i++ {
		// 随机树: 与参照实现对拍
		tree := genRandomTree(1 + rand.Intn(6))
		want := tree == nil || isSymRef(tree.Left, tree.Right)
		if got := lc.IsSymmetric(tree); got != want {
			t.Fatalf("IsSymmetric(随机树) = %v, want %v", got, want)
		}
		// 镜像构造树: 必须为 true
		mt := genMirrorTree(1 + rand.Intn(5))
		if got := lc.IsSymmetric(mt); !got {
			t.Fatalf("IsSymmetric(镜像构造树) = false, want true")
		}
	}
}
