table "book_reviews" {
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
    type = character_varying(300)
  }
  column "author" {
    null = true
    type = character_varying(200)
  }
  column "isbn" {
    null = true
    type = character_varying(20)
  }
  column "rating" {
    null = false
    type = bigint
  }
  column "review" {
    null = true
    type = text
  }
  column "total_pages" {
    null    = true
    type    = bigint
    default = 0
  }
  column "current_page" {
    null    = true
    type    = bigint
    default = 0
  }
  column "image_url" {
    null = true
    type = character_varying(2000)
  }
  column "status" {
    null    = true
    type    = text
    default = "not_started"
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
  column "deleted_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_book_reviews_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_book_reviews_deleted_at" {
    columns = [column.deleted_at]
  }
  index "idx_book_reviews_user_id" {
    columns = [column.user_id]
  }
}
