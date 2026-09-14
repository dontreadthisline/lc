#!/bin/bash
# 测试扩容和数据迁移

set -e

echo "=== 阶段 1: 启动 2 个节点 ==="
go run . -id=node1 -port=8001 &
node1_pid=$!
go run . -id=node2 -port=8002 &
node2_pid=$!

sleep 2

echo ""
echo "=== 写入 100 条数据 ==="
for i in {1..100}; do
    curl -s -X PUT -d "value-$i" "http://localhost:8001/key/key-$i" > /dev/null
done
echo "写入完成"

echo ""
echo "=== 查看扩容前数据分布 ==="
echo "Node1 数据量:"
curl -s "http://localhost:8001/status" | jq '.data_count'
echo "Node2 数据量:"
curl -s "http://localhost:8002/status" | jq '.data_count'

# 停止节点2，模拟下线
echo ""
echo "=== 阶段 2: 停止 node2，启动 node3（模拟扩容替换）==="
kill $node2_pid 2>/dev/null || true
sleep 1

# 使用新的配置文件（包含 node1 和 node3）
cat > cluster_new.json << 'EOF'
{
  "nodes": [
    {"id": "node1", "addr": "localhost:8001", "http_port": 8001},
    {"id": "node3", "addr": "localhost:8003", "http_port": 8003}
  ]
}
EOF

go run . -id=node3 -port=8003 -config=cluster_new.json &
node3_pid=$!

# 同时也要让 node1 使用新配置
echo "请手动重启 node1 使用新配置，或实现配置热加载"
echo "这里简化为: node1 继续使用旧配置，新写入会根据新环路由"

sleep 2

echo ""
echo "=== 从 node3 查询之前写入的数据 ==="
echo "（如果 node3 接管了某些 key，会转发到 node1 获取）"
curl -s "http://localhost:8003/key/key-1" | jq .
curl -s "http://localhost:8003/key/key-50" | jq .

echo ""
echo "=== 清理 ==="
kill $node1_pid $node3_pid 2>/dev/null || true
rm -f cluster_new.json

echo "完成！"
