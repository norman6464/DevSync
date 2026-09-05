table "you_tube_search_caches" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "query" {
    null = false
    type = character_varying(500)
  }
  column "language" {
    null    = true
    type    = character_varying(10)
    default = "ja"
  }
  column "video_ids" {
    null = true
    type = text
  }
  column "cache_expires" {
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
  index "idx_you_tube_search_caches_cache_expires" {
    columns = [column.cache_expires]
  }
  index "idx_you_tube_search_caches_query" {
    columns = [column.query]
  }
}
table "you_tube_videos" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "video_id" {
    null = false
    type = character_varying(20)
  }
  column "title" {
    null = false
    type = character_varying(500)
  }
  column "description" {
    null = true
    type = text
  }
  column "channel_id" {
    null = true
    type = character_varying(50)
  }
  column "channel_title" {
    null = true
    type = character_varying(200)
  }
  column "thumbnail_url" {
    null = true
    type = character_varying(500)
  }
  column "published_at" {
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
  index "idx_you_tube_videos_video_id" {
    unique  = true
    columns = [column.video_id]
  }
}
