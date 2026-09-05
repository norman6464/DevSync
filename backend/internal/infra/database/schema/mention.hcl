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
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_mentions_actor" {
    columns     = [column.actor_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "fk_mentions_post" {
    columns     = [column.post_id]
    ref_columns = [table.posts.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "fk_mentions_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_mentions_comment_id" {
    columns = [column.comment_id]
  }
  index "idx_mentions_comment_user" {
    unique  = true
    columns = [column.comment_id, column.user_id]
    where   = "(comment_id IS NOT NULL)"
  }
  index "idx_mentions_post_id" {
    columns = [column.post_id]
  }
  index "idx_mentions_post_user" {
    unique  = true
    columns = [column.post_id, column.user_id]
    where   = "((post_id IS NOT NULL) AND (comment_id IS NULL))"
  }
  index "idx_mentions_user_id" {
    columns = [column.user_id]
  }
}
