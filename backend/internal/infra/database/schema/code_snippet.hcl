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
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_snippet_favorite" {
    unique  = true
    columns = [column.user_id, column.snippet_id]
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
    null = true
    type = timestamptz
  }
  column "updated_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_posts_code_snippets" {
    columns     = [column.post_id]
    ref_columns = [table.posts.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
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
    null = true
    type = timestamptz
  }
  column "updated_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_snippet_comments_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_snippet_comments_snippet_id" {
    columns = [column.snippet_id]
  }
  index "idx_snippet_comments_user_id" {
    columns = [column.user_id]
  }
}
