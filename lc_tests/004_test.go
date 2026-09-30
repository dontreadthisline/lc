package lc_tests

import (
	"math/rand"
	"reflect"
	"slices"
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

// 0704. 二分查找
func TestBinarySearch(t *testing.T) {
	cases := []struct {
		name   string
		nums   []int
		target int
		want   int
	}{
		{"官方用例1_存在", []int{-1, 0, 3, 5, 9, 12}, 9, 4},
		{"官方用例2_不存在", []int{-1, 0, 3, 5, 9, 12}, 2, -1},
		{"单元素命中", []int{5}, 5, 0},
		{"单元素未命中", []int{5}, 4, -1},
		{"首元素", []int{1, 2, 3}, 1, 0},
		{"尾元素", []int{1, 2, 3}, 3, 2},
		{"小于全部", []int{2, 4, 6}, 1, -1},
		{"大于全部", []int{2, 4, 6}, 9, -1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := lc.BinarySearch(c.nums, c.target); got != c.want {
				t.Errorf("BinarySearch(%v, %d) = %d, want %d", c.nums, c.target, got, c.want)
			}
		})
	}
}

func TestBinarySearchRandomCrossCheck(t *testing.T) {
	rng := rand.New(rand.NewSource(31))
	for iter := 0; iter < 20000; iter++ {
		n := rng.Intn(15)
		nums := make([]int, n)
		v := rng.Intn(5) - 2
		for i := range nums {
			v += 1 + rng.Intn(3)
			nums[i] = v
		}
		target := rng.Intn(v+6) - 3
		want := -1
		for i, x := range nums {
			if x == target {
				want = i
				break
			}
		}
		if got := lc.BinarySearch(nums, target); got != want {
			t.Fatalf("iter=%d nums=%v target=%d: got=%d want=%d", iter, nums, target, got, want)
		}
	}
}

// 0034. 在排序数组中查找元素的第一个和最后一个位置
func TestSearchRange(t *testing.T) {
	cases := []struct {
		name   string
		nums   []int
		target int
		want   []int
	}{
		{"官方用例1", []int{5, 7, 7, 8, 8, 10}, 8, []int{3, 4}},
		{"官方用例2_不存在", []int{5, 7, 7, 8, 8, 10}, 6, []int{-1, -1}},
		{"官方用例3_空", []int{}, 0, []int{-1, -1}},
		{"全同元素", []int{7, 7, 7, 7}, 7, []int{0, 3}},
		{"单个命中", []int{1, 2, 3}, 2, []int{1, 1}},
		{"目标小于全部", []int{3, 5}, 1, []int{-1, -1}},
		{"目标大于全部", []int{3, 5}, 9, []int{-1, -1}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := lc.SearchRange(c.nums, c.target)
			if !slices.Equal(got, c.want) {
				t.Errorf("SearchRange(%v, %d) = %v, want %v", c.nums, c.target, got, c.want)
			}
		})
	}
}

func TestSearchRangeRandomCrossCheck(t *testing.T) {
	rng := rand.New(rand.NewSource(32))
	for iter := 0; iter < 20000; iter++ {
		n := rng.Intn(15)
		nums := make([]int, n)
		for i := range nums {
			nums[i] = rng.Intn(5)
		}
		slices.Sort(nums)
		target := rng.Intn(6) - 1
		first, last := -1, -1
		for i, x := range nums {
			if x == target {
				if first == -1 {
					first = i
				}
				last = i
			}
		}
		if got := lc.SearchRange(nums, target); !slices.Equal(got, []int{first, last}) {
			t.Fatalf("iter=%d nums=%v target=%d: got=%v want=[%d %d]", iter, nums, target, got, first, last)
		}
	}
}

// LowerBound：第一个 >= x 的下标
func TestLowerBound(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		x    int
		want int
	}{
		{"命中中间", []int{1, 3, 5, 7}, 5, 2},
		{"命中首", []int{1, 3, 5, 7}, 1, 0},
		{"插入位置在中间", []int{1, 3, 5, 7}, 4, 2},
		{"小于全部", []int{2, 4}, 1, 0},
		{"大于全部", []int{2, 4}, 9, 2},
		{"空表", []int{}, 3, 0},
		{"重复元素取最左", []int{2, 2, 2, 3}, 2, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := lc.LowerBound(c.nums, c.x); got != c.want {
				t.Errorf("LowerBound(%v, %d) = %d, want %d", c.nums, c.x, got, c.want)
			}
		})
	}
	rng := rand.New(rand.NewSource(33))
	for iter := 0; iter < 20000; iter++ {
		n := rng.Intn(15)
		nums := make([]int, n)
		for i := range nums {
			nums[i] = rng.Intn(6)
		}
		slices.Sort(nums)
		x := rng.Intn(8) - 1
		want := len(nums)
		for i, v := range nums {
			if v >= x {
				want = i
				break
			}
		}
		if got := lc.LowerBound(nums, x); got != want {
			t.Fatalf("iter=%d nums=%v x=%d: got=%d want=%d", iter, nums, x, got, want)
		}
	}
}

