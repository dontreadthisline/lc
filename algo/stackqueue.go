package algo

// Stack 是泛型 LIFO 栈，slice 实现。
type Stack[T any] struct {
	data []T
}

func NewStack[T any](items ...T) *Stack[T] {
	s := Stack[T]{
		data: make([]T, 0, len(items)), // 容量一次给足，Push 零搬运
	}
	for _, item := range items {
		s.Push(item)
	}
	return &s
}

func (s Stack[T]) Len() int { return len(s.data) } // 值接收者：只读，自表达不可变

func (s Stack[T]) Peek() T {
	return s.data[len(s.data)-1]
}

func (s *Stack[T]) Push(item T) {
	s.data = append(s.data, item)
}

func (s *Stack[T]) Pop() T {
	n := len(s.data)
	top := s.data[n-1]
	var zero T
	s.data[n-1] = zero
	s.data = s.data[0 : n-1]
	return top
}

func (s *Stack[T]) Clear() {
	clear(s.data)
	s.data = s.data[:0]
}

type Queue[T any] struct {
	data []T
	head int
}

// NewQueue 创建空队列。
func NewQueue[T any](items ...T) *Queue[T] {
	q := &Queue[T]{data: make([]T, 0, len(items))} // 对齐 Stack：容量预置，Push 零搬运
	for _, item := range items {
		q.Push(item)
	}
	return q
}

func (q Queue[T]) Len() int {
	return len(q.data) - q.head
}

func (q *Queue[T]) Push(x T) { q.data = append(q.data, x) }

// Peek 返回队头，不出队。空队列 panic。
func (q *Queue[T]) Peek() T {
	return q.data[q.head]
}

// Pop 出队并返回队头。空队列 panic。槽位置零 + 惰性压实，见类型注释。
func (q *Queue[T]) Pop() T {
	v := q.data[q.head]
	var zero T
	q.data[q.head] = zero
	q.head++
	if q.head*2 >= len(q.data) { // 死区过半：压实
		n := copy(q.data, q.data[q.head:])
		for i := n; i < len(q.data); i++ {
			q.data[i] = zero
		}
		q.data = q.data[:n]
		q.head = 0
	}
	return v
}

func (q *Queue[T]) Clear() {
	clear(q.data)
	q.data = q.data[:0]
	q.head = 0
}
