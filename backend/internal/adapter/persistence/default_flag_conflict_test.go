package persistence

import (
	"context"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNoteTemplateCreate_SwitchDefault_Succeeds は、既にデフォルトが設定されている
// ユーザーが「別のテンプレートへデフォルトを切り替える」という最も普通の操作が
// 成功することを検証する回帰テスト。デフォルト解除とCreateを1つのSQL文（データ変更CTE）で
// 行っていたときは、PostgreSQLのCTEが主文と同一スナップショットで動くために、この
// 普通の切り替え操作ですら毎回一意制約違反になっていた（2文に分けて修正）。
func TestNoteTemplateCreate_SwitchDefault_Succeeds(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()

	var userID int64
	err := pool.QueryRow(ctx, `
		INSERT INTO users (username, name, email, created_at, updated_at)
		VALUES ('switch_default_user', 'Switch Default', 'switch_default@example.com', now(), now())
		RETURNING id
	`).Scan(&userID)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	})

	_, err = pool.Exec(ctx, `
		INSERT INTO note_templates (user_id, name, content_template, is_default, created_at, updated_at)
		VALUES ($1, 'A', '本文', true, now(), now())
	`, userID)
	require.NoError(t, err)

	repo := &noteTemplateRepository{pool: pool, q: sqlcgen.New(pool)}
	templateB := &model.NoteTemplate{UserID: uint(userID), Name: "B", ContentTemplate: "本文", IsDefault: true}
	require.NoError(t, repo.Create(ctx, templateB))

	var defaultCount int
	err = pool.QueryRow(ctx, `SELECT count(*) FROM note_templates WHERE user_id = $1 AND is_default`, userID).Scan(&defaultCount)
	require.NoError(t, err)
	assert.Equal(t, 1, defaultCount, "デフォルトはユーザーごとに常に高々1件のはず")
	assert.True(t, templateB.IsDefault)
}

// TestNoteTemplateDefaultUniqueIndex_ConcurrentTransactions_Conflicts は、同一ユーザーの
// 2つのトランザクションが本当に並行して別々のテンプレートをデフォルトへ指定したとき、
// 部分UNIQUE索引（uq_note_templates_default）がそれでも「ユーザーごと高々1件」を守り、
// 敗者側の生のエラーがisUniqueViolationで検出できる（＝noteTemplateRepository.Create/Updateが
// domain.ErrConflictへ変換できる）ことを検証する。
//
// noteTemplateRepository.Create/Updateは「解除→作成」を1つのトランザクション内の2文として
// 実行するため、単一トランザクション内では常に安全（TestNoteTemplateCreate_SwitchDefault_Succeeds
// で検証済み）。しかし2つの独立したトランザクションが「片方のCOMMIT前にもう片方の解除UPDATEが
// 同じ行を見に行く」タイミングで重なると、後発側は「解除対象は既に無くなっていた
// （＝相手がまだ入れていない）」と誤認したまま自分の行だけをtrueにしてしまい、
// 結果として2行が同時にtrueになりかけて一意制約違反が起こる。Createはトランザクションを
// 自前で開始・コミットするため、外側から一時停止できない。そこでここではCreateの内部と
// 同じ2文（ClearOtherNoteTemplateDefaults→CreateNoteTemplate）を明示的な2本の
// トランザクションで直接実行し、そのタイミングを再現する。
func TestNoteTemplateDefaultUniqueIndex_ConcurrentTransactions_Conflicts(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()

	var userID int64
	err := pool.QueryRow(ctx, `
		INSERT INTO users (username, name, email, created_at, updated_at)
		VALUES ('conflict_user', 'Conflict User', 'conflict_user@example.com', now(), now())
		RETURNING id
	`).Scan(&userID)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	})

	// 既存のデフォルトテンプレートA（コミット済み）。
	_, err = pool.Exec(ctx, `
		INSERT INTO note_templates (user_id, name, content_template, is_default, created_at, updated_at)
		VALUES ($1, 'A', '本文', true, now(), now())
	`, userID)
	require.NoError(t, err)

	isDefault := true
	tx1, err := pool.Begin(ctx)
	require.NoError(t, err)
	defer tx1.Rollback(ctx) //nolint:errcheck
	q1 := sqlcgen.New(tx1)

	require.NoError(t, q1.ClearOtherNoteTemplateDefaults(ctx, sqlcgen.ClearOtherNoteTemplateDefaultsParams{
		UserID: userID, ID: 0,
	}))
	_, err = q1.CreateNoteTemplate(ctx, sqlcgen.CreateNoteTemplateParams{
		UserID: userID, Name: "B", ContentTemplate: "本文", IsDefault: isDefault,
	})
	require.NoError(t, err, "Tx1のデフォルト作成自体は成功するはず")

	// Tx2はTx1のCOMMIT前に開始する。AをクリアするUPDATEでTx1と行ロックが競合し、
	// Tx1のCOMMITまでブロックされるはず。
	tx2, err := pool.Begin(ctx)
	require.NoError(t, err)
	defer tx2.Rollback(ctx) //nolint:errcheck
	q2 := sqlcgen.New(tx2)

	var tx2Err error
	tx2Started := make(chan struct{})
	tx2Done := make(chan struct{})
	go func() {
		close(tx2Started)
		if cerr := q2.ClearOtherNoteTemplateDefaults(ctx, sqlcgen.ClearOtherNoteTemplateDefaultsParams{
			UserID: userID, ID: 0,
		}); cerr != nil {
			tx2Err = cerr
			close(tx2Done)
			return
		}
		_, tx2Err = q2.CreateNoteTemplate(ctx, sqlcgen.CreateNoteTemplateParams{
			UserID: userID, Name: "C", ContentTemplate: "本文", IsDefault: isDefault,
		})
		close(tx2Done)
	}()
	<-tx2Started
	// Tx2がAのロック待ちに入るのを確実にするための猶予。
	time.Sleep(300 * time.Millisecond)

	require.NoError(t, tx1.Commit(ctx))
	<-tx2Done

	require.Error(t, tx2Err, "後発トランザクションは一意制約違反になるはず")
	assert.True(t, isUniqueViolation(tx2Err), "noteTemplateRepository.Create/Updateがdomain.ErrConflictへ変換できる形のエラーのはず")
}

