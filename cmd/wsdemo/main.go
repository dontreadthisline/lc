// wsdemo 是 http-vs-websocket 学习 demo 的可执行入口。
//
// 用法:
//
//	go run ./cmd/wsdemo -mode server   # 启动 WS(:9090)+ HTTP 对照(:9091)
//	go run ./cmd/wsdemo -mode client   # 连接 WS 服务端,观察双向通信
//	go run ./cmd/wsdemo -mode raw      # 裸 TCP(:9092),观察握手前原始字节
package main

import "demo/http/ws"

func main() {
	ws.Run(nil)
}
