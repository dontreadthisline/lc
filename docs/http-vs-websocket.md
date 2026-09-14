# HTTP vs WebSocket:从协议层到 TCP 层的根本差异

> 配套 demo: `go run ./cmd/wsdemo -mode {server|client|raw}`
> 依赖: `golang.org/x/net/websocket`

## 0. 先纠正一个核心误解

很多人说"HTTP 是单向的,WebSocket 是双向的"。这句话**只对了一半**,而且容易让人误以为是 TCP 的限制。真相是:

> **HTTP 和 WebSocket 都跑在同一条全双工 TCP 连接上。TCP 从来就是双向的。"单向/双向"是应用层协议的语义约定,不是传输层的能力差异。**

- TCP 给你一根**双向管道**:客户端能写、服务端也能写,两边可以同时写,互不等待。
- HTTP/1.1 在这根管道上**强加了一个规则**:"客户端先说话(请求),服务端才能回话(响应),一回完这轮就结束"。规则是人定的,不是管道的物理限制。
- WebSocket 做的事:**先用 HTTP 的规则完成一次握手,然后宣布"从现在起我们不按 HTTP 规则玩了"**,把这根 TCP 管道还原成它本来的全双工样子,双方想什么时候发就什么时候发。

一句话总结:

| | 谁能主动发 | 受谁约束 |
|---|---|---|
| HTTP/1.1 | 只有客户端能发起请求,服务端只能应答 | **应用层语义**约束 |
| WebSocket | 双方都能随时发 | 握手后**摆脱** HTTP 语义,回归 TCP 全双工 |
| TCP(底层) | 双方都能随时发 | 物理上就是全双工,谁也不约束谁 |

所以你问的"同样是 TCP,为什么一个单向一个双向"——答案是:**TCP 一直是双向的,HTTP 是自己选择表现得像单向的,WebSocket 选择不装了。**

---

## 1. 自上而下:三层架构

```
┌─────────────────────────────────────────────────┐
│  应用层(你的业务代码)                          │
│  WS: websocket.Message.Send / Receive           │
│  HTTP: http.Get / http.HandleFunc               │
├─────────────────────────────────────────────────┤
│  协议层                                         │
│  WS: WebSocket 帧协议(opcode/mask/payload)      │
│  HTTP: HTTP/1.1 文本协议(请求行/头/体)         │
├─────────────────────────────────────────────────┤
│  传输层                                         │
│  TCP: 字节流,全双工,带 seq/ack 的可靠传输      │
└─────────────────────────────────────────────────┘
```

下面逐层往下钻,每层都配 demo 里的对应代码。

---

## 2. 协议层:HTTP 请求-响应 vs WebSocket 握手+帧

### 2.1 HTTP/1.1 的一次通信(对照 demo 的 `/poll`)

HTTP 是**文本行协议**,每行以 `\r\n` 结尾。一次完整的 HTTP 请求长这样:

```http
GET /poll HTTP/1.1\r\n
Host: localhost:9091\r\n
\r\n
```

服务端必须回一个响应,否则这个请求就"悬着":

```http
HTTP/1.1 200 OK\r\n
Content-Type: text/plain\r\n
Content-Length: 42\r\n
\r\n
msg at second 14 (HTTP: 客户端问一次,服务端才答一次)
```

关键点:**这一轮结束后,语义上对话就结束了。** 服务端想"通知"客户端?它没法主动开口——只能等客户端再发一个 GET。这就是 demo 里 `handlePoll` 的处境:

```
客户端                          服务端
  │  GET /poll                    │
  │ ─────────────────────────────>│
  │  200 OK (有消息/没消息)        │
  │ <─────────────────────────────│
  │                               │   ← 这轮结束,服务端想说话也说不出去
  │  GET /poll (再来问一次)        │
  │ ─────────────────────────────>│
  │  200 OK                       │
  │ <─────────────────────────────│
```

