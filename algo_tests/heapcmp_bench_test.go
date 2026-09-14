package algo_tests

import "demo/algo"

import (
	"cmp"
	"container/heap"
	"math/rand"
	"slices"
	"testing"

	jbaheap "github.com/jba/heap"
)

// ── 旧版 container/heap 的三件套仪式 ──
type oldIntHeap []int

func (h oldIntHeap) Len() int           { return len(h) }
func (h oldIntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h oldIntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *oldIntHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *oldIntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func heapBenchData() []int {
	r := rand.New(rand.NewSource(7))
	s := make([]int, 1000)
	for i := range s {
		s[i] = r.Int()
	}
	return s
}

// 正确性交叉验证：三方弹出顺序都应等于排序结果
func TestJbaHeapCrossCheck(t *testing.T) {
	data := heapBenchData()
	want := slices.Clone(data)
	slices.Sort(want)

	h1 := algo.NewHeap(func(a, b int) bool { return a < b })
	for _, v := range data {
		h1.Push(v)
	}
	got1 := make([]int, 0, len(data))
	for h1.Len() > 0 {
		got1 = append(got1, h1.Pop())
	}

	h2 := jbaheap.New(cmp.Compare[int])
	h2.Init(slices.Clone(data))
	got2 := make([]int, 0, len(data))
	for v := range h2.Drain() {
		got2 = append(got2, v)
	}

	if !slices.Equal(got1, want) || !slices.Equal(got2, want) {
		t.Error("弹出顺序与排序结果不一致")
	}
}

func BenchmarkStdlibContainerHeap(b *testing.B) {
	data := heapBenchData()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		h := &oldIntHeap{}
		for _, v := range data {
			heap.Push(h, v)
		}
		for h.Len() > 0 {
			_ = heap.Pop(h).(int)
		}
	}
}

func BenchmarkOursGenericHeap(b *testing.B) {
	data := heapBenchData()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		h := algo.NewHeap(func(a, b int) bool { return a < b })
		for _, v := range data {
			h.Push(v)
		}
		for h.Len() > 0 {
			_ = h.Pop()
		}
	}
}

func BenchmarkJbaHeap(b *testing.B) {
	data := heapBenchData()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		h := jbaheap.New(cmp.Compare[int])
		for _, v := range data {
			h.Insert(v)
		}
		for h.Len() > 0 {
			_ = h.TakeMin()
		}
	}
}

// 对照组：jba 堆 + 手写减法三路比较（拆开"实现差异"与"cmp.Compare 开销"）
// 注：减法比较仅对无溢出范围安全，这里 rand.Int() 非负可用
func BenchmarkJbaHeapSubCmp(b *testing.B) {
	data := heapBenchData()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		h := jbaheap.New(func(a, b int) int { return a - b })
		for _, v := range data {
			h.Insert(v)
		}
		for h.Len() > 0 {
			_ = h.TakeMin()
		}
	}
}

// 对照组：jba 堆 + 分支式手写三路比较（安全且无 cmp.Compare 的两次比较）
func BenchmarkJbaHeapBranchCmp(b *testing.B) {
	data := heapBenchData()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		h := jbaheap.New(func(a, b int) int {
			if a < b {
				return -1
			} else if a > b {
				return 1
			}
			return 0
		})
		for _, v := range data {
			h.Insert(v)
		}
		for h.Len() > 0 {
			_ = h.TakeMin()
		}
	}
}
