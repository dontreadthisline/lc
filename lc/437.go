package lc

// 437. 路径总和 III
// 题面: ../mother-problems/0437-路径总和III.md
// 链接: https://leetcode.cn/problems/path-sum-iii/

func PathSum(root *TreeNode, targetSum int) int {
	var cnt int
	//o(nlog(n))
	var dfs func(*TreeNode,[]int)
	dfs = func(root *TreeNode,path []int) {
		if root == nil {
			return
		}
		path = append(path,root.Val)
		for sum,i := 0,len(path)-1; i >= 0; i-- {
			sum += path[i]
			if sum == targetSum {
				cnt++
			}
		}
		dfs(root.Left,path)
		dfs(root.Right,path)
		//这里从来没有 path = path[0:len(path)-1]
	}
	dfs(root,[]int{})
	return cnt
}

func PathSumPrefixSumCache(root *TreeNode,targetSum int) int {
	var cnt int
	var dfs func(*TreeNode)
	dfs = func(root *TreeNode) {

	}
	dfs(root)
	return cnt
}
