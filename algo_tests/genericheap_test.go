package algo_tests

import "demo/algo"

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"
)

type task struct {
	priority int
	name     string
}

func byPriority(a, b task) bool { return a.priority < b.priority }

// drain 弹空堆，返回弹出顺序。
func drain[T any](h *algo.Heap[T]) []T {
	out := make([]T, 0, h.Len())
	for h.Len() > 0 {
		out = append(out, h.Pop())
	}
	return out
}

func TestGenericHeapOrder(t *testing.T) {
	h := algo.NewHeap(func(a, b int) bool { return a < b })
	for _, v := range []int{5, 3, 8, 1, 9, 2} {
		h.Push(v)
	}
	if got := drain(h); !slices.Equal(got, []int{1, 2, 3, 5, 8, 9}) {
		t.Fatalf("小顶堆弹出序 = %v, 期望 [1 2 3 5 8 9]", got)
	}
}

func TestGenericHeapMaxOrder(t *testing.T) {
	h := algo.NewHeap(func(a, b int) bool { return a > b })
	for _, v := range []int{5, 3, 8, 1} {
		h.Push(v)
	}
	if got := drain(h); !slices.Equal(got, []int{8, 5, 3, 1}) {
		t.Fatalf("大顶堆弹出序 = %v, 期望 [8 5 3 1]", got)
	}
}

func TestGenericHeapStruct(t *testing.T) {
	h := algo.NewHeap(byPriority)
	h.Push(task{3, "c"})
	h.Push(task{1, "a"})
	h.Push(task{2, "b"})
	got := fmt.Sprint(drain(h))
	if want := "[{1 a} {2 b} {3 c}]"; got != want {
		t.Fatalf("按 priority 弹出 = %s, 期望 %s", got, want)
	}
}

func TestGenericHeapFromSlice(t *testing.T) {
	vals := []int{9, 2, 7, 4, 1, 5}
	want := slices.Clone(vals)
	slices.Sort(want)

	h := algo.NewHeapFrom(vals, func(a, b int) bool { return a < b })
	if got := drain(h); !slices.Equal(got, want) {
		t.Fatalf("Floyd 建堆弹出序 = %v, 期望 %v", got, want)
	}
}

func TestGenericHeapBasics(t *testing.T) {
	h := algo.NewHeap(func(a, b int) bool { return a < b })
	if h.Len() != 0 {
		t.Fatal("空堆初始状态错误")
	}
	h.Push(7)
	h.Push(3)
	if h.Peek() != 3 || h.Len() != 2 {
		t.Fatalf("Peek=%d Len=%d, 期望 3/2", h.Peek(), h.Len())
	}
}

// Fix 两个方向都要覆盖：优先级调高要上滤，调低要下滤。
// 依次 Push 1,5,8,9,6 后内部布局是 [1 5 8 9 6]（合法小顶堆），
// 测试里的下标都基于这个布局。
func TestGenericHeapFix(t *testing.T) {
	newHeap := func() *algo.Heap[task] {
		h := algo.NewHeap(byPriority)
		for _, p := range []int{1, 5, 8, 9, 6} {
			h.Push(task{p, fmt.Sprintf("t%d", p)})
		}
		return h
	}

	// 调低：idx1 的 5 改成 7，比孩子 6 大，需下滤到 [1 6 8 9 7]
	h := newHeap()
	h.Update(1, task{7, "t5->7"})
	got := fmt.Sprint(drain(h))
	if want := "[{1 t1} {6 t6} {7 t5->7} {8 t8} {9 t9}]"; got != want {
		t.Fatalf("下调后弹出序 = %s, 期望 %s", got, want)
	}

	// 调高：idx4 的 6 改成 0，比父级都小，需上滤到 [0 1 8 9 5]
	h = newHeap()
	h.Update(4, task{0, "t6->0"})
	got = fmt.Sprint(drain(h))
	if want := "[{0 t6->0} {1 t1} {5 t5} {8 t8} {9 t9}]"; got != want {
		t.Fatalf("上调后弹出序 = %s, 期望 %s", got, want)
	}
}

