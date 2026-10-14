package lc

import "slices"

/*
字典 wordList 中从单词 beginWord* *到 endWord 的 转换序列 是一个按下述规格形成的序列 beginWord -> s_1 -> s_2 -> ... -> s_k：
//  family8_test.go:62: LadderLength("aaa", "zzz", [aaa bbb ccc zzz]) = 1, want 0
- 每一对相邻的单词只差一个字母。

-  对于 1 <= i <= k 时，每个 s_i 都在 wordList 中。注意， beginWord* *不需要在 wordList 中。

- s_k == endWord

给你两个单词* *beginWord* *和 endWord 和一个字典 wordList ，返回 *从 beginWord 到 endWord 的 最短转换序列 中的 单词数目* 。如果不存在这样的转换序列，返回 0 。

示例 1：

输入：beginWord = "hit", endWord = "cog", wordList = ["hot","dot","dog","lot","log","cog"]
输出：5
解释：一个最短转换序列是 "hit" -> "hot" -> "dot" -> "dog" -> "cog", 返回它的长度 5。

示例 2：

输入：beginWord = "hit", endWord = "cog", wordList = ["hot","dot","dog","lot","log"]
输出：0
解释：endWord "cog" 不在字典中，所以无法进行转换。

提示：

- 1 <= beginWord.length <= 10

- endWord.length == beginWord.length

- 1 <= wordList.length <= 5000

- wordList[i].length == beginWord.length

- beginWord、endWord 和 wordList[i] 由小写英文字母组成

- beginWord != endWord

- wordList 中的所有字符串 互不相同
*/

func LadderLength(begin string, end string, wordList []string) int {
	words := make(map[string]bool,len(wordList))
	for _,word := range wordList {
		words[word] = true
	}

	if !words[end] {
		return 0
	}

	visited := map[string]bool{}
	bfs := func(begin, end string)  int {
		cnt := 1//初始单词也算
		visited[begin] = true
		queue := []string{begin}
		for len(queue) > 0 {
			for range len(queue) {
				word := queue[0]
				if word == end {
					return cnt
				}
				//visited adj
				for i := range len(word) {
					for j := range 26 {
						tmp := []byte(word)
						tmp[i] = byte('a' + j)
						adj := string(tmp)
						if !visited[adj] && words[adj] {
							queue = append(queue,adj)
							visited[adj] = true
						}
					}
				}
				queue = queue[1:]
			}
			cnt += 1
	  }
		//如果没有路径,就返回0
		return 0
	}
	return bfs(begin,end)
}

/*
* 图类型, 普通图(类比阶乘),虚拟节点起点:起点单词,图的边构成:仅有单个位置不同的单词对
*
* 搜索类型:全路径,限制:不可重用节点。
* 套用模型: path_search:AllPathSimple
*/

func LadderLengthPath(begin string, end string, wordList []string) [][]string {
	words := make(map[string]bool,len(wordList))
	for _,word := range wordList {
		words[word] = true
	}

	var res [][]string
	if !words[end] {
		return  res
	}

	visited := map[string]bool{}
	var dfs func( string,string, []string)

	dfs = func(cur,to string,path []string) {
		if cur == to {
			res = append(res,slices.Clone(path))
			return
		}
		visited[cur] = true
		for i := range len(cur) {
			for j := range 26 {
				tmp := []byte(cur)
				tmp[i] = byte('a' + j)
				adj := string(tmp)
				if !visited[adj] && words[adj] {
					dfs(adj,to,append(path,adj))
				}
			}
		}
		visited[cur] = false
	}

	dfs(begin,end,[]string{begin})
	return res
}
