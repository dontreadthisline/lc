package lc

import "slices"

/*
给你一个整数数组 nums ，数组中的元素 互不相同 。返回该数组所有可能的子集（幂集）。

解集 不能 包含重复的子集。你可以按 任意顺序 返回解集。

示例 1：

输入：nums = [1,2,3]
输出：[[],[1],[2],[1,2],[3],[1,3],[2,3],[1,2,3]]

示例 2：

输入：nums = [0]
输出：[[],[0]]

提示：

- 1 <= nums.length <= 10

- -10 <= nums[i] <= 10

- nums 中的所有元素 互不相同
*/
func Subsets(nums []int) [][]int {
	var res [][]int
	var backtrack func(int, []int)
	backtrack = func(start int, path []int) {
		res = append(res, slices.Clone(path))    // 每个节点都收获(含根=空集); 必须拷贝: path 底层数组共享, 兄弟分支的 append 会覆写(cap 有余量时不扩容)
		for i := start + 1; i < len(nums); i++ { // 下标递增剪枝即循环起点; used 冗余已删(递增下标不可能重复); 无终止剪枝: 树叶 for 天然空转
			backtrack(i, append(path, nums[i])) // path 传参不回退; 零共享可变状态, 零撤销
		}
	}
	backtrack(-1, nil) // 下标哨兵: -1 = 空集(比最小下标 0 小一)
	return res
}
