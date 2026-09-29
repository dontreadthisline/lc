package lc_tests

import (
	"demo/lc"
	"math/rand"
	"slices"
	"testing"
)

// 覆盖 lc/003.go 中除 NoRepeateLongestSubStr（已有 003_test.go）外的全部滑窗解法。
// 用例来源：LeetCode 官方示例 + 针对收缩边界（i<j vs i<=j）与"无解"哨兵的边界用例。

// 209. 长度最小的子数组
func TestMinSubArrayLen(t *testing.T) {
	cases := []struct {
		name   string
		target int
		nums   []int
		want   int
	}{
		{"官方用例1", 7, []int{2, 3, 1, 2, 4, 3}, 2},
		{"官方用例2", 4, []int{1, 4, 4}, 1},
		{"官方用例3_无解", 11, []int{1, 1, 1, 1, 1, 1, 1, 1}, 0},
		{"总和恰好等于target", 15, []int{1, 2, 3, 4, 5}, 5},
		{"无解时必须返回0", 100, []int{1, 2, 3}, 0},
		{"单个元素即解", 5, []int{5}, 1},
		{"首元素即解", 1, []int{1, 2, 3}, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := lc.MinSubArrayLen(c.target, c.nums); got != c.want {
				t.Errorf("MinSubArrayLen(%d, %v) = %d, want %d", c.target, c.nums, got, c.want)
			}
		})
	}
}

// 438. 找到字符串中所有字母异位词
func TestFindAnagrams(t *testing.T) {
	cases := []struct {
		name string
		s    string
		p    string
		want []int
	}{
		{"官方用例1", "cbaebabacd", "abc", []int{0, 6}},
		{"官方用例2_重叠解", "abab", "ab", []int{0, 1, 2}},
		{"s与p相同", "a", "a", []int{0}},
		{"解在中间", "aab", "ab", []int{1}},
		{"p比s长", "ab", "abc", []int{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := lc.FindAnagrams(c.s, c.p)
			if !slices.Equal(got, c.want) {
				t.Errorf("FindAnagrams(%q, %q) = %v, want %v", c.s, c.p, got, c.want)
			}
		})
	}
}

// 904. 水果成篮
func TestTotalFruit(t *testing.T) {
	cases := []struct {
		name   string
		fruits []int
		want   int
	}{
		{"官方用例1", []int{1, 2, 1}, 3},
		{"官方用例2", []int{0, 1, 2, 2}, 3},
		{"官方用例3", []int{1, 2, 3, 2, 2}, 4},
		{"官方用例4", []int{3, 3, 3, 1, 2, 1, 1, 2, 3, 3, 4}, 5},
		{"单元素", []int{1}, 1},
		{"只有一种水果", []int{5, 5, 5}, 3},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := lc.TotalFruit(c.fruits); got != c.want {
				t.Errorf("TotalFruit(%v) = %d, want %d", c.fruits, got, c.want)
			}
		})
	}
}

// 424. 替换后的最长重复字符（两个实现共用同一组用例）
func TestCharacterReplacement(t *testing.T) {
	cases := []struct {
		name string
		s    string
		k    int
		want int
	}{
		{"官方用例1", "ABAB", 2, 4},
		{"官方用例2", "AABABBA", 1, 4},
		{"k为0_全相同", "AAAA", 0, 4},
		{"k足以覆盖整个串", "BAAAB", 2, 5},
		{"窗口被少数字符拖住", "ABBB", 2, 4},
		{"k不足以拼接两个字母", "ABCDE", 1, 2},
		{"空串", "", 5, 0},
	}
	impls := []struct {
		name string
		fn   func(string, int) int
	}{
		{"CharacterReplacement", lc.CharacterReplacement},
	}
	for _, impl := range impls {
		for _, c := range cases {
			t.Run(impl.name+"/"+c.name, func(t *testing.T) {
				if got := impl.fn(c.s, c.k); got != c.want {
					t.Errorf("%s(%q, %d) = %d, want %d", impl.name, c.s, c.k, got, c.want)
				}
			})
		}
	}
}

// 1004. 最大连续1的个数 III
func TestLongestOnes(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		{"官方用例1", []int{1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0}, 2, 6},
		{"官方用例2", []int{0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0, 1, 1, 1, 1}, 3, 10},
		{"k为0_取最长连续1", []int{1, 1, 0, 1, 0}, 0, 2},
		{"k为0_单个0", []int{0}, 0, 0},
		{"k为0_全0", []int{0, 0}, 0, 0},
		{"k为0_全1", []int{1, 1, 1}, 0, 3},
		{"k大于0的个数_全窗", []int{0, 0, 0}, 4, 3},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := lc.LongestOnes(c.nums, c.k); got != c.want {
				t.Errorf("LongestOnes(%v, %d) = %d, want %d", c.nums, c.k, got, c.want)
			}
		})
	}
}

