-- name: ClearOtherLearningLogTemplateDefaults :exec
-- 同一ユーザーの既存デフォルト（自分自身のid=$2は除く）を外す。Createからは
-- 存在しないid（0）を渡すことで全件が対象になる。呼び出し元（repository）が
-- 1つのトランザクション内でこのクエリの直後にCreate/Updateを実行することで、
-- 「解除」と「新しいデフォルトの確定」を1つの原子的な単位にする。
-- 【重要】1つのSQL文の中のCTEとしてUPDATE→INSERT/UPDATEを両方書く方式は使わない。
-- PostgreSQLのデータ変更CTEは主文と同一スナップショットで動作し、UPDATEによる解除の
-- 効果が同一文内の一意制約チェックに間に合わず、デフォルトが1件も無い状態からの
-- 切り替えでも毎回一意制約違反になることを確認済み（2文に分けて初めて正しく動く）。
-- uq_learning_log_templates_default（(user_id) WHERE is_default の部分UNIQUE索引）は
-- 本当に並行する別トランザクション同士の衝突に対する最終防衛。
UPDATE learning_log_templates SET is_default = false, updated_at = now()
WHERE learning_log_templates.user_id = $1 AND learning_log_templates.is_default IS TRUE AND learning_log_templates.id != $2;

-- name: CreateLearningLogTemplate :one
INSERT INTO learning_log_templates (user_id, name, default_title, default_content, default_category, default_duration, is_default, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())
RETURNING *;

-- name: UpdateLearningLogTemplate :one
-- GORMのSave（全カラム上書き）に相当。
UPDATE learning_log_templates SET
    name = $2,
    default_title = $3,
    default_content = $4,
    default_category = $5,
    default_duration = $6,
    is_default = $7,
    updated_at = now()
WHERE learning_log_templates.id = $1
RETURNING *;

-- name: DeleteLearningLogTemplate :exec
DELETE FROM learning_log_templates WHERE id = $1;

-- name: GetLearningLogTemplateByID :one
SELECT * FROM learning_log_templates WHERE id = $1;

-- name: ListLearningLogTemplatesByUser :many
SELECT * FROM learning_log_templates WHERE user_id = $1 ORDER BY created_at DESC;

-- name: GetDefaultLearningLogTemplateByUser :one
SELECT * FROM learning_log_templates WHERE user_id = $1 AND is_default = true;

-- name: CountLearningLogTemplatesByUser :one
SELECT COUNT(*) FROM learning_log_templates WHERE user_id = $1;
