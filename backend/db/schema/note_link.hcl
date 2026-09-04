table "note_links" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "source_note_id" {
    null = false
    type = bigint
  }
  column "target_note_id" {
    null = false
    type = bigint
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_note_links_source_note" {
    columns     = [column.source_note_id]
    ref_columns = [table.notes.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "fk_note_links_target_note" {
    columns     = [column.target_note_id]
    ref_columns = [table.notes.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_note_link_unique" {
    unique  = true
    columns = [column.source_note_id, column.target_note_id]
  }
  index "idx_note_links_source_note_id" {
    columns = [column.source_note_id]
  }
  index "idx_note_links_target_note_id" {
    columns = [column.target_note_id]
  }
}
