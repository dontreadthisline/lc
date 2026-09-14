// Package tcpduplex 用裸 TCP 演示"全双工"到底长什么样。
//
// 核心论点(承接 http/ws 的讨论):
//
//	TCP 本身就是全双工的——一条连接 = 一个 socket fd,
//	内核给它配了【独立的】发送缓冲区和接收缓冲区。
//	read 操作接收缓冲区,write 操作发送缓冲区,两者各走各的,
//	所以可以并发,互不阻塞。
//
//	但这个能力"看不见"——raw_tcp.go 只读不写,全双工就潜伏着。
//	只有当代码真的【同时读写两个方向】时,全双工才显现。
//
//	本 demo 就是把这件事做出来给你看:
//	  - 服务端:一个 goroutine 死读客户端消息,另一个 goroutine 每 2s 主动推一条
//	  - 客户端:一个 goroutine 死读服务端消息,主循环每 3s 主动发一条
//	两边都是"边发边收",且发和收跑在不同 goroutine,共用同一个 net.Conn。
//
// 运行:
//
//	终端1: go run ./cmd/tcpduplex -mode server   (监听 :9093)
//	终端2: go run ./cmd/tcpduplex -mode client   (连 127.0.0.1:9093)
//
// 观察日志:你会看到服务端/客户端的"发送"和"接收"日志交错出现,
// 且一方的发送不会卡住另一方的接收——这就是全双工。
package tcpduplex

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

const (
	serverAddr = "127.0.0.1:9093"
)

// Run 是 demo 的统一入口,args 为 nil 时读 os.Args[1:]。
func Run(args []string) {
	if args == nil {
		args = os.Args[1:]
	}
	fs := flag.NewFlagSet("tcp-duplex", flag.ExitOnError)
	var mode string
	fs.StringVar(&mode, "mode", "server", "server | client")
	_ = fs.Parse(args)

	switch mode {
	case "server":
		runServer()
	case "client":
		runClient()
	default:
		fmt.Printf("unknown mode: %s\n", mode)
	}
}

// ============================================================
// 服务端
// ============================================================

func runServer() {
	ln, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatalf("[server] listen %s: %v", serverAddr, err)
	}
	log.Printf("[server] listening on %s (裸 TCP 全双工)", serverAddr)
	log.Printf("[server] 另开终端: go run ./cmd/tcpduplex -mode client")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("[server] accept err: %v", err)
			continue
		}
		// 每个连接一个 goroutine,内部再拆成"读 goroutine + 写 goroutine"。
		go handleConn(conn, "server")
	}
}

// handleConn 是全双工的核心:把一条 conn 拆成读、写两条并发流水线。
// 服务端和客户端共用这套逻辑,只是日志前缀不同。
//
// 这里没有任何 HTTP 语义——没有请求-响应配对,没有"发完等收"。
// 读和写是两条独立的循环,各自想发就发、想读就读。
func handleConn(conn net.Conn, who string) {
	defer conn.Close()
	remote := conn.RemoteAddr()
	log.Printf("[%s] connected: %s", who, remote)

	// 全双工的代价:读、写两条流水线各自独立运行,互相看不见对方的状态。
	// 所以一方结束(对端关连接 / 出错)时,必须【显式】把另一方也拽出来,
	// 否则就会半死不活——比如读已经 EOF 了,写还在睡 ticker 等下一拍。
	//
	// done 就是这条联动线:谁先退谁 close(done);另一方的 select 会立刻命中
	// done 分支退出。close 用 once 保护,因为读、写都可能先退,重复 close 会 panic。
	var wg sync.WaitGroup
	wg.Add(2)
	done := make(chan struct{})
	var once sync.Once
	notifyDone := func() { once.Do(func() { close(done) }) }

	// ---- 读流水线:持续读对端发来的字节 ----
	go func() {
		defer wg.Done()
		defer notifyDone() // 读退出 → 通知写流水线;外层 defer conn.Close() 会顶开写的阻塞 Write
		r := bufio.NewReader(conn)
		for {
			line, err := r.ReadBytes('\n')
			if len(line) > 0 {
				// 注意:这里在"读",同时另一个 goroutine 可能在"写",
				// 两者操作同一个 conn 但不同的内核缓冲区,互不阻塞。
				log.Printf("[%s] recv from %s: %q", who, remote, string(line))
			}
			if err != nil {
				log.Printf("[%s] read end: %v", who, err)
				return
			}
		}
	}()

	// ---- 写流水线:持续主动往对端推数据 ----
	// 这是全双工的另一半:服务端不需要"被请求"才发,自己定时就发。
	go func() {
		defer wg.Done()
		defer notifyDone() // 写退出 → 通知读流水线;读若正阻塞在 ReadBytes,会被 conn.Close() 顶开
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		i := 0
		for {
			select {
			case <-done:
				// 对端(读流水线)已经结束,没必要再推了,立刻收尾。
				log.Printf("[%s] write end: peer closed", who)
				return
			case <-ticker.C:
			}
			i++
			msg := fmt.Sprintf("%s push #%d (主动推送,无需对端请求)\n", who, i)
			// 这里在"写",读 goroutine 同时可能也在读,互不阻塞。
			if _, err := conn.Write([]byte(msg)); err != nil {
				log.Printf("[%s] write end: %v", who, err)
				return
			}
			log.Printf("[%s] sent: %q", who, msg)
		}
	}()

	wg.Wait()
	log.Printf("[%s] connection closed: %s", who, remote)
}

// ============================================================
// 客户端
// ============================================================

func runClient() {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatalf("[client] dial %s: %v", serverAddr, err)
	}
	log.Printf("[client] connected to %s", serverAddr)

	// 客户端复用 handleConn:它同样起"读 goroutine + 写 goroutine"。
	// 于是客户端也能边发边收——这正是 HTTP client 抽象给不了你的东西。
	// http.Client 只有 Do(req) 这种"发完阻塞等收"的同步接口,
	// 它压根不暴露"给我一个句柄、我同时读写"的模型。
	//
	// 这里的 conn 就是你上一轮说的"句柄":拿到它,你自己开两个 goroutine
	// 边发边收,不再受请求-响应时序的约束。
	handleConn(conn, "client")

	// 注意:handleConn 内部的写流水线是"每 2s 推一条"。
	// 如果你想让客户端发得更频繁,可以单独写个客户端版,这里为了复用代码
	// 直接用同一套。观察日志时,服务端和客户端的收发会交错出现。
}
