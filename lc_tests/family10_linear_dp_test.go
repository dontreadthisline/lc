package lc_tests

import (
	"math/rand"
	"strings"
	"testing"

	"demo/lc"
)

// ---- 53. 最大子数组和 ----

func refMaxSubArray(nums []int) int {
	best, sum := nums[0], 0
	for _, v := range nums {
		sum += v
		if sum > best {
			best = sum
		}
		if sum < 0 {
			sum = 0
		}
	}
	return best
}

func TestMaxSubArray(t *testing.T) {
	cases := []struct {
		nums []int
		want int
	}{
		{[]int{-2, 1, -3, 4, -1, 2, 1, -5, 4}, 6},
		{[]int{1}, 1},
		{[]int{5, 4, -1, 7, 8}, 23},
		{[]int{-1}, -1},
		{[]int{-3, -2, -5}, -2},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.MaxSubArray(c.nums) })
		if got != c.want {
			t.Errorf("MaxSubArray(%v) = %d, want %d", c.nums, got, c.want)
		}
	}
}

func TestMaxSubArrayRandom(t *testing.T) {
	for i := 0; i < 5000; i++ {
		n := 1 + rand.Intn(50)
		nums := make([]int, n)
		for j := range nums {
			nums[j] = rand.Intn(21) - 10
		}
		want := refMaxSubArray(nums)
		var got int
		runSafely(t, func() { got = lc.MaxSubArray(nums) })
		if got != want {
			t.Fatalf("MaxSubArray(%v) = %d, want %d", nums, got, want)
		}
	}
}

// ---- 62. 不同路径 ----

func refUniquePaths(m, n int) int {
	// 独立 DP 实现
	dp := make([]int, n)
	for j := range dp {
		dp[j] = 1
	}
	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			dp[j] += dp[j-1]
		}
	}
	return dp[n-1]
}

func TestUniquePaths(t *testing.T) {
	cases := []struct{ m, n, want int }{
		{3, 7, 28}, {3, 2, 3}, {1, 1, 1}, {1, 10, 1}, {10, 1, 1},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.UniquePaths(c.m, c.n) })
		if got != c.want {
			t.Errorf("UniquePaths(%d, %d) = %d, want %d", c.m, c.n, got, c.want)
		}
	}
}

func TestUniquePathsRandom(t *testing.T) {
	for i := 0; i < 500; i++ {
		m := 1 + rand.Intn(10)
		n := 1 + rand.Intn(10)
		want := refUniquePaths(m, n)
		var got int
		runSafely(t, func() { got = lc.UniquePaths(m, n) })
		if got != want {
			t.Fatalf("UniquePaths(%d, %d) = %d, want %d", m, n, got, want)
		}
	}
}

// ---- 64. 最小路径和 ----

func refMinPathSum(grid [][]int) int {
	m, n := len(grid), len(grid[0])
	dp := make([][]int, m)
	for i := range dp {
		dp[i] = make([]int, n)
	}
	dp[0][0] = grid[0][0]
	for i := 1; i < m; i++ {
		dp[i][0] = dp[i-1][0] + grid[i][0]
	}
	for j := 1; j < n; j++ {
		dp[0][j] = dp[0][j-1] + grid[0][j]
	}
	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			dp[i][j] = min(dp[i-1][j], dp[i][j-1]) + grid[i][j]
		}
	}
	return dp[m-1][n-1]
}

func TestMinPathSum(t *testing.T) {
	cases := []struct {
		grid [][]int
		want int
	}{
		{[][]int{{1, 3, 1}, {1, 5, 1}, {4, 2, 1}}, 7},
		{[][]int{{1, 2, 3}, {4, 5, 6}}, 12},
		{[][]int{{5}}, 5},
	}
	for _, c := range cases {
		grid := make([][]int, len(c.grid))
		for i := range c.grid {
			grid[i] = append([]int(nil), c.grid[i]...)
		}
		var got int
		runSafely(t, func() { got = lc.MinPathSum(grid) })
		if got != c.want {
			t.Errorf("MinPathSum(%v) = %d, want %d", c.grid, got, c.want)
		}
	}
}

func TestMinPathSumRandom(t *testing.T) {
	for i := 0; i < 1000; i++ {
		m := 1 + rand.Intn(8)
		n := 1 + rand.Intn(8)
		grid := make([][]int, m)
		for r := range grid {
			grid[r] = make([]int, n)
			for c := range grid[r] {
				grid[r][c] = rand.Intn(20)
			}
		}
		want := refMinPathSum(grid)
		var got int
		runSafely(t, func() { got = lc.MinPathSum(grid) })
		if got != want {
			t.Fatalf("MinPathSum 随机网格 = %d, want %d", got, want)
		}
	}
}

