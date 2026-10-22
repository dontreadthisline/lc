package lc_tests

import (
	"math/rand"
	"testing"

	"demo/lc"
)

// ---- 二分答案合集 II (2064/2439/1292/2560/1898/1802) ----

// == 2064 分配给商店的最多商品的最小值 (min-max) ==

func TestMinimizedMaximum(t *testing.T) {
	cases := []struct {
		n          int
		quantities []int
		want       int
	}{
		{6, []int{11, 6}, 3},
		{7, []int{15, 10, 10}, 5},
		{1, []int{100000}, 100000},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.MinimizedMaximum(c.n, c.quantities) })
		if got != c.want {
			t.Errorf("MinimizedMaximum(%d, %v) = %d, want %d", c.n, c.quantities, got, c.want)
		}
	}
}

func TestMinimizedMaximumRandom(t *testing.T) {
	for i := 0; i < 500; i++ {
		n := 1 + rand.Intn(10)
		m := 1 + rand.Intn(5)
		q := make([]int, m)
		for j := range q {
			q[j] = 1 + rand.Intn(20)
		}
		// 暴力: 枚举上限 T, check = Σceil(q/T) ≤ n
		want := -1
		for T := 1; T <= 20; T++ {
			stores := 0
			for _, v := range q {
				stores += (v + T - 1) / T
			}
			if stores <= n {
				want = T
				break
			}
		}
		if want == -1 {
			continue
		}
		var got int
		runSafely(t, func() { got = lc.MinimizedMaximum(n, q) })
		if got != want {
			t.Fatalf("MinimizedMaximum(%d, %v) = %d, want %d", n, q, got, want)
		}
	}
}

// == 2439 最小化数组中的最大值 (min-max; 存在 O(n) 前缀均值公式作对照) ==

func TestMinimizeArrayValue(t *testing.T) {
	cases := []struct {
		nums []int
		want int
	}{
		{[]int{3, 7, 1, 6}, 5},
		{[]int{10, 1}, 10},
		{[]int{1}, 1},
		{[]int{1000000000, 1, 1, 1}, 250000000}, // ceil((1e9+3)/4)
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.MinimizeArrayValue(c.nums) })
		if got != c.want {
			t.Errorf("MinimizeArrayValue(%v) = %d, want %d", c.nums, got, c.want)
		}
	}
}

func TestMinimizeArrayValueRandom(t *testing.T) {
	// 独立参照: 数学公式 answer = max(ceil(prefix[i]/(i+1))) —— 与二分完全不同的方法
	for i := 0; i < 2000; i++ {
		n := 1 + rand.Intn(20)
		nums := make([]int, n)
		for j := range nums {
			nums[j] = rand.Intn(15)
		}
		want, pre := 0, 0
		for j, v := range nums {
			pre += v
			if c := (pre + j) / (j + 1); c > want {
				want = c
			}
		}
		var got int
		runSafely(t, func() { got = lc.MinimizeArrayValue(nums) })
		if got != want {
			t.Fatalf("MinimizeArrayValue(%v) = %d, want %d(公式)", nums, got, want)
		}
	}
}

// == 1292 元素和≤阈值的正方形最大边长 (二分边长 + 二维前缀和, 联动 304) ==

func TestMaxSideLength(t *testing.T) {
	cases := []struct {
		mat       [][]int
		threshold int
		want      int
	}{
		{[][]int{{1, 1, 3, 2, 4}, {3, 3, 1, 2, 5}, {2, 2, 3, 4, 5}, {1, 1, 3, 3, 4}, {2, 3, 1, 1, 3}}, 12, 2},
		{[][]int{{18, 70}, {61, 1}, {25, 85}, {14, 40}, {11, 96}, {97, 96}, {63, 45}}, 40184, 2},
		{[][]int{{1}}, 1, 1},
		{[][]int{{1}}, 0, 0},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.MaxSideLength(c.mat, c.threshold) })
		if got != c.want {
			t.Errorf("MaxSideLength(threshold=%d) = %d, want %d", c.threshold, got, c.want)
		}
	}
}

