package lc

// 114. 二叉树展开为链表
// 题面: ../mother-problems/0114-二叉树展开为链表.md
// 链接: https://leetcode.cn/problems/flatten-binary-tree-to-linked-list/
/*
```
    1               1
   / \               \
  2   5        →      2
 / \   \               \
3   4   6               3
                         \
                          4
                           \
                            5
                             \
                              6
```
*/

func Flatten(root *TreeNode) {
	if root == nil {
		return
	}
	Flatten(root.Left)
	Flatten(root.Right)
	//没有左子树,直接返回
	if node := root.Left; node != nil {
		for ;node.Right != nil; node = node.Right {
		}
		node.Right = root.Right
		root.Right = root.Left
		root.Left = nil
	}
}
