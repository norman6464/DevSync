package usecase

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/domain"
)

// errOwnedEntityNotFound は finder が「不在」を表す nil を返したときに ensureOwner が返すエラー。
// DomainError ではないため handler では 500 になり、不在を素の DB エラーとして扱っていた
// 移行前の挙動と一致する。
var errOwnedEntityNotFound = errors.New("対象が見つかりません")

// ensureOwner はエンティティを finder で取得し、userID が所有者であることを検証して返す。
// 取得に失敗した場合は finder のエラーをそのまま返し、所有者でなければ domain.ErrForbidden を返す。
// finder が不在を (nil, nil) で表す契約の場合に備え、nil を受け取ったら errOwnedEntityNotFound を返す
// （ownerOf に nil を渡して落ちるのを防ぐ）。
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
		return nil, errOwnedEntityNotFound
	}
	if ownerOf(entity) != userID {
		return nil, domain.ErrForbidden
	}
	return entity, nil
}
