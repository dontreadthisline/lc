package lc

// LRU 进行中（原稿备份 /tmp/146_wip.go.bak）
type LRUCache struct{}

func NewLRUCache(capacity int) *LRUCache { return &LRUCache{} }
func (c *LRUCache) Get(key int) int { return -1 }
func (c *LRUCache) Put(key int, value int) {}
