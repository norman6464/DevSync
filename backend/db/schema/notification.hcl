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
    null    = true
    type    = boolean
    default = false
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_notifications_actor" {
    columns     = [column.actor_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "fk_notifications_post" {
    columns     = [column.post_id]
    ref_columns = [table.posts.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "fk_notifications_question" {
    columns     = [column.question_id]
    ref_columns = [table.questions.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "fk_notifications_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
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
}
