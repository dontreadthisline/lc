// Package ws 演示 WebSocket 与 HTTP 的根本差异。
//
// 运行: go run ./http/ws -mode server   (启动 WS + HTTP 对照服务)
//       go run ./http/ws -mode client   (连接 WS 服务端)
//       go run ./http/ws -mode raw      (裸 TCP 对照,观察握手前字节流)
package ws

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// Run 是 demo 的统一入口。
// args 为 nil 时读取 os.Args[1:],便于直接 go run ./cmd/wsdemo -mode xxx。
func Run(args []string) {
	if args == nil {
		args = os.Args[1:]
	}
	fs := flag.NewFlagSet("ws-demo", flag.ExitOnError)
	var mode string
	fs.StringVar(&mode, "mode", "server", "server | client | raw")
	_ = fs.Parse(args)

	switch mode {
	case "server":
		runServer()
	case "client":
		runClient()
	case "raw":
		runRawTCP()
	default:
		fmt.Printf("unknown mode: %s\n", mode)
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
