// Package main はDevSyncバックエンドアプリケーションのエントリーポイント。
// DB接続・WebSocket Hub起動・HTTPサーバー起動を行う。
// スキーマ管理は Atlas（internal/infra/database/schema/*.hcl の宣言的適用）に委譲しており、ここでは行わない。
package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/handler"
	"github.com/norman6464/devsync/backend/internal/infra/config"
	"github.com/norman6464/devsync/backend/internal/infra/ws"
)

// fixRoadmapStepCounts はテンプレート初期登録の二重加算で誤った step_count を
// roadmap_steps の実数へ補正し、補正した行だけ progress と自動完了を再計算する。
// 正しい行には触れない冪等な補正のため、毎起動で実行してよい。
func fixRoadmapStepCounts(ctx context.Context, pool *pgxpool.Pool) {
	// PostgreSQL は 1 文内で同一行を 2 回更新できないため、進捗と自動完了も同じ UPDATE で行う。
	if _, err := pool.Exec(ctx, `WITH actual AS (
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
	WHERE r.id = a.roadmap_id AND r.step_count <> a.cnt`); err != nil {
		log.Printf("roadmap step_count 補正失敗: %v", err)
	}
}

func main() {
	// .envファイルから環境変数を読み込み（存在しなくてもエラーにしない）
	_ = godotenv.Load()

	cfg := config.Load()
	ctx := context.Background()

	// PostgreSQLに接続。スキーマ（テーブル・索引・外部キー）は Atlas の宣言的管理に一本化しており、
	// internal/infra/database/schema/*.hcl が正。デプロイ/ローカル起動前に `make -C backend db-schema-apply` を実行する。
	sqlPool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database (pgx): %v", err)
	}
	defer sqlPool.Close()

	// テンプレート初期登録の二重加算で誤った step_count と進捗を補正する
	fixRoadmapStepCounts(ctx, sqlPool)

	// 既存ユーザーのオンボーディング完了フラグを初期化
	if _, err := sqlPool.Exec(ctx, `UPDATE users SET onboarding_completed = true WHERE onboarding_completed = false`); err != nil {
		log.Printf("オンボーディング完了フラグ初期化失敗: %v", err)
	}

	// WebSocket Hubをバックグラウンドで起動
	hub := ws.NewHub(persistence.NewRoomMemberLookup(sqlcgen.New(sqlPool)))
	go hub.Run()

	// ルーターを構築しサーバーを起動
	r := handler.NewRouter(sqlPool, cfg, hub)

	log.Printf("Server starting on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
