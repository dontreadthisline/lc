package lc_tests

import (
	"demo/lc"
	"math/rand"
	"slices"
	"testing"
)

// ============ 家族 4：二分 ============

// 074. 搜索二维矩阵
func TestSearchSortedMatrix(t *testing.T) {
	m := [][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}}
	for _, c := range []struct {
		target int
		want   bool
	}{{3, true}, {13, false}, {1, true}, {60, true}, {0, false}, {61, false}} {
		if got := lc.SearchSortedMatrix(m, c.target); got != c.want {
			t.Errorf("SearchSortedMatrix(%d) = %v, want %v", c.target, got, c.want)
		}
	}
	if got := lc.SearchSortedMatrix([][]int{{1}}, 1); !got {
		t.Error("单元素命中")
	}
	if got := lc.SearchSortedMatrix([][]int{{1}}, 2); got {
		t.Error("单元素未命中")
	}
}

// 153. 寻找旋转排序数组中的最小值
func TestFindMin(t *testing.T) {
	cases := []struct {
		nums []int
		want int
	}{
		{[]int{3, 4, 5, 1, 2}, 1},
		{[]int{4, 5, 6, 7, 0, 1, 2}, 0},
		{[]int{11, 13, 15, 17}, 11},
		{[]int{1}, 1},
		{[]int{2, 1}, 1}, // n=2 边界
		{[]int{1, 2}, 1}, // 未旋转 n=2
	}
	for _, c := range cases {
		var got int
		panicked := safeCall(func() { got = lc.FindMin(slices.Clone(c.nums)) })
		if panicked {
			t.Errorf("FindMin(%v) panic", c.nums)
			continue
		}
		if got != c.want {
			t.Errorf("FindMin(%v) = %d, want %d", c.nums, got, c.want)
		}
	}
	rng := rand.New(rand.NewSource(71))
	for iter := 0; iter < 5000; iter++ {
		n := 1 + rng.Intn(20)
		sorted := make([]int, n)
		for i := range sorted {
			sorted[i] = rng.Intn(100)
		}
		slices.Sort(sorted)
		k := rng.Intn(n)
		rotated := append(slices.Clone(sorted[k:]), sorted[:k]...)
		want := slices.Min(sorted)
		got := lc.FindMin(rotated)
		if got != want {
			t.Fatalf("iter=%d nums=%v: got=%d want=%d", iter, rotated, got, want)
		}
	}
}

// 278. 第一个错误的版本（需要 mock isBadVersion，这里只测逻辑框架）
// 注意：isBadVersion 是包内函数，无法直接 mock。跳过自动测试。

// ============ 家族 5：栈 ============

// 739. 每日温度
func TestDailyTemperatures(t *testing.T) {
	cases := []struct {
		nums []int
		want []int
	}{
		{[]int{73, 74, 75, 71, 69, 72, 76, 73}, []int{1, 1, 4, 2, 1, 1, 0, 0}},
		{[]int{30, 40, 50, 60}, []int{1, 1, 1, 0}},
		{[]int{30, 60, 90}, []int{1, 1, 0}},
		{[]int{89, 62, 70, 58, 47, 47, 46, 76, 100, 70}, []int{8, 1, 5, 4, 3, 2, 1, 1, 0, 0}},
		{[]int{50}, []int{0}},
	}
	impls := []struct {
		name string
		fn   func([]int) []int
	}{
		{"DailyTemperatures", lc.DailyTemperatures},
		{"DailyTemperatures1", lc.DailyTemperatures1},
	}
	for _, impl := range impls {
		for _, c := range cases {
			t.Run(impl.name, func(t *testing.T) {
				var got []int
				panicked := safeCall(func() { got = impl.fn(slices.Clone(c.nums)) })
				if panicked {
					t.Errorf("%s(%v) panic", impl.name, c.nums)
					return
				}
				if !slices.Equal(got, c.want) {
					t.Errorf("%s(%v) = %v, want %v", impl.name, c.nums, got, c.want)
				}
			})
		}
	}
	rng := rand.New(rand.NewSource(72))
	// 随机对拍
	brute := func(nums []int) []int {
		n := len(nums)
		res := make([]int, n)
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if nums[j] > nums[i] {
					res[i] = j - i
					break
				}
			}
		}
		return res
	}
	for iter := 0; iter < 5000; iter++ {
		n := 1 + rng.Intn(30)
		nums := make([]int, n)
		for i := range nums {
			nums[i] = 30 + rng.Intn(71)
		}
		want := brute(nums)
		got := lc.DailyTemperatures(slices.Clone(nums))
		if !slices.Equal(got, want) {
			t.Fatalf("iter=%d nums=%v: got=%v want=%v", iter, nums, got, want)
		}
	}
}

