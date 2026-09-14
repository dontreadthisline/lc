// tcpduplex 是裸 TCP 全双工 demo 的可执行入口。
//
// 用法:
//
//	go run ./cmd/tcpduplex -mode server   # 启动服务端,监听 127.0.0.1:9093
//	go run ./cmd/tcpduplex -mode client   # 连接服务端,观察边发边收
package main

import "demo/http/tcpduplex"

func main() {
	tcpduplex.Run(nil)
}
