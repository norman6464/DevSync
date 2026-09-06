-- name: ClearOtherNoteTemplateDefaults :exec
-- 同一ユーザーの既存デフォルト（自分自身のid=$2は除く）を外す。Createからは
-- 存在しないid（0）を渡すことで全件が対象になる。呼び出し元（repository）が
-- 1つのトランザクション内でこのクエリの直後にCreate/Updateを実行することで、
-- 「解除」と「新しいデフォルトの確定」を1つの原子的な単位にする。
-- 【重要】1つのSQL文の中のCTEとしてUPDATE→INSERT/UPDATEを両方書く方式は使わない。
-- PostgreSQLのデータ変更CTEは主文と同一スナップショットで動作し、UPDATEによる解除の
-- 効果が同一文内の一意制約チェックに間に合わず、デフォルトが1件も無い状態からの
-- 切り替えでも毎回一意制約違反になることを確認済み（2文に分けて初めて正しく動く）。
-- uq_note_templates_default（(user_id) WHERE is_default の部分UNIQUE索引）は
-- 本当に並行する別トランザクション同士の衝突に対する最終防衛。
UPDATE note_templates SET is_default = false, updated_at = now()
WHERE note_templates.user_id = $1 AND note_templates.is_default IS TRUE AND note_templates.id != $2;

-- name: CreateNoteTemplate :one
INSERT INTO note_templates (user_id, name, description, default_title, content_template, default_tags, is_default, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())
RETURNING *;

-- name: UpdateNoteTemplate :one
UPDATE note_templates
SET name = $2,
    description = $3,
    default_title = $4,
    content_template = $5,
    default_tags = $6,
    is_default = $7,
    updated_at = now()
WHERE note_templates.id = $1
RETURNING *;

-- name: DeleteNoteTemplate :exec
DELETE FROM note_templates
WHERE id = $1;

-- name: GetNoteTemplateByID :one
SELECT * FROM note_templates
WHERE id = $1;

-- name: ListNoteTemplatesByUser :many
SELECT * FROM note_templates
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetDefaultNoteTemplateByUser :one
SELECT * FROM note_templates
WHERE user_id = $1 AND is_default = true
ORDER BY id ASC
LIMIT 1;

-- name: CountNoteTemplatesByUser :one
SELECT count(*) FROM note_templates
WHERE user_id = $1;
