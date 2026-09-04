table "learning_goals" {
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
  column "description" {
    null = true
    type = text
  }
  column "category" {
    null    = true
    type    = text
    default = "other"
  }
  column "target_date" {
    null = true
    type = timestamptz
  }
  column "progress" {
    null    = true
    type    = bigint
    default = 0
  }
  column "target_hours" {
    null    = true
    type    = bigint
    default = 0
  }
  column "status" {
    null    = true
    type    = text
    default = "active"
  }
  column "is_public" {
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
  column "completed_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_learning_goals_user_id" {
    columns = [column.user_id]
  }
}