func TestMaxSideLengthRandom(t *testing.T) {
	for i := 0; i < 400; i++ {
		m := 1 + rand.Intn(6)
		n := 1 + rand.Intn(6)
		mat := make([][]int, m)
		for r := range mat {
			mat[r] = make([]int, n)
			for c := range mat[r] {
				mat[r][c] = 1 + rand.Intn(9)
			}
		}
		threshold := rand.Intn(m * n * 9)
		// 暴力: 枚举边长, 每个位置暴力求和
		want := 0
		for side := 1; side <= min(m, n); side++ {
			found := false
			for r := 0; r+side <= m && !found; r++ {
				for c := 0; c+side <= n; c++ {
					s := 0
					for dr := 0; dr < side; dr++ {
						for dc := 0; dc < side; dc++ {
							s += mat[r+dr][c+dc]
						}
					}
					if s <= threshold {
						found = true
						break
					}
				}
			}
			if found {
				want = side
			} else {
				break // 更大的边长和更大, 不可能可行
			}
		}
		var got int
		runSafely(t, func() { got = lc.MaxSideLength(mat, threshold) })
		if got != want {
			t.Fatalf("MaxSideLength(随机 %dx%d, %d) = %d, want %d", m, n, threshold, got, want)
		}
	}
}

// == 2560 打家劫舍 IV (min-max + 贪心判定, 联动 198) ==

func TestMinCapability(t *testing.T) {
	cases := []struct {
		nums []int
		k    int
		want int
	}{
		{[]int{2, 3, 5, 9}, 2, 5},
		{[]int{2, 7, 9, 3, 1}, 2, 2},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.MinCapability(c.nums, c.k) })
		if got != c.want {
			t.Errorf("MinCapability(%v, %d) = %d, want %d", c.nums, c.k, got, c.want)
		}
	}
}

func TestMinCapabilityRandom(t *testing.T) {
	// 暴力: 枚举能力值 cap, check = 贪心数不相邻的 ≤cap 元素个数 ≥ k
	for i := 0; i < 500; i++ {
		n := 1 + rand.Intn(12)
		nums := make([]int, n)
		for j := range nums {
			nums[j] = 1 + rand.Intn(12)
		}
		k := 1 + rand.Intn((n+1)/2)
		want := -1
		for cap := 1; cap <= 12; cap++ {
			cnt, prev := 0, -2
			for j, v := range nums {
				if v <= cap && j > prev+1 {
					cnt++
					prev = j
				}
			}
			if cnt >= k {
				want = cap
				break
			}
		}
		if want == -1 {
			continue
		}
		var got int
		runSafely(t, func() { got = lc.MinCapability(nums, k) })
		if got != want {
			t.Fatalf("MinCapability(%v, %d) = %d, want %d", nums, k, got, want)
		}
	}
}

// == 1898 可移除字符的最大数目 ==

func TestMaximumRemovals(t *testing.T) {
	cases := []struct {
		s, p      string
		removable []int
		want      int
	}{
		{"abcacb", "ab", []int{3, 1, 0}, 2},
		{"abcbddddd", "abcd", []int{3, 2, 1, 4, 5, 6}, 1},
		{"abcab", "abc", []int{0, 1, 2, 3, 4}, 0},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.MaximumRemovals(c.s, c.p, c.removable) })
		if got != c.want {
			t.Errorf("MaximumRemovals(%q, %q, %v) = %d, want %d", c.s, c.p, c.removable, got, c.want)
		}
	}
}

func TestMaximumRemovalsRandom(t *testing.T) {
	// 暴力: k 从 0 递增, 每次构建删除后字符串做子序列判定
	isSubseq := func(s, p string) bool {
		i := 0
		for j := 0; j < len(s) && i < len(p); j++ {
			if s[j] == p[i] {
				i++
			}
		}
		return i == len(p)
	}
	for iter := 0; iter < 300; iter++ {
		n := 1 + rand.Intn(10)
		sb := make([]byte, n)
		for j := range sb {
			sb[j] = byte('a' + rand.Intn(3))
		}
		s := string(sb)
		ln := rand.Intn(n + 1)
		pb := make([]byte, ln)
		for j := range pb {
			pb[j] = byte('a' + rand.Intn(3))
		}
		p := string(pb)
		perm := rand.Perm(n)
		rm := perm[:rand.Intn(n+1)]
		want := 0
		for k := 0; k <= len(rm); k++ {
			del := make([]bool, n)
			for j := 0; j < k; j++ {
				del[rm[j]] = true
			}
			var t2 []byte
			for j := 0; j < n; j++ {
				if !del[j] {
					t2 = append(t2, s[j])
				}
			}
			if isSubseq(string(t2), p) {
				want = k
			} else {
				break
			}
		}
		var got int
		runSafely(t, func() { got = lc.MaximumRemovals(s, p, rm) })
		if got != want {
			t.Fatalf("MaximumRemovals(%q, %q, %v) = %d, want %d", s, p, rm, got, want)
		}
	}
}

