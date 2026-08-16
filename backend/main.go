// Package main はDevSyncバックエンドアプリケーションのエントリーポイント。
// DB接続・マイグレーション・WebSocket Hub起動・HTTPサーバー起動を行う。
package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence"
	"github.com/norman6464/devsync/backend/internal/config"
	"github.com/norman6464/devsync/backend/internal/infra/dbschema"
	"github.com/norman6464/devsync/backend/internal/infra/ws"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/router"
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

// preMigrateDedupeForUniqueIndexes は AutoMigrate が複合ユニーク索引を作る前に既存の重複行を除去する。
// 重複が残っていると索引の作成に失敗して起動できないため。最小 id の行を正とし、それ以外を削除する。
func preMigrateDedupeForUniqueIndexes(db *gorm.DB) {
	targets := []struct{ table, dedupe string }{
		// idx_bookmark_collection_post (collection_id, post_id)
		{"bookmark_collection_items", `DELETE FROM bookmark_collection_items a USING bookmark_collection_items b
			WHERE a.collection_id = b.collection_id AND a.post_id = b.post_id AND a.id > b.id`},
		// idx_streak_freeze_user_date (user_id, used_date)
		{"streak_freezes", `DELETE FROM streak_freezes a USING streak_freezes b
			WHERE a.user_id = b.user_id AND a.used_date = b.used_date AND a.id > b.id`},
	}
	for _, t := range targets {
		var exists bool
		db.Raw(`SELECT to_regclass(?) IS NOT NULL`, t.table).Scan(&exists)
		if exists {
			db.Exec(t.dedupe)
		}
	}
}

// fixRoadmapStepCounts はテンプレート初期登録の二重加算で誤った step_count を
// roadmap_steps の実数へ補正し、補正した行だけ progress と自動完了を再計算する。
// 正しい行には触れない冪等な補正のため、毎起動で実行してよい。
func fixRoadmapStepCounts(db *gorm.DB) {
	// PostgreSQL は 1 文内で同一行を 2 回更新できないため、進捗と自動完了も同じ UPDATE で行う。
	db.Exec(`WITH actual AS (
		SELECT roadmap_id, count(*) AS cnt FROM roadmap_steps GROUP BY roadmap_id
	)
	UPDATE roadmaps r
	SET step_count = a.cnt,
	    progress = CASE WHEN a.cnt > 0 THEN LEAST(r.completed_step_count * 100 / a.cnt, 100) ELSE 0 END,
	    status = CASE WHEN a.cnt > 0 AND r.completed_step_count >= a.cnt AND r.status = 'active'
	                  THEN 'completed' ELSE r.status END,
	    completed_at = CASE WHEN a.cnt > 0 AND r.completed_step_count >= a.cnt AND r.status = 'active'
	                        THEN COALESCE(r.completed_at, now()) ELSE r.completed_at END
	FROM actual a
	WHERE r.id = a.roadmap_id AND r.step_count <> a.cnt`)
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

	// 複合ユニーク索引の作成前に既存の重複行を除去（残っていると索引作成が失敗するため）
	preMigrateDedupeForUniqueIndexes(db)

	// 全モデルのAutoMigrationを実行
	if err := db.AutoMigrate(
		&model.User{},
		&model.Follow{},
		&model.GitHubContribution{},
		&model.GitHubLanguageStat{},
		&model.GitHubRepository{},
		&model.Post{},
		&model.PostTag{},
		&model.PostView{},
		&model.PostPin{},
		&model.PostTemplate{},
		&model.PostCollection{},
		&model.PostCollectionItem{},
		&model.PostSeries{},
		&model.PostSeriesItem{},
		&model.Like{},
		&model.Bookmark{},
		&model.Reaction{},
		&model.Comment{},
		&model.CommentLike{},
		&model.Message{},
		&model.Notification{},
		&model.NotificationSettings{},
		&model.ReminderSettings{},
		&model.PasswordResetToken{},
		&model.ZennArticle{},
		&model.QiitaArticle{},
		&model.LearningGoal{},
		&model.Project{},
		&model.LearningResource{},
		&model.ResourceLike{},
		&model.ResourceReview{},
		&model.ResourceSave{},
		&model.BookReview{},
		&model.Question{},
		&model.QuestionVote{},
		&model.QuestionBookmark{},
		&model.Answer{},
		&model.AnswerVote{},
		&model.Roadmap{},
		&model.RoadmapStep{},
		&model.ChatRoom{},
		&model.ChatRoomMember{},
		&model.GroupMessage{},
		&model.LearningLog{},
		&model.LearningLogTemplate{},
		&model.CodeSnippet{},
		&model.CodeSnippetFavorite{},
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
		&model.NoteFolder{},
		&model.Note{},
		&model.NoteVersion{},
		&model.NoteLink{},
		&model.NoteTemplate{},
		&model.ResourceProgress{},
		&model.ProjectMilestone{},
		&model.UserActivity{},
		&model.WeeklyGoal{},
		&model.WidgetSettings{},
	); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// GORM のタグで表現できないインデックス（git_hub_id の部分ユニーク）を補正する
	if err := dbschema.EnsureUserIndexes(db); err != nil {
		log.Fatalf("failed to ensure user indexes: %v", err)
	}
	if err := dbschema.EnsureMentionIndexes(db); err != nil {
		log.Fatalf("failed to ensure mention indexes: %v", err)
	}

	// テンプレート初期登録の二重加算で誤った step_count と進捗を補正する
	fixRoadmapStepCounts(db)

	// 既存ユーザーのオンボーディング完了フラグを初期化
	db.Model(&model.User{}).Where("onboarding_completed = ?", false).Update("onboarding_completed", true)

	// WebSocket Hubをバックグラウンドで起動
	hub := ws.NewHub(persistence.NewRoomMemberLookup(db))
	go hub.Run()

	// ルーターを構築しサーバーを起動
	r := router.Setup(db, cfg, hub)

	log.Printf("Server starting on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
