table "qiita_articles" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "qiita_id" {
    null = false
    type = text
  }
  column "title" {
    null = false
    type = text
  }
  column "url" {
    null = false
    type = text
  }
  column "likes_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "comments_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "tags" {
    null = true
    type = text
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
  index "idx_qiita_articles_qiita_id" {
    unique  = true
    columns = [column.qiita_id]
  }
  index "idx_qiita_articles_user_id" {
    columns = [column.user_id]
  }
}
