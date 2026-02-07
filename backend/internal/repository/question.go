package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

type QuestionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) Create(question *model.Question) error {
	return r.db.Create(question).Error
}

func (r *QuestionRepository) FindByID(id uint) (*model.Question, error) {
	var question model.Question
	err := r.db.Preload("User").First(&question, id).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *QuestionRepository) FindAll(limit, offset int, tag string, sort string) ([]model.Question, int64, error) {
	var questions []model.Question
	var total int64

	query := r.db.Model(&model.Question{})

	if tag != "" {
		query = query.Where("tags ILIKE ?", "%\""+tag+"\"%")
	}

	query.Count(&total)

	orderClause := "created_at DESC"
	switch sort {
	case "votes":
		orderClause = "vote_count DESC, created_at DESC"
	case "unanswered":
		query = query.Where("answer_count = 0")
		query.Count(&total) // recount after filter
	}

	err := query.Preload("User").
		Order(orderClause).
		Limit(limit).Offset(offset).
		Find(&questions).Error

	return questions, total, err
}

func (r *QuestionRepository) Search(q string, limit, offset int) ([]model.Question, int64, error) {
	var questions []model.Question
	var total int64

	searchQuery := "%" + q + "%"
	dbQuery := r.db.Model(&model.Question{}).
		Where("title ILIKE ? OR body ILIKE ? OR tags ILIKE ?", searchQuery, searchQuery, searchQuery)

	dbQuery.Count(&total)

	err := dbQuery.Preload("User").
		Order("vote_count DESC, created_at DESC").
		Limit(limit).Offset(offset).
		Find(&questions).Error

	return questions, total, err
}

func (r *QuestionRepository) FindByUserID(userID uint) ([]model.Question, error) {
	var questions []model.Question
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&questions).Error
	return questions, err
}

func (r *QuestionRepository) Update(question *model.Question) error {
	return r.db.Save(question).Error
}

func (r *QuestionRepository) Delete(id uint) error {
	return r.db.Delete(&model.Question{}, id).Error
}

// Vote operations
func (r *QuestionRepository) Vote(userID, questionID uint, value int) error {
	var existing model.QuestionVote
	err := r.db.Where("user_id = ? AND question_id = ?", userID, questionID).First(&existing).Error

	if err == nil {
		// Update existing vote
		oldValue := existing.Value
		existing.Value = value
		if err := r.db.Save(&existing).Error; err != nil {
			return err
		}
		diff := value - oldValue
		return r.db.Model(&model.Question{}).Where("id = ?", questionID).
			UpdateColumn("vote_count", gorm.Expr("vote_count + ?", diff)).Error
	}

	// Create new vote
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

func (r *QuestionRepository) GetUserVote(userID, questionID uint) (int, error) {
	var vote model.QuestionVote
	err := r.db.Where("user_id = ? AND question_id = ?", userID, questionID).First(&vote).Error
	if err != nil {
		return 0, nil // no vote
	}
	return vote.Value, nil
}
