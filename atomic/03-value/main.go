// 03-value: atomic.Value —— 任意类型的原子存储
//
// atomic.Value 可以对任意 Go 类型做原子 Load/Store。
// 和 atomic.Pointer[T] 不同，atomic.Value 存储的是 interface{}，
// 不需要暴露指针。
//
// 关键规则（违反会 panic）：
//   1. 第一次 Store 确定了类型，后续 Store 必须是同一具体类型
//   2. 不能 Store nil（第一次也不行）
//   3. 不能 Store 不同具体类型的值
//
// 适用场景：
//   - 配置热更新（读多写少）
//   - 缓存快照替换
//   - 需要原子替换整个值但又不想暴露指针的场合
//
// 运行: go run ./atomic/03-value/

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	fmt.Print("=== 03 atomic.Value ===\n\n")

	basicUsage()
	configHotReload()
	cacheSnapshot()
	typeSafetyDemo()
}

// ---- 基本用法 ----

func basicUsage() {
	var v atomic.Value

	// Store: 原子写入。第一次 Store 确定类型。
	v.Store("hello")
	// Load: 原子读取，返回 interface{}，需要类型断言。
	if s, ok := v.Load().(string); ok {
		fmt.Printf("Load: %s\n", s)
	}

	// Swap: 原子交换 (Go 1.17+)，返回旧值
	old := v.Swap("world")
	fmt.Printf("Swap: 旧值=%s, 新值=%s\n", old, v.Load())

	// CompareAndSwap: 原子 CAS (Go 1.17+)
	swapped := v.CompareAndSwap("world", "goodbye")
	fmt.Printf("CAS world→goodbye 成功? %v, 当前: %s\n", swapped, v.Load())

	// 这个 CAS 会失败，因为当前值不是 "hello"
	swapped = v.CompareAndSwap("hello", "nope")
	fmt.Printf("CAS hello→nope 成功? %v, 当前: %s\n", swapped, v.Load())

	fmt.Println()
}

// ---- 模式 1: 配置热更新 ----

type AppConfig struct {
	Host    string
	Port    int
	Debug   bool
	Timeout time.Duration
}

func configHotReload() {
	var cfg atomic.Value

	// 初始配置（第一次 Store 不能为 nil）
	cfg.Store(&AppConfig{
		Host:    "localhost",
		Port:    8080,
		Debug:   true,
		Timeout: 30 * time.Second,
	})

	// 多个 goroutine 并发读——无锁，极快
	var wg sync.WaitGroup
	for i := range 5 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			c := cfg.Load().(*AppConfig)
			fmt.Printf("读者 %d 读到: %s:%d (debug=%v)\n",
				id, c.Host, c.Port, c.Debug)
		}(i)
	}

	// 写者：原子替换整个配置
	// 注意：Load 返回的 *AppConfig 不应被修改——
	// 修改通过 Load 拿到的指针指向的对象会产生 data race。
	// 正确的做法是：创建新对象 → Store 新指针。
	//
	// 错误示例：
	//   c := cfg.Load().(*AppConfig)
	//   c.Port = 9090  // data race！其他 goroutine 可能正在读同一个对象
	//
	// 正确做法——immutable 风格：
	newCfg := &AppConfig{
		Host:    "0.0.0.0",
		Port:    9090,
		Debug:   false,
		Timeout: 60 * time.Second,
	}
	cfg.Store(newCfg)

	wg.Wait()

	// 验证更新
	c := cfg.Load().(*AppConfig)
	fmt.Printf("配置已更新: %s:%d (debug=%v)\n\n", c.Host, c.Port, c.Debug)
}

// ---- 模式 2: 缓存快照替换 ----

type CacheEntry struct {
	Data      map[string]string
	UpdatedAt time.Time
}

func cacheSnapshot() {
	var cache atomic.Value

	// 初始化缓存
	cache.Store(&CacheEntry{
		Data:      map[string]string{"key1": "val1"},
		UpdatedAt: time.Now(),
	})

	// 每隔一段时间，全量重建缓存
	go func() {
		for i := range 2 {
			time.Sleep(50 * time.Millisecond)
			newData := map[string]string{
				"key1": fmt.Sprintf("val1_v%d", i+2),
				"key2": fmt.Sprintf("val2_v%d", i+2),
			}
			cache.Store(&CacheEntry{
				Data:      newData,
				UpdatedAt: time.Now(),
			})
			fmt.Printf("缓存第 %d 次刷新\n", i+1)
		}
	}()

	// 读者
	time.Sleep(10 * time.Millisecond)
	for i := range 3 {
		entry := cache.Load().(*CacheEntry)
		fmt.Printf("读取缓存 (第 %d 次): data=%v, updated=%s\n",
			i+1, entry.Data, entry.UpdatedAt.Format("15:04:05.000"))
		time.Sleep(60 * time.Millisecond)
	}

	fmt.Println()
}

// ---- 类型安全注意事项 ----

func typeSafetyDemo() {
	var v atomic.Value

	// 规则 1: 第一次 Store 不能是 nil
	// v.Store(nil)  // panic: sync/atomic: store of nil value into Value

	v.Store("initial string")

	// 规则 2: 后续 Store 必须是同一具体类型，否则 panic
	// v.Store(42)  // panic: sync/atomic: store of inconsistently typed value into Value

	// 规则 3: 从 atomic.Value 取出的值，其内部的引用类型字段仍然可变。
	//         只保护了"哪个对象"的原子性，不保护对象内部字段的并发安全。
	//
	// 这意味着：如果你 Store 了一个 map/slice，然后修改它——data race。
	// atomic.Value 适合 immutable 快照模式，不适合用作可变共享状态的容器。

	// ---- 类型断言健壮写法 ----
	// 在生产代码中，建议用带检查的类型断言，以防因某些原因类型不匹配导致 panic。

	val := v.Load()
	switch val := val.(type) {
	case string:
		fmt.Printf("值类型正确: %s\n", val)
	default:
		fmt.Printf("意外类型: %T\n", val)
	}

	fmt.Println()

	// ---- atomic.Value vs atomic.Pointer[T] 的取舍 ----
	//
	// atomic.Pointer[T]:
	//   + 编译期类型安全，不需要类型断言
	//   + Load 返回 *T，不是 interface{}
	//   - 暴露指针——调用者拿到 *T 后可以修改内部字段
	//
	// atomic.Value:
	//   + 不暴露指针——对调用者更像"值传递"
	//   + 易于封装（可以在 Store 前做校验）
	//   - 运行时类型检查，出错 panic
	//   - Load 需要类型断言
	//
	// 推荐：能暴露指针的场景优先用 atomic.Pointer[T]，
	//       需要类型擦除或不想暴露指针的场景用 atomic.Value。
}
