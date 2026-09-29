package lc

/*
给你一个链表，删除链表的倒数第 n* *个结点，并且返回链表的头结点。

示例 1：

输入：head = [1,2,3,4,5], n = 2
输出：[1,2,3,5]

示例 2：

输入：head = [1], n = 1
输出：[]

示例 3：

输入：head = [1,2], n = 1
输出：[1]

提示：

- 链表中结点的数目为 sz

- 1 <= sz <= 30

- 0 <= Node.val <= 100

- 1 <= n <= sz

进阶：你能尝试使用一趟扫描实现吗？
*/

func removeNthFromEnd(head *ListNode, n int) *ListNode {
	node,l := head,0
	for ;node != nil; node,l = node.Next, l+1 {
	}
	node = head
	for i := 0; i < l - n - 1; i++ {
		node = node.Next
	}
	node.Next = node.Next.Next //删除后边的这个节点
	if l == n {
		return head.Next
	}
	//需要特殊处理一下头部
	return head
}
