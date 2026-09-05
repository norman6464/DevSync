-- name: CreateBookReview :one
INSERT INTO book_reviews (
    user_id, title, author, isbn, rating, review,
    total_pages, current_page, image_url, status, is_archived, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, now(), now()
) RETURNING *;

-- name: GetBookReviewWithUserByID :one
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
-- book_reviewsは論理削除があるため、削除済みは除外する（GORM Firstの自動スコープ相当）。
SELECT sqlc.embed(book_reviews), sqlc.embed(users)
FROM book_reviews
JOIN users ON users.id = book_reviews.user_id
WHERE book_reviews.id = $1 AND book_reviews.deleted_at IS NULL;

-- name: ListBookReviewsByUser :many
SELECT * FROM book_reviews
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountBookReviewsByUser :one
SELECT COUNT(*) FROM book_reviews WHERE user_id = $1 AND deleted_at IS NULL;

-- name: ListAllBookReviewsWithUser :many
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(book_reviews), sqlc.embed(users)
FROM book_reviews
JOIN users ON users.id = book_reviews.user_id
WHERE book_reviews.deleted_at IS NULL
ORDER BY book_reviews.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountAllBookReviews :one
SELECT COUNT(*) FROM book_reviews WHERE deleted_at IS NULL;

-- name: ListBookReviewsByRating :many
SELECT * FROM book_reviews
WHERE user_id = $1 AND rating >= $2 AND rating <= $3 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: SearchBookReviews :many
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(book_reviews), sqlc.embed(users)
FROM book_reviews
JOIN users ON users.id = book_reviews.user_id
WHERE (book_reviews.title ILIKE $1 OR book_reviews.author ILIKE $1 OR book_reviews.isbn ILIKE $1)
    AND book_reviews.deleted_at IS NULL
ORDER BY book_reviews.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountSearchBookReviews :one
SELECT COUNT(*) FROM book_reviews
WHERE (title ILIKE $1 OR author ILIKE $1 OR isbn ILIKE $1) AND deleted_at IS NULL;

-- name: UpdateBookReview :one
-- GORMのSave（全カラム上書き）に相当。
UPDATE book_reviews SET
    title = $2, author = $3, isbn = $4, rating = $5, review = $6,
    total_pages = $7, current_page = $8, image_url = $9, status = $10, is_archived = $11,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteBookReview :exec
-- GORMのDelete（論理削除）に相当。
UPDATE book_reviews SET deleted_at = now() WHERE id = $1;
