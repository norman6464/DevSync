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
    null    = true
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
    null    = true
    type    = boolean
    default = true
  }
  column "enable_email" {
    null    = true
    type    = boolean
    default = false
  }
  column "last_reminded_at" {
    null = true
    type = timestamptz
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
  foreign_key "fk_reminder_settings_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_reminder_settings_user_id" {
    unique  = true
    columns = [column.user_id]
  }
}
