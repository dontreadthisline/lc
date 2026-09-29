package lc_tests

import (
	"reflect"
	"testing"

	"demo/lc"
)

func TestFindMedian(t *testing.T) {
	cases := []struct {
		a, b []int
		want float64
	}{
		{[]int{1, 3}, []int{2}, 2},
		{[]int{1, 2}, []int{3, 4}, 2.5},
		{[]int{0, 0}, []int{0, 0}, 0},
		{[]int{}, []int{1}, 1},
		{[]int{2}, []int{}, 2},
	}
	for _, c := range cases {
		if got := lc.FindMedianMergeSortAlgs(c.a, c.b); got != c.want {
			t.Errorf("Merge(%v, %v) = %v, want %v", c.a, c.b, got, c.want)
		}
	}
}

func TestSearchPeak(t *testing.T) {
	cases := []struct {
		nums  []int
		peaks []int
	}{
		{[]int{1, 2, 3, 1}, []int{2}},
		{[]int{1, 2, 1, 3, 5, 6, 4}, []int{1, 5}},
		{[]int{1, 2}, []int{1}},
		{[]int{3, 1, 2}, []int{0, 2}},
	}
	isPeak := func(nums []int, i int) bool {
		left := i == 0 || nums[i-1] < nums[i]
		right := i == len(nums)-1 || nums[i] > nums[i+1]
		return left && right
	}
	for _, c := range cases {
		got := lc.SearchPeakBinarySearchAlgs(c.nums)
		if got < 0 || got >= len(c.nums) || !isPeak(c.nums, got) {
			t.Errorf("SearchPeak(%v) = %d, 不是合法峰", c.nums, got)
		}
	}
}

func TestMoveZeroes(t *testing.T) {
	nums := []int{0, 1, 0, 3, 12}
	lc.MoveZeroes(nums)
	if !reflect.DeepEqual(nums, []int{1, 3, 12, 0, 0}) {
		t.Errorf("MoveZeroes = %v, want [1,3,12,0,0]", nums)
	}
	nums = []int{0}
	lc.MoveZeroes(nums)
	if !reflect.DeepEqual(nums, []int{0}) {
		t.Errorf("MoveZeroes([0]) = %v, want [0]", nums)
	}
}

func TestMaxArea(t *testing.T) {
	if got := lc.MaxArea([]int{1, 8, 6, 2, 5, 4, 8, 3, 7}); got != 49 {
		t.Errorf("MaxArea = %d, want 49", got)
	}
	if got := lc.MaxArea([]int{1, 1}); got != 1 {
		t.Errorf("MaxArea([1,1]) = %d, want 1", got)
	}
}

func TestTransmite(t *testing.T) {
	cases := []struct {
		weights []int
		days    int
		want    int
	}{
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 5, 15},
		{[]int{3, 2, 2, 4, 1, 4}, 3, 6},
		{[]int{1, 2, 3, 1, 1}, 4, 3},
	}
	for _, c := range cases {
		if got := lc.Transmite(c.weights, c.days); got != c.want {
			t.Errorf("Transmite(%v, %d) = %d, want %d", c.weights, c.days, got, c.want)
		}
	}
}

func TestEtaFuckBananas(t *testing.T) {
	if got := lc.EtaFuckBananas([]int{3, 6, 7, 11}, 8); got != 4 {
		t.Errorf("EtaFuckBananas = %d, want 4", got)
	}
	if got := lc.EtaFuckBananas([]int{30, 11, 23, 4, 20}, 5); got != 30 {
		t.Errorf("EtaFuckBananas = %d, want 30", got)
	}
}

func TestSqrt(t *testing.T) {
	if got := lc.Sqrt(4); got != 2 {
		t.Errorf("Sqrt(4) = %d, want 2", got)
	}
	if got := lc.Sqrt(8); got != 2 {
		t.Errorf("Sqrt(8) = %d, want 2", got)
	}
}

func TestSqrtSearchAlgs(t *testing.T) {
	if got := lc.SqrtSearchAlgs(4); got != 2 {
		t.Errorf("SqrtSearchAlgs(4) = %d, want 2", got)
	}
	if got := lc.SqrtSearchAlgs(8); got != 2 {
		t.Errorf("SqrtSearchAlgs(8) = %d, want 2", got)
	}
}

func TestSearchMatrixSearchWay(t *testing.T) {
	m := [][]int{{1, 2}, {1, 3}}
	for k, want := range map[int]int{1: 1, 2: 1, 3: 2, 4: 3} {
		if got := lc.SearchMatrixSearchWay(m, k); got != want {
			t.Errorf("SearchMatrixSearchWay(k=%d) = %d, want %d", k, got, want)
		}
	}
}

func TestMinWindow(t *testing.T) {
	cases := []struct {
		s, t, want string
	}{
		{"ADOBECODEBANC", "ABC", "BANC"},
		{"a", "a", "a"},
		{"a", "aa", ""},
	}
	for _, c := range cases {
		if got := lc.MinWindow(c.s, c.t); got != c.want {
			t.Errorf("MinWindow(%q, %q) = %q, want %q", c.s, c.t, got, c.want)
		}
	}
}
