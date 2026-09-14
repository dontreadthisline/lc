#!/bin/bash
# 完整扩容测试

set -e
cd "$(dirname "$0")"

echo "=== 步骤1: 3节点配置 ==="
cat > cluster.json << 'EOF'
{"nodes":[{"id":"node1","addr":"localhost:8001","http_port":8001},{"id":"node2","addr":"localhost:8002","http_port":8002},{"id":"node3","addr":"localhost:8003","http_port":8003}]}
EOF
cat cluster.json
echo ""

echo "=== 步骤2: 启动3个节点 ==="
./hashnode -id=node1 -port=8001 > /tmp/node1.log 2>&1 &
echo "n1 pid: $!"
sleep 1
./hashnode -id=node2 -port=8002 > /tmp/node2.log 2>&1 &
echo "n2 pid: $!"
sleep 1
./hashnode -id=node3 -port=8003 > /tmp/node3.log 2>&1 &
echo "n3 pid: $!"
sleep 2

echo ""
echo "=== 步骤3: 写入20条数据 ==="
for i in {1..20}; do
    curl -s -X PUT -d "val-$i" "http://localhost:8001/key/key-$i" > /dev/null
done

echo "当前分布:"
echo -n "N1: "; curl -s "http://localhost:8001/status" | grep -o '"data_count":[0-9]*'
echo -n "N2: "; curl -s "http://localhost:8002/status" | grep -o '"data_count":[0-9]*'
echo -n "N3: "; curl -s "http://localhost:8003/status" | grep -o '"data_count":[0-9]*'

echo ""
echo "=== 步骤4: 修改配置为4节点 ==="
cat > cluster.json << 'EOF'
{"nodes":[{"id":"node1","addr":"localhost:8001","http_port":8001},{"id":"node2","addr":"localhost:8002","http_port":8002},{"id":"node3","addr":"localhost:8003","http_port":8003},{"id":"node4","addr":"localhost:8004","http_port":8004}]}
EOF
echo "配置已更新"

echo ""
echo "=== 步骤5: 启动node4 ==="
./hashnode -id=node4 -port=8004 > /tmp/node4.log 2>&1 &
echo "n4 pid: $!"

echo ""
echo "=== 步骤6: 等待配置同步（6秒）==="
sleep 6

echo ""
echo "=== 步骤7: 检查迁移日志 ==="
echo "Node1:"
grep "配置同步\|迁移" /tmp/node1.log 2>/dev/null | tail -5 || echo "(无日志)"
echo ""
echo "Node2:"
grep "配置同步\|迁移" /tmp/node2.log 2>/dev/null | tail -5 || echo "(无日志)"
echo ""
echo "Node3:"
grep "配置同步\|迁移" /tmp/node3.log 2>/dev/null | tail -5 || echo "(无日志)"

echo ""
echo "=== 步骤8: 扩容后分布 ==="
echo -n "N1: "; curl -s "http://localhost:8001/status" | grep -o '"data_count":[0-9]*'
echo -n "N2: "; curl -s "http://localhost:8002/status" | grep -o '"data_count":[0-9]*'
echo -n "N3: "; curl -s "http://localhost:8003/status" | grep -o '"data_count":[0-9]*'
echo -n "N4: "; curl -s "http://localhost:8004/status" | grep -o '"data_count":[0-9]*'

echo ""
echo "=== 步骤9: 查询测试 ==="
echo "查询 key-1:"
curl -s "http://localhost:8001/key/key-1"
echo ""

echo ""
echo "=== 清理 ==="
pkill -f hashnode 2>/dev/null || true
echo "测试完成！"
