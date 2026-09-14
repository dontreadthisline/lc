package algo_tests

import "testing"

// CountTo 是 Go 1.23 的函数迭代器：for x := range CountTo(5)
func CountTo(n int) func(yield func(int) bool) {
	return func(yield func(int) bool) {
		for i := 1; i <= n; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

func TestForLoopIdioms(t *testing.T) {
	// ① range int
	sum := 0
	for i := range 10 {
		sum += i
	}
	if sum != 45 {
		t.Errorf("range 10 sum = %d, want 45", sum)
	}
	for range 0 {
		t.Error("range 0 不该执行")
	} // n <= 0 直接零次，不 panic

	// ② range string 得到 rune
	runes := 0
	for _, r := range "héllo" {
		if r < 0 {
			t.Error("rune 不该是负的")
		}
		runes++
	}
	if runes != 5 {
		t.Errorf("runes = %d, want 5（len 是 6 字节）", runes)
	}

	// ③ gcd 状态机
	a, b := 48, 18
	for b != 0 {
		a, b = b, a%b
	}
	if a != 6 {
		t.Errorf("gcd = %d, want 6", a)
	}

	// ④ fib：range int + 元组赋值
	x, y := 0, 1
	for range 10 {
		x, y = y, x+y
	}
	if x != 55 {
		t.Errorf("fib10 = %d, want 55", x)
	}

	// ⑤ 就地过滤，零分配
	s := []int{1, 2, 3, 4, 5, 6}
	dst := s[:0]
	for _, v := range s {
		if v%2 == 0 {
			dst = append(dst, v)
		}
	}
	if len(dst) != 3 || dst[0] != 2 || dst[2] != 6 {
		t.Errorf("filter = %v, want [2 4 6]", dst)
	}

	// ⑥ 函数迭代器
	last := 0
	for v := range CountTo(5) {
		last = v
	}
	if last != 5 {
		t.Errorf("iterator last = %d, want 5", last)
	}
}
