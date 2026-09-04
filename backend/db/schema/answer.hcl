table "answer_votes" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "answer_id" {
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
  index "idx_answer_vote" {
    unique  = true
    columns = [column.user_id, column.answer_id]
  }
}
table "answers" {
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
  column "body" {
    null = false
    type = text
  }
  column "vote_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "is_best" {
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
  foreign_key "fk_answers_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_answers_deleted_at" {
    columns = [column.deleted_at]
  }
  index "idx_answers_question_id" {
    columns = [column.question_id]
  }
  index "idx_answers_user_id" {
    columns = [column.user_id]
  }
}
