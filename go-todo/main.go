package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-todo-app/config"
	"go-todo-app/database"
	"go-todo-app/models"
)

// HTMLテンプレート（PostgreSQL対応版）
const htmlTemplate = `
<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go Todo App - PostgreSQL版</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 700px;
            margin: 50px auto;
            padding: 20px;
            background-color: #f5f7fa;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 12px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        h1 {
            color: #2d3748;
            text-align: center;
            margin-bottom: 10px;
            font-size: 28px;
        }
        .subtitle {
            text-align: center;
            color: #4a5568;
            margin-bottom: 30px;
            font-size: 14px;
            background: #e6fffa;
            padding: 8px 16px;
            border-radius: 20px;
            display: inline-block;
            margin-left: 50%;
            transform: translateX(-50%);
        }
        .stats {
            display: flex;
            justify-content: space-around;
            margin-bottom: 30px;
            padding: 20px;
            background: #f7fafc;
            border-radius: 8px;
            border: 1px solid #e2e8f0;
        }
        .stat-item {
            text-align: center;
        }
        .stat-number {
            font-size: 24px;
            font-weight: bold;
            color: #3182ce;
        }
        .stat-label {
            font-size: 12px;
            color: #718096;
            margin-top: 4px;
        }
        .add-form {
            margin-bottom: 30px;
            display: flex;
            gap: 12px;
        }
        .add-form input[type="text"] {
            flex: 1;
            padding: 12px;
            border: 2px solid #e2e8f0;
            border-radius: 8px;
            font-size: 16px;
            transition: border-color 0.2s;
        }
        .add-form input[type="text"]:focus {
            outline: none;
            border-color: #3182ce;
        }
        .add-form button {
            padding: 12px 24px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            border-radius: 8px;
            cursor: pointer;
            font-size: 16px;
            font-weight: 600;
            transition: transform 0.2s;
        }
        .add-form button:hover {
            transform: translateY(-1px);
        }
        .todo-list {
            list-style: none;
            padding: 0;
        }
        .todo-item {
            display: flex;
            align-items: center;
            padding: 16px;
            margin-bottom: 12px;
            background: #ffffff;
            border: 1px solid #e2e8f0;
            border-radius: 8px;
            border-left: 4px solid #3182ce;
            transition: all 0.2s;
        }
        .todo-item:hover {
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        }
        .todo-item.completed {
            opacity: 0.7;
            border-left-color: #38a169;
            background: #f0fff4;
        }
        .todo-item.completed .todo-title {
            text-decoration: line-through;
        }
        .todo-checkbox {
            margin-right: 16px;
            transform: scale(1.2);
            cursor: pointer;
        }
        .todo-content {
            flex: 1;
        }
        .todo-title {
            font-size: 16px;
            color: #2d3748;
            margin-bottom: 4px;
        }
        .todo-meta {
            font-size: 12px;
            color: #718096;
        }
        .todo-actions {
            display: flex;
            gap: 8px;
        }
        .delete-btn {
            padding: 6px 12px;
            background: #e53e3e;
            color: white;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-size: 12px;
            transition: background-color 0.2s;
        }
        .delete-btn:hover {
            background: #c53030;
        }
        .empty-message {
            text-align: center;
            color: #718096;
            font-style: italic;
            padding: 60px 20px;
            background: #f7fafc;
            border-radius: 8px;
            border: 2px dashed #cbd5e0;
        }
        .empty-icon {
            font-size: 48px;
            margin-bottom: 16px;
        }
        .footer {
            margin-top: 30px;
            text-align: center;
            font-size: 12px;
            color: #a0aec0;
            padding-top: 20px;
            border-top: 1px solid #e2e8f0;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>📝 Go Todo App</h1>
        <div class="subtitle">
            🐘 PostgreSQL + GORM で永続化対応
        </div>
        
        <!-- 統計情報 -->
        <div class="stats">
            <div class="stat-item">
                <div class="stat-number">{{.Stats.Total}}</div>
                <div class="stat-label">総タスク</div>
            </div>
            <div class="stat-item">
                <div class="stat-number">{{.Stats.Completed}}</div>
                <div class="stat-label">完了済み</div>
            </div>
            <div class="stat-item">
                <div class="stat-number">{{.Stats.Pending}}</div>
                <div class="stat-label">未完了</div>
            </div>
        </div>
        
        <!-- 新規追加フォーム -->
        <form class="add-form" action="/add" method="POST">
            <input type="text" name="title" placeholder="新しいタスクを入力してください..." required maxlength="200">
            <button type="submit">追加</button>
        </form>
        
        <!-- Todo一覧 -->
        {{if .Todos}}
        <ul class="todo-list">
            {{range .Todos}}
            <li class="todo-item{{if .Completed}} completed{{end}}">
                <form action="/toggle/{{.ID}}" method="POST" style="display: inline;">
                    <input type="checkbox" class="todo-checkbox" {{if .Completed}}checked{{end}} 
                           onchange="this.form.submit()">
                </form>
                <div class="todo-content">
                    <div class="todo-title">{{.Title}}</div>
                    <div class="todo-meta">
                        ID: {{.ID}} | 作成: {{.CreatedAt.Format "2006/01/02 15:04"}}
                        {{if ne .CreatedAt .UpdatedAt}} | 更新: {{.UpdatedAt.Format "2006/01/02 15:04"}}{{end}}
                    </div>
                </div>
                <div class="todo-actions">
                    <form action="/delete/{{.ID}}" method="POST" style="display: inline;">
                        <button type="submit" class="delete-btn" 
                                onclick="return confirm('「{{.Title}}」を削除しますか？')">削除</button>
                    </form>
                </div>
            </li>
            {{end}}
        </ul>
        {{else}}
        <div class="empty-message">
            <div class="empty-icon">📝</div>
            <div>まだタスクがありません</div>
            <div>上のフォームから新しいタスクを追加してみましょう！</div>
        </div>
        {{end}}
        
        <div class="footer">
            Powered by Go + PostgreSQL + GORM + Docker Compose<br>
            エンジニア4ヶ月目の学習記録 💪
        </div>
    </div>
</body>
</html>
`

