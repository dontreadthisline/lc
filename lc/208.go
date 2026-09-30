package lc

/*
Trie（发音类似 "try"）或者说 前缀树 是一种树形数据结构，用于高效地存储和检索字符串数据集中的键。这一数据结构有相当多的应用情景，例如自动补全和拼写检查。

请你实现 Trie 类：

- Trie() 初始化前缀树对象。

- void insert(String word) 向前缀树中插入字符串 word 。

- boolean search(String word) 如果字符串 word 在前缀树中，返回 true（即，在检索之前已经插入）；否则，返回 false 。

- boolean startsWith(String prefix) 如果之前已经插入的字符串 word 的前缀之一为 prefix ，返回 true ；否则，返回 false 。

示例：

输入
["Trie", "insert", "search", "search", "startsWith", "insert", "search"]
[[], ["apple"], ["apple"], ["app"], ["app"], ["app"], ["app"]]
输出
[null, null, true, false, true, null, true]

解释
Trie trie = new Trie();
trie.insert("apple");
trie.search("apple");   // 返回 True
trie.search("app");     // 返回 False
trie.startsWith("app"); // 返回 True
trie.insert("app");
trie.search("app");     // 返回 True

提示：

- 1 <= word.length, prefix.length <= 2000

- word 和 prefix 仅由小写英文字母组成

- insert、search 和 startsWith 调用次数 总计 不超过 3 * 10^4 次
*/

//Trie trie节点
type Trie struct {
	children [26]*Trie
	end bool
}

func NewTrie() *Trie {
	return &Trie{}
}

func (t *Trie) Insert(word string) {
	insert(t,word)
}


func (t *Trie) Search(word string) bool {
	node := search(t,word)
	return node != nil && node.end
}

func (t *Trie) StartsWith(prefix string) bool {
	node := search(t, prefix )
	return node != nil
}

func insert(t *Trie,word string) {
	if word == "" {
		return
	}
	for i := range len(word) {
		c := word[i]
		if t.children[c] == nil {
			t.children[c] = new(Trie)
		}
		t = t.children[c]
	}

	t.end = true
}

func search(t *Trie, word string) *Trie {
	for i := range len(word) {
		if t == nil  || t.children[word[i]] == nil {
			return nil
		}
		t = t.children[word[i]]
	}
	return t
}
