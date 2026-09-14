package ws

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

// wsServerURL 是客户端要连的服务端地址,对应 server.go 的 wsAddr。
const wsServerURL = "ws://localhost:9090/ws"

// runClient 连上 WebSocket 服务端,然后:
//   - 起一个 goroutine 专门读服务端推来的消息(包括心跳推送和 echo)
//   - 主循环每 3 秒主动发一条消息给服务端
//
// 注意:客户端发消息和服务端推消息在时间上是交错的、各自独立的。
// 没有任何"客户端先问,服务端才能答"的约束 —— 这就是双向。
func runClient() {
	// websocket.Dial 内部做了三件事:
	//  1. 建立一条普通 TCP 连接
	//  2. 在这条 TCP 上发一个 HTTP GET 请求,带 Upgrade: websocket 头
	//  3. 等服务端回 101 Switching Protocols,握手成功
	// 之后返回的 conn 就是一条升级过的全双工连接。
	conn, err := websocket.Dial(wsServerURL, "", "http://localhost/")
	if err != nil {
		log.Fatalf("dial failed: %v", err)
	}
	defer conn.Close()
	log.Printf("[client] connected to %s", wsServerURL)

	// 读 goroutine:持续收服务端消息。
	// 这些消息里,有的是我主动发消息后服务端 echo 回来的,
	// 有的是服务端定时主动推的(我没发任何请求它也推)。
	go func() {
		for {
			var msg string
			if err := websocket.Message.Receive(conn, &msg); err != nil {
				log.Printf("[client] receive err: %v", err)
				return
			}
			// 用 [push]/[echo] 前缀帮你在日志里区分两类消息的来源。
			tag := "[echo]"
			if len(msg) > 10 && msg[:6] == "server" {
				tag = "[push]"
			}
			fmt.Printf("[client] recv %s: %s\n", tag, msg)
		}
	}()

	// 写循环:每 3 秒主动发一条。这个写和上面的读互不阻塞。
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	i := 0
	for range ticker.C {
		i++
		msg := fmt.Sprintf("hello from client #%d", i)
		if err := websocket.Message.Send(conn, msg); err != nil {
			log.Printf("[client] send err: %v", err)
			return
		}
		fmt.Printf("[client] sent: %s\n", msg)
	}
}
