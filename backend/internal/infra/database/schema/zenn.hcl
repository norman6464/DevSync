table "zenn_articles" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "zenn_id" {
    null = false
    type = bigint
  }
  column "title" {
    null = false
    type = text
  }
  column "slug" {
    null = false
    type = text
  }
  column "emoji" {
    null = true
    type = text
  }
  column "article_type" {
    null = true
    type = text
  }
  column "liked_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "comments_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "published_at" {
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
  index "idx_zenn_articles_user_id" {
    columns = [column.user_id]
  }
  index "idx_zenn_articles_zenn_id" {
    unique  = true
    columns = [column.zenn_id]
  }
}
