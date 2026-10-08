package lc

// 450. 删除二叉搜索树中的节点
// 题面: ../mother-problems/0450-删除二叉搜索树中的节点.md
// 链接: https://leetcode.cn/problems/delete-node-in-a-bst/
//二叉搜索树的插入和删除

func DeleteNode(root *TreeNode, key int) *TreeNode {
	return nil
}

func InsertNode(root *TreeNode, key int) *TreeNode {
	return nil
}

func SearchNode(root *TreeNode,key int) *TreeNode {
	if root == nil {
		return nil
	}
	if root.Val == key {
		return root
	}
	if root.Val < key {
		return SearchNode(root.Right, key)
	} else {
		return SearchNode(root.Left, key)
	}
}

