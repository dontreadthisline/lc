package lc_tests

import (
	"demo/lc"
	"math/rand"
	"slices"
	"testing"
)

// 179. 最大数（拼接序）
func TestLargestNumberV2(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want string
	}{
		{"官方1", []int{10, 2}, "210"},
		{"官方2", []int{3, 30, 34, 5, 9}, "9534330"},
		{"单元素", []int{10}, "10"},
		{"前缀对1", []int{121, 12}, "12121"},
		{"前缀对2", []int{4, 45}, "454"},
		{"前缀陷阱_单字符不够判", []int{21, 211}, "21211"},
		{"前缀陷阱2", []int{21, 212}, "21221"},
		{"前缀陷阱3", []int{12, 121, 1211}, "121211211"},
		{"含0和正数", []int{0, 1, 2, 3}, "3210"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := lc.LargestNumber(slices.Clone(c.nums))
			if got != c.want {
				t.Errorf("LargestNumber(%v) = %q, want %q", c.nums, got, c.want)
			}
		})
	}
}

// 268. 丢失的数字（原地交换法）
func TestMissingNumber(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want int
	}{
		{"官方1", []int{3, 0, 1}, 2},
		{"官方2", []int{0, 1}, 2},
		{"官方3", []int{9, 6, 4, 2, 3, 5, 7, 0, 1}, 8},
		{"缺零", []int{1}, 0},
		{"缺末尾", []int{0}, 1},
		{"两个元素缺中间", []int{0, 2}, 1},
		{"两个元素缺零", []int{1, 2}, 0},
		{"两个元素缺末尾", []int{0, 1}, 2},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got int
			panicked := safeCall(func() { got = lc.MissingNumber(slices.Clone(c.nums)) })
			if panicked {
				t.Errorf("MissingNumber(%v) panic", c.nums)
				return
			}
			if got != c.want {
				t.Errorf("MissingNumber(%v) = %d, want %d", c.nums, got, c.want)
			}
		})
	}
	// 随机对拍：从 0..n 中随机去掉一个，用求和法做参照
	rng := rand.New(rand.NewSource(61))
	for iter := 0; iter < 10000; iter++ {
		n := rng.Intn(20)
		full := make([]int, n+1)
		for i := range full {
			full[i] = i
		}
		missing := rng.Intn(n + 1)
		nums := slices.Delete(full, missing, missing+1)
		want := n*(n+1)/2 - sum(nums)
		got := lc.MissingNumber(slices.Clone(nums))
		if got != want {
			t.Fatalf("iter=%d nums=%v: got=%d want=%d", iter, nums, got, want)
		}
	}
}

func sum(nums []int) int {
	s := 0
	for _, x := range nums {
		s += x
	}
	return s
}
