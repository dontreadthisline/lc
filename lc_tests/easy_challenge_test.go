package lc_tests

import (
	"demo/lc"
	"math/rand"
	"slices"
	"testing"
)

// easy_challenge 五连的验收测试。
// 约定：实现仍是 panic("TODO") 时自动 Skip，不污染全套件；实现后测试自动生效。

func runOrSkip(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Skipf("尚未实现: %v", r)
		}
	}()
	f()
}

// 2908. 元素和最小的山形三元组 I
func TestMinMountainSum(t *testing.T) {
	brute := func(nums []int) int {
		n := len(nums)
		best := 1 << 30
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				for k := j + 1; k < n; k++ {
					if nums[i] < nums[j] && nums[k] < nums[j] {
						best = min(best, nums[i]+nums[j]+nums[k])
					}
				}
			}
		}
		if best == 1<<30 {
			return -1
		}
		return best
	}
	for _, c := range [][2]any{{[]int{8, 6, 1, 5, 3}, 9}, {[]int{5, 4, 8, 7, 10, 2}, 13}, {[]int{6, 5, 4, 3, 4, 5}, -1}} {
		nums, want := c[0].([]int), c[1].(int)
		runOrSkip(t, func() {
			if got := lc.MinMountainSum(nums); got != want {
				t.Errorf("MinMountainSum(%v) = %d, want %d", nums, got, want)
			}
		})
	}
	rng := rand.New(rand.NewSource(51))
	for iter := 0; iter < 5000; iter++ {
		n := 3 + rng.Intn(10)
		nums := make([]int, n)
		for i := range nums {
			nums[i] = 1 + rng.Intn(50)
		}
		want := brute(nums)
		runOrSkip(t, func() {
			if got := lc.MinMountainSum(nums); got != want {
				t.Fatalf("iter=%d nums=%v: got=%d want=%d", iter, nums, got, want)
			}
		})
	}
}

// 2644. 语义：返回能整除 nums 中最多元素的 divisors 之值，并列取最小
func TestMaxDivisibilityScoreDivisor(t *testing.T) {
	brute := func(nums, divisors []int) int {
		bestCnt, bestD := -1, 0
		for _, d := range divisors {
			cnt := 0
			for _, x := range nums {
				if x%d == 0 {
					cnt++
				}
			}
			if cnt > bestCnt || (cnt == bestCnt && d < bestD) {
				bestCnt, bestD = cnt, d
			}
		}
		return bestD
	}
	for _, c := range []struct {
		nums, divisors []int
		want           int
	}{
		{[]int{4, 7, 8, 16}, []int{8, 4, 2}, 2}, // 8:2 个，4:3 个，2:3 个 → 并列取最小 2
		{[]int{4, 7, 8, 16}, []int{8, 4}, 4},    // 8:2，4:3
		{[]int{3}, []int{2}, 2},                 // 全零计数取最小
		{[]int{9, 18, 27}, []int{9, 3, 6}, 3},   // 9:3, 3:3, 6:0 → 并列取 3
		{[]int{100}, []int{100, 50, 25}, 25},    // 三者都整除 100 → 最小 25
	} {
		runOrSkip(t, func() {
			if got := lc.MaxDivisibilityScoreDivisor(c.nums, c.divisors); got != c.want {
				t.Errorf("MaxDivisibilityScoreDivisor(%v, %v) = %d, want %d", c.nums, c.divisors, got, c.want)
			}
		})
	}
	rng := rand.New(rand.NewSource(52))
	for iter := 0; iter < 3000; iter++ {
		nums := make([]int, 1+rng.Intn(8))
		for i := range nums {
			nums[i] = 1 + rng.Intn(30)
		}
		divs := make([]int, 1+rng.Intn(5))
		for i := range divs {
			divs[i] = 1 + rng.Intn(12)
		}
		want := brute(nums, divs)
		runOrSkip(t, func() {
			if got := lc.MaxDivisibilityScoreDivisor(nums, divs); got != want {
				t.Fatalf("iter=%d nums=%v divs=%v: got=%d want=%d", iter, nums, divs, got, want)
			}
		})
	}
}

