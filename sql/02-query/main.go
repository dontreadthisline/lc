// 02-query: Query / QueryRow / Exec 的基本用法
//
// 演示 database/sql 三种核心操作：
//   - db.Exec — 执行不返回行的 SQL（INSERT/UPDATE/DELETE/DDL）
//   - db.Query — 返回多行结果，需遍历 rows
//   - db.QueryRow — 返回最多一行，不需要手动关闭
//
// 运行: go run ./sql/02-query/

package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	db := mustOpen()
	defer db.Close()

	ctx := context.Background()

	// ---- 建表 ----
	mustExec(ctx, db, `
		DROP TABLE IF EXISTS users;
		CREATE TABLE users (
			id    SERIAL PRIMARY KEY,
			name  TEXT NOT NULL,
			age   INT  NOT NULL
		);
	`)

	// ---- Exec: 插入数据 ----
	// Exec 返回 sql.Result，通过它获取 LastInsertId 和 RowsAffected。
	result, err := db.ExecContext(ctx,
		"INSERT INTO users (name, age) VALUES ($1, $2), ($3, $4), ($5, $6)",
		"张三", 28,
		"李四", 35,
		"王五", 22,
	)
	check(err)
	n, _ := result.RowsAffected()
	fmt.Printf("[Exec] 插入了 %d 行\n\n", n)

	// ---- Query: 多行查询 ----
	rows, err := db.QueryContext(ctx, "SELECT id, name, age FROM users ORDER BY id")
	check(err)
	defer rows.Close() // 重要：遍历完也要 Close，否则连接不会归还池

	fmt.Println("[Query] 遍历结果集:")
	for rows.Next() {
		var id, age int
		var name string
		check(rows.Scan(&id, &name, &age))
		fmt.Printf("  id=%d  name=%s  age=%d\n", id, name, age)
	}
	// 检查遍历过程中是否有错误
	check(rows.Err())
	fmt.Println()

	// ---- QueryRow: 单行查询 ----
	// QueryRow 不返回 error，Scan 时才返回 sql.ErrNoRows。
	var name string
	var age int
	err = db.QueryRowContext(ctx,
		"SELECT name, age FROM users WHERE id = $1", 2,
	).Scan(&name, &age)

	switch {
	case err == sql.ErrNoRows:
		fmt.Println("[QueryRow] 未找到该用户")
	case err != nil:
		check(err)
	default:
		fmt.Printf("[QueryRow] id=2 → name=%s age=%d\n\n", name, age)
	}

	// ---- 清理 ----
	mustExec(ctx, db, "DROP TABLE IF EXISTS users;")
	fmt.Println("[清理] 已删除 users 表")
}

func mustOpen() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://didi@localhost:5432/postgres?sslmode=disable"
	}
	db, err := sql.Open("pgx", dsn)
	check(err)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	check(db.PingContext(ctx))
	return db
}

func mustExec(ctx context.Context, db *sql.DB, query string, args ...any) {
	_, err := db.ExecContext(ctx, query, args...)
	check(err)
}

func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}
