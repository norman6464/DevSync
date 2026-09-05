-- name: GetContributionRanking :many
-- 指定期間（weekly/monthly）のGitHubコントリビューションランキング。
-- コントリビューション数の合計で降順ソートし、上位50件を返す。
SELECT u.id AS user_id, u.username, u.name, u.avatar_url, COALESCE(SUM(gc.count), 0)::bigint AS score
FROM users u
JOIN git_hub_contributions gc ON gc.user_id = u.id
WHERE gc.date >= $1
GROUP BY u.id
HAVING SUM(gc.count) > 0
ORDER BY score DESC
LIMIT 50;

-- name: GetLanguageRanking :many
-- 指定プログラミング言語のバイト数ランキング。上位50件を返す。
SELECT u.id AS user_id, u.username, u.name, u.avatar_url, gls.bytes AS score
FROM users u
JOIN git_hub_language_stats gls ON gls.user_id = u.id
WHERE gls.language = $1
ORDER BY gls.bytes DESC
LIMIT 50;

-- name: GetLevelRanking :many
-- ユーザーのXP合計に基づくレベルランキング。上位50件を返す。
-- scoreはSELECTの別名なのでWHERE/HAVINGからは参照できない（PostgreSQLの仕様）ため、
-- 内側で合算してから外側で絞り込む。
SELECT * FROM (
    SELECT u.id AS user_id, u.username, u.name, u.avatar_url,
        COALESCE(ll.xp, 0) + COALESCE(p.xp, 0) + COALESCE(gh.xp, 0) +
        COALESCE(g.xp, 0) + COALESCE(c.xp, 0) + COALESCE(lk.xp, 0) AS score
    FROM users u
    LEFT JOIN (
        SELECT user_id, COUNT(*) * 10 + COALESCE(SUM(duration), 0)::bigint / 2 AS xp
        FROM learning_logs GROUP BY user_id
    ) ll ON ll.user_id = u.id
    LEFT JOIN (
        SELECT user_id, COUNT(*) * 30 AS xp
        FROM posts GROUP BY user_id
    ) p ON p.user_id = u.id
    LEFT JOIN (
        SELECT user_id, COUNT(DISTINCT date) * 5 AS xp
        FROM git_hub_contributions WHERE count > 0 GROUP BY user_id
    ) gh ON gh.user_id = u.id
    LEFT JOIN (
        SELECT user_id, COUNT(*) * 50 AS xp
        FROM learning_goals WHERE status = 'completed' GROUP BY user_id
    ) g ON g.user_id = u.id
    LEFT JOIN (
        SELECT user_id, COUNT(*) * 5 AS xp
        FROM comments GROUP BY user_id
    ) c ON c.user_id = u.id
    LEFT JOIN (
        SELECT user_id, COALESCE(SUM(like_count), 0)::bigint * 3 AS xp
        FROM posts GROUP BY user_id
    ) lk ON lk.user_id = u.id
) ranked
WHERE score > 0
ORDER BY score DESC
LIMIT 50;

-- name: ListAvailableLanguages :many
-- ランキング対象となるプログラミング言語の一覧をアルファベット順で返す。
SELECT DISTINCT language FROM git_hub_language_stats ORDER BY language;