// TestLearningLogTemplateCreate_SwitchDefault_Succeeds は note_templates と同じ検証を
// learning_log_templates に対して行う（同じCTEの不具合を同じ2文構成で修正しているため）。
func TestLearningLogTemplateCreate_SwitchDefault_Succeeds(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()

	var userID int64
	err := pool.QueryRow(ctx, `
		INSERT INTO users (username, name, email, created_at, updated_at)
		VALUES ('switch_default_log_user', 'Switch Default Log', 'switch_default_log@example.com', now(), now())
		RETURNING id
	`).Scan(&userID)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	})

	_, err = pool.Exec(ctx, `
		INSERT INTO learning_log_templates (user_id, name, default_category, default_duration, is_default, created_at, updated_at)
		VALUES ($1, 'A', 'coding', 30, true, now(), now())
	`, userID)
	require.NoError(t, err)

	repo := &learningLogTemplateRepository{pool: pool, q: sqlcgen.New(pool)}
	templateB := &model.LearningLogTemplate{
		UserID: uint(userID), Name: "B", DefaultCategory: model.LogCategoryCoding, DefaultDuration: 30, IsDefault: true,
	}
	require.NoError(t, repo.Create(ctx, templateB))

	var defaultCount int
	err = pool.QueryRow(ctx, `SELECT count(*) FROM learning_log_templates WHERE user_id = $1 AND is_default`, userID).Scan(&defaultCount)
	require.NoError(t, err)
	assert.Equal(t, 1, defaultCount, "デフォルトはユーザーごとに常に高々1件のはず")
	assert.True(t, templateB.IsDefault)
}

// TestLearningLogTemplateDefaultUniqueIndex_ConcurrentTransactions_Conflicts は
// note_templates と同じ検証を learning_log_templates に対して行う。
func TestLearningLogTemplateDefaultUniqueIndex_ConcurrentTransactions_Conflicts(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()

	var userID int64
	err := pool.QueryRow(ctx, `
		INSERT INTO users (username, name, email, created_at, updated_at)
		VALUES ('conflict_log_user', 'Conflict Log User', 'conflict_log_user@example.com', now(), now())
		RETURNING id
	`).Scan(&userID)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	})

	_, err = pool.Exec(ctx, `
		INSERT INTO learning_log_templates (user_id, name, default_category, default_duration, is_default, created_at, updated_at)
		VALUES ($1, 'A', 'coding', 30, true, now(), now())
	`, userID)
	require.NoError(t, err)

	isDefault := true
	category := "coding"
	duration := int64(30)
	tx1, err := pool.Begin(ctx)
	require.NoError(t, err)
	defer tx1.Rollback(ctx) //nolint:errcheck
	q1 := sqlcgen.New(tx1)

	require.NoError(t, q1.ClearOtherLearningLogTemplateDefaults(ctx, sqlcgen.ClearOtherLearningLogTemplateDefaultsParams{
		UserID: userID, ID: 0,
	}))
	_, err = q1.CreateLearningLogTemplate(ctx, sqlcgen.CreateLearningLogTemplateParams{
		UserID: userID, Name: "B", DefaultCategory: &category, DefaultDuration: &duration, IsDefault: isDefault,
	})
	require.NoError(t, err, "Tx1のデフォルト作成自体は成功するはず")

	tx2, err := pool.Begin(ctx)
	require.NoError(t, err)
	defer tx2.Rollback(ctx) //nolint:errcheck
	q2 := sqlcgen.New(tx2)

	var tx2Err error
	tx2Started := make(chan struct{})
	tx2Done := make(chan struct{})
	go func() {
		close(tx2Started)
		if cerr := q2.ClearOtherLearningLogTemplateDefaults(ctx, sqlcgen.ClearOtherLearningLogTemplateDefaultsParams{
			UserID: userID, ID: 0,
		}); cerr != nil {
			tx2Err = cerr
			close(tx2Done)
			return
		}
		_, tx2Err = q2.CreateLearningLogTemplate(ctx, sqlcgen.CreateLearningLogTemplateParams{
			UserID: userID, Name: "C", DefaultCategory: &category, DefaultDuration: &duration, IsDefault: isDefault,
		})
		close(tx2Done)
	}()
	<-tx2Started
	time.Sleep(300 * time.Millisecond)

	require.NoError(t, tx1.Commit(ctx))
	<-tx2Done

	require.Error(t, tx2Err, "後発トランザクションは一意制約違反になるはず")
	assert.True(t, isUniqueViolation(tx2Err), "learningLogTemplateRepository.Create/Updateがdomain.ErrConflictへ変換できる形のエラーのはず")
}
