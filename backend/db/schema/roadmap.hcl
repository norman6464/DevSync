table "roadmap_steps" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "roadmap_id" {
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
  column "order_index" {
    null    = false
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
  column "resource_url" {
    null = true
    type = character_varying(500)
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
  foreign_key "fk_roadmaps_steps" {
    columns     = [column.roadmap_id]
    ref_columns = [table.roadmaps.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_roadmap_steps_roadmap_id" {
    columns = [column.roadmap_id]
  }
}
table "roadmaps" {
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
    type = character_varying(200)
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
  column "is_public" {
    null    = true
    type    = boolean
    default = false
  }
  column "is_template" {
    null    = true
    type    = boolean
    default = false
  }
  column "step_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "completed_step_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "progress" {
    null    = true
    type    = bigint
    default = 0
  }
  column "status" {
    null    = true
    type    = text
    default = "active"
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
  foreign_key "fk_roadmaps_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_roadmaps_is_public" {
    columns = [column.is_public]
  }
  index "idx_roadmaps_is_template" {
    columns = [column.is_template]
  }
  index "idx_roadmaps_user_id" {
    columns = [column.user_id]
  }
}
