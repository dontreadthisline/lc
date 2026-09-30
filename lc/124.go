package lc

import "math"

// MaxPathSum
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
		dfs(root.Left,path)
		path = path[:len(path)-1]
	}
	dfs(root,[]int{})
	return res
}
