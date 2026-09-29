package lc

/*
给你一个整数数组 nums 和一个整数 k ，请你统计并返回 *该数组中和为 k 的子数组的个数 *。

子数组是数组中元素的连续非空序列。

示例 1：

输入：nums = [1,1,1], k = 2
输出：2

示例 2：

输入：nums = [1,2,3], k = 3
输出：2

提示：

- 1 <= nums.length <= 2 * 10^4

- -1000 <= nums[i] <= 1000

- -10^7 <= k <= 10^7
*/
//没有nums[i] >= 0的约束了,所以没法用atMost(n) - atMost(n-1)的滑动窗口来做了

//pre[i] = {0 where i = 0,nums[0] + ... + nums[i-1]  where i in 1..n}
//sum(i,j) = pre[j+1] - pre[i] = k =>整理有 k + pre[i] == pre[j+1]  where j in 0..n-1 and i <= j
//sum(arr[i:j]) = pre[j+1] - pre[i] where j in 0..n-1 and i <= j

func SubarraySum(nums []int, k int) int {
	n := len(nums)
	prefix := make([]int, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] + nums[i-1]
	}
	var res int
	for j := range n {
		for i := j; i >= 0; i-- {
			if prefix[j+1]-prefix[i] == k {
				res += 1
			}
		}
	}
	return res
}

func SubarraySumHashway(nums []int, k int) int {
	n, res := len(nums), 0
	prefix, m := make([]int, n+1), make(map[int]int, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] + nums[i-1]
	}
	for j := range n {
		/*
			  sum(i,j) = pre[j+1] - pre[i] = k => 整理有 pre[i] == pre[j+1] - k  where j in 0..n-1
				for i := j; i >= 0; i-- {
					if prefix[j+1]-k == prefix[i] {
						res += 1
					}
				}*/
		m[prefix[j]] += 1
		res += m[prefix[j+1]-k]
	}
	return res
}
