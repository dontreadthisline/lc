package lc

// 235. 二叉搜索树的最近公共祖先
// 题面: ../mother-problems/0235-二叉搜索树的最近公共祖先.md
// 链接: https://leetcode.cn/problems/lowest-common-ancestor-of-a-binary-search-tree/

func LowestCommonAncestorBST(root, p, q *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	if root.Val > p.Val && root.Val > q.Val {
		return LowestCommonAncestorBST(root.Left,p,q)
	}
	if root.Val < p.Val && root.Val < q.Val {
		return LowestCommonAncestorBST(root.Right,p,q)

	}
	return root
}
