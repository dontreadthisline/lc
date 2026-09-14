# database/sql 设计模式深度解析

> 基于 Go 1.26.1 标准库源码分析，结合 pgx v5.7.5 驱动实现

---

## 目录

1. [整体架构：三层模型](#1-整体架构三层模型)
2. [面向用户的接口](#2-面向用户的接口-databasesql-公开-api)
3. [驱动接口：插件的契约](#3-驱动接口-databasesqldriver)
4. [插件的注册与发现](#4-插件的注册与发现机制)
5. [连接池的内部实现](#5-连接池的内部实现)
6. [查询执行全流程](#6-查询执行全流程)
7. [pgx 如何实现 driver 接口](#7-pgx-如何实现-driver-接口)
8. [设计模式拆解](#8-设计模式拆解)
9. [与其他语言的范式对比](#9-与其他语言的范式对比)
10. [总结：这种设计的好处](#10-总结这种设计的好处)
11. [抽象如何恰到好处地应对变更](#11-抽象如何恰到好处地应对变更)

---

## 1. 整体架构：三层模型

```
+------------------------------------------------------------------+
|                        业务代码                                    |
|  db.Query(...)   db.Exec(...)   tx.Commit()   rows.Scan(...)     |
+------------------------------------------------------------------+
                              |
                    面向用户的一层抽象
                     (database/sql 包)
                              |
+------------------------------------------------------------------+
|                      database/sql 核心                            |
|                                                                   |
|  DB (连接池)   Stmt (预编译)   Tx (事务)   Rows (结果集)            |
|  Register()    连接管理       参数转换      类型系统                |
+------------------------------------------------------------------+
                              |
                       driver 接口契约
                     (database/sql/driver 包)
                              |
+------------------------------------------------------------------+
|  pgx/stdlib     go-sql-driver/mysql     lib/pq      sqlite3       |
|  (PostgreSQL)   (MySQL)                 (PG 旧)     (SQLite)      |
+------------------------------------------------------------------+
```

database/sql 的设计核心是 **依赖倒置**：
- 上层（业务代码）依赖抽象（database/sql），不依赖具体驱动
- 下层（驱动）实现抽象接口（driver.Conn/Driver/Stmt...），注入到框架中
- 框架（database/sql）提供连接池、参数绑定、类型转换等通用能力

---

## 2. 面向用户的接口（database/sql 公开 API）

### 2.1 核心类型一览

```
sql.Open(driverName, dsn) ──→ *sql.DB ──→ 连接池句柄
                                    │
                    ┌───────────────┼───────────────┐
                    │               │               │
               Query/Exec      Prepare          Begin
                    │               │               │
                    ▼               ▼               ▼
               *sql.Rows        *sql.Stmt       *sql.Tx
               *sql.Row                        (Commit/Rollback)
                    │
                    ▼
               sql.Result
               (LastInsertId, RowsAffected)
```

### 2.2 DB：连接池句柄

```go
// 源码：sql.go:507
type DB struct {
    connector   driver.Connector   // 驱动提供的连接工厂
    freeConn    []*driverConn      // 空闲连接列表（按归还时间排序）
    connRequests connRequestSet    // 等待连接的请求队列
    numOpen     int                // 当前打开连接数
    maxOpen     int                // MaxOpenConns 上限，0=无限制
    maxIdleCount int               // MaxIdleConns 上限
    maxLifetime time.Duration     // 连接最大存活时间
    maxIdleTime time.Duration     // 空闲连接最大闲置时间
    // ...
}
```

`DB` 本身不直接持有数据库连接，它管理一个连接池。每个用户可见的操作（Query/Exec/Ping）都从池中取连接，用完归还。

### 2.3 Rows：结果集迭代器

```go
rows, err := db.Query("SELECT id, name FROM users")
defer rows.Close()  // 必须 Close，否则连接不归还

for rows.Next() {
    rows.Scan(&id, &name)  // 逐行扫描
}
// 检查遍历过程中的错误
if err := rows.Err(); err != nil { ... }
```

关键设计点：
- `Next()` 通知驱动取下一行数据
- `Scan()` 将当前行的列值拷贝到用户变量
- `Close()` 释放底层连接（**即使遍历完也必须调用**）
- `Err()` 返回遍历过程中遇到的错误

### 2.4 Stmt：预编译语句

```go
stmt, err := db.Prepare("SELECT name FROM users WHERE id = $1")
defer stmt.Close()

// 多次执行，数据库只解析一次 SQL
stmt.QueryRow(1).Scan(&name)
stmt.QueryRow(2).Scan(&name)
```

### 2.5 Tx：事务

```go
tx, err := db.Begin()
defer tx.Rollback() // 安全网：未 Commit 前退出会自动回滚

tx.Exec("UPDATE accounts SET balance = balance - 100 WHERE id = 1")
tx.Exec("UPDATE accounts SET balance = balance + 100 WHERE id = 2")

tx.Commit() // 提交，之后 Rollback 是 no-op
```

---

## 3. 驱动接口（database/sql/driver）

这是插件体系的核心契约。database/sql 只依赖这些接口，不依赖任何具体实现。

### 3.1 接口层次结构

```
driver.Driver                  ← 入口：驱动的"工厂"
    │
    ├── Open(name) → Conn
    │
    └── (DriverContext)        ← 可选：两步创建连接，支持 context
         └── OpenConnector(name) → Connector
              └── Connect(ctx) → Conn

driver.Conn                    ← 一条数据库连接
    ├── Prepare(query) → Stmt
    ├── Close()
    └── Begin() → Tx

    (可选增强接口，通过类型断言检测)
    ├── Pinger          → Ping(ctx)
    ├── ExecerContext   → ExecContext(ctx, query, args)
    ├── QueryerContext  → QueryContext(ctx, query, args)
    ├── ConnPrepareContext → PrepareContext(ctx, query)
    ├── ConnBeginTx     → BeginTx(ctx, opts)
    ├── SessionResetter → ResetSession(ctx)  ← 连接复用时重置会话
    └── Validator       → IsValid()          ← 放回池前检查连接

driver.Stmt                    ← 预编译语句
    ├── NumInput() → int
    ├── Exec(args) / ExecContext(ctx, args)
    ├── Query(args) / QueryContext(ctx, args)
    └── Close()

driver.Rows                    ← 结果集迭代
    ├── Columns() → []string
    ├── Next(dest []Value) → error
    └── Close()

driver.Tx                      ← 事务
    ├── Commit()
    └── Rollback()
```

### 3.2 核心接口定义

**Driver — 驱动的入口：**
```go
// driver.go:85
type Driver interface {
    Open(name string) (Conn, error)
}
```

**Connector — 更现代的两步连接方式（Go 1.10+）：**
```go
// driver.go:122
type Connector interface {
    Connect(context.Context) (Conn, error)
    Driver() Driver
}
```

**Conn — 一条数据库连接：**
```go
// driver.go:234
type Conn interface {
    Prepare(query string) (Stmt, error)
    Close() error
    Begin() (Tx, error)
}
```

### 3.3 可选接口的渐进增强

这是 database/sql 设计中最巧妙的点。驱动 **不需要** 实现所有接口，它通过 **类型断言（type assertion）** 在运行时发现能力：

```go
// database/sql 内部逻辑（简化）：
func (db *DB) pingDC(ctx context.Context, dc *driverConn) error {
    // 检查连接是否支持 Pinger 接口
    if pinger, ok := dc.ci.(driver.Pinger); ok {
        return pinger.Ping(ctx)
    }
    // 不支持 Pinger 时，降级为首选一条空闲连接
    return nil
}
```

可选接口按功能维度叠加：

| 接口 | 作用 | 不支持时的降级行为 |
|------|------|-------------------|
| Pinger | 真正的数据库 Ping | 只检查连接池中是否有连接 |
| ExecerContext | 直接执行（跳过 Prepare） | 先 Prepare → Exec → Close Stmt |
| QueryerContext | 直接查询（跳过 Prepare） | 先 Prepare → Query → Close Stmt |
| ConnPrepareContext | 带 context 的预编译 | 用 Prepare（无 context）降级 |
| ConnBeginTx | 带隔离级别的事务 | 用 Begin（无选项）降级 |
| SessionResetter | 连接复用前重置状态 | 不重置，可能带脏状态 |
| Validator | 归还前校验连接是否有效 | 总是认为有效 |

这种设计实现了 **渐进增强（Progressive Enhancement）**：驱动开发者可以从最小实现开始，逐步添加优化路径。

### 3.4 Value 类型系统

```go
type Value any  // 驱动内部使用的值类型

// 合法类型：nil, int64, float64, bool, []byte, string, time.Time
```

database/sql 在用户类型和 driver.Value 之间做了双向转换：

```
用户代码                           database/sql                    驱动
  string       ──→ convertAssign ──→ driver.Value(string)  ──→ 网络协议
  int          ──→ convertAssign ──→ driver.Value(int64)   ──→ 网络协议
  sql.NullString ─→ Valuer.Value() ──→ nil or string         ──→ 网络协议
  time.Time     ──→ convertAssign ──→ driver.Value(time.Time) ──→ 网络协议
```

### 3.5 扩展接口一览

驱动还可以实现更细粒度的接口来暴露数据库特性：

- `RowsColumnTypeScanType` — 告知每列的 Go 类型
- `RowsColumnTypeDatabaseTypeName` — 告知数据库类型名（用于 ORM）
- `RowsColumnTypeLength` — 列长度信息
- `RowsColumnTypeNullable` — 列是否可空
- `RowsColumnTypePrecisionScale` — decimal 的精度和小数位
- `RowsNextResultSet` — 多结果集支持
- `NamedValueChecker` — 自定义参数类型检查

---

## 4. 插件的注册与发现机制

### 4.1 注册流程

```
编译时：import _ "github.com/jackc/pgx/v5/stdlib"
         ↓
       执行 stdlib 包的 init() 函数
         ↓
       调用 sql.Register("pgx", driverInstance)
         ↓
       将驱动存入全局 map[driverName]driver.Driver
```

**database/sql 侧（sql.go:56）：**
```go
var drivers = make(map[string]driver.Driver)

func Register(name string, driver driver.Driver) {
    driversMu.Lock()
    defer driversMu.Unlock()
    if driver == nil {
        panic("sql: Register driver is nil")
    }
    if _, dup := drivers[name]; dup {
        panic("sql: Register called twice for driver " + name)
    }
    drivers[name] = driver
}
```

**pgx 侧（stdlib/sql.go init）：**
```go
func init() {
    pgxDriver = &Driver{configs: make(map[string]*pgx.ConnConfig)}
    if !slices.Contains(sql.Drivers(), "pgx") {
        sql.Register("pgx", pgxDriver)
    }
    sql.Register("pgx/v5", pgxDriver)
}
```

### 4.2 发现流程

```go
// sql.Open 内部
func Open(driverName, dataSourceName string) (*DB, error) {
    // 1. 从全局 map 查找驱动
    driveri, ok := drivers[driverName]
    if !ok {
        return nil, fmt.Errorf("sql: unknown driver %q (forgotten import?)", driverName)
    }

    // 2. 如果驱动实现了 DriverContext，走两步创建（Connector → Connect）
    if driverCtx, ok := driveri.(driver.DriverContext); ok {
        connector, err := driverCtx.OpenConnector(dataSourceName)
        return OpenDB(connector), nil
    }

    // 3. 否则走旧的一步创建（Driver.Open）
    return OpenDB(dsnConnector{dsn: dataSourceName, driver: driveri}), nil
}
```

### 4.3 为什么用空白导入（`import _`）

```go
import (
    "database/sql"
    _ "github.com/jackc/pgx/v5/stdlib"  // 仅执行 init()，注册驱动
)
```

- 业务代码不需要直接引用驱动包的任何类型
- 驱动通过副作用（注册到全局 map）完成集成
- 编译时链接，运行时通过字符串名查找
- 这是 Go 语言 **编译时插件** 模式的经典应用

---

## 5. 连接池的内部实现

### 5.1 DB 连接池的数据结构

```
DB
├── freeConn []*driverConn       ← 空闲连接（切片，按归还时间排序）
├── connRequests                 ← 等待队列（当 MaxOpenConns 耗尽时）
├── numOpen int                  ← 当前打开的连接数（含使用中 + 空闲）
├── openerCh chan struct{}       ← 通知后台 goroutine 创建新连接
└── connectionOpener goroutine   ← 后台 goroutine，监听 openerCh
```

### 5.2 获取连接：`db.conn()`

```
db.conn(ctx, strategy)

  1. 尝试从 freeConn 取一个空闲连接
     ├── 找到 → 检查是否过期（maxLifetime）
     │         ├── 过期 → 关闭，返回 ErrBadConn（上层会重试）
     │         └── 未过期 → ResetSession() → 返回连接
     │
     └── 没有空闲连接
          ├── numOpen < maxOpen → 直接调用 connector.Connect(ctx) 创建新连接
          │
          └── numOpen >= maxOpen → 排队等待
               ├── 将 connRequest 加入等待队列
               ├── 阻塞等待
               │   ├── 有连接归还 → 拿到连接 → 返回
               │   └── ctx 超时   → 取消等待 → 返回 context.DeadlineExceeded
               └── ...
```

关键代码（sql.go:1316-1450）：

```go
func (db *DB) conn(ctx context.Context, strategy connReuseStrategy) (*driverConn, error) {
    db.mu.Lock()
    // ...检查是否已关闭，检查 context...

    // 首选空闲连接
    if strategy == cachedOrNewConn && len(db.freeConn) > 0 {
        conn := db.freeConn[len(db.freeConn)-1]  // 取最后（最久的空闲）
        db.freeConn = db.freeConn[:len(db.freeConn)-1]
        conn.inUse = true
        // 检查过期
        if conn.expired(db.maxLifetime) {
            conn.Close()
            return nil, driver.ErrBadConn  // 上层会重试
        }
        db.mu.Unlock()
        conn.resetSession(ctx)  // 重置会话状态
        return conn, nil
    }

    // 已到上限，排队等待
    if db.maxOpen > 0 && db.numOpen >= db.maxOpen {
        req := make(chan connRequest, 1)
        db.connRequests.Add(req)
        db.waitCount++
        db.mu.Unlock()

        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        case ret := <-req:
            return ret.conn, ret.err
        }
    }

    // 还有余量，创建新连接
    db.numOpen++
    db.mu.Unlock()
    ci, err := db.connector.Connect(ctx)
    // ...包装成 driverConn 返回...
}
```

### 5.3 归还连接：`db.putConn()`

```
db.putConn(dc, err)

  1. 如果连接出错（ErrBadConn）
     → 关闭连接，不归还池

  2. 如果有等待者（connRequests 不为空）
     → 直接交给第一个等待者

  3. 如果空闲连接数 >= maxIdleCount
     → 关闭多余的空闲连接

  4. 否则
     → 加入 freeConn 列表
     → 如果实现了 Validator，调用 IsValid() 检查
```

### 5.4 连接清理器：`connectionCleaner`

database/sql 内部有一个定时运行的 goroutine（`connectionCleaner`），负责：
- 清理超过 `maxLifetime` 的连接
- 清理空闲超过 `maxIdleTime` 的连接
- 清理超过 `maxIdleCount` 的多余空闲连接

```go
func (db *DB) connectionCleanerRunLocked(d time.Duration) (time.Duration, []*driverConn) {
    // 逐一检查每个空闲连接
    for i := 0; i < len(db.freeConn); i++ {
        c := db.freeConn[i]
        if c.expired(db.maxLifetime) {
            // 超过了连接最大存活时间
            // ...
        } else if c.expired(db.maxIdleTime) {
            // 超过了空闲最大时间
            // ...
        }
    }
    // 如果空闲连接数超过 maxIdleCount，关闭最旧的
    for len(db.freeConn) > db.maxIdleCount {
        // ...
    }
}
```

---

## 6. 查询执行全流程

以一个简单的 `db.Query("SELECT * FROM users WHERE age > $1", 18)` 为例：

```
db.QueryContext(ctx, query, args...)
  │
  ├── 1. 获取连接
  │     db.conn(ctx, cachedOrNewConn)
  │       ├── 从 freeConn 取，或创建新连接
  │       └── ResetSession() 清理上一轮会话状态
  │
  ├── 2. 参数转换
  │     args 从 any → driver.NamedValue
  │     (checkNamedValue 检查类型合法性)
  │
  ├── 3. 执行查询
  │     ├── 如果 conn 实现了 QueryerContext → ci.QueryContext(ctx, query, nvargs)
  │     │   （pgx 走这条快速路径）
  │     │
  │     └── 否则降级：Prepare → Query → 扫完数据 → Close Stmt
  │         （兼容旧驱动）
  │
  ├── 4. 包装结果
  │     返回 *sql.Rows（内部持有 driver.Rows + driverConn）
  │     此时结果集还未被消费
  │
  └── 5. 连接所有权转移
       连接跟随 Rows 生命周期
       rows.Close() 时才归还连接
```

**Rows.Close() 的归还逻辑：**
```go
func (rs *Rows) Close() error {
    // 1. 确保所有行已被消费（驱动可能需要发送 cancel 消息）
    // 2. 关闭底层 driver.Rows
    rs.rowsi.Close()
    // 3. 归还 driverConn 到连接池
    rs.releaseConn(err)
    return rs.lasterr
}
```

**关键设计选择：连接跟随结果集**

Rows 持有底层连接的所有权。这是出于正确性的考虑：在 PostgreSQL 中，查询结果可能还在服务端未取完，如果提前归还连接给池，下一个查询会收到上一个查询的残留数据。因此你 **必须** `rows.Close()`。

---

## 7. pgx 如何实现 driver 接口

### 7.1 架构映射

```
database/sql 接口              pgx/stdlib 实现
─────────────────────────────────────────────────
driver.Driver              →  stdlib.Driver
driver.Connector           →  stdlib.connector
driver.Conn                →  stdlib.Conn（包装 *pgx.Conn）
driver.Stmt                →  stdlib.Stmt（包装 *pgconn.StatementDescription）
driver.Rows                →  stdlib.Rows（包装 pgx.Rows）
driver.Tx                  →  stdlib.wrapTx（包装 pgx.Tx）
```

### 7.2 pgx 提供的多条接入路径

pgx 的 stdlib 支持比标准更灵活的接入方式：

```go
// 方式1：DSN 字符串（标准方式）
db, _ := sql.Open("pgx", "postgres://user@localhost/db")

// 方式2：从 pgxpool.Pool 创建（共享连接池）
pool, _ := pgxpool.New(ctx, dsn)
db := stdlib.OpenDBFromPool(pool)

// 方式3：通过 ConnConfig 预配置
config, _ := pgx.ParseConfig(dsn)
config.Tracer = &myTracer{}
connStr := stdlib.RegisterConnConfig(config)
db, _ := sql.Open("pgx", connStr)

// 方式4：获取原生 *pgx.Conn 做高级操作
conn, _ := db.Conn(ctx)
conn.Raw(func(driverConn any) error {
    pgxConn := driverConn.(*stdlib.Conn).Conn()
    pgxConn.CopyFrom(...)   // 使用 database/sql 不提供的 PG 特有功能
    return nil
})
```

### 7.3 pgx 实现了哪些可选接口

从源码中可以看到 pgx 实现的可选接口：

| 接口 | 是否实现 | 说明 |
|------|---------|------|
| DriverContext | 是 | OpenConnector 支持 context |
| Pinger | 是 | 真正的 PostgreSQL ping |
| ExecerContext | 是 | 直接执行，跳过 Prepare |
| QueryerContext | 是 | 直接查询，跳过 Prepare |
| ConnPrepareContext | 是 | 带 context 的预编译 |
| ConnBeginTx | 是 | 支持隔离级别设置 |
| SessionResetter | 是 | ResetSession 做 ping 检查 |
| NamedValueChecker | 是 | CheckNamedValue 接受所有类型 |
| RowsColumnTypeScanType | 是 | 返回列的 Go 反射类型 |
| RowsColumnTypeDatabaseTypeName | 是 | 返回 PG 类型名 |
| RowsColumnTypeLength | 是 | varchar 等变长类型长度 |
| RowsColumnTypePrecisionScale | 是 | numeric 精度信息 |

### 7.4 连接引用的核心技巧：prepared statement 引用计数

pgx 有一个巧妙的细节：pgx 的 Prepare 会为相同的 SQL 文本生成 **确定性的** 语句名。这意味着多次 `db.Prepare("SELECT ...")` 会得到同一个底层 prepared statement。

但 database/sql 的 `Stmt.Close()` 语义要求关闭底层语句。如果两个 `*sql.Stmt` 共享同一个底层 `*pgconn.StatementDescription`，一个 Close 会破坏另一个。

pgx 的解决方案是 **引用计数**：

```go
// stdlib/sql.go Conn struct
type Conn struct {
    psRefCounts map[*pgconn.StatementDescription]int  // 引用计数
}

// Stmt.Close 中的逻辑：
func (s *Stmt) Close() error {
    refCount := s.conn.psRefCounts[s.sd]
    if refCount == 1 {
        delete(s.conn.psRefCounts, s.sd)  // 最后一个引用，真正关闭
    } else {
        s.conn.psRefCounts[s.sd]--         // 减小引用计数
        return nil
    }
    return s.conn.conn.Deallocate(ctx, s.sd.SQL)
}
```

### 7.5 pgx 如何预处理首行数据

再看一个细节优化。pgx 在 `QueryContext` 中会在返回结果前 **预取第一行**：

```go
func (c *Conn) QueryContext(ctx context.Context, query string, argsV []driver.NamedValue) (driver.Rows, error) {
    rows, err := c.conn.Query(ctx, query, args...)
    // 预取第一行！
    more := rows.Next()
    return &Rows{
        conn: c, rows: rows,
        skipNext: true,        // 标记第一行已取
        skipNextMore: more,    // 是否有第一行
    }, nil
}
```

为什么这么做？因为 database/sql 框架在拿到 `driver.Rows` 后会立即调用 `Rows.Columns()`。在 pgx 的底层协议中，列信息在收到第一行响应时才能确定。所以 pgx 必须先取一行才能返回正确的列信息——而 `skipNext` 标记确保 `Next()` 不会跳过这行数据。

---

## 8. 设计模式拆解

database/sql 使用了多种经典设计模式：

### 8.1 抽象工厂模式（Abstract Factory）

```
sql.Open("pgx", dsn)          ← 工厂方法
     ↓
driver.Driver.OpenConnector   ← 抽象工厂
     ↓
driver.Connector.Connect()    ← 具体产品创建
```

`sql.Open` 不需要知道具体产品的类型，只需要驱动名。不同的驱动名产生不同的连接器。

### 8.2 策略模式（Strategy）

不同的数据库驱动是可替换的策略：

```
业务代码
    ↓          使用 database/sql 统一接口
    ├── pgx        （PostgreSQL 策略）
    ├── go-mysql   （MySQL 策略）
    └── go-sqlite3 （SQLite 策略）
```

切换数据库只需改 import 和连接字符串，业务代码零修改。

### 8.3 代理模式（Proxy）

`sql.DB`/`sql.Stmt`/`sql.Tx`/`sql.Rows` 都是代理：

```
用户代码 → sql.DB.Query() → 获取连接 → 参数转换 → ci.QueryContext()
                                          ↑
                                   代理层：连接池管理、错误重试、类型转换
```

代理层在不改变驱动接口的前提下增加了：
- 连接池管理
- 透明重试（遇到 `ErrBadConn` 自动重试最多一次）
- 参数类型校验与转换
- goroutine 安全

### 8.4 模板方法模式（Template Method）

`db.retry()` 是一个模板方法：

```go
func (db *DB) retry(fn func(strategy connReuseStrategy) error) error {
    for i := 0; ; i++ {
        err := fn(cachedOrNewConn)
        if err == nil || !errors.Is(err, driver.ErrBadConn) {
            return err  // 成功或非连接错误，直接返回
        }
        // ErrBadConn，重试一次，强制新连接
        err = fn(alwaysNewConn)
        if err == nil || !errors.Is(err, driver.ErrBadConn) {
            return err
        }
        return err  // 两次都失败，放弃
    }
}
```

所有查询方法（`query`/`exec`/`ping` 等）都通过 `retry` 执行，获得统一的错误重试策略。

### 8.5 适配器模式（Adapter）

pgx 的 `stdlib.Conn` 就是适配器：

```
pgx.Conn（原生 PG 连接）  ──适配──→  driver.Conn（database/sql 连接接口）
pgx.Tx   （原生 PG 事务）  ──适配──→  driver.Tx  （database/sql 事务接口）
pgx.Rows （原生 PG 结果集） ──适配──→  driver.Rows（database/sql 结果集接口）
```

适配器包装了 pgx 的丰富功能，只暴露 database/sql 需要的子集。用户仍然可以通过 `conn.Raw()` 访问原始对象。

### 8.6 注册表模式（Registry）

```go
var drivers = make(map[string]driver.Driver)  // 全局注册表
```

驱动通过 `init()` + `sql.Register()` 自注册到全局表，框架通过名称查找。这是 Go 中最常见的插件模式。

---

## 9. 与其他语言的范式对比

### 9.1 Rust

Rust 生态中没有一个类似 database/sql 那样统一的标准抽象层。不同数据库有各自的异步驱动：

| 库 | 范式 | 特点 |
|----|------|------|
| sqlx | 异步 trait 抽象 | 最接近 database/sql 的设计 |
| diesel | ORM + 编译时类型检查 | 编译期验证 SQL |
| tokio-postgres | 原生异步 | 直接映射 PG 协议 |
| sea-orm | 异步 ORM | 类似 ActiveRecord |

**sqlx 的设计对比：**

```rust
// Rust: trait 定义（编译时多态）
#[async_trait]
pub trait Database: Send + Debug {
    type Connection: Connection<Database = Self>;
    type Row: Row<Database = Self>;
    // ...
}

// Rust: 泛型约束
async fn query<DB: Database>(pool: &Pool<DB>) -> Result<Vec<DB::Row>> {
    sqlx::query("SELECT * FROM users").fetch_all(pool).await
}
```

**Go 与 Rust 的核心差异：**

| 维度 | Go database/sql | Rust sqlx |
|------|----------------|-----------|
| 多态方式 | interface（运行时） | trait + generic（编译时） |
| 异步模型 | 同步 + context | async/await |
| 零成本抽象 | 有接口调度的虚拟调用开销 | 单态化（monomorphization），无虚调用 |
| 错误处理 | 返回 error | Result<T, E> |
| 编译时检查 | 无 | SQL 编译检查（sqlx::query!） |
| 连接池 | 内置 | 内置 |

**为什么 Go 选择运行时接口而非泛型？**

Go 的 interface 允许在运行时动态发现驱动能力（类型断言 `if p, ok := conn.(driver.Pinger)`）。Rust 的 trait 是编译时确定的，理论上也能做，但需要更复杂的 trait 继承层次。Go 的方式更灵活——驱动只需实现它支持的可选接口即可。

### 9.2 Python

Python 的 **DB-API 2.0（PEP 249）** 是 database/sql 精神上的前辈：

```python
# Python DB-API 2.0
import psycopg2

conn = psycopg2.connect("postgresql://user@localhost/db")
cur = conn.cursor()
cur.execute("SELECT * FROM users WHERE age > %s", (18,))
rows = cur.fetchall()
cur.close()
conn.close()
```

**Go 与 Python 的核心差异：**

| 维度 | Go database/sql | Python DB-API 2.0 |
|------|----------------|-------------------|
| 接口定义 | 编译时 interface | 约定（duck typing） |
| 遵守程度 | 强制（编译器检查） | 靠自觉（运行时可能 AttributeError） |
| 驱动注册 | import _ + init() | 需要显式 import 驱动模块 |
| 连接池 | 内置 | PEP 249 不定义，各驱动自行实现 |
| 参数风格 | 驱动自定义（$1 / ? / :name） | 五种 paramstyle（qmark/numeric/named/format/pyformat） |
| 异步 | context.Context | asyncio + asyncpg/aiomysql |
| 事务 | Begin/Commit/Rollback | commit()/rollback() |

**Python 的鸭子类型问题：**

Python DB-API 没有编译器强制接口实现。如果一个驱动漏了某个方法，只有运行时调用时才报 `AttributeError`。Go 的 `driver.Conn` interface 在编译时就能确保实现完整性。

**Python 异步方案的特殊设计：**

```python
# asyncpg（Python 最快的 PG 驱动，类似 Go 的 pgx 原生接口）
import asyncpg

async def main():
    pool = await asyncpg.create_pool("postgresql://user@localhost/db")
    async with pool.acquire() as conn:
        rows = await conn.fetch("SELECT * FROM users WHERE age > $1", 18)
```

Python 的异步驱动（asyncpg/aiomysql）完全绕过了 DB-API 2.0（它是同步的），各自定义了专有接口。这导致生态分裂——SQLAlchemy 需要为每个异步驱动写适配层。

### 9.3 对比总结

```
            接口定义方式      类型安全      扩展方式           异步支持
Go          编译时 interface   强           类型断言             context
Python      PEP 约定          弱（运行时）  鸭子类型            asyncio + 驱动专有
Rust        编译时 trait       极强         泛型 + trait bound  async/await
Java        JDBC interface     中           实现接口            Thread（JDBC 同步）
```

Go 的方案在 **简单性** 和 **安全性** 之间取得了平衡：
- 比 Python 的类型安全强（编译时检查接口实现）
- 比 Rust 的门槛低（不用理解 trait 和 lifetime）
- 比 Java 的 JDBC 更轻量（JDBC 接口层次更深，抽象更多）

---

## 10. 总结：这种设计的好处

### 10.1 对业务开发者

1. **一个 API 对接所有数据库** — 学一次 `sql.DB`，PostgreSQL/MySQL/SQLite 通用
2. **连接池开箱即用** — 不需要自己管理连接生命周期
3. **类型转换自动化** — `Scan(&name)` 自动处理 Go 类型到数据库类型的映射
4. **透明重试** — 遇到坏连接自动换新连接重试，业务代码无感知
5. **goroutine 安全** — `DB` 可以全局共享，并发安全

### 10.2 对驱动开发者

1. **接口最小化** — `driver.Conn` 只需 3 个方法即可起步
2. **可选接口渐进** — 逐步实现 Pinger/ExecerContext 等获取性能提升
3. **框架处理脏活** — 连接池、参数转换、重试逻辑由 database/sql 负责
4. **单一职责** — 驱动只需关心"如何与数据库通信"，不用关心线程安全、连接复用

### 10.3 设计哲学上的启示

database/sql 的设计体现了 Go 语言的核心哲学：

- **组合优于继承** — 可选接口通过类型断言组合，不需要复杂的接口继承树
- **约定优于配置** — 驱动注册用 `init()` + 空白导入，零配置
- **显性优于隐性** — `rows.Close()` 必须显式调用，连接所有权清晰
- **简单优于复杂** — 整个 core API 不到 10 个类型，driver 接口也只有十几个方法

这个设计经受了 13 年（Go 1.0 至今）的考验，支撑了 Go 生态中几乎所有数据库访问。它证明了一个好的抽象不需要很复杂——关键是把边界划对，让每一层只做自己该做的事。

---

## 11. 抽象如何恰到好处地应对变更

一个抽象的优劣，最终不是看它今天多优雅，而是看它未来被什么东西打脸。database/sql 在 14 年间核心接口几乎不变，却支撑了 PostgreSQL、MySQL、SQLite、Oracle、Snowflake、BigQuery 等数十种数据库驱动。以下是让这种抽象**恰到好处**的五个关键决策。

### 11.1 最小接口 + 可选增强 = 版本演化的安全网

`driver.Conn` 只需要 3 个方法即可起步：

```go
type Conn interface {
    Prepare(query string) (Stmt, error)
    Close() error
    Begin() (Tx, error)
}
```

后来需要加 `Ping`、带 Context 的查询、跳过 Prepare 的快速路径——这些怎么加？**不在接口里加方法，用类型断言**：

```go
// 运行时发现能力，不支持就降级
if pinger, ok := conn.(driver.Pinger); ok {
    return pinger.Ping(ctx)
}
// 降级行为：只检查连接池里有没有连接（不真正 ping）
```

这是整个设计中最关键的一手。假设用传统 OOP 思路，每加一个能力就要在接口里加方法——所有已有驱动全部编译失败。Go 的隐式接口 + 类型断言，让 database/sql 可以在**不修改接口定义、不破坏已有驱动**的前提下，持续添加新能力。

可选接口按功能维度叠加，不支持时有明确的降级行为：

| 接口 | 作用 | 不支持时的降级 |
|------|------|---------------|
| Pinger | 真正的数据库 Ping | 只检查连接池中是否有连接 |
| ExecerContext | 直接执行（跳过 Prepare） | Prepare → Exec → Close Stmt |
| QueryerContext | 直接查询（跳过 Prepare） | Prepare → Query → Close Stmt |
| ConnPrepareContext | 带 context 的预编译 | 用 Prepare（无 context）降级 |
| ConnBeginTx | 带隔离级别的事务 | 用 Begin（无选项）降级 |
| SessionResetter | 连接复用前重置状态 | 不重置，可能带脏状态 |
| Validator | 归还前校验连接有效性 | 总是认为有效 |

这就回答了"怎么应对外界变更"：**新增能力永远是可选的。** 驱动开发者从最小实现起步，逐步添加优化路径，每一步都是纯增量，零 breaking change。

### 11.2 双层入口设计：Driver（旧）→ Connector（新），平滑迁移

Go 1.0 的驱动入口是同步的一步式接口：

```go
type Driver interface {
    Open(name string) (Conn, error)
}
```

后来发现两个问题：
1. 没有 `context.Context`，超时控制不了
2. 像 pgx 这类驱动需要 `ConnConfig` 对象，不想每次都 parse DSN

于是 Go 1.10 新增了更现代的两步式接口：

```go
type DriverContext interface {
    OpenConnector(name string) (Connector, error)
}

type Connector interface {
    Connect(context.Context) (Conn, error)
    Driver() Driver
}
```

但 `sql.Open` 的签名没变。内部做的是能力检测：

```go
func Open(driverName, dataSourceName string) (*DB, error) {
    driveri := drivers[driverName]
    // 新驱动走两步创建
    if driverCtx, ok := driveri.(driver.DriverContext); ok {
        connector, _ := driverCtx.OpenConnector(dataSourceName)
        return OpenDB(connector), nil
    }
    // 旧驱动适配后走旧逻辑
    return OpenDB(dsnConnector{dsn: dataSourceName, driver: driveri}), nil
}
```

效果：用户代码一行不动。老驱动不用改。新驱动获得新能力。**接口没变，实现升了级。** 这是"开放封闭原则"在 Go 语言中的典范实现——对扩展开放（新增 Connector 路径），对修改封闭（旧 Driver 路径继续工作）。

### 11.3 `driver.Value` —— 极窄的类型交集作为通用货币

database/sql 定义的内部值类型只有 6 种：

```go
type Value any  // 实际合法类型: nil, int64, float64, bool, []byte, string, time.Time
```

为什么这么"抠"？因为这是一个**所有数据库的最大公约数**。MySQL 没有 PostgreSQL 的 `inet` 类型，SQLite 没有原生时间戳。如果在抽象层引入数据库特有类型，只有两种结果：

- 驱动做映射（增加复杂度，且映射可能不精确）
- 用户代码和特定数据库绑定（失去可移植性）

取舍策略：**标准路径用窄类型保证可移植性；特有类型通过 escape hatch 暴露。**

```go
// 标准路径：可移植，跨数据库
db.QueryRow("SELECT name, created_at FROM users WHERE id=$1", 1).Scan(&name, &createdAt)

// escape hatch：拿到原生连接，做 database/sql 不支持的事
conn, _ := db.Conn(ctx)
conn.Raw(func(driverConn any) error {
    pgxConn := driverConn.(*stdlib.Conn).Conn()
    pgxConn.CopyFrom(...)  // PostgreSQL 的 COPY 协议
    return nil
})
```

这个设计的启示：**抽象的边界不在"覆盖所有功能"，而在"覆盖所有数据库的共同功能"。** 把差异化能力留给 escape hatch，不在抽象层解决——database/sql 不拦你，但也不替你做。

### 11.4 连接池作为透明的中间层

用户视角的一行调用：

```
db.Query("SELECT ...")
```

背后实际发生的事情：

```
db.Query → 从池取连接 → 参数转换 → 可选接口检测 → 驱动查询
              ↑                                            ↓
         连接过期自动重试                           rows.Close() 时归还
      (ErrBadConn → 重取新连接)
```

这层代理是透明的——用户不关心连接从哪来、有没有过期、要不要重试。更重要的是：**如果驱动实现变了（比如 pgx v4 → v5 的连接创建方式完全不同），database/sql 的连接池逻辑完全不用动。** 它只管 `driver.Connector.Connect(ctx)` 拿到一个 `driver.Conn`，后面怎么来的不关心。

这就是依赖倒置的威力：上层依赖抽象（`driver.Connector` 接口），下层变化不影响上层。

### 11.5 把"谁负责什么"画了一条极清晰的线

| 层 | 职责 | 不负责 |
|----|------|--------|
| database/sql | 连接池、重试、参数转换、并发安全 | 任何 SQL 方言、数据库协议 |
| driver 接口 | 定义契约：Prepare/Query/Exec/Begin | 任何实现细节 |
| 驱动实现 | 数据库协议、SQL 方言、类型映射 | 连接管理、重试策略 |

这条线的效果：
- **驱动开发者的心智负担极低**——只需关心"怎么和数据库通信"，连接池、并发安全、重试全是 framework 的事
- **框架不会过度侵入**——它不假设任何 SQL 方言特性，不在参数绑定上强加规则
- **替换成本极低**——切换数据库只需换 import 和 DSN，业务代码一行不动（你 `02-query/main.go` 里改一行 `sql.Open("mysql", ...)` 就能切到 MySQL）

作为对比，Java 的 JDBC 也有类似分层，但接口树深得多（`DataSource` → `Connection` → `Statement` → `PreparedStatement` → `CallableStatement` → `ResultSet`），每加一个数据库特性都得在这些层次里找到合适的位置。database/sql 用一层很薄的必需接口 + 可选类型断言，避免了接口膨胀。

### 11.6 总结

database/sql 的抽象"恰到好处"，本质上是**边界划得对**：

```
                 database/sql（框架层）
                连接池 / 重试 / 参数转换 / 并发安全
                       ↑
              driver 接口（抽象契约层）
         driver.Conn / Stmt / Rows / Tx
              + 可选增强接口（Pinger / ExecerContext / ...）
                       ↑
              pgx / go-mysql / sqlite3（驱动实现层）
             数据库协议 / SQL 方言 / 类型映射
```

在框架和驱动之间，只传递三种东西——连接、语句、结果集。所有数据库都有这些，所以接口稳定。数据库之间的差异（特有类型、协议、优化），要么通过可选接口渐进增强，要么通过 escape hatch 留给用户自己处理。

14 年，核心接口零破坏性变更，同时能力持续增长。一个好的抽象不需要很复杂——**关键是把边界划对，让每一层只做自己该做的事。**

---

## 参考资料

- [Go database/sql 源码](https://cs.opensource.google/go/go/+/master:src/database/sql/)
- [Go database/sql/driver 源码](https://cs.opensource.google/go/go/+/master:src/database/sql/driver/)
- [pgx stdlib 源码](https://github.com/jackc/pgx/tree/master/stdlib)
- [PEP 249 — Python Database API Specification v2.0](https://peps.python.org/pep-0249/)
- [Rust sqlx](https://github.com/launchbadge/sqlx)
