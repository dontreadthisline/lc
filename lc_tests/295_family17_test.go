package lc_tests

import (
	"math/rand"
	"sort"
	"testing"

	"demo/lc"
)

// ---- 295. 数据流的中位数 ----

func TestMedianFinder(t *testing.T) {
	var mf lc.MedianFinder
	runSafely(t, func() { mf = lc.Constructor() })
	// 操作序列 + 排序参照
	stream := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	buffer := []int{}
	for _, v := range stream {
		runSafely(t, func() { mf.AddNum(v) })
		buffer = append(buffer, v)
		sorted := append([]int(nil), buffer...)
		sort.Ints(sorted)
		var want float64
		if len(sorted)%2 == 1 {
			want = float64(sorted[len(sorted)/2])
		} else {
			want = float64(sorted[len(sorted)/2-1]+sorted[len(sorted)/2]) / 2
		}
		var got float64
		runSafely(t, func() { got = mf.FindMedian() })
		if got != want {
			t.Fatalf("AddNum 流 %v 后 FindMedian() = %v, want %v", buffer, got, want)
		}
	}
}

func TestMedianFinderRandom(t *testing.T) {
	var mf lc.MedianFinder
	runSafely(t, func() { mf = lc.Constructor() })
	buffer := []int{}
	for i := 0; i < 2000; i++ {
		v := rand.Intn(2001) - 1000
		runSafely(t, func() { mf.AddNum(v) })
		buffer = append(buffer, v)
		if i%7 != 0 {
			continue
		}
		sorted := append([]int(nil), buffer...)
		sort.Ints(sorted)
		var want float64
		if len(sorted)%2 == 1 {
			want = float64(sorted[len(sorted)/2])
		} else {
			want = float64(sorted[len(sorted)/2-1]+sorted[len(sorted)/2]) / 2
		}
		var got float64
		runSafely(t, func() { got = mf.FindMedian() })
		if got != want {
			t.Fatalf("随机流长度 %d: FindMedian() = %v, want %v", len(buffer), got, want)
		}
	}
}

// ---- 48. 旋转图像 ----

// refRotate90: 期望矩阵直接按坐标变换生成 —— rotated[i][j] = orig[n-1-j][i]
func refRotate90(matrix [][]int) [][]int {
	n := len(matrix)
	out := make([][]int, n)
	for i := range out {
		out[i] = make([]int, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			out[i][j] = matrix[n-1-j][i]
		}
	}
	return out
}

func cloneMatrix(m [][]int) [][]int {
	cp := make([][]int, len(m))
	for i := range m {
		cp[i] = append([]int(nil), m[i]...)
	}
	return cp
}

