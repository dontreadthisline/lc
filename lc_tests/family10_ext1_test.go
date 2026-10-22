package lc_tests

import (
	"math/rand"
	"testing"

	"demo/lc"
)

// ---- 家族 10 扩展训练营 I: 巩固组(63/120/746) + 划分/双状态组(91/1567/978) ----

// == 63 不同路径 II ==

func refUniquePathsWithObstacles(g [][]int) int {
	m, n := len(g), len(g[0])
	dp := make([][]int, m)
	for i := range dp {
		dp[i] = make([]int, n)
	}
	for i := 0; i < m; i++ {
		if g[i][0] == 1 {
			break
		}
		dp[i][0] = 1
	}
	for j := 0; j < n; j++ {
		if g[0][j] == 1 {
			break
		}
		dp[0][j] = 1
	}
	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			if g[i][j] == 1 {
				dp[i][j] = 0
			} else {
				dp[i][j] = dp[i-1][j] + dp[i][j-1]
			}
		}
	}
	return dp[m-1][n-1]
}

func TestUniquePathsWithObstacles(t *testing.T) {
	cases := []struct {
		grid [][]int
		want int
	}{
		{[][]int{{0, 0, 0}, {0, 1, 0}, {0, 0, 0}}, 2},
		{[][]int{{0, 1}, {0, 0}}, 1},
		{[][]int{{1}}, 0},
		{[][]int{{0}}, 1},
		{[][]int{{0, 0}, {1, 1}, {0, 0}}, 0},
	}
	for _, c := range cases {
		grid := make([][]int, len(c.grid))
		for i := range c.grid {
			grid[i] = append([]int(nil), c.grid[i]...)
		}
		var got int
		runSafely(t, func() { got = lc.UniquePathsWithObstacles(grid) })
		if got != c.want {
			t.Errorf("UniquePathsWithObstacles(%v) = %d, want %d", c.grid, got, c.want)
		}
	}
}

func TestUniquePathsWithObstaclesRandom(t *testing.T) {
	for i := 0; i < 1000; i++ {
		m := 1 + rand.Intn(7)
		n := 1 + rand.Intn(7)
		g := make([][]int, m)
		for r := range g {
			g[r] = make([]int, n)
			for c := range g[r] {
				if rand.Intn(8) == 0 {
					g[r][c] = 1
				}
			}
		}
		want := refUniquePathsWithObstacles(g)
		var got int
		runSafely(t, func() { got = lc.UniquePathsWithObstacles(g) })
		if got != want {
			t.Fatalf("UniquePathsWithObstacles 随机网格 = %d, want %d", got, want)
		}
	}
}

// == 120 三角形最小路径和 ==

func refMinimumTotal(tri [][]int) int {
	// 从底往上
	dp := append([]int(nil), tri[len(tri)-1]...)
	for i := len(tri) - 2; i >= 0; i-- {
		for j := 0; j <= i; j++ {
			dp[j] = tri[i][j] + min(dp[j], dp[j+1])
		}
	}
	return dp[0]
}

func TestMinimumTotal(t *testing.T) {
	cases := []struct {
		tri  [][]int
		want int
	}{
		{[][]int{{2}, {3, 4}, {6, 5, 7}, {4, 1, 8, 3}}, 11},
		{[][]int{{-10}}, -10},
		{[][]int{{1}, {2, 3}}, 3},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.MinimumTotal(c.tri) })
		if got != c.want {
			t.Errorf("MinimumTotal(%v) = %d, want %d", c.tri, got, c.want)
		}
	}
}

func TestMinimumTotalRandom(t *testing.T) {
	for i := 0; i < 1000; i++ {
		n := 1 + rand.Intn(8)
		tri := make([][]int, n)
		for r := range tri {
			tri[r] = make([]int, r+1)
			for c := range tri[r] {
				tri[r][c] = rand.Intn(19) - 9
			}
		}
		want := refMinimumTotal(tri)
		var got int
		runSafely(t, func() { got = lc.MinimumTotal(tri) })
		if got != want {
			t.Fatalf("MinimumTotal 随机三角 = %d, want %d", got, want)
		}
	}
}

// == 746 使用最小花费爬楼梯 ==

func TestMinCostClimbingStairs(t *testing.T) {
	cases := []struct {
		cost []int
		want int
	}{
		{[]int{10, 15, 20}, 15},
		{[]int{1, 100, 1, 1, 1, 100, 1, 1, 100, 1}, 6},
		{[]int{5}, 0},   // 只有一阶: 起点可为 1, 直接到顶
		{[]int{0, 0}, 0},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.MinCostClimbingStairs(c.cost) })
		if got != c.want {
			t.Errorf("MinCostClimbingStairs(%v) = %d, want %d", c.cost, got, c.want)
		}
	}
}

