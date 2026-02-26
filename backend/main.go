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

// preMigrateUsername はAutoMigrate前にusernameカラムを安全に追加する。
// 既存行にNULL値があるとNOT NULL制約の追加が失敗するため、
// 先にカラムを作成しEmailから一意のユーザー名を生成して埋める。
func preMigrateUsername(db *gorm.DB) {
	var exists bool
	db.Raw(`SELECT EXISTS (
		SELECT 1 FROM information_schema.columns
		WHERE table_name = 'users' AND column_name = 'username'
	)`).Scan(&exists)
	if exists {
		return
	}

	// NULLを許容してカラムを追加
	db.Exec(`ALTER TABLE users ADD COLUMN username text`)

	// 既存ユーザーにEmailのローカルパート（@より前）をユーザー名として設定
	db.Exec(`UPDATE users SET username = split_part(email, '@', 1) WHERE username IS NULL`)

	// 重複がある場合はIDを末尾に付与して一意にする
	db.Exec(`UPDATE users SET username = username || id::text
		WHERE id IN (
			SELECT u.id FROM users u
			JOIN (SELECT username, MIN(id) AS min_id FROM users GROUP BY username HAVING COUNT(*) > 1) dup
			ON u.username = dup.username AND u.id != dup.min_id
		)`)
}

func main() {
	// .envファイルから環境変数を読み込み（存在しなくてもエラーにしない）
	_ = godotenv.Load()

	cfg := config.Load()

	// PostgreSQLに接続
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// usernameカラムのプレマイグレーション: NOT NULL制約を追加する前に既存行にデフォルト値を設定
	preMigrateUsername(db)

	// 全モデルのAutoMigrationを実行
	if err := db.AutoMigrate(
		&model.User{},
		&model.Follow{},
		&model.GitHubContribution{},
		&model.GitHubLanguageStat{},
		&model.GitHubRepository{},
		&model.Post{},
		&model.Like{},
		&model.Bookmark{},
		&model.Reaction{},
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
		&model.StudyCircle{},
		&model.StudyCircleMember{},
		&model.StudyCircleStep{},
		&model.StudyCircleMemberProgress{},
		&model.StudyCircleCheckin{},
		&model.Mention{},
		&model.YouTubeVideo{},
		&model.YouTubeSearchCache{},
		&model.SpotifyRecentTrack{},
		&model.StreakFreeze{},
		&model.BookmarkCollection{},
		&model.BookmarkCollectionItem{},
		&model.WeeklyChallenge{},
		&model.NoteVersion{},
		&model.ResourceProgress{},
		&model.ProjectMilestone{},
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
