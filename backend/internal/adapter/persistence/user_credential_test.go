package persistence

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserCredentials_CreateWithPassword は、パスワード付きで作成したユーザーだけが
// user_credentials行を持ち、GitHub認証のみを想定した空パスワードのユーザーは
// user_credentials行を持たないことを検証する（DEVSYNC-159でusersから分離）。
func TestUserCredentials_CreateWithPassword(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()
	authRepo := NewAuthUserRepository(pool)

	withPassword := &model.User{
		Username: "uc_with_password",
		Name:     "uc_with_password",
		Email:    "uc_with_password@example.com",
		Password: "hashed-secret",
	}
	require.NoError(t, authRepo.Create(ctx, withPassword))
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, withPassword.ID)
	})
	assert.Equal(t, "hashed-secret", withPassword.Password, "Createの戻り値でもPasswordが保たれているべき")

	var count int64
	require.NoError(t, pool.QueryRow(ctx, `SELECT COUNT(*) FROM user_credentials WHERE user_id = $1`, withPassword.ID).Scan(&count))
	assert.EqualValues(t, 1, count, "パスワードを渡した場合はuser_credentials行が作られるべき")

	withoutPassword := &model.User{
		Username: "uc_without_password",
		Name:     "uc_without_password",
		Email:    "uc_without_password@example.com",
	}
	require.NoError(t, authRepo.Create(ctx, withoutPassword))
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, withoutPassword.ID)
	})

	require.NoError(t, pool.QueryRow(ctx, `SELECT COUNT(*) FROM user_credentials WHERE user_id = $1`, withoutPassword.ID).Scan(&count))
	assert.EqualValues(t, 0, count, "パスワードを渡さなかった場合はuser_credentials行を作らないべき（GitHub認証のみのユーザーを想定）")
}

// TestUserCredentials_FindMethodsReturnPassword は、認証専用の各Findメソッドが
// パスワードハッシュを正しく返すことを検証する。
func TestUserCredentials_FindMethodsReturnPassword(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()
	authRepo := NewAuthUserRepository(pool)

	user := &model.User{
		Username: "uc_find_user",
		Name:     "uc_find_user",
		Email:    "uc_find_user@example.com",
		Password: "hashed-secret",
	}
	require.NoError(t, authRepo.Create(ctx, user))
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, user.ID)
	})

	byEmail, err := authRepo.FindByEmail(ctx, user.Email)
	require.NoError(t, err)
	require.NotNil(t, byEmail)
	assert.Equal(t, "hashed-secret", byEmail.Password)

	byID, err := authRepo.FindByIDWithPassword(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, byID)
	assert.Equal(t, "hashed-secret", byID.Password)

	// 生成系の汎用Read（UserRepository側）はパスワードを持ち込まない。
	genericRepo := NewUserRepository(pool)
	byUsername, err := genericRepo.FindByUsername(ctx, user.Username)
	require.NoError(t, err)
	require.NotNil(t, byUsername)
	assert.Empty(t, byUsername.Password, "汎用のFindByUsernameはパスワードハッシュを持ち込まないべき")
}

// TestUserCredentials_UpdatePasswordUpserts は、user_credentials行が無い状態からでも
// UpdatePasswordが新規作成として機能し、既存行があれば上書きすることを検証する
// （GitHubのみで登録したユーザーが後からパスワードを設定するケース）。
func TestUserCredentials_UpdatePasswordUpserts(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()
	authRepo := NewAuthUserRepository(pool)

	user := &model.User{
		Username: "uc_upsert_user",
		Name:     "uc_upsert_user",
		Email:    "uc_upsert_user@example.com",
	}
	require.NoError(t, authRepo.Create(ctx, user))
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, user.ID)
	})

	var count int64
	require.NoError(t, pool.QueryRow(ctx, `SELECT COUNT(*) FROM user_credentials WHERE user_id = $1`, user.ID).Scan(&count))
	require.EqualValues(t, 0, count)

	require.NoError(t, authRepo.UpdatePassword(ctx, user.ID, "first-hash"))
	found, err := authRepo.FindByIDWithPassword(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "first-hash", found.Password)

	require.NoError(t, authRepo.UpdatePassword(ctx, user.ID, "second-hash"))
	found, err = authRepo.FindByIDWithPassword(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "second-hash", found.Password, "既存のuser_credentials行は上書きされるべき")
}

// TestUserCredentials_UpdateDoesNotTouchPassword は、プロフィール更新（Update）が
// user_credentialsに一切触れないことを検証する。かつてはusersにpassword列を持ち、
// Updateが常に全カラム上書きしていたため、Password未設定のmodel.Userで呼ぶと
// パスワードが消える危険があった。分離後はUpdateの対象自体にできないため構造的に防げる。
func TestUserCredentials_UpdateDoesNotTouchPassword(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()
	authRepo := NewAuthUserRepository(pool)

	user := &model.User{
		Username: "uc_update_user",
		Name:     "uc_update_user",
		Email:    "uc_update_user@example.com",
		Password: "hashed-secret",
	}
	require.NoError(t, authRepo.Create(ctx, user))
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, user.ID)
	})

	// プロフィール更新のusecaseが典型的にそうするように（読み取り→一部フィールドだけ
	// 変更→保存）、作成直後のuserを複製してNameだけ変更する。PasswordはUpdateの対象
	// 外のため、値を持っていても持っていなくても結果に影響しないはずである。
	profileUpdate := *user
	profileUpdate.Name = "更新後の名前"
	require.NoError(t, authRepo.Update(ctx, &profileUpdate))

	found, err := authRepo.FindByIDWithPassword(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "hashed-secret", found.Password, "プロフィール更新はパスワードハッシュを消さないべき")
	assert.Equal(t, "更新後の名前", found.Name)
}

// TestUserCredentials_CascadesOnUserDelete は、ユーザー削除でuser_credentials行も
// FKのON DELETE CASCADEで削除されることを検証する。
func TestUserCredentials_CascadesOnUserDelete(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()
	authRepo := NewAuthUserRepository(pool)

	user := &model.User{
		Username: "uc_delete_user",
		Name:     "uc_delete_user",
		Email:    "uc_delete_user@example.com",
		Password: "hashed-secret",
	}
	require.NoError(t, authRepo.Create(ctx, user))

	require.NoError(t, authRepo.DeleteWithRelatedData(ctx, user.ID))

	var count int64
	require.NoError(t, pool.QueryRow(ctx, `SELECT COUNT(*) FROM user_credentials WHERE user_id = $1`, user.ID).Scan(&count))
	assert.EqualValues(t, 0, count)
}
