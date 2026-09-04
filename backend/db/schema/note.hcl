table "note_folders" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "parent_id" {
    null = true
    type = bigint
  }
  column "name" {
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
  foreign_key "fk_note_folders_parent" {
    columns     = [column.parent_id]
    ref_columns = [table.note_folders.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "fk_note_folders_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_note_folders_parent_id" {
    columns = [column.parent_id]
  }
  index "idx_note_folders_user_id" {
    columns = [column.user_id]
  }
}
table "notes" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "folder_id" {
    null = true
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
  column "is_favorite" {
    null    = true
    type    = boolean
    default = false
  }
  column "is_archived" {
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
  foreign_key "fk_notes_folder" {
    columns     = [column.folder_id]
    ref_columns = [table.note_folders.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "fk_notes_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_notes_folder_id" {
    columns = [column.folder_id]
  }
  index "idx_notes_user_id" {
    columns = [column.user_id]
  }
}
