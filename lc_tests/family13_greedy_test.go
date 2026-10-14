package lc_tests

import (
	"math/rand"
	"slices"
	"testing"

	"demo/lc"
)

// ---- 32. 最长有效括号 ----

// refLongestValidParentheses: O(n^2) 暴力 —— 每个子串用计数法验证
func validParen(s string) bool {
	bal := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '(' {
			bal++
		} else {
			bal--
			if bal < 0 {
				return false
			}
		}
	}
	return bal == 0
}

func refLongestValidParentheses(s string) int {
	best := 0
	for i := 0; i < len(s); i++ {
		for j := i + 2; j <= len(s); j += 2 {
			if validParen(s[i:j]) && j-i > best {
				best = j - i
			}
		}
	}
	return best
}

func TestLongestValidParentheses(t *testing.T) {
	cases := []struct {
		s    string
		want int
	}{
		{"(()", 2},
		{")()())", 4},
		{"", 0},
		{"()", 2},
		{"()(()", 2},
		{"()(())", 6},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.LongestValidParentheses(c.s) })
		if got != c.want {
			t.Errorf("LongestValidParentheses(%q) = %d, want %d", c.s, got, c.want)
		}
	}
}

func TestLongestValidParenthesesRandom(t *testing.T) {
	for i := 0; i < 2000; i++ {
		n := rand.Intn(20)
		b := make([]byte, n)
		for j := range b {
			if rand.Intn(2) == 0 {
				b[j] = '('
			} else {
				b[j] = ')'
			}
		}
		s := string(b)
		want := refLongestValidParentheses(s)
		var got int
		runSafely(t, func() { got = lc.LongestValidParentheses(s) })
		if got != want {
			t.Fatalf("LongestValidParentheses(%q) = %d, want %d", s, got, want)
		}
	}
}

// ---- 45. 跳跃游戏 II ----

// refJump: BFS 最短路(每个位置一步可达区间为邻居)
func refJump(nums []int) int {
	n := len(nums)
	if n == 1 {
		return 0
	}
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	dist[0] = 0
	queue := []int{0}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for step := 1; step <= nums[cur] && cur+step < n; step++ {
			next := cur + step
			if dist[next] == -1 {
				dist[next] = dist[cur] + 1
				if next == n-1 {
					return dist[next]
				}
				queue = append(queue, next)
			}
		}
	}
	return dist[n-1]
}

func TestJump(t *testing.T) {
	cases := []struct {
		nums []int
		want int
	}{
		{[]int{2, 3, 1, 1, 4}, 2},
		{[]int{2, 3, 0, 1, 4}, 2},
		{[]int{1}, 0},
		{[]int{1, 2}, 1},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.Jump(c.nums) })
		if got != c.want {
			t.Errorf("Jump(%v) = %d, want %d", c.nums, got, c.want)
		}
	}
}

func TestJumpRandom(t *testing.T) {
	for i := 0; i < 3000; i++ {
		n := 1 + rand.Intn(20)
		nums := make([]int, n)
		for j := range nums {
			nums[j] = 1 + rand.Intn(4)
		}
		want := refJump(nums)
		var got int
		runSafely(t, func() { got = lc.Jump(nums) })
		if got != want {
			t.Fatalf("Jump(%v) = %d, want %d(BFS)", nums, got, want)
		}
	}
}

// ---- 55. 跳跃游戏 ----

func refCanJump(nums []int) bool {
	reach := 0
	for i, v := range nums {
		if i > reach {
			return false
		}
		if i+v > reach {
			reach = i + v
		}
	}
	return true
}

func TestCanJump(t *testing.T) {
	cases := []struct {
		nums []int
		want bool
	}{
		{[]int{2, 3, 1, 1, 4}, true},
		{[]int{3, 2, 1, 0, 4}, false},
		{[]int{0}, true},
		{[]int{1, 0}, true},
	}
	for _, c := range cases {
		var got bool
		runSafely(t, func() { got = lc.CanJump(c.nums) })
		if got != c.want {
			t.Errorf("CanJump(%v) = %v, want %v", c.nums, got, c.want)
		}
	}
}

