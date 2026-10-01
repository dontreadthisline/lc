package lc

import "math"

// MaxPathSum 最大路径和
func MaxPathSum(root *TreeNode) int {
	res := math.MinInt
	var dfs func(*TreeNode,[]int)
	dfs = func(root *TreeNode,path []int) {
		if root == nil {
			return
		}
		path = append(path,root.Val)
		sum := 0
		for i := len(path) - 1; i >= 0; i-- {
			sum += path[i]
			res = max(res,sum)
		}
		dfs(root.Left,path)
		dfs(root.Right,path)
	}
	dfs(root,[]int{})
	return res
}
