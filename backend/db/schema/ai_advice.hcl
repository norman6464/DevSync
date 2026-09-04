table "ai_advices" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "type" {
    null = false
    type = text
  }
  column "priority" {
    null    = false
    type    = bigint
    default = 3
  }
  column "title_key" {
    null = false
    type = text
  }
  column "message_key" {
    null = false
    type = text
  }
  column "params" {
    null = true
    type = text
  }
  column "action_url" {
    null = true
    type = text
  }
  column "is_read" {
    null    = true
    type    = boolean
    default = false
  }
  column "expires_at" {
    null = true
    type = timestamptz
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
  index "idx_ai_advices_user_id" {
    columns = [column.user_id]
  }
}
table "ai_conversations" {
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
  index "idx_ai_conversations_user_id" {
    columns = [column.user_id]
  }
}
table "ai_messages" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "conversation_id" {
    null = false
    type = bigint
  }
  column "role" {
    null = false
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "tokens_used" {
    null    = true
    type    = bigint
    default = 0
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_ai_conversations_messages" {
    columns     = [column.conversation_id]
    ref_columns = [table.ai_conversations.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_ai_messages_conversation_id" {
    columns = [column.conversation_id]
  }
}
