table "note_versions" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "note_id" {
    null = false
    type = bigint
  }
  column "version_number" {
    null = false
    type = bigint
  }
  column "title" {
    null = false
    type = text
  }
  column "content" {
    null = true
    type = text
  }
  column "tags" {
    null = true
    type = text
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_note_versions_note_id" {
    columns = [column.note_id]
  }
}
