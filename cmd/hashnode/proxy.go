package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// ========== 查询转发逻辑 ==========

// isMyKey 判断 key 是否属于当前节点
func (s *Server) isMyKey(key string) bool {
	targetNode := s.hashRing.Get(key)
	return targetNode == s.self.ID
}

// forwardRequest 把请求转发给目标节点
func (s *Server) forwardRequest(node Node, method, key string, body []byte) ([]byte, error) {
	url := fmt.Sprintf("http://%s/internal/%s", node.Addr, key)
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// handleGet 处理查询请求
func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/key/"):]
	if key == "" {
		http.Error(w, "key required", http.StatusBadRequest)
		return
	}

	// 算出该 key 属于哪个节点
	targetNodeID := s.hashRing.Get(key)

	// 如果属于自己，直接返回
	if targetNodeID == s.self.ID {
		s.dataMu.RLock()
		value, exists := s.data[key]
		s.dataMu.RUnlock()

		if !exists {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"key":   key,
			"value": value,
			"node":  s.self.ID,
		})
		return
	}

	// 不属于自己，转发给目标节点
	node, ok := s.getNodeByID(targetNodeID)
	if !ok {
		http.Error(w, "node not found: "+targetNodeID, http.StatusInternalServerError)
		return
	}

	log.Printf("[转发 GET] key=%s 从 %s 转发到 %s", key, s.self.ID, targetNodeID)
	body, err := s.forwardRequest(node, "GET", key, nil)
	if err != nil {
		http.Error(w, "forward failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// handlePut 处理写入请求
type PutRequest struct {
	Value string `json:"value"`
}

func (s *Server) handlePut(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/key/"):]
	if key == "" {
		http.Error(w, "key required", http.StatusBadRequest)
		return
	}

	body, _ := io.ReadAll(r.Body)

	// 算出该 key 属于哪个节点
	targetNodeID := s.hashRing.Get(key)

	// 如果属于自己，直接存储
	if targetNodeID == s.self.ID {
		var req PutRequest
		if err := json.Unmarshal(body, &req); err != nil {
			// 也支持纯文本
			req.Value = string(body)
		}

		s.dataMu.Lock()
		s.data[key] = req.Value
		s.dataMu.Unlock()

		log.Printf("[存储] key=%s 存储在 %s", key, s.self.ID)
		json.NewEncoder(w).Encode(map[string]string{
			"key":    key,
			"value":  req.Value,
			"node":   s.self.ID,
			"status": "stored",
		})
		return
	}

	// 不属于自己，转发给目标节点
	node, ok := s.getNodeByID(targetNodeID)
	if !ok {
		http.Error(w, "node not found: "+targetNodeID, http.StatusInternalServerError)
		return
	}

	log.Printf("[转发 PUT] key=%s 从 %s 转发到 %s", key, s.self.ID, targetNodeID)
	respBody, err := s.forwardRequest(node, "PUT", key, body)
	if err != nil {
		http.Error(w, "forward failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBody)
}

// handleInternalGet 内部接口，直接返回本地数据（不转发）
func (s *Server) handleInternalGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/internal/"):]

	s.dataMu.RLock()
	value, exists := s.data[key]
	s.dataMu.RUnlock()

	if !exists {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"key":   key,
		"value": value,
		"node":  s.self.ID,
	})
}

// handleInternalPut 内部接口，直接存储到本地（不转发）
func (s *Server) handleInternalPut(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/internal/"):]
	body, _ := io.ReadAll(r.Body)

	var req PutRequest
	if err := json.Unmarshal(body, &req); err != nil {
		req.Value = string(body)
	}

	s.dataMu.Lock()
	s.data[key] = req.Value
	s.dataMu.Unlock()

	json.NewEncoder(w).Encode(map[string]string{
		"key":    key,
		"value":  req.Value,
		"node":   s.self.ID,
		"status": "stored",
	})
}
