-- name: CreateNoteLink :exec
INSERT INTO note_links (source_note_id, target_note_id, created_at)
VALUES ($1, $2, now());

-- name: ListNoteLinksBySource :many
-- GORMのPreload("TargetNote")に相当。リンク先ノートをsqlc.embedで一緒に取得する。
SELECT sqlc.embed(note_links), sqlc.embed(notes)
FROM note_links
JOIN notes ON notes.id = note_links.target_note_id
WHERE note_links.source_note_id = $1
ORDER BY note_links.created_at DESC;

-- name: ListNoteLinksByTarget :many
-- GORMのPreload("SourceNote")に相当。リンク元ノートをsqlc.embedで一緒に取得する。
SELECT sqlc.embed(note_links), sqlc.embed(notes)
FROM note_links
JOIN notes ON notes.id = note_links.source_note_id
WHERE note_links.target_note_id = $1
ORDER BY note_links.created_at DESC;

-- name: DeleteNoteLink :exec
DELETE FROM note_links
WHERE source_note_id = $1 AND target_note_id = $2;

-- name: CountNoteLinksBetween :one
SELECT count(*) FROM note_links
WHERE source_note_id = $1 AND target_note_id = $2;

-- name: CountNoteLinksBySource :one
SELECT count(*) FROM note_links
WHERE source_note_id = $1;

-- name: CountNoteLinksByTarget :one
SELECT count(*) FROM note_links
WHERE target_note_id = $1;
