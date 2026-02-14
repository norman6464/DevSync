-- ============================================================================
-- usersテーブルにusernameカラムを追加
-- プロフィールURL用の一意のユーザー名フィールド
-- ============================================================================

-- 1. username カラムを追加（一時的にNULLを許可）
ALTER TABLE users ADD COLUMN IF NOT EXISTS username TEXT;

-- 2. 既存ユーザーにデフォルトのusernameを生成（user{id}形式）
UPDATE users
SET username = 'user' || id::TEXT
WHERE username IS NULL OR username = '';

-- 3. username カラムに NOT NULL 制約を追加
ALTER TABLE users ALTER COLUMN username SET NOT NULL;

-- 4. username カラムに UNIQUE インデックスを追加
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users (username);
