package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

// hashnode 是一致性哈希分布式节点的可执行入口。
// 算法实现见 demo/algo；节点逻辑按职责拆分：
//   - server.go     Server 结构、生命周期、HTTP 路由
//   - proxy.go      查询转发（对外 /key 与内部 /internal 接口）
//   - membership.go 节点上下线、配置同步、数据迁移
func main() {
	var (
		nodeID     = flag.String("id", "", "节点 ID")
		port       = flag.Int("port", 8001, "HTTP 端口")
		configFile = flag.String("config", "cluster.json", "集群配置文件")
	)
	flag.Parse()

	if *nodeID == "" {
		log.Fatal("请指定节点 ID: -id=node1")
	}

	// 读取集群配置
	configData, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	var config ClusterConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	// 找到自己
	var self Node
	found := false
	for _, node := range config.Nodes {
		if node.ID == *nodeID {
			self = node
			found = true
			break
		}
	}
	if !found {
		// 如果没找到，用命令行参数创建
		self = Node{
			ID:       *nodeID,
			Addr:     fmt.Sprintf("localhost:%d", *port),
			HTTPPort: *port,
		}
	}

	server := NewServer(self, config, *configFile)
	server.Run()
}
