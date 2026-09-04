package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// passwordResetRepository は [repository.PasswordResetTokenRepository] の sqlc(pgx) 実装。
type passwordResetRepository struct {
	q *sqlcgen.Queries
}

// NewPasswordResetTokenRepository は PasswordResetTokenRepository の sqlc(pgx) 実装を返す。
func NewPasswordResetTokenRepository(q *sqlcgen.Queries) repository.PasswordResetTokenRepository {
	return &passwordResetRepository{q: q}
}

var _ repository.PasswordResetTokenRepository = (*passwordResetRepository)(nil)

// toModelPasswordResetToken は sqlc の生成行を model.PasswordResetToken へ変換する。
// User は GORM 実装でも Preload していなかったためゼロ値のまま返す。
func toModelPasswordResetToken(row sqlcgen.PasswordResetToken) model.PasswordResetToken {
	return model.PasswordResetToken{
		ID:        uint(row.ID),
		UserID:    uint(row.UserID),
		Token:     row.Token,
		ExpiresAt: row.ExpiresAt.Time,
		Used:      fromBoolPtr(row.Used),
		CreatedAt: row.CreatedAt.Time,
	}
}

// Create はリセットトークンを保存する。
func (r *passwordResetRepository) Create(ctx context.Context, token *model.PasswordResetToken) error {
	row, err := r.q.CreatePasswordResetToken(ctx, sqlcgen.CreatePasswordResetTokenParams{
		UserID:    int64(token.UserID),
		Token:     token.Token,
		ExpiresAt: toTimestamptzNotNull(token.ExpiresAt),
		Used:      &token.Used,
	})
	if err != nil {
		return err
	}
	*token = toModelPasswordResetToken(row)
	return nil
}

// FindByToken はハッシュ済みトークンで検索する。不在の場合は (nil, nil) を返す。
func (r *passwordResetRepository) FindByToken(ctx context.Context, hashedToken string) (*model.PasswordResetToken, error) {
	row, err := r.q.GetPasswordResetTokenByToken(ctx, hashedToken)
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	token := toModelPasswordResetToken(row)
	return &token, nil
}

// MarkAsUsed はトークンを使用済みにする。
func (r *passwordResetRepository) MarkAsUsed(ctx context.Context, id uint) error {
	return r.q.MarkPasswordResetTokenAsUsed(ctx, int64(id))
}

// InvalidateUserTokens は指定ユーザーの未使用トークンをすべて無効化する。
func (r *passwordResetRepository) InvalidateUserTokens(ctx context.Context, userID uint) error {
	return r.q.InvalidateUserPasswordResetTokens(ctx, int64(userID))
}
