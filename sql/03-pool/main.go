// 03-pool: 连接池配置与行为观察
//
// 演示 database/sql 连接池的核心参数：
//   - SetMaxOpenConns — 最大打开连接数（含使用中 + 空闲）
//   - SetMaxIdleConns — 最大空闲连接数
//   - SetConnMaxLifetime — 连接最大存活时间
//   - SetConnMaxIdleTime  — 空闲连接最大空闲时间
//
// 运行: go run ./sql/03-pool/

package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://didi@localhost:5432/postgres?sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)
	check(err)
	defer db.Close()

	// ---- 连接池配置 ----
	// MaxOpenConns ≤ 0 表示无限制（默认就是无限制）。
	// 实际环境中一定要设上限，避免打爆数据库。
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	printStats := func(label string) {
		s := db.Stats()
		fmt.Printf("[%-8s] Open=%d  InUse=%d  Idle=%d  WaitCount=%d  WaitDuration=%v\n",
			label, s.OpenConnections, s.InUse, s.Idle, s.WaitCount, s.WaitDuration)
	}

	printStats("初始")

	// ---- 模拟并发查询 ----
	// 同时发起 10 个查询，但 MaxOpenConns=5，所以后 5 个会排队等待。
	var wg sync.WaitGroup
	for i := range 10 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// pg_sleep(0.5) 让每个查询占用连接 0.5 秒，方便观察排队。
			var result int
			err := db.QueryRowContext(ctx, "SELECT $1", id).Scan(&result)
			if err != nil {
				fmt.Printf("  查询 %d 失败: %v\n", id, err)
			}
		}(i)
	}

	// 在并发执行期间多次采样连接池状态
	for range 4 {
		time.Sleep(50 * time.Millisecond)
		printStats("运行中")
	}

	wg.Wait()
	printStats("结束后")

	fmt.Println("\n要点:")
	fmt.Println("  - MaxOpenConns=5 限制同时最多 5 个连接")
	fmt.Println("  - 当全部占满时，新请求进入 WaitCount 排队")
	fmt.Println("  - 查询完成后连接回到 Idle 池，而不是立即关闭")
	fmt.Println("  - SetConnMaxLifetime 防止使用过久的连接（比如网络变化后失效）")
}

func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}
