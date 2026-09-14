package algo

type Heap[T any] struct {
	data []T
	less func(a, b T) bool
}

// NewHeap 创建空堆。
func NewHeap[T any](less func(a, b T) bool) *Heap[T] {
	return &Heap[T]{
		less: less,
	}
}

// NewHeapFrom Floyd 算法
func NewHeapFrom[T any](s []T, less func(a, b T) bool) *Heap[T] {
	h := &Heap[T]{data: s, less: less}
	for i := len(s)/2 - 1; i >= 0; i-- {
		h.down(i)
	}
	return h
}

// Len 返回元素个数。
func (h *Heap[T]) Len() int {
	return len(h.data)
}

func (h *Heap[T]) Peek() T {
	return h.data[0]
}

func (h *Heap[T]) Push(x T) {
	h.data = append(h.data, x)
	h.up(len(h.data) - 1)
}

func (h *Heap[T]) Pop() T {
	min := h.data[0]
	n := len(h.data)
	h.data[0] = h.data[n-1] //末尾元素提升至首位
	var zero T
	h.data[n-1] = zero //release memory
	h.data = h.data[:n-1]
	if len(h.data) > 0 {
		h.down(0)
	}
	return min
}

func (h *Heap[T]) At(i int) T { return h.data[i] }

func (h *Heap[T]) Update(i int, v T) {
	h.data[i] = v
	h.Fix(i)
}

func (h *Heap[T]) Fix(i int) {
	if !h.down(i) {
		h.up(i)
	}
}

func (h *Heap[T]) Remove(i int) T {
	x := h.data[i]
	n := len(h.data)
	if n-1 != i {
		h.data[i] = h.data[n-1]
	}
	var zero T
	h.data[n-1] = zero
	h.data = h.data[:n-1]
	if n-1 != i {
		h.Fix(i)
	}
	return x
}

// up 让第 i 个元素向上滤到合适位置：与父节点比，优先级更高则交换。
func (h *Heap[T]) up(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if !h.less(h.data[i], h.data[parent]) { //不满足交互的条件
			break
		}
		h.data[i], h.data[parent] = h.data[parent], h.data[i]
		i = parent
	}
}

// down 让第 i 个元素向下滤到合适位置：与较小的孩子比，返回是否发生过下沉。
func (h *Heap[T]) down(i int) bool {
	n, origin := len(h.data), i
	for {
		small := i
		if l := 2*i + 1; l < n && h.less(h.data[l], h.data[small]) {
			small = l
		}
		if r := 2*i + 2; r < n && h.less(h.data[r], h.data[small]) {
			small = r
		}
		if small == i {
			return i > origin
		}
		h.data[i], h.data[small] = h.data[small], h.data[i]
		i = small
	}
}
