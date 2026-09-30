package lc_tests

import (
	"demo/lc"
	"slices"
	"testing"
)

// 19. 删除链表的倒数第 N 个结点
func TestRemoveNthFromEnd(t *testing.T) {
	cases := []struct {
		name string
		vals []int
		n    int
		want []int
	}{
		{"官方用例_删中间", []int{1, 2, 3, 4, 5}, 2, []int{1, 2, 3, 5}},
		{"官方用例_单节点删自身", []int{1}, 1, []int{}},
		{"官方用例_两节点删尾", []int{1, 2}, 1, []int{1}},
		{"两节点删头", []int{1, 2}, 2, []int{2}},
		{"五节点删头", []int{1, 2, 3, 4, 5}, 5, []int{2, 3, 4, 5}},
		{"三节点删头", []int{1, 2, 3}, 3, []int{2, 3}},
		{"删正数第二个", []int{1, 2, 3, 4}, 3, []int{1, 3, 4}},
	}
	impls := []struct {
		name string
		fn   func(*lc.ListNode, int) *lc.ListNode
	}{
		{"RemoveNthFromEnd", lc.RemoveNthFromEnd},
		{"RemoveNthFromEndTwoPass", lc.RemoveNthFromEndTwoPass},
	}
	for _, impl := range impls {
		for _, c := range cases {
			t.Run(impl.name+"/"+c.name, func(t *testing.T) {
				var got []int
				panicked := safeCall(func() {
					got = listToSlice(impl.fn(lc.NewLinkList(c.vals), c.n))
				})
				if panicked {
					t.Errorf("%s(%v, n=%d) panic", impl.name, c.vals, c.n)
					return
				}
				if !slices.Equal(got, c.want) {
					t.Errorf("%s(%v, n=%d) = %v, want %v", impl.name, c.vals, c.n, got, c.want)
				}
			})
		}
	}
}

// 21. 合并两个有序链表
func TestMergeTwoLists(t *testing.T) {
	cases := []struct {
		name             string
		a, b, want       []int
	}{
		{"官方用例1", []int{1, 2, 4}, []int{1, 3, 4}, []int{1, 1, 2, 3, 4, 4}},
		{"官方用例2_均空", nil, nil, []int{}},
		{"官方用例3_一空", nil, []int{0}, []int{0}},
		{"长度不等", []int{1, 3, 5, 7}, []int{2, 4}, []int{1, 2, 3, 4, 5, 7}},
		{"单元素各一", []int{2}, []int{1}, []int{1, 2}},
		{"负数", []int{-3, -1}, []int{-2, 0}, []int{-3, -2, -1, 0}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got []int
			panicked := safeCall(func() {
				got = listToSlice(lc.MergeTwoLists(lc.NewLinkList(c.a), lc.NewLinkList(c.b)))
			})
			if panicked {
				t.Errorf("MergeTwoLists(%v, %v) panic", c.a, c.b)
				return
			}
			if !slices.Equal(got, c.want) {
				t.Errorf("MergeTwoLists(%v, %v) = %v, want %v", c.a, c.b, got, c.want)
			}
		})
	}
}

// 142. 环形链表 II
func TestDetectCycle(t *testing.T) {
	cases := []struct {
		name    string
		vals    []int
		pos     int
		wantNil bool
	}{
		{"官方用例1_尾连下标1", []int{3, 2, 0, -4}, 1, false},
		{"官方用例2_尾连头", []int{1, 2}, 0, false},
		{"官方用例3_无环", []int{1}, -1, true},
		{"无环多节点", []int{1, 2, 3, 4}, -1, true},
		{"空表", nil, -1, true},
		{"尾节点自环", []int{1, 2, 3}, 2, false},
	}
	impls := []struct {
		name string
		fn   func(*lc.ListNode) *lc.ListNode
	}{
		{"DetectCycle", lc.DetectCycle},
		{"DetectCycleHashWay", lc.DetectCycleHashWay},
	}
	for _, impl := range impls {
		for _, c := range cases {
			t.Run(impl.name+"/"+c.name, func(t *testing.T) {
				head := buildCycle(c.vals, c.pos)
				var wantNode *lc.ListNode
				if c.pos >= 0 {
					wantNode = head
					for i := 0; i < c.pos; i++ {
						wantNode = wantNode.Next
					}
				}
				var got *lc.ListNode
				panicked := safeCall(func() { got = impl.fn(head) })
				if panicked {
					t.Errorf("%s(%v, pos=%d) panic", impl.name, c.vals, c.pos)
					return
				}
				if c.wantNil {
					if got != nil {
						t.Errorf("%s(%v, pos=%d) = %v, want nil", impl.name, c.vals, c.pos, got)
					}
					return
				}
				if got != wantNode {
					t.Errorf("%s(%v, pos=%d) 入环点 = %v, want %v", impl.name, c.vals, c.pos, got, wantNode)
				}
			})
		}
	}
}
