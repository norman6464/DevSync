# DevSync のスキーマ正本（Atlas 宣言的スキーマ）。
#
# このファイルが唯一の正本。schema_generated.sql はこのファイルから
# `make db-schema-sql`（backend/Makefile）が機械生成する、DO NOT EDIT な副産物。
#
# 元は internal/infra/database/schema/*.hcl としてドメインごとに41ファイルへ分割していたが、
# 姉妹プロダクトFreStyleの単一ファイル運用（internal/infra/database/schema/schema.hcl）に
# 合わせて統合した。以下の見出しコメントは、統合前のファイル名（ドメイン単位）をそのまま
# 節区切りとして残したもの。

schema "public" {
  comment = "standard public schema"
}

# =====================================================================
# ai_advice
# =====================================================================

table "ai_advices" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "type" {
    null = false
    type = text
  }
  column "priority" {
    null    = false
    type    = bigint
    default = 3
  }
  column "title_key" {
    null = false
    type = text
  }
  column "message_key" {
    null = false
    type = text
  }
  column "params" {
    null = true
    type = text
  }
  column "action_url" {
    null = true
    type = text
  }
  column "is_read" {
    null    = false
    type    = boolean
    default = false
  }
  column "expires_at" {
    null = true
    type = timestamptz
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_ai_advices_user_id" {
    columns = [column.user_id]
  }
  foreign_key "fk_ai_advices_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  check "ck_ai_advices_type" {
    expr = "type IN ('streak_recovery', 'stalled_roadmap', 'goal_overdue', 'tech_suggestion', 'goal_suggestion', 'difficulty_up', 'praise', 'general')"
  }
  check "ck_ai_advices_priority" {
    expr = "priority BETWEEN 1 AND 5"
  }
}
table "ai_conversations" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "title" {
    null = false
    type = character_varying(200)
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_ai_conversations_user_id" {
    columns = [column.user_id]
  }
  foreign_key "fk_ai_conversations_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "ai_messages" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "conversation_id" {
    null = false
    type = bigint
  }
  column "role" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "tokens_used" {
    null    = true
    type    = bigint
    default = 0
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_ai_conversations_messages" {
    columns     = [column.conversation_id]
    ref_columns = [table.ai_conversations.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_ai_messages_conversation_id" {
    columns = [column.conversation_id]
  }
  check "ck_ai_messages_role" {
    expr = "role IN ('user', 'assistant', 'system')"
  }
}

# =====================================================================
# answer
# =====================================================================

table "answer_votes" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "answer_id" {
    null = false
    type = bigint
  }
  column "value" {
    null = false
    type = bigint
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "uq_answer_vote" {
    unique  = true
    columns = [column.user_id, column.answer_id]
  }
  foreign_key "fk_answer_votes_answer" {
    columns     = [column.answer_id]
    ref_columns = [table.answers.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "fk_answer_votes_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  check "ck_answer_votes_value" {
    expr = "value IN (1, -1)"
  }
}
table "answers" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "question_id" {
    null = false
    type = bigint
  }
  column "body" {
    null = false
    type = text
  }
  column "vote_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "is_best" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_answers_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_answers_question_id" {
    columns = [column.question_id]
  }
  index "idx_answers_user_id" {
    columns = [column.user_id]
  }
  foreign_key "fk_answers_question" {
    columns     = [column.question_id]
    ref_columns = [table.questions.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}

# =====================================================================
# book_review
# =====================================================================

table "book_reviews" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "title" {
    null = false
    type = character_varying(300)
  }
  column "author" {
    null = true
    type = character_varying(200)
  }
  column "isbn" {
    null = true
    type = character_varying(20)
  }
  column "rating" {
    null = false
    type = bigint
  }
  column "review" {
    null = true
    type = text
  }
  column "total_pages" {
    null    = true
    type    = bigint
    default = 0
  }
  column "current_page" {
    null    = true
    type    = bigint
    default = 0
  }
  column "image_url" {
    null = true
    type = character_varying(2000)
  }
  column "status" {
    null    = true
    type    = text
    default = "not_started"
  }
  column "is_archived" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_book_reviews_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_book_reviews_user_id" {
    columns = [column.user_id]
  }
  check "ck_book_reviews_rating" {
    expr = "rating BETWEEN 1 AND 5"
  }
  check "ck_book_reviews_status" {
    expr = "status IN ('not_started', 'reading', 'completed')"
  }
}

# =====================================================================
# bookmark_collection
# =====================================================================

