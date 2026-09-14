// http/chi-demo: chi 路由框架 + net/http 标准库协作演示
//
// 演示 net/http 的设计模式：
//   - http.Handler 接口作为统一契约
//   - chi 的中间件链 = database/sql 的 driver.Conn 可选接口链
//   - chi.Router 实现 http.Handler，可以嵌套挂载
//   - 对比 database/sql 的设计异同
//
// 运行: go run ./http/

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// ================================================================
	// 对比点 1: 创建入口
	// database/sql: db, _ := sql.Open("pgx", dsn)    → 返回 *sql.DB（具体结构体）
	// net/http:     r := chi.NewRouter()              → 返回 chi.Router（接口 + 具体实现）
	//              r 本身实现了 http.Handler 接口
	// ================================================================

	r := chi.NewRouter()

	// ================================================================
	// 对比点 2: 中间件 = database/sql 的可选接口
	// database/sql: driver.Conn 可选实现 Pinger/SessionResetter/Validator
	//              → 通过类型断言在运行时发现能力
	// net/http:    中间件包装 http.Handler → http.Handler
	//              → 装饰器模式，编译时确定调用链
	// ================================================================

	r.Use(middleware.RequestID)            // 每个请求分配 ID（类似 connection 的标识）
	r.Use(middleware.RealIP)               // 解析真实 IP
	r.Use(middleware.Logger)               // 请求日志
	r.Use(middleware.Recoverer)            // panic 恢复
	r.Use(middleware.Timeout(30 * time.Second)) // 超时控制（类似 context.Context）

	// ================================================================
	// 对比点 3: 路由注册 = database/sql 的驱动注册
	// database/sql: sql.Register("pgx", driver)  → 全局 map 注册
	// net/http:     r.Get("/path", handler)      → 路由器内部 trie 树注册
	// ================================================================

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("chi demo 运行中\n"))
	})

	// RESTful 路由组
	r.Route("/users", func(r chi.Router) {
		// GET /users
		r.Get("/", listUsers)
		// GET /users/{id}
		r.Get("/{id}", getUser)
		// POST /users
		r.Post("/", createUser)
	})

	// ================================================================
	// 对比点 4: 挂载点 = database/sql 的连接池
	// database/sql: *sql.DB 管理 driver.Conn 的生命周期（获取/归还/清理）
	// net/http:     *http.Server 管理 net.Conn 的 accept/serve/close
	//
	// chi.Router 只负责路由，Server 负责 I/O。
	// 这和 database/sql 的分工不同——sql.DB 既管路由（驱动选择）又管 I/O（连接池）。
	// http 把路由放到了 Handler 层，I/O 留在了 Server 层。
	// ================================================================

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,               // r 实现了 http.Handler
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 优雅关闭
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit
		fmt.Println("\n正在关闭服务器...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()

	fmt.Println("chi demo 启动在 http://localhost:8080")
	fmt.Println("试试: curl http://localhost:8080/users")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
	fmt.Println("服务器已关闭")
}

// 请求处理器
func listUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]`))
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"id":%s,"name":"用户 %s"}`, id, id)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"created"}`))
}
