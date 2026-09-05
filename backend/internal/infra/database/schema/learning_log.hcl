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
    null    = true
    type    = boolean
    default = false
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
  index "idx_learning_logs_goal_id" {
    columns = [column.goal_id]
  }
  index "idx_learning_logs_user_id" {
    columns = [column.user_id]
  }
}
