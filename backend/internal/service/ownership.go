package service

// findAndCheckOwnership はエンティティを取得し、所有権を検証する汎用ヘルパー関数。
// finder でエンティティを取得し、getOwnerID で所有者IDを取得して userID と比較する。
// 所有者が一致しない場合は ErrForbidden を返す。
func checkOwnership[T any](
	finder func(uint) (*T, error),
	id, userID uint,
	getOwnerID func(*T) uint,
) (*T, error) {
	entity, err := finder(id)
	if err != nil {
		return nil, err
	}
	if getOwnerID(entity) != userID {
		return nil, ErrForbidden
	}
	return entity, nil
}
