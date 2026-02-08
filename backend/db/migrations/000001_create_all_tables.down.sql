-- ============================================================================
-- DevSync データベーススキーマ ロールバック
-- 全29テーブルのDROP TABLE文（依存関係の逆順）
-- ============================================================================

-- グループチャット関連（29→27）
DROP TABLE IF EXISTS group_messages;
DROP TABLE IF EXISTS chat_room_members;
DROP TABLE IF EXISTS chat_rooms;

-- 通知（26）
DROP TABLE IF EXISTS notifications;

-- ロードマップ関連（25→24）
DROP TABLE IF EXISTS roadmap_steps;
DROP TABLE IF EXISTS roadmaps;

-- Q&A関連（23→20）
DROP TABLE IF EXISTS answer_votes;
DROP TABLE IF EXISTS answers;
DROP TABLE IF EXISTS question_votes;
DROP TABLE IF EXISTS questions;

-- 書籍レビュー（19）
DROP TABLE IF EXISTS book_reviews;

-- 学習リソース関連（18→16）
DROP TABLE IF EXISTS resource_saves;
DROP TABLE IF EXISTS resource_likes;
DROP TABLE IF EXISTS learning_resources;

-- プロジェクト（15）
DROP TABLE IF EXISTS projects;

-- 学習記録・目標（14→13）
DROP TABLE IF EXISTS learning_logs;
DROP TABLE IF EXISTS learning_goals;

-- 外部サービス連携（12→11）
DROP TABLE IF EXISTS qiita_articles;
DROP TABLE IF EXISTS zenn_articles;

-- パスワードリセット（10）
DROP TABLE IF EXISTS password_reset_tokens;

-- メッセージ（9）
DROP TABLE IF EXISTS messages;

-- 投稿関連（8→6）
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS likes;
DROP TABLE IF EXISTS posts;

-- GitHub関連（5→3）
DROP TABLE IF EXISTS github_repositories;
DROP TABLE IF EXISTS github_language_stats;
DROP TABLE IF EXISTS github_contributions;

-- フォロー（2）
DROP TABLE IF EXISTS follows;

-- ユーザー（1 — 最後に削除）
DROP TABLE IF EXISTS users;
