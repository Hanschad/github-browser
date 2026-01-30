#!/bin/bash

echo "🧪 Testing GitHub Browser Service"
echo "=================================="
echo ""

BASE_URL="http://localhost:9527"

# 测试健康检查
echo "1️⃣  Testing health check..."
curl -s $BASE_URL/health | jq .
echo ""

# 测试解析仓库 URL
echo "2️⃣  Testing repository URL parsing..."
echo "URL: https://github.com/golang/go"
curl -s -X POST $BASE_URL/open \
  -H "Content-Type: application/json" \
  -d '{"url": "https://github.com/golang/go"}' | jq .
echo ""

# 测试解析文件 URL
echo "3️⃣  Testing file URL parsing..."
echo "URL: https://github.com/golang/go/blob/master/src/runtime/proc.go#L123"
curl -s -X POST $BASE_URL/open \
  -H "Content-Type: application/json" \
  -d '{"url": "https://github.com/golang/go/blob/master/src/runtime/proc.go#L123"}' | jq .
echo ""

# 测试解析 PR URL
echo "4️⃣  Testing PR URL parsing..."
echo "URL: https://github.com/golang/go/pull/12345"
curl -s -X POST $BASE_URL/open \
  -H "Content-Type: application/json" \
  -d '{"url": "https://github.com/golang/go/pull/12345"}' | jq .
echo ""

# 查看缓存
echo "5️⃣  Listing cache..."
curl -s $BASE_URL/cache | jq .
echo ""

# 获取配置
echo "6️⃣  Getting config..."
curl -s $BASE_URL/config | jq .
echo ""

echo "✅ Tests complete!"
