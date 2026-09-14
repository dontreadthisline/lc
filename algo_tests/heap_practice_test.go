package algo_tests

import "demo/algo"

import (
	"math"
	"math/rand"
	"slices"
	"testing"
)

// ── LC 703 algo.KthLargest ──

func TestKthLargest(t *testing.T) {
	kl := algo.NewKthLargest(3, []int{4, 5, 8, 2})
	for _, c := range []struct {
		add, want int
	}{
		{3, 4}, {5, 5}, {10, 5}, {9, 8}, {4, 8},
	} {
		if got := kl.Add(c.add); got != c.want {
			t.Errorf("Add(%d) = %d, want %d", c.add, got, c.want)
		}
	}
	// k = 1 极端：堆里永远只有最大值
	kl1 := algo.NewKthLargest(1, []int{})
	if got := kl1.Add(-3); got != -3 {
		t.Errorf("k=1 Add(-3) = %d, want -3", got)
	}
	if got := kl1.Add(5); got != 5 {
		t.Errorf("k=1 Add(5) = %d, want 5", got)
	}
	if got := kl1.Add(-1); got != 5 {
		t.Errorf("k=1 Add(-1) = %d, want 5", got)
	}
}

// ── LC 1046 最后一块石头的重量 ──

func TestLastStoneWeight(t *testing.T) {
	if got := algo.LastStoneWeight([]int{2, 7, 4, 1, 8, 1}); got != 1 {
		t.Errorf("got %d, want 1", got)
	}
	if got := algo.LastStoneWeight([]int{1}); got != 1 {
		t.Errorf("单石头: got %d, want 1", got)
	}
	if got := algo.LastStoneWeight([]int{3, 3}); got != 0 {
		t.Errorf("两石同重: got %d, want 0", got)
	}
}

// ── LC 215 第 K 个最大元素 ──

func TestFindKthLargestHeap(t *testing.T) {
	if got := algo.FindKthLargestHeap([]int{3, 2, 1, 5, 6, 4}, 2); got != 5 {
		t.Errorf("got %d, want 5", got)
	}
	if got := algo.FindKthLargestHeap([]int{3, 2, 3, 1, 2, 4, 5, 5, 6}, 4); got != 4 {
		t.Errorf("got %d, want 4", got)
	}
	if got := algo.FindKthLargestHeap([]int{1}, 1); got != 1 {
		t.Errorf("单元素: got %d, want 1", got)
	}

	// 随机对照：标准库排序取第 k 大
	r := rand.New(rand.NewSource(7))
	for round := 0; round < 300; round++ {
		n := 1 + r.Intn(50)
		nums := make([]int, n)
		for i := range nums {
			nums[i] = r.Intn(20) // 小值域制造重复，专打"第k个最大≠第k个不同"
		}
		k := 1 + r.Intn(n)
		want := slices.Clone(nums)
		slices.Sort(want)
		wantVal := want[len(want)-k]
		if got := algo.FindKthLargestHeap(nums, k); got != wantVal {
			t.Fatalf("round %d: nums=%v k=%d got %d want %d", round, nums, k, got, wantVal)
		}
	}
}

// ── LC 973 最接近原点的 K 个点 ──

func TestKClosest(t *testing.T) {
	// 平面点集合比对：排序后比较（答案任意序）
	normalize := func(ps [][]int) [][]int {
		out := slices.Clone(ps)
		slices.SortFunc(out, func(a, b []int) int {
			if a[0] != b[0] {
				return a[0] - b[0]
			}
			return a[1] - b[1]
		})
		return out
	}

	got := algo.KClosest([][]int{{1, 3}, {-2, 2}}, 1)
	want := [][]int{{-2, 2}}
	if !slices.EqualFunc(normalize(got), want, slices.Equal[[]int]) {
		t.Errorf("got %v, want %v", got, want)
	}

	// 随机对照：按距离平方稳定排序取前 k
	r := rand.New(rand.NewSource(42))
	for round := 0; round < 200; round++ {
		n := 1 + r.Intn(30)
		points := make([][]int, n)
		for i := range points {
			points[i] = []int{r.Intn(21) - 10, r.Intn(21) - 10}
		}
		k := 1 + r.Intn(n)
		wantPts := slices.Clone(points)
		slices.SortStableFunc(wantPts, func(a, b []int) int {
			da, db := a[0]*a[0]+a[1]*a[1], b[0]*b[0]+b[1]*b[1]
			return da - db
		})
		wantPts = wantPts[:k]

		gotPts := algo.KClosest(points, k)
		if len(gotPts) != k {
			t.Fatalf("round %d: 返回 %d 个点, want %d", round, len(gotPts), k)
		}
		// 距离并列时任选合法：比对距离序列而非点位
		dists := func(ps [][]int) []int {
			ds := make([]int, len(ps))
			for i, p := range ps {
				ds[i] = p[0]*p[0] + p[1]*p[1]
			}
			slices.Sort(ds)
			return ds
		}
		if !slices.Equal(dists(gotPts), dists(wantPts)) {
			t.Fatalf("round %d: 距离序列不等 got %v want %v", round, dists(gotPts), dists(wantPts))
		}
	}
}

