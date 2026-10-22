package lc_tests

import (
	"math/rand"
	"testing"

	"demo/lc"
)

// ---- 家族 10 扩展训练营 II: 值域/状态机组(740/309) + 双序列(718) + 家族联动(132) ----

// == 740 删除并获得点数 ==

func refDeleteAndEarn(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	maxV := nums[0]
	for _, v := range nums {
		if v > maxV {
			maxV = v
		}
	}
	// 值域统计 + 打家劫舍
	cnt := make([]int, maxV+1)
	for _, v := range nums {
		cnt[v] += v
	}
	dp := make([]int, maxV+1)
	if maxV >= 1 {
		dp[1] = cnt[1]
	}
	for i := 2; i <= maxV; i++ {
		dp[i] = max(dp[i-1], dp[i-2]+cnt[i])
	}
	return dp[maxV]
}

func TestDeleteAndEarn(t *testing.T) {
	cases := []struct {
		nums []int
		want int
	}{
		{[]int{3, 4, 2}, 6},
		{[]int{2, 2, 3, 3, 3, 4}, 9},
		{[]int{1}, 1},
		{[]int{}, 0},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.DeleteAndEarn(c.nums) })
		if got != c.want {
			t.Errorf("DeleteAndEarn(%v) = %d, want %d", c.nums, got, c.want)
		}
	}
}

func TestDeleteAndEarnRandom(t *testing.T) {
	for i := 0; i < 2000; i++ {
		n := rand.Intn(12)
		nums := make([]int, n)
		for j := range nums {
			nums[j] = 1 + rand.Intn(10)
		}
		want := refDeleteAndEarn(nums)
		var got int
		runSafely(t, func() { got = lc.DeleteAndEarn(nums) })
		if got != want {
			t.Fatalf("DeleteAndEarn(%v) = %d, want %d", nums, got, want)
		}
	}
}

// == 309 买卖股票的最佳时机含冷冻期 ==

func refMaxProfitCooldown(prices []int) int {
	// 状态机: hold / sold(刚卖, 冷冻) / rest(可买)
	const neg = -1 << 30
	hold, sold, rest := neg, 0, 0
	for _, p := range prices {
		newHold := max(hold, rest-p)
		newSold := hold + p
		newRest := max(rest, sold)
		hold, sold, rest = newHold, newSold, newRest
	}
	return max(sold, rest)
}

func TestMaxProfitCooldown(t *testing.T) {
	cases := []struct {
		prices []int
		want   int
	}{
		{[]int{1, 2, 3, 0, 2}, 3},
		{[]int{1}, 0},
		{[]int{2, 1}, 0},
		{[]int{1, 2, 4}, 3},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.MaxProfitCooldown(c.prices) })
		if got != c.want {
			t.Errorf("MaxProfitCooldown(%v) = %d, want %d", c.prices, got, c.want)
		}
	}
}

func TestMaxProfitCooldownRandom(t *testing.T) {
	// 暴力参照: 记忆化枚举 (day, holding, cooldown) 三状态
	var search func(prices []int, day, holding int, cooldown bool, memo map[[3]int]int) int
	search = func(prices []int, day, holding int, cooldown bool, memo map[[3]int]int) int {
		if day == len(prices) {
			return 0
		}
		key := [3]int{day, holding, 0}
		if cooldown {
			key[2] = 1
		}
		if v, ok := memo[key]; ok {
			return v
		}
		best := search(prices, day+1, holding, false, memo) // 不动
		if cooldown {
			// 冷冻期不能买
		} else if holding == 1 {
			best = max(best, prices[day]+search(prices, day+1, 0, true, memo)) // 卖 → 冷冻
		} else {
			best = max(best, -prices[day]+search(prices, day+1, 1, false, memo)) // 买
		}
		memo[key] = best
		return best
	}
	for i := 0; i < 1500; i++ {
		n := 1 + rand.Intn(8)
		prices := make([]int, n)
		for j := range prices {
			prices[j] = rand.Intn(10)
		}
		want := search(prices, 0, 0, false, map[[3]int]int{})
		var got int
		runSafely(t, func() { got = lc.MaxProfitCooldown(prices) })
		if got != want {
			t.Fatalf("MaxProfitCooldown(%v) = %d, want %d", prices, got, want)
		}
	}
}

// == 718 最长重复子数组(连续) ==

func refFindLength(a, b []int) int {
	best := 0
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			k := 0
			for i+k < len(a) && j+k < len(b) && a[i+k] == b[j+k] {
				k++
			}
			if k > best {
				best = k
			}
		}
	}
	return best
}

func TestFindLength(t *testing.T) {
	cases := []struct {
		a, b []int
		want int
	}{
		{[]int{1, 2, 3, 2, 1}, []int{3, 2, 1, 4, 7}, 3},
		{[]int{0, 0, 0, 0, 0}, []int{0, 0, 0, 0, 0}, 5},
		{[]int{1, 2, 3}, []int{4, 5, 6}, 0},
		{[]int{1}, []int{1}, 1},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.FindLength(c.a, c.b) })
		if got != c.want {
			t.Errorf("FindLength(%v, %v) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestFindLengthRandom(t *testing.T) {
	for i := 0; i < 500; i++ {
		la, lb := 1+rand.Intn(10), 1+rand.Intn(10)
		a, b := make([]int, la), make([]int, lb)
		for j := range a {
			a[j] = rand.Intn(4)
		}
		for j := range b {
			b[j] = rand.Intn(4)
		}
		want := refFindLength(a, b)
		var got int
		runSafely(t, func() { got = lc.FindLength(a, b) })
		if got != want {
			t.Fatalf("FindLength 随机 = %d, want %d", got, want)
		}
	}
}

// == 132 分割回文串 II(最少分割次数) ==

func refMinCut(s string) int {
	n := len(s)
	isPal := make([][]bool, n)
	for i := range isPal {
		isPal[i] = make([]bool, n)
	}
	for i := 0; i < n; i++ {
		isPal[i][i] = true
	}
	for length := 2; length <= n; length++ {
		for i := 0; i+length <= n; i++ {
			j := i + length - 1
			if s[i] == s[j] && (length == 2 || isPal[i+1][j-1]) {
				isPal[i][j] = true
			}
		}
	}
	// dp[i] = 前 i 个字符的最少分割
	const inf = 1 << 30
	dp := make([]int, n+1)
	dp[0] = -1
	for i := 1; i <= n; i++ {
		dp[i] = inf
		for j := 0; j < i; j++ {
			if isPal[j][i-1] && dp[j]+1 < dp[i] {
				dp[i] = dp[j] + 1
			}
		}
	}
	return dp[n]
}

func TestMinCut(t *testing.T) {
	cases := []struct {
		s    string
		want int
	}{
		{"aab", 1},
		{"a", 0},
		{"ab", 1},
		{"aaa", 0},
		{"abcba", 0},
		{"leet", 2},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.MinCut(c.s) })
		if got != c.want {
			t.Errorf("MinCut(%q) = %d, want %d", c.s, got, c.want)
		}
	}
}

func TestMinCutRandom(t *testing.T) {
	for i := 0; i < 800; i++ {
		n := 1 + rand.Intn(10)
		b := make([]byte, n)
		for j := range b {
			b[j] = byte('a' + rand.Intn(3))
		}
		s := string(b)
		want := refMinCut(s)
		var got int
		runSafely(t, func() { got = lc.MinCut(s) })
		if got != want {
			t.Fatalf("MinCut(%q) = %d, want %d", s, got, want)
		}
	}
}
