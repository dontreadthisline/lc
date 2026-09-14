package algo_tests

import "demo/algo"

import (
	"container/heap"
	"testing"
)

func TestInfiniteNumSet(t *testing.T) {
	s := algo.NewInfiniteNumSet()

	// 正常流：pop 出 1,2,3
	t.Log(s.PopSmallests()) // 期望 1
	t.Log(s.PopSmallests()) // 期望 2

	// 放回 1 和 2
	s.AddBack(1)
	s.AddBack(2)
	s.AddBack(2) // 重复 add，靠 m 去重

	// 堆里现在有 [1,2]，再 pop 应该依次是 1, 2, 3...
	for range 4 {
		t.Log(s.PopSmallests())
	}

	// 直接戳 heap 的 bug：堆里只放 1 个元素时 Pop
	h := &algo.IntHeap{5}
	heap.Init(h)
	t.Log("about to heap.Pop on single-element heap...")
	t.Log(heap.Pop(h)) // 期望 5，实际会 panic 或返回错值
}