func TestMinCostClimbingStairsRandom(t *testing.T) {
	// 暴力参照: 递归从位置 i 到顶的最小花费
	var ref func(cost []int, i int, memo map[int]int) int
	ref = func(cost []int, i int, memo map[int]int) int {
		if i >= len(cost) {
			return 0
		}
		if v, ok := memo[i]; ok {
			return v
		}
		r := cost[i] + min(ref(cost, i+1, memo), ref(cost, i+2, memo))
		memo[i] = r
		return r
	}
	for i := 0; i < 2000; i++ {
		n := 1 + rand.Intn(20)
		cost := make([]int, n)
		for j := range cost {
			cost[j] = rand.Intn(100)
		}
		want := min(ref(cost, 0, map[int]int{}), ref(cost, 1, map[int]int{}))
		var got int
		runSafely(t, func() { got = lc.MinCostClimbingStairs(cost) })
		if got != want {
			t.Fatalf("MinCostClimbingStairs(%v) = %d, want %d", cost, got, want)
		}
	}
}

// == 91 解码方法 ==

func refNumDecodings(s string) int {
	n := len(s)
	dp := make([]int, n+1)
	dp[0] = 1
	if n >= 1 && s[0] != '0' {
		dp[1] = 1
	}
	for i := 2; i <= n; i++ {
		if s[i-1] != '0' {
			dp[i] += dp[i-1]
		}
		if s[i-2] == '1' || (s[i-2] == '2' && s[i-1] <= '6') {
			dp[i] += dp[i-2]
		}
	}
	return dp[n]
}

func TestNumDecodings(t *testing.T) {
	cases := []struct {
		s    string
		want int
	}{
		{"12", 2},
		{"226", 3},
		{"06", 0},
		{"0", 0},
		{"10", 1},
		{"27", 1},
		{"1111", 5},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.NumDecodings(c.s) })
		if got != c.want {
			t.Errorf("NumDecodings(%q) = %d, want %d", c.s, got, c.want)
		}
	}
}

func TestNumDecodingsRandom(t *testing.T) {
	for i := 0; i < 2000; i++ {
		n := rand.Intn(10)
		b := make([]byte, n)
		for j := range b {
			b[j] = byte('0' + rand.Intn(10))
		}
		s := string(b)
		want := refNumDecodings(s)
		var got int
		runSafely(t, func() { got = lc.NumDecodings(s) })
		if got != want {
			t.Fatalf("NumDecodings(%q) = %d, want %d", s, got, want)
		}
	}
}

// == 1567 乘积为正数的最长子数组长度 ==

func refGetMaxLen(nums []int) int {
	best := 0
	for i := 0; i < len(nums); i++ {
		prod := 1
		for j := i; j < len(nums); j++ {
			prod *= nums[j]
			if prod > 0 && j-i+1 > best {
				best = j - i + 1
			}
		}
	}
	return best
}

func TestGetMaxLen(t *testing.T) {
	cases := []struct {
		nums []int
		want int
	}{
		{[]int{1, -2, -3, 4}, 4},
		{[]int{0, 1, -2, -3, -4}, 3},
		{[]int{-1, -2, -3, 0, 1}, 2},
		{[]int{-1, 2}, 1},
		{[]int{1, 2, 3, 5, -6, 4, 0, 10}, 4},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.GetMaxLen(c.nums) })
		if got != c.want {
			t.Errorf("GetMaxLen(%v) = %d, want %d", c.nums, got, c.want)
		}
	}
}

func TestGetMaxLenRandom(t *testing.T) {
	for i := 0; i < 2000; i++ {
		n := 1 + rand.Intn(12)
		nums := make([]int, n)
		for j := range nums {
			nums[j] = rand.Intn(7) - 3
		}
		want := refGetMaxLen(nums)
		var got int
		runSafely(t, func() { got = lc.GetMaxLen(nums) })
		if got != want {
			t.Fatalf("GetMaxLen(%v) = %d, want %d", nums, got, want)
		}
	}
}

// == 978 最长湍流子数组 ==

func refMaxTurbulenceSize(arr []int) int {
	n := len(arr)
	best := 1
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			// 两种形态都试: > 开头 / < 开头
			for _, first := range []int{1, -1} {
				ok := true
				for k := i; k < j; k++ {
					c := 0
					if arr[k+1] > arr[k] {
						c = 1
					} else if arr[k+1] < arr[k] {
						c = -1
					}
					wantSign := first
					if (k-i)%2 == 1 {
						wantSign = -first
					}
					if c != wantSign {
						ok = false
						break
					}
				}
				if ok && j-i+1 > best {
					best = j - i + 1
				}
			}
		}
	}
	return best
}

func TestMaxTurbulenceSize(t *testing.T) {
	cases := []struct {
		arr  []int
		want int
	}{
		{[]int{9, 4, 2, 10, 7, 8, 8, 1, 9}, 5},
		{[]int{4, 8, 12, 16}, 2},
		{[]int{100}, 1},
		{[]int{1, 1, 1}, 1},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.MaxTurbulenceSize(c.arr) })
		if got != c.want {
			t.Errorf("MaxTurbulenceSize(%v) = %d, want %d", c.arr, got, c.want)
		}
	}
}

func TestMaxTurbulenceSizeRandom(t *testing.T) {
	for i := 0; i < 2000; i++ {
		n := 1 + rand.Intn(12)
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rand.Intn(5)
		}
		want := refMaxTurbulenceSize(arr)
		var got int
		runSafely(t, func() { got = lc.MaxTurbulenceSize(arr) })
		if got != want {
			t.Fatalf("MaxTurbulenceSize(%v) = %d, want %d", arr, got, want)
		}
	}
}
