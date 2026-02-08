-- ============================================================================
-- DevSync データベーススキーマ定義
-- 全29テーブルのCREATE TABLE文（PostgreSQL）
-- GORMモデル定義（backend/internal/model/）から忠実に変換
-- ============================================================================

-- ============================================================================
-- 1. users — ユーザーアカウント情報（全テーブルの基盤）
-- ============================================================================
CREATE TABLE IF NOT EXISTS users (
    id                   BIGSERIAL    PRIMARY KEY,
    name                 TEXT         NOT NULL,
    email                TEXT         NOT NULL,
    password             TEXT,
    avatar_url           TEXT,
    bio                  TEXT,
    git_hub_id           BIGINT,
    git_hub_username     TEXT,
    git_hub_token        TEXT,
    git_hub_connected    BOOLEAN      NOT NULL DEFAULT FALSE,
    zenn_username         TEXT,
    qiita_username        TEXT,
    skills_languages     TEXT,
    skills_frameworks    TEXT,
    onboarding_completed BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at           TIMESTAMPTZ,
    updated_at           TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_git_hub_id ON users (git_hub_id);

-- ============================================================================
-- 2. follows — ユーザー間のフォロー関係
-- ============================================================================
CREATE TABLE IF NOT EXISTS follows (
    id          BIGSERIAL    PRIMARY KEY,
    follower_id BIGINT       NOT NULL,
    followee_id BIGINT       NOT NULL,
    created_at  TIMESTAMPTZ,
    CONSTRAINT fk_follows_follower FOREIGN KEY (follower_id) REFERENCES users (id),
    CONSTRAINT fk_follows_followee FOREIGN KEY (followee_id) REFERENCES users (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_follower_following ON follows (follower_id, followee_id);

-- ============================================================================
-- 3. github_contributions — GitHub日別コントリビューション（草）データ
-- ============================================================================
CREATE TABLE IF NOT EXISTS github_contributions (
    id         BIGSERIAL    PRIMARY KEY,
    user_id    BIGINT       NOT NULL,
    date       TIMESTAMPTZ  NOT NULL,
    count      INTEGER      NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_date ON github_contributions (user_id, date);

-- ============================================================================
-- 4. github_language_stats — GitHub言語別統計データ
-- ============================================================================
CREATE TABLE IF NOT EXISTS github_language_stats (
    id         BIGSERIAL    PRIMARY KEY,
    user_id    BIGINT       NOT NULL,
    language   TEXT         NOT NULL,
    bytes      BIGINT       NOT NULL DEFAULT 0,
    repo_count INTEGER      NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_lang ON github_language_stats (user_id, language);

-- ============================================================================
-- 5. github_repositories — GitHubリポジトリ情報
-- ============================================================================
CREATE TABLE IF NOT EXISTS github_repositories (
    id              BIGSERIAL    PRIMARY KEY,
    user_id         BIGINT       NOT NULL,
    git_hub_repo_id BIGINT       NOT NULL,
    name            TEXT         NOT NULL,
    full_name       TEXT,
    description     TEXT,
    language        TEXT,
    stars           INTEGER,
    forks           INTEGER,
    is_private      BOOLEAN,
    updated_at      TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_github_repositories_user_id ON github_repositories (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_github_repositories_git_hub_repo_id ON github_repositories (git_hub_repo_id);

-- ============================================================================
-- 6. posts — ユーザー投稿
-- ============================================================================
CREATE TABLE IF NOT EXISTS posts (
    id            BIGSERIAL    PRIMARY KEY,
    user_id       BIGINT       NOT NULL,
    title         TEXT         NOT NULL,
    content       TEXT         NOT NULL,
    image_urls    TEXT,
    like_count    INTEGER      NOT NULL DEFAULT 0,
    comment_count INTEGER      NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ,
    CONSTRAINT fk_posts_user FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id);

-- ============================================================================
-- 7. likes — 投稿へのいいね
-- ============================================================================
CREATE TABLE IF NOT EXISTS likes (
    id         BIGSERIAL    PRIMARY KEY,
    user_id    BIGINT       NOT NULL,
    post_id    BIGINT       NOT NULL,
    created_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_post_like ON likes (user_id, post_id);
CREATE INDEX IF NOT EXISTS idx_likes_post_id ON likes (post_id);

-- ============================================================================
-- 8. comments — 投稿へのコメント
-- ============================================================================
CREATE TABLE IF NOT EXISTS comments (
    id         BIGSERIAL    PRIMARY KEY,
    user_id    BIGINT       NOT NULL,
    post_id    BIGINT       NOT NULL,
    content    TEXT         NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT fk_comments_user FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_comments_post FOREIGN KEY (post_id) REFERENCES posts (id)
);

CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments (user_id);
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments (post_id);

-- ============================================================================
-- 9. messages — ダイレクトメッセージ
-- ============================================================================
CREATE TABLE IF NOT EXISTS messages (
    id          BIGSERIAL    PRIMARY KEY,
    sender_id   BIGINT       NOT NULL,
    receiver_id BIGINT       NOT NULL,
    content     TEXT         NOT NULL,
    read        BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ,
    CONSTRAINT fk_messages_sender FOREIGN KEY (sender_id) REFERENCES users (id),
    CONSTRAINT fk_messages_receiver FOREIGN KEY (receiver_id) REFERENCES users (id)
);

CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON messages (sender_id);
CREATE INDEX IF NOT EXISTS idx_messages_receiver_id ON messages (receiver_id);

-- ============================================================================
-- 10. password_reset_tokens — パスワードリセット用トークン
-- ============================================================================
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id         BIGSERIAL    PRIMARY KEY,
    user_id    BIGINT       NOT NULL,
    token      TEXT         NOT NULL,
    expires_at TIMESTAMPTZ  NOT NULL,
    used       BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ,
    CONSTRAINT fk_password_reset_tokens_user FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_user_id ON password_reset_tokens (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_password_reset_tokens_token ON password_reset_tokens (token);

-- ============================================================================
-- 11. zenn_articles — Zenn同期記事データ
-- ============================================================================
CREATE TABLE IF NOT EXISTS zenn_articles (
    id             BIGSERIAL    PRIMARY KEY,
    user_id        BIGINT       NOT NULL,
    zenn_id        BIGINT       NOT NULL,
    title          TEXT         NOT NULL,
    slug           TEXT         NOT NULL,
    emoji          TEXT,
    article_type   TEXT,
    liked_count    INTEGER      NOT NULL DEFAULT 0,
    comments_count INTEGER      NOT NULL DEFAULT 0,
    published_at   TIMESTAMPTZ,
    updated_at     TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_zenn_articles_user_id ON zenn_articles (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_zenn_articles_zenn_id ON zenn_articles (zenn_id);

-- ============================================================================
-- 12. qiita_articles — Qiita同期記事データ
-- ============================================================================
CREATE TABLE IF NOT EXISTS qiita_articles (
    id             BIGSERIAL    PRIMARY KEY,
    user_id        BIGINT       NOT NULL,
    qiita_id       TEXT         NOT NULL,
    title          TEXT         NOT NULL,
    url            TEXT         NOT NULL,
    likes_count    INTEGER      NOT NULL DEFAULT 0,
    comments_count INTEGER      NOT NULL DEFAULT 0,
    tags           TEXT,
    published_at   TIMESTAMPTZ,
    updated_at     TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_qiita_articles_user_id ON qiita_articles (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_qiita_articles_qiita_id ON qiita_articles (qiita_id);

-- ============================================================================
-- 13. learning_goals — 学習目標
-- ============================================================================
CREATE TABLE IF NOT EXISTS learning_goals (
    id           BIGSERIAL    PRIMARY KEY,
    user_id      BIGINT       NOT NULL,
    title        TEXT         NOT NULL,
    description  TEXT,
    category     TEXT         DEFAULT 'other',
    target_date  TIMESTAMPTZ,
    progress     INTEGER      NOT NULL DEFAULT 0,
    status       TEXT         DEFAULT 'active',
    created_at   TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ,
    completed_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_learning_goals_user_id ON learning_goals (user_id);

-- ============================================================================
-- 14. learning_logs — 日々の学習記録
-- ============================================================================
CREATE TABLE IF NOT EXISTS learning_logs (
    id         BIGSERIAL    PRIMARY KEY,
    user_id    BIGINT       NOT NULL,
    title      TEXT         NOT NULL,
    content    TEXT         NOT NULL,
    category   TEXT         DEFAULT 'other',
    duration   INTEGER      NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_learning_logs_user_id ON learning_logs (user_id);

-- ============================================================================
-- 15. projects — プロジェクトショーケース（論理削除対応）
-- ============================================================================
CREATE TABLE IF NOT EXISTS projects (
    id             BIGSERIAL      PRIMARY KEY,
    user_id        BIGINT         NOT NULL,
    title          VARCHAR(200)   NOT NULL,
    description    TEXT,
    tech_stack     TEXT,
    demo_url       VARCHAR(500),
    github_url     VARCHAR(500),
    image_url      VARCHAR(500),
    role           VARCHAR(100),
    start_date     TIMESTAMPTZ,
    end_date       TIMESTAMPTZ,
    featured       BOOLEAN        NOT NULL DEFAULT FALSE,
    github_repo_id BIGINT,
    created_at     TIMESTAMPTZ,
    updated_at     TIMESTAMPTZ,
    deleted_at     TIMESTAMPTZ,
    CONSTRAINT fk_projects_user FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_projects_github_repo FOREIGN KEY (github_repo_id) REFERENCES github_repositories (id)
);

CREATE INDEX IF NOT EXISTS idx_projects_user_id ON projects (user_id);
CREATE INDEX IF NOT EXISTS idx_projects_deleted_at ON projects (deleted_at);

-- ============================================================================
-- 16. learning_resources — 学習リソース（論理削除対応）
-- ============================================================================
CREATE TABLE IF NOT EXISTS learning_resources (
    id          BIGSERIAL      PRIMARY KEY,
    user_id     BIGINT         NOT NULL,
    title       VARCHAR(300)   NOT NULL,
    description TEXT,
    url         VARCHAR(500),
    category    VARCHAR(50)    NOT NULL,
    difficulty  VARCHAR(50),
    tags        TEXT,
    image_url   VARCHAR(500),
    is_public   BOOLEAN        NOT NULL DEFAULT TRUE,
    like_count  INTEGER        NOT NULL DEFAULT 0,
    save_count  INTEGER        NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,
    CONSTRAINT fk_learning_resources_user FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX IF NOT EXISTS idx_learning_resources_user_id ON learning_resources (user_id);
CREATE INDEX IF NOT EXISTS idx_learning_resources_deleted_at ON learning_resources (deleted_at);

-- ============================================================================
-- 17. resource_likes — 学習リソースへのいいね
-- ============================================================================
CREATE TABLE IF NOT EXISTS resource_likes (
    id          BIGSERIAL    PRIMARY KEY,
    user_id     BIGINT       NOT NULL,
    resource_id BIGINT       NOT NULL,
    created_at  TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_resource_like ON resource_likes (user_id, resource_id);

-- ============================================================================
-- 18. resource_saves — 学習リソースのブックマーク
-- ============================================================================
CREATE TABLE IF NOT EXISTS resource_saves (
    id          BIGSERIAL    PRIMARY KEY,
    user_id     BIGINT       NOT NULL,
    resource_id BIGINT       NOT NULL,
    created_at  TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_resource_save ON resource_saves (user_id, resource_id);

-- ============================================================================
-- 19. book_reviews — 書籍レビュー（論理削除対応）
-- ============================================================================
CREATE TABLE IF NOT EXISTS book_reviews (
    id         BIGSERIAL      PRIMARY KEY,
    user_id    BIGINT         NOT NULL,
    title      VARCHAR(300)   NOT NULL,
    author     VARCHAR(200),
    isbn       VARCHAR(20),
    rating     INTEGER        NOT NULL,
    review     TEXT,
    image_url  VARCHAR(500),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_book_reviews_user FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX IF NOT EXISTS idx_book_reviews_user_id ON book_reviews (user_id);
CREATE INDEX IF NOT EXISTS idx_book_reviews_deleted_at ON book_reviews (deleted_at);

-- ============================================================================
-- 20. questions — Q&A質問（論理削除対応）
-- ============================================================================
CREATE TABLE IF NOT EXISTS questions (
    id           BIGSERIAL      PRIMARY KEY,
    user_id      BIGINT         NOT NULL,
    title        VARCHAR(500)   NOT NULL,
    body         TEXT           NOT NULL,
    tags         TEXT,
    vote_count   INTEGER        NOT NULL DEFAULT 0,
    answer_count INTEGER        NOT NULL DEFAULT 0,
    is_solved    BOOLEAN        NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ,
    deleted_at   TIMESTAMPTZ,
    CONSTRAINT fk_questions_user FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX IF NOT EXISTS idx_questions_user_id ON questions (user_id);
CREATE INDEX IF NOT EXISTS idx_questions_deleted_at ON questions (deleted_at);

-- ============================================================================
-- 21. question_votes — 質問への投票
-- ============================================================================
CREATE TABLE IF NOT EXISTS question_votes (
    id          BIGSERIAL    PRIMARY KEY,
    user_id     BIGINT       NOT NULL,
    question_id BIGINT       NOT NULL,
    value       INTEGER      NOT NULL,
    created_at  TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_question_vote ON question_votes (user_id, question_id);

-- ============================================================================
-- 22. answers — Q&A回答（論理削除対応）
-- ============================================================================
CREATE TABLE IF NOT EXISTS answers (
    id          BIGSERIAL      PRIMARY KEY,
    user_id     BIGINT         NOT NULL,
    question_id BIGINT         NOT NULL,
    body        TEXT           NOT NULL,
    vote_count  INTEGER        NOT NULL DEFAULT 0,
    is_best     BOOLEAN        NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,
    CONSTRAINT fk_answers_user FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_answers_question FOREIGN KEY (question_id) REFERENCES questions (id)
);

CREATE INDEX IF NOT EXISTS idx_answers_user_id ON answers (user_id);
CREATE INDEX IF NOT EXISTS idx_answers_question_id ON answers (question_id);
CREATE INDEX IF NOT EXISTS idx_answers_deleted_at ON answers (deleted_at);

-- ============================================================================
-- 23. answer_votes — 回答への投票
-- ============================================================================
CREATE TABLE IF NOT EXISTS answer_votes (
    id         BIGSERIAL    PRIMARY KEY,
    user_id    BIGINT       NOT NULL,
    answer_id  BIGINT       NOT NULL,
    value      INTEGER      NOT NULL,
    created_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_answer_vote ON answer_votes (user_id, answer_id);

-- ============================================================================
-- 24. roadmaps — 学習ロードマップ
-- ============================================================================
CREATE TABLE IF NOT EXISTS roadmaps (
    id                   BIGSERIAL      PRIMARY KEY,
    user_id              BIGINT         NOT NULL,
    title                VARCHAR(200)   NOT NULL,
    description          TEXT,
    category             TEXT           DEFAULT 'other',
    is_public            BOOLEAN        NOT NULL DEFAULT FALSE,
    step_count           INTEGER        NOT NULL DEFAULT 0,
    completed_step_count INTEGER        NOT NULL DEFAULT 0,
    progress             INTEGER        NOT NULL DEFAULT 0,
    status               TEXT           DEFAULT 'active',
    created_at           TIMESTAMPTZ,
    updated_at           TIMESTAMPTZ,
    completed_at         TIMESTAMPTZ,
    CONSTRAINT fk_roadmaps_user FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX IF NOT EXISTS idx_roadmaps_user_id ON roadmaps (user_id);
CREATE INDEX IF NOT EXISTS idx_roadmaps_is_public ON roadmaps (is_public);

-- ============================================================================
-- 25. roadmap_steps — ロードマップ内のステップ（CASCADE削除）
-- ============================================================================
CREATE TABLE IF NOT EXISTS roadmap_steps (
    id           BIGSERIAL      PRIMARY KEY,
    roadmap_id   BIGINT         NOT NULL,
    title        VARCHAR(200)   NOT NULL,
    description  TEXT,
    order_index  INTEGER        NOT NULL DEFAULT 0,
    is_completed BOOLEAN        NOT NULL DEFAULT FALSE,
    completed_at TIMESTAMPTZ,
    resource_url VARCHAR(500),
    created_at   TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ,
    CONSTRAINT fk_roadmap_steps_roadmap FOREIGN KEY (roadmap_id) REFERENCES roadmaps (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_roadmap_steps_roadmap_id ON roadmap_steps (roadmap_id);

-- ============================================================================
-- 26. notifications — 通知
-- ============================================================================
CREATE TABLE IF NOT EXISTS notifications (
    id          BIGSERIAL    PRIMARY KEY,
    user_id     BIGINT       NOT NULL,
    type        TEXT         NOT NULL,
    actor_id    BIGINT       NOT NULL,
    post_id     BIGINT,
    question_id BIGINT,
    badge_id    VARCHAR(50),
    read        BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ,
    CONSTRAINT fk_notifications_user FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_notifications_actor FOREIGN KEY (actor_id) REFERENCES users (id),
    CONSTRAINT fk_notifications_post FOREIGN KEY (post_id) REFERENCES posts (id),
    CONSTRAINT fk_notifications_question FOREIGN KEY (question_id) REFERENCES questions (id)
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications (user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_post_id ON notifications (post_id);
CREATE INDEX IF NOT EXISTS idx_notifications_question_id ON notifications (question_id);

-- ============================================================================
-- 27. chat_rooms — グループチャットルーム
-- ============================================================================
CREATE TABLE IF NOT EXISTS chat_rooms (
    id          BIGSERIAL      PRIMARY KEY,
    name        VARCHAR(100)   NOT NULL,
    description VARCHAR(500),
    owner_id    BIGINT         NOT NULL,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    CONSTRAINT fk_chat_rooms_owner FOREIGN KEY (owner_id) REFERENCES users (id)
);

CREATE INDEX IF NOT EXISTS idx_chat_rooms_owner_id ON chat_rooms (owner_id);

-- ============================================================================
-- 28. chat_room_members — チャットルームメンバー
-- ============================================================================
CREATE TABLE IF NOT EXISTS chat_room_members (
    id           BIGSERIAL    PRIMARY KEY,
    chat_room_id BIGINT       NOT NULL,
    user_id      BIGINT       NOT NULL,
    joined_at    TIMESTAMPTZ,
    CONSTRAINT fk_chat_room_members_room FOREIGN KEY (chat_room_id) REFERENCES chat_rooms (id),
    CONSTRAINT fk_chat_room_members_user FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX IF NOT EXISTS idx_chat_room_members_chat_room_id ON chat_room_members (chat_room_id);
CREATE INDEX IF NOT EXISTS idx_chat_room_members_user_id ON chat_room_members (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_room_user ON chat_room_members (chat_room_id, user_id);

-- ============================================================================
-- 29. group_messages — グループチャットメッセージ
-- ============================================================================
CREATE TABLE IF NOT EXISTS group_messages (
    id           BIGSERIAL    PRIMARY KEY,
    chat_room_id BIGINT       NOT NULL,
    sender_id    BIGINT       NOT NULL,
    content      TEXT         NOT NULL,
    created_at   TIMESTAMPTZ,
    CONSTRAINT fk_group_messages_room FOREIGN KEY (chat_room_id) REFERENCES chat_rooms (id),
    CONSTRAINT fk_group_messages_sender FOREIGN KEY (sender_id) REFERENCES users (id)
);

CREATE INDEX IF NOT EXISTS idx_group_messages_chat_room_id ON group_messages (chat_room_id);
CREATE INDEX IF NOT EXISTS idx_group_messages_sender_id ON group_messages (sender_id);
