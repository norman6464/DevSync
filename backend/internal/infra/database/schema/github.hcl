table "git_hub_contributions" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "date" {
    null = false
    type = timestamptz
  }
  column "count" {
    null    = false
    type    = bigint
    default = 0
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
  index "idx_user_date" {
    unique  = true
    columns = [column.user_id, column.date]
  }
}
table "git_hub_language_stats" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "language" {
    null = false
    type = text
  }
  column "bytes" {
    null    = false
    type    = bigint
    default = 0
  }
  column "repo_count" {
    null    = false
    type    = bigint
    default = 0
  }
  column "updated_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_user_lang" {
    unique  = true
    columns = [column.user_id, column.language]
  }
}
table "git_hub_repositories" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "git_hub_repo_id" {
    null = false
    type = bigint
  }
  column "name" {
    null = false
    type = text
  }
  column "full_name" {
    null = true
    type = text
  }
  column "description" {
    null = true
    type = text
  }
  column "language" {
    null = true
    type = text
  }
  column "stars" {
    null = true
    type = bigint
  }
  column "forks" {
    null = true
    type = bigint
  }
  column "is_private" {
    null = true
    type = boolean
  }
  column "updated_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_git_hub_repositories_git_hub_repo_id" {
    unique  = true
    columns = [column.git_hub_repo_id]
  }
  index "idx_git_hub_repositories_user_id" {
    columns = [column.user_id]
  }
}
