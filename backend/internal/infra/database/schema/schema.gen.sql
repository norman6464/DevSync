-- Add new schema named "public"
CREATE SCHEMA IF NOT EXISTS "public";
-- Set comment to schema: "public"
COMMENT ON SCHEMA "public" IS 'standard public schema';
-- Create "ai_advices" table
CREATE TABLE "public"."ai_advices" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "type" text NOT NULL,
  "priority" bigint NOT NULL DEFAULT 3,
  "title_key" text NOT NULL,
  "message_key" text NOT NULL,
  "params" text NULL,
  "action_url" text NULL,
  "is_read" boolean NULL DEFAULT false,
  "expires_at" timestamptz NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_ai_advices_user_id" to table: "ai_advices"
CREATE INDEX "idx_ai_advices_user_id" ON "public"."ai_advices" ("user_id");
-- Create "ai_conversations" table
CREATE TABLE "public"."ai_conversations" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "title" character varying(200) NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_ai_conversations_user_id" to table: "ai_conversations"
CREATE INDEX "idx_ai_conversations_user_id" ON "public"."ai_conversations" ("user_id");
-- Create "answer_votes" table
CREATE TABLE "public"."answer_votes" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "answer_id" bigint NOT NULL,
  "value" bigint NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_answer_vote" to table: "answer_votes"
CREATE UNIQUE INDEX "idx_answer_vote" ON "public"."answer_votes" ("user_id", "answer_id");
-- Create "bookmark_collections" table
CREATE TABLE "public"."bookmark_collections" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "name" text NOT NULL,
  "description" text NULL,
  "color" text NULL DEFAULT 'blue',
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_bookmark_collections_user_id" to table: "bookmark_collections"
CREATE INDEX "idx_bookmark_collections_user_id" ON "public"."bookmark_collections" ("user_id");
-- Create "code_snippet_favorites" table
CREATE TABLE "public"."code_snippet_favorites" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "snippet_id" bigint NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_snippet_favorite" to table: "code_snippet_favorites"
CREATE UNIQUE INDEX "idx_snippet_favorite" ON "public"."code_snippet_favorites" ("user_id", "snippet_id");
-- Create "comment_likes" table
CREATE TABLE "public"."comment_likes" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "comment_id" bigint NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_comment_likes_comment_id" to table: "comment_likes"
CREATE INDEX "idx_comment_likes_comment_id" ON "public"."comment_likes" ("comment_id");
-- Create index "idx_user_comment_like" to table: "comment_likes"
CREATE UNIQUE INDEX "idx_user_comment_like" ON "public"."comment_likes" ("user_id", "comment_id");
-- Create "git_hub_contributions" table
CREATE TABLE "public"."git_hub_contributions" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "date" timestamptz NOT NULL,
  "count" bigint NOT NULL DEFAULT 0,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_user_date" to table: "git_hub_contributions"
CREATE UNIQUE INDEX "idx_user_date" ON "public"."git_hub_contributions" ("user_id", "date");
-- Create "git_hub_language_stats" table
CREATE TABLE "public"."git_hub_language_stats" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "language" text NOT NULL,
  "bytes" bigint NOT NULL DEFAULT 0,
  "repo_count" bigint NOT NULL DEFAULT 0,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_user_lang" to table: "git_hub_language_stats"
CREATE UNIQUE INDEX "idx_user_lang" ON "public"."git_hub_language_stats" ("user_id", "language");
-- Create "learning_goals" table
CREATE TABLE "public"."learning_goals" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "title" text NOT NULL,
  "description" text NULL,
  "category" text NULL DEFAULT 'other',
  "target_date" timestamptz NULL,
  "progress" bigint NULL DEFAULT 0,
  "target_hours" bigint NULL DEFAULT 0,
  "status" text NULL DEFAULT 'active',
  "is_public" boolean NULL DEFAULT false,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "completed_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_learning_goals_user_id" to table: "learning_goals"
CREATE INDEX "idx_learning_goals_user_id" ON "public"."learning_goals" ("user_id");
-- Create "learning_log_templates" table
CREATE TABLE "public"."learning_log_templates" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "name" character varying(100) NOT NULL,
  "default_title" character varying(200) NULL,
  "default_content" text NULL,
  "default_category" character varying(50) NULL DEFAULT 'other',
  "default_duration" bigint NULL DEFAULT 0,
  "is_default" boolean NULL DEFAULT false,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_learning_log_templates_user_id" to table: "learning_log_templates"
CREATE INDEX "idx_learning_log_templates_user_id" ON "public"."learning_log_templates" ("user_id");
-- Create index "idx_log_templates_is_default" to table: "learning_log_templates"
CREATE INDEX "idx_log_templates_is_default" ON "public"."learning_log_templates" ("is_default");
-- Create "learning_logs" table
CREATE TABLE "public"."learning_logs" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "title" text NOT NULL,
  "content" text NOT NULL,
  "category" text NULL DEFAULT 'other',
  "duration" bigint NULL DEFAULT 0,
  "goal_id" bigint NULL,
  "source" text NULL DEFAULT 'manual',
  "is_favorite" boolean NULL DEFAULT false,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_learning_logs_goal_id" to table: "learning_logs"
CREATE INDEX "idx_learning_logs_goal_id" ON "public"."learning_logs" ("goal_id");
-- Create index "idx_learning_logs_user_id" to table: "learning_logs"
CREATE INDEX "idx_learning_logs_user_id" ON "public"."learning_logs" ("user_id");
-- Create "likes" table
CREATE TABLE "public"."likes" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "post_id" bigint NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_likes_post_id" to table: "likes"
CREATE INDEX "idx_likes_post_id" ON "public"."likes" ("post_id");
-- Create index "idx_user_post_like" to table: "likes"
CREATE UNIQUE INDEX "idx_user_post_like" ON "public"."likes" ("user_id", "post_id");
-- Create "note_templates" table
CREATE TABLE "public"."note_templates" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "name" character varying(100) NOT NULL,
  "description" character varying(500) NULL,
  "default_title" character varying(200) NULL,
  "content_template" text NOT NULL,
  "default_tags" character varying(255) NULL,
  "is_default" boolean NULL DEFAULT false,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_note_templates_is_default" to table: "note_templates"
