package lc

import "math"

/*
给你一个二叉树的根节点 root ，判断其是否是一个有效的二叉搜索树。

有效 二叉搜索树定义如下：

- 节点的左子树只包含 严格小于 当前节点的数。

- 节点的右子树只包含 严格大于 当前节点的数。

- 所有左子树和右子树自身必须也是二叉搜索树。

示例 1：

输入：root = [2,1,3]
输出：true

示例 2：

输入：root = [5,1,4,null,null,3,6]
输出：false
解释：根节点的值是 5 ，但是右子节点的值是 4 。

提示：

- 树中节点数目范围在[1, 10^4] 内

- -2^31 <= Node.val <= 2^31 - 1
*/

func IsValidBST(root *TreeNode) bool {
	pre := math.MinInt
	var dfs func(*TreeNode) bool
	dfs = func(root *TreeNode) bool{
		if root == nil {
			return true
		}
		if !dfs(root.Left) {
			return false
		}
		if root.Val <= pre {
			return false
		}
		pre = root.Val
		return dfs(root.Right)
	}
	return dfs(root)
}

func IsValidBSTRec(root *TreeNode) bool {
	var checkBST func(*TreeNode,int,int) bool
	checkBST = func(root *TreeNode,lo int,hi int) bool {
		if root == nil {
			return true
		}
		if root.Val >= hi || root.Val <= lo {
			return false
		}
		return checkBST(root.Left,lo,root.Val) && checkBST(root.Right,root.Val,hi)
	}
	return checkBST(root,math.MinInt,math.MaxInt)
}
