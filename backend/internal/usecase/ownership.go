package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
)

// ensureOwner はエンティティを finder で取得し、userID が所有者であることを検証して返す。
// 取得に失敗した場合は finder のエラーをそのまま返し、不在（finder が nil を返した場合）は
// domain.ErrNotFound、所有者でなければ domain.ErrForbidden を返す。
// 所有権チェックを持つ複数スライスで共有する汎用 helper。
func ensureOwner[T any](
	ctx context.Context,
	finder func(context.Context, uint) (*T, error),
	id, userID uint,
	ownerOf func(*T) uint,
) (*T, error) {
	entity, err := finder(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, domain.ErrNotFound
	}
	if ownerOf(entity) != userID {
		return nil, domain.ErrForbidden
	}
	return entity, nil
}
