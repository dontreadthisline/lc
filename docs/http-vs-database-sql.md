# database/sql vs net/http：实现机制对比

> 两个模块都是"标准库定义接口 + 第三方实现"的插件架构，但入口、分层、控制流的走向截然不同。

---

## 一、一眼看穿的核心差异

```go
// database/sql —— 入口是「具体结构体」
db, _ := sql.Open("pgx", dsn)   // *sql.DB，是个 struct
db.Query("SELECT ...")          // 方法在 DB 结构体上
db.Exec("INSERT ...")           // 用户面向的是 DB 这个"管家"

// net/http —— 入口是「接口」
var h http.Handler = chi.NewRouter()  // 接口类型
h.ServeHTTP(w, r)                     // 用户实现 Handler，Server 来调你
```

database/sql 的入口是**具体类型**，用户拿到的 `*sql.DB` 是一个结构体，方法都在它上面。
net/http 的入口是**接口**，用户实现 `http.Handler`，框架（Server）来调用你。

---

## 二、架构层次对比

```
database/sql                          net/http
─────────────                         ────────

用户代码                               用户代码
  │                                      │
  ▼                                      ▼
*sql.DB  (具体 struct)                 chi.Router  (实现 http.Handler 接口)
  │  - Query/Exec/Prepare               │  - Get/Post/Mount
  │  - Begin/Close                      │  - Use(middleware)
  │  - 连接池管理                        │  - Route
  │                                      │
  ▼                                      ▼
driver.Conn  (接口)                    http.Handler  (接口)
  │  - Prepare/Close/Begin              │  - ServeHTTP(w, r)
  │  - 可选: Pinger                     │
  │  - 可选: SessionResetter            │
  │                                      │
  ▼                                      ▼
pgx/stdlib  (实现)                    业务 handler  (你写的)
  - PG wire 协议编解码                   - listUsers / getUser / createUser
  - 连接建立/关闭
```

关键差异：**sql.DB 是中间层（管家），http.Server 才是中间层**。

```
sql:  用户 → DB(管家,提供连接池) → driver.Conn(接口) → pgx(实现)
http: Server(管家,提供连接管理) → Handler(接口) → chi/gin/业务代码(实现)
```

---

## 三、结构体 vs 接口：为什么一个用 struct，一个用 interface

### database/sql 必须用 struct（`*sql.DB`）

因为 `*sql.DB` 需要**管理状态**：

```go
type DB struct {
    connector   driver.Connector   // 持有的驱动连接器
    freeConn    []*driverConn      // 空闲连接（可变状态）
    mu          sync.Mutex         // 保护并发访问
    numOpen     int                // 当前打开连接数
    maxOpen     int                // 上限
    // ... 大量状态字段
}
```

连接池是有状态的——有空闲列表、有等待队列、有计数器。这些状态必须存在某个地方。`*sql.DB` 就是这个"地方"。暴露成 struct 很自然。

### net/http 用 interface（`http.Handler`）

因为 handler **不需要管理状态**（至少框架不要求它管理）：

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

一个 handler 只需要"处理请求，写响应"。它不需要维护连接池——那是 `http.Server` 的事。框架（Server）持有状态（TCP listener、连接计数），handler 是无状态的函数。

**实际上 http.Server 才是 database/sql 中 DB 的对应物**：

| database/sql | net/http | 职责 |
|-------------|----------|------|
| `*sql.DB` | `*http.Server` | 连接管理、生命周期 |
| `driver.Connector` | `net.Listener` | 底层 I/O 入口 |
| `driver.Conn` | `net.Conn` | 单条连接 |
| 用户操作 DB.Query/Exec | 用户实现 Handler.ServeHTTP | 业务逻辑 |

---

## 四、控制流方向：推 vs 拉

这是最本质的区别。

### database/sql：拉模式（Pull）

