package ws

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// runRawTCP 开一个最朴素的 net.Listen,不做任何协议处理,
// 把客户端发来的原始字节原样打印出来。
//
// 这个 demo 的目的:让你亲眼看到 ——
// 一个 WebSocket 客户端连过来时,第一波字节其实是【一个 HTTP 请求】。
// 也就是说 WebSocket 并没有另起炉灶,它的握手就是"借 HTTP 的壳"。
//
// 用法:
//   1. 先跑: go run ./http/ws -mode raw      (监听 :9092)
//   2. 另开终端,用任意 ws 客户端连 ws://localhost:9092/anything
//      或者直接用 nc: printf 'GET / HTTP/1.1\r\nHost: x\r\n\r\n' | nc localhost 9092
//   3. 观察打印出的原始字节,你会看到 Upgrade/Connection/Sec-WebSocket-Key 这些头。
const rawTCPAddr = ":9092"

func runRawTCP() {
	ln, err := net.Listen("tcp", rawTCPAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[raw-tcp] listening on %s (把收到的原始字节原样打印)", rawTCPAddr)
	log.Printf("[raw-tcp] 现在用 ws 客户端连 ws://localhost%s/any 或 nc 连过来", rawTCPAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("[raw-tcp] accept err: %v", err)
			continue
		}
		go handleRaw(conn)
	}
}

// handleRaw 用 bufio 逐行读,把每一行的字节十六进制 + 可见字符都打出来。
// 这样你能看清 \r\n 这种控制字符,理解"HTTP 是文本行协议"。
func handleRaw(conn net.Conn) {
	defer conn.Close()
	remote := conn.RemoteAddr()
	log.Printf("[raw-tcp] connection from %s", remote)

	r := bufio.NewReader(conn)
	lineNum := 0
	for {
		// ReadBytes('\n') 读到换行为止,保留 \r\n。
		// 这就是 HTTP/1.1 的逐行解析方式。
		line, err := r.ReadBytes('\n')
		if len(line) > 0 {
			lineNum++
			// 同时打印可见文本和十六进制,让你看到 \r (0x0d) \n (0x0a)。
			fmt.Printf("[raw-tcp] line#%d from %s:\n  text: %q\n  hex : %x\n",
				lineNum, remote, string(line), line)
		}
		if err != nil {
			log.Printf("[raw-tcp] %s read end: %v", remote, err)
			return
		}
		// 空行(\r\n)标志 HTTP 请求头结束。
		if len(line) == 2 && line[0] == '\r' && line[1] == '\n' {
			fmt.Printf("[raw-tcp] >>> 空行,HTTP 请求头结束。后续若有 body/帧数据会继续打印 <<<\n")
		}
	}
}
