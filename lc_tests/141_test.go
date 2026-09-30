package lc_tests

import (
	"demo/lc"
	"testing"
)

// 141. 环形链表
func TestHasCycle(t *testing.T) {
	cases := []struct {
		name string
		vals []int
		pos  int
		want bool
	}{
		{"官方用例1_尾连下标1", []int{3, 2, 0, -4}, 1, true},
		{"官方用例2_尾连头", []int{1, 2}, 0, true},
		{"官方用例3_单节点无环", []int{1}, -1, false},
		{"无环四节点", []int{1, 2, 3, 4}, -1, false},
		{"空表", nil, -1, false},
		{"单节点自环", []int{1}, 0, true},
		{"尾节点自环", []int{1, 2, 3}, 2, true},
	}
	impls := []struct {
		name string
		fn   func(*lc.ListNode) bool
	}{
		{"HasCycle", lc.HasCycle},
		{"HasCycleHashWay", lc.HasCycleHashWay},
	}
	for _, impl := range impls {
		for _, c := range cases {
			t.Run(impl.name+"/"+c.name, func(t *testing.T) {
				head := buildCycle(c.vals, c.pos)
				var got bool
				panicked := safeCall(func() { got = impl.fn(head) })
				if panicked {
					t.Errorf("%s(%v, pos=%d) panic", impl.name, c.vals, c.pos)
					return
				}
				if got != c.want {
					t.Errorf("%s(%v, pos=%d) = %v, want %v", impl.name, c.vals, c.pos, got, c.want)
				}
			})
		}
	}
}
