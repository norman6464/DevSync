table "project_milestones" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "project_id" {
    null = false
    type = bigint
  }
  column "title" {
    null = false
    type = character_varying(200)
  }
  column "description" {
    null = true
    type = text
  }
  column "status" {
    null    = true
    type    = character_varying(20)
    default = "not_started"
  }
  column "due_date" {
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
  index "idx_project_milestones_project_id" {
    columns = [column.project_id]
  }
}
