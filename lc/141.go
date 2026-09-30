package lc

/*
给你一个链表的头节点 head ，判断链表中是否有环。

如果链表中有某个节点，可以通过连续跟踪 next 指针再次到达，则链表中存在环。 为了表示给定链表中的环，评测系统内部使用整数 pos 来表示链表尾连接到链表中的位置（索引从 0 开始）。注意：pos 不作为参数进行传递 。仅仅是为了标识链表的实际情况。

*如果链表中存在环* ，则返回 true 。 否则，返回 false 。

示例 1：

输入：head = [3,2,0,-4], pos = 1
输出：true
解释：链表中有一个环，其尾部连接到第二个节点。

示例 2：

输入：head = [1,2], pos = 0
输出：true
解释：链表中有一个环，其尾部连接到第一个节点。

示例 3：

输入：head = [1], pos = -1
输出：false
解释：链表中没有环。

提示：

- 链表中节点的数目范围是 [0, 10^4]

- -10^5 <= Node.val <= 10^5

- pos 为 -1 或者链表中的一个 有效索引 。

进阶：你能用 O(1)（即，常量）内存解决此问题吗？
*/

func HasCycle(head *ListNode) bool {
	//如果指示判断有无环,那么快慢指针足够，因为进入圈以后,一定会相遇
	fast, slow := head, head
	//先修panic,先判slow再判fast
	for fast != nil && fast.Next != nil {
		fast = fast.Next.Next
		slow = slow.Next
		//先移动再判相等,避免初始化判断成了head
		if fast == slow {
			return true
		}
	}
	//能走到链表尾部,无环
	return false
}

func HasCycleHashWay(head *ListNode) bool {
	m := make(map[*ListNode]struct{}, 0)

	for ; head != nil; head = head.Next {
		if _, ok := m[head]; ok {
			return true
		}
		m[head] = struct{}{}
	}

	return false
}