// 1493. 删掉一个元素以后全为 1 的最长子数组
func TestLongestSubarray(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want int
	}{
		{"官方用例1", []int{1, 1, 0, 1}, 3},
		{"官方用例2", []int{0, 1, 1, 1, 0, 1, 1, 0, 1}, 5},
		{"官方用例3_全1必须删一个", []int{1, 1, 1}, 2},
		{"全0", []int{0, 0, 0}, 0},
		{"单个0", []int{0}, 0},
		{"单个1", []int{1}, 0},
		{"0夹在中间", []int{1, 1, 0, 1, 1, 1}, 5},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := lc.LongestSubarray(c.nums); got != c.want {
				t.Errorf("LongestSubarray(%v) = %d, want %d", c.nums, got, c.want)
			}
		})
	}
}

// 713. 乘积小于 K 的子数组
func TestNumSubarrayProductLessThanK(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		{"官方用例1", []int{10, 5, 2, 6}, 100, 8},
		{"官方用例2_k为0", []int{1, 2, 3}, 0, 0},
		{"k为1_乘积1不小于1", []int{1, 1, 1}, 1, 0},
		{"单个元素恰好等于k", []int{5}, 5, 0},
		{"k大于全部乘积", []int{1, 2, 3}, 10, 6},
		{"收缩中途退出", []int{2, 2, 2}, 6, 5},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := lc.NumSubarrayProductLessThanK(c.nums, c.k); got != c.want {
				t.Errorf("NumSubarrayProductLessThanK(%v, %d) = %d, want %d", c.nums, c.k, got, c.want)
			}
		})
	}
}

// 暴力对拍：O(n^2) 纯枚举所有子数组 vs 滑窗实现，验证 res += (j-i+1) 的计数不重不漏
func bruteForceProductLessThanK(nums []int, k int) int {
	cnt := 0
	for i := 0; i < len(nums); i++ {
		p := 1
		for j := i; j < len(nums); j++ {
			p *= nums[j]
			if p < k {
				cnt++
			}
		}
	}
	return cnt
}

func TestNumSubarrayProductLessThanKBruteCrossCheck(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for iter := 0; iter < 20000; iter++ {
		n := rng.Intn(13)
		nums := make([]int, n)
		for x := range nums {
			nums[x] = 1 + rng.Intn(10)
		}
		k := rng.Intn(301)
		want := bruteForceProductLessThanK(nums, k)
		if got := lc.NumSubarrayProductLessThanK(nums, k); got != want {
			t.Fatalf("iter=%d nums=%v k=%d: 滑窗=%d 暴力=%d", iter, nums, k, got, want)
		}
	}
}

// 930. 和相同的二元子数组（含非负泛化用例）
func TestNumSubarraysWithSum(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		goal int
		want int
	}{
		{"官方用例1", []int{1, 0, 1, 0, 1}, 2, 4},
		{"官方用例2_全零goal为零", []int{0, 0, 0, 0, 0}, 0, 15},
		{"零带来的重叠起点", []int{0, 0, 1}, 1, 3},
		{"goal为零_数全零段", []int{1, 0, 0, 1}, 0, 3},
		{"单元素恰等于goal", []int{1}, 1, 1},
		{"单元素大于goal", []int{5}, 5, 1},
		{"goal超出总和", []int{1, 1, 1}, 5, 0},
		{"单个零", []int{0}, 0, 1},
		{"泛化_非零一元素", []int{2, 1, 2}, 3, 2},
		{"泛化_单元素等于goal且大于减一", []int{5}, 5, 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := lc.NumSubarraysWithSum(c.nums, c.goal); got != c.want {
				t.Errorf("NumSubarraysWithSum(%v, %d) = %d, want %d", c.nums, c.goal, got, c.want)
			}
		})
	}
}

// 暴力对拍：O(n^2) 纯枚举 vs atMost 相减实现
func bruteForceSubarraysWithSum(nums []int, goal int) int {
	cnt := 0
	for i := 0; i < len(nums); i++ {
		s := 0
		for j := i; j < len(nums); j++ {
			s += nums[j]
			if s == goal {
				cnt++
			}
		}
	}
	return cnt
}

func TestNumSubarraysWithSumBruteCrossCheck(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	for iter := 0; iter < 20000; iter++ {
		n := rng.Intn(13)
		nums := make([]int, n)
		for x := range nums {
			nums[x] = rng.Intn(4)
		}
		goal := rng.Intn(7)
		want := bruteForceSubarraysWithSum(nums, goal)
		if got := lc.NumSubarraysWithSum(nums, goal); got != want {
			t.Fatalf("iter=%d nums=%v goal=%d: 实现=%d 暴力=%d", iter, nums, goal, got, want)
		}
	}
}

// 395. 至少有 K 个重复字符的最长子串
func TestLongestSubstring(t *testing.T) {
	cases := []struct {
		name string
		s    string
		k    int
		want int
	}{
		{"官方用例1", "aaabb", 3, 3},
		{"官方用例2", "ababbc", 2, 5},
		{"两段各成合法块_整串也合法", "aaabbb", 3, 6},
		{"最优解需要等待而非收缩", "aabbcc", 2, 6},
		{"单字符_k为1", "a", 1, 1},
		{"单字符_k超过出现次数", "a", 2, 0},
		{"k为1_整串即合法", "abcd", 1, 4},
		{"空串", "", 5, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := lc.LongestSubstring(c.s, c.k); got != c.want {
				t.Errorf("LongestSubstring(%q, %d) = %d, want %d", c.s, c.k, got, c.want)
			}
		})
	}
}
