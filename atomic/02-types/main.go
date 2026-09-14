// 02-types: Go 1.19+ 类型化原子变量
//
// Go 1.19 引入了 atomic.Int32、atomic.Int64、atomic.Bool 等类型，
// Go 1.20 加入了 atomic.Pointer[T]（1.19 也有，但 API 有调整），
// Go 1.23 加入了 atomic.Uintptr。
//
// 相比旧的 atomic.AddInt64(&x, 1) 函数式 API，新 API 的优势：
//   - 类型安全：atomic.Int64.Swap(0) 编译期检查，不会传错指针类型
//   - 更易读：counter.Add(1) vs atomic.AddInt64(&counter, 1)
//   - 零值可用：var c atomic.Int64 不需要额外初始化
//   - 防止误用：不会把普通 int64 直接当原子变量用（必须通过方法访问）
//
// 可用的类型化原子变量：
//   atomic.Bool             布尔值
//   atomic.Int32            有符号 32 位
//   atomic.Uint32           无符号 32 位
//   atomic.Int64            有符号 64 位
//   atomic.Uint64           无符号 64 位
//   atomic.Uintptr          指针大小的无符号整数 (Go 1.23+)
//   atomic.Pointer[T]       任意类型指针 (Go 1.19+)
//   atomic.Value            任意类型值（见 03-value）
//
// 运行: go run ./atomic/02-types/

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	fmt.Print("=== 02 类型化原子变量 ===\n\n")

	boolType()
	int64Type()
	pointerType()
	zeroValueSafety()
}

// ---- atomic.Bool ----

func boolType() {
	var flag atomic.Bool

	fmt.Printf("Bool 零值: %v\n", flag.Load())

	flag.Store(true)
	fmt.Printf("Store(true) 后: %v\n", flag.Load())

	// CAS 风格的状态切换
	swapped := flag.CompareAndSwap(true, false)
	fmt.Printf("CAS true→false 成功? %v, 当前: %v\n", swapped, flag.Load())

	// Swap：设新值，返回旧值
	old := flag.Swap(true)
	fmt.Printf("Swap(true): 旧值=%v\n", old)

	fmt.Println()
}

// ---- atomic.Int64 ----

func int64Type() {
	var counter atomic.Int64

	// Add: 原生加减，返回新值
	fmt.Printf("Add(1): %d\n", counter.Add(1))
	fmt.Printf("Add(10): %d\n", counter.Add(10))

	// Load / Store
	fmt.Printf("Load: %d\n", counter.Load())
	counter.Store(0)
	fmt.Printf("Store(0) 后: %d\n", counter.Load())

	// CAS 自旋累加到目标值
	const target = int64(50)
	for {
		old := counter.Load()
		if old >= target {
			break
		}
		counter.CompareAndSwap(old, old+1)
	}
	fmt.Printf("CAS 自旋到 %d: %d\n", target, counter.Load())

	// Swap
	old := counter.Swap(999)
	fmt.Printf("Swap(999): 旧值=%d\n", old)

	fmt.Println()
}

// ---- atomic.Pointer[T] ----

type Config struct {
	Host string
	Port int
}

func pointerType() {
	var cfgPtr atomic.Pointer[Config]

	// 零值: nil
	if cfgPtr.Load() == nil {
		fmt.Println("Pointer 零值为 nil")
	}

	// Store 和 Load
	cfgPtr.Store(&Config{Host: "localhost", Port: 8080})
	cfg := cfgPtr.Load()
	fmt.Printf("Load: %s:%d\n", cfg.Host, cfg.Port)

	// CAS 原子替换
	oldCfg := cfgPtr.Load()
	newCfg := &Config{Host: "0.0.0.0", Port: 9090}
	if cfgPtr.CompareAndSwap(oldCfg, newCfg) {
		fmt.Println("CAS 替换配置成功")
	}
	fmt.Printf("替换后: %s:%d\n", cfgPtr.Load().Host, cfgPtr.Load().Port)

	// Swap
	swappedOut := cfgPtr.Swap(&Config{Host: "prod.example.com", Port: 443})
	fmt.Printf("Swap 替换掉的旧配置: %s:%d\n", swappedOut.Host, swappedOut.Port)

	// ---- 重要：atomic.Pointer 不保护指向的对象 ----
	// Store 只是原子地替换了指针本身，不保护 *Config 内部的字段。
	// 如果在 Store 之后其他 goroutine 修改了 Config.Host，仍然会有 data race。
	//
	// 如果需要保护结构体内部字段，要么：
	//   1. 每次写整体替换（immutable 风格）——推荐
	//   2. 用 sync.Mutex
	//   3. 用 atomic.Value（见 03-value）

	fmt.Println()
}

// ---- 零值可用性 ----

func zeroValueSafety() {
	// 所有类型化原子变量都是零值可用的，无需初始化
	var (
		b   atomic.Bool
		i32 atomic.Int32
		i64 atomic.Int64
		u32 atomic.Uint32
		u64 atomic.Uint64
		p   atomic.Pointer[string]
	)

	// 零值行为
	fmt.Printf("Bool 零值: %v\n", b.Load())
	fmt.Printf("Int32 零值: %d\n", i32.Load())
	fmt.Printf("Int64 零值: %d\n", i64.Load())
	fmt.Printf("Uint32 零值: %d\n", u32.Load())
	fmt.Printf("Uint64 零值: %d\n", u64.Load())
	fmt.Printf("Pointer 零值: %v\n", p.Load())

	// 可以直接使用，不需要 &b 取地址
	b.Store(true)
	i64.Add(100)

	fmt.Printf("使用后: b=%v, i64=%d\n", b.Load(), i64.Load())

	// ---- 注意：atomic 类型不能复制 ----
	// var b2 atomic.Bool = b   // 错误！atomic 类型包含 noCopy 字段，vet 会报警
	//                           // go vet 会提示: "assignment copies lock value"
	//
	// 正确的做法：始终通过指针传递
	// func useFlag(f *atomic.Bool) { f.Store(true) }

	fmt.Println()
}

// ---- 常见并发模式 ----

// 状态机：用 atomic.Int32 做轻量状态
type Connection struct {
	state atomic.Int32
}

const (
	stateDisconnected = iota
	stateConnecting
	stateConnected
	stateClosed
)

func (c *Connection) TryConnect() bool {
	// 只有 Disconnected 状态才能开始连接
	return c.state.CompareAndSwap(stateDisconnected, stateConnecting)
}

func (c *Connection) MarkConnected() {
	c.state.Store(stateConnected)
}

func (c *Connection) Close() {
	c.state.Store(stateClosed)
}

func (c *Connection) IsConnected() bool {
	return c.state.Load() == stateConnected
}

// 这个模式的局限性：
//   - 状态多了之后 CAS 条件变复杂，容易出错
//   - 没有等待/通知机制（不像 channel 可以阻塞等待状态变更）
//   - 如果状态之间有依赖关系，建议用 sync.Mutex + sync.Cond 或 channel

// ---- 生产建议 ----

// goroutine 安全关闭标志 —— atomic.Bool 的经典用法
type Worker struct {
	closed atomic.Bool
	wg     sync.WaitGroup
}

func (w *Worker) Stop() {
	w.closed.Store(true)
	w.wg.Wait() // 等待工作协程退出
}

func (w *Worker) run() {
	defer w.wg.Done()
	for {
		if w.closed.Load() {
			return
		}
		// 做实际工作...
	}
}
