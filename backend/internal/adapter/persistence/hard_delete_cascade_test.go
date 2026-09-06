package persistence

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/require"
)

// TestDeleteQuestion_CascadesToDependents は、論理削除を全廃した質問の物理削除が、
// 回答・回答への投票・質問への投票・ブックマーク・通知までFKのON DELETE CASCADEで
// 正しく連鎖削除することを検証する（DEVSYNC-159）。
func TestDeleteQuestion_CascadesToDependents(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()

	var askerID, answererID, otherID int64
	for _, dst := range []struct {
		id       *int64
		username string
	}{
		{&askerID, "hd_asker"},
		{&answererID, "hd_answerer"},
		{&otherID, "hd_other"},
	} {
		err := pool.QueryRow(ctx, `
			INSERT INTO users (username, name, email, created_at, updated_at)
			VALUES ($1, $1, $1 || '@example.com', now(), now())
			RETURNING id
		`, dst.username).Scan(dst.id)
		require.NoError(t, err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id IN ($1, $2, $3)`, askerID, answererID, otherID)
	})

	questionRepo := NewQuestionRepository(pool)
	answerRepo := NewAnswerRepository(pool)

	question := &model.Question{UserID: uint(askerID), Title: "質問", Body: "本文", Tags: "[]"}
	require.NoError(t, questionRepo.Create(ctx, question))

	answer := &model.Answer{UserID: uint(answererID), QuestionID: question.ID, Body: "回答"}
	require.NoError(t, answerRepo.Create(ctx, answer))

	_, err := pool.Exec(ctx, `INSERT INTO answer_votes (user_id, answer_id, value, created_at) VALUES ($1, $2, 1, now())`, otherID, answer.ID)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO question_votes (user_id, question_id, value, created_at) VALUES ($1, $2, 1, now())`, otherID, question.ID)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO question_bookmarks (user_id, question_id, created_at) VALUES ($1, $2, now())`, otherID, question.ID)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `
		INSERT INTO notifications (user_id, type, actor_id, question_id, read, created_at)
		VALUES ($1, 'answer', $2, $3, false, now())
	`, askerID, answererID, question.ID)
	require.NoError(t, err)

	require.NoError(t, questionRepo.Delete(ctx, question.ID))

	assertZero(t, ctx, pool, `SELECT count(*) FROM questions WHERE id = $1`, question.ID)
	assertZero(t, ctx, pool, `SELECT count(*) FROM answers WHERE id = $1`, answer.ID)
	assertZero(t, ctx, pool, `SELECT count(*) FROM answer_votes WHERE answer_id = $1`, answer.ID)
	assertZero(t, ctx, pool, `SELECT count(*) FROM question_votes WHERE question_id = $1`, question.ID)
	assertZero(t, ctx, pool, `SELECT count(*) FROM question_bookmarks WHERE question_id = $1`, question.ID)
	assertZero(t, ctx, pool, `SELECT count(*) FROM notifications WHERE question_id = $1`, question.ID)
}

// TestDeleteLearningResource_CascadesToDependents は、学習リソースの物理削除が
// いいね・保存・進捗記録までFKのON DELETE CASCADEで連鎖削除することを検証する。
func TestDeleteLearningResource_CascadesToDependents(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()

	var ownerID, otherID int64
	for _, dst := range []struct {
		id       *int64
		username string
	}{
		{&ownerID, "hd_res_owner"},
		{&otherID, "hd_res_other"},
	} {
		err := pool.QueryRow(ctx, `
			INSERT INTO users (username, name, email, created_at, updated_at)
			VALUES ($1, $1, $1 || '@example.com', now(), now())
			RETURNING id
		`, dst.username).Scan(dst.id)
		require.NoError(t, err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id IN ($1, $2)`, ownerID, otherID)
	})

	q := sqlcgen.New(pool)
	resourceRepo := NewLearningResourceRepository(q)

	resource := &model.LearningResource{
		UserID: uint(ownerID), Title: "リソース",
		Category: model.ResourceCategoryBook, Difficulty: model.ResourceDifficultyBeginner, IsPublic: true,
	}
	require.NoError(t, resourceRepo.Create(ctx, resource))

	_, err := pool.Exec(ctx, `INSERT INTO resource_likes (user_id, resource_id, created_at) VALUES ($1, $2, now())`, otherID, resource.ID)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `INSERT INTO resource_saves (user_id, resource_id, created_at) VALUES ($1, $2, now())`, otherID, resource.ID)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `
		INSERT INTO resource_progresses (user_id, resource_id, status, created_at, updated_at)
		VALUES ($1, $2, 'in_progress', now(), now())
	`, otherID, resource.ID)
	require.NoError(t, err)

	require.NoError(t, resourceRepo.Delete(ctx, resource.ID))

	assertZero(t, ctx, pool, `SELECT count(*) FROM learning_resources WHERE id = $1`, resource.ID)
	assertZero(t, ctx, pool, `SELECT count(*) FROM resource_likes WHERE resource_id = $1`, resource.ID)
	assertZero(t, ctx, pool, `SELECT count(*) FROM resource_saves WHERE resource_id = $1`, resource.ID)
	assertZero(t, ctx, pool, `SELECT count(*) FROM resource_progresses WHERE resource_id = $1`, resource.ID)
}
