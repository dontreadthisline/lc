package lc_tests

import (
	"demo/lc"
	"testing"
)

// 098. 验证二叉搜索树（中序版 + 上下界版）
func TestIsValidBSTBoth(t *testing.T) {
	impls := []struct {
		name string
		fn   func(*lc.TreeNode) bool
	}{
		{"IsValidBST", lc.IsValidBST},
		{"IsValidBSTRec", lc.IsValidBSTRec},
	}
	for _, impl := range impls {
		cases := []struct {
			name string
			vals []*int
			want bool
		}{
			{"官方1", []*int{ip(2), ip(1), ip(3)}, true},
			{"官方2", []*int{ip(5), ip(1), ip(4), nil, nil, ip(3), ip(6)}, false},
			{"孙辈越界", []*int{ip(5), ip(4), ip(6), nil, nil, ip(3), ip(7)}, false},
			{"单节点", []*int{ip(1)}, true},
			{"重复", []*int{ip(2), ip(2), ip(2)}, false},
		}
		for _, c := range cases {
			t.Run(impl.name+"/"+c.name, func(t *testing.T) {
				var got bool
				panicked := safeCall(func() { got = impl.fn(buildTree(c.vals)) })
				if panicked {
					t.Errorf("%s panic", impl.name)
					return
				}
				if got != c.want {
					t.Errorf("%s(%v) = %v, want %v", impl.name, c.vals, got, c.want)
				}
			})
		}
		t.Run(impl.name+"/空树", func(t *testing.T) {
			if got := impl.fn(nil); !got {
				t.Errorf("%s(空树) = false, want true", impl.name)
			}
		})
	}
}

// 236. 二叉树的最近公共祖先（双返回值版）
func TestLowestCommonAncestor(t *testing.T) {
	// 构建官方用例树
	root := buildTree([]*int{
		ip(3),
		ip(5), ip(1),
		ip(6), ip(2), ip(0), ip(8),
		nil, nil, ip(7), ip(4),
	})
	// 收集节点指针
	var find func(*lc.TreeNode, int) *lc.TreeNode
	find = func(n *lc.TreeNode, val int) *lc.TreeNode {
		if n == nil {
			return nil
		}
		if n.Val == val {
			return n
		}
		if l := find(n.Left, val); l != nil {
			return l
		}
		return find(n.Right, val)
	}
	p, q := find(root, 5), find(root, 1)
	if p == nil || q == nil {
		t.Fatal("建树失败")
	}

	tests := []struct {
		name string
		fn   func(*lc.TreeNode, *lc.TreeNode, *lc.TreeNode) *lc.TreeNode
	}{
		{"双返回值版", lc.LowestCommonAncestor},
		{"压缩版", lc.LowestCommonAncestor},
	}
	_ = tests // 只测导出的双返回值版

	got := lc.LowestCommonAncestor(root, p, q)
	if got == nil || got.Val != 3 {
		t.Errorf("LCA(5,1) = %v, want 3", got)
	}

	p2, q2 := find(root, 5), find(root, 4)
	got2 := lc.LowestCommonAncestor(root, p2, q2)
	if got2 == nil || got2.Val != 5 {
		t.Errorf("LCA(5,4) = %v, want 5", got2)
	}
}

// 114. 二叉树展开为链表
func TestFlatten(t *testing.T) {
	t.Run("官方用例", func(t *testing.T) {
		root := buildTree([]*int{
			ip(1), ip(2), ip(5), ip(3), ip(4), nil, ip(6),
		})
		var panicked bool
		safeCall(func() { lc.Flatten(root) })
		if panicked {
			t.Fatal("Flatten panic")
			return
		}
		// 展开后应为 1→2→3→4→5→6 全右链
		want := []int{1, 2, 3, 4, 5, 6}
		cur := root
		for _, v := range want {
			if cur == nil || cur.Val != v {
				t.Errorf("展开后位置值 = %v, want %d", cur, v)
				return
			}
			if cur.Left != nil {
				t.Errorf("位置 %d 的左子树非 nil", v)
				return
			}
			cur = cur.Right
		}
		if cur != nil {
			t.Errorf("展开后还有多余节点")
		}
	})
	t.Run("空树", func(t *testing.T) {
		var panicked bool
		safeCall(func() { lc.Flatten(nil) })
		if panicked {
			t.Error("Flatten(nil) panic")
		}
	})
	t.Run("单节点", func(t *testing.T) {
		root := &lc.TreeNode{Val: 1}
		lc.Flatten(root)
		if root.Left != nil || root.Right != nil {
			t.Error("单节点展开后不应有子节点")
		}
	})
}

// 235. BST 的最近公共祖先
func TestLowestCommonAncestorBST(t *testing.T) {
	root := buildTree([]*int{
		ip(6), ip(2), ip(8), ip(0), ip(4), ip(7), ip(9), nil, nil, ip(3), ip(5),
	})
	var find func(*lc.TreeNode, int) *lc.TreeNode
	find = func(n *lc.TreeNode, val int) *lc.TreeNode {
		if n == nil {
			return nil
		}
		if n.Val == val {
			return n
		}
		if l := find(n.Left, val); l != nil {
			return l
		}
		return find(n.Right, val)
	}
	cases := []struct {
		name    string
		pv, qv  int
		wantVal int
	}{
		{"官方1", 2, 8, 6},
		{"官方2", 2, 4, 2},
		{"深嵌套", 3, 5, 4},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p, q := find(root, c.pv), find(root, c.qv)
			got := lc.LowestCommonAncestorBST(root, p, q)
			if got == nil || got.Val != c.wantVal {
				t.Errorf("LCA(%d,%d) = %v, want %d", c.pv, c.qv, got, c.wantVal)
			}
		})
	}
}

// 437. 路径总和 III
func TestPathSum(t *testing.T) {
	cases := []struct {
		name   string
		vals   []*int
		target int
		want   int
	}{
		{"官方1", []*int{
			ip(10), ip(5), ip(-3), ip(3), ip(2), nil, ip(11),
			ip(3), ip(-2), nil, ip(1),
		}, 8, 3},
		{"官方2", []*int{
			ip(5), ip(4), ip(8), ip(11), nil, ip(13), ip(4),
			ip(7), ip(2), nil, nil, ip(5), ip(1),
		}, 22, 3},
		{"空树", nil, 10, 0},
		{"单节点等于目标", []*int{ip(5)}, 5, 1},
		{"单节点不等于", []*int{ip(5)}, 3, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got int
			panicked := safeCall(func() { got = lc.PathSum(buildTree(c.vals), c.target) })
			if panicked {
				t.Errorf("PathSum panic")
				return
			}
			if got != c.want {
				t.Errorf("PathSum(%v, %d) = %d, want %d", c.vals, c.target, got, c.want)
			}
		})
	}
}
