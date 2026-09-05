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
  index "idx_weekly_goal_unique" {
    unique  = true
    columns = [column.user_id, column.category]
  }
}