table "bookmark_collection_items" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "collection_id" {
    null = false
    type = bigint
  }
  column "post_id" {
    null = false
    type = bigint
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_bookmark_collection_items_post" {
    columns     = [column.post_id]
    ref_columns = [table.posts.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_bookmark_collection_items_post_id" {
    columns = [column.post_id]
  }
  index "uq_bookmark_collection_post" {
    unique  = true
    columns = [column.collection_id, column.post_id]
  }
  foreign_key "fk_bookmark_collection_items_collection" {
    columns     = [column.collection_id]
    ref_columns = [table.bookmark_collections.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "bookmark_collections" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "name" {
    null = false
    type = text
  }
  column "description" {
    null = true
    type = text
  }
  column "color" {
    null    = true
    type    = text
    default = "blue"
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_bookmark_collections_user_id" {
    columns = [column.user_id]
  }
  foreign_key "fk_bookmark_collections_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}

# =====================================================================
# chat_room
# =====================================================================

table "chat_room_members" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "chat_room_id" {
    null = false
    type = bigint
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "joined_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_chat_room_members_chat_room" {
    columns     = [column.chat_room_id]
    ref_columns = [table.chat_rooms.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "fk_chat_room_members_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_chat_room_members_chat_room_id" {
    columns = [column.chat_room_id]
  }
  index "idx_chat_room_members_user_id" {
    columns = [column.user_id]
  }
  index "uq_room_user" {
    unique  = true
    columns = [column.chat_room_id, column.user_id]
  }
}
table "chat_rooms" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "name" {
    null = false
    type = character_varying(100)
  }
  column "description" {
    null = true
    type = character_varying(500)
  }
  column "owner_id" {
    null = false
    type = bigint
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_chat_rooms_owner" {
    columns     = [column.owner_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_chat_rooms_owner_id" {
    columns = [column.owner_id]
  }
}
table "group_messages" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "chat_room_id" {
    null = false
    type = bigint
  }
  column "sender_id" {
    null = false
    type = bigint
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_group_messages_chat_room" {
    columns     = [column.chat_room_id]
    ref_columns = [table.chat_rooms.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "fk_group_messages_sender" {
    columns     = [column.sender_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_group_messages_chat_room_id" {
    columns = [column.chat_room_id]
  }
  index "idx_group_messages_sender_id" {
    columns = [column.sender_id]
  }
}

# =====================================================================
# code_snippet
# =====================================================================

table "code_snippet_favorites" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "snippet_id" {
    null = false
    type = bigint
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "uq_snippet_favorite" {
    unique  = true
    columns = [column.user_id, column.snippet_id]
  }
  foreign_key "fk_code_snippet_favorites_snippet" {
    columns     = [column.snippet_id]
    ref_columns = [table.code_snippets.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "fk_code_snippet_favorites_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "code_snippets" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "post_id" {
    null = false
    type = bigint
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "language" {
    null = false
    type = text
  }
  column "file_name" {
    null = true
    type = text
  }
  column "code" {
    null = false
    type = text
  }
  column "comment_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "forked_from_id" {
    null = true
    type = bigint
  }
  column "fork_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_code_snippets" {
    columns     = [column.post_id]
    ref_columns = [table.posts.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_code_snippets_forked_from_id" {
    columns = [column.forked_from_id]
  }
  index "idx_code_snippets_post_id" {
    columns = [column.post_id]
  }
  index "idx_code_snippets_user_id" {
    columns = [column.user_id]
  }
  foreign_key "fk_code_snippets_forked_from" {
    columns     = [column.forked_from_id]
    ref_columns = [table.code_snippets.column.id]
    on_update   = NO_ACTION
    on_delete   = SET_NULL
  }
  foreign_key "fk_code_snippets_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "snippet_comments" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "snippet_id" {
    null = false
    type = bigint
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "line_number" {
    null = false
    type = bigint
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_snippet_comments_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_snippet_comments_snippet_id" {
    columns = [column.snippet_id]
  }
  index "idx_snippet_comments_user_id" {
    columns = [column.user_id]
  }
  foreign_key "fk_snippet_comments_snippet" {
    columns     = [column.snippet_id]
    ref_columns = [table.code_snippets.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}

# =====================================================================
# follow
# =====================================================================

table "follows" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "follower_id" {
    null = false
    type = bigint
  }
  column "followee_id" {
    null = false
    type = bigint
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_follows_followee" {
    columns     = [column.followee_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "fk_follows_follower" {
    columns     = [column.follower_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "uq_follower_following" {
    unique  = true
    columns = [column.follower_id, column.followee_id]
  }
  check "ck_follows_not_self" {
    expr = "follower_id <> followee_id"
  }
}

# =====================================================================
# github
# =====================================================================

table "git_hub_contributions" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "contributed_on" {
    null = false
    type = date
  }
  column "count" {
    null    = false
    type    = bigint
    default = 0
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "uq_git_hub_contributions_user_date" {
    unique  = true
    columns = [column.user_id, column.contributed_on]
  }
  foreign_key "fk_git_hub_contributions_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "git_hub_language_stats" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "language" {
    null = false
    type = text
  }
  column "bytes" {
    null    = false
    type    = bigint
    default = 0
  }
  column "repo_count" {
    null    = false
    type    = bigint
    default = 0
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "uq_user_lang" {
    unique  = true
    columns = [column.user_id, column.language]
  }
  foreign_key "fk_git_hub_language_stats_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "git_hub_repositories" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "git_hub_repo_id" {
    null = false
    type = bigint
  }
  column "name" {
    null = false
    type = text
  }
  column "full_name" {
    null = true
    type = text
  }
  column "description" {
    null = true
    type = text
  }
  column "language" {
    null = true
    type = text
  }
  column "stars" {
    null = true
    type = bigint
  }
  column "forks" {
    null = true
    type = bigint
  }
  column "is_private" {
    null = false
    type = boolean
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "uq_git_hub_repositories_user_repo" {
    unique  = true
    columns = [column.user_id, column.git_hub_repo_id]
  }
  index "idx_git_hub_repositories_user_id" {
    columns = [column.user_id]
  }
  foreign_key "fk_git_hub_repositories_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}

# =====================================================================
# learning_goal
# =====================================================================

table "learning_goals" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "title" {
    null = false
    type = text
  }
  column "description" {
    null = true
    type = text
  }
  column "category" {
    null    = true
    type    = text
    default = "other"
  }
  column "target_date" {
    null = true
    type = timestamptz
  }
  column "progress" {
    null    = true
    type    = bigint
    default = 0
  }
  column "target_hours" {
    null    = true
    type    = bigint
    default = 0
  }
  column "status" {
    null    = true
    type    = text
    default = "active"
  }
  column "is_public" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  column "completed_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_learning_goals_user_id" {
    columns = [column.user_id]
  }
  foreign_key "fk_learning_goals_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  check "ck_learning_goals_progress" {
    expr = "progress BETWEEN 0 AND 100"
  }
  check "ck_learning_goals_status" {
    expr = "status IN ('active', 'completed', 'paused')"
  }
  check "ck_learning_goals_category" {
    expr = "category IN ('language', 'framework', 'skill', 'project', 'other')"
  }
}

# =====================================================================
# learning_log
# =====================================================================

table "learning_logs" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "category" {
    null    = true
    type    = text
    default = "other"
  }
  column "duration" {
    null    = true
    type    = bigint
    default = 0
  }
  column "goal_id" {
    null = true
    type = bigint
  }
  column "source" {
    null    = true
    type    = text
    default = "manual"
  }
  column "is_favorite" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_learning_logs_goal_id" {
    columns = [column.goal_id]
  }
  index "idx_learning_logs_user_id" {
    columns = [column.user_id]
  }
  foreign_key "fk_learning_logs_goal" {
    columns     = [column.goal_id]
    ref_columns = [table.learning_goals.column.id]
    on_update   = NO_ACTION
    on_delete   = SET_NULL
  }
  foreign_key "fk_learning_logs_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  check "ck_learning_logs_category" {
    expr = "category IN ('coding', 'reading', 'course', 'meetup', 'other')"
  }
}

# =====================================================================
# learning_log_template
# =====================================================================

table "learning_log_templates" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "name" {
    null = false
    type = character_varying(100)
  }
  column "default_title" {
    null = true
    type = character_varying(200)
  }
  column "default_content" {
    null = true
    type = text
  }
  column "default_category" {
    null    = true
    type    = character_varying(50)
    default = "other"
  }
  column "default_duration" {
    null    = true
    type    = bigint
    default = 0
  }
  column "is_default" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_learning_log_templates_user_id" {
    columns = [column.user_id]
  }
  index "idx_log_templates_is_default" {
    columns = [column.is_default]
  }
  index "uq_learning_log_templates_default" {
    unique  = true
    columns = [column.user_id]
    where   = "is_default"
  }
  foreign_key "fk_learning_log_templates_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}

# =====================================================================
# learning_resource
# =====================================================================

table "learning_resources" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "title" {
    null = false
    type = character_varying(300)
  }
  column "description" {
    null = true
    type = text
  }
  column "url" {
    null = true
    type = character_varying(500)
  }
  column "category" {
    null = false
    type = character_varying(50)
  }
  column "difficulty" {
    null = true
    type = character_varying(50)
  }
  column "tags" {
    null = true
    type = text
  }
  column "image_url" {
    null = true
    type = character_varying(500)
  }
  column "is_public" {
    null = false
    type = boolean
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_learning_resources_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_learning_resources_user_id" {
    columns = [column.user_id]
  }
  check "ck_learning_resources_category" {
    expr = "category IN ('book', 'video', 'article', 'course', 'tutorial', 'podcast', 'tool', 'other')"
  }
  check "ck_learning_resources_difficulty" {
    expr = "difficulty IN ('beginner', 'intermediate', 'advanced')"
  }
}
table "learning_resource_metrics" {
  schema = schema.public
  # 一覧のORDER BYに使うlike_count/save_countだけをここへ出す（DEVSYNC-159、
  # postsのpost_metricsと同じ設計）。ファクトのINSERT/DELETEと同一トランザクションの
  # UPSERTで更新し、夜次reconcileジョブで実カウントとの乖離を補正する。行は
  # 最初のいいね/保存で遅延生成される（読み取り側はLEFT JOIN + COALESCEで0扱いにする）。
  column "resource_id" {
    null = false
    type = bigint
  }
  column "like_count" {
    null    = false
    type    = bigint
    default = 0
  }
  column "save_count" {
    null    = false
    type    = bigint
    default = 0
  }
  primary_key {
    columns = [column.resource_id]
  }
  foreign_key "fk_learning_resource_metrics_resource" {
    columns     = [column.resource_id]
    ref_columns = [table.learning_resources.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "resource_likes" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "resource_id" {
    null = false
    type = bigint
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "uq_resource_like" {
    unique  = true
    columns = [column.user_id, column.resource_id]
  }
  foreign_key "fk_resource_likes_resource" {
    columns     = [column.resource_id]
    ref_columns = [table.learning_resources.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "fk_resource_likes_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "resource_reviews" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "resource_id" {
    null = false
    type = bigint
  }
  column "rating" {
    null = false
    type = bigint
  }
  column "comment" {
    null = true
    type = text
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_resource_reviews_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "uq_resource_review" {
    unique  = true
    columns = [column.user_id, column.resource_id]
  }
  index "idx_resource_reviews_resource_id" {
    columns = [column.resource_id]
  }
  foreign_key "fk_resource_reviews_resource" {
    columns     = [column.resource_id]
    ref_columns = [table.learning_resources.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  check "ck_resource_reviews_rating" {
    expr = "rating BETWEEN 1 AND 5"
  }
}
table "resource_saves" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "resource_id" {
    null = false
    type = bigint
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "uq_resource_save" {
    unique  = true
    columns = [column.user_id, column.resource_id]
  }
  foreign_key "fk_resource_saves_resource" {
    columns     = [column.resource_id]
    ref_columns = [table.learning_resources.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "fk_resource_saves_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}

# =====================================================================
# mention
# =====================================================================

table "mentions" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "actor_id" {
    null = false
    type = bigint
  }
  column "post_id" {
    null = true
    type = bigint
  }
  column "comment_id" {
    null = true
    type = bigint
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_mentions_actor" {
    columns     = [column.actor_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "fk_mentions_post" {
    columns     = [column.post_id]
    ref_columns = [table.posts.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "fk_mentions_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_mentions_comment_id" {
    columns = [column.comment_id]
  }
  index "uq_mentions_comment_user" {
    unique  = true
    columns = [column.comment_id, column.user_id]
    where   = "(comment_id IS NOT NULL)"
  }
  index "idx_mentions_post_id" {
    columns = [column.post_id]
  }
  index "uq_mentions_post_user" {
    unique  = true
    columns = [column.post_id, column.user_id]
    where   = "((post_id IS NOT NULL) AND (comment_id IS NULL))"
  }
  index "idx_mentions_user_id" {
    columns = [column.user_id]
  }
  foreign_key "fk_mentions_comment" {
    columns     = [column.comment_id]
    ref_columns = [table.comments.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  check "ck_mentions_not_self" {
    expr = "user_id <> actor_id"
  }
}

# =====================================================================
# message
# =====================================================================

table "messages" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "sender_id" {
    null = false
    type = bigint
  }
  column "receiver_id" {
    null = false
    type = bigint
  }
  column "content" {
    null = false
    type = text
  }
  column "read" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_messages_receiver" {
    columns     = [column.receiver_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "fk_messages_sender" {
    columns     = [column.sender_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_messages_receiver_id" {
    columns = [column.receiver_id]
  }
  index "idx_messages_sender_id" {
    columns = [column.sender_id]
  }
  check "ck_messages_not_self" {
    expr = "sender_id <> receiver_id"
  }
}

# =====================================================================
# note
# =====================================================================

table "note_folders" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "parent_id" {
    null = true
    type = bigint
  }
  column "name" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_note_folders_parent" {
    columns     = [column.parent_id]
    ref_columns = [table.note_folders.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "fk_note_folders_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_note_folders_parent_id" {
    columns = [column.parent_id]
  }
  index "idx_note_folders_user_id" {
    columns = [column.user_id]
  }
}
table "notes" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "folder_id" {
    null = true
    type = bigint
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = true
    type = text
  }
  column "tags" {
    null = true
    type = text
  }
  column "is_favorite" {
    null    = false
    type    = boolean
    default = false
  }
  column "is_archived" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_notes_folder" {
    columns     = [column.folder_id]
    ref_columns = [table.note_folders.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "fk_notes_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_notes_folder_id" {
    columns = [column.folder_id]
  }
  index "idx_notes_user_id" {
    columns = [column.user_id]
  }
}

# =====================================================================
# note_link
# =====================================================================

table "note_links" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "source_note_id" {
    null = false
    type = bigint
  }
  column "target_note_id" {
    null = false
    type = bigint
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_note_links_source_note" {
    columns     = [column.source_note_id]
    ref_columns = [table.notes.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "fk_note_links_target_note" {
    columns     = [column.target_note_id]
    ref_columns = [table.notes.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "uq_note_link_unique" {
    unique  = true
    columns = [column.source_note_id, column.target_note_id]
  }
  index "idx_note_links_source_note_id" {
    columns = [column.source_note_id]
  }
  index "idx_note_links_target_note_id" {
    columns = [column.target_note_id]
  }
  check "ck_note_links_not_self" {
    expr = "source_note_id <> target_note_id"
  }
}

# =====================================================================
# note_template
# =====================================================================

table "note_templates" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "name" {
    null = false
    type = character_varying(100)
  }
  column "description" {
    null = true
    type = character_varying(500)
  }
  column "default_title" {
    null = true
    type = character_varying(200)
  }
  column "content_template" {
    null = false
    type = text
  }
  column "default_tags" {
    null = true
    type = character_varying(255)
  }
  column "is_default" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_note_templates_is_default" {
    columns = [column.is_default]
  }
  index "idx_note_templates_user_id" {
    columns = [column.user_id]
  }
  index "uq_note_templates_default" {
    unique  = true
    columns = [column.user_id]
    where   = "is_default"
  }
  foreign_key "fk_note_templates_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}

# =====================================================================
# note_version
# =====================================================================

table "note_versions" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "note_id" {
    null = false
    type = bigint
  }
  column "version_number" {
    null = false
    type = bigint
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = true
    type = text
  }
  column "tags" {
    null = true
    type = text
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_note_versions_note_id" {
    columns = [column.note_id]
  }
  foreign_key "fk_note_versions_note" {
    columns     = [column.note_id]
    ref_columns = [table.notes.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}

# =====================================================================
# notification
# =====================================================================

table "notifications" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "type" {
    null = false
    type = text
  }
  column "actor_id" {
    null = false
    type = bigint
  }
  column "post_id" {
    null = true
    type = bigint
  }
  column "question_id" {
    null = true
    type = bigint
  }
  column "badge_id" {
    null = true
    type = character_varying(50)
  }
  column "read" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_notifications_actor" {
    columns     = [column.actor_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "fk_notifications_post" {
    columns     = [column.post_id]
    ref_columns = [table.posts.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "fk_notifications_question" {
    columns     = [column.question_id]
    ref_columns = [table.questions.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "fk_notifications_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "fk_notifications_type" {
    columns     = [column.type]
    ref_columns = [table.notification_verbs.column.code]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_notifications_post_id" {
    columns = [column.post_id]
  }
  index "idx_notifications_question_id" {
    columns = [column.question_id]
  }
  index "idx_notifications_user_id" {
    columns = [column.user_id]
  }
  check "ck_notifications_exclusive_target" {
    expr = "num_nonnulls(post_id, question_id, badge_id) <= 1"
  }
}

table "notification_verbs" {
  schema = schema.public
  column "code" {
    null = false
    type = text
  }
  primary_key {
    columns = [column.code]
  }
}

# =====================================================================
# notification_settings
# =====================================================================

table "notification_settings" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "enable_likes" {
    null    = false
    type    = boolean
    default = true
  }
  column "enable_comments" {
    null    = false
    type    = boolean
    default = true
  }
  column "enable_follows" {
    null    = false
    type    = boolean
    default = true
  }
  column "enable_messages" {
    null    = false
    type    = boolean
    default = true
  }
  column "enable_mentions" {
    null    = false
    type    = boolean
    default = true
  }
  column "enable_web_push" {
    null    = false
    type    = boolean
    default = true
  }
  column "enable_email" {
    null    = false
    type    = boolean
    default = true
  }
  column "enable_sound" {
    null    = false
    type    = boolean
    default = true
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_notification_settings_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "uq_notification_settings_user_id" {
    unique  = true
    columns = [column.user_id]
  }
}

# =====================================================================
# password_reset
# =====================================================================

table "password_reset_tokens" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "token" {
    null = false
    type = text
  }
  column "expires_at" {
    null = false
    type = timestamptz
  }
  column "used" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_password_reset_tokens_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "uq_password_reset_tokens_token" {
    unique  = true
    columns = [column.token]
  }
  index "idx_password_reset_tokens_user_id" {
    columns = [column.user_id]
  }
}

# =====================================================================
# post
# =====================================================================

table "comment_likes" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "comment_id" {
    null = false
    type = bigint
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_comment_likes_comment_id" {
    columns = [column.comment_id]
  }
  index "uq_user_comment_like" {
    unique  = true
    columns = [column.user_id, column.comment_id]
  }
  foreign_key "fk_comment_likes_comment" {
    columns     = [column.comment_id]
    ref_columns = [table.comments.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "fk_comment_likes_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "comments" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "post_id" {
    null = false
    type = bigint
  }
  column "parent_id" {
    null = true
    type = bigint
  }
  column "content" {
    null = false
    type = text
  }
  column "like_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "is_hidden" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_comments_replies" {
    columns     = [column.parent_id]
    ref_columns = [table.comments.column.id]
    on_update   = NO_ACTION
    on_delete   = SET_NULL
  }
  foreign_key "fk_comments_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_comments_parent_id" {
    columns = [column.parent_id]
  }
  index "idx_comments_post_id" {
    columns = [column.post_id]
  }
  index "idx_comments_user_id" {
    columns = [column.user_id]
  }
  foreign_key "fk_comments_post" {
    columns     = [column.post_id]
    ref_columns = [table.posts.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "post_reactions" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "post_id" {
    null = false
    type = bigint
  }
  column "kind" {
    null = false
    type = text
  }
  column "value" {
    null    = false
    type    = text
    default = ""
  }
  column "pin_order" {
    null = true
    type = bigint
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_post_reactions_post_id" {
    columns = [column.post_id]
  }
  index "idx_post_reactions_user_id" {
    columns = [column.user_id]
  }
  index "uq_post_reactions_user_post_kind_value" {
    unique  = true
    columns = [column.user_id, column.post_id, column.kind, column.value]
  }
  foreign_key "fk_post_reactions_post" {
    columns     = [column.post_id]
    ref_columns = [table.posts.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "fk_post_reactions_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  check "ck_post_reactions_kind" {
    expr = "kind IN ('like', 'bookmark', 'pin', 'view', 'emoji')"
  }
}
table "post_collection_items" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "collection_id" {
    null = false
    type = bigint
  }
  column "post_id" {
    null = false
    type = bigint
  }
  column "note" {
    null = true
    type = text
  }
  column "order_index" {
    null    = false
    type    = bigint
    default = 0
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_post_collection_items_post" {
    columns     = [column.post_id]
    ref_columns = [table.posts.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "uq_post_collection_item" {
    unique  = true
    columns = [column.collection_id, column.post_id]
  }
  index "idx_post_collection_items_collection_id" {
    columns = [column.collection_id]
  }
  index "idx_post_collection_items_post_id" {
    columns = [column.post_id]
  }
  foreign_key "fk_post_collection_items_collection" {
    columns     = [column.collection_id]
    ref_columns = [table.post_collections.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "post_collections" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "title" {
    null = false
    type = text
  }
  column "description" {
    null = true
    type = text
  }
  column "is_public" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_post_collections_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_post_collections_user_id" {
    columns = [column.user_id]
  }
}
table "post_series" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "title" {
    null = false
    type = text
  }
  column "description" {
    null = true
    type = text
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_post_series_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_post_series_user_id" {
    columns = [column.user_id]
  }
}
table "post_series_items" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "series_id" {
    null = false
    type = bigint
  }
  column "post_id" {
    null = false
    type = bigint
  }
  column "order_index" {
    null    = false
    type    = bigint
    default = 0
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_post_series_items_post" {
    columns     = [column.post_id]
    ref_columns = [table.posts.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_post_series_items_post_id" {
    columns = [column.post_id]
  }
  index "idx_post_series_items_series_id" {
    columns = [column.series_id]
  }
  index "uq_series_post" {
    unique  = true
    columns = [column.series_id, column.post_id]
  }
  foreign_key "fk_post_series_items_series" {
    columns     = [column.series_id]
    ref_columns = [table.post_series.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "post_tags" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "post_id" {
    null = false
    type = bigint
  }
  column "tag" {
    null = false
    type = character_varying(50)
  }
  primary_key {
    columns = [column.id]
  }
  index "uq_post_tag" {
    unique  = true
    columns = [column.post_id, column.tag]
  }
  index "idx_post_tags_post_id" {
    columns = [column.post_id]
  }
  index "idx_post_tags_tag" {
    columns = [column.tag]
  }
  foreign_key "fk_post_tags_post" {
    columns     = [column.post_id]
    ref_columns = [table.posts.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "posts" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "image_urls" {
    null = true
    type = text
  }
  column "is_draft" {
    null    = false
    type    = boolean
    default = false
  }
  column "estimated_read_time" {
    null    = true
    type    = bigint
    default = 0
  }
  column "scheduled_at" {
    null = true
    type = timestamptz
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_posts_is_draft" {
    columns = [column.is_draft]
  }
  index "idx_posts_scheduled_at" {
    columns = [column.scheduled_at]
  }
  index "idx_posts_user_id" {
    columns = [column.user_id]
  }
}
table "post_metrics" {
  schema = schema.public
  # 一覧のORDER BYに使うlike_count/comment_count/view_countだけをここへ出す
  # （DEVSYNC-159）。ファクトのINSERT/DELETEと同一トランザクションのUPSERTで更新し、
  # 夜次reconcileジョブ（ReconcileAllPostMetrics）で実カウントとの乖離を補正する。
  # 行は投稿作成時ではなく最初のいいね/コメント/閲覧で遅延生成される
  # （読み取り側は全てLEFT JOIN + COALESCEで0扱いにする）。
  column "post_id" {
    null = false
    type = bigint
  }
  column "like_count" {
    null    = false
    type    = bigint
    default = 0
  }
  column "comment_count" {
    null    = false
    type    = bigint
    default = 0
  }
  column "view_count" {
    null    = false
    type    = bigint
    default = 0
  }
  primary_key {
    columns = [column.post_id]
  }
  foreign_key "fk_post_metrics_post" {
    columns     = [column.post_id]
    ref_columns = [table.posts.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
# =====================================================================
# post_template
# =====================================================================

table "post_templates" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "name" {
    null = false
    type = character_varying(100)
  }
  column "title_template" {
    null = true
    type = character_varying(200)
  }
  column "content_template" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_post_templates_user_id" {
    columns = [column.user_id]
  }
  foreign_key "fk_post_templates_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}

# =====================================================================
# project
# =====================================================================

table "projects" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "title" {
    null = false
    type = character_varying(200)
  }
  column "description" {
    null = true
    type = text
  }
  column "tech_stack" {
    null = true
    type = text
  }
  column "demo_url" {
    null = true
    type = character_varying(500)
  }
  column "github_url" {
    null = true
    type = character_varying(500)
  }
  column "image_url" {
    null = true
    type = character_varying(500)
  }
  column "role" {
    null = true
    type = character_varying(100)
  }
  column "start_date" {
    null = true
    type = timestamptz
  }
  column "end_date" {
    null = true
    type = timestamptz
  }
  column "featured" {
    null    = false
    type    = boolean
    default = false
  }
  column "is_archived" {
    null    = false
    type    = boolean
    default = false
  }
  column "github_repo_id" {
    null = true
    type = bigint
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_projects_github_repo" {
    columns     = [column.github_repo_id]
    ref_columns = [table.git_hub_repositories.column.id]
    on_update   = NO_ACTION
    on_delete   = SET_NULL
  }
  foreign_key "fk_projects_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_projects_user_id" {
    columns = [column.user_id]
  }
}

# =====================================================================
# project_milestone
# =====================================================================

table "project_milestones" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "project_id" {
    null = false
    type = bigint
  }
  column "title" {
    null = false
    type = character_varying(200)
  }
  column "description" {
    null = true
    type = text
  }
  column "status" {
    null    = true
    type    = character_varying(20)
    default = "not_started"
  }
  column "due_date" {
    null = true
    type = timestamptz
  }
  column "completed_at" {
    null = true
    type = timestamptz
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_project_milestones_project_id" {
    columns = [column.project_id]
  }
  foreign_key "fk_project_milestones_project" {
    columns     = [column.project_id]
    ref_columns = [table.projects.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  check "ck_project_milestones_status" {
    expr = "status IN ('not_started', 'in_progress', 'completed')"
  }
}

# =====================================================================
# qiita
# =====================================================================

table "qiita_articles" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "qiita_id" {
    null = false
    type = text
  }
  column "title" {
    null = false
    type = text
  }
  column "url" {
    null = false
    type = text
  }
  column "likes_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "comments_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "tags" {
    null = true
    type = text
  }
  column "published_at" {
    null = true
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "uq_qiita_articles_user_article" {
    unique  = true
    columns = [column.user_id, column.qiita_id]
  }
  index "idx_qiita_articles_user_id" {
    columns = [column.user_id]
  }
  foreign_key "fk_qiita_articles_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}

# =====================================================================
# question
# =====================================================================

table "question_bookmarks" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "question_id" {
    null = false
    type = bigint
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_question_bookmarks_question" {
    columns     = [column.question_id]
    ref_columns = [table.questions.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "uq_question_bookmark" {
    unique  = true
    columns = [column.user_id, column.question_id]
  }
  foreign_key "fk_question_bookmarks_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "question_votes" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "question_id" {
    null = false
    type = bigint
  }
  column "value" {
    null = false
    type = bigint
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "uq_question_vote" {
    unique  = true
    columns = [column.user_id, column.question_id]
  }
  foreign_key "fk_question_votes_question" {
    columns     = [column.question_id]
    ref_columns = [table.questions.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "fk_question_votes_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  check "ck_question_votes_value" {
    expr = "value IN (1, -1)"
  }
}
table "questions" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "title" {
    null = false
    type = character_varying(500)
  }
  column "body" {
    null = false
    type = text
  }
  column "tags" {
    null = true
    type = text
  }
  column "vote_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "answer_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "is_solved" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_questions_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_questions_user_id" {
    columns = [column.user_id]
  }
}

# =====================================================================
# reminder_settings
# =====================================================================

table "reminder_settings" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "enabled" {
    null    = false
    type    = boolean
    default = true
  }
  column "frequency" {
    null    = true
    type    = text
    default = "daily"
  }
  column "notification_time" {
    null    = true
    type    = text
    default = "09:00"
  }
  column "inactive_days" {
    null    = true
    type    = bigint
    default = 3
  }
  column "enable_web" {
    null    = false
    type    = boolean
    default = true
  }
  column "enable_email" {
    null    = false
    type    = boolean
    default = false
  }
  column "last_reminded_at" {
    null = true
    type = timestamptz
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_reminder_settings_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "uq_reminder_settings_user_id" {
    unique  = true
    columns = [column.user_id]
  }
  check "ck_reminder_settings_frequency" {
    expr = "frequency IN ('daily', 'weekly')"
  }
}

# =====================================================================
# resource_progress
# =====================================================================

table "resource_progresses" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "resource_id" {
    null = false
    type = bigint
  }
  column "status" {
    null    = true
    type    = character_varying(20)
    default = "not_started"
  }
  column "completion_percent" {
    null    = true
    type    = bigint
    default = 0
  }
  column "note" {
    null = true
    type = text
  }
  column "started_at" {
    null = true
    type = timestamptz
  }
  column "completed_at" {
    null = true
    type = timestamptz
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_resource_progresses_resource" {
    columns     = [column.resource_id]
    ref_columns = [table.learning_resources.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "uq_resource_progress" {
    unique  = true
    columns = [column.user_id, column.resource_id]
  }
  foreign_key "fk_resource_progresses_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  check "ck_resource_progresses_status" {
    expr = "status IN ('not_started', 'in_progress', 'completed')"
  }
}

# =====================================================================
# roadmap
# =====================================================================

table "roadmap_steps" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "roadmap_id" {
    null = false
    type = bigint
  }
  column "title" {
    null = false
    type = character_varying(200)
  }
  column "description" {
    null = true
    type = text
  }
  column "order_index" {
    null    = false
    type    = bigint
    default = 0
  }
  column "is_completed" {
    null    = false
    type    = boolean
    default = false
  }
  column "completed_at" {
    null = true
    type = timestamptz
  }
  column "resource_url" {
    null = true
    type = character_varying(500)
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_roadmaps_steps" {
    columns     = [column.roadmap_id]
    ref_columns = [table.roadmaps.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_roadmap_steps_roadmap_id" {
    columns = [column.roadmap_id]
  }
}
table "roadmaps" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "title" {
    null = false
    type = character_varying(200)
  }
  column "description" {
    null = true
    type = text
  }
  column "category" {
    null    = true
    type    = text
    default = "other"
  }
  column "is_public" {
    null    = false
    type    = boolean
    default = false
  }
  column "is_template" {
    null    = false
    type    = boolean
    default = false
  }
  column "status" {
    null    = true
    type    = text
    default = "active"
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  column "completed_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_roadmaps_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_roadmaps_is_public" {
    columns = [column.is_public]
  }
  index "idx_roadmaps_is_template" {
    columns = [column.is_template]
  }
  index "idx_roadmaps_user_id" {
    columns = [column.user_id]
  }
  check "ck_roadmaps_status" {
    expr = "status IN ('active', 'completed')"
  }
}

table "roadmap_metrics" {
  schema = schema.public
  column "roadmap_id" {
    null = false
    type = bigint
  }
  column "step_count" {
    null    = false
    type    = bigint
    default = 0
  }
  column "completed_step_count" {
    null    = false
    type    = bigint
    default = 0
  }
  primary_key {
    columns = [column.roadmap_id]
  }
  foreign_key "fk_roadmap_metrics_roadmap" {
    columns     = [column.roadmap_id]
    ref_columns = [table.roadmaps.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  check "ck_roadmap_metrics_non_negative" {
    expr = "step_count >= 0 AND completed_step_count >= 0"
  }
}

# =====================================================================
# spotify
# =====================================================================

table "spotify_recent_tracks" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "track_name" {
    null = false
    type = text
  }
  column "artist_name" {
    null = false
    type = text
  }
  column "album_name" {
    null = true
    type = text
  }
  column "album_image" {
    null = true
    type = text
  }
  column "track_url" {
    null = true
    type = text
  }
  column "played_at" {
    null = false
    type = timestamptz
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_spotify_recent_tracks_user_id" {
    columns = [column.user_id]
  }
  foreign_key "fk_spotify_recent_tracks_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}

# =====================================================================
# streak_freeze
# =====================================================================

table "streak_freezes" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "used_on" {
    null = false
    type = date
  }
  column "month" {
    null = false
    type = bigint
  }
  column "year" {
    null = false
    type = bigint
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "uq_streak_freeze_user_date" {
    unique  = true
    columns = [column.user_id, column.used_on]
  }
  index "idx_streak_freezes_user_id" {
    columns = [column.user_id]
  }
  foreign_key "fk_streak_freezes_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}

# =====================================================================
# study_circle
# =====================================================================

table "study_circle_checkins" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "circle_id" {
    null = false
    type = bigint
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "checked_on" {
    null = false
    type = date
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_study_circle_checkins_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "uq_checkin_unique" {
    unique  = true
    columns = [column.circle_id, column.user_id, column.checked_on]
  }
  foreign_key "fk_study_circle_checkins_circle" {
    columns     = [column.circle_id]
    ref_columns = [table.study_circles.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "study_circle_member_progresses" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "circle_id" {
    null = false
    type = bigint
  }
  column "step_id" {
    null = false
    type = bigint
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "is_completed" {
    null    = false
    type    = boolean
    default = false
  }
  column "completed_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_study_circle_member_progresses_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "uq_circle_step_user" {
    unique  = true
    columns = [column.circle_id, column.step_id, column.user_id]
  }
  foreign_key "fk_study_circle_member_progresses_circle" {
    columns     = [column.circle_id]
    ref_columns = [table.study_circles.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "fk_study_circle_member_progresses_step" {
    columns     = [column.step_id]
    ref_columns = [table.study_circle_steps.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "study_circle_members" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "circle_id" {
    null = false
    type = bigint
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "role" {
    null    = true
    type    = text
    default = "member"
  }
  column "joined_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_study_circle_members_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "fk_study_circles_members" {
    columns     = [column.circle_id]
    ref_columns = [table.study_circles.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "uq_circle_user" {
    unique  = true
    columns = [column.circle_id, column.user_id]
  }
  check "ck_study_circle_members_role" {
    expr = "role IN ('owner', 'member')"
  }
}
table "study_circle_steps" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "circle_id" {
    null = false
    type = bigint
  }
  column "title" {
    null = false
    type = character_varying(200)
  }
  column "description" {
    null = true
    type = text
  }
  column "order_index" {
    null    = true
    type    = bigint
    default = 0
  }
  column "resource_url" {
    null = true
    type = character_varying(2000)
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_study_circles_steps" {
    columns     = [column.circle_id]
    ref_columns = [table.study_circles.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_study_circle_steps_circle_id" {
    columns = [column.circle_id]
  }
}
table "study_circles" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "name" {
    null = false
    type = character_varying(200)
  }
  column "topic" {
    null = false
    type = character_varying(200)
  }
  column "description" {
    null = true
    type = text
  }
  column "owner_id" {
    null = false
    type = bigint
  }
  column "max_members" {
    null    = false
    type    = bigint
    default = 5
  }
  column "status" {
    null    = true
    type    = text
    default = "active"
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_study_circles_owner" {
    columns     = [column.owner_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_study_circles_owner_id" {
    columns = [column.owner_id]
  }
  check "ck_study_circles_status" {
    expr = "status IN ('active', 'completed', 'archived')"
  }
}

# =====================================================================
# user
# =====================================================================

table "users" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "username" {
    null = false
    type = text
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "avatar_url" {
    null = true
    type = text
  }
  column "bio" {
    null = true
    type = text
  }
  column "git_hub_id" {
    null = true
    type = bigint
  }
  column "git_hub_username" {
    null = true
    type = text
  }
  column "git_hub_token" {
    null = true
    type = text
  }
  column "git_hub_connected" {
    null    = false
    type    = boolean
    default = false
  }
  column "spotify_connected" {
    null    = false
    type    = boolean
    default = false
  }
  column "spotify_token" {
    null = true
    type = text
  }
  column "spotify_refresh_token" {
    null = true
    type = text
  }
  column "spotify_token_expiry" {
    null = true
    type = timestamptz
  }
  column "zenn_username" {
    null = true
    type = text
  }
  column "qiita_username" {
    null = true
    type = text
  }
  column "at_coder_username" {
    null = true
    type = text
  }
  column "paiza_rank" {
    null = true
    type = text
  }
  column "skills_languages" {
    null = true
    type = text
  }
  column "skills_frameworks" {
    null = true
    type = text
  }
  column "onboarding_completed" {
    null    = false
    type    = boolean
    default = false
  }
  column "email_weekly_report" {
    null    = false
    type    = boolean
    default = true
  }
  column "email_language" {
    null    = true
    type    = text
    default = "ja"
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "uq_users_email" {
    unique  = true
    columns = [column.email]
  }
  index "uq_users_git_hub_id_linked" {
    unique  = true
    columns = [column.git_hub_id]
    where   = "(git_hub_id <> 0)"
  }
  index "uq_users_username" {
    unique  = true
    columns = [column.username]
  }
  check "ck_users_email_language" {
    expr = "email_language IN ('ja', 'en', 'ko', 'zh-CN', 'zh-TW', 'es', 'fr', 'de', 'pt', 'ru')"
  }
}

table "user_credentials" {
  schema = schema.public
  column "user_id" {
    null = false
    type = bigint
  }
  column "password_hash" {
    null = false
    type = text
  }
  primary_key {
    columns = [column.user_id]
  }
  foreign_key "fk_user_credentials_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}

# =====================================================================
# user_activity
# =====================================================================

table "user_activities" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "activity_type" {
    null = false
    type = character_varying(50)
  }
  column "target_type" {
    null = false
    type = character_varying(50)
  }
  column "target_id" {
    null = false
    type = bigint
  }
  column "metadata" {
    null = true
    type = text
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_user_activities_activity_type" {
    columns = [column.activity_type]
  }
  index "idx_user_activities_created_at" {
    columns = [column.created_at]
  }
  index "idx_user_activities_user_id" {
    columns = [column.user_id]
  }
  foreign_key "fk_user_activities_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  check "ck_user_activities_activity_type" {
    expr = "activity_type IN ('post_created', 'comment_created', 'resource_shared', 'goal_completed', 'project_created', 'snippet_created')"
  }
}

# =====================================================================
# weekly_challenge
# =====================================================================

table "weekly_challenges" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "year" {
    null = false
    type = bigint
  }
  column "week" {
    null = false
    type = bigint
  }
  column "challenge_type" {
    null = false
    type = text
  }
  column "description" {
    null = false
    type = text
  }
  column "target_value" {
    null = false
    type = bigint
  }
  column "current_value" {
    null    = true
    type    = bigint
    default = 0
  }
  column "is_completed" {
    null    = false
    type    = boolean
    default = false
  }
  column "completed_at" {
    null = true
    type = timestamptz
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_weekly_challenges_user_id" {
    columns = [column.user_id]
  }
  foreign_key "fk_weekly_challenges_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  check "ck_weekly_challenges_challenge_type" {
    expr = "challenge_type IN ('duration_total', 'streak_days', 'category_count', 'log_count')"
  }
}

# =====================================================================
# weekly_goal
# =====================================================================

table "weekly_goals" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "category" {
    null = false
    type = text
  }
  column "target_minutes" {
    null    = false
    type    = bigint
    default = 0
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "uq_weekly_goal_unique" {
    unique  = true
    columns = [column.user_id, column.category]
  }
  foreign_key "fk_weekly_goals_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  check "ck_weekly_goals_category" {
    expr = "category IN ('coding', 'reading', 'course', 'meetup', 'other')"
  }
}

# =====================================================================
# widget_settings
# =====================================================================

table "widget_settings" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "settings" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "uq_widget_settings_user_id" {
    unique  = true
    columns = [column.user_id]
  }
  foreign_key "fk_widget_settings_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}

# =====================================================================
# youtube_video
# =====================================================================

table "you_tube_search_caches" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "query" {
    null = false
    type = character_varying(500)
  }
  column "language" {
    null    = true
    type    = character_varying(10)
    default = "ja"
  }
  column "video_ids" {
    null = true
    type = text
  }
  column "cache_expires" {
    null = true
    type = timestamptz
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_you_tube_search_caches_cache_expires" {
    columns = [column.cache_expires]
  }
  index "idx_you_tube_search_caches_query" {
    columns = [column.query]
  }
}
table "you_tube_videos" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "video_id" {
    null = false
    type = character_varying(20)
  }
  column "title" {
    null = false
    type = character_varying(500)
  }
  column "description" {
    null = true
    type = text
  }
  column "channel_id" {
    null = true
    type = character_varying(50)
  }
  column "channel_title" {
    null = true
    type = character_varying(200)
  }
  column "thumbnail_url" {
    null = true
    type = character_varying(500)
  }
  column "published_at" {
    null = true
    type = timestamptz
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "uq_you_tube_videos_video_id" {
    unique  = true
    columns = [column.video_id]
  }
}

# =====================================================================
# zenn
# =====================================================================

table "zenn_articles" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "zenn_id" {
    null = false
    type = bigint
  }
  column "title" {
    null = false
    type = text
  }
  column "slug" {
    null = false
    type = text
  }
  column "emoji" {
    null = true
    type = text
  }
  column "article_type" {
    null = true
    type = text
  }
  column "liked_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "comments_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "published_at" {
    null = true
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_zenn_articles_user_id" {
    columns = [column.user_id]
  }
  index "uq_zenn_articles_user_article" {
    unique  = true
    columns = [column.user_id, column.zenn_id]
  }
  foreign_key "fk_zenn_articles_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}

