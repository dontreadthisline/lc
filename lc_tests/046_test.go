package lc_tests

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"demo/lc"
)

// refPermute: 参照全排列生成器(used 数组回溯)
func refPermute(nums []int) [][]int {
	n := len(nums)
	used := make([]bool, n)
	cur := make([]int, 0, n)
	var res [][]int
	var dfs func()
	dfs = func() {
		if len(cur) == n {
			res = append(res, append([]int(nil), cur...))
			return
		}
		for i := 0; i < n; i++ {
			if used[i] {
				continue
			}
			used[i] = true
			cur = append(cur, nums[i])
			dfs()
			cur = cur[:len(cur)-1]
			used[i] = false
		}
	}
	dfs()
	return res
}

// permKey: 规范化结果(行内序列有义, 行间无序 -> 序列化后排序)
func permKey(res [][]int) string {
	keys := make([]string, len(res))
	for i, row := range res {
		keys[i] = fmt.Sprint(row)
	}
	sort.Strings(keys)
	return fmt.Sprint(keys)
}

func TestPermute(t *testing.T) {
	cases := [][]int{
		{1, 2, 3},
		{0, 1},
		{1},
		{-5, 10},
	}
	for _, nums := range cases {
		want := permKey(refPermute(nums))
		if got := permKey(lc.Permute(nums)); got != want {
			t.Errorf("Permute(%v)\n got  %s\n want %s", nums, got, want)
		}
	}
}

func TestPermuteRandom(t *testing.T) {
	for i := 0; i < 2000; i++ {
		n := 1 + rand.Intn(6)
		vals := randomShuffled(n)
		for j := range vals {
			vals[j] = vals[j] - 10 // 值域平移到 [-9,..] 随机化
		}
		want := permKey(refPermute(vals))
		if got := permKey(lc.Permute(vals)); got != want {
			t.Fatalf("Permute(%v)\n got  %s\n want %s", vals, got, want)
		}
	}
}

func TestFactorial(t *testing.T) {
	cases := []struct{ n, want int }{{0, 1}, {1, 1}, {5, 120}, {6, 720}}
	for _, c := range cases {
		if got := lc.Factorial(c.n); got != c.want {
			t.Errorf("Factorial(%d) = %d, want %d", c.n, got, c.want)
		}
	}
}
