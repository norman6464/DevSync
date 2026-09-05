-- name: SumGitHubContributionsByUser :one
SELECT COALESCE(SUM(count), 0)::bigint FROM git_hub_contributions WHERE user_id = $1;

-- name: ListGitHubContributionsByUser :many
-- GitHub連続コントリビューション日数の算出に使う（countが正の日のみ、新しい順）。
SELECT date, count FROM git_hub_contributions
WHERE user_id = $1 AND count > 0
ORDER BY date DESC;

-- name: CountAnswersByUserIncludingDeleted :one
-- qa_stats.sql の CountAnswersByUser は deleted_at IS NULL で絞るが、
-- こちらは既存の GORM Raw SQL 実装（db.Raw、GORMのsoft-deleteスコープ非適用）と
-- 挙動を変えないため、論理削除された回答も含めてカウントする。
SELECT COUNT(*) FROM answers WHERE user_id = $1;
