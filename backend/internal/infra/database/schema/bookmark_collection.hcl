table "bookmark_collection_items" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "collection_id" {
    null = false
    type = bigint
  }
  column "post_id" {
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
  foreign_key "fk_bookmark_collection_items_post" {
    columns     = [column.post_id]
    ref_columns = [table.posts.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_bookmark_collection_items_post_id" {
    columns = [column.post_id]
  }
  index "idx_bookmark_collection_post" {
    unique  = true
    columns = [column.collection_id, column.post_id]
  }
}
table "bookmark_collections" {
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
    type = text
  }
  column "description" {
    null = true
    type = text
  }
  column "color" {
    null    = true
    type    = text
    default = "blue"
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
  index "idx_bookmark_collections_user_id" {
    columns = [column.user_id]
  }
}
