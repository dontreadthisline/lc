package algo_tests

import "demo/algo"

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

// 构造 n 台机器名
func makeMachines(n int) []string {
	ms := make([]string, n)
	for i := 0; i < n; i++ {
		ms[i] = fmt.Sprintf("machine-%02d", i)
	}
	return ms
}

// 生成 count 个随机 key
func makeKeys(count int) []string {
	keys := make([]string, count)
	for i := 0; i < count; i++ {
		keys[i] = fmt.Sprintf("key:%d:%d", i, rand.Int63())
	}
	return keys
}

// 统计每个 key 落到哪台机器，返回 机器->key数 的分布
func distribute(c *algo.ConsistentHash, keys []string) map[string]int {
	dist := make(map[string]int)
	for _, k := range keys {
		dist[c.Get(k)]++
	}
	return dist
}

// TestDistribution 验证虚拟节点让分布大致均匀。
// 10 机器、150 vnodes、100 万 key，每台机器分到的 key 数应在均值的 ±20% 内。
func TestDistribution(t *testing.T) {
	const (
		machineN = 10
		replicas = 150
		keyN     = 1_000_000
	)
	c := algo.NewConsistentHash(replicas)
	c.Add(makeMachines(machineN)...)
	keys := makeKeys(keyN)
	dist := distribute(c, keys)

	avg := float64(keyN) / float64(machineN)
	lo, hi := avg*0.8, avg*1.2
	var worst float64
	for m, n := range dist {
		deviation := math.Abs(float64(n) - avg)
		if deviation > worst {
			worst = deviation
		}
		if float64(n) < lo || float64(n) > hi {
			t.Errorf("machine %s got %d keys, want within [%.0f, %.0f]", m, n, lo, hi)
		}
	}
	t.Logf("分布：均值=%.0f，最大偏差=%.0f (%.2f%%)", avg, worst, worst/avg*100)
}

// TestMigrationAddOne 验证加一台机器后，迁移比例约 1/(N+1)。
// 10 台加到 11 台，理论上新机器分到约 1/11≈9.09% 的 key，
// 其余约 90.9% 的 key 不应变动。这是“一致性”的核心：只有落到新机器段的那部分 key 迁移。
func TestMigrationAddOne(t *testing.T) {
	const (
		machineN = 10
		replicas = 150
		keyN     = 1_000_000
	)
	ms := makeMachines(machineN)
	c := algo.NewConsistentHash(replicas)
	c.Add(ms...)
	keys := makeKeys(keyN)

	// 记录每个 key 原来落到哪台机器
	before := make(map[string]string, keyN)
	for _, k := range keys {
		before[k] = c.Get(k)
	}

	// 加一台新机器
	newMachine := "machine-new"
	c.Add(newMachine)

	// 重查，统计迁移
	migrated := 0
	migratedToNew := 0
	for _, k := range keys {
		now := c.Get(k)
		if now != before[k] {
			migrated++
			if now == newMachine {
				migratedToNew++
			} else {
				// 迁移的 key 必须都落到新机器，否则逻辑有 bug
				t.Errorf("key %s 迁移到 %s，应只迁到新机器 %s", k, now, newMachine)
			}
		}
	}
	ratio := float64(migrated) / float64(keyN)
	expected := 1.0 / 11.0
	t.Logf("加 1 台：迁移 %d/%d = %.2f%%，期望约 %.2f%%，全部迁到新机器=%t",
		migrated, keyN, ratio*100, expected*100, migrated == migratedToNew)
	// 允许 ±3% 波动（虚拟节点的随机性）
	if math.Abs(ratio-expected) > 0.03 {
		t.Errorf("迁移比例 %.4f 偏离期望 %.4f 超过 3%%", ratio, expected)
	}
}

// TestMigrationRemoveOne 验证删一台机器后，迁移比例约 1/N。
// 10 台删到 9 台，被删机器的 key（约 1/10≈10%）应分散到其余机器，其余 key 不动。
func TestMigrationRemoveOne(t *testing.T) {
	const (
		machineN = 10
		replicas = 150
		keyN     = 1_000_000
	)
	ms := makeMachines(machineN)
	c := algo.NewConsistentHash(replicas)
	c.Add(ms...)
	keys := makeKeys(keyN)

	before := make(map[string]string, keyN)
	for _, k := range keys {
		before[k] = c.Get(k)
	}

	// 删第一台机器
	removed := ms[0]
	c.Remove(removed)

	migrated := 0
	migratedFromRemoved := 0
	for _, k := range keys {
		now := c.Get(k)
		if now != before[k] {
			migrated++
			// 只有原来在 removed 上的 key 才该迁移
			if before[k] == removed {
				migratedFromRemoved++
			} else {
				t.Errorf("key %s 原在 %s（非被删机器）却迁移了，一致性被破坏", k, before[k])
			}
		} else {
			// 原来在被删机器上的 key，删后必然迁移，没迁说明有 bug
			if before[k] == removed {
				t.Errorf("key %s 原在被删机器 %s 上，删除后未迁移", k, removed)
			}
		}
	}
	ratio := float64(migrated) / float64(keyN)
	expected := 1.0 / 10.0
	t.Logf("删 1 台：迁移 %d/%d = %.2f%%，期望约 %.2f%%，迁移全部来自被删机器=%t",
		migrated, keyN, ratio*100, expected*100, migrated == migratedFromRemoved)
	if math.Abs(ratio-expected) > 0.03 {
		t.Errorf("迁移比例 %.4f 偏离期望 %.4f 超过 3%%", ratio, expected)
	}
}

// TestEmptyRing 环空时 Get 返回空串。
func TestEmptyRing(t *testing.T) {
	c := algo.NewConsistentHash(150)
	if got := c.Get("anykey"); got != "" {
		t.Errorf("空环 Get 返回 %q，期望空串", got)
	}
}

// --- benchmark ---

// BenchmarkGet 测纯查询吞吐。10 机器 / 150 vnodes，每次迭代查一个 key。
func BenchmarkGet(b *testing.B) {
	c := algo.NewConsistentHash(150)
	c.Add(makeMachines(10)...)
	keys := makeKeys(1_000_000)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.Get(keys[i%len(keys)])
	}
}

// BenchmarkAdd 测加机器（含环重排）的开销。
func BenchmarkAdd(b *testing.B) {
	ms := makeMachines(10)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := algo.NewConsistentHash(150)
		c.Add(ms...)
	}
}

// BenchmarkRemove 测删机器的开销。
func BenchmarkRemove(b *testing.B) {
	ms := makeMachines(10)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		c := algo.NewConsistentHash(150)
		c.Add(ms...)
		b.StartTimer()
		c.Remove(ms[0])
	}
}