func TestCanJumpRandom(t *testing.T) {
	for i := 0; i < 5000; i++ {
		n := 1 + rand.Intn(20)
		nums := make([]int, n)
		for j := range nums {
			nums[j] = rand.Intn(4)
		}
		want := refCanJump(nums)
		var got bool
		runSafely(t, func() { got = lc.CanJump(nums) })
		if got != want {
			t.Fatalf("CanJump(%v) = %v, want %v", nums, got, want)
		}
	}
}

// ---- 122. 买卖股票的最佳时机 II ----

func refMaxProfitII(prices []int) int {
	profit := 0
	for i := 1; i < len(prices); i++ {
		if d := prices[i] - prices[i-1]; d > 0 {
			profit += d
		}
	}
	return profit
}

func TestMaxProfitII(t *testing.T) {
	cases := []struct {
		prices []int
		want   int
	}{
		{[]int{7, 1, 5, 3, 6, 4}, 7},
		{[]int{1, 2, 3, 4, 5}, 4},
		{[]int{7, 6, 4, 3, 1}, 0},
		{[]int{5}, 0},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.MaxProfitII(c.prices) })
		if got != c.want {
			t.Errorf("MaxProfitII(%v) = %d, want %d", c.prices, got, c.want)
		}
	}
}

func TestMaxProfitIIRandom(t *testing.T) {
	for i := 0; i < 20000; i++ {
		n := 1 + rand.Intn(60)
		prices := make([]int, n)
		for j := range prices {
			prices[j] = rand.Intn(50)
		}
		want := refMaxProfitII(prices)
		var got int
		runSafely(t, func() { got = lc.MaxProfitII(prices) })
		if got != want {
			t.Fatalf("MaxProfitII(%v) = %d, want %d", prices, got, want)
		}
	}
}

// ---- 134. 加油站 ----

func refCanCompleteCircuit(gas, cost []int) int {
	n := len(gas)
	for start := 0; start < n; start++ {
		tank, ok := 0, true
		for d := 0; d < n; d++ {
			i := (start + d) % n
			tank += gas[i] - cost[i]
			if tank < 0 {
				ok = false
				break
			}
		}
		if ok {
			return start
		}
	}
	return -1
}

func TestCanCompleteCircuit(t *testing.T) {
	cases := []struct {
		gas, cost []int
		want      int
	}{
		{[]int{1, 2, 3, 4, 5}, []int{3, 4, 5, 1, 2}, 3},
		{[]int{2, 3, 4}, []int{3, 4, 3}, -1},
		{[]int{5}, []int{5}, 0},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.CanCompleteCircuit(c.gas, c.cost) })
		if got != c.want {
			t.Errorf("CanCompleteCircuit(%v, %v) = %d, want %d", c.gas, c.cost, got, c.want)
		}
	}
}

func TestCanCompleteCircuitRandom(t *testing.T) {
	for i := 0; i < 2000; i++ {
		n := 1 + rand.Intn(10)
		gas := make([]int, n)
		cost := make([]int, n)
		for j := range gas {
			gas[j] = rand.Intn(10)
			cost[j] = rand.Intn(10)
		}
		want := refCanCompleteCircuit(gas, cost)
		var got int
		runSafely(t, func() { got = lc.CanCompleteCircuit(gas, cost) })
		if got != want {
			t.Fatalf("CanCompleteCircuit(%v, %v) = %d, want %d", gas, cost, got, want)
		}
	}
}

// ---- 452. 用最少数量的箭引爆气球 ----

// refFindMinArrowShots: 事件法独立实现(与典型「右端排序贪心」写法不同)
func refFindMinArrowShots(points [][]int) int {
	type ev struct {
		x    int
		kind int // -1 气球开始, +1 气球结束
	}
	events := []ev{}
	for _, p := range points {
		events = append(events, ev{p[0], -1}, ev{p[1], +1})
	}
	// 排序: 结束事件优先于开始事件(触点可同时引爆)
	slices.SortFunc(events, func(a, b ev) int {
		if a.x != b.x {
			return a.x - b.x
		}
		return a.kind - b.kind
	})
	active, arrows := 0, 0
	for _, e := range events {
		if e.kind == -1 {
			active++
		} else {
			if active > 0 {
				arrows++ // 在该结束点射一箭, 清空所有活跃气球
				active = 0
			}
		}
	}
	return arrows
}

