package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// questionRepository は [repository.QuestionRepository] の GORM 実装。
//
// 旧 repository パッケージにも同じテーブルを扱う実装が残っている。answer がまだそちらを
// 使っているため、移行が一巡するまで新旧のアダプタが並存する。
type questionRepository struct {
	db *gorm.DB
}

// NewQuestionRepository は QuestionRepository の GORM 実装を返す。
func NewQuestionRepository(db *gorm.DB) repository.QuestionRepository {
	return &questionRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.QuestionRepository = (*questionRepository)(nil)

// Create は新しい質問を作成する。
func (r *questionRepository) Create(ctx context.Context, question *model.Question) error {
	return r.db.WithContext(ctx).Create(question).Error
}

// FindByID は指定 ID の質問をユーザー情報付きで取得する。不在の場合は (nil, nil) を返す。
func (r *questionRepository) FindByID(ctx context.Context, id uint) (*model.Question, error) {
	var question model.Question
	err := r.db.WithContext(ctx).Preload("User").First(&question, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &question, nil
}

// Update は既存の質問を更新する。
func (r *questionRepository) Update(ctx context.Context, question *model.Question) error {
	return r.db.WithContext(ctx).Save(question).Error
}

// Delete は質問を論理削除する。
func (r *questionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Question{}, id).Error
}

// FindAll は質問一覧をタグ絞り込み・ソート・ページネーション付きで取得する。
// sort が "votes" なら投票数降順、"unanswered" なら回答が 0 件のものだけに絞る。
func (r *questionRepository) FindAll(ctx context.Context, limit, offset int, tag, sort string) ([]model.Question, int64, error) {
	scope := r.db.WithContext(ctx).Model(&model.Question{})
	if tag != "" {
		// タグは JSON 配列形式の文字列なので、引用符ごと含む部分一致で絞り込む。
		scope = scope.Where("tags ILIKE ?", "%\""+escapeLikeChars(tag)+"\"%")
	}

	orderClause := "created_at DESC"
	switch sort {
	case "votes":
		orderClause = "vote_count DESC, created_at DESC"
	case "unanswered":
		scope = scope.Where("answer_count = 0")
	}

	return r.paginatedQuestions(scope, orderClause, limit, offset)
}

// Search は質問をタイトル・本文・タグで部分一致検索する（投票数降順）。
func (r *questionRepository) Search(ctx context.Context, query string, limit, offset int) ([]model.Question, int64, error) {
	pattern := escapeLikePattern(query)
	scope := r.db.WithContext(ctx).Model(&model.Question{}).
		Where("title ILIKE ? OR body ILIKE ? OR tags ILIKE ?", pattern, pattern, pattern)

	return r.paginatedQuestions(scope, "vote_count DESC, created_at DESC", limit, offset)
}

// FindByUserID は指定ユーザーの質問をページネーション付きで取得する（新しい順）。
// 一覧系の中でこれだけユーザー情報をプリロードしない（移行前からの挙動）。
func (r *questionRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Question, int64, error) {
	scope := r.db.WithContext(ctx).Model(&model.Question{}).Where("user_id = ?", userID)

	var total int64
	if err := scope.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var questions []model.Question
	err := scope.Session(&gorm.Session{}).
		Order("created_at DESC").Limit(limit).Offset(offset).
		Find(&questions).Error
	return questions, total, err
}

// FindSolved は解決済みの質問をページネーション付きで取得する（新しい順）。
func (r *questionRepository) FindSolved(ctx context.Context, limit, offset int) ([]model.Question, int64, error) {
	scope := r.db.WithContext(ctx).Model(&model.Question{}).Where("is_solved = ?", true)
	return r.paginatedQuestions(scope, "created_at DESC", limit, offset)
}

// FindUnanswered は回答が 0 件の質問をページネーション付きで取得する（新しい順）。
func (r *questionRepository) FindUnanswered(ctx context.Context, limit, offset int) ([]model.Question, int64, error) {
	scope := r.db.WithContext(ctx).Model(&model.Question{}).Where("answer_count = 0")
	return r.paginatedQuestions(scope, "created_at DESC", limit, offset)
}

// FindBookmarkedByUserID は指定ユーザーがブックマークした質問を取得する（新しい順）。
func (r *questionRepository) FindBookmarkedByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Question, int64, error) {
	bookmarked := r.db.WithContext(ctx).Model(&model.QuestionBookmark{}).
		Select("question_id").Where("user_id = ?", userID)
	scope := r.db.WithContext(ctx).Model(&model.Question{}).Where("id IN (?)", bookmarked)
	return r.paginatedQuestions(scope, "created_at DESC", limit, offset)
}

