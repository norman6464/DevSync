package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// AnswerRepository はQ&A回答データへのアクセスを提供するリポジトリ実装。
type AnswerRepository struct {
	db *gorm.DB
}

// NewAnswerRepository は新しいAnswerRepositoryインスタンスを生成する。
func NewAnswerRepository(db *gorm.DB) *AnswerRepository {
	return &AnswerRepository{db: db}
}

// Create は質問への回答を作成し、質問のanswer_countをインクリメントする。
func (r *AnswerRepository) Create(answer *model.Answer) error {
	if err := r.db.Create(answer).Error; err != nil {
		return err
	}
	return r.db.Model(&model.Question{}).Where("id = ?", answer.QuestionID).
		UpdateColumn("answer_count", gorm.Expr("answer_count + 1")).Error
}

// FindByQuestionID は指定質問の全回答を取得する。
// ベストアンサー→投票数降順→作成日昇順でソートされる。
func (r *AnswerRepository) FindByQuestionID(questionID uint) ([]model.Answer, error) {
	var answers []model.Answer
	err := r.db.Preload("User").
		Where("question_id = ?", questionID).
		Order("is_best DESC, vote_count DESC, created_at ASC").
		Find(&answers).Error
	return answers, err
}

// FindByID は指定IDの回答をユーザー情報付きで取得する。
func (r *AnswerRepository) FindByID(id uint) (*model.Answer, error) {
	var answer model.Answer
	err := r.db.Preload("User").First(&answer, id).Error
	if err != nil {
		return nil, err
	}
	return &answer, nil
}

// Update は既存の回答を更新する。
func (r *AnswerRepository) Update(answer *model.Answer) error {
	return r.db.Save(answer).Error
}

// Delete は回答を削除し、質問のanswer_countをデクリメントする。
func (r *AnswerRepository) Delete(answer *model.Answer) error {
	if err := r.db.Delete(answer).Error; err != nil {
		return err
	}
	return r.db.Model(&model.Question{}).Where("id = ?", answer.QuestionID).
		UpdateColumn("answer_count", gorm.Expr("GREATEST(answer_count - 1, 0)")).Error
}

// SetBestAnswer は指定回答をベストアンサーに設定する。
// 既存のベストアンサーを解除してから新しいベストアンサーを設定し、
// 質問をis_solved=trueに更新する。
func (r *AnswerRepository) SetBestAnswer(questionID, answerID uint) error {
	// 既存のベストアンサーを解除
	r.db.Model(&model.Answer{}).
		Where("question_id = ? AND is_best = ?", questionID, true).
		Update("is_best", false)

	// 新しいベストアンサーを設定
	if err := r.db.Model(&model.Answer{}).Where("id = ?", answerID).
		Update("is_best", true).Error; err != nil {
		return err
	}

	// 質問を解決済みにマーク
	return r.db.Model(&model.Question{}).Where("id = ?", questionID).
		Update("is_solved", true).Error
}

// Vote は回答に投票する。既存投票がある場合は値を更新し、ない場合は新規作成する。
func (r *AnswerRepository) Vote(userID, answerID uint, value int) error {
	var existing model.AnswerVote
	err := r.db.Where("user_id = ? AND answer_id = ?", userID, answerID).First(&existing).Error

	if err == nil {
		// 既存の投票を更新
		oldValue := existing.Value
		existing.Value = value
		if err := r.db.Save(&existing).Error; err != nil {
			return err
		}
		diff := value - oldValue
		return r.db.Model(&model.Answer{}).Where("id = ?", answerID).
			UpdateColumn("vote_count", gorm.Expr("vote_count + ?", diff)).Error
	}

	// 新規投票を作成
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

// RemoveVote は回答への投票を取り消す。
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

// GetUserVotes は指定ユーザーの複数回答への投票値をマップで取得する。
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
