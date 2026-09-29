package lc_tests

import (
	"demo/lc"
	"math/rand"
	"testing"
)

// 1854. 人口最多的年份（差分，区间 [birth, death)）
func TestMaximumPopulation(t *testing.T) {
	cases := []struct {
		name string
		logs [][]int
		want int
	}{
		{"官方用例1", [][]int{{1993, 1999}, {2000, 2010}}, 1993},
		{"官方用例2", [][]int{{1950, 1961}, {1960, 1971}, {1970, 1981}}, 1960},
		{"只活一年", [][]int{{1950, 1951}}, 1950},
		{"并列取最早", [][]int{{1950, 1955}, {1955, 1960}}, 1950},
		{"死亡年衔接_半开不算", [][]int{{1950, 1960}, {1960, 1970}}, 1950},
		{"death到2050边界", [][]int{{2000, 2050}}, 2000},
		{"单人到世纪末", [][]int{{1999, 2000}}, 1999},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := lc.MaximumPopulation(c.logs); got != c.want {
				t.Errorf("MaximumPopulation(%v) = %d, want %d", c.logs, got, c.want)
			}
		})
	}
}

// 暴力对拍：逐年数活人，取人口最多且最早
func brutePopulation(logs [][]int) int {
	best, year := 0, 0
	for y := 1950; y <= 2049; y++ {
		alive := 0
		for _, lg := range logs {
			if lg[0] <= y && y < lg[1] {
				alive++
			}
		}
		if alive > best {
			best, year = alive, y
		}
	}
	return year
}

func TestMaximumPopulationBruteCrossCheck(t *testing.T) {
	rng := rand.New(rand.NewSource(17))
	for iter := 0; iter < 20000; iter++ {
		m := 1 + rng.Intn(8)
		logs := make([][]int, m)
		for i := range logs {
			b := 1950 + rng.Intn(100)
			logs[i] = []int{b, b + 1 + rng.Intn(2050-b)}
		}
		if got, want := lc.MaximumPopulation(logs), brutePopulation(logs); got != want {
			t.Fatalf("iter=%d logs=%v: 差分=%d 暴力=%d", iter, logs, got, want)
		}
		if got, want := lc.MaximumPopulationBruteForce(logs), brutePopulation(logs); got != want {
			t.Fatalf("iter=%d logs=%v: 仓库暴力=%d 测试暴力=%d", iter, logs, got, want)
		}
	}
}