func TestGenericHeapRemove(t *testing.T) {
	h := algo.NewHeap(byPriority)
	for _, p := range []int{1, 5, 8, 9, 6} {
		h.Push(task{p, fmt.Sprintf("t%d", p)})
	} // 布局 [1 5 8 9 6]

	// 摘中间：末尾 t6 补位到 idx2，无需移动（介于父 1 和无孩子之间）
	if removed := h.Remove(2); removed.name != "t8" {
		t.Fatalf("Remove(2) 返回 %v, 期望 t8", removed)
	}
	got := fmt.Sprint(drain(h))
	if want := "[{1 t1} {5 t5} {6 t6} {9 t9}]"; got != want {
		t.Fatalf("删除后弹出序 = %s, 期望 %s", got, want)
	}

	// 删堆顶等价于 Pop
	h2 := algo.NewHeap(func(a, b int) bool { return a < b })
	for _, v := range []int{4, 2, 6} {
		h2.Push(v)
	}
	if v := h2.Remove(0); v != 2 {
		t.Fatalf("Remove(0) = %d, 期望 2", v)
	}
	if got := drain(h2); !slices.Equal(got, []int{4, 6}) {
		t.Fatalf("删堆顶后弹出序 = %v", got)
	}

	// 删最后一个元素：无补位、无重排
	h3 := algo.NewHeap(func(a, b int) bool { return a < b })
	for _, v := range []int{3, 1} {
		h3.Push(v)
	}
	if v := h3.Remove(h3.Len() - 1); v != 3 {
		t.Fatalf("Remove(last) = %d, 期望 3", v)
	}
	if got := drain(h3); !slices.Equal(got, []int{1}) {
		t.Fatalf("删尾后弹出序 = %v", got)
	}
}

// Fix 的另一半用途：存指针时，元素持有者绕过堆直接改字段，改完只需重排。
// 这也是 Fix 不能收进私有的原因——值语义走 Update，指针语义走 Fix。
func TestGenericHeapFixPointer(t *testing.T) {
	type item struct {
		priority int
		name     string
	}
	h := algo.NewHeap(func(a, b *item) bool { return a.priority < b.priority })
	a := &item{1, "a"}
	b := &item{5, "b"}
	c := &item{8, "c"}
	for _, it := range []*item{a, b, c} {
		h.Push(it)
	}

	// b 还在堆里排队，业务侧直接降它的优先级到 0
	b.priority = 0
	// 找到 b 的下标（此处线性扫；生产上由元素自带 index 字段跟踪，见 go/types initorder.go）
	for i := range h.Len() {
		if h.At(i) == b {
			h.Fix(i)
			break
		}
	}

	if top := h.Pop(); top != b {
		t.Fatalf("改优先级后堆顶 = %v, 期望 b", top)
	}
	if second := h.Pop(); second != a {
		t.Fatalf("第二个弹出 = %v, 期望 a", second)
	}
}

// 随机对照：任意一批数据，堆的弹出序必须与全序排序一致。
// 这是堆的核心不变量，比逐个手算用例更能兜住 up/down 的边界错误。
func TestGenericHeapRandomCrossCheck(t *testing.T) {
	rnd := rand.New(rand.NewPCG(1, 2))
	for round := range 100 {
		vals := make([]int, rnd.IntN(200))
		for i := range vals {
			vals[i] = rnd.IntN(1000)
		}
		want := slices.Clone(vals)
		slices.Sort(want)

		h := algo.NewHeapFrom(vals, func(a, b int) bool { return a < b })
		if got := drain(h); !slices.Equal(got, want) {
			t.Fatalf("round %d: 堆弹出序 != 排序序\n got  %v\n want %v", round, got, want)
		}
	}
}
