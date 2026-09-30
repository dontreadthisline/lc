package lc_tests

import (
	"demo/lc"
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"testing"
)

// 215. 数组中的第K个最大元素
func TestFindKthLargest(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		{"官方用例1", []int{3, 2, 1, 5, 6, 4}, 2, 5},
		{"官方用例2_含重复", []int{3, 2, 3, 1, 2, 4, 5, 5, 6}, 4, 4},
		{"单元素", []int{1}, 1, 1},
		{"k等于长度", []int{2, 1}, 2, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got int
			panicked := safeCall(func() { got = lc.FindKthLargest(c.nums, c.k) })
			if panicked {
				t.Errorf("FindKthLargest(%v, %d) panic", c.nums, c.k)
				return
			}
			if got != c.want {
				t.Errorf("FindKthLargest(%v, %d) = %d, want %d", c.nums, c.k, got, c.want)
			}
		})
	}
	rng := rand.New(rand.NewSource(41))
	for iter := 0; iter < 5000; iter++ {
		n := 1 + rng.Intn(30)
		nums := make([]int, n)
		for i := range nums {
			nums[i] = rng.Intn(20)
		}
		k := 1 + rng.Intn(n)
		sorted := slices.Clone(nums)
		slices.Sort(sorted)
		want := sorted[n-k]
		if got := lc.FindKthLargest(nums, k); got != want {
			t.Fatalf("iter=%d nums=%v k=%d: got=%d want=%d", iter, nums, k, got, want)
		}
	}
}

// 347. 前 K 个高频元素（返回集合不受顺序影响）
func TestTopKFrequent(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		k    int
		want []int
	}{
		{"官方用例1", []int{1, 1, 1, 2, 2, 3}, 2, []int{1, 2}},
		{"官方用例2", []int{1}, 1, []int{1}},
		{"无并列", []int{3, 3, 3, 3, 5, 5, 7}, 2, []int{3, 5}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got []int
			panicked := safeCall(func() { got = lc.TopKFrequent(c.nums, c.k) })
			if panicked {
				t.Errorf("TopKFrequent(%v, %d) panic", c.nums, c.k)
				return
			}
			if len(got) != c.k || !containsAll(got, c.want) {
				t.Errorf("TopKFrequent(%v, %d) = %v, want 集合 %v", c.nums, c.k, got, c.want)
			}
		})
	}
	rng := rand.New(rand.NewSource(42))
	for iter := 0; iter < 3000; iter++ {
		n := 1 + rng.Intn(20)
		nums := make([]int, n)
		for i := range nums {
			nums[i] = rng.Intn(6)
		}
		distinct := len(freqsOf(nums))
		k := 1 + rng.Intn(distinct) // 题目约束 k <= distinct
		got := lc.TopKFrequent(nums, k)
		if len(got) != k {
			t.Fatalf("iter=%d nums=%v k=%d: 返回长度 %d", iter, nums, k, len(got))
		}
		freqs := freqsOf(nums)
		cnt := map[int]int{}
		for _, x := range nums {
			cnt[x]++
		}
		threshold := freqs[len(freqs)-k] // 第 k 高频次的下界
		for _, x := range got {
			if cnt[x] < threshold {
				t.Fatalf("iter=%d nums=%v k=%d: %d 频次 %d < 门槛 %d", iter, nums, k, x, cnt[x], threshold)
			}
		}
	}
}

func containsAll(got, want []int) bool {
	for _, w := range want {
		if !slices.Contains(got, w) {
			return false
		}
	}
	return true
}

// 23. 合并 K 个升序链表
func TestMergeKLists(t *testing.T) {
	t.Run("官方用例", func(t *testing.T) {
		lists := []*lc.ListNode{
			lc.NewLinkList([]int{1, 4, 5}),
			lc.NewLinkList([]int{1, 3, 4}),
			lc.NewLinkList([]int{2, 6}),
		}
		var got []int
		panicked := safeCall(func() { got = listToSlice(lc.MergeKLists(lists)) })
		if panicked {
			t.Fatalf("MergeKLists panic")
		}
		want := []int{1, 1, 2, 3, 4, 4, 5, 6}
		if !slices.Equal(got, want) {
			t.Errorf("MergeKLists = %v, want %v", got, want)
		}
	})
	t.Run("含空链表", func(t *testing.T) {
		lists := []*lc.ListNode{nil, lc.NewLinkList([]int{1}), nil}
		got := listToSlice(lc.MergeKLists(lists))
		if !slices.Equal(got, []int{1}) {
			t.Errorf("含空链表合并 = %v, want [1]", got)
		}
	})
	t.Run("空列表", func(t *testing.T) {
		if got := lc.MergeKLists(nil); got != nil {
			t.Errorf("空输入应 nil, got %v", got)
		}
	})
}

// 136. 只出现一次的数字
func TestSingleNumber(t *testing.T) {
	if got := lc.SingleNumber([]int{2, 2, 1}); got != 1 {
		t.Errorf("SingleNumber([2,2,1]) = %d, want 1", got)
	}
	if got := lc.SingleNumber([]int{4, 1, 2, 1, 2}); got != 4 {
		t.Errorf("SingleNumber = %d, want 4", got)
	}
	if got := lc.SingleNumber([]int{7}); got != 7 {
		t.Errorf("单元素 = %d, want 7", got)
	}
	rng := rand.New(rand.NewSource(43))
	for iter := 0; iter < 5000; iter++ {
		pairs := rng.Intn(20)
		nums := make([]int, 0, 2*pairs+1)
		used := map[int]bool{}
		for i := 0; i < pairs; i++ {
			v := rng.Intn(50)
			for used[v] {
				v = rng.Intn(50)
			}
			used[v] = true
			nums = append(nums, v, v)
		}
		solo := rng.Intn(1000)
		for used[solo] {
			solo++
		}
		nums = append(nums, solo)
		rng.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
		if got := lc.SingleNumber(nums); got != solo {
			t.Fatalf("iter=%d: got=%d want=%d", iter, got, solo)
		}
	}
}

