table "messages" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "sender_id" {
    null = false
    type = bigint
  }
  column "receiver_id" {
    null = false
    type = bigint
  }
  column "content" {
    null = false
    type = text
  }
  column "read" {
    null    = true
    type    = boolean
    default = false
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_messages_receiver" {
    columns     = [column.receiver_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "fk_messages_sender" {
    columns     = [column.sender_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_messages_receiver_id" {
    columns = [column.receiver_id]
  }
  index "idx_messages_sender_id" {
    columns = [column.sender_id]
  }
}
