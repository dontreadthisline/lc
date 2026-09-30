package lc_tests

import (
	"demo/lc"
	"testing"
)

// 876. 链表的中间结点（实现位于 lc/234.go）
func TestMiddleOfLinkList(t *testing.T) {
	cases := []struct {
		name    string
		vals    []int
		wantMid int // 期望中点节点的 Val；空表期望 nil
		isEmpty bool
	}{
		{"空表", nil, 0, true},
		{"单节点", []int{1}, 1, false},
		{"两节点_取后者", []int{1, 2}, 2, false},
		{"三节点_取正中", []int{1, 2, 3}, 2, false},
		{"四节点_取后者", []int{1, 2, 3, 4}, 3, false},
		{"五节点_取正中", []int{1, 2, 3, 4, 5}, 3, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got *lc.ListNode
			panicked := safeCall(func() { got = lc.MiddleOfLinkList(lc.NewLinkList(c.vals)) })
			if panicked {
				t.Errorf("MiddleOfLinkList(%v) panic", c.vals)
				return
			}
			if c.isEmpty {
				if got != nil {
					t.Errorf("MiddleOfLinkList(空) = %v, want nil", got)
				}
				return
			}
			if got == nil || got.Val != c.wantMid {
				t.Errorf("MiddleOfLinkList(%v) 中点 = %v, want Val=%d", c.vals, got, c.wantMid)
			}
		})
	}
}

// 234. 回文链表
func TestIsPalindrome(t *testing.T) {
	cases := []struct {
		name string
		vals []int
		want bool
	}{
		{"官方用例1_偶回文", []int{1, 2, 2, 1}, true},
		{"官方用例2_非回文", []int{1, 2}, false},
		{"空表", nil, true},
		{"单节点", []int{1}, true},
		{"奇回文", []int{1, 2, 1}, true},
		{"奇非回文", []int{1, 2, 3}, false},
		{"长奇回文", []int{1, 2, 3, 2, 1}, true},
		{"长偶回文", []int{1, 2, 3, 3, 2, 1}, true},
		{"头尾不同", []int{1, 2, 2, 2}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got bool
			panicked := safeCall(func() { got = lc.IsPalindrome(lc.NewLinkList(c.vals)) })
			if panicked {
				t.Errorf("IsPalindrome(%v) panic", c.vals)
				return
			}
			if got != c.want {
				t.Errorf("IsPalindrome(%v) = %v, want %v", c.vals, got, c.want)
			}
		})
	}
}
