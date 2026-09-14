package algo_tests

import "demo/algo"

import (
	"math/rand"
	"sort"
	"testing"
)

func TestFindKthSmallest(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for round := range 2000 {
		n := rng.Intn(50) + 1
		nums := make([]int, n)
		for i := range nums {
			//值域压到 10 以内，制造大量重复元素，专测重复值边界
			nums[i] = rng.Intn(10)
		}
		k := rng.Intn(n) + 1

		sorted := append([]int(nil), nums...)
		sort.Ints(sorted)
		want := sorted[k-1]

		got := algo.FindKthSmallest(append([]int(nil), nums...), k)
		if got != want {
			t.Fatalf("round %d: nums=%v k=%d got=%d want=%d", round, nums, k, got, want)
		}
	}
}

func TestFindKthLargest(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for round := range 2000 {
		n := rng.Intn(50) + 1
		nums := make([]int, n)
		for i := range nums {
			nums[i] = rng.Intn(10)
		}
		k := rng.Intn(n) + 1

		sorted := append([]int(nil), nums...)
		sort.Ints(sorted)
		//升序下标 n-k 即降序的第 k 个
		want := sorted[n-k]

		got := algo.FindKthLargest(append([]int(nil), nums...), k)
		if got != want {
			t.Fatalf("round %d: nums=%v k=%d got=%d want=%d", round, nums, k, got, want)
		}
	}
}
