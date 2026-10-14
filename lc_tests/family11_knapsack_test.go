package lc_tests

import (
	"math/rand"
	"testing"

	"demo/lc"
)

// ---- 300. 最长递增子序列 ----

func refLengthOfLIS(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	dp := make([]int, n)
	best := 1
	for i := range dp {
		dp[i] = 1
	}
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			if nums[j] < nums[i] && dp[j]+1 > dp[i] {
				dp[i] = dp[j] + 1
			}
		}
		if dp[i] > best {
			best = dp[i]
		}
	}
	return best
}

func TestLengthOfLIS(t *testing.T) {
	cases := []struct {
		nums []int
		want int
	}{
		{[]int{10, 9, 2, 5, 3, 7, 101, 18}, 4},
		{[]int{0, 1, 0, 3, 2, 3}, 4},
		{[]int{7, 7, 7, 7}, 1},
		{[]int{5}, 1},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.LengthOfLIS(c.nums) })
		if got != c.want {
			t.Errorf("LengthOfLIS(%v) = %d, want %d", c.nums, got, c.want)
		}
	}
}

func TestLengthOfLISRandom(t *testing.T) {
	for i := 0; i < 3000; i++ {
		n := 1 + rand.Intn(30)
		nums := make([]int, n)
		for j := range nums {
			nums[j] = rand.Intn(20)
		}
		want := refLengthOfLIS(nums)
		var got int
		runSafely(t, func() { got = lc.LengthOfLIS(nums) })
		if got != want {
			t.Fatalf("LengthOfLIS(%v) = %d, want %d", nums, got, want)
		}
	}
}

// ---- 322. 零钱兑换 ----

func refCoinChange(coins []int, amount int) int {
	// 独立 DP
	const inf = 1 << 30
	dp := make([]int, amount+1)
	for i := 1; i <= amount; i++ {
		dp[i] = inf
	}
	for i := 1; i <= amount; i++ {
		for _, c := range coins {
			if c <= i && dp[i-c]+1 < dp[i] {
				dp[i] = dp[i-c] + 1
			}
		}
	}
	if dp[amount] >= inf {
		return -1
	}
	return dp[amount]
}

func TestCoinChange(t *testing.T) {
	cases := []struct {
		coins  []int
		amount int
		want   int
	}{
		{[]int{1, 2, 5}, 11, 3},
		{[]int{2}, 3, -1},
		{[]int{1}, 0, 0},
		{[]int{1}, 1, 1},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.CoinChange(c.coins, c.amount) })
		if got != c.want {
			t.Errorf("CoinChange(%v, %d) = %d, want %d", c.coins, c.amount, got, c.want)
		}
	}
}

func TestCoinChangeRandom(t *testing.T) {
	for i := 0; i < 2000; i++ {
		k := 1 + rand.Intn(4)
		coins := make([]int, k)
		for j := range coins {
			coins[j] = 1 + rand.Intn(10)
		}
		amount := rand.Intn(40)
		want := refCoinChange(coins, amount)
		var got int
		runSafely(t, func() { got = lc.CoinChange(coins, amount) })
		if got != want {
			t.Fatalf("CoinChange(%v, %d) = %d, want %d", coins, amount, got, want)
		}
	}
}

// ---- 416. 分割等和子集 ----

func refCanPartition(nums []int) bool {
	sum := 0
	for _, v := range nums {
		sum += v
	}
	if sum%2 != 0 {
		return false
	}
	half := sum / 2
	n := len(nums)
	for mask := 0; mask < 1<<n; mask++ {
		s := 0
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				s += nums[i]
			}
		}
		if s == half {
			return true
		}
	}
	return false
}

func TestCanPartition(t *testing.T) {
	cases := []struct {
		nums []int
		want bool
	}{
		{[]int{1, 5, 11, 5}, true},
		{[]int{1, 2, 3, 5}, false},
		{[]int{2}, false},
		{[]int{1, 1}, true},
	}
	for _, c := range cases {
		var got bool
		runSafely(t, func() { got = lc.CanPartition(c.nums) })
		if got != c.want {
			t.Errorf("CanPartition(%v) = %v, want %v", c.nums, got, c.want)
		}
	}
}

func TestCanPartitionRandom(t *testing.T) {
	for i := 0; i < 2000; i++ {
		n := 1 + rand.Intn(13)
		nums := make([]int, n)
		for j := range nums {
			nums[j] = 1 + rand.Intn(15)
		}
		want := refCanPartition(nums)
		var got bool
		runSafely(t, func() { got = lc.CanPartition(nums) })
		if got != want {
			t.Fatalf("CanPartition(%v) = %v, want %v", nums, got, want)
		}
	}
}

