// 06-nullable: 可空类型处理
//
// database/sql 提供的可空类型：
//   - sql.NullString
//   - sql.NullInt64 (注意：没有 NullInt，int 类型的 NULL 也用 NullInt64)
//   - sql.NullFloat64
//   - sql.NullBool
//   - sql.NullTime
//   - sql.Null[T]  (Go 1.22+ 泛型版本)
//
// 运行: go run ./sql/06-nullable/

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

	// ---- 建表：包含 NULL 值的字段 ----
	mustExec(ctx, db, `
		DROP TABLE IF EXISTS employees;
		CREATE TABLE employees (
			id         SERIAL PRIMARY KEY,
			name       TEXT NOT NULL,
			email      TEXT,           -- 可空
			manager_id INT,            -- 可空
			salary     NUMERIC,        -- 可空
			left_at    TIMESTAMPTZ     -- 可空（离职时间）
		);
		INSERT INTO employees (name, email, manager_id, salary, left_at) VALUES
			('Alice', 'alice@example.com', NULL, 80000, NULL),
			('Bob',   NULL,                1,    NULL,    '2024-06-15'),
			('Carol', 'carol@example.com', 1,    60000,   NULL);
	`)

	rows, err := db.QueryContext(ctx,
		"SELECT name, email, manager_id, salary, left_at FROM employees ORDER BY id")
	check(err)
	defer rows.Close()

	fmt.Println("员工列表:")
	for rows.Next() {
		var name string
		var email sql.NullString
		var managerID sql.NullInt64
		var salary sql.NullFloat64
		var leftAt sql.NullTime

		check(rows.Scan(&name, &email, &managerID, &salary, &leftAt))

		fmt.Printf("  %s\n", name)
		fmt.Printf("    email:      %s\n", nullableStr(email))
		fmt.Printf("    manager_id: %s\n", nullableInt(managerID))
		fmt.Printf("    salary:     %s\n", nullableFloat(salary))
		fmt.Printf("    left_at:    %s\n", nullableTime(leftAt))
		fmt.Println()
	}
	check(rows.Err())

	// ---- 写入 NULL ----
	// 要写入 NULL，需要显式传 sql.NullXxx{Valid: false}
	_, err = db.ExecContext(ctx,
		"INSERT INTO employees (name, email, manager_id, salary) VALUES ($1, $2, $3, $4)",
		"Dave",
		sql.NullString{Valid: false}, // NULL
		sql.NullInt64{Valid: false},  // NULL
		sql.NullFloat64{Valid: false}, // NULL
	)
	check(err)
	fmt.Println("[写入 NULL] 插入 Dave（所有可空字段为 NULL）")

	// ---- 清理 ----
	mustExec(ctx, db, "DROP TABLE IF EXISTS employees;")
}

func nullableStr(n sql.NullString) string {
	if !n.Valid {
		return "NULL"
	}
	return n.String
}

func nullableInt(n sql.NullInt64) string {
	if !n.Valid {
		return "NULL"
	}
	return fmt.Sprintf("%d", n.Int64)
}

func nullableFloat(n sql.NullFloat64) string {
	if !n.Valid {
		return "NULL"
	}
	return fmt.Sprintf("%.2f", n.Float64)
}

func nullableTime(n sql.NullTime) string {
	if !n.Valid {
		return "NULL"
	}
	return n.Time.Format(time.RFC3339)
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
