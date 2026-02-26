package handler

import (
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// LearningLogServiceInterface はLearningLogHandlerが依存するサービスのインターフェース。
type LearningLogServiceInterface interface {
	Create(log *model.LearningLog) error
	BatchCreate(userID uint, logs []model.LearningLog) ([]model.LearningLog, error)
	ImportCSV(userID uint, data []byte) ([]model.LearningLog, error)
	GetByID(id, userID uint) (*model.LearningLog, error)
	GetByUserID(userID uint, limit, offset int) ([]model.LearningLog, int64, error)
	Update(id, userID uint, updates *model.LearningLog) (*model.LearningLog, error)
	Delete(id, userID uint) error
	GetStreakInfo(userID uint) (*model.StreakInfo, error)
	GetCalendarData(userID uint) ([]model.CalendarEntry, error)
	ExportCSV(userID uint, days int) ([]byte, error)
	ExportJSON(userID uint, days int) ([]byte, error)
	GetByCategory(userID uint, category string) ([]model.LearningLog, error)
	GetBySource(userID uint, source string) ([]model.LearningLog, error)
	GetWeeklyDuration(userID uint) (int, error)
	FavoriteLog(id, userID uint) error
	UnfavoriteLog(id, userID uint) error
	GetRecentCategories(userID uint) ([]string, error)
	GetLinkedLogs(goalID, userID uint, limit, offset int) ([]model.LearningLog, int64, error)
	GetGoalProgress(goalID, userID uint) (*model.GoalProgress, error)
	GetFavorites(userID uint, limit, offset int) ([]model.LearningLog, int64, error)
	GetMonthlySummary(userID uint, months int) ([]model.MonthlySummary, error)
}

// LearningLogHandler は学習ログ関連のHTTPハンドラ。
// 学習ログのCRUD・ストリーク・カレンダーデータの取得を処理する。
type LearningLogHandler struct {
	service LearningLogServiceInterface
}

// NewLearningLogHandler は新しいLearningLogHandlerインスタンスを生成する。
func NewLearningLogHandler(s LearningLogServiceInterface) *LearningLogHandler {
	return &LearningLogHandler{service: s}
}

// Create は新しい学習ログを作成する。
func (h *LearningLogHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	input := bindJSON[dto.CreateLearningLogRequest](c)
	if input == nil {
		return
	}

	log := &model.LearningLog{
		UserID:   userID,
		Title:    input.Title,
		Content:  input.Content,
		Category: model.LogCategory(input.Category),
		Duration: input.Duration,
		Source:   model.LogSource(input.Source),
		GoalID:   input.GoalID,
	}

	if err := h.service.Create(log); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, log)
}

// BatchCreate は複数の学習ログを一括作成する。
func (h *LearningLogHandler) BatchCreate(c *gin.Context) {
	userID := c.GetUint("userID")

	input := bindJSON[dto.BatchCreateLearningLogRequest](c)
	if input == nil {
		return
	}

	logs := make([]model.LearningLog, len(input.Logs))
	for i, l := range input.Logs {
		logs[i] = model.LearningLog{
			Title:    l.Title,
			Content:  l.Content,
			Category: model.LogCategory(l.Category),
			Duration: l.Duration,
			Source:   model.LogSource(l.Source),
		}
	}

	results, err := h.service.BatchCreate(userID, logs)
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, results)
}

// Update は指定された学習ログを更新する。
func (h *LearningLogHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	logID, ok := parseID(c, "id")
	if !ok {
		return
	}

	input := bindJSON[dto.UpdateLearningLogRequest](c)
	if input == nil {
		return
	}

	updates := &model.LearningLog{}
	if input.Title != nil {
		updates.Title = *input.Title
	}
	if input.Content != nil {
		updates.Content = *input.Content
	}
	if input.Category != nil {
		updates.Category = model.LogCategory(*input.Category)
	}
	if input.Duration != nil {
		updates.Duration = *input.Duration
	}

	log, err := h.service.Update(logID, userID, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, log)
}

