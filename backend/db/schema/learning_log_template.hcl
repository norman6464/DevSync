table "learning_log_templates" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "name" {
    null = false
    type = character_varying(100)
  }
  column "default_title" {
    null = true
    type = character_varying(200)
  }
  column "default_content" {
    null = true
    type = text
  }
  column "default_category" {
    null    = true
    type    = character_varying(50)
    default = "other"
  }
  column "default_duration" {
    null    = true
    type    = bigint
    default = 0
  }
  column "is_default" {
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
  index "idx_learning_log_templates_user_id" {
    columns = [column.user_id]
  }
  index "idx_log_templates_is_default" {
    columns = [column.is_default]
  }
}
