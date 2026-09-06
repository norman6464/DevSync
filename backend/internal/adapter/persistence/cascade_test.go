package persistence

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/stretchr/testify/require"
)

// cascadeTestDB は ATLAS_DB_URL が設定されているときだけ実PostgreSQLへ接続する。
// backend-test.yml では schema.hcl の apply 後にこの環境変数が既に設定されているため、
// CIでは常に実行される。ローカルでは `make db-schema-apply` 済みのDBに対して
// ATLAS_DB_URL=... go test ./internal/adapter/persistence/ -run TestCascade のように実行する。
func cascadeTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("ATLAS_DB_URL")
	if url == "" {
		t.Skip("ATLAS_DB_URL が未設定のためスキップ（実PostgreSQLが必要な結合テスト）")
	}
	pool, err := pgxpool.New(context.Background(), url)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	// notifications.typeはnotification_verbsをFK参照する（DEVSYNC-159）。本番では
	// NewContainerが起動時に同期的にシードするが、結合テストはコンテナを経由しない
	// ため、ここで同じ既知コード一覧を冪等に登録しておく。
	_, err = pool.Exec(context.Background(), `
		INSERT INTO notification_verbs (code)
		VALUES ('post'), ('message'), ('like'), ('comment'), ('follow'), ('answer'), ('badge'), ('level_up'), ('mention')
		ON CONFLICT (code) DO NOTHING
	`)
	require.NoError(t, err)

	return pool
}

// TestDeleteWithRelatedData_NoOrphans は、ユーザー退会（DeleteWithRelatedData）が
// 現在カバーしている範囲のテーブルで、削除後に孤児行が残らないことを検証する結合テスト。
//
// 【重要】この一覧はまだ全テーブルを網羅していない。DeleteWithRelatedData
// （internal/adapter/persistence/user_repository.go）が現時点でカバーしている
// テーブルだけを対象にしており、widget_settings・notification_settings・
// reminder_settings・notes・questions・roadmaps 等、多数のユーザー所有テーブルは
// まだこの関数でも本テストでもカバーされていない（DEVSYNC-156でFKにon_delete CASCADEを
// 宣言し、DEVSYNC-157で手書きの本関数をFK任せへ置き換える過程で、順次このテストへ
// 対象テーブルを足していく）。
//
// 新しいユーザー所有テーブルを追加したときは、このテストへ対応するINSERT/検証を
// 1件足すこと。足し忘れて退会処理の対応も忘れると、孤児行が静かに残り続ける
// （DEVSYNC-114で一度発生した不具合と同じ形）。
func TestDeleteWithRelatedData_NoOrphans(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()
	// NewUserRepository は参照系ポート repository.UserRepository を返すため、
	// DeleteWithRelatedData（repository.AuthUserRepository 側の契約）を呼ぶには
	// 具象型 *userRepository を直接構築する。
	repo := &userRepository{q: sqlcgen.New(pool)}

	// 退会対象ユーザー
	var targetID int64
	err := pool.QueryRow(ctx, `
		INSERT INTO users (username, name, email, created_at, updated_at)
		VALUES ('cascade_target', 'Cascade Target', 'cascade_target@example.com', now(), now())
		RETURNING id
	`).Scan(&targetID)
	require.NoError(t, err)

	// 他ユーザー（退会対象の投稿へコメント・いいね・フォローする側）
	var otherID int64
	err = pool.QueryRow(ctx, `
		INSERT INTO users (username, name, email, created_at, updated_at)
		VALUES ('cascade_other', 'Cascade Other', 'cascade_other@example.com', now(), now())
		RETURNING id
	`).Scan(&otherID)
	require.NoError(t, err)
	t.Cleanup(func() {
		// targetID は削除される想定だが、アサーション失敗時に残っても次回実行を汚さないよう
		// otherID とあわせて後始末する（重複エラーは無視してよい範囲）。
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id IN ($1, $2)`, targetID, otherID)
	})

	// 退会対象が投稿し、他ユーザーがそれにコメント・いいねする
	var postID int64
	err = pool.QueryRow(ctx, `
		INSERT INTO posts (user_id, title, content, created_at, updated_at)
		VALUES ($1, 'title', 'content', now(), now()) RETURNING id
	`, targetID).Scan(&postID)
	require.NoError(t, err)

	var commentID int64
	err = pool.QueryRow(ctx, `
		INSERT INTO comments (user_id, post_id, content, created_at, updated_at)
		VALUES ($1, $2, 'nice post', now(), now()) RETURNING id
	`, otherID, postID).Scan(&commentID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `INSERT INTO likes (user_id, post_id, created_at) VALUES ($1, $2, now())`, otherID, postID)
	require.NoError(t, err)

	// 相互フォロー（follower/followeeの両方向を1本のクエリで消せているか検証するため双方向作る）
	_, err = pool.Exec(ctx, `INSERT INTO follows (follower_id, followee_id, created_at) VALUES ($1, $2, now())`, targetID, otherID)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO follows (follower_id, followee_id, created_at) VALUES ($1, $2, now())`, otherID, targetID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO password_reset_tokens (user_id, token, expires_at, created_at)
		VALUES ($1, 'tok', $2, now())
	`, targetID, time.Now().Add(time.Hour))
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO git_hub_contributions (user_id, contributed_on, count, created_at, updated_at) VALUES ($1, current_date, 3, now(), now())
	`, targetID)
	require.NoError(t, err)

	// --- 退会処理を実行 ---
	err = repo.DeleteWithRelatedData(ctx, uint(targetID))
	require.NoError(t, err, "退会処理自体が失敗した（FK違反等でエラーになっていないか確認）")

	// --- 孤児行が残っていないことを検証 ---
	assertZero(t, ctx, pool, `SELECT count(*) FROM users WHERE id = $1`, targetID)
	assertZero(t, ctx, pool, `SELECT count(*) FROM posts WHERE user_id = $1`, targetID)
	assertZero(t, ctx, pool, `SELECT count(*) FROM comments WHERE id = $1`, commentID)
	assertZero(t, ctx, pool, `SELECT count(*) FROM likes WHERE post_id = $1`, postID)
	assertZero(t, ctx, pool, `SELECT count(*) FROM follows WHERE follower_id = $1 OR followee_id = $1`, targetID)
	assertZero(t, ctx, pool, `SELECT count(*) FROM password_reset_tokens WHERE user_id = $1`, targetID)
	assertZero(t, ctx, pool, `SELECT count(*) FROM git_hub_contributions WHERE user_id = $1`, targetID)

	// 他ユーザー自身は削除されていないこと（巻き添えにしていないか）
	assertZero(t, ctx, pool, `SELECT count(*) FROM users WHERE id = $1 AND username != 'cascade_other'`, otherID)
}

func assertZero(t *testing.T, ctx context.Context, pool *pgxpool.Pool, query string, args ...interface{}) {
	t.Helper()
	var n int
	err := pool.QueryRow(ctx, query, args...).Scan(&n)
	require.NoError(t, err)
	require.Zero(t, n, "孤児行が残っている: %s", query)
}