// Delete は指定された学習ログを削除する。
func (h *LearningLogHandler) Delete(c *gin.Context) {
	handleDelete(c, h.service.Delete)
}

// GetByID は指定されたIDの学習ログを取得する。
func (h *LearningLogHandler) GetByID(c *gin.Context) {
	handleGetByID(c, h.service.GetByID)
}

// GetMyLogs は認証ユーザー自身の学習ログ一覧を取得する。
func (h *LearningLogHandler) GetMyLogs(c *gin.Context) {
	userID := c.GetUint("userID")
	limit, offset := parseLimitOffset(c)

	logs, total, err := h.service.GetByUserID(userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.LearningLogListResponse{
		Logs:   logs,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// GetByUserID は指定されたユーザーの学習ログ一覧を取得する。
func (h *LearningLogHandler) GetByUserID(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}
	limit, offset := parseLimitOffset(c)

	logs, total, err := h.service.GetByUserID(userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.LearningLogListResponse{
		Logs:   logs,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// GetStreakInfo は指定されたユーザーのストリーク情報を取得する。
func (h *LearningLogHandler) GetStreakInfo(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	info, err := h.service.GetStreakInfo(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, info)
}

// GetCalendarData はカレンダー表示用の日別学習ログ件数を取得する。
func (h *LearningLogHandler) GetCalendarData(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	entries, err := h.service.GetCalendarData(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(entries))
}

// GetByCategory はカテゴリで学習ログをフィルタリングして取得する。
func (h *LearningLogHandler) GetByCategory(c *gin.Context) {
	userID := c.GetUint("userID")
	category := c.Param("category")

	logs, err := h.service.GetByCategory(userID, category)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(logs))
}

// GetBySource はソース（manual/pomodoro）で学習ログをフィルタリングして取得する。
func (h *LearningLogHandler) GetBySource(c *gin.Context) {
	userID := c.GetUint("userID")
	source := c.Param("source")

	logs, err := h.service.GetBySource(userID, source)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(logs))
}

// ExportLogs は学習ログをCSVまたはJSON形式でダウンロードする。
// クエリパラメータ: period=7|30|90|all（デフォルト30日）、format=csv|json（デフォルトcsv）
func (h *LearningLogHandler) ExportLogs(c *gin.Context) {
	userID := c.GetUint("userID")

	days, ok := parseExportPeriod(c)
	if !ok {
		return
	}

	format := c.DefaultQuery("format", "csv")
	timestamp := time.Now().Format("20060102")

	switch format {
	case "json":
		jsonBytes, err := h.service.ExportJSON(userID, days)
		if err != nil {
			respondError(c, err)
			return
		}
		filename := fmt.Sprintf("learning-logs-%ddays-%s.json", days, timestamp)
		if days == 0 {
			filename = fmt.Sprintf("learning-logs-all-%s.json", timestamp)
		}
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Data(200, "application/json; charset=utf-8", jsonBytes)
	case "csv":
		csvBytes, err := h.service.ExportCSV(userID, days)
		if err != nil {
			respondError(c, err)
			return
		}
		filename := fmt.Sprintf("learning-logs-%ddays-%s.csv", days, timestamp)
		if days == 0 {
			filename = fmt.Sprintf("learning-logs-all-%s.csv", timestamp)
		}
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Data(200, "text/csv; charset=utf-8", csvBytes)
	default:
		respondBadRequest(c, "formatはcsv/jsonのいずれかを指定してください")
	}
}

// GetWeeklyDuration は指定ユーザーの過去7日間の学習時間合計を返す。
func (h *LearningLogHandler) GetWeeklyDuration(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	duration, err := h.service.GetWeeklyDuration(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.WeeklyDurationResponse{Duration: duration})
}

// GetMonthlySummary は指定ユーザーの月別学習サマリーを返す。
func (h *LearningLogHandler) GetMonthlySummary(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	months := parseQueryIntSilent(c, "months", 12)

	summaries, err := h.service.GetMonthlySummary(userID, months)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(summaries))
}

// Favorite は学習ログをお気に入りに設定する。
func (h *LearningLogHandler) Favorite(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.service.FavoriteLog(id, userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("学習ログをお気に入りに追加しました"))
}

// GetRecentCategories はユーザーの最近よく使うカテゴリを返す。
func (h *LearningLogHandler) GetRecentCategories(c *gin.Context) {
	userID := c.GetUint("userID")

	categories, err := h.service.GetRecentCategories(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(categories))
}

// GetFavorites はお気に入りに設定した学習ログ一覧を取得する。
func (h *LearningLogHandler) GetFavorites(c *gin.Context) {
	userID := c.GetUint("userID")
	limit, offset := parseLimitOffset(c)

	logs, total, err := h.service.GetFavorites(userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.LearningLogListResponse{
		Logs:   logs,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// Unfavorite は学習ログのお気に入りを解除する。
func (h *LearningLogHandler) Unfavorite(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.service.UnfavoriteLog(id, userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("学習ログのお気に入りを解除しました"))
}

// ImportCSV はCSVファイルから学習ログを一括インポートする。
func (h *LearningLogHandler) ImportCSV(c *gin.Context) {
	userID := c.GetUint("userID")

	file, err := c.FormFile("file")
	if err != nil {
		respondBadRequest(c, "CSVファイルが必要です")
		return
	}

	// ファイルサイズ制限（1MB）
	if file.Size > 1*1024*1024 {
		respondBadRequest(c, "ファイルサイズは1MB以下にしてください")
		return
	}

	src, err := file.Open()
	if err != nil {
		respondInternalError(c, "ファイルの読み取りに失敗しました")
		return
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		respondInternalError(c, "ファイルの読み取りに失敗しました")
		return
	}

	logs, err := h.service.ImportCSV(userID, data)
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, dto.ImportCSVResponse{
		Imported: len(logs),
		Logs:     logs,
	})
}

// GetLinkedLogs は指定ゴールに紐付いた学習ログ一覧を取得する。
func (h *LearningLogHandler) GetLinkedLogs(c *gin.Context) {
	userID := c.GetUint("userID")
	goalID, ok := parseID(c, "id")
	if !ok {
		return
	}

	limit, offset := parseLimitOffset(c)
	logs, total, err := h.service.GetLinkedLogs(goalID, userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.LearningLogListResponse{
		Logs:   logs,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// GetGoalProgress は指定ゴールの実績時間 vs 目標時間の進捗情報を返す。
func (h *LearningLogHandler) GetGoalProgress(c *gin.Context) {
	userID := c.GetUint("userID")
	goalID, ok := parseID(c, "id")
	if !ok {
		return
	}

	progress, err := h.service.GetGoalProgress(goalID, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, progress)
}

// GetMyStreakInfo は認証ユーザー自身のストリーク情報を取得する。
func (h *LearningLogHandler) GetMyStreakInfo(c *gin.Context) {
	userID := c.GetUint("userID")

	info, err := h.service.GetStreakInfo(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, info)
}

// GetMyCalendarData は認証ユーザー自身のカレンダーデータを取得する。
func (h *LearningLogHandler) GetMyCalendarData(c *gin.Context) {
	userID := c.GetUint("userID")

	entries, err := h.service.GetCalendarData(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(entries))
}

// GetMyWeeklyDuration は認証ユーザー自身の過去7日間の学習時間合計を返す。
func (h *LearningLogHandler) GetMyWeeklyDuration(c *gin.Context) {
	userID := c.GetUint("userID")

	duration, err := h.service.GetWeeklyDuration(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.WeeklyDurationResponse{Duration: duration})
}
