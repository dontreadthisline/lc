package lc_tests

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"

	"demo/lc"
)

// ---- 22. 括号生成 ----

func refGenerateParenthesis(n int) []string {
	var res []string
	var gen func(s string, open, close int)
	gen = func(s string, open, close int) {
		if len(s) == 2*n {
			res = append(res, s)
			return
		}
		if open < n {
			gen(s+"(", open+1, close)
		}
		if close < open {
			gen(s+")", open, close+1)
		}
	}
	gen("", 0, 0)
	return res
}

func TestGenerateParenthesis(t *testing.T) {
	cases := []struct {
		n    int
		want int // 结果数量(卡特兰数)
	}{
		{1, 1}, {2, 2}, {3, 5},
	}
	for _, c := range cases {
		var got []string
		runSafely(t, func() { got = lc.GenerateParenthesis(c.n) })
		want := refGenerateParenthesis(c.n)
		if len(got) != c.want || !equalStrSet(got, want) {
			t.Errorf("GenerateParenthesis(%d) = %v, want %v", c.n, got, want)
		}
	}
}

func equalStrSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	as := append([]string(nil), a...)
	bs := append([]string(nil), b...)
	sort.Strings(as)
	sort.Strings(bs)
	for i := range as {
		if as[i] != bs[i] {
			return false
		}
	}
	return true
}

// ---- 77. 组合 ----

func refCombine(n, k int) [][]int {
	var res [][]int
	var cur []int
	var gen func(start int)
	gen = func(start int) {
		if len(cur) == k {
			res = append(res, append([]int(nil), cur...))
			return
		}
		for i := start; i <= n; i++ {
			cur = append(cur, i)
			gen(i + 1)
			cur = cur[:len(cur)-1]
		}
	}
	gen(1)
	return res
}

func TestCombine(t *testing.T) {
	cases := []struct{ n, k, want int }{
		{4, 2, 6}, {1, 1, 1}, {5, 5, 1}, {5, 0, 1}, // C(n,0)=1, 空组合
	}
	for _, c := range cases {
		var got [][]int
		runSafely(t, func() { got = lc.Combine(c.n, c.k) })
		if len(got) != c.want {
			t.Errorf("Combine(%d, %d) 数量 = %d, want %d", c.n, c.k, len(got), c.want)
			continue
		}
		want := refCombine(c.n, c.k)
		if setKey(got) != setKey(want) {
			t.Errorf("Combine(%d, %d) 集合不一致", c.n, c.k)
		}
	}
}

func TestCombineRandom(t *testing.T) {
	for i := 0; i < 300; i++ {
		n := 1 + rand.Intn(12)
		k := rand.Intn(n + 1)
		var got [][]int
		runSafely(t, func() { got = lc.Combine(n, k) })
		want := refCombine(n, k)
		if setKey(got) != setKey(want) {
			t.Fatalf("Combine(%d, %d) 集合不一致", n, k)
		}
	}
}

// ---- 78. 子集 ----

func refSubsets(nums []int) [][]int {
	n := len(nums)
	var res [][]int
	for mask := 0; mask < 1<<n; mask++ {
		cur := []int{}
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				cur = append(cur, nums[i])
			}
		}
		res = append(res, cur)
	}
	return res
}

func TestSubsets(t *testing.T) {
	cases := [][]int{
		{1, 2, 3},
		{0},
		{},
	}
	for _, nums := range cases {
		var got [][]int
		runSafely(t, func() { got = lc.Subsets(nums) })
		fmt.Println(got)
		want := refSubsets(nums)
		if setKey(got) != setKey(want) {
			t.Errorf("Subsets(%v) 集合不一致", nums)
		}
	}
}

func TestSubsetsRandom(t *testing.T) {
	for i := 0; i < 500; i++ {
		n := rand.Intn(11)
		nums := randomShuffled(n + 1)[:n]
		var got [][]int
		runSafely(t, func() { got = lc.Subsets(nums) })
		want := refSubsets(nums)
		if setKey(got) != setKey(want) {
			t.Fatalf("Subsets(%v) 集合不一致", nums)
		}
	}
}

// ---- 79. 单词搜索 ----

// refExist: 独立 DFS 参照(与典型解法写法刻意不同: 复制棋盘标记)
func refExist(board [][]byte, word string) bool {
	m, n := len(board), len(board[0])
	var dfs func(b [][]byte, r, c, idx int) bool
	dfs = func(b [][]byte, r, c, idx int) bool {
		if idx == len(word) {
			return true
		}
		if r < 0 || c < 0 || r >= m || c >= n || b[r][c] != word[idx] {
			return false
		}
		nb := make([][]byte, len(b))
		for i := range b {
			nb[i] = append([]byte(nil), b[i]...)
		}
		nb[r][c] = '#'
		return dfs(nb, r+1, c, idx+1) || dfs(nb, r-1, c, idx+1) || dfs(nb, r, c+1, idx+1) || dfs(nb, r, c-1, idx+1)
	}
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if dfs(board, r, c, 0) {
				return true
			}
		}
	}
	return false
}

