table "learning_resources" {
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
  column "description" {
    null = true
    type = text
  }
  column "url" {
    null = true
    type = character_varying(500)
  }
  column "category" {
    null = false
    type = character_varying(50)
  }
  column "difficulty" {
    null = true
    type = character_varying(50)
  }
  column "tags" {
    null = true
    type = text
  }
  column "image_url" {
    null = true
    type = character_varying(500)
  }
  column "is_public" {
    null = true
    type = boolean
  }
  column "like_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "save_count" {
    null    = true
    type    = bigint
    default = 0
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
  foreign_key "fk_learning_resources_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_learning_resources_deleted_at" {
    columns = [column.deleted_at]
  }
  index "idx_learning_resources_user_id" {
    columns = [column.user_id]
  }
}
table "resource_likes" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "resource_id" {
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
  index "idx_resource_like" {
    unique  = true
    columns = [column.user_id, column.resource_id]
  }
}
table "resource_reviews" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "resource_id" {
    null = false
    type = bigint
  }
  column "rating" {
    null = false
    type = bigint
  }
  column "comment" {
    null = true
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
  foreign_key "fk_resource_reviews_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_resource_review" {
    unique  = true
    columns = [column.user_id, column.resource_id]
  }
  index "idx_resource_reviews_resource_id" {
    columns = [column.resource_id]
  }
}
table "resource_saves" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "resource_id" {
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
  index "idx_resource_save" {
    unique  = true
    columns = [column.user_id, column.resource_id]
  }
}
