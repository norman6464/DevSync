package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

type AnswerRepository struct {
	db *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) *AnswerRepository {
	return &AnswerRepository{db: db}
}

func (r *AnswerRepository) Create(answer *model.Answer) error {
	if err := r.db.Create(answer).Error; err != nil {
		return err
	}
	return r.db.Model(&model.Question{}).Where("id = ?", answer.QuestionID).
		UpdateColumn("answer_count", gorm.Expr("answer_count + 1")).Error
}

func (r *AnswerRepository) FindByQuestionID(questionID uint) ([]model.Answer, error) {
	var answers []model.Answer
	err := r.db.Preload("User").
		Where("question_id = ?", questionID).
		Order("is_best DESC, vote_count DESC, created_at ASC").
		Find(&answers).Error
	return answers, err
}

func (r *AnswerRepository) FindByID(id uint) (*model.Answer, error) {
	var answer model.Answer
	err := r.db.Preload("User").First(&answer, id).Error
	if err != nil {
		return nil, err
	}
	return &answer, nil
}

func (r *AnswerRepository) Update(answer *model.Answer) error {
	return r.db.Save(answer).Error
}

func (r *AnswerRepository) Delete(answer *model.Answer) error {
	if err := r.db.Delete(answer).Error; err != nil {
		return err
	}
	return r.db.Model(&model.Question{}).Where("id = ?", answer.QuestionID).
		UpdateColumn("answer_count", gorm.Expr("GREATEST(answer_count - 1, 0)")).Error
}

func (r *AnswerRepository) SetBestAnswer(questionID, answerID uint) error {
	// Remove existing best answer
	r.db.Model(&model.Answer{}).
		Where("question_id = ? AND is_best = ?", questionID, true).
		Update("is_best", false)

	// Set new best answer
	if err := r.db.Model(&model.Answer{}).Where("id = ?", answerID).
		Update("is_best", true).Error; err != nil {
		return err
	}

	// Mark question as solved
	return r.db.Model(&model.Question{}).Where("id = ?", questionID).
		Update("is_solved", true).Error
}

// Vote operations
func (r *AnswerRepository) Vote(userID, answerID uint, value int) error {
	var existing model.AnswerVote
	err := r.db.Where("user_id = ? AND answer_id = ?", userID, answerID).First(&existing).Error

	if err == nil {
		oldValue := existing.Value
		existing.Value = value
		if err := r.db.Save(&existing).Error; err != nil {
			return err
		}
		diff := value - oldValue
		return r.db.Model(&model.Answer{}).Where("id = ?", answerID).
			UpdateColumn("vote_count", gorm.Expr("vote_count + ?", diff)).Error
	}

	vote := &model.AnswerVote{
		UserID:   userID,
		AnswerID: answerID,
		Value:    value,
	}
	if err := r.db.Create(vote).Error; err != nil {
		return err
	}
	return r.db.Model(&model.Answer{}).Where("id = ?", answerID).
		UpdateColumn("vote_count", gorm.Expr("vote_count + ?", value)).Error
}

func (r *AnswerRepository) RemoveVote(userID, answerID uint) error {
	var existing model.AnswerVote
	err := r.db.Where("user_id = ? AND answer_id = ?", userID, answerID).First(&existing).Error
	if err != nil {
		return err
	}

	oldValue := existing.Value
	if err := r.db.Delete(&existing).Error; err != nil {
		return err
	}
	return r.db.Model(&model.Answer{}).Where("id = ?", answerID).
		UpdateColumn("vote_count", gorm.Expr("vote_count - ?", oldValue)).Error
}

func (r *AnswerRepository) GetUserVotes(userID uint, answerIDs []uint) (map[uint]int, error) {
	var votes []model.AnswerVote
	err := r.db.Where("user_id = ? AND answer_id IN ?", userID, answerIDs).Find(&votes).Error
	if err != nil {
		return nil, err
	}
	result := make(map[uint]int)
	for _, v := range votes {
		result[v.AnswerID] = v.Value
	}
	return result, nil
}
