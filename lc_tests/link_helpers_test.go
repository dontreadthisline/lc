package lc_tests

import "demo/lc"

// 家族 6 链表测试的公共辅助

func listToSlice(head *lc.ListNode) []int {
	out := []int{}
	for ; head != nil; head = head.Next {
		out = append(out, head.Val)
		if len(out) > 2048 {
			break // 防御：实现若把链表改出环，避免死循环
		}
	}
	return out
}

// safeCall 捕获 panic 转为返回值，保证单个用例的崩溃不炸掉整个测试进程
func safeCall(f func()) (panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
		}
	}()
	f()
	return false
}

// buildCycle 构造带环链表：尾节点指向下标 pos 的节点；pos < 0 表示无环
func buildCycle(vals []int, pos int) *lc.ListNode {
	head := lc.NewLinkList(vals)
	if pos < 0 || head == nil {
		return head
	}
	tail := head
	for tail.Next != nil {
		tail = tail.Next
	}
	p := head
	for i := 0; i < pos; i++ {
		p = p.Next
	}
	tail.Next = p
	return head
}

// buildIntersect 构造相交链表：A = aVals + common，B = bVals + common，公共段节点共享
func buildIntersect(aVals, bVals, common []int) (*lc.ListNode, *lc.ListNode, *lc.ListNode) {
	commonHead := lc.NewLinkList(common)
	attach := func(prefix []int) *lc.ListNode {
		h := lc.NewLinkList(prefix)
		if h == nil {
			return commonHead
		}
		tail := h
		for tail.Next != nil {
			tail = tail.Next
		}
		tail.Next = commonHead
		return h
	}
	return attach(aVals), attach(bVals), commonHead
}
