table "chat_room_members" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "chat_room_id" {
    null = false
    type = bigint
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "joined_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_chat_room_members_chat_room" {
    columns     = [column.chat_room_id]
    ref_columns = [table.chat_rooms.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "fk_chat_room_members_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_chat_room_members_chat_room_id" {
    columns = [column.chat_room_id]
  }
  index "idx_chat_room_members_user_id" {
    columns = [column.user_id]
  }
  index "idx_room_user" {
    unique  = true
    columns = [column.chat_room_id, column.user_id]
  }
}
table "chat_rooms" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "name" {
    null = false
    type = character_varying(100)
  }
  column "description" {
    null = true
    type = character_varying(500)
  }
  column "owner_id" {
    null = false
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
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_chat_rooms_owner" {
    columns     = [column.owner_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_chat_rooms_owner_id" {
    columns = [column.owner_id]
  }
}
table "group_messages" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "chat_room_id" {
    null = false
    type = bigint
  }
  column "sender_id" {
    null = false
    type = bigint
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
  foreign_key "fk_group_messages_chat_room" {
    columns     = [column.chat_room_id]
    ref_columns = [table.chat_rooms.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "fk_group_messages_sender" {
    columns     = [column.sender_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_group_messages_chat_room_id" {
    columns = [column.chat_room_id]
  }
  index "idx_group_messages_sender_id" {
    columns = [column.sender_id]
  }
}