// 0033. 搜索旋转排序数组
func TestSearchRotated(t *testing.T) {
	t.Run("官方用例", func(t *testing.T) {
		if got := lc.SearchRotated([]int{4, 5, 6, 7, 0, 1, 2}, 0); got != 4 {
			t.Errorf("got %d, want 4", got)
		}
		if got := lc.SearchRotated([]int{4, 5, 6, 7, 0, 1, 2}, 3); got != -1 {
			t.Errorf("got %d, want -1", got)
		}
		if got := lc.SearchRotated([]int{1}, 0); got != -1 {
			t.Errorf("got %d, want -1", got)
		}
	})
	t.Run("未旋转", func(t *testing.T) {
		nums := []int{1, 2, 3, 4, 5}
		for j, v := range nums {
			if got := lc.SearchRotated(nums, v); got != j {
				t.Errorf("未旋转查 %d: got %d, want %d", v, got, j)
			}
		}
	})
}

func TestSearchRotatedExhaustiveCrossCheck(t *testing.T) {
	// 对每个旋转量 k、每个元素穷举 + 若干未命中值
	base := []int{2, 4, 6, 8, 10, 13, 17}
	n := len(base)
	for k := 0; k < n; k++ {
		rotated := make([]int, n)
		for i := range rotated {
			rotated[i] = base[(i+k)%n]
		}
		for j, v := range base {
			want := (j - k + n) % n
			if got := lc.SearchRotated(rotated, v); got != want {
				t.Fatalf("k=%d 查 %d: got %d, want %d", k, v, got, want)
			}
		}
		for _, miss := range []int{1, 3, 5, 99} {
			if got := lc.SearchRotated(rotated, miss); got != -1 {
				t.Fatalf("k=%d 未命中 %d: got %d, want -1", k, miss, got)
			}
		}
	}
}

// 0004. 寻找两个正序数组的中位数（二分切分版）
func TestFindMedianBinarySearchAlgs(t *testing.T) {
	cases := []struct {
		a, b []int
		want float64
	}{
		{[]int{1, 3}, []int{2}, 2},
		{[]int{1, 2}, []int{3, 4}, 2.5},
		{[]int{0, 0}, []int{0, 0}, 0},
		{[]int{}, []int{1}, 1},
		{[]int{2}, []int{}, 2},
		{[]int{1, 2, 3, 4}, []int{5}, 3},
		{[]int{1, 2, 3}, []int{4, 5, 6, 7}, 4},
		{[]int{1, 1, 1}, []int{1, 1}, 1},
	}
	for _, c := range cases {
		if got := lc.FindMedianBinarySearchAlgs(c.a, c.b); got != c.want {
			t.Errorf("BinarySearchAlgs(%v, %v) = %v, want %v", c.a, c.b, got, c.want)
		}
		if got := lc.FindMedianBinarySearchAlgs(c.b, c.a); got != c.want {
			t.Errorf("BinarySearchAlgs(%v, %v) = %v, want %v (参数顺序换边)", c.b, c.a, got, c.want)
		}
	}
}

func TestFindMedianBinarySearchRandomCrossCheck(t *testing.T) {
	rng := rand.New(rand.NewSource(34))
	sortedArr := func() []int {
		m := rng.Intn(12)
		arr := make([]int, m)
		v := rng.Intn(5) - 2
		for i := range arr {
			v += rng.Intn(3)
			arr[i] = v
		}
		return arr
	}
	for iter := 0; iter < 20000; iter++ {
		a, b := sortedArr(), sortedArr()
		if len(a) == 0 && len(b) == 0 {
			b = []int{rng.Intn(5)} // 题目约束 1 <= m+n，双空不合法，补一个
		}
		want := lc.FindMedianMergeSortAlgs(a, b)
		if got := lc.FindMedianBinarySearchAlgs(a, b); got != want {
			t.Fatalf("iter=%d a=%v b=%v: 二分=%v 归并=%v", iter, a, b, got, want)
		}
		if got := lc.FindMedianBinarySearchAlgs(b, a); got != want {
			t.Fatalf("iter=%d 换边 a=%v b=%v: 二分=%v 归并=%v", iter, a, b, got, want)
		}
	}
}

// 0162. 寻找峰值（暴力版属性测试，与二分版同标准）
func isPeakIdx(nums []int, i int) bool {
	left := i == 0 || nums[i-1] < nums[i]
	right := i == len(nums)-1 || nums[i] > nums[i+1]
	return left && right
}

func TestSearchPeakBruteForceAlgs(t *testing.T) {
	cases := [][]int{{1, 2, 3, 1}, {1, 2, 1, 3, 5, 6, 4}, {1, 2}, {3, 1}, {1}, {5}}
	for _, nums := range cases {
		var got int
		panicked := safeCall(func() { got = lc.SearchPeakBruteForceAlgs(nums) })
		if panicked {
			t.Errorf("SearchPeakBruteForceAlgs(%v) panic", nums)
			continue
		}
		if !isPeakIdx(nums, got) {
			t.Errorf("SearchPeakBruteForceAlgs(%v) = %d 不是合法峰值", nums, got)
		}
	}
	rng := rand.New(rand.NewSource(35))
	for iter := 0; iter < 10000; iter++ {
		n := 2 + rng.Intn(12)
		nums := make([]int, n)
		prev := rng.Intn(10)
		for i := range nums {
			dir := 1
			if rng.Intn(2) == 0 {
				dir = -1
			}
			prev += dir * (1 + rng.Intn(3))
			nums[i] = prev
		}
		got := lc.SearchPeakBruteForceAlgs(nums)
		if !isPeakIdx(nums, got) {
			t.Fatalf("iter=%d nums=%v: got=%d 不是合法峰值", iter, nums, got)
		}
	}
}
