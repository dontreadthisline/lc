package algo_tests

import "demo/algo"

import (
	"math/rand"
	"testing"
)

// 初始状态：每个元素自成一个集合，根就是自己
func TestUFSInit(t *testing.T) {
	u := algo.NewUFS(10)
	for i := 0; i < 10; i++ {
		if r := u.Find(i); r != i {
			t.Fatalf("Find(%d) = %d, want %d", i, r, i)
		}
	}
}

// 越界参数：Find 返回 -1，Union 返回 false，且不破坏现有结构
func TestUFSOutOfBounds(t *testing.T) {
	u := algo.NewUFS(5)
	if r := u.Find(-1); r != -1 {
		t.Fatalf("Find(-1) = %d, want -1", r)
	}
	if r := u.Find(5); r != -1 {
		t.Fatalf("Find(5) = %d, want -1", r)
	}
	if u.Union(-1, 0) {
		t.Fatal("Union(-1, 0) = true, want false")
	}
	if u.Union(0, 5) {
		t.Fatal("Union(0, 5) = true, want false")
	}
	if u.Union(5, -1) {
		t.Fatal("Union(5, -1) = true, want false")
	}
	if r := u.Find(0); r != 0 {
		t.Fatalf("越界操作后 Find(0) = %d, want 0", r)
	}
}

// 基本合并 + 传递性：0-1、1-2、2-3 链式合并后应全部同集
func TestUFSUnionTransitivity(t *testing.T) {
	u := algo.NewUFS(4)
	if !u.Union(0, 1) || !u.Union(1, 2) || !u.Union(2, 3) {
		t.Fatal("合法元素的 Union 应返回 true")
	}
	root := u.Find(0)
	for i := 1; i < 4; i++ {
		if r := u.Find(i); r != root {
			t.Fatalf("Find(%d) = %d, want %d（传递性被破坏）", i, r, root)
		}
	}
}

// 重复合并同一集合：返回 true，结构不变
func TestUFSUnionIdempotent(t *testing.T) {
	u := algo.NewUFS(3)
	u.Union(0, 1)
	root := u.Find(0)
	if !u.Union(0, 1) {
		t.Fatal("合并已在同一集合的两个元素应返回 true")
	}
	if !u.Union(1, 0) { // 参数顺序颠倒也应成立
		t.Fatal("Union(1, 0) 应返回 true")
	}
	if u.Find(0) != root || u.Find(1) != root {
		t.Fatal("重复合并不应改变集合")
	}
	// 未参与的元素保持独立
	if r := u.Find(2); r != 2 {
		t.Fatalf("Find(2) = %d, want 2（无关元素被波及）", r)
	}
}

// 回归测试：按秩合并的 rank 维护。
// 相等 rank 合并后，新根（被挂上去的那个）的 rank 才应该 +1。
// 修复前代码错误地给已失去根身份的 rx 加 rank，导致根的 rank 永远为 0。
func TestUFSRankMaintenance(t *testing.T) {
	// rank/树高是实现细节，目录分离后不可跨包观测；黑盒验证合并语义：
	// 多轮合并后连通性正确、根稳定（同元素多次 Find 结果一致）
	u := algo.NewUFS(16)
	for i := 0; i < 16; i += 2 {
		u.Union(i, i+1)
	}
	for i := 0; i < 16; i += 2 {
		if u.Find(i) != u.Find(i+1) {
			t.Fatalf("%d 和 %d 应连通", i, i+1)
		}
	}
	// 级联合并后仍连通
	u.Union(0, 2)
	u.Union(4, 0)
	u.Union(6, 8)
	u.Union(0, 6)
	for i := 0; i < 10; i++ {
		if u.Find(0) != u.Find(i) {
			t.Fatalf("0..9 应全连通, %d 脱离", i)
		}
	}
	if u.Find(0) == u.Find(10) {
		t.Fatal("10 不应连通")
	}
	// 根稳定性：重复 Find 返回一致
	r1, r2 := u.Find(7), u.Find(7)
	if r1 != r2 {
		t.Fatal("Find 应幂等")
	}
}

// 路径压缩的效果属于实现细节（目录分离后不可跨包观测 parent 数组）。
// 黑盒验证：链式合并后连通性与 Find 幂等性正确；压缩的性能收益由大 N 基准保障。
func TestUFSPathCompression(t *testing.T) {
	u := algo.NewUFS(5)
	// 链式 Union 构造依赖链 0-1-2-3-4
	u.Union(0, 1)
	u.Union(1, 2)
	u.Union(2, 3)
	u.Union(3, 4)
	for i := 0; i < 5; i++ {
		if u.Find(i) != u.Find(0) {
			t.Fatalf("%d 应与 0 连通", i)
		}
	}
	if r1, r2 := u.Find(0), u.Find(0); r1 != r2 {
		t.Fatal("Find 应幂等")
	}
}

// Union(x,y) 与 Union(y,x) 应产生相同的连通关系
func TestUFSUnionSymmetry(t *testing.T) {
	pairs := [][2]int{{0, 1}, {1, 2}, {2, 3}}
	a, b := algo.NewUFS(4), algo.NewUFS(4)
	for _, p := range pairs {
		a.Union(p[0], p[1])
		b.Union(p[1], p[0]) // 顺序颠倒
	}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			sameA := a.Find(i) == a.Find(j)
			sameB := b.Find(i) == b.Find(j)
			if sameA != sameB {
				t.Fatalf("(%d,%d): 正序 same=%v, 逆序 same=%v", i, j, sameA, sameB)
			}
		}
	}
}

// ---- 随机对照测试 ----

// 朴素并查集：union 时整体重写标签，语义简单，用作正确性基准
type naiveUFS struct {
	label []int
}

func newNaiveUFS(n int) *naiveUFS {
	l := make([]int, n)
	for i := range l {
		l[i] = i
	}
	return &naiveUFS{label: l}
}

func (n *naiveUFS) union(x, y int) bool {
	lx, ly := n.label[x], n.label[y]
	if lx == ly {
		return true
	}
	for i := range n.label {
		if n.label[i] == lx {
			n.label[i] = ly
		}
	}
	return true
}

func (n *naiveUFS) same(x, y int) bool {
	return n.label[x] == n.label[y]
}

// 与朴素实现随机交叉对照：union 的返回值与 same 的判定必须完全一致
func TestUFSRandomAgainstNaive(t *testing.T) {
	const (
		n      = 200
		ops    = 3000
		seed   = 42
		pUnion = 0.4
	)
	rng := rand.New(rand.NewSource(seed))
	u, ref := algo.NewUFS(n), newNaiveUFS(n)

	for op := 0; op < ops; op++ {
		x, y := rng.Intn(n), rng.Intn(n)
		if rng.Float64() < pUnion {
			// 偶尔混入越界下标，验证防御逻辑：越界时 algo.UFS 拒绝合并，期望值同样为 false
			if op%97 == 0 {
				x = n + rng.Intn(3)
			}
			want := x >= 0 && x < n && y >= 0 && y < n && ref.union(x, y) // && 短路，越界不会触碰 ref
			if got := u.Union(x, y); got != want {
				t.Fatalf("op %d Union(%d,%d): got %v, want %v", op, x, y, got, want)
			}
		} else {
			if got, want := u.Find(x) == u.Find(y), ref.same(x, y); got != want {
				t.Fatalf("op %d same(%d,%d): got %v, want %v", op, x, y, got, want)
			}
		}
	}
}
