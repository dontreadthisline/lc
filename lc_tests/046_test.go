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

// ---- 三种实现的性能对照（回溯版本仅存在于测试包, 不侵入 lc/046.go）----

// permuteBacktrack: 经典回溯(used 数组 + 选择/撤销)
func permuteBacktrack(nums []int) [][]int {
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

// permuteSwap: 原地交换回溯(cur 就是 nums 自身, 零路径数组)
func permuteSwap(nums []int) [][]int {
	var res [][]int
	var dfs func(k int, cur []int)
	dfs = func(k int, cur []int) {
		if k == len(cur) {
			res = append(res, append([]int(nil), cur...))
			return
		}
		for i := k; i < len(cur); i++ {
			cur[k], cur[i] = cur[i], cur[k]
			dfs(k+1, cur)
			cur[k], cur[i] = cur[i], cur[k]
		}
	}
	dfs(0, nums)
	return res
}

var permSink int

func benchmarkPermute(b *testing.B, fn func([]int) [][]int, n int) {
	nums := randomShuffled(n)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := fn(nums)
		permSink += len(res)
	}
}

func BenchmarkPermuteInsert8(b *testing.B)     { benchmarkPermute(b, lc.Permute, 8) }
func BenchmarkPermuteBacktrack8(b *testing.B)  { benchmarkPermute(b, permuteBacktrack, 8) }
func BenchmarkPermuteSwap8(b *testing.B)       { benchmarkPermute(b, permuteSwap, 8) }
func BenchmarkPermuteInsert9(b *testing.B)     { benchmarkPermute(b, lc.Permute, 9) }
func BenchmarkPermuteBacktrack9(b *testing.B)  { benchmarkPermute(b, permuteBacktrack, 9) }
func BenchmarkPermuteSwap9(b *testing.B)       { benchmarkPermute(b, permuteSwap, 9) }

// permuteBacktrackPre: 回溯 + 结果集精确预分配(n! 可预知) + 叶子显式 copy
func permuteBacktrackPre(nums []int) [][]int {
	n := len(nums)
	res := make([][]int, 0, lc.Factorial(n))
	used := make([]bool, n)
	cur := make([]int, 0, n)
	var dfs func()
	dfs = func() {
		if len(cur) == n {
			row := make([]int, n)
			copy(row, cur)
			res = append(res, row)
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

// permuteInsertCopy: 插入法的常数优化 —— 三段 append 换成一次 make + 两次 copy
func permuteInsertCopy(nums []int) [][]int {
	res := [][]int{[]int{}}
	for k := range nums {
		size := len(res[0]) + 1
		tmp := make([][]int, 0, len(res)*size)
		for _, item := range res {
			for i := range size {
				arr := make([]int, size)
				copy(arr, item[:i])
				arr[i] = nums[k]
				copy(arr[i+1:], item[i:])
				tmp = append(tmp, arr)
			}
		}
		res = tmp
	}
	return res
}

func BenchmarkPermuteBacktrackPre8(b *testing.B)  { benchmarkPermute(b, permuteBacktrackPre, 8) }
func BenchmarkPermuteInsertCopy8(b *testing.B)    { benchmarkPermute(b, permuteInsertCopy, 8) }
func BenchmarkPermuteBacktrackPre9(b *testing.B)  { benchmarkPermute(b, permuteBacktrackPre, 9) }
func BenchmarkPermuteInsertCopy9(b *testing.B)    { benchmarkPermute(b, permuteInsertCopy, 9) }