```
用户代码（主动发起）
  │ db.Query("SELECT ...")
  ▼
*sql.DB（从池中拉连接）
  │ conn, _ := db.conn(ctx, cachedOrNewConn)
  │ 从 freeConn 取一条，或创建新的，或排队等待
  ▼
driver.Conn（执行查询）
  │ ci.QueryContext(ctx, query, args)
  ▼
pgx（用 PG 协议发送/接收）
```

控制流是**从用户往下推，用户主动发起**。用户说"给我查询"，DB 去连接池拿连接，发给驱动。

### net/http：推模式（Push）

```
net.Conn（数据到达）
  │ 收到 HTTP 请求字节
  ▼
http.Server（解析 HTTP 协议）
  │ 解析 method/URL/headers/body → http.Request
  │ 创建 ResponseWriter
  ▼
http.Handler（调用用户的 handler）
  │ handler.ServeHTTP(w, r)
  ▼
你的 Handler（业务逻辑）
  │ w.Write(response)
```

控制流是**从底层往上推，Server 调用用户**。Server 收到请求后"推"给 handler。handler 是被调用的，不是主动拉取的。

### 形象类比

```
database/sql 像「遥控器」
  你按按钮(db.Query) → 电视内部的电路(连接池)工作 → 信号发出去(驱动发送)

net/http 像「门铃」
  有人按门铃(TCP 数据到达) → 你(handler)被叫去开门 → 你来应门(写响应)
```

---

## 五、可扩展性机制对比

两个模块都用**类型断言**发现可选能力，但应用的层级不同。

### database/sql：在驱动连接层

```go
// sql.DB 内部的判定逻辑
if pinger, ok := dc.ci.(driver.Pinger); ok {
    pinger.Ping(ctx)  // 驱动支持 Ping
}
if resetter, ok := dc.ci.(driver.SessionResetter); ok {
    resetter.ResetSession(ctx)  // 驱动支持会话重置
}
```

对接口的增强是在**单条连接**上。DB 拿到一条连接后，用类型断言检查它有没有额外能力。

### net/http：在响应写入层

```go
// Handler 内部可以检查 ResponseWriter 是否有额外能力
func myHandler(w http.ResponseWriter, r *http.Request) {
    // HTTP/2 Server Push
    if pusher, ok := w.(http.Pusher); ok {
        pusher.Push("/style.css", nil)
    }
    // 流式响应
    if flusher, ok := w.(http.Flusher); ok {
        flusher.Flush()  // 立即将缓冲数据发给客户端
    }
    // WebSocket 升级
    if hijacker, ok := w.(http.Hijacker); ok {
        conn, _, _ := hijacker.Hijack()  // 接管原始 TCP 连接
    }
}
```

对接口的增强是在**ResponseWriter**上。handler 在处理请求时，检查 ResponseWriter 是否支持特殊操作。

### 区别总结

| 维度 | database/sql | net/http |
|------|-------------|----------|
| 类型断言对象 | driver.Conn | http.ResponseWriter |
| 谁做断言 | sql.DB 内部 | 用户的 Handler |
| 增强的能力 | Ping/ResetSession | Push/Flush/Hijack |
| 目的 | 连接优化和健康检查 | 协议级别的特殊处理 |

---

## 六、中间件 vs 可选接口

这是另一个有趣的对比——同一个需求（在不改核心逻辑的前提下叠加能力），两套不同的实现路径。

### net/http 的中间件：包装模式（Decorator）

```go
// 中间件就是一个"Handler 的转换函数"
func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)  // 包一层，前后加逻辑
    })
}

// 使用时层层包裹
handler = Logger(Recoverer(Timeout(myHandler)))
```

每个中间件都实现了 `http.Handler`，形成**洋葱模型**。编译时就确定了调用顺序。

### database/sql 的可选接口：能力发现模式（Capability Detection）

```go
// 不是包装，而是运行时检查"你有没有这个能力"
func (db *DB) pingDC(ctx context.Context, dc *driverConn) error {
    if pinger, ok := dc.ci.(driver.Pinger); ok {
        return pinger.Ping(ctx)  // 有 Ping 能力就用
    }
    return nil  // 没有就算了，不报错
}
```

