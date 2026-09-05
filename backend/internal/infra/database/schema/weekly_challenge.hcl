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
    null    = true
    type    = boolean
    default = false
  }
  column "completed_at" {
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
  index "idx_weekly_challenges_user_id" {
    columns = [column.user_id]
  }
}
