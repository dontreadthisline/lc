// package algo some algo for
package algo

import "container/heap"

type InfiniteNumSet struct {
	h     *IntHeap
	m     map[int]bool
	upper int
}

func NewInfiniteNumSet() *InfiniteNumSet {
	return &InfiniteNumSet{
		h:     &IntHeap{},
		m:     make(map[int]bool, 0),
		upper: 1,
	}
}

func (i *InfiniteNumSet) PopSmallests() int {
	if i.h.Len() > 0 {
		x := heap.Pop(i.h).(int)
		delete(i.m, x)
		return x
	}
	x := i.upper
	i.upper += 1
	return x
}

func (i *InfiniteNumSet) AddBack(num int) {
	if num >= i.upper {
		return
	}

	if ok := i.m[num]; !ok {
		heap.Push(i.h, num)
		i.m[num] = true
	}
}

type IntHeap []int

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h IntHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h *IntHeap) Pop() any {
	n := len(*h)
	x := (*h)[n-1]
	*h = (*h)[:n-1]
	return x
}

func (h *IntHeap) Push(x any) {
	*h = append(*h, x.(int))
}
