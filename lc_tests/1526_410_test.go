package lc_tests

import (
	"fmt"
	"math/rand"
	"testing"

	"demo/lc"
)

// ---- 1526. 形成目标数组的子数组最少增加次数 ----

// refMinNumberOperations: BFS 暴力 —— 从 target 每次「选一段全 >0 的连续子区间 -1」, 到全 0 的最短步数
func refMinNumberOperations(target []int) int {
	cur := append([]int(nil), target...)
	type state []int
	seen := map[string]int{fmt.Sprint(cur): 0}
	queue := []state{cur}
	for len(queue) > 0 {
		s := queue[0]
		queue = queue[1:]
		steps := seen[fmt.Sprint(s)]
		done := true
		for _, v := range s {
			if v != 0 {
				done = false
				break
			}
		}
		if done {
			return steps
		}
		// 枚举所有极大正段之外, 也枚举所有子段(值域小, 直接全枚举)
		n := len(s)
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				ok := true
				for k := i; k <= j; k++ {
					if s[k] <= 0 {
						ok = false
						break
					}
				}
				if !ok {
					break
				}
				ns := append([]int(nil), s...)
				for k := i; k <= j; k++ {
					ns[k]--
				}
				key := fmt.Sprint(ns)
				if _, vis := seen[key]; !vis {
					seen[key] = steps + 1
					queue = append(queue, ns)
				}
			}
		}
	}
	return -1
}

func TestMinNumberOperations(t *testing.T) {
	cases := []struct {
		target []int
		want   int
	}{
		{[]int{1, 2, 3, 2, 1}, 3},
		{[]int{3, 1, 1, 2}, 4},
		{[]int{3, 1, 1, 2, 1}, 4}, // 题面示例(修正版): 3,1,1,2,1 -> 分段 3 | 1,1,2 -> 1 | 2 递增两次
		{[]int{1}, 1},
		{[]int{5}, 5},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.MinNumberOperations(c.target) })
		if got != c.want {
			t.Errorf("MinNumberOperations(%v) = %d, want %d", c.target, got, c.want)
		}
	}
}

func TestMinNumberOperationsRandom(t *testing.T) {
	for i := 0; i < 300; i++ {
		n := 1 + rand.Intn(5)
		target := make([]int, n)
		for j := range target {
			target[j] = rand.Intn(4)
		}
		want := refMinNumberOperations(target)
		var got int
		runSafely(t, func() { got = lc.MinNumberOperations(append([]int(nil), target...)) })
		if got != want {
			t.Fatalf("MinNumberOperations(%v) = %d, want %d(BFS)", target, got, want)
		}
	}
}

// ---- 410. 分割数组的最大值 ----

// refSplitArray: 枚举全部 2^(n-1) 种分割, 最小化各段和的最大值
func refSplitArray(nums []int, k int) int {
	n := len(nums)
	best := 1 << 60
	for mask := 0; mask < 1<<(n-1); mask++ {
		// mask 的每个 1 表示一个分割点
		cuts := []int{}
		for i := 0; i < n-1; i++ {
			if mask&(1<<i) != 0 {
				cuts = append(cuts, i+1)
			}
		}
		if len(cuts)+1 != k {
			continue
		}
		maxSum, prev := 0, 0
		for _, c := range append(cuts, n) {
			s := 0
			for _, v := range nums[prev:c] {
				s += v
			}
			if s > maxSum {
				maxSum = s
			}
			prev = c
		}
		if maxSum < best {
			best = maxSum
		}
	}
	return best
}

func TestSplitArray(t *testing.T) {
	cases := []struct {
		nums []int
		k    int
		want int
	}{
		{[]int{7, 2, 5, 10, 8}, 2, 18},
		{[]int{1, 2, 3, 4, 5}, 2, 9},
		{[]int{1, 4, 4}, 3, 4},
		{[]int{5}, 1, 5},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.SplitArray(c.nums, c.k) })
		if got != c.want {
			t.Errorf("SplitArray(%v, %d) = %d, want %d", c.nums, c.k, got, c.want)
		}
	}
}

func TestSplitArrayRandom(t *testing.T) {
	for i := 0; i < 500; i++ {
		n := 1 + rand.Intn(11)
		k := 1 + rand.Intn(n)
		nums := make([]int, n)
		for j := range nums {
			nums[j] = 1 + rand.Intn(20)
		}
		want := refSplitArray(nums, k)
		var got int
		runSafely(t, func() { got = lc.SplitArray(append([]int(nil), nums...), k) })
		if got != want {
			t.Fatalf("SplitArray(%v, %d) = %d, want %d(枚举)", nums, k, got, want)
		}
	}
}
