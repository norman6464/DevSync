table "question_bookmarks" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "question_id" {
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
  foreign_key "fk_question_bookmarks_question" {
    columns     = [column.question_id]
    ref_columns = [table.questions.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_question_bookmark" {
    unique  = true
    columns = [column.user_id, column.question_id]
  }
}
table "question_votes" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "question_id" {
    null = false
    type = bigint
  }
  column "value" {
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
  index "idx_question_vote" {
    unique  = true
    columns = [column.user_id, column.question_id]
  }
}
table "questions" {
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
    type = character_varying(500)
  }
  column "body" {
    null = false
    type = text
  }
  column "tags" {
    null = true
    type = text
  }
  column "vote_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "answer_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "is_solved" {
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
  column "deleted_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_questions_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_questions_deleted_at" {
    columns = [column.deleted_at]
  }
  index "idx_questions_user_id" {
    columns = [column.user_id]
  }
}
