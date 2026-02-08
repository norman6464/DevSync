-- ============================================================================
-- DevSync 初期マイグレーション（UP）
-- 全29テーブルのCREATE TABLE文（PostgreSQL 15）
-- GORMのAutoMigrateが実際に生成するスキーマと完全一致
-- ============================================================================

-- 1. users
CREATE TABLE IF NOT EXISTS users (
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
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email            ON users (email);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_git_hub_id       ON users (git_hub_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_git_hub_username ON users (git_hub_username);

-- 2. follows
CREATE TABLE IF NOT EXISTS follows (
    id          BIGSERIAL   PRIMARY KEY,
    follower_id BIGINT      NOT NULL,
    followee_id BIGINT      NOT NULL,
    created_at  TIMESTAMPTZ,
    CONSTRAINT fk_follows_follower FOREIGN KEY (follower_id) REFERENCES users (id),
    CONSTRAINT fk_follows_followee FOREIGN KEY (followee_id) REFERENCES users (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_follower_following ON follows (follower_id, followee_id);

-- 3. git_hub_contributions
CREATE TABLE IF NOT EXISTS git_hub_contributions (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL,
    date       TIMESTAMPTZ NOT NULL,
    count      BIGINT      NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_date ON git_hub_contributions (user_id, date);

-- 4. git_hub_language_stats
CREATE TABLE IF NOT EXISTS git_hub_language_stats (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL,
    language   TEXT        NOT NULL,
    bytes      BIGINT      NOT NULL DEFAULT 0,
    repo_count BIGINT      NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_lang ON git_hub_language_stats (user_id, language);

-- 5. git_hub_repositories
CREATE TABLE IF NOT EXISTS git_hub_repositories (
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
CREATE INDEX        IF NOT EXISTS idx_git_hub_repositories_user_id         ON git_hub_repositories (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_git_hub_repositories_git_hub_repo_id ON git_hub_repositories (git_hub_repo_id);

-- 6. posts
CREATE TABLE IF NOT EXISTS posts (
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
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id);

-- 7. likes
CREATE TABLE IF NOT EXISTS likes (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL,
    post_id    BIGINT      NOT NULL,
    created_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_post_like ON likes (user_id, post_id);
CREATE INDEX        IF NOT EXISTS idx_likes_post_id  ON likes (post_id);

-- 8. comments
CREATE TABLE IF NOT EXISTS comments (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL,
    post_id    BIGINT      NOT NULL,
    content    TEXT        NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT fk_comments_user FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments (user_id);
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments (post_id);

-- 9. messages
CREATE TABLE IF NOT EXISTS messages (
    id          BIGSERIAL   PRIMARY KEY,
    sender_id   BIGINT      NOT NULL,
    receiver_id BIGINT      NOT NULL,
    content     TEXT        NOT NULL,
    read        BOOLEAN     DEFAULT false,
    created_at  TIMESTAMPTZ,
    CONSTRAINT fk_messages_sender   FOREIGN KEY (sender_id)   REFERENCES users (id),
    CONSTRAINT fk_messages_receiver FOREIGN KEY (receiver_id) REFERENCES users (id)
);
CREATE INDEX IF NOT EXISTS idx_messages_sender_id   ON messages (sender_id);
CREATE INDEX IF NOT EXISTS idx_messages_receiver_id ON messages (receiver_id);

-- 10. password_reset_tokens
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL,
    token      TEXT        NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used       BOOLEAN     DEFAULT false,
    created_at TIMESTAMPTZ,
    CONSTRAINT fk_password_reset_tokens_user FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE INDEX        IF NOT EXISTS idx_password_reset_tokens_user_id ON password_reset_tokens (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_password_reset_tokens_token   ON password_reset_tokens (token);

-- 11. zenn_articles
CREATE TABLE IF NOT EXISTS zenn_articles (
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
CREATE INDEX        IF NOT EXISTS idx_zenn_articles_user_id ON zenn_articles (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_zenn_articles_zenn_id ON zenn_articles (zenn_id);

-- 12. qiita_articles
CREATE TABLE IF NOT EXISTS qiita_articles (
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
CREATE INDEX        IF NOT EXISTS idx_qiita_articles_user_id  ON qiita_articles (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_qiita_articles_qiita_id ON qiita_articles (qiita_id);

-- 13. learning_goals
CREATE TABLE IF NOT EXISTS learning_goals (
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
CREATE INDEX IF NOT EXISTS idx_learning_goals_user_id ON learning_goals (user_id);

-- 14. learning_logs
CREATE TABLE IF NOT EXISTS learning_logs (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL,
    title      TEXT        NOT NULL,
    content    TEXT        NOT NULL,
    category   TEXT        DEFAULT 'other'::text,
    duration   BIGINT      DEFAULT 0,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_learning_logs_user_id ON learning_logs (user_id);

-- 15. projects（論理削除対応）
CREATE TABLE IF NOT EXISTS projects (
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
CREATE INDEX IF NOT EXISTS idx_projects_user_id    ON projects (user_id);
CREATE INDEX IF NOT EXISTS idx_projects_deleted_at ON projects (deleted_at);

-- 16. learning_resources（論理削除対応）
CREATE TABLE IF NOT EXISTS learning_resources (
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
CREATE INDEX IF NOT EXISTS idx_learning_resources_user_id    ON learning_resources (user_id);
CREATE INDEX IF NOT EXISTS idx_learning_resources_deleted_at ON learning_resources (deleted_at);

-- 17. resource_likes
CREATE TABLE IF NOT EXISTS resource_likes (
    id          BIGSERIAL   PRIMARY KEY,
    user_id     BIGINT      NOT NULL,
    resource_id BIGINT      NOT NULL,
    created_at  TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_resource_like ON resource_likes (user_id, resource_id);

-- 18. resource_saves
CREATE TABLE IF NOT EXISTS resource_saves (
    id          BIGSERIAL   PRIMARY KEY,
    user_id     BIGINT      NOT NULL,
    resource_id BIGINT      NOT NULL,
    created_at  TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_resource_save ON resource_saves (user_id, resource_id);

-- 19. book_reviews（論理削除対応）
CREATE TABLE IF NOT EXISTS book_reviews (
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
CREATE INDEX IF NOT EXISTS idx_book_reviews_user_id    ON book_reviews (user_id);
CREATE INDEX IF NOT EXISTS idx_book_reviews_deleted_at ON book_reviews (deleted_at);

-- 20. questions（論理削除対応）
CREATE TABLE IF NOT EXISTS questions (
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
CREATE INDEX IF NOT EXISTS idx_questions_user_id    ON questions (user_id);
CREATE INDEX IF NOT EXISTS idx_questions_deleted_at ON questions (deleted_at);

-- 21. question_votes
CREATE TABLE IF NOT EXISTS question_votes (
    id          BIGSERIAL   PRIMARY KEY,
    user_id     BIGINT      NOT NULL,
    question_id BIGINT      NOT NULL,
    value       BIGINT      NOT NULL,
    created_at  TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_question_vote ON question_votes (user_id, question_id);

-- 22. answers（論理削除対応）
CREATE TABLE IF NOT EXISTS answers (
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
CREATE INDEX IF NOT EXISTS idx_answers_user_id     ON answers (user_id);
CREATE INDEX IF NOT EXISTS idx_answers_question_id ON answers (question_id);
CREATE INDEX IF NOT EXISTS idx_answers_deleted_at  ON answers (deleted_at);

-- 23. answer_votes
CREATE TABLE IF NOT EXISTS answer_votes (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL,
    answer_id  BIGINT      NOT NULL,
    value      BIGINT      NOT NULL,
    created_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_answer_vote ON answer_votes (user_id, answer_id);

-- 24. roadmaps
CREATE TABLE IF NOT EXISTS roadmaps (
    id                   BIGSERIAL    PRIMARY KEY,
    user_id              BIGINT       NOT NULL,
    title                VARCHAR(200) NOT NULL,
    description          TEXT,
    category             TEXT         DEFAULT 'other'::text,
    is_public            BOOLEAN      DEFAULT false,
    step_count           BIGINT       DEFAULT 0,
    completed_step_count BIGINT       DEFAULT 0,
    progress             BIGINT       DEFAULT 0,
    status               TEXT         DEFAULT 'active'::text,
    created_at           TIMESTAMPTZ,
    updated_at           TIMESTAMPTZ,
    completed_at         TIMESTAMPTZ,
    CONSTRAINT fk_roadmaps_user FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE INDEX IF NOT EXISTS idx_roadmaps_user_id   ON roadmaps (user_id);
CREATE INDEX IF NOT EXISTS idx_roadmaps_is_public ON roadmaps (is_public);

-- 25. roadmap_steps（CASCADE削除）
CREATE TABLE IF NOT EXISTS roadmap_steps (
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
CREATE INDEX IF NOT EXISTS idx_roadmap_steps_roadmap_id ON roadmap_steps (roadmap_id);

-- 26. notifications
CREATE TABLE IF NOT EXISTS notifications (
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
CREATE INDEX IF NOT EXISTS idx_notifications_user_id     ON notifications (user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_post_id     ON notifications (post_id);
CREATE INDEX IF NOT EXISTS idx_notifications_question_id ON notifications (question_id);

-- 27. chat_rooms
CREATE TABLE IF NOT EXISTS chat_rooms (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    owner_id    BIGINT       NOT NULL,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    CONSTRAINT fk_chat_rooms_owner FOREIGN KEY (owner_id) REFERENCES users (id)
);
CREATE INDEX IF NOT EXISTS idx_chat_rooms_owner_id ON chat_rooms (owner_id);

-- 28. chat_room_members
CREATE TABLE IF NOT EXISTS chat_room_members (
    id           BIGSERIAL   PRIMARY KEY,
    chat_room_id BIGINT      NOT NULL,
    user_id      BIGINT      NOT NULL,
    joined_at    TIMESTAMPTZ,
    CONSTRAINT fk_chat_room_members_chat_room FOREIGN KEY (chat_room_id) REFERENCES chat_rooms (id),
    CONSTRAINT fk_chat_room_members_user      FOREIGN KEY (user_id)      REFERENCES users (id)
);
CREATE INDEX        IF NOT EXISTS idx_chat_room_members_chat_room_id ON chat_room_members (chat_room_id);
CREATE INDEX        IF NOT EXISTS idx_chat_room_members_user_id      ON chat_room_members (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_room_user                      ON chat_room_members (chat_room_id, user_id);

-- 29. group_messages
CREATE TABLE IF NOT EXISTS group_messages (
    id           BIGSERIAL   PRIMARY KEY,
    chat_room_id BIGINT      NOT NULL,
    sender_id    BIGINT      NOT NULL,
    content      TEXT        NOT NULL,
    created_at   TIMESTAMPTZ,
    CONSTRAINT fk_group_messages_chat_room FOREIGN KEY (chat_room_id) REFERENCES chat_rooms (id),
    CONSTRAINT fk_group_messages_sender    FOREIGN KEY (sender_id)    REFERENCES users (id)
);
CREATE INDEX IF NOT EXISTS idx_group_messages_chat_room_id ON group_messages (chat_room_id);
CREATE INDEX IF NOT EXISTS idx_group_messages_sender_id    ON group_messages (sender_id);
