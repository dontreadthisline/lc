package lc_tests

import (
	"demo/lc"
	"testing"
)

// 树测试辅助：层序构建（*int 为 nil）
func ip(i int) *int { return &i }

func buildTree(nodes []*int) *lc.TreeNode {
	if len(nodes) == 0 || nodes[0] == nil {
		return nil
	}
	root := &lc.TreeNode{Val: *nodes[0]}
	queue := []*lc.TreeNode{root}
	idx := 1
	for len(queue) > 0 && idx < len(nodes) {
		cur := queue[0]
		queue = queue[1:]
		if idx < len(nodes) {
			if nodes[idx] != nil {
				cur.Left = &lc.TreeNode{Val: *nodes[idx]}
				queue = append(queue, cur.Left)
			}
			idx++
		}
		if idx < len(nodes) {
			if nodes[idx] != nil {
				cur.Right = &lc.TreeNode{Val: *nodes[idx]}
				queue = append(queue, cur.Right)
			}
			idx++
		}
	}
	return root
}

// 098. 验证二叉搜索树
func TestIsValidBST(t *testing.T) {
	cases := []struct {
		name string
		vals []*int
		want bool
	}{
		{"官方用例1", []*int{ip(2), ip(1), ip(3)}, true},
		{"官方用例2_直接子违例", []*int{ip(5), ip(1), ip(4), nil, nil, ip(3), ip(6)}, false},
		{"孙辈越界_仅查父子会漏", []*int{ip(5), ip(4), ip(6), nil, nil, ip(3), ip(7)}, false},
		{"单节点", []*int{ip(1)}, true},
		{"重复值非法", []*int{ip(2), ip(2), ip(2)}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got bool
			panicked := safeCall(func() { got = lc.IsValidBST(buildTree(c.vals)) })
			if panicked {
				t.Errorf("IsValidBST panic")
				return
			}
			if got != c.want {
				t.Errorf("IsValidBST(%v) = %v, want %v", c.vals, got, c.want)
			}
		})
	}
	var got bool
	if safeCall(func() { got = lc.IsValidBST(nil) }); !got {
		t.Errorf("空树应为 true")
	}
}

// 102. 二叉树的层序遍历
func TestLevelOrder(t *testing.T) {
	a := lc.LevelOrder(buildTree([]*int{ip(3), ip(9), ip(20), nil, nil, ip(15), ip(7)}))
	want := [][]int{{3}, {9, 20}, {15, 7}}
	if len(a) != len(want) {
		t.Fatalf("LevelOrder 层数 = %d, want %d", len(a), len(want))
	}
	for i := range want {
		for j, v := range want[i] {
			if a[i][j] != v {
				t.Errorf("LevelOrder[%d][%d] = %d, want %d", i, j, a[i][j], v)
			}
		}
	}
	if got := lc.LevelOrder(nil); len(got) != 0 {
		t.Errorf("空树应返回空, got %v", got)
	}
	if got := lc.LevelOrder(&lc.TreeNode{Val: 1}); len(got) != 1 || got[0][0] != 1 {
		t.Errorf("单节点应 [[1]], got %v", got)
	}
}

// 104. 二叉树的最大深度
func TestMaxDepth(t *testing.T) {
	if got := lc.MaxDepth(buildTree([]*int{ip(3), ip(9), ip(20), nil, nil, ip(15), ip(7)})); got != 3 {
		t.Errorf("MaxDepth 官方树 = %d, want 3", got)
	}
	if got := lc.MaxDepth(nil); got != 0 {
		t.Errorf("空树 = %d, want 0", got)
	}
	skew := buildTree([]*int{ip(1), ip(2), nil, ip(3)})
	if got := lc.MaxDepth(skew); got != 3 {
		t.Errorf("左斜树 = %d, want 3", got)
	}
}

// 124. 二叉树中的最大路径和
func TestMaxPathSum(t *testing.T) {
	cases := []struct {
		name string
		vals []*int
		want int
	}{
		{"官方用例1", []*int{ip(-10), ip(9), ip(20), nil, nil, ip(15), ip(7)}, 42},
		{"官方用例2_单负节点", []*int{ip(-3)}, -3},
		{"官方用例3", []*int{ip(2), ip(-1)}, 2},
		{"全正小树", []*int{ip(1), ip(2), ip(3)}, 6},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got int
			panicked := safeCall(func() { got = lc.MaxPathSum(buildTree(c.vals)) })
			if panicked {
				t.Errorf("MaxPathSum(%v) panic", c.vals)
				return
			}
			if got != c.want {
				t.Errorf("MaxPathSum(%v) = %d, want %d", c.vals, got, c.want)
			}
		})
	}
}

// 226. 翻转二叉树
func TestInvertTree(t *testing.T) {
	root := buildTree([]*int{ip(4), ip(2), ip(7), ip(1), ip(3), ip(6), ip(9)})
	got := lc.InvertTree(root)
	a := lc.LevelOrder(got)
	want := [][]int{{4}, {7, 2}, {9, 6, 3, 1}}
	for i := range want {
		for j, v := range want[i] {
			if a[i][j] != v {
				t.Errorf("InvertTree[%d][%d] = %d, want %d", i, j, a[i][j], v)
			}
		}
	}
	if lc.InvertTree(nil) != nil {
		t.Errorf("空树应返回 nil")
	}
}
