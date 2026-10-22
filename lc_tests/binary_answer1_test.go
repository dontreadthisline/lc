package lc_tests

import (
	"math/rand"
	"sort"
	"testing"

	"demo/lc"
)

// ---- 二分答案合集 I (1283/1870/1482/2594/475/1760) ----
// 全部走同一台机器: 答案空间 [lo,hi] + 单调判定 check + 找最小可行

// == 1283 使结果不超过阈值的最小除数 ==

func TestSmallestDivisor(t *testing.T) {
	cases := []struct {
		nums      []int
		threshold int
		want      int
	}{
		{[]int{1, 2, 5, 9}, 6, 5},
		{[]int{44, 22, 33, 11, 1}, 5, 44},
		{[]int{1}, 1, 1},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.SmallestDivisor(c.nums, c.threshold) })
		if got != c.want {
			t.Errorf("SmallestDivisor(%v, %d) = %d, want %d", c.nums, c.threshold, got, c.want)
		}
	}
}

func TestSmallestDivisorRandom(t *testing.T) {
	for i := 0; i < 500; i++ {
		n := 1 + rand.Intn(10)
		nums := make([]int, n)
		for j := range nums {
			nums[j] = 1 + rand.Intn(20)
		}
		threshold := n + rand.Intn(30)
		// 暴力: 枚举除数
		want := -1
		for d := 1; d <= 25; d++ {
			s := 0
			for _, v := range nums {
				s += (v + d - 1) / d
			}
			if s <= threshold {
				want = d
				break
			}
		}
		if want == -1 {
			continue // 随机数据无可行解, 跳过
		}
		var got int
		runSafely(t, func() { got = lc.SmallestDivisor(nums, threshold) })
		if got != want {
			t.Fatalf("SmallestDivisor(%v, %d) = %d, want %d", nums, threshold, got, want)
		}
	}
}

// == 1870 准时到达的列车最小时速 ==

func TestMinSpeedOnTime(t *testing.T) {
	cases := []struct {
		dist []int
		hour float64
		want int
	}{
		{[]int{1, 3, 2}, 6, 1},
		{[]int{1, 3, 2}, 2.7, 3},
		{[]int{1, 3, 2}, 1.9, -1},
		// 浮点边界活标本: 正向模拟语义下答案恰好 10^7(total==hour 恰好相等, IEEE 确定);
		// 官方 expected = -1 ⟹ 官方判题器必然不是纯正向 double 语义(逻辑排他),
		// 具体机制(反向推导/整数化截断/恰等判负)未经实证确认。
		// 本测试锁定正向模拟语义 → 10^7; 待 LeetCode 提交实证后回填。
		{[]int{1, 1, 100000}, 2.01, 10000000},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.MinSpeedOnTime(c.dist, c.hour) })
		if got != c.want {
			t.Errorf("MinSpeedOnTime(%v, %v) = %d, want %d", c.dist, c.hour, got, c.want)
		}
	}
}

// == 1482 制作 m 束花所需的最少天数 ==

func TestMinDays(t *testing.T) {
	cases := []struct {
		bloomDay []int
		m, k     int
		want     int
	}{
		{[]int{1, 10, 3, 10, 2}, 3, 1, 3},
		{[]int{1, 10, 3, 10, 2}, 3, 2, -1},
		{[]int{7, 7, 7, 7, 12, 7, 7}, 2, 3, 12},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.MinDays(c.bloomDay, c.m, c.k) })
		if got != c.want {
			t.Errorf("MinDays(%v, %d, %d) = %d, want %d", c.bloomDay, c.m, c.k, got, c.want)
		}
	}
}

func TestMinDaysRandom(t *testing.T) {
	for i := 0; i < 500; i++ {
		n := 1 + rand.Intn(10)
		bloom := make([]int, n)
		for j := range bloom {
			bloom[j] = 1 + rand.Intn(8)
		}
		m := 1 + rand.Intn(3)
		k := 1 + rand.Intn(3)
		// 暴力: 枚举天数
		want := -1
		for d := 1; d <= 8; d++ {
			bouquets, streak := 0, 0
			for _, b := range bloom {
				if b <= d {
					streak++
					if streak == k {
						bouquets++
						streak = 0
					}
				} else {
					streak = 0
				}
			}
			if bouquets >= m {
				want = d
				break
			}
		}
		var got int
		runSafely(t, func() { got = lc.MinDays(bloom, m, k) })
		if got != want {
			t.Fatalf("MinDays(%v, %d, %d) = %d, want %d", bloom, m, k, got, want)
		}
	}
}

// == 2594 修车的最少时间 ==

