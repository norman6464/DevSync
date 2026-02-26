package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// QuestionRepository はQ&A質問データへのアクセスを提供するリポジトリ実装。
type QuestionRepository struct {
	db *gorm.DB
}

// NewQuestionRepository は新しいQuestionRepositoryインスタンスを生成する。
func NewQuestionRepository(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

// Create は新しい質問をデータベースに作成する。
func (r *QuestionRepository) Create(question *model.Question) error {
	return r.db.Create(question).Error
}

// FindByID は指定IDの質問をユーザー情報付きで取得する。
func (r *QuestionRepository) FindByID(id uint) (*model.Question, error) {
	var question model.Question
	err := r.db.Preload("User").First(&question, id).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

// FindAll は質問一覧をフィルタ・ソート・ページネーション付きで取得する。
// tagが指定された場合、そのタグを含む質問にフィルタする。
// sortが"votes"の場合は投票数降順、"unanswered"の場合は未回答のみに絞る。
func (r *QuestionRepository) FindAll(limit, offset int, tag string, sort string) ([]model.Question, int64, error) {
	var questions []model.Question
	var total int64

	query := r.db.Model(&model.Question{})

	if tag != "" {
		query = query.Where("tags ILIKE ?", "%\""+EscapeLikeChars(tag)+"\"%")
	}

	query.Count(&total)

	orderClause := "created_at DESC"
	switch sort {
	case "votes":
		orderClause = "vote_count DESC, created_at DESC"
	case "unanswered":
		query = query.Where("answer_count = 0")
		query.Count(&total) // フィルタ後に再カウント
	}

	err := query.Preload("User").
		Order(orderClause).
		Limit(limit).Offset(offset).
		Find(&questions).Error

	return questions, total, err
}

// Search は質問をタイトル・本文・タグで全文検索する。
func (r *QuestionRepository) Search(q string, limit, offset int) ([]model.Question, int64, error) {
	var questions []model.Question
	var total int64

	searchQuery := EscapeLikePattern(q)
	dbQuery := r.db.Model(&model.Question{}).
		Where("title ILIKE ? OR body ILIKE ? OR tags ILIKE ?", searchQuery, searchQuery, searchQuery)

	dbQuery.Count(&total)

	err := dbQuery.Preload("User").
		Order("vote_count DESC, created_at DESC").
		Limit(limit).Offset(offset).
		Find(&questions).Error

	return questions, total, err
}

// FindByUserID は指定ユーザーの質問をページネーション付きで取得する（新しい順）。
func (r *QuestionRepository) FindByUserID(userID uint, limit, offset int) ([]model.Question, int64, error) {
	var questions []model.Question
	var total int64
	q := r.db.Where("user_id = ?", userID)
	if err := q.Model(&model.Question{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&questions).Error
	return questions, total, err
}

// Update は既存の質問を更新する。
func (r *QuestionRepository) Update(question *model.Question) error {
	return r.db.Save(question).Error
}

// Delete は指定IDの質問を削除する。
func (r *QuestionRepository) Delete(id uint) error {
	return r.db.Delete(&model.Question{}, id).Error
}

// Vote は質問に投票する。既存投票がある場合は値を更新し、ない場合は新規作成する。
// 質問のvote_countも同時に更新する。
func (r *QuestionRepository) Vote(userID, questionID uint, value int) error {
	var existing model.QuestionVote
	err := r.db.Where("user_id = ? AND question_id = ?", userID, questionID).First(&existing).Error

	if err == nil {
		// 既存の投票を更新
		oldValue := existing.Value
		existing.Value = value
		if err := r.db.Save(&existing).Error; err != nil {
			return err
		}
		diff := value - oldValue
		return r.db.Model(&model.Question{}).Where("id = ?", questionID).
			UpdateColumn("vote_count", gorm.Expr("vote_count + ?", diff)).Error
	}

	// 新規投票を作成
	vote := &model.QuestionVote{
		UserID:     userID,
		QuestionID: questionID,
		Value:      value,
	}
	if err := r.db.Create(vote).Error; err != nil {
		return err
	}
	return r.db.Model(&model.Question{}).Where("id = ?", questionID).
		UpdateColumn("vote_count", gorm.Expr("vote_count + ?", value)).Error
}

// RemoveVote は質問への投票を取り消し、質問のvote_countから元の値を減算する。
func (r *QuestionRepository) RemoveVote(userID, questionID uint) error {
	var existing model.QuestionVote
	err := r.db.Where("user_id = ? AND question_id = ?", userID, questionID).First(&existing).Error
	if err != nil {
		return err
	}

	oldValue := existing.Value
	if err := r.db.Delete(&existing).Error; err != nil {
		return err
	}
	return r.db.Model(&model.Question{}).Where("id = ?", questionID).
		UpdateColumn("vote_count", gorm.Expr("vote_count - ?", oldValue)).Error
}

// FindSolved は解決済みの質問一覧をページネーション付きで取得する。
func (r *QuestionRepository) FindSolved(limit, offset int) ([]model.Question, int64, error) {
	var questions []model.Question
	var total int64

	query := r.db.Model(&model.Question{}).Where("is_solved = ?", true)
	query.Count(&total)

	err := query.Preload("User").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&questions).Error

	return questions, total, err
}

// FindUnanswered は回答が0件の質問をページネーション付きで取得する（新しい順）。
func (r *QuestionRepository) FindUnanswered(limit, offset int) ([]model.Question, int64, error) {
	var questions []model.Question
	var total int64

	query := r.db.Model(&model.Question{}).Where("answer_count = 0")
	query.Count(&total)

	err := query.Preload("User").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&questions).Error

	return questions, total, err
}

// GetUserVote は指定ユーザーの指定質問への投票値を取得する（未投票の場合は0を返す）。
func (r *QuestionRepository) GetUserVote(userID, questionID uint) (int, error) {
	var vote model.QuestionVote
	err := r.db.Where("user_id = ? AND question_id = ?", userID, questionID).First(&vote).Error
	if err != nil {
		return 0, nil // 未投票
	}
	return vote.Value, nil
}

// Bookmark は質問をブックマークする。
func (r *QuestionRepository) Bookmark(userID, questionID uint) error {
	bookmark := &model.QuestionBookmark{
		UserID:     userID,
		QuestionID: questionID,
	}
	return r.db.Create(bookmark).Error
}

// Unbookmark は質問のブックマークを解除する。
func (r *QuestionRepository) Unbookmark(userID, questionID uint) error {
	return r.db.Where("user_id = ? AND question_id = ?", userID, questionID).
		Delete(&model.QuestionBookmark{}).Error
}

// HasBookmarked は指定ユーザーが指定質問をブックマークしているかを返す。
func (r *QuestionRepository) HasBookmarked(userID, questionID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.QuestionBookmark{}).
		Where("user_id = ? AND question_id = ?", userID, questionID).
		Count(&count).Error
	return count > 0, err
}

// FindBookmarkedByUserID は指定ユーザーのブックマーク済み質問をページネーション付きで取得する。
func (r *QuestionRepository) FindBookmarkedByUserID(userID uint, limit, offset int) ([]model.Question, int64, error) {
	var questions []model.Question
	var total int64

	subQuery := r.db.Model(&model.QuestionBookmark{}).
		Select("question_id").
		Where("user_id = ?", userID)

	query := r.db.Model(&model.Question{}).Where("id IN (?)", subQuery)
	query.Count(&total)

	err := query.Preload("User").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&questions).Error

	return questions, total, err
}
