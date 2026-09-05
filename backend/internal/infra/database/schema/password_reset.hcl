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
  foreign_key "fk_password_reset_tokens_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_password_reset_tokens_token" {
    unique  = true
    columns = [column.token]
  }
  index "idx_password_reset_tokens_user_id" {
    columns = [column.user_id]
  }
}
