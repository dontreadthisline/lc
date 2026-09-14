// 04-prepared: 预编译语句
//
// 演示 prepared statement 的用法和原理：
//   - db.Prepare 让数据库提前解析 SQL，后续执行只需传参数
//   - 适合同一 SQL 重复执行多次的场景
//   - pgx 在 database/sql 层默认也会使用服务端预编译
//
// 运行: go run ./sql/04-prepared/

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
		DROP TABLE IF EXISTS products;
		CREATE TABLE products (
			id    SERIAL PRIMARY KEY,
			name  TEXT    NOT NULL,
			price NUMERIC NOT NULL
		);
	`)

	// ---- Prepare: 预编译 INSERT ----
	// 内部流程：数据库解析 SQL → 生成执行计划 → 返回语句句柄。
	// 后续每次 Exec 只需要传参数，跳过解析阶段。
	insertStmt, err := db.PrepareContext(ctx,
		"INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id",
	)
	check(err)
	defer insertStmt.Close() // 释放数据库端的预编译资源

	products := []struct {
		name  string
		price float64
	}{
		{"键盘", 299.00},
		{"鼠标", 149.50},
		{"显示器", 1899.00},
	}

	for _, p := range products {
		var id int
		err := insertStmt.QueryRowContext(ctx, p.name, p.price).Scan(&id)
		check(err)
		fmt.Printf("[Prepare INSERT] %s → id=%d\n", p.name, id)
	}
	fmt.Println()

	// ---- Prepare: 预编译 SELECT ----
	queryStmt, err := db.PrepareContext(ctx,
		"SELECT id, name, price FROM products WHERE price > $1 ORDER BY price DESC",
	)
	check(err)
	defer queryStmt.Close()

	fmt.Println("[Prepare SELECT] 价格 > 200 的商品:")
	rows, err := queryStmt.QueryContext(ctx, 200.0)
	check(err)
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var price float64
		check(rows.Scan(&id, &name, &price))
		fmt.Printf("  %s  ¥%.2f\n", name, price)
	}
	check(rows.Err())

	// ---- 性能对比 ----
	// 预编译 vs 直接 Exec 的简单对比
	fmt.Println("\n[性能对比] 执行 1000 次相同结构的 INSERT:")
	benchmark(ctx, db, insertStmt)

	// ---- 清理 ----
	mustExec(ctx, db, "DROP TABLE IF EXISTS products;")
}

// benchmark 对比预编译语句和直接 Exec 的性能
func benchmark(ctx context.Context, db *sql.DB, stmt *sql.Stmt) {
	iterations := 1000

	// 使用预编译语句
	start := time.Now()
	for i := range iterations {
		var id int
		err := stmt.QueryRowContext(ctx, fmt.Sprintf("item-%d", i), float64(i)).Scan(&id)
		check(err)
	}
	preparedDuration := time.Since(start)
	fmt.Printf("  预编译语句: %v  (%d 次/秒)\n",
		preparedDuration,
		int(float64(iterations)/preparedDuration.Seconds()))

	// 直接用 db.QueryRow（每次都要解析 SQL）
	start = time.Now()
	for i := range iterations {
		var id int
		err := db.QueryRowContext(ctx,
			"INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id",
			fmt.Sprintf("item-direct-%d", i), float64(i),
		).Scan(&id)
		check(err)
	}
	directDuration := time.Since(start)
	fmt.Printf("  直接 QueryRow: %v  (%d 次/秒)\n",
		directDuration,
		int(float64(iterations)/directDuration.Seconds()))
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
