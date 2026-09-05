table "users" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "username" {
    null = false
    type = text
  }
  column "name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "password" {
    null = true
    type = text
  }
  column "avatar_url" {
    null = true
    type = text
  }
  column "bio" {
    null = true
    type = text
  }
  column "git_hub_id" {
    null = true
    type = bigint
  }
  column "git_hub_username" {
    null = true
    type = text
  }
  column "git_hub_token" {
    null = true
    type = text
  }
  column "git_hub_connected" {
    null    = true
    type    = boolean
    default = false
  }
  column "spotify_connected" {
    null    = true
    type    = boolean
    default = false
  }
  column "spotify_token" {
    null = true
    type = text
  }
  column "spotify_refresh_token" {
    null = true
    type = text
  }
  column "spotify_token_expiry" {
    null = true
    type = timestamptz
  }
  column "zenn_username" {
    null = true
    type = text
  }
  column "qiita_username" {
    null = true
    type = text
  }
  column "at_coder_username" {
    null = true
    type = text
  }
  column "paiza_rank" {
    null = true
    type = text
  }
  column "skills_languages" {
    null = true
    type = text
  }
  column "skills_frameworks" {
    null = true
    type = text
  }
  column "onboarding_completed" {
    null    = true
    type    = boolean
    default = false
  }
  column "email_weekly_report" {
    null    = true
    type    = boolean
    default = true
  }
  column "email_language" {
    null    = true
    type    = text
    default = "ja"
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
  index "idx_users_email" {
    unique  = true
    columns = [column.email]
  }
  index "idx_users_git_hub_id_linked" {
    unique  = true
    columns = [column.git_hub_id]
    where   = "(git_hub_id <> 0)"
  }
  index "idx_users_username" {
    unique  = true
    columns = [column.username]
  }
}
