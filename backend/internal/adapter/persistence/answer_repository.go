package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// answerRepository は [repository.AnswerRepository] の GORM 実装。
type answerRepository struct {
	db *gorm.DB
}

// NewAnswerRepository は AnswerRepository の GORM 実装を返す。
func NewAnswerRepository(db *gorm.DB) repository.AnswerRepository {
	return &answerRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.AnswerRepository = (*answerRepository)(nil)

// Create は回答の作成と質問の回答数の加算を同一トランザクションで行う。
func (r *answerRepository) Create(ctx context.Context, answer *model.Answer) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(answer).Error; err != nil {
			return err
		}
		return tx.Model(&model.Question{}).Where("id = ?", answer.QuestionID).
			UpdateColumn("answer_count", gorm.Expr("answer_count + 1")).Error
	})
}

// FindByID は指定 ID の回答をユーザー情報付きで取得する。不在の場合は (nil, nil) を返す。
func (r *answerRepository) FindByID(ctx context.Context, id uint) (*model.Answer, error) {
	var answer model.Answer
	err := r.db.WithContext(ctx).Preload("User").First(&answer, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &answer, nil
}

// Update は既存の回答を更新する。
func (r *answerRepository) Update(ctx context.Context, answer *model.Answer) error {
	return r.db.WithContext(ctx).Save(answer).Error
}

// Delete は回答の論理削除と質問の回答数の減算（0 未満にはしない）を
// 同一トランザクションで行う。
func (r *answerRepository) Delete(ctx context.Context, answer *model.Answer) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(answer).Error; err != nil {
			return err
		}
		return tx.Model(&model.Question{}).Where("id = ?", answer.QuestionID).
			UpdateColumn("answer_count", gorm.Expr("GREATEST(answer_count - 1, 0)")).Error
	})
}

// FindByQuestionID は指定質問の回答を取得する（ベストアンサー優先 → 投票数降順 → 作成日昇順）。
func (r *answerRepository) FindByQuestionID(ctx context.Context, questionID uint) ([]model.Answer, error) {
	var answers []model.Answer
	err := r.db.WithContext(ctx).Preload("User").
		Where("question_id = ?", questionID).
		Order("is_best DESC, vote_count DESC, created_at ASC").
		Find(&answers).Error
	return answers, err
}

// FindByVoteRange は指定質問の回答を投票数の範囲で絞り込んで取得する（投票数降順 → 作成日昇順）。
func (r *answerRepository) FindByVoteRange(ctx context.Context, questionID uint, minVote, maxVote int) ([]model.Answer, error) {
	var answers []model.Answer
	err := r.db.WithContext(ctx).Preload("User").
		Where("question_id = ? AND vote_count >= ? AND vote_count <= ?", questionID, minVote, maxVote).
		Order("vote_count DESC, created_at ASC").
		Find(&answers).Error
	return answers, err
}

// SetBestAnswer は既存のベストアンサーの解除・指定回答の設定・質問の解決済み化を
// 同一トランザクションで行う。質問行を FOR UPDATE でロックして直列化し、
// 同時設定で複数の回答に is_best が立つ競合を防ぐ。
func (r *answerRepository) SetBestAnswer(ctx context.Context, questionID, answerID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var question model.Question
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&question, questionID).Error; err != nil {
			return err
		}

		// 既存のベストアンサーを解除する。対象が無くてもエラーにはならない。
		if err := tx.Model(&model.Answer{}).
			Where("question_id = ? AND is_best = ?", questionID, true).
			Update("is_best", false).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Answer{}).Where("id = ?", answerID).
			Update("is_best", true).Error; err != nil {
			return err
		}

		return tx.Model(&model.Question{}).Where("id = ?", questionID).
			Update("is_solved", true).Error
	})
}

// Vote は投票行の作成・更新と vote_count の増減を同一トランザクションで行う。
// 既に投票済みなら値を更新し、回答の投票数も差分だけ増減させる。
func (r *answerRepository) Vote(ctx context.Context, userID, answerID uint, value int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing model.AnswerVote
		err := tx.Where("user_id = ? AND answer_id = ?", userID, answerID).First(&existing).Error
		switch {
		case err == nil:
			diff := value - existing.Value
			existing.Value = value
			if err := tx.Save(&existing).Error; err != nil {
				return err
			}
			return tx.Model(&model.Answer{}).Where("id = ?", answerID).
				UpdateColumn("vote_count", gorm.Expr("vote_count + ?", diff)).Error
		case errors.Is(err, gorm.ErrRecordNotFound):
			vote := &model.AnswerVote{UserID: userID, AnswerID: answerID, Value: value}
			if err := tx.Create(vote).Error; err != nil {
				return err
			}
			return tx.Model(&model.Answer{}).Where("id = ?", answerID).
				UpdateColumn("vote_count", gorm.Expr("vote_count + ?", value)).Error
		default:
			return err
		}
	})
}

// RemoveVote は投票の削除と vote_count の減算を同一トランザクションで行う。
func (r *answerRepository) RemoveVote(ctx context.Context, userID, answerID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing model.AnswerVote
		if err := tx.Where("user_id = ? AND answer_id = ?", userID, answerID).First(&existing).Error; err != nil {
			return err
		}
		if err := tx.Delete(&existing).Error; err != nil {
			return err
		}
		return tx.Model(&model.Answer{}).Where("id = ?", answerID).
			UpdateColumn("vote_count", gorm.Expr("vote_count - ?", existing.Value)).Error
	})
}