CREATE INDEX "idx_note_templates_is_default" ON "public"."note_templates" ("is_default");
-- Create index "idx_note_templates_user_id" to table: "note_templates"
CREATE INDEX "idx_note_templates_user_id" ON "public"."note_templates" ("user_id");
-- Create "note_versions" table
CREATE TABLE "public"."note_versions" (
  "id" bigserial NOT NULL,
  "note_id" bigint NOT NULL,
  "version_number" bigint NOT NULL,
  "title" text NOT NULL,
  "content" text NULL,
  "tags" text NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_note_versions_note_id" to table: "note_versions"
CREATE INDEX "idx_note_versions_note_id" ON "public"."note_versions" ("note_id");
-- Create "post_tags" table
CREATE TABLE "public"."post_tags" (
  "id" bigserial NOT NULL,
  "post_id" bigint NOT NULL,
  "tag" character varying(50) NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_post_tag" to table: "post_tags"
CREATE UNIQUE INDEX "idx_post_tag" ON "public"."post_tags" ("post_id", "tag");
-- Create index "idx_post_tags_post_id" to table: "post_tags"
CREATE INDEX "idx_post_tags_post_id" ON "public"."post_tags" ("post_id");
-- Create index "idx_post_tags_tag" to table: "post_tags"
CREATE INDEX "idx_post_tags_tag" ON "public"."post_tags" ("tag");
-- Create "post_templates" table
CREATE TABLE "public"."post_templates" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "name" character varying(100) NOT NULL,
  "title_template" character varying(200) NULL,
  "content_template" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_post_templates_user_id" to table: "post_templates"
CREATE INDEX "idx_post_templates_user_id" ON "public"."post_templates" ("user_id");
-- Create "post_views" table
CREATE TABLE "public"."post_views" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "post_id" bigint NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_post_views_post_id" to table: "post_views"
CREATE INDEX "idx_post_views_post_id" ON "public"."post_views" ("post_id");
-- Create index "idx_post_views_user_id" to table: "post_views"
CREATE INDEX "idx_post_views_user_id" ON "public"."post_views" ("user_id");
-- Create index "idx_user_post_view" to table: "post_views"
CREATE UNIQUE INDEX "idx_user_post_view" ON "public"."post_views" ("user_id", "post_id");
-- Create "project_milestones" table
CREATE TABLE "public"."project_milestones" (
  "id" bigserial NOT NULL,
  "project_id" bigint NOT NULL,
  "title" character varying(200) NOT NULL,
  "description" text NULL,
  "status" character varying(20) NULL DEFAULT 'not_started',
  "due_date" timestamptz NULL,
  "completed_at" timestamptz NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_project_milestones_project_id" to table: "project_milestones"
CREATE INDEX "idx_project_milestones_project_id" ON "public"."project_milestones" ("project_id");
-- Create "qiita_articles" table
CREATE TABLE "public"."qiita_articles" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "qiita_id" text NOT NULL,
  "title" text NOT NULL,
  "url" text NOT NULL,
  "likes_count" bigint NULL DEFAULT 0,
  "comments_count" bigint NULL DEFAULT 0,
  "tags" text NULL,
  "published_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_qiita_articles_qiita_id" to table: "qiita_articles"
CREATE UNIQUE INDEX "idx_qiita_articles_qiita_id" ON "public"."qiita_articles" ("qiita_id");
-- Create index "idx_qiita_articles_user_id" to table: "qiita_articles"
CREATE INDEX "idx_qiita_articles_user_id" ON "public"."qiita_articles" ("user_id");
-- Create "question_votes" table
CREATE TABLE "public"."question_votes" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "question_id" bigint NOT NULL,
  "value" bigint NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_question_vote" to table: "question_votes"
CREATE UNIQUE INDEX "idx_question_vote" ON "public"."question_votes" ("user_id", "question_id");
-- Create "reactions" table
CREATE TABLE "public"."reactions" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "post_id" bigint NOT NULL,
  "emoji" character varying(10) NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_reactions_post_id" to table: "reactions"
CREATE INDEX "idx_reactions_post_id" ON "public"."reactions" ("post_id");
-- Create index "idx_user_post_emoji" to table: "reactions"
CREATE UNIQUE INDEX "idx_user_post_emoji" ON "public"."reactions" ("user_id", "post_id", "emoji");
-- Create "resource_likes" table
CREATE TABLE "public"."resource_likes" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "resource_id" bigint NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_resource_like" to table: "resource_likes"
CREATE UNIQUE INDEX "idx_resource_like" ON "public"."resource_likes" ("user_id", "resource_id");
-- Create "resource_saves" table
CREATE TABLE "public"."resource_saves" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "resource_id" bigint NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_resource_save" to table: "resource_saves"
CREATE UNIQUE INDEX "idx_resource_save" ON "public"."resource_saves" ("user_id", "resource_id");
-- Create "spotify_recent_tracks" table
CREATE TABLE "public"."spotify_recent_tracks" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "track_name" text NOT NULL,
  "artist_name" text NOT NULL,
  "album_name" text NULL,
  "album_image" text NULL,
  "track_url" text NULL,
  "played_at" timestamptz NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_spotify_recent_tracks_user_id" to table: "spotify_recent_tracks"
CREATE INDEX "idx_spotify_recent_tracks_user_id" ON "public"."spotify_recent_tracks" ("user_id");
-- Create "streak_freezes" table
CREATE TABLE "public"."streak_freezes" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "used_date" character varying(10) NOT NULL,
  "month" bigint NOT NULL,
  "year" bigint NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_streak_freeze_user_date" to table: "streak_freezes"
CREATE UNIQUE INDEX "idx_streak_freeze_user_date" ON "public"."streak_freezes" ("user_id", "used_date");
-- Create index "idx_streak_freezes_user_id" to table: "streak_freezes"
CREATE INDEX "idx_streak_freezes_user_id" ON "public"."streak_freezes" ("user_id");
-- Create "user_activities" table
CREATE TABLE "public"."user_activities" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "activity_type" character varying(50) NOT NULL,
  "target_type" character varying(50) NOT NULL,
  "target_id" bigint NOT NULL,
  "metadata" text NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_user_activities_activity_type" to table: "user_activities"
CREATE INDEX "idx_user_activities_activity_type" ON "public"."user_activities" ("activity_type");
-- Create index "idx_user_activities_created_at" to table: "user_activities"
CREATE INDEX "idx_user_activities_created_at" ON "public"."user_activities" ("created_at");
-- Create index "idx_user_activities_user_id" to table: "user_activities"
CREATE INDEX "idx_user_activities_user_id" ON "public"."user_activities" ("user_id");
-- Create "weekly_challenges" table
CREATE TABLE "public"."weekly_challenges" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "year" bigint NOT NULL,
  "week" bigint NOT NULL,
  "challenge_type" text NOT NULL,
  "description" text NOT NULL,
  "target_value" bigint NOT NULL,
  "current_value" bigint NULL DEFAULT 0,
  "is_completed" boolean NULL DEFAULT false,
  "completed_at" timestamptz NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_weekly_challenges_user_id" to table: "weekly_challenges"
CREATE INDEX "idx_weekly_challenges_user_id" ON "public"."weekly_challenges" ("user_id");
-- Create "weekly_goals" table
CREATE TABLE "public"."weekly_goals" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "category" text NOT NULL,
  "target_minutes" bigint NOT NULL DEFAULT 0,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_weekly_goal_unique" to table: "weekly_goals"
CREATE UNIQUE INDEX "idx_weekly_goal_unique" ON "public"."weekly_goals" ("user_id", "category");
-- Create "widget_settings" table
CREATE TABLE "public"."widget_settings" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "settings" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_widget_settings_user_id" to table: "widget_settings"
CREATE UNIQUE INDEX "idx_widget_settings_user_id" ON "public"."widget_settings" ("user_id");
-- Create "you_tube_search_caches" table
CREATE TABLE "public"."you_tube_search_caches" (
  "id" bigserial NOT NULL,
  "query" character varying(500) NOT NULL,
  "language" character varying(10) NULL DEFAULT 'ja',
  "video_ids" text NULL,
  "cache_expires" timestamptz NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_you_tube_search_caches_cache_expires" to table: "you_tube_search_caches"
CREATE INDEX "idx_you_tube_search_caches_cache_expires" ON "public"."you_tube_search_caches" ("cache_expires");
-- Create index "idx_you_tube_search_caches_query" to table: "you_tube_search_caches"
CREATE INDEX "idx_you_tube_search_caches_query" ON "public"."you_tube_search_caches" ("query");
-- Create "you_tube_videos" table
CREATE TABLE "public"."you_tube_videos" (
  "id" bigserial NOT NULL,
  "video_id" character varying(20) NOT NULL,
  "title" character varying(500) NOT NULL,
  "description" text NULL,
  "channel_id" character varying(50) NULL,
  "channel_title" character varying(200) NULL,
  "thumbnail_url" character varying(500) NULL,
  "published_at" timestamptz NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_you_tube_videos_video_id" to table: "you_tube_videos"
CREATE UNIQUE INDEX "idx_you_tube_videos_video_id" ON "public"."you_tube_videos" ("video_id");
-- Create "zenn_articles" table
CREATE TABLE "public"."zenn_articles" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "zenn_id" bigint NOT NULL,
  "title" text NOT NULL,
  "slug" text NOT NULL,
  "emoji" text NULL,
  "article_type" text NULL,
  "liked_count" bigint NULL DEFAULT 0,
  "comments_count" bigint NULL DEFAULT 0,
  "published_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_zenn_articles_user_id" to table: "zenn_articles"
CREATE INDEX "idx_zenn_articles_user_id" ON "public"."zenn_articles" ("user_id");
-- Create index "idx_zenn_articles_zenn_id" to table: "zenn_articles"
CREATE UNIQUE INDEX "idx_zenn_articles_zenn_id" ON "public"."zenn_articles" ("zenn_id");
-- Create "ai_messages" table
CREATE TABLE "public"."ai_messages" (
  "id" bigserial NOT NULL,
  "conversation_id" bigint NOT NULL,
  "role" text NOT NULL,
  "content" text NOT NULL,
  "tokens_used" bigint NULL DEFAULT 0,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_ai_conversations_messages" FOREIGN KEY ("conversation_id") REFERENCES "public"."ai_conversations" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_ai_messages_conversation_id" to table: "ai_messages"
CREATE INDEX "idx_ai_messages_conversation_id" ON "public"."ai_messages" ("conversation_id");
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "username" text NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "password" text NULL,
  "avatar_url" text NULL,
  "bio" text NULL,
  "git_hub_id" bigint NULL,
  "git_hub_username" text NULL,
  "git_hub_token" text NULL,
  "git_hub_connected" boolean NULL DEFAULT false,
  "spotify_connected" boolean NULL DEFAULT false,
  "spotify_token" text NULL,
  "spotify_refresh_token" text NULL,
  "spotify_token_expiry" timestamptz NULL,
  "zenn_username" text NULL,
  "qiita_username" text NULL,
  "at_coder_username" text NULL,
  "paiza_rank" text NULL,
  "skills_languages" text NULL,
  "skills_frameworks" text NULL,
  "onboarding_completed" boolean NULL DEFAULT false,
  "email_weekly_report" boolean NULL DEFAULT true,
  "email_language" text NULL DEFAULT 'ja',
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "public"."users" ("email");
-- Create index "idx_users_git_hub_id_linked" to table: "users"
CREATE UNIQUE INDEX "idx_users_git_hub_id_linked" ON "public"."users" ("git_hub_id") WHERE (git_hub_id <> 0);
-- Create index "idx_users_username" to table: "users"
CREATE UNIQUE INDEX "idx_users_username" ON "public"."users" ("username");
-- Create "answers" table
CREATE TABLE "public"."answers" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "question_id" bigint NOT NULL,
  "body" text NOT NULL,
  "vote_count" bigint NULL DEFAULT 0,
  "is_best" boolean NULL DEFAULT false,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_answers_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_answers_deleted_at" to table: "answers"
CREATE INDEX "idx_answers_deleted_at" ON "public"."answers" ("deleted_at");
-- Create index "idx_answers_question_id" to table: "answers"
CREATE INDEX "idx_answers_question_id" ON "public"."answers" ("question_id");
-- Create index "idx_answers_user_id" to table: "answers"
CREATE INDEX "idx_answers_user_id" ON "public"."answers" ("user_id");
-- Create "book_reviews" table
CREATE TABLE "public"."book_reviews" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "title" character varying(300) NOT NULL,
  "author" character varying(200) NULL,
  "isbn" character varying(20) NULL,
  "rating" bigint NOT NULL,
  "review" text NULL,
  "total_pages" bigint NULL DEFAULT 0,
  "current_page" bigint NULL DEFAULT 0,
  "image_url" character varying(2000) NULL,
  "status" text NULL DEFAULT 'not_started',
  "is_archived" boolean NULL DEFAULT false,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_book_reviews_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_book_reviews_deleted_at" to table: "book_reviews"
CREATE INDEX "idx_book_reviews_deleted_at" ON "public"."book_reviews" ("deleted_at");
-- Create index "idx_book_reviews_user_id" to table: "book_reviews"
CREATE INDEX "idx_book_reviews_user_id" ON "public"."book_reviews" ("user_id");
-- Create "posts" table
CREATE TABLE "public"."posts" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "title" text NOT NULL,
  "content" text NOT NULL,
  "image_urls" text NULL,
  "is_draft" boolean NULL DEFAULT false,
  "like_count" bigint NULL DEFAULT 0,
  "comment_count" bigint NULL DEFAULT 0,
  "bookmark_count" bigint NULL DEFAULT 0,
  "view_count" bigint NULL DEFAULT 0,
  "estimated_read_time" bigint NULL DEFAULT 0,
  "scheduled_at" timestamptz NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_posts_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_posts_is_draft" to table: "posts"
CREATE INDEX "idx_posts_is_draft" ON "public"."posts" ("is_draft");
-- Create index "idx_posts_scheduled_at" to table: "posts"
CREATE INDEX "idx_posts_scheduled_at" ON "public"."posts" ("scheduled_at");
-- Create index "idx_posts_user_id" to table: "posts"
CREATE INDEX "idx_posts_user_id" ON "public"."posts" ("user_id");
-- Create "bookmark_collection_items" table
CREATE TABLE "public"."bookmark_collection_items" (
  "id" bigserial NOT NULL,
  "collection_id" bigint NOT NULL,
  "post_id" bigint NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_bookmark_collection_items_post" FOREIGN KEY ("post_id") REFERENCES "public"."posts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_bookmark_collection_items_post_id" to table: "bookmark_collection_items"
CREATE INDEX "idx_bookmark_collection_items_post_id" ON "public"."bookmark_collection_items" ("post_id");
-- Create index "idx_bookmark_collection_post" to table: "bookmark_collection_items"
CREATE UNIQUE INDEX "idx_bookmark_collection_post" ON "public"."bookmark_collection_items" ("collection_id", "post_id");
-- Create "bookmarks" table
CREATE TABLE "public"."bookmarks" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "post_id" bigint NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_bookmarks_post" FOREIGN KEY ("post_id") REFERENCES "public"."posts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_bookmarks_post_id" to table: "bookmarks"
CREATE INDEX "idx_bookmarks_post_id" ON "public"."bookmarks" ("post_id");
-- Create index "idx_user_post_bookmark" to table: "bookmarks"
CREATE UNIQUE INDEX "idx_user_post_bookmark" ON "public"."bookmarks" ("user_id", "post_id");
-- Create "chat_rooms" table
CREATE TABLE "public"."chat_rooms" (
  "id" bigserial NOT NULL,
  "name" character varying(100) NOT NULL,
  "description" character varying(500) NULL,
  "owner_id" bigint NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_chat_rooms_owner" FOREIGN KEY ("owner_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_chat_rooms_owner_id" to table: "chat_rooms"
CREATE INDEX "idx_chat_rooms_owner_id" ON "public"."chat_rooms" ("owner_id");
-- Create "chat_room_members" table
CREATE TABLE "public"."chat_room_members" (
  "id" bigserial NOT NULL,
  "chat_room_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "joined_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_chat_room_members_chat_room" FOREIGN KEY ("chat_room_id") REFERENCES "public"."chat_rooms" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_chat_room_members_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_chat_room_members_chat_room_id" to table: "chat_room_members"
CREATE INDEX "idx_chat_room_members_chat_room_id" ON "public"."chat_room_members" ("chat_room_id");
-- Create index "idx_chat_room_members_user_id" to table: "chat_room_members"
CREATE INDEX "idx_chat_room_members_user_id" ON "public"."chat_room_members" ("user_id");
-- Create index "idx_room_user" to table: "chat_room_members"
CREATE UNIQUE INDEX "idx_room_user" ON "public"."chat_room_members" ("chat_room_id", "user_id");
-- Create "code_snippets" table
CREATE TABLE "public"."code_snippets" (
  "id" bigserial NOT NULL,
  "post_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "language" text NOT NULL,
  "file_name" text NULL,
  "code" text NOT NULL,
  "comment_count" bigint NULL DEFAULT 0,
  "forked_from_id" bigint NULL,
  "fork_count" bigint NULL DEFAULT 0,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_posts_code_snippets" FOREIGN KEY ("post_id") REFERENCES "public"."posts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_code_snippets_forked_from_id" to table: "code_snippets"
CREATE INDEX "idx_code_snippets_forked_from_id" ON "public"."code_snippets" ("forked_from_id");
-- Create index "idx_code_snippets_post_id" to table: "code_snippets"
CREATE INDEX "idx_code_snippets_post_id" ON "public"."code_snippets" ("post_id");
-- Create index "idx_code_snippets_user_id" to table: "code_snippets"
CREATE INDEX "idx_code_snippets_user_id" ON "public"."code_snippets" ("user_id");
-- Create "comments" table
CREATE TABLE "public"."comments" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "post_id" bigint NOT NULL,
  "parent_id" bigint NULL,
  "content" text NOT NULL,
  "like_count" bigint NULL DEFAULT 0,
  "is_hidden" boolean NULL DEFAULT false,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_comments_replies" FOREIGN KEY ("parent_id") REFERENCES "public"."comments" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_comments_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_comments_parent_id" to table: "comments"
CREATE INDEX "idx_comments_parent_id" ON "public"."comments" ("parent_id");
-- Create index "idx_comments_post_id" to table: "comments"
CREATE INDEX "idx_comments_post_id" ON "public"."comments" ("post_id");
-- Create index "idx_comments_user_id" to table: "comments"
CREATE INDEX "idx_comments_user_id" ON "public"."comments" ("user_id");
-- Create "follows" table
CREATE TABLE "public"."follows" (
  "id" bigserial NOT NULL,
  "follower_id" bigint NOT NULL,
  "followee_id" bigint NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_follows_followee" FOREIGN KEY ("followee_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_follows_follower" FOREIGN KEY ("follower_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_follower_following" to table: "follows"
CREATE UNIQUE INDEX "idx_follower_following" ON "public"."follows" ("follower_id", "followee_id");
-- Create "group_messages" table
CREATE TABLE "public"."group_messages" (
  "id" bigserial NOT NULL,
  "chat_room_id" bigint NOT NULL,
  "sender_id" bigint NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_group_messages_chat_room" FOREIGN KEY ("chat_room_id") REFERENCES "public"."chat_rooms" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_group_messages_sender" FOREIGN KEY ("sender_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_group_messages_chat_room_id" to table: "group_messages"
CREATE INDEX "idx_group_messages_chat_room_id" ON "public"."group_messages" ("chat_room_id");
-- Create index "idx_group_messages_sender_id" to table: "group_messages"
CREATE INDEX "idx_group_messages_sender_id" ON "public"."group_messages" ("sender_id");
-- Create "learning_resources" table
CREATE TABLE "public"."learning_resources" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "title" character varying(300) NOT NULL,
  "description" text NULL,
  "url" character varying(500) NULL,
  "category" character varying(50) NOT NULL,
  "difficulty" character varying(50) NULL,
  "tags" text NULL,
  "image_url" character varying(500) NULL,
  "is_public" boolean NULL,
  "like_count" bigint NULL DEFAULT 0,
  "save_count" bigint NULL DEFAULT 0,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_learning_resources_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_learning_resources_deleted_at" to table: "learning_resources"
CREATE INDEX "idx_learning_resources_deleted_at" ON "public"."learning_resources" ("deleted_at");
-- Create index "idx_learning_resources_user_id" to table: "learning_resources"
CREATE INDEX "idx_learning_resources_user_id" ON "public"."learning_resources" ("user_id");
-- Create "mentions" table
CREATE TABLE "public"."mentions" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "actor_id" bigint NOT NULL,
  "post_id" bigint NULL,
  "comment_id" bigint NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_mentions_actor" FOREIGN KEY ("actor_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_mentions_post" FOREIGN KEY ("post_id") REFERENCES "public"."posts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_mentions_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_mentions_comment_id" to table: "mentions"
CREATE INDEX "idx_mentions_comment_id" ON "public"."mentions" ("comment_id");
-- Create index "idx_mentions_comment_user" to table: "mentions"
CREATE UNIQUE INDEX "idx_mentions_comment_user" ON "public"."mentions" ("comment_id", "user_id") WHERE (comment_id IS NOT NULL);
-- Create index "idx_mentions_post_id" to table: "mentions"
CREATE INDEX "idx_mentions_post_id" ON "public"."mentions" ("post_id");
-- Create index "idx_mentions_post_user" to table: "mentions"
CREATE UNIQUE INDEX "idx_mentions_post_user" ON "public"."mentions" ("post_id", "user_id") WHERE ((post_id IS NOT NULL) AND (comment_id IS NULL));
-- Create index "idx_mentions_user_id" to table: "mentions"
CREATE INDEX "idx_mentions_user_id" ON "public"."mentions" ("user_id");
-- Create "messages" table
CREATE TABLE "public"."messages" (
  "id" bigserial NOT NULL,
  "sender_id" bigint NOT NULL,
  "receiver_id" bigint NOT NULL,
  "content" text NOT NULL,
  "read" boolean NULL DEFAULT false,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_messages_receiver" FOREIGN KEY ("receiver_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_messages_sender" FOREIGN KEY ("sender_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_messages_receiver_id" to table: "messages"
CREATE INDEX "idx_messages_receiver_id" ON "public"."messages" ("receiver_id");
-- Create index "idx_messages_sender_id" to table: "messages"
CREATE INDEX "idx_messages_sender_id" ON "public"."messages" ("sender_id");
-- Create "note_folders" table
CREATE TABLE "public"."note_folders" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "parent_id" bigint NULL,
  "name" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_note_folders_parent" FOREIGN KEY ("parent_id") REFERENCES "public"."note_folders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_note_folders_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_note_folders_parent_id" to table: "note_folders"
CREATE INDEX "idx_note_folders_parent_id" ON "public"."note_folders" ("parent_id");
-- Create index "idx_note_folders_user_id" to table: "note_folders"
CREATE INDEX "idx_note_folders_user_id" ON "public"."note_folders" ("user_id");
-- Create "notes" table
CREATE TABLE "public"."notes" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "folder_id" bigint NULL,
  "title" text NOT NULL,
  "content" text NULL,
  "tags" text NULL,
  "is_favorite" boolean NULL DEFAULT false,
  "is_archived" boolean NULL DEFAULT false,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_notes_folder" FOREIGN KEY ("folder_id") REFERENCES "public"."note_folders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_notes_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_notes_folder_id" to table: "notes"
CREATE INDEX "idx_notes_folder_id" ON "public"."notes" ("folder_id");
-- Create index "idx_notes_user_id" to table: "notes"
CREATE INDEX "idx_notes_user_id" ON "public"."notes" ("user_id");
-- Create "note_links" table
CREATE TABLE "public"."note_links" (
  "id" bigserial NOT NULL,
  "source_note_id" bigint NOT NULL,
  "target_note_id" bigint NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_note_links_source_note" FOREIGN KEY ("source_note_id") REFERENCES "public"."notes" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_note_links_target_note" FOREIGN KEY ("target_note_id") REFERENCES "public"."notes" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_note_link_unique" to table: "note_links"
CREATE UNIQUE INDEX "idx_note_link_unique" ON "public"."note_links" ("source_note_id", "target_note_id");
-- Create index "idx_note_links_source_note_id" to table: "note_links"
CREATE INDEX "idx_note_links_source_note_id" ON "public"."note_links" ("source_note_id");
-- Create index "idx_note_links_target_note_id" to table: "note_links"
CREATE INDEX "idx_note_links_target_note_id" ON "public"."note_links" ("target_note_id");
-- Create "notification_settings" table
CREATE TABLE "public"."notification_settings" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "enable_likes" boolean NULL DEFAULT true,
  "enable_comments" boolean NULL DEFAULT true,
  "enable_follows" boolean NULL DEFAULT true,
  "enable_messages" boolean NULL DEFAULT true,
  "enable_mentions" boolean NULL DEFAULT true,
  "enable_web_push" boolean NULL DEFAULT true,
  "enable_email" boolean NULL DEFAULT true,
  "enable_sound" boolean NULL DEFAULT true,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_notification_settings_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_notification_settings_user_id" to table: "notification_settings"
CREATE UNIQUE INDEX "idx_notification_settings_user_id" ON "public"."notification_settings" ("user_id");
-- Create "questions" table
CREATE TABLE "public"."questions" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "title" character varying(500) NOT NULL,
  "body" text NOT NULL,
  "tags" text NULL,
  "vote_count" bigint NULL DEFAULT 0,
  "answer_count" bigint NULL DEFAULT 0,
  "is_solved" boolean NULL DEFAULT false,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_questions_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_questions_deleted_at" to table: "questions"
CREATE INDEX "idx_questions_deleted_at" ON "public"."questions" ("deleted_at");
-- Create index "idx_questions_user_id" to table: "questions"
CREATE INDEX "idx_questions_user_id" ON "public"."questions" ("user_id");
-- Create "notifications" table
CREATE TABLE "public"."notifications" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "type" text NOT NULL,
  "actor_id" bigint NOT NULL,
  "post_id" bigint NULL,
  "question_id" bigint NULL,
  "badge_id" character varying(50) NULL,
  "read" boolean NULL DEFAULT false,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_notifications_actor" FOREIGN KEY ("actor_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_notifications_post" FOREIGN KEY ("post_id") REFERENCES "public"."posts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_notifications_question" FOREIGN KEY ("question_id") REFERENCES "public"."questions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_notifications_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_notifications_post_id" to table: "notifications"
CREATE INDEX "idx_notifications_post_id" ON "public"."notifications" ("post_id");
-- Create index "idx_notifications_question_id" to table: "notifications"
CREATE INDEX "idx_notifications_question_id" ON "public"."notifications" ("question_id");
-- Create index "idx_notifications_user_id" to table: "notifications"
CREATE INDEX "idx_notifications_user_id" ON "public"."notifications" ("user_id");
-- Create "password_reset_tokens" table
CREATE TABLE "public"."password_reset_tokens" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "token" text NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "used" boolean NULL DEFAULT false,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_password_reset_tokens_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_password_reset_tokens_token" to table: "password_reset_tokens"
CREATE UNIQUE INDEX "idx_password_reset_tokens_token" ON "public"."password_reset_tokens" ("token");
-- Create index "idx_password_reset_tokens_user_id" to table: "password_reset_tokens"
CREATE INDEX "idx_password_reset_tokens_user_id" ON "public"."password_reset_tokens" ("user_id");
-- Create "post_collection_items" table
CREATE TABLE "public"."post_collection_items" (
  "id" bigserial NOT NULL,
  "collection_id" bigint NOT NULL,
  "post_id" bigint NOT NULL,
  "note" text NULL,
  "order_index" bigint NOT NULL DEFAULT 0,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_post_collection_items_post" FOREIGN KEY ("post_id") REFERENCES "public"."posts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_post_collection_item" to table: "post_collection_items"
CREATE UNIQUE INDEX "idx_post_collection_item" ON "public"."post_collection_items" ("collection_id", "post_id");
-- Create index "idx_post_collection_items_collection_id" to table: "post_collection_items"
CREATE INDEX "idx_post_collection_items_collection_id" ON "public"."post_collection_items" ("collection_id");
-- Create index "idx_post_collection_items_post_id" to table: "post_collection_items"
CREATE INDEX "idx_post_collection_items_post_id" ON "public"."post_collection_items" ("post_id");
-- Create "post_collections" table
CREATE TABLE "public"."post_collections" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "title" text NOT NULL,
  "description" text NULL,
  "is_public" boolean NULL DEFAULT false,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_post_collections_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_post_collections_user_id" to table: "post_collections"
CREATE INDEX "idx_post_collections_user_id" ON "public"."post_collections" ("user_id");
-- Create "post_pins" table
CREATE TABLE "public"."post_pins" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "post_id" bigint NOT NULL,
  "pin_order" bigint NOT NULL DEFAULT 0,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_post_pins_post" FOREIGN KEY ("post_id") REFERENCES "public"."posts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_post_pins_post_id" to table: "post_pins"
CREATE INDEX "idx_post_pins_post_id" ON "public"."post_pins" ("post_id");
-- Create index "idx_post_pins_user_id" to table: "post_pins"
CREATE INDEX "idx_post_pins_user_id" ON "public"."post_pins" ("user_id");
-- Create index "idx_user_post_pin" to table: "post_pins"
CREATE UNIQUE INDEX "idx_user_post_pin" ON "public"."post_pins" ("user_id", "post_id");
-- Create "post_series" table
CREATE TABLE "public"."post_series" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "title" text NOT NULL,
  "description" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_post_series_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_post_series_user_id" to table: "post_series"
CREATE INDEX "idx_post_series_user_id" ON "public"."post_series" ("user_id");
-- Create "post_series_items" table
CREATE TABLE "public"."post_series_items" (
  "id" bigserial NOT NULL,
  "series_id" bigint NOT NULL,
  "post_id" bigint NOT NULL,
  "order_index" bigint NOT NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_post_series_items_post" FOREIGN KEY ("post_id") REFERENCES "public"."posts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_post_series_items_post_id" to table: "post_series_items"
CREATE INDEX "idx_post_series_items_post_id" ON "public"."post_series_items" ("post_id");
-- Create index "idx_post_series_items_series_id" to table: "post_series_items"
CREATE INDEX "idx_post_series_items_series_id" ON "public"."post_series_items" ("series_id");
-- Create index "idx_series_post" to table: "post_series_items"
CREATE UNIQUE INDEX "idx_series_post" ON "public"."post_series_items" ("series_id", "post_id");
-- Create "git_hub_repositories" table
CREATE TABLE "public"."git_hub_repositories" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "git_hub_repo_id" bigint NOT NULL,
  "name" text NOT NULL,
  "full_name" text NULL,
  "description" text NULL,
  "language" text NULL,
  "stars" bigint NULL,
  "forks" bigint NULL,
  "is_private" boolean NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_git_hub_repositories_git_hub_repo_id" to table: "git_hub_repositories"
CREATE UNIQUE INDEX "idx_git_hub_repositories_git_hub_repo_id" ON "public"."git_hub_repositories" ("git_hub_repo_id");
-- Create index "idx_git_hub_repositories_user_id" to table: "git_hub_repositories"
CREATE INDEX "idx_git_hub_repositories_user_id" ON "public"."git_hub_repositories" ("user_id");
-- Create "projects" table
CREATE TABLE "public"."projects" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "title" character varying(200) NOT NULL,
  "description" text NULL,
  "tech_stack" text NULL,
  "demo_url" character varying(500) NULL,
  "github_url" character varying(500) NULL,
  "image_url" character varying(500) NULL,
  "role" character varying(100) NULL,
  "start_date" timestamptz NULL,
  "end_date" timestamptz NULL,
  "featured" boolean NULL DEFAULT false,
  "is_archived" boolean NULL DEFAULT false,
  "github_repo_id" bigint NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_projects_github_repo" FOREIGN KEY ("github_repo_id") REFERENCES "public"."git_hub_repositories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_projects_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_projects_deleted_at" to table: "projects"
CREATE INDEX "idx_projects_deleted_at" ON "public"."projects" ("deleted_at");
-- Create index "idx_projects_user_id" to table: "projects"
CREATE INDEX "idx_projects_user_id" ON "public"."projects" ("user_id");
-- Create "question_bookmarks" table
CREATE TABLE "public"."question_bookmarks" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "question_id" bigint NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_question_bookmarks_question" FOREIGN KEY ("question_id") REFERENCES "public"."questions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_question_bookmark" to table: "question_bookmarks"
CREATE UNIQUE INDEX "idx_question_bookmark" ON "public"."question_bookmarks" ("user_id", "question_id");
-- Create "reminder_settings" table
CREATE TABLE "public"."reminder_settings" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "enabled" boolean NULL DEFAULT true,
  "frequency" text NULL DEFAULT 'daily',
  "notification_time" text NULL DEFAULT '09:00',
  "inactive_days" bigint NULL DEFAULT 3,
  "enable_web" boolean NULL DEFAULT true,
  "enable_email" boolean NULL DEFAULT false,
  "last_reminded_at" timestamptz NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_reminder_settings_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_reminder_settings_user_id" to table: "reminder_settings"
CREATE UNIQUE INDEX "idx_reminder_settings_user_id" ON "public"."reminder_settings" ("user_id");
-- Create "resource_progresses" table
CREATE TABLE "public"."resource_progresses" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "resource_id" bigint NOT NULL,
  "status" character varying(20) NULL DEFAULT 'not_started',
  "completion_percent" bigint NULL DEFAULT 0,
  "note" text NULL,
  "started_at" timestamptz NULL,
  "completed_at" timestamptz NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_resource_progresses_resource" FOREIGN KEY ("resource_id") REFERENCES "public"."learning_resources" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_resource_progress" to table: "resource_progresses"
CREATE UNIQUE INDEX "idx_resource_progress" ON "public"."resource_progresses" ("user_id", "resource_id");
-- Create "resource_reviews" table
CREATE TABLE "public"."resource_reviews" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "resource_id" bigint NOT NULL,
  "rating" bigint NOT NULL,
  "comment" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_resource_reviews_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_resource_review" to table: "resource_reviews"
CREATE UNIQUE INDEX "idx_resource_review" ON "public"."resource_reviews" ("user_id", "resource_id");
-- Create index "idx_resource_reviews_resource_id" to table: "resource_reviews"
CREATE INDEX "idx_resource_reviews_resource_id" ON "public"."resource_reviews" ("resource_id");
-- Create "roadmaps" table
CREATE TABLE "public"."roadmaps" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "title" character varying(200) NOT NULL,
  "description" text NULL,
  "category" text NULL DEFAULT 'other',
  "is_public" boolean NULL DEFAULT false,
  "is_template" boolean NULL DEFAULT false,
  "step_count" bigint NULL DEFAULT 0,
  "completed_step_count" bigint NULL DEFAULT 0,
  "progress" bigint NULL DEFAULT 0,
  "status" text NULL DEFAULT 'active',
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "completed_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_roadmaps_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_roadmaps_is_public" to table: "roadmaps"
CREATE INDEX "idx_roadmaps_is_public" ON "public"."roadmaps" ("is_public");
-- Create index "idx_roadmaps_is_template" to table: "roadmaps"
CREATE INDEX "idx_roadmaps_is_template" ON "public"."roadmaps" ("is_template");
-- Create index "idx_roadmaps_user_id" to table: "roadmaps"
CREATE INDEX "idx_roadmaps_user_id" ON "public"."roadmaps" ("user_id");
-- Create "roadmap_steps" table
CREATE TABLE "public"."roadmap_steps" (
  "id" bigserial NOT NULL,
  "roadmap_id" bigint NOT NULL,
  "title" character varying(200) NOT NULL,
  "description" text NULL,
  "order_index" bigint NOT NULL DEFAULT 0,
  "is_completed" boolean NULL DEFAULT false,
  "completed_at" timestamptz NULL,
  "resource_url" character varying(500) NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_roadmaps_steps" FOREIGN KEY ("roadmap_id") REFERENCES "public"."roadmaps" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_roadmap_steps_roadmap_id" to table: "roadmap_steps"
CREATE INDEX "idx_roadmap_steps_roadmap_id" ON "public"."roadmap_steps" ("roadmap_id");
-- Create "snippet_comments" table
CREATE TABLE "public"."snippet_comments" (
  "id" bigserial NOT NULL,
  "snippet_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "line_number" bigint NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_snippet_comments_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_snippet_comments_snippet_id" to table: "snippet_comments"
CREATE INDEX "idx_snippet_comments_snippet_id" ON "public"."snippet_comments" ("snippet_id");
-- Create index "idx_snippet_comments_user_id" to table: "snippet_comments"
CREATE INDEX "idx_snippet_comments_user_id" ON "public"."snippet_comments" ("user_id");
-- Create "study_circle_checkins" table
CREATE TABLE "public"."study_circle_checkins" (
  "id" bigserial NOT NULL,
  "circle_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "date" character varying(10) NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_study_circle_checkins_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_checkin_unique" to table: "study_circle_checkins"
CREATE UNIQUE INDEX "idx_checkin_unique" ON "public"."study_circle_checkins" ("circle_id", "user_id", "date");
-- Create "study_circle_member_progresses" table
CREATE TABLE "public"."study_circle_member_progresses" (
  "id" bigserial NOT NULL,
  "circle_id" bigint NOT NULL,
  "step_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "is_completed" boolean NULL DEFAULT false,
  "completed_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_study_circle_member_progresses_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_circle_step_user" to table: "study_circle_member_progresses"
CREATE UNIQUE INDEX "idx_circle_step_user" ON "public"."study_circle_member_progresses" ("circle_id", "step_id", "user_id");
-- Create "study_circles" table
CREATE TABLE "public"."study_circles" (
  "id" bigserial NOT NULL,
  "name" character varying(200) NOT NULL,
  "topic" character varying(200) NOT NULL,
  "description" text NULL,
  "owner_id" bigint NOT NULL,
  "max_members" bigint NOT NULL DEFAULT 5,
  "status" text NULL DEFAULT 'active',
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_study_circles_owner" FOREIGN KEY ("owner_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_study_circles_owner_id" to table: "study_circles"
CREATE INDEX "idx_study_circles_owner_id" ON "public"."study_circles" ("owner_id");
-- Create "study_circle_members" table
CREATE TABLE "public"."study_circle_members" (
  "id" bigserial NOT NULL,
  "circle_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "role" text NULL DEFAULT 'member',
  "joined_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_study_circle_members_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_study_circles_members" FOREIGN KEY ("circle_id") REFERENCES "public"."study_circles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_circle_user" to table: "study_circle_members"
CREATE UNIQUE INDEX "idx_circle_user" ON "public"."study_circle_members" ("circle_id", "user_id");
-- Create "study_circle_steps" table
CREATE TABLE "public"."study_circle_steps" (
  "id" bigserial NOT NULL,
  "circle_id" bigint NOT NULL,
  "title" character varying(200) NOT NULL,
  "description" text NULL,
  "order_index" bigint NULL DEFAULT 0,
  "resource_url" character varying(2000) NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_study_circles_steps" FOREIGN KEY ("circle_id") REFERENCES "public"."study_circles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_study_circle_steps_circle_id" to table: "study_circle_steps"
CREATE INDEX "idx_study_circle_steps_circle_id" ON "public"."study_circle_steps" ("circle_id");
