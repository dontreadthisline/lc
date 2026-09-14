package algo_tests

import "demo/algo"

import (
	"math/rand"
	"testing"
)

// 暴力 oracle：每轮线性扫前 candidate 和后 candidate，取最小（平局取下标小）
func pickMinCostWorkersNaive(costs []int, k, candidate int) int64 {
	a := append([]int(nil), costs...)
	var res int64
	for range k {
		n := len(a)
		w := candidate
		if 2*w > n {
			w = n // 窗口重叠时全体可选
		}
		best := 0
		for i := 1; i < n; i++ {
			inWin := i < w || i >= n-w
			bestIn := best < w || best >= n-w
			if inWin && (!bestIn || a[i] < a[best]) {
				best = i
			}
		}
		res += int64(a[best])
		a = append(a[:best], a[best+1:]...)
	}
	return res
}

func TestPickMinCostWorkers(t *testing.T) {
	// 手算反例集：每一题当年都抓过一个 bug
	cases := []struct {
		name    string
		costs   []int
		k, cand int
		want    int64
	}{
		{"例1", []int{17, 12, 10, 2, 7, 2, 11, 20, 8}, 3, 4, 11},
		{"例2", []int{1, 2, 4, 1}, 3, 3, 4},
		{"右种子计数", []int{5, 9, 9, 9, 1, 5}, 1, 2, 1}, // 右窗口少推 → 1@4 饿死
		{"左种子多推", []int{10, 1, 5, 5, 5}, 1, 1, 5},   // 1@1 不在前1窗口
		{"中间人饿死", []int{9, 9, 1, 9}, 3, 2, 19},      // 汇合点 l==r
		{"candidate=1右空", []int{100, 1}, 1, 1, 1},   // 右窗口推 0 个
		{"单元素", []int{42}, 1, 1, 42},                // n=1
		{"全雇", []int{9, 9, 1, 9, 9}, 5, 2, 37},      // k=n
	}
	for _, c := range cases {
		if got := algo.PickMinCostWorkers(c.costs, c.k, c.cand); got != c.want {
			t.Errorf("%s: got %d, want %d", c.name, got, c.want)
		}
	}

	// 随机交叉验证 vs 暴力 oracle
	r := rand.New(rand.NewSource(2024))
	for round := 0; round < 500; round++ {
		n := 1 + r.Intn(30)
		costs := make([]int, n)
		for i := range costs {
			costs[i] = r.Intn(20) // 小值域制造平局，专打 tie-break
		}
		k := 1 + r.Intn(n)
		cand := 1 + r.Intn(n)
		want := pickMinCostWorkersNaive(costs, k, cand)
		if got := algo.PickMinCostWorkers(costs, k, cand); got != want {
			t.Fatalf("round %d: costs=%v k=%d cand=%d got %d want %d",
				round, costs, k, cand, got, want)
		}
	}
}
