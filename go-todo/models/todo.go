package models

import (
	"errors"
	"time"
)

// Todo はタスクを表現する構造体
type Todo struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null;size:200"`
	Completed bool      `json:"completed" gorm:"default:false"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName はテーブル名を明示的に指定
func (Todo) TableName() string {
	return "todos"
}

// Validate はTodoのバリデーションを行う
func (t *Todo) Validate() error {
	if t.Title == "" {
		return errors.New("title is required")
	}
	if len(t.Title) > 200 {
		return errors.New("title is too long (max 200 characters)")
	}
	return nil
}
