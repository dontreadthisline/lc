// 04-patterns: atomic 实战模式与最佳实践
//
// 汇集生产中最常用的 atomic 模式：
//   1. 高性能计数器（Prometheus Counter 风格）
//   2. 无锁环形队列索引
//   3. 优雅关闭标志
//   4. 单次初始化（sync.Once 的 atomic 底层原理）
//   5. 非阻塞读写的配置管理器
//   6. atomic + channel 混合模式
//
// 运行: go run ./atomic/04-patterns/

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	fmt.Print("=== 04 atomic 实战模式 ===\n\n")

	highPerfCounter()
	gracefulShutdown()
	singleInit()
	nonBlockingConfig()
	atomicChannel()
}

// ---- 模式 1: 高性能计数器 ----

type Metrics struct {
	requestCount atomic.Int64
	errorCount   atomic.Int64
	totalLatency atomic.Int64 // 纳秒累计
}

func (m *Metrics) RecordRequest(latency time.Duration, isError bool) {
	m.requestCount.Add(1)
	m.totalLatency.Add(int64(latency))
	if isError {
		m.errorCount.Add(1)
	}
}

func (m *Metrics) Snapshot() (requests, errors int64, avgLatency time.Duration) {
	requests = m.requestCount.Load()
	errors = m.errorCount.Load()
	totalLat := m.totalLatency.Load()
	if requests > 0 {
		avgLatency = time.Duration(totalLat / requests)
	}
	return
}

func highPerfCounter() {
	var m Metrics

	var wg sync.WaitGroup
	n := 100

	for i := range n {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			latency := time.Duration(1+i%10) * time.Millisecond
			isErr := i%20 == 0 // 5% 错误率
			m.RecordRequest(latency, isErr)
		}(i)
	}
	wg.Wait()

	reqs, errs, avgLat := m.Snapshot()
	fmt.Printf("计数器: 请求=%d, 错误=%d, 平均延迟=%v\n", reqs, errs, avgLat)
	fmt.Println()
}

// ---- 模式 2: 优雅关闭 ----

func gracefulShutdown() {
	var closed atomic.Bool
	var wg sync.WaitGroup

	// 启动 worker
	const workerCount = 3
	for i := range workerCount {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				if closed.Load() {
					fmt.Printf("  worker %d 收到关闭信号，退出\n", id)
					return
				}
				// 模拟工作
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	time.Sleep(30 * time.Millisecond)
	fmt.Println("发送关闭信号...")
	closed.Store(true)
	wg.Wait()
	fmt.Println("所有 worker 已退出")
	fmt.Println()
}

// ---- 模式 3: 单次初始化（sync.Once 原理） ----

// sync.Once 的内部实现就是 atomic + Mutex 的组合：
//   1. 用 atomic 快速检查是否已完成
//   2. 如果未完成，用 Mutex 保护初始化逻辑
//   3. 初始化完成后用 atomic 标记 done
//
// 下面的代码是 sync.Once 的精简版实现，展示其原理。

type Once struct {
	done atomic.Uint32
	mu   sync.Mutex
}

func (o *Once) Do(f func()) {
	if o.done.Load() == 1 {
		return // 快速路径：已经初始化过，直接返回
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.done.Load() == 0 { // 双重检查
		f()
		o.done.Store(1) // 先完成初始化，再设标志
	}
}

func singleInit() {
	var once Once
	var counter atomic.Int32

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			once.Do(func() {
				counter.Add(1)
				fmt.Println("  初始化逻辑执行了")
			})
		}()
	}
	wg.Wait()

	fmt.Printf("  counter=%d (期望 1，只初始化一次)\n\n", counter.Load())
}

// ---- 模式 4: 非阻塞读写配置管理器 ----

// 当配置读取频率极高（百万 QPS），写频率极低（分钟级）时，
// 用 RWMutex 的读锁开销也值得消除。此时用 atomic.Value 实现无锁读。

type NonBlockingConfigManager struct {
	current atomic.Value // 存储 *ConfigSnapshot
}

type ConfigSnapshot struct {
	Settings map[string]string
	Version  int64
}

func NewNonBlockingConfigManager(initial map[string]string) *NonBlockingConfigManager {
	m := &NonBlockingConfigManager{}
	m.current.Store(&ConfigSnapshot{
		Settings: initial,
		Version:  1,
	})
	return m
}

// Get: 无锁读取。在高并发读场景下性能远优于 RWMutex。
func (m *NonBlockingConfigManager) Get() *ConfigSnapshot {
	return m.current.Load().(*ConfigSnapshot)
}

// Update: 写入。创建新快照后原子替换。
func (m *NonBlockingConfigManager) Update(newSettings map[string]string) {
	old := m.current.Load().(*ConfigSnapshot)
	m.current.Store(&ConfigSnapshot{
		Settings: newSettings,
		Version:  old.Version + 1,
	})
}

func nonBlockingConfig() {
	mgr := NewNonBlockingConfigManager(map[string]string{
		"host": "localhost",
		"port": "8080",
	})

	// 并发读——无锁
	var wg sync.WaitGroup
	for i := range 5 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cfg := mgr.Get()
			fmt.Printf("  读者 %d: version=%d, host=%s\n",
				id, cfg.Version, cfg.Settings["host"])
		}(i)
	}

	wg.Wait()

	// 原子更新配置
	mgr.Update(map[string]string{
		"host": "0.0.0.0",
		"port": "9090",
	})
	cfg := mgr.Get()
	fmt.Printf("  更新后: version=%d, host=%s\n\n",
		cfg.Version, cfg.Settings["host"])
}

