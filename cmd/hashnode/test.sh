#!/bin/bash
# 测试一致性哈希分布式节点

set -e

echo "=== 启动 3 个节点 ==="
go run . -id=node1 -port=8001 &
node1_pid=$!
go run . -id=node2 -port=8002 &
node2_pid=$!
go run . -id=node3 -port=8003 &
node3_pid=$!
echo "$node1_pid $node2_pid $node3_pid" > /tmp/hashnode_pids

sleep 2

echo ""
echo "=== 测试 1: 写入 10 条数据 ==="
for i in {1..10}; do
    curl -s -X PUT -d "value-data-$i" "http://localhost:8001/key/key-$i" | jq .
done

echo ""
echo "=== 测试 2: 查询数据（会转发到正确节点） ==="
curl -s "http://localhost:8002/key/key-1" | jq .
curl -s "http://localhost:8003/key/key-5" | jq .

echo ""
echo "=== 测试 3: 查看各节点状态 ==="
echo "Node1:"
curl -s "http://localhost:8001/status" | jq .
echo "Node2:"
curl -s "http://localhost:8002/status" | jq .
echo "Node3:"
curl -s "http://localhost:8003/status" | jq .

echo ""
echo "=== 测试完成，清理进程 ==="
kill $(cat /tmp/hashnode_pids) 2>/dev/null || true
rm -f /tmp/hashnode_pids

echo "完成！"