func TestRotate(t *testing.T) {
	cases := [][][]int{
		{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
		{{5, 1, 9, 11}, {2, 4, 8, 10}, {13, 3, 6, 7}, {15, 14, 12, 16}},
		{{1}},
	}
	for _, c := range cases {
		orig := cloneMatrix(c)
		want := refRotate90(orig)
		got := cloneMatrix(c)
		runSafely(t, func() { lc.Rotate(got) })
		for i := range want {
			for j := range want[i] {
				if got[i][j] != want[i][j] {
					t.Errorf("Rotate(%v) 得到 %v, want %v", orig, got, want)
					break
				}
			}
		}
	}
}

func TestRotateRandom(t *testing.T) {
	for i := 0; i < 300; i++ {
		n := 1 + rand.Intn(10)
		m := make([][]int, n)
		for r := range m {
			m[r] = make([]int, n)
			for c := range m[r] {
				m[r][c] = rand.Intn(100)
			}
		}
		want := refRotate90(m)
		got := cloneMatrix(m)
		runSafely(t, func() { lc.Rotate(got) })
		for r := range want {
			for c := range want[r] {
				if got[r][c] != want[r][c] {
					t.Fatalf("Rotate 随机矩阵 (%d,%d) = %d, want %d", r, c, got[r][c], want[r][c])
				}
			}
		}
	}
}

// ---- 54. 螺旋矩阵 ----

// refSpiralOrder: 独立按层模拟(与典型解法结构不同)
func refSpiralOrder(matrix [][]int) []int {
	m, n := len(matrix), len(matrix[0])
	var out []int
	top, bottom, left, right := 0, m-1, 0, n-1
	for top <= bottom && left <= right {
		for j := left; j <= right; j++ {
			out = append(out, matrix[top][j])
		}
		top++
		for i := top; i <= bottom; i++ {
			out = append(out, matrix[i][right])
		}
		right--
		if top <= bottom {
			for j := right; j >= left; j-- {
				out = append(out, matrix[bottom][j])
			}
			bottom--
		}
		if left <= right {
			for i := bottom; i >= top; i-- {
				out = append(out, matrix[i][left])
			}
			left++
		}
	}
	return out
}

func TestSpiralOrder(t *testing.T) {
	cases := []struct {
		matrix [][]int
		want   []int
	}{
		{[][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, []int{1, 2, 3, 6, 9, 8, 7, 4, 5}},
		{[][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}}, []int{1, 2, 3, 4, 8, 12, 11, 10, 9, 5, 6, 7}},
		{[][]int{{1}}, []int{1}},
		{[][]int{{1}, {2}}, []int{1, 2}},
	}
	for _, c := range cases {
		var got []int
		runSafely(t, func() { got = lc.SpiralOrder(cloneMatrix(c.matrix)) })
		if !equalInts(got, c.want) {
			t.Errorf("SpiralOrder(%v) = %v, want %v", c.matrix, got, c.want)
		}
	}
}

func TestSpiralOrderRandom(t *testing.T) {
	for i := 0; i < 500; i++ {
		m := 1 + rand.Intn(8)
		n := 1 + rand.Intn(8)
		matrix := make([][]int, m)
		for r := range matrix {
			matrix[r] = make([]int, n)
			for c := range matrix[r] {
				matrix[r][c] = rand.Intn(100)
			}
		}
		want := refSpiralOrder(matrix)
		var got []int
		runSafely(t, func() { got = lc.SpiralOrder(cloneMatrix(matrix)) })
		if !equalInts(got, want) {
			t.Fatalf("SpiralOrder 随机 %dx%d 不一致", m, n)
		}
	}
}

// ---- 73. 矩阵置零 ----

// refSetZeroes: 复制标记法(非原地, 独立参照)
func refSetZeroes(matrix [][]int) [][]int {
	m, n := len(matrix), len(matrix[0])
	out := cloneMatrix(matrix)
	zeroRows := map[int]bool{}
	zeroCols := map[int]bool{}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if matrix[i][j] == 0 {
				zeroRows[i] = true
				zeroCols[j] = true
			}
		}
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if zeroRows[i] || zeroCols[j] {
				out[i][j] = 0
			}
		}
	}
	return out
}

func TestSetZeroes(t *testing.T) {
	cases := [][][]int{
		{{1, 1, 1}, {1, 0, 1}, {1, 1, 1}},
		{{0, 1, 2, 0}, {3, 4, 5, 2}, {1, 3, 1, 5}},
		{{1}},
		{{0}},
	}
	for _, c := range cases {
		want := refSetZeroes(c)
		got := cloneMatrix(c)
		runSafely(t, func() { lc.SetZeroes(got) })
		for i := range want {
			for j := range want[i] {
				if got[i][j] != want[i][j] {
					t.Errorf("SetZeroes(%v) 得到 %v, want %v", c, got, want)
					break
				}
			}
		}
	}
}

func TestSetZeroesRandom(t *testing.T) {
	for i := 0; i < 500; i++ {
		m := 1 + rand.Intn(8)
		n := 1 + rand.Intn(8)
		matrix := make([][]int, m)
		for r := range matrix {
			matrix[r] = make([]int, n)
			for c := range matrix[r] {
				if rand.Intn(8) == 0 {
					matrix[r][c] = 0
				} else {
					matrix[r][c] = 1 + rand.Intn(9)
				}
			}
		}
		want := refSetZeroes(matrix)
		got := cloneMatrix(matrix)
		runSafely(t, func() { lc.SetZeroes(got) })
		for r := range want {
			for c := range want[r] {
				if got[r][c] != want[r][c] {
					t.Fatalf("SetZeroes 随机矩阵 (%d,%d) = %d, want %d", r, c, got[r][c], want[r][c])
				}
			}
		}
	}
}