// == 1802 有界数组中指定下标处的最大值 ==

func TestMaxValue(t *testing.T) {
	cases := []struct {
		n, index, maxSum int
		want             int
	}{
		{4, 2, 6, 2},
		{6, 1, 10, 3},
		{1, 0, 0, 0},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.MaxValue(c.n, c.index, c.maxSum) })
		if got != c.want {
			t.Errorf("MaxValue(%d, %d, %d) = %d, want %d", c.n, c.index, c.maxSum, got, c.want)
		}
	}
}

func TestMaxValueRandom(t *testing.T) {
	// 独立参照: 数学法 —— 指定下标取 v 时两侧等差衰减的最小总和 O(1) 可算, 枚举 v
	for i := 0; i < 500; i++ {
		n := 1 + rand.Intn(12)
		index := rand.Intn(n)
		maxSum := rand.Intn(n*n + 10)
		minSum := func(v int) int64 {
			total := int64(v)
			left := index
			if int64(left) > int64(v-1) {
				left = v - 1
			}
			total += int64(left)*int64(v) - int64(left)*int64(left+1)/2
			if index > v-1 {
				total += int64(index - (v - 1))
			}
			right := n - 1 - index
			if int64(right) > int64(v-1) {
				right = v - 1
			}
			total += int64(right)*int64(v) - int64(right)*int64(right+1)/2
			if n-1-index > v-1 {
				total += int64(n - 1 - index - (v - 1))
			}
			return total
		}
		want := 0
		for v := 0; ; v++ {
			if minSum(v) > int64(maxSum) {
				want = v - 1
				break
			}
		}
		var got int
		runSafely(t, func() { got = lc.MaxValue(n, index, maxSum) })
		if got != want {
			t.Fatalf("MaxValue(%d, %d, %d) = %d, want %d", n, index, maxSum, got, want)
		}
	}
}

// == B13: 2226 每个小孩最多能分到多少糖果 (max-min: 410 的镜像) ==

func TestMaximumCandies(t *testing.T) {
	cases := []struct {
		candies []int
		k       int64
		want    int64
	}{
		{[]int{5, 8, 6}, 3, 5},
		{[]int{2, 5}, 11, 0},
		{[]int{9}, 1, 9},
	}
	for _, c := range cases {
		var got int64
		runSafely(t, func() { got = lc.MaximumCandies(c.candies, c.k) })
		if got != c.want {
			t.Errorf("MaximumCandies(%v, %d) = %d, want %d", c.candies, c.k, got, c.want)
		}
	}
}

func TestMaximumCandiesRandom(t *testing.T) {
	for i := 0; i < 500; i++ {
		n := 1 + rand.Intn(8)
		candies := make([]int, n)
		for j := range candies {
			candies[j] = 1 + rand.Intn(15)
		}
		k := int64(1 + rand.Intn(12))
		// 暴力: 枚举每人份数 x(从大到小找最大可行 = 从小到大找最后可行)
		var want int64 = 0
		for x := int64(1); x <= 15; x++ {
			piles := int64(0)
			for _, c := range candies {
				piles += int64(c) / x
			}
			if piles >= k {
				want = x // 继续找更大的
			}
		}
		var got int64
		runSafely(t, func() { got = lc.MaximumCandies(candies, k) })
		if got != want {
			t.Fatalf("MaximumCandies(%v, %d) = %d, want %d", candies, k, got, want)
		}
	}
}