// ---- 139. 单词拆分 ----

func refWordBreak(s string, dict []string) bool {
	words := map[string]bool{}
	for _, w := range dict {
		words[w] = true
	}
	n := len(s)
	// 枚举 2^(n-1) 分割
	for mask := 0; mask < 1<<(n-1); mask++ {
		ok := true
		prev := 0
		for i := 0; i < n-1; i++ {
			if mask&(1<<i) != 0 {
				if !words[s[prev:i+1]] {
					ok = false
					break
				}
				prev = i + 1
			}
		}
		if ok && words[s[prev:]] {
			return true
		}
	}
	return false
}

func TestWordBreak(t *testing.T) {
	cases := []struct {
		s    string
		dict []string
		want bool
	}{
		{"leetcode", []string{"leet", "code"}, true},
		{"applepenapple", []string{"apple", "pen"}, true},
		{"catsandog", []string{"cats", "dog", "sand", "and", "cat"}, false},
		{"a", []string{"a"}, true},
		{"a", []string{"b"}, false},
	}
	for _, c := range cases {
		var got bool
		runSafely(t, func() { got = lc.WordBreak(c.s, c.dict) })
		if got != c.want {
			t.Errorf("WordBreak(%q, %v) = %v, want %v", c.s, c.dict, got, c.want)
		}
	}
}

func TestWordBreakRandom(t *testing.T) {
	dictPool := []string{"a", "ab", "abc", "bc", "c", "bca", "cab"}
	for i := 0; i < 1500; i++ {
		n := 1 + rand.Intn(11)
		var sb strings.Builder
		for j := 0; j < n; j++ {
			sb.WriteByte(byte('a' + rand.Intn(3)))
		}
		s := sb.String()
		dict := []string{}
		for _, w := range dictPool {
			if rand.Intn(2) == 0 {
				dict = append(dict, w)
			}
		}
		want := refWordBreak(s, dict)
		var got bool
		runSafely(t, func() { got = lc.WordBreak(s, dict) })
		if got != want {
			t.Fatalf("WordBreak(%q, %v) = %v, want %v", s, dict, got, want)
		}
	}
}

// ---- 152. 乘积最大子数组 ----

func refMaxProduct(nums []int) int {
	best := nums[0]
	for i := 0; i < len(nums); i++ {
		p := 1
		for j := i; j < len(nums); j++ {
			p *= nums[j]
			if p > best {
				best = p
			}
		}
	}
	return best
}

func TestMaxProduct(t *testing.T) {
	cases := []struct {
		nums []int
		want int
	}{
		{[]int{2, 3, -2, 4}, 6},
		{[]int{-2, 0, -1}, 0},
		{[]int{-2}, -2},
		{[]int{-4, -3}, 12},
		{[]int{0, 2}, 2},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.MaxProduct(c.nums) })
		if got != c.want {
			t.Errorf("MaxProduct(%v) = %d, want %d", c.nums, got, c.want)
		}
	}
}

func TestMaxProductRandom(t *testing.T) {
	for i := 0; i < 5000; i++ {
		n := 1 + rand.Intn(25)
		nums := make([]int, n)
		for j := range nums {
			nums[j] = rand.Intn(11) - 5
		}
		want := refMaxProduct(nums)
		var got int
		runSafely(t, func() { got = lc.MaxProduct(nums) })
		if got != want {
			t.Fatalf("MaxProduct(%v) = %d, want %d", nums, got, want)
		}
	}
}

// ---- 198. 打家劫舍 ----

func refRob(nums []int) int {
	n := len(nums)
	best := 0
	for mask := 0; mask < 1<<n; mask++ {
		ok, sum := true, 0
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				if i > 0 && mask&(1<<(i-1)) != 0 {
					ok = false
					break
				}
				sum += nums[i]
			}
		}
		if ok && sum > best {
			best = sum
		}
	}
	return best
}

func TestRob(t *testing.T) {
	cases := []struct {
		nums []int
		want int
	}{
		{[]int{1, 2, 3, 1}, 4},
		{[]int{2, 7, 9, 3, 1}, 12},
		{[]int{5}, 5},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.Rob(c.nums) })
		if got != c.want {
			t.Errorf("Rob(%v) = %d, want %d", c.nums, got, c.want)
		}
	}
}