// 191. 位1的个数
func TestHammingWeight(t *testing.T) {
	cases := []struct {
		n    uint32
		want int
	}{{11, 3}, {128, 1}, {4294967293, 31}, {0, 0}, {1, 1}}
	for _, c := range cases {
		if got := lc.HammingWeight(c.n); got != c.want {
			t.Errorf("HammingWeight(%d) = %d, want %d", c.n, got, c.want)
		}
	}
}

// 338. 比特位计数
func TestCountBits(t *testing.T) {
	cases := []struct {
		n    int
		want []int
	}{{2, []int{0, 1, 1}}, {5, []int{0, 1, 1, 2, 1, 2}}, {0, []int{0}}, {1, []int{0, 1}}}
	for _, c := range cases {
		var got []int
		panicked := safeCall(func() { got = lc.CountBits(c.n) })
		if panicked {
			t.Errorf("CountBits(%d) panic", c.n)
			continue
		}
		if !slices.Equal(got, c.want) {
			t.Errorf("CountBits(%d) = %v, want %v", c.n, got, c.want)
		}
	}
}

// 208. 实现 Trie（前缀树）
func TestTrie(t *testing.T) {
	var tr *lc.Trie
	panicked := safeCall(func() { tr = lc.NewTrie() })
	if panicked {
		t.Fatal("NewTrie panic")
	}
	check := func(name string, fn func(), expectPanic bool) {
		t.Helper()
		if p := safeCall(fn); p != expectPanic {
			t.Errorf("%s panic=%v, want %v", name, p, expectPanic)
		}
	}
	var b bool
	check("Insert apple", func() { tr.Insert("apple") }, false)
	check("Search apple", func() { b = tr.Search("apple") }, false)
	if !b {
		t.Errorf("Search(apple) = false, want true")
	}
	check("Search app", func() { b = tr.Search("app") }, false)
	if b {
		t.Errorf("Search(app) = true, want false")
	}
	check("StartsWith app", func() { b = tr.StartsWith("app") }, false)
	if !b {
		t.Errorf("StartsWith(app) = false, want true")
	}
	check("Insert app", func() { tr.Insert("app") }, false)
	check("Search app again", func() { b = tr.Search("app") }, false)
	if !b {
		t.Errorf("Search(app) after insert = false, want true")
	}
	check("StartsWith b", func() { b = tr.StartsWith("b") }, false)
	if b {
		t.Errorf("StartsWith(b) = true, want false")
	}
}

// 912. 排序数组
func TestSortArray(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want []int
	}{
		{"官方用例1", []int{5, 2, 3, 1}, []int{1, 2, 3, 5}},
		{"官方用例2", []int{5, 1, 1, 2, 0, 0}, []int{0, 0, 1, 1, 2, 5}},
		{"单元素", []int{1}, []int{1}},
		{"已排序", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"逆序", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got []int
			panicked := safeCall(func() {
				got = lc.SortArray(slices.Clone(c.nums))
			})
			if panicked {
				t.Errorf("SortArray(%v) panic", c.nums)
				return
			}
			if !slices.Equal(got, c.want) {
				t.Errorf("SortArray(%v) = %v, want %v", c.nums, got, c.want)
			}
		})
	}
	rng := rand.New(rand.NewSource(44))
	for iter := 0; iter < 3000; iter++ {
		n := rng.Intn(60)
		nums := make([]int, n)
		for i := range nums {
			nums[i] = rng.Intn(200) - 100
		}
		want := slices.Clone(nums)
		slices.Sort(want)
		if got := lc.SortArray(nums); !slices.Equal(got, want) {
			t.Fatalf("iter=%d nums=%v: got %v", iter, nums, got)
		}
	}
}

// 179. 最大数
func TestLargestNumber(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want string
	}{
		{"官方用例1", []int{10, 2}, "210"},
		{"官方用例2", []int{3, 30, 34, 5, 9}, "9534330"},
		{"全零", []int{0, 0}, "0"},
		{"单元素", []int{10}, "10"},
		{"前缀对_直比会错", []int{121, 12}, "12121"},
		{"前缀对2", []int{4, 45}, "454"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got string
			panicked := safeCall(func() { got = lc.LargestNumber(slices.Clone(c.nums)) })
			if panicked {
				t.Errorf("LargestNumber(%v) panic", c.nums)
				return
			}
			if got != c.want {
				t.Errorf("LargestNumber(%v) = %q, want %q", c.nums, got, c.want)
			}
		})
	}
	_ = strings.Compare
	_ = fmt.Sprintf
}

func freqsOf(nums []int) []int {
	cnt := map[int]int{}
	for _, x := range nums {
		cnt[x]++
	}
	freqs := make([]int, 0, len(cnt))
	for _, c := range cnt {
		freqs = append(freqs, c)
	}
	slices.Sort(freqs)
	return freqs
}
