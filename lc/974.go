package lc

/*
给定一个整数数组 nums 和一个整数 k ，返回其中元素之和可被 k 整除的非空 子数组 的数目。

子数组 是数组中 连续 的部分。

示例 1：

输入：nums = [4,5,0,-2,-3,1], k = 5
输出：7
解释：
有 7 个子数组满足其元素之和可被 k = 5 整除：
[4, 5, 0, -2, -3, 1], [5], [5, 0], [5, 0, -2, -3], [0], [0, -2, -3], [-2, -3]

示例 2:

输入: nums = [5], k = 9
输出: 0

提示:

- 1 <= nums.length <= 3 * 10^4

- -10^4 <= nums[i] <= 10^4

- 2 <= k <= 10^4
*/

//sum[i,j] % 5 == 0
//(pre[j+1] - pre[i] ) % K == 0

func SubarraysDivByK(nums []int, k int) int {
	n, res := len(nums), 0
	prefix, m := make([]int, n+1), make(map[int]int, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] + nums[i-1]
	}
	r := func(x int) int {
		return (x%k + k) % k
	}
	for j := range n {
		/*
			for i := j; i >= 0; i-- {
				if prefix[j+1] % k ==  prefix[i] % k {
					res += 1
				}
			}*/
		m[r(prefix[j])] += 1
		res += m[r(prefix[j+1])]
	}
	return res
}
