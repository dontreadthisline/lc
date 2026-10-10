package lc_tests

import (
	"math/rand"
	"slices"
	"testing"

	"demo/lc"
)

// ---- 450. 删除二叉搜索树中的节点 ----

// TestSearchNode: SearchNode 有实现, 应当为绿
func TestSearchNode(t *testing.T) {
	for i := 0; i < 200; i++ {
		n := 1 + rand.Intn(50)
		vals := randomShuffled(n)
		var root *lc.TreeNode
		for _, v := range vals {
			root = bstNaiveInsert(root, v)
		}
		// 每个存在的 key 都能找到且值匹配
		for _, v := range vals {
			node := lc.SearchNode(root, v)
			if node == nil || node.Val != v {
				t.Fatalf("SearchNode(%v, %d) 未找到", vals, v)
			}
		}
		// 不存在的 key 返回 nil
		for _, miss := range []int{0, -1, 1000, 9999} {
			if slices.Contains(vals, miss) {
				continue
			}
			if node := lc.SearchNode(root, miss); node != nil {
				t.Fatalf("SearchNode(%v, %d) = %d, want nil", vals, miss, node.Val)
			}
		}
	}
}

// TestDeleteNode: 删除后中序遍历应等于原中序去掉 key, 且仍有序(BST 合法)
func TestDeleteNode(t *testing.T) {
	for i := 0; i < 300; i++ {
		n := 1 + rand.Intn(40)
		vals := randomShuffled(n)
		var root *lc.TreeNode
		for _, v := range vals {
			root = bstNaiveInsert(root, v)
		}
		before := inorderVals(root)

		// 删除一个存在的 key
		key := vals[rand.Intn(n)]
		want := slices.DeleteFunc(slices.Clone(before), func(v int) bool { return v == key })
		got := inorderVals(lc.DeleteNode(root, key))
		if !slices.Equal(got, want) {
			t.Fatalf("DeleteNode(key=%d)\n before=%v\n got   =%v\n want  =%v", key, before, got, want)
		}
		if !slices.IsSorted(got) {
			t.Fatalf("删除后中序无序: %v", got)
		}

		// 删除不存在的 key: 树应不变
		for _, miss := range []int{0, -1, 1000} {
			if slices.Contains(before, miss) {
				continue
			}
			if got := inorderVals(lc.DeleteNode(root, miss)); !slices.Equal(got, before) {
				t.Fatalf("DeleteNode(不存在 key=%d) 改变了树: %v -> %v", miss, before, got)
			}
		}
	}
}
