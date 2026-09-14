// 01-connect: database/sql 连接生命周期
//
// 演示 sql.Open 的工作原理：
//   - sql.Open 不建立实际连接，只验证 driver 名和 DSN 格式
//   - db.Ping 才会真正连接数据库
//   - db.Close 释放连接池中所有连接
//
// 运行: go run ./sql/01-connect/

package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // 注册 pgx 驱动
)

func main() {
	dsn := dsnFromEnv()

	// sql.Open 只是创建一个 DB 句柄，不会建立任何连接。
	// 内部流程：查找已注册的 "pgx" 驱动 → 调用 driver.Open(dsn) → 返回 driver.Connector
	// 这个设计允许延迟连接，避免启动时就占用数据库资源。
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		fmt.Printf("sql.Open 失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	fmt.Println("[OK] sql.Open 成功（此时还未连接数据库）")

	// Ping 建立第一条连接，验证数据库可达。
	// 内部：从连接池取一条空闲连接（没有就新建）→ 调用 driver.Conn.Ping()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		fmt.Printf("Ping 失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("[OK] Ping 成功（连接已建立）")

	// 查看连接池状态
	stats := db.Stats()
	fmt.Printf("\n连接池状态:\n")
	fmt.Printf("  空闲连接: %d\n", stats.Idle)
	fmt.Printf("  使用中:   %d\n", stats.InUse)
	fmt.Printf("  最大打开: %d (0=无限制)\n", stats.MaxOpenConnections)
	fmt.Printf("  等待计数: %d\n", stats.WaitCount)

	fmt.Println("\n程序退出时 defer db.Close() 会释放所有连接。")
}

func dsnFromEnv() string {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn
	}
	return "postgres://didi@localhost:5432/postgres?sslmode=disable"
}
