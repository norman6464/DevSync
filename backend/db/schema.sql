-- ============================================================================
-- DevSync データベーススキーマ定義（schema.sql）
-- PostgreSQL 15 | GORMのAutoMigrateが実際に生成するスキーマと完全一致
-- 生成元: backend/internal/model/*.go（GORMモデル定義）
-- テーブル数: 29 | カラム数: 244
-- ============================================================================

-- ============================================================================
-- 1. users（ユーザーアカウント）
-- 全テーブルの基盤。認証情報・プロフィール・外部サービス連携を保持
-- ============================================================================
CREATE TABLE users (
    id                   BIGSERIAL PRIMARY KEY,
    name                 TEXT      NOT NULL,
    email                TEXT      NOT NULL,
    password             TEXT      NOT NULL,
    avatar_url           TEXT,
    bio                  TEXT,
    git_hub_id           BIGINT,
    git_hub_username     TEXT,
    git_hub_token        TEXT,
    git_hub_connected    BOOLEAN   DEFAULT false,
    zenn_username         TEXT,
    qiita_username        TEXT,
    at_coder_username    TEXT,
    paiza_rank           TEXT,
    skills_languages     TEXT,
    skills_frameworks    TEXT,
    onboarding_completed BOOLEAN   DEFAULT false,
    created_at           TIMESTAMPTZ,
    updated_at           TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_users_email            ON users (email);
CREATE UNIQUE INDEX idx_users_git_hub_id       ON users (git_hub_id);
CREATE UNIQUE INDEX idx_users_git_hub_username ON users (git_hub_username);

-- ============================================================================
-- 2. follows（フォロー関係）
-- ユーザー間のフォロー/フォロワー関係を管理
-- ============================================================================
CREATE TABLE follows (
    id          BIGSERIAL   PRIMARY KEY,
    follower_id BIGINT      NOT NULL,
    followee_id BIGINT      NOT NULL,
    created_at  TIMESTAMPTZ,
    CONSTRAINT fk_follows_follower FOREIGN KEY (follower_id) REFERENCES users (id),
    CONSTRAINT fk_follows_followee FOREIGN KEY (followee_id) REFERENCES users (id)
);

CREATE UNIQUE INDEX idx_follower_following ON follows (follower_id, followee_id);

-- ============================================================================
-- 3. git_hub_contributions（GitHub日別コントリビューション）
-- GitHubの草データを日単位で同期・保存
-- ============================================================================
CREATE TABLE git_hub_contributions (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL,
    date       TIMESTAMPTZ NOT NULL,
    count      BIGINT      NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_user_date ON git_hub_contributions (user_id, date);

-- ============================================================================
-- 4. git_hub_language_stats（GitHub言語別統計）
-- ユーザーのリポジトリ言語統計データ
-- ============================================================================
CREATE TABLE git_hub_language_stats (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL,
    language   TEXT        NOT NULL,
    bytes      BIGINT      NOT NULL DEFAULT 0,
    repo_count BIGINT      NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_user_lang ON git_hub_language_stats (user_id, language);

-- ============================================================================
-- 5. git_hub_repositories（GitHubリポジトリ情報）
-- GitHubから同期したリポジトリのメタデータ
-- ============================================================================
CREATE TABLE git_hub_repositories (
    id              BIGSERIAL   PRIMARY KEY,
    user_id         BIGINT      NOT NULL,
    git_hub_repo_id BIGINT      NOT NULL,
    name            TEXT        NOT NULL,
    full_name       TEXT,
    description     TEXT,
    language        TEXT,
    stars           BIGINT,
    forks           BIGINT,
    is_private      BOOLEAN,
    updated_at      TIMESTAMPTZ
);

CREATE INDEX        idx_git_hub_repositories_user_id         ON git_hub_repositories (user_id);
CREATE UNIQUE INDEX idx_git_hub_repositories_git_hub_repo_id ON git_hub_repositories (git_hub_repo_id);

-- ============================================================================
-- 6. posts（ユーザー投稿）
-- 学習報告や進捗共有の投稿
-- ============================================================================
CREATE TABLE posts (
    id            BIGSERIAL   PRIMARY KEY,
    user_id       BIGINT      NOT NULL,
    title         TEXT        NOT NULL,
    content       TEXT        NOT NULL,
    image_urls    TEXT,
    like_count    BIGINT      DEFAULT 0,
    comment_count BIGINT      DEFAULT 0,
    created_at    TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ,
    CONSTRAINT fk_posts_user FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX idx_posts_user_id ON posts (user_id);

-- ============================================================================
-- 7. likes（投稿へのいいね）
-- ユーザーごとに1投稿1いいねをユニーク制約で保証
-- ============================================================================
CREATE TABLE likes (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL,
    post_id    BIGINT      NOT NULL,
    created_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_user_post_like ON likes (user_id, post_id);
CREATE INDEX        idx_likes_post_id  ON likes (post_id);

-- ============================================================================
-- 8. comments（投稿へのコメント）
-- ============================================================================
CREATE TABLE comments (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL,
    post_id    BIGINT      NOT NULL,
    content    TEXT        NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT fk_comments_user FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX idx_comments_user_id ON comments (user_id);
CREATE INDEX idx_comments_post_id ON comments (post_id);

-- ============================================================================
-- 9. messages（ダイレクトメッセージ）
-- ユーザー間の1対1メッセージ
-- ============================================================================
CREATE TABLE messages (
    id          BIGSERIAL   PRIMARY KEY,
    sender_id   BIGINT      NOT NULL,
    receiver_id BIGINT      NOT NULL,
    content     TEXT        NOT NULL,
    read        BOOLEAN     DEFAULT false,
    created_at  TIMESTAMPTZ,
    CONSTRAINT fk_messages_sender   FOREIGN KEY (sender_id)   REFERENCES users (id),
    CONSTRAINT fk_messages_receiver FOREIGN KEY (receiver_id) REFERENCES users (id)
);

CREATE INDEX idx_messages_sender_id   ON messages (sender_id);
CREATE INDEX idx_messages_receiver_id ON messages (receiver_id);

-- ============================================================================
-- 10. password_reset_tokens（パスワードリセットトークン）
-- 1時間有効のワンタイムトークン
-- ============================================================================
CREATE TABLE password_reset_tokens (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL,
    token      TEXT        NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used       BOOLEAN     DEFAULT false,
    created_at TIMESTAMPTZ,
    CONSTRAINT fk_password_reset_tokens_user FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX        idx_password_reset_tokens_user_id ON password_reset_tokens (user_id);
CREATE UNIQUE INDEX idx_password_reset_tokens_token   ON password_reset_tokens (token);

-- ============================================================================
-- 11. zenn_articles（Zenn記事データ）
-- Zennから同期した記事のメタデータ
-- ============================================================================
CREATE TABLE zenn_articles (
    id             BIGSERIAL   PRIMARY KEY,
    user_id        BIGINT      NOT NULL,
    zenn_id        BIGINT      NOT NULL,
    title          TEXT        NOT NULL,
    slug           TEXT        NOT NULL,
    emoji          TEXT,
    article_type   TEXT,
    liked_count    BIGINT      DEFAULT 0,
    comments_count BIGINT      DEFAULT 0,
    published_at   TIMESTAMPTZ,
    updated_at     TIMESTAMPTZ
);

CREATE INDEX        idx_zenn_articles_user_id ON zenn_articles (user_id);
CREATE UNIQUE INDEX idx_zenn_articles_zenn_id ON zenn_articles (zenn_id);

-- ============================================================================
-- 12. qiita_articles（Qiita記事データ）
-- Qiitaから同期した記事のメタデータ
-- ============================================================================
CREATE TABLE qiita_articles (
    id             BIGSERIAL   PRIMARY KEY,
    user_id        BIGINT      NOT NULL,
    qiita_id       TEXT        NOT NULL,
    title          TEXT        NOT NULL,
    url            TEXT        NOT NULL,
    likes_count    BIGINT      DEFAULT 0,
    comments_count BIGINT      DEFAULT 0,
    tags           TEXT,
    published_at   TIMESTAMPTZ,
    updated_at     TIMESTAMPTZ
);

CREATE INDEX        idx_qiita_articles_user_id  ON qiita_articles (user_id);
CREATE UNIQUE INDEX idx_qiita_articles_qiita_id ON qiita_articles (qiita_id);

-- ============================================================================
-- 13. learning_goals（学習目標）
-- 進捗率0〜100%、ステータス: active/completed/paused
-- ============================================================================
CREATE TABLE learning_goals (
    id           BIGSERIAL   PRIMARY KEY,
    user_id      BIGINT      NOT NULL,
    title        TEXT        NOT NULL,
    description  TEXT,
    category     TEXT        DEFAULT 'other'::text,
    target_date  TIMESTAMPTZ,
    progress     BIGINT      DEFAULT 0,
    status       TEXT        DEFAULT 'active'::text,
    created_at   TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ,
    completed_at TIMESTAMPTZ
);

CREATE INDEX idx_learning_goals_user_id ON learning_goals (user_id);

-- ============================================================================
-- 14. learning_logs（学習記録）
-- 日々の学習活動ログ。カテゴリ: coding/reading/course/meetup/other
-- ============================================================================
CREATE TABLE learning_logs (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL,
    title      TEXT        NOT NULL,
    content    TEXT        NOT NULL,
    category   TEXT        DEFAULT 'other'::text,
    duration   BIGINT      DEFAULT 0,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE INDEX idx_learning_logs_user_id ON learning_logs (user_id);

-- ============================================================================
-- 15. projects（プロジェクトショーケース）★論理削除対応
-- ユーザーのプロジェクト作品を展示。GitHubリポジトリとの紐付け可能
-- ============================================================================
CREATE TABLE projects (
    id             BIGSERIAL    PRIMARY KEY,
    user_id        BIGINT       NOT NULL,
    title          VARCHAR(200) NOT NULL,
    description    TEXT,
    tech_stack     TEXT,
    demo_url       VARCHAR(500),
    github_url     VARCHAR(500),
    image_url      VARCHAR(500),
    role           VARCHAR(100),
    start_date     TIMESTAMPTZ,
    end_date       TIMESTAMPTZ,
    featured       BOOLEAN      DEFAULT false,
    github_repo_id BIGINT,
    created_at     TIMESTAMPTZ,
    updated_at     TIMESTAMPTZ,
    deleted_at     TIMESTAMPTZ,
    CONSTRAINT fk_projects_user        FOREIGN KEY (user_id)        REFERENCES users (id),
    CONSTRAINT fk_projects_github_repo FOREIGN KEY (github_repo_id) REFERENCES git_hub_repositories (id)
);

CREATE INDEX idx_projects_user_id    ON projects (user_id);
CREATE INDEX idx_projects_deleted_at ON projects (deleted_at);

-- ============================================================================
-- 16. learning_resources（学習リソース）★論理削除対応
-- 書籍・動画・記事等の学習リソース共有。公開/非公開切り替え可能
-- ============================================================================
CREATE TABLE learning_resources (
    id          BIGSERIAL    PRIMARY KEY,
    user_id     BIGINT       NOT NULL,
    title       VARCHAR(300) NOT NULL,
    description TEXT,
    url         VARCHAR(500),
    category    VARCHAR(50)  NOT NULL,
    difficulty  VARCHAR(50),
    tags        TEXT,
    image_url   VARCHAR(500),
    is_public   BOOLEAN      DEFAULT true,
    like_count  BIGINT       DEFAULT 0,
    save_count  BIGINT       DEFAULT 0,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,
    CONSTRAINT fk_learning_resources_user FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX idx_learning_resources_user_id    ON learning_resources (user_id);
CREATE INDEX idx_learning_resources_deleted_at ON learning_resources (deleted_at);

-- ============================================================================
-- 17. resource_likes（学習リソースへのいいね）
-- ============================================================================
CREATE TABLE resource_likes (
    id          BIGSERIAL   PRIMARY KEY,
    user_id     BIGINT      NOT NULL,
    resource_id BIGINT      NOT NULL,
    created_at  TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_resource_like ON resource_likes (user_id, resource_id);

-- ============================================================================
-- 18. resource_saves（学習リソースのブックマーク）
-- ============================================================================
CREATE TABLE resource_saves (
    id          BIGSERIAL   PRIMARY KEY,
    user_id     BIGINT      NOT NULL,
    resource_id BIGINT      NOT NULL,
    created_at  TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_resource_save ON resource_saves (user_id, resource_id);

-- ============================================================================
-- 19. book_reviews（書籍レビュー）★論理削除対応
-- 技術書のレビュー・評価（1〜5段階）
-- ============================================================================
CREATE TABLE book_reviews (
    id         BIGSERIAL    PRIMARY KEY,
    user_id    BIGINT       NOT NULL,
    title      VARCHAR(300) NOT NULL,
    author     VARCHAR(200),
    isbn       VARCHAR(20),
    rating     BIGINT       NOT NULL,
    review     TEXT,
    image_url  VARCHAR(500),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_book_reviews_user FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX idx_book_reviews_user_id    ON book_reviews (user_id);
CREATE INDEX idx_book_reviews_deleted_at ON book_reviews (deleted_at);

-- ============================================================================
-- 20. questions（Q&A質問）★論理削除対応
-- 質問投稿。投票数・回答数・解決済みフラグを保持
-- ============================================================================
CREATE TABLE questions (
    id           BIGSERIAL    PRIMARY KEY,
    user_id      BIGINT       NOT NULL,
    title        VARCHAR(500) NOT NULL,
    body         TEXT         NOT NULL,
    tags         TEXT,
    vote_count   BIGINT       DEFAULT 0,
    answer_count BIGINT       DEFAULT 0,
    is_solved    BOOLEAN      DEFAULT false,
    created_at   TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ,
    deleted_at   TIMESTAMPTZ,
    CONSTRAINT fk_questions_user FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX idx_questions_user_id    ON questions (user_id);
CREATE INDEX idx_questions_deleted_at ON questions (deleted_at);

-- ============================================================================
-- 21. question_votes（質問への投票）
-- value: +1（賛成）または -1（反対）
-- ============================================================================
CREATE TABLE question_votes (
    id          BIGSERIAL   PRIMARY KEY,
    user_id     BIGINT      NOT NULL,
    question_id BIGINT      NOT NULL,
    value       BIGINT      NOT NULL,
    created_at  TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_question_vote ON question_votes (user_id, question_id);

-- ============================================================================
-- 22. answers（Q&A回答）★論理削除対応
-- 質問への回答。ベストアンサーフラグで最良回答を明示
-- ============================================================================
CREATE TABLE answers (
    id          BIGSERIAL   PRIMARY KEY,
    user_id     BIGINT      NOT NULL,
    question_id BIGINT      NOT NULL,
    body        TEXT        NOT NULL,
    vote_count  BIGINT      DEFAULT 0,
    is_best     BOOLEAN     DEFAULT false,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,
    CONSTRAINT fk_answers_user FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX idx_answers_user_id     ON answers (user_id);
CREATE INDEX idx_answers_question_id ON answers (question_id);
CREATE INDEX idx_answers_deleted_at  ON answers (deleted_at);

-- ============================================================================
-- 23. answer_votes（回答への投票）
-- value: +1（賛成）または -1（反対）
-- ============================================================================
CREATE TABLE answer_votes (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL,
    answer_id  BIGINT      NOT NULL,
    value      BIGINT      NOT NULL,
    created_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_answer_vote ON answer_votes (user_id, answer_id);

-- ============================================================================
-- 24. roadmaps（学習ロードマップ）
-- 学習計画。公開/非公開切り替え、進捗率の自動計算
-- ============================================================================
CREATE TABLE roadmaps (
    id                   BIGSERIAL    PRIMARY KEY,
    user_id              BIGINT       NOT NULL,
    title                VARCHAR(200) NOT NULL,
    description          TEXT,
    category             TEXT         DEFAULT 'other'::text,
    is_public            BOOLEAN      DEFAULT false,
    is_template          BOOLEAN      DEFAULT false,
    step_count           BIGINT       DEFAULT 0,
    completed_step_count BIGINT       DEFAULT 0,
    progress             BIGINT       DEFAULT 0,
    status               TEXT         DEFAULT 'active'::text,
    created_at           TIMESTAMPTZ,
    updated_at           TIMESTAMPTZ,
    completed_at         TIMESTAMPTZ,
    CONSTRAINT fk_roadmaps_user FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX idx_roadmaps_user_id   ON roadmaps (user_id);
CREATE INDEX idx_roadmaps_is_public ON roadmaps (is_public);

-- ============================================================================
-- 25. roadmap_steps（ロードマップステップ）
-- ロードマップ内の個別学習ステップ。親削除時にCASCADE削除
-- ============================================================================
CREATE TABLE roadmap_steps (
    id           BIGSERIAL    PRIMARY KEY,
    roadmap_id   BIGINT       NOT NULL,
    title        VARCHAR(200) NOT NULL,
    description  TEXT,
    order_index  BIGINT       NOT NULL DEFAULT 0,
    is_completed BOOLEAN      DEFAULT false,
    completed_at TIMESTAMPTZ,
    resource_url VARCHAR(500),
    created_at   TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ,
    CONSTRAINT fk_roadmaps_steps FOREIGN KEY (roadmap_id) REFERENCES roadmaps (id) ON DELETE CASCADE
);

CREATE INDEX idx_roadmap_steps_roadmap_id ON roadmap_steps (roadmap_id);

-- ============================================================================
-- 26. notifications（通知）
-- いいね・コメント・フォロー・回答・バッジ等の通知
-- ============================================================================
CREATE TABLE notifications (
    id          BIGSERIAL   PRIMARY KEY,
    user_id     BIGINT      NOT NULL,
    type        TEXT        NOT NULL,
    actor_id    BIGINT      NOT NULL,
    post_id     BIGINT,
    question_id BIGINT,
    badge_id    VARCHAR(50),
    read        BOOLEAN     DEFAULT false,
    created_at  TIMESTAMPTZ,
    CONSTRAINT fk_notifications_user     FOREIGN KEY (user_id)     REFERENCES users (id),
    CONSTRAINT fk_notifications_actor    FOREIGN KEY (actor_id)    REFERENCES users (id),
    CONSTRAINT fk_notifications_post     FOREIGN KEY (post_id)     REFERENCES posts (id),
    CONSTRAINT fk_notifications_question FOREIGN KEY (question_id) REFERENCES questions (id)
);

CREATE INDEX idx_notifications_user_id     ON notifications (user_id);
CREATE INDEX idx_notifications_post_id     ON notifications (post_id);
CREATE INDEX idx_notifications_question_id ON notifications (question_id);

-- ============================================================================
-- 27. chat_rooms（グループチャットルーム）
-- ============================================================================
CREATE TABLE chat_rooms (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    owner_id    BIGINT       NOT NULL,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    CONSTRAINT fk_chat_rooms_owner FOREIGN KEY (owner_id) REFERENCES users (id)
);

CREATE INDEX idx_chat_rooms_owner_id ON chat_rooms (owner_id);

-- ============================================================================
-- 28. chat_room_members（チャットルームメンバー）
-- ルームとユーザーの多対多リレーション
-- ============================================================================
CREATE TABLE chat_room_members (
    id           BIGSERIAL   PRIMARY KEY,
    chat_room_id BIGINT      NOT NULL,
    user_id      BIGINT      NOT NULL,
    joined_at    TIMESTAMPTZ,
    CONSTRAINT fk_chat_room_members_chat_room FOREIGN KEY (chat_room_id) REFERENCES chat_rooms (id),
    CONSTRAINT fk_chat_room_members_user      FOREIGN KEY (user_id)      REFERENCES users (id)
);

CREATE INDEX        idx_chat_room_members_chat_room_id ON chat_room_members (chat_room_id);
CREATE INDEX        idx_chat_room_members_user_id      ON chat_room_members (user_id);
CREATE UNIQUE INDEX idx_room_user                      ON chat_room_members (chat_room_id, user_id);

-- ============================================================================
-- 29. group_messages（グループチャットメッセージ）
-- ============================================================================
CREATE TABLE group_messages (
    id           BIGSERIAL   PRIMARY KEY,
    chat_room_id BIGINT      NOT NULL,
    sender_id    BIGINT      NOT NULL,
    content      TEXT        NOT NULL,
    created_at   TIMESTAMPTZ,
    CONSTRAINT fk_group_messages_chat_room FOREIGN KEY (chat_room_id) REFERENCES chat_rooms (id),
    CONSTRAINT fk_group_messages_sender    FOREIGN KEY (sender_id)    REFERENCES users (id)
);

CREATE INDEX idx_group_messages_chat_room_id ON group_messages (chat_room_id);
CREATE INDEX idx_group_messages_sender_id    ON group_messages (sender_id);
