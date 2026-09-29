package lc

/*
给你二叉树的根节点 root ，返回其节点值的 层序遍历 。 （即逐层地，从左到右访问所有节点）。

示例 1：

输入：root = [3,9,20,null,null,15,7]
输出：[[3],[9,20],[15,7]]

示例 2：

输入：root = [1]
输出：[[1]]

示例 3：

输入：root = []
输出：[]

提示：

- 树中节点数目在范围 [0, 2000] 内

- -1000 <= Node.val <= 1000
*/

func LevelOrder(root *TreeNode) [][]int {
	res ,queue:= make([][]int,0), make([]*TreeNode,0)

	if root != nil {
		queue = append(queue,root)
	}

	for len(queue) > 0 {
		sz := len(queue)
		level := []int{}
		for range sz {
			node := queue[0]
			level  = append(level,node.Val)
			if node.Left != nil {
				queue = append(queue,node.Left)
			}
			if node.Right != nil {
				queue = append(queue ,node.Right)
			}
			queue = queue[1:]
		}
		res = append(res,level)
	}
	return res
}
