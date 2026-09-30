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
//6,2
//dummy,1,2,3,4,5

func RemoveNthFromEndTwoPass(head *ListNode, n int) *ListNode {
	dummy, l := &ListNode{
		Next: head,
	}, 0
	node := dummy
	for ; node != nil; node, l = node.Next, l+1 {
	}
	node = dummy
	//为什么这里不用变?
	for i := 0; i < l-n-1; i++ {
		node = node.Next
	}
	node.Next = node.Next.Next //删除后边的这个节点
	return dummy.Next
}

func RemoveNthFromEnd(head *ListNode, n int) *ListNode {
	dummy := &ListNode{
		Next: head,
	}
	fast, slow := dummy, dummy
	//关键的引入dummy节点(如果不引入dummy节点,走n步)
	for range n + 1 {
		fast = fast.Next
	}
	for ; fast != nil; fast, slow = fast.Next, slow.Next {
	}
	//不引入dummy节点,怎么判断是删除头部节点的形式?
	//可以用slow == head 来判断
	//因为引入了dummy节点,原来的head节点就是普通的内部节点,之前我们讨论过,内部节点的删除是不需要有任何特殊处理的
	slow.Next = slow.Next.Next //删除
	return dummy.Next
}
