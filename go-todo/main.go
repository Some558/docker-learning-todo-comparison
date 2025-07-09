package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Todo構造体の定義（C#版と同等）
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"createdAt"`
}

// In-Memoryデータストア
var (
	todos  []Todo
	nextID = 1
	mu     sync.RWMutex // 並行安全性のため
)

// CRUD操作関数

// 全てのTodoを取得
func getTodos() []Todo {
	mu.RLock()
	defer mu.RUnlock()

	// コピーを返して安全性を確保
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
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	todos = append(todos, todo)
	nextID++

	return todo
}

// IDでTodoを検索
func getTodoByID(id int) (Todo, bool) {
	mu.RLock()
	defer mu.RUnlock()

	for _, todo := range todos {
		if todo.ID == id {
			return todo, true
		}
	}
	return Todo{}, false
}

// Todoを更新
func updateTodo(id int, title string, completed bool) (Todo, error) {
	mu.Lock()
	defer mu.Unlock()

	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Title = title
			todos[i].Completed = completed
			return todos[i], nil
		}
	}
	return Todo{}, fmt.Errorf("todo with id %d not found", id)
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

// HTTPハンドラー（基本的なもの）
func todosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		todos := getTodos()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todos)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"healthy","message":"Go Todo Server","todos":%d}`, len(getTodos()))
}

func main() {
	// 初期データを追加（テスト用）
	createTodo("Go言語を学習する")
	createTodo("Docker化を実装する")

	// ルート設定
	http.HandleFunc("/api/todos", todosHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Go Todo App - Server Running!\nTodos: %d", len(getTodos()))
	})

	fmt.Println("Go Todo Server starting on :8080")
	fmt.Printf("Initial todos: %d\n", len(getTodos()))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
