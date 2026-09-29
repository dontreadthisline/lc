// 07-pgx-native: pgx 原生接口 vs database/sql
//
// 演示 pgx 的两种使用方式及其差异：
//   1. database/sql 接口 — 通用抽象，所有数据库统一 API
//   2. pgx 原生接口 — PostgreSQL 专属，支持 COPY、LISTEN/NOTIFY、数组类型等
//
// 关键区别：
//   - database/sql 抽象了所有数据库，pgx 原生接口暴露 PostgreSQL 特性
//   - pgx 原生类型更丰富：pgtype 支持 PG 特有类型（数组、JSONB、UUID 等）
//   - pgxpool 是 pgx 的连接池，性能比 database/sql 的连接池略好
//
// 运行: go run ./sql/07-pgx-native/

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://didi@localhost:5432/postgres?sslmode=disable"
	}

	ctx := context.Background()

	// ==========================================
	// 方式一：pgxpool（推荐用于生产环境）
	// ==========================================
	pool, err := pgxpool.New(ctx, dsn)
	check(err)
	defer pool.Close()

	// 建表
	mustExecPool(ctx, pool, `
		DROP TABLE IF EXISTS pgx_demo;
		CREATE TABLE pgx_demo (
			id    SERIAL PRIMARY KEY,
			name  TEXT NOT NULL,
			tags  TEXT[],
			data  JSONB
		);
	`)

	// ---- 插入数据：展示 pgx 对 PG 特有类型的支持 ----
	// database/sql 的 Exec 只接受 ...any，对数组/JSONB 需要自己序列化。
	// pgx 原生接口直接支持 Go 切片 → PG 数组的转换。
	tags := []string{"go", "postgres", "pgx"}
	data := map[string]any{"version": "v5", "features": []string{"COPY", "LISTEN"}}

	var id int
	err = pool.QueryRow(ctx,
		"INSERT INTO pgx_demo (name, tags, data) VALUES ($1, $2, $3) RETURNING id",
		"pgx 原生演示",
		tags,   // pgx 直接将 []string 转为 text[]
		data,   // pgx 直接将 map 转为 jsonb
	).Scan(&id)
	check(err)
	fmt.Printf("[pgx 原生] 插入成功 id=%d\n\n", id)

	// ---- 批量 COPY（database/sql 没有的 PG 专属功能） ----
	fmt.Println("[COPY] 批量导入 3 行数据（database/sql 不支持）:")
	copyCount, err := pool.CopyFrom(
		ctx,
		pgx.Identifier{"pgx_demo"},
		[]string{"name", "tags", "data"},
		pgx.CopyFromSlice(3, func(i int) ([]any, error) {
			return []any{
				fmt.Sprintf("item-%d", i),
				[]string{fmt.Sprintf("tag-%d", i)},
				map[string]any{"index": i},
			}, nil
		}),
	)
	check(err)
	fmt.Printf("  复制了 %d 行\n\n", copyCount)

	// ---- 查询 ----
	rows, err := pool.Query(ctx, "SELECT id, name, tags, data FROM pgx_demo ORDER BY id")
	check(err)
	defer rows.Close()

	fmt.Println("[查询] 所有行:")
	for rows.Next() {
		var rid int
		var rname string
		var rtags []string // pgx 自动将 PG text[] 扫入 []string
		var rdata map[string]any
		check(rows.Scan(&rid, &rname, &rtags, &rdata))
		fmt.Printf("  id=%-2d name=%-15s tags=%v data=%v\n", rid, rname, rtags, rdata)
	}

	// ==========================================
	// 方式二：pgconn（底层连接，演示连接级别操作）
	// ==========================================
	conn, err := pgx.Connect(ctx, dsn)
	check(err)
	defer conn.Close(ctx)

	fmt.Println("\n[pgconn 直接连接] 演示 NOTIFY/LISTEN:")
	// pgx 原生支持 PostgreSQL 的异步通知机制
	// database/sql 没有对应的抽象

	// 先发一个 NOTIFY
	_, err = conn.Exec(ctx, "NOTIFY demo_channel, 'hello from pgx'")
	check(err)
	fmt.Println("  已发送 NOTIFY")

	// pgconn 层可以获取命令状态标签
	var pgConn *pgconn.PgConn = conn.PgConn()
	fmt.Printf("  连接状态: %t\n", pgConn.IsClosed())

	// ---- 清理 ----
	mustExecPool(ctx, pool, "DROP TABLE IF EXISTS pgx_demo;")
	fmt.Println("\n[清理] 完成")

	// ==========================================
	// 总结
	// ==========================================
	fmt.Println("\n═══════════════════════════════════════")
	fmt.Println("database/sql vs pgx 原生 选择指南:")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println()
	fmt.Println("用 database/sql 的场景:")
	fmt.Println("  - 需要支持多种数据库（MySQL + PG 切换）")
	fmt.Println("  - 只使用标准 SQL，不需要 PG 特有功能")
	fmt.Println("  - 使用 ORM（GORM/sqlx 等）配合")
	fmt.Println()
	fmt.Println("用 pgx 原生接口的场景:")
	fmt.Println("  - 只使用 PostgreSQL")
	fmt.Println("  - 需要 PG 特有类型（数组、JSONB、UUID、inet 等）")
	fmt.Println("  - 需要 COPY 协议批量导入")
	fmt.Println("  - 需要 LISTEN/NOTIFY 异步通知")
	fmt.Println("  - 追求极致性能（pgx 有更多零分配优化）")
}

func mustExecPool(ctx context.Context, pool *pgxpool.Pool, sql string) {
	_, err := pool.Exec(ctx, sql)
	check(err)
}

func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}
