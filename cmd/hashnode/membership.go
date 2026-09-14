package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"demo/algo"
)

// ========== 节点上下线 / 配置同步 / 数据迁移 ==========

// startConfigSync 定期从配置文件同步集群配置
// 当检测到节点列表变化时，自动重建哈希环并触发数据迁移
func (s *Server) startConfigSync() {
	ticker := time.NewTicker(s.syncInterval)
	defer ticker.Stop()

	for range ticker.C {
		s.syncConfig()
	}
}

// syncConfig 读取配置文件并同步配置
func (s *Server) syncConfig() {
	data, err := os.ReadFile(s.configFile)
	if err != nil {
		log.Printf("[配置同步] 读取配置文件失败: %v", err)
		return
	}

	var newConfig ClusterConfig
	if err := json.Unmarshal(data, &newConfig); err != nil {
		log.Printf("[配置同步] 解析配置文件失败: %v", err)
		return
	}

	s.configMu.Lock()
	defer s.configMu.Unlock()

	// 检查节点列表是否有变化
	if s.configChanged(s.config.Nodes, newConfig.Nodes) {
		log.Printf("[配置同步] 检测到节点变化: %v -> %v",
			s.nodeIDs(s.config.Nodes), s.nodeIDs(newConfig.Nodes))

		oldConfig := s.config
		s.config = newConfig

		// 重建哈希环
		s.rebuildRingUnlocked()

		// 检查自己是否还在集群中
		if !s.isInCluster() {
			log.Printf("[配置同步] 当前节点 %s 已不在集群中，准备下线", s.self.ID)
			return
		}

		// 触发数据迁移（异步）
		go s.doMigrate(&oldConfig)
	}
}

// configChanged 比较两个节点列表是否有变化
func (s *Server) configChanged(old, new []Node) bool {
	if len(old) != len(new) {
		return true
	}
	oldSet := make(map[string]bool)
	for _, n := range old {
		oldSet[n.ID] = true
	}
	for _, n := range new {
		if !oldSet[n.ID] {
			return true
		}
	}
	return false
}

// nodeIDs 提取节点 ID 列表
func (s *Server) nodeIDs(nodes []Node) []string {
	ids := make([]string, len(nodes))
	for i, n := range nodes {
		ids[i] = n.ID
	}
	return ids
}

// isInCluster 检查当前节点是否还在集群配置中
func (s *Server) isInCluster() bool {
	for _, n := range s.config.Nodes {
		if n.ID == s.self.ID {
			return true
		}
	}
	return false
}

// getNodeByID 需要在读锁保护下访问配置
func (s *Server) getNodeByID(id string) (Node, bool) {
	s.configMu.RLock()
	defer s.configMu.RUnlock()
	for _, node := range s.config.Nodes {
		if node.ID == id {
			return node, true
		}
	}
	return Node{}, false
}

// rebuildRingUnlocked 根据当前配置重建一致性哈希环（无锁版本，调用者需保证线程安全）
func (s *Server) rebuildRingUnlocked() {
	s.hashRing = algo.NewConsistentHash(150) // 150 虚拟节点
	for _, node := range s.config.Nodes {
		s.hashRing.Add(node.ID)
	}
	log.Printf("[重构环] 节点数: %d, 虚拟节点数: %d", len(s.config.Nodes), len(s.config.Nodes)*150)
}

// rebuildRing 根据当前配置重建一致性哈希环（带锁版本，供外部调用）
func (s *Server) rebuildRing() {
	s.configMu.RLock()
	nodes := s.config.Nodes
	s.configMu.RUnlock()

	s.hashRing = algo.NewConsistentHash(150) // 150 虚拟节点
	for _, node := range nodes {
		s.hashRing.Add(node.ID)
	}
	log.Printf("[重构环] 节点数: %d, 虚拟节点数: %d", len(nodes), len(nodes)*150)
}

// handleMigrate 触发数据迁移
func (s *Server) handleMigrate(w http.ResponseWriter, r *http.Request) {
	go s.doMigrate(nil)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "migration started",
	})
}

// doMigrate 扫描本地数据，把不属于当前节点的数据迁移出去
// oldConfig 是扩容前的配置，用于构建旧环判断数据归属
func (s *Server) doMigrate(oldConfig *ClusterConfig) {
	log.Printf("[迁移开始] %s", s.self.ID)

	// 构建旧环：基于扩容前的配置
	oldRing := algo.NewConsistentHash(150)
	if oldConfig != nil {
		for _, node := range oldConfig.Nodes {
			oldRing.Add(node.ID)
		}
	}

	// 第一步：收集需要迁移的 key 列表（持有读锁，时间很短）
	s.dataMu.RLock()
	type item struct {
		key   string
		value string
	}
	toMigrate := make([]item, 0)
	for key, value := range s.data {
		// 用旧环判断：之前属于我，现在不属于我，才需要迁移
		oldTarget := oldRing.Get(key)
		newTarget := s.hashRing.Get(key)
		if oldTarget == s.self.ID && newTarget != s.self.ID {
			toMigrate = append(toMigrate, item{key, value})
		}
	}
	s.dataMu.RUnlock()

	if len(toMigrate) == 0 {
		log.Printf("[迁移完成] %s 无需迁移", s.self.ID)
		return
	}

	log.Printf("[迁移进度] %s 需要迁移 %d 条数据", s.self.ID, len(toMigrate))

	// 第二步：逐个迁移（不持有锁，允许并发）
	migrated := 0
	failed := 0
	for _, it := range toMigrate {
		targetNodeID := s.hashRing.Get(it.key)
		if targetNodeID == s.self.ID {
			continue // 已经不需要迁移了（可能配置又变了）
		}

		node, ok := s.getNodeByID(targetNodeID)
		if !ok {
			log.Printf("[迁移错误] 找不到节点 %s", targetNodeID)
			failed++
			continue
		}

		body, _ := json.Marshal(PutRequest{Value: it.value})
		_, err := s.forwardRequest(node, "PUT", it.key, body)
		if err != nil {
			log.Printf("[迁移失败] key=%s 到 %s: %v", it.key, targetNodeID, err)
			failed++
			continue
		}

		// 迁移成功，删除本地数据（持有写锁，时间很短）
		s.dataMu.Lock()
		delete(s.data, it.key)
		s.dataMu.Unlock()

		migrated++
		if migrated%100 == 0 {
			log.Printf("[迁移进度] %s 已迁移 %d/%d", s.self.ID, migrated, len(toMigrate))
		}
	}

	log.Printf("[迁移完成] %s 成功:%d 失败:%d", s.self.ID, migrated, failed)
}
