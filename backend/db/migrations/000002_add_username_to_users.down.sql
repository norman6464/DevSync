-- ============================================================================
-- usersテーブルからusernameカラムを削除（ロールバック）
-- ============================================================================

-- 1. UNIQUE インデックスを削除
DROP INDEX IF EXISTS idx_users_username;

-- 2. username カラムを削除
ALTER TABLE users DROP COLUMN IF EXISTS username;
