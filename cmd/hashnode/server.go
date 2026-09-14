package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"demo/algo"
)

// Node 表示集群中的一个节点
type Node struct {
	ID       string `json:"id"`
	Addr     string `json:"addr"`     // 对外服务的 HTTP 地址
	HTTPPort int    `json:"http_port"`
}

// ClusterConfig 集群配置
type ClusterConfig struct {
	Nodes []Node `json:"nodes"`
}

// Server 是一个一致性哈希节点
type Server struct {
	self         Node
	config       ClusterConfig
	configMu     sync.RWMutex
	hashRing     *algo.ConsistentHash
	configFile   string        // 配置文件路径
	syncInterval time.Duration // 配置同步间隔

	// 本地存储
	data   map[string]string
	dataMu sync.RWMutex

	// HTTP 客户端，用于转发请求
	client *http.Client
}

func NewServer(self Node, config ClusterConfig, configFile string) *Server {
	s := &Server{
		self:         self,
		config:       config,
		configFile:   configFile,
		syncInterval: 5 * time.Second, // 默认5秒同步一次
		data:         make(map[string]string),
		client:       &http.Client{Timeout: 5 * time.Second},
	}
	s.rebuildRing()
	return s
}

// handleStatus 返回节点状态
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	s.dataMu.RLock()
	dataCount := len(s.data)
	s.dataMu.RUnlock()

	s.configMu.RLock()
	nodes := s.config.Nodes
	s.configMu.RUnlock()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"node_id":    s.self.ID,
		"addr":       s.self.Addr,
		"data_count": dataCount,
		"nodes":      nodes,
	})
}

func (s *Server) Run() {
	// 启动配置同步 goroutine
	go s.startConfigSync()

	mux := http.NewServeMux()

	// 对外接口
	mux.HandleFunc("/key/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.handleGet(w, r)
		case http.MethodPut, http.MethodPost:
			s.handlePut(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// 内部接口（直接操作本地数据，不转发）
	mux.HandleFunc("/internal/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.handleInternalGet(w, r)
		case http.MethodPut, http.MethodPost:
			s.handleInternalPut(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// 管理接口
	mux.HandleFunc("/migrate", s.handleMigrate)
	mux.HandleFunc("/status", s.handleStatus)

	log.Printf("[启动] 节点 %s 监听 %s", s.self.ID, s.self.Addr)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.self.HTTPPort), mux))
}
