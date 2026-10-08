package lc

/*
给你一个二叉树的根节点 root ， 检查它是否轴对称。

示例 1：

输入：root = [1,2,2,3,4,4,3]
输出：true

示例 2：

输入：root = [1,2,2,null,3,null,3]
输出：false

提示：

- 树中节点数目在范围 [1, 1000] 内

- -100 <= Node.val <= 100

进阶：你可以运用递归和迭代两种方法解决这个问题吗？
*/

func IsSymmetric(root *TreeNode) bool {
	var symmetricSame func(*TreeNode,*TreeNode) bool
	symmetricSame = func(p *TreeNode,q *TreeNode) bool {
		if p == nil && q == nil {
			return true
		}
		if p != nil && q == nil || p == nil && q != nil {
			return false
		}
		if p.Val != q.Val {
			return false
		}
		return symmetricSame(p.Left,q.Right) && symmetricSame(p.Right,q.Left)
	}
	if root == nil {
		return true
	}
	return symmetricSame(root.Left,root.Right)
}

//层序遍历

func IsSymmetricLevelIter(root *TreeNode) bool {
	return false
}
