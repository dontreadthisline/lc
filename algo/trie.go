package algo

type Trie struct {
	end    bool
	val    int
	childs [26]*Trie
}

func NewTrie() *Trie {
	return &Trie{}
}

func insert(t *Trie, key string, value int) {
	if key == "" {
		return
	}
	for i := range len(key) {
		idx := key[i] - 'a'
		if t.childs[idx] == nil {
			t.childs[idx] = new(Trie)
		}
		t = t.childs[idx]
	}
	t.end = true
	t.val = value
}

func search(t *Trie, key string) *Trie {

	for i := range len(key) {
		c := key[i]
		if t == nil || t.childs[c] == nil {
			return nil
		}
		t = t.childs[c]
	}
	return t
}

type CollectFunc func(string, *Trie)

func dfs(t *Trie, word []byte, collect CollectFunc) {
	if t == nil {
		return
	}

	if t.end {
		collect(string(word), t) //copy 而不是share
	}

	for i, child := range t.childs {
		if child == nil {
			continue
		}
		word = append(word, byte(i+'a'))
		dfs(child, word, collect)
		word = word[:len(word)-1]
	}
}

func startWith(t *Trie, prefix string) ([]string, []int) {
	words := []string{}
	vals := []int{}
	collect := func(word string, node *Trie) {
		words = append(words, word)
		vals = append(vals, node.val)
	}
	node := search(t, prefix)
	dfs(node, []byte(prefix), collect)
	return words, vals
}
