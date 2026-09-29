package lc

/*
给定一个数组 nums，编写一个函数将所有 0 移动到数组的末尾，同时保持非零元素的相对顺序。

请注意 ，必须在不复制数组的情况下原地对数组进行操作。

示例 1:

输入: nums = [0,1,0,3,12]
输出: [1,3,12,0,0]

示例 2:

输入: nums = [0]
输出: [0]

提示:

- 1 <= nums.length <= 10^4

- -2^31 <= nums[i] <= 2^31 - 1

进阶：你能尽量减少完成的操作次数吗？
*/

// 1:自定义排序
// 2:开辟新数组,然后copy
// 4:就地移动的方法
//
//3:partation的方法,就地算法,且是O(n)
func MoveZeroes(nums []int) {
	PartationAlgs(nums)
}

func PartationAlgs(nums []int) {
	//idx < i where nums[idx] meet the condition
	i, j, n := 0, 0, len(nums)
	for j < n {
		if nums[j] != 0 {
			nums[i], nums[j] = nums[j], nums[i]
			i += 1
		}
		j += 1
	}
}

func InplaceMoveZero(nums []int) {
	n := len(nums)
	i, j := 0, 0
	for ; j < n; j++ {
		if nums[j] != 0 {
			nums[i] = nums[j]
			i += 1
		}
	}
	for ; i < n; i++ {
		nums[i] = 0
	}
}
