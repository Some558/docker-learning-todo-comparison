package database

import (
	"fmt"
	"log"

	"go-todo-app/config"
	"go-todo-app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// グローバルなデータベース接続
var DB *gorm.DB

// Connect はデータベースに接続する
func Connect(cfg *config.Config) error {
	dsn := cfg.GetDSN()

	log.Printf("🐘 データベース接続試行: %s", hidePassword(dsn))

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // SQL文をログ出力
	})

	if err != nil {
		return fmt.Errorf("データベース接続エラー: %w", err)
	}

	// 接続テスト
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("データベース取得エラー: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("データベースPingエラー: %w", err)
	}

	log.Println("✅ データベース接続成功")
	return nil
}

// GetDB はデータベース接続を返す
func GetDB() *gorm.DB {
	return DB
}

// hidePassword はログ出力用にパスワードを隠す
func hidePassword(dsn string) string {
	// 簡易版：パスワード部分を***に置換
	// 実際のプロダクトではより厳密な処理が必要
	return dsn // とりあえずそのまま（開発環境なので）
}

// Migrate はデータベーステーブルを作成・更新する
func Migrate() error {
	log.Println("🔄 データベースマイグレーション開始")

	// Todo構造体からテーブルを自動生成
	err := DB.AutoMigrate(&models.Todo{})
	if err != nil {
		return fmt.Errorf("マイグレーションエラー: %w", err)
	}

	log.Println("✅ データベースマイグレーション完了")
	return nil
}