func TestRobRandom(t *testing.T) {
	for i := 0; i < 3000; i++ {
		n := 1 + rand.Intn(15)
		nums := make([]int, n)
		for j := range nums {
			nums[j] = rand.Intn(20)
		}
		want := refRob(nums)
		var got int
		runSafely(t, func() { got = lc.Rob(nums) })
		if got != want {
			t.Fatalf("Rob(%v) = %d, want %d", nums, got, want)
		}
	}
}

// ---- 213. 打家劫舍 II (环形) ----

func refRobII(nums []int) int {
	n := len(nums)
	best := 0
	for mask := 0; mask < 1<<n; mask++ {
		ok, sum := true, 0
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				if i > 0 && mask&(1<<(i-1)) != 0 {
					ok = false
					break
				}
				if i == n-1 && mask&1 != 0 { // 环形: 首尾相邻
					ok = false
					break
				}
				sum += nums[i]
			}
		}
		if ok && sum > best {
			best = sum
		}
	}
	return best
}

func TestRobII(t *testing.T) {
	cases := []struct {
		nums []int
		want int
	}{
		{[]int{2, 3, 2}, 3},
		{[]int{1, 2, 3, 1}, 4},
		{[]int{1, 2, 3}, 3},
		{[]int{5}, 5},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.RobII(c.nums) })
		if got != c.want {
			t.Errorf("RobII(%v) = %d, want %d", c.nums, got, c.want)
		}
	}
}

func TestRobIIRandom(t *testing.T) {
	for i := 0; i < 3000; i++ {
		n := 1 + rand.Intn(14)
		nums := make([]int, n)
		for j := range nums {
			nums[j] = rand.Intn(20)
		}
		want := refRobII(nums)
		var got int
		runSafely(t, func() { got = lc.RobII(nums) })
		if got != want {
			t.Fatalf("RobII(%v) = %d, want %d", nums, got, want)
		}
	}
}

// ---- 70. 爬楼梯 ----

func TestClimbStairs(t *testing.T) {
	// 期望值: 迭代斐波那契 (want[0] 占位, n 从 1 起)
	want := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55}
	for n := 1; n <= 10; n++ {
		var got int
		runSafely(t, func() { got = lc.ClimbStairs(n) })
		if got != want[n] {
			t.Errorf("ClimbStairs(%d) = %d, want %d", n, got, want[n])
		}
	}
}

// ---- 118. 杨辉三角 ----

func TestGenerate(t *testing.T) {
	var got [][]int
	runSafely(t, func() { got = lc.Generate(5) })
	if len(got) != 5 {
		t.Fatalf("Generate(5) 行数 = %d, want 5", len(got))
	}
	// 独立构造: C(n,k)
	for i, row := range got {
		if len(row) != i+1 {
			t.Fatalf("Generate(5) 第 %d 行长度 = %d, want %d", i, len(row), i+1)
		}
		for k, v := range row {
			c := 1
			for j := 0; j < k; j++ {
				c = c * (i - j) / (j + 1)
			}
			if v != c {
				t.Errorf("Generate(5)[%d][%d] = %d, want C(%d,%d)=%d", i, k, v, i, k, c)
			}
		}
	}
}

// ---- 72. 编辑距离 ----

func refMinDistance(a, b string) int {
	// 无 memo 递归, 限小字符串
	var dfs func(i, j int) int
	dfs = func(i, j int) int {
		if i == len(a) {
			return len(b) - j
		}
		if j == len(b) {
			return len(a) - i
		}
		if a[i] == b[j] {
			return dfs(i+1, j+1)
		}
		return 1 + min(dfs(i+1, j+1), min(dfs(i+1, j), dfs(i, j+1)))
	}
	return dfs(0, 0)
}

func TestMinDistance(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"horse", "ros", 3},
		{"intention", "execution", 5},
		{"", "", 0},
		{"", "abc", 3},
		{"abc", "", 3},
		{"same", "same", 0},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.MinDistance(c.a, c.b) })
		if got != c.want {
			t.Errorf("MinDistance(%q, %q) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestMinDistanceRandom(t *testing.T) {
	for i := 0; i < 800; i++ {
		la, lb := rand.Intn(9), rand.Intn(9)
		var sa, sb strings.Builder
		for j := 0; j < la; j++ {
			sa.WriteByte(byte('a' + rand.Intn(4)))
		}
		for j := 0; j < lb; j++ {
			sb.WriteByte(byte('a' + rand.Intn(4)))
		}
		a, b := sa.String(), sb.String()
		want := refMinDistance(a, b)
		var got int
		runSafely(t, func() { got = lc.MinDistance(a, b) })
		if got != want {
			t.Fatalf("MinDistance(%q, %q) = %d, want %d", a, b, got, want)
		}
	}
}