// ── LC 373 查找和最小的 K 对数字 ──

func TestKSmallestPairs(t *testing.T) {
	got := algo.KSmallestPairs([]int{1, 7, 11}, []int{2, 4, 6}, 3)
	want := [][]int{{1, 2}, {1, 4}, {1, 6}}
	if !slices.EqualFunc(got, want, slices.Equal[[]int]) {
		t.Errorf("got %v, want %v", got, want)
	}
	if got := algo.KSmallestPairs([]int{1, 2}, []int{3}, 3); len(got) != 2 {
		t.Errorf("k 超过对数应截断: got %v", got)
	}

	// 随机对照：全生成 + 排序取前 k（小规模）
	r := rand.New(rand.NewSource(11))
	for round := 0; round < 200; round++ {
		n1, n2 := 1+r.Intn(8), 1+r.Intn(8)
		a := make([]int, n1)
		b := make([]int, n2)
		for i := range a {
			a[i] = r.Intn(30)
		}
		for i := range b {
			b[i] = r.Intn(30)
		}
		slices.Sort(a)
		slices.Sort(b)
		k := 1 + r.Intn(n1*n2)

		all := make([][]int, 0, n1*n2)
		for _, x := range a {
			for _, y := range b {
				all = append(all, []int{x, y})
			}
		}
		slices.SortFunc(all, func(p, q []int) int { return (p[0] + p[1]) - (q[0] + q[1]) })
		want := all[:k]

		gotPairs := algo.KSmallestPairs(a, b, k)
		if len(gotPairs) != k {
			t.Fatalf("round %d: len=%d want %d", round, len(gotPairs), k)
		}
		for i := range want {
			if gotPairs[i][0]+gotPairs[i][1] != want[i][0]+want[i][1] {
				t.Fatalf("round %d: 第 %d 对和 %d != %d（集合内容可能错）", round, i,
					gotPairs[i][0]+gotPairs[i][1], want[i][0]+want[i][1])
			}
		}
	}
}

// ── LC 378 有序矩阵中第 K 小 ──

func TestKthSmallestMatrix(t *testing.T) {
	m := [][]int{{1, 5, 9}, {10, 11, 13}, {12, 13, 15}}
	if got := algo.KthSmallest(m, 8); got != 13 {
		t.Errorf("got %d, want 13", got)
	}
	if got := algo.KthSmallest([][]int{{-5}}, 1); got != -5 {
		t.Errorf("单元素: got %d, want -5", got)
	}

	// 随机对照：构造每行每列升序的矩阵（值 = 行序 + 列序 + 扰动？直接用 i*n+j 保证有序）
	r := rand.New(rand.NewSource(3))
	for round := 0; round < 100; round++ {
		n := 1 + r.Intn(10)
		mat := make([][]int, n)
		flat := []int{}
		base := r.Intn(100)
		for i := 0; i < n; i++ {
			mat[i] = make([]int, n)
			for j := 0; j < n; j++ {
				mat[i][j] = base + i*n + j // 行列均升序
				flat = append(flat, mat[i][j])
			}
		}
		slices.Sort(flat)
		k := 1 + r.Intn(n*n)
		if got := algo.KthSmallest(mat, k); got != flat[k-1] {
			t.Fatalf("round %d: n=%d k=%d got %d want %d", round, n, k, got, flat[k-1])
		}
	}
	_ = math.MaxInt
}

// ── 练习：合并 K 个有序数组前 m 个 ──

