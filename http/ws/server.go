package ws

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

// wsAddr / httpAddr 是对照服务的监听地址。
const (
	wsAddr   = ":9090"
	httpAddr = ":9091"
)

// runServer 同时启动两个服务:
//   - :9090  WebSocket 服务端(全双工)
//   - :9091  普通 HTTP 服务端(请求-响应,做对照)
//
// 同样的"服务端主动给客户端推一条消息"这件事,在两个服务里表现完全不同。
func runServer() {
	// ---- WebSocket 服务端 ----
	// websocket.Handler 是 golang.org/x/net/websocket 提供的最简封装。
	// 注意它的签名:func(ws *websocket.Conn)。一旦握手完成,这个 goroutine
	// 就拿到了一个全双工的 ws.Conn,读写可以并发进行 —— 这正是和 HTTP 的本质区别。
	http.Handle("/ws", websocket.Handler(echoAndPush))

	// ---- HTTP 对照服务端 ----
	// 同样的"服务端想主动通知客户端"需求,HTTP 只能用两种笨办法:
	//   1. 客户端轮询(poll): 客户端反复 GET /poll,服务端有消息才返回
	//   2. 长轮询(long poll): 服务端 hold 住请求直到有消息或超时
	// 这里用最直白的轮询做对照。
	http.HandleFunc("/poll", handlePoll)

	go func() {
		log.Printf("WebSocket server listening on %s/ws", wsAddr)
		if err := http.ListenAndServe(wsAddr, nil); err != nil {
			log.Fatal(err)
		}
	}()

	log.Printf("HTTP compare server listening on %s/poll", httpAddr)
	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		log.Fatal(err)
	}
}

// echoAndPush 演示 WebSocket 的双向能力:
//   - 读循环: 收客户端发来的消息(echo 回去)
//   - 写循环: 服务端每 2 秒主动推一条心跳(无需客户端先请求)
//
// 这两条逻辑跑在不同 goroutine 里,共用同一个 ws.Conn,
// 这在 HTTP/1.1 的请求-响应模型下是做不到的。
func echoAndPush(ws *websocket.Conn) {
	// 重要:握手已经在进入这个函数之前完成。ws.Conn 现在是一条
	// 被升级过的全双工 TCP 连接,不再受 HTTP 语义约束。
	remote := ws.Request().RemoteAddr
	log.Printf("[ws] client connected: %s", remote)

	// 服务端主动推送 goroutine —— 这是 HTTP 做不到的核心能力。
	// 它不依赖任何客户端请求,自己想发就发。
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		i := 0
		for {
			select {
			case <-ticker.C:
				i++
				msg := fmt.Sprintf("server-push #%d (服务端主动推送,无请求触发)", i)
				// websocket.Message.Send 会把字符串封装成 WebSocket 文本帧。
				// 一帧 = 一个带 opcode 的二进制块,详见 docs/http-vs-websocket.md。
				if err := websocket.Message.Send(ws, msg); err != nil {
					log.Printf("[ws] push send err: %v", err)
					return
				}
			}
		}
	}()

	// 读循环:阻塞读客户端发来的消息。注意这个读和上面的推送写
	// 是并发的,互不阻塞 —— 这就是"全双工"在代码里的样子。
	for {
		var msg string
		// websocket.Message.Receive 负责解帧、拼接、去掉 mask。
		if err := websocket.Message.Receive(ws, &msg); err != nil {
			if err == io.EOF {
				log.Printf("[ws] client %s disconnected", remote)
			} else {
				log.Printf("[ws] receive err: %v", err)
			}
			return
		}
		log.Printf("[ws] recv from %s: %q", remote, msg)
		// echo 回去。注意:此时推送 goroutine 可能也在写同一个 ws.Conn。
		// x/net/websocket 的 Conn 内部有写锁,并发写是安全的。
		if err := websocket.Message.Send(ws, "echo: "+msg); err != nil {
			log.Printf("[ws] echo send err: %v", err)
			return
		}
	}
}

// handlePoll 是 HTTP 对照:服务端想"通知"客户端,只能等客户端来问。
// 客户端必须反复 GET /poll,服务端才能把消息"塞"进响应里返回。
// 消息的发起方始终是客户端,服务端永远被动应答。
func handlePoll(w http.ResponseWriter, r *http.Request) {
	// 模拟"服务端当前是否有新消息"。
	// 真实场景里这通常是个队列/时间戳比较,这里用偶数秒简化。
	now := time.Now().Second()
	if now%2 == 0 {
		fmt.Fprintf(w, "msg at second %d (HTTP: 客户端问一次,服务端才答一次)\n", now)
		return
	}
	// 没消息也得返回一个响应 —— HTTP 语义要求每个请求必须有响应。
	fmt.Fprintln(w, "no new message (HTTP: 即使没消息也要回一个响应,白白浪费一次往返)")
}