驱动不需要实现所有接口，框架在运行时发现能力。顺序由框架的调用逻辑决定，驱动只需要声明自己能做什么。

### 为什么不同？

**中间件做的是"行为叠加"，可选接口做的是"能力声明"。**

- 中间件是**主动的**——"在执行你的逻辑之前，我先做点事"。适合对请求/响应做统一的前处理和后处理。
- 可选接口是**被动的**——"如果你有 Ping 的能力，我就用；没有就算了"。适合驱动根据自身数据库特性提供不同的优化路径。

database/sql 不需要中间件模式，因为它的核心操作（Query/Exec/Prepare）都是不可拆分的原子操作，不适合"包一层"。net/http 的请求处理是自由的业务逻辑，需要拦截和增强。

---

## 七、全局注册表：sql.Register vs http.Handle

```go
// database/sql: 全局注册表
sql.Register("pgx", pgxDriver)          // 驱动自注册
db, _ := sql.Open("pgx", dsn)           // 按名查找

// net/http: 全局注册表（仅 DefaultServeMux 使用）
http.HandleFunc("/users", listUsers)    // 注册到全局默认路由器
http.ListenAndServe(":8080", nil)       // nil 使用 DefaultServeMux

// net/http: 局部注册表（第三方路由器，更常见的做法）
r := chi.NewRouter()
r.Get("/users", listUsers)              // 注册到 chi 的路由树
http.ListenAndServe(":8080", r)         // 传入具体 router
```

database/sql 的注册表总是全局的（因为一个进程通常只连一种数据库的少数几个实例）。
net/http 的路由从 Go 1.22 开始也支持了标准库内的模式匹配，但第三方路由器（chi/gin）仍然更灵活。而且 http 的"注册"是可选的——你可以直接把 handler 传给 Server，完全不用注册表。

---

## 八、复杂度模型：协议异构 vs 业务异构

前面讲了那么多机制差异，根因只有一个：**两者面对的复杂度模型完全不同**。

### database/sql：协议异构，操作统一

```
                ┌─────────────┐
                │  用户代码     │
                │ db.Query()   │  ← 不管连什么库，写法一样
                │ db.Exec()    │
                └──────┬──────┘
                       │ 统一接口（struct）
                ┌──────┴──────┐
                │   *sql.DB   │  ← 连接池 + 统一流程
                │              │     prepare → bind → exec → scan
                └──────┬──────┘
                       │ driver.Conn（接口）
          ┌────────────┼────────────┐
          ▼            ▼            ▼
      ┌──────┐    ┌──────┐    ┌──────┐
      │ pgx  │    │mysql │    │sqlite│   ← 协议编解码各不相同
      │ PG   │    │MySQL │    │ C    │
      │ wire │    │protocol│   │ API  │
      └──────┘    └──────┘    └──────┘
```

SQL 操作有一个天然的优势：**CRUD 的流程是高度模板化的**。不管底层连的是 PostgreSQL 还是 MySQL 还是 SQLite，`prepare → bind parameters → execute → scan rows` 这个生命周期是不变的。不同数据库的差异集中在「怎么把参数编码成协议字节流」——这是一个可以下沉到驱动的职责。

所以 database/sql 的策略是：**把复杂性藏在接口下面**。标准库把连接池、重试、上下文管理这些通用能力做成 `*sql.DB` struct，把协议适配下沉到 `driver.Conn` 接口。普通用户永远只跟 `*sql.DB` 打交道——「换数据库」是部署时的事，不是业务逻辑的事。

### net/http：协议统一，业务异构

```
                ┌─────────────┐
                │  业务逻辑     │
                │ listUsers()  │  ← 每个接口的处理逻辑完全不同
                │ getUser()    │
                │ createUser() │
                └──────┬──────┘
                       │ 业务多样性 ↑
                       │ http.Handler（接口）
                ┌──────┴──────┐
                │ http.Server │  ← 协议解析 + 连接管理
                │              │     HTTP/1.1 和 HTTP/2 都是标准化的
                └──────┬──────┘
                       │ 协议统一 ↓
                ┌──────┴──────┐
                │  net.Conn   │  ← TCP 字节流
                └─────────────┘
```