// 統計情報構造体
type TodoStats struct {
	Total     int
	Completed int
	Pending   int
}

// 統計情報を計算
func calculateStats(todos []models.Todo) TodoStats {
	stats := TodoStats{
		Total: len(todos),
	}

	for _, todo := range todos {
		if todo.Completed {
			stats.Completed++
		} else {
			stats.Pending++
		}
	}

	return stats
}

// HTTPハンドラー

// メインページ（Todo一覧表示）
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// データベースから全てのTodoを取得（作成日時の降順）
	var todos []models.Todo
	result := database.GetDB().Order("created_at desc").Find(&todos)
	if result.Error != nil {
		log.Printf("データベースエラー: %v", result.Error)
		http.Error(w, "データベースエラーが発生しました", http.StatusInternalServerError)
		return
	}

	log.Printf("Todoを%d件取得しました", len(todos))

	// テンプレートを解析
	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		log.Printf("テンプレート解析エラー: %v", err)
		http.Error(w, "テンプレートエラー", http.StatusInternalServerError)
		return
	}

	// テンプレートに渡すデータ
	data := struct {
		Todos []models.Todo
		Stats TodoStats
	}{
		Todos: todos,
		Stats: calculateStats(todos),
	}

	// テンプレートを実行してレスポンスに書き込み
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("テンプレート実行エラー: %v", err)
		http.Error(w, "テンプレート実行エラー", http.StatusInternalServerError)
		return
	}
}

