package lc_tests

import (
	"math/rand"
	"testing"

	"demo/lc"
)

// ---- 42. 接雨水 ----

// refTrap: 每个位置的水 = min(左侧最大, 右侧最大) - 自身高度
func refTrap(height []int) int {
	n := len(height)
	total := 0
	for i := 0; i < n; i++ {
		lMax, rMax := 0, 0
		for j := 0; j <= i; j++ {
			if height[j] > lMax {
				lMax = height[j]
			}
		}
		for j := i; j < n; j++ {
			if height[j] > rMax {
				rMax = height[j]
			}
		}
		if m := min(lMax, rMax) - height[i]; m > 0 {
			total += m
		}
	}
	return total
}

func TestTrap(t *testing.T) {
	cases := []struct {
		height []int
		want   int
	}{
		{[]int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}, 6},
		{[]int{4, 2, 0, 3, 2, 5}, 9},
		{[]int{}, 0},
		{[]int{5}, 0},
		{[]int{3, 3, 3}, 0},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.Trap(c.height) })
		if got != c.want {
			t.Errorf("Trap(%v) = %d, want %d", c.height, got, c.want)
		}
	}
}

func TestTrapRandom(t *testing.T) {
	for i := 0; i < 2000; i++ {
		n := rand.Intn(30)
		height := make([]int, n)
		for j := range height {
			height[j] = rand.Intn(15)
		}
		want := refTrap(height)
		var got int
		runSafely(t, func() { got = lc.Trap(height) })
		if got != want {
			t.Fatalf("Trap(%v) = %d, want %d", height, got, want)
		}
	}
}

// ---- 84. 柱状图中最大的矩形 ----

// refLargestRectangleArea: 枚举每对 i<=j, 高度取区间最小值
func refLargestRectangleArea(heights []int) int {
	n := len(heights)
	best := 0
	for i := 0; i < n; i++ {
		h := heights[i]
		for j := i; j < n; j++ {
			if heights[j] < h {
				h = heights[j]
			}
			if area := h * (j - i + 1); area > best {
				best = area
			}
		}
	}
	return best
}

func TestLargestRectangleArea(t *testing.T) {
	cases := []struct {
		heights []int
		want    int
	}{
		{[]int{2, 1, 5, 6, 2, 3}, 10},
		{[]int{2, 4}, 4},
		{nil, 0},
		{[]int{1}, 1},
		{[]int{5, 5, 5}, 15},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.LargestRectangleArea(c.heights) })
		if got != c.want {
			t.Errorf("LargestRectangleArea(%v) = %d, want %d", c.heights, got, c.want)
		}
	}
}

func TestLargestRectangleAreaRandom(t *testing.T) {
	for i := 0; i < 2000; i++ {
		n := 1 + rand.Intn(40)
		heights := make([]int, n)
		for j := range heights {
			heights[j] = 1 + rand.Intn(20)
		}
		want := refLargestRectangleArea(heights)
		var got int
		runSafely(t, func() { got = lc.LargestRectangleArea(heights) })
		if got != want {
			t.Fatalf("LargestRectangleArea(%v) = %d, want %d", heights, got, want)
		}
	}
}
