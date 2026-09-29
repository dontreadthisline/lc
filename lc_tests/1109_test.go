package lc_tests

import (
	"demo/lc"
	"math/rand"
	"testing"
)

// 1109. 航班预订统计（差分）
func TestCorpFlightBookings(t *testing.T) {
	cases := []struct {
		name     string
		bookings [][]int
		n        int
		want     []int
	}{
		{"官方用例1", [][]int{{1, 2, 10}, {2, 3, 20}, {2, 5, 25}}, 5, []int{10, 55, 45, 25, 25}},
		{"官方用例2", [][]int{{1, 2, 10}, {2, 2, 15}}, 2, []int{10, 25}},
		{"单站预订", [][]int{{3, 3, 7}}, 4, []int{0, 0, 7, 0}},
		{"区间到最后一站_last+1越界", [][]int{{1, 5, 10}}, 5, []int{10, 10, 10, 10, 10}},
		{"宽区间中间站平", [][]int{{2, 4, 10}, {3, 5, 20}}, 5, []int{0, 10, 30, 30, 20}},
		{"重叠同起点", [][]int{{2, 5, 10}, {2, 5, 10}}, 5, []int{0, 20, 20, 20, 20}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := lc.CorpFlightBookings(c.bookings, c.n)
			for i := range got {
				if got[i] != c.want[i] {
					t.Errorf("CorpFlightBookings(%v, %d) = %v, want %v", c.bookings, c.n, got, c.want)
					return
				}
			}
		})
	}
}

// 暴力对拍：每站扫全部预订
func bruteBookings(bookings [][]int, n int) []int {
	res := make([]int, n)
	for _, b := range bookings {
		for i := b[0]; i <= b[1]; i++ {
			res[i-1] += b[2]
		}
	}
	return res
}

func TestCorpFlightBookingsBruteCrossCheck(t *testing.T) {
	rng := rand.New(rand.NewSource(11))
	for iter := 0; iter < 20000; iter++ {
		n := 1 + rng.Intn(10)
		m := rng.Intn(8)
		bookings := make([][]int, m)
		for i := range bookings {
			a := 1 + rng.Intn(n)
			b := a + rng.Intn(n-a+1)
			bookings[i] = []int{a, b, rng.Intn(50)}
		}
		want := bruteBookings(bookings, n)
		got := lc.CorpFlightBookings(bookings, n)
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("iter=%d bookings=%v n=%d: got=%v want=%v", iter, bookings, n, got, want)
			}
		}
	}
}
