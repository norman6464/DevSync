-- name: CreateNote :one
INSERT INTO notes (user_id, folder_id, title, content, tags, is_favorite, is_archived, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())
RETURNING *;

-- name: UpdateNote :one
-- GORMのSave（全カラム上書き）に相当。
UPDATE notes
SET title = $2, content = $3, tags = $4, folder_id = $5, is_favorite = $6, is_archived = $7, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteNote :exec
DELETE FROM notes WHERE id = $1;

-- name: GetNoteByID :one
-- GORMのPreload("User").Preload("Folder")に相当。
-- folder_idはNULL許容のためLEFT JOIN。note_foldersはsqlc.embedを使わず個別カラム選択にすることで、
-- sqlcのJOIN文脈依存のnull推論（LEFT JOIN側は全カラムがnullableになる）を効かせる。
-- sqlc.embedは対象テーブル自身のスキーマ上のnull許容性をそのまま使うため、LEFT JOINで欠落しうる
-- 行に対しては非NULL型（例: NoteFolder.ID int64）のままとなりスキャン時にエラーとなる。
SELECT sqlc.embed(notes), sqlc.embed(users),
    note_folders.id AS folder_id_2,
    note_folders.user_id AS folder_user_id,
    note_folders.parent_id AS folder_parent_id,
    note_folders.name AS folder_name,
    note_folders.created_at AS folder_created_at,
    note_folders.updated_at AS folder_updated_at
FROM notes
JOIN users ON users.id = notes.user_id
LEFT JOIN note_folders ON note_folders.id = notes.folder_id
WHERE notes.id = $1;

-- name: ListNotesByUser :many
-- GORMのPreload("Folder")に相当。folder_idはNULL許容のためLEFT JOIN（理由は GetNoteByID と同じ）。
SELECT sqlc.embed(notes),
    note_folders.id AS folder_id_2,
    note_folders.user_id AS folder_user_id,
    note_folders.parent_id AS folder_parent_id,
    note_folders.name AS folder_name,
    note_folders.created_at AS folder_created_at,
    note_folders.updated_at AS folder_updated_at
FROM notes
LEFT JOIN note_folders ON note_folders.id = notes.folder_id
WHERE notes.user_id = $1 AND notes.is_archived = false
ORDER BY notes.updated_at DESC
LIMIT $2 OFFSET $3;

-- name: ListFavoriteNotesByUser :many
SELECT sqlc.embed(notes),
    note_folders.id AS folder_id_2,
    note_folders.user_id AS folder_user_id,
    note_folders.parent_id AS folder_parent_id,
    note_folders.name AS folder_name,
    note_folders.created_at AS folder_created_at,
    note_folders.updated_at AS folder_updated_at
FROM notes
LEFT JOIN note_folders ON note_folders.id = notes.folder_id
WHERE notes.user_id = $1 AND notes.is_favorite = true
ORDER BY notes.updated_at DESC
LIMIT $2 OFFSET $3;

-- name: ListArchivedNotesByUser :many
SELECT sqlc.embed(notes),
    note_folders.id AS folder_id_2,
    note_folders.user_id AS folder_user_id,
    note_folders.parent_id AS folder_parent_id,
    note_folders.name AS folder_name,
    note_folders.created_at AS folder_created_at,
    note_folders.updated_at AS folder_updated_at
FROM notes
LEFT JOIN note_folders ON note_folders.id = notes.folder_id
WHERE notes.user_id = $1 AND notes.is_archived = true
ORDER BY notes.updated_at DESC
LIMIT $2 OFFSET $3;

-- name: ListNotesByFolder :many
SELECT * FROM notes
WHERE folder_id = $1 AND user_id = $2
ORDER BY updated_at DESC;

-- name: SearchNotes :many
SELECT sqlc.embed(notes),
    note_folders.id AS folder_id_2,
    note_folders.user_id AS folder_user_id,
    note_folders.parent_id AS folder_parent_id,
    note_folders.name AS folder_name,
    note_folders.created_at AS folder_created_at,
    note_folders.updated_at AS folder_updated_at
FROM notes
LEFT JOIN note_folders ON note_folders.id = notes.folder_id
WHERE notes.user_id = $1 AND notes.is_archived = false
    AND (notes.title LIKE $2 OR notes.content LIKE $2)
ORDER BY notes.updated_at DESC
LIMIT $3 OFFSET $4;

-- name: CountSearchNotes :one
SELECT COUNT(*) FROM notes
WHERE user_id = $1 AND is_archived = false
    AND (title LIKE $2 OR content LIKE $2);

-- name: CountActiveNotesByUser :one
-- note_stats.sql の CountNotesByUser はアーカイブ済みも含めた全件数のため名前を分けている。
SELECT COUNT(*) FROM notes WHERE user_id = $1 AND is_archived = false;

-- name: ToggleNoteFavorite :exec
UPDATE notes SET is_favorite = NOT is_favorite, updated_at = now() WHERE id = $1;

-- name: ArchiveNote :exec
UPDATE notes SET is_archived = true, updated_at = now() WHERE id = $1;

-- name: UnarchiveNote :exec
UPDATE notes SET is_archived = false, updated_at = now() WHERE id = $1;