func TestFindMinArrowShots(t *testing.T) {
	cases := []struct {
		points [][]int
		want   int
	}{
		{[][]int{{10, 16}, {2, 8}, {1, 6}, {7, 12}}, 2},
		{[][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}}, 4},
		{[][]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}}, 2}, // 触点可同箭
		{[][]int{{1, 1}}, 1},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.FindMinArrowShots(cloneIntervals(c.points)) })
		if got != c.want {
			t.Errorf("FindMinArrowShots(%v) = %d, want %d", c.points, got, c.want)
		}
	}
}

func TestFindMinArrowShotsRandom(t *testing.T) {
	for i := 0; i < 2000; i++ {
		n := 1 + rand.Intn(15)
		points := make([][]int, n)
		for j := range points {
			a := rand.Intn(30)
			b := a + rand.Intn(10)
			points[j] = []int{a, b}
		}
		want := refFindMinArrowShots(points)
		var got int
		runSafely(t, func() { got = lc.FindMinArrowShots(cloneIntervals(points)) })
		if got != want {
			t.Fatalf("FindMinArrowShots(%v) = %d, want %d", points, got, want)
		}
	}
}

// ---- 763. 划分字母区间 ----

// refPartitionLabels: 枚举分割点 2^(n-1), 要求每段字母集两两不相交
func refPartitionLabels(s string) []int {
	n := len(s)
	segMask := func(seg string) int {
		m := 0
		for i := 0; i < len(seg); i++ {
			m |= 1 << (seg[i] - 'a')
		}
		return m
	}
	var best []int
	for mask := 0; mask < 1<<(n-1); mask++ {
		prev, parts := 0, []string{}
		for i := 0; i < n-1; i++ {
			if mask&(1<<i) != 0 {
				parts = append(parts, s[prev:i+1])
				prev = i + 1
			}
		}
		parts = append(parts, s[prev:])
		ok := true
		masks := make([]int, len(parts))
		for i, p := range parts {
			masks[i] = segMask(p)
		}
		for i := 0; i < len(masks) && ok; i++ {
			for j := i + 1; j < len(masks); j++ {
				if masks[i]&masks[j] != 0 {
					ok = false
					break
				}
			}
		}
		if ok {
			sizes := make([]int, len(parts))
			for i, p := range parts {
				sizes[i] = len(p)
			}
			// 题意: 尽可能多的段 —— 分割点最多者优先
			if len(sizes) > len(best) {
				best = sizes
			}
		}
	}
	if best == nil {
		return []int{n}
	}
	return best
}

func TestPartitionLabels(t *testing.T) {
	cases := []struct {
		s    string
		want []int
	}{
		{"ababcbacadefegdehijhklij", []int{9, 7, 8}},
		{"eccbbbbdec", []int{10}},
		{"a", []int{1}},
	}
	for _, c := range cases {
		var got []int
		runSafely(t, func() { got = lc.PartitionLabels(c.s) })
		if !equalInts(got, c.want) {
			t.Errorf("PartitionLabels(%q) = %v, want %v", c.s, got, c.want)
		}
	}
}

func TestPartitionLabelsRandom(t *testing.T) {
	for i := 0; i < 400; i++ {
		n := 1 + rand.Intn(13)
		b := make([]byte, n)
		for j := range b {
			b[j] = byte('a' + rand.Intn(3))
		}
		s := string(b)
		want := refPartitionLabels(s)
		var got []int
		runSafely(t, func() { got = lc.PartitionLabels(s) })
		if !equalInts(got, want) {
			t.Fatalf("PartitionLabels(%q) = %v, want %v(枚举)", s, got, want)
		}
	}
}
