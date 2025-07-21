package main

import (
	"fmt"
	"go-todo-app/config"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Todo構造体（シンプル版）
type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// In-Memoryデータストア
var (
	todos  []Todo
	nextID = 1
	mu     sync.RWMutex
)

// CRUD操作関数

// 全てのTodoを取得
func getTodos() []Todo {
	mu.RLock()
	defer mu.RUnlock()

	result := make([]Todo, len(todos))
	copy(result, todos)
	return result
}

// 新しいTodoを作成
func createTodo(title string) Todo {
	mu.Lock()
	defer mu.Unlock()

	todo := Todo{
		ID:        nextID,
		Title:     strings.TrimSpace(title),
		Completed: false,
	}

	todos = append(todos, todo)
	nextID++

	return todo
}

// Todoの完了状態を切り替え
func toggleTodo(id int) error {
	mu.Lock()
	defer mu.Unlock()

	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Completed = !todos[i].Completed
			return nil
		}
	}
	return fmt.Errorf("todo with id %d not found", id)
}

// Todoを削除
func deleteTodo(id int) error {
	mu.Lock()
	defer mu.Unlock()

	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("todo with id %d not found", id)
}

// HTMLテンプレート（シンプル版）
const htmlTemplate = `
<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go Todo App</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 600px;
            margin: 50px auto;
            padding: 20px;
            background-color: #f9f9f9;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            text-align: center;
            margin-bottom: 30px;
        }
        .add-form {
            margin-bottom: 30px;
            display: flex;
            gap: 10px;
        }
        .add-form input[type="text"] {
            flex: 1;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 16px;
        }
        .add-form button {
            padding: 10px 20px;
            background: #007bff;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
        }
        .add-form button:hover {
            background: #0056b3;
        }
        .todo-list {
            list-style: none;
            padding: 0;
        }
        .todo-item {
            display: flex;
            align-items: center;
            padding: 15px;
            margin-bottom: 10px;
            background: #f8f9fa;
            border-radius: 4px;
            border-left: 3px solid #007bff;
        }
        .todo-item.completed {
            opacity: 0.7;
            border-left-color: #28a745;
        }
        .todo-item.completed .todo-title {
            text-decoration: line-through;
        }
        .todo-checkbox {
            margin-right: 15px;
        }
        .todo-title {
            flex: 1;
            font-size: 16px;
            color: #333;
        }
        .delete-btn {
            padding: 5px 10px;
            background: #dc3545;
            color: white;
            border: none;
            border-radius: 3px;
            cursor: pointer;
            font-size: 12px;
        }
        .delete-btn:hover {
            background: #c82333;
        }
        .empty-message {
            text-align: center;
            color: #666;
            font-style: italic;
            padding: 40px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>📝 Go Todo App</h1>
        
        <!-- 新規追加フォーム -->
        <form class="add-form" action="/add" method="POST">
            <input type="text" name="title" placeholder="新しいタスクを入力..." required>
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
                <span class="todo-title">{{.Title}}</span>
                <form action="/delete/{{.ID}}" method="POST" style="display: inline;">
                    <button type="submit" class="delete-btn" 
                            onclick="return confirm('本当に削除しますか？')">削除</button>
                </form>
            </li>
            {{end}}
        </ul>
        {{else}}
        <div class="empty-message">
            タスクがありません。<br>
            上のフォームから新しいタスクを追加してください。
        </div>
        {{end}}
    </div>
</body>
</html>
`

// HTTPハンドラー

// メインページ（Todo一覧表示）
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl := template.Must(template.New("index").Parse(htmlTemplate))

	data := struct {
		Todos []Todo
	}{
		Todos: getTodos(),
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Todo追加
func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	if strings.TrimSpace(title) == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	createTodo(title)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Todo完了切り替え
func toggleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URLから ID を抽出
	path := strings.TrimPrefix(r.URL.Path, "/toggle/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	if err := toggleTodo(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Todo削除
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URLから ID を抽出
	path := strings.TrimPrefix(r.URL.Path, "/delete/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	if err := deleteTodo(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// メイン関数
func main() {
	cfg := config.LoadConfig()

	// 初期データを追加（デモ用）
	createTodo("Go言語を学習する")
	createTodo("シンプルなTodoアプリを作る")

	// ルート設定
	http.HandleFunc("/", indexHandler)         // メインページ
	http.HandleFunc("/add", addHandler)        // Todo追加
	http.HandleFunc("/toggle/", toggleHandler) // 完了切り替え
	http.HandleFunc("/delete/", deleteHandler) // Todo削除
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}) // ヘルスチェック（Docker用）

	port := cfg.Port

	fmt.Printf("🚀 Go Todo Server starting on :%s\n", port)
	fmt.Printf("📝 Initial todos: %d\n", len(getTodos()))
	fmt.Printf("🌐 Open http://localhost:%s in your browser\n", port)
	fmt.Printf("⚙️  Using config: DB=%s@%s:%s/%s\n",
		cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// 設定されたポートでサーバー起動
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
