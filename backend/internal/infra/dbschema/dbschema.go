// Package dbschema は GORM のタグでは表現できないスキーマ定義を補う。
// AutoMigrate の後に呼び出して、インデックスの形を正しい状態へ揃える。
package dbschema

import "gorm.io/gorm"

// EnsureUserIndexes は users のインデックスを AutoMigrate 後の状態から補正する。
//
// git_hub_id は GitHub 未連携のユーザーがゼロ値 0 を共有するため、全行を対象にした
// ユニークインデックスでは 2 人目の登録が一意制約に衝突する。連携済みの行だけを
// 対象にした部分インデックスへ置き換えることで、未連携ユーザーは何人でも作れて、
// 実在する GitHub ID の重複は引き続き防げる。
//
// 部分インデックスは GORM のタグで表現できないため、ここで明示的に作成する。
func EnsureUserIndexes(db *gorm.DB) error {
	if err := db.Exec(`DROP INDEX IF EXISTS idx_users_git_hub_id`).Error; err != nil {
		return err
	}
	return db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_git_hub_id_linked
		ON users (git_hub_id) WHERE git_hub_id <> 0`).Error
}
