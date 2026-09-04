-- name: CountNotesByUser :one
SELECT count(*) FROM notes
WHERE user_id = $1;

-- name: CountArchivedNotesByUser :one
SELECT count(*) FROM notes
WHERE user_id = $1 AND is_archived = true;

-- name: CountFavoriteNotesByUser :one
SELECT count(*) FROM notes
WHERE user_id = $1 AND is_favorite = true;

-- name: CountNoteFoldersByUser :one
SELECT count(*) FROM note_folders
WHERE user_id = $1;

-- name: CountNotesByUserSince :one
SELECT count(*) FROM notes
WHERE user_id = $1 AND created_at >= $2;
