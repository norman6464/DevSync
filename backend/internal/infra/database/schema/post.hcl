table "bookmarks" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
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
  foreign_key "fk_bookmarks_post" {
    columns     = [column.post_id]
    ref_columns = [table.posts.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_bookmarks_post_id" {
    columns = [column.post_id]
  }
  index "idx_user_post_bookmark" {
    unique  = true
    columns = [column.user_id, column.post_id]
  }
}
table "comment_likes" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "comment_id" {
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
  index "idx_comment_likes_comment_id" {
    columns = [column.comment_id]
  }
  index "idx_user_comment_like" {
    unique  = true
    columns = [column.user_id, column.comment_id]
  }
}
table "comments" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "post_id" {
    null = false
    type = bigint
  }
  column "parent_id" {
    null = true
    type = bigint
  }
  column "content" {
    null = false
    type = text
  }
  column "like_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "is_hidden" {
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
  foreign_key "fk_comments_replies" {
    columns     = [column.parent_id]
    ref_columns = [table.comments.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "fk_comments_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_comments_parent_id" {
    columns = [column.parent_id]
  }
  index "idx_comments_post_id" {
    columns = [column.post_id]
  }
  index "idx_comments_user_id" {
    columns = [column.user_id]
  }
}
table "likes" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
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
  index "idx_likes_post_id" {
    columns = [column.post_id]
  }
  index "idx_user_post_like" {
    unique  = true
    columns = [column.user_id, column.post_id]
  }
}
table "post_collection_items" {
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
  column "note" {
    null = true
    type = text
  }
  column "order_index" {
    null    = false
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
  foreign_key "fk_post_collection_items_post" {
    columns     = [column.post_id]
    ref_columns = [table.posts.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_post_collection_item" {
    unique  = true
    columns = [column.collection_id, column.post_id]
  }
  index "idx_post_collection_items_collection_id" {
    columns = [column.collection_id]
  }
  index "idx_post_collection_items_post_id" {
    columns = [column.post_id]
  }
}
table "post_collections" {
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
    type = text
  }
  column "description" {
    null = true
    type = text
  }
  column "is_public" {
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
  foreign_key "fk_post_collections_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_post_collections_user_id" {
    columns = [column.user_id]
  }
}
table "post_pins" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "post_id" {
    null = false
    type = bigint
  }
  column "pin_order" {
    null    = false
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
  foreign_key "fk_post_pins_post" {
    columns     = [column.post_id]
    ref_columns = [table.posts.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_post_pins_post_id" {
    columns = [column.post_id]
  }
  index "idx_post_pins_user_id" {
    columns = [column.user_id]
  }
  index "idx_user_post_pin" {
    unique  = true
    columns = [column.user_id, column.post_id]
  }
}
table "post_series" {
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
    type = text
  }
  column "description" {
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
  foreign_key "fk_post_series_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_post_series_user_id" {
    columns = [column.user_id]
  }
}
table "post_series_items" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "series_id" {
    null = false
    type = bigint
  }
  column "post_id" {
    null = false
    type = bigint
  }
  column "order_index" {
    null    = false
    type    = bigint
    default = 0
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_post_series_items_post" {
    columns     = [column.post_id]
    ref_columns = [table.posts.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_post_series_items_post_id" {
    columns = [column.post_id]
  }
  index "idx_post_series_items_series_id" {
    columns = [column.series_id]
  }
  index "idx_series_post" {
    unique  = true
    columns = [column.series_id, column.post_id]
  }
}
table "post_tags" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "post_id" {
    null = false
    type = bigint
  }
  column "tag" {
    null = false
    type = character_varying(50)
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_post_tag" {
    unique  = true
    columns = [column.post_id, column.tag]
  }
  index "idx_post_tags_post_id" {
    columns = [column.post_id]
  }
  index "idx_post_tags_tag" {
    columns = [column.tag]
  }
}
table "post_views" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
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
  index "idx_post_views_post_id" {
    columns = [column.post_id]
  }
  index "idx_post_views_user_id" {
    columns = [column.user_id]
  }
  index "idx_user_post_view" {
    unique  = true
    columns = [column.user_id, column.post_id]
  }
}
table "posts" {
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
    type = text
  }
  column "content" {
    null = false
    type = text
  }
  column "image_urls" {
    null = true
    type = text
  }
  column "is_draft" {
    null    = true
    type    = boolean
    default = false
  }
  column "like_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "comment_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "bookmark_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "view_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "estimated_read_time" {
    null    = true
    type    = bigint
    default = 0
  }
  column "scheduled_at" {
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
  foreign_key "fk_posts_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_posts_is_draft" {
    columns = [column.is_draft]
  }
  index "idx_posts_scheduled_at" {
    columns = [column.scheduled_at]
  }
  index "idx_posts_user_id" {
    columns = [column.user_id]
  }
}
table "reactions" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "post_id" {
    null = false
    type = bigint
  }
  column "emoji" {
    null = false
    type = character_varying(10)
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_reactions_post_id" {
    columns = [column.post_id]
  }
  index "idx_user_post_emoji" {
    unique  = true
    columns = [column.user_id, column.post_id, column.emoji]
  }
}
