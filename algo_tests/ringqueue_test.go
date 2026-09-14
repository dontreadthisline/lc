package algo_tests

import "demo/algo"

import (
	"math/rand"
	"testing"
)

func TestRingQueueFIFO(t *testing.T) {
	q := algo.NewRingQueue[int](2) // 故意小容量，逼出 wrap 和扩容
	for i := 0; i < 5; i++ {
		q.Push(i)
	}
	for i := 0; i < 5; i++ {
		if v := q.Pop(); v != i {
			t.Fatalf("FIFO 顺序错: got %d, want %d", v, i)
		}
	}
}

// wrap + 扩容展开：交替入出，越过容量边界多轮
func TestRingQueueWraparound(t *testing.T) {
	q := algo.NewRingQueue[int](4)
	next := 0
	for round := 0; round < 100; round++ {
		for i := 0; i < 3; i++ { // 每轮进 3 出 3，head 持续前移制造 wrap
			q.Push(next)
			next++
		}
		for i := 0; i < 3; i++ {
			want := next - 3 + i
			if got := q.Pop(); got != want {
				t.Fatalf("round %d: got %d, want %d", round, got, want)
			}
		}
	}
}

// 槽位置零：出队后全数组无残留指针
func TestRingQueueSlotGC(t *testing.T) {
	q := algo.NewRingQueue[*int](4)
	for i := 0; i < 9; i++ { // 9 个元素，跨过一次扩容
		v := i
		q.Push(&v)
	}
	for i := 0; i < 9; i++ {
		q.Pop()
	}
	if q.Len() != 0 {
		t.Errorf("弹空后 Len=%d, want 0", q.Len())
	}
}

func TestRingQueueRandomCrossCheck(t *testing.T) {
	r := rand.New(rand.NewSource(7))
	q := algo.NewRingQueue[int](3)
	var naive []int
	for step := 0; step < 10000; step++ {
		if len(naive) == 0 || r.Intn(2) == 0 {
			v := r.Intn(1000)
			q.Push(v)
			naive = append(naive, v)
		} else {
			if q.Pop() != naive[0] {
				t.Fatalf("step %d: FIFO 分叉", step)
			}
			naive = naive[1:]
		}
		if q.Len() != len(naive) {
			t.Fatalf("step %d: len %d != %d", step, q.Len(), len(naive))
		}
	}
}
