package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// msgInvalidAtCoderUsername は連携するユーザー名の検証に失敗したときのメッセージ。
const msgInvalidAtCoderUsername = "invalid AtCoder username"

// msgLinkedUserNotFound は連携対象のユーザーが見つからないときのメッセージ。
const msgLinkedUserNotFound = "user not found"

// GetAtCoderRatingUseCase は AtCoder のレーティング情報を取得する。
type GetAtCoderRatingUseCase struct {
	ratings repository.AtCoderRatingFetcher
}

// NewGetAtCoderRatingUseCase は GetAtCoderRatingUseCase を生成する。
func NewGetAtCoderRatingUseCase(ratings repository.AtCoderRatingFetcher) *GetAtCoderRatingUseCase {
	return &GetAtCoderRatingUseCase{ratings: ratings}
}

// Execute はユーザー名を検証し、最新のレーティングから色とランクを求めて返す。
// 履歴が 1 件も無い場合はレーティング 0（灰）として返す。
func (uc *GetAtCoderRatingUseCase) Execute(ctx context.Context, username string) (*model.AtCoderRatingInfo, error) {
	if err := domain.ValidateExternalUsername(username); err != nil {
		return nil, err
	}

	history, err := uc.ratings.FetchRatingHistory(ctx, username)
	if err != nil {
		return nil, err
	}

	info := &model.AtCoderRatingInfo{
		Username: username,
		Rating:   0,
		Color:    "gray",
		Rank:     "灰",
	}
	if len(history) > 0 {
		info.Rating = history[len(history)-1].NewRating
		info.Color = AtCoderRatingColor(info.Rating)
		info.Rank = AtCoderRatingRank(info.Rating)
	}
	return info, nil
}

// AtCoderRatingColor はレーティング値を AtCoder の色名に変換する。
func AtCoderRatingColor(rating int) string {
	switch {
	case rating >= 2800:
		return "red"
	case rating >= 2400:
		return "orange"
	case rating >= 2000:
		return "yellow"
	case rating >= 1600:
		return "blue"
	case rating >= 1200:
		return "cyan"
	case rating >= 800:
		return "green"
	case rating >= 400:
		return "brown"
	default:
		return "gray"
	}
}

// AtCoderRatingRank はレーティング値を AtCoder のランク名に変換する。
func AtCoderRatingRank(rating int) string {
	switch {
	case rating >= 2800:
		return "赤"
	case rating >= 2400:
		return "橙"
	case rating >= 2000:
		return "黄"
	case rating >= 1600:
		return "青"
	case rating >= 1200:
		return "水色"
	case rating >= 800:
		return "緑"
	case rating >= 400:
		return "茶"
	default:
		return "灰"
	}
}

// ConnectAtCoderUseCase は AtCoder のユーザー名をプロフィールへ保存する。
type ConnectAtCoderUseCase struct {
	users   repository.ExternalAccountLinker
	ratings repository.AtCoderRatingFetcher
}

// NewConnectAtCoderUseCase は ConnectAtCoderUseCase を生成する。
func NewConnectAtCoderUseCase(
	users repository.ExternalAccountLinker,
	ratings repository.AtCoderRatingFetcher,
) *ConnectAtCoderUseCase {
	return &ConnectAtCoderUseCase{users: users, ratings: ratings}
}

// Execute はユーザー名の形式と AtCoder 上の存在を確認し、プロフィールへ保存する。
func (uc *ConnectAtCoderUseCase) Execute(ctx context.Context, userID uint, username string) (*model.User, error) {
	if domain.ValidateExternalUsername(username) != nil || !uc.ratings.UserExists(ctx, username) {
		return nil, domain.NewError(domain.ErrCodeBadRequest, msgInvalidAtCoderUsername, nil)
	}

	user, err := findLinkedUser(ctx, uc.users, userID)
	if err != nil {
		return nil, err
	}

	user.AtCoderUsername = username
	if err := uc.users.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// DisconnectAtCoderUseCase は AtCoder のユーザー名をプロフィールから消す。
type DisconnectAtCoderUseCase struct {
	users repository.ExternalAccountLinker
}

// NewDisconnectAtCoderUseCase は DisconnectAtCoderUseCase を生成する。
func NewDisconnectAtCoderUseCase(users repository.ExternalAccountLinker) *DisconnectAtCoderUseCase {
	return &DisconnectAtCoderUseCase{users: users}
}

// Execute は AtCoder のユーザー名を空にして保存する。
func (uc *DisconnectAtCoderUseCase) Execute(ctx context.Context, userID uint) (*model.User, error) {
	user, err := findLinkedUser(ctx, uc.users, userID)
	if err != nil {
		return nil, err
	}

	user.AtCoderUsername = ""
	if err := uc.users.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// findLinkedUser は外部サービス連携の対象ユーザーを取得する。
// 不在・取得失敗はどちらも 404 として扱う（移行前の挙動）。
func findLinkedUser(ctx context.Context, users repository.UserSkillsReader, userID uint) (*model.User, error) {
	user, err := users.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, msgLinkedUserNotFound, err)
	}
	return user, nil
}