func TestMergeKSorted(t *testing.T) {
	if got := algo.MergeKSorted([][]int{{1, 4, 7, 10}, {2, 5, 8}, {0, 3, 6, 9}}, 7); !slices.Equal(got, []int{0, 1, 2, 3, 4, 5, 6}) {
		t.Errorf("got %v, want [0 1 2 3 4 5 6]", got)
	}
	// m 超过总元素数：全部返回
	if got := algo.MergeKSorted([][]int{{1, 3}, {2}}, 10); !slices.Equal(got, []int{1, 2, 3}) {
		t.Errorf("m 超限: got %v, want [1 2 3]", got)
	}
	// 空路
	if got := algo.MergeKSorted([][]int{{}, {5}, {}}, 1); !slices.Equal(got, []int{5}) {
		t.Errorf("空路: got %v, want [5]", got)
	}
	// 单路退化
	if got := algo.MergeKSorted([][]int{{5}}, 1); !slices.Equal(got, []int{5}) {
		t.Errorf("单路单元素: got %v, want [5]", got)
	}

	// 随机对照：全部拼接 + 排序取前 m
	r := rand.New(rand.NewSource(99))
	for round := 0; round < 200; round++ {
		ways := 1 + r.Intn(6)
		arrs := make([][]int, ways)
		total := 0
		for i := range arrs {
			n := r.Intn(8) // 允许空路
			arrs[i] = make([]int, n)
			for j := range arrs[i] {
				arrs[i][j] = r.Intn(100)
			}
			slices.Sort(arrs[i])
			total += n
		}
		m := 1 + r.Intn(total+1) // m ∈ [1, total+1]，含超限
		all := []int{}
		for _, a := range arrs {
			all = append(all, a...)
		}
		slices.Sort(all)
		want := all[:min(m, len(all))]

		if got := algo.MergeKSorted(arrs, m); !slices.Equal(got, want) {
			t.Fatalf("round %d: arrs=%v m=%d got %v want %v", round, arrs, m, got, want)
		}
	}
}

// ── LC 23 合并 K 个升序链表 ──

