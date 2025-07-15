# ===== 7. 簡単な起動スクリプト =====
# ファイル名: start.sh
#!/bin/bash

echo "🐳 Starting Go Todo App with Docker..."

# ビルド & 起動
docker-compose up --build -d

echo "✅ Go Todo App is running at http://localhost:8080"
echo "📊 Health check: docker-compose ps"
echo "📋 Logs: docker-compose logs -f go-todo"
echo "🛑 Stop: docker-compose down"
