package algo_tests

import "demo/algo"

import (
	"math/rand"
	"slices"
	"strings"
	"testing"
)

func TestMaxKFreqWords(t *testing.T) {
	cases := []struct {
		name  string
		words []string
		k     int
		want  []string
	}{
		{"示例1", []string{"i", "love", "leetcode", "i", "love", "coding"}, 2, []string{"i", "love"}},
		{"示例2", []string{"the", "day", "is", "sunny", "the", "the", "the", "sunny", "is", "is"}, 4, []string{"the", "is", "sunny", "day"}},
		{"单词", []string{"a"}, 1, []string{"a"}},
		// tiebreak 专项：五个词全同频 1，前 3 应为字典序最小的 a,b,c
		{"同频tiebreak", []string{"e", "d", "c", "b", "a"}, 3, []string{"a", "b", "c"}},
		// 部分同频：aa=2, bb=2, cc=2, zz=1 → 前3 = aa,bb,cc（zz 频次低先出局）
		{"部分同频", []string{"zz", "aa", "aa", "bb", "bb", "cc", "cc"}, 3, []string{"aa", "bb", "cc"}},
	}
	for _, c := range cases {
		if got := algo.MaxKFreqWords(c.words, c.k); !slices.Equal(got, c.want) {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}

// 随机对照：map 遍历顺序天然乱序，专抓顺序敏感的堆 bug
func TestMaxKFreqWordsRandom(t *testing.T) {
	r := rand.New(rand.NewSource(7))
	vocab := []string{"apple", "banana", "cherry", "date", "elder", "fig", "grape", "honey"}
	for round := 0; round < 300; round++ {
		n := 1 + r.Intn(40)
		words := make([]string, n)
		for i := range words {
			words[i] = vocab[r.Intn(len(vocab))]
		}
		// oracle：计数 + 稳定排序（频次降序，同频字典序升序）取前 k
		cnt := map[string]int{}
		for _, w := range words {
			cnt[w]++
		}
		keys := make([]string, 0, len(cnt))
		for w := range cnt {
			keys = append(keys, w)
		}
		slices.SortFunc(keys, func(a, b string) int {
			if cnt[a] != cnt[b] {
				return cnt[b] - cnt[a] // 频次降序
			}
			return strings.Compare(a, b) // 同频字典序升序
		})
		k := 1 + r.Intn(len(keys))
		want := keys[:k]

		if got := algo.MaxKFreqWords(words, k); !slices.Equal(got, want) {
			t.Fatalf("round %d: words=%v k=%d\n got %v\nwant %v", round, words, k, got, want)
		}
	}
}
