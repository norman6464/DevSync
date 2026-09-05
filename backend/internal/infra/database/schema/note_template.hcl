table "note_templates" {
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
  column "description" {
    null = true
    type = character_varying(500)
  }
  column "default_title" {
    null = true
    type = character_varying(200)
  }
  column "content_template" {
    null = false
    type = text
  }
  column "default_tags" {
    null = true
    type = character_varying(255)
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
  index "idx_note_templates_is_default" {
    columns = [column.is_default]
  }
  index "idx_note_templates_user_id" {
    columns = [column.user_id]
  }
}
