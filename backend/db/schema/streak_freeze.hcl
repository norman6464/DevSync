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
  column "used_date" {
    null = false
    type = character_varying(10)
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
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_streak_freeze_user_date" {
    unique  = true
    columns = [column.user_id, column.used_date]
  }
  index "idx_streak_freezes_user_id" {
    columns = [column.user_id]
  }
}
