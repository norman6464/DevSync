// Package main はDevSyncバックエンドアプリケーションのエントリーポイント。
// DB接続・マイグレーション・WebSocket Hub起動・HTTPサーバー起動を行う。
package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/norman6464/devsync/backend/internal/config"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/router"
	"github.com/norman6464/devsync/backend/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// .envファイルから環境変数を読み込み（存在しなくてもエラーにしない）
	_ = godotenv.Load()

	cfg := config.Load()

	// PostgreSQLに接続
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// 全モデルのAutoMigrationを実行
	if err := db.AutoMigrate(
		&model.User{},
		&model.Follow{},
		&model.GitHubContribution{},
		&model.GitHubLanguageStat{},
		&model.GitHubRepository{},
		&model.Post{},
		&model.Like{},
		&model.Comment{},
		&model.Message{},
		&model.Notification{},
		&model.PasswordResetToken{},
		&model.ZennArticle{},
		&model.QiitaArticle{},
		&model.LearningGoal{},
		&model.Project{},
		&model.LearningResource{},
		&model.ResourceLike{},
		&model.ResourceSave{},
		&model.BookReview{},
		&model.Question{},
		&model.QuestionVote{},
		&model.Answer{},
		&model.AnswerVote{},
		&model.Roadmap{},
		&model.RoadmapStep{},
		&model.ChatRoom{},
		&model.ChatRoomMember{},
		&model.GroupMessage{},
		&model.LearningLog{},
		&model.CodeSnippet{},
		&model.SnippetComment{},
		&model.AIAdvice{},
		&model.AIConversation{},
		&model.AIMessage{},
	); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// 既存ユーザーのオンボーディング完了フラグを初期化
	db.Model(&model.User{}).Where("onboarding_completed = ?", false).Update("onboarding_completed", true)

	// WebSocket Hubをバックグラウンドで起動
	hub := service.NewHub()
	go hub.Run()

	// ルーターを構築しサーバーを起動
	r := router.Setup(db, cfg, hub)

	log.Printf("Server starting on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
