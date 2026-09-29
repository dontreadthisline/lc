package lc_tests

// demo：泛型二分全家桶 —— slices.BinarySearch / BinarySearchFunc（Go 1.21+）
//
// 这就是标准库里"有泛型的最新接口"的全部：就这两个。
//
//	BinarySearch(s, target)            → 有序类型切片，返回 (插入点, 是否命中)
//	BinarySearchFunc(s, target, cmp)   → 自定义类型/自定义比较，返回 (插入点, 是否命中)
//
// 边界说明：谓词版的"看不见的数组"（875 那种在速度值域上搜的），
// 标准库没有泛型替代——那是你手写 lowerBound 收敛式的领域（lc/004.go 已写）。
//
// —— 大坑：cmp 的参数顺序是 (元素, target)，元素在前；写反就是死循环或漏解 ——

import (
	"cmp"
	"slices"
	"testing"
)

type user struct {
	id   int
	name string
}

func TestSlicesBinaryDemo(t *testing.T) {
	a := []int{1, 3, 5, 7}

	// A. BinarySearch：返回 (插入点, 是否真命中)。
	//    插入点永远有值——它是"第一个 >= target 的位置"，found 负责区分两种含义。
	//    （等于 704/34 题手写版的合体：一个返回值顶你手写的 lowerBound + 一次命中校验）
	pos, found := slices.BinarySearch(a, 5)
	if pos != 2 || !found {
		t.Errorf("BinarySearch 5 = (%d, %v), want (2, true)", pos, found)
	}
	pos, found = slices.BinarySearch(a, 4) // 4 不存在 → 插入点 2，found=false
	if pos != 2 || found {
		t.Errorf("BinarySearch 4 = (%d, %v), want (2, false)", pos, found)
	}
	pos, found = slices.BinarySearch([]int{}, 1) // 空切片 → (0, false)
	if pos != 0 || found {
		t.Errorf("BinarySearch 空切片 = (%d, %v), want (0, false)", pos, found)
	}

	// B. 自定义类型：BinarySearchFunc，cmp(元素, target)，元素在前
	//    返回负数 = 元素应排在 target 前；0 = 命中；正数 = 元素在 target 后
	users := []user{{1, "ann"}, {5, "bob"}, {9, "cid"}}
	i, found := slices.BinarySearchFunc(users, 5, func(u user, target int) int {
		return cmp.Compare(u.id, target)
	})
	if !found || users[i].name != "bob" {
		t.Errorf("BinarySearchFunc id=5 = (%d, %v), want 命中 bob", i, found)
	}
	i, found = slices.BinarySearchFunc(users, 6, func(u user, target int) int {
		return cmp.Compare(u.id, target)
	})
	if found || i != 2 { // 6 不存在 → 插入在 id=5 和 id=9 之间
		t.Errorf("BinarySearchFunc id=6 = (%d, %v), want (2, false)", i, found)
	}

	// C-. upperBound 惯用法：标准库没有独立接口，惯用法是"找 x+1 的插入点"
	//    （整数切片特供；非整数/自定义序用 BinarySearchFunc 或手写收敛式）
	//    "最后一个等于 x 的位置" = upperBound(x) - 1 —— 正是 34 题手写过的技巧
	ub := func(s []int, x int) int {
		p, _ := slices.BinarySearch(s, x+1)
		return p
	}
	nums := []int{5, 7, 7, 8, 8, 10}
	if ub(nums, 8) != 5 || ub(nums, 7) != 3 || ub(nums, 11) != 6 {
		t.Errorf("upperBound 惯用法结果错误: ub(8)=%d ub(7)=%d ub(11)=%d", ub(nums, 8), ub(nums, 7), ub(nums, 11))
	}

	// C. 一致性自检：对同一批 target，插入点必须等于手写 lowerBound 的返回值
	//    （lc/004.go 里的 lowerBound 思路：第一个 >= x 的下标）
	lowerBound := func(s []int, x int) int {
		l, r := 0, len(s)
		for l < r {
			mid := l + (r-l)/2
			if s[mid] >= x {
				r = mid
			} else {
				l = mid + 1
			}
		}
		return l
	}
	for _, x := range []int{0, 1, 2, 4, 5, 6, 7, 8} {
		want := posOf(a, x)
		if got := lowerBound(a, x); got != want {
			t.Errorf("target=%d: 手写 lowerBound=%d, slices.BinarySearch=%d", x, got, want)
		}
	}
}

// posOf：直接取 BinarySearch 的插入点，供 C 段对拍
func posOf(s []int, x int) int {
	p, _ := slices.BinarySearch(s, x)
	return p
}
