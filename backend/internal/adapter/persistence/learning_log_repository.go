package persistence

import (
	"context"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// learningLogRepository は [repository.LearningLogRepository] の sqlc(pgx) 実装。
// CreateBatch は複数件の作成を1トランザクションで行うため、
// Queries だけでなくトランザクションを開始できる *pgxpool.Pool を直接保持する。
type learningLogRepository struct {
	pool *pgxpool.Pool
	q    *sqlcgen.Queries
}

// NewLearningLogRepository は LearningLogRepository の sqlc(pgx) 実装を返す。
func NewLearningLogRepository(pool *pgxpool.Pool) repository.LearningLogRepository {
	return &learningLogRepository{pool: pool, q: sqlcgen.New(pool)}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.LearningLogRepository = (*learningLogRepository)(nil)

func toModelLearningLog(row sqlcgen.LearningLog) model.LearningLog {
	return model.LearningLog{
		ID:         uint(row.ID),
		UserID:     uint(row.UserID),
		Title:      row.Title,
		Content:    row.Content,
		Category:   model.LogCategory(fromStringPtr(row.Category)),
		Duration:   int(fromInt64PtrValue(row.Duration)),
		GoalID:     fromInt64PtrToUintPtr(row.GoalID),
		Source:     model.LogSource(fromStringPtr(row.Source)),
		IsFavorite: row.IsFavorite,
		CreatedAt:  timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:  timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

func toModelLearningLogs(rows []sqlcgen.LearningLog) []model.LearningLog {
	logs := make([]model.LearningLog, len(rows))
	for i, row := range rows {
		logs[i] = toModelLearningLog(row)
	}
	return logs
}

// createLearningLogWith は指定の Queries（トランザクション経由も可）で学習ログを作成する。
// Sourceは移行前のGORMの `gorm:"default:'manual'"` に相当し、未指定（ゼロ値）なら manual を補う。
func createLearningLogWith(ctx context.Context, q *sqlcgen.Queries, log *model.LearningLog) error {
	source := log.Source
	if source == "" {
		source = model.LogSourceManual
	}

	row, err := q.CreateLearningLog(ctx, sqlcgen.CreateLearningLogParams{
		UserID:     int64(log.UserID),
		Title:      log.Title,
		Content:    log.Content,
		Category:   (*string)(&log.Category),
		Duration:   toInt64Ptr(log.Duration),
		GoalID:     toInt64PtrFromUintPtr(log.GoalID),
		Source:     (*string)(&source),
		IsFavorite: log.IsFavorite,
	})
	if err != nil {
		return err
	}
	*log = toModelLearningLog(row)
	return nil
}

// Create は新しい学習ログを作成する。
func (r *learningLogRepository) Create(ctx context.Context, log *model.LearningLog) error {
	return createLearningLogWith(ctx, r.q, log)
}

// CreateBatch は複数の学習ログを一括作成する。
// GORMのCreate(&slice)（単一トランザクション相当）に合わせ、1トランザクション内で
// ループ挿入することでアトミック性を維持する。
func (r *learningLogRepository) CreateBatch(ctx context.Context, logs []model.LearningLog) error {
	if len(logs) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	for i := range logs {
		if err := createLearningLogWith(ctx, q, &logs[i]); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// Update は既存の学習ログを更新する（GORMのSave＝全カラム上書きに相当）。
func (r *learningLogRepository) Update(ctx context.Context, log *model.LearningLog) error {
	row, err := r.q.UpdateLearningLog(ctx, sqlcgen.UpdateLearningLogParams{
		ID:         int64(log.ID),
		Title:      log.Title,
		Content:    log.Content,
		Category:   (*string)(&log.Category),
		Duration:   toInt64Ptr(log.Duration),
		GoalID:     toInt64PtrFromUintPtr(log.GoalID),
		Source:     (*string)(&log.Source),
		IsFavorite: log.IsFavorite,
	})
	if err != nil {
		return err
	}
	*log = toModelLearningLog(row)
	return nil
}

// Delete は所有者本人の学習ログを削除する。
func (r *learningLogRepository) Delete(ctx context.Context, id, userID uint) error {
	return r.q.DeleteLearningLog(ctx, sqlcgen.DeleteLearningLogParams{
		ID:     int64(id),
		UserID: int64(userID),
	})
}

// FindByID は指定 ID の学習ログを取得する。不在の場合は (nil, nil) を返す。
func (r *learningLogRepository) FindByID(ctx context.Context, id uint) (*model.LearningLog, error) {
	row, err := r.q.GetLearningLogByID(ctx, int64(id))
	if isNoRows(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	log := toModelLearningLog(row)
	return &log, nil
}

// GetByUserID はユーザーの学習ログを作成日の新しい順でページ取得し、総数も返す。
func (r *learningLogRepository) GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	total, err := r.q.CountLearningLogsByUser(ctx, int64(userID))
	if err != nil {
		return nil, 0, err
	}
	rows, err := r.q.ListLearningLogsByUser(ctx, sqlcgen.ListLearningLogsByUserParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	return toModelLearningLogs(rows), total, nil
}

// GetFavorites はお気に入りの学習ログを作成日の新しい順でページ取得し、総数も返す。
func (r *learningLogRepository) GetFavorites(ctx context.Context, userID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	total, err := r.q.CountFavoriteLearningLogs(ctx, int64(userID))
	if err != nil {
		return nil, 0, err
	}
	rows, err := r.q.ListFavoriteLearningLogs(ctx, sqlcgen.ListFavoriteLearningLogsParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	return toModelLearningLogs(rows), total, nil
}

// GetByGoalID は指定ゴールに紐付いた学習ログを作成日の新しい順でページ取得し、総数も返す。
func (r *learningLogRepository) GetByGoalID(ctx context.Context, goalID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	goalID64 := int64(goalID)

	total, err := r.q.CountLearningLogsByGoal(ctx, &goalID64)
	if err != nil {
		return nil, 0, err
	}
	rows, err := r.q.ListLearningLogsByGoal(ctx, sqlcgen.ListLearningLogsByGoalParams{
		GoalID: &goalID64,
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	return toModelLearningLogs(rows), total, nil
}

// GetByCategory はカテゴリで絞り込んだ学習ログを作成日の新しい順で取得する。
func (r *learningLogRepository) GetByCategory(ctx context.Context, userID uint, category string) ([]model.LearningLog, error) {
	rows, err := r.q.ListLearningLogsByCategory(ctx, sqlcgen.ListLearningLogsByCategoryParams{
		UserID:   int64(userID),
		Category: &category,
	})
	if err != nil {
		return nil, err
	}
	return toModelLearningLogs(rows), nil
}

// GetBySource はソースで絞り込んだ学習ログを作成日の新しい順で取得する。
func (r *learningLogRepository) GetBySource(ctx context.Context, userID uint, source string) ([]model.LearningLog, error) {
	rows, err := r.q.ListLearningLogsBySource(ctx, sqlcgen.ListLearningLogsBySourceParams{
		UserID: int64(userID),
		Source: &source,
	})
	if err != nil {
		return nil, err
	}
	return toModelLearningLogs(rows), nil
}

// GetByPeriod は直近 days 日の学習ログを作成日の新しい順で取得する。days が 0 以下なら全期間。
func (r *learningLogRepository) GetByPeriod(ctx context.Context, userID uint, days int) ([]model.LearningLog, error) {
	var since *time.Time
	if days > 0 {
		t := time.Now().AddDate(0, 0, -days)
		since = &t
	}

	rows, err := r.q.ListLearningLogsByPeriod(ctx, sqlcgen.ListLearningLogsByPeriodParams{
		UserID: int64(userID),
		Since:  toTimestamptz(since),
	})
	if err != nil {
		return nil, err
	}
	return toModelLearningLogs(rows), nil
}

// SumDurationByPeriod は直近 days 日の学習時間合計（分）を返す。
func (r *learningLogRepository) SumDurationByPeriod(ctx context.Context, userID uint, days int) (int, error) {
	since := time.Now().AddDate(0, 0, -days)
	total, err := r.q.SumLearningLogDurationSince(ctx, sqlcgen.SumLearningLogDurationSinceParams{
		UserID:    int64(userID),
		CreatedAt: toTimestamptzNotNull(since),
	})
	return int(total), err
}

// SumDurationByGoalID は指定ゴールに紐付いた学習ログの学習時間合計（分）を返す。
func (r *learningLogRepository) SumDurationByGoalID(ctx context.Context, goalID uint) (int, error) {
	goalID64 := int64(goalID)
	total, err := r.q.SumLearningLogDurationByGoal(ctx, &goalID64)
	return int(total), err
}

// CountByUserID はユーザーの学習ログ総数を返す。
func (r *learningLogRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountLearningLogsByUser(ctx, int64(userID))
}

// GetStreakInfo は学習ログの日付から連続学習情報を算出する。
func (r *learningLogRepository) GetStreakInfo(ctx context.Context, userID uint) (*model.StreakInfo, error) {
	dateRows, err := r.q.ListDistinctLearningLogDates(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	dates := make([]time.Time, len(dateRows))
	for i, d := range dateRows {
		dates[i] = d.Time
	}
	return calcStreakInfo(dates, time.Now()), nil
}

// calcStreakInfo は学習日の一覧から連続学習情報を組み立てる。
// 直近のログが今日または昨日のときだけ現在の連続日数を数える。
func calcStreakInfo(dates []time.Time, now time.Time) *model.StreakInfo {
	info := &model.StreakInfo{TotalDays: len(dates)}
	if len(dates) == 0 {
		return info
	}

	sort.Slice(dates, func(i, j int) bool { return dates[i].After(dates[j]) })
	info.LastLogDate = dates[0].Format("2006-01-02")

	if isTodayOrYesterday(normalizeToCalendarDay(dates[0]), normalizeToCalendarDay(now)) {
		info.CurrentStreak = 1
		for i := 1; i < len(dates); i++ {
			// 同日の重複は「翌暦日」にならないため連続とみなさない。
			if !isNextCalendarDay(normalizeToCalendarDay(dates[i-1]), normalizeToCalendarDay(dates[i])) {
				break
			}
			info.CurrentStreak++
		}
	}

	longest, streak := 1, 1
	for i := 1; i < len(dates); i++ {
		if isNextCalendarDay(normalizeToCalendarDay(dates[i-1]), normalizeToCalendarDay(dates[i])) {
			streak++
			if streak > longest {
				longest = streak
			}
			continue
		}
		streak = 1
	}
	info.LongestStreak = longest

	return info
}

// GetRecentCategories は使用回数の多い順にカテゴリを limit 件返す。
func (r *learningLogRepository) GetRecentCategories(ctx context.Context, userID uint, limit int) ([]string, error) {
	rows, err := r.q.ListRecentLearningLogCategories(ctx, sqlcgen.ListRecentLearningLogCategoriesParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
	})
	if err != nil {
		return nil, err
	}
	categories := make([]string, len(rows))
	for i, row := range rows {
		categories[i] = fromStringPtr(row)
	}
	return categories, nil
}

// GetCalendarData はカレンダー表示用の日別ログ件数を取得する。
func (r *learningLogRepository) GetCalendarData(ctx context.Context, userID uint) ([]model.CalendarEntry, error) {
	rows, err := r.q.ListLearningLogCalendarData(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	entries := make([]model.CalendarEntry, len(rows))
	for i, row := range rows {
		entries[i] = model.CalendarEntry{
			Date:  row.Date.Time.Format("2006-01-02"),
			Count: int(row.Count),
		}
	}
	return entries, nil
}

// GetMonthlySummary は直近 months ヶ月の月別サマリー（合計時間・ログ件数）を取得する。
func (r *learningLogRepository) GetMonthlySummary(ctx context.Context, userID uint, months int) ([]model.MonthlySummary, error) {
	startDate := time.Now().AddDate(0, -months, 0)

	rows, err := r.q.ListLearningLogMonthlySummary(ctx, sqlcgen.ListLearningLogMonthlySummaryParams{
		UserID:    int64(userID),
		CreatedAt: toTimestamptzNotNull(startDate),
	})
	if err != nil {
		return nil, err
	}
	summaries := make([]model.MonthlySummary, len(rows))
	for i, row := range rows {
		summaries[i] = model.MonthlySummary{
			Month:        row.Month,
			TotalMinutes: int(row.TotalMinutes),
			LogCount:     int(row.LogCount),
		}
	}
	return summaries, nil
}
