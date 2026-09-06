-- name: GetBookReviewStats :one
-- レビュー0件でもCOALESCEにより全項目0を返す
-- （GORM実装のtotal_reviews==0での早期returnと同じ結果になる）。
SELECT
  count(*) AS total_reviews,
  COALESCE(AVG(rating), 0)::float8 AS average_rating,
  COALESCE(MAX(rating), 0)::bigint AS max_rating,
  COALESCE(MIN(rating), 0)::bigint AS min_rating,
  count(*) FILTER (WHERE rating = 5) AS five_star_count
FROM book_reviews
WHERE user_id = $1;
