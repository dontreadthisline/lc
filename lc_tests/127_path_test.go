package lc_tests

import (
	"math/rand"
	"testing"

	"demo/lc"
)

// ---- 127 变体: 枚举 begin→end 的所有转换路径 ----

// refLadderPaths: 独立参照 —— 预建邻接表(被测是逐位变异) + 每层复制 visited 的纯函数式 DFS(被测是撤销式)
// 结构刻意不同, 对拍才有意义
func refLadderPaths(begin, end string, wordList []string) [][]string {
	wordSet := map[string]bool{}
	nodes := []string{begin}
	for _, w := range wordList {
		if w != begin && !wordSet[w] {
			wordSet[w] = true
			nodes = append(nodes, w)
		}
	}
	if !wordSet[end] {
		return nil
	}
	diff1 := func(a, b string) bool {
		d := 0
		for i := 0; i < len(a); i++ {
			if a[i] != b[i] {
				d++
			}
		}
		return d == 1
	}
	adj := map[string][]string{}
	for i, a := range nodes {
		for j, b := range nodes {
			if i != j && diff1(a, b) {
				adj[a] = append(adj[a], b)
			}
		}
	}
	var res [][]string
	var walk func(cur string, path []string, visited map[string]bool)
	walk = func(cur string, path []string, visited map[string]bool) {
		if cur == end {
			res = append(res, append([]string(nil), path...))
			return
		}
		for _, nb := range adj[cur] {
			if !visited[nb] {
				nv := make(map[string]bool, len(visited)+1)
				for k := range visited {
					nv[k] = true
				}
				nv[nb] = true
				walk(nb, append(path, nb), nv)
			}
		}
	}
	if begin == end {
		return [][]string{{begin}}
	}
	walk(begin, []string{begin}, map[string]bool{begin: true})
	return res
}

func TestLadderLengthPath(t *testing.T) {
	cases := []struct {
		name     string
		begin    string
		end      string
		wordList []string
		want     [][]string
	}{
		{
			"官方词典_全部8条",
			"hit", "cog", []string{"hot", "dot", "dog", "lot", "log", "cog"},
			[][]string{
				// 注意: dot↔lot、dog↔log 也互为邻居(差一字母), 图比直觉稠密
				{"hit", "hot", "dot", "dog", "cog"},
				{"hit", "hot", "dot", "dog", "log", "cog"},
				{"hit", "hot", "dot", "lot", "log", "cog"},
				{"hit", "hot", "dot", "lot", "log", "dog", "cog"},
				{"hit", "hot", "lot", "dot", "dog", "cog"},
				{"hit", "hot", "lot", "dot", "dog", "log", "cog"},
				{"hit", "hot", "lot", "log", "cog"},
				{"hit", "hot", "lot", "log", "dog", "cog"},
			},
		},
		{
			"一步直达",
			"aa", "ab", []string{"ab"},
			[][]string{{"aa", "ab"}},
		},
		{
			"无路径_邻居均不可达",
			"aa", "zz", []string{"zz", "by"},
			nil,
		},
		{
			"end不在词典",
			"aa", "zz", []string{"ab"},
			nil,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got [][]string
			runSafely(t, func() { got = lc.LadderLengthPath(c.begin, c.end, c.wordList) })
			if setKey(got) != setKey(c.want) {
				t.Errorf("LadderLengthPath(%q, %q, %v)\n got  %s\n want %s",
					c.begin, c.end, c.wordList, setKey(got), setKey(c.want))
			}
		})
	}
}

func TestLadderLengthPathRandom(t *testing.T) {
	// 2 字母小词典: 词数 <=8, 简单路径总量可控, 覆盖各种共享中转结构
	for i := 0; i < 2000; i++ {
		nWords := 3 + rand.Intn(6)
		pool := map[string]bool{}
		for len(pool) < nWords {
			pool[string([]byte{byte('a' + rand.Intn(4)), byte('a' + rand.Intn(4))})] = true
		}
		wordList := make([]string, 0, len(pool))
		for w := range pool {
			wordList = append(wordList, w)
		}
		begin := string([]byte{byte('a' + rand.Intn(4)), byte('a' + rand.Intn(4))})
		end := wordList[rand.Intn(len(wordList))]
		for end == begin {
			end = string([]byte{byte('a' + rand.Intn(4)), byte('a' + rand.Intn(4))})
			if !pool[end] {
				wordList = append(wordList, end)
				pool[end] = true
			}
		}
		want := refLadderPaths(begin, end, wordList)
		var got [][]string
		runSafely(t, func() { got = lc.LadderLengthPath(begin, end, append([]string(nil), wordList...)) })
		if setKey(got) != setKey(want) {
			t.Fatalf("LadderLengthPath(%q, %q, %v)\n got  %s\n want %s",
				begin, end, wordList, setKey(got), setKey(want))
		}
	}
}