要实现"服务端推送",HTTP 只能靠**轮询**(客户端不停问)或**长轮询**(服务端 hold 住请求)。两者都在和 HTTP 的"请求-响应"语义作对,效率低、延迟高。

### 2.2 WebSocket 的握手:借 HTTP 的壳

WebSocket 没有另起炉灶建新连接,它**复用 HTTP 来完成握手**。客户端发的是一个长得像 HTTP 的请求:

```http
GET /ws HTTP/1.1\r\n
Host: localhost:9090\r\n
Upgrade: websocket\r\n
Connection: Upgrade\r\n
Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==\r\n
Sec-WebSocket-Version: 13\r\n
\r\n
```

服务端如果同意升级,回一个 **101 Switching Protocols**(注意状态码不是 200):

```http
HTTP/1.1 101 Switching Protocols\r\n
Upgrade: websocket\r\n
Connection: Upgrade\r\n
Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=\r\n
\r\n
```

`Sec-WebSocket-Accept` 是服务端用 `Key + 一个固定魔法字符串` 做 SHA-1 再 base64,证明"我确实懂 WebSocket 协议,不是随便回的"。

**101 响应发完,这条 TCP 连接的性质就变了**:它不再是 HTTP 连接,而是一条裸的、双方都能随时写的全双工通道。从这一刻起,HTTP 的请求-响应规则不再适用。

> 这正是 demo 里 `echoAndPush` 函数能同时跑"读循环"和"推送 goroutine"的前提——握手已完成,连接升级完毕。

### 2.3 WebSocket 帧:握手之后怎么传数据

握手之后,双方不再发 HTTP 文本,而是发**二进制的 WebSocket 帧**。一帧的结构(RFC 6455):

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-------+-+-------------+-------------------------------+
|F|R|R|R| opcode|M| Payload len |    Extended payload length    |
|I|S|S|S|  (4)  |A|     (7)     |             (16/64)           |
|N|V|V|V|       |S|             |   (if payload len==126/127)   |
| |1|2|3|       |K|             |                               |
+-+-+-+-+-------+-+-------------+ - - - - - - - - - - - - - - - +
|     Extended payload length continued, if payload len == 127  |
+ - - - - - - - - - - - - - - - +-------------------------------+
|                               |Masking-key, if MASK set to 1  |
+-------------------------------+-------------------------------+
| Masking-key (continued)       |          Payload Data         |
+-------------------------------- - - - - - - - - - - - - - - - +
:                     Payload Data continued ...                :
+ - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - +
|                     Payload Data continued ...                |
+---------------------------------------------------------------+
```

几个关键字段:

- **FIN (1 bit)**:这帧是不是一个完整消息的最后一帧。WebSocket 支持把一个大消息拆成多帧。
- **opcode (4 bits)**:这帧是什么类型。
  - `0x1` 文本帧 / `0x2` 二进制帧 / `0x8` 关闭帧 / `0x9` ping / `0xA` pong
- **MASK (1 bit)**:客户端发的帧必须带 mask(用 masking-key 异或 payload),服务端发的帧不带。这是为了防止中间代理被混淆,不是加密。
- **Payload length**:7 位 / 7+16 位 / 7+64 位三档,自适应消息大小。

demo 里你调用的 `websocket.Message.Send(ws, msg)`,库做的就是:把 `msg` 编码成 UTF-8,套上 `opcode=0x1`、算好长度、(客户端)生成 mask 异或 payload,拼成一帧写进 TCP。`Receive` 则反过来拆帧。

### 2.4 对照时序图

**HTTP 轮询(服务端推送的笨办法):**

```
客户端                                        服务端
  │                                            │
  │── GET /poll ──────────────────────────────>│  客户端主动问
  │<── 200 OK "no new message" ───────────────│  服务端被动答
  │                                            │
  │   (服务端此刻有消息了,但发不出去!)        │
  │                                            │
  │── GET /poll ──────────────────────────────>│  客户端再来问
  │<── 200 OK "msg at second 14" ─────────────│  这才把消息带回去
  │                                            │
  延迟 = 轮询间隔;大量空请求浪费带宽
