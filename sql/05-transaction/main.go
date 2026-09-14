// 05-transaction: 事务操作
//
// 演示 database/sql 事务：
//   - db.Begin 开始事务
//   - tx.Commit 提交
//   - tx.Rollback 回滚
//   - 事务隔离级别的概念
//
// 运行: go run ./sql/05-transaction/

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
		DROP TABLE IF EXISTS accounts;
		CREATE TABLE accounts (
			id       SERIAL PRIMARY KEY,
			name     TEXT    NOT NULL,
			balance  INT     NOT NULL DEFAULT 0
		);
		INSERT INTO accounts (name, balance) VALUES
			('Alice', 1000),
			('Bob',   500);
	`)

	showBalances := func(label string) {
		rows, err := db.QueryContext(ctx, "SELECT name, balance FROM accounts ORDER BY id")
		check(err)
		defer rows.Close()
		fmt.Printf("[%s] ", label)
		for rows.Next() {
			var name string
			var bal int
			check(rows.Scan(&name, &bal))
			fmt.Printf("%s=%d  ", name, bal)
		}
		fmt.Println()
	}

	showBalances("初始")

	// ---- 成功事务：转账 ----
	// Begin 从连接池取一条连接，在该连接上开启事务。
	// 后续所有操作都在同一条连接上执行。
	tx, err := db.BeginTx(ctx, nil)
	check(err)

	// Alice 转 200 给 Bob
	if _, err := tx.ExecContext(ctx,
		"UPDATE accounts SET balance = balance - 200 WHERE name = 'Alice'"); err != nil {
		tx.Rollback()
		check(err)
	}
	if _, err := tx.ExecContext(ctx,
		"UPDATE accounts SET balance = balance + 200 WHERE name = 'Bob'"); err != nil {
		tx.Rollback()
		check(err)
	}

	// Commit 将修改持久化。失败则所有操作都不会生效。
	check(tx.Commit())
	showBalances("转账后")

	// ---- 回滚演示 ----
	tx2, err := db.BeginTx(ctx, nil)
	check(err)

	if _, err := tx2.ExecContext(ctx,
		"UPDATE accounts SET balance = balance - 9999 WHERE name = 'Alice'"); err != nil {
		tx2.Rollback()
		check(err)
	}
	// 故意回滚，模拟业务逻辑判断不通过
	fmt.Println("\n[回滚] 检测到异常，执行 Rollback")
	check(tx2.Rollback())
	showBalances("回滚后")

	// ---- 清理 ----
	mustExec(ctx, db, "DROP TABLE IF EXISTS accounts;")
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
