table "user_activities" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "activity_type" {
    null = false
    type = character_varying(50)
  }
  column "target_type" {
    null = false
    type = character_varying(50)
  }
  column "target_id" {
    null = false
    type = bigint
  }
  column "metadata" {
    null = true
    type = text
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_user_activities_activity_type" {
    columns = [column.activity_type]
  }
  index "idx_user_activities_created_at" {
    columns = [column.created_at]
  }
  index "idx_user_activities_user_id" {
    columns = [column.user_id]
  }
}
