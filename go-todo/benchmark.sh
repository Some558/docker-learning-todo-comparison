# ===== 6. 性能テスト用スクリプト =====
# ファイル名: benchmark.sh
#!/bin/bash

echo "🚀 Go vs C# Todo App Performance Comparison"
echo "=========================================="

# Go Todo起動
echo "📦 Building and starting Go Todo..."
docker-compose up -d go-todo

# 起動待機
echo "⏳ Waiting for services to start..."
sleep 15

# Go版テスト
echo "🔥 Testing Go Todo App..."
docker run --rm --network go-todo-app_todo-network williamyeh/wrk \
  -t12 -c400 -d30s --latency http://go-todo:8080/ > go-results.txt

echo "📊 Go Results:"
cat go-results.txt

# C#版テスト (別途起動が必要)
# echo "🔥 Testing C# Todo App..."
# docker run --rm --network csharp-todo_default williamyeh/wrk \
#   -t12 -c400 -d30s --latency http://csharp-todo:8080/ > csharp-results.txt

echo "✅ Performance test completed!"
echo "Go results saved to: go-results.txt"
# echo "C# results saved to: csharp-results.txt"