func TestRepairCars(t *testing.T) {
	cases := []struct {
		ranks []int
		cars  int
		want  int64
	}{
		{[]int{4, 2, 3, 1}, 10, 16},
		{[]int{5, 1, 8}, 6, 16},
		{[]int{1}, 1, 1},
	}
	for _, c := range cases {
		var got int64
		runSafely(t, func() { got = lc.RepairCars(c.ranks, c.cars) })
		if got != c.want {
			t.Errorf("RepairCars(%v, %d) = %d, want %d", c.ranks, c.cars, got, c.want)
		}
	}
}

func TestRepairCarsRandom(t *testing.T) {
	for i := 0; i < 300; i++ {
		n := 1 + rand.Intn(5)
		ranks := make([]int, n)
		for j := range ranks {
			ranks[j] = 1 + rand.Intn(5)
		}
		cars := 1 + rand.Intn(12)
		// 小数据暴力: 枚举时间
		var want int64 = -1
		for t := int64(1); t <= 800; t++ {
			capSum := 0
			for _, r := range ranks {
				// r * c^2 <= t → c <= sqrt(t/r)
				c := 0
				for int64(r)*int64(c+1)*int64(c+1) <= t {
					c++
				}
				capSum += c
			}
			if capSum >= cars {
				want = t
				break
			}
		}
		if want == -1 {
			continue
		}
		var got int64
		runSafely(t, func() { got = lc.RepairCars(ranks, cars) })
		if got != want {
			t.Fatalf("RepairCars(%v, %d) = %d, want %d", ranks, cars, got, want)
		}
	}
}

// == 475 供暖器 ==

func TestFindRadius(t *testing.T) {
	cases := []struct {
		houses, heaters []int
		want            int
	}{
		{[]int{1, 2, 3}, []int{2}, 1},
		{[]int{1, 2, 3, 4}, []int{1, 4}, 1},
		{[]int{1, 5}, []int{2}, 3},
		{[]int{1}, []int{1}, 0},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.FindRadius(c.houses, c.heaters) })
		if got != c.want {
			t.Errorf("FindRadius(%v, %v) = %d, want %d", c.houses, c.heaters, got, c.want)
		}
	}
}

func TestFindRadiusRandom(t *testing.T) {
	for i := 0; i < 500; i++ {
		nh := 1 + rand.Intn(8)
		nhe := 1 + rand.Intn(4)
		houses := rand.Perm(40)[:nh]
		heaters := rand.Perm(40)[:nhe]
		sort.Ints(heaters)
		// 暴力: 每房子到最近供暖器距离的最大值
		want := 0
		for _, h := range houses {
			best := 1 << 30
			for _, he := range heaters {
				if d := abs(h - he); d < best {
					best = d
				}
			}
			if best > want {
				want = best
			}
		}
		var got int
		runSafely(t, func() { got = lc.FindRadius(houses, heaters) })
		if got != want {
			t.Fatalf("FindRadius(%v, %v) = %d, want %d", houses, heaters, got, want)
		}
	}
}

// == 1760 袋子里最少数目的球 (min-max) ==

func TestMinimumSize(t *testing.T) {
	cases := []struct {
		nums          []int
		maxOperations int
		want          int
	}{
		{[]int{9}, 2, 3},
		{[]int{2, 4, 8, 2}, 4, 2},
		{[]int{7, 17}, 2, 7},
	}
	for _, c := range cases {
		var got int
		runSafely(t, func() { got = lc.MinimumSize(c.nums, c.maxOperations) })
		if got != c.want {
			t.Errorf("MinimumSize(%v, %d) = %d, want %d", c.nums, c.maxOperations, got, c.want)
		}
	}
}

func TestMinimumSizeRandom(t *testing.T) {
	for i := 0; i < 500; i++ {
		n := 1 + rand.Intn(8)
		nums := make([]int, n)
		for j := range nums {
			nums[j] = 1 + rand.Intn(15)
		}
		maxOp := rand.Intn(8)
		// 暴力: 枚举上限 T, check = Σ(ceil(b/T)-1) ≤ maxOp
		want := -1
		for T := 1; T <= 15; T++ {
			ops := 0
			for _, b := range nums {
				ops += (b + T - 1) / T - 1
			}
			if ops <= maxOp {
				want = T
				break
			}
		}
		if want == -1 {
			continue
		}
		var got int
		runSafely(t, func() { got = lc.MinimumSize(nums, maxOp) })
		if got != want {
			t.Fatalf("MinimumSize(%v, %d) = %d, want %d", nums, maxOp, got, want)
		}
	}
}
