package service

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// LearningLogService は学習ログのビジネスロジックを提供する。
// 学習記録のCRUD操作と、ストリーク・カレンダーデータの取得を担当する。
type LearningLogService struct {
	repo repository.LearningLogRepositoryInterface
}

// NewLearningLogService は新しいLearningLogServiceインスタンスを生成する。
func NewLearningLogService(repo repository.LearningLogRepositoryInterface) *LearningLogService {
	return &LearningLogService{repo: repo}
}

// Create は新しい学習ログを作成する。
// Duration、Category、Sourceのバリデーションを行う。
func (s *LearningLogService) Create(log *model.LearningLog) error {
	// Duration: 0以上1440以下（24時間）
	if log.Duration < 0 || log.Duration > 1440 {
		return ErrBadRequest
	}
	// Category: 空文字（デフォルト適用）または有効な値のみ許可
	if log.Category != "" && !model.ValidCategories[log.Category] {
		return ErrBadRequest
	}
	// Source: 空文字（デフォルト"manual"）または有効な値のみ許可
	if log.Source != "" && !model.ValidSources[log.Source] {
		return ErrBadRequest
	}
	return s.repo.Create(log)
}

// GetByID は指定IDの学習ログを取得する。
func (s *LearningLogService) GetByID(id uint) (*model.LearningLog, error) {
	return s.repo.FindByID(id)
}

// GetByUserID は指定ユーザーの全学習ログを取得する。
func (s *LearningLogService) GetByUserID(userID uint) ([]model.LearningLog, error) {
	return s.repo.GetByUserID(userID)
}

// findAndCheckOwnership は学習ログを取得し、指定ユーザーが所有者かを検証する。
func (s *LearningLogService) findAndCheckOwnership(id, userID uint) (*model.LearningLog, error) {
	log, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if log.UserID != userID {
		return nil, ErrForbidden
	}
	return log, nil
}

// Update は所有権を検証した後、学習ログを更新する。
func (s *LearningLogService) Update(id, userID uint, updates *model.LearningLog) (*model.LearningLog, error) {
	log, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	if updates.Title != "" {
		log.Title = updates.Title
	}
	if updates.Content != "" {
		log.Content = updates.Content
	}
	if updates.Category != "" {
		log.Category = updates.Category
	}
	if updates.Duration != 0 {
		if updates.Duration < 0 || updates.Duration > 1440 {
			return nil, ErrBadRequest
		}
		log.Duration = updates.Duration
	}

	if err := s.repo.Update(log); err != nil {
		return nil, err
	}
	return log, nil
}

// Delete は所有権を検証した後、学習ログを削除する。
func (s *LearningLogService) Delete(id, userID uint) error {
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.Delete(id, userID)
}

// GetStreakInfo は指定ユーザーの学習ストリーク情報を取得する。
func (s *LearningLogService) GetStreakInfo(userID uint) (*model.StreakInfo, error) {
	return s.repo.GetStreakInfo(userID)
}

// GetCalendarData はカレンダー表示用の日別学習ログ集計データを取得する。
func (s *LearningLogService) GetCalendarData(userID uint) ([]model.CalendarEntry, error) {
	return s.repo.GetCalendarData(userID)
}

// ExportCSV は指定ユーザーの学習ログをCSV形式でエクスポートする。
// days: 取得する過去の日数（0は全期間）。負の値はバリデーションエラー。
func (s *LearningLogService) ExportCSV(userID uint, days int) ([]byte, error) {
	if days < 0 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "期間は0以上の値を指定してください", nil)
	}

	logs, err := s.repo.GetByPeriod(userID, days)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	// ヘッダー行（BOM付きUTF-8でExcel互換）
	buf.WriteString("\xef\xbb\xbf")
	if err := w.Write([]string{"日付", "カテゴリ", "タイトル", "学習時間(分)", "メモ"}); err != nil {
		return nil, err
	}

	for _, log := range logs {
		row := []string{
			log.CreatedAt.Format("2006-01-02"),
			string(log.Category),
			log.Title,
			fmt.Sprintf("%d", log.Duration),
			log.Content,
		}
		if err := w.Write(row); err != nil {
			return nil, err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
