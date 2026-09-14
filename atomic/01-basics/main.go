// 01-basics: sync/atomic 基础操作
//
// 演示原子操作的核心函数：
//   - AddInt64 / AddUint64  原子加减
//   - LoadInt64 / LoadUint64 原子读取
//   - StoreInt64            原子写入
//   - SwapInt64             原子交换
//   - CompareAndSwapInt64   原子比较并交换 (CAS)
//
// 关键认知：
//   - 原子操作保证单个变量的读写是"不可分割"的，不会出现读到一半被修改的情况
//   - 原子操作是 lock-free 的，比 sync.Mutex 轻量得多，适合简单计数/状态切换场景
//   - 原子操作不保证多个变量之间的顺序一致性——那需要 Mutex 或 Channel
//   - 在 x86/x86_64 上，对齐的 64 位操作本身就是原子的，但编译器可能重排指令，
//     所以仍然必须用 atomic 函数
//
// 运行: go run ./atomic/01-basics/

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	fmt.Print("=== 01 sync/atomic 基础操作 ===\n\n")

	basicArithmetic()
	casPattern()
	swapPattern()
}

// ---- 基础加减与读写 ----

func basicArithmetic() {
	var counter int64

	// AddInt64: 原子加。返回新值。注意和 LoadInt64 的区别——Add 是读-改-写一体。
	// 如果你只是需要 +1 然后拿到新值，用 Add 一步完成，不需要 Load+Store。
	newVal := atomic.AddInt64(&counter, 1)
	fmt.Printf("AddInt64 +1 后: %d\n", newVal)

	atomic.AddInt64(&counter, 10)
	fmt.Printf("AddInt64 +10 后: %d\n", atomic.LoadInt64(&counter))

	// LoadInt64: 原子读取。即使是 int64，在 32 位平台上也可能被拆成两次 32 位读，
	// 不用 atomic.Load 就可能读到"半个"修改中的值（tearing）。
	val := atomic.LoadInt64(&counter)
	fmt.Printf("LoadInt64: %d\n", val)

	// StoreInt64: 原子写入。对称于 Load。
	atomic.StoreInt64(&counter, 0)
	fmt.Printf("StoreInt64 重置后: %d\n", atomic.LoadInt64(&counter))

	// ---- 常见错误：以为 += 是原子的 ----
	// counter += 1   // 错误！这是读-改-写三步，不是原子的
	// 必须用 atomic.AddInt64(&counter, 1)

	fmt.Println()
}

// ---- CAS: Compare And Swap ----

func casPattern() {
	var state int64

	// CompareAndSwapInt64(addr, old, new): 如果 *addr == old，就设为 new，返回 true；
	// 否则不做任何事，返回 false。整个操作是原子的。
	//
	// 这是实现 lock-free 数据结构的基础原语。

	swapped := atomic.CompareAndSwapInt64(&state, 0, 1)
	fmt.Printf("CAS 0→1 成功? %v, 当前值: %d\n", swapped, atomic.LoadInt64(&state))

	// 这次会失败，因为当前值是 1，不是 0
	swapped = atomic.CompareAndSwapInt64(&state, 0, 2)
	fmt.Printf("CAS 0→2 成功? %v, 当前值: %d\n", swapped, atomic.LoadInt64(&state))

	// 常见 CAS 循环模式：自旋重试
	// 适用于竞争不激烈的场景。竞争激烈时用 Mutex 更合适。
	var count int64
	const target = int64(100)

	for {
		old := atomic.LoadInt64(&count)
		if old >= target {
			break
		}
		// 尝试 CAS：只有 count 还是 old 时才更新
		if atomic.CompareAndSwapInt64(&count, old, old+1) {
			// 成功，继续下一轮
		}
		// 失败说明被别的 goroutine 抢先了，循环重试
	}
	fmt.Printf("CAS 自旋循环结果: count=%d\n", atomic.LoadInt64(&count))

	fmt.Println()
}

// ---- Swap: 原子交换 ----

func swapPattern() {
	var current int64 = 42

	// SwapInt64: 原子地把 new 写入，返回旧值。是 Store+Load 的原子版。
	old := atomic.SwapInt64(&current, 100)
	fmt.Printf("Swap: 旧值=%d, 新值=%d\n", old, atomic.LoadInt64(&current))

	// Swap 常用于"取走并重置"场景
	var pending int64 = 5
	taken := atomic.SwapInt64(&pending, 0)
	fmt.Printf("取走并重置: 取走=%d, 剩余=%d\n", taken, atomic.LoadInt64(&pending))

	fmt.Println()
}

// ---- 并发对比：atomic vs Mutex ----

// 以下代码不在 main 中执行，仅作演示说明。
// 可以看到 concurrentAddWithAtomic 和 concurrentAddWithMutex 的执行方式差异。

func concurrentAddWithAtomic() {
	var counter int64
	var wg sync.WaitGroup
	n := 1000

	for range n {
		wg.Add(1)
		go func() {
			atomic.AddInt64(&counter, 1) // 无锁
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("atomic 并发累加 %d 次: %d\n", n, counter)
}

func concurrentAddWithMutex() {
	var counter int64
	var mu sync.Mutex
	var wg sync.WaitGroup
	n := 1000

	for range n {
		wg.Add(1)
		go func() {
			mu.Lock()
			counter++ // 有锁
			mu.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("Mutex 并发累加 %d 次: %d\n", n, counter)
}
