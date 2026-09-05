table "study_circle_checkins" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "circle_id" {
    null = false
    type = bigint
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "date" {
    null = false
    type = character_varying(10)
  }
  column "content" {
    null = false
    type = text
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_study_circle_checkins_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_checkin_unique" {
    unique  = true
    columns = [column.circle_id, column.user_id, column.date]
  }
}
table "study_circle_member_progresses" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "circle_id" {
    null = false
    type = bigint
  }
  column "step_id" {
    null = false
    type = bigint
  }
  column "user_id" {
    null = false
    type = bigint
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
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_study_circle_member_progresses_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_circle_step_user" {
    unique  = true
    columns = [column.circle_id, column.step_id, column.user_id]
  }
}
table "study_circle_members" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "circle_id" {
    null = false
    type = bigint
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "role" {
    null    = true
    type    = text
    default = "member"
  }
  column "joined_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_study_circle_members_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "fk_study_circles_members" {
    columns     = [column.circle_id]
    ref_columns = [table.study_circles.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_circle_user" {
    unique  = true
    columns = [column.circle_id, column.user_id]
  }
}
table "study_circle_steps" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "circle_id" {
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
    null    = true
    type    = bigint
    default = 0
  }
  column "resource_url" {
    null = true
    type = character_varying(2000)
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
  foreign_key "fk_study_circles_steps" {
    columns     = [column.circle_id]
    ref_columns = [table.study_circles.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_study_circle_steps_circle_id" {
    columns = [column.circle_id]
  }
}
table "study_circles" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "name" {
    null = false
    type = character_varying(200)
  }
  column "topic" {
    null = false
    type = character_varying(200)
  }
  column "description" {
    null = true
    type = text
  }
  column "owner_id" {
    null = false
    type = bigint
  }
  column "max_members" {
    null    = false
    type    = bigint
    default = 5
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
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_study_circles_owner" {
    columns     = [column.owner_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_study_circles_owner_id" {
    columns = [column.owner_id]
  }
}
