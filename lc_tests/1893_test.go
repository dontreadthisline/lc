package lc_tests

import (
	"demo/lc"
	"math/rand"
	"testing"
)

// 1893. 检查是否区域内的所有整数都被覆盖（值域差分）
func TestIsCovered(t *testing.T) {
	cases := []struct {
		name   string
		ranges [][]int
		left   int
		right  int
		want   bool
	}{
		{"题面示例1", [][]int{{1, 3}, {2, 6}}, 2, 5, true},
		{"题面示例2", [][]int{{1, 2}, {5, 6}}, 2, 5, false},
		{"大区间起点在left之前", [][]int{{1, 50}}, 2, 50, true},
		{"区间恰从left开始", [][]int{{2, 3}}, 2, 3, true},
		{"零起点区间", [][]int{{0, 1}}, 1, 1, true},
		{"中间缺口", [][]int{{1, 2}, {4, 5}}, 1, 5, false},
		{"三段拼满", [][]int{{1, 2}, {3, 4}, {5, 6}}, 1, 6, true},
		{"窗口右端恰衔接", [][]int{{1, 2}, {2, 3}}, 1, 3, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := lc.IsCovered(c.ranges, c.left, c.right); got != c.want {
				t.Errorf("IsCovered(%v, %d, %d) = %v, want %v", c.ranges, c.left, c.right, got, c.want)
			}
		})
	}
}

// 双暴力对拍：仓库暴力 + 测试内独立暴力
func TestIsCoveredBruteCrossCheck(t *testing.T) {
	independent := func(ranges [][]int, left, right int) bool {
		for x := left; x <= right; x++ {
			ok := false
			for _, rg := range ranges {
				if rg[0] <= x && x <= rg[1] {
					ok = true
					break
				}
			}
			if !ok {
				return false
			}
		}
		return true
	}
	rng := rand.New(rand.NewSource(21))
	for iter := 0; iter < 20000; iter++ {
		m := rng.Intn(8)
		ranges := make([][]int, m)
		for i := range ranges {
			a := rng.Intn(51)
			ranges[i] = []int{a, a + rng.Intn(51-a)}
		}
		left := rng.Intn(50)
		right := left + rng.Intn(51-left)
		want := independent(ranges, left, right)
		if got := lc.IsCovered(ranges, left, right); got != want {
			t.Fatalf("iter=%d ranges=%v [%d,%d]: 差分=%v 独立暴力=%v", iter, ranges, left, right, got, want)
		}
		if got := lc.IsCoveredBruteForce(ranges, left, right); got != want {
			t.Fatalf("iter=%d ranges=%v [%d,%d]: 仓库暴力=%v 独立暴力=%v", iter, ranges, left, right, got, want)
		}
	}
}
