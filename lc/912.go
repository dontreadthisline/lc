package lc

/*
给你一个整数数组 nums，请你将该数组升序排列。

你必须在 不使用任何内置函数 的情况下解决问题，时间复杂度为 O(nlog(n))，并且空间复杂度尽可能小。

示例 1：

输入：nums = [5,2,3,1]
输出：[1,2,3,5]
解释：数组排序后，某些数字的位置没有改变（例如，2 和 3），而其他数字的位置发生了改变（例如，1 和 5）。

示例 2：

输入：nums = [5,1,1,2,0,0]
输出：[0,0,1,1,2,5]
解释：请注意，nums 的值不一定唯一。

提示：

- 1 <= nums.length <= 5 * 10^4

- -5 * 10^4 <= nums[i] <= 5 * 10^4
*/

func SortArray(nums []int) []int {
	quickSort(nums,0,len(nums) -1)
	return nums
}

//quickSort 手写快排
func quickSort(nums []int,l,r int) {
	if l >= r {
		return
	}
	p := partation(nums,l,r)
	quickSort(nums,l,p-1)
	quickSort(nums,p+1,r)
}

func partation(nums []int,l,r int) int {
	p1,p2 := l,l
	pivot := nums[r]
	for ; p2 <= r; p2++ {
		if nums[p2] < pivot {
			nums[p1],nums[p2] = nums[p2],nums[p1]
			p1++
		}
	}
	nums[p1],nums[r] = nums[r],nums[p1]
	return p1
}
