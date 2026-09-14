package algo

func PickMinCostWorkers(costs []int, k int, candidate int) int64 {
	var res int
	less := func(a, b []int) bool {
		if a[0] != b[0] {
			return a[0] < b[0]
		}
		return a[1] < b[1]
	}
	h := NewHeap(less)
	n := len(costs)
	l, r := 0, n-1

	for l < r && l < candidate {
		h.Push([]int{costs[l], l})
		l += 1
	}

	for l <= r && n-r-1 < candidate {
		h.Push([]int{costs[r], r})
		r -= 1
	}
	for range k {
		cost := h.Pop()
		res += cost[0]
		if cost[1] < l && l <= r {
			h.Push([]int{costs[l], l})
			l += 1
		} else if cost[1] > r && l <= r {
			h.Push([]int{costs[r], r})
			r -= 1
		}
	}
	return int64(res)
}
