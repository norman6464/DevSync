package handler

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// LearningLogServiceInterface はLearningLogHandlerが依存するサービスのインターフェース。
type LearningLogServiceInterface interface {
	Create(log *model.LearningLog) error
	GetByID(id uint) (*model.LearningLog, error)
	GetByUserID(userID uint) ([]model.LearningLog, error)
	Update(id, userID uint, updates *model.LearningLog) (*model.LearningLog, error)
	Delete(id, userID uint) error
	GetStreakInfo(userID uint) (*model.StreakInfo, error)
	GetCalendarData(userID uint) ([]model.CalendarEntry, error)
	ExportCSV(userID uint, days int) ([]byte, error)
	GetByCategory(userID uint, category string) ([]model.LearningLog, error)
	GetBySource(userID uint, source string) ([]model.LearningLog, error)
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
	}

	// カテゴリが未指定の場合はデフォルト値を設定
	if input.Category == "" {
		log.Category = model.LogCategoryOther
	}

	if err := h.service.Create(log); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, log)
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
	userID := c.GetUint("userID")
	logID, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.service.Delete(logID, userID); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}

// GetByID は指定されたIDの学習ログを取得する。
func (h *LearningLogHandler) GetByID(c *gin.Context) {
	logID, ok := parseID(c, "id")
	if !ok {
		return
	}

	log, err := h.service.GetByID(logID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, log)
}

// GetMyLogs は認証ユーザー自身の学習ログ一覧を取得する。
func (h *LearningLogHandler) GetMyLogs(c *gin.Context) {
	userID := c.GetUint("userID")

	logs, err := h.service.GetByUserID(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, logs)
}

// GetByUserID は指定されたユーザーの学習ログ一覧を取得する。
func (h *LearningLogHandler) GetByUserID(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	logs, err := h.service.GetByUserID(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, logs)
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

	respondOK(c, entries)
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
	if logs == nil {
		logs = []model.LearningLog{}
	}

	respondOK(c, logs)
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
	if logs == nil {
		logs = []model.LearningLog{}
	}

	respondOK(c, logs)
}

// ExportLogs は学習ログをCSV形式でダウンロードする。
// クエリパラメータ: period=7|30|90|all（デフォルト30日）
func (h *LearningLogHandler) ExportLogs(c *gin.Context) {
	userID := c.GetUint("userID")

	days, ok := parseExportPeriod(c)
	if !ok {
		return
	}

	csvBytes, err := h.service.ExportCSV(userID, days)
	if err != nil {
		respondError(c, err)
		return
	}

	filename := fmt.Sprintf("learning-logs-%ddays-%s.csv", days, time.Now().Format("20060102"))
	if days == 0 {
		filename = fmt.Sprintf("learning-logs-all-%s.csv", time.Now().Format("20060102"))
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(200, "text/csv; charset=utf-8", csvBytes)
}
