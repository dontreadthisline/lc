package algo_tests

import (
	"math/rand"
	"slices"
	"sort"
	"testing"
)

const sortN = 100_000

// Go 1.0 风格：为实现 sort.Interface 而生的三件套
type intSlice []int

func (s intSlice) Len() int           { return len(s) }
func (s intSlice) Less(i, j int) bool { return s[i] < s[j] }
func (s intSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

var sortSink []int

func benchmarkSort(b *testing.B, sortFn func([]int)) {
	r := rand.New(rand.NewSource(42))
	data := make([]int, sortN)
	for i := range data {
		data[i] = r.Int()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := slices.Clone(data) // 三家都含克隆，公平
		sortFn(s)
		sortSink = s
	}
}

func BenchmarkSortInterface(b *testing.B) { benchmarkSort(b, func(s []int) { sort.Sort(intSlice(s)) }) }
func BenchmarkSortSliceReflect(b *testing.B) {
	benchmarkSort(b, func(s []int) { sort.Slice(s, func(i, j int) bool { return s[i] < s[j] }) })
}
func BenchmarkSlicesSort(b *testing.B) { benchmarkSort(b, slices.Sort[[]int]) }
