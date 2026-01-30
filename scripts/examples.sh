#!/bin/bash

# GitHub Browser - Usage Examples

BASE_URL="http://localhost:9527"

echo "🧪 GitHub Browser - Usage Examples"
echo "==================================="
echo ""

# 示例 1: 打开仓库
echo "Example 1: Open a repository"
echo "-----------------------------"
echo "URL: https://github.com/golang/go"
echo ""
curl -X POST $BASE_URL/open \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://github.com/golang/go",
    "ide": "code"
  }' | jq .
echo ""
echo ""

# 示例 2: 打开文件（带行号）
echo "Example 2: Open a file with line number"
echo "----------------------------------------"
echo "URL: https://github.com/golang/go/blob/master/src/runtime/proc.go#L123"
echo ""
curl -X POST $BASE_URL/open \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://github.com/golang/go/blob/master/src/runtime/proc.go#L123",
    "ide": "code"
  }' | jq .
echo ""
echo ""

# 示例 3: 打开 Pull Request
echo "Example 3: Open a Pull Request"
echo "-------------------------------"
echo "URL: https://github.com/golang/go/pull/12345"
echo ""
echo "Note: This will fail if PR #12345 doesn't exist, but demonstrates the API"
curl -X POST $BASE_URL/open \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://github.com/golang/go/pull/12345",
    "ide": "code"
  }' | jq .
echo ""
echo ""

# 示例 4: 查看缓存
echo "Example 4: List cached repositories"
echo "------------------------------------"
curl -s $BASE_URL/cache | jq .
echo ""
echo ""

# 示例 5: 健康检查
echo "Example 5: Health check"
echo "-----------------------"
curl -s $BASE_URL/health | jq .
echo ""
echo ""

# 示例 6: 获取配置
echo "Example 6: Get configuration"
echo "----------------------------"
curl -s $BASE_URL/config | jq .
echo ""
echo ""

echo "✅ Examples complete!"
echo ""
echo "💡 Tips:"
echo "  - Replace 'code' with your preferred IDE (zed, idea, etc.)"
echo "  - Use real PR numbers for PR examples"
echo "  - Check GUIDE.md for more information"
