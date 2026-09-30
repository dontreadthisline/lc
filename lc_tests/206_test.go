package lc_tests

import (
	"demo/lc"
	"slices"
	"testing"
)

// 206. 反转链表（递归版 + 迭代版）
func TestReverseList(t *testing.T) {
	cases := []struct {
		name string
		vals []int
		want []int
	}{
		{"官方用例", []int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{"两节点", []int{1, 2}, []int{2, 1}},
		{"单节点", []int{1}, []int{1}},
		{"空表", nil, []int{}},
		{"偶数长", []int{1, 2, 3, 4}, []int{4, 3, 2, 1}},
	}
	impls := []struct {
		name string
		fn   func(*lc.ListNode) *lc.ListNode
	}{
		{"ReverseList", lc.ReverseList},
		{"ReverseListIterWay", lc.ReverseListIterWay},
	}
	for _, impl := range impls {
		for _, c := range cases {
			t.Run(impl.name+"/"+c.name, func(t *testing.T) {
				var got []int
				panicked := safeCall(func() {
					got = listToSlice(impl.fn(lc.NewLinkList(c.vals)))
				})
				if panicked {
					t.Errorf("%s(%v) panic", impl.name, c.vals)
					return
				}
				if !slices.Equal(got, c.want) {
					t.Errorf("%s(%v) = %v, want %v", impl.name, c.vals, got, c.want)
				}
			})
		}
	}
}
