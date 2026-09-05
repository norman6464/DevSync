table "projects" {
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
  column "tech_stack" {
    null = true
    type = text
  }
  column "demo_url" {
    null = true
    type = character_varying(500)
  }
  column "github_url" {
    null = true
    type = character_varying(500)
  }
  column "image_url" {
    null = true
    type = character_varying(500)
  }
  column "role" {
    null = true
    type = character_varying(100)
  }
  column "start_date" {
    null = true
    type = timestamptz
  }
  column "end_date" {
    null = true
    type = timestamptz
  }
  column "featured" {
    null    = true
    type    = boolean
    default = false
  }
  column "is_archived" {
    null    = true
    type    = boolean
    default = false
  }
  column "github_repo_id" {
    null = true
    type = bigint
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  column "updated_at" {
    null = true
    type = timestamptz
  }
  column "deleted_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_projects_github_repo" {
    columns     = [column.github_repo_id]
    ref_columns = [table.git_hub_repositories.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "fk_projects_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_projects_deleted_at" {
    columns = [column.deleted_at]
  }
  index "idx_projects_user_id" {
    columns = [column.user_id]
  }
}