```

**WebSocket(服务端推送的自然办法):**

```
客户端                                        服务端
  │                                            │
  │── HTTP GET (Upgrade: websocket) ─────────>│  握手(借 HTTP 的壳)
  │<── 101 Switching Protocols ───────────────│  握手完成,连接升级
  │                                            │
  │ ============ 从此全双工,无请求-响应约束 ============ │
  │                                            │
  │<── server-push #1 (主动推送) ─────────────│  服务端想发就发
  │── hello from client #1 ──────────────────>│  客户端想发就发
  │<── echo: hello from client #1 ────────────│  两边交错,互不等待
  │<── server-push #2 (主动推送) ─────────────│
  │── hello from client #2 ──────────────────>│
  │                                            │
  延迟 ≈ 网络往返;无空请求;服务端可真正主动推送
```

---

## 3. TCP 层:为什么"全双工"不是 WebSocket 的功劳

### 3.1 TCP 本来就是全双工的

TCP 连接建立后(三次握手),内核为这条连接维护**两个方向的独立数据流**:

```
客户端                              服务端
  │  ──── 发送队列 ────>            │
  │  <─── 接收队列 ────             │
  │                                │
  │  <─── 发送队列 ────             │
  │  ──── 接收队列 ────>            │
```

两个方向各有自己的序号(seq)、确认号(ack)、窗口大小。服务端往客户端写数据,根本不需要客户端"先请求"——TCP 层面它随时能写。

### 3.2 那 HTTP 为什么显得"单向"?

因为 HTTP/1.1 的**应用层状态机**规定:一条(逻辑上的)对话由客户端的请求触发,服务端的响应结束。在这个状态机里:

- 服务端不能在"没有收到请求"的时候主动往连接里写响应数据(写了客户端也解析不了,它正等着自己的请求的响应呢)。
- 虽然底层 TCP 允许服务端随时写,但 HTTP 解析器在客户端那边不认这种"无请求的响应"。

所以"单向"是**应用层解析规则的约束**,不是 TCP 的限制。HTTP/1.1 还有个 keep-alive,允许在一条 TCP 上连续发多个请求-响应,但依然是"请求驱动响应"的模式,不是真正的双向。

### 3.3 WebSocket 如何"解放"这条 TCP

WebSocket 握手成功后,双方约定:**之后这条连接上的字节流,不再按 HTTP 规则解析,改按 WebSocket 帧规则解析。** 而 WebSocket 帧规则是**对称的**——客户端能发帧,服务端也能发帧,谁也不用等谁。

所以与其说"WebSocket 实现了全双工",不如说"**WebSocket 没有像 HTTP 那样人为限制全双工**"。全双工能力 TCP 一直给着,WebSocket 只是把它暴露出来用了。

---

## 4. 裸 TCP 对照实验(回答"握手前是什么样")

跑 `go run ./cmd/wsdemo -mode raw`,它用 `net.Listen` 监听 :9092,把收到的字节原样打印。然后用任意 ws 客户端连 `ws://localhost:9092/any`。

你会看到服务端打印出的第一波字节**是一个 HTTP 请求**:

```
[raw-tcp] line#1 from 127.0.0.1:xxxxx:
  text: "GET /any HTTP/1.1\r\n"
  hex : 474554202f616e7920485454502f312e310d0a
[raw-tcp] line#2 from 127.0.0.1:xxxxx:
  text: "Host: localhost:9092\r\n"
  hex : 486f73743a206c6f63616c686f73743a393039320d0a
...
[raw-tcp] line#N:
  text: "Upgrade: websocket\r\n"
...
[raw-tcp] line#N+1:
  text: "Sec-WebSocket-Key: xxxxx==\r\n"
...
[raw-tcp] >>> 空行,HTTP 请求头结束 <<<
```