func listFrom(vals []int) *algo.ListNode {
	dummy := &algo.ListNode{}
	cur := dummy
	for _, v := range vals {
		cur.Next = &algo.ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

func listTo(h *algo.ListNode) []int {
	out := []int{}
	for ; h != nil; h = h.Next {
		out = append(out, h.Val)
	}
	return out
}

func TestMergeKLists(t *testing.T) {
	// 示例 1
	lists := []*algo.ListNode{listFrom([]int{1, 4, 5}), listFrom([]int{1, 3, 4}), listFrom([]int{2, 6})}
	if got := listTo(algo.MergeKLists(lists)); !slices.Equal(got, []int{1, 1, 2, 3, 4, 4, 5, 6}) {
		t.Errorf("示例1: got %v", got)
	}
	// 示例 2/3：空数组、空链表
	if got := algo.MergeKLists(nil); got != nil {
		t.Errorf("空数组: got %v", listTo(got))
	}
	if got := algo.MergeKLists([]*algo.ListNode{nil}); got != nil {
		t.Errorf("[[]]: got %v", listTo(got))
	}
	// 混合空链表
	lists2 := []*algo.ListNode{nil, listFrom([]int{2}), nil, listFrom([]int{1})}
	if got := listTo(algo.MergeKLists(lists2)); !slices.Equal(got, []int{1, 2}) {
		t.Errorf("混合空路: got %v", got)
	}

	// 随机对照：全部拼接 + 排序
	r := rand.New(rand.NewSource(23))
	for round := 0; round < 200; round++ {
		k := r.Intn(6) // 允许 0 路
		var lists []*algo.ListNode
		var all []int
		for i := 0; i < k; i++ {
			n := r.Intn(8) // 允许空路
			vals := make([]int, n)
			for j := range vals {
				vals[j] = r.Intn(100)
			}
			slices.Sort(vals)
			lists = append(lists, listFrom(vals))
			all = append(all, vals...)
		}
		slices.Sort(all)

		if got := listTo(algo.MergeKLists(lists)); !slices.Equal(got, all) {
			t.Fatalf("round %d: got %v want %v", round, got, all)
		}
	}
}

// ── LC 632 最小区间 ──

func TestSmallestRange(t *testing.T) {
	if got := algo.SmallestRange([][]int{{4, 10, 15, 24, 26}, {0, 9, 12, 20}, {5, 18, 22, 30}}); !slices.Equal(got, []int{20, 24}) {
		t.Errorf("示例1: got %v, want [20 24]", got)
	}
	// 题面其他示例
	if got := algo.SmallestRange([][]int{{1, 2, 3}, {1, 2, 3}, {1, 2, 3}}); !slices.Equal(got, []int{1, 1}) {
		t.Errorf("示例2: got %v, want [1 1]", got)
	}

	// 随机对照：暴力枚举所有选法（乘积小）
	r := rand.New(rand.NewSource(632))
	for round := 0; round < 200; round++ {
		k := 1 + r.Intn(4)
		nums := make([][]int, k)
		for i := range nums {
			n := 1 + r.Intn(5)
			nums[i] = make([]int, n)
			for j := range nums[i] {
				nums[i][j] = r.Intn(30)
			}
			slices.Sort(nums[i])
		}
		// 暴力：每路各取一个，枚举全部组合，取最小差
		bestDiff := math.MaxInt
		var want []int
		var pick func(row int, curMin, curMax int)
		var cur []int
		pick = func(row, curMin, curMax int) {
			if row == k {
				if d := curMax - curMin; d < bestDiff {
					bestDiff = d
					want = []int{curMin, curMax}
				}
				return
			}
			for _, v := range nums[row] {
				cur = append(cur, v)
				pick(row+1, min(curMin, v), max(curMax, v))
				cur = cur[:len(cur)-1]
			}
		}
		pick(0, math.MaxInt, math.MinInt)

		if got := algo.SmallestRange(nums); !slices.Equal(got, want) {
			t.Fatalf("round %d: nums=%v got %v want %v", round, nums, got, want)
		}
	}
}

// ── LC 295 数据流的中位数 ──

func TestMedianFinder(t *testing.T) {
	mf := algo.NewMedianFinder()
	mf.AddNum(1)
	mf.AddNum(2)
	if got := mf.FindMedian(); got != 1.5 {
		t.Errorf("got %v, want 1.5", got)
	}
	mf.AddNum(3)
	if got := mf.FindMedian(); got != 2 {
		t.Errorf("got %v, want 2", got)
	}
	// 交叉场景：新数比右顶大
	mf2 := algo.NewMedianFinder()
	mf2.AddNum(3)
	mf2.AddNum(8)
	mf2.AddNum(10) // left={3,10} right={8} → 应触发交叉修复
	if got := mf2.FindMedian(); got != 8 {
		t.Errorf("交叉: got %v, want 8", got)
	}

	// 随机对照：排序取中位
	r := rand.New(rand.NewSource(295))
	for round := 0; round < 300; round++ {
		n := 1 + r.Intn(60)
		mf := algo.NewMedianFinder()
		stream := make([]int, 0, n)
		for i := 0; i < n; i++ {
			v := r.Intn(200) - 100 // 含负数
			mf.AddNum(v)
			stream = append(stream, v)
			want := func() float64 {
				s := slices.Clone(stream)
				slices.Sort(s)
				if len(s)%2 == 1 {
					return float64(s[len(s)/2])
				}
				return float64(s[len(s)/2-1]+s[len(s)/2]) / 2
			}()
			if got := mf.FindMedian(); got != want {
				t.Fatalf("round %d step %d: got %v want %v (stream=%v)", round, i, got, want, stream)
			}
		}
	}
}

// ── LC 480 滑动窗口中位数 ──

func TestMedianSlidingWindow(t *testing.T) {
	got := algo.MedianSlidingWindow([]int{1, 3, -1, -3, 5, 3, 6, 7}, 3)
	want := []float64{1, -1, -1, 3, 5, 6}
	for i := range want {
		if math.Abs(got[i]-want[i]) > 1e-9 {
			t.Fatalf("题面用例: got %v want %v", got, want)
		}
	}
	// k=1：窗口单元素，中位数就是它自己
	got1 := algo.MedianSlidingWindow([]int{4, 2, 9}, 1)
	for i, v := range []float64{4, 2, 9} {
		if got1[i] != v {
			t.Fatalf("k=1: got %v", got1)
		}
	}
	// k=n：整窗一个中位数
	if g := algo.MedianSlidingWindow([]int{2, 1, 5, 7}, 4)[0]; g != 3.5 {
		t.Fatalf("k=n: got %v want 3.5", g)
	}

	// 随机对照：暴力每窗排序
	r := rand.New(rand.NewSource(480))
	for round := 0; round < 200; round++ {
		n := 1 + r.Intn(40)
		nums := make([]int, n)
		for i := range nums {
			nums[i] = r.Intn(20) - 10 // 小值域制造重复值，专打等值删除
		}
		k := 1 + r.Intn(n)
		wantBrute := make([]float64, 0, n-k+1)
		for i := k - 1; i < n; i++ {
			win := slices.Clone(nums[i-k+1 : i+1])
			slices.Sort(win)
			m := len(win) / 2
			if len(win)%2 == 1 {
				wantBrute = append(wantBrute, float64(win[m]))
			} else {
				wantBrute = append(wantBrute, float64(win[m-1]+win[m])/2)
			}
		}
		gotR := algo.MedianSlidingWindow(nums, k)
		if len(gotR) != len(wantBrute) {
			t.Fatalf("round %d: len %d != %d", round, len(gotR), len(wantBrute))
		}
		for i := range wantBrute {
			if math.Abs(gotR[i]-wantBrute[i]) > 1e-9 {
				t.Fatalf("round %d: nums=%v k=%d idx %d: got %v want %v",
					round, nums, k, i, gotR[i], wantBrute[i])
			}
		}
	}
}
