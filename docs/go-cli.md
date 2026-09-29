# Go 命令行教程：从跑通到日用

> 基于本仓库真实操作整理（demo 模块，Go 1.24）。`go test` 详见《go test 命令详解》部分（对话记录）与 `go help testflag`，本文不重复。

## 0. 全景：19 个命令分四组

| 组 | 命令 | 用途 |
|---|---|---|
| 跑代码 | run / build / install | 编译执行、产出二进制 |
| 管依赖 | get / mod / work | go.mod 的增删改查 |
| 保质量 | fmt / vet / test | 格式化、静态检查、测试 |
| 查信息 | doc / list / env / version | 文档、包清单、环境 |

其余（clean / generate / fix / tool / bug / telemetry）见第 6 节速览。

前置概念一句话：**一个仓库 = 一个 module（`go.mod` 定义），module 下每个目录是一个 package**。`./...` 通配"当前 module 的所有包"。查全部包：`go list ./...`（本仓库当前 22 个包）。

## 1. 跑代码：run / build / install

### go run —— 编译+运行一步走，不留二进制

```bash
go run main.go        # 单文件（仅限练习，文件多了就错）
go run .              # 当前目录整个 main 包（推荐）
go run demo/cmd/wsdemo  # 跑子包
```

### go build —— 只编译，验错或出产物

```bash
go build ./...        # 编译所有包，不产文件（main 包产物在当前目录）
go build -o /tmp/demo-bin .        # 指定输出路径
GOOS=linux GOARCH=arm64 go build -o /tmp/demo-linux .   # 交叉编译（实测试过，产出 ELF）
go tool dist list     # 查全部可交叉的平台组合
```

- `go build ./...` 是最快的"能不能编译过"检查，失败必须有输出，成功静默
- 验证产物架构：`file /tmp/demo-linux`（实测输出 `ELF 64-bit... aarch64`）

### go install —— 编译并安装到 $GOPATH/bin（或 $GOBIN）

```bash
go install .                    # 安装自己写的 main 包
go install golang.org/x/tools/gopls@latest   # 安装第三方工具（@版本，不依赖当前 module）
```

装完的命令直接敲名字用，前提 `$GOBIN` 在 PATH 里。查路径：`go env GOPATH`。

## 2. 管依赖：go.mod 的增删改查

```bash
go mod init demo          # 建模块（本仓库 module 名就叫 demo，所以包路径都是 demo/xxx）
go get github.com/jackc/pgx/v5        # 加依赖（写进 go.mod）
go get github.com/jackc/pgx/v5@v5.6.0 # 锁定版本
go get -u ./...           # 全部升到最新次级版本
go get xxx@none           # 删依赖
go mod tidy               # 大扫除：补齐缺的、删掉没用的（代码改完必跑）
```

- `go.mod` 是你声明要什么，`go.sum` 是这些依赖的哈希校验（自动维护，别手改）
- 版本查询/依赖树：`go list -m all`（全量版本）、`go mod graph`（依赖图）、`go mod why xxx`（为什么需要它）
- 换源：`go env -w GOPROXY=https://goproxy.cn,direct`

### go work（多仓库联动，暂时用不上先知道）

```bash
go work init ./learn-go ./other-repo   # 生成 go.work，本地多个 module 互相引用不依赖发布
```

## 3. 保质量：fmt / vet / test 闭环（今天刚实战过）

```bash
go fmt ./...    # 格式化（内部就是 gofmt -l -w），无脑跑
go vet ./...    # 静态检查
go test ./...   # 测试
```

**实战记录（2026-09-15，本仓库）**：`go vet ./...` 一把抓出两个真 bug——

1. `main.go` 导入 `"fmt"` 未使用 → 这在 Go 里是**编译错误**（不是警告），整个 `demo` 包直接编译不过
2. `sql/07-pgx-native/main.go` 用 `%s` 打印 bool → 编译不报（Printf 收 interface{}），**运行时**输出 `%!s(bool=true)` 这种乱码 → 应为 `%t`

教训：vet 抓的正是"能编译但有错"的灰区。三件套固定顺序 `fmt → vet → test`，写完就跑。

常用 flag 补充：

```bash
go vet -vettool=none ./...   # 关闭 vet
go test -race ./...          # 竞态检测（并发代码必开）
go test -cover ./...         # 覆盖率
```

## 4. 查信息：doc / list / env

### go doc —— 查 API（本地文档系统里唯一写得像人的部分）

```bash
go doc testing                # 包概览
go doc testing.T              # 类型+全部方法
go doc testing.T.Run          # 单个方法
go doc -all strings           # 整包展开
go doc -src sort.Slice        # 直接看源码
```

### go list —— 包结构的 X 光

```bash
go list ./...                 # 所有包路径
go list -m all                # 全部依赖及版本
go list -deps ./...           # 连标准库依赖一起列
go list -f '{{.Name}} {{.Dir}}' ./...   # 自定义模板输出
```

### go env —— 环境真相

```bash
go env GOPATH GOROOT GOPROXY  # 看单个
go env -changed               # 只看被改过的项（挑重点用）
go env -w GOFLAGS=-count=1    # 永久设置（go env -w GOPROXY=... 同理）
```

## 5. 日常闭环（背下这一行）

```bash
go fmt ./... && go vet ./... && go test ./... && go run .
```

写码 → 四连 → 提交。今天这个闭环刚在你仓库里抓了两个真 bug，就是这么用的。

## 6. 其余命令速览

| 命令 | 一句话 | 什么时候用 |
|---|---|---|
| `go clean -cache` | 清构建缓存 | 磁盘紧张/缓存玄学问题（`-testcache` 只清测试缓存） |
| `go generate ./...` | 执行源码里 `//go:generate` 注释的命令 | 用代码生成器（stringer、mockgen）时 |
| `go fix` | 自动改写废弃 API 用法 | 升级 Go 版本后 |
| `go tool` | 调内置工具（`go tool dist list`、`pprof`、`trace`） | 性能分析时 |
| `go version -m ./二进制` | 查看某二进制的 Go 版本和依赖 | 排查线上产物 |
| `go bug` | 打开 bug 报告页 | —— |
| `go telemetry` | 遥测开关 | 在意隐私就 `go telemetry off` |

## 7. 高频坑位清单

1. `import` 了没用 = 编译错误（不是警告），删掉或用 `_ "pkg"` 显式忽略
2. `go build` 成功静默，失败才说话——别等输出，看退出码
3. `go test ./...` 的 `(cached)` 是特性不是 bug，要真跑 `-count=1`
4. `go run main.go` 只编译单文件，包里多个文件互相调用会报 undefined——用 `go run .`
5. 依赖拉不下来先查 `GOPROXY`（国内 `go env -w GOPROXY=https://goproxy.cn,direct`）
6. 交叉编译的纯 Go 包零配置；包里有 cgo 时要配交叉编译器（先 `CGO_ENABLED=0` 试试）