// ---- 模式 5: atomic + channel 混合模式 ----

// atomic 擅长"状态标志"，channel 擅长"通知和同步"。
// 两者结合可以实现高效且正确的并发模式。

// 例子：速率限制器
//   - atomic 记录当前令牌数（无锁）
//   - channel 发送"需要令牌"的信号

type RateLimiter struct {
	burst  int64          // 最大令牌数
	tokens atomic.Int64   // 当前令牌数
	addCh  chan struct{}  // 补充令牌信号
	stopCh chan struct{}  // 停止信号
}

func NewRateLimiter(burst int64, ratePerSec int64) *RateLimiter {
	rl := &RateLimiter{
		burst:  burst,
		addCh:  make(chan struct{}, 1),
		stopCh: make(chan struct{}),
	}
	rl.tokens.Store(burst)

	// 后台协程定时补充令牌
	go func() {
		interval := time.Second / time.Duration(ratePerSec)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				// 原子补充令牌，不超过 burst
				for {
					current := rl.tokens.Load()
					if current >= burst {
						break
					}
					if rl.tokens.CompareAndSwap(current, min(current+1, burst)) {
						break
					}
				}
			case <-rl.stopCh:
				return
			}
		}
	}()

	return rl
}

func (rl *RateLimiter) Allow() bool {
	for {
		current := rl.tokens.Load()
		if current <= 0 {
			return false
		}
		if rl.tokens.CompareAndSwap(current, current-1) {
			return true
		}
	}
}

func (rl *RateLimiter) Stop() {
	close(rl.stopCh)
}

func atomicChannel() {
	rl := NewRateLimiter(5, 10) // burst=5, 每秒补充 10 个令牌

	allowed := 0
	denied := 0
	for range 10 {
		if rl.Allow() {
			allowed++
		} else {
			denied++
		}
	}

	fmt.Printf("速率限制器: 允许=%d, 拒绝=%d (burst=%d)\n", allowed, denied, 5)
	rl.Stop()
	fmt.Println()
}

// ---- 最佳实践总结 ----

// 本文件演示的各个模式在选择时遵循以下原则：
//
// 用 atomic 的时机：
//   1. 单个变量的并发读写（计数器、标志位）
//   2. 读多写少，且写是"整体替换"（配置管理）
//   3. 需要极致性能，且逻辑简单到 CAS 循环能搞定
//   4. 和 channel 配合，atomic 管状态，channel 管通知
//
// 用 Mutex 的时机：
//   1. 需要保护多个相关变量的原子性（转账：扣 A 加 B）
//   2. 临界区代码比较复杂（多行操作）
//   3. 并发竞争激烈（CAS 自旋浪费 CPU）
//   4. 需要结合条件变量（sync.Cond）的等待/通知
//
// 用 Channel 的时机：
//   1. goroutine 之间需要通信和同步
//   2. 需要阻塞等待
//   3. 需要传递数据所有权
//   4. 限流、信号量模式
//
// 经典名言：
//   "Don't communicate by sharing memory; share memory by communicating."
//   但 atomic 就是"通过共享内存来通信"的高效原语——在合适的场景使用即可。

// ---- 补充：常见坑与注意事项 ----

// 1. 原子操作不保证顺序
//    atomic.Add 只保证单个操作是原子的，不保证多个操作之间的 happens-before 关系。
//    如果 goroutine A 先 Add counter，goroutine B 后 Add counter，
//    不能保证 A 的 Add 在 B 之前被其他 goroutine 观察到。

// 2. 64 位对齐
//    在 32 位平台上，int64/uint64 必须 8 字节对齐，否则 atomic 操作会 panic。
//    结构体中把 64 位字段放在最前面可以避免这个问题。
//    在 64 位平台上（amd64/arm64）不会有这个问题。

// 3. atomic.Value 的类型锁定
//    第一次 Store 之后类型就固定了——想存别的类型需要用新的 atomic.Value。

// 4. 不要用 atomic 替代 Mutex 做复杂操作
//    如果临界区超过 2-3 行代码，或者涉及多个变量，Mutex 更安全且代码更清晰。

// 5. Go 1.19+ 优先用类型化原子变量
//    atomic.Int64 替代 atomic.AddInt64(&x, 1)，
//    atomic.Bool 替代用 int32 模拟布尔值。
//    更好的 API 设计意味着更少的 bug。

// 6. noCopy 检查
//    类型化原子变量包含 noCopy 字段（通过 sync.Locker 接口）。
//    复制含有 atomic 字段的结构体会触发 go vet 警告。
//    通过指针传递可以避免。

// 7. 32 位平台上 AddUint64 减法技巧
//    对于 uint32/uint64，没有原生的减法操作。
//    AddUint64(&x, ^uint64(0)) 是减 1 的惯用写法（利用二进制补码）。
//    Go 1.19+ 的 atomic.Uint64 同样没有 Sub 方法。
//    示例: atomic.AddUint64(&x, ^uint64(10-1)) 等价于 x -= 10
func demoUintSubtraction() {
	var x atomic.Uint64
	x.Store(100)
	// 减法：Add(^uint64(delta-1)) 等价于 Add(-delta)
	x.Add(^uint64(10 - 1)) // 减 10
	fmt.Printf("100 - 10 = %d\n", x.Load())
}
