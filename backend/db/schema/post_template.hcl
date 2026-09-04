table "post_templates" {
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
  column "title_template" {
    null = true
    type = character_varying(200)
  }
  column "content_template" {
    null = false
    type = text
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
  index "idx_post_templates_user_id" {
    columns = [column.user_id]
  }
}
