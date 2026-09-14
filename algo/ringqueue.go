package algo

// RingQueue 是泛型环形 FIFO：head/count + 模运算下标，出队空间立即可复用，
// 没有压实的摊还拷贝。对照 head+压实的 Queue：
//
//   - 优点：零压实拷贝、容量固定时零再分配（channel 内部即此结构）
//   - 代价：扩容要手工展开两段 wrap 数据；活区不连续，不能当 slice 直接用
//
// GC 契约不变：环形不免除槽位置零——出队的槽位同样要写零值，
// 否则残留引用照样钉住对象。
type RingQueue[T any] struct {
	data  []T
	head  int // 队头下标（模 len(data)）
	count int // 活元素个数；满/空判别靠它，不牺牲格子
}

// NewRingQueue 创建容量为 cap 的环形队列（cap < 1 视为 1）。
func NewRingQueue[T any](capacity int) *RingQueue[T] {
	if capacity < 1 {
		capacity = 1
	}
	return &RingQueue[T]{data: make([]T, capacity)}
}

func (q *RingQueue[T]) Len() int      { return q.count }
func (q *RingQueue[T]) IsEmpty() bool { return q.count == 0 }

// Push 入队，摊还 O(1)。满了倍增扩容：wrap 成两段的数据按模展开搬到新数组。
func (q *RingQueue[T]) Push(x T) {
	if q.count == len(q.data) {
		grown := make([]T, len(q.data)*2)
		for i := 0; i < q.count; i++ {
			grown[i] = q.data[(q.head+i)%len(q.data)]
		}
		q.data = grown
		q.head = 0
	}
	q.data[(q.head+q.count)%len(q.data)] = x
	q.count++
}

// Peek 返回队头。空队列 panic。
func (q *RingQueue[T]) Peek() T {
	if q.count == 0 {
		panic("algo: RingQueue.Peek on empty queue")
	}
	return q.data[q.head]
}

// Pop 出队并返回队头。空队列 panic。槽位置零（GC 契约）。
func (q *RingQueue[T]) Pop() T {
	if q.count == 0 {
		panic("algo: RingQueue.Pop on empty queue")
	}
	v := q.data[q.head]
	var zero T
	q.data[q.head] = zero
	q.head = (q.head + 1) % len(q.data)
	q.count--
	if q.count == 0 {
		q.head = 0 // 排空归零，下轮从数组头开始，保持简单
	}
	return v
}

// Clear 清空并释放全部引用，保留容量。
func (q *RingQueue[T]) Clear() {
	clear(q.data)
	q.head = 0
	q.count = 0
}
