package algo

func FindKthSmallest(nums []int, k int) int {
	n := len(nums)
	if k < 1 || k > n {
		panic("k out of range")
	}
	op := func(a, b int) bool {
		return a < b
	}
	return findKthOrder(nums, 0, n-1, k, op)
}

func FindKthLargest(nums []int, k int) int {
	n := len(nums)
	if k < 1 || k > n {
		panic("k out of range")
	}
	op := func(a, b int) bool {
		return a > b
	}
	return findKthOrder(nums, 0, n-1, k, op)
}

func findKthOrder(nums []int, l, r, k int, op func(a, b int) bool) int {
	pivot := nums[r]
	b := partation(l, r, nums, func(a int) bool {
		return op(a, pivot)
	})
	nums[b], nums[r] = nums[r], nums[b]
	//出口不需要显式写：不变式 1 <= k <= r-l+1 保持，且每层区间严格缩小，
	//区间缩到 1 个元素时必有 n == k == 1 触发返回
	n := b - l + 1
	if n == k {
		return nums[b]
	} else if n > k {
		return findKthOrder(nums, l, b-1, k, op)
	} else {
		return findKthOrder(nums, b+1, r, k-n, op)
	}
}

// [l,b-1],[b,j-1],[j,r]
func partation[T any](l, r int, items []T, p func(T) bool) int {
	b := l
	for j := l; j <= r; j++ {
		if p(items[j]) {
			items[b], items[j] = items[j], items[b]
			b += 1
		}
	}
	return b
}
