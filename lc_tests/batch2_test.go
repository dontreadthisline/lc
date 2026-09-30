package lc_tests

import (
	"demo/lc"
	"slices"
	"testing"
)

// 992. K 个不同整数的子数组（恰好型）
func TestGoodArrays(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		{"官方用例1", []int{1, 2, 1, 2, 3}, 2, 7},
		{"官方用例2", []int{1, 2, 1, 3, 4}, 3, 3},
		{"恰好一种", []int{1, 1, 1}, 1, 6},
		{"k等于全部种数", []int{1, 2, 3}, 3, 1},
		{"k超出种数", []int{1, 1}, 2, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := lc.GoodArrays(c.nums, c.k); got != c.want {
				t.Errorf("GoodArrays(%v, %d) = %d, want %d", c.nums, c.k, got, c.want)
			}
		})
	}
}

// 200. 岛屿数量
func TestNumIslands(t *testing.T) {
	cases := []struct {
		name string
		grid [][]byte
		want int
	}{
		{"官方用例1", [][]byte{
			{'1', '1', '1', '1', '0'},
			{'1', '1', '0', '1', '0'},
			{'1', '1', '0', '0', '0'},
			{'0', '0', '0', '0', '0'}}, 1},
		{"官方用例2", [][]byte{
			{'1', '1', '0', '0', '0'},
			{'1', '1', '0', '0', '0'},
			{'0', '0', '1', '0', '0'},
			{'0', '0', '0', '1', '1'}}, 3},
		{"全水", [][]byte{{'0', '0'}, {'0', '0'}}, 0},
		{"单格岛", [][]byte{{'1'}}, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got int
			panicked := safeCall(func() { got = lc.NumIslands(c.grid) })
			if panicked {
				t.Errorf("NumIslands panic")
				return
			}
			if got != c.want {
				t.Errorf("NumIslands = %d, want %d", got, c.want)
			}
		})
	}
}

// 695. 岛屿的最大面积
func TestMaxAreaOfIsland(t *testing.T) {
	cases := []struct {
		name string
		grid [][]int
		want int
	}{
		{"官方用例1", [][]int{
			{0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0},
			{0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0},
			{0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0}}, 6},
		{"全零", [][]int{{0, 0, 0}, {0, 0, 0}}, 0},
		{"单格岛", [][]int{{1}}, 1},
		{"双格岛", [][]int{{1, 1}}, 2},
		{"两岛比大小", [][]int{{1, 1, 0}, {0, 0, 1}}, 2},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got int
			panicked := safeCall(func() { got = lc.MaxAreaOfIsland(c.grid) })
			if panicked {
				t.Errorf("MaxAreaOfIsland panic")
				return
			}
			if got != c.want {
				t.Errorf("MaxAreaOfIsland = %d, want %d", got, c.want)
			}
		})
	}
}

// 020. 有效的括号
func TestIssValid(t *testing.T) {
	cases := []struct {
		name string
		s    string
		want bool
	}{
		{"官方用例1", "()", true},
		{"官方用例2", "()[]{}", true},
		{"官方用例3", "(]", false},
		{"官方用例4_交叉", "([)]", false},
		{"官方用例5_嵌套", "{[]}", true},
		{"单个左括号", "(", false},
		{"单个右括号", "]", false},
		{"空串", "", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got bool
			panicked := safeCall(func() { got = lc.IssValid(c.s) })
			if panicked {
				t.Errorf("IssValid(%q) panic", c.s)
				return
			}
			if got != c.want {
				t.Errorf("IssValid(%q) = %v, want %v", c.s, got, c.want)
			}
		})
	}
}

// 155. 最小栈
func TestMinStack(t *testing.T) {
	var s *lc.MinStack
	panicked := safeCall(func() { s = lc.NewMinStack() })
	if panicked {
		t.Fatal("NewMinStack panic")
	}
	type op struct {
		name string
		do   func() (int, bool)
	}
	ops := []op{
		{"首次Push", func() (int, bool) { s.Push(2); return 0, false }},
		{"Push", func() (int, bool) { s.Push(0); return 0, false }},
		{"Push", func() (int, bool) { s.Push(3); return 0, false }},
		{"GetMin", func() (int, bool) { return s.GetMin(), true }},
		{"Pop", func() (int, bool) { s.Pop(); return 0, false }},
		{"Top", func() (int, bool) { return s.Top(), true }},
		{"GetMin", func() (int, bool) { return s.GetMin(), true }},
		{"Pop", func() (int, bool) { s.Pop(); return 0, false }},
		{"Top", func() (int, bool) { return s.Top(), true }},
		{"GetMin", func() (int, bool) { return s.GetMin(), true }},
		{"Pop", func() (int, bool) { s.Pop(); return 0, false }},
	}
	wants := []int{0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 0}
	for i, o := range ops {
		var v int
		hasV := false
		p := safeCall(func() { v, hasV = o.do() })
		if p {
			t.Errorf("op[%d] %s panic", i, o.name)
			return
		}
		if hasV && v != wants[i] {
			t.Errorf("op[%d] %s = %d, want %d", i, o.name, v, wants[i])
		}
	}
}

// 304. 二维区域和检索
func TestNumMatrixSumRegion(t *testing.T) {
	matrix := [][]int{
		{3, 0, 1, 4, 2},
		{5, 6, 3, 2, 1},
		{1, 2, 0, 1, 5},
		{4, 1, 0, 1, 7},
		{1, 0, 3, 0, 5},
	}
	var nm *lc.NumMatrix
	panicked := safeCall(func() { nm = lc.NewNumMatrix(matrix) })
	if panicked {
		t.Fatal("NewNumMatrix panic")
	}
	queries := []struct {
		r1, c1, r2, c2, want int
	}{{2, 1, 4, 3, 8}, {1, 1, 2, 2, 11}, {1, 2, 2, 4, 12}, {0, 0, 0, 0, 3}, {0, 0, 4, 4, 58}}
	for _, q := range queries {
		var got int
		p := safeCall(func() { got = nm.SumRegion(q.r1, q.c1, q.r2, q.c2) })
		if p {
			t.Errorf("SumRegion(%d,%d,%d,%d) panic", q.r1, q.c1, q.r2, q.c2)
			continue
		}
		if got != q.want {
			t.Errorf("SumRegion(%d,%d,%d,%d) = %d, want %d", q.r1, q.c1, q.r2, q.c2, got, q.want)
		}
	}
	_ = slices.Max([]int{1})
}