// paginatedQuestions は絞り込み済みクエリに対して総件数とページを取得する共通処理。
func (r *questionRepository) paginatedQuestions(scope *gorm.DB, orderClause string, limit, offset int) ([]model.Question, int64, error) {
	var total int64
	if err := scope.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var questions []model.Question
	err := scope.Session(&gorm.Session{}).
		Preload("User").
		Order(orderClause).Limit(limit).Offset(offset).
		Find(&questions).Error
	return questions, total, err
}

// CountByUserID は指定ユーザーの質問総数を返す。
func (r *questionRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Question{}).
		Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// Vote は質問に投票する。既に投票済みなら値を更新し、質問の投票数も差分だけ増減させる。
// Vote は投票行の作成・更新と vote_count の増減を同一トランザクションで行う。
// 片方だけが反映されて投票行と vote_count がずれた状態を残さない。
func (r *questionRepository) Vote(ctx context.Context, userID, questionID uint, value int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing model.QuestionVote
		err := tx.Where("user_id = ? AND question_id = ?", userID, questionID).First(&existing).Error
		switch {
		case err == nil:
			diff := value - existing.Value
			existing.Value = value
			if err := tx.Save(&existing).Error; err != nil {
				return err
			}
			return tx.Model(&model.Question{}).Where("id = ?", questionID).
				UpdateColumn("vote_count", gorm.Expr("vote_count + ?", diff)).Error
		case errors.Is(err, gorm.ErrRecordNotFound):
			vote := &model.QuestionVote{UserID: userID, QuestionID: questionID, Value: value}
			if err := tx.Create(vote).Error; err != nil {
				return err
			}
			return tx.Model(&model.Question{}).Where("id = ?", questionID).
				UpdateColumn("vote_count", gorm.Expr("vote_count + ?", value)).Error
		default:
			return err
		}
	})
}

// RemoveVote は投票を取り消し、質問の投票数から元の値を差し引く。
// RemoveVote は投票を取り消す。未投票の場合は何もせず成功する（冪等）。
// 削除と vote_count の減算は同一トランザクションで行い、ずれを残さない。
func (r *questionRepository) RemoveVote(ctx context.Context, userID, questionID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing model.QuestionVote
		if err := tx.Where("user_id = ? AND question_id = ?", userID, questionID).First(&existing).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}
		if err := tx.Delete(&existing).Error; err != nil {
			return err
		}
		return tx.Model(&model.Question{}).Where("id = ?", questionID).
			UpdateColumn("vote_count", gorm.Expr("vote_count - ?", existing.Value)).Error
	})
}

// GetUserVote は指定ユーザーの投票値を返す。未投票の場合は 0 を返す。
func (r *questionRepository) GetUserVote(ctx context.Context, userID, questionID uint) (int, error) {
	var vote model.QuestionVote
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND question_id = ?", userID, questionID).First(&vote).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return vote.Value, nil
}

// Bookmark は質問をブックマークする。
func (r *questionRepository) Bookmark(ctx context.Context, userID, questionID uint) error {
	bookmark := &model.QuestionBookmark{UserID: userID, QuestionID: questionID}
	return r.db.WithContext(ctx).Create(bookmark).Error
}

// Unbookmark は質問のブックマークを解除する。
func (r *questionRepository) Unbookmark(ctx context.Context, userID, questionID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND question_id = ?", userID, questionID).
		Delete(&model.QuestionBookmark{}).Error
}

// HasBookmarked は指定ユーザーがブックマーク済みかを返す。
func (r *questionRepository) HasBookmarked(ctx context.Context, userID, questionID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.QuestionBookmark{}).
		Where("user_id = ? AND question_id = ?", userID, questionID).
		Count(&count).Error
	return count > 0, err
}