func cloneBoard(board [][]byte) [][]byte {
	cp := make([][]byte, len(board))
	for i := range board {
		cp[i] = append([]byte(nil), board[i]...)
	}
	return cp
}

func TestExist(t *testing.T) {
	board := [][]byte{
		{'A', 'B', 'C', 'E'},
		{'S', 'F', 'C', 'S'},
		{'A', 'D', 'E', 'E'},
	}
	cases := []struct {
		word string
		want bool
	}{
		{"ABCCED", true},
		{"SEE", true},
		{"ABCB", false},
	}
	for _, c := range cases {
		var got bool
		runSafely(t, func() { got = lc.Exist(cloneBoard(board), c.word) })
		if got != c.want {
			t.Errorf("Exist(word=%s) = %v, want %v", c.word, got, c.want)
		}
	}
}

func TestExistRandom(t *testing.T) {
	for i := 0; i < 500; i++ {
		m := 1 + rand.Intn(4)
		n := 1 + rand.Intn(4)
		board := make([][]byte, m)
		for r := range board {
			board[r] = make([]byte, n)
			for c := range board[r] {
				board[r][c] = byte('A' + rand.Intn(3))
			}
		}
		wl := 1 + rand.Intn(5)
		var sb strings.Builder
		for j := 0; j < wl; j++ {
			sb.WriteByte(byte('A' + rand.Intn(3)))
		}
		word := sb.String()
		want := refExist(board, word)
		var got bool
		runSafely(t, func() { got = lc.Exist(cloneBoard(board), word) })
		if got != want {
			t.Fatalf("Exist(board=%v, word=%s) = %v, want %v", board, word, got, want)
		}
	}
}

// ---- 131. 分割回文串 ----

func isPalindromeStr(s string) bool {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
			return false
		}
	}
	return true
}

func refPartition(s string) [][]string {
	n := len(s)
	var res [][]string
	for mask := 0; mask < 1<<(n-1); mask++ {
		parts := []string{}
		prev := 0
		for i := 0; i < n-1; i++ {
			if mask&(1<<i) != 0 {
				parts = append(parts, s[prev:i+1])
				prev = i + 1
			}
		}
		parts = append(parts, s[prev:])
		ok := true
		for _, p := range parts {
			if !isPalindromeStr(p) {
				ok = false
				break
			}
		}
		if ok {
			res = append(res, parts)
		}
	}
	return res
}

func TestPartition(t *testing.T) {
	cases := []struct {
		s    string
		want int
	}{
		{"aab", 2},
		{"a", 1},
		{"aa", 2},
		{"abc", 1},
	}
	for _, c := range cases {
		var got [][]string
		runSafely(t, func() { got = lc.Partition(c.s) })
		if len(got) != c.want {
			t.Errorf("Partition(%q) 数量 = %d, want %d", c.s, len(got), c.want)
		}
	}
}

func TestPartitionRandom(t *testing.T) {
	for i := 0; i < 400; i++ {
		n := 1 + rand.Intn(9)
		var sb strings.Builder
		for j := 0; j < n; j++ {
			sb.WriteByte(byte('a' + rand.Intn(3)))
		}
		s := sb.String()
		var got [][]string
		runSafely(t, func() { got = lc.Partition(s) })
		want := refPartition(s)
		if setKey(got) != setKey(want) {
			t.Fatalf("Partition(%q) 集合不一致", s)
		}
	}
}

// ---- 51. N 皇后 ----

func refSolveNQueens(n int) [][]string {
	var res [][]string
	cols := make([]int, n)
	for i := range cols {
		cols[i] = i
	}
	var perm func(idx int)
	perm = func(idx int) {
		if idx == n {
			ok := true
			for r1 := 0; r1 < n && ok; r1++ {
				for r2 := r1 + 1; r2 < n; r2++ {
					if abs(r1-r2) == abs(cols[r1]-cols[r2]) {
						ok = false
						break
					}
				}
			}
			if ok {
				b := make([]string, n)
				for r := 0; r < n; r++ {
					row := make([]byte, n)
					for c := 0; c < n; c++ {
						if cols[r] == c {
							row[c] = 'Q'
						} else {
							row[c] = '.'
						}
					}
					b[r] = string(row)
				}
				res = append(res, b)
			}
			return
		}
		for i := idx; i < n; i++ {
			cols[idx], cols[i] = cols[i], cols[idx]
			perm(idx + 1)
			cols[idx], cols[i] = cols[i], cols[idx]
		}
	}
	perm(0)
	return res
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func TestSolveNQueens(t *testing.T) {
	cases := []struct{ n, want int }{
		{1, 1}, {4, 2}, {5, 10}, {6, 4},
	}
	for _, c := range cases {
		var got [][]string
		runSafely(t, func() { got = lc.SolveNQueens(c.n) })
		if len(got) != c.want {
			t.Errorf("SolveNQueens(%d) 数量 = %d, want %d", c.n, len(got), c.want)
		}
		want := refSolveNQueens(c.n)
		if setKey(got) != setKey(want) {
			t.Errorf("SolveNQueens(%d) 棋盘集合不一致", c.n)
		}
	}
}
