package lc_tests

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"demo/lc"
)

// ---- 056. 合并区间 ----

// refMerge: 经典合并(排序 + end 取 max)
func refMerge(intervals [][]int) [][]int {
	cp := make([][]int, len(intervals))
	for i, iv := range intervals {
		cp[i] = append([]int(nil), iv...)
	}
	sort.Slice(cp, func(i, j int) bool { return cp[i][0] < cp[j][0] })
	var res [][]int
	for _, iv := range cp {
		if len(res) == 0 || iv[0] > res[len(res)-1][1] {
			res = append(res, append([]int(nil), iv...))
		} else if iv[1] > res[len(res)-1][1] {
			res[len(res)-1][1] = iv[1]
		}
	}
	return res
}

func intervalKey(res [][]int) string {
	keys := make([]string, len(res))
	for i, iv := range res {
		keys[i] = fmt.Sprint(iv)
	}
	sort.Strings(keys)
	return fmt.Sprint(keys)
}

func cloneIntervals(in [][]int) [][]int {
	cp := make([][]int, len(in))
	for i, iv := range in {
		cp[i] = append([]int(nil), iv...)
	}
	return cp
}

func TestMergeRange(t *testing.T) {
	cases := []struct {
		in   [][]int
		want [][]int
	}{
		{[][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}, [][]int{{1, 6}, {8, 10}, {15, 18}}},
		{[][]int{{1, 4}, {4, 5}}, [][]int{{1, 5}}},
		{[][]int{{4, 7}, {1, 4}}, [][]int{{1, 7}}},
		{[][]int{{1}}, [][]int{{1}}},                                  // 防御: 非法输入不炸即可, 期望按原样
		{[][]int{{1, 10}, {2, 3}}, [][]int{{1, 10}}},                   // 大区间吃小区间: end 必须取 max
		{[][]int{{1, 10}, {2, 12}}, [][]int{{1, 12}}},                  // 真扩展
		{[][]int{{5, 6}, {1, 2}, {3, 4}}, [][]int{{1, 2}, {3, 4}, {5, 6}}},
	}
	for _, c := range cases {
		if len(c.in[0]) < 2 { // 跳过非法行
			continue
		}
		if got := intervalKey(lc.MergeRange(cloneIntervals(c.in))); got != intervalKey(c.want) {
			t.Errorf("MergeRange(%v)\n got  %s\n want %s", c.in, got, intervalKey(c.want))
		}
	}
}

func TestMergeRangeRandom(t *testing.T) {
	for i := 0; i < 3000; i++ {
		n := 1 + rand.Intn(25)
		in := make([][]int, n)
		for j := range in {
			a := rand.Intn(21)
			b := a + 1 + rand.Intn(21-a)
			in[j] = []int{a, b}
		}
		want := intervalKey(refMerge(in))
		// MergeRange 会原地排序入参, 必须传副本
		if got := intervalKey(lc.MergeRange(cloneIntervals(in))); got != want {
			t.Fatalf("MergeRange(%v)\n got  %s\n want %s", in, got, want)
		}
	}
}

// ---- 435. 无重叠区间 ----

// overlap: 触点不算重叠([1,2] 与 [2,3] 不重叠)
func intervalsOverlap(a, b []int) bool {
	return max(a[0], b[0]) < min(a[1], b[1])
}

// refErase: 位掩码枚举保留子集, 求最大两两不重叠保留数
func refErase(intervals [][]int) int {
	n := len(intervals)
	best := 0
	for mask := 1; mask < 1<<n; mask++ {
		idx := []int{}
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				idx = append(idx, i)
			}
		}
		ok := true
		for x := 0; x < len(idx) && ok; x++ {
			for y := x + 1; y < len(idx); y++ {
				if intervalsOverlap(intervals[idx[x]], intervals[idx[y]]) {
					ok = false
					break
				}
			}
		}
		if ok && len(idx) > best {
			best = len(idx)
		}
	}
	return n - best
}

func TestEraseOverlapIntervals(t *testing.T) {
	cases := []struct {
		in   [][]int
		want int
	}{
		{[][]int{{1, 2}, {2, 3}, {3, 4}, {1, 3}}, 1},
		{[][]int{{1, 2}, {1, 2}, {1, 2}}, 2},
		{[][]int{{1, 2}, {2, 3}}, 0},
	}
	for _, c := range cases {
		if got := lc.EraseOverlapIntervals(cloneIntervals(c.in)); got != c.want {
			t.Errorf("EraseOverlapIntervals(%v) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestEraseOverlapIntervalsRandom(t *testing.T) {
	for i := 0; i < 3000; i++ {
		n := 1 + rand.Intn(10)
		in := make([][]int, n)
		for j := range in {
			a := rand.Intn(16)
			b := a + 1 + rand.Intn(16-a)
			in[j] = []int{a, b}
		}
		want := refErase(in)
		if got := lc.EraseOverlapIntervals(cloneIntervals(in)); got != want {
			t.Fatalf("EraseOverlapIntervals(%v) = %d, want %d", in, got, want)
		}
	}
}

// ---- 121. 买卖股票的最佳时机 ----

func TestMaxProfit(t *testing.T) {
	cases := []struct {
		prices []int
		want   int
	}{
		{[]int{7, 1, 5, 3, 6, 4}, 5},
		{[]int{7, 6, 4, 3, 1}, 0},
		{[]int{1}, 0},
		{[]int{2, 4, 1}, 2},
		{[]int{3, 3, 3}, 0},
	}
	for _, c := range cases {
		if got := lc.MaxProfit(c.prices); got != c.want {
			t.Errorf("MaxProfit(%v) = %d, want %d", c.prices, got, c.want)
		}
	}
}

func TestMaxProfitRandom(t *testing.T) {
	for i := 0; i < 20000; i++ {
		n := 1 + rand.Intn(80)
		prices := make([]int, n)
		for j := range prices {
			prices[j] = rand.Intn(51)
		}
		want := lc.MaxProfitBruteSearch(prices)
		if got := lc.MaxProfit(prices); got != want {
			t.Fatalf("MaxProfit(%v) = %d, want %d(brute)", prices, got, want)
		}
	}
}