// 225. 用队列实现栈
func TestMyStack(t *testing.T) {
	s := lc.NewMyStack()
	if !s.Empty() {
		t.Error("新栈应为空")
	}
	s.Push(1)
	s.Push(2)
	if got := s.Top(); got != 2 {
		t.Errorf("Top = %d, want 2", got)
	}
	if got := s.Pop(); got != 2 {
		t.Errorf("Pop = %d, want 2", got)
	}
	if s.Empty() {
		t.Error("还有一个元素")
	}
	if got := s.Pop(); got != 1 {
		t.Errorf("Pop = %d, want 1", got)
	}
	if !s.Empty() {
		t.Error("空栈")
	}
	s.Push(3)
	if got := s.Top(); got != 3 {
		t.Errorf("Top = %d, want 3", got)
	}
}

// 232. 用栈实现队列
func TestMyQueue(t *testing.T) {
	q := lc.NewMyQueue()
	if !q.Empty() {
		t.Error("新队列应为空")
	}
	q.Push(1)
	q.Push(2)
	if got := q.Peek(); got != 1 {
		t.Errorf("Peek = %d, want 1", got)
	}
	if got := q.Pop(); got != 1 {
		t.Errorf("Pop = %d, want 1", got)
	}
	if q.Empty() {
		t.Error("还有一个元素")
	}
	q.Push(3)
	if got := q.Pop(); got != 2 {
		t.Errorf("Pop = %d, want 2", got)
	}
	if got := q.Pop(); got != 3 {
		t.Errorf("Pop = %d, want 3", got)
	}
	if !q.Empty() {
		t.Error("空队列")
	}
}

// 622. 设计循环队列
func TestMyCircularQueue(t *testing.T) {
	q := lc.NewMyCircularQueue(3)
	if !q.Empty() {
		t.Error("新队列应为空")
	}
	if !q.EnQueue(1) {
		t.Error("EnQueue(1) 应成功")
	}
	if !q.EnQueue(2) {
		t.Error("EnQueue(2) 应成功")
	}
	if !q.EnQueue(3) {
		t.Error("EnQueue(3) 应成功")
	}
	if q.EnQueue(4) {
		t.Error("EnQueue(4) 满了应失败")
	}
	if got := q.Head(); got != 1 {
		t.Errorf("Head = %d, want 1", got)
	}
	if got := q.Tail(); got != 3 {
		t.Errorf("Tail = %d, want 3", got)
	}
	if !q.DeQueue() {
		t.Error("DeQueue 应成功")
	}
	if !q.EnQueue(4) {
		t.Error("EnQueue(4) 空出一个位置应成功")
	}
	if got := q.Head(); got != 2 {
		t.Errorf("Head = %d, want 2", got)
	}
	if got := q.Tail(); got != 4 {
		t.Errorf("Tail = %d, want 4", got)
	}
	// 清空再验证
	q.DeQueue()
	q.DeQueue()
	q.DeQueue()
	if !q.Empty() {
		t.Error("全部出队后应为空")
	}
	if q.DeQueue() {
		t.Error("空队列 DeQueue 应失败")
	}
}

// ============ 随机对拍用的共享 rng ============
var _ = rand.New(rand.NewSource(99))
