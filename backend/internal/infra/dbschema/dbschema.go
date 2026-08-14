// Package dbschema は GORM のタグでは表現できないスキーマ定義を補う。
// AutoMigrate の後に呼び出して、インデックスの形を正しい状態へ揃える。
package dbschema

import (
	"strings"

	"gorm.io/gorm"
)

const (
	// legacyUserGitHubIndex は全行を対象にしていた旧ユニークインデックス。
	// GitHub 未連携のユーザーがゼロ値 0 を共有するため、2 人目の登録が衝突していた。
	legacyUserGitHubIndex = "idx_users_git_hub_id"
	// userGitHubIndex は連携済みの行だけを対象にする部分ユニークインデックス。
	userGitHubIndex = "idx_users_git_hub_id_linked"
	// userGitHubIndexPredicate は pg_indexes が返す定義に含まれる述語。
	userGitHubIndexPredicate = "(git_hub_id <> 0)"
)

// EnsureUserIndexes は users のインデックスを AutoMigrate 後の状態から補正する。
//
// git_hub_id は GitHub 未連携のユーザーがゼロ値 0 を共有するため、全行を対象にした
// ユニークインデックスでは 2 人目の登録が一意制約に衝突する。連携済みの行だけを
// 対象にした部分インデックスへ置き換えることで、未連携ユーザーは何人でも作れて、
// 実在する GitHub ID の重複は引き続き防げる。
//
// 部分インデックスは GORM のタグで表現できないため、ここで明示的に作成する。
// 起動のたびに呼ばれるため冪等にしてあり、張り替えは 1 トランザクションで行う
// （途中で失敗しても一意性の無い状態を残さない）。
func EnsureUserIndexes(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`DROP INDEX IF EXISTS ` + legacyUserGitHubIndex).Error; err != nil {
			return err
		}

		// 同名で定義の違うインデックスがあると CREATE ... IF NOT EXISTS は
		// 中身を確かめずに成功するため、定義を突き合わせて必要なら作り直す。
		var def string
		if err := tx.Raw(
			`SELECT COALESCE(max(indexdef), '') FROM pg_indexes
			 WHERE schemaname = current_schema() AND tablename = 'users' AND indexname = ?`,
			userGitHubIndex,
		).Scan(&def).Error; err != nil {
			return err
		}
		if def != "" && !isExpectedUserGitHubIndex(def) {
			if err := tx.Exec(`DROP INDEX IF EXISTS ` + userGitHubIndex).Error; err != nil {
				return err
			}
		}

		return tx.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS ` + userGitHubIndex + `
			ON users (git_hub_id) WHERE git_hub_id <> 0`).Error
	})
}

// isExpectedUserGitHubIndex は pg_indexes の定義が意図した形かを判定する。
// ユニークであること・git_hub_id を対象にしていること・未連携を除く述語を
// 持つことの 3 点を確認する。
func isExpectedUserGitHubIndex(indexDef string) bool {
	normalized := strings.ToUpper(indexDef)
	return strings.Contains(normalized, "CREATE UNIQUE INDEX") &&
		strings.Contains(normalized, "(GIT_HUB_ID)") &&
		strings.Contains(normalized, strings.ToUpper("WHERE "+userGitHubIndexPredicate))
}
