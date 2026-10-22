package lc

import "slices"

/*
给定一个非负整数数组 nums 和一个整数 k ，你需要将这个数组分成 k* *个非空的连续子数组，使得这 k* *个子数组各自和的最大值 最小。

返回分割后最小的和的最大值。

子数组 是数组中连续的部分。

示例 1：

输入：nums = [7,2,5,10,8], k = 2
输出：18
解释：
一共有四种方法将 nums 分割为 2 个子数组。
其中最好的方式是将其分为 [7,2,5] 和 [10,8] 。
因为此时这两个子数组各自的和的最大值为18，在所有情况中最小。

示例 2：

输入：nums = [1,2,3,4,5], k = 2
输出：9

示例 3：

输入：nums = [1,4,4], k = 3
输出：4

提示：

- 1 <= nums.length <= 1000

- 0 <= nums[i] <= 10^6

- 1 <= k <= min(50, nums.length)
*/

func SplitArray(nums []int, k int) int {
	sum,n := 0,len(nums)
	for _,num := range nums {
		sum += num
	}

	meet := func(T int) bool {
		cnt,i,j := 0,0,0
		for i < n {
			s := 0
			for j = i; j < n && s + nums[j] <= T; j++ {
				s += nums[j]
			}
			cnt += 1
			i = j
		}
		return cnt <= k
	}
  //当meet满足时,如果求最大值最小,显然可以在搜索空间里顺序遍历. 从l,r到遍历,找第一个不满足的,返回t-1那就是 "最小"
	//但是问题来了，我怎么确定这样切，在这个条件下，他是最大呢？
	//以及，这样的搜索，一定是在搜索空间里没漏呢？
	l,r := slices.Max(nums),sum
	for l < r {
		mid := l + (r - l) / 2
		if meet(mid) {
			r = mid
		} else {
			l = mid + 1
		}
	}
	return l
}

// SplitArrayMaxMin 变式: 把 nums 分成恰好 k 个非空连续段, 最大化「最小段和」
// (410 的镜像: min-max → max-min)
func SplitArrayMaxMin(nums []int, k int) int {
	panic("TODO")
}
