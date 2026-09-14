package algo_tests

import "demo/algo"

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"
)

func ages(users []*algo.User) []int {
	out := make([]int, len(users))
	for i, u := range users {
		out[i] = u.Age
	}
	return out
}

// 乱序基本用例：升序 + 稳定性（同龄 id 保持原序）
func TestBucketSort(t *testing.T) {
	users := []*algo.User{
		{Id: 1, Age: 20}, {Id: 2, Age: 18}, {Id: 3, Age: 20},
		{Id: 4, Age: 22}, {Id: 5, Age: 18}, {Id: 6, Age: 20},
	}
	algo.BucketSort(users)
	wantAge := []int{18, 18, 20, 20, 20, 22}
	if got := ages(users); !slices.Equal(got, wantAge) {
		t.Errorf("年龄序 = %v, want %v", got, wantAge)
	}
	// 同龄段内 id 应保持原相对顺序（稳定性）
	for _, seg := range [][2]int{{0, 2}, {2, 5}} {
		for i := seg[0] + 1; i < seg[1]; i++ {
			if users[i-1].Id >= users[i].Id {
				t.Errorf("同龄段 id 乱序: %v", users)
			}
		}
	}
	fmt.Println("基本用例:", ages(users))
}

// 值域上界：algo.MaxAge=120 是合法年龄，必须在值域内
func TestBucketSortBoundary(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Age=120 触发 panic: %v", r)
		}
	}()
	users := []*algo.User{{Id: 1, Age: 120}, {Id: 2, Age: 0}, {Id: 3, Age: 60}}
	algo.BucketSort(users)
	if got := ages(users); !slices.Equal(got, []int{0, 60, 120}) {
		t.Errorf("边界用例 = %v, want [0 60 120]", got)
	}
}

// countSort 返回新切片，不改动 users
func TestCountSort(t *testing.T) {
	users := []*algo.User{
		{Id: 1, Age: 20}, {Id: 2, Age: 18}, {Id: 3, Age: 20},
		{Id: 4, Age: 22}, {Id: 5, Age: 18}, {Id: 6, Age: 20},
	}
	got := algo.CountSort(users)
	wantAge := []int{18, 18, 20, 20, 20, 22}
	ga := make([]int, len(got))
	for i, u := range got {
		ga[i] = u.Age
	}
	if !slices.Equal(ga, wantAge) {
		t.Errorf("年龄序 = %v, want %v", ga, wantAge)
	}
	for _, seg := range [][2]int{{0, 2}, {2, 5}} {
		for i := seg[0] + 1; i < seg[1]; i++ {
			if got[i-1].Id >= got[i].Id {
				t.Errorf("同龄段 id 乱序（稳定性）: %v", ga)
			}
		}
	}
}

// Age=0 是合法值（婴儿），不能 panic
func TestCountSortZeroAge(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Age=0 触发 panic: %v", r)
		}
	}()
	got := algo.CountSort([]*algo.User{{Id: 1, Age: 0}, {Id: 2, Age: 1}, {Id: 3, Age: 0}})
	for i, want := range []int{1, 3, 2} { // 两个 0 岁 id 1,3 在前，然后 1 岁 id 2
		if got[i].Id != want {
			t.Errorf("got id=%d at %d, want %d", got[i].Id, i, want)
		}
	}
}

func TestCountSortRandomCrossCheck(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	for round := 0; round < 300; round++ {
		n := 1 + r.Intn(50)
		users := make([]*algo.User, n)
		for i := range users {
			users[i] = &algo.User{Id: i, Age: r.Intn(algo.MaxAge + 1)} // 含 Age=0
		}
		want := slices.Clone(users)
		slices.SortStableFunc(want, func(a, b *algo.User) int { return a.Age - b.Age })

		got := algo.CountSort(users)
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("round %d: 不稳定或错序", round)
			}
		}
	}
}

// 随机交叉验证：对照标准库稳定排序（按 Age 升序）
func TestBucketSortRandomCrossCheck(t *testing.T) {
	r := rand.New(rand.NewSource(2024))
	for round := 0; round < 300; round++ {
		n := 1 + r.Intn(50)
		users := make([]*algo.User, n)
		for i := range users {
			users[i] = &algo.User{Id: i, Age: r.Intn(algo.MaxAge + 1)} // Age ∈ [0, 120]
		}
		want := slices.Clone(users)
		slices.SortStableFunc(want, func(a, b *algo.User) int { return a.Age - b.Age })

		algo.BucketSort(users)
		for i := range want {
			if users[i] != want[i] { // 指针比对：位置+身份+稳定性全覆盖
				t.Fatalf("round %d: got %+v, want %+v", round, ages(users), ages(want))
			}
		}
	}
}
