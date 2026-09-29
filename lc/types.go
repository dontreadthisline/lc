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
			buf.WriteString(fmt.Sprintf("%d",l.Val))
		} else {
			buf.WriteString(fmt.Sprintf("%d->",l.Val))
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
	for _,val := range vals {
		node := ListNode{
			Val:val,
		}
		tmp.Next = &node
		tmp = tmp.Next
	}

	return dummy.Next
}

//TreeNode 二叉树节点
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}