这个实验直观证明三件事:

1. **WebSocket 握手就是 HTTP**:客户端发的第一波字节,和你用 `curl` 发的 HTTP 请求是同一种东西,只是多了几个 `Upgrade`/`Sec-WebSocket-*` 头。
2. **TCP 不认识 HTTP 也不认识 WebSocket**:对 `net.Listen` 来说,它收到的就是一堆字节。是**两端的应用层代码**约定了怎么解析这些字节。TCP 只负责把字节可靠地、按序地搬过去。
3. **"协议"本质上是双方对字节流的解释约定**:同一个 TCP 连接,前半段按 HTTP 解释(握手),后半段按 WebSocket 帧解释(数据)。底层 TCP 全程无感。

你也可以用 nc 直接发个普通 HTTP 请求验证 TCP 的"协议无关":

```bash
printf 'GET / HTTP/1.1\r\nHost: x\r\n\r\n' | nc localhost 9092
```

TCP 照单全收,raw 服务端照样打印——因为它根本不关心你发的是 HTTP 还是别的。

---

## 5. 三层对照总表

| 维度 | HTTP/1.1 | WebSocket | TCP |
|---|---|---|---|
| 谁能主动发数据 | 仅客户端发起请求 | 双方随时可发 | 双方随时可发 |
| 通信模型 | 请求-响应(一问一答) | 消息帧(对称双向) | 字节流(对称双向) |
| 数据格式 | 文本行(`\r\n` 分隔) | 二进制帧(opcode+mask+payload) | 无格式字节流 |
| 服务端推送能力 | 无(需轮询/长轮询/SSE) | 原生支持 | 不适用(传输层无此概念) |
| 连接生命周期 | 一轮请求-响应后可复用(keep-alive)但语义仍是一问一答 | 握手后长期保持全双工 | 建立后一直双向可用,直到关闭 |
| 连接复用 | keep-alive 复用 TCP,但仍请求驱动 | 握手后单条 TCP 持续全双工 | 一条连接天然持续可用 |
| demo 对应代码 | `handlePoll` / `http.Get` | `echoAndPush` / `websocket.Message` | `net.Listen` / `raw_tcp.go` |

---

## 6. 何时用哪个

- **用 HTTP**:请求-响应模式天然契合的场景——拉取资源、提交表单、调用 API、CRUD。多数业务都是这种"客户端主动、偶发、无状态"的交互,HTTP 的语义约束反而是优点(无状态、可缓存、可中间代理)。
- **用 WebSocket**:需要**服务端主动、低延迟、高频双向**的场景——即时聊天、实时协作、行情推送、多人游戏。这类场景用 HTTP 轮询会浪费大量空请求、引入不可控延迟。

记住:**WebSocket 不是"更好的 HTTP",而是为不同交互模式设计的协议。** 它用 HTTP 握手是为了复用现有的 Web 基础设施(端口、代理、防火墙都认识 HTTP),握手完就切到自己的帧协议。

---

## 7. 跑 demo 的完整流程

开三个终端:

```bash
# 终端1:启动 WS + HTTP 对照服务
go run ./cmd/wsdemo -mode server

# 终端2:启动客户端,观察双向通信
go run ./cmd/wsdemo -mode client

# 终端3(可选):裸 TCP,观察握手前字节
go run ./cmd/wsdemo -mode raw
# 然后另开: printf 'GET / HTTP/1.1\r\nHost: x\r\n\r\n' | nc localhost 9092
```

观察重点:
- 终端2 里 `[push]` 和 `[echo]` 消息交错出现,且 `[push]` 不依赖你发任何请求——这就是服务端主动推送。
- 对照 `curl http://localhost:9091/poll` 反复请求,体会 HTTP"问一次答一次"的被动。
- 终端3 里看到的原始字节,证明 WS 握手就是 HTTP、TCP 对协议无感。
