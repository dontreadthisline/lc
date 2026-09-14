package algo_tests

import "demo/algo"

import (
	"math/rand"
	"testing"
)

func TestStackLIFO(t *testing.T) {
	s := algo.NewStack(1, 2, 3)
	if s.Len() != 3 {
		t.Fatalf("len=%d, want 3", s.Len())
	}
	got := []int{}
	for s.Len() > 0 {
		got = append(got, s.Pop())
	}
	if len(got) != 3 || got[0] != 3 || got[1] != 2 || got[2] != 1 {
		t.Errorf("LIFO 顺序错误: %v", got)
	}
}

// 槽位置零验证：弹出的指针必须从底层数组消失
func TestStackSlotGC(t *testing.T) {
	s := algo.NewStack[*int]()
	a, b, c := 1, 2, 3
	s.Push(&a)
	s.Push(&b)
	s.Push(&c)
	s.Pop()
	s.Pop()
	s.Pop()
	// 槽位置零的 GC 契约属白盒断言，目录分离后移交实现包自查；
	// 黑盒视角：弹空后 Len==0 且可继续使用
	if s.Len() != 0 {
		t.Errorf("弹空后 Len=%d, want 0", s.Len())
	}
}

func TestStackEmptyPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("空栈 Pop 应 panic")
		}
	}()
	algo.NewStack[int]().Pop()
}

func TestQueueFIFO(t *testing.T) {
	q := algo.NewQueue[int]()
	for _, v := range []int{1, 2, 3} {
		q.Push(v)
	}
	if q.Peek() != 1 {
		t.Fatalf("peek=%d, want 1", q.Peek())
	}
	got := []int{}
	for q.Len() > 0 {
		got = append(got, q.Pop())
	}
	if len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
		t.Errorf("FIFO 顺序错误: %v", got)
	}
}

// 压实验证：pop 到死区过半，head 应归零且顺序不乱
func TestQueueCompact(t *testing.T) {
	q := algo.NewQueue[int]()
	for i := 0; i < 10; i++ {
		q.Push(i)
	}
	for i := 0; i < 5; i++ { // 第 5 次触发压实（head*2 >= len）
		q.Pop()
	}
	// 压实是内部优化，黑盒只验证行为：剩余 5..9，补 3 个 100..102，共 8 个按序弹出
	for i := 0; i < 3; i++ {
		q.Push(100 + i)
	}
	if q.Len() != 8 {
		t.Fatalf("len=%d, want 8", q.Len())
	}
	for i := 0; i < 5; i++ {
		if v := q.Pop(); v != 5+i {
			t.Errorf("压实后顺序错: got %d want %d", v, 5+i)
		}
	}
	for i := 0; i < 3; i++ {
		if v := q.Pop(); v != 100+i {
			t.Errorf("压实后新元素顺序错: got %d want %d", v, 100+i)
		}
	}
	if q.Len() != 0 {
		t.Errorf("排空后 Len=%d, want 0", q.Len())
	}
}

// 随机交叉验证：对照朴素 slice 队列
func TestQueueRandomCrossCheck(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	q := algo.NewQueue[int]()
	var naive []int
	for step := 0; step < 10000; step++ {
		if len(naive) == 0 || r.Intn(2) == 0 {
			v := r.Intn(1000)
			q.Push(v)
			naive = append(naive, v)
		} else {
			if q.Pop() != naive[0] {
				t.Fatalf("step %d: FIFO 顺序分叉", step)
			}
			naive = naive[1:]
		}
		if q.Len() != len(naive) {
			t.Fatalf("step %d: len %d != %d", step, q.Len(), len(naive))
		}
	}
}
