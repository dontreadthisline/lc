package lc

/*
给你一个单链表的头节点 head ，请你判断该链表是否为回文链表。如果是，返回 true ；否则，返回 false 。

示例 1：

输入：head = [1,2,2,1]
输出：true

示例 2：

输入：head = [1,2]
输出：false

提示：

- 链表中节点数目在范围[1, 10^5] 内

- 0 <= Node.val <= 9

进阶：你能否用 O(n) 时间复杂度和 O(1) 空间复杂度解决此题？
*/

func IsPalindrome(head *ListNode) bool {
	mid := MiddleOfLinkList(head)
	mid = ReverseList(mid)
	for ; mid != nil && head != nil; mid = mid.Next {
		if mid.Val != head.Val {
			return false
		}
		head = head.Next
	}
	return true
}

//补充一道链表找中点的题目
//双指针一遍的方法

func MiddleOfLinkList(head *ListNode) *ListNode {
	fast, slow := head, head
	for ; fast != nil; slow = slow.Next {
		if fast.Next == nil {
			return slow
		}
		fast = fast.Next.Next
	}
	return slow
}