HTTP 协议是标准化的——HTTP/1.1 的请求行、头部、body 解析逻辑对所有应用都一样，HTTP/2 的帧格式也是固定的。没有「协议方言」需要适配。但「一个请求来了该怎么处理」是彻底随业务变化的：每个路由的鉴权逻辑、参数校验、业务规则、返回格式都不一样，不存在模板化可能。

所以 net/http 的策略是：**把灵活性留在接口上面**。标准库把协议解析和连接管理自己做了（`http.Server`），把请求处理抽象成一个极薄的 `http.Handler` 接口，由用户/第三方框架实现。这个接口只有一个方法——标准库不关心你在里面做什么。

### 扩展点的位置：底部 vs 顶部

```
database/sql                         net/http
─────────────                        ────────

用户代码（普通用户）                  用户代码（所有用户）
  │                                     │
  ▼                                     ▼
*sql.DB  ← 公共 API                  http.Handler  ← 扩展点在这（顶部）
  │                                     │
  │  扩展点在这（底部）                   ▼
  ▼                                  http.Server  ← 公共实现
driver.Conn  ← 只暴露给驱动作者         │
  │                                     ▼
  ▼                                  net.Conn
pgx / mysql / sqlite
```

两种设计的取舍逻辑：

| 维度 | database/sql | net/http |
|------|-------------|----------|
| 协议面 | 异构（PG wire / MySQL / SQLite C API） | 统一（HTTP/1.1、HTTP/2 都是标准） |
| 业务面 | 统一（CRUD 流程模板化） | 异构（每个路由逻辑完全不同） |
| 扩展点位置 | 底部（驱动层），只给驱动作者 | 顶部（Handler 层），给所有用户 |
| 扩展点的受众 | 少数驱动开发者 | 所有应用开发者 |
| 设计策略 | 把复杂性藏在接口下面 | 把灵活性留在接口上面 |

这就是为什么两者看起来都是「标准库定义接口 + 第三方实现」，但具体做法截然相反：**database/sql 把统一的做成 struct、把差异的做成 interface（往下藏）；net/http 把统一的自己做完、把差异的做成 interface（往上露）**。

### 这也能解释「推 vs 拉」的根源

控制流方向的差异不是偶然的，它直接来源于复杂度模型：

- database/sql 是拉模式，因为**用户掌握「什么时候操作」的主动权**——「查什么、什么时候查」是业务逻辑决定的。数据库在那边等着，你发起请求它才干活。
- net/http 是推模式，因为**框架掌握「什么时候请求到达」的主动权**——「请求什么时候来」是外部客户端决定的。你的 handler 是被动等待调用的，数据来了才触发。

所以推/拉不是设计者的偏好选择，而是问题域本身的属性决定了控制流方向。控制流方向又进一步决定了入口类型：主动调用适合 struct（方法挂在上面），被动回调适合 interface（框架来调你）。

---

## 九、总结

```
                database/sql              net/http
─────────────────────────────────────────────────────
入口             *sql.DB (struct)         http.Handler (interface)
核心操作         用户主动调用 Query/Exec   框架回调 ServeHTTP
状态持有者       sql.DB                    http.Server
连接管理         DB.freeConn 连接池        Server 的 net.Conn 池
能力扩展         类型断言 driver.Conn      类型断言 ResponseWriter
行为叠加         N/A（不需要）             中间件装饰器模式
注册机制         全局 map                  DefaultServeMux 全局 + 路由器局部
下层可替换       driver.Connector          net.Listener
```

**一句话**：database/sql 是你手里的工具（你主动用它），net/http 是你写的回调（等着被框架调用）。两种模式覆盖了 I/O 框架设计的两大流派——它们的内核思维是对称的，但控制流方向决定了入口是 struct 还是 interface。
