package lc_tests

import (
	"math/rand"
	"reflect"
	"testing"

	"demo/lc"
)

func bruteNextGreaterElements(nums []int) []int {
	n := len(nums)
	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = -1
		for d := 1; d < n; d++ {
			if nums[(i+d)%n] > nums[i] {
				res[i] = nums[(i+d)%n]
				break
			}
		}
	}
	return res
}

func TestNextGreaterElements(t *testing.T) {
	cases := []struct {
		nums []int
		want []int
	}{
		{[]int{1, 2, 1}, []int{2, -1, 2}},
		{[]int{1, 2, 3, 4, 3}, []int{2, 3, 4, -1, 4}},
		{[]int{5}, []int{-1}},
		{[]int{2, 2, 2}, []int{-1, -1, -1}},
	}
	for _, c := range cases {
		if got := lc.NextGreaterElements(c.nums); !reflect.DeepEqual(got, c.want) {
			t.Errorf("NextGreaterElements(%v) = %v, want %v", c.nums, got, c.want)
		}
	}
}

func TestNextGreaterElementsRandom(t *testing.T) {
	for i := 0; i < 5000; i++ {
		n := 1 + rand.Intn(60)
		nums := make([]int, n)
		for j := range nums {
			nums[j] = rand.Intn(41) - 20
		}
		want := bruteNextGreaterElements(nums)
		if got := lc.NextGreaterElements(nums); !reflect.DeepEqual(got, want) {
			t.Fatalf("NextGreaterElements(%v) = %v, want %v", nums, got, want)
		}
	}
}
