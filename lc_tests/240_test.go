package lc_tests

import (
	"math/rand"
	"testing"

	"demo/lc"
)

// ---- 240. 搜索二维矩阵 II ----

// genSortedMatrix: 生成行、列均严格递增的随机矩阵
func genSortedMatrix(m, n, span int) [][]int {
	mat := make([][]int, m)
	for i := range mat {
		mat[i] = make([]int, n)
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			base := 0
			if i > 0 && mat[i-1][j] > base {
				base = mat[i-1][j]
			}
			if j > 0 && mat[i][j-1] > base {
				base = mat[i][j-1]
			}
			mat[i][j] = base + 1 + rand.Intn(span)
		}
	}
	return mat
}

func TestSearchMatrixII(t *testing.T) {
	mat := [][]int{
		{1, 4, 7, 11, 15},
		{2, 5, 8, 12, 19},
		{3, 6, 9, 16, 22},
		{10, 13, 14, 17, 24},
		{18, 21, 23, 26, 30},
	}
	cases := []struct {
		target int
		want   bool
	}{
		{5, true},
		{20, false},
		{1, true},   // 左上角
		{30, true},  // 右下角
		{15, true},  // 右上角
		{18, true},  // 左下角
		{0, false},
		{31, false},
	}
	for _, c := range cases {
		var got bool
		runSafely(t, func() { got = lc.SearchMatrix(mat, c.target) })
		if got != c.want {
			t.Errorf("SearchMatrix(%d) = %v, want %v", c.target, got, c.want)
		}
	}
}

func TestSearchMatrixIIRandom(t *testing.T) {
	for i := 0; i < 1000; i++ {
		m := 1 + rand.Intn(8)
		n := 1 + rand.Intn(8)
		mat := genSortedMatrix(m, n, 5)
		var target int
		if rand.Intn(2) == 0 {
			target = mat[rand.Intn(m)][rand.Intn(n)] // 必在
		} else {
			target = rand.Intn(80) // 未必在
		}
		// 朴素参照: 全扫
		want := false
		for r := range mat {
			for _, v := range mat[r] {
				if v == target {
					want = true
				}
			}
		}
		var got bool
		runSafely(t, func() { got = lc.SearchMatrix(mat, target) })
		if got != want {
			t.Fatalf("SearchMatrix(矩阵 %dx%d, %d) = %v, want %v", m, n, target, got, want)
		}
	}
}
