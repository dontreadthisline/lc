package lc

import (
	"fmt"
	"strings"
)

// ListNode 链表节点
type ListNode struct {
	Val  int
	Next *ListNode
}

func (l *ListNode) String() string {
	if l == nil {
		return ""
	}
	buf := strings.Builder{}
	for ; l != nil; l = l.Next {
		if l.Next == nil {
			buf.WriteString(fmt.Sprintf("%d", l.Val))
		} else {
			buf.WriteString(fmt.Sprintf("%d->", l.Val))
		}
	}
	return buf.String()
}

func NewLinkList(vals []int) *ListNode {
	if len(vals) <= 0 {
		return nil
	}
	dummy := ListNode{}
	tmp := &dummy
	for _, val := range vals {
		node := ListNode{
			Val: val,
		}
		tmp.Next = &node
		tmp = tmp.Next
	}

	return dummy.Next
}

// TreeNode 二叉树节点
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

type Stack[T any] struct{
	items []T
}

func NewStack[T any]() *Stack[T] {

	return &Stack[T]{
		items:make([]T,0),
	}
}

func (s Stack[T]) Len() int {
	return len(s.items)
}

func (s Stack[T]) Empty() bool {
	return len(s.items) == 0
}

func (s Stack[T]) Peek() T {
	return s.items[s.Len()-1]
}

func (s *Stack[T]) Pop() T {
	item := s.Peek()
	s.items = s.items[0:s.Len()-1]
	return item
}

func (s *Stack[T]) Push(item T)  {
	s.items = append(s.items,item)
}

type node[T any] struct {
	item T
	next *node[T]
}

type Queue[T any] struct{
	head *node[T]
	tail *node[T]
	n int
}

func NewQueue[T any]() *Queue[T]{
	return &Queue[T]{
		head:nil,
		tail:nil,
		n:0,
	}
}

func (q Queue[T]) Len() int {
	return q.n
}

func (q Queue[T]) Empty() bool {
	return q.n == 0
}

func (q Queue[T]) Head() T {
	return q.head.item
}

func (q Queue[T]) Tail() T {
	return q.tail.item
}

func (q *Queue[T]) Enque(item T) {
	if q.n == 0 {
		q.head = &node[T]{
			item:item,
		}
		q.tail = q.head
	} else {
		q.tail.next = &node[T]{
			item:item,
		}
		q.tail = q.tail.next
	}
	q.n += 1
}

func (q *Queue[T]) Deque() T {
	peek := q.Head()
	head := q.head
	q.head = q.head.next
	//memory free?
	head.next = nil
	head = nil
	q.n -= 1
	return peek
}

