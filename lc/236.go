package lc

/*
给定一个二叉树, 找到该树中两个指定节点的最近公共祖先。

百度百科中最近公共祖先的定义为：“对于有根树 T 的两个节点 p、q，最近公共祖先表示为一个节点 x，满足 x 是 p、q 的祖先且 x 的深度尽可能大（一个节点也可以是它自己的祖先）。”

示例 1：

输入：root = [3,5,1,6,2,0,8,null,null,7,4], p = 5, q = 1
输出：3
解释：节点 5 和节点 1 的最近公共祖先是节点 3 。

示例 2：

输入：root = [3,5,1,6,2,0,8,null,null,7,4], p = 5, q = 4
输出：5
解释：节点 5 和节点 4 的最近公共祖先是节点 5 。因为根据定义最近公共祖先节点可以为节点本身。

示例 3：

输入：root = [1,2], p = 1, q = 2
输出：1

提示：

- 树中节点数目在范围 [2, 10^5] 内。

- -10^9 <= Node.val <= 10^9

- 所有 Node.val 互不相同 。

- p != q

- p 和 q 均存在于给定的二叉树中。
*/

func lowestCommonAncestor(root *TreeNode, p *TreeNode, q *TreeNode) *TreeNode {
	//exists定义,如果以node为根节点的树,出现过任意节点,那么返回这个根节点
	var exists func(*TreeNode,*TreeNode,*TreeNode) *TreeNode
	exists = func(root *TreeNode,p *TreeNode,q *TreeNode) *TreeNode {
		if root == nil || root == p || root == q {
			return root
		}
		left := exists(root.Left, p,q)
		right := exists(root.Right, p,q)
		if left != nil && right != nil {
			return root
		}
		if left != nil {
			return left
		}
		return right
	}
	return exists(root,p,q)
}
//  1 (1,1,2)
// /
//2

func LowestCommonAncestor(root *TreeNode, p *TreeNode, q *TreeNode) *TreeNode {
	var lcs func(*TreeNode,*TreeNode,*TreeNode) (int,*TreeNode)
	lcs = func(root *TreeNode,q *TreeNode,p*TreeNode) (int,*TreeNode) {
		var total int
		if root == nil {
			return 0,root
		}
		found := 0
		if root == p || root == q {
			found += 1
		}
		cnt1,left := lcs(root.Left,p,q) //2,1,2
		cnt2,right := lcs(root.Right,p,q)
		total = found + cnt1 + cnt2
		if cnt1 == 1 && cnt2 == 1 && left == nil  && right == nil {
			return total,root
		} else if found == 1 && total == 2 && left == nil && right == nil {//对应p是q的直接祖先的情况
			return total,root
		} else if left != nil {
			return total,left
		} else if right != nil {
			return total,right
		}
		return total,nil
	}
	_,node := lcs(root,p,q)
	return node
}