// ---- 494. 目标和 ----

func refFindTargetSumWays(nums []int, target int) int {
	cnt := 0
	n := len(nums)
	for mask := 0; mask < 1<<n; mask++ {
		s := 0
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				s += nums[i]
			} else {
				s -= nums[i]
			}
		}
		if s == target {
			cnt++
		}
	}
	return cnt
}

func TestFindTargetSumWays(t *testing.T) {
	cases := []struct {
		nums   []int
		target int
		want   int
	}{
		{[]int{1, 1, 1, 1, 1}, 3, 5},
		{[]int{1}, 1, 1},
		{[]int{1}, 2, 0},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.FindTargetSumWays(c.nums, c.target) })
		if got != c.want {
			t.Errorf("FindTargetSumWays(%v, %d) = %d, want %d", c.nums, c.target, got, c.want)
		}
	}
}

func TestFindTargetSumWaysRandom(t *testing.T) {
	for i := 0; i < 2000; i++ {
		n := 1 + rand.Intn(13)
		nums := make([]int, n)
		for j := range nums {
			nums[j] = rand.Intn(6)
		}
		target := rand.Intn(21) - 10
		want := refFindTargetSumWays(nums, target)
		var got int
		runSafely(t, func() { got = lc.FindTargetSumWays(nums, target) })
		if got != want {
			t.Fatalf("FindTargetSumWays(%v, %d) = %d, want %d", nums, target, got, want)
		}
	}
}

// ---- 518. 零钱兑换 II ----

func refChange(amount int, coins []int) int {
	// 独立 DP(组合计数)
	dp := make([]int, amount+1)
	dp[0] = 1
	for _, c := range coins {
		for i := c; i <= amount; i++ {
			dp[i] += dp[i-c]
		}
	}
	return dp[amount]
}

func TestChange(t *testing.T) {
	cases := []struct {
		amount int
		coins  []int
		want   int
	}{
		{5, []int{1, 2, 5}, 4},
		{3, []int{2}, 0},
		{10, []int{10}, 1},
		{0, []int{1}, 1},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.Change(c.amount, c.coins) })
		if got != c.want {
			t.Errorf("Change(%d, %v) = %d, want %d", c.amount, c.coins, got, c.want)
		}
	}
}

func TestChangeRandom(t *testing.T) {
	for i := 0; i < 2000; i++ {
		k := 1 + rand.Intn(4)
		coins := make([]int, k)
		for j := range coins {
			coins[j] = 1 + rand.Intn(8)
		}
		amount := rand.Intn(30)
		want := refChange(amount, coins)
		var got int
		runSafely(t, func() { got = lc.Change(amount, coins) })
		if got != want {
			t.Fatalf("Change(%d, %v) = %d, want %d", amount, coins, got, want)
		}
	}
}

// ---- 1143. 最长公共子序列 ----

func refLCS(a, b string) int {
	// 无 memo 递归, 限小字符串
	var dfs func(i, j int) int
	dfs = func(i, j int) int {
		if i == len(a) || j == len(b) {
			return 0
		}
		if a[i] == b[j] {
			return 1 + dfs(i+1, j+1)
		}
		return max(dfs(i+1, j), dfs(i, j+1))
	}
	return dfs(0, 0)
}

func TestLongestCommonSubsequence(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"abcde", "ace", 3},
		{"abc", "abc", 3},
		{"abc", "def", 0},
		{"", "abc", 0},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.LongestCommonSubsequence(c.a, c.b) })
		if got != c.want {
			t.Errorf("LongestCommonSubsequence(%q, %q) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestLongestCommonSubsequenceRandom(t *testing.T) {
	for i := 0; i < 800; i++ {
		la, lb := rand.Intn(10), rand.Intn(10)
		var sa, sb []byte
		for j := 0; j < la; j++ {
			sa = append(sa, byte('a'+rand.Intn(3)))
		}
		for j := 0; j < lb; j++ {
			sb = append(sb, byte('a'+rand.Intn(3)))
		}
		a, b := string(sa), string(sb)
		want := refLCS(a, b)
		var got int
		runSafely(t, func() { got = lc.LongestCommonSubsequence(a, b) })
		if got != want {
			t.Fatalf("LongestCommonSubsequence(%q, %q) = %d, want %d", a, b, got, want)
		}
	}
}