// Todo追加ハンドラー
func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POSTメソッドのみ対応", http.StatusMethodNotAllowed)
		return
	}

	// フォームデータを解析
	if err := r.ParseForm(); err != nil {
		log.Printf("フォーム解析エラー: %v", err)
		http.Error(w, "フォームデータが無効です", http.StatusBadRequest)
		return
	}

	// タイトルを取得・バリデーション
	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		http.Error(w, "タイトルは必須です", http.StatusBadRequest)
		return
	}

	// 新しいTodoを作成
	todo := models.Todo{
		Title:     title,
		Completed: false,
	}

	// アプリケーションレベルのバリデーション
	if err := todo.Validate(); err != nil {
		log.Printf("バリデーションエラー: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// データベースに保存
	result := database.GetDB().Create(&todo)
	if result.Error != nil {
		log.Printf("Todo作成エラー: %v", result.Error)
		http.Error(w, "データベースエラー", http.StatusInternalServerError)
		return
	}

	log.Printf("新しいTodoを作成しました: ID=%d, Title=%s", todo.ID, todo.Title)

	// メインページにリダイレクト
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Todo完了切り替えハンドラー
func toggleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POSTメソッドのみ対応", http.StatusMethodNotAllowed)
		return
	}

	// URLからIDを抽出
	path := strings.TrimPrefix(r.URL.Path, "/toggle/")
	id, err := strconv.ParseUint(path, 10, 32)
	if err != nil {
		log.Printf("無効なTodo ID: %s", path)
		http.Error(w, "無効なTodo IDです", http.StatusBadRequest)
		return
	}

	// データベースからTodoを取得
	var todo models.Todo
	result := database.GetDB().First(&todo, id)
	if result.Error != nil {
		log.Printf("Todo取得エラー (ID=%d): %v", id, result.Error)
		http.Error(w, "Todoが見つかりません", http.StatusNotFound)
		return
	}

	// 完了状態を切り替え
	todo.Completed = !todo.Completed

	// データベースに保存
	result = database.GetDB().Save(&todo)
	if result.Error != nil {
		log.Printf("Todo更新エラー (ID=%d): %v", id, result.Error)
		http.Error(w, "データベースエラー", http.StatusInternalServerError)
		return
	}

	log.Printf("Todoの完了状態を更新しました: ID=%d, Completed=%t", todo.ID, todo.Completed)

	// メインページにリダイレクト
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Todo削除ハンドラー
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POSTメソッドのみ対応", http.StatusMethodNotAllowed)
		return
	}

	// URLからIDを抽出
	path := strings.TrimPrefix(r.URL.Path, "/delete/")
	id, err := strconv.ParseUint(path, 10, 32)
	if err != nil {
		log.Printf("無効なTodo ID: %s", path)
		http.Error(w, "無効なTodo IDです", http.StatusBadRequest)
		return
	}

	// データベースから削除
	result := database.GetDB().Delete(&models.Todo{}, id)
	if result.Error != nil {
		log.Printf("Todo削除エラー (ID=%d): %v", id, result.Error)
		http.Error(w, "データベースエラー", http.StatusInternalServerError)
		return
	}

	// 実際に削除されたかチェック
	if result.RowsAffected == 0 {
		log.Printf("削除対象のTodoが見つかりませんでした: ID=%d", id)
		http.Error(w, "Todoが見つかりません", http.StatusNotFound)
		return
	}

	log.Printf("Todoを削除しました: ID=%d", id)

	// メインページにリダイレクト
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ヘルスチェックハンドラー
func healthHandler(w http.ResponseWriter, r *http.Request) {
	// データベース接続確認
	sqlDB, err := database.GetDB().DB()
	if err != nil {
		log.Printf("データベース接続取得エラー: %v", err)
		http.Error(w, "データベース接続エラー", http.StatusInternalServerError)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		log.Printf("データベースPingエラー: %v", err)
		http.Error(w, "データベースに接続できません", http.StatusInternalServerError)
		return
	}

	// Todo件数を取得
	var count int64
	database.GetDB().Model(&models.Todo{}).Count(&count)

	// JSON形式でレスポンス
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{
	"status": "healthy",
	"database": "connected",
	"message": "PostgreSQL Todo App",
	"todo_count": %d,
	"timestamp": "%s"
}`, count, fmt.Sprintf("%v", time.Now().Format("2006-01-02 15:04:05")))
}

// メイン関数
func main() {
	log.Println("🚀 Go Todo Server (PostgreSQL版) を起動しています...")

	// 1. 設定読み込み
	cfg := config.LoadConfig()
	log.Printf("⚙️  設定読み込み完了: DB=%s@%s:%s/%s",
		cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// 2. データベース接続
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("❌ データベース接続失敗: %v", err)
	}

	// 3. マイグレーション実行
	if err := database.Migrate(); err != nil {
		log.Fatalf("❌ マイグレーション失敗: %v", err)
	}

	// 4. 初期データの確認
	var todoCount int64
	database.GetDB().Model(&models.Todo{}).Count(&todoCount)
	log.Printf("📊 現在のTodo件数: %d件", todoCount)

	// 初期データがない場合のサンプル作成（オプション）
	if todoCount == 0 {
		log.Println("📝 初期データを作成します...")
		sampleTodos := []models.Todo{
			{Title: "Go + PostgreSQL環境の構築", Completed: true},
			{Title: "GORMでのCRUD操作実装", Completed: true},
			{Title: "Docker Composeによる環境管理", Completed: false},
			{Title: "次のステップ: テスト実装", Completed: false},
		}

		for _, todo := range sampleTodos {
			database.GetDB().Create(&todo)
		}
		log.Printf("✅ サンプルデータを%d件作成しました", len(sampleTodos))
	}

	// 5. ルート設定
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/toggle/", toggleHandler)
	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/health", healthHandler)

	// 6. サーバー起動
	port := cfg.Port
	log.Printf("🌐 サーバーを起動しました")
	log.Printf("   📍 ポート: %s", port)
	log.Printf("   🔗 URL: http://localhost:%s", port)
	log.Printf("   🏥 ヘルスチェック: http://localhost:%s/health", port)
	log.Printf("   🐘 データベース: PostgreSQL (%s:%s/%s)", cfg.DBHost, cfg.DBPort, cfg.DBName)
	log.Println("🎉 準備完了！ブラウザでアクセスしてください")

	// HTTPサーバー起動
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("❌ サーバー起動エラー: %v", err)
	}
}
