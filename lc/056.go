package lc

import "slices"

/*
以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [start_i, end_i] 。请你合并所有重叠的区间，并返回 *一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间* 。

示例 1：

输入：intervals = [[1,3],[2,6],[8,10],[15,18]]
输出：[[1,6],[8,10],[15,18]]
解释：区间 [1,3] 和 [2,6] 重叠, 将它们合并为 [1,6].

示例 2：

输入：intervals = [[1,4],[4,5]]
输出：[[1,5]]
解释：区间 [1,4] 和 [4,5] 可被视为重叠区间。

示例 3：

输入：intervals = [[4,7],[1,4]]
输出：[[1,7]]
解释：区间 [1,4] 和 [4,7] 可被视为重叠区间。

提示：

- 1 <= intervals.length <= 10^4

- intervals[i].length == 2

- 0 <= start_i <= end_i <= 10^4
*/

func MergeRange(intervals [][]int) [][]int {
	slices.SortFunc(intervals,func(a,b []int) int {
		return a[0] - b[0]
	})
	i,j,n := 1,0,len(intervals)
	res := make([][]int,0,n)
	res = append(res, intervals[0])
	for ; i < n; i++ {
		if intervals[i][0] <= res[j][1] {
			res[j] = []int{res[j][0],max(intervals[i][1],res[j][1])}
			//res[j] = []int{res[j][0],intervals[i][1]}
		} else {
			res = append(res, intervals[i])
			j++
		}
	}
	return res
}
