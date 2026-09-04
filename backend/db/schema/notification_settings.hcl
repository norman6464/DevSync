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
    null    = true
    type    = boolean
    default = true
  }
  column "enable_comments" {
    null    = true
    type    = boolean
    default = true
  }
  column "enable_follows" {
    null    = true
    type    = boolean
    default = true
  }
  column "enable_messages" {
    null    = true
    type    = boolean
    default = true
  }
  column "enable_mentions" {
    null    = true
    type    = boolean
    default = true
  }
  column "enable_web_push" {
    null    = true
    type    = boolean
    default = true
  }
  column "enable_email" {
    null    = true
    type    = boolean
    default = true
  }
  column "enable_sound" {
    null    = true
    type    = boolean
    default = true
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
  foreign_key "fk_notification_settings_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_notification_settings_user_id" {
    unique  = true
    columns = [column.user_id]
  }
}
