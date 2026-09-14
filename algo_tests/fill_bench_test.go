package algo_tests

import (
	"slices"
	"testing"
)

const benchN = 1 << 16 // 65536 个 int，512KB

func BenchmarkFillLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := make([]int, benchN)
		for j := range s {
			s[j] = -1
		}
		sink = s
	}
}

func BenchmarkFillRepeat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sink = slices.Repeat([]int{-1}, benchN)
	}
}

func BenchmarkFillDoubling(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := make([]int, benchN)
		s[0] = -1
		for j := 1; j < benchN; j *= 2 {
			copy(s[j:], s[:j])
		}
		sink = s
	}
}

var sink []int

// 只分配：语言保证返回零值（处女页跳过清零 / 复用页 memclr）
func BenchmarkAllocOnly(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sink = make([]int, benchN)
	}
}

// 复用同一块内存，纯测“Go 层循环填 -1”的成本
func BenchmarkRefillLoop(b *testing.B) {
	s := make([]int, benchN)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := range s {
			s[j] = -1
		}
	}
	sink = s
}

// 复用同一块内存，纯测 clear（memclr）的成本
func BenchmarkRezeroClear(b *testing.B) {
	s := make([]int, benchN)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		clear(s)
	}
	sink = s
}
