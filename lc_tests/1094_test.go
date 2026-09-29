package lc_tests

import (
	"demo/lc"
	"math/rand"
	"testing"
)

// 1094. 拼车（差分）
func TestCarPooling(t *testing.T) {
	cases := []struct {
		name     string
		trips    [][]int
		capacity int
		want     bool
	}{
		{"官方用例1_超载", [][]int{{2, 1, 5}, {3, 3, 7}}, 4, false},
		{"官方用例2_恰好", [][]int{{2, 1, 5}, {3, 3, 7}}, 5, true},
		{"换乘衔接_前客下车点等于后客上车点", [][]int{{2, 1, 5}, {3, 5, 7}}, 4, true},
		{"换乘衔接_真超载", [][]int{{2, 1, 5}, {3, 5, 7}}, 2, false},
		{"单行程_满载", [][]int{{4, 1, 3}}, 4, true},
		{"单行程_超载", [][]int{{5, 1, 3}}, 4, false},
		{"相邻上下车", [][]int{{2, 0, 1}, {3, 1, 2}}, 5, true},
		{"位置零上车且超载", [][]int{{3, 0, 2}}, 2, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := lc.CarPooling(c.trips, c.capacity); got != c.want {
				t.Errorf("CarPooling(%v, %d) = %v, want %v", c.trips, c.capacity, got, c.want)
			}
		})
	}
}

// 暴力对拍：逐站数在车人数（占用区间 [from, to)）
func bruteCarPooling(trips [][]int, capacity int) bool {
	for loc := 0; loc <= 1000; loc++ {
		onboard := 0
		for _, tr := range trips {
			if tr[1] <= loc && loc < tr[2] {
				onboard += tr[0]
			}
		}
		if onboard > capacity {
			return false
		}
	}
	return true
}

func TestCarPoolingBruteCrossCheck(t *testing.T) {
	rng := rand.New(rand.NewSource(14))
	for iter := 0; iter < 20000; iter++ {
		m := rng.Intn(10)
		trips := make([][]int, m)
		for i := range trips {
			from := rng.Intn(15)
			to := from + 1 + rng.Intn(4)
			trips[i] = []int{1 + rng.Intn(6), from, to}
		}
		capacity := 1 + rng.Intn(20)
		if got, want := lc.CarPooling(trips, capacity), bruteCarPooling(trips, capacity); got != want {
			t.Fatalf("iter=%d trips=%v cap=%d: got=%v want=%v", iter, trips, capacity, got, want)
		}
	}
}