// 3127. 构造相同颜色的正方形
func TestMakeSameSquare(t *testing.T) {
	ref := func(grid [][]byte) bool {
		for r := 0; r < 2; r++ {
			for c := 0; c < 2; c++ {
				b, w := 0, 0
				for dr := 0; dr < 2; dr++ {
					for dc := 0; dc < 2; dc++ {
						if grid[r+dr][c+dc] == 'B' {
							b++
						} else {
							w++
						}
					}
				}
				if b >= 3 || w >= 3 {
					return true
				}
			}
		}
		return false
	}
	cases := []struct {
		grid [][]byte
		want bool
	}{
		{[][]byte{{'B', 'W', 'B'}, {'B', 'W', 'W'}, {'B', 'W', 'B'}}, true},
		{[][]byte{{'B', 'W', 'B'}, {'B', 'W', 'B'}, {'B', 'W', 'B'}}, false},
		{[][]byte{{'B', 'B', 'B'}, {'B', 'B', 'B'}, {'B', 'B', 'B'}}, true},
		{[][]byte{{'W', 'B', 'W'}, {'B', 'W', 'B'}, {'W', 'B', 'W'}}, false},
	}
	for _, c := range cases {
		runOrSkip(t, func() {
			if got := lc.MakeSameSquare(c.grid); got != c.want {
				t.Errorf("MakeSameSquare = %v, want %v", got, c.want)
			}
		})
	}
	rng := rand.New(rand.NewSource(53))
	for iter := 0; iter < 5000; iter++ {
		grid := make([][]byte, 3)
		for r := range grid {
			grid[r] = make([]byte, 3)
			for c := range grid[r] {
				if rng.Intn(2) == 0 {
					grid[r][c] = 'B'
				} else {
					grid[r][c] = 'W'
				}
			}
		}
		want := ref(grid)
		cp := make([][]byte, 3)
		for r := range grid {
			cp[r] = slices.Clone(grid[r])
		}
		runOrSkip(t, func() {
			if got := lc.MakeSameSquare(cp); got != want {
				t.Fatalf("iter=%d grid=%v: got=%v want=%v", iter, grid, got, want)
			}
		})
	}
}

// 806. 写字符串需要的行数
func TestNumberOfLines806(t *testing.T) {
	ref := func(widths []int, s string) []int {
		lines, cur := 1, 0
		for i := 0; i < len(s); i++ {
			w := widths[s[i]-'a']
			if cur+w > 100 {
				lines++
				cur = w
			} else {
				cur += w
			}
		}
		return []int{lines, cur}
	}
	uniform := make([]int, 26)
	for i := range uniform {
		uniform[i] = 10
	}
	w4 := slices.Clone(uniform)
	w4[0] = 4
	for _, c := range []struct {
		widths []int
		s      string
		want   []int
	}{
		{uniform, "abcdefghijklmnopqrstuvwxyz", []int{3, 60}},
		{w4, "bbbcccdddaaa", []int{2, 4}},
		{uniform, "a", []int{1, 10}},
		{uniform, "aaaaaaaaaaa", []int{2, 10}},
	} {
		runOrSkip(t, func() {
			got := lc.NumberOfLines806(c.widths, c.s)
			if !slices.Equal(got, c.want) {
				t.Errorf("NumberOfLines806(%q) = %v, want %v", c.s, got, c.want)
			}
		})
	}
	rng := rand.New(rand.NewSource(54))
	for iter := 0; iter < 5000; iter++ {
		widths := make([]int, 26)
		for i := range widths {
			widths[i] = 2 + rng.Intn(9)
		}
		n := rng.Intn(30)
		b := make([]byte, n)
		for i := range b {
			b[i] = byte('a' + rng.Intn(26))
		}
		want := ref(widths, string(b))
		runOrSkip(t, func() {
			got := lc.NumberOfLines806(widths, string(b))
			if !slices.Equal(got, want) {
				t.Fatalf("iter=%d s=%q: got=%v want=%v", iter, string(b), got, want)
			}
		})
	}
}

// 3232. 判断是否可以赢得数字游戏（参照：直接模拟对局）
func TestDigitGameWin(t *testing.T) {
	ref := func(nums []int) bool {
		single, multi := 0, 0
		for _, x := range nums {
			if x < 10 {
				single++
			} else {
				multi++
			}
		}
		// 模拟：Alice 只能拿 single，Bob 只能拿 multi，轮到无牌者输，Alice 先手
		turnAlice := true
		for {
			if turnAlice {
				if single == 0 {
					return false
				}
				single--
			} else {
				if multi == 0 {
					return true
				}
				multi--
			}
			turnAlice = !turnAlice
		}
	}
	for _, c := range []struct {
		nums []int
		want bool
	}{
		{[]int{1, 2, 3, 4, 10}, true},    // 官方 1
		{[]int{1, 2, 3, 4, 5, 14}, true}, // 官方 2
		{[]int{5, 5, 5, 25}, true},       // 官方 3
		{[]int{10, 20}, false},           // Alice 无 single 可拿
		{[]int{1, 100}, false},           // 1:1 平，Alice 先耗尽
		{[]int{5}, true},                 // Bob 直接无 multi
		{[]int{100}, false},
	} {
		runOrSkip(t, func() {
			if got := lc.DigitGameWin(c.nums); got != c.want {
				t.Errorf("DigitGameWin(%v) = %v, want %v", c.nums, got, c.want)
			}
		})
	}
	rng := rand.New(rand.NewSource(55))
	for iter := 0; iter < 5000; iter++ {
		nums := make([]int, rng.Intn(10))
		for i := range nums {
			nums[i] = 1 + rng.Intn(99)
		}
		want := ref(nums)
		runOrSkip(t, func() {
			if got := lc.DigitGameWin(nums); got != want {
				t.Fatalf("iter=%d nums=%v: got=%v want=%v", iter, nums, got, want)
			}
		})
	}
}
