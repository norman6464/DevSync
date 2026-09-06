-- name: SumGitHubContributionsByUser :one
SELECT COALESCE(SUM(count), 0)::bigint FROM git_hub_contributions WHERE user_id = $1;

-- name: ListGitHubContributionsByUser :many
-- GitHub連続コントリビューション日数の算出に使う（countが正の日のみ、新しい順）。
SELECT contributed_on, count FROM git_hub_contributions
WHERE user_id = $1 AND count > 0
ORDER BY contributed_on DESC;
