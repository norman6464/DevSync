table "resource_progresses" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "resource_id" {
    null = false
    type = bigint
  }
  column "status" {
    null    = true
    type    = character_varying(20)
    default = "not_started"
  }
  column "completion_percent" {
    null    = true
    type    = bigint
    default = 0
  }
  column "note" {
    null = true
    type = text
  }
  column "started_at" {
    null = true
    type = timestamptz
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
  foreign_key "fk_resource_progresses_resource" {
    columns     = [column.resource_id]
    ref_columns = [table.learning_resources.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_resource_progress" {
    unique  = true
    columns = [column.user_id, column.resource_id]
  }
}
