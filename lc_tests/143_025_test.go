package lc_tests

import (
	"math/rand"
	"testing"

	"demo/lc"
)

// ---- 143. 重排链表 ----
// L0 → Ln → L1 → Ln-1 → L2 → Ln-2 → …

// refReorderSeq: 直接构造期望序列
func refReorderSeq(vals []int) []int {
	n := len(vals)
	out := make([]int, 0, n)
	for i, j := 0, n-1; i <= j; i, j = i+1, j-1 {
		out = append(out, vals[i])
		if i != j {
			out = append(out, vals[j])
		}
	}
	return out
}

func TestReorderList(t *testing.T) {
	cases := [][]int{
		{1, 2, 3, 4},
		{1, 2, 3, 4, 5},
		{1},
		{},
		{1, 2},
	}
	for _, vals := range cases {
		head := lc.NewLinkList(vals)
		runSafely(t, func() { lc.ReorderList(head) })
		got := listToSlice(head)
		want := refReorderSeq(vals)
		if len(got) != len(want) {
			t.Errorf("ReorderList(%v) 长度 = %d, want %d", vals, len(got), len(want))
			continue
		}
		for i := range got {
			if got[i] != want[i] {
				t.Errorf("ReorderList(%v) = %v, want %v", vals, got, want)
				break
			}
		}
	}
}

func TestReorderListRandom(t *testing.T) {
	for i := 0; i < 1000; i++ {
		n := rand.Intn(50)
		vals := make([]int, n)
		for j := range vals {
			vals[j] = rand.Intn(100)
		}
		head := lc.NewLinkList(vals)
		runSafely(t, func() { lc.ReorderList(head) })
		got := listToSlice(head)
		want := refReorderSeq(vals)
		if len(got) != len(want) {
			t.Fatalf("ReorderList(%v) 长度错误", vals)
		}
		for j := range got {
			if got[j] != want[j] {
				t.Fatalf("ReorderList(%v) = %v, want %v", vals, got, want)
			}
		}
	}
}

// ---- 25. K 个一组翻转链表 ----

// refReverseKGroupSeq: 切片分组反转构造期望序列
func refReverseKGroupSeq(vals []int, k int) []int {
	out := append([]int(nil), vals...)
	for i := 0; i+k <= len(out); i += k {
		for a, b := i, i+k-1; a < b; a, b = a+1, b-1 {
			out[a], out[b] = out[b], out[a]
		}
	}
	return out
}

func TestReverseKGroup(t *testing.T) {
	cases := []struct {
		vals []int
		k    int
	}{
		{[]int{1, 2, 3, 4, 5}, 2},
		{[]int{1, 2, 3, 4, 5}, 3},
		{[]int{1, 2, 3, 4, 5}, 1},
		{[]int{1}, 1},
		{[]int{1, 2, 3, 4, 5, 6, 7}, 4}, // 尾部不足一组保持原序
	}
	for _, c := range cases {
		var gotSeq []int
		runSafely(t, func() {
			head := lc.ReverseKGroup(lc.NewLinkList(c.vals), c.k)
			gotSeq = listToSlice(head)
		})
		want := refReverseKGroupSeq(c.vals, c.k)
		if len(gotSeq) != len(want) {
			t.Errorf("ReverseKGroup(%v, %d) 长度 = %d, want %d", c.vals, c.k, len(gotSeq), len(want))
			continue
		}
		for i := range gotSeq {
			if gotSeq[i] != want[i] {
				t.Errorf("ReverseKGroup(%v, %d) = %v, want %v", c.vals, c.k, gotSeq, want)
				break
			}
		}
	}
}

func TestReverseKGroupRandom(t *testing.T) {
	for i := 0; i < 1000; i++ {
		n := rand.Intn(50)
		k := 1 + rand.Intn(8)
		vals := make([]int, n)
		for j := range vals {
			vals[j] = rand.Intn(100)
		}
		var gotSeq []int
		runSafely(t, func() {
			head := lc.ReverseKGroup(lc.NewLinkList(append([]int(nil), vals...)), k)
			gotSeq = listToSlice(head)
		})
		want := refReverseKGroupSeq(vals, k)
		if len(gotSeq) != len(want) {
			t.Fatalf("ReverseKGroup(%v, %d) 长度错误(节点丢失?)", vals, k)
		}
		for j := range gotSeq {
			if gotSeq[j] != want[j] {
				t.Fatalf("ReverseKGroup(%v, %d) = %v, want %v", vals, k, gotSeq, want)
			}
		}
	}
}
